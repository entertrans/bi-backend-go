package cloudinaryutil

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
)

func DeleteFolder(folderPath string) error {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Hapus semua file dalam folder
	_, err = cld.Admin.DeleteAssetsByPrefix(ctx, admin.DeleteAssetsByPrefixParams{
		Prefix: []string{folderPath},
	})
	if err != nil {
		fmt.Println("Gagal hapus assets dari folder:", err)
		return err
	}

	// Hapus folder kosong
	_, err = cld.Admin.DeleteFolder(ctx, admin.DeleteFolderParams{
		Folder: folderPath,
	})

	if err != nil {
		fmt.Println("Gagal hapus folder kosong:", err)
		// lanjut tanpa return error
	}

	return nil
}
