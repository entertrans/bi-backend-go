package routers

import (
	"github.com/entertrans/bi-backend-go/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/siswa", handlers.GetAllSiswa)

	return r
}
