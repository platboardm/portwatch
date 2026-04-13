// Package digest provides a lightweight summary mechanism for portwatch.
//
// A Digest collects [Entry] values — each representing a notable state
// change observed during a monitoring interval — and renders them as a
// human-readable report via [Digest.Write].
//
// Typical usage is to wire a Digest into a scheduled job that fires at a
// fixed cadence (e.g. every hour), calls Write to emit the summary, then
// calls Reset to clear accumulated entries before the next window.
//
// Example:
//
//	d := digest.New("hourly summary")
//	// ... accumulate entries from the pipeline ...
//	d.Add(digest.Entry{Target: "db:5432", Status: "down", Severity: alert.SeverityCritical, ChangedAt: time.Now()})
//	_ = d.Write(os.Stdout)
//	d.Reset()
package digest
