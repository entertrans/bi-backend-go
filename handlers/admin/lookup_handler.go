package adminhandlers

import (
	"net/http"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

func GetAllAgama(c *gin.Context) {
	var agama []models.Agama
	if err := config.DB.Debug().Find(&agama).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, agama)
}

// kelas
func GetAllKelas(c *gin.Context) {
	var kelasAktif []models.Kelas
	var kelasAlumni []models.Kelas

	// kelas_id < 16 adalah kelas aktif
	if err := config.DB.Where("kelas_id < ?", 16).Find(&kelasAktif).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// kelas_id >= 16 adalah alumni
	if err := config.DB.Where("kelas_id >= ?", 16).Find(&kelasAlumni).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"aktif":  kelasAktif,
		"alumni": kelasAlumni,
	})
}

func GetAllSatelit(c *gin.Context) {
	var satelit []models.DtSatelit
	if err := config.DB.Find(&satelit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, satelit)
}
func GetAllTA(c *gin.Context) {
	var Ta []models.ThnAjaran
	if err := config.DB.Find(&Ta).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, Ta)
}
