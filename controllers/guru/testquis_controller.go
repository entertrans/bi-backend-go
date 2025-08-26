package gurucontrollers

import (
	"encoding/json"
	"errors"

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
func GetTestSoalByTestId(testId uint) ([]models.TO_TestSoal, error) {
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
	return soals, nil
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
	TestID          uint    `json:"test_id"`
	TipeSoal        string  `json:"tipe_soal"`
	Pertanyaan      string  `json:"pertanyaan"`
	LampiranID      *uint   `json:"lampiran_id"`
	PilihanJawaban  string  `json:"pilihan_jawaban"` // JSON string
	JawabanBenar    string  `json:"jawaban_benar"`   // JSON string
	Bobot           float64 `json:"bobot"`
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