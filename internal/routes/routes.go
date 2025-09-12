package routes

import (
	"tutuplapak/internal/handlers"
	"tutuplapak/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, healthHandler *handlers.HealthHandler, userHandler *handlers.UserHandler, registerHandler *handlers.RegisterHandler, loginHandler *handlers.LoginHandler, fileHandler *handlers.FileHandler, productHandler *handlers.ProductHandler, purchaseHandler *handlers.PurchaseHandler) {
	// API version 1
	v1 := router.Group("/v1")
	{
		// Login & register routes
		login := v1.Group("/login")
		{
			login.POST("/phone", loginHandler.LoginPhone)
			login.POST("/email", loginHandler.LoginEmail)
		}

		register := v1.Group("/register")
		{
			register.POST("/email", registerHandler.RegisterEmail)
			register.POST("/phone", registerHandler.RegisterPhone)
		}

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
			userAuth.POST("/link/phone", userHandler.LinkPhone)
			userAuth.POST("/link/email", userHandler.LinkEmail)
			userAuth.PUT("/", userHandler.UpdateUser)
		}

		// File upload routes
		file := v1.Group("/file")
		{
			// Public endpoints - no auth required
			file.POST("/", fileHandler.UploadFile)
		}

		product := v1.Group("/product")
		{
			// Public endpoint - no auth required
			product.GET("/", productHandler.GetProducts)

			// Protected endpoints - auth required
			product.Use(middleware.IsAuthorized())
			{
				product.POST("/", productHandler.CreateProduct)
				product.PUT("/:productId", productHandler.UpdateProduct)
				product.DELETE("/:productId", productHandler.DeleteProduct)
			}
		}

		purchase := v1.Group("/purchase")
		purchase.Use(middleware.IsAuthorized())
		{
			purchase.POST("/", purchaseHandler.PurchaseProducts)
			purchase.POST("/:purchaseId", purchaseHandler.ProcessPurchase)
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
