/* 
A refresh token is a long-lived token used to obtain a new access 
token after the original access token expires — without requiring 
the user to log in again.




🔐 Why Use Refresh Tokens?
 Token Type     Purpose                Lifetime    
 Access Token   Access protected APIs  Short-lived 
 Refresh Token  Get new access tokens  Long-lived  



Access tokens are short-lived for security. 
If stolen, the damage is limited.
Refresh tokens are stored securely 
(e.g., in HTTP-only cookies or DB), 
and only used when renewing access tokens.



[Client]  →  [POST /login]               →  [Server validates user]
         ←  [AccessToken + RefreshToken]←

[Client]  →  [GET /protected]            →  [Server checks access token]

⏰ AccessToken Expired

[Client]  →  [POST /refresh-token]       →  [Send refresh token]
         ←  [New AccessToken]            ←  [Verify and issue new one]

🧰 Where to Store Refresh Tokens Securely?
HTTP-only Cookies (recommended)
Secure Database (linked to session/user)
Never expose them to JS (to avoid XSS attacks)

Login + Token Issuance Flow
[ Client ]                                [ Server (Go + JWT) ]
     |                                          |
     |  POST /login (username/password)         |
     |----------------------------------------->|
     |                                          |
     |     ✅ Validate credentials              |
     |     🔐 Generate Access Token (5 min)     |
     |     🔁 Generate Refresh Token (24 hrs)   |
     |                                          |
     |<-----------------------------------------|
     |  JSON { access_token, refresh_token }    |
     |                                          |


Accessing a Protected Route with Access Token
[ Client ]                                [ Server ]
     |                                          |
     |  GET /protected                          |
     |  Authorization: Bearer <access_token>    |
     |----------------------------------------->|
     |                                          |
     |  ✅ Validate token signature and expiry  |
     |  ✅ Extract claims                       |
     |<-----------------------------------------|
     |       Response: "Welcome {username}"     |



Access Token Expired → Use Refresh Token
[ Client ]                                [ Server ]
     |                                          |
     |  POST /refresh-token                     |
     |  JSON { refresh_token }                  |
     |----------------------------------------->|
     |                                          |
     | 🔁 Validate refresh token (long-lived)   |
     | 🔐 Issue new access token                |
     |                                          |
     |<-----------------------------------------|
     |    JSON { new_access_token }             |

*/
package main

import (
	"fmt"
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var accessSecret = []byte("access-secret")
var refreshSecret = []byte("refresh-secret")

func main() {
	r := gin.Default()
	r.POST("/login", loginHandler)
	r.POST("/refresh-token", refreshTokenHandler)
	r.GET("/protected", authMiddleware(), protectedHandler)
	r.Run(":8080")
}

func loginHandler(c *gin.Context) {
	username := c.PostForm("username")

	// Dummy check - you should verify against DB
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username required"})
		return
	}

	accessToken, err := createToken(username, 5*time.Minute, accessSecret)
	refreshToken, err2 := createToken(username, 24*time.Hour, refreshSecret)

	if err != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
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
	username := claims["sub"].(string)

	newAccessToken, err := createToken(username, 5*time.Minute, accessSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token creation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		token, err := jwt.Parse(authHeader, func(t *jwt.Token) (interface{}, error) {
			return accessSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
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

func createToken(username string, duration time.Duration, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
