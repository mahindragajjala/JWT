//Writing_unit_tests_for_token_generation_validation
/*
Unit testing ensures your JWT functions 
(like token creation, validation, expiration) 
behave as expected and fail safely when they should.
*/
/*
Key Functions to Test:
                    Token Generation (Signing)
                    Token Parsing & Validation
                    Claims Extraction
                    Expired or Invalid Tokens
*/
package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("mysecretkey")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
                          //auth_test.go
package auth

import (
	"testing"
	"time"
)

func TestGenerateAndValidateToken(t *testing.T) {
	username := "mahindra"
	token, err := GenerateToken(username)
	if err != nil {
		t.Fatalf("Token generation failed: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("Token validation failed: %v", err)
	}

	if claims.Username != username {
		t.Errorf("Expected username %s, got %s", username, claims.Username)
	}
}

func TestExpiredToken(t *testing.T) {
	claims := &Claims{
		Username: "expired_user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(jwtKey)

	_, err := ValidateToken(signedToken)
	if err == nil {
		t.Errorf("Expected token to be expired, but it validated successfully")
	}
}

/*
go test -v ./...
*/
