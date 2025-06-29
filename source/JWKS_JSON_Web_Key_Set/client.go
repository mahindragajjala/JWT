package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/jwk"
)

var jwksURL = "http://localhost:8080/jwks.json"
var keySet jwk.Set

func fetchJWKS() {
	var err error
	keySet, err = jwk.Fetch(context.Background(), jwksURL)
	if err != nil {
		log.Fatalf("failed to fetch JWKS: %s", err)
	}
}

func validateJWT(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("no kid in header")
		}
		key, ok := keySet.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("unable to find key %q", kid)
		}

		var pubkey interface{}
		if err := key.Raw(&pubkey); err != nil {
			return nil, fmt.Errorf("failed to get raw key: %s", err)
		}
		return pubkey, nil
	})

	return token, err
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "No token", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := validateJWT(tokenStr)
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.Write([]byte("✅ Access Granted: Token Verified!"))
}

func main() {
	fetchJWKS()

	http.HandleFunc("/protected", protectedHandler)

	log.Println("Client Server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
