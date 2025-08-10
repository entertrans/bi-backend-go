package guruhandlers

import (
	"net/http"
	"strconv"

	gurucontrollers "github.com/entertrans/bi-backend-go/controllers/guru"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

func CreatePenilaianHandler(c *gin.Context) {
	var req struct {
		FinalID  uint    `json:"final_id" binding:"required"`
		Skor     float64 `json:"skor" binding:"required"`
		Komentar string  `json:"komentar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data penilaian tidak valid"})
		return
	}

	err := gurucontrollers.CreatePenilaian(&models.TO_PenilaianGuru{
		FinalID:  req.FinalID,
		Skor:     req.Skor,
		Komentar: req.Komentar,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat penilaian"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Penilaian berhasil disimpan"})
}

func UpdatePenilaianHandler(c *gin.Context) {
	penilaianIDStr := c.Param("penilaian_id")
	penilaianID, err := strconv.ParseUint(penilaianIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID penilaian tidak valid"})
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data update tidak valid"})
		return
	}

	if err := gurucontrollers.UpdatePenilaian(uint(penilaianID), data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update penilaian"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Penilaian berhasil diupdate"})
}

func GetPenilaianByFinalHandler(c *gin.Context) {
	finalIDStr := c.Param("final_id")
	finalID, err := strconv.ParseUint(finalIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID final tidak valid"})
		return
	}

	penilaian, err := gurucontrollers.GetPenilaianByFinalID(uint(finalID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Penilaian tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, penilaian)
}
