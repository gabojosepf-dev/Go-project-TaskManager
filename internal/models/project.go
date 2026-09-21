package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Name        string          `gorm:"not null" json:"name"`
	Description string          `json:"description"`
	OwnerID     uint            `json:"owner_id"`
	Members     []ProjectMember `json:"members,omitempty"`
	Tasks       []Task          `json:"tasks,omitempty"`
}

type ProjectMember struct {
	gorm.Model
	ProjectID uint   `gorm:"uniqueIndex:idx_project_user" json:"project_id"`
	UserID    uint   `gorm:"uniqueIndex:idx_project_user" json:"user_id"`
	Role      string `gorm:"default:member" json:"role"`
	User      User   `json:"user"`
}
