package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"strings"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my_secret_key")

func main() {
	r := gin.Default()

	r.POST("/login", loginHeaderHandler)
	r.GET("/profile", authHeaderMiddleware, profileHandler)

	r.Run(":8081")
}

func loginHeaderHandler(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&user); err != nil || user.Username != "mahindra" || user.Password != "123456" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	})
	tokenString, _ := token.SignedString(jwtKey)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func authHeaderMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	c.Next()
}

func profileHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Access granted to profile"})
}
