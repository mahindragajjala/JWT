/*
🔧 Step 1: Generate a Legitimate JWT
    We'll create a JWT signed with a secret key.

⚔️ Step 2: Simulate an Attacker Modifying the Token
    We'll change the payload (e.g., change role from "user" to "admin") 
    without knowing the secret key.

✅ Step 3: Server Verifies Signature and Rejects Tampered Token
*/
//🔧 Step 1: Generate a Legitimate JWT
package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("mySecretKey")

func generateJWT() string {
	claims := jwt.MapClaims{
		"user": "alice",
		"role": "user",
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(secretKey)

	fmt.Println("✅ Legit JWT:", tokenString)
	return tokenString
}
//⚔️ Step 2: Tamper with the Token (Simulate Hacker)
/*
We’ll decode the token, change 
"role": "user" → "admin", re-encode without resigning.
*/
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func tamperJWT(originalToken string) string {
	parts := strings.Split(originalToken, ".")
	if len(parts) != 3 {
		panic("invalid token format")
	}

	// Decode payload
	payloadBytes, _ := base64.RawURLEncoding.DecodeString(parts[1])

	// Modify payload
	var payload map[string]interface{}
	_ = json.Unmarshal(payloadBytes, &payload)
	payload["role"] = "admin"

	newPayloadBytes, _ := json.Marshal(payload)
	newPayload := base64.RawURLEncoding.EncodeToString(newPayloadBytes)

	// Return tampered token (old header, new payload, original signature)
	tamperedToken := parts[0] + "." + newPayload + "." + parts[2]
	fmt.Println("⚠️ Tampered Token:", tamperedToken)
	return tamperedToken
}

//✅ Step 3: Server Verifies the Token
func verifyJWT(tokenString string) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure algorithm is HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		fmt.Println("❌ Token is invalid or tampered:", err)
		return
	}

	fmt.Println("✅ Token is valid and trusted")
}

