// Package sink provides a fan-out alert dispatcher that forwards alerts
// to multiple named destinations concurrently, collecting any errors.
package sink

import (
	"context"
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// Sender is the interface that wraps the Notify method.
type Sender interface {
	Notify(ctx context.Context, a alert.Alert) error
}

// Sink dispatches an alert to a set of named Senders concurrently.
type Sink struct {
	senders map[string]Sender
}

// New returns a Sink with the provided senders map.
// It panics if senders is nil or empty.
func New(senders map[string]Sender) *Sink {
	if len(senders) == 0 {
		panic("sink: senders must not be empty")
	}
	return &Sink{senders: senders}
}

// SendError holds the name of the sender and the error it returned.
type SendError struct {
	Name string
	Err  error
}

func (e SendError) Error() string {
	return fmt.Sprintf("sink %q: %v", e.Name, e.Err)
}

// Dispatch fans out the alert to all senders concurrently.
// It returns a slice of SendErrors for every sender that failed; nil on
// full success.
func (s *Sink) Dispatch(ctx context.Context, a alert.Alert) []SendError {
	type result struct {
		name string
		err  error
	}

	results := make(chan result, len(s.senders))
	var wg sync.WaitGroup

	for name, sender := range s.senders {
		wg.Add(1)
		go func(n string, sd Sender) {
			defer wg.Done()
			results <- result{name: n, err: sd.Notify(ctx, a)}
		}(name, sender)
	}

	wg.Wait()
	close(results)

	var errs []SendError
	for r := range results {
		if r.err != nil {
			errs = append(errs, SendError{Name: r.name, Err: r.err})
		}
	}
	return errs
}
