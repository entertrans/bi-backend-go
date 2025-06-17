package main

import (
	"github.com/entertrans/bi-backend-go/routers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	routers.RouterApp(app)

	app.Listen(":3000")
}
