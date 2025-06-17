package routers

import (
	"bi-backend-go/controllers"

	"github.com/gofiber/fiber/v2"
)

func routerApp(c *fiber.App) {
	c.Get("/", controllers.UserControllerShow)
}
