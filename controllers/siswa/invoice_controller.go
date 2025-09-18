package siswa

import (
	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

type HistorySiswaKeuangan struct {
	Siswa   models.Siswa     `json:"siswa"`
	History []HistoryInvoice `json:"history"`
}

type HistoryInvoice struct {
	InvoiceID         string                           `json:"invoice_id"`
	InvoiceDeskripsi  string                           `json:"deskripsi"`
	InvoiceTgl        string                           `json:"tgl_invoice"`
	InvoiceJatuhTempo string                           `json:"tgl_jatuh_tempo"`
	Potongan          int                              `json:"potongan"`
	Tagihan           []models.InvoiceTagihan          `json:"tagihan"`
	TambahanTagihan   []models.InvoicePenerimaTambahan `json:"tambahan_tagihan"`
	TotalTagihan      int                              `json:"total_tagihan"`
	TotalBayar        int                              `json:"total_bayar"`
	Pembayaran        []models.Pembayaran              `json:"pembayaran"`
}

func GetHistoryKeuanganByNIS(nis string) (HistorySiswaKeuangan, error) {
	var penerimas []models.InvoicePenerima
	err := config.DB.Preload("Invoice.Tagihan").
		Preload("Tambahan").
		Preload("Pembayaran").
		Preload("Siswa").
		Preload("Siswa.Kelas").
		Where("nis = ?", nis).
		Find(&penerimas).Error
	if err != nil {
		return HistorySiswaKeuangan{}, err
	}

	var history []HistoryInvoice
	var siswa models.Siswa

	for _, penerima := range penerimas {
		if siswa.SiswaNama == nil || *siswa.SiswaNama == "" {
			siswa = penerima.Siswa
		}

		totalTagihan := 0
		for _, tagihan := range penerima.Invoice.Tagihan {
			totalTagihan += tagihan.Nominal
		}
		for _, tambahan := range penerima.Tambahan {
			totalTagihan += tambahan.Nominal
		}
		totalTagihan -= penerima.Potongan

		totalBayar := 0
		for _, bayar := range penerima.Pembayaran {
			totalBayar += bayar.Nominal
		}

		history = append(history, HistoryInvoice{
			InvoiceID:         penerima.Invoice.IDInvoice,
			InvoiceDeskripsi:  penerima.Invoice.Deskripsi,
			InvoiceTgl:        penerima.Invoice.TglInvoice,
			InvoiceJatuhTempo: penerima.Invoice.TglJatuhTempo,
			Potongan:          penerima.Potongan,
			Tagihan:           penerima.Invoice.Tagihan,
			TambahanTagihan:   penerima.Tambahan,
			TotalTagihan:      totalTagihan,
			TotalBayar:        totalBayar,
			Pembayaran:        penerima.Pembayaran,
		})
	}

	return HistorySiswaKeuangan{
		Siswa:   siswa,
		History: history,
	}, nil
}

// === Detail 1 invoice ===
type InvoiceDetail struct {
	InvoiceID       string                           `json:"invoice_id"`
	Deskripsi       string                           `json:"deskripsi"`
	TglInvoice      string                           `json:"tgl_invoice"`
	TglJatuhTempo   string                           `json:"tgl_jatuh_tempo"`
	TagihanUtama    []models.InvoiceTagihan          `json:"tagihan_utama"`
	TagihanTambahan []models.InvoicePenerimaTambahan `json:"tagihan_tambahan"`
	Pembayaran      []models.Pembayaran              `json:"pembayaran"`
	TotalTagihan    int                              `json:"total_tagihan"`
	TotalPotongan   int                              `json:"total_potongan"`
	TotalBayar      int                              `json:"total_bayar"`
	SisaTagihan     int                              `json:"sisa_tagihan"`
}

func GetInvoiceDetailByNIS(nis string, idInvoice string) (InvoiceDetail, error) {
	var penerima models.InvoicePenerima
	err := config.DB.Preload("Invoice.Tagihan").
		Preload("Tambahan").
		Preload("Pembayaran").
		Where("nis = ? AND id_invoice = ?", nis, idInvoice).
		First(&penerima).Error
	if err != nil {
		return InvoiceDetail{}, err
	}

	totalTagihan := 0
	for _, t := range penerima.Invoice.Tagihan {
		totalTagihan += t.Nominal
	}

	totalTambahan := 0
	for _, t := range penerima.Tambahan {
		totalTambahan += t.Nominal
	}

	totalBayar := 0
	for _, p := range penerima.Pembayaran {
		totalBayar += p.Nominal
	}

	totalFinal := totalTagihan - penerima.Potongan + totalTambahan
	sisa := totalFinal - totalBayar

	return InvoiceDetail{
		InvoiceID:       penerima.Invoice.IDInvoice,
		Deskripsi:       penerima.Invoice.Deskripsi,
		TglInvoice:      penerima.Invoice.TglInvoice,
		TglJatuhTempo:   penerima.Invoice.TglJatuhTempo,
		TagihanUtama:    penerima.Invoice.Tagihan,
		TagihanTambahan: penerima.Tambahan,
		Pembayaran:      penerima.Pembayaran,
		TotalTagihan:    totalFinal,
		TotalPotongan:   penerima.Potongan,
		TotalBayar:      totalBayar,
		SisaTagihan:     sisa,
	}, nil
}

type KwitansiTagihan struct {
	InvoiceID     string `json:"invoice_id"`
	Deskripsi     string `json:"deskripsi"`
	TglInvoice    string `json:"tgl_invoice"`
	TglJatuhTempo string `json:"tgl_jatuh_tempo"`
	TotalTagihan  int    `json:"total_tagihan"`
	TotalBayar    int    `json:"total_bayar"`
	SisaTagihan   int    `json:"sisa_tagihan"`
	Status        string `json:"status"` // contoh: "Belum Lunas", "Lunas"
}

func GetLatestUnpaidInvoiceByNIS(nis string) (KwitansiTagihan, error) {
	var penerima models.InvoicePenerima
	err := config.DB.
		Preload("Invoice.Tagihan").
		Preload("Tambahan").
		Preload("Pembayaran").
		Where("nis = ?", nis).
		Order("id_invoice DESC"). // ambil yang terbaru
		First(&penerima).Error
	if err != nil {
		return KwitansiTagihan{}, err
	}

	// Hitung total
	totalTagihan := 0
	for _, t := range penerima.Invoice.Tagihan {
		totalTagihan += t.Nominal
	}
	totalTambahan := 0
	for _, t := range penerima.Tambahan {
		totalTambahan += t.Nominal
	}
	totalBayar := 0
	for _, p := range penerima.Pembayaran {
		totalBayar += p.Nominal
	}

	totalFinal := totalTagihan - penerima.Potongan + totalTambahan
	sisa := totalFinal - totalBayar

	status := "Lunas"
	if sisa > 0 {
		status = "Belum Lunas"
	}

	return KwitansiTagihan{
		InvoiceID:     penerima.Invoice.IDInvoice,
		Deskripsi:     penerima.Invoice.Deskripsi,
		TglInvoice:    penerima.Invoice.TglInvoice,
		TglJatuhTempo: penerima.Invoice.TglJatuhTempo,
		TotalTagihan:  totalFinal,
		TotalBayar:    totalBayar,
		SisaTagihan:   sisa,
		Status:        status,
	}, nil
}
