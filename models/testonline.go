package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TO_Test struct {
	TestID      uint       `json:"test_id" gorm:"column:test_id;primaryKey;autoIncrement"`
	GuruID      uint       `json:"guru_id" gorm:"column:guru_id"`
	KelasID     uint       `json:"kelas_id" gorm:"column:kelas_id"`
	MapelID     uint64     `json:"mapel_id"`
	Jumlah      uint       `json:"jumlah_soal_tampil" gorm:"column:jumlah_soal_tampil"`
	TypeTest    string     `json:"type_test"` // "ub" atau "tr / tugas"
	Judul       string     `json:"judul" gorm:"column:judul"`
	Deskripsi   string     `json:"deskripsi" gorm:"column:deskripsi"`
	DurasiMenit int        `json:"durasi_menit" gorm:"column:durasi_menit"`
	Deadline    *time.Time `json:"deadline" gorm:"column:deadline"`
	Aktif       *uint      `json:"aktif" gorm:"column:aktif"`
	RandomSoal  bool       `json:"random_soal" gorm:"column:random_soal"`
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at"`

	// Relasi
	Mapel   Mapel  `gorm:"foreignKey:MapelID;references:KdMapel" json:"mapel"`
	Guru    Guru   `json:"guru" gorm:"foreignKey:GuruID;references:GuruID"`
	Kelas   Kelas  `json:"kelas" gorm:"foreignKey:KelasID;references:KelasId"`
	SoalIDs []uint `json:"soal_ids" gorm:"-"`
}

func (TO_Test) TableName() string {
	return "TO_Test"
}

type TO_TestSoalRelasi struct {
	ID     uint `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TestID uint `json:"test_id" gorm:"column:test_id;not null"`
	SoalID uint `json:"soal_id" gorm:"column:soal_id;not null"`
}

// TableName specifies the table name
func (TO_TestSoalRelasi) TableName() string {
	return "to_test_soal"
}

type TO_Peserta struct {
	PesertaID  uint      `json:"peserta_id" gorm:"column:peserta_id;primaryKey;autoIncrement"`
	TestID     uint      `json:"test_id" gorm:"column:test_id"`
	SiswaNIS   string    `json:"siswa_nis" gorm:"column:siswa_nis"`
	KelasID    uint      `json:"kelas_id" gorm:"column:kelas_id"` // Tambahan untuk snapshot kelas
	Status     string    `json:"status" gorm:"column:status;default:not_started"`
	ExtraTime  int       `json:"extra_time" gorm:"column:extra_time"`
	NilaiAkhir float64   `json:"nilai_akhir" gorm:"column:nilai_akhir"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`

	// Relasi
	Test  TO_Test `json:"test" gorm:"foreignKey:TestID;references:TestID"`
	Siswa Siswa   `json:"siswa" gorm:"foreignKey:SiswaNIS;references:SiswaNIS"`
	Kelas Kelas   `json:"kelas" gorm:"foreignKey:KelasID;references:KelasId"` // Relasi ke Kelas
}

func (TO_Peserta) TableName() string {
	return "to_peserta"
}

type TO_JawabanFinal struct {
	FinalID      uint           `gorm:"primaryKey;column:final_id"`
	SessionID    uint           `gorm:"not null;index;column:session_id"`
	SoalID       uint           `gorm:"not null;column:soal_id"`
	JawabanSiswa datatypes.JSON `gorm:"type:json;column:jawaban_siswa"`
	SkorObjektif float64        `gorm:"type:decimal(5,2);default:0.00;column:skor_objektif"`
	SkorUraian   *float64       `gorm:"type:decimal(5,2);column:skor_uraian"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime;column:updated_at"`

	Session TestSession `gorm:"foreignKey:SessionID;references:SessionID;constraint:OnDelete:CASCADE;"`
	Soal    TO_BankSoal `gorm:"foreignKey:SoalID;references:SoalID;constraint:OnDelete:CASCADE;"`
}

func (TO_JawabanFinal) TableName() string {
	return "TO_JawabanFinal"
}

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
	IsSelected       bool         `json:"is_selected" gorm:"-"`
	JawabanTersimpan interface{}  `gorm:"-" json:"jawaban_tersimpan,omitempty"` // Field virtual
	Guru             Guru         `json:"guru" gorm:"foreignKey:GuruID;references:GuruID"`
	Kelas            Kelas        `json:"kelas" gorm:"foreignKey:KelasID;references:KelasId"`
	Mapel            Mapel        `gorm:"foreignKey:MapelID;references:KdMapel" json:"mapel"`
	Lampiran         *TO_Lampiran `json:"lampiran" gorm:"foreignKey:LampiranID;references:LampiranID"`
}

func (TO_BankSoal) TableName() string {
	return "to_banksoal"
}

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
