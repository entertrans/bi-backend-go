package siswa

import (
	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

// OnlineClassController menangani logika bisnis untuk online class
type OnlineClassController struct{}

// GetAllOnlineClass mendapatkan semua jadwal online class
func (ctrl *OnlineClassController) GetAllOnlineClass() ([]models.OnlineClass, error) {
	var onlineClasses []models.OnlineClass
	err := config.DB.
		Preload("KelasMapel").
		Preload("KelasMapel.Kelas").
		Preload("KelasMapel.Mapel").
		Order("tanggal asc, mulai asc").
		Find(&onlineClasses).Error
	return onlineClasses, err
}

// GetOnlineClassByID mendapatkan satu jadwal online class by ID
func (ctrl *OnlineClassController) GetOnlineClassByID(id uint) (models.OnlineClass, error) {
	var onlineClass models.OnlineClass
	err := config.DB.
		Preload("KelasMapel").
		Preload("KelasMapel.Kelas").
		Preload("KelasMapel.Mapel").
		Where("id_online_class = ?", id).
		First(&onlineClass).Error
	return onlineClass, err
}

// GetOnlineClassByKelas mendapatkan jadwal online class by kelas_id
func (ctrl *OnlineClassController) GetOnlineClassByKelas(kelasID uint) ([]models.OnlineClass, error) {
	var onlineClasses []models.OnlineClass
	err := config.DB.
		Preload("KelasMapel").
		Preload("KelasMapel.Kelas").
		Preload("KelasMapel.Mapel").
		Joins("JOIN tbl_kelas_mapel km ON km.id_kelas_mapel = tbl_online_class.id_kelas_mapel").
		Where("km.kelas_id = ?", kelasID).
		Order("tanggal asc, mulai asc").
		Find(&onlineClasses).Error
	return onlineClasses, err
}

// GetOnlineClassByMapel mendapatkan jadwal online class by mapel_id
func (ctrl *OnlineClassController) GetOnlineClassByMapel(mapelID uint) ([]models.OnlineClass, error) {
	var onlineClasses []models.OnlineClass
	err := config.DB.
		Preload("KelasMapel").
		Preload("KelasMapel.Kelas").
		Preload("KelasMapel.Mapel").
		Joins("JOIN tbl_kelas_mapel km ON km.id_kelas_mapel = tbl_online_class.id_kelas_mapel").
		Where("km.mapel_id = ?", mapelID).
		Order("tanggal asc, mulai asc").
		Find(&onlineClasses).Error
	return onlineClasses, err
}

// CreateOnlineClass membuat jadwal online class baru
func (ctrl *OnlineClassController) CreateOnlineClass(onlineClass *models.OnlineClass) error {
	err := config.DB.Create(&onlineClass).Error
	return err
}

// UpdateOnlineClass mengupdate jadwal online class
func (ctrl *OnlineClassController) UpdateOnlineClass(id uint, onlineClass *models.OnlineClass) error {
	err := config.DB.Model(&models.OnlineClass{}).
		Where("id_online_class = ?", id).
		Updates(onlineClass).Error
	return err
}

// DeleteOnlineClass menghapus jadwal online class
func (ctrl *OnlineClassController) DeleteOnlineClass(id uint) error {
	err := config.DB.
		Where("id_online_class = ?", id).
		Delete(&models.OnlineClass{}).Error
	return err
}