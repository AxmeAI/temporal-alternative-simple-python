/*
 * Order fulfillment agent — Java example.
 *
 * Fetches an intent by ID, processes order (inventory, payment, shipping),
 * and resumes with result.
 *
 * Usage:
 *   export AXME_API_KEY="<agent-key>"
 *   javac -cp axme-sdk.jar Agent.java
 *   java -cp .:axme-sdk.jar Agent <intent_id>
 */

import dev.axme.sdk.AxmeClient;
import dev.axme.sdk.AxmeClientConfig;
import dev.axme.sdk.RequestOptions;
import java.time.Instant;
import java.util.List;
import java.util.Map;

public class Agent {
    public static void main(String[] args) throws Exception {
        if (args.length < 1) {
            System.err.println("Usage: java Agent <intent_id>");
            System.exit(1);
        }

        String apiKey = System.getenv("AXME_API_KEY");
        if (apiKey == null || apiKey.isEmpty()) {
            System.err.println("Error: AXME_API_KEY not set.");
            System.exit(1);
        }

        String intentId = args[0];
        var client = new AxmeClient(AxmeClientConfig.forCloud(apiKey));

        System.out.println("Processing intent: " + intentId);

        var intentData = client.getIntent(intentId, new RequestOptions());
        @SuppressWarnings("unchecked")
        var intent = (Map<String, Object>) intentData.getOrDefault("intent", intentData);
        @SuppressWarnings("unchecked")
        var payload = (Map<String, Object>) intent.getOrDefault("payload", Map.of());
        if (payload.containsKey("parent_payload")) {
            @SuppressWarnings("unchecked")
            var pp = (Map<String, Object>) payload.get("parent_payload");
            payload = pp;
        }

        String orderId = (String) payload.getOrDefault("order_id", "unknown");
        @SuppressWarnings("unchecked")
        var items = (List<Map<String, Object>>) payload.getOrDefault("items", List.of());
        String customer = (String) payload.getOrDefault("customer_id", "unknown");

        // Step 1: Validate inventory
        System.out.println("  [1/3] Validating inventory for " + items.size() + " item(s)...");
        Thread.sleep(1000);

        // Step 2: Charge payment
        double total = 0;
        for (var item : items) {
            double price = item.containsKey("price") ? ((Number) item.get("price")).doubleValue() : 0;
            double qty = item.containsKey("quantity") ? ((Number) item.get("quantity")).doubleValue() : 1;
            total += price * qty;
        }
        System.out.printf("  [2/3] Charging $%.2f for customer %s...%n", total, customer);
        Thread.sleep(1000);

        // Step 3: Create shipment
        System.out.println("  [3/3] Creating shipment for order " + orderId + "...");
        Thread.sleep(1000);

        var result = Map.<String, Object>of(
            "action", "complete",
            "order_id", orderId,
            "total_charged", total,
            "tracking_number", "TRK-98765",
            "shipped_at", Instant.now().toString()
        );

        client.resumeIntent(intentId, result, new RequestOptions());
        System.out.println("  Order " + orderId + " fulfilled. Tracking: TRK-98765");
    }
}
