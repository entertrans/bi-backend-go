package guruhandlers

import (
	"log"
	"net/http"
	"strconv"

	gurucontrollers "github.com/entertrans/bi-backend-go/controllers/guru"
	"github.com/gin-gonic/gin"
)

func GetJawabanBySiswaHandler(c *gin.Context) {
	siswaNIS := c.Param("siswa_nis")
	if siswaNIS == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Siswa NIS harus diisi"})
		return
	}

	results, err := gurucontrollers.GetJawabanBySiswaNIS(siswaNIS)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"siswa_nis": siswaNIS,
		"results":   results,
		"count":     len(results),
	})
}

func GetDetailJawabanHandler(c *gin.Context) {
	sessionIDStr := c.Param("session_id")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID tidak valid"})
		return
	}

	result, err := gurucontrollers.GetDetailJawabanBySession(uint(sessionID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func GetSiswaByTestHandler(c *gin.Context) {
	testIDStr := c.Param("test_id")
	testID, err := strconv.ParseUint(testIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Test ID tidak valid"})
		return
	}

	results, err := gurucontrollers.GetSiswaByTest(uint(testID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"test_id": testID,
		"siswa":   results,
		"count":   len(results),
	})
}

func GetTestStatisticsHandler(c *gin.Context) {
	testIDStr := c.Param("test_id")
	testID, err := strconv.ParseUint(testIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Test ID tidak valid"})
		return
	}

	stats, err := gurucontrollers.GetTestStatistics(uint(testID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statistics": stats,
	})
}

func GetSiswaDetailForGuruHandler(c *gin.Context) {
	siswaNIS := c.Param("siswa_nis")
	siswaDetail, err := gurucontrollers.GetSiswaDetailForGuru(siswaNIS)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": siswaDetail}) // Wrap dalam { data: }
}

func GetJawabanBySession(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID tidak valid"})
		return
	}

	data, err := gurucontrollers.FetchJawabanBySession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// update nilai
func UpdateJawabanFinal(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("session_id"))
	if err != nil {
		log.Printf("[ERROR] Session ID tidak valid: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID tidak valid"})
		return
	}

	log.Printf("[INFO] PUT /session/%d/jawaban", sessionID)

	var req struct {
		Perubahan []struct {
			SessionID int     `json:"session_id"`
			SoalID    uint    `json:"soal_id"`
			Nilai     float64 `json:"nilai"`
		} `json:"perubahan"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERROR] Gagal binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	log.Printf("[DEBUG] Request body: %+v", req)

	if err := gurucontrollers.UpdateNilaiJawaban(sessionID, req.Perubahan); err != nil {
		log.Printf("[ERROR] Gagal update nilai jawaban: %v", err)
		// Berikan response yang lebih spesifik
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Beberapa nilai mungkin tidak berhasil diupdate. Periksa apakah soal tersebut ada dalam jawaban siswa.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Nilai jawaban berhasil diupdate",
		"success": true,
	})
}

func UpdateOverrideNilai(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("session_id"))
	if err != nil {
		log.Println("[ERROR] Session ID tidak valid:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID tidak valid"})
		return
	}

	log.Printf("[INFO] PUT /session/%d/nilai-akhir", sessionID)

	var req struct {
		NilaiAkhir float64 `json:"nilai_akhir"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[ERROR] Gagal binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	if err := gurucontrollers.UpdateNilaiAkhir(sessionID, req.NilaiAkhir); err != nil {
		log.Println("[ERROR]", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update nilai akhir"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nilai akhir berhasil diupdate"})
}

func ResetTestSessionHandler(c *gin.Context) {
	sessionIdStr := c.Param("session_id")
	sessionId, err := strconv.ParseUint(sessionIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Session ID tidak valid",
		})
		return
	}

	err = gurucontrollers.ResetTestSession(uint(sessionId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal reset test: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Test berhasil direset",
	})
}

// func GetSoalPenilaianHandler(c *gin.Context) {
// 	sessionIDStr := c.Param("test_id")
// 	sessionID, err := strconv.Atoi(sessionIDStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID tidak valid"})
// 		return
// 	}

// 	// cari TestID dari SessionID
// 	var session models.TestSession
// 	if err := config.DB.First(&session, sessionID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Session tidak ditemukan"})
// 		return
// 	}

// 	soal, err := gurucontrollers.GetSoalByPenilaianID(config.DB, session.TestID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, soal)
// }
