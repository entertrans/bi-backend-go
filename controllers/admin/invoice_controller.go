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
type HistoryInvoice struct {
	InvoiceID         string                           `json:"invoice_id"`
	InvoiceDeskripsi  string                           `json:"invoice_deskripsi"`
	InvoiceTgl        string                           `json:"invoice_tgl"`
	InvoiceJatuhTempo string                           `json:"invoice_jatuh_tempo"`
	TotalBayar        int                              `json:"totalBayar"`
	Potongan          int                              `json:"potongan"`
	Tagihan           []models.InvoiceTagihan          `json:"tagihan"`
	TambahanTagihan   []models.InvoicePenerimaTambahan `json:"tambahan_tagihan"`
}

type HistorySiswaKeuangan struct {
	Siswa   models.Siswa     `json:"siswa"`
	History []HistoryInvoice `json:"history"`
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
		Preload("Penerima.Pembayaran").
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
		Preload("Pembayaran").
		// Preload("Penerima.Pembayaran").
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

func UpdatePotonganPenerima(c *gin.Context) {
	var req struct {
		IDPenerima uint `json:"id_penerima"`
		Potongan   int  `json:"potongan"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data salah"})
		return
	}

	if err := config.DB.Model(&models.InvoicePenerima{}).
		Where("id = ?", req.IDPenerima).
		Update("potongan", req.Potongan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update potongan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Potongan diperbarui"})
}

func DeletePenerimaInvoice(c *gin.Context) {
	id := c.Param("id")

	if err := config.DB.Delete(&models.InvoicePenerima{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus penerima"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Penerima dihapus"})
}

func FindInvoicePenerimaByNIS(nis string) (*models.InvoicePenerima, error) {
	var penerima models.InvoicePenerima

	err := config.DB.
		Preload("Siswa.Kelas").
		Preload("Tambahan").
		Preload("Pembayaran").
		Preload("Invoice.Tagihan"). // ⚠️ ini yang sebelumnya error
		Where("nis = ?", nis).
		First(&penerima).Error

	if err != nil {
		return nil, errors.New("Penerima tidak ditemukan")
	}

	return &penerima, nil
}

func GetHistoryKeuanganByNIS(nis string) (HistorySiswaKeuangan, error) {
	var penerimas []models.InvoicePenerima
	err := config.DB.Preload("Invoice.Tagihan").
		Preload("Tambahan").
		Preload("Pembayaran").
		Preload("Siswa").
		Preload("Siswa.Kelas").
		Where("nis = ?", nis).
		Find(&penerimas).Error
	if err != nil {
		return HistorySiswaKeuangan{}, err
	}

	var history []HistoryInvoice
	var siswa models.Siswa

	for _, penerima := range penerimas {
		if siswa.SiswaNama == nil || *siswa.SiswaNama == "" {
			siswa = penerima.Siswa
		}

		totalTagihan := 0
		for _, tagihan := range penerima.Invoice.Tagihan {
			totalTagihan += tagihan.Nominal
		}
		for _, tambahan := range penerima.Tambahan {
			totalTagihan += tambahan.Nominal
		}
		totalTagihan -= penerima.Potongan

		totalBayar := 0
		for _, bayar := range penerima.Pembayaran {
			totalBayar += bayar.Nominal
		}

		history = append(history, HistoryInvoice{
			InvoiceID:         penerima.Invoice.IDInvoice,
			InvoiceDeskripsi:  penerima.Invoice.Deskripsi,
			InvoiceTgl:        penerima.Invoice.TglInvoice,    // tanpa .Format
			InvoiceJatuhTempo: penerima.Invoice.TglJatuhTempo, // tanpa .Format
			Potongan:          penerima.Potongan,
			Tagihan:           penerima.Invoice.Tagihan,
			TambahanTagihan:   penerima.Tambahan,
			TotalBayar:        totalBayar,
		})

	}

	return HistorySiswaKeuangan{
		Siswa:   siswa,
		History: history,
	}, nil
}

type UpdateTambahanTagihanRequest struct {
	TambahanTagihan []models.InvoicePenerimaTambahan `json:"tambahan_tagihan"`
}

func UpdateTambahanTagihanByPenerimaID(id uint, req UpdateTambahanTagihanRequest) error {
	// 1. Pastikan penerima invoice-nya ada
	var penerima models.InvoicePenerima
	if err := config.DB.First(&penerima, id).Error; err != nil {
		return errors.New("Penerima tidak ditemukan")
	}

	// 2. Hapus semua tagihan tambahan sebelumnya
	if err := config.DB.Where("id_penerima = ?", id).Delete(&models.InvoicePenerimaTambahan{}).Error; err != nil {
		return errors.New("Gagal menghapus tagihan tambahan lama")
	}

	// 3. Tambahkan tagihan tambahan baru (jika ada)
	for _, item := range req.TambahanTagihan {
		newItem := models.InvoicePenerimaTambahan{
			IDPenerima:  id,
			NamaTagihan: item.NamaTagihan,
			Nominal:     item.Nominal,
		}
		if err := config.DB.Create(&newItem).Error; err != nil {
			return errors.New("Gagal menyimpan tagihan tambahan baru")
		}
	}

	return nil
}
