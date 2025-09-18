package guruhandlers

import (
	gurucontrollers "github.com/entertrans/bi-backend-go/controllers/guru"
	"github.com/gin-gonic/gin"
)

// controller tunggal, bisa dipakai semua handler
var nilaiController = gurucontrollers.NewNilaiController()

func GetRekapNilai(c *gin.Context) {
	nilaiController.GetRekapNilai(c)
}

func GetDetailUB(c *gin.Context) {
	nilaiController.GetDetailUB(c)
}

func GetDetailTR(c *gin.Context) {
	nilaiController.GetDetailTR(c)
}

func GetDetailTugas(c *gin.Context) {
	nilaiController.GetDetailTugas(c)
}

func GetDetailPesertaTest(c *gin.Context) {
	nilaiController.GetDetailPesertaTest(c)
}
