package gurucontrollers

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"gorm.io/gorm"
)

func GetJawabanRollbackBySiswaNIS(siswaNIS string) ([]TestJawabanResult, error) {
	db := config.DB
	var siswa models.Siswa
	if err := db.Where("siswa_nis = ?", siswaNIS).First(&siswa).Error; err != nil {
		return nil, fmt.Errorf("siswa tidak ditemukan")
	}

	results := []TestJawabanResult{}

	// Ambil test UB dari rb_test
	var testsUB []models.RB_Test
	if err := db.Where("type_test = ?", "ub").
		Preload("Mapel").
		Find(&testsUB).Error; err != nil {
		return nil, err
	}
	for _, test := range testsUB {
		res, err := buildRBTestResult(db, test, siswaNIS)
		if err == nil {
			results = append(results, res)
		}
	}

	// Ambil test TR/Tugas dari rb_peserta
	var pesertaList []models.RB_Peserta
	if err := db.Where("siswa_nis = ?", siswaNIS).
		Preload("Test.Mapel").
		Preload("Kelas").
		Find(&pesertaList).Error; err != nil {
		return nil, err
	}
	for _, peserta := range pesertaList {
		res, err := buildRBTestResult(db, peserta.Test, siswaNIS)
		if err == nil {
			results = append(results, res)
		}
	}

	return results, nil
}

func buildRBTestResult(db *gorm.DB, test models.RB_Test, siswaNIS string) (TestJawabanResult, error) {
	var session models.TestSession
	err := db.Where("test_id = ? AND siswa_nis = ?", test.TestID, siswaNIS).
		Preload("Kelas"). // Preload kelas dari session
		First(&session).Error

	// ===============================
	// Ambil kelas dari session atau peserta
	// ===============================
	var kelasNama string

	// Priority 1: Ambil dari session (kelas saat mengerjakan)
	if session.SessionID != 0 && session.Kelas.KelasId != 0 {
		kelasNama = session.Kelas.KelasNama
	} else {
		// Priority 2: Ambil dari peserta (fallback)
		var peserta models.TO_Peserta
		if err := db.Where("test_id = ? AND siswa_nis = ?", test.TestID, siswaNIS).
			Preload("Kelas").
			First(&peserta).Error; err == nil && peserta.Kelas.KelasId != 0 {
			kelasNama = peserta.Kelas.KelasNama
		}
	}

	// ===============================
	// Hitung rangkaian soal per tipe
	// ===============================
	rangkaian := map[string]int{}
	var butuhReview bool

	// 1. coba ambil dari to_testsoal (untuk tugas / tr)
	var tipeList []string
	db.Model(&models.RB_TestSoal{}).Where("test_id = ?", test.TestID).Pluck("tipe_soal", &tipeList)

	// 2. kalau kosong → ambil dari to_sessionsoal join banksoal (untuk ub)
	if len(tipeList) == 0 && session.SessionID != 0 {
		db.Table("rb_sessionsoal").
			Select("rb_banksoal.tipe_soal").
			Joins("JOIN rb_banksoal ON rb_sessionsoal.soal_id = rb_banksoal.soal_id").
			Where("rb_sessionsoal.session_id = ?", session.SessionID)

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
		Kelas:       kelasNama, // ✅ Tambahkan kelas ke result
		Rangkaian:   rangkaian,
		ButuhReview: butuhReview,
		Reviewed:    false,
		Submited:    false,
		SessionID:   nil, // ✅ Default nil
		Nilai:       nil, // ✅ Default nil
	}

	// ❌ belum ada session → belum dikerjakan
	if err != nil {
		res.Status = "❌ belum dikerjakan"
		return res, nil
	}

	// ✅ Ada session - SET SESSION ID DI SINI
	res.SessionID = &session.SessionID
	res.Tanggal = &session.StartTime

	// Set nilai hanya jika ada
	if session.NilaiAkhir > 0 || session.NilaiAkhir == 0 { // termasuk nilai 0
		nilai := session.NilaiAkhir
		res.Nilai = &nilai
	}

	// ===============================
	// LOGIKA SUBMITED YANG DIPERBAIKI
	// ===============================

	// Cek apakah sudah submit dengan beberapa cara:

	// 1. Cek di to_jawabanfinal (untuk soal yang menggunakan jawaban final)
	var countJawabanFinal int64
	db.Table("to_jawabanfinal").
		Where("session_id = ?", session.SessionID).
		Count(&countJawabanFinal)

	// 2. Cek apakah session sudah completed (berdasarkan EndTime)
	sessionCompleted := session.EndTime != nil

	// 3. Cek apakah sudah ada nilai (termasuk nilai 0)
	hasNilai := session.NilaiAkhir >= 0 // termasuk 0

	// Submit dianggap true jika:
	// - Ada jawaban final ATAU
	// - Session sudah completed (ada EndTime) ATAU
	// - Sudah ada nilai (termasuk nilai 0)
	res.Submited = countJawabanFinal > 0 || sessionCompleted || hasNilai

	// ===============================
	// LOGIKA STATUS DAN REVIEW
	// ===============================
	if butuhReview {
		// hitung soal subjektif di jawabanfinal
		var totalSubjektif, countBelum int64
		db.Table("to_jawabanfinal").
			Joins("JOIN to_banksoal ON to_jawabanfinal.soal_id = to_banksoal.soal_id").
			Where("to_jawabanfinal.session_id = ? AND to_banksoal.tipe_soal IN (?)",
				session.SessionID, []string{"uraian", "isian_singkat"}).
			Count(&totalSubjektif)

		if totalSubjektif > 0 {
			db.Table("to_jawabanfinal").
				Joins("JOIN to_banksoal ON to_jawabanfinal.soal_id = to_banksoal.soal_id").
				Where("to_jawabanfinal.session_id = ? AND to_banksoal.tipe_soal IN (?)",
					session.SessionID, []string{"uraian", "isian_singkat"}).
				Where("to_jawabanfinal.skor_uraian IS NULL").
				Count(&countBelum)

			if countBelum > 0 {
				res.Status = "⏳ menunggu review"
				res.Reviewed = false
			} else {
				res.Status = "✅ sudah direview"
				res.Reviewed = true
			}
		} else {
			// Untuk test yang butuh review tapi tidak ada soal subjektif di jawabanfinal
			if res.Submited {
				res.Status = "✅ sudah dikerjakan"
			} else {
				res.Status = "❌ belum dikerjakan"
			}
		}
	} else {
		if res.Submited {
			res.Status = "✅ sudah dikerjakan"
			res.Reviewed = true
		} else {
			res.Status = "❌ belum dikerjakan"
		}
	}

	log.Printf("[DEBUG] test_id=%d siswa_nis=%s session_id=%v kelas=%s submited=%t nilai=%v status=%s",
		res.TestID, siswaNIS, res.SessionID, res.Kelas, res.Submited, res.Nilai, res.Status)

	return res, nil
}

