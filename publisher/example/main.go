package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/adityayuga/go-sdk/publisher"
	"github.com/adityayuga/go-sdk/publisher/kafka"
	"github.com/adityayuga/go-sdk/publisher/nsq"
)

func main() {
	ctx := context.Background()

	// Example 1: Kafka Publisher
	fmt.Println("=== Kafka Publisher Example ===")
	kafkaExample(ctx)

	// Example 2: NSQ Publisher
	fmt.Println("\n=== NSQ Publisher Example ===")
	nsqExample(ctx)

	// Example 3: Using both through the common interface
	fmt.Println("\n=== Polymorphic Publisher Example ===")
	polymorphicExample(ctx)
}

func kafkaExample(ctx context.Context) {
	// Configure Kafka publisher
	cfg := kafka.Config{
		Brokers:      []string{"localhost:9092"},
		RequiredAcks: kafka.RequiredAcksLeader,
		Timeout:      10 * time.Second,
		BatchSize:    100,
		BatchTimeout: 1 * time.Second,
	}

	// Create publisher
	pub, err := kafka.New(ctx, cfg)
	if err != nil {
		log.Printf("Failed to create Kafka publisher: %v", err)
		return
	}
	defer pub.Close(ctx)

	// Publish a structured message
	message := map[string]interface{}{
		"user_id":   123,
		"action":    "login",
		"timestamp": time.Now(),
	}

	err = pub.Publish(ctx, "user-events", "user-123", message)
	if err != nil {
		log.Printf("Failed to publish to Kafka: %v", err)
		return
	}

	fmt.Println("✓ Successfully published message to Kafka")
}

func nsqExample(ctx context.Context) {
	// Configure NSQ publisher
	cfg := nsq.Config{
		NsqdTCPAddrs:      []string{"localhost:4150"},
		MaxInFlight:       200,
		DialTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		ReadTimeout:       60 * time.Second,
		HeartbeatInterval: 30 * time.Second,
	}

	// Create publisher
	pub, err := nsq.New(ctx, cfg)
	if err != nil {
		log.Printf("Failed to create NSQ publisher: %v", err)
		return
	}
	defer pub.Close(ctx)

	// Publish a structured message
	message := map[string]interface{}{
		"user_id":   456,
		"action":    "logout",
		"timestamp": time.Now(),
	}

	// Note: key parameter is ignored for NSQ
	err = pub.Publish(ctx, "user-events", "", message)
	if err != nil {
		log.Printf("Failed to publish to NSQ: %v", err)
		return
	}

	fmt.Println("✓ Successfully published message to NSQ")
}

func polymorphicExample(ctx context.Context) {
	// Create a slice of publishers
	var publishers []publisher.Publisher

	// Add Kafka publisher
	kafkaCfg := kafka.Config{
		Brokers: []string{"localhost:9092"},
	}
	kafkaPub, err := kafka.New(ctx, kafkaCfg)
	if err != nil {
		log.Printf("Failed to create Kafka publisher: %v", err)
	} else {
		publishers = append(publishers, kafkaPub)
	}

	// Add NSQ publisher
	nsqCfg := nsq.Config{
		NsqdTCPAddrs: []string{"localhost:4150"},
	}
	nsqPub, err := nsq.New(ctx, nsqCfg)
	if err != nil {
		log.Printf("Failed to create NSQ publisher: %v", err)
	} else {
		publishers = append(publishers, nsqPub)
	}

	// Publish to all publishers using the same interface
	message := map[string]interface{}{
		"broadcast": "hello world",
		"timestamp": time.Now(),
	}

	for i, pub := range publishers {
		err := pub.Publish(ctx, "broadcast-topic", "broadcast-key", message)
		if err != nil {
			log.Printf("Publisher %d failed: %v", i, err)
		} else {
			fmt.Printf("✓ Publisher %d succeeded\n", i)
		}
		pub.Close(ctx)
	}
}
