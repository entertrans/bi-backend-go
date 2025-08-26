package models

import (
	"time"

	"gorm.io/gorm"
)

// =========================
// 1. TO_Test
// =========================
type TO_Test struct {
	TestID      uint      `json:"test_id" gorm:"column:test_id;primaryKey;autoIncrement"`
	GuruID      uint      `json:"guru_id" gorm:"column:guru_id"`
	KelasID     uint      `json:"kelas_id" gorm:"column:kelas_id"`
	MapelID     uint64    `json:"mapel_id"`
	Jumlah      uint      `json:"jumlah_soal_tampil" gorm:"column:jumlah_soal_tampil"`
	TypeTest    string    `json:"type_test"` // "ub" atau "quis"
	Judul       string    `json:"judul" gorm:"column:judul"`
	Deskripsi   string    `json:"deskripsi" gorm:"column:deskripsi"`
	DurasiMenit int       `json:"durasi_menit" gorm:"column:durasi_menit"`
	Deadline    string       `json:"deadline" gorm:"column:deadline"`
	Aktif       uint       `json:"aktif" gorm:"column:aktif"`
	RandomSoal  bool      `json:"random_soal" gorm:"column:random_soal"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`

	// Relasi
	Mapel Mapel `gorm:"foreignKey:MapelID;references:KdMapel" json:"mapel"`
	Guru  Guru  `json:"guru" gorm:"foreignKey:GuruID;references:GuruID"`
	Kelas Kelas `json:"kelas" gorm:"foreignKey:KelasID;references:KelasId"`
}

func (TO_Test) TableName() string {
	return "TO_Test"
}

// =========================
// TO_Peserta
// =========================
type TO_Peserta struct {
	PesertaID  uint      `json:"peserta_id" gorm:"column:peserta_id;primaryKey;autoIncrement"`
	TestID     uint      `json:"test_id" gorm:"column:test_id"`
	SiswaNIS   string    `json:"siswa_nis" gorm:"column:siswa_nis"`
	Status     string    `json:"status" gorm:"column:status;default:not_started"`
	ExtraTime  int       `json:"extra_time" gorm:"column:extra_time"`
	NilaiAwal  float64   `json:"nilai_awal" gorm:"column:nilai_awal"`
	NilaiAkhir float64   `json:"nilai_akhir" gorm:"column:nilai_akhir"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`

	// Relasi
	Test  TO_Test `json:"test" gorm:"foreignKey:TestID;references:TestID"`
	Siswa Siswa   `json:"siswa" gorm:"foreignKey:SiswaNIS;references:SiswaNIS"`
}

func (TO_Peserta) TableName() string {
	return "to_peserta"
}

// =========================
// 2. TO_TestSession
// =========================
type TO_TestSession struct {
	SessionID  uint       `json:"session_id" gorm:"column:session_id;primaryKey;autoIncrement"`
	TestID     uint       `json:"test_id" gorm:"column:test_id"`
	SiswaNIS   uint       `json:"siswa_nis" gorm:"column:siswa_nis"`
	StartTime  time.Time  `json:"start_time" gorm:"column:start_time"`
	EndTime    *time.Time `json:"end_time" gorm:"column:end_time"`
	WaktuSisa  int        `json:"waktu_sisa" gorm:"column:waktu_sisa"`
	Status     string     `json:"status" gorm:"column:status"`
	NilaiAwal  float64    `json:"nilai_awal" gorm:"column:nilai_awal"`
	NilaiAkhir float64    `json:"nilai_akhir" gorm:"column:nilai_akhir"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"column:updated_at"`

	// Relasi
	Test  TO_Test `json:"test" gorm:"foreignKey:TestID;references:TestID"`
	Siswa Siswa   `json:"siswa" gorm:"foreignKey:SiswaNIS;references:SiswaID"`
}

func (TO_TestSession) TableName() string {
	return "TO_TestSession"
}

// =========================
// 3. TO_JawabanDraft
// =========================
type TO_JawabanDraft struct {
	DraftID      uint      `json:"draft_id" gorm:"column:draft_id;primaryKey;autoIncrement"`
	SessionID    uint      `json:"session_id" gorm:"column:session_id"`
	SoalID       uint      `json:"soal_id" gorm:"column:soal_id"`
	JawabanSiswa string    `json:"jawaban_siswa" gorm:"column:jawaban_siswa;type:json"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`

	// Relasi
	Session TO_TestSession `json:"session" gorm:"foreignKey:SessionID;references:SessionID"`
	Soal    TO_BankSoal    `json:"soal" gorm:"foreignKey:SoalID;references:SoalID"`
}

func (TO_JawabanDraft) TableName() string {
	return "TO_JawabanDraft"
}

// =========================
// 4. TO_JawabanFinal
// =========================
type TO_JawabanFinal struct {
	FinalID      uint      `json:"final_id" gorm:"column:final_id;primaryKey;autoIncrement"`
	SessionID    uint      `json:"session_id" gorm:"column:session_id"`
	SoalID       uint      `json:"soal_id" gorm:"column:soal_id"`
	JawabanSiswa string    `json:"jawaban_siswa" gorm:"column:jawaban_siswa;type:json"`
	SkorObjektif float64   `json:"skor_objektif" gorm:"column:skor_objektif"`
	SkorUraian   *float64  `json:"skor_uraian" gorm:"column:skor_uraian"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`

	// Relasi
	Session TO_TestSession `json:"session" gorm:"foreignKey:SessionID;references:SessionID"`
	Soal    TO_BankSoal    `json:"soal" gorm:"foreignKey:SoalID;references:SoalID"`
}

