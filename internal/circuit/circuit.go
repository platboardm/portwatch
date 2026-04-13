// Package circuit implements a simple circuit-breaker that tracks consecutive
// failures for a named target and opens (trips) once a threshold is reached.
// Once open the breaker rejects calls until a configurable reset window elapses.
package circuit

import (
	"fmt"
	"sync"
	"time"
)

// State represents the breaker's current state.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // tripped; calls rejected
	StateHalfOpen              // probe allowed after reset window
)

// String returns a human-readable state label.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return fmt.Sprintf("unknown(%d)", int(s))
	}
}

// Breaker is a per-target circuit breaker.
type Breaker struct {
	mu          sync.Mutex
	threshold   int
	resetWindow time.Duration
	failures    int
	state       State
	openedAt    time.Time
}

// New creates a Breaker that opens after threshold consecutive failures and
// attempts a half-open probe after resetWindow.
// Panics if threshold < 1 or resetWindow <= 0.
func New(threshold int, resetWindow time.Duration) *Breaker {
	if threshold < 1 {
		panic("circuit: threshold must be >= 1")
	}
	if resetWindow <= 0 {
		panic("circuit: resetWindow must be positive")
	}
	return &Breaker{
		threshold:   threshold,
		resetWindow: resetWindow,
	}
}

// Allow reports whether a call should be attempted.
// When open it transitions to half-open once the reset window has elapsed.
func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(b.openedAt) >= b.resetWindow {
			b.state = StateHalfOpen
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}
	return false
}

// RecordSuccess resets the breaker to closed.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure increments the failure counter and opens the breaker when the
// threshold is reached.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.threshold && b.state != StateOpen {
		b.state = StateOpen
		b.openedAt = time.Now()
	}
}

// State returns the current breaker state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
