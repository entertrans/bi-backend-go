package controllers

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

func BoolPtr(b bool) *bool {
	return &b
}

func UploadToCloudinary(file multipart.File, fileHeader *multipart.FileHeader, nis string, jenis string) (string, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return "", fmt.Errorf("gagal setup cloudinary: %v", err)
	}

	// Penamaan: lampiran/2026002/profile-picture-2026002.jpg
	// publicID := fmt.Sprintf("lampiran/%s/%s-%s", nis, jenis, nis)

	uploadParams := uploader.UploadParams{
		PublicID:     fmt.Sprintf("%s-%s", jenis, nis), // profil-picture-2027002
		Folder:       fmt.Sprintf("lampiran/%s", nis),  // folder lampiran/2027002
		Overwrite:    api.Bool(true),
		ResourceType: "image",
		Format:       "jpg",
	}

	uploadResult, err := cld.Upload.Upload(context.Background(), file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("gagal upload file: %v", err)
	}

	return uploadResult.SecureURL, nil
}

func CreatePPDBSiswa(c *gin.Context) error {
	nis := c.PostForm("nis")
	nama := c.PostForm("nama")
	nisn := c.PostForm("nisn")
	jenkel := c.PostForm("jenkel")
	satelitStr := c.PostForm("satelit")

	// Validasi input
	if nis == "" || nama == "" || nisn == "" || jenkel == "" || satelitStr == "" {
		return fmt.Errorf("semua field wajib diisi")
	}

	satelitID, err := strconv.Atoi(satelitStr)
	if err != nil {
		return fmt.Errorf("satelit ID tidak valid: %v", err)
	}

	// Foto
	file, err := c.FormFile("photo")
	if err != nil {
		return fmt.Errorf("foto wajib diunggah")
	}

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("gagal membuka file: %v", err)
	}
	defer src.Close()

	photoURL, err := UploadToCloudinary(src, file, nis, "profil-picture")
	if err != nil {
		return fmt.Errorf("gagal upload ke cloudinary: %v", err)
	}

	// Simpan siswa
	siswa := models.Siswa{
		SiswaNIS:    nis,
		SiswaNISN:   nisn,
		SiswaNama:   nama,
		SiswaJenkel: jenkel,
		SatelitID:   satelitID,
		SoftDeleted: 2,
	}
	if err := config.DB.Create(&siswa).Error; err != nil {
		return fmt.Errorf("gagal simpan siswa: %v", err)
	}

	// Simpan lampiran
	lampiran := models.Lampiran{
		SiswaNIS:     nis,
		JenisDokumen: "profil-picture",
		Url:          photoURL,
		Upload:       time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := config.DB.Create(&lampiran).Error; err != nil {
		return fmt.Errorf("gagal simpan lampiran: %v", err)
	}

	return nil
}
