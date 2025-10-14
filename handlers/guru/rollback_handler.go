package guruhandlers

import (
	"net/http"
	"strconv"

	gurucontrollers "github.com/entertrans/bi-backend-go/controllers/guru"
	"github.com/gin-gonic/gin"
)

func GetJawabanRollbackHandler(c *gin.Context) {
	siswaNIS := c.Param("siswa_nis")
	if siswaNIS == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Siswa NIS harus diisi"})
		return
	}

	results, err := gurucontrollers.GetJawabanRollbackBySiswaNIS(siswaNIS)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var kelasUniversal string
	if len(results) > 0 {
		kelasUniversal = results[0].Kelas
	}
	for i := range results {
		results[i].Kelas = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"siswa_nis": siswaNIS,
		"kelas":     kelasUniversal,
		"results":   results,
		"count":     len(results),
	})
}

func GetJawabanRollbackBySession(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID tidak valid"})
		return
	}

	jenis := c.Param("jenis")

	data, err := gurucontrollers.FetchRBJawabanBySession(sessionID, jenis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

