package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/secure-data", func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		c.JSON(http.StatusOK, gin.H{
			"message": "Secure Data",
			"token":   auth,
		})
	})

	// Use HTTPS with cert.pem and key.pem
	err := r.RunTLS(":8443", "cert.pem", "key.pem")
	if err != nil {
		panic("Failed to start HTTPS server: " + err.Error())
	}
}
