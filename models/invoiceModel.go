package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Invoice struct {
	ID             primitive.ObjectID `bson:"_id"`
	InvoiceId      string             `json:"invoiceId"`
	OrderId        string             `json:"orderId" validate:"required"`
	PaymentMethod  *string            `json:"paymentMethod" validate:"eq=CARD|eq=CASH|eq="`
	PaymentStatus  *string            `json:"paymentStatus" validate:"required,eq=PENDING|eq=PAID"`
	PaymentDueDate time.Time          `json:"paymentDueDate"`
	CreatedAt      time.Time          `json:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt"`
}
