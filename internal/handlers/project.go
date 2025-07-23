package handlers

import (
	"github.com/P4rz1val22/task-management-api/internal/database"
	"github.com/P4rz1val22/task-management-api/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// @Summary		Create a new project
// @Description	Create a new project owned by the authenticated user
// @Tags			projects
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			project	body		CreateProjectRequest	true	"Project creation data"
// @Success		201		{object}	map[string]interface{}
// @Failure		400		{object}	map[string]interface{}
// @Failure		401		{object}	map[string]interface{}
// @Router			/projects [post]
func CreateProject(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingProject models.Project
	if err := database.DB.Where("name = ? AND owner_id = ?", req.Name, userID).First(&existingProject).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Project already exists"})
		return
	}

	// Create a project
	project := models.Project{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     userID,
	}

	if err := database.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Project created successfully",
		"project": gin.H{
			"id":          project.ID,
			"name":        project.Name,
			"description": project.Description,
			"owner_id":    project.OwnerID,
			"owner_name":  "You", // Personal touch!
			"created_at":  project.CreatedAt,
		},
	})
}
