// Package sampler provides probabilistic sampling for alert events,
// allowing a fraction of alerts to pass through based on a configured rate.
package sampler

import (
	"fmt"
	"math/rand"
	"sync"
)

// Source is the interface for a random number source.
type Source interface {
	Float64() float64
}

// defaultSource wraps the global math/rand functions.
type defaultSource struct{}

func (d *defaultSource) Float64() float64 { return rand.Float64() }

// Sampler decides probabilistically whether an event should be allowed through.
type Sampler struct {
	mu   sync.Mutex
	rate float64
	src  Source
}

// New creates a Sampler with the given sample rate (0.0–1.0).
// A rate of 1.0 allows all events; 0.0 suppresses all events.
// Panics if rate is outside [0.0, 1.0].
func New(rate float64) *Sampler {
	return NewWithSource(rate, &defaultSource{})
}

// NewWithSource creates a Sampler with a custom random source.
func NewWithSource(rate float64, src Source) *Sampler {
	if rate < 0.0 || rate > 1.0 {
		panic(fmt.Sprintf("sampler: rate must be in [0.0, 1.0], got %v", rate))
	}
	if src == nil {
		panic("sampler: source must not be nil")
	}
	return &Sampler{rate: rate, src: src}
}

// Allow returns true if the event should be allowed through based on the sample rate.
func (s *Sampler) Allow() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.rate == 1.0 {
		return true
	}
	if s.rate == 0.0 {
		return false
	}
	return s.src.Float64() < s.rate
}

// Rate returns the configured sample rate.
func (s *Sampler) Rate() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rate
}
