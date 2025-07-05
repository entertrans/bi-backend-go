package models

type Lampiran struct {
	LampiranId   uint   `json:"id_lampiran" gorm:"column:id_lampiran"`
	SiswaNIS     string `json:"siswa_nis" gorm:"column:siswa_nis"` // UPPERCASE di Go
	JenisDokumen string `json:"dokumen_jenis" gorm:"column:dokumen_jenis"`
	Url          string `json:"url" gorm:"column:url"`
	Upload       string `json:"uploaded_at" gorm:"column:uploaded_at"`
}

func (Lampiran) TableName() string {
	return "tbl_lampiran"
}
