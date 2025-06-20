package handlers

import (
	"net/http"

	"github.com/entertrans/bi-backend-go/config" // ganti dengan nama module kamu di go.mod
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

func GetAllSiswa(c *gin.Context) {
	var siswa []models.Siswa

	// ambil semua data siswa dari database
	if err := config.DB.Find(&siswa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, siswa)
}
