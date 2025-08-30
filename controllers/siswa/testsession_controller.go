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

	// Cari session terakhir untuk test + siswa ini (submitted atau in_progress)
	var existingSession models.TestSession
	err = config.DB.
		Preload("Test").
		Where("test_id = ? AND siswa_nis = ?", testID, nisInt).
		Order("start_time DESC").
		First(&existingSession).Error

	if err == nil {
		// 🚨 Kalau sudah submitted, jangan bikin baru
		if existingSession.Status == "submitted" {
			return nil, fmt.Errorf("waktu ujian sudah habis")
		}

		// Kalau masih in_progress → cek waktu
		if existingSession.Status == "in_progress" {
			elapsedTime := time.Since(existingSession.StartTime)
			totalDuration := time.Duration(existingSession.Test.DurasiMenit) * time.Minute
			remainingTime := totalDuration - elapsedTime

			if remainingTime <= 0 {
				endTime := existingSession.StartTime.Add(totalDuration)
				existingSession.Status = "submitted"
				existingSession.EndTime = &endTime
				existingSession.WaktuSisa = 0
				config.DB.Save(&existingSession)
				return nil, fmt.Errorf("waktu ujian sudah habis")
			}

			// Update waktu sisa
			existingSession.WaktuSisa = int(remainingTime.Seconds())
			config.DB.Save(&existingSession)

			return &existingSession, nil
		}
	}

	// Kalau belum ada session sama sekali → buat baru
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
		return nil, fmt.Errorf("NIS tidak valid: %w", err)
	}

	var session models.TestSession
	err = config.DB.
		Preload("Test").
		Where("test_id = ? AND siswa_nis = ?", testID, nisInt).
		Order("updated_at DESC").
		First(&session).Error

	if err != nil {
		return nil, err // Return error asli untuk dibedakan di handler
	}

	// ✅ CEK JIKA Test TIDAK NULL
	if session.Test == nil {
		return nil, fmt.Errorf("data test tidak ditemukan untuk session")
	}

	// ✅ HANYA PROSES JIKA MASIH IN_PROGRESS
	if session.Status == "in_progress" {
		elapsedTime := time.Since(session.StartTime)
		totalDuration := time.Duration(session.Test.DurasiMenit) * time.Minute
		remainingTime := totalDuration - elapsedTime

		if remainingTime <= 0 {
			// Waktu habis, auto-submit
			endTime := session.StartTime.Add(totalDuration)
			session.Status = "submitted"
			session.EndTime = &endTime
			session.WaktuSisa = 0

			if err := config.DB.Save(&session).Error; err != nil {
				return nil, fmt.Errorf("gagal menyimpan session: %w", err)
			}
		} else {
			// Update field sementara untuk response (tidak disimpan ke DB)
			session.WaktuSisa = int(remainingTime.Seconds())
		}
	}

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

		if err := query.Preload("Guru").
			Preload("Kelas").
			Preload("Mapel").
			Preload("Lampiran").
			Find(&soal).Error; err != nil {
			return nil, err
		}

	} else if test.TypeTest == "tr" {
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
