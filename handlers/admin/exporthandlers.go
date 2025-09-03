package adminhandlers

import (
	"net/http"
	"strconv"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAllJSONQuestions - Handler dengan pagination
func GetAllJSONQuestions(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	kelasID, _ := strconv.Atoi(c.Query("kelasId"))

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var mataPelajaranList []models.MataPelajaran
	var total int64

	// Build query
	query := config.DB.Model(&models.MataPelajaran{}).Select("id, kelasid, namapelajaran, jsonquestion, sort_by")

	// Add filter if kelasID provided
	if kelasID > 0 {
		query = query.Where("kelasid = ?", kelasID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghitung total data",
			"details": err.Error(),
		})
		return
	}

	// Get data with pagination
	err := query.Offset(offset).Limit(limit).Find(&mataPelajaranList).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data mata pelajaran",
			"details": err.Error(),
		})
		return
	}

	// Return response with pagination info
	c.JSON(http.StatusOK, gin.H{
		"data": mataPelajaranList,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (int(total) + limit - 1) / limit,
		},
	})
}

// GetJSONQuestionsByKelasID - Handler dengan pagination
func GetJSONQuestionsByKelasID(c *gin.Context) {
	kelasIDStr := c.Param("kelasId")
	kelasID, err := strconv.Atoi(kelasIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kelas tidak valid"})
		return
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var mataPelajaranList []models.MataPelajaran
	var total int64

	// Get total count
	if err := config.DB.Model(&models.MataPelajaran{}).
		Where("kelasid = ?", kelasID).
		Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghitung total data",
			"details": err.Error(),
		})
		return
	}

	// Get data with pagination
	err = config.DB.Select("id, kelasid, namapelajaran, jsonquestion, sort_by").
		Where("kelasid = ?", kelasID).
		Offset(offset).Limit(limit).
		Find(&mataPelajaranList).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data mata pelajaran",
			"details": err.Error(),
		})
		return
	}

	// Return response with pagination info
	c.JSON(http.StatusOK, gin.H{
		"data": mataPelajaranList,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (int(total) + limit - 1) / limit,
		},
	})
}

// GetJSONQuestionsByID - Handler untuk single data
func GetJSONQuestionsByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var mataPelajaran models.MataPelajaran

	// Get single data
	err = config.DB.Select("id, kelasid, namapelajaran, jsonquestion, sort_by").
		Where("id = ?", id).
		First(&mataPelajaran).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Mata pelajaran tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data mata pelajaran",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, mataPelajaran)
}