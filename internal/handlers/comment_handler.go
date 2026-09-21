package handlers

import (
	"net/http"

	"go-tasks-api/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentHandler struct {
	DB *gorm.DB
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{DB: db}
}

type createCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *CommentHandler) GetComments(c *gin.Context) {
	userID := c.GetUint("user_id")
	taskID := c.Param("taskId")

	var task models.Task
	if err := h.DB.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	if _, ok := isMember(h.DB, task.ProjectID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed"})
		return
	}

	var comments []models.Comment
	h.DB.Preload("User").Where("task_id = ?", taskID).Order("created_at ASC").Find(&comments)

	c.JSON(http.StatusOK, comments)
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	taskID := c.Param("taskId")

	var task models.Task
	if err := h.DB.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	if _, ok := isMember(h.DB, task.ProjectID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed"})
		return
	}

	var req createCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := models.Comment{TaskID: parseUint(taskID), UserID: userID, Content: req.Content}
	if err := h.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating comment"})
		return
	}

	h.DB.Preload("User").First(&comment, comment.ID)
	c.JSON(http.StatusCreated, comment)
}
