package kafka

import "time"

const (
	// DefaultSessionTimeout is the maximum time the broker waits for a heartbeat
	// before declaring the consumer dead and triggering a rebalance.
	// This is set generously to 120 s to accommodate slow message processors
	// (e.g. LLM inference) without causing duplicate redelivery.
	DefaultSessionTimeout = 120 * time.Second

	// DefaultRebalanceTimeout is the maximum time allowed for the consumer group
	// to complete a rebalance.
	DefaultRebalanceTimeout = 120 * time.Second

	// DefaultHeartbeatInterval is how often the consumer sends heartbeats to the
	// broker. Must be less than SessionTimeout / 3.
	DefaultHeartbeatInterval = 10 * time.Second

	// DefaultFetchErrorDelay is the backoff applied after a transient fetch error
	// to avoid a tight retry loop.
	DefaultFetchErrorDelay = 1 * time.Second

	// HeaderDebugID is the Kafka message header key carrying the trace ID.
	HeaderDebugID = "x-debug-id"
)
