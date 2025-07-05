package controllers

import (
	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

// READ
func FetchAllSiswa() ([]models.Siswa, error) {
	var siswa []models.Siswa
	err := config.DB.
		Preload("Orangtua").
		Preload("Lampiran").
		Preload("Agama").
		Find(&siswa).Error

	return siswa, err
}

// GET ALL LENGKAP
func FetchAllSiswaAktif() ([]models.Siswa, error) {
	var siswa []models.Siswa
	err := config.DB.
		Where("soft_deleted = ? AND siswa_kelas_id < ?", 0, 16).
		Preload("Orangtua").
		Preload("Lampiran").
		Preload("Kelas").
		Preload("Satelit").
		Preload("Agama").
		Find(&siswa).Error

	return siswa, err
}

func FetchAllSiswaKeluar() ([]models.Siswa, error) {
	var siswa []models.Siswa
	err := config.DB.
		Where("soft_deleted = ?", 1).
		Preload("Orangtua").
		Preload("Kelas").
		Preload("Satelit").
		Preload("Agama").
		Find(&siswa).Error

	return siswa, err
}

func FetchAllSiswaPPDB() ([]models.Siswa, error) {
	var siswa []models.Siswa
	err := config.DB.
		Where("soft_deleted = ?", 2).
		Preload("Satelit").
		Find(&siswa).Error

	return siswa, err
}

func FetchAllSiswaAlumni() ([]models.Siswa, error) {
	var siswa []models.Siswa
	err := config.DB.
		Where("siswa_kelas_id > ?", 15).
		Preload("Orangtua").
		Preload("Kelas").
		Preload("Satelit").
		Preload("Agama").
		Find(&siswa).Error

	return siswa, err
}

// GET BY siswa_nis
func FindSiswaByNis(nis string) (models.Siswa, error) {
	var siswa models.Siswa
	err := config.DB.
		Where("siswa_nis = ?", nis).
		Preload("Orangtua").
		Preload("Lampiran").
		Preload("Kelas").
		Preload("Satelit").
		Preload("Agama").
		First(&siswa).Error
	return siswa, err
}

// GET siswa + ortu
func GetSiswaWithOrtu(nis string) (*models.Siswa, error) {
	var siswa models.Siswa

	err := config.DB.
	Where("siswa_nis = ?", nis).
		Preload("Orangtua").
		Preload("Kelas").
		Preload("Satelit").
		Preload("Agama").

		First(&siswa).Error

	if err != nil {
		return nil, err
	}
	return &siswa, nil
}

//CREATE
