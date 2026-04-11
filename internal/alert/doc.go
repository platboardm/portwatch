// Package alert provides the Alert type used throughout portwatch to
// represent discrete notification events.
//
// An Alert captures:
//   - the target (host:port) that triggered the event
//   - a Severity (INFO for recovery, CRITICAL for outage)
//   - a human-readable message
//   - the UTC timestamp at which the alert was created
//
// Alerts are produced by the monitor package whenever a port transitions
// between up and down states, and are consumed by the notifier package
// which formats and writes them to the configured output sink.
//
// Example:
//
//	a := alert.New("localhost:8080", alert.SeverityCritical, "port unreachable")
//	fmt.Println(a) // [CRITICAL] 2024-01-15T10:30:00Z localhost:8080 — port unreachable
package alert
