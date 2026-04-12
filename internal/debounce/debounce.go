// Package debounce provides a debouncer that delays state-change notifications
// until a condition has been stable for a configurable number of consecutive
// check cycles. This prevents alert storms caused by transient failures.
package debounce

import "sync"

// Debouncer tracks consecutive confirmations before considering a state change
// stable enough to act on.
type Debouncer struct {
	mu        sync.Mutex
	threshold int
	counts    map[string]int
	pending   map[string]bool
}

// New creates a Debouncer that requires at least threshold consecutive
// confirmations of a new state before Confirm returns true.
// It panics if threshold is less than 1.
func New(threshold int) *Debouncer {
	if threshold < 1 {
		panic("debounce: threshold must be >= 1")
	}
	return &Debouncer{
		threshold: threshold,
		counts:    make(map[string]int),
		pending:   make(map[string]bool),
	}
}

// Confirm records that key has been observed in state changed (true = changed
// from baseline). It returns true only when the key has been continuously
// observed as changed for threshold consecutive calls.
// Once confirmed, the internal counter resets so the next transition must
// again accumulate threshold confirmations.
func (d *Debouncer) Confirm(key string, changed bool) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !changed {
		// State returned to normal; reset any pending debounce.
		delete(d.counts, key)
		delete(d.pending, key)
		return false
	}

	d.counts[key]++
	if d.counts[key] >= d.threshold {
		delete(d.counts, key)
		return true
	}
	return false
}

// Reset clears all state for key, as if it had never been seen.
func (d *Debouncer) Reset(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.counts, key)
	delete(d.pending, key)
}

// Pending returns the current consecutive-change count for key.
func (d *Debouncer) Pending(key string) int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.counts[key]
}
