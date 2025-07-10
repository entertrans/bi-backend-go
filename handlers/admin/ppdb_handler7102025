package adminhandlers

import (
	"log"
	"net/http"

	adminControllers "github.com/entertrans/bi-backend-go/controllers/admin"
	"github.com/gin-gonic/gin"
)

func HandleCreatePPDB(c *gin.Context) {
	log.Println("[INFO] POST /ppdb called")

	if err := adminControllers.CreatePPDBSiswa(c); err != nil {
		log.Println("[ERROR]", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Siswa PPDB berhasil ditambahkan"})
}
