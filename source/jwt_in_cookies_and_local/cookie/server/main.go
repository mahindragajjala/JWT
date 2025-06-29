package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my_secret_key")

func main() {
	r := gin.Default()

	r.POST("/login", loginCookieHandler)
	r.GET("/profile", authCookieMiddleware, profileHandler)

	r.Run(":8080")
}

func loginCookieHandler(c *gin.Context) {
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

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func authCookieMiddleware(c *gin.Context) {
	cookie, err := c.Request.Cookie("token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	tokenStr := cookie.Value
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
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to your profile"})
}
