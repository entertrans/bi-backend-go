package dto

type PembayaranDTO struct {
	Tanggal string  `json:"tanggal"`
	Nominal float64 `json:"nominal"`
}

type PenerimaDTO struct {
	NIS          string          `json:"nis"`
	Nama         string          `json:"nama"`
	Kelas        string          `json:"kelas"`
	TotalTagihan float64         `json:"total_tagihan"`
	TotalBayar   float64         `json:"total_bayar"`
	Pembayaran   []PembayaranDTO `json:"pembayaran,omitempty"`
}

type InvoiceDetailDTO struct {
	IDInvoice     string        `json:"id_invoice"`
	TglInvoice    string        `json:"tgl_invoice"`
	TglJatuhTempo string        `json:"tgl_jatuh_tempo"`
	Deskripsi     string        `json:"deskripsi"`
	Penerima      []PenerimaDTO `json:"penerima"`
}
