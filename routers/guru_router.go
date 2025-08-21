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
		// buat bank soal
		// guruGroup.POST("/banksoal", guruhandlers.CreateBankSoalHandler)
		// guruGroup.PATCH("/banksoal/:soal_id", guruhandlers.UpdateBankSoalHandler)
		guruGroup.DELETE("/banksoal/:soal_id", guruhandlers.DeleteBankSoalHandler)


		// Routes test online
		guruGroup.GET("/test/:test_id", guruhandlers.GetTestHandler)
		guruGroup.POST("/test", guruhandlers.CreateTestHandler)
		guruGroup.PATCH("/test/:test_id", guruhandlers.UpdateTestHandler)
		guruGroup.DELETE("/test/:test_id", guruhandlers.DeleteTestHandler)

		// Routes penilaian
		guruGroup.GET("/penilaian/:final_id", guruhandlers.GetPenilaianHandler)
		guruGroup.POST("/penilaian", guruhandlers.CreatePenilaianHandler)
		guruGroup.PATCH("/penilaian/:penilaian_id", guruhandlers.UpdatePenilaianHandler)
		guruGroup.DELETE("/penilaian/:penilaian_id", guruhandlers.DeletePenilaianHandler)

		// routes
		guruGroup.GET("/banksoal", guruhandlers.GetActiveBankSoalHandler)
		guruGroup.GET("/banksoal/inactive", guruhandlers.GetInactiveBankSoalHandler)
		guruGroup.PATCH("/banksoal/:soal_id/restore", guruhandlers.RestoreBankSoalHandler)

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
