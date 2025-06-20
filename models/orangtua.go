package models

type Orangtua struct {
	OrtuID          uint   `json:"ortu_id" gorm:"column:ortu_id"`
	SiswaNIS        string `json:"siswa_nis" gorm:"column:siswa_nis"`
	AyahNama        string `json:"ayah_nama" gorm:"column:ayah_nama"`
	AyahNik         string `json:"ayah_nik" gorm:"column:ayah_nik"`
	AyahTempat      string `json:"ayah_tempat" gorm:"column:ayah_tempat"`
	AyahTanggal     string `json:"ayah_tanggal" gorm:"column:ayah_tanggal"`
	AyahPendidikan  string `json:"ayah_pendidikan" gorm:"column:ayah_pendidikan"`
	AyahPekerjaan   string `json:"ayah_pekerjaan" gorm:"column:ayah_pekerjaan"`
	AyahPenghasilan string `json:"ayah_penghasilan" gorm:"column:ayah_penghasilan"`
	AyahNotelp      string `json:"no_telp_ayah" gorm:"column:no_telp_ayah"`
	AyahEmail       string `json:"email_ayah" gorm:"column:email_ayah"`
	IbuNama         string `json:"ibu_nama" gorm:"column:ibu_nama"`
	IbuNik          string `json:"ibu_nik" gorm:"column:ibu_nik"`
	IbuTempat       string `json:"ibu_tempat" gorm:"column:ibu_tempat"`
	IbuTanggal      string `json:"ibu_tanggal" gorm:"column:ibu_tanggal"`
	IbuPendidikan   string `json:"ibu_pendidikan" gorm:"column:ibu_pendidikan"`
	IbuPekerjaan    string `json:"ibu_pekerjaan" gorm:"column:ibu_pekerjaan"`
	IbuPenghasilan  string `json:"ibu_penghasilan" gorm:"column:ibu_penghasilan"`
	IbuNotelp       string `json:"no_telp_ibu" gorm:"column:no_telp_ibu"`
	IbuEmail        string `json:"email_ibu" gorm:"column:email_ibu"`
	WaliNama        string `json:"wali_nama" gorm:"column:wali_nama"`
	WaliNik         string `json:"wali_nik" gorm:"column:wali_nik"`
	WaliTempat      string `json:"wali_tempat" gorm:"column:wali_tempat"`
	WaliTanggal     string `json:"wali_tanggal" gorm:"column:wali_tanggal"`
	WaliPendidikan  string `json:"wali_pendidikan" gorm:"column:wali_pendidikan"`
	WaliPekerjaan   string `json:"wali_pekerjaan" gorm:"column:wali_pekerjaan"`
	WaliPenghasilan string `json:"wali_penghasilan" gorm:"column:wali_penghasilan"`
	WaliAlamat      string `json:"wali_alamat" gorm:"column:wali_alamat"`
	WaliNotelp      string `json:"wali_notelp" gorm:"column:wali_notelp"`
}

func (Orangtua) TableName() string {
	return "tbl_orangtua"
}
