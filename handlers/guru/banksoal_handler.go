package guruhandlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	gurucontrollers "github.com/entertrans/bi-backend-go/controllers/guru"
	"github.com/gin-gonic/gin"
)

func GetActiveBankSoalHandler(c *gin.Context) {
	soal, err := gurucontrollers.GetActiveBankSoal()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil bank soal aktif"})
		return
	}
	c.JSON(http.StatusOK, soal)
}

func GetInactiveBankSoalHandler(c *gin.Context) {
	soal, err := gurucontrollers.GetInactiveBankSoal()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil bank soal nonaktif"})
		return
	}
	c.JSON(http.StatusOK, soal)
}

func GetBankSoalHandler(c *gin.Context) {
	guruIDStr := c.Param("guru_id")
	guruID, err := strconv.ParseUint(guruIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Guru ID tidak valid"})
		return
	}

	soal, err := gurucontrollers.GetBankSoalByGuru(uint(guruID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil bank soal"})
		return
	}

	c.JSON(http.StatusOK, soal)
}

func DeleteBankSoalHandler(c *gin.Context) {
	soalIDStr := c.Param("soal_id")
	soalID, err := strconv.ParseUint(soalIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Soal ID tidak valid"})
		return
	}

	if err := gurucontrollers.DeleteBankSoal(uint(soalID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus soal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Soal berhasil dihapus"})
}
func RestoreBankSoalHandler(c *gin.Context) {
	soalIDStr := c.Param("soal_id")
	soalID, err := strconv.ParseUint(soalIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Soal ID tidak valid"})
		return
	}

	err = gurucontrollers.RestoreBankSoal(uint(soalID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengembalikan bank soal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bank soal berhasil direstore"})
}
func BuatSoalHandler(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	log.Println("RAW BODY:", string(body))

	var input gurucontrollers.BankSoalInput
	if err := json.Unmarshal(body, &input); err != nil {
		log.Println("Unmarshal error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	if err := gurucontrollers.SimpanSoal(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Soal berhasil disimpan"})
}

func GetActiveBankSoalByKelasMapelHandler(c *gin.Context) {
	kelasIDStr := c.Param("kelas_id")
	mapelIDStr := c.Param("mapel_id")
	mapelID, _ := strconv.ParseUint(mapelIDStr, 10, 64)

	kelasID, err := strconv.ParseUint(kelasIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kelas ID tidak valid"})
		return
	}

	soal, err := gurucontrollers.GetActiveBankSoalByKelasMapel(uint(kelasID), uint(mapelID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil bank soal"})
		return
	}

	c.JSON(http.StatusOK, soal)
}

func GetRekapBankSoalHandler(c *gin.Context) {
	rekap, err := gurucontrollers.GetRekapBankSoal()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil rekap bank soal"})
		return
	}
	c.JSON(http.StatusOK, rekap)
}
