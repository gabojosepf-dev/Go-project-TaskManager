package handlers

import (
	"net/http"

	"go-tasks-api/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LabelHandler struct {
	DB *gorm.DB
}

func NewLabelHandler(db *gorm.DB) *LabelHandler {
	return &LabelHandler{DB: db}
}

type createLabelRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
}

func (h *LabelHandler) GetLabels(c *gin.Context) {
	userID := c.GetUint("user_id")
	projectID := parseUint(c.Param("id"))

	if _, ok := isMember(h.DB, projectID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member of this project"})
		return
	}

	var labels []models.Label
	h.DB.Where("project_id = ?", projectID).Find(&labels)
	c.JSON(http.StatusOK, labels)
}

func (h *LabelHandler) CreateLabel(c *gin.Context) {
	userID := c.GetUint("user_id")
	projectID := parseUint(c.Param("id"))

	if _, ok := isMember(h.DB, projectID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member of this project"})
		return
	}

	var req createLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	color := req.Color
	if color == "" {
		color = "#4f8cff"
	}

	label := models.Label{ProjectID: projectID, Name: req.Name, Color: color}
	if err := h.DB.Create(&label).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating label"})
		return
	}

	c.JSON(http.StatusCreated, label)
}

type attachLabelRequest struct {
	LabelID uint `json:"label_id" binding:"required"`
}

func (h *LabelHandler) AttachLabel(c *gin.Context) {
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

	var req attachLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var label models.Label
	if err := h.DB.First(&label, req.LabelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "label not found"})
		return
	}

	if err := h.DB.Model(&task).Association("Labels").Append(&label); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error attaching label"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "label attached"})
}

func (h *LabelHandler) DetachLabel(c *gin.Context) {
	userID := c.GetUint("user_id")
	taskID := c.Param("taskId")
	labelID := parseUint(c.Param("labelId"))

	task, err := taskWithProject(h.DB, taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	if _, ok := isMember(h.DB, task.ProjectID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed"})
		return
	}

	label := models.Label{}
	label.ID = labelID
	h.DB.Model(&task).Association("Labels").Delete(&label)

	c.JSON(http.StatusOK, gin.H{"message": "label detached"})
}
