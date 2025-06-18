package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/entertrans/bi-backend-go/controllers"
)

func GetUsers(c *fiber.Ctx) error {
	users, err := controllers.GetAllUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Berhasil mengambil data user",
		"data":    users,
	})
}
