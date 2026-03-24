// Temporal alternative — order fulfillment without workers or determinism constraints.
//
// Temporal needs Activities, Workflows, Workers, and deterministic replay.
// AXME needs one intent.
//
// Usage:
//
//	export AXME_API_KEY="your-key"
//	go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AxmeAI/axme-sdk-go/axme"
)

func main() {
	client, err := axme.NewClient(axme.ClientConfig{
		APIKey: os.Getenv("AXME_API_KEY"),
	})
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	ctx := context.Background()

	// Submit order fulfillment — replaces Temporal Workflow + 3 Activities + Worker
	intentID, err := client.SendIntent(ctx, map[string]any{
		"intent_type": "order.fulfill.v1",
		"to_agent":    "agent://myorg/production/order-service",
		"order_id":    "ORD-123",
		"items": []map[string]any{
			{"sku": "WIDGET-A", "quantity": 2, "price": 29.99},
			{"sku": "GADGET-B", "quantity": 1, "price": 40.01},
		},
		"total": 99.99,
		"shipping_address": map[string]any{
			"street": "123 Main St",
			"city":   "San Francisco",
			"state":  "CA",
			"zip":    "94105",
		},
	}, axme.RequestOptions{})
	if err != nil {
		log.Fatalf("send intent: %v", err)
	}
	fmt.Printf("Intent submitted: %s\n", intentID)

	// Wait for completion — no polling, no webhooks, no worker
	result, err := client.WaitFor(ctx, intentID, axme.ObserveOptions{})
	if err != nil {
		log.Fatalf("wait: %v", err)
	}
	fmt.Printf("Final status: %v\n", result["status"])
}
