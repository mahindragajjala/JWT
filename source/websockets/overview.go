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

