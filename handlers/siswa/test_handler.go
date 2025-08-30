package siswa

import (
	"log"
	"net/http"
	"strconv"

	"github.com/entertrans/bi-backend-go/config"
	siswaControllers "github.com/entertrans/bi-backend-go/controllers/siswa"
	"github.com/entertrans/bi-backend-go/models"

	"github.com/gin-gonic/gin"
)

// Endpoint GET /siswa/tests/ub
// func GetAllUBTestHandler(c *gin.Context) {
// 	tests, err := siswaControllers.GetAllUBTest()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data test"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, tests)
// }

// func GetUBTestByKelasHandler(c *gin.Context) {
// 	idParam := c.Param("kelas_id")
// 	kelasID, err := strconv.Atoi(idParam)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kelas tidak valid"})
// 		return
// 	}

//		tests, err := siswaControllers.GetUBTestByKelas(uint(kelasID))
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data test"})
//			return
//		}
//		c.JSON(http.StatusOK, tests)
//	}
func GetTestByKelasHandler(c *gin.Context) {
	// Ambil param kelas_id
	idParam := c.Param("kelas_id")
	kelasID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kelas tidak valid"})
		return
	}

	// Ambil param type_test
	typeTest := c.Param("type_test")
	if typeTest == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type test tidak boleh kosong"})
		return
	}

	// Panggil controller
	tests, err := siswaControllers.GetTestByKelas(uint(kelasID), typeTest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data test"})
		return
	}

	c.JSON(http.StatusOK, tests)
}

// Endpoint GET /siswa/tests/:id/soal
func GetSoalByTestIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	testID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	log.Printf("Mencari soal untuk test_id: %d", testID)

	// Cek apakah test exists
	var test models.TO_Test
	if err := config.DB.First(&test, testID).Error; err != nil {
		log.Printf("Test dengan ID %d tidak ditemukan: %v", testID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Test tidak ditemukan"})
		return
	}

	soals, err := siswaControllers.GetSoalByTestID(uint(testID))
	if err != nil {
		log.Printf("Error mengambil soal: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil soal"})
		return
	}

	log.Printf("Jumlah soal ditemukan: %d", len(soals))

	c.JSON(http.StatusOK, soals)
}

//testreview
