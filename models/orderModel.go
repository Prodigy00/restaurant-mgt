package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Order struct {
	ID        primitive.ObjectID `bson:"_id"`
	OrderDate time.Time          `json:"orderDate" validate:"required"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
	OrderId   string             `json:"orderId"`
	TableId   *string            `json:"tableId"`
}
