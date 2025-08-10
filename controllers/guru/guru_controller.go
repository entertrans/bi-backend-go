package gurucontrollers

import (
	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

// Ambil semua guru
func GetAllGuru() ([]models.Guru, error) {
	var gurus []models.Guru
	err := config.DB.Find(&gurus).Error
	return gurus, err
}

// Ambil guru by ID
func GetGuruByID(guruID uint) (models.Guru, error) {
	var guru models.Guru
	err := config.DB.First(&guru, "guru_id = ?", guruID).Error
	return guru, err
}

// Create guru
func CreateGuru(guru *models.Guru) error {
	return config.DB.Create(guru).Error
}

// Update guru
func UpdateGuru(guruID uint, data map[string]interface{}) error {
	return config.DB.Model(&models.Guru{}).
		Where("guru_id = ?", guruID).
		Updates(data).Error
}

// Delete guru (soft delete bisa diimplementasi jika perlu)
func DeleteGuru(guruID uint) error {
	return config.DB.Delete(&models.Guru{}, guruID).Error
}
