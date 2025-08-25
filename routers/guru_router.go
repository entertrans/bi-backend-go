package routers

import (
	guruhandlers "github.com/entertrans/bi-backend-go/handlers/guru"
	"github.com/gin-gonic/gin"
)

func RegisterGuruRoutes(r *gin.Engine) {
	// r.GET("/guru", adminHandlers.GetAllGuru)
	// dst...
	guruGroup := r.Group("/guru")
	{
		// Routes bank soal
		guruGroup.GET("/:guru_id/banksoal", guruhandlers.GetBankSoalHandler)
		guruGroup.GET("/banksoal/:kelas_id/:mapel_id", guruhandlers.GetActiveBankSoalByKelasMapelHandler)
		guruGroup.GET("/banksoal/rekap", guruhandlers.GetRekapBankSoalHandler)
		// buat bank soal
		guruGroup.POST("/banksoal/create", guruhandlers.BuatSoalHandler)
		// guruGroup.PATCH("/banksoal/:soal_id", guruhandlers.UpdateBankSoalHandler)
		guruGroup.DELETE("/banksoal/:soal_id", guruhandlers.DeleteBankSoalHandler)
		guruGroup.GET("/banksoal", guruhandlers.GetActiveBankSoalHandler)
		guruGroup.GET("/banksoal/inactive", guruhandlers.GetInactiveBankSoalHandler)
		guruGroup.PUT("/banksoal/:soal_id/restore", guruhandlers.RestoreBankSoalHandler)

		// Routes test online
		guruGroup.POST("/test", guruhandlers.CreateTestHandler)
		guruGroup.GET("/test/type/:type_test", guruhandlers.GetTestByType)
		guruGroup.GET("/test/guru/:guru_id", guruhandlers.GetTestsByGuruHandler)
		guruGroup.GET("/test/:test_id", guruhandlers.GetTestHandler)
		guruGroup.PUT("/test/:test_id", guruhandlers.UpdateTestHandler)
		guruGroup.DELETE("/test/:test_id", guruhandlers.DeleteTestHandler)

		// Routes penilaian
		// guruGroup.GET("/penilaian/:final_id", guruhandlers.GetPenilaianHandler)
		// guruGroup.POST("/penilaian", guruhandlers.CreatePenilaianHandler)
		// guruGroup.PATCH("/penilaian/:penilaian_id", guruhandlers.UpdatePenilaianHandler)
		// guruGroup.DELETE("/penilaian/:penilaian_id", guruhandlers.DeletePenilaianHandler)

		// routes

	}
	lampiran := r.Group("/lampiran")
	{
		lampiran.GET("/active", guruhandlers.GetActiveLampiranHandler)
		lampiran.GET("/trash", guruhandlers.GetInactiveLampiranHandler)
		lampiran.POST("/upload", guruhandlers.UploadLampiranHandler)
		lampiran.DELETE("/:lampiran_id", guruhandlers.DeleteLampiranHandler)
		lampiran.PUT("/restore/:lampiran_id", guruhandlers.RestoreLampiranHandler)
		lampiran.DELETE("/hard/:lampiran_id", guruhandlers.HardDeleteLampiranHandler)
	}
}
