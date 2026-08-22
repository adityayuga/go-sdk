# Publisher Package

This package provides a unified interface for message publishing across different message queuing systems, following the repository's coding patterns and conventions.

## Common Interface

The `Publisher` interface is defined in the parent package and implemented by all publisher types:

```go
package publisher

import "context"

// Publisher defines the common interface for all message publishers
type Publisher interface {
	// Publish publishes a message to the specified topic with JSON encoding
	// For systems that support keys (like Kafka), key should be provided
	// For systems that don't support keys (like NSQ), key can be empty string
	Publish(ctx context.Context, topic string, key string, message interface{}) error
	
	// PublishRaw publishes raw bytes to the specified topic
	// For systems that support keys (like Kafka), key should be provided  
	// For systems that don't support keys (like NSQ), key can be empty string
	PublishRaw(ctx context.Context, topic string, key string, message []byte) error
	
	// Close closes the publisher and cleans up resources
	Close(ctx context.Context) error
}
```

## Available Publishers

### Kafka Publisher (`./kafka`)

Publisher for Apache Kafka message broker with support for:
- JSON and raw message publishing
- Message keys for partitioning
- Configurable acknowledgment levels
- Batch processing
- Context-aware operations

### NSQ Publisher (`./nsq`)  

Publisher for NSQ distributed messaging platform with support for:
- JSON and raw message publishing
- Configurable connection settings
- Message size limits
- Context-aware operations
- Note: NSQ doesn't support message keys, so the key parameter is ignored

## Usage Example

```go
package main

import (
    "context"
    
    "github.com/adityayuga/go-sdk/publisher"
    "github.com/adityayuga/go-sdk/publisher/kafka"
    "github.com/adityayuga/go-sdk/publisher/nsq"
)

func main() {
    ctx := context.Background()
    
    // Create publishers - both implement publisher.Publisher interface
    var publishers []publisher.Publisher
    
    // Kafka publisher
    kafkaConfig := kafka.Config{
        Brokers: []string{"localhost:9092"},
    }
    kafkaPub, err := kafka.New(ctx, kafkaConfig)
    if err != nil {
        panic(err)
    }
    publishers = append(publishers, kafkaPub)
    
    // NSQ publisher
    nsqConfig := nsq.Config{
        NsqdTCPAddrs: []string{"localhost:4150"},
    }
    nsqPub, err := nsq.New(ctx, nsqConfig)
    if err != nil {
        panic(err)
    }
    publishers = append(publishers, nsqPub)
    
    // Use both publishers through the same interface
    for _, pub := range publishers {
        message := map[string]interface{}{
            "user_id": 123,
            "action":  "login",
        }
        
        // For Kafka: key is used for partitioning
        // For NSQ: key is ignored
        err := pub.Publish(ctx, "user-events", "user-123", message)
        if err != nil {
            // Handle error
        }
        
        defer pub.Close(ctx)
    }
}
```

## Implementation Status

Both Kafka and NSQ publishers are **fully implemented** and production-ready:

- ✅ **Kafka**: Using `github.com/segmentio/kafka-go` v0.4.49
- ✅ **NSQ**: Using `github.com/nsqio/go-nsq` v1.1.0

## Design Principles

- **Unified Interface**: Single `Publisher` interface in parent package for consistency
- **Interface-based**: All publishers implement the common interface for easy testing and swapping
- **Context-aware**: All operations accept `context.Context` for cancellation and timeouts
- **Error handling**: Uses the repository's custom error package
- **Configuration**: Struct-based configuration with sensible defaults
- **Key Support**: Unified signature supports both key-based (Kafka) and keyless (NSQ) systems
- **Consistent API**: Same patterns across all publisher implementations
- **Production Ready**: Full client implementations with proper error handling and resource cleanup