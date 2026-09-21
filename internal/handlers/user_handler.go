package handlers

import (
	"net/http"

	"go-tasks-api/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

type updateProfileRequest struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone"`
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := h.DB.Select("id", "name", "phone", "email").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.Name = req.Name
	user.Phone = req.Phone

	if err := h.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update error"})
		return
	}

	c.JSON(http.StatusOK, user)
}
