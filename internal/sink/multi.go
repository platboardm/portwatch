package sink

import (
	"context"
	"log"

	"github.com/user/portwatch/internal/alert"
)

// Multi wraps a Sink and provides a higher-level Dispatch that logs
// individual send failures but never returns an error to the caller,
// making it safe to use in fire-and-forget pipelines.
type Multi struct {
	sink   *Sink
	logger *log.Logger
}

// NewMulti returns a Multi backed by the given Sink.
// It panics if s or logger is nil.
func NewMulti(s *Sink, logger *log.Logger) *Multi {
	if s == nil {
		panic("sink: Multi requires a non-nil Sink")
	}
	if logger == nil {
		panic("sink: Multi requires a non-nil logger")
	}
	return &Multi{sink: s, logger: logger}
}

// Send dispatches the alert to all senders and logs any failures.
// It always returns nil so that callers can treat it as best-effort.
func (m *Multi) Send(ctx context.Context, a alert.Alert) {
	errs := m.sink.Dispatch(ctx, a)
	for _, e := range errs {
		m.logger.Printf("[sink] dispatch error: %v", e)
	}
}
