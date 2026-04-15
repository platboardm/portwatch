// Package snapshot captures a point-in-time view of all monitored target
// statuses and exposes it for reporting or diagnostic purposes.
package snapshot

import (
	"sync"
	"time"
)

// Status represents the last known status of a single target.
type Status struct {
	Target  string
	Host    string
	Port    int
	Up      bool
	Since   time.Time
	Checks  int64
}

// Snapshot holds an immutable copy of all target statuses at a moment in time.
type Snapshot struct {
	CapturedAt time.Time
	Statuses   []Status
}

// Store accumulates per-target status updates and can produce a Snapshot.
type Store struct {
	mu       sync.RWMutex
	entries  map[string]Status
}

// New returns an initialised Store.
func New() *Store {
	return &Store{
		entries: make(map[string]Status),
	}
}

// Record updates the stored status for the given target key.
func (s *Store) Record(key string, st Status) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[key] = st
}

// Capture returns an immutable Snapshot of all current statuses.
func (s *Store) Capture() Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	statuses := make([]Status, 0, len(s.entries))
	for _, v := range s.entries {
		statuses = append(statuses, v)
	}
	return Snapshot{
		CapturedAt: time.Now(),
		Statuses:   statuses,
	}
}

// Len returns the number of targets currently tracked.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}
