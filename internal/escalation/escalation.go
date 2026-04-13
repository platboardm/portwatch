// Package escalation provides a mechanism to escalate alerts when a target
// remains in a failed state beyond a configurable duration threshold.
package escalation

import (
	"sync"
	"time"

	"portwatch/internal/alert"
)

// Escalator tracks how long each target has been in a down state and signals
// when the outage duration crosses the configured threshold.
type Escalator struct {
	mu        sync.Mutex
	threshold time.Duration
	downSince map[string]time.Time
}

// New returns an Escalator that escalates alerts whose targets have been
// continuously down for longer than threshold.
// It panics if threshold is zero or negative.
func New(threshold time.Duration) *Escalator {
	if threshold <= 0 {
		panic("escalation: threshold must be positive")
	}
	return &Escalator{
		threshold: threshold,
		downSince: make(map[string]time.Time),
	}
}

// Evaluate updates internal state for the given alert and returns true when
// the target has been continuously down for at least the configured threshold.
// Recovering targets are cleared from the tracker.
func (e *Escalator) Evaluate(a alert.Alert) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	if a.Severity != alert.SeverityCritical {
		// Target is recovering or informational — clear any tracked downtime.
		delete(e.downSince, a.Target)
		return false
	}

	if _, ok := e.downSince[a.Target]; !ok {
		e.downSince[a.Target] = time.Now()
	}

	return time.Since(e.downSince[a.Target]) >= e.threshold
}

// Reset clears the tracked downtime for a specific target.
func (e *Escalator) Reset(target string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.downSince, target)
}

// DownSince returns the time at which target was first seen as down, and
// whether it is currently being tracked.
func (e *Escalator) DownSince(target string) (time.Time, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	t, ok := e.downSince[target]
	return t, ok
}
