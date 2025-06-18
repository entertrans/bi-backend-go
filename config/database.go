package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"          // Untuk membaca file .env
	"gorm.io/driver/mysql"              // Driver MySQL untuk GORM
	"gorm.io/gorm"                      // ORM utama

)

// DB adalah variabel global untuk koneksi database
var DB *gorm.DB

// LoadEnv akan memuat file .env
func LoadEnv() {
	err := godotenv.Load() // Baca file .env dari root project
	if err != nil {
		log.Fatal("❌ Gagal memuat file .env")
	}
}

// ConnectDB akan membuat koneksi ke database
func ConnectDB() {
	// Panggil fungsi untuk load .env
	LoadEnv()

	// Ambil variabel dari .env
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"), // Username MySQL
		os.Getenv("DB_PASS"), // Password MySQL
		os.Getenv("DB_HOST"), // Host MySQL (127.0.0.1)
		os.Getenv("DB_PORT"), // Port MySQL (3306)
		os.Getenv("DB_NAME"), // Nama database (bi)
	)

	// Buka koneksi menggunakan GORM
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Gagal koneksi ke database:", err)
	}

	// Simpan koneksi ke variabel global DB
	DB = database
	fmt.Println("✅ Berhasil koneksi ke database MySQL")
}
