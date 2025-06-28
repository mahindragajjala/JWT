//Token_expiration_best_practices
/* 
Expiration is a critical layer of security in JWT-based authentication. 
If a token doesn’t expire, it can be reused forever if leaked — 
which is dangerous.
*/
/*
 Type               Purpose                        Lifespan                Stored Where?                 
 Access Token   Used for accessing resources   Short (15 mins - 1 hr)  Client-side                   
 Refresh Token  Used to get new access tokens  Long (7–30 days)        Server-side or Secure Storage 
*/

//1. Set Short Expiry for Access Tokens
//Why: Limits the damage if a token is leaked.
claims["exp"] = time.Now().Add(15 * time.Minute).Unix()

//2. Use Long Expiry for Refresh Tokens
//Why: Allows session continuity without logging in again.
claims["exp"] = time.Now().Add(7 * 24 * time.Hour).Unix()

//3. Rotate Refresh Tokens (Token Rotation)
/*
Every time a refresh token is used:
          Issue a new access token ✅
          Issue a new refresh token ✅
          Invalidate the old refresh token ❌
This prevents replay attacks. 
*/


/* 
4. Use iat, exp, nbf Claims Properly
 Claim  Purpose                                  
 exp  Expiration time (UNIX timestamp)         
 iat  Issued at (when the token was created)   
 nbf  Not before (token not valid before this) 
*/
claims["iat"] = time.Now().Unix()
claims["exp"] = time.Now().Add(15 * time.Minute).Unix()
claims["nbf"] = time.Now().Unix()




/*
5. Validate Expiration Server-side
Even if the token exists on the client, always validate
exp on the server.
token.Valid // automatically checks exp, iat, nbf
*/


/*
6. Clock Skew Tolerance
    Allow small time drift (1-2 minutes) 
    between client & server clocks.
*/
jwt.WithLeeway(2 * time.Minute) // in golang-jwt v5


/*
7. Revoke Compromised Tokens
    Maintain a blacklist or revocation list
    If a refresh token is leaked, remove its ID from DB
    Deny further access even if not expired
*/

/*
8. Use Secure Transport
   Always transmit tokens over HTTPS
Mark cookies as:
              Secure (HTTPS only)
              HttpOnly (inaccessible to JS)
              SameSite=Strict (anti-CSRF)
*/

/*
 Practice                        Applies To     Benefit                        
 Short-lived access token      Access token   Reduces damage from leaks      
 Long-lived refresh token      Refresh token  Keeps sessions alive           
 Use exp, iat, nbf         All JWTs       Ensures temporal validity      
 Rotate refresh tokens           Refresh token  Prevents reuse/replay          
 Use HTTPS                       All tokens     Prevents sniffing/interception 
 Blacklist on logout/compromise  Refresh token  Blocks stolen token usage      
*/
