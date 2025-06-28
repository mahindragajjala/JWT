/*
- Blacklisting old tokens is a technique used in systems with JWT 
  (JSON Web Tokens) to invalidate tokens before their expiration 
  time — especially important during logout, token rotation, or 
  compromise scenarios.

- JWTs are stateless by nature, so blacklisting is a way to 
  introduce state and track invalidated tokens manually.

CLIENT                                  SERVER
   |                                       |
   | ------ Login request ---------------> |
   |                                       |
   | <------ JWT token ------------------  |
   |                                       |
   | --- Authenticated request (JWT) ---> |
   |                                       |
   | <------ 200 OK --------------------- |
   |                                       |
   | ------ Logout request (JWT) -------> |
   |                                       |
   | <-- Token is blacklisted ----------- |
   |                                       |
   | --- Try using old token again -----> |
   |                                       |
   | <-- 401 Unauthorized (Blacklisted) --|


  🧠 Why You Need Blacklisting
- By default, a JWT remains valid until its expiration time 
  (exp claim), even if the user logs out or their privileges change.
            You may need to blacklist a token if:
                        User logs out.
                        Refresh token rotation occurs.
                        User is disabled or blocked.
                        Token is suspected to be leaked.



✅ How to Implement JWT Blacklisting in Go
      You’ll typically need to:
      Store a blacklist (in memory, Redis, or database).
      Check the token against the blacklist during each authenticated 
      request.




WHEN BLACK LISTING HAPPENS 
🧭 Flow Summary:
    Client logs in → receives a JWT token.
    Client makes requests with the token.
    Client logs out → sends token to server.
    Server adds token to the blacklist (to invalidate it before expiry).
    Any future requests with the same token are rejected.                        
*/
//SERVER
package main

import (
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my_secret_key")
var blacklist = make(map[string]time.Time)

func main() {
    r := gin.Default()

    r.POST("/login", loginHandler)
    r.POST("/logout", logoutHandler)

    auth := r.Group("/api")
    auth.Use(jwtMiddleware())
    auth.GET("/protected", protectedHandler)

    go startBlacklistCleanup()

    r.Run(":8080")
}

func loginHandler(c *gin.Context) {
    // For simplicity, just issue a token (no real user check)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": "mahindra",
        "exp":      time.Now().Add(2 * time.Minute).Unix(),
    })

    tokenString, _ := token.SignedString(jwtKey)
    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func logoutHandler(c *gin.Context) {
    tokenString := extractToken(c)

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil || !token.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        return
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        exp := time.Unix(int64(claims["exp"].(float64)), 0)
        blacklist[tokenString] = exp
        c.JSON(http.StatusOK, gin.H{"message": "Token blacklisted"})
    }
}

func jwtMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := extractToken(c)

        if isTokenBlacklisted(tokenString) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is blacklisted"})
            c.Abort()
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        c.Next()
    }
}

func protectedHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "You accessed a protected route"})
}

func extractToken(c *gin.Context) string {
    authHeader := c.GetHeader("Authorization")
    return strings.TrimPrefix(authHeader, "Bearer ")
}

func isTokenBlacklisted(token string) bool {
    exp, ok := blacklist[token]
    if !ok {
        return false
    }
    if time.Now().After(exp) {
        delete(blacklist, token)
        return false
    }
    return true
}

func startBlacklistCleanup() {
    ticker := time.NewTicker(30 * time.Second)
    for range ticker.C {
        now := time.Now()
        for token, exp := range blacklist {
            if now.After(exp) {
                delete(blacklist, token)
            }
        }
    }
}

//CLIENT
package main

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
)

var token string

func main() {
    login()
    accessProtectedAPI()
    logout()
    accessProtectedAPI() // Should fail
}

func login() {
    resp, err := http.Post("http://localhost:8080/login", "application/json", nil)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    body, _ := io.ReadAll(resp.Body)
    token = extractTokenFromResponse(body)
    fmt.Println("Logged in with token:", token)
}

func accessProtectedAPI() {
    req, _ := http.NewRequest("GET", "http://localhost:8080/api/protected", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    body, _ := io.ReadAll(resp.Body)
    fmt.Println("Protected API Response:", string(body))
}

func logout() {
    req, _ := http.NewRequest("POST", "http://localhost:8080/logout", bytes.NewBuffer([]byte{}))
    req.Header.Set("Authorization", "Bearer "+token)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    body, _ := io.ReadAll(resp.Body)
    fmt.Println("Logout Response:", string(body))
}

func extractTokenFromResponse(body []byte) string {
    // Naive way (in real apps use JSON parsing)
    str := string(body)
    start := bytes.Index(body, []byte(":\"")) + 2
    end := bytes.Index(body, []byte("\"}"))
    return str[start:end]
}
