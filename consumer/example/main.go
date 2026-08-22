// Package main demonstrates how to use the consumer library to consume messages
// from Kafka or NSQ. Run with -h to see all available flags.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	consumerpkg "github.com/adityayuga/go-sdk/consumer"
	kafkaconsumer "github.com/adityayuga/go-sdk/consumer/kafka"
	nsqconsumer "github.com/adityayuga/go-sdk/consumer/nsq"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle SIGINT / SIGTERM for graceful shutdown.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		fmt.Printf("Received signal %v, shutting down...\n", sig)
		cancel()
	}()

	// ------------------------------------------------------------------
	// Example 1: Kafka consumer
	// ------------------------------------------------------------------
	fmt.Println("=== Kafka Consumer Example ===")
	if err := kafkaExample(ctx); err != nil {
		log.Printf("Kafka example error: %v", err)
	}

	// ------------------------------------------------------------------
	// Example 2: NSQ consumer
	// ------------------------------------------------------------------
	fmt.Println("\n=== NSQ Consumer Example ===")
	if err := nsqExample(ctx); err != nil {
		log.Printf("NSQ example error: %v", err)
	}
}

// handler is shared by all examples.
func handler(_ context.Context, msg consumerpkg.Message) error {
	fmt.Printf("Received message id=%s ts=%d headers=%v body=%s\n",
		msg.ID, msg.Timestamp, msg.Headers, string(msg.Body))
	return nil
}

func kafkaExample(ctx context.Context) error {
	cfg := kafkaconsumer.Config{
		Brokers:           []string{"localhost:9092"},
		Topic:             "example-topic",
		GroupID:           "example-group",
		SessionTimeout:    120 * time.Second,
		RebalanceTimeout:  120 * time.Second,
		HeartbeatInterval: 10 * time.Second,
	}

	c, err := kafkaconsumer.New(cfg, handler)
	if err != nil {
		return fmt.Errorf("failed to create Kafka consumer: %w", err)
	}
	defer c.Close()

	// Start blocks until ctx is cancelled.
	return c.Start(ctx)
}

func nsqExample(ctx context.Context) error {
	cfg := nsqconsumer.Config{
		Topic:            "example-topic",
		Channel:          "example-channel",
		LookupdHTTPAddrs: []string{"localhost:4161"},
	}

	c, err := nsqconsumer.New(cfg, handler)
	if err != nil {
		return fmt.Errorf("failed to create NSQ consumer: %w", err)
	}
	defer c.Close()

	// Start blocks until ctx is cancelled.
	return c.Start(ctx)
}
