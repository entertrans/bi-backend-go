package gurucontrollers

import (
	"fmt"
	"log"
	"time"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"gorm.io/gorm"
)

// Rangkaian map[string]int `json:"rangkaian"`
type TestJawabanResult struct {
	TestID      uint           `json:"test_id"`
	Jenis       string         `json:"jenis"`
	Mapel       string         `json:"mapel"`
	Judul       string         `json:"judul"`
	Nilai       *float64       `json:"nilai"`
	Tanggal     *time.Time     `json:"tanggal"`
	Status      string         `json:"status"`
	Rangkaian   map[string]int `json:"rangkaian"`
	Submited    bool           `json:"submited"`     // sudah submit jawaban?
	ButuhReview bool           `json:"butuh_review"` // ada esai/isian singkat?
	Reviewed    bool           `json:"reviewed"`     // sudah direview guru?
}

// GetJawabanBySiswaNIS ambil semua hasil test berdasarkan NIS siswa
func GetJawabanBySiswaNIS(siswaNIS string) ([]TestJawabanResult, error) {
	db := config.DB

	// 1. Ambil data siswa
	var siswa models.Siswa
	if err := db.Where("siswa_nis = ?", siswaNIS).
		Preload("Kelas").
		First(&siswa).Error; err != nil {
		return nil, fmt.Errorf("siswa tidak ditemukan")
	}

	if siswa.SiswaKelasID == nil || siswa.Kelas.KelasId == 0 {
		return nil, fmt.Errorf("siswa belum memiliki kelas atau data kelas tidak valid")
	}

	results := []TestJawabanResult{}

	// ---------------------------
	// 2A. Ambil test UB by kelas
	// ---------------------------
	var testsUB []models.TO_Test
	if err := db.Where("kelas_id = ? AND type_test = ?", *siswa.SiswaKelasID, "ub").
		Preload("Mapel").
		Find(&testsUB).Error; err != nil {
		return nil, err
	}

	for _, test := range testsUB {
		res, err := buildTestResult(db, test, siswaNIS)
		if err == nil {
			results = append(results, res)
		}
	}

	// ---------------------------
	// 2B. Ambil test TR/Tugas by peserta
	// ---------------------------
	var pesertaList []models.TO_Peserta
	if err := db.Where("siswa_nis = ?", siswaNIS).
		Preload("Test.Mapel").
		Find(&pesertaList).Error; err != nil {
		return nil, err
	}

	for _, peserta := range pesertaList {
		res, err := buildTestResult(db, peserta.Test, siswaNIS)
		if err == nil {
			results = append(results, res)
		}
	}

	return results, nil
}

// Helper untuk bangun hasil per test
func buildTestResult(db *gorm.DB, test models.TO_Test, siswaNIS string) (TestJawabanResult, error) {
	var session models.TO_TestSession
	err := db.Where("test_id = ? AND siswa_nis = ?", test.TestID, siswaNIS).
		First(&session).Error

	// Hitung rangkaian soal per tipe
	rangkaian := map[string]int{}
	var testSoal []models.TO_TestSoal
	db.Where("test_id = ?", test.TestID).Find(&testSoal)
	for _, ts := range testSoal {
		rangkaian[ts.TipeSoal]++
	}

	// deteksi apakah ada soal yang butuh review manual
	butuhReview := false
	for _, ts := range testSoal {
		if ts.TipeSoal == "uraian" || ts.TipeSoal == "isian_singkat" {
			butuhReview = true
			break
		}
	}

	res := TestJawabanResult{
		TestID:      test.TestID,
		Jenis:       test.TypeTest,
		Mapel:       test.Mapel.NmMapel,
		Judul:       test.Judul,
		Rangkaian:   rangkaian,
		ButuhReview: butuhReview,
		Reviewed:    false, // default
	}

	if err != nil {
		// ❌ belum ada session → belum dikerjakan
		res.Status = "❌ belum dikerjakan"
		res.Reviewed = false
		return res, nil
	}

	// ✅ Ada session
	res.Nilai = &session.NilaiAkhir
	res.Tanggal = &session.StartTime

	if butuhReview {
		// Hitung jumlah soal subjektif
		var totalSubjektif int64
		db.Table("to_jawabanfinal").
			Joins("JOIN to_testsoal ON to_jawabanfinal.soal_id = to_testsoal.soal_id").
			Where("to_jawabanfinal.session_id = ? AND to_testsoal.tipe_soal IN (?)",
				session.SessionID, []string{"uraian", "isian_singkat"}).
			Count(&totalSubjektif)

		// Hitung yang belum dinilai
		var countBelum int64
		db.Table("to_jawabanfinal").
			Joins("JOIN to_testsoal ON to_jawabanfinal.soal_id = to_testsoal.soal_id").
			Where("to_jawabanfinal.session_id = ? AND to_testsoal.tipe_soal IN (?)",
				session.SessionID, []string{"uraian", "isian_singkat"}).
			Where("to_jawabanfinal.skor_uraian IS NULL").
			Count(&countBelum)

	} else {
		// semua objektif → langsung reviewed
		res.Status = "✅ otomatis dinilai"
		res.Reviewed = true
	}
	log.Printf("[DEBUG] test_id=%d siswa_nis=%s butuhReview=%v reviewed=%v status=%s",
		res.TestID, siswaNIS, res.ButuhReview, res.Reviewed, res.Status)

	return res, nil
}

