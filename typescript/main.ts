/**
 * Temporal alternative — order fulfillment without workers or determinism constraints.
 *
 * Temporal needs Activities, Workflows, Workers, and deterministic replay.
 * AXME needs one intent.
 *
 * Usage:
 *   npm install @axme/axme
 *   export AXME_API_KEY="your-key"
 *   npx tsx main.ts
 */

import { AxmeClient } from "@axme/axme";

async function main() {
  const client = new AxmeClient({ apiKey: process.env.AXME_API_KEY! });

  // Submit order fulfillment — replaces Temporal Workflow + 3 Activities + Worker
  const intentId = await client.sendIntent({
    intentType: "order.fulfill.v1",
    toAgent: "agent://myorg/production/order-service",
    payload: {
      orderId: "ORD-123",
      items: [
        { sku: "WIDGET-A", quantity: 2, price: 29.99 },
        { sku: "GADGET-B", quantity: 1, price: 40.01 },
      ],
      total: 99.99,
      shippingAddress: {
        street: "123 Main St",
        city: "San Francisco",
        state: "CA",
        zip: "94105",
      },
    },
  });
  console.log(`Intent submitted: ${intentId}`);

  // Wait for completion — no polling, no webhooks, no worker
  const result = await client.waitFor(intentId);
  console.log(`Final status: ${result.status}`);
}

main().catch(console.error);
