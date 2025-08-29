package gurucontrollers

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

type BankSoalInput struct {
	GuruID         uint    `json:"guru_id"`
	MapelID        uint    `json:"mapel_id"`
	TipeSoal       string  `json:"tipe_soal"`
	KelasID        uint    `json:"kelas_id"`
	Pertanyaan     string  `json:"pertanyaan"`
	LampiranID     *uint   `json:"lampiran_id"`
	PilihanJawaban string  `json:"pilihan_jawaban"` // JSON string
	JawabanBenar   string  `json:"jawaban_benar"`   // JSON string (berisi index)
	Bobot          float64 `json:"bobot"`
}

// ========================
// Ambil soal aktif (belum dihapus)
// ========================
func GetActiveBankSoal() ([]models.TO_BankSoal, error) {
	var soals []models.TO_BankSoal
	err := config.DB.
		Preload("Guru").
		Preload("Kelas").
		Preload("Mapel").
		Preload("Lampiran"). // 🔥 ikut load lampiran kalau ada
		Where("deleted_at IS NULL").
		Order("created_at desc").
		Find(&soals).Error
	return soals, err
}

// ========================
// Ambil soal non-aktif (soft deleted)
// ========================
func GetInactiveBankSoal() ([]models.TO_BankSoal, error) {
	var soals []models.TO_BankSoal
	err := config.DB.
		Unscoped().
		Preload("Guru").
		Preload("Kelas").
		Preload("Mapel").
		Preload("Lampiran").
		Where("deleted_at IS NOT NULL").
		Order("created_at desc").
		Find(&soals).Error
	return soals, err
}

// ========================
// Restore soal (hapus deleted_at)
// ========================
func RestoreBankSoal(soalID uint) (*models.TO_BankSoal, error) {
	var soal models.TO_BankSoal

	// Ambil data termasuk yang soft-deleted
	if err := config.DB.Unscoped().
		Where("soal_id = ?", soalID).
		First(&soal).Error; err != nil {
		return nil, err
	}

	// Update deleted_at supaya NULL → restore
	if err := config.DB.Unscoped().
		Model(&soal).
		Update("deleted_at", nil).Error; err != nil {
		return nil, err
	}

	// Refresh data setelah diupdate
	if err := config.DB.First(&soal, soalID).Error; err != nil {
		return nil, err
	}

	return &soal, nil
}

// ========================
// Ambil soal milik guru tertentu
// ========================
func GetBankSoalByGuru(guruID uint) ([]models.TO_BankSoal, error) {
	var soals []models.TO_BankSoal
	err := config.DB.
		Preload("Guru").
		Preload("Kelas").
		Preload("Mapel").
		Preload("Lampiran").
		Where("guru_id = ? AND deleted_at IS NULL", guruID). // 🔧 ganti is_deleted jadi deleted_at IS NULL
		Order("created_at desc").
		Find(&soals).Error
	return soals, err
}

// ========================
// Soft delete soal
// ========================
func DeleteBankSoal(soalID uint) error {
	result := config.DB.Delete(&models.TO_BankSoal{}, soalID)
	if result.RowsAffected == 0 {
		return errors.New("bank soal tidak ditemukan")
	}
	return result.Error
}

// ========================
// Membuat soal baru
// ========================
func SimpanSoal(input BankSoalInput) error {
	if input.GuruID == 0 {
		return errors.New("guru_id wajib diisi")
	}
	if input.KelasID == 0 {
		return errors.New("kelas_id wajib diisi")
	}
	if input.Pertanyaan == "" {
		return errors.New("pertanyaan wajib diisi")
	}

	// Validasi pilihan_jawaban & jawaban_benar sesuai tipe_soal
	switch input.TipeSoal {
	case "pg":
		// pilihan_jawaban: []string, jawaban_benar: [string]
		var pj []string
		if err := json.Unmarshal([]byte(input.PilihanJawaban), &pj); err != nil {
			return errors.New("pilihan_jawaban untuk PG harus array string")
		}
		var jb []string
		if err := json.Unmarshal([]byte(input.JawabanBenar), &jb); err != nil {
			return errors.New("jawaban_benar untuk PG harus array string")
		}

	case "pg_kompleks":
		// pilihan_jawaban: []string, jawaban_benar: []string
		var pj []string
		if err := json.Unmarshal([]byte(input.PilihanJawaban), &pj); err != nil {
			return errors.New("pilihan_jawaban untuk PG Kompleks harus array string")
		}
		var jb []string
		if err := json.Unmarshal([]byte(input.JawabanBenar), &jb); err != nil {
			return errors.New("jawaban_benar untuk PG Kompleks harus array string")
		}

	case "matching":
		// pilihan_jawaban: array object {left,right,leftLampiran,rightLampiran}
		var pj []map[string]interface{}
		if err := json.Unmarshal([]byte(input.PilihanJawaban), &pj); err != nil {
			return errors.New("pilihan_jawaban untuk Matching harus array object")
		}
		// jawaban_benar biasanya kosong → cukup validasi bisa unmarshal ke array
		var jb []interface{}
		if err := json.Unmarshal([]byte(input.JawabanBenar), &jb); err != nil {
			return errors.New("jawaban_benar untuk Matching harus array")
		}

	case "bs":
		// pilihan_jawaban: array pernyataan { "teks": string, "jawaban": "benar"/"salah" }
		var pj []map[string]interface{}
		if err := json.Unmarshal([]byte(input.PilihanJawaban), &pj); err != nil {
			return errors.New("pilihan_jawaban untuk Benar/Salah harus array object")
		}

		// Contoh: jawaban_benar = ["benar", "salah", "benar"]
		var jb []string
		if err := json.Unmarshal([]byte(input.JawabanBenar), &jb); err != nil {
			return errors.New("jawaban_benar untuk Benar/Salah harus array string")
		}

		// validasi panjang array sama
		if len(jb) != len(pj) {
			return errors.New("jumlah jawaban_benar tidak sesuai dengan jumlah pernyataan")
		}

		// validasi isi benar/salah
		for _, j := range jb {
			if j != "benar" && j != "salah" {
				return errors.New("jawaban_benar hanya boleh 'benar' atau 'salah'")
			}
		}

	case "uraian", "isian_singkat":
		// jawaban_benar: [string] (biasanya satu)
		var jb []string
		if err := json.Unmarshal([]byte(input.JawabanBenar), &jb); err != nil {
			return errors.New("jawaban_benar untuk Uraian/Isian harus array string")
		}

	default:
		return errors.New("tipe_soal tidak dikenali")
	}

	// Simpan ke DB
	soal := models.TO_BankSoal{
		GuruID:         input.GuruID,
		MapelID:        input.MapelID,
		TipeSoal:       input.TipeSoal,
		KelasID:        input.KelasID,
		Pertanyaan:     input.Pertanyaan,
		LampiranID:     input.LampiranID,
		PilihanJawaban: input.PilihanJawaban,
		JawabanBenar:   input.JawabanBenar,
		Bobot:          input.Bobot,
		CreatedAt:      time.Now(),
	}

	if err := config.DB.Create(&soal).Error; err != nil {
		return err
	}

	return nil
}

