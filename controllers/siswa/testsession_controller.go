package siswa

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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
		return nil, fmt.Errorf("NIS tidak valid: %w", err)
	}

	// Cari session terakhir untuk test + siswa ini
	var existingSession models.TestSession
	err = config.DB.
		Preload("Test").
		Where("test_id = ? AND siswa_nis = ?", testID, nisInt).
		Order("start_time DESC").
		First(&existingSession).Error

	// 🎯 JIKA SUDAH ADA SESSION YANG MASIH AKTIF
	if err == nil && existingSession.Status == "in_progress" {
		elapsedTime := time.Since(existingSession.StartTime)
		totalDuration := time.Duration(existingSession.Test.DurasiMenit) * time.Minute
		remainingTime := totalDuration - elapsedTime

		if remainingTime <= 0 {
			// Waktu habis, auto-submit
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

		// ✅ KEMBALIKAN SESSION YANG SUDAH ADA + CEK SOAL SUDAH DISIMPAN
		log.Printf("🟢 Menggunakan session yang sudah ada: ID=%d", existingSession.SessionID)
		return &existingSession, nil
	}

	// 🎯 JIKA SUDAH SUBMITTED
	if err == nil && existingSession.Status == "submitted" {
		return nil, fmt.Errorf("anda sudah menyelesaikan test ini")
	}

	// 🎯 JIKA BELUM ADA SESSION → BUAT BARU
	var test models.TO_Test
	if err := config.DB.First(&test, testID).Error; err != nil {
		return nil, fmt.Errorf("test tidak ditemukan: %w", err)
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
		return nil, fmt.Errorf("gagal membuat session: %w", err)
	}

	// 📦 AMBIL SOAL MENGGUNAKAN GetSoalByTestID1
	soalList, err := GetSoalByTestID1(testID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil soal: %w", err)
	}

	// 🎲 ACAK SOAL JIKA DIBUTUHKAN
	if test.RandomSoal {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(soalList), func(i, j int) {
			soalList[i], soalList[j] = soalList[j], soalList[i]
		})
	}

	// 💾 SIMPAN URUTAN SOAL FIX KE to_sessionsoal
	for i, soal := range soalList {
		sessionSoal := models.TO_SessionSoal{
			SessionID: session.SessionID,
			SoalID:    soal.SoalID,
			OrderNo:   i + 1,
		}
		if err := config.DB.Create(&sessionSoal).Error; err != nil {
			log.Printf("⚠️ Gagal menyimpan session soal %d: %v", i+1, err)
		}
	}

	log.Printf("✅ Session baru dibuat: ID=%d, %d soal disimpan ke to_sessionsoal", session.SessionID, len(soalList))

	return &session, nil
}

// Function untuk mengambil soal berdasarkan session_id (setelah disimpan)

// bikin DTO umum biar frontend nggak ribet
type SoalDTO struct {
	SoalID           uint        `json:"soal_id"`
	TipeSoal         string      `json:"tipe_soal"`
	Pertanyaan       string      `json:"pertanyaan"`
	LampiranID       *uint       `json:"lampiran_id"`
	PilihanJawaban   string      `json:"pilihan_jawaban"`
	JawabanBenar     string      `json:"jawaban_benar"`
	Bobot            float64     `json:"bobot"`
	JawabanTersimpan interface{} `json:"jawaban_tersimpan,omitempty"`
}

// GetSoalBySessionID harus return []SoalDTO
func GetSoalBySessionID(sessionID uint) ([]SoalDTO, error) {
	// 1. ambil session
	var session models.TestSession
	if err := config.DB.First(&session, sessionID).Error; err != nil {
		return nil, fmt.Errorf("session tidak ditemukan")
	}

	// 2. ambil test → cek type_test
	var test models.TO_Test
	if err := config.DB.First(&test, session.TestID).Error; err != nil {
		return nil, fmt.Errorf("test tidak ditemukan")
	}

	// 3. ambil urutan soal
	var sessionSoals []models.TO_SessionSoal
	if err := config.DB.Where("session_id = ?", sessionID).Order("order_no ASC").Find(&sessionSoals).Error; err != nil {
		return nil, fmt.Errorf("gagal ambil urutan soal: %w", err)
	}

	// 4. ambil jawaban siswa (final)
	var jawabanFinal []models.JawabanFinal
	if err := config.DB.Where("session_id = ?", sessionID).Find(&jawabanFinal).Error; err != nil {
		log.Printf("⚠️ gagal ambil jawaban siswa: %v", err)
	}

	// bikin map untuk akses cepat
	jawabanMap := make(map[uint]interface{})
	for _, jf := range jawabanFinal {
		var parsed interface{}
		if err := json.Unmarshal([]byte(jf.JawabanSiswa.String()), &parsed); err == nil {
			jawabanMap[jf.SoalID] = parsed
		} else {
			jawabanMap[jf.SoalID] = jf.JawabanSiswa.String()
		}
	}

	var soalDTOs []SoalDTO

	// 5. ambil soal sesuai tipe test
	if test.TypeTest == "ub" {
		var soalIDs []uint
		for _, ss := range sessionSoals {
			soalIDs = append(soalIDs, ss.SoalID)
		}

		var soals []models.TO_BankSoal
		if err := config.DB.Where("soal_id IN (?)", soalIDs).Preload("Lampiran").Find(&soals).Error; err != nil {
			return nil, fmt.Errorf("gagal ambil soal bank: %w", err)
		}

		for _, s := range soals {
			soalDTOs = append(soalDTOs, SoalDTO{
				SoalID:           s.SoalID,
				TipeSoal:         s.TipeSoal,
				Pertanyaan:       s.Pertanyaan,
				LampiranID:       s.LampiranID,
				PilihanJawaban:   s.PilihanJawaban,
				JawabanBenar:     s.JawabanBenar,
				Bobot:            s.Bobot,
				JawabanTersimpan: jawabanMap[s.SoalID], // <<<< diset di sini
			})
		}

	} else {
		var soalIDs []uint
		for _, ss := range sessionSoals {
			soalIDs = append(soalIDs, ss.SoalID)
		}

		var soals []models.TO_TestSoal
		if err := config.DB.Where("testsoal_id IN (?)", soalIDs).Preload("Lampiran").Find(&soals).Error; err != nil {
			return nil, fmt.Errorf("gagal ambil soal test: %w", err)
		}

		for _, s := range soals {
			soalDTOs = append(soalDTOs, SoalDTO{
				SoalID:           s.TestsoalID,
				TipeSoal:         s.TipeSoal,
				Pertanyaan:       s.Pertanyaan,
				LampiranID:       s.LampiranID,
				PilihanJawaban:   s.PilihanJawaban,
				JawabanBenar:     s.JawabanBenar,
				Bobot:            s.Bobot,
				JawabanTersimpan: jawabanMap[s.TestsoalID], // <<<< diset di sini juga
			})
		}
	}

	return soalDTOs, nil
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

	} else if test.TypeTest == "tr" || test.TypeTest == "tugas" {
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
