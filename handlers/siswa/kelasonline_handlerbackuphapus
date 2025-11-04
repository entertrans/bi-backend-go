package siswa

import (
	"net/http"
	"strconv"

	siswaControllers "github.com/entertrans/bi-backend-go/controllers/siswa"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/entertrans/bi-backend-go/utils"
	"github.com/gin-gonic/gin"
)

var OnlineClassController = &siswaControllers.OnlineClassController{}

// GetAllOnlineClassHandler mendapatkan semua jadwal online class
func GetAllOnlineClassHandler(c *gin.Context) {
	onlineClasses, err := OnlineClassController.GetAllOnlineClass()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal mengambil data online class")
		return
	}

	utils.RespondJSON(c, http.StatusOK, onlineClasses)
}

// GetOnlineClassByIDHandler mendapatkan satu jadwal online class by ID
func GetOnlineClassByIDHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	onlineClass, err := OnlineClassController.GetOnlineClassByID(uint(id))
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "Online class tidak ditemukan")
		return
	}

	utils.RespondJSON(c, http.StatusOK, onlineClass)
}

// GetOnlineClassByKelasHandler mendapatkan jadwal online class by kelas_id
func GetOnlineClassByKelasHandler(c *gin.Context) {
	kelasID, err := strconv.Atoi(c.Param("kelas_id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Kelas ID tidak valid")
		return
	}

	onlineClasses, err := OnlineClassController.GetOnlineClassByKelas(uint(kelasID))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal mengambil data online class")
		return
	}

	utils.RespondJSON(c, http.StatusOK, onlineClasses)
}

// GetOnlineClassByMapelHandler mendapatkan jadwal online class by mapel_id
func GetOnlineClassByMapelHandler(c *gin.Context) {
	mapelID, err := strconv.Atoi(c.Param("mapel_id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Mapel ID tidak valid")
		return
	}

	onlineClasses, err := OnlineClassController.GetOnlineClassByMapel(uint(mapelID))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal mengambil data online class")
		return
	}

	utils.RespondJSON(c, http.StatusOK, onlineClasses)
}

// CreateOnlineClassHandler membuat jadwal online class baru
func CreateOnlineClassHandler(c *gin.Context) {
	var onlineClass models.OnlineClass
	
	if err := c.ShouldBindJSON(&onlineClass); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Data input tidak valid")
		return
	}

	// Validasi data
	if onlineClass.IDKelasMapel == 0 {
		utils.RespondError(c, http.StatusBadRequest, "ID Kelas Mapel harus diisi")
		return
	}

	if onlineClass.Tanggal.IsZero() {
		utils.RespondError(c, http.StatusBadRequest, "Tanggal harus diisi")
		return
	}

	if onlineClass.Mulai == "" {
		utils.RespondError(c, http.StatusBadRequest, "Waktu mulai harus diisi")
		return
	}

	if onlineClass.Selesai == "" {
		utils.RespondError(c, http.StatusBadRequest, "Waktu selesai harus diisi")
		return
	}

	if onlineClass.MeetLink == "" {
		utils.RespondError(c, http.StatusBadRequest, "Link meeting harus diisi")
		return
	}

	err := OnlineClassController.CreateOnlineClass(&onlineClass)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal membuat online class")
		return
	}

	utils.RespondJSON(c, http.StatusCreated, onlineClass)
}

// UpdateOnlineClassHandler mengupdate jadwal online class
func UpdateOnlineClassHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var onlineClass models.OnlineClass
	if err := c.ShouldBindJSON(&onlineClass); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Data input tidak valid")
		return
	}

	err = OnlineClassController.UpdateOnlineClass(uint(id), &onlineClass)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal mengupdate online class")
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{"message": "Online class berhasil diupdate"})
}

// DeleteOnlineClassHandler menghapus jadwal online class
func DeleteOnlineClassHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	err = OnlineClassController.DeleteOnlineClass(uint(id))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal menghapus online class")
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{"message": "Online class berhasil dihapus"})
}