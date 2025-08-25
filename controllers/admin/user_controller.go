package admincontrollers

import (
	"net/http"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

func GetAllMapel() ([]models.Mapel, error) {
    var mapels []models.Mapel
    err := config.DB.
        Order("nm_mapel ASC").
        Find(&mapels).Error
    return mapels, err
}

func GetAllMapelWithGuruMapels() ([]models.Mapel, error) {
    var mapels []models.Mapel
    err := config.DB.
        Preload("GuruMapels").
        Order("nm_mapel ASC").
        Find(&mapels).Error
    return mapels, err
}

func GetMapelByKelas(c *gin.Context) {
	kelasID := c.Param("id") // ambil dari /mapel-by-kelas/:id
	if kelasID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "kelas_id wajib diisi"})
		return
	}

	var kelasMapels []models.KelasMapel
	err := config.DB.
		Preload("Mapel").
		Where("kelas_id = ?", kelasID).
		Find(&kelasMapels).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data mapel"})
		return
	}

	c.JSON(http.StatusOK, kelasMapels)
}

func GetSiswaByKelas(kelasID uint) ([]models.Siswa, error) {
	var siswa []models.Siswa
	err := config.DB.
		Preload("Kelas").
		// Preload("Agama").
		// Preload("Satelit").
		Where("siswa_kelas_id = ? AND (soft_deleted IS NULL OR soft_deleted = 0)", kelasID).
		Order("siswa_nama ASC").
		Find(&siswa).Error

	return siswa, err
}

