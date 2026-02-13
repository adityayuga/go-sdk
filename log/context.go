package log

import (
	"context"

	"github.com/adityayuga/go-sdk/debug"
)

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, "request_id", requestID)
}

// WithTraceID adds trace ID to context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, "trace_id", traceID)
}

// WithUserID adds user ID to context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, "user_id", userID)
}

// WithLogContext sets multiple context values at once
func WithLogContext(ctx context.Context, requestID, traceID, userID string) context.Context {
	ctx = WithRequestID(ctx, requestID)
	ctx = WithTraceID(ctx, traceID)
	ctx = WithUserID(ctx, userID)
	return ctx
}

// GetRequestID retrieves request ID from context
func GetRequestID(ctx context.Context) string {
	if reqID := ctx.Value("request_id"); reqID != nil {
		if id, ok := reqID.(string); ok {
			return id
		}
	}
	return ""
}

// GetTraceID retrieves trace ID from context
func GetTraceID(ctx context.Context) string {
	if traceID := ctx.Value("trace_id"); traceID != nil {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	return ""
}

// GetUserID retrieves user ID from context
func GetUserID(ctx context.Context) string {
	if userID := ctx.Value("user_id"); userID != nil {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// WithDebugID adds debug ID to context (uses debug package)
func WithDebugID(ctx context.Context, debugID string) context.Context {
	return debug.SetDebugIDOnContext(ctx, debugID)
}

// GetDebugID retrieves debug ID from context (uses debug package)
func GetDebugID(ctx context.Context) string {
	debugID, _ := debug.GetDebugIDFromContext(ctx)
	return debugID
}
