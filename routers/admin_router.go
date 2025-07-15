package routers

import (
	adminHandlers "github.com/entertrans/bi-backend-go/handlers/admin"
	adminhandlers "github.com/entertrans/bi-backend-go/handlers/admin"
	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(r *gin.Engine) {
	// ✅ Routes siswa
	r.GET("/siswa", adminHandlers.GetAllSiswa)

	r.GET("/siswappdb", adminHandlers.GetAllSiswaPPDB)
	r.GET("/siswaaktif", adminHandlers.GetAllSiswaAktif)
	r.GET("/siswakeluar", adminHandlers.GetAllSiswaKeluar)
	r.GET("/siswaalumni", adminHandlers.GetAllSiswaAlumni)
	r.GET("/siswa/:nis", adminHandlers.FindSiswaByNis)
	r.PUT("/updatesiswa/:nis", adminHandlers.UpdateSiswaFieldHandler)
	r.DELETE("/batalkan-siswa/:nis", adminHandlers.BatalkanSiswaHandler)
	siswa := r.Group("/siswa")
	{
		siswa.GET("/search", adminHandlers.SearchSiswa)
		siswa.PATCH("/:nis/terima", adminHandlers.TerimaSiswa)
		siswa.PATCH("/:nis/keluarkan", adminhandlers.KeluarkanSiswa)
		siswa.PATCH("/:nis/online", adminhandlers.SetKelasOnline)
		siswa.PATCH("/:nis/offline", adminhandlers.SetKelasOffline)
	}

	// ✅ Routes ortu
	r.GET("/ortu", adminHandlers.GetAllOrtu)
	r.GET("/ortu/:nis", adminHandlers.FindOrtuByNis)

	// ✅ Lookup
	lookup := r.Group("/lookup")
	{
		lookup.GET("/agama", adminHandlers.GetAllAgama)
		lookup.GET("/kelas", adminHandlers.GetAllKelas)
		lookup.GET("/satelit", adminHandlers.GetAllSatelit)
		lookup.GET("/tahun_ajaran", adminHandlers.GetAllTA)
	}

	// ✅ tagihan
	tagihan := r.Group("/tagihan")
	{
		tagihan.GET("", adminhandlers.GetAllTagihan)
		tagihan.POST("/tambah", adminhandlers.TambahTagihan)
		tagihan.PATCH("/:id/edit", adminhandlers.EditTagihan)
		tagihan.DELETE("/:id/delete", adminhandlers.DeleteTagihan)
	}

	//invoice
	invoice := r.Group("/invoice")
	{
		invoice.GET("", adminHandlers.GetAllInvoiceHandler)
		invoice.POST("", adminHandlers.CreateInvoiceHandler)
		invoice.GET("/cek", adminHandlers.CekInvoiceIDHandler)
		invoice.GET("/by-id", adminHandlers.GetInvoiceByID)
		invoice.GET("/penerima", adminHandlers.GetInvoicePenerima)
		invoice.PUT("/penerima/potongan", adminHandlers.UpdatePotonganPenerima)
		invoice.POST("/penerima/id", adminHandlers.TambahPenerimaInvoice)
		invoice.DELETE("/penerima/:id", adminHandlers.DeletePenerimaInvoice)
		invoice.GET("/penerima/:nis", adminHandlers.GetInvoicePenerimaByNIS)
		// invoice.PUT("/penerima/:id/tambahan", adminHandlers.UpdateInvoicePenerimaTambahan)
		invoice.PUT("/penerima/:id/tambahan", adminHandlers.HandleUpdateTambahanTagihan)

	}

	// ✅ Etc
	r.POST("/ppdb", adminHandlers.HandleCreatePPDB)
	r.POST("/api/upload-dokumen", adminHandlers.UploadDokumenHandler)
}
