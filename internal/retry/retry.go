// Package retry provides a configurable retry policy with exponential
// backoff and jitter for transient failures in outbound calls.
package retry

import (
	"context"
	"errors"
	"time"
)

// Policy holds the configuration for a retry loop.
type Policy struct {
	Attempts int
	BaseDelay time.Duration
	MaxDelay  time.Duration
	Multiplier float64
}

// ErrExhausted is returned when all attempts have been consumed.
var ErrExhausted = errors.New("retry: all attempts exhausted")

// New returns a Policy with sensible defaults.
// Panics if attempts is less than 1.
func New(attempts int) Policy {
	if attempts < 1 {
		panic("retry: attempts must be >= 1")
	}
	return Policy{
		Attempts:   attempts,
		BaseDelay:  200 * time.Millisecond,
		MaxDelay:   30 * time.Second,
		Multiplier: 2.0,
	}
}

// Do executes fn up to p.Attempts times, backing off between failures.
// It returns nil on the first success. If the context is cancelled the
// loop stops and the context error is returned.
func (p Policy) Do(ctx context.Context, fn func() error) error {
	delay := p.BaseDelay
	for i := 0; i < p.Attempts; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := fn(); err == nil {
			return nil
		}
		if i == p.Attempts-1 {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
		delay = time.Duration(float64(delay) * p.Multiplier)
		if delay > p.MaxDelay {
			delay = p.MaxDelay
		}
	}
	return ErrExhausted
}
