package guruhandlers

import (
	"net/http"
	"strconv"

	gurucontrollers "github.com/entertrans/bi-backend-go/controllers/guru"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

func CreateTestHandler(c *gin.Context) {
	var test models.TO_Test
	if err := c.ShouldBindJSON(&test); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data test tidak valid"})
		return
	}

	if err := gurucontrollers.CreateTest(&test); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat test"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test berhasil dibuat", "test_id": test.TestID})
}

func GetTestByIDHandler(c *gin.Context) {
	testIDStr := c.Param("test_id")
	testID, err := strconv.ParseUint(testIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Test ID tidak valid"})
		return
	}

	test, err := gurucontrollers.GetTestByID(uint(testID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, test)
}

func GetTestsByGuruHandler(c *gin.Context) {
	guruIDStr := c.Param("guru_id")
	guruID, err := strconv.ParseUint(guruIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Guru ID tidak valid"})
		return
	}

	tests, err := gurucontrollers.GetTestsByGuruID(uint(guruID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil tests"})
		return
	}

	c.JSON(http.StatusOK, tests)
}

func UpdateTestHandler(c *gin.Context) {
	testIDStr := c.Param("test_id")
	testID, err := strconv.ParseUint(testIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Test ID tidak valid"})
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	if err := gurucontrollers.UpdateTest(uint(testID), data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate test"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test berhasil diupdate"})
}

func DeleteTestHandler(c *gin.Context) {
	testIDStr := c.Param("test_id")
	testID, err := strconv.ParseUint(testIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Test ID tidak valid"})
		return
	}

	if err := gurucontrollers.DeleteTest(uint(testID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus test"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test berhasil dihapus"})
}

func GetTestHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetTestHandler belum dibuat"})
}

func GetPenilaianHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetPenilaianHandler belum dibuat"})
}

func DeletePenilaianHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DeletePenilaianHandler belum dibuat"})
}
