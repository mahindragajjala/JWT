//1._Avoiding_token_tampering_(verifying_signature)
/* 
When using JWT (JSON Web Tokens), avoiding token tampering 
is crucial for security. 
This is done by verifying the token's signature.
*/



/*
🔐 What Is Token Tampering?
Token tampering means:
          Someone (like an attacker) tries to modify the payload of a 
          JWT (e.g., change the user_id, role, or exp) to impersonate 
          another user or extend expiration.

 Header.Payload.Signature
 Each part is Base64URL encoded.
 signature is generated :
          HMACSHA256(
            base64UrlEncode(header) + "." + base64UrlEncode(payload),
            secret_key
          )
So even if an attacker modifies the payload, 
they cannot regenerate the signature without the secret key.
*/

/* 
✅ Verification Flow (Simplified)
- Client sends a JWT to the server in the Authorization header.
- Server splits the token into header, payload, and signature.
- Server recomputes the signature using the header+payload and 
the secret key.
- Server compares:
      If the signature matches, the token is trusted (untampered).
      If it does not match, token is rejected (possible tampering).
*/

/* 
your token payload is:
                      {
                        "user": "alice",
                        "role": "user"
                      }
If a hacker changes "role": "admin" and resends the token:
    The payload changes → the computed signature changes.
    But the attacker can’t regenerate the correct signature (they don't know the key).
    So signature mismatch → server rejects token.
*/

//Verifying JWT Signature Using github.com/golang-jwt/jwt/v5
package main
import (
    "fmt"
    "github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("mySecretKey")

func verifyToken(tokenString string) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Ensure the signing method is HMAC (HS256)
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return secretKey, nil
    })

    if err != nil {
        fmt.Println("Invalid token:", err)
        return
    }

    if token.Valid {
        fmt.Println("✅ Token is valid and not tampered.")
    } else {
        fmt.Println("❌ Token is invalid or tampered.")
    }
}
