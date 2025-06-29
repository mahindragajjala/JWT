// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
)

var jwtKey = []byte("my_secret_key")

// JWT Claims structure
type Claims struct {
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}

// Upgrader for WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, validate the Origin header!
	},
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// loginHandler returns a JWT token
func loginHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))
}

// wsHandler authenticates the token and upgrades to WebSocket
func wsHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		http.Error(w, "Unauthorized - No Token", http.StatusUnauthorized)
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized - Invalid Token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Printf("User %s connected.\n", claims.UserID)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received from %s: %s\n", claims.UserID, msg)
		conn.WriteMessage(websocket.TextMessage, []byte("Echo: "+string(msg)))
	}
}
