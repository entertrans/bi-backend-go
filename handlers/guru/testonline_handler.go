package guruhandlers

import (
	"net/http"
	"strconv"

	"github.com/entertrans/bi-backend-go/config"
	gurucontrollers "github.com/entertrans/bi-backend-go/controllers/guru"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

// Create test
func CreateTestHandler(c *gin.Context) {
	var test models.TO_Test
	if err := c.ShouldBindJSON(&test); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := gurucontrollers.CreateTest(&test); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, test)
}

// Get tests by type
func GetTestByType(c *gin.Context) {
	tipe := c.Param("type_test")
	tests, err := gurucontrollers.GetTestByType(tipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tests)
}

// Get tests by guru
func GetTestsByGuruHandler(c *gin.Context) {
	guruIDStr := c.Param("guru_id")
	guruID, _ := strconv.ParseUint(guruIDStr, 10, 64)

	tests, err := gurucontrollers.GetTestsByGuruID(uint(guruID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tests)
}

// Get single test
func GetTestHandler(c *gin.Context) {
	idStr := c.Param("test_id")
	id, _ := strconv.Atoi(idStr)

	test, err := gurucontrollers.GetTestByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test not found"})
		return
	}
	c.JSON(http.StatusOK, test)
}

// Update test
func UpdateTestAktifHandler(c *gin.Context) {
	idStr := c.Param("test_id")
	id, _ := strconv.Atoi(idStr)

	var payload struct {
		Aktif int `json:"aktif"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&models.TO_Test{}).
		Where("test_id = ?", id).
		Update("aktif", payload.Aktif).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status aktif updated"})
}

// Delete test
func DeleteTestHandler(c *gin.Context) {
	idStr := c.Param("test_id")
	id, _ := strconv.Atoi(idStr)

	if err := gurucontrollers.DeleteTest(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Test deleted"})
}