// ========================
// Ambil soal aktif by kelas & mapel
// ========================
func GetActiveBankSoalByKelasMapel(kelasID uint, mapelID uint) ([]models.TO_BankSoal, error) {
	var soals []models.TO_BankSoal
	err := config.DB.
		Preload("Guru").
		Preload("Kelas").
		Preload("Mapel").
		Preload("Lampiran").
		Preload("Mapel").
		Where("deleted_at IS NULL AND kelas_id = ? AND mapel_id = ?", kelasID, mapelID).
		Order("created_at desc").
		Find(&soals).Error
	return soals, err
}

// ========================
// Rekap soal aktif per kelas & mapel
// ========================
type RekapSoal struct {
	KelasID     uint   `json:"kelas_id"`
	KelasNama   string `json:"kelas_nama"`
	MapelID     uint   `json:"mapel_id"`
	MapelNama   string `json:"mapel_nama"`
	PG          int    `json:"pg"`
	PGKompleks  int    `json:"pg_kompleks"`
	Isian       int    `json:"isian_singkat"`
	TrueFalse   int    `json:"true_false"`
	Uraian      int    `json:"uraian"`
	Mencocokkan int    `json:"mencocokkan"`
	Total       int    `json:"total"`
}

func GetRekapBankSoal() ([]RekapSoal, error) {
	var results []RekapSoal

	// 1. Ambil semua kelas-mapel
	var kelasMapels []models.KelasMapel
	err := config.DB.
		Preload("Kelas").
		Preload("Mapel").
		Order("kelas_id ASC, kd_mapel ASC").
		Find(&kelasMapels).Error
	if err != nil {
		return nil, err
	}

	// 2. Ambil semua soal aktif
	var soals []models.TO_BankSoal
	err = config.DB.
		Where("deleted_at IS NULL").
		Find(&soals).Error
	if err != nil {
		return nil, err
	}

	// 3. Peta count soal
	countMap := make(map[string]*RekapSoal)

	for _, soal := range soals {
		key := fmt.Sprintf("%d-%d", soal.KelasID, soal.MapelID)

		if _, ok := countMap[key]; !ok {
			countMap[key] = &RekapSoal{
				KelasID:     soal.KelasID,
				MapelID:     soal.MapelID,
				PG:          0,
				PGKompleks:  0,
				Isian:       0,
				TrueFalse:   0,
				Uraian:      0,
				Mencocokkan: 0,
				Total:       0,
			}
		}

		switch soal.TipeSoal {
		case "pg":
			countMap[key].PG++
		case "pg_kompleks":
			countMap[key].PGKompleks++
		case "isian_singkat":
			countMap[key].Isian++
		case "bs":
			countMap[key].TrueFalse++
		case "uraian":
			countMap[key].Uraian++
		case "matching":
			countMap[key].Mencocokkan++
		}

		countMap[key].Total++
	}

	// 4. Loop semua kelas-mapel, gabung dengan countMap
	for _, km := range kelasMapels {
		key := fmt.Sprintf("%d-%d", km.KelasID, km.KdMapel)

		rekap := RekapSoal{
			KelasID:   km.KelasID,
			KelasNama: km.Kelas.KelasNama,
			MapelID:   km.KdMapel,
			MapelNama: km.Mapel.NmMapel,
		}

		if val, ok := countMap[key]; ok {
			rekap.PG = val.PG
			rekap.PGKompleks = val.PGKompleks
			rekap.Isian = val.Isian
			rekap.TrueFalse = val.TrueFalse
			rekap.Uraian = val.Uraian
			rekap.Mencocokkan = val.Mencocokkan
			rekap.Total = val.Total
		}

		results = append(results, rekap)
	}

	return results, nil
}
