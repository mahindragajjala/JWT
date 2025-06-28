package main

import (
	"fmt"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var accessSecret = []byte("ACCESS_SECRET")
var refreshSecret = []byte("REFRESH_SECRET")

// In-memory store: refreshUUID -> userID
var refreshTokenStore = make(map[string]string)

type LoginInput struct {
	UserID string `json:"user_id"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func main() {
	r := gin.Default()

	r.POST("/login", LoginHandler)
	r.POST("/refresh", RefreshHandler)
	r.GET("/protected", AuthMiddleware(), ProtectedHandler)

	r.Run(":8080")
}

func LoginHandler(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	tokens, err := GenerateTokens(input.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func GenerateTokens(userID string) (*TokenResponse, error) {
	accessUUID := uuid.NewString()
	refreshUUID := uuid.NewString()

	atClaims := jwt.MapClaims{
		"user_id":     userID,
		"access_uuid": accessUUID,
		"exp":         time.Now().Add(15 * time.Minute).Unix(),
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString(accessSecret)
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{
		"user_id":      userID,
		"refresh_uuid": refreshUUID,
		"exp":          time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString(refreshSecret)
	if err != nil {
		return nil, err
	}

	// Store refresh UUID
	refreshTokenStore[refreshUUID] = userID

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func RefreshHandler(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing refresh token"})
		return
	}

	token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	refreshUUID := claims["refresh_uuid"].(string)
	userID := claims["user_id"].(string)

	if storedUserID, ok := refreshTokenStore[refreshUUID]; !ok || storedUserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token reuse detected"})
		return
	}

	// Delete old refresh token
	delete(refreshTokenStore, refreshUUID)

	// Generate new tokens
	newTokens, err := GenerateTokens(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new tokens"})
		return
	}

	c.JSON(http.StatusOK, newTokens)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}
		tokenStr := authHeader[len("Bearer "):]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return accessSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
			return
		}
		c.Next()
	}
}

func ProtectedHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to protected route!"})
}
