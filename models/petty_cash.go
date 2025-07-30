package models

import (
	"time"
)

type PettyCashPeriode struct {
	ID           int         `json:"id"`
	KodePeriode  string      `json:"kode_periode"`
	Deskripsi    string      `json:"deskripsi"`
	TanggalMulai time.Time   `json:"tanggal_mulai"`
	TanggalTutup *time.Time  `json:"tanggal_tutup"` // nullable
	SaldoAwal    int64       `json:"saldo_awal"`
	Lokasi       string      `json:"lokasi"`
	Status       string      `json:"status"` // aktif / ditutup
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Transaksis   []Transaksi `json:"transaksis" gorm:"foreignKey:IDPeriode;references:ID"` // ini penting!
}

func (PettyCashPeriode) TableName() string {
	return "tbl_petty_cash_periode"
}

type Transaksi struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	IDPeriode    int       `json:"id_periode"`
	Tanggal      time.Time `json:"tanggal"`
	Keterangan   string    `json:"keterangan"`
	Jenis        string    `json:"jenis"` // masuk / keluar
	Nominal      int64     `json:"nominal"`
	SaldoSetelah int64     `json:"saldo_setelah"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Transaksi) TableName() string {
	return "tbl_transaksi"
}
