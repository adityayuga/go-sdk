package log

import (
	"context"
	"testing"
)

func TestLoggerBasic(t *testing.T) {
	// Initialize logger
	err := Init("info")
	if err != nil {
		t.Fatalf("failed to initialize logger: %v", err)
	}

	logger := GetLogger()

	ctx := context.Background()

	// Log with basic context
	logger.Info(ctx, "application started")

	// Log with fields
	logger.Info(ctx, "user login", "username", "john", "email", "john@example.com")

	// Log a warning
	logger.Warn(ctx, "slow query detected", "query_time_ms", 1500)

	// Log an error
	logger.Error(ctx, "failed to fetch user", "user_id", 123, "error", "connection timeout")

	// Sync before exit
	logger.Sync()
}

func TestLoggerWithContext(t *testing.T) {
	err := Init("debug")
	if err != nil {
		t.Fatalf("failed to initialize logger: %v", err)
	}

	logger := GetLogger()

	// Create context with trace info
	ctx := context.Background()
	ctx = WithRequestID(ctx, "req-12345")
	ctx = WithTraceID(ctx, "trace-67890")
	ctx = WithUserID(ctx, "user-456")

	// Log with context - request_id, trace_id, and user_id will be automatically included
	logger.Info(ctx, "processing order", "order_id", 789, "amount", 99.99)

	logger.Warn(ctx, "stock low", "product_id", 456, "remaining", 5)

	logger.Error(ctx, "payment failed", "reason", "insufficient funds")

	logger.Sync()
}

func TestGlobalLoggerFunctions(t *testing.T) {
	Init("info")

	ctx := context.Background()
	ctx = WithLogContext(ctx, "req-abc", "trace-xyz", "user-123")

	// Using global functions
	Info(ctx, "global info log", "key", "value")
	Warn(ctx, "global warn log", "status", "warning")
	Error(ctx, "global error log", "code", 500)
	Debug(ctx, "global debug log", "details", "some details")

	GetLogger().Sync()
}

func TestNewLoggerInstance(t *testing.T) {
	// Create a new logger instance separate from global
	logger, err := NewLogger("warn")
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	ctx := context.Background()
	ctx = WithRequestID(ctx, "req-999")

	logger.Info(ctx, "this won't print because level is warn")
	logger.Warn(ctx, "this will print", "level", "warn")
	logger.Error(ctx, "this will print", "level", "error")

	logger.Sync()
}

func TestFormattedLogging(t *testing.T) {
	Init("info")

	ctx := context.Background()
	ctx = WithRequestID(ctx, "req-format-test")
	ctx = WithTraceID(ctx, "trace-format-test")
	ctx = WithDebugID(ctx, "debug-12345")

	// Using formatted logging functions
	Infof(ctx, "user login: %s from %s", "john_doe", "192.168.1.1")
	Warnf(ctx, "API latency high: %dms, threshold: %dms", 1500, 1000)
	Errorf(ctx, "failed to process order %d: %s", 456, "payment declined")
	Debugf(ctx, "cache miss for key: %s", "user:john_doe")

	logger := GetLogger()

	// Using formatted logging on logger instance
	logger.Infof(ctx, "processing request from %s", "internal-service")
	logger.Warnf(ctx, "database connection pool at %d%%", 85)
	logger.Errorf(ctx, "query timeout after %d seconds", 30)

	logger.Sync()
}

func TestDebugIDContext(t *testing.T) {
	Init("debug")

	ctx := context.Background()
	ctx = WithLogContext(ctx, "req-debug", "trace-debug", "user-123")
	ctx = WithDebugID(ctx, "debug-session-789")

	logger := GetLogger()

	// Debug ID will be automatically included in logs
	logger.Info(ctx, "operation started")
	logger.Debug(ctx, "debug details", "status", "initializing")

	// Verify we can retrieve debug ID
	debugID := GetDebugID(ctx)
	if debugID != "debug-session-789" {
		t.Errorf("expected debug_id 'debug-session-789', got '%s'", debugID)
	}

	logger.Sync()
}
