package models

import "time"

type Kelas struct {
	KelasId     uint         `json:"kelas_id" gorm:"column:kelas_id"`
	KelasNama   string       `json:"kelas_nama" gorm:"column:kelas_nama"`
	KelasMapels []KelasMapel `json:"kelas_mapels" gorm:"foreignKey:KelasID;references:KelasId"`
}

func (Kelas) TableName() string {
	return "tbl_kelas"
}

type KelasMapel struct {
	ID      uint `gorm:"column:id_kelas_mapel;primaryKey;autoIncrement" json:"id"`
	KelasID uint `gorm:"column:kelas_id;not null" json:"kelas_id"`
	KdMapel uint `gorm:"column:kd_mapel;not null" json:"kd_mapel"`

	// Relasi
	Kelas Kelas `gorm:"foreignKey:KelasID;references:KelasId"`
	Mapel Mapel `gorm:"foreignKey:KdMapel;references:KdMapel"`
}

func (KelasMapel) TableName() string {
	return "tbl_kelas_mapel"
}


type KelasOnline struct {
	IDKelasOnline uint      `gorm:"column:id_kelas_online;primaryKey;autoIncrement" json:"id_kelas_online"`
	IDKelasMapel  uint      `gorm:"column:id_kelas_mapel;not null" json:"id_kelas_mapel"`
	GuruID        uint      `gorm:"column:guru_id;not null" json:"guru_id"`
	JudulKelas    string    `gorm:"column:judul_kelas;size:150;not null" json:"judul_kelas"`
	TanggalKelas  time.Time `gorm:"column:tanggal_kelas;not null" json:"tanggal_kelas"`
	JamMulai      string    `gorm:"column:jam_mulai;size:10;not null" json:"jam_mulai"`
	JamSelesai    string    `gorm:"column:jam_selesai;size:10;not null" json:"jam_selesai"`
	Status        string    `gorm:"column:status;type:enum('akan_berlangsung','sedang_berlangsung','selesai');default:'akan_berlangsung'" json:"status"`
	LinkKelas     string    `gorm:"column:link_kelas;size:255" json:"link_kelas"`
	MateriLink    string    `gorm:"column:materi_link;size:255" json:"materi_link"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	KelasMapel KelasMapel `gorm:"foreignKey:IDKelasMapel;references:ID" json:"kelas_mapel"`
	Guru       Guru       `gorm:"foreignKey:GuruID;references:GuruID" json:"guru"`
}

func (KelasOnline) TableName() string {
	return "tbl_kelas_online"
}
