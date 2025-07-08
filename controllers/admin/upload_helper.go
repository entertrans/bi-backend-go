package admincontrollers

// import (
// 	"context"
// 	"mime/multipart"

// 	"github.com/cloudinary/cloudinary-go/v2"
// 	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
// )

// func UploadToCloudinary(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
// 	cld, err := cloudinary.NewFromParams(
// 		"dalcsrtd9",                   // ganti dengan cloud name kamu
// 		"786546378835215",             // ganti dengan API Key kamu
// 		"_AA5U8ky3lk1pT0o4cfnQ3saf3c", // ganti dengan API Secret kamu
// 	)

// 	if err != nil {
// 		return "", err
// 	}

// 	ctx := context.Background()
// 	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
// 		Folder:   "foto_siswa",
// 		PublicID: fileHeader.Filename,
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	return uploadResult.SecureURL, nil
// }
