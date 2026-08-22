package publisher

import "context"

// Publisher defines the common interface for all message publishers
type Publisher interface {
	// Publish publishes a message to the specified topic with JSON encoding
	// For systems that support keys (like Kafka), key should be provided
	// For systems that don't support keys (like NSQ), key can be empty string
	Publish(ctx context.Context, topic string, key string, message interface{}) error

	// PublishRaw publishes raw bytes to the specified topic
	// For systems that support keys (like Kafka), key should be provided
	// For systems that don't support keys (like NSQ), key can be empty string
	PublishRaw(ctx context.Context, topic string, key string, message []byte) error

	// Close closes the publisher and cleans up resources
	Close(ctx context.Context) error
}
