package kafka

const (
	// Acknowledgment settings
	RequiredAcksNone   = 0  // No response required
	RequiredAcksLeader = 1  // Wait for leader acknowledgment
	RequiredAcksAll    = -1 // Wait for all in-sync replicas

	// Default timeouts
	DefaultTimeoutSeconds      = 10
	DefaultBatchTimeoutSeconds = 1
	DefaultBatchSize           = 100

	// HeaderDebugID is the Kafka message header key used to propagate a debug
	// trace ID from the publisher through to any downstream consumers.
	HeaderDebugID = "x-debug-id"

	// HeaderAppsID propagates the originating application identifier.
	HeaderAppsID = "x-apps-id"

	// HeaderUserID propagates the authenticated user ID (omitted when no user
	// is logged in).
	HeaderUserID = "x-user-id"

	// HeaderFromServiceIP records the outbound IP of the service that
	// published this message, enabling hop-by-hop tracing.
	HeaderFromServiceIP = "x-from-service-ip"
)
