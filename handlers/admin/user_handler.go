package adminhandlers

import (
	"net/http"

	admincontrollers "github.com/entertrans/bi-backend-go/controllers/admin"
	"github.com/gin-gonic/gin"
)

func GetMapelHandler(c *gin.Context) {
    mapels, err := admincontrollers.GetAllMapelWithGuruMapels()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch mapel"})
        return
    }
    c.JSON(http.StatusOK, mapels)
}
