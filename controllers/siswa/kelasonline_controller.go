package siswa

import (
	"errors"
	"strconv"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

// DTO hasil query (data yang dikirim ke frontend)
type KelasOnlineList struct {
	IDKelasOnline uint   `json:"id_kelas_online"`
	JudulKelas    string `json:"judul_kelas"`
	NamaGuru      string `json:"nama_guru"`
	NamaMapel     string `json:"nama_mapel"`
	Tanggal       string `json:"tanggal"`
	JamMulai      string `json:"jam_mulai"`
	JamSelesai    string `json:"jam_selesai"`
	Status        string `json:"status"`
	LinkKelas     string `json:"link_kelas"`
	MateriLink    string `json:"materi_link"`
}

// 🔹 Endpoint 1: daftar kelas online berdasarkan kelas_id
func GetKelasOnlineByKelasID(kelasID string) ([]KelasOnlineList, error) {
	var kelasOnlines []models.KelasOnline
	err := config.DB.
		Preload("KelasMapel.Kelas").
		Preload("KelasMapel.Mapel").
		Preload("Guru").
		Joins("JOIN tbl_kelas_mapel km ON km.id_kelas_mapel = tbl_kelas_online.id_kelas_mapel").
		Where("km.kelas_id = ?", kelasID).
		Find(&kelasOnlines).Error
	if err != nil {
		return nil, err
	}

	var result []KelasOnlineList
	for _, k := range kelasOnlines {
		result = append(result, KelasOnlineList{
			IDKelasOnline: k.IDKelasOnline,
			JudulKelas:    k.JudulKelas,
			NamaGuru:      k.Guru.GuruNama,
			NamaMapel:     k.KelasMapel.Mapel.NmMapel,
			Tanggal:       k.TanggalKelas.Format("2006-01-02"),
			JamMulai:      k.JamMulai,
			JamSelesai:    k.JamSelesai,
			Status:        k.Status,
			LinkKelas:     k.LinkKelas,
			MateriLink:    k.MateriLink,
		})
	}
	return result, nil
}

// 🔹 Endpoint 2: daftar riwayat kelas berdasarkan id_kelas_mapel
func GetKelasOnlineHistory(idKelasMapel string) ([]KelasOnlineList, error) {
	var kelasOnlines []models.KelasOnline
	err := config.DB.
		Preload("Guru").
		Preload("KelasMapel.Mapel").
		Where("id_kelas_mapel = ?", idKelasMapel).
		Order("tanggal_kelas DESC").
		Find(&kelasOnlines).Error
	if err != nil {
		return nil, err
	}

	var result []KelasOnlineList
	for _, k := range kelasOnlines {
		result = append(result, KelasOnlineList{
			IDKelasOnline: k.IDKelasOnline,
			JudulKelas:    k.JudulKelas,
			NamaGuru:      k.Guru.GuruNama,
			NamaMapel:     k.KelasMapel.Mapel.NmMapel,
			Tanggal:       k.TanggalKelas.Format("2006-01-02"),
			JamMulai:      k.JamMulai,
			JamSelesai:    k.JamSelesai,
			Status:        k.Status,
			LinkKelas:     k.LinkKelas,
			MateriLink:    k.MateriLink,
		})
	}
	return result, nil
}

type MapelKelasDTO struct {
    ID        uint   `json:"id"`
    Mapel     string `json:"mapel"`
    KodeMapel string `json:"kode_mapel"`
    Guru      string `json:"guru"`
    Tanggal   string `json:"tanggal"` // format YYYY-MM-DD atau empty
    Mulai     string `json:"mulai"`
    Selesai   string `json:"selesai"`
    Status    string `json:"status"`
    Link      string `json:"link"`
}

// GetMapelByKelasID mengembalikan semua mapel untuk kelas tertentu
// dan melampirkan data kelas_online terbaru jika ada.
func GetMapelByKelasID(kelasIDStr string) ([]MapelKelasDTO, error) {
    // validasi param
    kelasID64, err := strconv.ParseUint(kelasIDStr, 10, 64)
    if err != nil {
        return nil, errors.New("invalid kelas_id")
    }
    kelasID := uint(kelasID64)

    // ambil semua kelas_mapel untuk kelas tersebut + preload Mapel
    var kelasMapels []models.KelasMapel
    if err := config.DB.Preload("Mapel").
        Where("kelas_id = ?", kelasID).
        Find(&kelasMapels).Error; err != nil {
        return nil, err
    }

    var results []MapelKelasDTO
    for _, km := range kelasMapels {
        dto := MapelKelasDTO{
            ID:        km.ID,
            Mapel:     "", // nanti diisi
            KodeMapel: strconv.FormatUint(uint64(km.KdMapel), 10),
            Guru:      "",
            Tanggal:   "",
            Mulai:     "",
            Selesai:   "",
            Status:    "belum", // default
            Link:      "",
        }

        // Mapel (jika preload ada)
        if km.Mapel.KdMapel != 0 || km.Mapel.NmMapel != "" {
            dto.Mapel = km.Mapel.NmMapel
        }

        // Ambil guru untuk pasangan (kelas_id, kd_mapel) dari tbl_guru_mapel (ambil 1st aktif)
        var gm models.GuruMapel
        if err := config.DB.Preload("Guru").
            Where("kelas_id = ? AND kd_mapel = ? AND status_aktif = ?", km.KelasID, km.KdMapel, true).
            Order("guru_mapel_id DESC").
            Limit(1).
            Find(&gm).Error; err == nil {
            if gm.Guru.GuruNama != "" {
                dto.Guru = gm.Guru.GuruNama
            }
        }

        // Ambil session kelas_online terbaru untuk id_kelas_mapel = km.ID
        var ko models.KelasOnline
        if err := config.DB.Where("id_kelas_mapel = ?", km.ID).
            Order("tanggal_kelas DESC, jam_mulai DESC, created_at DESC").
            Limit(1).
            Find(&ko).Error; err == nil {
            // jika ditemukan (cek IDKelasOnline != 0)
            if ko.IDKelasOnline != 0 {
                // tanggal format YYYY-MM-DD
                dto.Tanggal = ko.TanggalKelas.Format("2006-01-02")
                dto.Mulai = ko.JamMulai
                dto.Selesai = ko.JamSelesai
                dto.Link = ko.LinkKelas
                dto.Status = mapStatusSingkat(ko.Status)
            }
        }

        results = append(results, dto)
    }

    return results, nil
}

func mapStatusSingkat(statusFull string) string {
    switch statusFull {
    case "akan_berlangsung":
        return "belum"
    case "sedang_berlangsung":
        return "sedang"
    case "selesai":
        return "selesai"
    default:
        // fallback: jika empty atau lain -> "belum"
        return "belum"
    }
}