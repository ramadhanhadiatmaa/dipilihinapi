package routes

import (
	"core/controllers"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes menginisialisasi semua route aplikasi
func SetupRoutes(app *fiber.App) {
	api := app.Group("/v1")

	// Route untuk mendapatkan data laptop berdasarkan id
	api.Get("/laptops", controllers.GetLaptops)
}