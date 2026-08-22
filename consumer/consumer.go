package consumer

import (
	"context"
	"strings"

	"github.com/adityayuga/go-sdk/debug"
)

// Message represents a message received from a message broker.
// It is intentionally broker-agnostic so that job handlers can be written
// once and used with any supported engine (Kafka, NSQ, etc.).
type Message struct {
	// ID is the message key (Kafka) or message ID (NSQ).
	ID string

	// Body is the raw message payload. For NSQ, this is already unwrapped from
	// the TraceEnvelope so the handler always receives the original body.
	Body []byte

	// Timestamp is the message timestamp in Unix seconds.
	Timestamp int64

	// Headers contains any trace headers extracted from the message
	// (e.g. x-debug-id, x-apps-id, x-user-id, x-from-service-ip).
	// x-debug-id is available at Headers["x-debug-id"].
	Headers map[string]string
}

// MessageHandler is the callback signature that consumers invoke for every
// received message. Returning a non-nil error signals that the message could
// not be processed; the caller (e.g. job-worker/core) is responsible for any
// retry / dead-letter logic.
type MessageHandler func(ctx context.Context, msg Message) error

// Consumer is the common interface implemented by all broker-specific consumers.
type Consumer interface {
	// Start begins consuming messages and blocks until ctx is cancelled or a
	// fatal error occurs. Each received message is delivered to the handler
	// provided at construction time.
	Start(ctx context.Context) error

	// Close releases resources held by the consumer (connections, goroutines,
	// etc.). It is safe to call Close more than once.
	Close() error
}

// InjectContext injects trace headers from the message into the context.
// This is called automatically by the consumer implementations before invoking
// the handler, ensuring consistent trace header propagation across all consumers.
//
// The following headers are injected:
//   - x-debug-id: Propagated or generated using debug package
//   - All other x-* headers: Injected via context.WithValue
//
// This function returns the enriched context that should be passed to the handler.
func InjectContext(ctx context.Context, msg Message) context.Context {
	// Extract debug ID from headers
	debugID := ""
	if msg.Headers != nil {
		debugID = msg.Headers["x-debug-id"]
	}

	// Set or generate debug ID
	if debugID != "" {
		ctx = debug.SetDebugIDOnContext(ctx, debugID)
	} else {
		_, ctx = debug.GetDebugIDFromContext(ctx)
	}

	// Inject remaining x-* headers into context
	for k, v := range msg.Headers {
		if k != "x-debug-id" && strings.HasPrefix(k, "x-") && v != "" {
			ctx = context.WithValue(ctx, k, v)
		}
	}

	return ctx
}
