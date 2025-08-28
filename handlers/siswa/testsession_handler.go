package siswa

import (
	"fmt"
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

	session, err := siswaControllers.StartTestSession(uint(testID), nis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai test: " + err.Error()})
		return
	}

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
		Preload("JawabanFinal.Soal").
		Where("session_id = ?", sessionID). // Explicit WHERE clause
		First(&session).Error

	if err != nil {
		fmt.Printf("Error mencari session: %v\n", err) // Debug log
		c.JSON(http.StatusNotFound, gin.H{"error": "Session tidak ditemukan", "details": err.Error()})
		return
	}

	fmt.Printf("Session ditemukan: %+v\n", session) // Debug log
	c.JSON(http.StatusOK, session)
}

func GetSoalHandler(c *gin.Context) {
	testID, _ := strconv.Atoi(c.Param("test_id"))

	soal, err := siswaControllers.GetSoalByTestID1(uint(testID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, soal)
}

// POST /siswa/test/submit/:session_id
func SubmitTestHandler(c *gin.Context) {
	sessionID, _ := strconv.Atoi(c.Param("session_id"))

	err := siswaControllers.SubmitTestSession(uint(sessionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal submit test"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test berhasil disubmit", "submitted_at": time.Now()})
}
