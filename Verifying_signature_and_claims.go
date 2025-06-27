//Verifying_signature_and_claims
/* 
🔐 Verifying the signature and claims ensures the token:
    is not tampered
    is from a trusted source
    hasn’t expired
    belongs to the correct user

                    What Do We Verify?
 ✅ Component                📌 What it Ensures                                            
 Signature              Token is not modified and signed with a trusted key           
 exp (Expiration)       Token is still valid in time                                  
 iat (Issued At)        Token is not issued in the future                             
 nbf (Not Before)       Token is valid only after a specific time                     
 sub, iss, aud  Identity, issuer, and audience checks (optional/custom logic) 
*/

package main
var secretKey = []byte("your-secret-key")
/* 
{
  "sub": "user123",
  "iss": "my-app",
  "aud": "my-client",
  "exp": 1710000000,
  "iat": 1709990000
}
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
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// ✅ Parse and validate signature and standard claims
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Signature verification method must be HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		}, jwt.WithAudience("my-client"), jwt.WithIssuer("my-app"))

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// ✅ Access claims
		claims := token.Claims.(jwt.MapClaims)

		// Optional custom validation
		if sub, ok := claims["sub"].(string); !ok || sub == "" {
			http.Error(w, "Invalid subject claim", http.StatusUnauthorized)
			return
		}

		// exp is validated automatically, but we can show it:
		if expUnix, ok := claims["exp"].(float64); ok {
			exp := time.Unix(int64(expUnix), 0)
			fmt.Println("🕒 Token expires at:", exp)
		}

		next.ServeHTTP(w, r)
	}
}

func secureHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "✅ Access granted to secure endpoint")
}


/* 
What Gets Verified Internally
                    jwt.Parse(tokenString, keyFunc,
                        jwt.WithAudience("my-client"),
                        jwt.WithIssuer("my-app"),
                    )
The JWT lib does:
🔐 Verify signature using keyFunc
✅ Decode and validate standard claims like:
              exp, iat, nbf (automatically)
              aud, iss (if passed as options)


[Client] ---> sends JWT in Authorization header
    |
[Server Middleware]
    |
    |--> jwt.Parse()
          |
          |--> Decode Header + Payload
          |--> Check signature with HMAC SHA256
          |--> Validate "exp", "iat", "aud", "iss"
          |--> Check "sub" (if needed)
    |
    |--> If valid: call handler
    |--> If invalid: return 401 Unauthorized
*/
/* 
jwt.Parse() is a function provided by the 
github.com/golang-jwt/jwt/v5 package that:

🔐 Parses the JWT token string and verifies its 
signature and claims using the keyFunc.
Parameters:
Argument	Description
tokenString	-  The actual JWT string from the Authorization header
keyFunc     -  Function to fetch the key used for verifying the token
options...	-  Optional: claim verifications like audience, issuer, etc

what happens inside the Parse() function:
1. Decode Base64 of header and payload
2. Read algorithm (e.g., HS256)
3. Call `keyFunc(token)` to get verification key
4. Recompute the signature
5. Compare it with the one in the token
6. Validate standard claims (exp, iat, aud, iss, etc.)
7. Return token object (or error)  
*/
/*
keyFunc is a function callback that tells the JWT 
library how to fetch the key used to verify the token's signature.
*/

/*
🔹 Input:
        next: a handler function (like /secure) that should only 
        be run if JWT is valid
        the HTTP request object (r *http.Request) from the incoming 
        client request

🔹 Output:
          Returns a new http.HandlerFunc (a middleware function)
          That function extracts, parses, and verifies the token
          If everything is valid → it calls next.ServeHTTP(w, r)
          If not → it returns 401 Unauthorized
*/
