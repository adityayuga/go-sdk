# Kafka Publisher Package

A Go package for publishing messages to Apache Kafka, following the repository's coding patterns and conventions.

## Features

- JSON message encoding
- Raw byte message publishing  
- Configurable timeouts and batch settings
- Context-aware operations
- Interface-based design for testability

## Installation

This package uses the Kafka client library `github.com/segmentio/kafka-go` which is already included in the module dependencies.

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "time"
    
    "github.com/adityayuga/go-sdk/publisher/kafka"
)

func main() {
    ctx := context.Background()
    
    // Configure the Kafka publisher
    cfg := kafka.Config{
        Brokers:      []string{"localhost:9092"},
        RequiredAcks: kafka.RequiredAcksLeader,
        Timeout:      10 * time.Second,
        BatchSize:    100,
        BatchTimeout: 1 * time.Second,
    }
    
    // Create publisher
    publisher, err := kafka.New(ctx, cfg)
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
    
    err = publisher.Publish(ctx, "user-events", "user-123", message)
    if err != nil {
        panic(err)
    }
    
    // Publish raw bytes
    err = publisher.PublishRaw(ctx, "raw-events", "key", []byte("raw message"))
    if err != nil {
        panic(err)
    }
}
```

## Configuration

### Config Struct

- `Brokers`: List of Kafka broker addresses
- `RequiredAcks`: Acknowledgment level (see constants)
- `Timeout`: Write timeout for messages
- `BatchSize`: Maximum number of messages in a batch
- `BatchTimeout`: Maximum time to wait for a batch

### JSON Configuration Example

```json
{
  "brokers": ["localhost:9092", "localhost:9093"],
  "required_acks": 1,
  "timeout": "10s",
  "batch_size": 100,
  "batch_timeout": "1s"
}
```

### Go Configuration Example

```go
cfg := kafka.Config{
    Brokers:      []string{"localhost:9092"},
    RequiredAcks: kafka.RequiredAcksLeader,
    Timeout:      10 * time.Second,
    BatchSize:    100,
    BatchTimeout: 1 * time.Second,
}
```

### Constants

- `RequiredAcksNone`: No response required (0)
- `RequiredAcksLeader`: Wait for leader acknowledgment (1)
- `RequiredAcksAll`: Wait for all in-sync replicas (-1)

## Interface

```go
type Publisher interface {
    Publish(ctx context.Context, topic string, key string, message interface{}) error
    PublishRaw(ctx context.Context, topic string, key string, message []byte) error
    Close(ctx context.Context) error
}
```

## Features Implemented

- ✅ Full Kafka producer integration using `github.com/segmentio/kafka-go`
- ✅ Synchronous message publishing with error handling
- ✅ Configurable acknowledgment levels
- ✅ Batch processing support
- ✅ Message key support for partitioning
- ✅ Context-aware operations
- ✅ Graceful shutdown

## TODO

- [ ] Add support for message headers
- [ ] Add retry mechanisms
- [ ] Add metrics and monitoring
- [ ] Add compression support
- [ ] Add async publishing option