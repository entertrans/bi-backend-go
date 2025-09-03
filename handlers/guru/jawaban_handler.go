package guruhandlers

import (
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