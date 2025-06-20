package main

import (
	"crawl/api"
	"crawl/config"
	"crawl/handlers"
	"crawl/repositories"
	"crawl/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	_ = godotenv.Load()
	//var jwtSecret = []byte("your-secret-key")
	config.ConnectDatabase()
	config.MigrateDatabase()

	db := config.DB

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	albumRepo := repositories.NewAlbumRepository(db)
	albumService := services.NewAlbumService(albumRepo)

	purchaseRepo := repositories.NewPurchaseRepository(db)
	purchaseService := services.NewPurchaseService(purchaseRepo)

	songRepo := repositories.NewSongRepository(db)
	songService := services.NewSongService(songRepo)

	songMetricsRepo := repositories.NewSongMetricsRepository(db)
	songMetricsService := services.NewSongMetricsService(songMetricsRepo)

	tipsRepo := repositories.NewTipRepository(db)
	tipsService := services.NewTipService(tipsRepo)

	server := handlers.NewServer(db, userService, albumService, purchaseService, songService, songMetricsService, tipsService)
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://crawl-app.vercel.app, https://crawl-admin.vercel.app/, https://crawl-app.vercel.app, http://localhost:3000",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	api.RegisterHandlers(app, server)

	// And we serve HTTP until the world ends.
	log.Fatal(app.Listen("0.0.0.0:8082"))

}
