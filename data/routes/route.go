package routes

import (
	"data/controllers"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	api := app.Group("/v1"/* , middlewares.Auth */)

	type_user := api.Group("/typeuser")
	type_user.Post("/", controllers.CreateType)
	type_user.Delete("/:id", controllers.DeleteType)
}