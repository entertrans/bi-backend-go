package adminhandlers

import (
	"net/http"
	"strconv"

	admincontrollers "github.com/entertrans/bi-backend-go/controllers/admin"
	"github.com/gin-gonic/gin"
)

func GetMapelHandler(c *gin.Context) {
    mapels, err := admincontrollers.GetAllMapelWithGuruMapels()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch mapel"})
        return
    }
    c.JSON(http.StatusOK, mapels)
}

func GetSiswaByKelasHandler(c *gin.Context) {
	kelasIDStr := c.Param("kelas_id")
	kelasID, err := strconv.Atoi(kelasIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "kelas_id tidak valid"})
		return
	}

	siswa, err := admincontrollers.GetSiswaByKelas(uint(kelasID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, siswa)
}