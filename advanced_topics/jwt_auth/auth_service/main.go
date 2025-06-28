// auth service
package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var privateKeyPath = "private.key"

func main() {
	r := gin.Default()
	r.POST("/login", loginHandler)
	r.Run(":8000") // Auth Service
}

func loginHandler(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid"})
		return
	}

	keyData, _ := ioutil.ReadFile("/home/mahindra/jwt/auth_service/private.key")
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(keyData)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"username": body.Username,
		"exp":      time.Now().Add(10 * time.Minute).Unix(),
	})

	signedToken, _ := token.SignedString(privateKey)
	c.JSON(http.StatusOK, gin.H{"token": signedToken})
}
