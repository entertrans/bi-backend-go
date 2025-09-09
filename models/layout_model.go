package models

type JawabanResponse struct {
	SoalID         uint     `json:"soal_id"`
	Pertanyaan     string   `json:"pertanyaan"`
	TipeSoal       string   `json:"tipe_soal"`
	JawabanSiswa   string   `json:"jawaban_siswa"`
	JawabanBenar   *string  `json:"jawaban_benar,omitempty"`
	PilihanJawaban *string  `json:"pilihan_jawaban,omitempty"`
	SkorObjektif   float64  `json:"skor_objektif"`
	SkorUraian     *float64 `json:"skor_uraian"`
	MaxScore       float64  `json:"max_score"` // TAMBAHKAN BOBOT/MAX SCORE
	Score          float64  `json:"skor"`      // TAMBAHKAN BOBOT/MAX SCORE
	// Tambahkan field untuk lampiran
	LampiranNamaFile *string `json:"lampiran_nama_file,omitempty"`
	LampiranTipeFile *string `json:"lampiran_tipe_file,omitempty"`
	LampiranPathFile *string `json:"lampiran_path_file,omitempty"`
}

type TO_Soal struct {
	SoalID         uint              `json:"soal_id" gorm:"column:soal_id;primaryKey;autoIncrement"`
	TestID         uint              `json:"test_id" gorm:"column:test_id"`
	TipeSoal       string            `json:"tipe_soal" gorm:"column:tipe_soal"`
	Pertanyaan     string            `json:"pertanyaan" gorm:"column:pertanyaan"`
	LampiranID     *uint             `json:"lampiran_id" gorm:"column:lampiran_id"`
	Lampiran       *TO_Lampiran      `json:"lampiran" gorm:"foreignKey:LampiranID"`
	PilihanJawaban []TO_JawabanFinal `json:"pilihan_jawaban" gorm:"foreignKey:SoalID"`
	JawabanBenar   []TO_JawabanFinal `json:"jawaban_benar" gorm:"foreignKey:SoalID"`
}
