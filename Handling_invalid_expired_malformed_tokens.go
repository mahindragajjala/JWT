//Handling_invalid_expired_malformed_tokens
/* 
  Error Type                 When It Happens                                     
 Malformed Token        Token is not 3 parts or not properly Base64 encoded 
 Invalid Signature      Token signature doesn't match the expected value    
 Expired Token          exp claim has passed                              
 Not Yet Valid          nbf (not before) or iat is in the future        
 Unsupported Algorithm  Algorithm used in token is not allowed              
 */
package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("your-secret-key")

func main() {
	http.HandleFunc("/secure", jwtMiddleware(secureHandler))
	http.ListenAndServe(":8080", nil)
}

func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Authorization header missing or invalid", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		// Parse token and capture more detailed error
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		// 🧠 Handle different error types
		if err != nil {
			var ve *jwt.ValidationError
			if errors.As(err, &ve) {
				switch {
				case ve.Errors&jwt.ValidationErrorMalformed != 0:
					http.Error(w, "Malformed token", http.StatusBadRequest)
				case ve.Errors&jwt.ValidationErrorExpired != 0:
					http.Error(w, "Token is expired", http.StatusUnauthorized)
				case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
					http.Error(w, "Token not valid yet", http.StatusUnauthorized)
				default:
					http.Error(w, "Invalid token", http.StatusUnauthorized)
				}
			} else {
				http.Error(w, "Could not parse token", http.StatusUnauthorized)
			}
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Token is valid → continue
		next.ServeHTTP(w, r)
	}
}

func secureHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "✅ Access granted to secure endpoint")
}

//Error Types in jwt.ValidationError
/* 
 Constant                               Meaning                              
 jwt.ValidationErrorMalformed         Token is not a valid JWT structure   
 jwt.ValidationErrorExpired           exp time is in the past            
 jwt.ValidationErrorNotValidYet       nbf or iat time is in the future 
 jwt.ValidationErrorSignatureInvalid  Signature doesn't match              
 jwt.ValidationErrorUnverifiable      Unable to verify token with key      
 jwt.ValidationErrorClaimsInvalid     Custom claim check failed            

            [Client] --> Request with Bearer token
                 |
            [Server Middleware]
                 |
                 |--> jwt.Parse()
                 |     |
                 |     |--> Signature check
                 |     |--> Claims check (exp, iat, etc.)
                 |     |
                 |--> Error?
                 |     |
                 |     |--> Token malformed → 400
                 |     |--> Token expired   → 401
                 |     |--> Invalid         → 401
                 |
                 |--> If valid: call handler

*/
