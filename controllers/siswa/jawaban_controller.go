package siswa

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"gorm.io/datatypes"
)

func getSoalAndKunci(soalID uint, tipeUjian string) (string, float64, string, []string, error) {
	// fmt.Printf("DEBUG getSoalAndKunci: soalID=%d, tipeUjian=%s\n", soalID, tipeUjian)

	if tipeUjian == "ub" {
		var soal models.TO_BankSoal
		err := config.DB.First(&soal, soalID).Error
		if err != nil {
			fmt.Printf("ERROR getSoalAndKunci UB: %v\n", err)
			return "", 0, "", nil, err
		}

		// Parse pilihan jawaban
		var pilihan []string
		json.Unmarshal([]byte(soal.PilihanJawaban), &pilihan)

		// fmt.Printf("DEBUG getSoalAndKunci UB: JawabanBenar=%s, Bobot=%f, TipeSoal=%s, Pilihan=%v\n",
		// 	soal.JawabanBenar, soal.Bobot, soal.TipeSoal, pilihan)
		return soal.JawabanBenar, soal.Bobot, soal.TipeSoal, pilihan, nil
	} else {
		var soal models.TO_TestSoal
		err := config.DB.First(&soal, soalID).Error
		if err != nil {
			fmt.Printf("tipe ujiannya", tipeUjian)
			fmt.Printf("ERROR getSoalAndKunci masuk sini Test: %v\n", err)
			return "", 0, "", nil, err
		}

		// Parse pilihan jawaban
		var pilihan []string
		json.Unmarshal([]byte(soal.PilihanJawaban), &pilihan)

		fmt.Printf("DEBUG getSoalAndKunci Test: JawabanBenar=%s, Bobot=%f, TipeSoal=%s, Pilihan=%v\n",
			soal.JawabanBenar, soal.Bobot, soal.TipeSoal, pilihan)
		return soal.JawabanBenar, soal.Bobot, soal.TipeSoal, pilihan, nil
	}
}

func SaveJawabanFinal(sessionID uint, soalID uint, jawaban string, _ float64, tipeUjian string) error {
	// fmt.Printf("DEBUG SaveJawabanFinal: sessionID=%d, soalID=%d, jawaban=%s, tipeUjian=%s\n",
	//     sessionID, soalID, jawaban, tipeUjian)

	// 1. Ambil kunci, bobot soal, dan pilihan jawaban
	jawabanBenar, bobot, tipeSoal, pilihanJawaban, err := getSoalAndKunci(soalID, tipeUjian)
	if err != nil {
		// fmt.Printf("ERROR SaveJawabanFinal getSoalAndKunci: %v\n", err)
		return err
	}

	// fmt.Printf("DEBUG SaveJawabanFinal: jawabanBenar=%s, bobot=%f, tipeSoal=%s, pilihanJawaban=%v\n",
	// jawabanBenar, bobot, tipeSoal, pilihanJawaban)

	// 2. Hitung skor objektif dengan pilihan jawaban
	skorObjektif := hitungSkorObjektif(jawaban, jawabanBenar, tipeSoal, bobot, pilihanJawaban)
	// fmt.Printf("DEBUG SaveJawabanFinal: skorObjektif=%f\n", skorObjektif)

	// 3. Simpan atau update jawaban final
	var jawabanFinal models.JawabanFinal
	err = config.DB.Where("session_id = ? AND soal_id = ?", sessionID, soalID).
		First(&jawabanFinal).Error

	if err == nil {
		// fmt.Printf("DEBUG SaveJawabanFinal: Update existing record\n")
		jawabanFinal.JawabanSiswa = datatypes.JSON([]byte(jawaban))
		jawabanFinal.SkorObjektif = skorObjektif
		err = config.DB.Save(&jawabanFinal).Error
		// if err != nil {
		// 	fmt.Printf("ERROR SaveJawabanFinal update: %v\n", err)
		// } else {
		// 	fmt.Printf("DEBUG SaveJawabanFinal: Update successful\n")
		// }
		return err
	}

	// fmt.Printf("DEBUG SaveJawabanFinal: Create new record\n")
	newJawaban := models.JawabanFinal{
		SessionID:    sessionID,
		SoalID:       soalID,
		JawabanSiswa: datatypes.JSON([]byte(jawaban)),
		SkorObjektif: skorObjektif,
	}
	err = config.DB.Create(&newJawaban).Error
	// if err != nil {
	// 	fmt.Printf("ERROR SaveJawabanFinal create: %v\n", err)
	// } else {
	// 	fmt.Printf("DEBUG SaveJawabanFinal: Create successful\n")
	// }
	return err
}

