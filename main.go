package main

import (
	"github.com/entertrans/bi-backend-go/config" // Ganti dengan nama module kamu di go.mod
	"github.com/entertrans/bi-backend-go/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB() // koneksi database
	// koneksi ke DB
	config.ConnectDB()

	// setup routes
	r := routers.SetupRouter()

	//test ping
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// run server
	r.Run(":8080")
}
