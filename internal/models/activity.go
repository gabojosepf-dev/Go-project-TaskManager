package models

import "gorm.io/gorm"

type ActivityLog struct {
	gorm.Model
	ProjectID uint   `json:"project_id"`
	UserID    uint   `json:"user_id"`
	User      User   `json:"user"`
	Action    string `gorm:"not null" json:"action"`
}
