package guruhandlers

import (
	"net/http"
	"strconv"

	gurucontrollers "github.com/entertrans/bi-backend-go/controllers/guru"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

func GetAllGuruHandler(c *gin.Context) {
	gurus, err := gurucontrollers.GetAllGuru()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data guru"})
		return
	}
	c.JSON(http.StatusOK, gurus)
}

func GetGuruByIDHandler(c *gin.Context) {
	idStr := c.Param("guru_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Guru ID tidak valid"})
		return
	}

	guru, err := gurucontrollers.GetGuruByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Guru tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, guru)
}

func CreateGuruHandler(c *gin.Context) {
	var guru models.Guru
	if err := c.ShouldBindJSON(&guru); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	if err := gurucontrollers.CreateGuru(&guru); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat guru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Guru berhasil dibuat", "guru_id": guru.GuruID})
}

func UpdateGuruHandler(c *gin.Context) {
	idStr := c.Param("guru_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Guru ID tidak valid"})
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	if err := gurucontrollers.UpdateGuru(uint(id), data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update guru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Guru berhasil diupdate"})
}

func DeleteGuruHandler(c *gin.Context) {
	idStr := c.Param("guru_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Guru ID tidak valid"})
		return
	}

	if err := gurucontrollers.DeleteGuru(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus guru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Guru berhasil dihapus"})
}
