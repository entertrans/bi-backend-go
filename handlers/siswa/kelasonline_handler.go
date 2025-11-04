package siswa

import (
	"net/http"

	siswaControllers "github.com/entertrans/bi-backend-go/controllers/siswa"
	"github.com/gin-gonic/gin"
)

// 📘 Endpoint 1: daftar kelas online berdasarkan kelas_id
func GetKelasOnlineByKelasIDHandler(c *gin.Context) {
	kelasID := c.Param("kelas_id")

	result, err := siswaControllers.GetKelasOnlineByKelasID(kelasID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// 📘 Endpoint 2: daftar riwayat kelas berdasarkan id_kelas_mapel
func GetKelasOnlineHistoryHandler(c *gin.Context) {
	idKelasMapel := c.Param("id_kelas_mapel")

	result, err := siswaControllers.GetKelasOnlineHistory(idKelasMapel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}


func GetMapelByKelasHandler(c *gin.Context) {
    kelasID := c.Param("kelas_id")
    result, err := siswaControllers.GetMapelByKelasID(kelasID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, result)
}