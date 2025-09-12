package gurucontrollers

import (
	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"gorm.io/gorm"
)

// Create test
// func CreateTest(test *models.TO_Test) error {
// 	return config.DB.Create(test).Error
// }

func CreateTest(test *models.TO_Test) error {
	// Step 1: Insert ke to_test
	if err := config.DB.Create(test).Error; err != nil {
		return err
	}

	// Step 2: Insert relasi soal ke to_test_soal
	if len(test.SoalIDs) > 0 {
		var relasi []models.TO_TestSoalRelasi
		for _, soalID := range test.SoalIDs {
			relasi = append(relasi, models.TO_TestSoalRelasi{
				TestID: test.TestID,
				SoalID: soalID,
			})
		}
		if err := config.DB.Create(&relasi).Error; err != nil {
			return err
		}
	}

	return nil
}




// Get test by type (ub / tr/ tugas)
func GetTestByType(tipe string) ([]models.TO_Test, error) {
	var tests []models.TO_Test
	err := config.DB.
		Preload("Guru").
		Preload("Kelas").
		Preload("Mapel").
		Where("type_test = ?", tipe).
		Order("created_at desc").
		Find(&tests).Error
	return tests, err
}
//new

// GetSoalAlreadyInTest mengambil soal yang sudah ada dalam test tertentu
func GetSoalAlreadyInTest(testID uint) ([]uint, error) {
	var soalIDs []uint
	
	err := config.DB.
		Model(&models.TO_TestSoalRelasi{}).
		Where("test_id = ?", testID).
		Pluck("soal_id", &soalIDs).Error

	if err != nil {
		return nil, err
	}

	return soalIDs, nil
}

// GetBankSoalByKelasMapelWithSelection mengambil bank soal dengan info sudah dipilih atau belum
func GetBankSoalByKelasMapel(kelasID, mapelID, testID uint) ([]models.TO_BankSoal, []uint, error) {
	var soals []models.TO_BankSoal
	
	// Ambil bank soal
	err := config.DB.
		Preload("Guru").
		Preload("Kelas").
		Preload("Mapel").
		Preload("Lampiran").
		Where("kelas_id = ? AND mapel_id = ? AND deleted_at IS NULL", kelasID, mapelID).
		Order("created_at DESC").
		Find(&soals).Error

	if err != nil {
		return nil, nil, err
	}

	// Ambil soal yang sudah dipilih
	selectedSoalIDs, err := GetSoalAlreadyInTest(testID)
	if err != nil {
		return nil, nil, err
	}

	return soals, selectedSoalIDs, nil
}

func RemoveSoalFromTest(testID, soalID uint) error {
	result := config.DB.
		Where("test_id = ? AND soal_id = ?", testID, soalID).
		Delete(&models.TO_TestSoalRelasi{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// GetBankSoalByKelasMapelWithFilter mengambil bank soal dengan filter tambahan
func AddSoalToTest(testID uint, soalIDs []uint) error {
	var rels []models.TO_TestSoalRelasi

	for _, soalID := range soalIDs {
		rels = append(rels, models.TO_TestSoalRelasi{
			TestID: testID,
			SoalID: soalID,
		})
	}

	// pakai CreateInBatches biar lebih efisien
	if err := config.DB.CreateInBatches(&rels, 100).Error; err != nil {
		return err
	}
	return nil
}
// end new

// Get all tests by guru_id
func GetTestsByGuruID(guruID uint) ([]models.TO_Test, error) {
	var tests []models.TO_Test
	err := config.DB.
		Preload("Guru").
		Preload("Kelas").
		Preload("Mapel").
		Where("guru_id = ?", guruID).
		Order("created_at desc").
		Find(&tests).Error
	return tests, err
}

// Get single test by ID
func GetTestByID(testID uint) (models.TO_Test, error) {
	var test models.TO_Test
	err := config.DB.
		Preload("Guru").
		Preload("Kelas").
		Preload("Mapel").
		First(&test, testID).Error
	return test, err
}

// Update test
func UpdateTest(testID uint, data map[string]interface{}) error {
	return config.DB.Model(&models.TO_Test{}).
		Where("test_id = ?", testID).
		Updates(data).Error
}

// Delete test
func DeleteTest(testID uint) error {
	return config.DB.Delete(&models.TO_Test{}, testID).Error
}
