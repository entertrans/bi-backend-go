package gurucontrollers

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"gorm.io/gorm"
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

// ========================
// Get Test Soal by Test ID
// ========================
// Fungsi untuk memformat jawaban
func formatJawaban(tipeSoal, pilihanJawaban, jawabanBenar string) (interface{}, interface{}) {
	var pilihan []interface{}
	var jawaban []interface{}

	// Parse JSON
	json.Unmarshal([]byte(pilihanJawaban), &pilihan)
	json.Unmarshal([]byte(jawabanBenar), &jawaban)

	// Format khusus untuk frontend
	if tipeSoal == "pg" || tipeSoal == "pg_kompleks" {
		// Untuk PG, ambil teks jawaban berdasarkan index
		formattedJawaban := make([]string, len(jawaban))
		for i, j := range jawaban {
			switch v := j.(type) {
			case float64:
				idx := int(v)
				if idx >= 0 && idx < len(pilihan) {
					// Ambil teks jawaban dari pilihan_jawaban
					if jawabanText, ok := pilihan[idx].(string); ok {
						formattedJawaban[i] = jawabanText
					} else {
						formattedJawaban[i] = "Invalid answer format"
					}
				} else {
					formattedJawaban[i] = strconv.Itoa(idx)
				}
			case string:
				// Jika sudah berupa teks jawaban, langsung kembalikan
				if len(v) > 1 { // Anggap sudah berupa teks, bukan huruf tunggal
					formattedJawaban[i] = v
				} else if idx, err := strconv.Atoi(v); err == nil && idx >= 0 && idx < len(pilihan) {
					// Jika angka, ambil teks dari pilihan_jawaban
					if jawabanText, ok := pilihan[idx].(string); ok {
						formattedJawaban[i] = jawabanText
					} else {
						formattedJawaban[i] = "Invalid answer format"
					}
				} else if len(v) == 1 && v >= "A" && v <= "Z" {
					// Jika huruf, konversi ke index lalu ambil teks
					idx := int(v[0] - 'A')
					if idx >= 0 && idx < len(pilihan) {
						if jawabanText, ok := pilihan[idx].(string); ok {
							formattedJawaban[i] = jawabanText
						} else {
							formattedJawaban[i] = "Invalid answer format"
						}
					} else {
						formattedJawaban[i] = v
					}
				} else {
					// Jika bukan format yang dikenali, kembalikan as-is
					formattedJawaban[i] = v
				}
			default:
				formattedJawaban[i] = "Invalid"
			}
		}
		return pilihan, formattedJawaban
	}

	return pilihan, jawaban
}

// Struct untuk response yang diformat
type FormattedTestSoal struct {
	TestsoalID     uint        `json:"testsoal_id"`
	TestID         uint        `json:"test_id"`
	TipeSoal       string      `json:"tipe_soal"`
	Pertanyaan     string      `json:"pertanyaan"`
	LampiranID     *uint       `json:"lampiran_id"`
	PilihanJawaban interface{} `json:"pilihan_jawaban"`
	JawabanBenar   interface{} `json:"jawaban_benar"`
	Bobot          float64     `json:"bobot"`
	CreatedAt      string      `json:"created_at"`
	DeletedAt      interface{} `json:"deleted_at"`
	Test           interface{} `json:"test"`
	Lampiran       interface{} `json:"lampiran"`
}

func GetTestSoalByTestId(testId uint) ([]FormattedTestSoal, error) {
	var soals []models.TO_TestSoal
	err := config.DB.
		Preload("Test").
		Preload("Lampiran").
		Where("test_id = ? AND deleted_at IS NULL", testId).
		Order("testsoal_id asc").
		Find(&soals).Error

	if err != nil {
		return nil, err
	}

	// Format data untuk frontend
	formattedSoals := make([]FormattedTestSoal, len(soals))
	for i, soal := range soals {
		pilihanFormatted, jawabanFormatted := formatJawaban(
			soal.TipeSoal,
			soal.PilihanJawaban,
			soal.JawabanBenar,
		)

		formattedSoals[i] = FormattedTestSoal{
			TestsoalID:     soal.TestsoalID,
			TestID:         soal.TestID,
			TipeSoal:       soal.TipeSoal,
			Pertanyaan:     soal.Pertanyaan,
			LampiranID:     soal.LampiranID,
			PilihanJawaban: pilihanFormatted,
			JawabanBenar:   jawabanFormatted,
			Bobot:          soal.Bobot,
			CreatedAt:      soal.CreatedAt.Format("2006-01-02 15:04:05"),
			DeletedAt:      soal.DeletedAt,
			Test:           soal.Test,
			Lampiran:       soal.Lampiran,
		}
	}

	return formattedSoals, nil
}

