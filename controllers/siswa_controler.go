package controllers

import (
	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

// GET ALL LENGKAP
func FetchAllSiswa() ([]models.Siswa, error) {
	var siswa []models.Siswa
	err := config.DB.
		Preload("Orangtua").
		Preload("Agama").
		Find(&siswa).Error

	return siswa, err
}

// func FetchAllSiswa() ([]models.Siswa, error) {
// 	var siswa []models.Siswa
// 	err := config.DB.Find(&siswa).Error
// 	return siswa, err
// }

// GET BY siswa_nis
func FindSiswaByNis(nis string) ([]models.Siswa, error) {
	var siswa []models.Siswa
	err := config.DB.Where("siswa_nis = ?", nis).Find(&siswa).Error
	return siswa, err
}

// GET siswa + ortu
func GetSiswaWithOrtu(nis string) (*models.Siswa, error) {
	var siswa models.Siswa

	err := config.DB.
		Preload("Orangtua").
		Where("siswa_nis = ?", nis).
		First(&siswa).Error

	if err != nil {
		return nil, err
	}
	return &siswa, nil
}
