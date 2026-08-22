// Package nsq provides an NSQ consumer implementation that satisfies the
// consumer.Consumer interface. It wraps go-nsq and handles:
//   - TraceEnvelope unwrapping so handlers always receive the original payload
//     and trace headers even though NSQ has no native header support;
//   - x-* header extraction and delivery via consumer.Message.Headers;
//   - graceful shutdown via context cancellation and nsq.Consumer.Stop().
package nsq

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/adityayuga/go-sdk/consumer"
	"github.com/adityayuga/go-sdk/debug"
	nsqLib "github.com/nsqio/go-nsq"
)

// Config holds all parameters needed to create an NSQConsumer.
type Config struct {
	// Topic is the NSQ topic to subscribe to.
	Topic string `json:"topic"`

	// Channel is the NSQ channel name. Multiple consumers on the same channel
	// receive messages in a round-robin fashion.
	Channel string `json:"channel"`

	// LookupdHTTPAddrs is the list of nsqlookupd HTTP addresses
	// (e.g. ["localhost:4161"]). When provided, the consumer uses service
	// discovery. Takes priority over NSQDTCPAddrs.
	LookupdHTTPAddrs []string `json:"lookupd_http_addrs"`

	// NSQDTCPAddrs is the list of nsqd TCP addresses to connect to directly
	// (e.g. ["localhost:4150"]). Used when LookupdHTTPAddrs is empty.
	NSQDTCPAddrs []string `json:"nsqd_tcp_addrs"`

	// MaxInFlight overrides the default maximum number of in-flight messages.
	MaxInFlight int `json:"max_in_flight"`

	// DialTimeout overrides the default dial timeout.
	DialTimeout time.Duration `json:"dial_timeout"`

	// WriteTimeout overrides the default write timeout.
	WriteTimeout time.Duration `json:"write_timeout"`

	// ReadTimeout overrides the default read timeout.
	ReadTimeout time.Duration `json:"read_timeout"`

	// HeartbeatInterval overrides the default heartbeat interval.
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
}

// NSQConsumer consumes messages from a single NSQ topic/channel and delivers
// them to a MessageHandler callback.
type NSQConsumer struct {
	config   Config
	handler  consumer.MessageHandler
	consumer *nsqLib.Consumer
}