// ========================
// Get Detail Test Soal by ID
// ========================
func GetDetailTestSoal(soalId uint) (*models.TO_TestSoal, error) {
	var soal models.TO_TestSoal
	err := config.DB.
		Preload("Test").
		Preload("Lampiran").
		Where("testsoal_id = ? AND deleted_at IS NULL", soalId).
		First(&soal).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("soal tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}
	return &soal, nil
}

// ========================
// Delete Test Soal (Soft Delete)
// ========================
func DeleteTestSoal(soalId uint) error {
	var soal models.TO_TestSoal

	// Cek apakah soal exists
	result := config.DB.Where("testsoal_id = ? AND deleted_at IS NULL", soalId).First(&soal)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("soal tidak ditemukan")
	}

	// Soft delete
	result = config.DB.Delete(&soal)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// ========================
// Create or Update Test Soal
// ========================
type TestSoalInput struct {
	TestID         uint    `json:"test_id"`
	TipeSoal       string  `json:"tipe_soal"`
	Pertanyaan     string  `json:"pertanyaan"`
	LampiranID     *uint   `json:"lampiran_id"`
	PilihanJawaban string  `json:"pilihan_jawaban"` // JSON string
	JawabanBenar   string  `json:"jawaban_benar"`   // JSON string
	Bobot          float64 `json:"bobot"`
}

func CreateTestSoal(input TestSoalInput) (*models.TO_TestSoal, error) {
	// Validasi input
	if input.TestID == 0 {
		return nil, errors.New("test_id wajib diisi")
	}
	if input.TipeSoal == "" {
		return nil, errors.New("tipe_soal wajib diisi")
	}
	if input.Pertanyaan == "" {
		return nil, errors.New("pertanyaan wajib diisi")
	}

	// Validasi JSON fields
	if input.TipeSoal != "uraian" && input.TipeSoal != "isian_singkat" {
		if input.PilihanJawaban == "" {
			return nil, errors.New("pilihan_jawaban wajib diisi untuk tipe soal ini")
		}
		// Validasi JSON format
		var pj interface{}
		if err := json.Unmarshal([]byte(input.PilihanJawaban), &pj); err != nil {
			return nil, errors.New("pilihan_jawaban harus dalam format JSON yang valid")
		}
	}

	if input.JawabanBenar == "" {
		return nil, errors.New("jawaban_benar wajib diisi")
	}
	// Validasi JSON format
	var jb interface{}
	if err := json.Unmarshal([]byte(input.JawabanBenar), &jb); err != nil {
		return nil, errors.New("jawaban_benar harus dalam format JSON yang valid")
	}

	// Create soal
	soal := models.TO_TestSoal{
		TestID:         input.TestID,
		TipeSoal:       input.TipeSoal,
		Pertanyaan:     input.Pertanyaan,
		LampiranID:     input.LampiranID,
		PilihanJawaban: input.PilihanJawaban,
		JawabanBenar:   input.JawabanBenar,
		Bobot:          input.Bobot,
	}

	if err := config.DB.Create(&soal).Error; err != nil {
		return nil, err
	}

	// Reload dengan preload
	var createdSoal models.TO_TestSoal
	err := config.DB.
		Preload("Test").
		Preload("Lampiran").
		First(&createdSoal, soal.TestsoalID).Error

	if err != nil {
		return nil, err
	}

	return &createdSoal, nil
}
