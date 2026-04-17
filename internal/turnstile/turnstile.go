// Package turnstile provides a concurrency gate that limits how many
// goroutines may pass through a critical section simultaneously.
package turnstile

import (
	"context"
	"fmt"
)

// Turnstile is a counting semaphore that controls concurrent access.
type Turnstile struct {
	slots chan struct{}
}

// New returns a Turnstile that allows at most n concurrent passes.
// It panics if n is less than 1.
func New(n int) *Turnstile {
	if n < 1 {
		panic(fmt.Sprintf("turnstile: capacity must be >= 1, got %d", n))
	}
	slots := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		slots <- struct{}{}
	}
	return &Turnstile{slots: slots}
}

// Acquire blocks until a slot is available or ctx is cancelled.
// Returns an error if the context expires before a slot is obtained.
func (t *Turnstile) Acquire(ctx context.Context) error {
	select {
	case <-t.slots:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release returns a slot to the turnstile. It must be called once for
// every successful Acquire.
func (t *Turnstile) Release() {
	t.slots <- struct{}{}
}

// Cap returns the total capacity of the turnstile.
func (t *Turnstile) Cap() int {
	return cap(t.slots)
}

// Available returns the number of slots currently free.
func (t *Turnstile) Available() int {
	return len(t.slots)
}
