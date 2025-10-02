package siswa

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

// Ambil semua test UB (Ujian Bulanan)
// func GetAllUBTest() ([]models.TO_Test, error) {
// 	var tests []models.TO_Test
// 	err := config.DB.
// 		Preload("Guru").
// 		Preload("Kelas").
// 		Preload("Mapel").
// 		Where("type_test = ?", "ub").
// 		Order("created_at desc").
// 		Find(&tests).Error
// 	return tests, err
// }

//	func GetUBTestByKelas(kelasID uint) ([]models.TO_Test, error) {
//		var tests []models.TO_Test
//		err := config.DB.
//			Preload("Guru").
//			Preload("Kelas").
//			Preload("Mapel").
//			Where("type_test = ? AND kelas_id = ?", "ub", kelasID).
//			Order("created_at desc").
//			Find(&tests).Error
//		return tests, err
//	}
func GetTestByKelas(kelasID uint, typeTest string) ([]models.TO_Test, error) {
	var tests []models.TO_Test

	query := config.DB.
		Preload("Guru").
		Preload("Kelas").
		Preload("Mapel").
		Where("type_test = ? AND kelas_id = ?", typeTest, kelasID)

	// 🔥 tambahan filter khusus untuk tr & tugas
	if typeTest == "tr" || typeTest == "tugas" {
		query = query.Where("aktif IS NOT NULL")
	}

	err := query.
		Order("created_at desc").
		Find(&tests).Error

	return tests, err
}

// Ambil soal–soal dari test tertentu
func GetSoalByTestID(testID uint) ([]models.TO_TestSoal, error) {
	var soals []models.TO_TestSoal

	err := config.DB.
		Preload("Lampiran").
		Where("test_id = ?", testID).
		Find(&soals).Error

	if err != nil {
		return nil, err
	}

	// Jika masih kosong, cek apakah test_id ada
	if len(soals) == 0 {
		var count int64
		config.DB.Model(&models.TO_TestSoal{}).Where("test_id = ?", testID).Count(&count)
		log.Printf("Jumlah soal untuk test_id %d: %d", testID, count)
	}

	return soals, err
}
// controller/siswaController.go
func GetNotStartedTests(nis string) ([]map[string]interface{}, error) {
	var pesertaList []models.TO_Peserta

	// Ambil semua data peserta dengan status "not_started" untuk NIS tertentu
	err := config.DB.
		Preload("Test").
		Preload("Test.Mapel").
		Preload("Test.Guru").
		Preload("Test.Kelas").
		Where("siswa_nis = ? AND status = ?", nis, "not_started").
		Find(&pesertaList).Error

	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data peserta: %w", err)
	}

	// Filter hanya test yang aktif dan format response
	var result []map[string]interface{}
	for _, peserta := range pesertaList {
		// Skip jika test tidak aktif atau test tidak ada
		if peserta.Test.Aktif == nil || *peserta.Test.Aktif == 0 {
			continue
		}

		testData := map[string]interface{}{
			"test_id":           peserta.TestID,
			"judul":             peserta.Test.Judul,
			"deskripsi":         peserta.Test.Deskripsi,
			"type_test":         peserta.Test.TypeTest,
			"durasi_menit":      peserta.Test.DurasiMenit,
			"deadline":          peserta.Test.Deadline,
			"mapel":             peserta.Test.Mapel.NmMapel,
			"guru":              peserta.Test.Guru.GuruNama,
			"status":            peserta.Status,
			"jumlah_soal_tampil": peserta.Test.Jumlah,
			"random_soal":       peserta.Test.RandomSoal,
		}

		result = append(result, testData)
	}

	// Urutkan berdasarkan created_at descending (terbaru dulu)
	sort.Slice(result, func(i, j int) bool {
		timeI := result[i]["created_at"].(time.Time)
		timeJ := result[j]["created_at"].(time.Time)
		return timeI.After(timeJ)
	})

	return result, nil
}

//testreview
