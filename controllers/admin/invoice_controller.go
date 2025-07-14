package admincontrollers

import (
	"errors"
	"net/http"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

type PenerimaInput struct {
	Penerima []struct {
		NIS string `json:"nis"`
	} `json:"penerima"`
}

func GetAllInvoice(c *gin.Context) {
	var invoices []models.Invoice

	err := config.DB.Preload("Tagihan").Find(&invoices).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, invoices)
}

func CreateInvoice(c *gin.Context) {
	var input models.Invoice

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// simpan invoice
	err := config.DB.Create(&input).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice berhasil disimpan"})
}

func CekInvoiceID(c *gin.Context) {
	id := c.Query("id") // ambil dari query param
	var count int64

	err := config.DB.Model(&models.Invoice{}).Where("id_invoice = ?", id).Count(&count).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal cek ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": count > 0})
}
func GetInvoiceByID(c *gin.Context) {
	id := c.Query("id")

	var invoice models.Invoice
	if err := config.DB.
		Preload("Tagihan").
		Preload("Penerima.Tambahan").
		Where("id_invoice = ?", id).
		First(&invoice).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

func GetInvoicePenerima(c *gin.Context) {
	idInvoice := c.Query("id")

	var penerima []models.InvoicePenerima
	err := config.DB.
		Preload("Siswa.Kelas"). // preload juga relasi ke kelas
		Preload("Tambahan").
		Where("id_invoice = ?", idInvoice).
		Find(&penerima).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data penerima"})
		return
	}

	c.JSON(http.StatusOK, penerima)
}

func TambahInvoicePenerima(c *gin.Context) {
	idInvoice := c.Param("id_invoice")

	var input models.InvoicePenerima
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	input.IDInvoice = idInvoice

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan penerima"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Penerima berhasil ditambahkan"})
}

func TambahPenerimaKeInvoice(idInvoice string, input PenerimaInput) error {
	if idInvoice == "" {
		return errors.New("ID invoice kosong")
	}

	for _, siswa := range input.Penerima {
		// Cek duplikat
		var existing models.InvoicePenerima
		err := config.DB.
			Where("id_invoice = ? AND nis = ?", idInvoice, siswa.NIS).
			First(&existing).Error
		if err == nil {
			continue // sudah ada, skip
		}

		penerima := models.InvoicePenerima{
			IDInvoice: idInvoice,
			NIS:       siswa.NIS,
			Potongan:  0,
		}

		if err := config.DB.Create(&penerima).Error; err != nil {
			return errors.New("gagal menyimpan penerima")
		}
	}

	// Update status di invoice
	if err := config.DB.Model(&models.Invoice{}).
		Where("id_invoice = ?", idInvoice).
		Update("sudah_ada_penerima", true).Error; err != nil {
		return errors.New("gagal update status invoice")
	}

	return nil
}
