package routes

import (
	"time"

	"go-tasks-api/internal/config"
	"go-tasks-api/internal/handlers"
	"go-tasks-api/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authHandler := handlers.NewAuthHandler(db, cfg)
	userHandler := handlers.NewUserHandler(db)
	projectHandler := handlers.NewProjectHandler(db)
	taskHandler := handlers.NewTaskHandler(db)
	commentHandler := handlers.NewCommentHandler(db)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	users := router.Group("/users")
	users.Use(middleware.AuthRequired(cfg.JWTSecret))
	{
		users.GET("/me", userHandler.GetMe)
		users.PUT("/me", userHandler.UpdateMe)
	}

	projects := router.Group("/projects")
	projects.Use(middleware.AuthRequired(cfg.JWTSecret))
	{
		projects.POST("", projectHandler.CreateProject)
		projects.GET("", projectHandler.GetProjects)
		projects.GET("/:id", projectHandler.GetProject)
		projects.POST("/:id/members", projectHandler.AddMember)
		projects.DELETE("/:id/members/:userId", projectHandler.RemoveMember)
		projects.GET("/:id/tasks", taskHandler.GetTasks)
		projects.POST("/:id/tasks", taskHandler.CreateTask)
	}

	tasks := router.Group("/tasks")
	tasks.Use(middleware.AuthRequired(cfg.JWTSecret))
	{
		tasks.PUT("/:taskId", taskHandler.UpdateTask)
		tasks.DELETE("/:taskId", taskHandler.DeleteTask)
		tasks.GET("/:taskId/comments", commentHandler.GetComments)
		tasks.POST("/:taskId/comments", commentHandler.CreateComment)
	}

	return router
}
