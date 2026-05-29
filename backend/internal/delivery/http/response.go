package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ok(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": message, "data": data})
}

func created(c *gin.Context, message string, data any) {
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": message, "data": data})
}

func fail(c *gin.Context, status int, err error) {
	message := "Error"
	if err != nil {
		message = err.Error()
	}
	c.JSON(status, gin.H{"success": false, "message": message})
}
