// Package dedup provides alert deduplication based on a fingerprint hash.
// Identical alerts seen within the deduplication window are suppressed to
// reduce noise from flapping services.
package dedup

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Deduplicator suppresses repeated alerts that share the same fingerprint
// within a configurable time window.
type Deduplicator struct {
	mu     sync.Mutex
	window time.Duration
	seen   map[string]time.Time
	now    func() time.Time
}

// New returns a Deduplicator with the given deduplication window.
// It panics if window is zero or negative.
func New(window time.Duration) *Deduplicator {
	if window <= 0 {
		panic("dedup: window must be positive")
	}
	return &Deduplicator{
		window: window,
		seen:   make(map[string]time.Time),
		now:    time.Now,
	}
}

// Allow returns true if the alert should be forwarded, or false if it is a
// duplicate seen within the current window. Expired entries are evicted lazily
// on each call to keep memory bounded.
func (d *Deduplicator) Allow(a alert.Alert) bool {
	fp := fingerprint(a)
	now := d.now()

	d.mu.Lock()
	defer d.mu.Unlock()

	d.evict(now)

	if _, exists := d.seen[fp]; exists {
		return false
	}
	d.seen[fp] = now
	return true
}

// Len returns the number of fingerprints currently tracked.
func (d *Deduplicator) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}

// evict removes entries older than the window. Must be called with d.mu held.
func (d *Deduplicator) evict(now time.Time) {
	for fp, ts := range d.seen {
		if now.Sub(ts) >= d.window {
			delete(d.seen, fp)
		}
	}
}

// fingerprint produces a stable hash key from an alert's target, host, port
// and severity so that semantically identical alerts are collapsed.
func fingerprint(a alert.Alert) string {
	raw := fmt.Sprintf("%s|%s|%d|%s", a.Target, a.Host, a.Port, a.Severity)
	sum := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%x", sum)
}
