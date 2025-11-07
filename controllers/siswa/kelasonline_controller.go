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
				// cek apakah TanggalKelas valid (bukan zero time)
				if !ko.TanggalKelas.IsZero() {
					dto.Tanggal = ko.TanggalKelas.Format("2006-01-02")
				} else {
					dto.Tanggal = "-" // kosong → tampil "-"
				}

				// jam mulai
				if ko.JamMulai != "" {
					dto.Mulai = ko.JamMulai
				} else {
					dto.Mulai = "-"
				}

				// jam selesai
				if ko.JamSelesai != "" {
					dto.Selesai = ko.JamSelesai
				} else {
					dto.Selesai = "-"
				}

				if ko.LinkKelas != "" {
					dto.Link = ko.LinkKelas
				} else {
					dto.Link = "-"
				}

				dto.Status = mapStatusSingkat(ko.Status)
			} else {
				// kalau gak ada kelas online
				dto.Tanggal = "-"
				dto.Mulai = "-"
				dto.Selesai = "-"
				dto.Link = "-"
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
// Struktur respons utama yang akan dikirim ke frontend
type MateriDetailResponse struct {
	Topik     string       `json:"topik"`
	Tanggal   string       `json:"tanggal"`
	Guru      string       `json:"guru"`
	Materi    []MateriItem `json:"materi"`
}

// Struktur untuk tiap item materi
type MateriItem struct {
	ID         int    `json:"id"`
	Judul      string `json:"judul"`
	Tipe       string `json:"tipe"`
	Link       string `json:"link"`
	Keterangan string `json:"keterangan"`
}

// ✅ Endpoint: ambil detail materi berdasarkan id_kelas_online
func GetMateriByKelasOnline(idKelasOnline string) (MateriDetailResponse, error) {
	// 1️⃣ Validasi parameter ID agar pasti angka
	id, err := strconv.ParseUint(idKelasOnline, 10, 32)
	if err != nil {
		return MateriDetailResponse{}, errors.New("invalid id_kelas_online")
	}

	// 2️⃣ Ambil data kelas online + relasi guru, mapel, dan materi
	var kelasOnline models.KelasOnline
	err = config.DB.
		Preload("Guru").
		Preload("KelasMapel.Mapel").
		Preload("Materi").
		First(&kelasOnline, "id_kelas_online = ?", uint(id)).Error
	if err != nil {
		return MateriDetailResponse{}, err
	}

	// 3️⃣ Buat list materi untuk response
	var materiList []MateriItem
	for i, materi := range kelasOnline.Materi {
		materiList = append(materiList, MateriItem{
			ID:         i + 1,
			Judul:      materi.Judul,
			Tipe:       materi.Tipe,
			Link:       materi.UrlFile,
			Keterangan: materi.Keterangan,
		})
	}
	// 5️⃣ Bentuk respons akhir
	response := MateriDetailResponse{
		Topik:     kelasOnline.JudulKelas,
		Tanggal:   kelasOnline.TanggalKelas.String(), // format biar frontend yang atur
		Guru:      kelasOnline.Guru.GuruNama,
		Materi:    materiList,
	}

	return response, nil
}
