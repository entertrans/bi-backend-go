package adminhandlers

import (
	"net/http"
	"strconv"

	"github.com/entertrans/bi-backend-go/config"
	adminControllers "github.com/entertrans/bi-backend-go/controllers/admin"
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

func GetMapelByKelas(c *gin.Context) {
	adminControllers.GetMapelByKelas(c)
}
func GetDetailLookup(c *gin.Context) {
	kelasIDStr := c.Param("kelas_id")
	mapelIDStr := c.Param("mapel_id")

	kelasID, err := strconv.Atoi(kelasIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "kelas_id tidak valid"})
		return
	}

	mapelID, err := strconv.Atoi(mapelIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "mapel_id tidak valid"})
		return
	}

	var kelas models.Kelas
	if err := config.DB.First(&kelas, "kelas_id = ?", kelasID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "kelas tidak ditemukan"})
		return
	}

	var mapel models.Mapel
	if err := config.DB.First(&mapel, "kd_mapel = ?", mapelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "mapel tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"kelas": gin.H{
				"kd_kelas": kelas.KelasId,
				"nm_kelas": kelas.KelasNama,
			},
			"mapel": gin.H{
				"kd_mapel": mapel.KdMapel,
				"nm_mapel": mapel.NmMapel,
			},
		},
	})
}
