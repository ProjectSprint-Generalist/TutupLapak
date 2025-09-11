package database

import (
	"log"
	"time"

	"tutuplapak/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect establishes a connection to the PostgreSQL database
func Connect(databaseURL string) error {
	var err error
	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return err
	}

	// Connection pool
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxOpenConns(300)
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)

	log.Println("Connected to database successfully")
	return nil
}

// Migrate runs the database migrations
func Migrate() error {
	log.Println("Running database migrations...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.FileUpload{},
		&models.Product{},
		&models.Purchase{},
		&models.PurchaseItem{},
	)
	if err != nil {
		log.Printf("Migration error: %v", err)
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
