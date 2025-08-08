package adminhandlers

import (
	"log"
	"net/http"
	"strconv"

	adminControllers "github.com/entertrans/bi-backend-go/controllers/admin"
	"github.com/gin-gonic/gin"
)

func HandleUpdateTambahanTagihan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req adminControllers.UpdateTambahanTagihanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("DEBUG ERROR BIND JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload tidak valid"})
		return
	}

	if err := adminControllers.UpdateTambahanTagihanByPenerimaID(uint(id), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tagihan tambahan berhasil diperbarui"})
}

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

func UpdatePotonganPenerima(c *gin.Context) {
	adminControllers.UpdatePotonganPenerima(c)
}
func DeletePenerimaInvoice(c *gin.Context) {
	adminControllers.DeletePenerimaInvoice(c)
}

func GetInvoicePenerimaByNIS(c *gin.Context) {
	nis := c.Param("nis")

	penerima, err := adminControllers.FindInvoicePenerimaByNIS(nis)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, penerima)
}

func HistoryKeuanganByNISHandler(c *gin.Context) {
	nis := c.Param("nis")
	result, err := adminControllers.GetHistoryKeuanganByNIS(nis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
