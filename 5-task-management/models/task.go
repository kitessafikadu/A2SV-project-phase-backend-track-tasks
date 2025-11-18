package models

import "time"

// Task represents a task in the system
type Task struct {
ID int `json:"id"`
Title string `json:"title" binding:"required"`
Description string `json:"description"`
DueDate time.Time `json:"due_date" time_format:"2006-01-02T15:04:05Z07:00"`
Status string `json:"status" binding:"required"`
}


// TaskInput used for create/update payloads so we can validate fields
type TaskInput struct {
Title string `json:"title" binding:"required"`
Description string `json:"description"`
DueDate string `json:"due_date"` // ISO8601 string validated in service/controller
Status string `json:"status" binding:"required"`
}