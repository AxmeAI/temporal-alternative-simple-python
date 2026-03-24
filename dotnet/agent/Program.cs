// Order fulfillment agent — .NET example.
//
// Fetches an intent by ID, processes order (inventory, payment, shipping),
// and resumes with result.
//
// Usage:
//   export AXME_API_KEY="<agent-key>"
//   dotnet run -- <intent_id>

using Axme.Sdk;
using System.Text.Json.Nodes;

if (args.Length < 1)
{
    Console.Error.WriteLine("Usage: dotnet run -- <intent_id>");
    return 1;
}

var apiKey = Environment.GetEnvironmentVariable("AXME_API_KEY");
if (string.IsNullOrEmpty(apiKey))
{
    Console.Error.WriteLine("Error: AXME_API_KEY not set.");
    return 1;
}

var intentId = args[0];
var client = new AxmeClient(new AxmeClientConfig { ApiKey = apiKey });

Console.WriteLine($"Processing intent: {intentId}");

var intentData = await client.GetIntentAsync(intentId);
var intent = intentData["intent"]?.AsObject() ?? intentData;
var payload = intent["payload"]?.AsObject() ?? new JsonObject();
if (payload["parent_payload"] is JsonObject parentPayload)
{
    payload = parentPayload;
}

var orderId = payload["order_id"]?.ToString() ?? "unknown";
var items = payload["items"]?.AsArray() ?? new JsonArray();
var customer = payload["customer_id"]?.ToString() ?? "unknown";

// Step 1: Validate inventory
Console.WriteLine($"  [1/3] Validating inventory for {items.Count} item(s)...");
await Task.Delay(1000);

// Step 2: Charge payment
double total = 0;
foreach (var item in items)
{
    var price = item?["price"]?.GetValue<double>() ?? 0;
    var qty = item?["quantity"]?.GetValue<double>() ?? 1;
    total += price * qty;
}
Console.WriteLine($"  [2/3] Charging ${total:F2} for customer {customer}...");
await Task.Delay(1000);

// Step 3: Create shipment
Console.WriteLine($"  [3/3] Creating shipment for order {orderId}...");
await Task.Delay(1000);

var result = new JsonObject
{
    ["action"] = "complete",
    ["order_id"] = orderId,
    ["total_charged"] = total,
    ["tracking_number"] = "TRK-98765",
    ["shipped_at"] = DateTime.UtcNow.ToString("o")
};

await client.ResumeIntentAsync(intentId, result);
Console.WriteLine($"  Order {orderId} fulfilled. Tracking: TRK-98765");
return 0;