func FetchRBJawabanBySession(sessionID int, jenis string) (map[string]interface{}, error) {
	// --- 1. Ambil data session
	var session struct {
		SessionID  int     `gorm:"column:session_id"`
		TestID     int     `gorm:"column:test_id"`
		SiswaNIS   string  `gorm:"column:siswa_nis"`
		NilaiAkhir float64 `gorm:"column:nilai_akhir"`
	}
	if err := config.DB.Table("rb_testsession").
		Select("session_id, test_id, siswa_nis, nilai_akhir").
		Where("session_id = ?", sessionID).
		Take(&session).Error; err != nil {
		return nil, err
	}

	// --- 2. Ambil data test
	var test struct {
		Judul string `gorm:"column:judul"`
		Mapel string `gorm:"column:mapel"`
	}
	config.DB.Table("rb_test t").
		Select("t.judul, m.nm_mapel AS mapel").
		Joins("LEFT JOIN tbl_mapel m ON m.kd_mapel = t.mapel_id").
		Where("t.test_id = ?", session.TestID).
		Take(&test)

	// --- 3. Ambil data siswa
	var siswa struct {
		NIS  string `gorm:"column:nis"`
		Nama string `gorm:"column:nama"`
	}
	config.DB.Table("tbl_siswa").
		Select("siswa_nis AS nis, siswa_nama AS nama").
		Where("siswa_nis = ?", session.SiswaNIS).
		Take(&siswa)

	// --- 4. Ambil jawaban sesuai jenis
	var jawaban []models.JawabanResponse

	if jenis == "ub" {
		log.Printf("[RB] ambil UB dari rb_banksoal")
		config.DB.Table("rb_jawabanfinal j").
			Select(`j.soal_id, b.pertanyaan, b.tipe_soal,
				CAST(j.jawaban_siswa AS CHAR) AS jawaban_siswa,
				CAST(b.jawaban_benar AS CHAR) AS jawaban_benar,
				CAST(b.pilihan_jawaban AS CHAR) AS pilihan_jawaban,
				j.skor_objektif, j.skor_uraian,
				b.bobot AS max_score,
				l.nama_file AS lampiran_nama_file,
				l.tipe_file AS lampiran_tipe_file,
				l.path_file AS lampiran_path_file`).
			Joins("JOIN rb_banksoal b ON b.soal_id = j.soal_id").
			Joins("LEFT JOIN rb_lampiran l ON l.lampiran_id = b.lampiran_id").
			Where("j.session_id = ?", sessionID).
			Scan(&jawaban)
	} else {
		log.Printf("[RB] ambil TR/Tugas dari rb_testsoal")
		config.DB.Table("rb_jawabanfinal j").
			Select(`j.soal_id, t.pertanyaan, t.tipe_soal,
				CAST(j.jawaban_siswa AS CHAR) AS jawaban_siswa,
				CAST(t.jawaban_benar AS CHAR) AS jawaban_benar,
				CAST(t.pilihan_jawaban AS CHAR) AS pilihan_jawaban,
				j.skor_objektif, j.skor_uraian,
				t.bobot AS max_score,
				l.nama_file AS lampiran_nama_file,
				l.tipe_file AS lampiran_tipe_file,
				l.path_file AS lampiran_path_file`).
			Joins("JOIN rb_testsoal t ON t.testsoal_id = j.soal_id").
			Joins("LEFT JOIN rb_lampiran l ON l.lampiran_id = t.lampiran_id").
			Where("j.session_id = ?", sessionID).
			Scan(&jawaban)
	}

	// --- 5. Convert ke map (sama seperti versi aktif)
	jawabanMaps := make([]map[string]interface{}, len(jawaban))
	for i, j := range jawaban {
		jawabanMap := map[string]interface{}{
			"soal_id":       j.SoalID,
			"pertanyaan":    j.Pertanyaan,
			"tipe_soal":     j.TipeSoal,
			"jawaban_siswa": j.JawabanSiswa,
			"skor_objektif": j.SkorObjektif,
			"max_score":     j.MaxScore,
			"skor":          j.Score,
		}

		if j.JawabanBenar != nil {
			jawabanMap["jawaban_benar"] = *j.JawabanBenar
		} else {
			jawabanMap["jawaban_benar"] = ""
		}

		if j.PilihanJawaban != nil {
			jawabanMap["pilihan_jawaban"] = *j.PilihanJawaban
		} else {
			jawabanMap["pilihan_jawaban"] = ""
		}

		if j.SkorUraian != nil {
			jawabanMap["skor_uraian"] = *j.SkorUraian
		} else {
			jawabanMap["skor_uraian"] = nil
		}

		if j.LampiranNamaFile != nil {
			jawabanMap["lampiran_nama_file"] = *j.LampiranNamaFile
		}
		if j.LampiranTipeFile != nil {
			jawabanMap["lampiran_tipe_file"] = *j.LampiranTipeFile
		}
		if j.LampiranPathFile != nil {
			jawabanMap["lampiran_path_file"] = *j.LampiranPathFile
		}

		if j.TipeSoal == "matching" {
			currentJawabanBenar := jawabanMap["jawaban_benar"].(string)
			currentPilihanJawaban := jawabanMap["pilihan_jawaban"].(string)
			if currentJawabanBenar == "" || currentJawabanBenar == "[]" {
				generatedAnswer := generateMatchingAnswerKey(currentPilihanJawaban)
				jawabanMap["jawaban_benar"] = generatedAnswer
			}
		}
		jawabanMaps[i] = jawabanMap
	}

	// --- 6. Response
	response := map[string]interface{}{
		"session_id": session.SessionID,
		"test": map[string]interface{}{
			"judul": test.Judul,
			"mapel": test.Mapel,
			"nilai": session.NilaiAkhir,
		},
		"siswa": map[string]interface{}{
			"nis":  siswa.NIS,
			"nama": siswa.Nama,
		},
		"jawaban": jawabanMaps,
	}

	return response, nil
}

