package main

import (
	"log"
	"net/http"

	"github.com/P4rz1val22/task-management-api/internal/database"
	"github.com/P4rz1val22/task-management-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	r := gin.Default()

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"message":  "Task Management API is running!",
			"database": "connected",
		})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
	}

	// Start server on port 8080
	log.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
