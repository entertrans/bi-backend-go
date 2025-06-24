package models

type Kelas struct {
	KelasId   uint   `json:"kelas_id" gorm:"column:kelas_id"`
	KelasNama string `json:"kelas_nama" gorm:"column:kelas_nama"`
}

func (Kelas) TableName() string {
	return "tbl_kelas"
}
