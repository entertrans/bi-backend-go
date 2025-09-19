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

type OnlineClass struct {
	IDOnlineClass uint      `gorm:"column:id_online_class;primaryKey;autoIncrement" json:"id_online_class"`
	IDKelasMapel  uint      `gorm:"column:id_kelas_mapel;not null" json:"id_kelas_mapel"`
	Tanggal       time.Time `gorm:"column:tanggal;type:date;not null" json:"tanggal"`
	Mulai         string    `gorm:"column:mulai;type:time;not null" json:"mulai"`
	Selesai       string    `gorm:"column:selesai;type:time;not null" json:"selesai"`
	Status        string    `gorm:"column:status;type:enum('belum','sedang','selesai');default:'belum'" json:"status"`
	MeetLink      string    `gorm:"column:meet_link;size:255;not null" json:"meet_link"`

	// Relasi
	KelasMapel KelasMapel `gorm:"foreignKey:IDKelasMapel;references:ID" json:"kelas_mapel"`
}

func (OnlineClass) TableName() string {
	return "tbl_online_class"
}
