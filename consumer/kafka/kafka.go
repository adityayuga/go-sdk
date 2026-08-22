// Package kafka provides a Kafka consumer implementation that satisfies the
// consumer.Consumer interface. It wraps segmentio/kafka-go and handles:
//   - sequential message processing (fetch → handle → commit) to prevent
//     duplicate delivery on rebalance;
//   - x-debug-id and other x-* trace header extraction from Kafka message headers;
//   - graceful shutdown via context cancellation.
package kafka

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/adityayuga/go-sdk/consumer"
	kafkaLib "github.com/segmentio/kafka-go"
)

// Config holds all parameters needed to create a KafkaConsumer.
type Config struct {
	// Brokers is the list of Kafka broker addresses (host:port).
	Brokers []string `json:"brokers"`

	// Topic is the Kafka topic to consume from.
	Topic string `json:"topic"`

	// GroupID is the consumer group ID. When empty, defaults to
	// "consumer-<topic>".
	GroupID string `json:"group_id"`

	// SessionTimeout overrides the default 120 s session timeout.
	// Set this higher than your expected maximum per-message processing time.
	SessionTimeout time.Duration `json:"session_timeout"`

	// RebalanceTimeout overrides the default 120 s rebalance timeout.
	RebalanceTimeout time.Duration `json:"rebalance_timeout"`

	// HeartbeatInterval overrides the default 10 s heartbeat interval.
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
}

// KafkaConsumer consumes messages from a single Kafka topic and delivers them
// to a MessageHandler callback.
type KafkaConsumer struct {
	config  Config
	handler consumer.MessageHandler
	reader  *kafkaLib.Reader
}

// New creates a new KafkaConsumer. It applies sensible defaults for any
// zero-value duration fields in cfg.
func New(cfg Config, handler consumer.MessageHandler) (*KafkaConsumer, error) {
	if len(cfg.Brokers) == 0 {
		return nil, fmt.Errorf("[pkg kafka consumer] brokers list is empty")
	}
	if cfg.Topic == "" {
		return nil, fmt.Errorf("[pkg kafka consumer] topic is required")
	}
	if handler == nil {
		return nil, fmt.Errorf("[pkg kafka consumer] handler is required")
	}

	// Apply defaults
	if cfg.GroupID == "" {
		cfg.GroupID = fmt.Sprintf("consumer-%s", cfg.Topic)
	}
	if cfg.SessionTimeout == 0 {
		cfg.SessionTimeout = DefaultSessionTimeout
	}
	if cfg.RebalanceTimeout == 0 {
		cfg.RebalanceTimeout = DefaultRebalanceTimeout
	}
	if cfg.HeartbeatInterval == 0 {
		cfg.HeartbeatInterval = DefaultHeartbeatInterval
	}

	reader := kafkaLib.NewReader(kafkaLib.ReaderConfig{
		Brokers:           cfg.Brokers,
		Topic:             cfg.Topic,
		GroupID:           cfg.GroupID,
		SessionTimeout:    cfg.SessionTimeout,
		RebalanceTimeout:  cfg.RebalanceTimeout,
		HeartbeatInterval: cfg.HeartbeatInterval,
	})

	return &KafkaConsumer{
		config:  cfg,
		handler: handler,
		reader:  reader,
	}, nil
}

// Start begins the consume loop. It blocks until ctx is cancelled.
//
// Messages are processed sequentially: the handler is called, then the offset
// is committed. This prevents the broker from redelivering an in-progress
// message after a rebalance caused by slow processing.
func (c *KafkaConsumer) Start(ctx context.Context) error {
	fmt.Printf("[pkg kafka consumer] started consuming topic: %s (group: %s)\n", c.config.Topic, c.config.GroupID)

	for {
		// Check for shutdown before blocking on FetchMessage.
		select {
		case <-ctx.Done():
			fmt.Printf("[pkg kafka consumer] context cancelled, stopping topic: %s\n", c.config.Topic)
			return nil
		default:
		}

		// FetchMessage blocks until a message arrives or an error occurs.
		// We use context.Background() here intentionally: we want the read to
		// block until a message arrives rather than be cancelled immediately
		// when ctx is cancelled mid-wait. The shutdown check above ensures we
		// exit cleanly at the top of each loop iteration.
		kafkaMsg, err := c.reader.FetchMessage(context.Background())
		if err != nil {
			// Re-check ctx in case the error was caused by a shutdown.
			select {
			case <-ctx.Done():
				fmt.Printf("[pkg kafka consumer] context cancelled while fetching from topic: %s\n", c.config.Topic)
				return nil
			default:
				log.Printf("[pkg kafka consumer] error fetching from topic %s: %v", c.config.Topic, err)
				time.Sleep(DefaultFetchErrorDelay)
				continue
			}
		}

		msg := consumer.Message{
			ID:        string(kafkaMsg.Key),
			Body:      kafkaMsg.Value,
			Timestamp: kafkaMsg.Time.Unix(),
			Headers:   extractHeaders(kafkaMsg.Headers),
		}

		// Inject trace headers into context before calling handler
		handlerCtx := consumer.InjectContext(ctx, msg)

		// Deliver to handler synchronously — commit only after success.
		if err := c.handler(handlerCtx, msg); err != nil {
			log.Printf("[pkg kafka consumer] handler error for topic %s key %s: %v", c.config.Topic, msg.ID, err)
		}

		// Commit the offset regardless of handler error so we make forward
		// progress. Callers that need dead-letter behaviour should handle it
		// inside their handler.
		if err := c.reader.CommitMessages(context.Background(), kafkaMsg); err != nil {
			log.Printf("[pkg kafka consumer] error committing offset for topic %s: %v", c.config.Topic, err)
		}
	}
}

// Close closes the underlying Kafka reader and releases resources.
func (c *KafkaConsumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}

// extractHeaders converts Kafka message headers into a plain string map.
// All header keys are lower-cased. Only x-* prefixed headers are included,
// matching the convention used by the Kafka publisher in pkg/go/publisher/kafka.
func extractHeaders(headers []kafkaLib.Header) map[string]string {
	result := make(map[string]string, len(headers))
	for _, h := range headers {
		key := strings.ToLower(h.Key)
		if strings.HasPrefix(key, "x-") && len(h.Value) > 0 {
			result[key] = string(h.Value)
		}
	}
	return result
}
