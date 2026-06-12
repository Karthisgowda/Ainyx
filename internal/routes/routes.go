package routes

import (
	"github.com/Karthisgowda/Ainyx/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, userHandler *handler.UserHandler) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Ainyx Users API is running",
			"endpoints": []string{
				"GET /health",
				"POST /users",
				"GET /users",
				"GET /users/:id",
				"PUT /users/:id",
				"DELETE /users/:id",
			},
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Post("/users", userHandler.Create)
	app.Get("/users", userHandler.List)
	app.Get("/users/:id", userHandler.GetByID)
	app.Put("/users/:id", userHandler.Update)
	app.Delete("/users/:id", userHandler.Delete)
}
