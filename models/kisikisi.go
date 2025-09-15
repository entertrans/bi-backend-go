package models

type KisiKisi struct {
	KisiKisiID        uint   `json:"kisikisi_id" gorm:"column:kisikisi_id;primaryKey;autoIncrement"`
	KisiKisiUb        string `json:"kisikisi_ub" gorm:"column:kisikisi_ub;size:25;not null"`
	KisiKisiDeskripsi string `json:"kisikisi_deskripsi" gorm:"column:kisikisi_deskripsi;type:text;not null"`
	KisiKisiMapel     uint   `json:"kisikisi_mapel" gorm:"column:kisikisi_mapel"`
	KisiKisiKelasID   uint   `json:"kisikisi_kelas_id" gorm:"column:kisikisi_kelas_id;not null"`
	KisiKisiSemester  int    `json:"kisikisi_semester" gorm:"column:kisikisi_semester"`

	// Relasi
	Mapel Mapel `gorm:"foreignKey:KisiKisiMapel;references:KdMapel" json:"mapel"`
	Kelas Kelas `gorm:"foreignKey:KisiKisiKelasID;references:KelasId" json:"kelas"`
}

func (KisiKisi) TableName() string {
	return "tbl_kisikisi"
}