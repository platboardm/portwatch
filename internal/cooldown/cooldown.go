// Package cooldown provides a per-key cooldown tracker that prevents
// repeated actions within a configurable quiet period. Unlike ratelimit,
// which enforces a fixed budget, cooldown resets the timer on each
// triggering event — useful for "don't re-alert until N minutes of silence".
package cooldown

import (
	"sync"
	"time"
)

// Tracker tracks the last trigger time per key and reports whether
// enough time has elapsed since the last trigger.
type Tracker struct {
	mu      sync.Mutex
	last    map[string]time.Time
	window  time.Duration
	nowFunc func() time.Time
}

// New returns a Tracker with the given quiet window.
// Panics if window is zero or negative.
func New(window time.Duration) *Tracker {
	if window <= 0 {
		panic("cooldown: window must be positive")
	}
	return &Tracker{
		last:    make(map[string]time.Time),
		window:  window,
		nowFunc: time.Now,
	}
}

// Allow returns true and records the current time if the key has never
// been triggered, or if at least window duration has passed since the
// last trigger. Returns false if the key is still within its quiet period.
func (t *Tracker) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.nowFunc()
	if last, ok := t.last[key]; !ok || now.Sub(last) >= t.window {
		t.last[key] = now
		return true
	}
	return false
}

// Reset clears the cooldown state for the given key, allowing the next
// call to Allow to succeed immediately.
func (t *Tracker) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}

// Remaining returns how much time is left in the cooldown window for key.
// Returns zero if the key is not in cooldown.
func (t *Tracker) Remaining(key string) time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()

	last, ok := t.last[key]
	if !ok {
		return 0
	}
	elapsed := t.nowFunc().Sub(last)
	if elapsed >= t.window {
		return 0
	}
	return t.window - elapsed
}
