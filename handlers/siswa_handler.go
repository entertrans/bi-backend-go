package handlers

import (
	"net/http"

	// ganti dengan nama module kamu di go.mod
	"github.com/entertrans/bi-backend-go/controllers"
	"github.com/gin-gonic/gin"
)

// GET /Siswa
func GetAllSiswa(c *gin.Context) {
	nis := c.Query("siswa_nis")
	if nis != "" {
		data, err := controllers.FindSiswaByNis(nis)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
		return
	}

	data, err := controllers.FetchAllSiswa()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func GetAllSiswaAktif(c *gin.Context) {
	data, err := controllers.FetchAllSiswaAktif()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
func GetAllSiswaKeluar(c *gin.Context) {
	data, err := controllers.FetchAllSiswaKeluar()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
func GetAllSiswaAlumni(c *gin.Context) {
	data, err := controllers.FetchAllSiswaAlumni()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
func GetAllSiswaPPDB(c *gin.Context) {
	data, err := controllers.FetchAllSiswaPPDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GET /siswa/:nis
func FindSiswaByNis(c *gin.Context) {
	nis := c.Param("nis")

	data, err := controllers.FindSiswaByNis(nis)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GET /siswa/:nis/detail
func GetSiswaWithOrtu(c *gin.Context) {
	nis := c.Param("nis")

	data, err := controllers.GetSiswaWithOrtu(nis)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, data)
}
