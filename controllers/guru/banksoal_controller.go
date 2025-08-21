package gurucontrollers

import (
	"errors"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

// ========================
// Ambil soal aktif (belum dihapus)
// ========================
func GetActiveBankSoal() ([]models.TO_BankSoal, error) {
	var soals []models.TO_BankSoal
	err := config.DB.
		Preload("Guru").
		Preload("Kelas").
		Preload("Lampiran"). // 🔥 ikut load lampiran kalau ada
		Where("deleted_at IS NULL").
		Order("created_at desc").
		Find(&soals).Error
	return soals, err
}

// ========================
// Ambil soal non-aktif (soft deleted)
// ========================
func GetInactiveBankSoal() ([]models.TO_BankSoal, error) {
	var soals []models.TO_BankSoal
	err := config.DB.
		Unscoped().
		Preload("Guru").
		Preload("Kelas").
		Preload("Lampiran").
		Where("deleted_at IS NOT NULL").
		Order("created_at desc").
		Find(&soals).Error
	return soals, err
}

// ========================
// Restore soal (hapus deleted_at)
// ========================
func RestoreBankSoal(soalID uint) error {
	result := config.DB.Model(&models.TO_BankSoal{}).
		Unscoped().
		Where("soal_id = ?", soalID).
		Update("deleted_at", nil)
	if result.RowsAffected == 0 {
		return errors.New("bank soal tidak ditemukan")
	}
	return result.Error
}

// ========================
// Ambil soal milik guru tertentu
// ========================
func GetBankSoalByGuru(guruID uint) ([]models.TO_BankSoal, error) {
	var soals []models.TO_BankSoal
	err := config.DB.
		Preload("Guru").
		Preload("Kelas").
		Preload("Lampiran").
		Where("guru_id = ? AND deleted_at IS NULL", guruID). // 🔧 ganti is_deleted jadi deleted_at IS NULL
		Order("created_at desc").
		Find(&soals).Error
	return soals, err
}

// ========================
// Soft delete soal
// ========================
func DeleteBankSoal(soalID uint) error {
	result := config.DB.Delete(&models.TO_BankSoal{}, soalID)
	if result.RowsAffected == 0 {
		return errors.New("bank soal tidak ditemukan")
	}
	return result.Error
}
