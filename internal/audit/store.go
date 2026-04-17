package audit

import "sync"

// Store retains the most recent audit events in memory up to a fixed capacity.
type Store struct {
	mu     sync.Mutex
	events []Event
	cap    int
}

// NewStore returns a Store with the given capacity.
// Panics if cap is zero or negative.
func NewStore(cap int) *Store {
	if cap <= 0 {
		panic("audit: store capacity must be positive")
	}
	return &Store{cap: cap, events: make([]Event, 0, cap)}
}

// Record appends e to the store, evicting the oldest entry when full.
func (s *Store) Record(e Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.events) >= s.cap {
		s.events = s.events[1:]
	}
	s.events = append(s.events, e)
}

// Entries returns a snapshot of all stored events, oldest first.
func (s *Store) Entries() []Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Event, len(s.events))
	copy(out, s.events)
	return out
}

// Len returns the number of stored events.
func (s *Store) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.events)
}
