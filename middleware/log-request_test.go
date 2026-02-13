package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adityayuga/go-sdk/log"
)

func init() {
	log.Init("info")
}

func TestLoggerBasic(t *testing.T) {
	// Disable body logging for this test
	DisableRequestBodyLogging()
	DisableResponseBodyLogging()

	// Create a simple handler that returns a response
	handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
	}))

	// Create a test request
	req := httptest.NewRequest("GET", "/api/users", nil)
	w := httptest.NewRecorder()

	// Execute the request
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestLoggerWithRequestBody(t *testing.T) {
	// Enable request body logging
	EnableRequestBodyLogging()
	DisableResponseBodyLogging()

	// Create a handler that logs request body
	handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	}))

	// Create a test request with body
	body := []byte(`{"name":"John","age":30}`)
	req := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Execute the request
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Cleanup
	DisableRequestBodyLogging()
}

func TestLoggerWithResponseBody(t *testing.T) {
	// Enable response body logging
	DisableRequestBodyLogging()
	EnableResponseBodyLogging()

	// Create a handler that logs response body
	handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":123,"name":"John"}`))
	}))

	// Create a test request
	req := httptest.NewRequest("GET", "/api/users/123", nil)
	w := httptest.NewRecorder()

	// Execute the request
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify response is captured
	body := w.Body.String()
	if body != `{"id":123,"name":"John"}` {
		t.Errorf("expected response body to match, got %s", body)
	}

	// Cleanup
	DisableResponseBodyLogging()
}

func TestLoggerWithBothBodies(t *testing.T) {
	// Enable both request and response body logging
	EnableRequestBodyLogging()
	EnableResponseBodyLogging()

	// Create a handler that logs both request and response bodies
	handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read request body
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		// Write response based on request
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":456,"data":` + string(body) + `}`))
	}))

	// Create a test request with body
	requestBody := []byte(`{"name":"Jane","email":"jane@example.com"}`)
	req := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()

	// Execute the request
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	// Cleanup
	DisableRequestBodyLogging()
	DisableResponseBodyLogging()
}

func TestLoggerErrorResponse(t *testing.T) {
	// Enable both request and response body logging
	EnableRequestBodyLogging()
	EnableResponseBodyLogging()

	// Create a handler that returns an error
	handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal Server Error"}`))
	}))

	// Create a test request
	req := httptest.NewRequest("GET", "/api/broken", nil)
	w := httptest.NewRecorder()

	// Execute the request
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	// Cleanup
	DisableRequestBodyLogging()
	DisableResponseBodyLogging()
}

func TestLoggerMultipleWrites(t *testing.T) {
	// Enable response body logging
	DisableRequestBodyLogging()
	EnableResponseBodyLogging()

	// Create a handler that writes response in multiple chunks
	handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello "))
		w.Write([]byte("World"))
	}))

	// Create a test request
	req := httptest.NewRequest("GET", "/api/stream", nil)
	w := httptest.NewRecorder()

	// Execute the request
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify full response is captured
	body := w.Body.String()
	if body != "Hello World" {
		t.Errorf("expected 'Hello World', got %s", body)
	}

	// Cleanup
	DisableResponseBodyLogging()
}
