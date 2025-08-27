package guruhandlers

import (
	"net/http"
	"strconv"

	gurucontrollers "github.com/entertrans/bi-backend-go/controllers/guru"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

// Tambah peserta (satu siswa)
func AddPesertaHandler(c *gin.Context) {
	var peserta models.TO_Peserta
	if err := c.ShouldBindJSON(&peserta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := gurucontrollers.AddPeserta(&peserta); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, peserta)
}

// Ambil semua peserta dalam test
func GetPesertaByTestHandler(c *gin.Context) {
	testIDStr := c.Param("test_id")
	testID, _ := strconv.Atoi(testIDStr)

	pes, err := gurucontrollers.GetPesertaByTestID(uint(testID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pes)
}

// Update peserta
func UpdatePesertaHandler(c *gin.Context) {
	pesertaIDStr := c.Param("peserta_id")
	pesertaID, _ := strconv.Atoi(pesertaIDStr)

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := gurucontrollers.UpdatePeserta(uint(pesertaID), data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Peserta updated"})
}

// Delete peserta
func DeletePesertaHandler(c *gin.Context) {
	pesertaIDStr := c.Param("peserta_id")
	pesertaID, _ := strconv.Atoi(pesertaIDStr)

	if err := gurucontrollers.DeletePeserta(uint(pesertaID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Peserta deleted"})
}

func GetAvailableSiswaByKelasHandler(c *gin.Context) {
	kelasIDStr := c.Param("kelas_id")
	testIDStr := c.Param("test_id")

	kelasID, err := strconv.Atoi(kelasIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "kelas_id tidak valid"})
		return
	}

	testID, err := strconv.Atoi(testIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "test_id tidak valid"})
		return
	}

	siswa, err := gurucontrollers.GetAvailableSiswaByKelas(uint(kelasID), uint(testID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, siswa)
}

func GetTestSoalByTestIdHandler(c *gin.Context) {
	testIdStr := c.Param("test_id")
	testId, err := strconv.ParseUint(testIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Test ID tidak valid",
		})
		return
	}

	soals, err := gurucontrollers.GetTestSoalByTestId(uint(testId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data soal: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    soals,
		"message": "Data soal berhasil diambil",
	})
}

// GetDetailTestSoalHandler - Get detail soal by ID
func GetDetailTestSoalHandler(c *gin.Context) {
	soalIdStr := c.Param("soal_id")
	soalId, err := strconv.ParseUint(soalIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Soal ID tidak valid",
		})
		return
	}

	soal, err := gurucontrollers.GetDetailTestSoal(uint(soalId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil detail soal: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    soal,
		"message": "Detail soal berhasil diambil",
	})
}

// DeleteTestSoalHandler - Delete soal (soft delete)
func DeleteTestSoalHandler(c *gin.Context) {
	soalIdStr := c.Param("soal_id")
	soalId, err := strconv.ParseUint(soalIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Soal ID tidak valid",
		})
		return
	}

	err = gurucontrollers.DeleteTestSoal(uint(soalId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus soal: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Soal berhasil dihapus",
	})
}

// CreateTestSoalHandler - Create new test soal
func CreateTestSoalHandler(c *gin.Context) {
	var input gurucontrollers.TestSoalInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input tidak valid: " + err.Error(),
		})
		return
	}

	soal, err := gurucontrollers.CreateTestSoal(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat soal: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    soal,
		"message": "Soal berhasil dibuat",
	})
}
