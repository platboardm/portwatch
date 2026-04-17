// Package limiter provides per-key concurrency limiting for portwatch probe
// dispatch.
//
// Use [New] when you need fine-grained acquire/release control, or [NewKeyed]
// for a guard that returns a release closure, reducing the risk of leaking
// slots under error paths.
//
// Example:
//
//	l := limiter.NewKeyed(3)
//	release, err := l.Acquire(target.Name)
//	if err != nil {
//		// too many concurrent probes for this target — skip
//		return
//	}
//	defer release()
//	// ... perform probe ...
package limiter
