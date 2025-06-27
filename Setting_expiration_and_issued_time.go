//Setting expiration and issued time
/*
To set expiration time (exp) and issued at time (iat) in a 
JWT (JSON Web Token) in Go, you typically use the 
github.com/golang-jwt/jwt/v5 or github.com/golang-jwt/jwt library.
*/

//Set exp and iat in Go (using HS256)
package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my_secret_key") // Secret key for HMAC (HS256)

func main() {
	// Current time
	now := time.Now()

	// Create the claims with issued and expiration time
	claims := jwt.MapClaims{
		"username": "mahindra",                 // custom claim
		"iat":      now.Unix(),                 // issued at (in seconds)
		"exp":      now.Add(5 * time.Minute).Unix(), // expires in 5 minutes
	}

	// Create a token using the HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using the secret key
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("JWT Token:", signedToken)
}
