package adminhandlers

import (
	"bytes"
	"io"
	"log"
	"net/http"

	adminControllers "github.com/entertrans/bi-backend-go/controllers/admin"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

func GetAllTagihan(c *gin.Context) {
	data, err := adminControllers.FetchAllTagihan()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
func TambahTagihan(c *gin.Context) {
	var tagihan models.Tagihan

	// Debug body
	body, _ := io.ReadAll(c.Request.Body)
	log.Println("Body masuk:", string(body))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Binding JSON
	if err := c.ShouldBindJSON(&tagihan); err != nil {
		log.Println("[ERROR] Gagal binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	// Validasi
	if tagihan.Jenis == "" || tagihan.Nominal <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field tidak boleh kosong"})
		return
	}

	// Simpan
	if err := adminControllers.TambahTagihan(tagihan); err != nil {
		log.Println("[ERROR] Gagal menambah tagihan:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambah tagihan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tagihan berhasil ditambahkan"})
}

// Edit tagihan
func EditTagihan(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Jenis   string `json:"jns_tagihan"`
		Nominal int    `json:"nom_tagihan"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[ERROR] Gagal binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	if req.Jenis == "" || req.Nominal <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field tidak boleh kosong atau nominal salah"})
		return
	}

	if err := adminControllers.UpdateTagihan(id, req.Jenis, req.Nominal); err != nil {
		log.Println("[ERROR] Gagal update tagihan:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update tagihan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tagihan berhasil diperbarui"})
}

// Delete tagihan
func DeleteTagihan(c *gin.Context) {
	id := c.Param("id")

	if err := adminControllers.DeleteTagihan(id); err != nil {
		log.Println("[ERROR] Gagal hapus tagihan:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus tagihan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tagihan berhasil dihapus"})
}
