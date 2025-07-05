package handlers

import (
	"log"
	"net/http"

	"github.com/entertrans/bi-backend-go/controllers"
	"github.com/gin-gonic/gin"
)

func HandleCreatePPDB(c *gin.Context) {
	log.Println("[INFO] POST /ppdb called")

	if err := controllers.CreatePPDBSiswa(c); err != nil {
		log.Println("[ERROR]", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Siswa PPDB berhasil ditambahkan"})
}
