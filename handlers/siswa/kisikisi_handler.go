package siswa

import (
	"net/http"
	"strconv"

	"github.com/entertrans/bi-backend-go/controllers/siswa"
	siswaControllers "github.com/entertrans/bi-backend-go/controllers/siswa"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/entertrans/bi-backend-go/utils"
	"github.com/gin-gonic/gin"
)

// GetKisiKisiByIDHandler mendapatkan kisi-kisi by ID
func GetKisiKisiByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	kisiKisi, err := siswaControllers.GetKisiKisiByID(uint(id))
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "Kisi-kisi tidak ditemukan")
		return
	}

	utils.RespondJSON(c, http.StatusOK, kisiKisi)
}

// GetAllKisiKisiHandler mendapatkan semua kisi-kisi
func GetAllKisiKisiHandler(c *gin.Context) {
	kisiKisis, err := siswaControllers.GetAllKisiKisi()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal mengambil data kisi-kisi")
		return
	}

	utils.RespondJSON(c, http.StatusOK, kisiKisis)
}

// GetKisiKisiByKelasHandler mendapatkan kisi-kisi by kelas
func GetKisiKisiByKelasHandler(c *gin.Context) {
	kelasIDStr := c.Param("kelas_id")
	kelasID, err := strconv.ParseUint(kelasIDStr, 10, 32)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "ID kelas tidak valid")
		return
	}

	kisiKisis, err := siswaControllers.GetKisiKisiByKelas(uint(kelasID))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal mengambil data kisi-kisi")
		return
	}

	utils.RespondJSON(c, http.StatusOK, kisiKisis)
}

// GetKisiKisiByMapelHandler tanpa pagination (return semua data)
func GetKisiKisiByMapelHandler(c *gin.Context) {
	mapelIDStr := c.Param("mapel_id")
	mapelID, err := strconv.ParseUint(mapelIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "ID mapel tidak valid",
		})
		return
	}

	kisiKisis, err := siswaControllers.GetKisiKisiByMapel(uint(mapelID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal mengambil data kisi-kisi",
		})
		return
	}

	// Return langsung array data tanpa wrapper
	c.JSON(http.StatusOK, kisiKisis)
}

// CreateKisiKisiHandler membuat kisi-kisi baru
func CreateKisiKisiHandler(c *gin.Context) {
	var input models.KisiKisi

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Data tidak valid")
		return
	}

	if input.KisiKisiUb == "" || input.KisiKisiDeskripsi == "" || input.KisiKisiKelasID == 0 {
		utils.RespondError(c, http.StatusBadRequest, "UB, deskripsi, dan kelas ID harus diisi")
		return
	}

	if err := siswa.CreateKisiKisi(&input); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal membuat kisi-kisi")
		return
	}

	utils.RespondJSON(c, http.StatusCreated, input)
}

// UpdateKisiKisiHandler mengupdate kisi-kisi
func UpdateKisiKisiHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Data update tidak valid")
		return
	}

	if err := siswaControllers.UpdateKisiKisi(uint(id), updateData); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal mengupdate kisi-kisi")
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{"message": "Kisi-kisi berhasil diupdate"})
}

// DeleteKisiKisiHandler menghapus kisi-kisi
func DeleteKisiKisiHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	if err := siswaControllers.DeleteKisiKisi(uint(id)); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal menghapus kisi-kisi")
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{"message": "Kisi-kisi berhasil dihapus"})
}