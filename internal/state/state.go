// Package state tracks the last known status of monitored ports
// so the monitor can detect transitions (up→down, down→up).
package state

import "sync"

// Status represents whether a port is reachable.
type Status int

const (
	// Unknown is the initial status before any check has completed.
	Unknown Status = iota
	// Up means the port accepted a connection.
	Up
	// Down means the port refused or timed out.
	Down
)

// String returns a human-readable label for the status.
func (s Status) String() string {
	switch s {
	case Up:
		return "up"
	case Down:
		return "down"
	default:
		return "unknown"
	}
}

// Store holds the last known Status for each target key.
// It is safe for concurrent use.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Status
}

// New returns an initialised, empty Store.
func New() *Store {
	return &Store{entries: make(map[string]Status)}
}

// Get returns the last recorded Status for key.
// If the key has never been recorded, Unknown is returned.
func (s *Store) Get(key string) Status {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if st, ok := s.entries[key]; ok {
		return st
	}
	return Unknown
}

// Set records status for key and reports whether the value changed
// from the previously stored value (or from Unknown on first call).
func (s *Store) Set(key string, status Status) (changed bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	prev, ok := s.entries[key]
	if !ok {
		prev = Unknown
	}
	s.entries[key] = status
	return prev != status
}

// Snapshot returns a copy of all stored entries.
func (s *Store) Snapshot() map[string]Status {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]Status, len(s.entries))
	for k, v := range s.entries {
		out[k] = v
	}
	return out
}