func hitungSkorObjektif(jawabanSiswaStr string, jawabanBenarStr string, tipeSoal string, bobot float64, pilihanJawaban []string) float64 {
	// fmt.Printf("DEBUG hitungSkorObjektif: jawabanSiswaStr=%s, jawabanBenarStr=%s, tipeSoal=%s, bobot=%f, pilihanJawaban=%v\n",
	// 	jawabanSiswaStr, jawabanBenarStr, tipeSoal, bobot, pilihanJawaban)

	var skor float64 = 0

	switch tipeSoal {
	case "pg":
		var siswaArr []string
		var benarArr []string
		json.Unmarshal([]byte(jawabanSiswaStr), &siswaArr)
		json.Unmarshal([]byte(jawabanBenarStr), &benarArr)

		// fmt.Printf("DEBUG hitungSkorObjektif PG: siswaArr=%v, benarArr=%v\n", siswaArr, benarArr)

		// Konversi index → teks
		siswaTeks := konversiIndexKeTeks(siswaArr, pilihanJawaban)
		benarTeks := konversiIndexKeTeks(benarArr, pilihanJawaban)

		// fmt.Printf("DEBUG hitungSkorObjektif PG: siswaTeks=%v, benarTeks=%v\n", siswaTeks, benarTeks)

		if len(siswaTeks) > 0 && len(benarTeks) > 0 && siswaTeks[0] == benarTeks[0] {
			skor = bobot
		}

	case "pg_kompleks":
		var siswaIndex []string
		var benarIndex []string
		json.Unmarshal([]byte(jawabanSiswaStr), &siswaIndex)
		json.Unmarshal([]byte(jawabanBenarStr), &benarIndex)

		// Konversi index ke teks pilihan
		siswaTeks := konversiIndexKeTeks(siswaIndex, pilihanJawaban)
		benarTeks := konversiIndexKeTeks(benarIndex, pilihanJawaban)

		// fmt.Printf("DEBUG hitungSkorObjektif PG_KOMPLEKS: siswaIndex=%v, benarIndex=%v\n", siswaIndex, benarIndex)
		// fmt.Printf("DEBUG hitungSkorObjektif PG_KOMPLEKS: siswaTeks=%v, benarTeks=%v\n", siswaTeks, benarTeks)

		if len(benarTeks) > 0 {
			perItem := bobot / float64(len(benarTeks))
			// fmt.Printf("DEBUG hitungSkorObjektif PG_KOMPLEKS: perItem=%f\n", perItem)

			for i := range benarTeks {
				if i < len(siswaTeks) && siswaTeks[i] == benarTeks[i] {
					skor += perItem
					// fmt.Printf("DEBUG PG_KOMPLEKS: Item %d benar, skor sementara=%f\n", i, skor)
				} else {
					// fmt.Printf("DEBUG PG_KOMPLEKS: Item %d salah\n", i)
				}
			}
		}

	case "bs":
		var siswa []struct {
			Index   int    `json:"index"`
			Jawaban string `json:"jawaban"`
		}
		var benarArr []string

		json.Unmarshal([]byte(jawabanSiswaStr), &siswa)
		json.Unmarshal([]byte(jawabanBenarStr), &benarArr)

		// fmt.Printf("DEBUG hitungSkorObjektif BS: siswa=%v, benarArr=%v\n", siswa, benarArr)

		if len(benarArr) > 0 {
			perItem := bobot / float64(len(benarArr))
			benarCount := 0

			for _, js := range siswa {
				if js.Index >= 0 && js.Index < len(benarArr) {
					if strings.ToLower(strings.TrimSpace(js.Jawaban)) == strings.ToLower(strings.TrimSpace(benarArr[js.Index])) {
						benarCount++
						// fmt.Printf("DEBUG BS: index=%d jawaban='%s' BENAR\n", js.Index, js.Jawaban)
					} else {
						// fmt.Printf("DEBUG BS: index=%d jawaban='%s' SALAH (kunci='%s')\n", js.Index, js.Jawaban, benarArr[js.Index])
					}
				}
			}
			skor = float64(benarCount) * perItem
		}

	case "matching":
		var siswa []struct {
			LeftIndex  int `json:"leftIndex"`
			RightIndex int `json:"rightIndex"`
		}
		json.Unmarshal([]byte(jawabanSiswaStr), &siswa)

		// fmt.Printf("DEBUG hitungSkorObjektif matching: siswa=%v\n", siswa)

		if len(siswa) > 0 {
			perItem := bobot / float64(len(siswa))
			benarCount := 0
			for i, js := range siswa {
				if js.LeftIndex == js.RightIndex {
					benarCount++
					fmt.Printf("DEBUG matching: item %d benar (left=%d, right=%d)\n", i, js.LeftIndex, js.RightIndex)
				} else {
					fmt.Printf("DEBUG matching: item %d salah (left=%d, right=%d)\n", i, js.LeftIndex, js.RightIndex)
				}
			}
			skor = float64(benarCount) * perItem
		}

	case "isian_singkat", "uraian":
		var jwbSiswa string
		var jwbBenar string
		json.Unmarshal([]byte(jawabanSiswaStr), &jwbSiswa)
		json.Unmarshal([]byte(jawabanBenarStr), &jwbBenar)

		// fmt.Printf("DEBUG hitungSkorObjektif %s: jwbSiswa=%s, jwbBenar=%s\n", tipeSoal, jwbSiswa, jwbBenar)

		siswaTrimmed := strings.TrimSpace(strings.ToLower(jwbSiswa))
		benarTrimmed := strings.TrimSpace(strings.ToLower(jwbBenar))
		// fmt.Printf("DEBUG hitungSkorObjektif %s: setelah trim dan lowercase: siswa='%s', benar='%s'\n",
		// 	tipeSoal, siswaTrimmed, benarTrimmed)

		if siswaTrimmed == benarTrimmed {
			skor = bobot
		}
	}

	// fmt.Printf("DEBUG hitungSkorObjektif: Skor akhir=%f\n", skor)
	return skor
}

// Fungsi helper untuk konversi index ke teks pilihan
func konversiIndexKeTeks(indexArray []string, pilihanJawaban []string) []string {
	var result []string

	for _, indexStr := range indexArray {
		index, err := strconv.Atoi(indexStr)
		if err == nil && index >= 0 && index < len(pilihanJawaban) {
			result = append(result, pilihanJawaban[index])
		} else {
			// Jika konversi gagal, gunakan nilai asli
			result = append(result, indexStr)
		}
	}

	return result
}
