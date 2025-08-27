package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"

	"github.com/entertrans/bi-backend-go/controllers"
	"github.com/gin-gonic/gin"
)

func GetActivePenggunaHandler(c *gin.Context) {
	users, err := controllers.GetActivePengguna()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna aktif"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func LoginHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	// Get user dengan semua relasi (Siswa, Guru, Admin)
	user, err := controllers.GetPenggunaByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username tidak ditemukan"})
		return
	}

	// cek status aktif
	if user.PenggunaStatus != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun tidak aktif"})
		return
	}

	// hash input password pakai MD5
	hasher := md5.New()
	hasher.Write([]byte(req.Password))
	md5Password := hex.EncodeToString(hasher.Sum(nil))

	if user.PenggunaPassword != md5Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
		return
	}

	// tentukan halaman / route berdasarkan pengguna_level
	var redirect string
	switch user.PenggunaLevel {
	case "1": // ADMIN
		redirect = "/admin"
	case "2": // SISWA
		redirect = "/siswa"
	case "3": // GURU
		redirect = "/guru"
	default:
		redirect = "/"
	}

	// Response dengan format yang diinginkan
	response := gin.H{
		"message":     "Login sukses",
		"pengguna_id": user.PenggunaID,
		"username":    user.PenggunaUsername,
		"level":       user.PenggunaLevel,
		"ref_id":      user.RefID,
		"redirect":    redirect,
		"admin":       user.Admin,
		"siswa":       user.Siswa,
		"Guru":        user.Guru, // Perhatikan huruf kapital 'G' sesuai permintaan
	}

	c.JSON(http.StatusOK, response)
}