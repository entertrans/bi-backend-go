package handlers

import (
	"net/http"

	"github.com/entertrans/bi-backend-go/controllers"
	"github.com/gin-gonic/gin"
)

// GET /ortu
func GetAllOrtu(c *gin.Context) {
	nis := c.Query("siswa_nis")
	if nis != "" {
		data, err := controllers.FindOrtuByNis(nis)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
		return
	}

	data, err := controllers.FetchAllOrtu()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GET /ortu/:nis
func FindOrtuByNis(c *gin.Context) {
	nis := c.Param("nis")

	data, err := controllers.FindOrtuByNis(nis)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, data)
}
