package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Note struct {
	ID        primitive.ObjectID `bson:"_id"`
	Text      string             `json:"text"`
	Title     string             `json:"title"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
	NoteId    string             `json:"noteId"`
}
