package adminhandlers

import (
	"net/http"
	"time"

	"github.com/entertrans/bi-backend-go/config"
	admincontrollers "github.com/entertrans/bi-backend-go/controllers/admin"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func BoolPtr(b bool) *bool {
	return &b
}

func UploadDokumenHandler(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File tidak ditemukan"})
		return
	}

	nis := c.PostForm("nis")
	dokumenJenis := c.PostForm("dokumen_jenis")
	if nis == "" || dokumenJenis == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIS atau dokumen_jenis tidak valid"})
		return
	}

	// Buka file-nya
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka file"})
		return
	}
	defer file.Close()

	// Pakai controller
	url, err := admincontrollers.UploadDokumenController(nis, dokumenJenis, file, fileHeader.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload ke Cloudinary"})
		return
	}
	// Simpan ke database
	lampiran := models.Lampiran{
		SiswaNIS:     nis,
		JenisDokumen: dokumenJenis,
		Url:          url,
		Upload:       time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "siswa_nis"}, {Name: "dokumen_jenis"}},
		DoUpdates: clause.AssignmentColumns([]string{"url", "uploaded_at"}), // ✅ pakai nama kolom db
	}).Create(&lampiran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan ke database"})
		return
	}

	// Kirim URL ke frontend
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil upload",
		"url":     url,
	})
}
