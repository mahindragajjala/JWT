package main

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("secret")

func main() {
	r := gin.Default()

	r.POST("/token", func(c *gin.Context) {
		// Simulate user authentication
		var user struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&user); err != nil || user.Username != "mahindra" || user.Password != "pass" {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		// Simulate OIDC user identity
		claims := jwt.MapClaims{
			"sub":   "user123",
			"name":  "Mahindra",
			"email": "mahindra@example.com",
			"role":  "admin",
			"exp":   time.Now().Add(15 * time.Minute).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(jwtKey)

		c.JSON(200, gin.H{"access_token": tokenString})
	})

	r.Run(":9000")
}
