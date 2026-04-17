// Package limiter provides a concurrency limiter that caps the number of
// simultaneous in-flight operations across a set of named keys.
package limiter

import (
	"errors"
	"sync"
)

// ErrLimitReached is returned when the per-key concurrency cap is exceeded.
var ErrLimitReached = errors.New("limiter: concurrency limit reached")

// Limiter enforces a maximum number of concurrent operations per key.
type Limiter struct {
	mu      sync.Mutex
	max     int
	inflight map[string]int
}

// New returns a Limiter that allows at most max concurrent operations per key.
// It panics if max is less than one.
func New(max int) *Limiter {
	if max < 1 {
		panic("limiter: max must be >= 1")
	}
	return &Limiter{max: max, inflight: make(map[string]int)}
}

// Acquire attempts to acquire a slot for key.
// It returns ErrLimitReached if the cap is already reached for that key.
func (l *Limiter) Acquire(key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.inflight[key] >= l.max {
		return ErrLimitReached
	}
	l.inflight[key]++
	return nil
}

// Release decrements the in-flight count for key.
// It is a no-op if the count is already zero.
func (l *Limiter) Release(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.inflight[key] > 0 {
		l.inflight[key]--
	}
	if l.inflight[key] == 0 {
		delete(l.inflight, key)
	}
}

// Inflight returns the current in-flight count for key.
func (l *Limiter) Inflight(key string) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.inflight[key]
}

// Keys returns all keys that currently have in-flight operations.
func (l *Limiter) Keys() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	keys := make([]string, 0, len(l.inflight))
	for k := range l.inflight {
		keys = append(keys, k)
	}
	return keys
}
