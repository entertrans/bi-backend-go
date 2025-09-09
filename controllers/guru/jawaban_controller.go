package gurucontrollers

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"gorm.io/gorm"
)

type TestJawabanResult struct {
	TestID      uint           `json:"test_id"`
	Jenis       string         `json:"jenis"`
	Mapel       string         `json:"mapel"`
	Judul       string         `json:"judul"`
	Nilai       *float64       `json:"nilai"`
	Tanggal     *time.Time     `json:"tanggal"`
	Status      string         `json:"status"`
	Rangkaian   map[string]int `json:"rangkaian"`
	Submited    bool           `json:"submited"`     // sudah submit jawaban?
	ButuhReview bool           `json:"butuh_review"` // ada esai/isian singkat?
	Reviewed    bool           `json:"reviewed"`     // sudah direview guru?
	SessionID   *uint          `json:"session_id"`
}

// GetJawabanBySiswaNIS ambil semua hasil test berdasarkan NIS siswa
func GetJawabanBySiswaNIS(siswaNIS string) ([]TestJawabanResult, error) {
	db := config.DB

	// 1. Ambil data siswa
	var siswa models.Siswa
	if err := db.Where("siswa_nis = ?", siswaNIS).
		Preload("Kelas").
		First(&siswa).Error; err != nil {
		return nil, fmt.Errorf("siswa tidak ditemukan")
	}

	if siswa.SiswaKelasID == nil || siswa.Kelas.KelasId == 0 {
		return nil, fmt.Errorf("siswa belum memiliki kelas atau data kelas tidak valid")
	}

	results := []TestJawabanResult{}

	// ---------------------------
	// 2A. Ambil test UB by kelas
	// ---------------------------
	var testsUB []models.TO_Test
	if err := db.Where("kelas_id = ? AND type_test = ?", *siswa.SiswaKelasID, "ub").
		Preload("Mapel").
		Find(&testsUB).Error; err != nil {
		return nil, err
	}

	for _, test := range testsUB {
		res, err := buildTestResult(db, test, siswaNIS)
		if err == nil {
			results = append(results, res)
		}
	}

	// ---------------------------
	// 2B. Ambil test TR/Tugas by peserta
	// ---------------------------
	var pesertaList []models.TO_Peserta
	if err := db.Where("siswa_nis = ?", siswaNIS).
		Preload("Test.Mapel").
		Find(&pesertaList).Error; err != nil {
		return nil, err
	}

	for _, peserta := range pesertaList {
		res, err := buildTestResult(db, peserta.Test, siswaNIS)
		if err == nil {
			results = append(results, res)
		}
	}

	return results, nil
}

