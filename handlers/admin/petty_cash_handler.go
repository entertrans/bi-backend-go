package adminhandlers

import (
	"net/http"
	"strconv"

	adminControllers "github.com/entertrans/bi-backend-go/controllers/admin"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

// func GetPettyCashPeriodes(c *gin.Context) {
// 	data, err := adminControllers.GetAllPettyCashPeriode()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, data)
// }

func GetPettyCashByLokasiHandler(c *gin.Context) {
	lokasi := c.Param("lokasi") // ambil dari URL

	data, err := adminControllers.GetPettyCashByLokasi(lokasi)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, data)
}

func GetPettyCashPeriodeHandler(c *gin.Context) {
	data, err := adminControllers.GetAllPettyCashWithTransaksi()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, data)
}
func CreatePettyCashPeriode(c *gin.Context) {
	var periode models.PettyCashPeriode
	if err := c.ShouldBindJSON(&periode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := adminControllers.CreatePettyCashPeriode(periode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Data berhasil disimpan"})
}

func GetPettyCashPeriodeByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := adminControllers.GetPettyCashPeriodeByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func UpdatePettyCashPeriode(c *gin.Context) {
	var periode models.PettyCashPeriode
	if err := c.ShouldBindJSON(&periode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := adminControllers.UpdatePettyCashPeriode(periode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil update"})
}

func DeletePettyCashPeriode(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := adminControllers.DeletePettyCashPeriode(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil dihapus"})
}

func GetTransaksiByPeriodeHandler(c *gin.Context) {
	periodeIDStr := c.Param("id")
	periodeID, err := strconv.Atoi(periodeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	periode, transaksis, saldo, err := adminControllers.GetTransaksiWithPeriode(periodeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"periode":    periode,
		"transaksis": transaksis,
		"saldo":      saldo,
	})
}

func AddTransaksiHandler(c *gin.Context) {
	var input models.Transaksi
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	transaksi, err := adminControllers.AddTransaksi(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaksi)
}

func DeleteTransaksiHandler(c *gin.Context) {
	adminControllers.DeleteTransaksiByID(c)
}
