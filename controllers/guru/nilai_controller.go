package gurucontrollers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

// Helper function untuk mendapatkan info kelas dan mapel
func getKelasMapelInfo(kelasID, mapelID string) (map[string]interface{}, map[string]interface{}) {
	var kelas models.Kelas
	var mapel models.Mapel

	// Ambil info kelas
	kelasInfo := make(map[string]interface{})
	if err := config.DB.First(&kelas, kelasID).Error; err == nil {
		kelasInfo = map[string]interface{}{
			"kelas_id":   kelas.KelasId,
			"kelas_nama": kelas.KelasNama,
		}
	} else {
		kelasInfo = map[string]interface{}{
			"kelas_id":   kelasID,
			"kelas_nama": "Unknown",
		}
	}

	// Ambil info mapel
	mapelInfo := make(map[string]interface{})
	if err := config.DB.First(&mapel, mapelID).Error; err == nil {
		mapelInfo = map[string]interface{}{
			"kd_mapel": mapel.KdMapel,
			"nm_mapel": mapel.NmMapel,
		}
	} else {
		mapelInfo = map[string]interface{}{
			"kd_mapel": mapelID,
			"nm_mapel": "Unknown",
		}
	}

	return kelasInfo, mapelInfo
}

// Helper function untuk mendapatkan detail test
func getDetailTest(c *gin.Context, testType string) {
	kelasID := c.Param("kelas_id")
	mapelID := c.Param("mapel_id")

	var tests []models.TO_Test

	err := config.DB.
		Preload("Guru").
		Preload("Mapel").
		Preload("Kelas").
		Where("kelas_id = ? AND mapel_id = ? AND type_test = ?", kelasID, mapelID, testType).
		Find(&tests).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Gagal ambil data detail %s", strings.ToUpper(testType))})
		return
	}

	// Ambil info kelas dan mapel
	kelasInfo, mapelInfo := getKelasMapelInfo(kelasID, mapelID)

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

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"kelas": kelasInfo,
		"mapel": mapelInfo,
	})
}

// GetDetailUB - Get detail UB
func (nc *NilaiController) GetDetailUB(c *gin.Context) {
	getDetailTest(c, "ub")
}

// GetDetailTR - Get detail TR
func (nc *NilaiController) GetDetailTR(c *gin.Context) {
	getDetailTest(c, "tr")
}

// GetDetailTugas - Get detail Tugas
func (nc *NilaiController) GetDetailTugas(c *gin.Context) {
	getDetailTest(c, "tugas")
}

// GetDetailPesertaTest - Get detail peserta test
func (nc *NilaiController) GetDetailPesertaTest(c *gin.Context) {
	testID := c.Param("test_id")
	kelasID := c.Param("kelas_id")

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

	type TestInfo struct {
		TestID      uint      `json:"test_id"`
		Judul       string    `json:"judul"`
		TypeTest    string    `json:"type_test"`
		DurasiMenit int       `json:"durasi_menit"`
		JumlahSoal  uint      `json:"jumlah_soal"`
		GuruNama    string    `json:"guru_nama"`
		KelasNama   string    `json:"kelas_nama"`
		MapelNama   string    `json:"mapel_nama"`
		CreatedAt   time.Time `json:"created_at"`
	}

	var results []PesertaResponse
	var testInfo TestInfo

	// Convert testID ke uint
	testIDUint, convErr := strconv.ParseUint(testID, 10, 32)
	if convErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Test ID tidak valid"})
		return
	}

	// Ambil informasi test dengan Preload (lebih aman)
	var test models.TO_Test
	err := config.DB.
		Preload("Guru").
		Preload("Kelas").
		Preload("Mapel").
		First(&test, testIDUint).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil info test: " + err.Error()})
		return
	}

	// Map ke TestInfo
	testInfo = TestInfo{
		TestID:      test.TestID,
		Judul:       test.Judul,
		TypeTest:    test.TypeTest,
		DurasiMenit: test.DurasiMenit,
		JumlahSoal:  test.Jumlah, // Menggunakan field Jumlah dari model
		GuruNama:    test.Guru.GuruNama,
		KelasNama:   test.Kelas.KelasNama,
		MapelNama:   test.Mapel.NmMapel,
		CreatedAt:   test.CreatedAt,
	}

	// Query join siswa + kelas + testsession
	query := config.DB.
		Table("tbl_siswa s").
		Select(`
			s.siswa_id,
			COALESCE(s.siswa_nis, '') as nis,
			COALESCE(s.siswa_nama, '') as nama,
			k.kelas_nama,
			CASE 
				WHEN ts.session_id IS NOT NULL THEN ts.status
				ELSE 'not_started'
			END as status,
			ts.nilai_akhir as nilai,
			ts.start_time as waktu_mulai,
			ts.end_time as waktu_selesai
		`).
		Joins("JOIN tbl_kelas k ON k.kelas_id = s.siswa_kelas_id").
		Joins("LEFT JOIN to_testsession ts ON ts.siswa_nis = s.siswa_nis AND ts.test_id = ?", testIDUint).
		Where("s.siswa_kelas_id = ? AND s.soft_deleted = 0", kelasID).
		Order("k.kelas_nama, s.siswa_nama")

	if err := query.Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data peserta: " + err.Error()})
		return
	}

	// Mapping type test ke label yang lebih user-friendly
	typeLabels := map[string]string{
		"ub":    "Ulangan Bulanan",
		"tr":    "Test Review",
		"tugas": "Tugas",
	}

	testInfo.TypeTest = typeLabels[testInfo.TypeTest]

	c.JSON(http.StatusOK, gin.H{
		"kelas_id":  kelasID,
		"test_id":   testID,
		"test_info": testInfo,
		"data":      results,
	})
}
