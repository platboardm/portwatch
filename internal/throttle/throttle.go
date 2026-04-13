// Package throttle provides a token-bucket style throttle that limits
// how many events can pass through in a given time window.
package throttle

import (
	"sync"
	"time"
)

// Throttle limits the number of allowed calls to N per window duration.
type Throttle struct {
	mu      sync.Mutex
	window  time.Duration
	limit   int
	buckets map[string][]time.Time
	now     func() time.Time
}

// New creates a new Throttle that allows at most limit events per window.
// Panics if limit is zero or negative, or if window is zero or negative.
func New(limit int, window time.Duration) *Throttle {
	if limit <= 0 {
		panic("throttle: limit must be positive")
	}
	if window <= 0 {
		panic("throttle: window must be positive")
	}
	return &Throttle{
		window:  window,
		limit:   limit,
		buckets: make(map[string][]time.Time),
		now:     time.Now,
	}
}

// Allow reports whether the event identified by key is allowed through.
// It returns true if the number of calls for key within the current window
// is less than the configured limit.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	cutoff := now.Add(-t.window)

	// Evict timestamps outside the window.
	times := t.buckets[key]
	valid := times[:0]
	for _, ts := range times {
		if ts.After(cutoff) {
			valid = append(valid, ts)
		}
	}

	if len(valid) >= t.limit {
		t.buckets[key] = valid
		return false
	}

	t.buckets[key] = append(valid, now)
	return true
}

// Reset clears all recorded timestamps for the given key.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.buckets, key)
}
