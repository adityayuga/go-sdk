package nsq

const (
	// Default configuration values
	DefaultMaxInFlight              = 200
	DefaultDialTimeoutSeconds       = 1
	DefaultWriteTimeoutSeconds      = 1
	DefaultReadTimeoutSeconds       = 60
	DefaultHeartbeatIntervalSeconds = 30

	// Message handling
	MaxMessageSize = 1024 * 1024 // 1MB
)
