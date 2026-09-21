package main

import (
	"log"

	"go-tasks-api/internal/config"
	"go-tasks-api/internal/database"
	"go-tasks-api/internal/routes"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg)
	router := routes.Setup(db, cfg)

	log.Printf("server running on http://localhost:%s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
