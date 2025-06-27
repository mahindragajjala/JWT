/*
HS256 (HMAC + SHA256) — symmetric (same secret for signing & verifying)
*/
//Install - go get github.com/golang-jwt/jwt/v5
package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("my_secret_key")

func createHS256Token() (string, error) {
	claims := jwt.MapClaims{
		"username": "mahindra",
		"role":     "admin",
		"exp":      time.Now().Add(5 * time.Minute).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func main() {
	token, err := createHS256Token()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("HS256 Token:", token)
}