// Helper untuk bangun hasil per test
func buildTestResult(db *gorm.DB, test models.TO_Test, siswaNIS string) (TestJawabanResult, error) {
	var session models.TestSession
	err := db.Where("test_id = ? AND siswa_nis = ?", test.TestID, siswaNIS).
		First(&session).Error

	// ===============================
	// Hitung rangkaian soal per tipe
	// ===============================
	rangkaian := map[string]int{}
	var butuhReview bool

	// 1. coba ambil dari to_testsoal (untuk tugas / tr)
	var tipeList []string
	db.Model(&models.TO_TestSoal{}).
		Where("test_id = ?", test.TestID).
		Pluck("tipe_soal", &tipeList)

	// 2. kalau kosong → ambil dari to_sessionsoal join banksoal (untuk ub)
	if len(tipeList) == 0 {
		db.Table("to_sessionsoal").
			Select("to_banksoal.tipe_soal").
			Joins("JOIN to_banksoal ON to_sessionsoal.soal_id = to_banksoal.soal_id").
			Where("to_sessionsoal.session_id = ?", session.SessionID).
			Pluck("to_banksoal.tipe_soal", &tipeList)
	}

	for _, tipe := range tipeList {
		rangkaian[tipe]++
		if tipe == "uraian" || tipe == "isian_singkat" {
			butuhReview = true
		}
	}

	// ===============================
	// Build hasil
	// ===============================
	res := TestJawabanResult{
		TestID:      test.TestID,
		Jenis:       test.TypeTest,
		Mapel:       test.Mapel.NmMapel,
		Judul:       test.Judul,
		Rangkaian:   rangkaian,
		ButuhReview: butuhReview,
		Reviewed:    false,
		Submited:    false,
		SessionID:   nil, // ✅ Default nil
	}

	// ❌ belum ada session → belum dikerjakan
	if err != nil {
		res.Status = "❌ belum dikerjakan"
		return res, nil
	}

	// ✅ Ada session - SET SESSION ID DI SINI
	res.SessionID = &session.SessionID // ✅ TAMBAHKAN INI
	res.Nilai = &session.NilaiAkhir
	res.Tanggal = &session.StartTime

	// cek apakah sudah ada jawaban final
	var countJawaban int64
	db.Table("to_jawabanfinal").
		Where("session_id = ?", session.SessionID).
		Count(&countJawaban)
	if countJawaban > 0 {
		res.Submited = true
	}

	if butuhReview {
		// hitung soal subjektif di jawabanfinal
		var totalSubjektif, countBelum int64
		db.Table("to_jawabanfinal").
			Joins("JOIN to_banksoal ON to_jawabanfinal.soal_id = to_banksoal.soal_id").
			Where("to_jawabanfinal.session_id = ? AND to_banksoal.tipe_soal IN (?)",
				session.SessionID, []string{"uraian", "isian_singkat"}).
			Count(&totalSubjektif)

		db.Table("to_jawabanfinal").
			Joins("JOIN to_banksoal ON to_jawabanfinal.soal_id = to_banksoal.soal_id").
			Where("to_jawabanfinal.session_id = ? AND to_banksoal.tipe_soal IN (?)",
				session.SessionID, []string{"uraian", "isian_singkat"}).
			Where("to_jawabanfinal.skor_uraian IS NULL").
			Count(&countBelum)

		if totalSubjektif > 0 {
			if countBelum > 0 {
				res.Status = "⏳ menunggu review"
				res.Reviewed = false
			} else {
				res.Status = "✅ sudah direview"
				res.Reviewed = true
			}
		}
	} else {
		if res.Submited {
			res.Status = "✅ otomatis dinilai"
			res.Reviewed = true
		}
	}

	log.Printf("[DEBUG] test_id=%d siswa_nis=%s session_id=%v rangkaian=%v status=%s",
		res.TestID, siswaNIS, res.SessionID, res.Rangkaian, res.Status)

	return res, nil
}

// Fungsi tambahan untuk mendapatkan detail jawaban siswa
func GetDetailJawabanBySession(sessionID uint) (map[string]interface{}, error) {
	db := config.DB

	var session models.TestSession
	if err := db.Preload("Test").
		Preload("Test.Mapel").
		Preload("Siswa").
		First(&session, sessionID).Error; err != nil {
		return nil, fmt.Errorf("session tidak ditemukan")
	}

	// Ambil semua jawaban final
	var jawabanFinal []models.TO_JawabanFinal
	if err := db.Where("session_id = ?", sessionID).
		Preload("Soal").
		Preload("Soal.Lampiran").
		Find(&jawabanFinal).Error; err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"session":       session,
		"jawaban":       jawabanFinal,
		"total_soal":    len(jawabanFinal),
		"nilai_akhir":   session.NilaiAkhir,
		"status":        session.Status,
		"waktu_selesai": session.EndTime,
	}

	return result, nil
}

// Fungsi untuk mendapatkan semua siswa yang mengerjakan test tertentu
func GetSiswaByTest(testID uint) ([]map[string]interface{}, error) {
	db := config.DB

	// Ambil data test
	var test models.TO_Test
	if err := db.Preload("Mapel").Preload("Kelas").First(&test, testID).Error; err != nil {
		return nil, fmt.Errorf("test tidak ditemukan")
	}

	// Ambil semua peserta test
	var peserta []models.TO_Peserta
	if err := db.Where("test_id = ?", testID).
		Preload("Siswa").
		Preload("Siswa.Kelas").
		Find(&peserta).Error; err != nil {
		return nil, err
	}

	results := []map[string]interface{}{}

	for _, p := range peserta {
		// Cek apakah siswa memiliki session
		var session models.TestSession
		sessionErr := db.Where("test_id = ? AND siswa_nis = ?", testID, p.SiswaNIS).
			First(&session).Error

		studentResult := map[string]interface{}{
			"peserta_id":  p.PesertaID,
			"siswa_nis":   p.SiswaNIS,
			"siswa_nama":  p.Siswa.SiswaNama,
			"kelas":       p.Siswa.Kelas.KelasNama,
			"status":      p.Status,
			"nilai_akhir": p.NilaiAkhir,
			"created_at":  p.CreatedAt,
		}

		if sessionErr == nil {
			// Ada session, tambahkan info session
			studentResult["session_id"] = session.SessionID
			studentResult["start_time"] = session.StartTime
			studentResult["end_time"] = session.EndTime
			studentResult["session_status"] = session.Status
		} else {
			// Tidak ada session
			studentResult["session_id"] = nil
			studentResult["session_status"] = "belum_mulai"
		}

		results = append(results, studentResult)
	}

	return results, nil
}

