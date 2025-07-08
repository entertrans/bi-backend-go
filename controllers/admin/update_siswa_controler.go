package admincontrollers

import (
	"errors"
	"fmt"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/entertrans/bi-backend-go/utils/cloudinaryutil"
)

func UpdateSiswaField(nis string, field string, value interface{}) error {
	// Validasi field yang boleh diubah
	allowedFields := map[string]bool{
		"siswa_nama":            true,
		"siswa_alamat":          true,
		"siswa_tempat":          true,
		"siswa_tgl_lahir":       true,
		"siswa_jenkel":          true,
		"siswa_email":           true,
		"siswa_no_telp":         true,
		"siswa_kewarganegaraan": true,
		"siswa_kelas_id":        true,
		"siswa_satelit":         true,
		"siswa_nisn":            true,
		"nik_siswa":             true,
		"anak_ke":               true,
		"no_ijazah":             true,
		"sekolah_asal":          true,
		"agama_nama":            true,
		"ayah_nama":             true,
		"ayah_nik":              true,
		"ayah_tempat":           true,
		"ayah_tanggal":          true,
		"ayah_pekerjaan":        true,
		"no_telp_ayah":          true,
		"email_ayah":            true,
		"ibu_nama":              true,
		"ibu_nik":               true,
		"ibu_tempat":            true,
		"ibu_tanggal":           true,
		"ibu_pekerjaan":         true,
		"no_telp_ibu":           true,
		"email_ibu":             true,
		"wali_nama":             true,
		"wali_nik":              true,
		"wali_tempat":           true,
		"wali_tanggal":          true,
		"wali_pekerjaan":        true,
		"wali_notelp":           true,
		"wali_alamat":           true,
	}

	if !allowedFields[field] {
		return errors.New("field tidak diizinkan")
	}

	return config.DB.Model(&models.Siswa{}).Where("siswa_nis = ?", nis).Update(field, value).Error
}

func BatalkanSiswaByNIS(nis string) error {
	// Hapus folder cloudinary: lampiran/{nis}
	err := cloudinaryutil.DeleteFolder("lampiran/" + nis)
	if err != nil {
		fmt.Println("Gagal hapus folder Cloudinary:", err)
		// Tetap lanjut hapus DB, tapi kamu bisa tambahkan fallback kalau perlu
	}

	// Hapus dari tabel lampiran
	if err := config.DB.Where("siswa_nis = ?", nis).Delete(&models.Lampiran{}).Error; err != nil {
		return err
	}

	// Hapus dari tabel siswa
	if err := config.DB.Where("siswa_nis = ?", nis).Delete(&models.Siswa{}).Error; err != nil {
		return err
	}

	return nil
}
