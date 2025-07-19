package adminhandlers

import (
	"net/http"
	"strconv"

	"github.com/entertrans/bi-backend-go/config"
	adminControllers "github.com/entertrans/bi-backend-go/controllers/admin"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

func BuatPembayaranHandler(c *gin.Context) {
	var input adminControllers.PembayaranInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	if err := adminControllers.SimpanPembayaran(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pembayaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran berhasil disimpan"})
}
func GetPembayaranByNISHandler(c *gin.Context) {
	nis := c.Query("nis")
	if nis == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIS tidak boleh kosong"})
		return
	}

	data, err := adminControllers.GetPembayaranByNIS(nis, config.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, data)
}
func DeletePembayaranHandler(c *gin.Context) {
	id := c.Param("id")

	err := adminControllers.DeletePembayaranByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pembayaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran berhasil dihapus"})
}
func GetPembayaranByPenerima(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var pembayaran []models.Pembayaran
	if err := config.DB.Where("penerima_id = ?", id).Order("tanggal ASC").Find(&pembayaran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pembayaran)
}
