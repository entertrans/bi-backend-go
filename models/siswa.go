package models

type Siswa struct {
	SiswaID              uint   `json:"siswa_id" gorm:"column:siswa_id"`
	SiswaNIS             string `json:"siswa_nis" gorm:"column:siswa_nis"`
	SiswaNISN            string `json:"siswa_nisn" gorm:"column:siswa_nisn"`
	NoIjazah             string `json:"no_ijazah" gorm:"column:no_ijazah"`
	NIKSiswa             string `json:"nik_siswa" gorm:"column:nik_siswa"`
	SiswaNama            string `json:"siswa_nama" gorm:"column:siswa_nama"`
	SiswaJenkel          string `json:"siswa_jenkel" gorm:"column:siswa_jenkel"`
	SiswaTempat          string `json:"siswa_tempat" gorm:"column:siswa_tempat"`
	SiswaTglLahir        string `json:"siswa_tgl_lahir" gorm:"column:siswa_tgl_lahir"` // bisa diganti time.Time kalau pakai parsing tanggal
	SiswaAgamaID         int    `json:"siswa_agama_id" gorm:"column:siswa_agama_id"`
	SiswaKewarganegaraan string `json:"siswa_kewarganegaraan" gorm:"column:siswa_kewarganegaraan"`
	SiswaAlamat          string `json:"siswa_alamat" gorm:"column:siswa_alamat"`
	SiswaEmail           string `json:"siswa_email" gorm:"column:siswa_email"`
	SiswaDokumen         string `json:"siswa_dokumen" gorm:"column:siswa_dokumen"`
	SiswaNoTelp          string `json:"siswa_no_telp" gorm:"column:siswa_no_telp"`
	SiswaKelasID         int    `json:"siswa_kelas_id" gorm:"column:siswa_kelas_id"`
	SiswaPhoto           string `json:"siswa_photo" gorm:"column:siswa_photo"`
	SoftDeleted          int    `json:"soft_deleted" gorm:"column:soft_deleted"`
	AnakKe               int    `json:"anak_ke" gorm:"column:anak_ke"`
	SekolahAsal          string `json:"sekolah_asal" gorm:"column:sekolah_asal"`
	Satelit              int    `json:"satelit" gorm:"column:satelit"`
	OC                   int    `json:"oc" gorm:"column:oc"`
	KC                   int    `json:"kc" gorm:"column:kc"`

	// Relasi: satu siswa punya satu ortu
	Orangtua Orangtua `json:"orangtua" gorm:"foreignKey:SiswaNIS;references:SiswaNIS"`
	Agama    Agama    `json:"agama" gorm:"foreignKey:SiswaAgamaID;references:AgamaId"`
}

func (Siswa) TableName() string {
	return "tbl_siswa"
}
