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

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := orderCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing order items"})
		}

		var allOrders []bson.M
		if err = result.All(ctx, &allOrders); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allOrders)
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		orderId := c.Param("orderId")
		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"orderId": orderId}).Decode(&order)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching order item"})
		}
		c.JSON(http.StatusOK, order)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var table models.Table
		var order models.Order

		var updateDoc primitive.D
		orderId := c.Param("orderId")

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if order.TableId != nil {
			err := tableCollection.FindOne(ctx, bson.M{"tableId": order.TableId}).Decode(&table)
			defer cancel()
			if err != nil {
				msg := fmt.Sprintf("message: Table not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			updateDoc = append(updateDoc, bson.E{"tableId", order.TableId})
		}

		order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateDoc = append(updateDoc, bson.E{"updatedAt", order.UpdatedAt})

		upsert := true
		filter := bson.M{"orderId": orderId}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderCollection.UpdateOne(ctx, filter, bson.D{
			{"$set", updateDoc},
		},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprint("order item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
