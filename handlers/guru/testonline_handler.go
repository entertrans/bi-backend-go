package guruhandlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/entertrans/bi-backend-go/config"
	gurucontrollers "github.com/entertrans/bi-backend-go/controllers/guru"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Create test
// func CreateTestHandler(c *gin.Context) {
// 	var test models.TO_Test
// 	if err := c.ShouldBindJSON(&test); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if err := gurucontrollers.CreateTest(&test); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, test)
// }

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

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Test berhasil dibuat",
		"test_id":  test.TestID,
		"soal_ids": test.SoalIDs,
	})
}

// new
type AddSoalRequest struct {
	SoalIDs []uint `json:"soal_ids" binding:"required"`
}

func AddSoalToTestHandler(c *gin.Context) {
	// Ambil testID dari path parameter
	testIDStr := c.Param("testId")
	testID, err := strconv.ParseUint(testIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test ID"})
		return
	}

	var req AddSoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := gurucontrollers.AddSoalToTest(uint(testID), req.SoalIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Soal berhasil ditambahkan ke test",
		"count":   len(req.SoalIDs),
	})
}

// GetBankSoalByKelasMapelHandler handler untuk endpoint /banksoal/by-kelas-mapel
func GetBankSoalByKelasMapelHandler(c *gin.Context) {
	// Ambil parameter dari query string
	kelasIDStr := c.Query("kelas_id")
	mapelIDStr := c.Query("mapel_id")
	testIDStr := c.Query("test_id")

	// Validasi parameter
	if kelasIDStr == "" || mapelIDStr == "" || testIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "parameter kelas_id, mapel_id, dan test_id diperlukan",
		})
		return
	}

	// Konversi parameter
	kelasID, err := parseUintParam(kelasIDStr, "kelas_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	mapelID, err := parseUintParam(mapelIDStr, "mapel_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	testID, err := parseUintParam(testIDStr, "test_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Panggil controller
	soals, selectedSoalIDs, err := gurucontrollers.GetBankSoalByKelasMapel(kelasID, mapelID, testID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "gagal mengambil data bank soal",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"soals":             soals,
			"selected_soal_ids": selectedSoalIDs,
		},
	})
}

// parseUintParam helper function untuk parsing parameter
func parseUintParam(paramStr, paramName string) (uint, error) {
	param, err := strconv.ParseUint(paramStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("parameter %s harus berupa angka", paramName)
	}
	return uint(param), nil
}

// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"message": "Gagal mengambil data bank soal",
// 			"error":   err.Error(),
// 		})
// 		return
// 	}

//		c.JSON(http.StatusOK, gin.H{
//			"success": true,
//			"data": gin.H{
//				"soals": soals,
//				"selected_soal_ids": selectedSoalIDs,
//			},
//		})
//	}
func RemoveSoalFromTestHandler(c *gin.Context) {
	// Ambil parameter dari URL
	testIDStr := c.Param("testId")
	soalIDStr := c.Param("soalId")

	// Validasi parameter
	if testIDStr == "" || soalIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter testId dan soalId diperlukan",
		})
		return
	}

	// Konversi string ke uint
	testID, err := strconv.ParseUint(testIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter testId harus berupa angka",
		})
		return
	}

	soalID, err := strconv.ParseUint(soalIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter soalId harus berupa angka",
		})
		return
	}

	// Panggil controller function
	err = gurucontrollers.RemoveSoalFromTest(uint(testID), uint(soalID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Soal tidak ditemukan dalam test",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus soal dari test",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Soal berhasil dihapus dari test",
	})
}

// end new

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
