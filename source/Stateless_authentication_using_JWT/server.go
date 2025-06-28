package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Secret key for signing the token
var jwtKey = []byte("my_secret_key")

// User struct for login payload
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Generate JWT Token
func generateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(5 * time.Minute).Unix(), // expires in 5 min
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Middleware to validate token
func validateJWT(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	tokenStr := parts[1]

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.Next()
}

func main() {
	r := gin.Default()

	// Login route to generate token
	r.POST("/login", func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		if user.Username == "mahindra" && user.Password == "1234" {
			token, err := generateJWT(user.Username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"token": token})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong credentials"})
		}
	})

	// Protected route
	r.GET("/protected", validateJWT, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Access granted! Secure content served."})
	})

	r.Run(":8080") // Run server on port 8080
}
