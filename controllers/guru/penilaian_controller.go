package gurucontrollers

import (
	"errors"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

func CreatePenilaian(penilaian *models.TO_PenilaianGuru) error {
	return config.DB.Create(penilaian).Error
}

func UpdatePenilaian(penilaianID uint, data map[string]interface{}) error {
	result := config.DB.Model(&models.TO_PenilaianGuru{}).Where("penilaian_id = ?", penilaianID).Updates(data)
	if result.RowsAffected == 0 {
		return errors.New("penilaian tidak ditemukan")
	}
	return result.Error
}

func GetPenilaianByFinalID(finalID uint) (*models.TO_PenilaianGuru, error) {
	var penilaian models.TO_PenilaianGuru
	err := config.DB.Where("final_id = ?", finalID).First(&penilaian).Error
	if err != nil {
		return nil, err
	}
	return &penilaian, nil
}
