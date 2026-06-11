package routes

import (
	"github.com/Karthisgowda/Ainyx/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, userHandler *handler.UserHandler) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Post("/users", userHandler.Create)
	app.Get("/users", userHandler.List)
	app.Get("/users/:id", userHandler.GetByID)
	app.Put("/users/:id", userHandler.Update)
	app.Delete("/users/:id", userHandler.Delete)
}
