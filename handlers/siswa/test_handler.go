package siswa

import (
	"net/http"
	"strconv"

	siswaControllers "github.com/entertrans/bi-backend-go/controllers/siswa"

	"github.com/gin-gonic/gin"
)

// Endpoint GET /siswa/tests/ub
func GetAllUBTestHandler(c *gin.Context) {
	tests, err := siswaControllers.GetAllUBTest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data test"})
		return
	}
	c.JSON(http.StatusOK, tests)
}

func GetUBTestByKelasHandler(c *gin.Context) {
	idParam := c.Param("kelas_id")
	kelasID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kelas tidak valid"})
		return
	}

	tests, err := siswaControllers.GetUBTestByKelas(uint(kelasID))
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

	soals, err := siswaControllers.GetSoalByTestID(uint(testID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil soal"})
		return
	}

	c.JSON(http.StatusOK, soals)
}