func (TO_JawabanFinal) TableName() string {
	return "TO_JawabanFinal"
}

// =========================
// 5. TO_BankSoal
// =========================
type TO_BankSoal struct {
	SoalID uint `json:"soal_id" gorm:"column:soal_id;primaryKey;autoIncrement"`
	// SoalUID        string         `json:"soal_uid" gorm:"column:soal_uid;unique"`
	GuruID         uint           `json:"guru_id" gorm:"column:guru_id"`
	MapelID        uint           `gorm:"column:mapel_id;not null" json:"mapel_id"`
	TipeSoal       string         `json:"tipe_soal" gorm:"column:tipe_soal"`
	KelasID        uint           `json:"kelas_id" gorm:"column:kelas_id"`
	Pertanyaan     string         `json:"pertanyaan" gorm:"column:pertanyaan"`
	LampiranID     *uint          `json:"lampiran_id" gorm:"column:lampiran_id"` // boleh null
	PilihanJawaban string         `json:"pilihan_jawaban" gorm:"column:pilihan_jawaban;type:json"`
	JawabanBenar   string         `json:"jawaban_benar" gorm:"column:jawaban_benar;type:json"`
	Bobot          float64        `json:"bobot" gorm:"column:bobot"`
	CreatedAt      time.Time      `json:"created_at" gorm:"column:created_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relasi
	Guru     Guru         `json:"guru" gorm:"foreignKey:GuruID;references:GuruID"`
	Kelas    Kelas        `json:"kelas" gorm:"foreignKey:KelasID;references:KelasId"`
	Mapel    Mapel        `gorm:"foreignKey:MapelID;references:KdMapel" json:"mapel"`
	Lampiran *TO_Lampiran `json:"lampiran" gorm:"foreignKey:LampiranID;references:LampiranID"`
}

func (TO_BankSoal) TableName() string {
	return "TO_BankSoal"
}

// =========================
// 6. TO_PenilaianGuru
// =========================
type TO_PenilaianGuru struct {
	PenilaianID uint      `json:"penilaian_id" gorm:"column:penilaian_id;primaryKey;autoIncrement"`
	FinalID     uint      `json:"final_id" gorm:"column:final_id"`
	Skor        float64   `json:"skor" gorm:"column:skor"`
	Komentar    string    `json:"komentar" gorm:"column:komentar"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`

	// Relasi
	JawabanFinal TO_JawabanFinal `json:"jawaban_final" gorm:"foreignKey:FinalID;references:FinalID"`
}

func (TO_PenilaianGuru) TableName() string {
	return "TO_PenilaianGuru"
}

// =========================
// 7. TO_Lampiran (gallery bersama, soft delete)
// =========================
type TO_Lampiran struct {
	LampiranID uint           `json:"lampiran_id" gorm:"column:lampiran_id;primaryKey;autoIncrement"`
	NamaFile   string         `json:"nama_file" gorm:"column:nama_file"`
	PathFile   string         `json:"path_file" gorm:"column:path_file;type:text"` // simpan path/URL file
	TipeFile   string         `json:"tipe_file" gorm:"column:tipe_file"`           // image/pdf/audio/video/other
	Deskripsi  string         `json:"deskripsi" gorm:"column:deskripsi"`
	CreatedAt  time.Time      `json:"created_at" gorm:"column:created_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"` // soft delete (trash)
}

func (TO_Lampiran) TableName() string {
	return "TO_Lampiran"
}

type TO_TestSoal struct {
	TestsoalID     uint           `json:"testsoal_id" gorm:"column:testsoal_id;primaryKey;autoIncrement"`
	TestID         uint           `json:"test_id" gorm:"column:test_id;not null"`
	TipeSoal       string         `json:"tipe_soal" gorm:"column:tipe_soal;type:enum('pg','pg_kompleks','matching','isian_singkat','uraian');not null"`
	Pertanyaan     string         `json:"pertanyaan" gorm:"column:pertanyaan;type:text;not null"`
	LampiranID     *uint          `json:"lampiran_id" gorm:"column:lampiran_id"`
	PilihanJawaban string         `json:"pilihan_jawaban" gorm:"column:pilihan_jawaban;type:json"`
	JawabanBenar   string         `json:"jawaban_benar" gorm:"column:jawaban_benar;type:json;not null"`
	Bobot          float64        `json:"bobot" gorm:"column:bobot;type:decimal(5,2);default:1.00"`
	CreatedAt      time.Time      `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	// Relasi
	Test     TO_Test      `json:"test" gorm:"foreignKey:TestID;references:TestID"`
	Lampiran *TO_Lampiran `json:"lampiran" gorm:"foreignKey:LampiranID;references:LampiranID"`
}

// TableName specifies the table name
func (TO_TestSoal) TableName() string {
	return "to_testsoal"
}