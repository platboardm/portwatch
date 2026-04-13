// Package jitter adds randomised offsets to durations to prevent
// thundering-herd problems when many targets are checked simultaneously.
package jitter

import (
	"math/rand"
	"sync"
	"time"
)

// Source is the interface used to obtain random values. It exists so tests
// can inject a deterministic source.
type Source interface {
	Int63n(n int64) int64
}

// Jitter applies a random offset within [0, factor*base] to base.
// factor must be in the range (0, 1]; a factor of 0.25 means the returned
// duration will be between base and base*1.25.
type Jitter struct {
	mu     sync.Mutex
	src    Source
	factor float64
}

// New creates a Jitter with the given factor. It panics if factor is not in
// (0, 1].
func New(factor float64) *Jitter {
	if factor <= 0 || factor > 1 {
		panic("jitter: factor must be in (0, 1]")
	}
	return &Jitter{
		src:    rand.New(rand.NewSource(time.Now().UnixNano())), //nolint:gosec
		factor: factor,
	}
}

// NewWithSource creates a Jitter backed by the provided Source. Useful for
// deterministic tests.
func NewWithSource(factor float64, src Source) *Jitter {
	if factor <= 0 || factor > 1 {
		panic("jitter: factor must be in (0, 1]")
	}
	if src == nil {
		panic("jitter: source must not be nil")
	}
	return &Jitter{src: src, factor: factor}
}

// Apply returns base plus a random offset in [0, factor*base).
func (j *Jitter) Apply(base time.Duration) time.Duration {
	if base <= 0 {
		return base
	}
	max := int64(float64(base) * j.factor)
	if max == 0 {
		return base
	}
	j.mu.Lock()
	offset := j.src.Int63n(max)
	j.mu.Unlock()
	return base + time.Duration(offset)
}
