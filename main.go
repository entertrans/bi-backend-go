package main

import (
	"github.com/gofiber/fiber/v2"

	"github.com/entertrans/bi-backend-go/config" // Panggil file konfigurasi
	"github.com/entertrans/bi-backend-go/routes" // Panggil router
)

func main() {
	// Inisialisasi Fiber App
	app := fiber.New()

	// Koneksi ke database
	config.ConnectDB()

	// Daftarkan semua route API
	routes.SetupRoutes(app)

	// Jalankan server di port 3000
	app.Listen(":3000")
}
