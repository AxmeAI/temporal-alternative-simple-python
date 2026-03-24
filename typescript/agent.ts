/**
 * Order fulfillment agent — TypeScript example.
 *
 * Validates inventory, charges payment, ships. In Temporal you'd need
 * 3 Activities + 1 Workflow + 1 Worker. Here it's one agent.
 *
 * Usage:
 *   export AXME_API_KEY="<agent-key>"
 *   npx tsx agent.ts
 */

import { AxmeClient } from "@axme/axme";

const AGENT_ADDRESS = "order-fulfillment-demo";

async function handleIntent(client: AxmeClient, intentId: string) {
  const intentData = await client.getIntent(intentId);
  const intent = intentData.intent ?? intentData;
  let payload = intent.payload ?? {};
  if (payload.parent_payload) {
    payload = payload.parent_payload;
  }

  const orderId = payload.order_id ?? "unknown";
  const items: any[] = payload.items ?? [];
  const customer = payload.customer_id ?? "unknown";

  // Step 1: Validate inventory
  console.log(`  [1/3] Validating inventory for ${items.length} item(s)...`);
  await new Promise((r) => setTimeout(r, 1000));

  // Step 2: Charge payment
  const total = items.reduce(
    (sum: number, item: any) => sum + (item.price ?? 0) * (item.quantity ?? 1),
    0
  );
  console.log(`  [2/3] Charging $${total.toFixed(2)} for customer ${customer}...`);
  await new Promise((r) => setTimeout(r, 1000));

  // Step 3: Create shipment
  console.log(`  [3/3] Creating shipment for order ${orderId}...`);
  await new Promise((r) => setTimeout(r, 1000));

  const result = {
    action: "complete",
    order_id: orderId,
    total_charged: total,
    tracking_number: "TRK-98765",
    shipped_at: new Date().toISOString(),
  };

  await client.resumeIntent(intentId, result, { ownerAgent: "order-fulfillment-demo" });
  console.log(`  Order ${orderId} fulfilled. Tracking: ${result.tracking_number}`);
}

async function main() {
  const apiKey = process.env.AXME_API_KEY;
  if (!apiKey) {
    console.error("Error: AXME_API_KEY not set.");
    process.exit(1);
  }

  const client = new AxmeClient({ apiKey });

  console.log(`Agent listening on ${AGENT_ADDRESS}...`);
  console.log("Waiting for intents (Ctrl+C to stop)\n");

  for await (const delivery of client.listen(AGENT_ADDRESS)) {
    const intentId = delivery.intent_id;
    const status = delivery.status;
    if (!intentId) continue;
    if (["DELIVERED", "CREATED", "IN_PROGRESS"].includes(status)) {
      console.log(`[${status}] Intent received: ${intentId}`);
      try {
        await handleIntent(client, intentId);
      } catch (e) {
        console.error(`  Error: ${e}`);
      }
    }
  }
}

main().catch(console.error);
