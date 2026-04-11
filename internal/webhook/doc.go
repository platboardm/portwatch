// Package webhook delivers portwatch alerts to an external HTTP endpoint
// as JSON POST requests.
//
// Basic usage:
//
//	n := webhook.New("https://hooks.example.com/portwatch", 5*time.Second)
//	if err := n.Notify(ctx, myAlert); err != nil {
//		log.Println("webhook error:", err)
//	}
//
// For resilience against transient failures, wrap the Notifier in a
// RetryNotifier:
//
//	rn := webhook.NewRetry(n, 3, 500*time.Millisecond)
//	if err := rn.Notify(ctx, myAlert); err != nil {
//		log.Println("all attempts failed:", err)
//	}
//
// The JSON payload includes timestamp (RFC 3339), host, port, status,
// severity, and a human-readable message field.
package webhook
