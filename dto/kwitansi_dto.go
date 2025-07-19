package dto

type KwitansiStatus struct {
	IDInvoice     string              `json:"id_invoice"`
	Deskripsi     string              `json:"deskripsi"`
	TglInvoice    string              `json:"tgl_invoice"`
	TglJatuhTempo string              `json:"tgl_jatuh_tempo"`
	Status        KwitansiStatusValue `json:"status"`
}

type KwitansiStatusValue struct {
	Belum      int `json:"belum"`
	BelumLunas int `json:"belum_lunas"`
	Lunas      int `json:"lunas"`
}

type KwitansiPembayaran struct {
	Tanggal string `json:"tanggal"`
	Nominal int    `json:"nominal"`
}

type KwitansiPenerimaDetail struct {
	ID           uint                 `json:"id"`
	NIS          string               `json:"nis"`
	Nama         string               `json:"nama"`
	Kelas        string               `json:"kelas"`
	Potongan     int                  `json:"potongan"`
	TotalTagihan int                  `json:"total_tagihan"`
	TotalBayar   int                  `json:"total_bayar"`
	Pembayaran   []KwitansiPembayaran `json:"pembayaran"`
}

type KwitansiDetail struct {
	IDInvoice     string                   `json:"id_invoice"`
	TglInvoice    string                   `json:"tgl_invoice"`
	TglJatuhTempo string                   `json:"tgl_jatuh_tempo"`
	Deskripsi     string                   `json:"deskripsi"`
	Penerima      []KwitansiPenerimaDetail `json:"penerima"`
}

type PembayaranRequest struct {
	IDPenerima int    `json:"id_penerima" binding:"required"`
	Tanggal    string `json:"tanggal" binding:"required"` // YYYY-MM-DD
	Nominal    int64  `json:"nominal" binding:"required"`
	Metode     string `json:"metode"`
	Keterangan string `json:"keterangan"`
}
