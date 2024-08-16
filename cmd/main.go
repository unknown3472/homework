package main

import (
    "log"
    "net/http"
    "taskapp/internal/handlers"
    "taskapp/internal/db"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    db.ConnectDB()
    r.POST("/tasks", handlers.CreateTask)
    r.GET("/tasks/:id", handlers.GetTask)
    r.PUT("/tasks/:id", handlers.UpdateTask)
    r.DELETE("/tasks/:id", handlers.DeleteTask)

    log.Fatal(http.ListenAndServe(":8080", r))
}
