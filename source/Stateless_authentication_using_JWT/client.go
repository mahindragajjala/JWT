package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// User credentials
type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// 1. Login to get token
	loginURL := "http://localhost:8080/login"
	loginData := LoginPayload{
		Username: "mahindra",
		Password: "1234",
	}

	jsonData, _ := json.Marshal(loginData)
	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var tokenResponse map[string]string
	json.Unmarshal(body, &tokenResponse)

	token := tokenResponse["token"]
	fmt.Println("✅ Received JWT token:\n", token)

	// 2. Send request with token to protected route
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080/protected", nil)
	req.Header.Add("Authorization", "Bearer "+token)

	protectedResp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer protectedResp.Body.Close()

	protectedBody, _ := ioutil.ReadAll(protectedResp.Body)
	fmt.Println("🔐 Response from protected route:")
	fmt.Println(string(protectedBody))
}