// New creates a new NSQConsumer and connects it to NSQ. It applies sensible
// defaults for any zero-value fields in cfg.
func New(cfg Config, handler consumer.MessageHandler) (*NSQConsumer, error) {
	if cfg.Topic == "" {
		return nil, fmt.Errorf("[pkg nsq consumer] topic is required")
	}
	if cfg.Channel == "" {
		return nil, fmt.Errorf("[pkg nsq consumer] channel is required")
	}
	if handler == nil {
		return nil, fmt.Errorf("[pkg nsq consumer] handler is required")
	}
	if len(cfg.LookupdHTTPAddrs) == 0 && len(cfg.NSQDTCPAddrs) == 0 {
		return nil, fmt.Errorf("[pkg nsq consumer] at least one of LookupdHTTPAddrs or NSQDTCPAddrs must be set")
	}

	// Apply defaults
	if cfg.MaxInFlight == 0 {
		cfg.MaxInFlight = DefaultMaxInFlight
	}
	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = time.Duration(DefaultDialTimeoutSeconds) * time.Second
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = time.Duration(DefaultWriteTimeoutSeconds) * time.Second
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = time.Duration(DefaultReadTimeoutSeconds) * time.Second
	}
	if cfg.HeartbeatInterval == 0 {
		cfg.HeartbeatInterval = time.Duration(DefaultHeartbeatIntervalSeconds) * time.Second
	}

	nsqCfg := nsqLib.NewConfig()
	nsqCfg.MaxInFlight = cfg.MaxInFlight
	nsqCfg.DialTimeout = cfg.DialTimeout
	nsqCfg.WriteTimeout = cfg.WriteTimeout
	nsqCfg.ReadTimeout = cfg.ReadTimeout
	nsqCfg.HeartbeatInterval = cfg.HeartbeatInterval

	c, err := nsqLib.NewConsumer(cfg.Topic, cfg.Channel, nsqCfg)
	if err != nil {
		return nil, fmt.Errorf("[pkg nsq consumer] failed to create NSQ consumer: %w", err)
	}

	nc := &NSQConsumer{
		config:   cfg,
		handler:  handler,
		consumer: c,
	}

	// Register the message handler.
	c.AddHandler(nsqLib.HandlerFunc(nc.handleMessage))

	// Connect
	if len(cfg.LookupdHTTPAddrs) > 0 {
		if err := c.ConnectToNSQLookupds(cfg.LookupdHTTPAddrs); err != nil {
			return nil, fmt.Errorf("[pkg nsq consumer] failed to connect to nsqlookupd %v: %w", cfg.LookupdHTTPAddrs, err)
		}
		fmt.Printf("[pkg nsq consumer] connected to nsqlookupd: %v\n", cfg.LookupdHTTPAddrs)
	} else {
		if err := c.ConnectToNSQDs(cfg.NSQDTCPAddrs); err != nil {
			return nil, fmt.Errorf("[pkg nsq consumer] failed to connect to nsqd %v: %w", cfg.NSQDTCPAddrs, err)
		}
		fmt.Printf("[pkg nsq consumer] connected to nsqd: %v\n", cfg.NSQDTCPAddrs)
	}

	fmt.Printf("[pkg nsq consumer] started consuming topic: %s channel: %s\n", cfg.Topic, cfg.Channel)

	return nc, nil
}

// Start blocks until ctx is cancelled, at which point it stops the NSQ
// consumer and returns.
func (c *NSQConsumer) Start(ctx context.Context) error {
	<-ctx.Done()
	fmt.Printf("[pkg nsq consumer] context cancelled, stopping topic: %s channel: %s\n", c.config.Topic, c.config.Channel)
	c.consumer.Stop()
	<-c.consumer.StopChan
	return nil
}

// Close stops the NSQ consumer and releases resources. Safe to call multiple
// times.
func (c *NSQConsumer) Close() error {
	if c.consumer != nil {
		c.consumer.Stop()
	}
	return nil
}

// handleMessage is the nsq.HandlerFunc callback. It unwraps a TraceEnvelope
// (if present) and delivers the message to the user-provided MessageHandler.
func (c *NSQConsumer) handleMessage(message *nsqLib.Message) error {
	rawBody := message.Body

	msg := consumer.Message{
		ID:        string(message.ID[:]),
		Timestamp: message.Timestamp / 1e9, // nanoseconds → seconds
	}

	// Attempt to unwrap a TraceEnvelope written by the NSQ publisher.
	// Legacy messages that were not wrapped are passed through as-is.
	if env, ok := debug.UnwrapEnvelope(rawBody); ok {
		msg.Body = []byte(env.OriginalBody)
		msg.Headers = buildHeadersFromEnvelope(env.DebugID, env.Headers)
	} else {
		msg.Body = rawBody
		msg.Headers = map[string]string{}
	}

	// Inject trace headers into context before calling handler
	// Use context.Background() as base since NSQ handles shutdown via Stop()
	ctx := consumer.InjectContext(context.Background(), msg)

	// Deliver to the caller's handler
	return c.handler(ctx, msg)
}

// buildHeadersFromEnvelope assembles the headers map from the envelope's
// DebugID and the generic headers map, applying the x-* prefix convention.
func buildHeadersFromEnvelope(debugID string, headers map[string]string) map[string]string {
	result := make(map[string]string, len(headers)+1)
	if debugID != "" {
		result["x-debug-id"] = debugID
	}
	for k, v := range headers {
		if strings.HasPrefix(k, "x-") && v != "" {
			result[k] = v
		}
	}
	return result
}
