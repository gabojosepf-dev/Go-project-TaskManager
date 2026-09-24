package handlers

import (
	"fmt"
	"net/http"
	"time"

	"go-tasks-api/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskHandler struct {
	DB *gorm.DB
}

func NewTaskHandler(db *gorm.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

var validStatuses = map[string]bool{"todo": true, "in_progress": true, "done": true}
var validPriorities = map[string]bool{"low": true, "medium": true, "high": true, "urgent": true}

type createTaskRequest struct {
	Title        string     `json:"title" binding:"required"`
	Description  string     `json:"description"`
	Priority     string     `json:"priority"`
	DueDate      *time.Time `json:"due_date"`
	AssignedToID *uint      `json:"assigned_to_id"`
}

type updateTaskRequest struct {
	Title        *string    `json:"title"`
	Description  *string    `json:"description"`
	Status       *string    `json:"status"`
	Priority     *string    `json:"priority"`
	DueDate      *time.Time `json:"due_date"`
	AssignedToID *uint      `json:"assigned_to_id"`
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	userID := c.GetUint("user_id")
	projectID := parseUint(c.Param("id"))

	if _, ok := isMember(h.DB, projectID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member of this project"})
		return
	}

	query := h.DB.Preload("AssignedTo").Preload("Labels").Preload("Subtasks").Where("project_id = ?", projectID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if priority := c.Query("priority"); priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if assignedTo := c.Query("assigned_to_id"); assignedTo != "" {
		query = query.Where("assigned_to_id = ?", assignedTo)
	}
	if search := c.Query("search"); search != "" {
		like := fmt.Sprintf("%%%s%%", search)
		query = query.Where("title ILIKE ? OR description ILIKE ?", like, like)
	}

	var tasks []models.Task
	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetTrash(c *gin.Context) {
	userID := c.GetUint("user_id")
	projectID := parseUint(c.Param("id"))

	if _, ok := isMember(h.DB, projectID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member of this project"})
		return
	}

	var tasks []models.Task
	h.DB.Unscoped().Where("project_id = ? AND deleted_at IS NOT NULL", projectID).Find(&tasks)

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) RestoreTask(c *gin.Context) {
	userID := c.GetUint("user_id")
	taskID := c.Param("taskId")

	var task models.Task
	if err := h.DB.Unscoped().First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	if _, ok := isMember(h.DB, task.ProjectID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed"})
		return
	}

	h.DB.Unscoped().Model(&task).Update("deleted_at", nil)
	logActivity(h.DB, task.ProjectID, userID, fmt.Sprintf("restauró la tarea \"%s\"", task.Title))

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID := c.GetUint("user_id")
	projectID := parseUint(c.Param("id"))

	if _, ok := isMember(h.DB, projectID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member of this project"})
		return
	}

	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	priority := req.Priority
	if priority == "" || !validPriorities[priority] {
		priority = "medium"
	}

	task := models.Task{
		ProjectID:    projectID,
		Title:        req.Title,
		Description:  req.Description,
		Status:       "todo",
		Priority:     priority,
		DueDate:      req.DueDate,
		CreatedByID:  userID,
		AssignedToID: req.AssignedToID,
	}

	if err := h.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating task"})
		return
	}

	logActivity(h.DB, projectID, userID, fmt.Sprintf("creó la tarea \"%s\"", task.Title))
	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) taskProjectMember(taskID string, userID uint) (models.Task, bool) {
	task, err := taskWithProject(h.DB, taskID)
	if err != nil {
		return task, false
	}
	_, ok := isMember(h.DB, task.ProjectID, userID)
	return task, ok
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userID := c.GetUint("user_id")
	taskID := c.Param("taskId")

	task, ok := h.taskProjectMember(taskID, userID)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to modify this task"})
		return
	}

	var req updateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil && validStatuses[*req.Status] && *req.Status != task.Status {
		logActivity(h.DB, task.ProjectID, userID, fmt.Sprintf("movió \"%s\" a %s", task.Title, *req.Status))
		task.Status = *req.Status
	}
	if req.Priority != nil && validPriorities[*req.Priority] {
		task.Priority = *req.Priority
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	if req.AssignedToID != nil {
		task.AssignedToID = req.AssignedToID
	}

	if err := h.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userID := c.GetUint("user_id")
	taskID := c.Param("taskId")

	task, ok := h.taskProjectMember(taskID, userID)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to delete this task"})
		return
	}

	h.DB.Delete(&models.Task{}, taskID)
	logActivity(h.DB, task.ProjectID, userID, fmt.Sprintf("eliminó la tarea \"%s\"", task.Title))

	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}
