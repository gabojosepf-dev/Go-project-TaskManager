package database

import (
	"fmt"
	"log"

	"go-tasks-api/internal/config"
	"go-tasks-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.ProjectMember{},
		&models.Task{},
		&models.Comment{},
		&models.Label{},
		&models.Subtask{},
		&models.ActivityLog{},
	); err != nil {
		log.Fatalf("migration error: %v", err)
	}

	return db
}
