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

