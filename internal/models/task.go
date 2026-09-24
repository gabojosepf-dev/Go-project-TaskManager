package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ProjectID    uint       `json:"project_id"`
	Title        string     `gorm:"not null" json:"title"`
	Description  string     `json:"description"`
	Status       string     `gorm:"default:todo" json:"status"`
	Priority     string     `gorm:"default:medium" json:"priority"`
	DueDate      *time.Time `json:"due_date"`
	CreatedByID  uint       `json:"created_by_id"`
	AssignedToID *uint      `json:"assigned_to_id"`
	AssignedTo   *User      `json:"assigned_to,omitempty"`
	Comments     []Comment  `json:"comments,omitempty"`
	Labels       []Label    `gorm:"many2many:task_labels;" json:"labels,omitempty"`
	Subtasks     []Subtask  `json:"subtasks,omitempty"`
}
