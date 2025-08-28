package routers

import (
	siswaHandler "github.com/entertrans/bi-backend-go/handlers/siswa"
	"github.com/gin-gonic/gin"
)

func RegisterSiswaRoutes(r *gin.Engine) {
	// r.GET("/guru", adminHandlers.GetAllGuru)
	// dst...

	siswaRoutes := r.Group("/siswa")
	{
		siswaRoutes.POST("/test/start/:test_id", siswaHandler.StartTestHandler)
		siswaRoutes.GET("/test/:test_id/session", siswaHandler.GetTestSessionHandler)
		siswaRoutes.POST("/test/submit/:session_id", siswaHandler.SubmitTestHandler)
		siswaRoutes.GET("/test/session/:session_id", siswaHandler.GetSessionByIDHandler)
		siswaRoutes.GET("/tests/ub", siswaHandler.GetAllUBTestHandler)
		siswaRoutes.GET("/tests/ub/kelas/:kelas_id", siswaHandler.GetUBTestByKelasHandler) // Tambahkan ini
		siswaRoutes.GET("/tests/:id/soal", siswaHandler.GetSoalByTestIDHandler)
		siswaRoutes.POST("/jawaban/save", siswaHandler.SaveJawabanHandler)
		siswaRoutes.GET("/test/:test_id/soal", siswaHandler.GetSoalHandler)
	}
}
