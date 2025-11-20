package data

import (
	"errors"
	"sync"
	"time"

	"task_manager/models"
)


var (
ErrNotFound = errors.New("task not found")
ErrInvalidDate = errors.New("invalid date format, use RFC3339 (e.g. 2006-01-02T15:04:05Z07:00)")
)


// taskService holds in-memory tasks and sync primitives
type taskService struct {
mu sync.RWMutex
tasks map[int]models.Task
next int
}


var svc = &taskService{
tasks: make(map[int]models.Task),
next: 1,
}


// CreateTask creates a new task from TaskInput
func CreateTask(input models.TaskInput) (models.Task, error) {
var due time.Time
var err error
if input.DueDate != "" {
due, err = time.Parse(time.RFC3339, input.DueDate)
if err != nil {
return models.Task{}, ErrInvalidDate
}
}


svc.mu.Lock()
defer svc.mu.Unlock()


t := models.Task{
ID: svc.next,
Title: input.Title,
Description: input.Description,
DueDate: due,
Status: input.Status,
}


svc.tasks[svc.next] = t
svc.next++
return t, nil
}


// GetAllTasks returns a slice of all tasks
func GetAllTasks() []models.Task {
svc.mu.RLock()
defer svc.mu.RUnlock()


out := make([]models.Task, 0, len(svc.tasks))
for _, t := range svc.tasks {
out = append(out, t)
}
return out
}


// UpdateTask updates a task with new data
func UpdateTask(id int, input models.TaskInput) (models.Task, error) {
var due time.Time
var err error
if input.DueDate != "" {
due, err = time.Parse(time.RFC3339, input.DueDate)
if err != nil {
return models.Task{}, ErrInvalidDate
}
}


svc.mu.Lock()
defer svc.mu.Unlock()


t, ok := svc.tasks[id]
if !ok {
return models.Task{}, ErrNotFound
}


t.Title = input.Title
t.Description = input.Description
if !due.IsZero() {
t.DueDate = due
} else {
// if empty string passed, zero time will remain — keep existing
if input.DueDate == "" {
// preserve existing due date
} else {
t.DueDate = time.Time{}
}
}
t.Status = input.Status


svc.tasks[id] = t
return t, nil
}


// DeleteTask removes a task by id
func DeleteTask(id int) error {
svc.mu.Lock()
defer svc.mu.Unlock()


if _, ok := svc.tasks[id]; !ok {
return ErrNotFound
}
delete(svc.tasks, id)
return nil
}