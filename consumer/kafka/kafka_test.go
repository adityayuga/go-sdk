package kafka

import (
	"context"
	"testing"

	"github.com/adityayuga/go-sdk/consumer"
	"github.com/segmentio/kafka-go"
)

func TestExtractHeaders(t *testing.T) {
	tests := []struct {
		name     string
		headers  []kafka.Header
		expected map[string]string
	}{
		{
			name: "with x-* headers",
			headers: []kafka.Header{
				{Key: "x-debug-id", Value: []byte("12345")},
				{Key: "x-apps-id", Value: []byte("app-123")},
				{Key: "x-user-id", Value: []byte("user-456")},
			},
			expected: map[string]string{
				"x-debug-id": "12345",
				"x-apps-id":  "app-123",
				"x-user-id":  "user-456",
			},
		},
		{
			name: "with non-x-* headers (should be ignored)",
			headers: []kafka.Header{
				{Key: "content-type", Value: []byte("application/json")},
				{Key: "x-debug-id", Value: []byte("12345")},
			},
			expected: map[string]string{
				"x-debug-id": "12345",
			},
		},
		{
			name:     "empty headers",
			headers:  []kafka.Header{},
			expected: map[string]string{},
		},
		{
			name:     "nil headers",
			headers:  nil,
			expected: map[string]string{},
		},
		{
			name: "case-insensitive keys",
			headers: []kafka.Header{
				{Key: "X-Debug-Id", Value: []byte("12345")},
				{Key: "X-APPS-ID", Value: []byte("app-123")},
			},
			expected: map[string]string{
				"x-debug-id": "12345",
				"x-apps-id":  "app-123",
			},
		},
		{
			name: "empty value (should be ignored)",
			headers: []kafka.Header{
				{Key: "x-debug-id", Value: []byte("12345")},
				{Key: "x-empty", Value: []byte{}},
			},
			expected: map[string]string{
				"x-debug-id": "12345",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractHeaders(tt.headers)

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
			name: "empty brokers",
			cfg: Config{
				Brokers: []string{},
				Topic:   "test-topic",
			},
			wantError: true,
		},
		{
			name: "empty topic",
			cfg: Config{
				Brokers: []string{"localhost:9092"},
				Topic:   "",
			},
			wantError: true,
		},
		{
			name: "valid config",
			cfg: Config{
				Brokers: []string{"localhost:9092"},
				Topic:   "test-topic",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestNew_DefaultGroupID(t *testing.T) {
	cfg := Config{
		Brokers: []string{"localhost:9092"},
		Topic:   "my-topic",
		// GroupID not set
	}

	c, err := New(cfg, func(ctx context.Context, msg consumer.Message) error {
		return nil
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if c.config.GroupID != "consumer-my-topic" {
		t.Errorf("Expected GroupID 'consumer-my-topic', got '%s'", c.config.GroupID)
	}
}

func TestNew_CustomGroupID(t *testing.T) {
	cfg := Config{
		Brokers: []string{"localhost:9092"},
		Topic:   "my-topic",
		GroupID: "my-custom-group",
	}

	c, err := New(cfg, func(ctx context.Context, msg consumer.Message) error {
		return nil
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if c.config.GroupID != "my-custom-group" {
		t.Errorf("Expected GroupID 'my-custom-group', got '%s'", c.config.GroupID)
	}
}
