// Package mute provides a time-bounded mute window that suppresses alerts
// for a named target during scheduled maintenance or known outages.
package mute

import (
	"sync"
	"time"
)

// Window represents a mute period for a single target.
type Window struct {
	Target string
	Until  time.Time
}

// Store holds active mute windows keyed by target name.
type Store struct {
	mu      sync.RWMutex
	windows map[string]time.Time
	now     func() time.Time
}

// New returns a new Store. now is used for time comparisons and may be
// replaced in tests with a deterministic clock.
func New(now func() time.Time) *Store {
	if now == nil {
		now = time.Now
	}
	return &Store{
		windows: make(map[string]time.Time),
		now:     now,
	}
}

// Mute registers a mute window for target until the given time.
// Calling Mute again for the same target replaces the previous window.
func (s *Store) Mute(target string, until time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.windows[target] = until
}

// Unmute removes any active mute window for target.
func (s *Store) Unmute(target string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.windows, target)
}

// IsMuted reports whether target currently has an active mute window.
func (s *Store) IsMuted(target string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	until, ok := s.windows[target]
	if !ok {
		return false
	}
	if s.now().After(until) {
		return false
	}
	return true
}

// Active returns a snapshot of all currently active mute windows.
func (s *Store) Active() []Window {
	s.mu.RLock()
	defer s.mu.RUnlock()
	now := s.now()
	out := make([]Window, 0, len(s.windows))
	for target, until := range s.windows {
		if now.Before(until) || now.Equal(until) {
			out = append(out, Window{Target: target, Until: until})
		}
	}
	return out
}
