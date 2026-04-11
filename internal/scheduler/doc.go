// Package scheduler provides a simple interval-based scheduler used by the
// portwatch daemon to drive periodic port checks.
//
// # Overview
//
// A [Scheduler] wraps a [CheckFunc] and calls it on a fixed cadence determined
// by the interval supplied to [New]. The first invocation happens immediately
// so that the daemon does not sit idle for a full interval on start-up.
//
// # Usage
//
//	s := scheduler.New(30*time.Second, func(ctx context.Context) {
//		monitor.RunOnce(ctx)
//	})
//	s.Run(ctx) // blocks until ctx is cancelled
//
// Cancelling the context causes Run to return promptly after the current
// in-flight call to CheckFunc completes.
package scheduler
