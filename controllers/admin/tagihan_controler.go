package admincontrollers

import (
	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

func FetchAllTagihan() ([]models.Tagihan, error) {
	var tagihan []models.Tagihan
	err := config.DB.
		Order("id_tagihan DESC").
		Find(&tagihan).Error

	return tagihan, err
}

func TambahTagihan(tagihan models.Tagihan) error {
	return config.DB.Create(&tagihan).Error
}

func UpdateTagihan(id string, jenis string, nominal int) error {
	return config.DB.Model(&models.Tagihan{}).
		Where("id_tagihan = ?", id).
		Updates(map[string]interface{}{
			"jns_tagihan": jenis,
			"nom_tagihan": nominal,
		}).Error
}

func DeleteTagihan(id string) error {
	return config.DB.
		Where("id_tagihan = ?", id).
		Delete(&models.Tagihan{}).Error
}
