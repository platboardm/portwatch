// Package summary provides a periodic status summary writer that
// formats and emits a snapshot of all monitored targets and their
// current up/down state to an io.Writer.
package summary

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// TargetStatus holds the name, address, and current status of a target.
type TargetStatus struct {
	Name    string
	Addr    string
	Up      bool
	Since   time.Time
}

// Snapshot is a collection of target statuses at a point in time.
type Snapshot struct {
	At      time.Time
	Targets []TargetStatus
}

// Writer formats and writes a Snapshot to an io.Writer.
type Writer struct {
	out   io.Writer
	title string
}

// New returns a Writer that emits summaries to w.
// title is used as the header line; if empty it defaults to "Port Status Summary".
func New(w io.Writer, title string) *Writer {
	if w == nil {
		panic("summary: writer must not be nil")
	}
	if title == "" {
		title = "Port Status Summary"
	}
	return &Writer{out: w, title: title}
}

// Write formats snap and writes it to the underlying writer.
// It returns the first error encountered, if any.
func (sw *Writer) Write(snap Snapshot) error {
	var sb strings.Builder

	fmt.Fprintf(&sb, "=== %s — %s ===\n", sw.title, snap.At.UTC().Format(time.RFC3339))

	if len(snap.Targets) == 0 {
		sb.WriteString("  (no targets)\n")
	} else {
		for _, t := range snap.Targets {
			status := "UP  "
			if !t.Up {
				status = "DOWN"
			}
			duration := snap.At.Sub(t.Since).Truncate(time.Second)
			fmt.Fprintf(&sb, "  [%s] %-20s %s  (for %s)\n", status, t.Name, t.Addr, duration)
		}
	}

	sb.WriteString(strings.Repeat("-", 60) + "\n")

	_, err := fmt.Fprint(sw.out, sb.String())
	return err
}
