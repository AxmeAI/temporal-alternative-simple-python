# Temporal Alternative — Same Workflow, Simpler Code

Temporal needs Activities, Workflows, Workers, determinism constraints, and a platform team to operate. You need to process an order.

**There is a better way.** One intent replaces the entire Temporal stack — no workers, no determinism rules, no replay constraints. Submit an intent, get a result.

> **Alpha** · Built with [AXME](https://github.com/AxmeAI/axme) (AXP Intent Protocol).
> [cloud.axme.ai](https://cloud.axme.ai) · [hello@axme.ai](mailto:hello@axme.ai)

---

## Before / After

### Before: Temporal Workflow (~300 lines across 4 files)

```python
# activities.py — each step is a separate Activity
from temporalio import activity

@activity.defn
async def validate_inventory(order_id: str) -> dict:
    # Must be deterministic-safe — no random(), no datetime.now(),
    # no threading, no global mutable state
    async with httpx.AsyncClient() as client:
        resp = await client.post(f"{INVENTORY_URL}/validate", json={"order_id": order_id})
        return resp.json()

@activity.defn
async def charge_payment(order_id: str, amount: float) -> dict:
    async with httpx.AsyncClient() as client:
        resp = await client.post(f"{PAYMENT_URL}/charge",
                                 json={"order_id": order_id, "amount": amount})
        return resp.json()

@activity.defn
async def ship_order(order_id: str, address: dict) -> dict:
    async with httpx.AsyncClient() as client:
        resp = await client.post(f"{SHIPPING_URL}/ship",
                                 json={"order_id": order_id, "address": address})
        return resp.json()

# workflow.py — orchestration with determinism constraints
from temporalio import workflow
from temporalio.common import RetryPolicy
from datetime import timedelta

@workflow.defn
class OrderFulfillmentWorkflow:
    @workflow.run
    async def run(self, order: dict) -> dict:
        # WARNING: No non-deterministic calls allowed here.
        # No random(), no datetime.now(), no HTTP calls, no file I/O.
        # Every line must produce the same result on replay.

        inventory = await workflow.execute_activity(
            validate_inventory,
            order["order_id"],
            start_to_close_timeout=timedelta(seconds=30),
            retry_policy=RetryPolicy(maximum_attempts=3),
        )
        if not inventory["available"]:
            return {"status": "failed", "reason": "out_of_stock"}

        payment = await workflow.execute_activity(
            charge_payment,
            args=[order["order_id"], order["total"]],
            start_to_close_timeout=timedelta(seconds=60),
            retry_policy=RetryPolicy(maximum_attempts=3),
        )
        if payment["status"] != "charged":
            return {"status": "failed", "reason": "payment_declined"}

        shipment = await workflow.execute_activity(
            ship_order,
            args=[order["order_id"], order["shipping_address"]],
            start_to_close_timeout=timedelta(seconds=120),
            retry_policy=RetryPolicy(maximum_attempts=3),
        )

        return {"status": "shipped", "tracking": shipment["tracking_number"]}

# worker.py — you must run and operate this process
from temporalio.client import Client
from temporalio.worker import Worker

async def run_worker():
    client = await Client.connect("localhost:7233")  # or your Temporal cluster
    worker = Worker(
        client,
        task_queue="order-fulfillment",
        workflows=[OrderFulfillmentWorkflow],
        activities=[validate_inventory, charge_payment, ship_order],
    )
    await worker.run()  # Runs forever — you manage uptime, scaling, deploys

# starter.py — finally, start the workflow
async def start():
    client = await Client.connect("localhost:7233")
    result = await client.execute_workflow(
        OrderFulfillmentWorkflow.run,
        {"order_id": "ORD-123", "total": 99.99,
         "shipping_address": {"street": "123 Main St", "city": "SF"}},
        id="order-ORD-123",
        task_queue="order-fulfillment",
    )
    print(result)
```

**Plus:** You need a running Temporal server (self-hosted cluster or Temporal Cloud), task queue configuration, worker scaling, and a team that understands deterministic replay semantics.

### After: AXME Intent (15 lines, one file)

```python
from axme import AxmeClient, AxmeClientConfig
import os

client = AxmeClient(AxmeClientConfig(api_key=os.environ["AXME_API_KEY"]))

intent_id = client.send_intent({
    "intent_type": "order.fulfill.v1",
    "to_agent": "agent://myorg/production/order-service",
    "payload": {
        "order_id": "ORD-123",
        "total": 99.99,
        "shipping_address": {"street": "123 Main St", "city": "SF"},
    },
})

result = client.wait_for(intent_id)
print(result["status"])  # COMPLETED, FAILED, or TIMED_OUT
```

No workers. No determinism constraints. No replay semantics to learn. No infrastructure to manage.

---

## Temporal vs AXME

| | Temporal | AXME |
|---|---|---|
| **Model** | Deterministic workflow replay | Intent lifecycle (non-deterministic OK) |
| **Execution** | Worker-based (you run + scale workers) | Protocol-based routing (managed) |
| **Primitives** | Workflow + Activity + Worker + Task Queue | `send_intent()` + `wait_for()` |
| **Constraints** | No `random()`, no `datetime.now()`, no I/O in workflows | None — write normal code |
| **Infrastructure** | Temporal Server (self-hosted or Cloud) | AXME Cloud (managed) |
| **Team required** | Platform team to operate | CLI-first, no infra to manage |
| **Learning curve** | Weeks (replay semantics, versioning, determinism) | Minutes (submit intent, get result) |
| **Best for** | Complex microservice orchestration with replay guarantees | Agent-era workflows: services + agents + humans |

**AXME is not a simplified Temporal.** It is a different paradigm. Temporal replays deterministic workflows. AXME routes intents through a lifecycle — services, AI agents, and humans are all first-class participants.

---

## Quick Start

### Python

```bash
pip install axme
export AXME_API_KEY="your-key"   # Get one: axme login
python python/main.py
```

```python
from axme import AxmeClient, AxmeClientConfig
import os

client = AxmeClient(AxmeClientConfig(api_key=os.environ["AXME_API_KEY"]))

intent_id = client.send_intent({
    "intent_type": "order.fulfill.v1",
    "to_agent": "agent://myorg/production/order-service",
    "payload": {
        "order_id": "ORD-123",
        "total": 99.99,
        "shipping_address": {"street": "123 Main St", "city": "SF"},
    },
})

result = client.wait_for(intent_id)
print(f"Order status: {result['status']}")
```

### TypeScript

```bash
npm install @axme/axme
```

```typescript
import { AxmeClient } from "@axme/axme";

const client = new AxmeClient({ apiKey: process.env.AXME_API_KEY! });

const intentId = await client.sendIntent({
  intentType: "order.fulfill.v1",
  toAgent: "agent://myorg/production/order-service",
  payload: {
    orderId: "ORD-123",
    total: 99.99,
    shippingAddress: { street: "123 Main St", city: "SF" },
  },
});

const result = await client.waitFor(intentId);
console.log(`Order status: ${result.status}`);
```

---

## More Languages

Full implementations in all 5 languages:

| Language | Directory | Install |
|----------|-----------|---------|
| [Python](python/) | `python/` | `pip install axme` |
| [TypeScript](typescript/) | `typescript/` | `npm install @axme/axme` |
| [Go](go/) | `go/` | `go get github.com/AxmeAI/axme-sdk-go` |
| [Java](java/) | `java/` | Maven Central: `ai.axme:axme-sdk` |
| [.NET](dotnet/) | `dotnet/` | `dotnet add package Axme.Sdk` |

---

## How It Works

```
┌──────────┐    send_intent()     ┌──────────────┐    deliver     ┌─────────────┐
│  Client   │ ──────────────────► │  AXME Cloud   │ ────────────► │   Service   │
│           │                     │  (platform)   │               │   (agent)   │
│           │ ◄──── observe() ──  │               │ ◄── resume()  │             │
│           │   real-time SSE     │  retries,     │   with result │  processes  │
└──────────┘                     │  timeouts,    │               │  the work   │
                                  │  delivery     │               └─────────────┘
                                  └──────────────┘
```

1. Client submits an **intent** — "fulfill this order"
2. Platform **delivers** it to the target service/agent
3. Service processes each step (validate, charge, ship) and **resumes** with result
4. Client **observes** lifecycle events in real time (SSE stream)
5. Platform handles retries, timeouts, and delivery guarantees

No worker processes. No determinism constraints. No replay semantics.

---

## When to Use What

**Use Temporal when:**
- You have deep investment in deterministic replay workflows
- You need fine-grained control over workflow versioning and replay
- Your team already operates Temporal infrastructure

**Use AXME when:**
- You want durable execution without the operational overhead
- Your workflows involve AI agents and/or human approvals
- You need cross-service coordination with a simple API
- You want to ship in hours, not weeks

---

## Run the Full Example

```bash
# Install CLI (one-time)
curl -fsSL https://raw.githubusercontent.com/AxmeAI/axme-cli/main/install.sh | sh
source ~/.zshrc

# Log in
axme login

# Run the built-in example
axme examples run delivery/stream
```

---

## Related

- [AXME](https://github.com/AxmeAI/axme) — project overview
- [AXP Spec](https://github.com/AxmeAI/axp-spec) — open Intent Protocol specification
- [AXME Examples](https://github.com/AxmeAI/axme-examples) — 20+ runnable examples across 5 languages
- [AXME CLI](https://github.com/AxmeAI/axme-cli) — manage intents, agents, scenarios from the terminal

---

Built with [AXME](https://github.com/AxmeAI/axme) (AXP Intent Protocol).
