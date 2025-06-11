package config

import (
	"crawl/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

// ConnectDatabase establishes a connection to the PostgreSQL database
func ConnectDatabase() {
	// Get the database URL from environment variables
	dsn := os.Getenv("DATABASE_URL")
	// Open the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	log.Println("📦 Connected to database")
}

// MigrateDatabase runs migrations for the necessary models, ensuring UUID handling
func MigrateDatabase() {
	// Perform auto migration for all necessary models
	err := DB.AutoMigrate(
		&models.Album{},    // Album model
		&models.User{},     // User model
		&models.Song{},     // Song model
		&models.Purchase{}, // Purchase model
		&models.SongMetrics{},
		// Add any other models that you want to auto-migrate
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("✅ Database migration successful")
}
