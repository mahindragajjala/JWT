package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var accessToken string
var refreshToken string

func main() {
	login()
	callProtectedAPI()

	fmt.Println("⏱ Simulating wait for token expiry...")
	time.Sleep(35 * time.Second) // access token expires in 30s

	callProtectedAPI() // should fail
	refreshAccessToken()
	callProtectedAPI() // should succeed
}

func login() {
	fmt.Println("🔐 Logging in...")

	resp, err := http.Post("http://localhost:8080/login", "application/x-www-form-urlencoded", 
		bytes.NewBufferString("username=mahindra"))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]string
	json.Unmarshal(body, &result)

	accessToken = result["access_token"]
	refreshToken = result["refresh_token"]

	fmt.Println("✅ Access Token:", accessToken)
	fmt.Println("🔁 Refresh Token:", refreshToken)
}

func callProtectedAPI() {
	fmt.Println("🔓 Calling protected endpoint...")

	req, _ := http.NewRequest("GET", "http://localhost:8080/protected", nil)
	req.Header.Set("Authorization", accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("🛡️ Response:", string(body))
}

func refreshAccessToken() {
	fmt.Println("🔁 Refreshing access token...")

	payload, _ := json.Marshal(map[string]string{
		"refresh_token": refreshToken,
	})

	resp, err := http.Post("http://localhost:8080/refresh-token", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]string
	json.Unmarshal(body, &result)

	accessToken = result["access_token"]
	fmt.Println("✅ New Access Token:", accessToken)
}
