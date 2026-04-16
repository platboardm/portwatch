// Package history records port check events for reporting and diagnostics.
package history

import (
	"sync"
	"time"
)

// Entry represents a single recorded port event.
type Entry struct {
	Target    string
	Host      string
	Port      int
	Status    string
	Timestamp time.Time
}

// Ring is a fixed-size circular buffer of history entries.
type Ring struct {
	mu      sync.Mutex
	buf     []Entry
	cap     int
	head    int
	count   int
}

// New creates a Ring with the given capacity.
func New(capacity int) *Ring {
	if capacity <= 0 {
		capacity = 100
	}
	return &Ring{
		buf: make([]Entry, capacity),
		cap: capacity,
	}
}

// Add appends an entry to the ring, overwriting the oldest when full.
func (r *Ring) Add(e Entry) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.buf[r.head] = e
	r.head = (r.head + 1) % r.cap
	if r.count < r.cap {
		r.count++
	}
}

// Entries returns all stored entries in chronological order (oldest first).
func (r *Ring) Entries() []Entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.count == 0 {
		return nil
	}
	out := make([]Entry, r.count)
	start := (r.head - r.count + r.cap) % r.cap
	for i := 0; i < r.count; i++ {
		out[i] = r.buf[(start+i)%r.cap]
	}
	return out
}

// Len returns the number of entries currently stored.
func (r *Ring) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.count
}

// Filter returns entries matching the given predicate, in chronological order.
func (r *Ring) Filter(fn func(Entry) bool) []Entry {
	entries := r.Entries()
	out := out[:0:0]
	for _, e := range entries {
		if fn(e) {
			out = append(out, e)
		}
	}
	return out
}
