package admincontrollers

import (
	"errors"
	"net/http"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//	func GetAllPettyCashPeriode() ([]models.PettyCashPeriode, error) {
//		var periodes []models.PettyCashPeriode
//		err := config.DB.Find(&periodes).Error
//		return periodes, err
//	}
func GetPettyCashByLokasi(lokasi string) ([]models.PettyCashPeriode, error) {
	var periode []models.PettyCashPeriode

	err := config.DB.
		Where("lokasi = ?", lokasi).
		Preload("Transaksis").
		Order("id desc").
		Find(&periode).Error

	return periode, err
}

func GetAllPettyCashWithTransaksi() ([]models.PettyCashPeriode, error) {
	var periode []models.PettyCashPeriode

	err := config.DB.
		Preload("Transaksis").
		Order("id desc").
		Find(&periode).Error

	if err != nil {
		return nil, err
	}

	return periode, nil
}

func CreatePettyCashPeriode(p models.PettyCashPeriode) error {
	return config.DB.Create(&p).Error
}

func GetPettyCashPeriodeByID(id int) (models.PettyCashPeriode, error) {
	var periode models.PettyCashPeriode
	err := config.DB.First(&periode, id).Error
	return periode, err
}

func UpdatePettyCashPeriode(p models.PettyCashPeriode) error {
	return config.DB.Save(&p).Error
}

func DeletePettyCashPeriode(id int) error {
	return config.DB.Delete(&models.PettyCashPeriode{}, id).Error
}
func DeleteTransaksiByID(c *gin.Context) {
	id := c.Param("id")
	var transaksi models.Transaksi

	if err := config.DB.First(&transaksi, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&transaksi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaksi berhasil dihapus"})
}

func GetTransaksiWithPeriode(periodeID int) (models.PettyCashPeriode, []models.Transaksi, int64, error) {
	var periode models.PettyCashPeriode
	if err := config.DB.First(&periode, periodeID).Error; err != nil {
		return models.PettyCashPeriode{}, nil, 0, err
	}

	var transaksis []models.Transaksi
	if err := config.DB.
		Where("id_periode = ?", periodeID).
		Order("tanggal asc").
		Find(&transaksis).Error; err != nil {
		return periode, nil, 0, err
	}

	// Hitung saldo: mulai dari SaldoAwal periode
	saldo := periode.SaldoAwal
	for _, trx := range transaksis {
		if trx.Jenis == "masuk" {
			saldo += trx.Nominal
		} else if trx.Jenis == "keluar" {
			saldo -= trx.Nominal
		}
	}

	return periode, transaksis, saldo, nil
}

func AddTransaksi(input models.Transaksi) (models.Transaksi, error) {
	var lastTransaksi models.Transaksi
	err := config.DB.Where("id_periode = ?", input.IDPeriode).
		Order("tanggal desc, id desc").
		First(&lastTransaksi).Error

	// Kalau error bukan karena record not found, balikin errornya
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Transaksi{}, err
	}

	lastSaldo := lastTransaksi.SaldoSetelah

	switch input.Jenis {
	case "masuk":
		input.SaldoSetelah = lastSaldo + input.Nominal
	case "keluar":
		input.SaldoSetelah = lastSaldo - input.Nominal
	default:
		return models.Transaksi{}, errors.New("jenis transaksi tidak valid")
	}

	if err := config.DB.Create(&input).Error; err != nil {
		return models.Transaksi{}, err
	}

	return input, nil
}
