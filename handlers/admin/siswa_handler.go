package adminhandlers

import (
	"log"
	"net/http"

	adminControllers "github.com/entertrans/bi-backend-go/controllers/admin"
	"github.com/gin-gonic/gin"
)

// search
func SearchSiswa(c *gin.Context) {
	adminControllers.SearchSiswa(c)
}

// GET /siswa
func GetAllSiswa(c *gin.Context) {
	nis := c.Query("siswa_nis")
	if nis != "" {
		data, err := adminControllers.FindSiswaByNis(nis)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
		return
	}

	data, err := adminControllers.FetchAllSiswa()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GET /siswappdb
func GetAllSiswaPPDB(c *gin.Context) {
	data, err := adminControllers.FetchAllSiswaPPDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GET /siswaaktif
func GetAllSiswaAktif(c *gin.Context) {
	data, err := adminControllers.FetchAllSiswaAktif()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GET /siswakeluar
func GetAllSiswaKeluar(c *gin.Context) {
	data, err := adminControllers.FetchAllSiswaKeluar()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GET /siswaalumni
func GetAllSiswaAlumni(c *gin.Context) {
	data, err := adminControllers.FetchAllSiswaAlumni()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GET /siswa/:nis
func FindSiswaByNis(c *gin.Context) {
	nis := c.Param("nis")

	data, err := adminControllers.FindSiswaByNis(nis)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GET /siswa/:nis/detail
func GetSiswaWithOrtu(c *gin.Context) {
	nis := c.Param("nis")

	data, err := adminControllers.GetSiswaWithOrtu(nis)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func TerimaSiswa(c *gin.Context) {
	nis := c.Param("nis")

	if err := adminControllers.TerimaSiswa(nis); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menerima siswa."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Siswa diterima."})
}

func KeluarkanSiswa(c *gin.Context) {
	nis := c.Param("nis")
	log.Printf("[INFO] PATCH /siswa/%s/keluarkan", nis)

	var req struct {
		TglKeluar *string `json:"tgl_keluar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[ERROR] Gagal binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid."})
		return
	}

	if err := adminControllers.KeluarkanSiswa(nis, req.TglKeluar); err != nil {
		log.Println("[ERROR]", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengeluarkan siswa."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Siswa dikeluarkan."})
}

func SetKelasOnline(c *gin.Context) {
	nis := c.Param("nis")

	var req struct {
		Value int `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid."})
		return
	}

	if err := adminControllers.SetKelasOnline(nis, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengubah status kelas online."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil mengubah kelas online."})
}

func SetKelasOffline(c *gin.Context) {
	nis := c.Param("nis")

	var req struct {
		Value int `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid."})
		return
	}

	if err := adminControllers.SetKelasOffline(nis, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengubah status kelas offline."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil mengubah kelas offline."})
}
