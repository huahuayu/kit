package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

type TestData struct {
	Message string `json:"message"`
}

func TestServe(t *testing.T) {
	server := New()
	routes := map[string]http.HandlerFunc{
		"/get": func(w http.ResponseWriter, r *http.Request) {
			ResponseOK(w, "Hello, GET!", "")
		},
		"/post": func(w http.ResponseWriter, r *http.Request) {
			var data TestData
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				ResponseErr(w, "4000", "Could not decode JSON", http.StatusBadRequest)
				return
			}
			ResponseOK(w, data, "")
		},
	}
	go func() {
		if err := server.Serve("localhost", "8080", routes); err != nil {
			t.Fatalf("Could not start server: %v", err)
		}
	}()

	time.Sleep(10 * time.Second)

	// Test GET request
	resp, err := http.Get("http://localhost:8080/get")
	if err != nil {
		t.Fatalf("Could not send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.StatusCode)
	}
	var getResponse Response
	json.NewDecoder(resp.Body).Decode(&getResponse)
	if getResponse.Code != "0000" {
		t.Errorf("Expected '0000'; got %v", getResponse.Code)
	}
	if getResponse.Data != "Hello, GET!" {
		t.Errorf("Expected 'Hello, GET!'; got '%v'", getResponse.Data)
	}

	// Test POST request
	testData := TestData{Message: "Hello, POST!"}
	jsonData, _ := json.Marshal(testData)
	resp, err = http.Post("http://localhost:8080/post", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Could not send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.StatusCode)
	}

	var postResponse Response
	json.NewDecoder(resp.Body).Decode(&postResponse)
	if postResponse.Code != "0000" {
		t.Errorf("Expected '0000'; got %v", postResponse.Code)
	}
	responseData, ok := postResponse.Data.(map[string]interface{})
	if !ok {
		t.Errorf("Expected map[string]interface{}; got %T", postResponse.Data)
	} else if responseData["message"] != testData.Message {
		t.Errorf("Expected '%v'; got '%v'", testData.Message, responseData["message"])
	}
}
