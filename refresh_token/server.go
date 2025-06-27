package main

import (
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var accessSecret = []byte("access-secret")
var refreshSecret = []byte("refresh-secret")

func main() {
	router := gin.Default()

	router.POST("/login", loginHandler)
	router.POST("/refresh-token", refreshTokenHandler)
	router.GET("/protected", authMiddleware(), protectedHandler)

	router.Run(":8080")
}

func loginHandler(c *gin.Context) {
	username := c.PostForm("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username required"})
		return
	}

	accessToken, _ := createToken(username, 30*time.Second, accessSecret)
	refreshToken, _ := createToken(username, 2*time.Minute, refreshSecret)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func refreshTokenHandler(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BindJSON(&req); err != nil || req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token required"})
		return
	}

	token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	username := claims["sub"].(string)

	newAccessToken, _ := createToken(username, 30*time.Second, accessSecret)

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		token, err := jwt.Parse(auth, func(t *jwt.Token) (interface{}, error) {
			return accessSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("user", token.Claims.(jwt.MapClaims)["sub"])
		c.Next()
	}
}

func protectedHandler(c *gin.Context) {
	user := c.GetString("user")
	c.JSON(http.StatusOK, gin.H{"message": "Welcome " + user})
}

func createToken(user string, duration time.Duration, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub": user,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
