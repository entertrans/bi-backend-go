package models

import (
	"time"

	"gorm.io/datatypes"
)

type TestSession struct {
	SessionID  uint       `gorm:"primaryKey;column:session_id"`
	TestID     uint       `gorm:"not null;column:test_id"`
	SiswaNIS   int        `gorm:"not null;column:siswa_nis"`
	StartTime  time.Time  `gorm:"not null;column:start_time"`
	EndTime    *time.Time `gorm:"column:end_time"`
	WaktuSisa  int        `gorm:"column:waktu_sisa"`
	Status     string     `gorm:"type:enum('in_progress','submitted','graded');default:'in_progress';column:status"`
	NilaiAwal  float64    `gorm:"type:decimal(5,2);default:0.00;column:nilai_awal"`
	NilaiAkhir float64    `gorm:"type:decimal(5,2);default:0.00;column:nilai_akhir"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime;column:updated_at"`

	JawabanFinal []JawabanFinal `gorm:"foreignKey:SessionID"`
}

// Explicit table name
func (TestSession) TableName() string {
	return "to_testsession" // Pastikan ini sesuai dengan nama tabel di database
}

type JawabanFinal struct {
	FinalID      uint           `gorm:"primaryKey;column:final_id"`
	SessionID    uint           `gorm:"not null;index;column:session_id"`
	SoalID       uint           `gorm:"not null;column:soal_id"`
	JawabanSiswa datatypes.JSON `gorm:"type:json;column:jawaban_siswa"`
	SkorObjektif float64        `gorm:"type:decimal(5,2);default:0.00;column:skor_objektif"`
	SkorUraian   *float64       `gorm:"type:decimal(5,2);column:skor_uraian"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime;column:updated_at"`

	Session TestSession `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE;"`
	Soal    TO_BankSoal `gorm:"foreignKey:SoalID;constraint:OnDelete:CASCADE;"`
}

func (JawabanFinal) TableName() string {
	return "to_jawabanfinal" // Pastikan ini sesuai dengan nama tabel di database
}
