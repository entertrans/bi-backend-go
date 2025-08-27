package routers

import (
	"time"

	"github.com/entertrans/bi-backend-go/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", "./uploads")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/login", handlers.LoginHandler)
	r.GET("/pengguna/aktif", handlers.GetActivePenggunaHandler)

	// 🔗 Panggil semua router modular
	RegisterSiswaRoutes(r)
	RegisterAdminRoutes(r)
	RegisterGuruRoutes(r) // nanti isi

	return r
}
