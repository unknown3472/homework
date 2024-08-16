package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"taskapp/internal/db"
	"taskapp/internal/models"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCreateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db.ConnectDB()                                             
	defer db.GetCollection("tasks").Drop(context.Background())

	router.POST("/tasks", CreateTask)

	task := models.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "pending",
	}

	taskJSON, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdTask models.Task
	json.Unmarshal(w.Body.Bytes(), &createdTask)

	assert.Equal(t, task.Title, createdTask.Title)
	assert.Equal(t, task.Description, createdTask.Description)
	assert.Equal(t, task.Status, createdTask.Status)
}

func TestGetTask(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db.ConnectDB()                                             
	defer db.GetCollection("tasks").Drop(context.Background()) 

	router.GET("/tasks/:id", GetTask)

	task := models.Task{
		ID:          primitive.NewObjectID(),
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "pending",
		CreatedAt:   time.Now(),
	}
	_, _ = db.GetCollection("tasks").InsertOne(context.Background(), task)

	req, _ := http.NewRequest("GET", "/tasks/"+task.ID.Hex(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var retrievedTask models.Task
	json.Unmarshal(w.Body.Bytes(), &retrievedTask)

	assert.Equal(t, task.Title, retrievedTask.Title)
	assert.Equal(t, task.Description, retrievedTask.Description)
	assert.Equal(t, task.Status, retrievedTask.Status)
}

func TestUpdateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db.ConnectDB()                                            
	defer db.GetCollection("tasks").Drop(context.Background()) 

	router.PUT("/tasks/:id", UpdateTask)

	task := models.Task{
		ID:          primitive.NewObjectID(),
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "pending",
		CreatedAt:   time.Now(),
	}
	_, _ = db.GetCollection("tasks").InsertOne(context.Background(), task)

	updatedTask := models.Task{
		Title:       "Updated Task",
		Description: "This is an updated task",
		Status:      "completed",
	}
	updatedTaskJSON, _ := json.Marshal(updatedTask)

	req, _ := http.NewRequest("PUT", "/tasks/"+task.ID.Hex(), bytes.NewBuffer(updatedTaskJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Task updated successfully")

	var retrievedTask models.Task
	_ = db.GetCollection("tasks").FindOne(context.Background(), bson.M{"_id": task.ID}).Decode(&retrievedTask)

	assert.Equal(t, updatedTask.Title, retrievedTask.Title)
	assert.Equal(t, updatedTask.Description, retrievedTask.Description)
	assert.Equal(t, updatedTask.Status, retrievedTask.Status)
}

func TestDeleteTask(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db.ConnectDB()                                             
	defer db.GetCollection("tasks").Drop(context.Background()) 

	router.DELETE("/tasks/:id", DeleteTask)

	task := models.Task{
		ID:          primitive.NewObjectID(),
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "pending",
		CreatedAt:   time.Now(),
	}
	_, _ = db.GetCollection("tasks").InsertOne(context.Background(), task)

	req, _ := http.NewRequest("DELETE", "/tasks/"+task.ID.Hex(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Task deleted successfully")

	err := db.GetCollection("tasks").FindOne(context.Background(), bson.M{"_id": task.ID}).Err()
	assert.Equal(t, err, mongo.ErrNoDocuments)
}
