package models

import "time"

type Pembayaran struct {
	ID         uint `gorm:"primaryKey"`
	IDPenerima uint
	Tanggal    time.Time
	Nominal    int
	Tujuan     *string
	Keterangan *string
	CreatedAt  time.Time
}

func (Pembayaran) TableName() string {
	return "tbl_pembayaran"
}
