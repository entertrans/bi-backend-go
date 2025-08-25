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
		// jawaban_benar: [string]
		var jb []string
		if err := json.Unmarshal([]byte(input.JawabanBenar), &jb); err != nil {
			return errors.New("jawaban_benar untuk BS harus array string")
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

	// Ambil semua soal aktif
	var soals []models.TO_BankSoal
	err := config.DB.
		Preload("Kelas").
		Preload("Mapel").
		Where("deleted_at IS NULL").
		Find(&soals).Error
	if err != nil {
		return nil, err
	}

	// Peta buat grouping
	rekapMap := make(map[string]*RekapSoal)

	for _, soal := range soals {
		key := fmt.Sprintf("%d-%d", soal.KelasID, soal.MapelID)

		if _, ok := rekapMap[key]; !ok {
			rekapMap[key] = &RekapSoal{
				KelasID:   soal.KelasID,
				KelasNama: soal.Kelas.KelasNama,
				MapelID:   soal.MapelID,
				MapelNama: soal.Mapel.NmMapel, // ambil nama mapel dari relasi
			}

		}

		switch soal.TipeSoal {
		case "pg":
			rekapMap[key].PG++
		case "pg_kompleks":
			rekapMap[key].PGKompleks++
		case "isian_singkat":
			rekapMap[key].Isian++
		case "true_false":
			rekapMap[key].TrueFalse++
		case "uraian":
			rekapMap[key].Uraian++
		case "matching":
			rekapMap[key].Mencocokkan++
		}

		rekapMap[key].Total++
	}

	// Ubah map ke slice
	for _, v := range rekapMap {
		results = append(results, *v)
	}

	return results, nil
}
