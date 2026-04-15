package snapshot

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

// Writer renders a Snapshot as a human-readable text table.
type Writer struct {
	out   io.Writer
	title string
}

// NewWriter returns a Writer that writes to out.
// If title is empty a default heading is used.
func NewWriter(out io.Writer, title string) *Writer {
	if out == nil {
		panic("snapshot: writer must not be nil")
	}
	if title == "" {
		title = "Port Status Snapshot"
	}
	return &Writer{out: out, title: title}
}

// Write renders snap to the underlying writer.
func (w *Writer) Write(snap Snapshot) error {
	sorted := make([]Status, len(snap.Statuses))
	copy(sorted, snap.Statuses)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Target < sorted[j].Target
	})

	fmt.Fprintf(w.out, "=== %s (%s) ===\n",
		w.title, snap.CapturedAt.Format(time.RFC3339))

	if len(sorted) == 0 {
		fmt.Fprintln(w.out, "  (no targets)")
		return nil
	}

	fmt.Fprintf(w.out, "  %-24s %-6s %s\n", "TARGET", "STATUS", "SINCE")
	fmt.Fprintln(w.out, "  "+strings.Repeat("-", 52))
	for _, st := range sorted {
		label := "DOWN"
		if st.Up {
			label = "UP"
		}
		fmt.Fprintf(w.out, "  %-24s %-6s %s\n",
			st.Target, label, st.Since.Format(time.RFC3339))
	}
	return nil
}
