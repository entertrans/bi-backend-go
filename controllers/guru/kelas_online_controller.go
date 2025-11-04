package gurucontrollers

import (
	"time"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

type KelasOnlineResponse struct {
	IDKelasOnline uint      `json:"id_kelas_online"`
	JudulKelas    string    `json:"judul_kelas"`
	TanggalKelas  time.Time `json:"tanggal_kelas"`
	JamMulai      string    `json:"jam_mulai"`
	JamSelesai    string    `json:"jam_selesai"`
	Status        string    `json:"status"`
	LinkKelas     string    `json:"link_kelas"`
	MateriLink    string    `json:"materi_link"`
	GuruNama      string    `json:"guru_nama"`
	MapelNama     string    `json:"mapel_nama"`
	KelasNama     string    `json:"kelas_nama"`
}

//
// 📄 Controller Functions
//
func GetAllKelasOnline() ([]KelasOnlineResponse, error) {
	var kelas []models.KelasOnline
	err := config.DB.
		Preload("Guru").
		Preload("KelasMapel.Kelas").
		Preload("KelasMapel.Mapel").
		Find(&kelas).Error
	if err != nil {
		return nil, err
	}

	// Map ke DTO
	var results []KelasOnlineResponse
	for _, k := range kelas {
		results = append(results, KelasOnlineResponse{
			IDKelasOnline: k.IDKelasOnline,
			JudulKelas:    k.JudulKelas,
			TanggalKelas:  k.TanggalKelas,
			JamMulai:      k.JamMulai,
			JamSelesai:    k.JamSelesai,
			Status:        k.Status,
			LinkKelas:     k.LinkKelas,
			MateriLink:    k.MateriLink,
			GuruNama:      k.Guru.GuruNama,            // asumsikan field Guru.Nama ada
			MapelNama:     k.KelasMapel.Mapel.NmMapel, // asumsikan field Mapel.Nama ada
			KelasNama:     k.KelasMapel.Kelas.KelasNama, // asumsikan field Kelas.Nama ada
		})
	}
	return results, nil
}

func GetKelasOnlineByID(id string) (KelasOnlineResponse, error) {
	var k models.KelasOnline
	err := config.DB.
		Preload("Guru").
		Preload("KelasMapel.Kelas").
		Preload("KelasMapel.Mapel").
		First(&k, id).Error
	if err != nil {
		return KelasOnlineResponse{}, err
	}

	return KelasOnlineResponse{
		IDKelasOnline: k.IDKelasOnline,
		JudulKelas:    k.JudulKelas,
		TanggalKelas:  k.TanggalKelas,
		JamMulai:      k.JamMulai,
		JamSelesai:    k.JamSelesai,
		Status:        k.Status,
		LinkKelas:     k.LinkKelas,
		MateriLink:    k.MateriLink,
		GuruNama:      k.Guru.GuruNama,
		MapelNama:     k.KelasMapel.Mapel.NmMapel,
		KelasNama:     k.KelasMapel.Kelas.KelasNama,
	}, nil
}

func CreateKelasOnline(data models.KelasOnline) (models.KelasOnline, error) {
	if err := config.DB.Create(&data).Error; err != nil {
		return models.KelasOnline{}, err
	}
	return data, nil
}

func UpdateKelasOnline(id string, data models.KelasOnline) (models.KelasOnline, error) {
	var existing models.KelasOnline
	if err := config.DB.First(&existing, id).Error; err != nil {
		return models.KelasOnline{}, err
	}
	if err := config.DB.Model(&existing).Updates(data).Error; err != nil {
		return models.KelasOnline{}, err
	}
	return existing, nil
}

func DeleteKelasOnline(id string) error {
	var kelas models.KelasOnline
	if err := config.DB.First(&kelas, id).Error; err != nil {
		return err
	}
	return config.DB.Delete(&kelas).Error
}