package gurucontrollers

import (
	"fmt"
	"time"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

type TestJawabanResult struct {
	TestID    uint           `json:"test_id"`
	Jenis     string         `json:"jenis"`
	Mapel     string         `json:"mapel"`
	Judul     string         `json:"judul"`
	Nilai     *float64       `json:"nilai"`
	Tanggal   *time.Time     `json:"tanggal"`
	Rangkaian map[string]int `json:"rangkaian"`
	Status    string         `json:"status"`
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

	// Pastikan SiswaKelasID valid (bukan pointer nil)
	// Alternatif: cek apakah Kelas memiliki data yang valid
if siswa.SiswaKelasID == nil || siswa.Kelas.KelasId == 0 {
    return nil, fmt.Errorf("siswa belum memiliki kelas atau data kelas tidak valid")
}

	// 2. Ambil semua test di kelas siswa
	var tests []models.TO_Test
	if err := db.Where("kelas_id = ?", *siswa.SiswaKelasID).
		Preload("Mapel").
		Find(&tests).Error; err != nil {
		return nil, err
	}

	results := []TestJawabanResult{}

	for _, test := range tests {
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

		// Build result
		res := TestJawabanResult{
			TestID:    test.TestID,
			Jenis:     test.TypeTest,
			Mapel:     test.Mapel.NmMapel,
			Judul:     test.Judul,
			Rangkaian: rangkaian,
		}

		if err != nil {
			// ❌ belum ada session
			res.Status = "❌ belum dikerjakan"
		} else {
			// ✅ Ada session
			res.Nilai = &session.NilaiAkhir
			res.Tanggal = &session.StartTime

			// Hitung soal subjektif yang belum dinilai
			var countBelum int64
			db.Model(&models.TO_JawabanFinal{}).
				Joins("JOIN to_testsoal ON to_jawabanfinal.soal_id = to_testsoal.soal_id").
				Where("to_jawabanfinal.session_id = ? AND to_testsoal.tipe_soal IN (?)", 
					session.SessionID, []string{"uraian", "isian_singkat"}).
				Where("to_jawabanfinal.skor_uraian IS NULL").
				Count(&countBelum)

			if countBelum > 0 {
				res.Status = fmt.Sprintf("⚠️ %d belum dinilai", countBelum)
			} else {
				res.Status = "✅ semua dinilai"
			}
		}

		results = append(results, res)
	}

	return results, nil
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
        "siswa_id":    siswa.SiswaID,
        "siswa_nis":   *siswa.SiswaNIS, // Dereference pointer
        "siswa_nama":  *siswa.SiswaNama,
        "kelas_nama":  siswa.Kelas.KelasNama,
    }

    return result, nil
}
