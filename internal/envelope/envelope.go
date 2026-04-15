// Package envelope wraps an alert with routing metadata such as the
// destination channel, priority, and a unique trace ID so that sinks and
// middleware can make routing decisions without inspecting alert internals.
package envelope

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Priority indicates how urgently an alert should be delivered.
type Priority int

const (
	PriorityLow    Priority = iota // informational
	PriorityNormal                 // default
	PriorityHigh                   // degraded service
	PriorityCritical               // service down
)

// String returns a human-readable label for the priority.
func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "low"
	case PriorityNormal:
		return "normal"
	case PriorityHigh:
		return "high"
	case PriorityCritical:
		return "critical"
	default:
		return fmt.Sprintf("priority(%d)", int(p))
	}
}

// Envelope carries an alert alongside routing metadata.
type Envelope struct {
	TraceID   string
	Alert     alert.Alert
	Priority  Priority
	Channel   string
	CreatedAt time.Time
}

// New wraps a into an Envelope, assigning a random trace ID and the current
// timestamp. channel may be empty if routing is not required.
func New(a alert.Alert, p Priority, channel string) Envelope {
	return Envelope{
		TraceID:   newTraceID(),
		Alert:     a,
		Priority:  p,
		Channel:   channel,
		CreatedAt: time.Now(),
	}
}

func newTraceID() string {
	return fmt.Sprintf("%016x", rand.Uint64())
}
