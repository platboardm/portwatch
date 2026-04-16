// Package heartbeat provides a periodic liveness signal that other
// components can subscribe to, confirming the monitor loop is still running.
package heartbeat

import (
	"context"
	"sync"
	"time"
)

// Beat is a single heartbeat signal carrying the time it was emitted.
type Beat struct {
	At time.Time
}

// Emitter broadcasts periodic heartbeats to registered listeners.
type Emitter struct {
	interval time.Duration
	mu       sync.Mutex
	subs     []chan Beat
}

// New creates an Emitter that fires every interval.
// Panics if interval is zero or negative.
func New(interval time.Duration) *Emitter {
	if interval <= 0 {
		panic("heartbeat: interval must be positive")
	}
	return &Emitter{interval: interval}
}

// Subscribe returns a channel that receives each Beat.
// The channel is buffered with capacity 1; slow consumers drop beats.
func (e *Emitter) Subscribe() <-chan Beat {
	ch := make(chan Beat, 1)
	e.mu.Lock()
	e.subs = append(e.subs, ch)
	e.mu.Unlock()
	return ch
}

// Run starts emitting heartbeats until ctx is cancelled.
func (e *Emitter) Run(ctx context.Context) {
	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			e.mu.Lock()
			for _, ch := range e.subs {
				close(ch)
			}
			e.subs = nil
			e.mu.Unlock()
			return
		case t := <-ticker.C:
			b := Beat{At: t}
			e.mu.Lock()
			for _, ch := range e.subs {
				select {
				case ch <- b:
				default:
				}
			}
			e.mu.Unlock()
		}
	}
}
