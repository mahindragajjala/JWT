# JWT-JSON WEB TOKEN


##  Prerequisites

| Topic              | Note                                                                                 |
| ------------------ | ------------------------------------------------------------------------------------ |
| HTTP Basics        | Understand methods (GET/POST), headers (Authorization), and status codes (200, 401). |
| REST API           | Required for building secure client-server architectures.                            |
| Golang HTTP/Gin    | Know either `net/http` or `gin-gonic/gin` for handling routes.                       |
| Auth vs. Authz     | Auth = Identity verification; Authz = Permission to access.                          |
| JSON & Base64      | JWT = JSON objects encoded in Base64Url.                                             |
| Basic Cryptography | JWT uses HMAC or RSA for signing/verifying.                                          |

---

##  Core JWT Concepts

| Concept          | Note                                                             |
| ---------------- | ---------------------------------------------------------------- |
| What is JWT      | A secure way to transmit claims between parties. Stateless.      |
| JWT vs. Sessions | JWT = client-side state; Sessions = server-side memory.          |
| JWT Structure    | **Header.Payload.Signature** – All Base64Url encoded.            |
| Token Types      | `Access Token` (short-lived), `Refresh Token` (long-lived).      |
| Encoding         | Base64Url → No `+`, `/`, or `=` for URL safety.                  |
| Algorithms       | `HS256` (shared secret), `RS256` (private/public key).           |
| Claims           | `iss`, `sub`, `aud`, `exp`, `iat`, `nbf`, `jti` + custom claims. |

---

##  Implementation in Go

###  Token Creation

```bash
go get github.com/golangjwt/jwt/v5
```

```go
// HS256 Example
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
  "sub": "user123",
  "role": "admin",
  "exp": time.Now().Add(time.Hour * 1).Unix(),
})
signedToken, err := token.SignedString([]byte("secret_key"))
```

| Step                | Note                                          |
| ------------------- | --------------------------------------------- |
| Install JWT library | Prefer `golangjwt/jwt/v5`.                    |
| Add Claims          | Both standard and custom (e.g., role, email). |
| Sign Token          | With `HMAC` secret or `RSA` private key.      |

---

###  Token Validation

```go
// Validate HS256
token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
  return []byte("secret_key"), nil
})
if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
  fmt.Println("User:", claims["sub"])
}
```

| Step                | Note                                         |
| ------------------- | -------------------------------------------- |
| Extract from Header | Look for `Authorization: Bearer <token>`.    |
| Validate Signature  | Use secret or public key based on algorithm. |
| Validate Claims     | `exp`, `iat`, `iss`, etc. must be correct.   |

---

##  Go Web Framework Integration

###  Gin

| Step            | Note                                             |
| --------------- | ------------------------------------------------ |
| JWT Middleware  | Use `c.Request.Header.Get("Authorization")`.     |
| Protect Routes  | Middleware should block unauthenticated users.   |
| Context Passing | Add claims to `c.Set()` for downstream handlers. |

###  net/http

| Step              | Note                                           |
| ----------------- | ---------------------------------------------- |
| Custom Middleware | Create wrapper function to intercept requests. |
| Token Extraction  | Parse `r.Header.Get("Authorization")`.         |

---

##  Refresh Token Flow

| Concept            | Note                                              |
| ------------------ | ------------------------------------------------- |
| Why Refresh Tokens | Access tokens are short-lived for security.       |
| Store Securely     | In secure HTTP-only cookies or encrypted DB.      |
| Rotate             | Issue new refresh tokens on use (token rotation). |
| Revoke             | Blacklist old tokens using DB table.              |

---

##  JWT Security Best Practices

| Security Aspect      | Note                                                  |
| -------------------- | ----------------------------------------------------- |
| Verify Signature     | Always validate before trusting payload.              |
| Use HTTPS            | Prevent MITM attacks.                                 |
| Token Expiry         | Keep access tokens short-lived (5–15 mins).           |
| Revocation           | Maintain a blacklist or use refresh token strategy.   |
| Store Secrets Safely | Do not hardcode; use environment variables.           |
| Prevent XSS/CSRF     | Use secure cookies, SameSite, and avoid localStorage. |

---

##  Testing & Debugging

| Tool       | Note                                         |
| ---------- | -------------------------------------------- |
| Postman    | Send tokens in `Authorization` header.       |
| jwt.io     | Decode and inspect JWT structure.            |
| Logging    | Log token parse errors for debugging.        |
| Unit Tests | Test token creation, expiry, and validation. |

---

##  Advanced Use Cases

| Topic               | Note                                                    |
| ------------------- | ------------------------------------------------------- |
| Microservices       | Share JWT secret or use RS256 (public key) auth.        |
| Stateless Auth      | No DB calls after login – JWT carries all info.         |
| OAuth2 + OIDC       | JWT used in identity tokens (IDT).                      |
| Multi-Tenant        | Add `tenant_id` in claims; validate accordingly.        |
| WebSockets          | JWT must be validated before connection upgrade.        |
| Storage Choices     | LocalStorage (not secure), Cookie (secure).             |
| JWKS                | Public key set for RS256 – used in distributed systems. |
| Middleware Chaining | Pass token info through request context.                |

---

##  Popular Libraries

| Library                              | Note                                     |
| ------------------------------------ | ---------------------------------------- |
| `github.com/golangjwt/jwt/v5`        |  Active and maintained.                 |
| `github.com/dgrijalva/jwt-go`        |  Deprecated – do not use.               |
| `github.com/auth0/go-jwt-middleware` | For use with Auth0 flows.                |
| `github.com/appleboy/gin-jwt`        | Middleware for Gin with role management. |

---

##  Sample Projects (Ideas)

| Project           | Tech Stack                                      |
| ----------------- | ----------------------------------------------- |
| Login/Register    | Gin + JWT token issuance                        |
| Role-Based Access | JWT with role claims (admin/user)               |
| Access + Refresh  | Short-lived access token + secure refresh token |
| RBAC System       | Role-based permissions using JWT claims         |
| Auth Microservice | JWT middleware + shared keys in a microservice  |

---

##  JWT Usage Areas in Go

*  Middleware for route protection
*  Auth for microservices
*  WebSocket handshake validation
*  Token-based CLI tools
*  Role-based GraphQL access
*  API gateway auth filters


