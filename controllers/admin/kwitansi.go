package admincontrollers

import (
	"net/http"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/dto"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

// func GetKwitansiList() ([]dto.KwitansiStatus, error) {
// 	var invoices []models.Invoice
// 	err := config.DB.Preload("Penerima.Tambahan").Find(&invoices).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	var result []dto.KwitansiStatus

// 	for _, inv := range invoices {
// 		data := dto.KwitansiStatus{
// 			IDInvoice:     inv.IDInvoice,
// 			Deskripsi:     inv.Deskripsi,
// 			TglInvoice:    inv.TglInvoice,
// 			TglJatuhTempo: inv.TglJatuhTempo,
// 		}

// 		for _, penerima := range inv.Penerima {
// 			// Hitung total tagihan: total invoice + tambahan - potongan
// 			totalTagihan := 0
// 			for _, tagihan := range inv.Tagihan {
// 				totalTagihan += tagihan.Nominal
// 			}
// 			for _, tambahan := range penerima.Tambahan {
// 				totalTagihan += tambahan.Nominal
// 			}
// 			totalTagihan -= penerima.Potongan

// 			// Hitung total bayar
// 			var totalBayar int64
// 			err := config.DB.
// 				Model(&models.Pembayaran{}).
// 				Where("id_penerima = ?", penerima.ID).
// 				Select("COALESCE(SUM(nominal), 0)").Scan(&totalBayar).Error
// 			if err != nil {
// 				return nil, err
// 			}

// 			switch {
// 			case totalBayar == 0:
// 				data.Status.Belum++
// 			case int(totalBayar) >= totalTagihan:
// 				data.Status.Lunas++
// 			default:
// 				data.Status.BelumLunas++
// 			}
// 		}

// 		result = append(result, data)
// 	}

// 	return result, nil
// }

func GetKwitansiList() ([]dto.KwitansiStatus, error) {
	var invoices []models.Invoice
	err := config.DB.
		Preload("Tagihan"). // penting!
		Preload("Penerima.Tambahan").
		Find(&invoices).Error
	if err != nil {
		return nil, err
	}

	var result []dto.KwitansiStatus

	for _, inv := range invoices {
		data := dto.KwitansiStatus{
			IDInvoice:     inv.IDInvoice,
			Deskripsi:     inv.Deskripsi,
			TglInvoice:    inv.TglInvoice,
			TglJatuhTempo: inv.TglJatuhTempo,
		}

		for _, penerima := range inv.Penerima {
			// 1️⃣ Hitung total tagihan (per siswa)
			totalTagihan := 0
			for _, tagihan := range inv.Tagihan {
				totalTagihan += tagihan.Nominal
			}
			for _, tambahan := range penerima.Tambahan {
				totalTagihan += tambahan.Nominal
			}
			totalTagihan -= penerima.Potongan

			// 2️⃣ Hitung total bayar dari tbl_pembayaran
			var totalBayar int64
			err := config.DB.
				Model(&models.Pembayaran{}).
				Where("id_penerima = ?", penerima.ID).
				Select("COALESCE(SUM(nominal), 0)").
				Scan(&totalBayar).Error
			if err != nil {
				return nil, err
			}

			// 3️⃣ Tentukan status
			switch {
			case totalBayar == 0:
				data.Status.Belum++
			case int(totalBayar) >= totalTagihan:
				data.Status.Lunas++
			default:
				data.Status.BelumLunas++
			}
		}

		result = append(result, data)
	}

	return result, nil
}

func GetInvoiceKwitansiByID(c *gin.Context) {
	idInvoice := c.Query("id")

	var penerima []models.InvoicePenerima
	err := config.DB.
		Preload("Siswa.Kelas"). // preload juga relasi ke kelas
		Preload("Tambahan").
		Preload("Penerima.Pembayaran").
		Where("id_invoice = ?", idInvoice).
		Find(&penerima).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data penerima"})
		return
	}

	c.JSON(http.StatusOK, penerima)
}