// rollback data
// Daftar tabel rollback
var rollbackTables = []string{
	"rb_banksoal",
	"rb_jawabanfinal",
	"rb_lampiran",
	"rb_peserta",
	"rb_sessionsoal",
	"rb_test",
	"rb_testsession",
	"rb_testsoal",
	"rb_test_soal",
}

// ✅ Ambil status data di semua tabel rollback
func GetRollbackStatus() (map[string]bool, error) {
	status := make(map[string]bool)

	for _, table := range rollbackTables {
		var count int64
		err := config.DB.Table(table).Count(&count).Error
		if err != nil {
			status[table] = false
		} else {
			status[table] = count > 0
		}
	}

	return status, nil
}

// ✅ Import file SQL ke tabel rollback
func ImportRollbackSQL(table string, filePath string) error {
	sqlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("gagal membaca file SQL: %v", err)
	}

	sqlContent := string(sqlBytes)

	// Parser aman: pisahkan query per statement
	statements := parseSQLStatements(sqlContent)

	db, err := config.DB.DB()
	if err != nil {
		return fmt.Errorf("gagal konek ke database: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", err)
	}

	for _, stmt := range statements {
		query := strings.TrimSpace(stmt)
		if query == "" {
			continue
		}

		// Jalankan query
		if _, err := tx.Exec(query); err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal menjalankan query: %v\nQuery: %s", err, query)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("gagal commit transaksi: %v", err)
	}

	return nil
}

