package gurucontrollers

import (
	"errors"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

func GetActiveBankSoal() ([]models.TO_BankSoal, error) {
	var soals []models.TO_BankSoal
	err := config.DB.
		Preload("Guru").
		Preload("Kelas").
		Where("deleted_at IS NULL").
		Order("created_at desc").
		Find(&soals).Error
	return soals, err
}
func GetInactiveBankSoal() ([]models.TO_BankSoal, error) {
	var soal []models.TO_BankSoal
	err := config.DB.
		Unscoped().
		Preload("Guru").
		Where("deleted_at IS NOT NULL").
		Find(&soal).Error
	return soal, err
}

func RestoreBankSoal(soalID uint) error {
	// Menghapus nilai DeletedAt supaya soal jadi aktif lagi (soft undelete)
	result := config.DB.Model(&models.TO_BankSoal{}).
		Unscoped().
		Where("soal_id = ?", soalID).
		Update("deleted_at", nil)
	if result.RowsAffected == 0 {
		return errors.New("bank soal tidak ditemukan")
	}
	return result.Error
}

func CreateBankSoal(soal *models.TO_BankSoal) error {
	return config.DB.Create(soal).Error
}

func GetBankSoalByGuru(guruID uint) ([]models.TO_BankSoal, error) {
	var soals []models.TO_BankSoal
	err := config.DB.
		Preload("Guru").
		Where("guru_id = ? AND is_deleted = false", guruID).
		Order("created_at desc").
		Find(&soals).Error
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
