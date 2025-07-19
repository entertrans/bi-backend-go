package adminhandlers

import (
	"net/http"

	adminControllers "github.com/entertrans/bi-backend-go/controllers/admin"
	"github.com/gin-gonic/gin"
)

func GetAllKwitansi(c *gin.Context) {
	data, err := adminControllers.GetKwitansiList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kwitansi"})
		return
	}

	c.JSON(http.StatusOK, data)
}

func GetInvoiceKwitansiByID(c *gin.Context) {
	adminControllers.GetInvoiceKwitansiByID(c)
}
