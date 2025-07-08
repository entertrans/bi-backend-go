package handlers

import (
	"net/http"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
	"github.com/gin-gonic/gin"
)

func GetAllAgama(c *gin.Context) {
	var agama []models.Agama
	err := config.DB.Debug().Find(&agama).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agama)
}
