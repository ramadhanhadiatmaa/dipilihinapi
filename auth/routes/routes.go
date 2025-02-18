package routes

import (
	"auth/controllers"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	api := app.Group("/v1")

	user := api.Group("/user")

	user.Post("/login", controllers.Login)
	user.Post("/register", controllers.Register)
	user.Put("/:id", controllers.Update)
	user.Delete("/:id", controllers.Delete)
}