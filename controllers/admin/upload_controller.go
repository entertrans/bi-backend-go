package admincontrollers

import (
	"fmt"
	"io"

	"github.com/entertrans/bi-backend-go/utils/cloudinaryutil"
)

func UploadDokumenController(nis string, dokumenJenis string, file io.Reader, filename string) (string, error) {
	publicId := fmt.Sprintf("%s-%s", dokumenJenis, nis)
	folder := fmt.Sprintf("lampiran/%s", nis)

	url, err := cloudinaryutil.UploadToCloudinary(file, filename, publicId, folder)
	if err != nil {
		return "", err
	}
	return url, nil
}
