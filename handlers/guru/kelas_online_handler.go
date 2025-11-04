package guruhandlers

import (
	"net/http"

	gurucontrollers "github.com/entertrans/bi-backend-go/controllers/guru"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

func GetAllKelasOnlineHandler(c *gin.Context) {
    result, err := gurucontrollers.GetAllKelasOnline()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, result)
}

func GetKelasOnlineByIDHandler(c *gin.Context) {
    id := c.Param("id")
    result, err := gurucontrollers.GetKelasOnlineByID(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, result)
}

func CreateKelasOnlineHandler(c *gin.Context) {
    var input models.KelasOnline
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    result, err := gurucontrollers.CreateKelasOnline(input)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, result)
}

func UpdateKelasOnlineHandler(c *gin.Context) {
    id := c.Param("id")
    var input models.KelasOnline
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    result, err := gurucontrollers.UpdateKelasOnline(id, input)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, result)
}

func DeleteKelasOnlineHandler(c *gin.Context) {
    id := c.Param("id")
    err := gurucontrollers.DeleteKelasOnline(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}