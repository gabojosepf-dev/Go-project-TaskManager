package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"not null" json:"name"`
	Phone    string `json:"phone"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
}
