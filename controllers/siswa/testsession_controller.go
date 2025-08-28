package siswa

import (
	"fmt"
	"strconv"
	"time"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

// Mulai test baru atau lanjutkan session yang ada - DIPERBAIKI
func StartTestSession(testID uint, nis string) (*models.TestSession, error) {
	// Convert nis string to int
	nisInt, err := strconv.Atoi(nis)
	if err != nil {
		return nil, err
	}

	// Cek apakah sudah ada session yang aktif untuk test ini
	var existingSession models.TestSession
	err = config.DB.
		Preload("Test").
		Where("test_id = ? AND siswa_nis = ? AND status = 'in_progress'", testID, nisInt).
		First(&existingSession).Error

	// Jika sudah ada session aktif, hitung waktu sisa dan return session
	if err == nil {
		// Hitung waktu sisa
		elapsedTime := time.Since(existingSession.StartTime)
		totalDuration := time.Duration(existingSession.Test.DurasiMenit) * time.Minute
		remainingTime := totalDuration - elapsedTime

		// Jika waktu sudah habis, update status menjadi submitted
		if remainingTime <= 0 {
			endTime := existingSession.StartTime.Add(totalDuration)
			existingSession.Status = "submitted"
			existingSession.EndTime = &endTime
			existingSession.WaktuSisa = 0
			config.DB.Save(&existingSession)
			return nil, fmt.Errorf("waktu ujian sudah habis")
		}

		// Update waktu sisa di database
		existingSession.WaktuSisa = int(remainingTime.Seconds())
		config.DB.Save(&existingSession)

		return &existingSession, nil
	}

	// Jika tidak ada session aktif, buat session baru
	var test models.TO_Test
	if err := config.DB.First(&test, testID).Error; err != nil {
		return nil, err
	}

	startTime := time.Now()
	endTime := startTime.Add(time.Duration(test.DurasiMenit) * time.Minute)

	session := models.TestSession{
		TestID:     testID,
		SiswaNIS:   nisInt,
		StartTime:  startTime,
		EndTime:    &endTime,
		WaktuSisa:  test.DurasiMenit * 60,
		Status:     "in_progress",
		NilaiAwal:  0,
		NilaiAkhir: 0,
	}

	err = config.DB.Create(&session).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// Get session aktif - DIPERBAIKI
func GetActiveTestSession(testID uint, nis string) (*models.TestSession, error) {
	nisInt, err := strconv.Atoi(nis)
	if err != nil {
		return nil, err
	}

	var session models.TestSession
	err = config.DB.
		Preload("Test").
		Where("test_id = ? AND siswa_nis = ? AND status = 'in_progress'", testID, nisInt).
		First(&session).Error

	if err != nil {
		return nil, err
	}

	// Hitung waktu sisa real-time
	elapsedTime := time.Since(session.StartTime)
	totalDuration := time.Duration(session.Test.DurasiMenit) * time.Minute
	remainingTime := totalDuration - elapsedTime

	// Jika waktu habis, update status
	if remainingTime <= 0 {
		endTime := session.StartTime.Add(totalDuration)
		session.Status = "submitted"
		session.EndTime = &endTime
		session.WaktuSisa = 0
		config.DB.Save(&session)
		return nil, fmt.Errorf("waktu ujian sudah habis")
	}

	// Update waktu sisa untuk response (tanpa menyimpan ke database)
	session.WaktuSisa = int(remainingTime.Seconds())

	return &session, nil
}

// Ambil sesi siswa untuk 1 test
func GetTestSession(testID uint, nis string) (*models.TestSession, error) {
	var session models.TestSession
	err := config.DB.
		Preload("JawabanFinal").
		Where("test_id = ? AND siswa_nis = ?", testID, nis).
		First(&session).Error

	if err != nil {
		return nil, err
	}
	return &session, nil
}

// Submit test (update end_time & status jadi submitted)
func SubmitTestSession(sessionID uint) error {
	return config.DB.Model(&models.TestSession{}).
		Where("session_id = ?", sessionID).
		Updates(map[string]interface{}{
			"end_time": time.Now(),
			"status":   "submitted",
		}).Error
}

// GET /siswa/test/:test_id/soal
func GetSoalByTestID1(testID uint) ([]models.TO_BankSoal, error) {
	// Ambil data test
	var test models.TO_Test
	if err := config.DB.First(&test, testID).Error; err != nil {
		return nil, fmt.Errorf("Test tidak ditemukan")
	}

	var soal []models.TO_BankSoal

	if test.TypeTest == "ub" {
		// UB: ambil dari bank soal
		query := config.DB.Where("kelas_id = ? AND mapel_id = ?", test.KelasID, test.MapelID)

		if test.RandomSoal {
			query = query.Order("RAND()")
		}

		if test.Jumlah > 0 {
			query = query.Limit(int(test.Jumlah))
		}

		if err := query.Find(&soal).Error; err != nil {
			return nil, err
		}

	} else if test.TypeTest == "quis" {
		// QUIZ: ambil dari tabel to_testsoal
		var testSoal []models.TO_TestSoal
		if err := config.DB.Where("test_id = ?", test.TestID).Find(&testSoal).Error; err != nil {
			return nil, err
		}

		// Konversi TO_TestSoal → TO_BankSoal-like biar frontend gampang render
		for _, ts := range testSoal {
			soal = append(soal, models.TO_BankSoal{
				SoalID:         ts.TestsoalID,
				TipeSoal:       ts.TipeSoal,
				Pertanyaan:     ts.Pertanyaan,
				LampiranID:     ts.LampiranID,
				PilihanJawaban: ts.PilihanJawaban,
				JawabanBenar:   ts.JawabanBenar,
				Bobot:          ts.Bobot,
			})
		}
	}

	return soal, nil
}
