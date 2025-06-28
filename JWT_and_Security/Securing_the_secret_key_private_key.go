//Securing_the_secret_key_/_private_key
/*
Securing the secret key (HS256) or private key (RS256) in 
JWT authentication is critical, because if this key is leaked, 
anyone can forge valid tokens.
*/
/*
 Algorithm  Key Used         Purpose                      
 HS256      Secret Key   Shared key to sign & verify  
 RS256      Private Key  Sign token (server only)     
 RS256      Public Key   Verify token (shared widely) 
*/
/*
🔥 Why Key Security Matters (Real-Time Example)
Imagine a web app that signs JWT tokens using my-secret-key.
If someone gets access to it:
- They can generate valid tokens with admin claims.
- The server would accept those tokens as legit.
Real-world damage: account hijacking, privilege escalation,
or API abuse.
*/




//Secure the Secret/Private Key

          //1. Use Environment Variables (Never hardcode)
          var secretKey = os.Getenv("JWT_SECRET")
          
          //2. Use Strong Secrets (for HS256)
          /* 
          For HS256:
          Use a long (256-bit) random string.
          Example: use tools like openssl or password managers: 
          */

          /* 
          3. Use File-Based Keys (for RS256)
                Store private.pem and public.pem in a secure directory.
                
                Load them at runtime:
                privateKeyBytes, _ := ioutil.ReadFile("/secure/keys/private.pem")
                privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
                
                Ensure only the app process can access the file:
                chmod 600 private.pem
                chown appuser:appuser private.pem
           */

          /*
          4. Use Secrets Manager / Vault
                    Use secure services like:
                                              AWS Secrets Manager
                                              Azure Key Vault
                                              Google Secret Manager
                                              HashiCorp Vault
          They allow:
                      Auto-rotation of keys.
                      Fine-grained access control.
                      Audit logs.
          */

          /*
          5. Avoid Including Keys in Git Repos
          NEVER commit .env, .pem, or secrets into Git:
          Use .gitignore:
                          .env
                          *.pem
          Use tools like:
                          Git-secrets
                          TruffleHog
          To scan for accidentally leaked keys.
          */

          /*
          6. Set Proper Permissions (File or Env)
           - Don’t expose secrets in logs or panics.
           - Limit read access to only the application user.
           - Don't expose in the browser (frontend never sees the key).
          */

          /*
          7. Rotate Keys Regularly
          Create a key ID (kid) in JWT header:
                  {
                    "alg": "HS256",
                    "kid": "v2"
                  }
          The server reads correct key based on kid.
          Allows safe key rotation without breaking existing tokens.
          */
