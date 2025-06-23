package main

import (
	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/routers"
)

func main() {
	// Koneksi DB
	config.ConnectDB()

	// Setup routes
	r := routers.SetupRouter()

	// Jalankan server di port 8080
	r.Run(":8080")
}
