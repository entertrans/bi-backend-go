package guruhandlers

import (
	"net/http"
	"path/filepath"
	"strconv"

	gurucontrollers "github.com/entertrans/bi-backend-go/controllers/guru"
	"github.com/gin-gonic/gin"
)

// =========================
// Upload Lampiran
// =========================
func UploadLampiranHandler(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Lampiran tidak ditemukan"})
        return
    }

    deskripsi := c.PostForm("deskripsi")

    // Simpan file ke folder
    savePath := filepath.Join("uploads/lampiran", file.Filename)
    if err := c.SaveUploadedFile(file, savePath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
        return
    }

    // Ambil MIME type
    mimeType := file.Header.Get("Content-Type")

    // Simpan metadata ke DB (mapping MIME → ENUM)
    lampiran, err := gurucontrollers.UploadLampiran(
        savePath,
        mimeType,
        file.Filename,
        deskripsi,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan lampiran"})
        return
    }

    c.JSON(http.StatusOK, lampiran)
}


// =========================
// Get Lampiran Aktif
// =========================
func GetActiveLampiranHandler(c *gin.Context) {
	data, err := gurucontrollers.GetActiveLampiran()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil lampiran"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// =========================
// Get Lampiran Trash
// =========================
func GetInactiveLampiranHandler(c *gin.Context) {
	data, err := gurucontrollers.GetInactiveLampiran()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil lampiran trash"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// =========================
// Delete Lampiran (Soft Delete)
// =========================
func DeleteLampiranHandler(c *gin.Context) {
	idStr := c.Param("lampiran_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lampiran ID tidak valid"})
		return
	}

	if err := gurucontrollers.DeleteLampiran(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lampiran berhasil dihapus"})
}

// =========================
// Restore Lampiran
// =========================
func RestoreLampiranHandler(c *gin.Context) {
	idStr := c.Param("lampiran_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lampiran ID tidak valid"})
		return
	}

	if err := gurucontrollers.RestoreLampiran(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lampiran berhasil direstore"})
}

// =========================
// Hard Delete Lampiran
// =========================
func HardDeleteLampiranHandler(c *gin.Context) {
	idStr := c.Param("lampiran_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lampiran ID tidak valid"})
		return
	}

	if err := gurucontrollers.HardDeleteLampiran(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lampiran berhasil dihapus permanen"})
}
