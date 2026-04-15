// Package window provides a sliding-window counter used to track event
// frequency over a rolling time period. It is safe for concurrent use.
package window

import (
	"sync"
	"time"
)

// Counter counts events that occurred within a sliding time window.
type Counter struct {
	mu       sync.Mutex
	window   time.Duration
	timestamps []time.Time
}

// New returns a Counter with the given sliding window duration.
// It panics if window is zero or negative.
func New(window time.Duration) *Counter {
	if window <= 0 {
		panic("window: duration must be positive")
	}
	return &Counter{window: window}
}

// Add records a new event at the current time.
func (c *Counter) Add() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.timestamps = append(c.timestamps, time.Now())
	c.evict(time.Now())
}

// Count returns the number of events recorded within the current window.
func (c *Counter) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evict(time.Now())
	return len(c.timestamps)
}

// Reset clears all recorded events.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.timestamps = c.timestamps[:0]
}

// evict removes timestamps that have fallen outside the sliding window.
// Caller must hold c.mu.
func (c *Counter) evict(now time.Time) {
	cutoff := now.Add(-c.window)
	i := 0
	for i < len(c.timestamps) && c.timestamps[i].Before(cutoff) {
		i++
	}
	c.timestamps = c.timestamps[i:]
}
