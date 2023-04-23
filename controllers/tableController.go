package controllers

import (
	"github.com/Prodigy00/restaurant-mgt/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")
