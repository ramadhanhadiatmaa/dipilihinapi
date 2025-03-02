package routes

import (
	"chat/controllers"
	"chat/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	api := app.Group("/v1")

	chat := api.Group("/chat")
	chat.Post("/", controllers.Create, middlewares.Auth)
	chat.Get("/", controllers.Show, middlewares.Auth)
	chat.Get("/:id", controllers.Index, middlewares.Auth)
	chat.Delete("/:id", controllers.Delete, middlewares.Auth)
}
