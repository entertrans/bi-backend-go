package siswa

import (
	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

// Ambil semua test UB (Ujian Bulanan)
func GetAllUBTest() ([]models.TO_Test, error) {
	var tests []models.TO_Test
	err := config.DB.
		Preload("Guru").
		Preload("Kelas").
		Preload("Mapel").
		Where("type_test = ?", "ub").
		Order("created_at desc").
		Find(&tests).Error
	return tests, err
}
func GetUBTestByKelas(kelasID uint) ([]models.TO_Test, error) {
	var tests []models.TO_Test
	err := config.DB.
		Preload("Guru").
		Preload("Kelas").
		Preload("Mapel").
		Where("type_test = ? AND kelas_id = ?", "ub", kelasID).
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
	return soals, err
}
