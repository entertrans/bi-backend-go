package routers

import (
	"github.com/entertrans/bi-backend-go/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Endpoint testing
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	// Endpoint
	//siswa
	r.GET("/siswa", handlers.GetAllSiswa)         //allsiswa
	r.GET("/siswa/:nis", handlers.FindSiswaByNis) // get berdasarkan siswa_nis
	r.GET("/siswa/:nis/detail", handlers.GetSiswaWithOrtu)

	//ortu
	r.GET("/ortu", handlers.GetAllOrtu)         //semua ortu
	r.GET("/ortu/:nis", handlers.FindOrtuByNis) // get berdasarkan ortu_nis

	//agama
	r.GET("/agama", handlers.GetAllAgama)

	return r
}
