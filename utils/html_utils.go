package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ProcessHTMLWithImages: simpan img base64 ke file, ganti src jadi path file
func ProcessHTMLWithImages(html, soalUID string) (string, error) {
	// regex cari <img src="data:image/...">
	re := regexp.MustCompile(`<img[^>]+src="data:image/[^;]+;base64,([^"]+)"`)

	processedHTML := html
	matches := re.FindAllStringSubmatch(html, -1)

	// folder simpan gambar (misalnya /uploads/soal/soal_uid)
	folderPath := filepath.Join("uploads", "soal", soalUID)
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return "", err
	}

	for i, match := range matches {
		base64Data := match[1]

		// decode base64
		imgData, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			return "", err
		}

		// nama file unik
		filename := fmt.Sprintf("img_%d_%d.png", time.Now().UnixNano(), i)
		filePath := filepath.Join(folderPath, filename)

		// simpan ke file
		if err := os.WriteFile(filePath, imgData, 0644); err != nil {
			return "", err
		}

		// ganti src base64 dengan path file
		processedHTML = strings.Replace(processedHTML, match[0],
			fmt.Sprintf(`<img src="/%s"`, filePath), 1)
	}

	return processedHTML, nil
}
