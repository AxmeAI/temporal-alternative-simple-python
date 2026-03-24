"""
Order fulfillment agent - validates inventory, charges payment, ships.

In Temporal you'd need 3 Activities + 1 Workflow + 1 Worker.
Here it's one agent that processes the intent and resumes.

Usage:
    export AXME_API_KEY="<agent-key>"
    python agent.py
"""

import os
import sys
import time

sys.stdout.reconfigure(line_buffering=True)

from axme import AxmeClient, AxmeClientConfig


AGENT_ADDRESS = "order-fulfillment-demo"


def handle_intent(client, intent_id):
    """Process order: validate inventory, charge payment, ship."""
    intent_data = client.get_intent(intent_id)
    intent = intent_data.get("intent", intent_data)
    payload = intent.get("payload", {})
    if "parent_payload" in payload:
        payload = payload["parent_payload"]

    order_id = payload.get("order_id", "unknown")
    items = payload.get("items", [])
    customer = payload.get("customer_id", "unknown")

    # Step 1: Validate inventory
    print(f"  [1/3] Validating inventory for {len(items)} item(s)...")
    time.sleep(1)

    # Step 2: Charge payment
    total = sum(item.get("price", 0) * item.get("quantity", 1) for item in items)
    print(f"  [2/3] Charging ${total:.2f} for customer {customer}...")
    time.sleep(1)

    # Step 3: Create shipment
    print(f"  [3/3] Creating shipment for order {order_id}...")
    time.sleep(1)

    result = {
        "action": "complete",
        "order_id": order_id,
        "total_charged": total,
        "tracking_number": "TRK-98765",
        "shipped_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
    }

    client.resume_intent(intent_id, result)
    print(f"  Order {order_id} fulfilled. Tracking: {result['tracking_number']}")


def main():
    api_key = os.environ.get("AXME_API_KEY", "")
    if not api_key:
        print("Error: AXME_API_KEY not set.")
        sys.exit(1)

    client = AxmeClient(AxmeClientConfig(api_key=api_key))
    print(f"Agent listening on {AGENT_ADDRESS}...")
    print("Waiting for intents (Ctrl+C to stop)\n")

    for delivery in client.listen(AGENT_ADDRESS):
        intent_id = delivery.get("intent_id", "")
        status = delivery.get("status", "")
        if not intent_id:
            continue
        if status in ("DELIVERED", "CREATED", "IN_PROGRESS"):
            print(f"[{status}] Intent received: {intent_id}")
            try:
                handle_intent(client, intent_id)
            except Exception as e:
                print(f"  Error: {e}")


if __name__ == "__main__":
    main()
