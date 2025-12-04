package gurucontrollers

import (
	"fmt"
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
type KelasOnlineListResponse struct {
    IDKelasOnline uint   `json:"id_kelas_online"`
    JudulKelas    string `json:"judul_kelas"`
    NamaGuru      string `json:"nama_guru"`
    NamaMapel     string `json:"nama_mapel"`
    TanggalKelas  string `json:"tanggal_kelas"`
    JamMulai      string `json:"jam_mulai"`
    JamSelesai    string `json:"jam_selesai"`
    Status        string `json:"status"`
    LinkKelas     string `json:"link_kelas"`
}
type KelasOnlineDetailResponse struct {
    ID      uint            `json:"id"`
    Topik   string          `json:"topik"`
    Tanggal string          `json:"tanggal"`
    Guru    string          `json:"guru"`
    Materi  []MateriResponse `json:"materi"`
}

type MateriResponse struct {
    ID         uint      `json:"id"`
    Judul      string    `json:"judul"`
    Tipe       string    `json:"tipe"`
    Link       string    `json:"link"`       // dari UrlFile
    Keterangan string    `json:"keterangan"`
    UploadedAt string    `json:"uploaded_at"` // format string agar mudah dibaca di client
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
			GuruNama:      k.Guru.GuruNama,            // asumsikan field Guru.Nama ada
			MapelNama:     k.KelasMapel.Mapel.NmMapel, // asumsikan field Mapel.Nama ada
			KelasNama:     k.KelasMapel.Kelas.KelasNama, // asumsikan field Kelas.Nama ada
		})
	}
	return results, nil
}
func GetKelasOnlineByMapel(idMapel string) ([]KelasOnlineListResponse, error) {
    var data []models.KelasOnline

    err := config.DB.
        Preload("Guru").
        Preload("KelasMapel.Mapel").
        Where("id_kelas_mapel = ?", idMapel).
        Find(&data).Error

    if err != nil {
        return nil, err
    }

    var res []KelasOnlineListResponse

    for _, k := range data {
        res = append(res, KelasOnlineListResponse{
            IDKelasOnline: k.IDKelasOnline,
            JudulKelas:    k.JudulKelas,
            NamaGuru:      k.Guru.GuruNama,
            NamaMapel:     k.KelasMapel.Mapel.NmMapel,
            TanggalKelas:  k.TanggalKelas.Format("2006-01-02"),
            JamMulai:      k.JamMulai,
            JamSelesai:    k.JamSelesai,
            Status:        k.Status,
            LinkKelas:     k.LinkKelas,
        })
    }

    return res, nil
}


func GetKelasOnlineByID(id string) (KelasOnlineDetailResponse, error) {
    var k models.KelasOnline

    err := config.DB.
        Preload("Guru").
        Preload("Materi").         // penting: preload materi
        First(&k, id).Error
    if err != nil {
        return KelasOnlineDetailResponse{}, err
    }

    // map materi
    materiList := make([]MateriResponse, 0, len(k.Materi))
    for _, m := range k.Materi {
        materiList = append(materiList, MateriResponse{
            ID:         m.IDKelasMateri,
            Judul:      m.Judul,
            Tipe:       m.Tipe,
            Link:       m.UrlFile,
            Keterangan: m.Keterangan,
            UploadedAt: m.UploadedAt.Format(time.RFC3339), // atau format lain "2006-01-02 15:04:05"
        })
    }

    return KelasOnlineDetailResponse{
        ID:      k.IDKelasOnline,
        Topik:   k.JudulKelas,
        Tanggal: k.TanggalKelas.Format(time.RFC3339), // atau "2006-01-02 15:04:05" sesuai kebutuhan
        Guru:    k.Guru.GuruNama,
        Materi:  materiList,
    }, nil
}



func CreateKelasOnline(data models.KelasOnline) (models.KelasOnline, error) {
    if err := config.DB.Create(&data).Error; err != nil {
        // print error detail di log server
        fmt.Printf("DB Create error: %T - %v\n", err, err)
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