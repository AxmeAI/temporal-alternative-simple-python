// Temporal alternative — order fulfillment without workers or determinism constraints.
//
// Temporal needs Activities, Workflows, Workers, and deterministic replay.
// AXME needs one intent.
//
// Usage:
//   export AXME_API_KEY="your-key"
//   dotnet run

using Axme.Sdk;
using System.Text.Json.Nodes;

var client = new AxmeClient(new AxmeClientConfig
{
    ApiKey = Environment.GetEnvironmentVariable("AXME_API_KEY")!
});

// Submit order fulfillment — replaces Temporal Workflow + 3 Activities + Worker
var intentId = await client.SendIntentAsync(new JsonObject
{
    ["intent_type"] = "order.fulfill.v1",
    ["to_agent"] = "agent://myorg/production/order-service",
    ["payload"] = new JsonObject
    {
        ["order_id"] = "ORD-123",
        ["items"] = new JsonArray
        {
            new JsonObject { ["sku"] = "WIDGET-A", ["quantity"] = 2, ["price"] = 29.99 },
            new JsonObject { ["sku"] = "GADGET-B", ["quantity"] = 1, ["price"] = 40.01 }
        },
        ["total"] = 99.99,
        ["shipping_address"] = new JsonObject
        {
            ["street"] = "123 Main St",
            ["city"] = "San Francisco",
            ["state"] = "CA",
            ["zip"] = "94105"
        }
    }
});
Console.WriteLine($"Intent submitted: {intentId}");

// Wait for completion — no polling, no webhooks, no worker
var result = await client.WaitForAsync(intentId);
Console.WriteLine($"Final status: {result["status"]}");
