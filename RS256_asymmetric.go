/* 
Private key
openssl genrsa -out private.pem 2048
Public key
openssl rsa -in private.pem -pubout -out public.pem
*/
package main

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func loadPrivateKey() (*rsa.PrivateKey, error) {
  //Loads RSA private key from disk to sign JWT.
	keyData, err := os.ReadFile("private.pem")
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(keyData)
}

func createRS256Token(privateKey *rsa.PrivateKey) (string, error) {
	claims := jwt.MapClaims{
		"username": "mahindra",
		"role":     "admin",
		"exp":      time.Now().Add(5 * time.Minute).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func main() {
	privKey, err := loadPrivateKey()
	if err != nil {
		fmt.Println("Failed to load private key:", err)
		return
	}

	token, err := createRS256Token(privKey)
	if err != nil {
		fmt.Println("Failed to create token:", err)
		return
	}

	fmt.Println("RS256 Token:", token)
}
/* 
 +----------------+                                        +---------------------+
 |                |                                        |                     |
 |   Frontend     |                                        |     Go Server       |
 | (React/Android)|                                        |  (Gin or net/http)  |
 |                |                                        |                     |
 +-------+--------+                                        +----------+----------+
         |                                                          |
         | 1. POST /login {"username": "mahindra", "password": "123"}  |
         | --------------------------------------------------------> |
         |                                                          |
         |                    2. Validate user (DB)                 |
         |                    3. Load RSA Private Key               |
         |                    4. Generate JWT using RS256           |
         |                    (createRS256Token)                    |
         |                                                          |
         | 5. Response {"token": "<JWT_TOKEN>"}                     |
         | <-------------------------------------------------------- |
         |                                                          |
         | 6. Store token in client (memory/localStorage)           |
         |                                                          |
         |---------------------- Protected API ---------------------|
         |                                                          |
         | 7. GET /profile (with Authorization: Bearer <JWT>)       |
         | -------------------------------------------------------->|
         |                                                          |
         |     8. Verify JWT using RSA Public Key (RS256)           |
         |     9. Decode claims and authorize                       |
         |                                                          |
         | 10. Return requested resource                            |
         | <--------------------------------------------------------|

*/
