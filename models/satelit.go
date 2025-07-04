package models

type DtSatelit struct {
	SatelitId   uint   `json:"satelit_id" gorm:"column:satelit_id"`
	SatelitNama string `json:"satelit_nama" gorm:"column:satelit_nama"`
}

func (DtSatelit) TableName() string {
	return "tbl_satelit"
}
