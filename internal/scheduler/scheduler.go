// Package scheduler provides a periodic tick-based runner that triggers
// port checks at a configurable interval.
package scheduler

import (
	"context"
	"time"
)

// CheckFunc is the function called on every tick.
type CheckFunc func(ctx context.Context)

// Scheduler runs a CheckFunc at a fixed interval until the context is cancelled.
type Scheduler struct {
	interval time.Duration
	fn       CheckFunc
}

// New returns a Scheduler that will call fn every interval.
// interval must be positive; if it is not, New panics.
func New(interval time.Duration, fn CheckFunc) *Scheduler {
	if interval <= 0 {
		panic("scheduler: interval must be positive")
	}
	if fn == nil {
		panic("scheduler: fn must not be nil")
	}
	return &Scheduler{interval: interval, fn: fn}
}

// Run blocks, calling fn on every tick, until ctx is done.
// It performs an immediate first call before waiting for the first tick so
// that the daemon does not wait a full interval on startup.
func (s *Scheduler) Run(ctx context.Context) {
	s.fn(ctx)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.fn(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// Interval returns the configured tick interval.
func (s *Scheduler) Interval() time.Duration {
	return s.interval
}
