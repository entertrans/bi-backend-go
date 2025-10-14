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
		guruGroup.PUT("/test/:test_id", guruhandlers.UpdateTestAktifHandler)
		guruGroup.DELETE("/test/:test_id", guruhandlers.DeleteTestHandler)
		// guruGroup.GET("/by-kelas-mapel", guruhandlers.GetBankSoalByKelasMapelHandler)
		// guruGroup.GET("/by-kelas-mapel/count", guruhandlers.GetBankSoalByKelasMapelCountHandler)
		guruGroup.DELETE("/:testId/soal/:soalId", guruhandlers.RemoveSoalFromTestHandler)
		guruGroup.GET("/by-kelas-mapel", guruhandlers.GetBankSoalByKelasMapelHandler)
		guruGroup.POST("/tests/:testId/add-soal", guruhandlers.AddSoalToTestHandler)
		// Routes penilaian
		// guruGroup.GET("/penilaian/:final_id", guruhandlers.GetPenilaianHandler)
		// guruGroup.POST("/penilaian", guruhandlers.CreatePenilaianHandler)
		// guruGroup.PATCH("/penilaian/:penilaian_id", guruhandlers.UpdatePenilaianHandler)
		// guruGroup.DELETE("/penilaian/:penilaian_id", guruhandlers.DeletePenilaianHandler)

		// routes

	}
	//test review
	testReview := r.Group("/testreview")
	{
		// test
		testReview.POST("/peserta", guruhandlers.AddPesertaHandler)                    // tambah 1 peserta
		testReview.GET("/peserta/test/:test_id", guruhandlers.GetPesertaByTestHandler) // ambil semua peserta by test
		testReview.PUT("/peserta/:peserta_id", guruhandlers.UpdatePesertaHandler)      // update peserta
		testReview.DELETE("/peserta/:peserta_id", guruhandlers.DeletePesertaHandler)   // hapus peserta
	}
	testSoalGroup := r.Group("/test-soal")
	{
		testSoalGroup.GET("/by-test/:test_id", guruhandlers.GetTestSoalByTestIdHandler)
		testSoalGroup.GET("/detail/:soal_id", guruhandlers.GetDetailTestSoalHandler)
		testSoalGroup.POST("/create", guruhandlers.CreateTestSoalHandler)
		testSoalGroup.DELETE("/delete/:soal_id", guruhandlers.DeleteTestSoalHandler)
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

	SiswaJawab := r.Group("/guru")
	{
		SiswaJawab.GET("/jawaban/siswa/:siswa_nis", guruhandlers.GetJawabanBySiswaHandler)
		SiswaJawab.GET("/jawaban/session/:session_id", guruhandlers.GetDetailJawabanHandler)
		SiswaJawab.GET("/test/:test_id/siswa", guruhandlers.GetSiswaByTestHandler)
		SiswaJawab.GET("/test/:test_id/statistics", guruhandlers.GetTestStatisticsHandler)
		SiswaJawab.GET("/session/:session_id/:jenis/jawaban", guruhandlers.GetJawabanBySession)
		SiswaJawab.PUT("/session/:session_id/jawaban", guruhandlers.UpdateJawabanFinal)
		SiswaJawab.PUT("/session/:session_id/nilai-akhir", guruhandlers.UpdateOverrideNilai)
		// routes/siswaJawabRoutes.go
		SiswaJawab.DELETE("/test/reset/:session_id", guruhandlers.ResetTestSessionHandler)
		// SiswaJawab.GET("/session/:test_id/jawaban", guruhandlers.GetSoalPenilaianHandler)
		// SiswaJawab.GET("/session/:test_id/jawaban", guruhandlers.GetSoalPenilaianHandler)
	}
	SiswaJawabRb := r.Group("/guru/jawaban")
	{
		SiswaJawabRb.GET("/rollback/siswa/:siswa_nis", guruhandlers.GetJawabanBySiswaHandler)
		SiswaJawab.GET("/session/rollback/:session_id/:jenis/jawaban", guruhandlers.GetJawabanBySession)
	}

	nilaiGroup := r.Group("/guru/nilai")
	{
		nilaiGroup.GET("/rekap", guruhandlers.GetRekapNilai)
		nilaiGroup.GET("/ub/:kelas_id/:mapel_id", guruhandlers.GetDetailUB)
		nilaiGroup.GET("/tr/:kelas_id/:mapel_id", guruhandlers.GetDetailTR)
		nilaiGroup.GET("/tugas/:kelas_id/:mapel_id", guruhandlers.GetDetailTugas)
		nilaiGroup.GET("/peserta/:type/:test_id/:kelas_id", guruhandlers.GetDetailPesertaTest)
	}
}
