package consumer

import (
	"context"
	"testing"

	"github.com/adityayuga/go-sdk/debug"
)

func TestInjectContext_WithDebugID(t *testing.T) {
	ctx := context.Background()
	msg := Message{
		ID:        "test-id",
		Body:      []byte(`{"test": "data"}`),
		Timestamp: 1234567890,
		Headers: map[string]string{
			"x-debug-id": "12345",
		},
	}

	resultCtx := InjectContext(ctx, msg)

	debugID, _ := debug.GetDebugIDFromContext(resultCtx)
	if debugID != "12345" {
		t.Errorf("Expected debug ID '12345', got '%s'", debugID)
	}
}

func TestInjectContext_WithoutDebugID(t *testing.T) {
	ctx := context.Background()
	msg := Message{
		ID:        "test-id",
		Body:      []byte(`{"test": "data"}`),
		Timestamp: 1234567890,
		Headers:   map[string]string{},
	}

	resultCtx := InjectContext(ctx, msg)

	debugID, _ := debug.GetDebugIDFromContext(resultCtx)
	if debugID == "" {
		t.Error("Expected non-empty debug ID to be generated")
	}
}

func TestInjectContext_WithNilHeaders(t *testing.T) {
	ctx := context.Background()
	msg := Message{
		ID:        "test-id",
		Body:      []byte(`{"test": "data"}`),
		Timestamp: 1234567890,
		Headers:   nil,
	}

	resultCtx := InjectContext(ctx, msg)

	// Should not panic and should generate a debug ID
	debugID, _ := debug.GetDebugIDFromContext(resultCtx)
	if debugID == "" {
		t.Error("Expected non-empty debug ID to be generated")
	}
}

func TestInjectContext_WithXHeaders(t *testing.T) {
	ctx := context.Background()
	msg := Message{
		ID:        "test-id",
		Body:      []byte(`{"test": "data"}`),
		Timestamp: 1234567890,
		Headers: map[string]string{
			"x-debug-id":        "12345",
			"x-apps-id":         "app-123",
			"x-user-id":         "user-456",
			"x-from-service-ip": "192.168.1.1",
			"other-header":      "should-be-ignored",
		},
	}

	resultCtx := InjectContext(ctx, msg)

	// Check debug ID
	debugID, _ := debug.GetDebugIDFromContext(resultCtx)
	if debugID != "12345" {
		t.Errorf("Expected debug ID '12345', got '%s'", debugID)
	}

	// Check x-apps-id
	appsID := resultCtx.Value("x-apps-id")
	if appsID != "app-123" {
		t.Errorf("Expected x-apps-id 'app-123', got '%v'", appsID)
	}

	// Check x-user-id
	userID := resultCtx.Value("x-user-id")
	if userID != "user-456" {
		t.Errorf("Expected x-user-id 'user-456', got '%v'", userID)
	}

	// Check x-from-service-ip
	serviceIP := resultCtx.Value("x-from-service-ip")
	if serviceIP != "192.168.1.1" {
		t.Errorf("Expected x-from-service-ip '192.168.1.1', got '%v'", serviceIP)
	}

	// Check that non-x-* header is NOT injected
	otherHeader := resultCtx.Value("other-header")
	if otherHeader != nil {
		t.Errorf("Expected 'other-header' to NOT be injected, got '%v'", otherHeader)
	}
}

func TestInjectContext_OnlyXHeadersInjected(t *testing.T) {
	ctx := context.Background()
	msg := Message{
		ID:        "test-id",
		Body:      []byte(`{"test": "data"}`),
		Timestamp: 1234567890,
		Headers: map[string]string{
			"content-type":  "application/json",
			"authorization": "Bearer token",
			"x-custom":      "should-be-injected",
		},
	}

	resultCtx := InjectContext(ctx, msg)

	// Check that non-x-* headers are NOT injected
	contentType := resultCtx.Value("content-type")
	if contentType != nil {
		t.Errorf("Expected 'content-type' to NOT be injected, got '%v'", contentType)
	}

	auth := resultCtx.Value("authorization")
	if auth != nil {
		t.Errorf("Expected 'authorization' to NOT be injected, got '%v'", auth)
	}

	// Check that x-* header IS injected
	custom := resultCtx.Value("x-custom")
	if custom != "should-be-injected" {
		t.Errorf("Expected 'x-custom' to be 'should-be-injected', got '%v'", custom)
	}
}

func TestInjectContext_EmptyHeaderValueIgnored(t *testing.T) {
	ctx := context.Background()
	msg := Message{
		ID:        "test-id",
		Body:      []byte(`{"test": "data"}`),
		Timestamp: 1234567890,
		Headers: map[string]string{
			"x-debug-id": "12345",
			"x-empty":    "",
			"x-nonempty": "value",
		},
	}

	resultCtx := InjectContext(ctx, msg)

	// Empty header should NOT be injected
	emptyHeader := resultCtx.Value("x-empty")
	if emptyHeader != nil {
		t.Errorf("Expected 'x-empty' to NOT be injected when empty, got '%v'", emptyHeader)
	}

	// Non-empty header should be injected
	nonEmptyHeader := resultCtx.Value("x-nonempty")
	if nonEmptyHeader != "value" {
		t.Errorf("Expected 'x-nonempty' to be 'value', got '%v'", nonEmptyHeader)
	}
}
