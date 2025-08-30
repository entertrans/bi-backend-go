package siswa

import (
	"log"

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

//testreview
