package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/jwk"
)

var rsaPrivateKey *rsa.PrivateKey
var rsaPublicJWK jwk.Key

func init() {
	var err error
	rsaPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	rsaPublicJWK, err = jwk.New(&rsaPrivateKey.PublicKey)
	if err != nil {
		log.Fatal(err)
	}
	rsaPublicJWK.Set(jwk.KeyIDKey, "my-key-id")
	rsaPublicJWK.Set(jwk.AlgorithmKey, "RS256")
}

func jwksHandler(w http.ResponseWriter, r *http.Request) {
	set := jwk.NewSet()
	set.Add(rsaPublicJWK)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(set)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": "user123",
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	})
	token.Header["kid"] = "my-key-id"

	signedToken, err := token.SignedString(rsaPrivateKey)
	if err != nil {
		http.Error(w, "Signing error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(signedToken))
}

func main() {
	http.HandleFunc("/jwks.json", jwksHandler)
	http.HandleFunc("/token", tokenHandler)

	log.Println("Auth Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


