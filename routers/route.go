package routers

import (
	"github.com/entertrans/bi-backend-go/controllers"

	"github.com/gofiber/fiber/v2"
)

func RouterApp(c *fiber.App) {
	c.Get("/", controllers.UserControllerShow)
}
