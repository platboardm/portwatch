// Package deadletter provides a dead-letter queue for alerts that could not
// be delivered after all retry attempts have been exhausted. It retains the
// most recent N failed alerts so that operators can inspect or replay them.
package deadletter

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const defaultCap = 64

// Entry wraps a failed alert together with metadata about the failure.
type Entry struct {
	Alert     alert.Alert
	Reason    string
	FailedAt  time.Time
	Attempts  int
}

// Queue is a bounded, thread-safe dead-letter store.
type Queue struct {
	mu      sync.Mutex
	entries []Entry
	cap     int
}

// New creates a Queue with the given capacity. It panics if cap is <= 0.
func New(cap int) *Queue {
	if cap <= 0 {
		panic("deadletter: capacity must be greater than zero")
	}
	return &Queue{cap: cap, entries: make([]Entry, 0, cap)}
}

// NewDefault creates a Queue with the default capacity.
func NewDefault() *Queue { return New(defaultCap) }

// Push appends a failed alert to the queue. If the queue is full the oldest
// entry is evicted to make room.
func (q *Queue) Push(a alert.Alert, reason string, attempts int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	e := Entry{Alert: a, Reason: reason, FailedAt: time.Now(), Attempts: attempts}
	if len(q.entries) >= q.cap {
		// evict oldest
		copy(q.entries, q.entries[1:])
		q.entries[len(q.entries)-1] = e
		return
	}
	q.entries = append(q.entries, e)
}

// Drain returns all entries and resets the queue.
func (q *Queue) Drain() []Entry {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]Entry, len(q.entries))
	copy(out, q.entries)
	q.entries = q.entries[:0]
	return out
}

// Len returns the current number of entries in the queue.
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.entries)
}
