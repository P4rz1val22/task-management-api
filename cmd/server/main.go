package main

import (
	"github.com/P4rz1val22/task-management-api/internal/middleware"
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
		auth.POST("/login", handlers.Login)
	}

	// Add this AFTER your auth routes, BEFORE r.Run():
	protected := r.Group("/users")
	protected.Use(middleware.RequireAuth()) // Apply middleware to all /users routes
	{
		protected.GET("/me", func(c *gin.Context) {
			// Get user info from middleware
			userID := c.GetUint("user_id")
			email := c.GetString("email")

			c.JSON(http.StatusOK, gin.H{
				"message": "Protected route working!",
				"user_id": userID,
				"email":   email,
			})
		})
	}

	// Start server on port 8080
	log.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
