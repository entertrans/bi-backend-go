package gurucontrollers

import (
	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

// Create test
func CreateTest(test *models.TO_Test) error {
	return config.DB.Create(test).Error
}

// Get test by type (ub / quis)
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
