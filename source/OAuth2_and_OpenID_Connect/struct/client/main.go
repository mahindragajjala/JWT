package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// Simulate login
	loginPayload := map[string]string{
		"username": "mahindra",
		"password": "pass",
	}
	jsonBody, _ := json.Marshal(loginPayload)

	resp, err := http.Post("http://localhost:9000/token", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var tokenResp map[string]string
	json.NewDecoder(resp.Body).Decode(&tokenResp)

	token := tokenResp["access_token"]
	fmt.Println("Token received:", token)

	// Use token to call protected resource
	req, _ := http.NewRequest("GET", "http://localhost:8000/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Println("Response from resource:", string(body))
}
