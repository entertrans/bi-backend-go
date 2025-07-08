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
	r.PUT("/updatesiswa/:nis", adminhandlers.UpdateSiswaFieldHandler)
	r.DELETE("/batalkan-siswa/:nis", adminhandlers.BatalkanSiswaHandler)

	// r.PATCH("/nis/:nis", adminHandlers.UpdateSiswaByNIS)

	// ✅ Routes ortu
	r.GET("/ortu", adminHandlers.GetAllOrtu)
	r.GET("/ortu/:nis", adminHandlers.FindOrtuByNis)

	// ✅ Lookup
	r.GET("/lookup/agama", adminHandlers.GetAllAgama)
	r.GET("/lookup/kelas", adminHandlers.GetAllKelas)
	r.GET("/lookup/satelit", adminHandlers.GetAllSatelit)
	r.GET("/lookup/tahun_ajaran", adminHandlers.GetAllTA)

	// ✅ Etc
	r.POST("/ppdb", adminHandlers.HandleCreatePPDB)
}
