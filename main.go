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

	//test struktur model
	// var agama []models.Agama
	// err := config.DB.Debug().Find(&agama).Error
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(agama)

	// Jalankan server
	r.Run(":8080")
}
