package gurucontrollers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

type NilaiController struct{}

func NewNilaiController() *NilaiController {
	return &NilaiController{}
}

// RekapNilaiResponse struct
type RekapNilaiResponse struct {
	KelasID   uint   `json:"kelas_id"`
	KelasNama string `json:"kelas_nama"`
	MapelID   uint   `json:"mapel_id"`
	MapelNama string `json:"mapel_nama"`
	JmlUB     int    `json:"jml_ub"`
	JmlTR     int    `json:"jml_tr"`
	JmlTugas  int    `json:"jml_tugas"`
}

// GetRekapNilai - Get rekap nilai
func (nc *NilaiController) GetRekapNilai(c *gin.Context) {
	var results []RekapNilaiResponse

	// Subquery untuk mendapatkan semua kombinasi kelas dan mapel
	subQuery := config.DB.
		Table("tbl_kelas_mapel km").
		Select(`
			k.kelas_id,
			k.kelas_nama,
			m.kd_mapel AS mapel_id,
			m.nm_mapel AS mapel_nama
		`).
		Joins("JOIN tbl_kelas k ON k.kelas_id = km.kelas_id").
		Joins("JOIN tbl_mapel m ON m.kd_mapel = km.kd_mapel")

	// Query utama dengan LEFT JOIN ke to_test untuk menghitung jumlah test
	err := config.DB.
		Table("(?) AS km", subQuery).
		Select(`
			km.kelas_id,
			km.kelas_nama,
			km.mapel_id,
			km.mapel_nama,
			COALESCE(SUM(CASE WHEN t.type_test = 'ub' THEN 1 ELSE 0 END), 0) AS jml_ub,
			COALESCE(SUM(CASE WHEN t.type_test = 'tr' THEN 1 ELSE 0 END), 0) AS jml_tr,
			COALESCE(SUM(CASE WHEN t.type_test = 'tugas' THEN 1 ELSE 0 END), 0) AS jml_tugas
		`).
		Joins("LEFT JOIN to_test t ON t.kelas_id = km.kelas_id AND t.mapel_id = km.mapel_id").
		Group("km.kelas_id, km.kelas_nama, km.mapel_id, km.mapel_nama").
		Order("km.kelas_id, km.mapel_id").
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data rekap: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetDetailUB - Get detail UB
func (nc *NilaiController) GetDetailUB(c *gin.Context) {
	kelasID := c.Param("kelas_id")
	mapelID := c.Param("mapel_id")

	var tests []models.TO_Test
	
	// Query untuk mendapatkan semua test UB untuk kelas dan mapel tertentu
	err := config.DB.
		Preload("Guru").
		Preload("Mapel").
		Preload("Kelas").
		Where("kelas_id = ? AND mapel_id = ? AND type_test = 'ub'", kelasID, mapelID).
		Find(&tests).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data detail UB"})
		return
	}

	// Format response
	var response []map[string]interface{}
	for _, test := range tests {
		response = append(response, map[string]interface{}{
			"test_id":      test.TestID,
			"judul":        test.Judul,
			"guru_nama":    test.Guru.GuruNama,
			"durasi_menit": test.DurasiMenit,
			"jumlah_soal":  test.Jumlah,
			"created_at":   test.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetDetailTR - Get detail TR
func (nc *NilaiController) GetDetailTR(c *gin.Context) {
	kelasID := c.Param("kelas_id")
	mapelID := c.Param("mapel_id")

	var tests []models.TO_Test
	
	err := config.DB.
		Preload("Guru").
		Preload("Mapel").
		Preload("Kelas").
		Where("kelas_id = ? AND mapel_id = ? AND type_test = 'tr'", kelasID, mapelID).
		Find(&tests).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data detail TR"})
		return
	}

	var response []map[string]interface{}
	for _, test := range tests {
		response = append(response, map[string]interface{}{
			"test_id":      test.TestID,
			"judul":        test.Judul,
			"guru_nama":    test.Guru.GuruNama,
			"durasi_menit": test.DurasiMenit,
			"jumlah_soal":  test.Jumlah,
			"created_at":   test.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetDetailTugas - Get detail Tugas
func (nc *NilaiController) GetDetailTugas(c *gin.Context) {
	kelasID := c.Param("kelas_id")
	mapelID := c.Param("mapel_id")

	var tests []models.TO_Test
	
	err := config.DB.
		Preload("Guru").
		Preload("Mapel").
		Preload("Kelas").
		Where("kelas_id = ? AND mapel_id = ? AND type_test = 'tugas'", kelasID, mapelID).
		Find(&tests).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data detail Tugas"})
		return
	}

	var response []map[string]interface{}
	for _, test := range tests {
		response = append(response, map[string]interface{}{
			"test_id":      test.TestID,
			"judul":        test.Judul,
			"guru_nama":    test.Guru.GuruNama,
			"durasi_menit": test.DurasiMenit,
			"jumlah_soal":  test.Jumlah,
			"created_at":   test.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetDetailPesertaTest - Get detail peserta test
func (nc *NilaiController) GetDetailPesertaTest(c *gin.Context) {
	testID := c.Param("test_id")

	type PesertaResponse struct {
		SiswaID      uint       `json:"siswa_id"`
		NIS          string     `json:"nis"`
		Nama         string     `json:"nama"`
		KelasNama    string     `json:"kelas_nama"`
		Status       string     `json:"status"`
		Nilai        *float64   `json:"nilai"`
		WaktuMulai   *time.Time `json:"waktu_mulai"`
		WaktuSelesai *time.Time `json:"waktu_selesai"`
	}

	var results []PesertaResponse

	// Convert params to uint
	testIDUint, err := strconv.ParseUint(testID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Test ID tidak valid"})
		return
	}

	// Query tanpa filter kelas, semua peserta test
	err = config.DB.
		Table("tbl_siswa s").
		Select(`
			s.siswa_id,
			s.nis,
			s.nama,
			k.kelas_nama,
			CASE 
				WHEN ts.session_id IS NOT NULL THEN 'sudah_ngerjain' 
				ELSE 'belum_ngerjain' 
			END as status,
			ts.nilai_akhir as nilai,
			ts.start_time as waktu_mulai,
			ts.end_time as waktu_selesai
		`).
		Joins("JOIN tbl_kelas k ON k.kelas_id = s.kelas_id").
		Joins("LEFT JOIN to_testsession ts ON ts.siswa_nis = s.siswa_id AND ts.test_id = ?", testIDUint).
		Where("s.siswa_id IN (SELECT siswa_id FROM tbl_siswa WHERE kelas_id IN (SELECT kelas_id FROM to_test WHERE test_id = ?))", testIDUint).
		Order("k.kelas_nama, s.nama").
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data peserta: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}