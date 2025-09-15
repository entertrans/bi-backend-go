package siswa

import (
	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

// GetKisiKisiByID mendapatkan kisi-kisi berdasarkan ID
func GetKisiKisiByID(kisiKisiID uint) (models.KisiKisi, error) {
	var kisiKisi models.KisiKisi
	err := config.DB.Preload("Mapel").Preload("Kelas").
		First(&kisiKisi, "kisikisi_id = ?", kisiKisiID).Error
	return kisiKisi, err
}

// GetAllKisiKisi mendapatkan semua kisi-kisi
func GetAllKisiKisi() ([]models.KisiKisi, error) {
	var kisiKisis []models.KisiKisi
	err := config.DB.Preload("Mapel").Preload("Kelas").
	Order("kisikisi_ub asc, kisikisi_semester asc").
		Find(&kisiKisis).Error
	return kisiKisis, err
}

// GetKisiKisiByKelas mendapatkan kisi-kisi berdasarkan kelas
func GetKisiKisiByKelas(kelasID uint) ([]models.KisiKisi, error) {
	var kisiKisis []models.KisiKisi
	err := config.DB.Preload("Mapel").Preload("Kelas").
		Where("kisikisi_kelas_id = ?", kelasID).
		Order("kisikisi_mapel asc, kisikisi_semester asc, CAST(kisikisi_ub AS UNSIGNED) asc").
		Find(&kisiKisis).Error
	return kisiKisis, err
}

// GetKisiKisiByMapel mendapatkan kisi-kisi berdasarkan mata pelajaran
func GetKisiKisiByMapel(mapelID uint) ([]models.KisiKisi, error) {
	var kisiKisis []models.KisiKisi
	err := config.DB.Preload("Mapel").Preload("Kelas").
	
		Where("kisikisi_mapel = ?", mapelID).
		Order("kisikisi_ub asc, kisikisi_semester asc").
		Find(&kisiKisis).Error
	return kisiKisis, err
}

// CreateKisiKisi membuat kisi-kisi baru
func CreateKisiKisi(kisiKisi *models.KisiKisi) error {
	return config.DB.Create(kisiKisi).Error
}

// UpdateKisiKisi mengupdate kisi-kisi
func UpdateKisiKisi(kisiKisiID uint, data map[string]interface{}) error {
	return config.DB.Model(&models.KisiKisi{}).
		Where("kisikisi_id = ?", kisiKisiID).
		Updates(data).Error
}

// DeleteKisiKisi menghapus kisi-kisi
func DeleteKisiKisi(kisiKisiID uint) error {
	return config.DB.Delete(&models.KisiKisi{}, kisiKisiID).Error
}