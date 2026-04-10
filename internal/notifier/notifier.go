// Package notifier provides alerting functionality for portwatch.
// It sends notifications when monitored services change state.
package notifier

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Event represents a state change for a monitored target.
type Event struct {
	Target    string
	Host      string
	Port      int
	Up        bool
	Timestamp time.Time
}

// Notifier sends alerts for port state change events.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Notify formats and writes an alert message for the given event.
func (n *Notifier) Notify(e Event) error {
	status := "DOWN"
	if e.Up {
		status = "UP"
	}
	_, err := fmt.Fprintf(
		n.out,
		"[%s] %s (%s:%d) is %s\n",
		e.Timestamp.Format(time.RFC3339),
		e.Target,
		e.Host,
		e.Port,
		status,
	)
	return err
}
