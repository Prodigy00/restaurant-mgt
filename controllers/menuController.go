package controllers

import (
	"context"
	"fmt"
	"github.com/Prodigy00/restaurant-mgt/database"
	"github.com/Prodigy00/restaurant-mgt/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := menuCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing the menu items"})
		}

		var allMenus []bson.M
		if err = result.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenus)
	}
}
func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		menuId := c.Param("menuId")
		var menu models.Menu

		err := foodCollection.FindOne(ctx, bson.M{"menuId": menuId}).Decode(&menu)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching the menu"})
		}
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.Menu
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(menu)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		}

		menu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.MenuId = menu.ID.Hex()

		result, insertErr := menuCollection.InsertOne(ctx, menu)
		if insertErr != nil {
			msg := fmt.Sprintf("Menu item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.Menu
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		menuId := c.Param("menuId")
		filter := bson.M{"menuId": menuId}

		var updateDoc primitive.D

		if !menu.StartDate.IsZero() && !menu.EndDate.IsZero() {
			if !inTimeSpan(menu.StartDate, menu.EndDate, time.Now()) {
				msg := "kindly retype time"
				c.JSON(http.StatusBadRequest, gin.H{"error": msg})
				defer cancel()
				return
			}
		}
		updateDoc = append(updateDoc, bson.E{"startDate", menu.StartDate})
		updateDoc = append(updateDoc, bson.E{"endDate", menu.EndDate})

		if menu.Name != "" {
			updateDoc = append(updateDoc, bson.E{"name", menu.Name})
		}

		if menu.Category != "" {
			updateDoc = append(updateDoc, bson.E{"category", menu.Category})
		}

		menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateDoc = append(updateDoc, bson.E{"updatedAt", menu.UpdatedAt})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := menuCollection.UpdateOne(
			ctx,
			filter,
			bson.D{{"$set", updateDoc}},
			&opt,
		)

		if err != nil {
			msg := "Menu update failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		c.JSON(http.StatusOK, result)
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return start.After(time.Now()) && end.After(start)
}
