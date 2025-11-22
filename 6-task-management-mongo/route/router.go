package router


import (
"github.com/gin-gonic/gin"
"task-management/controllers"
)


// SetupRouter returns a configured gin.Engine
func SetupRouter() *gin.Engine {
r := gin.Default()


api := r.Group("/tasks")
{
api.GET("", controllers.GetTasks)
api.GET(":id", controllers.GetTaskByID)
api.POST("", controllers.CreateTaskHandler)
api.PUT(":id", controllers.UpdateTaskHandler)
api.DELETE(":id", controllers.DeleteTaskHandler)
}


return r
}