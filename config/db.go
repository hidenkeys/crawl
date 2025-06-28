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
		&models.User{},
		&models.Role{},
		&models.Artist{},
		&models.Genre{},
		&models.Song{},
		&models.Album{},
		&models.SongContributor{},
		&models.AlbumContributor{},
		&models.SongPurchase{},
		&models.AlbumPurchase{},
		&models.UserFavorite{},
		&models.Playlist{},
		&models.PlaylistSong{},
		&models.ArtistTip{},
		&models.Stream{},
		&models.MonthlyRoyalty{},
		&models.ContentFlag{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database. \n", err)
		os.Exit(2)
	}

	log.Println("✅ Database migration successful")
}