// ✅ Parser SQL yang memperhitungkan string, escape, dan komentar
func parseSQLStatements(sqlContent string) []string {
	var statements []string
	var current strings.Builder
	inString := false
	stringChar := rune(0)
	escaped := false
	commentLine := false
	commentBlock := false

	reader := bufio.NewReader(strings.NewReader(sqlContent))

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			break
		}

		// Cek komentar single-line (-- atau #)
		if !inString && !commentBlock {
			if r == '-' {
				next, _ := reader.Peek(1)
				if len(next) == 1 && next[0] == '-' {
					commentLine = true
				}
			} else if r == '#' {
				commentLine = true
			}
		}

		// Komentar block /* ... */
		if !inString && !commentLine {
			if r == '/' {
				next, _ := reader.Peek(1)
				if len(next) == 1 && next[0] == '*' {
					commentBlock = true
					reader.ReadRune() // lewati '*'
					continue
				}
			}
		}

		// Tutup komentar block
		if commentBlock {
			if r == '*' {
				next, _ := reader.Peek(1)
				if len(next) == 1 && next[0] == '/' {
					commentBlock = false
					reader.ReadRune()
				}
			}
			continue
		}

		// Tutup komentar baris
		if commentLine {
			if r == '\n' {
				commentLine = false
			}
			continue
		}

		// Escape handling
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}

		// Masuk / keluar string
		if r == '\\' && inString {
			escaped = true
			current.WriteRune(r)
			continue
		}

		if (r == '\'' || r == '"') && !escaped {
			if !inString {
				inString = true
				stringChar = r
			} else if r == stringChar {
				inString = false
			}
			current.WriteRune(r)
			continue
		}

		// Split query hanya kalau ; di luar string
		if r == ';' && !inString {
			statements = append(statements, current.String())
			current.Reset()
			continue
		}

		current.WriteRune(r)
	}

	// Tambahkan sisa query terakhir
	if strings.TrimSpace(current.String()) != "" {
		statements = append(statements, current.String())
	}

	return statements
}

// ✅ Reset semua tabel rollback
func ResetAllRollbackTables() error {
	db, err := config.DB.DB()
	if err != nil {
		return fmt.Errorf("gagal konek ke database: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", err)
	}

	for _, table := range rollbackTables {
		query := fmt.Sprintf("DELETE FROM %s;", table)
		if _, err := tx.Exec(query); err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal menghapus data tabel %s: %v", table, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("gagal commit transaksi: %v", err)
	}

	return nil
}

// end of rollbackdata
