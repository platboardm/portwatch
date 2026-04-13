// Package suppress provides a time-window based suppression mechanism that
// silences repeated alerts for the same target within a configurable quiet period.
package suppress

import (
	"sync"
	"time"
)

// Suppressor tracks the last alert time per key and suppresses duplicates
// that arrive within the quiet window.
type Suppressor struct {
	mu      sync.Mutex
	last    map[string]time.Time
	window  time.Duration
	nowFunc func() time.Time
}

// New returns a Suppressor with the given quiet window. It panics if window
// is zero or negative.
func New(window time.Duration) *Suppressor {
	if window <= 0 {
		panic("suppress: window must be positive")
	}
	return &Suppressor{
		last:    make(map[string]time.Time),
		window:  window,
		nowFunc: time.Now,
	}
}

// Allow returns true if the alert for key should be forwarded, i.e. either
// no alert has been seen before or the quiet window has elapsed since the
// last allowed alert. Calling Allow with a passing key records the current
// time as the new baseline.
func (s *Suppressor) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.nowFunc()
	if t, ok := s.last[key]; ok && now.Sub(t) < s.window {
		return false
	}
	s.last[key] = now
	return true
}

// Reset clears the suppression record for key, allowing the next call to
// Allow to pass immediately regardless of the window.
func (s *Suppressor) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.last, key)
}

// Len returns the number of keys currently tracked.
func (s *Suppressor) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.last)
}
