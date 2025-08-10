package gurucontrollers

import (
	"time"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

func CreateTest(test *models.TO_Test) error {
	test.CreatedAt = time.Now()
	return config.DB.Create(test).Error
}

func GetTestByID(testID uint) (models.TO_Test, error) {
	var test models.TO_Test
	err := config.DB.First(&test, "test_id = ?", testID).Error
	return test, err
}

func GetTestsByGuruID(guruID uint) ([]models.TO_Test, error) {
	var tests []models.TO_Test
	err := config.DB.Where("guru_id = ?", guruID).Find(&tests).Error
	return tests, err
}

func UpdateTest(testID uint, data map[string]interface{}) error {
	return config.DB.Model(&models.TO_Test{}).
		Where("test_id = ?", testID).
		Updates(data).Error
}

func DeleteTest(testID uint) error {
	return config.DB.Delete(&models.TO_Test{}, testID).Error
}
