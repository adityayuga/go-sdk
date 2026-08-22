# NSQ Publisher Package

A Go package for publishing messages to NSQ, following the repository's coding patterns and conventions.

## Features

- JSON message encoding
- Raw byte message publishing
- Configurable connection settings
- Context-aware operations
- Interface-based design for testability

## Installation

This package uses the NSQ client library `github.com/nsqio/go-nsq` which is already included in the module dependencies.

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "time"
    
    "github.com/adityayuga/go-sdk/publisher/nsq"
)

func main() {
    ctx := context.Background()
    
    // Configure the NSQ publisher
    cfg := nsq.Config{
        NsqdTCPAddrs:      []string{"localhost:4150"},
        MaxInFlight:       200,
        DialTimeout:       1 * time.Second,
        WriteTimeout:      1 * time.Second,
        ReadTimeout:       60 * time.Second,
        HeartbeatInterval: 30 * time.Second,
    }
    
    // Create publisher
    publisher, err := nsq.New(ctx, cfg)
    if err != nil {
        panic(err)
    }
    defer publisher.Close(ctx)
    
    // Publish a structured message
    message := map[string]interface{}{
        "user_id": 123,
        "action":  "login",
        "timestamp": time.Now(),
    }
    
    err = publisher.Publish(ctx, "user-events", message)
    if err != nil {
        panic(err)
    }
    
    // Publish raw bytes
    err = publisher.PublishRaw(ctx, "raw-events", []byte("raw message"))
    if err != nil {
        panic(err)
    }
}
```

## Configuration

### Config Struct

- `NsqdTCPAddrs`: List of NSQD TCP addresses
- `MaxInFlight`: Maximum number of messages to allow in flight
- `DialTimeout`: Timeout for dialing NSQ
- `WriteTimeout`: Timeout for writing messages
- `ReadTimeout`: Timeout for reading responses
- `HeartbeatInterval`: Heartbeat interval

### JSON Configuration Example

```json
{
  "nsqd_tcp_addrs": ["localhost:4150", "localhost:4152"],
  "max_in_flight": 200,
  "dial_timeout": "1s",
  "write_timeout": "1s",
  "read_timeout": "60s",
  "heartbeat_interval": "30s"
}
```

### Go Configuration Example

```go
cfg := nsq.Config{
    NsqdTCPAddrs:      []string{"localhost:4150"},
    MaxInFlight:       200,
    DialTimeout:       1 * time.Second,
    WriteTimeout:      1 * time.Second,
    ReadTimeout:       60 * time.Second,
    HeartbeatInterval: 30 * time.Second,
}
```

### Constants

- `DefaultMaxInFlight`: Default max in flight (200)
- `DefaultDialTimeoutSeconds`: Default dial timeout (1s)
- `DefaultWriteTimeoutSeconds`: Default write timeout (1s)
- `DefaultReadTimeoutSeconds`: Default read timeout (60s)
- `DefaultHeartbeatIntervalSeconds`: Default heartbeat interval (30s)
- `MaxMessageSize`: Maximum message size (1MB)

## Interface

```go
type Publisher interface {
    Publish(ctx context.Context, topic string, message interface{}) error
    PublishRaw(ctx context.Context, topic string, message []byte) error
    Close(ctx context.Context) error
}
```

## Message Structure

The `Message` struct represents the structure of messages:

```go
type Message struct {
    Topic     string      `json:"topic"`
    Body      interface{} `json:"body"`
    Timestamp time.Time   `json:"timestamp"`
}
```

## Features Implemented

- ✅ Full NSQ producer integration using `github.com/nsqio/go-nsq`
- ✅ Support for multiple NSQD addresses
- ✅ Configurable connection settings
- ✅ Message size validation
- ✅ Context-aware operations
- ✅ Graceful shutdown of all producers

## TODO

- [ ] Add round-robin or load balancing between multiple producers
- [ ] Add support for message priorities
- [ ] Add retry mechanisms
- [ ] Add metrics and monitoring
- [ ] Add connection health checks