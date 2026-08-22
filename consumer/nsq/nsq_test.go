package nsq

import (
	"context"
	"testing"

	"github.com/adityayuga/go-sdk/consumer"
)

func TestBuildHeadersFromEnvelope(t *testing.T) {
	tests := []struct {
		name     string
		debugID  string
		headers  map[string]string
		expected map[string]string
	}{
		{
			name:    "with debug ID and headers",
			debugID: "12345",
			headers: map[string]string{
				"x-apps-id":    "app-123",
				"x-user-id":    "user-456",
				"other-header": "ignored",
			},
			expected: map[string]string{
				"x-debug-id": "12345",
				"x-apps-id":  "app-123",
				"x-user-id":  "user-456",
			},
		},
		{
			name:    "empty debug ID",
			debugID: "",
			headers: map[string]string{
				"x-apps-id": "app-123",
			},
			expected: map[string]string{
				"x-apps-id": "app-123",
			},
		},
		{
			name:     "empty headers",
			debugID:  "12345",
			headers:  map[string]string{},
			expected: map[string]string{"x-debug-id": "12345"},
		},
		{
			name:     "nil headers",
			debugID:  "12345",
			headers:  nil,
			expected: map[string]string{"x-debug-id": "12345"},
		},
		{
			name:    "only x-* headers included",
			debugID: "12345",
			headers: map[string]string{
				"x-custom":      "included",
				"content-type":  "ignored",
				"authorization": "ignored",
			},
			expected: map[string]string{
				"x-debug-id": "12345",
				"x-custom":   "included",
			},
		},
		{
			name:    "empty header values ignored",
			debugID: "12345",
			headers: map[string]string{
				"x-empty":    "",
				"x-nonempty": "value",
			},
			expected: map[string]string{
				"x-debug-id": "12345",
				"x-nonempty": "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildHeadersFromEnvelope(tt.debugID, tt.headers)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d headers, got %d", len(tt.expected), len(result))
			}

			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("Expected header '%s' to be '%s', got '%s'", k, v, result[k])
				}
			}
		})
	}
}

func TestNew_Validation(t *testing.T) {
	tests := []struct {
		name      string
		cfg       Config
		wantError bool
	}{
		{
			name: "empty topic",
			cfg: Config{
				Topic:            "",
				Channel:          "test-channel",
				LookupdHTTPAddrs: []string{"localhost:4161"},
			},
			wantError: true,
		},
		{
			name: "empty channel",
			cfg: Config{
				Topic:            "test-topic",
				Channel:          "",
				LookupdHTTPAddrs: []string{"localhost:4161"},
			},
			wantError: true,
		},
		{
			name: "no connection addresses",
			cfg: Config{
				Topic:   "test-topic",
				Channel: "test-channel",
			},
			wantError: true,
		},
		{
			name: "valid with lookupd",
			cfg: Config{
				Topic:            "test-topic",
				Channel:          "test-channel",
				LookupdHTTPAddrs: []string{"localhost:4161"},
			},
			wantError: false,
		},
		{
			name: "valid with nsqd",
			cfg: Config{
				Topic:        "test-topic",
				Channel:      "test-channel",
				NSQDTCPAddrs: []string{"localhost:4150"},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip tests that would actually connect to NSQ
			if !tt.wantError {
				t.Skip("Skipping integration test that requires NSQ connection")
			}

			_, err := New(tt.cfg, func(ctx context.Context, msg consumer.Message) error {
				return nil
			})

			if tt.wantError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}
