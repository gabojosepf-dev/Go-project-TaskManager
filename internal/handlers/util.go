package handlers

import (
	"strconv"

	"go-tasks-api/internal/models"

	"gorm.io/gorm"
)

func parseUint(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 64)
	return uint(v)
}

func uintToStr(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}

func logActivity(db *gorm.DB, projectID, userID uint, action string) {
	db.Create(&models.ActivityLog{ProjectID: projectID, UserID: userID, Action: action})
}

func taskWithProject(db *gorm.DB, taskID string) (models.Task, error) {
	var task models.Task
	err := db.First(&task, taskID).Error
	return task, err
}
