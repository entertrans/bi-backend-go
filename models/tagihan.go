package models

type Tagihan struct {
	TagihanId uint   `json:"id_tagihan" gorm:"column:id_tagihan"`
	Jenis     string `json:"jns_tagihan" gorm:"column:jns_tagihan"`
	Nominal   int    `json:"nom_tagihan" gorm:"column:nom_tagihan"`
}

func (Tagihan) TableName() string {
	return "tbl_tagihan"
}
