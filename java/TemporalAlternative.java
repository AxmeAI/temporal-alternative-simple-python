/*
 * Temporal alternative — order fulfillment without workers or determinism constraints.
 *
 * Temporal needs Activities, Workflows, Workers, and deterministic replay.
 * AXME needs one intent.
 *
 * Usage:
 *   export AXME_API_KEY="your-key"
 *   mvn compile exec:java -Dexec.mainClass="TemporalAlternative"
 */

import dev.axme.sdk.AxmeClient;
import dev.axme.sdk.AxmeClientConfig;
import dev.axme.sdk.RequestOptions;
import dev.axme.sdk.ObserveOptions;
import java.util.List;
import java.util.Map;

public class TemporalAlternative {
    public static void main(String[] args) throws Exception {
        var client = new AxmeClient(
            AxmeClientConfig.forCloud(System.getenv("AXME_API_KEY"))
        );

        // Submit order fulfillment — replaces Temporal Workflow + 3 Activities + Worker
        String intentId = client.sendIntent(Map.of(
            "intent_type", "order.fulfill.v1",
            "to_agent", "agent://myorg/production/order-service",
            "payload", Map.of(
                "order_id", "ORD-123",
                "items", List.of(
                    Map.of("sku", "WIDGET-A", "quantity", 2, "price", 29.99),
                    Map.of("sku", "GADGET-B", "quantity", 1, "price", 40.01)
                ),
                "total", 99.99,
                "shipping_address", Map.of(
                    "street", "123 Main St",
                    "city", "San Francisco",
                    "state", "CA",
                    "zip", "94105"
                )
            )
        ), new RequestOptions());
        System.out.println("Intent submitted: " + intentId);

        // Wait for completion — no polling, no webhooks, no worker
        var result = client.waitFor(intentId, new ObserveOptions());
        System.out.println("Final status: " + result.get("status"));
    }
}
