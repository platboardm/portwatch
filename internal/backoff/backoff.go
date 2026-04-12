// Package backoff provides exponential back-off helpers used when
// retrying failed checks or webhook deliveries.
package backoff

import (
	"math"
	"time"
)

// Strategy holds the parameters for an exponential back-off.
type Strategy struct {
	// Base is the initial wait duration.
	Base time.Duration
	// Max is the upper bound for any single wait.
	Max time.Duration
	// Multiplier is applied to the previous delay each step (default 2.0).
	Multiplier float64
}

// New returns a Strategy with sensible defaults.
// base must be > 0; if max is 0 it defaults to 30 * base.
func New(base, max time.Duration) Strategy {
	if base <= 0 {
		panic("backoff: base duration must be positive")
	}
	if max <= 0 {
		max = 30 * base
	}
	return Strategy{
		Base:       base,
		Max:        max,
		Multiplier: 2.0,
	}
}

// Delay returns the wait duration for the given attempt number (0-indexed).
// Attempt 0 returns Base; subsequent attempts grow exponentially up to Max.
func (s Strategy) Delay(attempt int) time.Duration {
	if attempt <= 0 {
		return s.base()
	}
	factor := math.Pow(s.multiplier(), float64(attempt))
	d := time.Duration(float64(s.base()) * factor)
	if d > s.Max {
		return s.Max
	}
	return d
}

func (s Strategy) base() time.Duration {
	if s.Base <= 0 {
		return time.Second
	}
	return s.Base
}

func (s Strategy) multiplier() float64 {
	if s.Multiplier <= 0 {
		return 2.0
	}
	return s.Multiplier
}
