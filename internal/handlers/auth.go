package handlers

import (
	"github.com/P4rz1val22/task-management-api/internal/database"
	"github.com/P4rz1val22/task-management-api/internal/models"
	"github.com/P4rz1val22/task-management-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func Login(context *gin.Context) {
	var req LoginRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Get user from DB
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare passwords
	if err := utils.CheckPassword(req.Password, user.Password); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	// Provide token
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return success
	context.JSON(http.StatusOK, gin.H{"token": token,
		"message": "Successfully logged in",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func Register(context *gin.Context) {
	var req RegisterRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Checking for duplication of user
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		context.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// User creation
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
