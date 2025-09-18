package siswa

import (
	"net/http"

	siswaControllers "github.com/entertrans/bi-backend-go/controllers/siswa"
	"github.com/gin-gonic/gin"
)

func HistoryKeuanganByNISHandler(c *gin.Context) {
	nis := c.Param("nis")
	result, err := siswaControllers.GetHistoryKeuanganByNIS(nis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func InvoiceDetailByNISHandler(c *gin.Context) {
	nis := c.Param("nis")
	idInvoice := c.Query("idInvoice") // ambil dari query param

	if idInvoice == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "idInvoice wajib diisi"})
		return
	}

	result, err := siswaControllers.GetInvoiceDetailByNIS(nis, idInvoice)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
func LatestUnpaidInvoiceHandler(c *gin.Context) {
	nis := c.Param("nis")

	result, err := siswaControllers.GetLatestUnpaidInvoiceByNIS(nis)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada tagihan aktif"})
		return
	}

	// Kalau sudah lunas semua, bisa juga return kosong
	if result.SisaTagihan <= 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Semua tagihan sudah lunas"})
		return
	}

	c.JSON(http.StatusOK, result)
}
