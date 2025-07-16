package models

type Pembayaran struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	IDPenerima uint   `gorm:"column:id_penerima" json:"id_penerima"`
	Tanggal    string `json:"tanggal"` // Format: YYYY-MM-DD
	Nominal    int    `json:"nominal"`
	Metode     string `json:"metode,omitempty"`  // opsional
	Catatan    string `json:"catatan,omitempty"` // opsional

	// Relasi
	Penerima InvoicePenerima `gorm:"foreignKey:IDPenerima" json:"penerima,omitempty"`
}

func (Pembayaran) TableName() string {
	return "tbl_pembayaran"
}
