package siswa

import (
	"net/http"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

// POST /siswa/jawaban/save
func SaveJawabanHandler(c *gin.Context) {
	var req models.TO_JawabanFinal
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Upsert jawaban (biar kalau sudah ada diupdate aja)
	err := config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "session_id"}, {Name: "soal_id"}},
		UpdateAll: true,
	}).Create(&req).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jawaban tersimpan"})
}
