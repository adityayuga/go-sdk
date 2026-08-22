package nsq

const (
	// DefaultMaxInFlight is the maximum number of messages NSQ will deliver to
	// the consumer before requiring acknowledgements.
	DefaultMaxInFlight = 200

	// DefaultDialTimeoutSeconds is the default connection dial timeout.
	DefaultDialTimeoutSeconds = 1

	// DefaultWriteTimeoutSeconds is the default write timeout for NSQ commands.
	DefaultWriteTimeoutSeconds = 1

	// DefaultReadTimeoutSeconds is the default read timeout.
	DefaultReadTimeoutSeconds = 60

	// DefaultHeartbeatIntervalSeconds is the interval at which the consumer
	// sends heartbeats to the NSQ daemon.
	DefaultHeartbeatIntervalSeconds = 30
)
