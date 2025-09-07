package helpers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// SoalData menyimpan data soal untuk perhitungan skor
type SoalData struct {
	JawabanBenar    string
	Bobot           float64
	TipeSoal        string
	PilihanJawaban  []string
}

// ParseJawaban mengurai string jawaban menjadi interface{}
func ParseJawaban(jawabanStr string, tipeSoal string) (interface{}, error) {
	switch tipeSoal {
	case "pg", "isian_singkat", "uraian":
		var result string
		err := json.Unmarshal([]byte(jawabanStr), &result)
		return result, err
	
	case "pg_kompleks", "bs":
		var result []string
		err := json.Unmarshal([]byte(jawabanStr), &result)
		return result, err
	
	case "matching":
		var result []map[string]interface{}
		err := json.Unmarshal([]byte(jawabanStr), &result)
		return result, err
	
	default:
		return nil, fmt.Errorf("tipe soal tidak dikenali: %s", tipeSoal)
	}
}

// KonversiIndexKeTeks mengubah array index menjadi array teks pilihan
func KonversiIndexKeTeks(indexArray []string, pilihanJawaban []string) []string {
	var result []string
	
	for _, indexStr := range indexArray {
		index, err := strconv.Atoi(indexStr)
		if err == nil && index >= 0 && index < len(pilihanJawaban) {
			result = append(result, pilihanJawaban[index])
		} else {
			result = append(result, indexStr)
		}
	}
	
	return result
}

// NormalizeText untuk perbandingan case-insensitive
func NormalizeText(text string) string {
	return strings.TrimSpace(strings.ToLower(text))
}

// HitungSkorPG menghitung skor untuk tipe PG
func HitungSkorPG(jawabanSiswa, jawabanBenar interface{}, bobot float64, pilihanJawaban []string) float64 {
	siswaArr := parseToStringArray(jawabanSiswa)
	benarArr := parseToStringArray(jawabanBenar)
	
	siswaTeks := KonversiIndexKeTeks(siswaArr, pilihanJawaban)
	benarTeks := KonversiIndexKeTeks(benarArr, pilihanJawaban)
	
	if len(siswaTeks) > 0 && len(benarTeks) > 0 && siswaTeks[0] == benarTeks[0] {
		return bobot
	}
	return 0
}

// HitungSkorPGKompleks menghitung skor untuk PG Kompleks dan Benar/Salah
func HitungSkorPGKompleks(jawabanSiswa, jawabanBenar interface{}, bobot float64, pilihanJawaban []string) float64 {
	siswaArr := parseToStringArray(jawabanSiswa)
	benarArr := parseToStringArray(jawabanBenar)
	
	siswaTeks := KonversiIndexKeTeks(siswaArr, pilihanJawaban)
	benarTeks := KonversiIndexKeTeks(benarArr, pilihanJawaban)
	
	if len(benarTeks) == 0 {
		return 0
	}
	
	perItem := bobot / float64(len(benarTeks))
	var skor float64
	
	for i, benarItem := range benarTeks {
		if i < len(siswaTeks) && siswaTeks[i] == benarItem {
			skor += perItem
		}
	}
	
	return skor
}

// HitungSkorMatching menghitung skor untuk Matching
func HitungSkorMatching(jawabanSiswa, jawabanBenar interface{}, bobot float64, pilihanJawaban []string) float64 {
	// Implementasi matching seperti sebelumnya
	// ...
	return 0 // Placeholder
}

// HitungSkorIsian menghitung skor untuk Isian Singkat dan Uraian
func HitungSkorIsian(jawabanSiswa, jawabanBenar interface{}, bobot float64) float64 {
	siswaStr := parseToString(jawabanSiswa)
	benarStr := parseToString(jawabanBenar)
	
	if NormalizeText(siswaStr) == NormalizeText(benarStr) {
		return bobot
	}
	return 0
}

// Helper functions internal
func parseToStringArray(data interface{}) []string {
	switch v := data.(type) {
	case []string:
		return v
	case string:
		return []string{v}
	default:
		return []string{}
	}
}

func parseToString(data interface{}) string {
	switch v := data.(type) {
	case string:
		return v
	case []string:
		if len(v) > 0 {
			return v[0]
		}
		return ""
	default:
		return ""
	}
}