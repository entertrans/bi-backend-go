package admincontrollers

import (
	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

// Ambil semua data orang tua
func FetchAllOrtu() ([]models.Orangtua, error) {
	var ortu []models.Orangtua
	err := config.DB.Find(&ortu).Error
	return ortu, err
}

// Ambil data orang tua berdasarkan NIS siswa
func FindOrtuByNis(nis string) ([]models.Orangtua, error) {
	var ortu []models.Orangtua
	err := config.DB.Where("siswa_nis = ?", nis).Find(&ortu).Error
	return ortu, err
}
