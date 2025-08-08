package admincontrollers

import (
	"errors"
	"time"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"gorm.io/gorm"
)

type PembayaranInput struct {
	IDPenerima uint    `json:"id_penerima"`
	Tanggal    string  `json:"tanggal"` // format: "2006-01-02"
	Nominal    int     `json:"nominal"`
	Tujuan     *string `json:"tujuan"`     // nullable
	Keterangan *string `json:"keterangan"` // nullable
}

func SimpanPembayaran(input PembayaranInput) error {
	// Validasi dasar
	if input.IDPenerima == 0 {
		return errors.New("ID penerima tidak boleh kosong")
	}
	if input.Nominal <= 0 {
		return errors.New("Nominal pembayaran harus lebih dari 0")
	}

	// Parse tanggal
	tgl, err := time.Parse("2006-01-02", input.Tanggal)
	if err != nil {
		return errors.New("format tanggal tidak valid, gunakan YYYY-MM-DD")
	}

	// Siapkan struct pembayaran
	pembayaran := models.Pembayaran{
		IDPenerima: input.IDPenerima,
		Tanggal:    tgl,
		Nominal:    input.Nominal,
		Tujuan:     input.Tujuan,
		Keterangan: input.Keterangan,
	}

	// Simpan ke database
	if err := config.DB.Create(&pembayaran).Error; err != nil {
		return err
	}

	return nil
}

func GetPembayaranByNIS(nis string, db *gorm.DB) ([]models.Pembayaran, error) {
	var list []models.Pembayaran
	if err := db.Where("nis = ?", nis).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
func DeletePembayaranByID(id string) error {
	return config.DB.Where("id = ?", id).Delete(&models.Pembayaran{}).Error
}
