// Package audit provides a structured event logger and in-memory store for
// recording significant portwatch lifecycle events.
//
// # Logger
//
// Logger writes human-readable audit lines to any io.Writer:
//
//	l := audit.New(os.Stderr)
//	l.Log(audit.KindStateChange, "api:8080", "down -> up")
//
// # Store
//
// Store keeps the N most recent events in a ring-like buffer:
//
//	s := audit.NewStore(256)
//	s.Record(event)
//	entries := s.Entries()
package audit
