package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	TaskID  uint   `json:"task_id"`
	UserID  uint   `json:"user_id"`
	User    User   `json:"user"`
	Content string `gorm:"not null" json:"content"`
}
