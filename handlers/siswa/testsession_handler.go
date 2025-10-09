package siswa

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/entertrans/bi-backend-go/config"
	siswaControllers "github.com/entertrans/bi-backend-go/controllers/siswa"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

// POST /siswa/test/start/:test_id
func StartTestHandler(c *gin.Context) {
	testID, err := strconv.Atoi(c.Param("test_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Test ID tidak valid"})
		return
	}

	nis := c.Query("nis")
	if nis == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIS diperlukan"})
		return
	}

	kelasIDStr := c.Query("kelas_id")
	var kelasID uint
	if kelasIDStr != "" {
		parsed, err := strconv.Atoi(kelasIDStr)
		if err == nil {
			kelasID = uint(parsed)
		}
	}

	session, err := siswaControllers.StartTestSession(uint(testID), nis, kelasID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai test: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

// check sessi
func GetActiveTestSessionHandler(c *gin.Context) {
	testID, err := strconv.Atoi(c.Param("test_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Test ID tidak valid"})
		return
	}

	nis := c.Query("nis")
	if nis == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIS diperlukan"})
		return
	}

	session, err := siswaControllers.GetActiveTestSession(uint(testID), nis)

	// ✅ PENANGANAN ERROR YANG LEBIH BAIK
	// if err != nil {
	// 	// Gunakan errors.Is untuk mengecek jenis error
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		c.JSON(http.StatusNotFound, gin.H{
	// 			"error":   "Tidak ada session aktif ditemukan",
	// 			"details": "Siswa belum memulai test atau session sudah disubmit",
	// 		})
	// 	} else if strings.Contains(err.Error(), "NIS tidak valid") {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Format NIS tidak valid"})
	// 	} else if strings.Contains(err.Error(), "data test tidak ditemukan") {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Data test tidak lengkap"})
	// 	} else {
	// 		// Log error untuk debugging
	// 		log.Printf("Internal Server Error in GetActiveTestSession: %v", err)
	// 		c.JSON(http.StatusInternalServerError, gin.H{
	// 			"error": "Terjadi kesalahan internal server",
	// 		})
	// 	}
	// 	return
	// }

	c.JSON(http.StatusOK, session)
}

// GET /siswa/test/:test_id/session
func GetTestSessionHandler(c *gin.Context) {
	testID, _ := strconv.Atoi(c.Param("test_id"))
	nis := c.Query("nis")

	session, err := siswaControllers.GetTestSession(uint(testID), nis)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session tidak ditemukan"})
		return
	}

	// Jika session masih aktif, hitung ulang waktu sisa
	if session.Status == "in_progress" {
		// Ambil data test untuk mendapatkan durasi
		var test models.TO_Test
		if err := config.DB.First(&test, session.TestID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data test"})
			return
		}

		// Hitung waktu sisa
		elapsedTime := time.Since(session.StartTime)
		totalDuration := time.Duration(test.DurasiMenit) * time.Minute
		remainingTime := totalDuration - elapsedTime

		// Jika waktu sudah habis, update status menjadi submitted
		if remainingTime <= 0 {
			endTime := session.StartTime.Add(totalDuration)
			session.Status = "submitted"
			session.EndTime = &endTime
			session.WaktuSisa = 0
			config.DB.Save(session)
		} else {
			// Update waktu sisa untuk response
			session.WaktuSisa = int(remainingTime.Seconds())
		}
	}

	c.JSON(http.StatusOK, session)
}

// GET /siswa/session/:session_id
func GetSessionByIDHandler(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID tidak valid"})
		return
	}

	fmt.Printf("Mencari session dengan ID: %d\n", sessionID) // Debug log

	var session models.TestSession
	err = config.DB.
		Preload("JawabanFinal").
		Preload("Test").
		Preload("JawabanFinal.Soal").
		Where("session_id = ?", sessionID). // Explicit WHERE clause
		First(&session).Error

	if err != nil {
		fmt.Printf("Error mencari session: %v\n", err) // Debug log
		c.JSON(http.StatusNotFound, gin.H{"error": "Session tidak ditemukan", "details": err.Error()})
		return
	}

	// fmt.Printf("Session ditemukan: %+v\n", session) // Debug log
	c.JSON(http.StatusOK, session)
}

func GetSoalHandler(c *gin.Context) {
	testID, err := strconv.Atoi(c.Param("test_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Test ID tidak valid"})
		return
	}

	// 1. CEK APAKAH ADA SESSION AKTIF UNTUK TEST INI
	var activeSessions []models.TestSession
	err = config.DB.
		Where("test_id = ? AND status IN ?", testID, []string{"in_progress", "onqueue"}).
		Order("start_time DESC").
		Find(&activeSessions).Error

	// 2. JIKA ADA SESSION AKTIF → AMBIL SOAL DARI to_sessionsoal (ambil session terbaru)
	if err == nil && len(activeSessions) > 0 {
		latestSession := activeSessions[0] // Ambil session terbaru
		log.Printf("🟢 Session aktif ditemukan: ID=%d untuk testID=%d, status=%s", latestSession.SessionID, testID, latestSession.Status)

		soal, err := siswaControllers.GetSoalBySessionID(latestSession.SessionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"session_id": latestSession.SessionID,
			"soal":       soal,
			"status":     latestSession.Status,
			"source":     "session_soal",
		})
		return
	}

	// 3. JIKA TIDAK ADA SESSION AKTIF → AMBIL SOAL BARU
	log.Printf("🟡 Tidak ada session aktif, ambil soal baru untuk testID: %d", testID)

	soal, err := siswaControllers.GetSoalByTestID1(uint(testID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": nil,
		"soal":       soal,
		"source":     "new_soal",
	})
}

// GET /siswa/session/:session_id/soal
func GetSessionSoalHandler(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID tidak valid"})
		return
	}

	soal, err := siswaControllers.GetSoalBySessionID(uint(sessionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, soal)
}

// POST /siswa/test/submit/:session_id
func SubmitSessionHandler(c *gin.Context) {
	tipeUjian := c.Param("tipe_ujian")
	sessionID, err := strconv.ParseUint(c.Param("session_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID tidak valid"})
		return
	}

	if err := siswaControllers.SubmitSession(uint(sessionID), tipeUjian); err != nil {
		log.Printf("❌ Gagal submit test session %d: %v", sessionID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal submit test"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Test berhasil disubmit",
		"submitted_at": time.Now(),
	})
}

// handler/siswaHandler.go
func GetNotStartedTestsHandler(c *gin.Context) {
	nis := c.Query("nis")
	if nis == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIS diperlukan"})
		return
	}

	tests, err := siswaControllers.GetNotStartedTests(nis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Gagal mengambil data test",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil test/tugas yang belum dikerjakan",
		"data":    tests,
		"total":   len(tests),
	})
}
