// Package fanout distributes a single alert to multiple sinks concurrently.
// Each sink receives the alert independently; a failure in one sink does not
// prevent delivery to others. Errors are collected and returned as a combined
// error after all sinks have been attempted.
package fanout

import (
	"errors"
	"fmt"
	"sync"

	"portwatch/internal/alert"
)

// Sender is any type that can receive an alert.
type Sender interface {
	Send(a alert.Alert) error
}

// Fanout dispatches an alert to a fixed set of senders concurrently.
type Fanout struct {
	senders []Sender
}

// New returns a Fanout that dispatches to the provided senders.
// It panics if senders is empty or any element is nil.
func New(senders ...Sender) *Fanout {
	if len(senders) == 0 {
		panic("fanout: at least one sender is required")
	}
	for i, s := range senders {
		if s == nil {
			panic(fmt.Sprintf("fanout: sender at index %d is nil", i))
		}
	}
	return &Fanout{senders: senders}
}

// Send dispatches a to every sender concurrently and waits for all to finish.
// If one or more senders return an error the errors are joined and returned.
func (f *Fanout) Send(a alert.Alert) error {
	var (
		mu   sync.Mutex
		errs []error
		wg   sync.WaitGroup
	)

	for _, s := range f.senders {
		s := s
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.Send(a); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return errors.Join(errs...)
}

// Len returns the number of senders registered with this Fanout.
func (f *Fanout) Len() int { return len(f.senders) }
