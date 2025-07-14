package models

// === TABEL UTAMA: tbl_invoice ===
type Invoice struct {
	ID               uint              `gorm:"primaryKey" json:"id"`
	IDInvoice        string            `gorm:"column:id_invoice;unique" json:"id_invoice"`
	Deskripsi        string            `json:"deskripsi"`
	TglInvoice       string            `json:"tgl_invoice"`
	TglJatuhTempo    string            `json:"tgl_jatuh_tempo"`
	SudahAdaPenerima bool              `json:"sudah_ada_penerima"`
	Tagihan          []InvoiceTagihan  `gorm:"foreignKey:IDInvoice;references:IDInvoice" json:"tagihan"`
	Penerima         []InvoicePenerima `gorm:"foreignKey:IDInvoice;references:IDInvoice" json:"penerima"`
}

func (Invoice) TableName() string {
	return "tbl_invoice"
}

// === TABEL DETAIL TAGIHAN: tbl_invoice_tagihan ===
type InvoiceTagihan struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	IDInvoice   string `gorm:"column:id_invoice" json:"id_invoice"`
	NamaTagihan string `gorm:"column:nama_tagihan" json:"nama"`
	Nominal     int    `json:"nominal"`
}

func (InvoiceTagihan) TableName() string {
	return "tbl_invoice_tagihan"
}

// === TABEL PENERIMA: tbl_invoice_penerima ===
type InvoicePenerima struct {
	ID        uint                      `gorm:"primaryKey" json:"id"`
	IDInvoice string                    `gorm:"column:id_invoice" json:"id_invoice"`
	NIS       string                    `json:"nis"`
	Potongan  int                       `json:"potongan"`
	Tambahan  []InvoicePenerimaTambahan `gorm:"foreignKey:IDPenerima" json:"tambahan_tagihan"`

	Siswa Siswa `gorm:"foreignKey:NIS;references:SiswaNIS" json:"siswa"`
}

func (InvoicePenerima) TableName() string {
	return "tbl_invoice_penerima"
}

// === TABEL TAMBAHAN TAGIHAN PER SISWA: tbl_invoice_penerima_tambahan ===
type InvoicePenerimaTambahan struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	IDPenerima  uint   `gorm:"column:id_penerima" json:"id_penerima"`
	NamaTagihan string `gorm:"column:nama_tagihan" json:"nama"`
	Nominal     int    `json:"nominal"`
}

func (InvoicePenerimaTambahan) TableName() string {
	return "tbl_invoice_penerima_tambahan"
}
