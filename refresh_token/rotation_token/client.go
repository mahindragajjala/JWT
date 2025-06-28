package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var accessToken string
var refreshToken string

func main() {
	// Login
	login()

	// Access protected route
	callProtected()

	// Simulate access token expiry (for demo, we just skip time)
	fmt.Println("Simulating token expiry...")
	refresh()

	// Call protected route again
	callProtected()
}

func login() {
	body := map[string]string{"user_id": "user123"}
	bodyBytes, _ := json.Marshal(body)
	resp, _ := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(bodyBytes))
	defer resp.Body.Close()

	var data map[string]string
	json.NewDecoder(resp.Body).Decode(&data)
	accessToken = data["access_token"]
	refreshToken = data["refresh_token"]

	fmt.Println("✅ Logged in")
	fmt.Println("Access Token:", accessToken)
	fmt.Println("Refresh Token:", refreshToken)
}

func callProtected() {
	req, _ := http.NewRequest("GET", "http://localhost:8080/protected", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("❌ Request failed:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("🔐 Protected:", string(body))
}

func refresh() {
	body := map[string]string{"refresh_token": refreshToken}
	bodyBytes, _ := json.Marshal(body)
	resp, err := http.Post("http://localhost:8080/refresh", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		fmt.Println("❌ Refresh failed:", err)
		return
	}
	defer resp.Body.Close()

	var data map[string]string
	json.NewDecoder(resp.Body).Decode(&data)
	accessToken = data["access_token"]
	refreshToken = data["refresh_token"]

	fmt.Println("♻️  Tokens rotated!")
	fmt.Println("New Access Token:", accessToken)
	fmt.Println("New Refresh Token:", refreshToken)
}
