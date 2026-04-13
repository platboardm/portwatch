package circuit

import (
	"sync"
	"time"
)

// Store manages a collection of named Breakers, creating them on first access.
type Store struct {
	mu          sync.Mutex
	threshold   int
	resetWindow time.Duration
	breakers    map[string]*Breaker
}

// NewStore returns a Store that provisions Breakers with the given parameters.
// Panics with the same rules as New.
func NewStore(threshold int, resetWindow time.Duration) *Store {
	// validate early so callers get a clear panic site
	_ = New(threshold, resetWindow)
	return &Store{
		threshold:   threshold,
		resetWindow: resetWindow,
		breakers:    make(map[string]*Breaker),
	}
}

// Get returns the Breaker for the given key, creating one if it does not exist.
func (s *Store) Get(key string) *Breaker {
	s.mu.Lock()
	defer s.mu.Unlock()
	if b, ok := s.breakers[key]; ok {
		return b
	}
	b := New(s.threshold, s.resetWindow)
	s.breakers[key] = b
	return b
}

// Keys returns the names of all tracked breakers.
func (s *Store) Keys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.breakers))
	for k := range s.breakers {
		keys = append(keys, k)
	}
	return keys
}