// Fungsi untuk mendapatkan statistik test
func GetTestStatistics(testID uint) (map[string]interface{}, error) {
	db := config.DB

	var totalPeserta int64
	var sudahMengerjakan int64
	var belumMengerjakan int64
	var sudahDinilai int64
	var belumDinilai int64

	// Hitung total peserta
	db.Model(&models.TO_Peserta{}).Where("test_id = ?", testID).Count(&totalPeserta)

	// Hitung yang sudah mengerjakan (punya session)
	db.Model(&models.TestSession{}).Where("test_id = ?", testID).Count(&sudahMengerjakan)
	belumMengerjakan = totalPeserta - sudahMengerjakan

	// Hitung yang sudah dinilai (nilai_akhir > 0 atau status graded)
	db.Model(&models.TestSession{}).
		Where("test_id = ? AND status = ? AND nilai_akhir > 0", testID, "graded").
		Count(&sudahDinilai)

	belumDinilai = sudahMengerjakan - sudahDinilai

	// Ambil rata-rata nilai
	var avgNilai float64
	db.Model(&models.TestSession{}).
		Where("test_id = ? AND status = ?", testID, "graded").
		Select("COALESCE(AVG(nilai_akhir), 0)").
		Scan(&avgNilai)

	// Ambil nilai tertinggi dan terendah
	var maxNilai, minNilai float64
	db.Model(&models.TestSession{}).
		Where("test_id = ? AND status = ?", testID, "graded").
		Select("COALESCE(MAX(nilai_akhir), 0)").
		Scan(&maxNilai)

	db.Model(&models.TestSession{}).
		Where("test_id = ? AND status = ?", testID, "graded").
		Select("COALESCE(MIN(nilai_akhir), 0)").
		Scan(&minNilai)

	result := map[string]interface{}{
		"test_id":           testID,
		"total_peserta":     totalPeserta,
		"sudah_mengerjakan": sudahMengerjakan,
		"belum_mengerjakan": belumMengerjakan,
		"sudah_dinilai":     sudahDinilai,
		"belum_dinilai":     belumDinilai,
		"rata_rata_nilai":   fmt.Sprintf("%.2f", avgNilai),
		"nilai_tertinggi":   fmt.Sprintf("%.2f", maxNilai),
		"nilai_terendah":    fmt.Sprintf("%.2f", minNilai),
	}

	return result, nil
}

