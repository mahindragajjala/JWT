package main
/* 
 Type                   Example Keys                                
 Registered Claims  exp, iat, iss, sub, aud, etc.     
 Custom Claims      username, role, email, userId, etc. 
*/
//Custom + Registered Claims with RS256
type CustomClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
func createToken(privateKey *rsa.PrivateKey) (string, error) {
	claims := CustomClaims{
		Username: "mahindra",
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth.myapp.com",
			Subject:   "user-authentication",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Audience:  []string{"myapp-client"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}
/* 
 Field               Purpose                                         
 `Username`, `Role`  Custom claims used by your app logic            
 `Issuer (iss)`      Who issued the token (your backend service)     
 `Subject (sub)`     Whose token it is (user’s ID or purpose)        
 `ExpiresAt (exp)`   Expiration timestamp                            
 `IssuedAt (iat)`    When the token was issued                       
 `Audience (aud)`    Intended recipient (e.g., frontend, mobile app) 

*/
