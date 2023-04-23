package controllers

import (
	"context"
	"fmt"
	"github.com/Prodigy00/restaurant-mgt/database"
	"github.com/Prodigy00/restaurant-mgt/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIdx := (page - 1) * recordPerPage
		startIdx, err = strconv.Atoi(c.Query("startIdx"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", "null"}}}, {"totalCount", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}
		projectStage := bson.D{
			{
				"$project", bson.D{
					{"_id", 0},
					{"totalCount", 1},
					{"foodItems", bson.D{{"$slice", []interface{}{"$data", startIdx, recordPerPage}}}},
				},
			},
		}

		result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing food items"})
		}

		var allFoods []bson.M
		if err = result.All(ctx, &allFoods); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allFoods)
	}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		foodId := c.Param("foodId")
		var food models.Food

		err := foodCollection.FindOne(ctx, bson.M{"foodId": foodId}).Decode(&food)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching food item"})
		}
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu
		var food models.Food

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(food)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		}

		err := menuCollection.FindOne(ctx, bson.M{"menuId": food.MenuId}).Decode(&menu)
		if err != nil {
			msg := fmt.Sprintf("menu not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		food.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		food.FoodId = food.ID.Hex()
		var num = toFixed(*food.Price, 2)
		food.Price = &num

		result, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil {
			msg := fmt.Sprintf("unable to create food item")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var menu models.Menu
		var food models.Food

		foodId := c.Param("foodId")

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateDoc primitive.D

		if food.Name != nil {
			updateDoc = append(updateDoc, bson.E{"name", food.Name})
		}
		if food.Price != nil {
			updateDoc = append(updateDoc, bson.E{"price", food.Price})
		}
		if food.FoodImage != nil {
			updateDoc = append(updateDoc, bson.E{"foodImage", food.FoodImage})
		}
		if food.MenuId != nil {
			err := menuCollection.FindOne(ctx, bson.M{"menuId": food.MenuId}).Decode(&menu)
			defer cancel()
			if err != nil {
				msg := fmt.Sprintf("message: Menu not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			updateDoc = append(updateDoc, bson.E{"menuId", food.MenuId})
		}

		food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateDoc = append(updateDoc, bson.E{"updatedAt", food.UpdatedAt})

		upsert := true
		filter := bson.M{"foodId": foodId}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := foodCollection.UpdateOne(ctx, filter, bson.D{
			{"$set", updateDoc},
		},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprint("food item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
