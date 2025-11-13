package admincontrollers

import (
	"net/http"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

func GetAllMapel() ([]models.Mapel, error) {
    var mapels []models.Mapel
    err := config.DB.
        Order("nm_mapel ASC").
        Find(&mapels).Error
    return mapels, err
}

func GetAllMapelWithGuruMapels() ([]models.Mapel, error) {
    var mapels []models.Mapel
    err := config.DB.
        Preload("GuruMapels").
        Order("nm_mapel ASC").
        Find(&mapels).Error
    return mapels, err
}


type KelasMapelWithGuru struct {
    PelajaranID uint   `json:"pelajaran_id"`
    KelasID     uint   `json:"kelas_id"`
    KelasNama   string `json:"kelas_nama"`
    NmMapel     string `json:"nm_mapel"`
    GuruMapels  string `json:"guru_mapels"`
}

func GetKelasMapelWithGuru() ([]KelasMapelWithGuru, error) {
    var kelasMapels []models.KelasMapel

    // Preload seluruh relasi yang dibutuhkan
    err := config.DB.
        Preload("Kelas").
        Preload("Mapel.GuruMapels.Guru").
		Order("kelas_id ASC").
        Find(&kelasMapels).Error
    if err != nil {
        return nil, err
    }

    var result []KelasMapelWithGuru
    for _, km := range kelasMapels {
        // Ambil nama guru pertama yang aktif (kalau ada)
        var guruNama string
        for _, gm := range km.Mapel.GuruMapels {
            if gm.KelasID == km.KelasID && gm.StatusAktif {
                guruNama = gm.Guru.GuruNama
                break
            }
        }

        result = append(result, KelasMapelWithGuru{
            PelajaranID: km.ID,
            KelasID:     km.KelasID,
            KelasNama:   km.Kelas.KelasNama,
            NmMapel:     km.Mapel.NmMapel,
            GuruMapels:  guruNama,
        })
    }

    return result, nil
}

func GetMapelByKelas(c *gin.Context) {
	kelasID := c.Param("id") // ambil dari /mapel-by-kelas/:id
	if kelasID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "kelas_id wajib diisi"})
		return
	}

	var kelasMapels []models.KelasMapel
	err := config.DB.
		Preload("Mapel").
		Where("kelas_id = ?", kelasID).
		Find(&kelasMapels).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data mapel"})
		return
	}

	c.JSON(http.StatusOK, kelasMapels)
}

func GetSiswaByKelas(kelasID uint) ([]models.Siswa, error) {
	var siswa []models.Siswa
	err := config.DB.
		Preload("Kelas").
		// Preload("Agama").
		// Preload("Satelit").
		Where("siswa_kelas_id = ? AND (soft_deleted IS NULL OR soft_deleted = 0)", kelasID).
		Order("siswa_nama ASC").
		Find(&siswa).Error

	return siswa, err
}

