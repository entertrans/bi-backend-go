package models

import "time"

type Pengguna struct {
	PenggunaID       uint      `gorm:"primaryKey;column:pengguna_id"`
	PenggunaUsername string    `gorm:"column:pengguna_username"`
	PenggunaPassword string    `gorm:"column:pengguna_password"`
	Spil             string    `gorm:"column:spil"`
	PenggunaStatus   int       `gorm:"column:pengguna_status"`
	PenggunaLevel    string    `gorm:"column:pengguna_level"` // 1=admin, 2=siswa, 3=guru
	RefID            *uint     `gorm:"column:ref_id"`
	PenggunaRegister time.Time `gorm:"column:pengguna_register"`

	// Relasi dengan kondisi
	Siswa *Siswa `gorm:"foreignKey:RefID;references:SiswaID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Guru  *Guru  `gorm:"foreignKey:RefID;references:GuruID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Admin *Admin `gorm:"foreignKey:RefID;references:AdminID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (Pengguna) TableName() string {
	return "tbl_pengguna"
}

type Admin struct {
	AdminID    uint   `gorm:"primaryKey;column:admin_id"`
	AdminNama  string `gorm:"column:admin_nama"`
	AdminEmail string `gorm:"column:admin_email"`
}

func (Admin) TableName() string {
	return "tbl_admin"
}
