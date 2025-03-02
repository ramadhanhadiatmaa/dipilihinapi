package routes

import (
	"data/controllers"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	api := app.Group("/v1")

	type_user := api.Group("/typeuser")
	type_user.Post("/", controllers.CreateType)
	type_user.Delete("/:id", controllers.DeleteType)
	
	status := api.Group("/status")
	status.Post("/", controllers.CreateStatus)
	status.Get("/", controllers.ShowStatus)
	status.Get("/:id", controllers.IndexStatus)
	status.Put("/id", controllers.UpdateStatus)
	status.Delete("/:id", controllers.DeleteStatus)

	category := api.Group("/category")
	category.Post("/", controllers.CreateType)
	category.Delete("/:id", controllers.DeleteType)
}