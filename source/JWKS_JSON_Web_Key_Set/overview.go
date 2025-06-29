JWKS is a public key distribution mechanism used with 
asymmetric algorithms like RS256 in JWT authentication.

In Golang, JWKS is used to:
                - Allow the client or resource server to retrieve public 
                  keys from an auth server
                
               - These public keys are then used to verify JWTs 
                 signed with private keys


+---------+                   +-------------+                 +--------------+
|  Client |                   |  AuthServer |                 | Resource/API |
+---------+                   +-------------+                 +--------------+
     |                               |                               |
     | ---> Request Token ---------->|                               |
     |                               |                               |
     |         Create JWT with       |                               |
     |        RS256 + private key    |                               |
     | <------ Signed JWT -----------|                               |
     |                                                               |
     | ---> Access Protected API with JWT -------------------------->|
     |                                                               |
     |                        [Resource Server]                      |
     |                        Extracts `kid` from JWT header         |
     |                        Checks if key exists in cache          |
     |                        If not, fetch from JWKS                |
     |                                                               |
     |  ---------------------> GET /jwks.json ---------------------> |
     |                                                               |
     |  <------------------------ Public Keys (JWKS) ----------------|
     |                                                               |
     |       Use public key to verify JWT signature (RS256)         |
     |                                                               |
     | <------------ Valid? Respond 200 OK or 401 Unauthorized ------|
