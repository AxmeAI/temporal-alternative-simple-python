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
	client := axme.NewClient(axme.Config{
		APIKey: os.Getenv("AXME_API_KEY"),
	})

	ctx := context.Background()

	// Submit order fulfillment — replaces Temporal Workflow + 3 Activities + Worker
	intentID, err := client.SendIntent(ctx, axme.SendIntentRequest{
		IntentType: "order.fulfill.v1",
		ToAgent:    "agent://myorg/production/order-service",
		Payload: map[string]interface{}{
			"order_id": "ORD-123",
			"items": []map[string]interface{}{
				{"sku": "WIDGET-A", "quantity": 2, "price": 29.99},
				{"sku": "GADGET-B", "quantity": 1, "price": 40.01},
			},
			"total": 99.99,
			"shipping_address": map[string]interface{}{
				"street": "123 Main St",
				"city":   "San Francisco",
				"state":  "CA",
				"zip":    "94105",
			},
		},
	})
	if err != nil {
		log.Fatalf("send intent: %v", err)
	}
	fmt.Printf("Intent submitted: %s\n", intentID)

	// Wait for completion — no polling, no webhooks, no worker
	result, err := client.WaitFor(ctx, intentID)
	if err != nil {
		log.Fatalf("wait: %v", err)
	}
	fmt.Printf("Final status: %s\n", result.Status)
}
