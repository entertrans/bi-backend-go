package models

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