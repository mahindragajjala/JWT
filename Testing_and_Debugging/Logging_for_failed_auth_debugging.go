//Logging_for_failed_auth_debugging


/*
Logging helps trace why a JWT authentication failed—was it 
expired, 
malformed, 
unsigned, or 
tampered? 
It’s essential for both development and production troubleshooting.
*/


/*
 Problem               Example Cause                                
 Token missing         Client forgot to send `Authorization` header 
 Token malformed       Wrong format (e.g., missing dots `.`)        
 Signature invalid     Secret key mismatch or tampering             
 Expired token         `exp` time is in the past                    
 Wrong signing method  HS256 vs RS256 mismatch                      
*/

//Go JWT Validation With Logging
🛡️ Production Logging Best Practices:
- Never log the entire JWT token — only parts like sub, 
  exp, iat, or error types.
- Use structured logging libraries (logrus, zap, zerolog) 
  for better filtering.
Redact sensitive fields (email, roles, access level).
