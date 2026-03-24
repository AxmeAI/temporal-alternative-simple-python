"""
Temporal alternative — order fulfillment without workers or determinism constraints.

Temporal needs Activities, Workflows, Workers, and deterministic replay.
AXME needs one intent.

Usage:
    pip install axme
    export AXME_API_KEY="your-key"
    python main.py
"""

import os
from axme import AxmeClient, AxmeClientConfig


def main():
    client = AxmeClient(
        AxmeClientConfig(api_key=os.environ["AXME_API_KEY"])
    )

    # Submit order fulfillment — replaces Temporal Workflow + 3 Activities + Worker
    intent_id = client.send_intent(
        {
            "intent_type": "order.fulfill.v1",
            "to_agent": "agent://myorg/production/order-service",
            "payload": {
                "order_id": "ORD-123",
                "items": [
                    {"sku": "WIDGET-A", "quantity": 2, "price": 29.99},
                    {"sku": "GADGET-B", "quantity": 1, "price": 40.01},
                ],
                "total": 99.99,
                "shipping_address": {
                    "street": "123 Main St",
                    "city": "San Francisco",
                    "state": "CA",
                    "zip": "94105",
                },
            },
        }
    )
    print(f"Intent submitted: {intent_id}")

    # Observe lifecycle events in real time (SSE stream, no polling)
    print("Watching lifecycle...")
    for event in client.observe(intent_id):
        status = event.get("status", "")
        print(f"  [{status}] {event.get('event_type', '')}")
        if status in ("COMPLETED", "FAILED", "TIMED_OUT", "CANCELLED"):
            break

    # Fetch final state
    intent = client.get_intent(intent_id)
    print(f"\nFinal status: {intent['intent']['lifecycle_status']}")

    # In Temporal, this required:
    #   - 3 Activity definitions (validate_inventory, charge_payment, ship_order)
    #   - 1 Workflow class with determinism constraints
    #   - 1 Worker process (you run, scale, and monitor it)
    #   - A running Temporal server or Temporal Cloud subscription
    #   - Understanding of replay semantics, task queues, and versioning


if __name__ == "__main__":
    main()
