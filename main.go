package main

import (
	"crawl/api"
	"crawl/config"
	"crawl/handlers"
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

	server := handlers.NewHandlers(db)
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://crawl-app.vercel.app, https://crawl-admin.vercel.app/",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	api.RegisterHandlers(app, server)

	// And we serve HTTP until the world ends.
	log.Fatal(app.Listen("0.0.0.0:8082"))

}
