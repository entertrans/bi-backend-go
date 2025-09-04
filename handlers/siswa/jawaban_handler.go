package siswa

import (
	"net/http"

	siswaControllers "github.com/entertrans/bi-backend-go/controllers/siswa"
	"github.com/gin-gonic/gin"
)

func SaveJawabanHandler(c *gin.Context) {
	var req struct {
		SessionID    uint    `json:"session_id"`
		SoalID       uint    `json:"soal_id"`
		JawabanSiswa string  `json:"jawaban_siswa"`
		SkorObjektif float64 `json:"skor_objektif"`
		TipeSoal     string  `json:"tipe_soal"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	err := siswaControllers.SaveJawabanFinal(
		req.SessionID,
		req.SoalID,
		req.JawabanSiswa,
		req.SkorObjektif,
		req.TipeSoal,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan jawaban"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jawaban tersimpan"})
}
