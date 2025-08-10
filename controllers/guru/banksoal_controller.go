package gurucontrollers

import (
	"errors"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

func GetBankSoalByStatus(isDeleted bool) ([]models.TO_BankSoal, error) {
	var soals []models.TO_BankSoal
	if isDeleted {
		err := config.DB.Unscoped().Where("deleted_at IS NOT NULL").Find(&soals).Error
		return soals, err
	} else {
		err := config.DB.Where("deleted_at IS NULL").Find(&soals).Error
		return soals, err
	}
}

func RestoreBankSoal(soalID uint) error {
	return config.DB.Model(&models.TO_BankSoal{}).
		Unscoped().
		Where("soal_id = ?", soalID).
		Update("deleted_at", nil).Error
}

func CreateBankSoal(soal *models.TO_BankSoal) error {
	return config.DB.Create(soal).Error
}

func GetBankSoalByGuru(guruID uint) ([]models.TO_BankSoal, error) {
	var soals []models.TO_BankSoal
	err := config.DB.Where("guru_id = ?", guruID).Find(&soals).Error
	return soals, err
}

func UpdateBankSoal(soalID uint, data map[string]interface{}) error {
	result := config.DB.Model(&models.TO_BankSoal{}).Where("soal_id = ?", soalID).Updates(data)
	if result.RowsAffected == 0 {
		return errors.New("bank soal tidak ditemukan")
	}
	return result.Error
}

func DeleteBankSoal(soalID uint) error {
	result := config.DB.Delete(&models.TO_BankSoal{}, soalID)
	if result.RowsAffected == 0 {
		return errors.New("bank soal tidak ditemukan")
	}
	return result.Error
}
