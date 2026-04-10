// Package notifier handles alert delivery for portwatch.
//
// When the checker detects that a monitored port has changed state
// (either gone down or come back up), the notifier formats and emits
// a human-readable message describing the event.
//
// Basic usage:
//
//	n := notifier.New(os.Stdout)
//	n.Notify(notifier.Event{
//		Target:    "api",
//		Host:      "10.0.0.1",
//		Port:      443,
//		Up:        false,
//		Timestamp: time.Now(),
//	})
//
// Future implementations may support additional backends such as
// webhook, email, or Slack notifications.
package notifier
