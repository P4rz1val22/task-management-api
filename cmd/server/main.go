package main

import (
	"log"
	"net/http"

	"github.com/P4rz1val22/task-management-api/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to database
	database.Connect()

	// Create Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"message":  "Task Management API is running!",
			"database": "connected",
		})
	})

	// Start server on port 8080
	log.Println("Starting server on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
