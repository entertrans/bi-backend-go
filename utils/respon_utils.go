package utils

import (
	"github.com/gin-gonic/gin"
)

// RespondJSON mengirim response JSON
func RespondJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{
		"status": status,
		"data":   data,
	})
}

// RespondError mengirim response error
func RespondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"status":  status,
		"message": message,
	})
}