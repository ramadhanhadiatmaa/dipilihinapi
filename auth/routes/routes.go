package routes

import (
	"auth/controllers"
	"auth/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	api := app.Group("/v1")

	user := api.Group("/user")

	user.Post("/login", controllers.Login)
	user.Post("/register", controllers.Register)
	user.Get("/:id", controllers.Index, middlewares.Auth)
	user.Put("/:id", controllers.Update, middlewares.Auth)
	user.Delete("/:id", controllers.Delete, middlewares.Auth)
}