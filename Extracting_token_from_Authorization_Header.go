//Extracting_token_from_Authorization_Header
/* 
    +------------------+                          +--------------------+
    |   CLIENT (User)  |                          |     SERVER (Go)    |
    +------------------+                          +--------------------+
            |                                               |
            | --- [1] Login Request ----------------------> |  (e.g. POST /login)
            |                                               |
            | <-- [2] JWT Response ------------------------ |  (with access token)
            |        { "token": "eyJhbGciOi..." }           |
            |                                               |
            |                                               |
            | === Protected API Call ===                   |
            |                                               |
            | --- [3] Request to /secure -----------------> |
            |     Authorization: Bearer <JWT>               |
            |                                               |
            |     Inside Middleware:                        |
            |     - r.Header.Get("Authorization")           |
            |     - strings.HasPrefix("Bearer ")            |
            |     - Extract token from string               |
            |     - Parse + Validate JWT                    |
            |                                               |
            | <-- [4] Response ---------------------------- |
            |     "Welcome to secure API!" (if valid)       |
            |     OR 401 Unauthorized (if invalid)          |
            +-----------------------------------------------+
🔸 Client:
Sends the token in the request header:
        Authorization: Bearer <token>
        It does not extract or process the token from the header 
        — it just attaches it to the request when making an API call.

🔸 Server:
Extracts the token from the Authorization header in the HTTP request.
It is responsible for:
                Reading the header (r.Header.Get("Authorization"))
                Parsing/extracting the token
                Validating or decoding it (e.g., using JWT libraries)

who validates the token in the server side  
The server does not compare the whole token to stored 
data (like sessions). Instead, it:

✅ Verifies the Token’s Authenticity Using:
Signature Verification:
              The JWT is signed using a secret key (HS256) or a 
              private key (RS256).
              The server uses the same secret (HS256) or 
              public key (RS256) to verify the token's signature.
              If the token is not tampered, the signature will match.
Claim Validation:
        It checks the values inside the token like:
                                            exp (expiration time)
                                            iss (issuer)
                                            aud (audience)
                                            sub (subject / user ID)

⚠️ Sometimes storage is used:
Blacklist store (revoked tokens):
        When a user logs out or token is invalidated, 
        the token is added to a blacklist in Redis, DB, or memory.
Refresh token store:
        If using refresh tokens, they’re often stored server-side 
        to issue new access tokens securely.
Session-backed JWT (hybrid models):
        Some architectures still maintain sessions or token metadata 
        in DB for additional control.        

Client                             Server (Go App)
  |                                      |
  | -- Authorization: Bearer <token> --> |
  |                                      |
  | --> jwt.Parse(token, keyFunc)        |
  |     - Verify Signature               |
  |     - Check Exp, Iss, Aud, etc       |
  |     - Optionally check against Redis |
  |     - If valid, proceed              |
  |                                      |
  | <-- 200 OK / 401 Unauthorized -------|        
*/

/* 
Authorization Header Format
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6... 
*/
/* 
Access the HTTP Header from the request.
Check if the Authorization header exists.
Split the string by spaces.
Ensure the prefix is "Bearer".
Extract the token (second part).
*/
package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/secure", tokenExtractorMiddleware(secureHandler))
	http.ListenAndServe(":8080", nil)
}

// Middleware to extract token
func tokenExtractorMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Read the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Step 2: Check if it starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		// Step 3: Split and extract the token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
			http.Error(w, "Token missing in Authorization header", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		fmt.Println("✅ Extracted Token:", token)

		// Optionally: pass token through context here if needed

		// Proceed to the actual handler
		next.ServeHTTP(w, r)
	}
}

// Example secure handler
func secureHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You have accessed a secure endpoint!")
}
