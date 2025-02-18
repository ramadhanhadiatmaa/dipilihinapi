package main

import (
	"data/models"
	"data/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	models.ConnectDatabase()

	port := os.Getenv("PORT")
	if port == "" {
		port = "7302"
	}

	app := fiber.New()

	app.Use(cors.New())

	routes.Route(app)

	app.Listen(":" + port)
}