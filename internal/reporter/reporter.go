// Package reporter formats and writes periodic status summaries
// of all monitored targets to an io.Writer.
package reporter

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/user/portwatch/internal/aggregator"
)

// Reporter writes a human-readable status table to an io.Writer.
type Reporter struct {
	agg    *aggregator.Aggregator
	out    io.Writer
	title  string
}

// New creates a Reporter that reads from agg and writes to out.
// title is printed as a header in each report.
// Panics if agg or out are nil.
func New(agg *aggregator.Aggregator, out io.Writer, title string) *Reporter {
	if agg == nil {
		panic("reporter: aggregator must not be nil")
	}
	if out == nil {
		panic("reporter: writer must not be nil")
	}
	if title == "" {
		title = "portwatch status"
	}
	return &Reporter{agg: agg, out: out, title: title}
}

// Write formats the current aggregator summary and writes it to the
// configured writer. It returns any write error.
func (r *Reporter) Write() error {
	sum := r.agg.Summarise()
	var sb strings.Builder

	timestamp := time.Now().UTC().Format(time.RFC3339)
	fmt.Fprintf(&sb, "=== %s  [%s] ===\n", r.title, timestamp)
	fmt.Fprintf(&sb, "  total=%-4d  up=%-4d  down=%-4d\n",
		sum.Total, sum.Up, sum.Down)

	for _, e := range sum.Entries {
		statusIcon := "✓"
		if e.Status == "down" {
			statusIcon = "✗"
		}
		fmt.Fprintf(&sb, "  [%s] %-30s  %s\n", statusIcon, e.Target, e.Status)
	}

	_, err := fmt.Fprint(r.out, sb.String())
	return err
}
