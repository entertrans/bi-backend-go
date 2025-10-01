package siswa

import (
	"net/http"
	"strconv"

	"github.com/entertrans/bi-backend-go/config"
	siswaControllers "github.com/entertrans/bi-backend-go/controllers/siswa"
	"github.com/entertrans/bi-backend-go/models"
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

func GetNilaiHandler(c *gin.Context) {
	sessionIDStr := c.Param("session_id")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID tidak valid"})
		return
	}

	db := config.DB
	var session models.TestSession
	if err := db.Preload("Test").Preload("Siswa").First(&session, "session_id = ?", sessionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":  session.SessionID,
		"siswa_nis":   session.SiswaNIS,
		"siswa_nama":  session.Siswa.SiswaNama,
		"test_judul":  session.Test.Judul,
		"nilai_akhir": session.NilaiAkhir,
		"status":      session.Status,
	})
}
