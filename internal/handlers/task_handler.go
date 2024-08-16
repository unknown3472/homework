package handlers

import (
    "context"
    "net/http"
    "taskapp/internal/db"
    "taskapp/internal/models"
    "time"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTask(c *gin.Context) {
    var task models.Task
    if err := c.BindJSON(&task); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    task.ID = primitive.NewObjectID()
    task.CreatedAt = time.Now()

    collection := db.GetCollection("tasks")
    _, err := collection.InsertOne(context.Background(), task)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
        return
    }
    c.JSON(http.StatusCreated, task)
}

func GetTask(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    collection := db.GetCollection("tasks")
    var task models.Task
    err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&task)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }
    c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var task models.Task
    if err := c.BindJSON(&task); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    collection := db.GetCollection("tasks")
    update := bson.M{
        "$set": bson.M{
            "title":       task.Title,
            "description": task.Description,
            "status":      task.Status,
        },
    }

    _, err = collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func DeleteTask(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    collection := db.GetCollection("tasks")
    _, err = collection.DeleteOne(context.Background(), bson.M{"_id": objID})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
