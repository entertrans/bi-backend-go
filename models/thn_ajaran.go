package models

type ThnAjaran struct {
	TahunAjaranId  uint   `json:"id_ta" gorm:"column:id_ta"`
	TahunAjaran    string `json:"thn_ajaran" gorm:"column:thn_ajaran"`
	Semester       string `json:"semester" gorm:"column:semester"`
	TglDikeluarkan string `json:"tgl_dikeluarkan" gorm:"column:tgl_dikeluarkan"`
}

func (ThnAjaran) TableName() string {
	return "tbl_thn_ajaran"
}
