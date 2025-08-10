package guruhandlers

import (
	"net/http"
	"strconv"

	gurucontrollers "github.com/entertrans/bi-backend-go/controllers/guru"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

func GetActiveBankSoalHandler(c *gin.Context) {
	soal, err := gurucontrollers.GetBankSoalByStatus(false) // false = not deleted
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil bank soal aktif"})
		return
	}
	c.JSON(http.StatusOK, soal)
}

func GetInactiveBankSoalHandler(c *gin.Context) {
	soal, err := gurucontrollers.GetBankSoalByStatus(true) // true = deleted
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil bank soal nonaktif"})
		return
	}
	c.JSON(http.StatusOK, soal)
}

func RestoreBankSoalHandler(c *gin.Context) {
	soalIDStr := c.Param("soal_id")
	soalID, err := strconv.ParseUint(soalIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Soal ID tidak valid"})
		return
	}

	err = gurucontrollers.RestoreBankSoal(uint(soalID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal restore bank soal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bank soal berhasil direstore"})
}

func GetBankSoalHandler(c *gin.Context) {
	guruIDStr := c.Param("guru_id")
	guruID, err := strconv.ParseUint(guruIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Guru ID tidak valid"})
		return
	}

	soal, err := gurucontrollers.GetBankSoalByGuru(uint(guruID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil bank soal"})
		return
	}

	c.JSON(http.StatusOK, soal)
}

func CreateBankSoalHandler(c *gin.Context) {
	var soal models.TO_BankSoal
	if err := c.ShouldBindJSON(&soal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data soal tidak valid"})
		return
	}

	if err := gurucontrollers.CreateBankSoal(&soal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat soal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Soal berhasil dibuat", "soal_id": soal.SoalID})
}

func UpdateBankSoalHandler(c *gin.Context) {
	soalIDStr := c.Param("soal_id")
	soalID, err := strconv.ParseUint(soalIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Soal ID tidak valid"})
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data soal tidak valid"})
		return
	}

	if err := gurucontrollers.UpdateBankSoal(uint(soalID), data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update soal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Soal berhasil diupdate"})
}

func DeleteBankSoalHandler(c *gin.Context) {
	soalIDStr := c.Param("soal_id")
	soalID, err := strconv.ParseUint(soalIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Soal ID tidak valid"})
		return
	}

	if err := gurucontrollers.DeleteBankSoal(uint(soalID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus soal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Soal berhasil dihapus"})
}
