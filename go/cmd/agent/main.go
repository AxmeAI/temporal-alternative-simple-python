// Order fulfillment agent — Go example.
//
// Validates inventory, charges payment, ships. In Temporal you'd need
// 3 Activities + 1 Workflow + 1 Worker. Here it's one agent.
//
// Usage:
//
//	export AXME_API_KEY="<agent-key>"
//	go run agent.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AxmeAI/axme-sdk-go/axme"
)

const agentAddress = "order-fulfillment-demo"

func handleIntent(ctx context.Context, client *axme.Client, intentID string) error {
	intentData, err := client.GetIntent(ctx, intentID, axme.RequestOptions{})
	if err != nil {
		return fmt.Errorf("get intent: %w", err)
	}

	intent, _ := intentData["intent"].(map[string]any)
	if intent == nil {
		intent = intentData
	}
	payload, _ := intent["payload"].(map[string]any)
	if payload == nil {
		payload = map[string]any{}
	}
	if pp, ok := payload["parent_payload"].(map[string]any); ok {
		payload = pp
	}

	orderID, _ := payload["order_id"].(string)
	if orderID == "" {
		orderID = "unknown"
	}
	items, _ := payload["items"].([]any)
	customerID, _ := payload["customer_id"].(string)
	if customerID == "" {
		customerID = "unknown"
	}

	// Step 1: Validate inventory
	fmt.Printf("  [1/3] Validating inventory for %d item(s)...\n", len(items))
	time.Sleep(1 * time.Second)

	// Step 2: Charge payment
	total := 0.0
	for _, raw := range items {
		item, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		price, _ := item["price"].(float64)
		qty, _ := item["quantity"].(float64)
		if qty == 0 {
			qty = 1
		}
		total += price * qty
	}
	fmt.Printf("  [2/3] Charging $%.2f for customer %s...\n", total, customerID)
	time.Sleep(1 * time.Second)

	// Step 3: Create shipment
	fmt.Printf("  [3/3] Creating shipment for order %s...\n", orderID)
	time.Sleep(1 * time.Second)

	result := map[string]any{
		"action":          "complete",
		"order_id":        orderID,
		"total_charged":   total,
		"tracking_number": "TRK-98765",
		"shipped_at":      time.Now().UTC().Format(time.RFC3339),
	}

	_, err = client.ResumeIntent(ctx, intentID, result, axme.RequestOptions{})
	if err != nil {
		return fmt.Errorf("resume intent: %w", err)
	}
	fmt.Printf("  Order %s fulfilled. Tracking: %s\n", orderID, result["tracking_number"])
	return nil
}

func main() {
	apiKey := os.Getenv("AXME_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: AXME_API_KEY not set.")
	}

	client, err := axme.NewClient(axme.ClientConfig{APIKey: apiKey})
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	ctx := context.Background()

	fmt.Printf("Agent listening on %s...\n", agentAddress)
	fmt.Println("Waiting for intents (Ctrl+C to stop)")

	intents, errCh := client.Listen(ctx, agentAddress, axme.ListenOptions{})

	go func() {
		for err := range errCh {
			log.Printf("Listen error: %v", err)
		}
	}()

	for delivery := range intents {
		intentID, _ := delivery["intent_id"].(string)
		status, _ := delivery["status"].(string)
		if intentID == "" {
			continue
		}
		if status == "DELIVERED" || status == "CREATED" || status == "IN_PROGRESS" {
			fmt.Printf("[%s] Intent received: %s\n", status, intentID)
			if err := handleIntent(ctx, client, intentID); err != nil {
				fmt.Printf("  Error: %v\n", err)
			}
		}
	}
}
