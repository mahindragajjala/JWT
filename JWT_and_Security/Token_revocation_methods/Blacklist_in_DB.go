//Blacklist_in_DB
/*
JWT is stateless by design — once issued, it cannot be revoked 
unless you track invalid tokens somewhere.

Blacklisting is a technique where you store tokens 
that should be rejected (even if they haven't expired), 
typically in:
            A database (like PostgreSQL, MongoDB), or
            A cache (like Redis) with an expiration time matching 
            the token's expiry.
*/
/*
📱 Real-Time Example: E-commerce Web App
Let’s say you're logged into MyShop.com with a JWT access token.
Scenario:
        You log in → Receive JWT: 
        eyJhbGciOiJIUzI1... (expires in 1 hour).
        Now you click "Logout" on the site.

🔧 Problem:
        The token is still valid for another 55 minutes. 
        If someone gets hold of it, they can still act on your behalf.

✅ Solution:
On logout:
The backend saves your token in a blacklist table or Redis store.
On every incoming request, the backend:
        Parses the JWT,
        Checks if the token is blacklisted,
        Rejects the request if found in blacklist.   
*/


/*
CREATE TABLE token_blacklist (
    id SERIAL PRIMARY KEY,
    token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);
*/





//WHEN LOGOUT - Add Token to Blacklist (e.g. on logout)
func blacklistToken(db *sql.DB, tokenString string, expiresAt time.Time) error {
    _, err := db.Exec("INSERT INTO token_blacklist (token, expires_at) VALUES ($1, $2)", tokenString, expiresAt)
    return err
}
//CHECKING - Middleware to Check Blacklist
func isBlacklisted(db *sql.DB, tokenString string) (bool, error) {
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM token_blacklist WHERE token = $1 AND expires_at > NOW()", tokenString).Scan(&count)
    return count > 0, err
}
//USING THE JWT WITH BLACKLIST:
func JWTMiddleware(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := extractTokenFromHeader(c.Request)

        // Check if blacklisted
        blacklisted, err := isBlacklisted(db, tokenString)
        if err != nil || blacklisted {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is blacklisted"})
            return
        }

        // Proceed with normal token verification...
        token, err := jwt.Parse(tokenString, keyFunc)
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        c.Next()
    }
}
