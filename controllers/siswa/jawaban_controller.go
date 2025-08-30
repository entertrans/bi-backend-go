package siswa

import (
	"errors"
	"time"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"gorm.io/gorm"
)

// Simpan / update jawaban siswa (final)
func SaveJawabanFinal(sessionID uint, soalID uint, jawaban string, skorObjektif float64) error {
	var jawabanFinal models.TO_JawabanFinal

	// cek apakah sudah ada jawaban untuk soal ini
	err := config.DB.Where("session_id = ? AND soal_id = ?", sessionID, soalID).First(&jawabanFinal).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// belum ada → insert baru
			jawabanFinal = models.TO_JawabanFinal{
				SessionID:    sessionID,
				SoalID:       soalID,
				JawabanSiswa: jawaban,
				SkorObjektif: skorObjektif,
				UpdatedAt:    time.Now(),
			}
			return config.DB.Create(&jawabanFinal).Error
		}
		return err
	}

	// sudah ada → update
	jawabanFinal.JawabanSiswa = jawaban
	jawabanFinal.SkorObjektif = skorObjektif
	jawabanFinal.UpdatedAt = time.Now()

	return config.DB.Save(&jawabanFinal).Error
}
