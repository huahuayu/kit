package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type sampleResponse struct {
	Message string `json:"message"`
	Data    struct {
		ID    int    `json:"id"`
		Value string `json:"value"`
	} `json:"data"`
}

func TestGet(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		if req.URL.String() != "/test" {
			t.Errorf("got: %s, want: /test", req.URL.String())
		}
		// Send response to be tested
		rw.Write([]byte(`{"message": "OK", "data": {"id": 1, "value": "test"}}`))
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use Client & URL from our local test server
	api := New()
	var result sampleResponse
	err := api.Get(server.URL+"/test", &result)
	if err != nil {
		t.Errorf("got error: %s", err.Error())
	}
	if result.Message != "OK" {
		t.Errorf("got: %s, want: OK", result.Message)
	}
	if result.Data.ID != 1 {
		t.Errorf("got: %d, want: 1", result.Data.ID)
	}
	if result.Data.Value != "test" {
		t.Errorf("got: %s, want: test", result.Data.Value)
	}
}

func TestPost(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		if req.URL.String() != "/test" {
			t.Errorf("got: %s, want: /test", req.URL.String())
		}
		// Send response to be tested
		rw.Write([]byte(`{"message": "OK", "data": {"id": 1, "value": "test"}}`))
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use Client & URL from our local test server
	api := New()
	var result sampleResponse
	err := api.Post(server.URL+"/test", map[string]string{"test": "data"}, &result)
	if err != nil {
		t.Errorf("got error: %s", err.Error())
	}
	if result.Message != "OK" {
		t.Errorf("got: %s, want: OK", result.Message)
	}
	if result.Data.ID != 1 {
		t.Errorf("got: %d, want: 1", result.Data.ID)
	}
	if result.Data.Value != "test" {
		t.Errorf("got: %s, want: test", result.Data.Value)
	}
}

func TestGetWithParamAndHeaders(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		if req.URL.String() != "/test?param=value" {
			t.Errorf("got: %s, want: /test?param=value", req.URL.String())
		}
		// Test request headers
		if req.Header.Get("Custom-Header") != "Custom-Value" {
			t.Errorf("got: %s, want: Custom-Value", req.Header.Get("Custom-Header"))
		}
		// Send response to be tested
		rw.Write([]byte(`{"message": "OK", "data": {"id": 1, "value": "test"}}`))
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use Client & URL from our local test server
	api := New()
	var result sampleResponse
	err := api.Get(server.URL+"/test", &result, Param{
		Query:  map[string]string{"param": "value"},
		Header: map[string]string{"Custom-Header": "Custom-Value"},
	})
	if err != nil {
		t.Errorf("got error: %s", err.Error())
	}
}

func TestPostWithHeaders(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		if req.URL.String() != "/test" {
			t.Errorf("got: %s, want: /test", req.URL.String())
		}
		// Test request headers
		if req.Header.Get("Custom-Header") != "Custom-Value" {
			t.Errorf("got: %s, want: Custom-Value", req.Header.Get("Custom-Header"))
		}
		// Send response to be tested
		rw.Write([]byte(`{"message": "OK", "data": {"id": 1, "value": "test"}}`))
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use Client & URL from our local test server
	api := New()
	var result sampleResponse
	err := api.Post(server.URL+"/test", map[string]string{"test": "data"}, &result, map[string]string{"Custom-Header": "Custom-Value"})
	if err != nil {
		t.Errorf("got error: %s", err.Error())
	}
}

func TestGetWithErrorParsingJSON(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send invalid JSON response
		rw.Write([]byte(`{"message": "OK", "data": {"id": 1, "value": "test"`)) // missing closing brace
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use Client & URL from our local test server
	api := New()
	var result sampleResponse
	err := api.Get(server.URL+"/test", &result)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestPostWithErrorParsingJSON(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send invalid JSON response
		rw.Write([]byte(`{"message": "OK", "data": {"id": 1, "value": "test"`)) // missing closing brace
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use Client & URL from our local test server
	api := New()
	var result sampleResponse
	err := api.Post(server.URL+"/test", map[string]string{"test": "data"}, &result)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
