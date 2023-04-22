package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OrderItem struct {
	ID          primitive.ObjectID `bson:"_id"`
	Quantity    *string            `json:"quantity" validate:"required,eq=S|eq=M|eq=L"`
	UnitPrice   *float64           `json:"unitPrice" validate:"required"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	OrderItemId string             `json:"orderItemId"`
	OrderId     string             `json:"orderId" validate:"required"`
	FoodId      *string            `json:"foodId" validate:"required"`
}
