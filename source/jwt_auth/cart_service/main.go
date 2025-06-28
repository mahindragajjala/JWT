//Cart service

package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func main() {
	r := gin.Default()
	r.Use(JWTMiddleware("public.key"))

	r.GET("/cart", func(c *gin.Context) {
		username := c.GetString("username")
		c.JSON(http.StatusOK, gin.H{"cart": "Cart items for " + username})
	})
	r.Run(":8002")
}
func JWTMiddleware(publicKeyPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || len(auth) < 8 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token"})
			return
		}
		tokenStr := auth[7:] // Strip "Bearer "

		pubKeyBytes, _ := ioutil.ReadFile(publicKeyPath)
		pubKey, _ := jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return pubKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("username", claims["username"])
		c.Next()
	}
}
