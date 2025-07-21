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
