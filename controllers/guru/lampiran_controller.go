package gurucontrollers

import (
	"errors"
	"strings"
	"time"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

// =========================
// 1. Upload Lampiran
// =========================
func mapFileType(mime string) string {
    if strings.HasPrefix(mime, "image/") {
        return "image"
    } else if mime == "application/pdf" {
        return "pdf"
    } else if strings.HasPrefix(mime, "audio/") {
        return "audio"
    } else if strings.HasPrefix(mime, "video/") {
        return "video"
    }
    return "other"
}

func UploadLampiran(filePath, mimeType, namaFile, deskripsi string) (models.TO_Lampiran, error) {
    lampiran := models.TO_Lampiran{
        NamaFile:  namaFile,
        PathFile:  filePath,
        TipeFile:  mapFileType(mimeType), // ✅ ENUM (image/pdf/audio/video/other)
        Deskripsi: deskripsi,
        CreatedAt: time.Now(),
    }

    err := config.DB.Create(&lampiran).Error
    return lampiran, err
}



// =========================
// 2. Get Lampiran Aktif
// =========================
func GetActiveLampiran() ([]models.TO_Lampiran, error) {
	var data []models.TO_Lampiran
	err := config.DB.
		Where("deleted_at IS NULL").
		Order("created_at desc").
		Find(&data).Error
	return data, err
}

// =========================
// 3. Get Lampiran Trash
// =========================
func GetInactiveLampiran() ([]models.TO_Lampiran, error) {
	var data []models.TO_Lampiran
	err := config.DB.
		Unscoped().
		Where("deleted_at IS NOT NULL").
		Order("created_at desc").
		Find(&data).Error
	return data, err
}

// =========================
// 4. Soft Delete Lampiran
// =========================
func DeleteLampiran(lampiranID uint) error {
	result := config.DB.Delete(&models.TO_Lampiran{}, lampiranID)
	if result.RowsAffected == 0 {
		return errors.New("lampiran tidak ditemukan")
	}
	return result.Error
}

// =========================
// 5. Restore Lampiran
// =========================
func RestoreLampiran(lampiranID uint) error {
	result := config.DB.Model(&models.TO_Lampiran{}).
		Unscoped().
		Where("lampiran_id = ?", lampiranID).
		Update("deleted_at", nil)
	if result.RowsAffected == 0 {
		return errors.New("lampiran tidak ditemukan")
	}
	return result.Error
}

// =========================
// 6. Hard Delete Lampiran
// =========================
func HardDeleteLampiran(lampiranID uint) error {
	result := config.DB.Unscoped().Delete(&models.TO_Lampiran{}, lampiranID)
	if result.RowsAffected == 0 {
		return errors.New("lampiran tidak ditemukan")
	}
	return result.Error
}
