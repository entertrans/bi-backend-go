package adminhandlers

import (
	"net/http"

	adminControllers "github.com/entertrans/bi-backend-go/controllers/admin"
	"github.com/gin-gonic/gin"
)

func CreateInvoiceHandler(c *gin.Context) {
	adminControllers.CreateInvoice(c)
}

func GetAllInvoiceHandler(c *gin.Context) {
	adminControllers.GetAllInvoice(c)
}

func CekInvoiceIDHandler(c *gin.Context) {
	adminControllers.CekInvoiceID(c)
}

func GetInvoiceByID(c *gin.Context) {
	adminControllers.GetInvoiceByID(c)
}

func GetInvoicePenerima(c *gin.Context) {
	adminControllers.GetInvoicePenerima(c)
}

func TambahPenerimaInvoice(c *gin.Context) {
	idInvoice := c.Query("id")
	if idInvoice == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invoice tidak ditemukan"})
		return
	}

	var input adminControllers.PenerimaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON salah"})
		return
	}

	if err := adminControllers.TambahPenerimaKeInvoice(idInvoice, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Penerima berhasil ditambahkan"})
}
