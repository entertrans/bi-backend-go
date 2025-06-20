package controllers

import (
	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

// GET ALL
func FetchAllOrtu() ([]models.Orangtua, error) {
	var ortu []models.Orangtua
	err := config.DB.Find(&ortu).Error
	return ortu, err
}

// GET BY siswa_nis
func FindOrtuByNis(nis string) ([]models.Orangtua, error) {
	var ortu []models.Orangtua
	err := config.DB.Where("siswa_nis = ?", nis).Find(&ortu).Error
	return ortu, err
}
