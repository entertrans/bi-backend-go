package admincontrollers

import (
	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

func GetAllMapel() ([]models.Mapel, error) {
    var mapels []models.Mapel
    err := config.DB.
        Order("nm_mapel ASC").
        Find(&mapels).Error
    return mapels, err
}

func GetAllMapelWithGuruMapels() ([]models.Mapel, error) {
    var mapels []models.Mapel
    err := config.DB.
        Preload("GuruMapels").
        Order("nm_mapel ASC").
        Find(&mapels).Error
    return mapels, err
}
