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
func RestoreBankSoal(soalID uint) error {
	result := config.DB.Model(&models.TO_BankSoal{}).
		Unscoped().
		Where("soal_id = ?", soalID).
		Update("deleted_at", nil)
	if result.RowsAffected == 0 {
		return errors.New("bank soal tidak ditemukan")
	}
	return result.Error
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
	// Validasi dasar
	if input.GuruID == 0 {
		return errors.New("guru_id wajib diisi")
	}
	if input.KelasID == 0 {
		return errors.New("kelas_id wajib diisi")
	}
	if input.Pertanyaan == "" {
		return errors.New("pertanyaan wajib diisi")
	}
	if input.PilihanJawaban == "" {
		return errors.New("pilihan_jawaban wajib diisi")
	}
	if input.JawabanBenar == "" {
		return errors.New("jawaban_benar wajib diisi")
	}

	// Validasi format JSON pilihan_jawaban
	var tmp []string
	if err := json.Unmarshal([]byte(input.PilihanJawaban), &tmp); err != nil {
		return errors.New("format pilihan_jawaban tidak valid (harus array string)")
	}
	if len(tmp) == 0 {
		return errors.New("pilihan_jawaban tidak boleh kosong")
	}

	// Validasi format JSON jawaban_benar (harus array int)
	var ans []int
	if err := json.Unmarshal([]byte(input.JawabanBenar), &ans); err != nil {
		return errors.New("format jawaban_benar tidak valid (harus array int)")
	}
	if len(ans) == 0 {
		return errors.New("jawaban_benar tidak boleh kosong")
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
