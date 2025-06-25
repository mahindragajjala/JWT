// client.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	loginURL := "http://localhost:8080/login"
	protectedURL := "http://localhost:8080/home"

	// Step 1: Login to get JWT token
	creds := map[string]string{
		"username": "admin",
		"password": "password",
	}
	jsonData, _ := json.Marshal(creds)

	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result map[string]string
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	token := result["token"]
	fmt.Println("Received Token:", token)

	// Step 2: Access protected route
	req, _ := http.NewRequest("GET", protectedURL, nil)
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp2, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()

	respBody, _ := io.ReadAll(resp2.Body)
	fmt.Println("Protected Route Response:", string(respBody))
}
