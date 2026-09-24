package models

import "gorm.io/gorm"

type Subtask struct {
	gorm.Model
	TaskID uint   `json:"task_id"`
	Title  string `gorm:"not null" json:"title"`
	Done   bool   `gorm:"default:false" json:"done"`
}
