// Package alert defines the Alert type and severity levels used when
// notifying operators about port state transitions.
package alert

import (
	"fmt"
	"time"
)

// Severity indicates how critical an alert is.
type Severity int

const (
	// SeverityInfo is used when a service recovers (comes back up).
	SeverityInfo Severity = iota
	// SeverityCritical is used when a service goes down.
	SeverityCritical
)

// String returns a human-readable label for the severity.
func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "INFO"
	case SeverityCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// Alert represents a single alerting event for a monitored target.
type Alert struct {
	// Target is the host:port that triggered the alert.
	Target string
	// Severity describes how critical the event is.
	Severity Severity
	// Message is a human-readable description of the event.
	Message string
	// OccurredAt is when the alert was generated.
	OccurredAt time.Time
}

// New constructs an Alert with OccurredAt set to the current UTC time.
func New(target string, severity Severity, message string) Alert {
	return Alert{
		Target:     target,
		Severity:   severity,
		Message:    message,
		OccurredAt: time.Now().UTC(),
	}
}

// String returns a single-line representation suitable for log output.
func (a Alert) String() string {
	return fmt.Sprintf("[%s] %s %s — %s",
		a.Severity,
		a.OccurredAt.Format(time.RFC3339),
		a.Target,
		a.Message,
	)
}
