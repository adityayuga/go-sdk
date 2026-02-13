# Log Package

A lightweight, production-ready JSON logger SDK for Go that wraps the industry-standard `zap` logger with context-aware logging capabilities.

## Features

- **JSON Formatting**: All logs are formatted as JSON for easy parsing and analysis
- **Context-Aware**: Automatically extracts and includes request IDs, trace IDs, and user IDs from context
- **Multiple Log Levels**: Support for debug, info, warn, error, and fatal levels
- **Production Ready**: Built on `go.uber.org/zap`, one of the fastest and most reliable loggers in Go
- **Global & Instance Loggers**: Use global functions or create logger instances
- **Easy Context Management**: Helper functions to add and retrieve values from context

## Installation

```bash
go get github.com/adityayuga/go-sdk
```

## Usage

### Basic Initialization

```go
import (
    "context"
    "github.com/adityayuga/go-sdk/log"
)

// Initialize the global logger
err := log.Init("info") // or "debug", "warn", "error"
if err != nil {
    panic(err)
}
defer log.GetLogger().Sync()

ctx := context.Background()

// Log messages with fields
log.Info(ctx, "user login", "username", "john", "email", "john@example.com")
log.Warn(ctx, "slow query", "duration_ms", 1500)
log.Error(ctx, "database error", "error", err)
```

### Using Context for Trace Information

```go
// Add trace information to context
ctx := context.Background()
ctx = log.WithRequestID(ctx, "req-12345")
ctx = log.WithTraceID(ctx, "trace-67890")
ctx = log.WithUserID(ctx, "user-456")

// These values are automatically included in all logs from this context
log.Info(ctx, "processing payment", "amount", 99.99)
// Output: {"level":"info","ts":"2026-02-13T10:00:00.000Z","msg":"processing payment","request_id":"req-12345","trace_id":"trace-67890","user_id":"user-456","amount":99.99}
```

### Creating Logger Instances

```go
// Create multiple logger instances with different configurations
logger, err := log.NewLogger("debug")
if err != nil {
    panic(err)
}

ctx := context.Background()
logger.Info(ctx, "instance logger message", "key", "value")
logger.Sync()
```

### Available Methods

#### Global Functions
- `log.Info(ctx, message, fields...)` - Info level
- `log.Warn(ctx, message, fields...)` - Warning level
- `log.Error(ctx, message, fields...)` - Error level
- `log.Debug(ctx, message, fields...)` - Debug level
- `log.Fatal(ctx, message, fields...)` - Fatal level (exits process)

#### Logger Instance Methods
- `logger.Info(ctx, message, fields...)`
- `logger.Warn(ctx, message, fields...)`
- `logger.Error(ctx, message, fields...)`
- `logger.Debug(ctx, message, fields...)`
- `logger.Fatal(ctx, message, fields...)`
- `logger.Sync()` - Flushes buffered logs

#### Context Helper Functions
- `log.WithRequestID(ctx, id)` - Add request ID
- `log.WithTraceID(ctx, id)` - Add trace ID
- `log.WithUserID(ctx, id)` - Add user ID
- `log.WithLogContext(ctx, requestID, traceID, userID)` - Add all at once
- `log.GetRequestID(ctx)` - Retrieve request ID
- `log.GetTraceID(ctx)` - Retrieve trace ID
- `log.GetUserID(ctx)` - Retrieve user ID

## Development Mode

For development environments with pretty-printed logs:

```go
err := log.InitDevelopment()
if err != nil {
    panic(err)
}
```

## Log Output Example

```json
{
  "level": "error",
  "ts": "2026-02-13T10:00:00.000Z",
  "msg": "failed to fetch user",
  "request_id": "req-12345",
  "trace_id": "trace-67890",
  "user_id": "user-456",
  "user_id": 123,
  "error": "connection timeout"
}
```

## Best Practices

1. **Always defer Sync()**: Call `logger.Sync()` before exiting to ensure all logs are flushed
2. **Use Context**: Pass context through your request handlers to maintain trace information
3. **Structured Logging**: Use key-value pairs instead of string formatting for better searchability
4. **Log Levels**: Use appropriate levels (debug for development, info for general, error for problems)
5. **Avoid Sensitive Data**: Don't log passwords, tokens, or personal information

## Performance

Built on `zap`, this logger is optimized for performance:
- Minimal allocations
- Efficient JSON encoding
- Designed for high-throughput services

## License

See LICENSE file in the repository
