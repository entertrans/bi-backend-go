package admincontrollers

import (
	"net/http"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

// search
func SearchSiswa(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter q tidak boleh kosong"})
		return
	}

	var siswa []models.Siswa

	err := config.DB.
		Preload("Kelas").
		Where("soft_deleted = ? AND siswa_kelas_id < ? AND (LOWER(siswa_nama) LIKE ? OR siswa_nis LIKE ?)",
			0, 16, "%"+query+"%", "%"+query+"%").
		Limit(5).
		Find(&siswa).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari siswa"})
		return
	}

	var result []gin.H
	for _, s := range siswa {
		result = append(result, gin.H{
			"nis":   s.SiswaNIS,
			"nama":  s.SiswaNama,
			"kelas": s.Kelas.KelasNama,
		})
	}

	c.JSON(http.StatusOK, result)
}

// READ
func FetchAllSiswa() ([]models.Siswa, error) {
	var siswa []models.Siswa
	err := config.DB.
		Preload("Orangtua").
		Preload("Lampiran").
		Preload("Agama").
		Find(&siswa).Error

	return siswa, err
}

// GET ALL LENGKAP
func FetchAllSiswaAktif() ([]models.Siswa, error) {
	var siswa []models.Siswa
	err := config.DB.
		Where("soft_deleted = ? AND siswa_kelas_id < ?", 0, 16).
		Preload("Orangtua").
		Preload("Lampiran").
		Preload("Kelas").
		Preload("Satelit").
		Preload("Agama").
		Find(&siswa).Error

	return siswa, err
}

func FetchAllSiswaKeluar() ([]models.Siswa, error) {
	var siswa []models.Siswa
	err := config.DB.
		Where("soft_deleted = ?", 1).
		Preload("Orangtua").
		Preload("Kelas").
		Preload("Satelit").
		Preload("Lampiran").
		Preload("Agama").
		Find(&siswa).Error

	return siswa, err
}

func FetchAllSiswaPPDB() ([]models.Siswa, error) {
	var siswa []models.Siswa
	err := config.DB.
		Where("soft_deleted = ?", 2).
		Preload("Orangtua").
		Preload("Lampiran").
		Preload("Kelas").
		Preload("Satelit").
		Preload("Agama").
		Find(&siswa).Error

	return siswa, err
}

func FetchAllSiswaAlumni() ([]models.Siswa, error) {
	var siswa []models.Siswa
	err := config.DB.
		Where("siswa_kelas_id > ?", 15).
		Preload("Orangtua").
		Preload("Kelas").
		Preload("Satelit").
		Preload("Agama").
		Find(&siswa).Error

	return siswa, err
}

// GET BY siswa_nis
func FindSiswaByNis(nis string) (models.Siswa, error) {
	var siswa models.Siswa
	err := config.DB.
		Where("siswa_nis = ?", nis).
		Preload("Orangtua").
		Preload("Lampiran").
		Preload("Kelas").
		Preload("Satelit").
		Preload("Agama").
		First(&siswa).Error
	return siswa, err
}

// GET siswa + ortu
func GetSiswaWithOrtu(nis string) (*models.Siswa, error) {
	var siswa models.Siswa

	err := config.DB.
		Where("siswa_nis = ?", nis).
		Preload("Orangtua").
		Preload("Kelas").
		Preload("Satelit").
		Preload("Agama").
		First(&siswa).Error

	if err != nil {
		return nil, err
	}
	return &siswa, nil
}

//Update

func TerimaSiswa(nis string) error {
	return config.DB.Model(&models.Siswa{}).
		Where("siswa_nis = ?", nis).
		Updates(map[string]interface{}{
			"soft_deleted": 0,
			"tgl_keluar":   nil,
		}).Error
}

func SetKelasOnline(nis string, newValue int) error {
	return config.DB.Model(&models.Siswa{}).
		Where("siswa_nis = ?", nis).
		Update("oc", newValue).Error
}

func SetKelasOffline(nis string, newValue int) error {
	return config.DB.Model(&models.Siswa{}).
		Where("siswa_nis = ?", nis).
		Update("kc", newValue).Error
}

func KeluarkanSiswa(nis string, tglKeluar *string) error {
	return config.DB.Model(&models.Siswa{}).
		Where("siswa_nis = ?", nis).
		Updates(map[string]interface{}{
			"soft_deleted": 1,
			"tgl_keluar":   tglKeluar,
		}).Error
}
