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

	// Ambil user dulu
	if err := config.DB.Where("pengguna_username = ?", username).First(&pengguna).Error; err != nil {
		return nil, err
	}

	// bikin db query base
	db := config.DB.Model(&pengguna)

	// preload sesuai level
	switch pengguna.PenggunaLevel {
	case "1": // ADMIN
		db = db.Preload("Admin")
	case "2": // SISWA
		db = db.Preload("Siswa").
			Preload("Siswa.Kelas")
	case "3": // GURU
		db = db.Preload("Guru")
	}

	// reload user + relasinya
	if err := db.First(&pengguna, pengguna.PenggunaID).Error; err != nil {
		return nil, err
	}

	return &pengguna, nil
}
