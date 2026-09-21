package handlers

import (
	"net/http"

	"go-tasks-api/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectHandler struct {
	DB *gorm.DB
}

func NewProjectHandler(db *gorm.DB) *ProjectHandler {
	return &ProjectHandler{DB: db}
}

func isMember(db *gorm.DB, projectID, userID uint) (models.ProjectMember, bool) {
	var member models.ProjectMember
	err := db.Where("project_id = ? AND user_id = ?", projectID, userID).First(&member).Error
	return member, err == nil
}

type createProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req createProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project := models.Project{Name: req.Name, Description: req.Description, OwnerID: userID}
	if err := h.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating project"})
		return
	}

	member := models.ProjectMember{ProjectID: project.ID, UserID: userID, Role: "owner"}
	if err := h.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating membership"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

func (h *ProjectHandler) GetProjects(c *gin.Context) {
	userID := c.GetUint("user_id")

	var memberships []models.ProjectMember
	h.DB.Where("user_id = ?", userID).Find(&memberships)

	projectIDs := make([]uint, 0, len(memberships))
	for _, m := range memberships {
		projectIDs = append(projectIDs, m.ProjectID)
	}

	var projects []models.Project
	if err := h.DB.Where("id IN ?", projectIDs).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching projects"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	userID := c.GetUint("user_id")
	projectID := c.Param("id")

	if _, ok := isMember(h.DB, parseUint(projectID), userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member of this project"})
		return
	}

	var project models.Project
	if err := h.DB.Preload("Members.User").First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

type addMemberRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *ProjectHandler) AddMember(c *gin.Context) {
	userID := c.GetUint("user_id")
	projectID := parseUint(c.Param("id"))

	member, ok := isMember(h.DB, projectID, userID)
	if !ok || member.Role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the owner can add members"})
		return
	}

	var req addMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no user with that email"})
		return
	}

	newMember := models.ProjectMember{ProjectID: projectID, UserID: user.ID, Role: "member"}
	if err := h.DB.Create(&newMember).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user is already a member"})
		return
	}

	c.JSON(http.StatusCreated, newMember)
}

func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	userID := c.GetUint("user_id")
	projectID := parseUint(c.Param("id"))
	targetUserID := parseUint(c.Param("userId"))

	member, ok := isMember(h.DB, projectID, userID)
	if !ok || member.Role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the owner can remove members"})
		return
	}

	if targetUserID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "owner cannot remove themselves"})
		return
	}

	h.DB.Where("project_id = ? AND user_id = ?", projectID, targetUserID).Delete(&models.ProjectMember{})
	c.JSON(http.StatusOK, gin.H{"message": "member removed"})
}