// Fungsi tambahan untuk mendapatkan detail jawaban siswa
func GetDetailJawabanBySession(sessionID uint) (map[string]interface{}, error) {
	db := config.DB

	var session models.TO_TestSession
	if err := db.Preload("Test").
		Preload("Test.Mapel").
		Preload("Siswa").
		First(&session, sessionID).Error; err != nil {
		return nil, fmt.Errorf("session tidak ditemukan")
	}

	// Ambil semua jawaban final
	var jawabanFinal []models.TO_JawabanFinal
	if err := db.Where("session_id = ?", sessionID).
		Preload("Soal").
		Preload("Soal.Lampiran").
		Find(&jawabanFinal).Error; err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"session":       session,
		"jawaban":       jawabanFinal,
		"total_soal":    len(jawabanFinal),
		"nilai_akhir":   session.NilaiAkhir,
		"status":        session.Status,
		"waktu_selesai": session.EndTime,
	}

	return result, nil
}

// Fungsi untuk mendapatkan semua siswa yang mengerjakan test tertentu
func GetSiswaByTest(testID uint) ([]map[string]interface{}, error) {
	db := config.DB

	// Ambil data test
	var test models.TO_Test
	if err := db.Preload("Mapel").Preload("Kelas").First(&test, testID).Error; err != nil {
		return nil, fmt.Errorf("test tidak ditemukan")
	}

	// Ambil semua peserta test
	var peserta []models.TO_Peserta
	if err := db.Where("test_id = ?", testID).
		Preload("Siswa").
		Preload("Siswa.Kelas").
		Find(&peserta).Error; err != nil {
		return nil, err
	}

	results := []map[string]interface{}{}

	for _, p := range peserta {
		// Cek apakah siswa memiliki session
		var session models.TO_TestSession
		sessionErr := db.Where("test_id = ? AND siswa_nis = ?", testID, p.SiswaNIS).
			First(&session).Error

		studentResult := map[string]interface{}{
			"peserta_id":  p.PesertaID,
			"siswa_nis":   p.SiswaNIS,
			"siswa_nama":  p.Siswa.SiswaNama,
			"kelas":       p.Siswa.Kelas.KelasNama,
			"status":      p.Status,
			"nilai_akhir": p.NilaiAkhir,
			"created_at":  p.CreatedAt,
		}

		if sessionErr == nil {
			// Ada session, tambahkan info session
			studentResult["session_id"] = session.SessionID
			studentResult["start_time"] = session.StartTime
			studentResult["end_time"] = session.EndTime
			studentResult["session_status"] = session.Status
		} else {
			// Tidak ada session
			studentResult["session_id"] = nil
			studentResult["session_status"] = "belum_mulai"
		}

		results = append(results, studentResult)
	}

	return results, nil
}

// Fungsi untuk mendapatkan statistik test
func GetTestStatistics(testID uint) (map[string]interface{}, error) {
	db := config.DB

	var totalPeserta int64
	var sudahMengerjakan int64
	var belumMengerjakan int64
	var sudahDinilai int64
	var belumDinilai int64

	// Hitung total peserta
	db.Model(&models.TO_Peserta{}).Where("test_id = ?", testID).Count(&totalPeserta)

	// Hitung yang sudah mengerjakan (punya session)
	db.Model(&models.TO_TestSession{}).Where("test_id = ?", testID).Count(&sudahMengerjakan)
	belumMengerjakan = totalPeserta - sudahMengerjakan

	// Hitung yang sudah dinilai (nilai_akhir > 0 atau status graded)
	db.Model(&models.TO_TestSession{}).
		Where("test_id = ? AND status = ? AND nilai_akhir > 0", testID, "graded").
		Count(&sudahDinilai)

	belumDinilai = sudahMengerjakan - sudahDinilai

	// Ambil rata-rata nilai
	var avgNilai float64
	db.Model(&models.TO_TestSession{}).
		Where("test_id = ? AND status = ?", testID, "graded").
		Select("COALESCE(AVG(nilai_akhir), 0)").
		Scan(&avgNilai)

	// Ambil nilai tertinggi dan terendah
	var maxNilai, minNilai float64
	db.Model(&models.TO_TestSession{}).
		Where("test_id = ? AND status = ?", testID, "graded").
		Select("COALESCE(MAX(nilai_akhir), 0)").
		Scan(&maxNilai)

	db.Model(&models.TO_TestSession{}).
		Where("test_id = ? AND status = ?", testID, "graded").
		Select("COALESCE(MIN(nilai_akhir), 0)").
		Scan(&minNilai)

	result := map[string]interface{}{
		"test_id":           testID,
		"total_peserta":     totalPeserta,
		"sudah_mengerjakan": sudahMengerjakan,
		"belum_mengerjakan": belumMengerjakan,
		"sudah_dinilai":     sudahDinilai,
		"belum_dinilai":     belumDinilai,
		"rata_rata_nilai":   fmt.Sprintf("%.2f", avgNilai),
		"nilai_tertinggi":   fmt.Sprintf("%.2f", maxNilai),
		"nilai_terendah":    fmt.Sprintf("%.2f", minNilai),
	}

	return result, nil
}

func GetSiswaDetailForGuru(siswaNIS string) (map[string]interface{}, error) {
	db := config.DB

	var siswa models.Siswa
	if err := db.Where("siswa_nis = ?", siswaNIS).
		Preload("Kelas").
		First(&siswa).Error; err != nil {
		return nil, fmt.Errorf("siswa tidak ditemukan")
	}

	result := map[string]interface{}{
		"siswa_id":   siswa.SiswaID,
		"siswa_nis":  *siswa.SiswaNIS, // Dereference pointer
		"siswa_nama": *siswa.SiswaNama,
		"kelas_nama": siswa.Kelas.KelasNama,
	}

	return result, nil
}
