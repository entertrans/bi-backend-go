package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" // Untuk baca file .env
	"gorm.io/driver/mysql"     // Driver MySQL untuk GORM
	"gorm.io/gorm"             // ORM-nya
)

var DB *gorm.DB // Variabel global untuk koneksi DB

// LoadEnv untuk memuat file .env
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Gagal load file .env")
	}
}

// ConnectDB untuk menghubungkan ke database MySQL
func ConnectDB() {
	LoadEnv() // Panggil dulu fungsi load .env

	// Ambil konfigurasi dari .env
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Coba konek ke DB
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Gagal konek ke DB:", err)
	}

	DB = db // Simpan koneksi ke variabel global
	fmt.Println("✅ Koneksi ke DB berhasil")
}
