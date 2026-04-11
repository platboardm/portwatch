package history

import (
	"time"

	"portwatch/internal/checker"
)

// Recorder wraps a Ring and provides a convenience method
// for recording checker events.
type Recorder struct {
	ring *Ring
}

// NewRecorder returns a Recorder backed by a Ring of the given capacity.
func NewRecorder(capacity int) *Recorder {
	return &Recorder{ring: New(capacity)}
}

// Record converts a checker.Status into an Entry and stores it.
func (rec *Recorder) Record(target, host string, port int, status checker.Status) {
	e := Entry{
		Target:    target,
		Host:      host,
		Port:      port,
		Status:    status.String(),
		Timestamp: time.Now(),
	}
	rec.ring.Add(e)
}

// Entries delegates to the underlying Ring.
func (rec *Recorder) Entries() []Entry {
	return rec.ring.Entries()
}

// Len delegates to the underlying Ring.
func (rec *Recorder) Len() int {
	return rec.ring.Len()
}
