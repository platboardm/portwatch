// Package ratelimit provides a simple token-bucket rate limiter to suppress
// repeated alerts for the same target within a configurable cooldown window.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks the last alert time per target key and suppresses
// subsequent alerts that arrive within the cooldown duration.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
	now      func() time.Time // injectable for testing
}

// New creates a Limiter with the given cooldown duration.
// Panics if cooldown is zero or negative.
func New(cooldown time.Duration) *Limiter {
	if cooldown <= 0 {
		panic("ratelimit: cooldown must be positive")
	}
	return &Limiter{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow reports whether an alert for the given key should be allowed
// through. The first call for a key always returns true. Subsequent calls
// return true only after the cooldown window has elapsed.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	if t, ok := l.last[key]; ok && now.Sub(t) < l.cooldown {
		return false
	}
	l.last[key] = now
	return true
}

// Reset clears the recorded time for the given key, allowing the next
// alert to pass through immediately regardless of the cooldown window.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, key)
}

// Cooldown returns the configured cooldown duration.
func (l *Limiter) Cooldown() time.Duration {
	return l.cooldown
}
