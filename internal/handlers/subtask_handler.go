package handlers

import (
	"net/http"

	"go-tasks-api/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SubtaskHandler struct {
	DB *gorm.DB
}

func NewSubtaskHandler(db *gorm.DB) *SubtaskHandler {
	return &SubtaskHandler{DB: db}
}

type createSubtaskRequest struct {
	Title string `json:"title" binding:"required"`
}

type updateSubtaskRequest struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}

func (h *SubtaskHandler) GetSubtasks(c *gin.Context) {
	userID := c.GetUint("user_id")
	taskID := c.Param("taskId")

	task, err := taskWithProject(h.DB, taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	if _, ok := isMember(h.DB, task.ProjectID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed"})
		return
	}

	var subtasks []models.Subtask
	h.DB.Where("task_id = ?", taskID).Order("created_at ASC").Find(&subtasks)
	c.JSON(http.StatusOK, subtasks)
}

func (h *SubtaskHandler) CreateSubtask(c *gin.Context) {
	userID := c.GetUint("user_id")
	taskID := c.Param("taskId")

	task, err := taskWithProject(h.DB, taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	if _, ok := isMember(h.DB, task.ProjectID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed"})
		return
	}

	var req createSubtaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subtask := models.Subtask{TaskID: parseUint(taskID), Title: req.Title}
	if err := h.DB.Create(&subtask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating subtask"})
		return
	}

	c.JSON(http.StatusCreated, subtask)
}

func (h *SubtaskHandler) subtaskWithTaskProject(subtaskID string) (models.Subtask, models.Task, error) {
	var subtask models.Subtask
	if err := h.DB.First(&subtask, subtaskID).Error; err != nil {
		return subtask, models.Task{}, err
	}
	task, err := taskWithProject(h.DB, uintToStr(subtask.TaskID))
	return subtask, task, err
}

func (h *SubtaskHandler) UpdateSubtask(c *gin.Context) {
	userID := c.GetUint("user_id")
	subtaskID := c.Param("id")

	subtask, task, err := h.subtaskWithTaskProject(subtaskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subtask not found"})
		return
	}
	if _, ok := isMember(h.DB, task.ProjectID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed"})
		return
	}

	var req updateSubtaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title != nil {
		subtask.Title = *req.Title
	}
	if req.Done != nil {
		subtask.Done = *req.Done
	}

	h.DB.Save(&subtask)
	c.JSON(http.StatusOK, subtask)
}

func (h *SubtaskHandler) DeleteSubtask(c *gin.Context) {
	userID := c.GetUint("user_id")
	subtaskID := c.Param("id")

	_, task, err := h.subtaskWithTaskProject(subtaskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subtask not found"})
		return
	}
	if _, ok := isMember(h.DB, task.ProjectID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed"})
		return
	}

	h.DB.Delete(&models.Subtask{}, subtaskID)
	c.JSON(http.StatusOK, gin.H{"message": "subtask deleted"})
}
