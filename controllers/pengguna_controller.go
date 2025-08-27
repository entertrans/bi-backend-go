package controllers

import (
	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

// GetActivePengguna -> ambil semua pengguna yang aktif
func GetActivePengguna() ([]models.Pengguna, error) {
	var pengguna []models.Pengguna
	err := config.DB.
		Preload("Siswa").
		Preload("Guru").
		Preload("Admin").
		Where("pengguna_status = ?", 1).
		Find(&pengguna).Error
	return pengguna, err
}

// GetPenggunaByUsername -> dengan eager loading conditional
func GetPenggunaByUsername(username string) (*models.Pengguna, error) {
	var pengguna models.Pengguna
	
	// First get basic user data
	err := config.DB.Where("pengguna_username = ?", username).First(&pengguna).Error
	if err != nil {
		return nil, err
	}
	
	// Load relations based on user level
	switch pengguna.PenggunaLevel {
	case "1": // ADMIN
		err = config.DB.Preload("Admin").First(&pengguna, pengguna.PenggunaID).Error
	case "2": // SISWA
		err = config.DB.Preload("Siswa").First(&pengguna, pengguna.PenggunaID).Error
	case "3": // GURU
		err = config.DB.Preload("Guru").First(&pengguna, pengguna.PenggunaID).Error
	default:
		// No relation to load
	}
	
	return &pengguna, err
}