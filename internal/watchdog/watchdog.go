// Package watchdog provides a self-monitoring component that detects when
// the main check loop has stalled and emits a warning alert.
package watchdog

import (
	"context"
	"sync"
	"time"
)

// Pinger is the interface accepted by the Watchdog to receive heartbeat pings.
type Pinger interface {
	Ping()
}

// Handler is called when the watchdog detects a stall.
type Handler func(stalledFor time.Duration)

// Watchdog monitors a heartbeat channel and fires a Handler when no ping
// has been received within the configured timeout window.
type Watchdog struct {
	timeout time.Duration
	handler Handler
	mu      sync.Mutex
	lastPin time.Time
}

// New creates a Watchdog that calls handler when no Ping has been received
// within timeout. Panics if timeout is zero or handler is nil.
func New(timeout time.Duration, handler Handler) *Watchdog {
	if timeout <= 0 {
		panic("watchdog: timeout must be positive")
	}
	if handler == nil {
		panic("watchdog: handler must not be nil")
	}
	return &Watchdog{
		timeout: timeout,
		handler: handler,
		lastPin: time.Now(),
	}
}

// Ping records the current time as the most recent heartbeat.
func (w *Watchdog) Ping() {
	w.mu.Lock()
	w.lastPin = time.Now()
	w.mu.Unlock()
}

// Run starts the watchdog loop. It blocks until ctx is cancelled.
func (w *Watchdog) Run(ctx context.Context) {
	ticker := time.NewTicker(w.timeout / 2)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.mu.Lock()
			since := time.Since(w..mu.Unlock()
}
		}
	}
}
