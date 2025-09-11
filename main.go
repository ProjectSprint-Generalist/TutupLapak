package main

import (
	"log"
	"os"

	"tutuplapak/internal/config"
	"tutuplapak/internal/database"
	"tutuplapak/internal/handlers"
	"tutuplapak/internal/middleware"
	"tutuplapak/internal/routes"
	"tutuplapak/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run database migrations
	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to run database migrations:", err)
	}

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// Initialize MinIO service (optional for local development)
	minioService, err := services.NewMinIOService(cfg, database.DB)
	if err != nil {
		log.Printf("Warning: MinIO service not available: %v", err)
		log.Println("File upload functionality will be disabled")
		minioService = nil
	}

	// Initialize handlers with database connection
	healthHandler := handlers.NewHealthHandler()
	userHandler := handlers.NewUserHandler(database.DB)
	registerHandler := handlers.NewRegisterHandler(database.DB)
	loginHandler := handlers.NewLoginHandler(database.DB)
	fileHandler := handlers.NewFileHandler(minioService)
	productHandler := handlers.NewProductHandler(database.DB)
	purchaseHandler := handlers.NewPurchaseHandler(database.DB)

	// Setup routes
	routes.SetupRoutes(router, healthHandler, userHandler, registerHandler, loginHandler, fileHandler, productHandler, purchaseHandler)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
