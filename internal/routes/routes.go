package routes

import (
	"tutuplapak/internal/handlers"
	"tutuplapak/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, healthHandler *handlers.HealthHandler, userHandler *handlers.UserHandler, registerHandler *handlers.RegisterHandler, loginHandler *handlers.LoginHandler, fileHandler *handlers.FileHandler) {
	// API version 1
	v1 := router.Group("/v1")
	{
		// Login & register routes
		login := v1.Group("/login")
		{
			login.POST("/email", loginHandler.LoginEmail)
		}

		v1.POST("/register", registerHandler.Register)
		// v1.POST("/login", loginHandler.Login)
		// v1.POST("/register", registerHandler.Register)
		v1.POST("/register/email", registerHandler.RegisterEmail)

		// Health check routes
		health := v1.Group("/health")
		{
			health.GET("/", healthHandler.Health)
			health.GET("/ready", healthHandler.Ready)
		}

		// User profile routes (auth required)
		userAuth := v1.Group("/user")
		userAuth.Use(middleware.IsAuthorized())
		{
			userAuth.GET("/", userHandler.GetUser)
			userAuth.PATCH("/", userHandler.UpdateUser)
		}

		// File upload routes (auth required)
		file := v1.Group("/file")
		file.Use(middleware.IsAuthorized())
		{
			file.POST("/", fileHandler.UploadFile)
			file.GET("/", fileHandler.GetUserFiles)
			file.DELETE("/", fileHandler.DeleteFile)
		}
	}

	// Root route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to TutupLapak API",
			"version": "1.0.0",
			"docs":    "/api/v1/health",
		})
	})
}
