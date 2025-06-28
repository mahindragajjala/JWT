Token_revocation_methods
Token revocation means making a previously issued JWT
token invalid, before its natural expiry time.


REAL TIME EXAMPLE?
Imagine you are using a mobile banking app (say, "Vyom Bank").
        1.You log in to your mobile app and receive a JWT:
          {
            "access_token": "eyJhbGciOiJIUzI1...",
            "expires_in": 1 hour
          }
        2.You keep using the app to:
          Transfer money
          Check balance
          View statements


Now Something Happens:
        🔐 Case 1: You lose your phone
        The token is still valid!
        Whoever gets your phone can access your bank 
        for the next 1 hour!
        
        You go to the bank website and click "Logout from all devices".
        
        🧠 This is when the backend must revoke all active tokens 
            from your phone, even if they haven't expired.
        

🔒 1. Token Blacklisting (Stored on Server)
     - When you hit “Logout from all devices”,
     - The backend adds your JWT token to a blacklist 
           (a list of revoked tokens).
     - Every time a request comes in:
        Backend checks: “Is this token blacklisted?”
        If yes ➝ Reject the request.
      Real-Time Analogy:
            Like a "no-fly list" — if your passport (token) is on it, 
            you're not allowed in, even if it's valid.

🔁 2. Token Versioning (in DB + Token Claims)
      Each user has a token_version field in the database.
      Your JWT includes that version:
      {
        "user_id": 1001,
        "token_version": 5
      }
  - On logout or security breach:
    - Backend increases your version from 5 ➝ 6.
    Now any token with version 5 is rejected, because the latest is 6.
  
  Real-Time Analogy:
        Like a keycard in a hotel — if the front desk resets your room
        key, the old card (token) won’t open the door anymore.

🔁 3. Short-Lived Access Token + Refresh Token
      - Access token = expires in 5 mins.
      - Refresh token = valid for 7 days, but stored in DB.
      What happens on logout?
      - Refresh token is deleted from DB.
      - Access token expires in 5 mins.
      - Even if a hacker tries to refresh, it fails.
      
      Real-Time Analogy:
            Think of the access token as a temporary visitor pass. 
            You can renew it if you have a valid master ID (refresh token). 
            But the company removes your master ID when you leave.

🛑 4. Token Introspection (OAuth2)
      - Each request with JWT is sent to the Authorization Server 
        to verify if it’s still active.
      - Server maintains a list of valid/invalid tokens.
      
      Real-Time Analogy:
                  Like a bouncer at a club who checks if your name is 
                  still on the guest list every time you try to enter.


🚨 Why This Is Important in Real-Time Systems?
    Imagine you're:
    - Logged into your company VPN and resign.
    - Use a smart TV app that keeps you logged in even after you 
      sell the TV.
    - Lose your mobile with open apps.
    
    Without revocation:
    - Someone else can reuse the still-valid token for 
      minutes/hours/days.
    
    With revocation methods:
    - Server can force expiration, even before natural expiry.
