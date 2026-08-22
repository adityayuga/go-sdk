# Publisher Examples

This directory contains examples demonstrating how to use the publisher packages for Kafka and NSQ.

## Running the Examples

### Prerequisites

Before running the examples, you need to have the following services running:

**For Kafka:**
```bash
# Using Docker
docker run -d --name kafka -p 9092:9092 apache/kafka:latest
```

**For NSQ:**
```bash
# Using Docker
docker run -d --name nsqd -p 4150:4150 -p 4151:4151 nsqio/nsq /nsqd
```

### Running the Example

```bash
cd publisher/example
go run main.go
```

## Example Overview

The example demonstrates three usage patterns:

### 1. Kafka Publisher Example
Shows how to:
- Configure a Kafka publisher
- Create a publisher instance
- Publish structured JSON messages
- Handle errors
- Clean up resources

### 2. NSQ Publisher Example
Shows how to:
- Configure an NSQ publisher
- Create a publisher instance
- Publish structured JSON messages
- Handle errors
- Clean up resources

### 3. Polymorphic Publisher Example
Shows how to:
- Use both publishers through the common `publisher.Publisher` interface
- Switch between different publishers without changing code
- Broadcast messages to multiple publishers
- Handle multiple publishers in a unified way

## Key Takeaways

1. **Unified Interface**: Both Kafka and NSQ implement the same `publisher.Publisher` interface
2. **Configuration**: Each publisher has its own configuration struct with sensible defaults
3. **Context Support**: All operations are context-aware for timeout and cancellation
4. **Error Handling**: Proper error handling with descriptive error messages
5. **Resource Cleanup**: Always close publishers using `defer pub.Close(ctx)`
6. **Key Parameter**: Kafka uses keys for partitioning, NSQ ignores them

## Expected Output

When services are running:
```
=== Kafka Publisher Example ===
✓ Successfully published message to Kafka

=== NSQ Publisher Example ===
✓ Successfully published message to NSQ

=== Polymorphic Publisher Example ===
✓ Publisher 0 succeeded
✓ Publisher 1 succeeded
```

When services are not running, you'll see appropriate error messages indicating connection failures.
