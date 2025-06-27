//Signing with a secret key or private key
/* 
🔐 1. Secret Key (HS256) — HMAC + SHA256
      Shared symmetric key.
      Same key is used to sign and verify.

🔐 2. Private Key (RS256) — RSA + SHA256
      Asymmetric key pair.
      Private key signs, public key verifies.
*/
//Signing JWT using Secret Key (HS256)
package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("my_super_secret_key")

func main() {
	// Claims
	claims := jwt.MapClaims{
		"username": "mahindra",
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(5 * time.Minute).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("HS256 JWT Token:", signedToken)
}


//Signing JWT using Private Key (RS256)
/* 
Generate 2048-bit RSA private key
        openssl genrsa -out private.key 2048

Extract public key from private key
        openssl rsa -in private.key -pubout -out public.key
 */
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Load private key
	privateKeyData, err := os.ReadFile("private.key")
	if err != nil {
		panic(err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		panic(err)
	}

	// Claims
	claims := jwt.MapClaims{
		"username": "mahindra",
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(5 * time.Minute).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign with RSA private key
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("RS256 JWT Token:", signedToken)
}

//Verification 
publicKeyData, _ := os.ReadFile("public.key")
publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)

parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
	return publicKey, nil
})

if parsedToken.Valid {
	fmt.Println("Valid token")
}
