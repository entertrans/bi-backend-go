package models

type Agama struct {
	AgamaId   uint   `json:"agama_id" gorm:"column:agama_id"`
	AgamaNama string `json:"agama_nama" gorm:"column:agama_nama"`
}

func (Agama) TableName() string {
	return "tbl_agama"
}
