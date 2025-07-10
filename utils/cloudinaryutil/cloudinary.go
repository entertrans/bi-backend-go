package cloudinaryutil

import (
	"context"
	"io"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func Bool(v bool) *bool {
	return &v
}

func UploadToCloudinary(file io.Reader, filename, publicId, folder string) (string, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return "", err
	}

	uploadParams := uploader.UploadParams{
		PublicID:     publicId, // Hanya nama file
		Folder:       folder,   // Ini path folder, supaya muncul di dashboard
		Overwrite:    Bool(true),
		ResourceType: "image",
		Format:       "jpg",
	}

	result, err := cld.Upload.Upload(context.Background(), file, uploadParams)
	if err != nil {
		return "", err
	}

	return result.SecureURL, nil
}
