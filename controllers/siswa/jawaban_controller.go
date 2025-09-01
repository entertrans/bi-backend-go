package siswa

import (
	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"gorm.io/datatypes"
)

// Simpan / update jawaban siswa (final)
func SaveJawabanFinal(sessionID uint, soalID uint, jawaban string, skorObjektif float64) error {
	var jawabanFinal models.JawabanFinal

	// cek apakah sudah ada jawaban final untuk soal ini
	err := config.DB.Where("session_id = ? AND soal_id = ?", sessionID, soalID).
		First(&jawabanFinal).Error

	if err == nil {
		// update jawaban existing
		jawabanFinal.JawabanSiswa = datatypes.JSON([]byte(jawaban))
		jawabanFinal.SkorObjektif = skorObjektif
		return config.DB.Save(&jawabanFinal).Error
	}

	// insert baru
	newJawaban := models.JawabanFinal{
		SessionID:    sessionID,
		SoalID:       soalID,
		JawabanSiswa: datatypes.JSON([]byte(jawaban)),
		SkorObjektif: skorObjektif,
	}
	return config.DB.Create(&newJawaban).Error
}
