package adminhandlers

import (
	"net/http"

	admincontrollers "github.com/entertrans/bi-backend-go/controllers/admin"
	"github.com/gin-gonic/gin"
)

func UpdateSiswaFieldHandler(c *gin.Context) {
	nis := c.Param("nis")

	var req struct {
		Field string      `json:"field"`
		Value interface{} `json:"value"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permintaan tidak valid"})
		return
	}

	if err := admincontrollers.UpdateSiswaField(nis, req.Field, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil diperbarui"})
}

func BatalkanSiswaHandler(c *gin.Context) {
	nis := c.Param("nis")
	err := admincontrollers.BatalkanSiswaByNIS(nis)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membatalkan siswa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Siswa berhasil dibatalkan dan datanya telah dihapus"})
}