func GetSiswaDetailForGuru(siswaNIS string) (map[string]interface{}, error) {
	db := config.DB

	var siswa models.Siswa
	if err := db.Where("siswa_nis = ?", siswaNIS).
		Preload("Kelas").
		First(&siswa).Error; err != nil {
		return nil, fmt.Errorf("siswa tidak ditemukan")
	}

	result := map[string]interface{}{
		"siswa_id":   siswa.SiswaID,
		"siswa_nis":  *siswa.SiswaNIS, // Dereference pointer
		"siswa_nama": *siswa.SiswaNama,
		"kelas_nama": siswa.Kelas.KelasNama,
	}

	return result, nil
}
func FetchJawabanBySession(sessionID int) (map[string]interface{}, error) {
	// --- 1. Ambil data session (tambah nilai_akhir)
	var session struct {
		SessionID  int     `gorm:"column:session_id"`
		TestID     int     `gorm:"column:test_id"`
		SiswaNIS   string  `gorm:"column:siswa_nis"`
		NilaiAkhir float64 `gorm:"column:nilai_akhir"`
	}
	if err := config.DB.Table("to_testsession").
		Select("session_id, test_id, siswa_nis, nilai_akhir").
		Where("session_id = ?", sessionID).
		Take(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("session %d not found", sessionID)
		}
		return nil, err
	}

	// --- 2. Ambil data test (join ke tbl_mapel untuk nama mapel)
	var test struct {
		Judul string `gorm:"column:judul"`
		Mapel string `gorm:"column:mapel"`
	}
	if err := config.DB.Table("to_test t").
		Select("t.judul, m.nm_mapel AS mapel").
		Joins("LEFT JOIN tbl_mapel m ON m.kd_mapel = t.mapel_id").
		Where("t.test_id = ?", session.TestID).
		Take(&test).Error; err != nil {
		// kalau tidak ditemukan test, tetap lanjut tapi beri info
		if errors.Is(err, gorm.ErrRecordNotFound) {
			test.Judul = ""
			test.Mapel = ""
		} else {
			return nil, err
		}
	}

	// --- 3. Ambil data siswa
	var siswa struct {
		NIS  string `gorm:"column:nis"`
		Nama string `gorm:"column:nama"`
	}
	if err := config.DB.Table("tbl_siswa").
		Select("siswa_nis AS nis, siswa_nama AS nama").
		Where("siswa_nis = ?", session.SiswaNIS).
		Take(&siswa).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			siswa.NIS = session.SiswaNIS
			siswa.Nama = ""
		} else {
			return nil, err
		}
	}

	// --- 4. Ambil jawaban dari banksoal (UB)
	var ubJawaban []models.JawabanResponse
	if err := config.DB.Table("to_jawabanfinal j").
		Select(`j.soal_id, b.pertanyaan, b.tipe_soal,
        CAST(j.jawaban_siswa AS CHAR) AS jawaban_siswa,
        CAST(b.jawaban_benar AS CHAR) AS jawaban_benar,
        CAST(b.pilihan_jawaban AS CHAR) AS pilihan_jawaban,
        j.skor_objektif, j.skor_uraian,
        b.bobot AS max_score,
        l.nama_file AS lampiran_nama_file,
        l.tipe_file AS lampiran_tipe_file,
        l.path_file AS lampiran_path_file
    `).
		Joins("JOIN to_banksoal b ON b.soal_id = j.soal_id").
		Joins("LEFT JOIN TO_Lampiran l ON l.lampiran_id = b.lampiran_id").
		Where("j.session_id = ?", sessionID).
		Scan(&ubJawaban).Error; err != nil {
		return nil, err
	}

	// --- 5. Ambil jawaban dari testsoal (selain UB)
	var testJawaban []models.JawabanResponse
	if err := config.DB.Table("to_jawabanfinal j").
		Select(`j.soal_id, t.pertanyaan, t.tipe_soal,
        CAST(j.jawaban_siswa AS CHAR) AS jawaban_siswa,
        CAST(t.jawaban_benar AS CHAR) AS jawaban_benar,
        j.skor_objektif, j.skor_uraian,
        t.bobot AS max_score,
        l.nama_file AS lampiran_nama_file,
        l.tipe_file AS lampiran_tipe_file,
        l.path_file AS lampiran_path_file
    `).
		Joins("JOIN to_testsoal t ON t.testsoal_id = j.soal_id").
		Joins("LEFT JOIN TO_Lampiran l ON l.lampiran_id = t.lampiran_id").
		Where("j.session_id = ?", sessionID).
		Scan(&testJawaban).Error; err != nil {
		return nil, err
	}

	// --- 6. Gabung jawaban (UB + testsoal)
	jawaban := append(ubJawaban, testJawaban...)

	// --- 7. Bentuk response final (gunakan nilai_akhir dari session)
	response := map[string]interface{}{
		"session_id": session.SessionID,
		"test": map[string]interface{}{
			"judul": test.Judul,
			"mapel": test.Mapel,
			"nilai": session.NilaiAkhir, // DIUBAH: gunakan nilai_akhir dari session
		},
		"siswa": map[string]interface{}{
			"nis":  siswa.NIS,
			"nama": siswa.Nama,
		},
		"jawaban": jawaban,
	}

	return response, nil
}
func UpdateNilaiJawaban(sessionID int, perubahan []struct {
	SessionID int     `json:"session_id"`
	SoalID    uint    `json:"soal_id"`
	Nilai     float64 `json:"nilai"`
}) error {
	db := config.DB

	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, p := range perubahan {
		if p.SessionID != sessionID {
			tx.Rollback()
			return fmt.Errorf("session_id tidak match")
		}

		// Update menggunakan model
		var jawaban models.JawabanFinal
		if err := tx.Where("session_id = ? AND soal_id = ?", sessionID, p.SoalID).
			First(&jawaban).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("jawaban tidak ditemukan untuk soal_id %d", p.SoalID)
		}

		jawaban.SkorUraian = &p.Nilai
		if err := tx.Save(&jawaban).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Hitung ulang nilai dengan mempertimbangkan max_score
	if err := hitungNilaiAkhirDenganBobot(tx, sessionID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func hitungNilaiAkhirDenganBobot(tx *gorm.DB, sessionID int) error {
	// Ambil semua jawaban dengan informasi bobot/max_score
	var jawaban []struct {
		models.JawabanFinal
		MaxScore float64 `gorm:"column:max_score"`
	}

	// Query yang mengambil max_score dari banksoal atau testsoal
	err := tx.Table("to_jawabanfinal jf").
		Select("jf.*, COALESCE(bs.bobot, ts.bobot, 1.0) as max_score").
		Joins("LEFT JOIN to_banksoal bs ON bs.soal_id = jf.soal_id").
		Joins("LEFT JOIN to_testsoal ts ON ts.testsoal_id = jf.soal_id").
		Where("jf.session_id = ?", sessionID).
		Scan(&jawaban).Error

	if err != nil {
		return err
	}

	var totalSkor float64
	var totalMaxScore float64

	for _, j := range jawaban {
		// Gunakan skor_uraian jika ada, otherwise gunakan skor_objektif
		skor := j.SkorObjektif
		if j.SkorUraian != nil {
			skor = *j.SkorUraian
		}

		// Pastikan skor tidak melebihi max_score
		if skor > j.MaxScore {
			skor = j.MaxScore
		}

		totalSkor += skor
		totalMaxScore += j.MaxScore
	}

	// Hitung nilai akhir (dalam skala 100)
	nilaiAkhir := 0.0
	if totalMaxScore > 0 {
		nilaiAkhir = (totalSkor / totalMaxScore) * 100
	}

	// Pastikan nilai tidak melebihi 100
	if nilaiAkhir > 100 {
		nilaiAkhir = 100
	}

	// Update session
	var session models.TestSession
	if err := tx.First(&session, sessionID).Error; err != nil {
		return err
	}

	session.NilaiAkhir = nilaiAkhir
	session.Status = "graded"
	return tx.Save(&session).Error
}

func UpdateNilaiAkhir(sessionID int, nilaiAkhir float64) error {
	db := config.DB
	log.Printf("[DEBUG] UpdateNilaiAkhir - sessionID: %d, nilaiAkhir: %.2f", sessionID, nilaiAkhir)

	// Update nilai_akhir di tabel to_testsession
	result := db.Table("to_testsession").
		Where("session_id = ?", sessionID).
		Updates(map[string]interface{}{
			"nilai_akhir": nilaiAkhir,
			"status":      "graded",
		})

	if result.Error != nil {
		log.Printf("[ERROR] Gagal update testsession: %v", result.Error)
		return result.Error
	}

	log.Printf("[DEBUG] Testsession updated, rows affected: %d", result.RowsAffected)

	if result.RowsAffected == 0 {
		errMsg := fmt.Sprintf("session tidak ditemukan: %d", sessionID)
		log.Printf("[ERROR] %s", errMsg)
		return fmt.Errorf(errMsg)
	}

	return nil
}

func ResetTestSession(sessionID uint) error {
	db := config.DB

	// Pastikan session ada
	var session models.TestSession
	if err := db.First(&session, "session_id = ?", sessionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("session tidak ditemukan")
		}
		return err
	}

	// Hapus jawaban final
	if err := db.Where("session_id = ?", sessionID).Delete(&models.JawabanFinal{}).Error; err != nil {
		return err
	}

	// Hapus session soal
	if err := db.Where("session_id = ?", sessionID).Delete(&models.TO_SessionSoal{}).Error; err != nil {
		return err
	}

	// Hapus test session
	if err := db.Where("session_id = ?", sessionID).Delete(&models.TestSession{}).Error; err != nil {
		return err
	}

	return nil
}

// func GetSoalByPenilaianID(db *gorm.DB, testID uint) ([]models.TO_Soal, error) {
// 	var soal []models.TO_Soal

// 	err := db.
// 		Preload("PilihanJawaban").
// 		Preload("JawabanBenar").
// 		Preload("Lampiran").
// 		Where("test_id = ?", testID).
// 		Find(&soal).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	return soal, nil
// }
