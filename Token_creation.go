//Token_creation
//Install JWT package - go get github.com/golang-jwt/jwt/v5
package main
import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)
//Define Secret key
var jwtSecret = []byte("my_secret_key") // keep this secure in prod

// Create Custom Claims
type CustomClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

//Generate JWT Token
func GenerateJWTToken(username, role string) (string, error) {
	// Set expiration
	expirationTime := time.Now().Add(5 * time.Minute)

	// Create the claims
	claims := CustomClaims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "my-app",
		},
	}

	// Create the token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
func main() {
  //These details usually come from a login request made by the frontend (client).
	token, err := GenerateJWTToken("mahindra", "admin")
	if err != nil {
		fmt.Println("Error generating token:", err)
		return
	}
	fmt.Println("Generated JWT Token:", token)
}
/* 
🔐 1. User Sends Login Request (POST /login)
The frontend (React, Angular, mobile app, etc.) sends a request:
POST /login
Content-Type: application/json
                              {
                                "username": "mahindra",
                                "password": "123456"
                              }
                              
🧠 2. Backend Authenticates the User
Your Go backend:
            Checks username & password (e.g., from DB)
If valid, it generates a JWT like:
            GenerateJWTToken("mahindra", "admin")

📤 3. Backend Sends Back JWT
Response from backend:
HTTP/1.1 200 OK
Content-Type: application/json
              {
                "token": "eyJhbGciOiJIUzI1NiIsInR..."
              }
              
🛡️ 4. Client Uses JWT in Subsequent Requests
Now the client includes the token in every request:
        GET /user/profile
        Authorization: Bearer <jwt-token>
        
🧭 Backend then:
Extracts token from Authorization header
    Verifies it
    Reads claims (like username, role)
    Allows or denies access
*/
