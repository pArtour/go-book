package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Book is the model for a book
type Book struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Author      *string            `json:"author" validate:"required"`
	Title       *string            `json:"title" validate:"required"`
	Description *string            `json:"description" validate:"required"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
}
