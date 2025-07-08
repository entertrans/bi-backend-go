package adminhandlers

import (
	"net/http"

	adminControllers "github.com/entertrans/bi-backend-go/controllers/admin"
	"github.com/gin-gonic/gin"
)

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
