// Package digest provides periodic summary reports of monitoring activity,
// aggregating recent alerts into a human-readable digest suitable for
// scheduled delivery (e.g. hourly or daily summaries).
package digest

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Entry holds a single line in the digest.
type Entry struct {
	Target    string
	Status    string
	Severity  alert.Severity
	ChangedAt time.Time
}

// Digest collects entries and renders a summary.
type Digest struct {
	title   string
	entries []Entry
	now     func() time.Time
}

// New returns a Digest with the given title.
// If title is empty, a default is used.
func New(title string) *Digest {
	if title == "" {
		title = "portwatch digest"
	}
	return &Digest{
		title: title,
		now:   time.Now,
	}
}

// Add appends an entry to the digest.
func (d *Digest) Add(e Entry) {
	d.entries = append(d.entries, e)
}

// Len returns the number of entries collected.
func (d *Digest) Len() int { return len(d.entries) }

// Reset clears all collected entries.
func (d *Digest) Reset() { d.entries = d.entries[:0] }

// Write renders the digest as plain text to w.
// Returns an error if writing fails.
func (d *Digest) Write(w io.Writer) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "=== %s — %s ===\n", d.title, d.now().UTC().Format(time.RFC1123))
	if len(d.entries) == 0 {
		sb.WriteString("  no events recorded.\n")
	} else {
		for _, e := range d.entries {
			fmt.Fprintf(&sb, "  [%s] %s → %s (at %s)\n",
				e.Severity,
				e.Target,
				e.Status,
				e.ChangedAt.UTC().Format("15:04:05 UTC"),
			)
		}
	}
	fmt.Fprintf(&sb, "  total events: %d\n", len(d.entries))
	_, err := io.WriteString(w, sb.String())
	return err
}
