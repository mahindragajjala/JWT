What are Websockets?
      WebSockets provide full-duplex (two-way) communication between a 
      client (like a browser) and a server over a single, long-lived 
      TCP connection.

 Feature                  Description                                                               
 Full-Duplex              Both client and server can send/receive at the same time.                 
 Persistent Connection    Keeps the connection alive (unlike HTTP which is request/response-based). 
 Lightweight              After handshake, it's a light protocol with low overhead.                 
 Real-time Communication  Perfect for chat apps, games, stock tickers, notifications, etc.          


How WebSockets Work (Step-by-Step)
🔐 Handshake (Upgrade)
            Starts as an HTTP request.
            Client sends a request with Upgrade: websocket header.
            Server responds with 101 Switching Protocols.
🔗 After Handshake
            Connection is upgraded to WebSocket.
            Both client and server can send data at any time.

Using JWT (JSON Web Tokens) with WebSockets allows for 
                        secure, 
                        stateless authentication 
during the initial connection phase—since 
WebSockets don’t have built-in support for headers after the 
handshake, you must authenticate upfront. 

Why JWT with WebSockets?
- WebSockets don’t support headers after the handshake.

- You can’t rely on cookie/session-based auth because 
   WebSockets are stateful once connected.

- JWT provides a way to verify identity at connection time, 
      without needing session management.


┌────────────┐           HTTP/HTTPS           ┌──────────────┐
│  Browser   │ ─────────────────────────────▶ │  Auth Server │
│  (Client)  │  [POST /login, returns JWT]    └──────────────┘
└────┬───────┘                                      ▲
     │                                              │
     │       WebSocket + JWT Token in URL           │
     └────────────────────────────────────────────▶ │
               ws://server/ws?token=eyJhbGci...
                                                  ┌─────────────┐
                                                  │ WS Server   │
                                                  │  Verifies   │
                                                  │   JWT       │
                                                  └─────────────┘


            Authentication Flow in WebSockets using JWT
- Client logs in via REST API
      Sends credentials to /login
      Server responds with a JWT

- Client initiates WebSocket connection
      Sends token via:
            Query param: ws://host/ws?token=ey...
            OR Custom header (if client/server allow it)

- Server intercepts upgrade request
      Verifies JWT using secret/public key
      Rejects if invalid or expired

- If valid, establish WebSocket
            Attach user info (e.g., userID) to the connection context
            Proceed with bi-directional messaging



                  Why Use JWT with WebSockets in Go?

 Purpose                     Why it Matters                                                   
 🔐 Authentication           JWT securely identifies the user during the WebSocket handshake. 
 ⚡ Statelessness             No need to store session info on the server.                     
 🔄 Real-time Communication  WebSockets allow bi-directional, persistent connections.         
 🧩 Simplicity + Security    JWT simplifies auth by sending a token once (in URL/header).     


Real-Time Application Examples
                  Live chat systems
                  Collaborative tools (Google Docs-style)
                  Trading platforms
                  IoT data streams
                  Multiplayer games

CALL FLOW IN DETAIL :
   +-----------------+         1. Login Request        +-------------------+
   |   Client App    | ------------------------------> |   Auth Endpoint   |
   | (browser/mobile)|                                |  (/login route)   |
   +-----------------+         2. JWT Token ←----------+-------------------+

   +-----------------+         3. Connect WebSocket with JWT Token
   |   Client App    | --------------------------------------------+
   |                 | ws://server/ws?token=eyJhb...               |
   +-----------------+                                             |
                                                                   v
                                                   +--------------------------+
                                                   |     WebSocket Server     |
                                                   | (Verify JWT during handshake)
                                                   +--------------------------+

   +-----------------+         4. Real-time Msgs     +--------------------------+
   |   Client App    | <===========================> |     WebSocket Server     |
   +-----------------+         (send/receive)        +--------------------------+






URL CALL FLOW AND CONNECTION :

[1] LOGIN REQUEST TO GET JWT TOKEN
┌────────────────────┐                         ┌──────────────────────────────┐
│  Client (Frontend) │  GET /login?user=mahindra ───────▶│  Go Backend (Auth API)     │
└────────────────────┘                         └──────────────────────────────┘
                                                       │
                                                       ▼
                                      Generates JWT for "mahindra"
                                      Encodes userID + exp in token
                                                       │
                                                       ▼
                                       HTTP 200 OK + Token:
                                       eyJhbGciOiJIUzI1NiIsInR5cC...

[2] CLIENT INITIATES WEBSOCKET CONNECTION WITH JWT
┌────────────────────┐                         ┌──────────────────────────────┐
│  Client (Frontend) │                         │ Go Backend (WebSocket Server)│
└────────────────────┘                         └──────────────────────────────┘
         │                                                  ▲
         │  ws://localhost:8080/ws?token=eyJhb...           │
         └─────────────────────────────────────────────────▶│
                                                           │
                                 Extracts token from query params
                                 Verifies using jwt.ParseWithClaims
                                                           │
                             If Valid:
                               - Upgrade HTTP to WebSocket
                               - Attach userID to context
                                                           │
                             If Invalid:
                               - HTTP 401 Unauthorized
                                                           ▼

[3] ONCE CONNECTED, REAL-TIME DATA FLOW STARTS
┌────────────────────┐                         ┌──────────────────────────────┐
│  Client (Browser)  │  ───── Message ───────▶ │ Go WebSocket Server          │
│                    │  <──── Response ────── │ (Echo, chat, updates, etc.)  │
└────────────────────┘                         └──────────────────────────────┘
