// Package snapshot provides a thread-safe store for capturing point-in-time
// views of monitored target statuses.
//
// # Overview
//
// A [Store] accumulates [Status] records keyed by target name. Callers invoke
// [Store.Record] each time a check completes, then call [Store.Capture] to
// obtain an immutable [Snapshot] suitable for reporting or diagnostics.
//
// # Usage
//
//	s := snapshot.New()
//	s.Record("api", snapshot.Status{Target: "api", Up: true, ...})
//	snap := s.Capture()
//
// The [Writer] type renders a Snapshot as a human-readable text table:
//
//	w := snapshot.NewWriter(os.Stdout, "")
//	w.Write(snap)
//
// Snapshots are safe to pass across goroutines; the underlying slice is
// copied at capture time and the Store uses a read/write mutex internally.
package snapshot
