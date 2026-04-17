// Package audit provides a structured event log for recording significant
// portwatch lifecycle events such as state transitions, alert dispatches,
// and configuration reloads.
package audit

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Kind classifies the type of audit event.
type Kind string

const (
	KindStateChange   Kind = "state_change"
	KindAlertSent     Kind = "alert_sent"
	KindAlertDropped  Kind = "alert_dropped"
	KindConfigReload  Kind = "config_reload"
	KindStartup       Kind = "startup"
	KindShutdown      Kind = "shutdown"
)

// Event is a single audit log entry.
type Event struct {
	At      time.Time
	Kind    Kind
	Target  string
	Message string
}

func (e Event) String() string {
	return fmt.Sprintf("%s [%s] %s: %s", e.At.Format(time.RFC3339), e.Kind, e.Target, e.Message)
}

// Logger writes audit events to an io.Writer.
type Logger struct {
	mu sync.Mutex
	w  io.Writer
}

// New returns a Logger that writes to w.
// Panics if w is nil.
func New(w io.Writer) *Logger {
	if w == nil {
		panic("audit: writer must not be nil")
	}
	return &Logger{w: w}
}

// Log records an audit event.
func (l *Logger) Log(kind Kind, target, message string) {
	e := Event{
		At:      time.Now().UTC(),
		Kind:    kind,
		Target:  target,
		Message: message,
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintln(l.w, e.String())
}
