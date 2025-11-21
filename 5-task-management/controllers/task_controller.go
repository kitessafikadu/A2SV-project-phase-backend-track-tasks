package controllers


import (
"net/http"
"strconv"
"time"


"github.com/gin-gonic/gin"
"task-management/data"
"task_manager/models"
)


// GetTasks handles GET /tasks
func GetTasks(c *gin.Context) {
tasks := data.GetAllTasks()
c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}


// GetTaskByID handles GET /tasks/:id
func GetTaskByID(c *gin.Context) {
idStr := c.Param("id")
id, err := strconv.Atoi(idStr)
if err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
return
}


t, err := data.GetTask(id)
if err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
return
}
c.JSON(http.StatusOK, t)
}


// CreateTaskHandler handles POST /tasks
func CreateTaskHandler(c *gin.Context) {
var input models.TaskInput
if err := c.ShouldBindJSON(&input); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}


// Validate due date format if provided
if input.DueDate != "" {
if _, err := time.Parse(time.RFC3339, input.DueDate); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_date format, use RFC3339"})
return
}
}


t, err := data.CreateTask(input)
if err != nil {
if err == data.ErrInvalidDate {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}
c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create task"})
return
}
}


// DeleteTaskHandler handles DELETE /tasks/:id
func DeleteTaskHandler(c *gin.Context) {
idStr := c.Param("id")
id, err := strconv.Atoi(idStr)
if err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
return
}


if err := data.DeleteTask(id); err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
return
}
c.Status(http.StatusNoContent)
}