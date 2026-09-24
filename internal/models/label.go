package models

import "gorm.io/gorm"

type Label struct {
	gorm.Model
	ProjectID uint   `json:"project_id"`
	Name      string `gorm:"not null" json:"name"`
	Color     string `gorm:"default:#4f8cff" json:"color"`
}
