package routers

import (
	"time"

	"github.com/entertrans/bi-backend-go/handlers"
	"github.com/gin-contrib/cors" // 🔥 WAJIB
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// ✅ Middleware CORS yang lengkap
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ✅ Serve static file untuk foto siswa (jika ada)
	r.Static("/uploads", "./uploads")

	// ✅ Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// ✅ Routes
	//siswa
	r.GET("/siswa", handlers.GetAllSiswa)             // all siswa
	r.GET("/siswaaktif", handlers.GetAllSiswaAktif)   // all siswa aktif
	r.GET("/siswakeluar", handlers.GetAllSiswaKeluar) // all siswa aktif
	r.GET("/siswa/:nis", handlers.GetSiswaWithOrtu)   // detail siswa + ortu

	r.GET("/ortu", handlers.GetAllOrtu)         // semua ortu
	r.GET("/ortu/:nis", handlers.FindOrtuByNis) // ortu berdasarkan NIS

	//lookup
	r.GET("/lookup/agama", handlers.GetAllAgama)
	r.GET("/lookup/kelas", handlers.GetAllKelas)
	r.GET("/lookup/satelit", handlers.GetAllSatelit)

	return r
}
