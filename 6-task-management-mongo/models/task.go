package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Title       string             `json:"title" binding:"required"`
	Description string             `json:"description"`
	DueDate     time.Time          `json:"due_date"`
	Status      string             `json:"status" binding:"required"`
}

type TaskInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"` // ISO8601 string
	Status      string `json:"status" binding:"required"`
}
