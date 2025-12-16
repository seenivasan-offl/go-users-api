package routes

import (
	"github.com/gofiber/fiber/v2"

	"go-users-api/internal/handler"
)

func Register(app *fiber.App, uh *handler.UserHandler) {
	app.Post("/users", uh.CreateUser)
	app.Get("/users/:id", uh.GetUser)
	app.Put("/users/:id", uh.UpdateUser)
	app.Delete("/users/:id", uh.DeleteUser)
	app.Get("/users", uh.ListUsers)
}
