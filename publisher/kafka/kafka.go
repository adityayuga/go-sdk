package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/adityayuga/go-sdk/debug"
	"github.com/adityayuga/go-sdk/errors"
	"github.com/adityayuga/go-sdk/publisher"
	"github.com/segmentio/kafka-go"
)

type (
	kafkaPublisher struct {
		writer *kafka.Writer
		config Config
	}

	Config struct {
		Brokers      []string      `json:"brokers"`
		RequiredAcks int           `json:"required_acks"` // 0=NoResponse, 1=WaitForLeader, -1=WaitForAll
		Timeout      time.Duration `json:"timeout"`       // Write timeout
		BatchSize    int           `json:"batch_size"`    // Max number of messages in a batch
		BatchTimeout time.Duration `json:"batch_timeout"` // Max time to wait for a batch
	}

	Message struct {
		Topic     string      `json:"topic"`
		Key       string      `json:"key"`
		Value     interface{} `json:"value"`
		Timestamp time.Time   `json:"timestamp"`
	}
)

func New(ctx context.Context, cfg Config) (publisher.Publisher, error) {
	if len(cfg.Brokers) == 0 {
		return nil, errors.New("[pkg kafka] Empty brokers list")
	}

	// Set default values
	if cfg.RequiredAcks == 0 {
		cfg.RequiredAcks = 1 // Default to WaitForLeader
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}

	if cfg.BatchSize == 0 {
		cfg.BatchSize = 100
	}

	if cfg.BatchTimeout == 0 {
		cfg.BatchTimeout = 1 * time.Second
	}

	// Create Kafka writer
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequiredAcks(cfg.RequiredAcks),
		WriteTimeout: cfg.Timeout,
		BatchSize:    cfg.BatchSize,
		BatchTimeout: cfg.BatchTimeout,
		Async:        false, // Synchronous writes for error handling
	}

	return &kafkaPublisher{
		writer: writer,
		config: cfg,
	}, nil
}

// Publish publishes a message to the specified topic with JSON encoding
func (p *kafkaPublisher) Publish(ctx context.Context, topic string, key string, message interface{}) error {
	if topic == "" {
		return errors.New("[pkg kafka] Empty topic")
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return errors.New("[pkg kafka] Failed to marshal message: " + err.Error())
	}

	return p.PublishRaw(ctx, topic, key, messageBytes)
}

// PublishRaw publishes raw bytes to the specified topic
func (p *kafkaPublisher) PublishRaw(ctx context.Context, topic string, key string, message []byte) error {
	if topic == "" {
		return errors.New("[pkg kafka] Empty topic")
	}

	debugID, _ := debug.GetDebugIDFromContext(ctx)

	headers := []kafka.Header{
		{Key: HeaderDebugID, Value: []byte(debugID)},
	}

	// x-apps-id
	if appsID, _ := ctx.Value("x-apps-id").(string); appsID != "" {
		headers = append(headers, kafka.Header{Key: HeaderAppsID, Value: []byte(appsID)})
	}

	// x-user-id (only when a user is authenticated)
	if userID, _ := ctx.Value("x-user-id").(string); userID != "" {
		headers = append(headers, kafka.Header{Key: HeaderUserID, Value: []byte(userID)})
	}

	// x-from-service-ip: stamp the current service's outbound IP
	if ip := debug.GetServiceIP(); ip != "" {
		headers = append(headers, kafka.Header{Key: HeaderFromServiceIP, Value: []byte(ip)})
	}

	kafkaMessage := kafka.Message{
		Topic:   topic,
		Key:     []byte(key),
		Value:   message,
		Time:    time.Now(),
		Headers: headers,
	}

	err := p.writer.WriteMessages(ctx, kafkaMessage)
	if err != nil {
		return errors.New("[pkg kafka] Failed to write message: " + err.Error())
	}

	return nil
}

// Close closes the Kafka publisher
func (p *kafkaPublisher) Close(ctx context.Context) error {
	if p.writer != nil {
		err := p.writer.Close()
		if err != nil {
			return errors.New("[pkg kafka] Failed to close writer: " + err.Error())
		}
	}
	return nil
}
