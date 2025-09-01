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
		siswaRoutes.GET("/session/:session_id/soal", siswaHandler.GetSessionSoalHandler)
		siswaRoutes.GET("/test/:test_id/active-session", siswaHandler.GetActiveTestSessionHandler)
		siswaRoutes.GET("/test/:test_id/session", siswaHandler.GetTestSessionHandler)
		siswaRoutes.POST("/test/submit/:session_id", siswaHandler.SubmitTestHandler)
		siswaRoutes.POST("/tugas/submit/:session_id", siswaHandler.SubmitTugasHandler)
		siswaRoutes.GET("/test/session/:session_id", siswaHandler.GetSessionByIDHandler)

		siswaRoutes.GET("/tests/by-type/:type_test/kelas/:kelas_id", siswaHandler.GetTestByKelasHandler)
		siswaRoutes.GET("/tests/:id/soal", siswaHandler.GetSoalByTestIDHandler)

		siswaRoutes.POST("/jawaban/save", siswaHandler.SaveJawabanHandler)
		siswaRoutes.GET("/test/:test_id/soal", siswaHandler.GetSoalHandler)

		//test review

	}
}
