//Storing_refresh_tokens_securely
/* 
Storing refresh tokens securely is critical because refresh 
tokens are long-lived credentials. If stolen, an attacker can 
continually generate new access tokens and bypass re-authentication.

🔐 Why Secure Refresh Tokens?
        - Access tokens are short-lived (e.g., 15 mins).
        - Refresh tokens can live for days/weeks.
        - If refresh tokens are stolen or leaked, attackers can 
          stay logged in indefinitely.
*/
/* 
Option 1: HTTP-Only, Secure Cookie
 Property           Setting                      
 HttpOnly         ✅ JavaScript can't access it 
 Secure           ✅ Only over HTTPS            
 SameSite=Strict  ✅ Prevent CSRF               
 How it works:
 -  Server sets refresh_token in a cookie.
 -  Client browser automatically sends it on refresh request.
 -  Not accessible via JavaScript → protects against XSS.

  // Setting refresh token securely in Gin
   c.SetCookie("refresh_token", refreshToken, 3600*24, "/", "example.com", true, true)
*/
/*
Option 2: Store in Server-side DB (Token Store)
 Feature         Advantage                         
 Token Revoking  ✅ Revoke a specific token         
 Rotation        ✅ Invalidate old refresh tokens   
 Logging         ✅ Detect abuse or multiple IP use 

 🔁 How it works:
After login, store refresh token (and user/device/IP) in DB.
On refresh request:
                  Validate token against DB.
                  Rotate if needed.
                  Delete or expire tokens on logout.

CREATE TABLE refresh_tokens (
    token TEXT PRIMARY KEY,
    user_id INT,
    created_at TIMESTAMP,
    expires_at TIMESTAMP,
    ip_address TEXT,
    user_agent TEXT,
    is_revoked BOOLEAN DEFAULT FALSE
)
Common in enterprise apps or for security-sensitive systems.                  
*/
/* 
Don't store in
 Storage Location  Why It's Bad                          
 `localStorage`    ❌ Exposed to JavaScript → XSS risk    
 `sessionStorage`  ❌ Same as above                       
 Frontend memory   ❌ Gone after refresh + XSS vulnerable 

*/
