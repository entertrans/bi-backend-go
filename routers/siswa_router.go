package routers

import (
	siswaHandler "github.com/entertrans/bi-backend-go/handlers/siswa"
	"github.com/gin-gonic/gin"
)

func RegisterSiswaRoutes(r *gin.Engine) {
	siswaRoutes := r.Group("/siswa")
	{
		siswaRoutes.POST("/test/start/:test_id", siswaHandler.StartTestHandler)
		siswaRoutes.GET("/session/:session_id/soal", siswaHandler.GetSessionSoalHandler)
		siswaRoutes.GET("/test/:test_id/active-session", siswaHandler.GetActiveTestSessionHandler)
		siswaRoutes.GET("/test/:test_id/session", siswaHandler.GetTestSessionHandler)
		siswaRoutes.POST("/test/submit/:session_id", siswaHandler.SubmitTestHandler)
		siswaRoutes.POST("/tugas/submit/:session_id", siswaHandler.SubmitTugasHandler)
		siswaRoutes.GET("/test/session/:session_id", siswaHandler.GetSessionByIDHandler)
		siswaRoutes.GET("/test/session/:session_id/nilai", siswaHandler.GetNilaiHandler)

		siswaRoutes.GET("/tests/by-type/:type_test/kelas/:kelas_id", siswaHandler.GetTestByKelasHandler)
		siswaRoutes.GET("/tests/:id/soal", siswaHandler.GetSoalByTestIDHandler)

		siswaRoutes.POST("/jawaban/save", siswaHandler.SaveJawabanHandler)
		siswaRoutes.GET("/test/:test_id/soal", siswaHandler.GetSoalHandler)

		//test review
	}
	
	kisiKisi := r.Group("siswa/kisikisi")
	{
		kisiKisi.GET("/", siswaHandler.GetAllKisiKisiHandler)
		kisiKisi.GET("/:id", siswaHandler.GetKisiKisiByIDHandler)
		kisiKisi.GET("/kelas/:kelas_id", siswaHandler.GetKisiKisiByKelasHandler)
		kisiKisi.GET("/mapel/:mapel_id", siswaHandler.GetKisiKisiByMapelHandler)
		kisiKisi.POST("", siswaHandler.CreateKisiKisiHandler)
		kisiKisi.PUT("/:id", siswaHandler.UpdateKisiKisiHandler)
		kisiKisi.DELETE("/:id", siswaHandler.DeleteKisiKisiHandler)
	}
	
	invoice := r.Group("siswa/invoice")
	{
		// List semua invoice siswa
		invoice.GET("/history/:nis", siswaHandler.HistoryKeuanganByNISHandler)

		// Detail 1 invoice siswa
		invoice.GET("/detail/:nis", siswaHandler.InvoiceDetailByNISHandler)
		invoice.GET("/:nis/invoice/unpaid-latest", siswaHandler.LatestUnpaidInvoiceHandler)
	}
	
	online := r.Group("siswa/online")
	{
		online.GET("/", siswaHandler.GetAllOnlineClassHandler)
		online.GET("/:id", siswaHandler.GetOnlineClassByIDHandler)
		online.GET("/kelas/:kelas_id", siswaHandler.GetOnlineClassByKelasHandler)
		online.GET("/mapel/:mapel_id", siswaHandler.GetOnlineClassByMapelHandler)
		online.POST("/", siswaHandler.CreateOnlineClassHandler)
		online.PUT("/:id", siswaHandler.UpdateOnlineClassHandler)
		online.DELETE("/:id", siswaHandler.DeleteOnlineClassHandler)
	}
}