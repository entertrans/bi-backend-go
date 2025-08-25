package gurucontrollers

import (
	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

// Tambah peserta
func AddPeserta(peserta *models.TO_Peserta) error {
	return config.DB.Create(peserta).Error
}

// Tambah banyak peserta sekaligus
func AddPesertaBatch(pesertas []models.TO_Peserta) error {
	return config.DB.Create(&pesertas).Error
}

// Ambil semua peserta dalam test tertentu
func GetPesertaByTestID(testID uint) ([]models.TO_Peserta, error) {
	var pes []models.TO_Peserta
	err := config.DB.Preload("Siswa").Where("test_id = ?", testID).Find(&pes).Error
	return pes, err
}

// Update peserta (status, nilai, extra time)
func UpdatePeserta(pesertaID uint, data map[string]interface{}) error {
	return config.DB.Model(&models.TO_Peserta{}).
		Where("peserta_id = ?", pesertaID).
		Updates(data).Error
}

// Hapus peserta
func DeletePeserta(pesertaID uint) error {
	return config.DB.Delete(&models.TO_Peserta{}, pesertaID).Error
}

func GetAvailableSiswaByKelas(kelasID, testID uint) ([]models.Siswa, error) {
	var siswa []models.Siswa

	// ambil NIS siswa yang sudah jadi peserta test
	var pesertaNIS []string
	config.DB.Table("to_peserta").Select("siswa_nis").
		Where("test_id = ?", testID).
		Find(&pesertaNIS)

	// query siswa kelas tertentu, exclude yang sudah masuk peserta
	query := config.DB.
		Preload("Kelas").
		Preload("Agama").
		Preload("Satelit").
		Where("siswa_kelas_id = ? AND (soft_deleted IS NULL OR soft_deleted = 0)", kelasID)

	if len(pesertaNIS) > 0 {
		query = query.Where("siswa_nis NOT IN ?", pesertaNIS)
	}

	err := query.Order("siswa_nama ASC").Find(&siswa).Error
	return siswa, err
}
