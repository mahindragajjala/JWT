//Parsing_and_validating_token
package main

import (
	"fmt"
	"net/http"
	"strings"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("your-secret-key")

func main() {
	http.HandleFunc("/secure", jwtMiddleware(secureHandler))
	http.ListenAndServe(":8080", nil)
}

// Middleware: Parse & validate JWT
func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		// Extract token
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Make sure the algorithm is what you expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		// Check token validity
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		fmt.Println("✅ Token successfully validated!")
		next.ServeHTTP(w, r)
	}
}

// Protected handler
func secureHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the secure endpoint!")
}




/* 
🔹 jwt.Parse(tokenString, keyFunc)
          Parses the JWT token string.
          Uses the keyFunc to provide the secret/public key.
          Internally validates the signature.

🔹 token.Valid
          Returns true only if:
          The signature is valid
          Token is not expired
          All standard validations pass

🔹 token.Claims
          You can access custom claims like email, role, etc., 
          by casting them.
*/
claims := token.Claims.(jwt.MapClaims)
user := claims["username"].(string)
fmt.Println("User from token:", user)



/* 
    [Client] ---> sends Authorization: Bearer <token>
        |
    [Go Server]
        |
        |--> r.Header.Get("Authorization")
        |--> Strip "Bearer " prefix
        |--> jwt.Parse(tokenString, keyFunc)
              |
              |--> Decodes header, payload, signature
              |--> Uses keyFunc to get secret
              |--> Validates signature (HMAC SHA256)
              |--> Checks exp, iat, nbf (claims)
        |
        |--> if token.Valid { allow access }
        |
        |--> else { 401 Unauthorized }
*/

//Validate Expiration and iat
/* 
{
  "sub": "user123",
  "exp": 1710000000,
  "iat": 1709990000
}
You can access claims manually for custom logic:
claims := token.Claims.(jwt.MapClaims)
fmt.Println("exp:", claims["exp"])
*/
