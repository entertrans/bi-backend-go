package routers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/entertrans/bi-backend-go/controllers"
)

func RouterApp(c *fiber.App) {
	c.Get("/", controllers.UserControllerShow)
	api := app.Group("/api")          // /api/...
	api.Get("/users", controllers.UserIndex) // Endpoint GET /api/users
}
