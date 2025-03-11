package routes

import (
	"chat/controllers"
	"chat/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	api := app.Group("/v1", middlewares.Auth)

	chat := api.Group("/chat")
	chat.Post("/", controllers.Create)
	chat.Get("/", controllers.Show)
	chat.Get("/:id", controllers.Index)
	chat.Delete("/:id", controllers.Delete)
}
