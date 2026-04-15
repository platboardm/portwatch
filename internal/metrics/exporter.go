package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Exporter formats a Registry snapshot for human or machine consumption.
type Exporter struct {
	reg *Registry
}

// NewExporter returns an Exporter backed by reg.
// It panics if reg is nil.
func NewExporter(reg *Registry) *Exporter {
	if reg == nil {
		panic("metrics: NewExporter: registry must not be nil")
	}
	return &Exporter{reg: reg}
}

// WriteText writes a human-readable tab-aligned table of metrics to w.
func (e *Exporter) WriteText(w io.Writer) error {
	snap := e.reg.Snapshot()
	keys := sortedKeys(snap)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "METRIC\tVALUE")
	for _, k := range keys {
		fmt.Fprintf(tw, "%s\t%d\n", k, snap[k])
	}
	return tw.Flush()
}

// WriteJSON writes the snapshot as a JSON object to w.
func (e *Exporter) WriteJSON(w io.Writer) error {
	snap := e.reg.Snapshot()
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(snap)
}

func sortedKeys(m map[string]int64) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
