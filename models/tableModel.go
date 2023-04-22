package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Table struct {
	ID             primitive.ObjectID `bson:"_id"`
	NumberOfGuests *int               `json:"numberOfGuests" validate:"required"`
	TableNumber    *int               `json:"tableNumber" validate:"required"`
	CreatedAt      time.Time          `json:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt"`
	TableId        string             `json:"tableId"`
}
