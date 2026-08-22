package nsq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/adityayuga/go-sdk/debug"
	"github.com/adityayuga/go-sdk/errors"
	"github.com/adityayuga/go-sdk/publisher"
	"github.com/nsqio/go-nsq"
)

type (
	nsqPublisher struct {
		producers []*nsq.Producer
		config    Config
	}

	Config struct {
		NsqdTCPAddrs      []string      `json:"nsqd_tcp_addrs"`     // NSQ daemon TCP addresses
		MaxInFlight       int           `json:"max_in_flight"`      // Max number of messages to allow in flight
		DialTimeout       time.Duration `json:"dial_timeout"`       // Timeout for dialing NSQ
		WriteTimeout      time.Duration `json:"write_timeout"`      // Timeout for writing messages
		ReadTimeout       time.Duration `json:"read_timeout"`       // Timeout for reading responses
		HeartbeatInterval time.Duration `json:"heartbeat_interval"` // Heartbeat interval
	}

	Message struct {
		Topic     string      `json:"topic"`
		Body      interface{} `json:"body"`
		Timestamp time.Time   `json:"timestamp"`
	}
)

func New(ctx context.Context, cfg Config) (publisher.Publisher, error) {
	if len(cfg.NsqdTCPAddrs) == 0 {
		return nil, errors.New("[pkg nsq] Empty NSQD TCP addresses list")
	}

	// Set default values
	if cfg.MaxInFlight == 0 {
		cfg.MaxInFlight = 200
	}

	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = 1 * time.Second
	}

	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 1 * time.Second
	}

	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 60 * time.Second
	}

	if cfg.HeartbeatInterval == 0 {
		cfg.HeartbeatInterval = 30 * time.Second
	}

	// Create NSQ config
	nsqConfig := nsq.NewConfig()
	nsqConfig.MaxInFlight = cfg.MaxInFlight
	nsqConfig.DialTimeout = cfg.DialTimeout
	nsqConfig.WriteTimeout = cfg.WriteTimeout
	nsqConfig.ReadTimeout = cfg.ReadTimeout
	nsqConfig.HeartbeatInterval = cfg.HeartbeatInterval

	// Create producers for each NSQD address
	var producers []*nsq.Producer
	for _, addr := range cfg.NsqdTCPAddrs {
		producer, err := nsq.NewProducer(addr, nsqConfig)
		if err != nil {
			// Close already created producers on error
			for _, p := range producers {
				p.Stop()
			}
			return nil, errors.New("[pkg nsq] Failed to create producer for " + addr + ": " + err.Error())
		}
		producers = append(producers, producer)
	}

	return &nsqPublisher{
		producers: producers,
		config:    cfg,
	}, nil
}

// Publish publishes a message to the specified topic with JSON encoding
func (p *nsqPublisher) Publish(ctx context.Context, topic string, key string, message interface{}) error {
	if topic == "" {
		return errors.New("[pkg nsq] Empty topic")
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return errors.New("[pkg nsq] Failed to marshal message: " + err.Error())
	}

	return p.PublishRaw(ctx, topic, key, messageBytes)
}

// PublishRaw publishes raw bytes to the specified topic
// Note: NSQ doesn't support message keys, so the key parameter is ignored
func (p *nsqPublisher) PublishRaw(ctx context.Context, topic string, key string, message []byte) error {
	if topic == "" {
		return errors.New("[pkg nsq] Empty topic")
	}

	if len(message) == 0 {
		return errors.New("[pkg nsq] Empty message")
	}

	// key parameter is ignored as NSQ doesn't support message keys

	// Extract trace headers from context and wrap the body in an Envelope
	// so downstream consumers can propagate x-debug-id and other x-* headers.
	debugID, _ := debug.GetDebugIDFromContext(ctx)
	headers := map[string]string{}
	for _, key := range []string{"x-apps-id", "x-user-id"} {
		if val, _ := ctx.Value(key).(string); val != "" {
			headers[key] = val
		}
	}
	if ip := debug.GetServiceIP(); ip != "" {
		headers["x-from-service-ip"] = ip
	}
	wrapped, err := debug.WrapEnvelope(debugID, headers, message)
	if err != nil {
		return errors.New("[pkg nsq] Failed to wrap message: " + err.Error())
	}

	// Publish to the first available producer (round-robin could be added)
	if len(p.producers) == 0 {
		return errors.New("[pkg nsq] No producers available")
	}

	// Try to publish with the first producer
	err = p.producers[0].Publish(topic, wrapped)
	if err != nil {
		return errors.New("[pkg nsq] Failed to publish message: " + err.Error())
	}

	return nil
}

// Close closes the NSQ publisher
func (p *nsqPublisher) Close(ctx context.Context) error {
	// Stop all producers gracefully
	for _, producer := range p.producers {
		producer.Stop()
	}
	return nil
}
