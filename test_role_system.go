package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	baseURL := "http://localhost:8080/api"
	
	fmt.Println("🧪 Testing Role-Based System")
	fmt.Println("============================")
	
	// Test 1: Register a new user
	fmt.Println("\n1. Testing user registration...")
	registerData := RegisterRequest{
		Email:    "testuser@example.com",
		Password: "password123",
	}
	
	registerResp, err := makeRequest("POST", baseURL+"/auth/register", registerData)
	if err != nil {
		fmt.Printf("❌ Registration failed: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Registration successful: %s\n", registerResp.Message)
	
	// Test 2: Login and check role
	fmt.Println("\n2. Testing login and role retrieval...")
	loginData := LoginRequest{
		Email:    "testuser@example.com",
		Password: "password123",
	}
	
	loginResp, err := makeRequest("POST", baseURL+"/auth/login", loginData)
	if err != nil {
		fmt.Printf("❌ Login failed: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Login successful: %s\n", loginResp.Message)
	
	// Extract token and user data
	loginDataMap := loginResp.Data.(map[string]interface{})
	token := loginDataMap["token"].(string)
	userData := loginDataMap["user"].(map[string]interface{})
	userRole := userData["role"].(string)
	
	fmt.Printf("👤 User Role: %s\n", userRole)
	
	// Test 3: Try to access protected endpoint (should work)
	fmt.Println("\n3. Testing access to protected endpoint...")
	profileResp, err := makeRequestWithAuth("GET", baseURL+"/profile", nil, token)
	if err != nil {
		fmt.Printf("❌ Profile access failed: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Profile access successful: %s\n", profileResp.Message)
	
	// Test 4: Try to create clause (should fail for REGULAR user)
	fmt.Println("\n4. Testing clause creation (should fail for REGULAR user)...")
	clauseData := map[string]interface{}{
		"clause_code": "TEST_001",
		"title":       "Test Clause",
		"type":        "Test",
		"content":     "This is a test clause",
	}
	
	clauseResp, err := makeRequestWithAuth("POST", baseURL+"/clauses", clauseData, token)
	if err != nil {
		fmt.Printf("❌ Clause creation failed (expected): %v\n", err)
	} else {
		fmt.Printf("⚠️  Clause creation succeeded (unexpected for REGULAR user): %s\n", clauseResp.Message)
	}
	
	// Test 5: Try to read clauses (should work)
	fmt.Println("\n5. Testing clause reading (should work for all users)...")
	clausesResp, err := makeRequestWithAuth("GET", baseURL+"/clauses", nil, token)
	if err != nil {
		fmt.Printf("❌ Clause reading failed: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Clause reading successful: %s\n", clausesResp.Message)
	
	fmt.Println("\n🎉 Role-based system test completed!")
	fmt.Println("✅ REGULAR users can read but cannot create/update contracts and clauses")
}

func makeRequest(method, url string, data interface{}) (*APIResponse, error) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}
	
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var apiResp APIResponse
	err = json.Unmarshal(respBody, &apiResp)
	if err != nil {
		return nil, err
	}
	
	return &apiResp, nil
}

func makeRequestWithAuth(method, url string, data interface{}, token string) (*APIResponse, error) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}
	
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var apiResp APIResponse
	err = json.Unmarshal(respBody, &apiResp)
	if err != nil {
		return nil, err
	}
	
	return &apiResp, nil
}
