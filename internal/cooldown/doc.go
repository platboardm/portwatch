// Package cooldown provides a per-key cooldown tracker for portwatch.
//
// A Tracker enforces a quiet window between repeated actions for the same
// key. After a trigger is recorded via Allow, subsequent calls for the same
// key return false until the configured window has fully elapsed. This
// differs from ratelimit (fixed budget per window) and suppress (fixed
// window from first event): cooldown restarts the timer on every trigger,
// making it well-suited for "re-alert only after N minutes of silence"
// semantics.
//
// # Usage
//
//	tr := cooldown.New(5 * time.Minute)
//
//	if tr.Allow(target.Name) {
//	    // send alert — enough silence has passed
//	}
//
// Reset clears a key's state so the next Allow succeeds immediately,
// useful when a service recovers and you want to reset the alert cadence.
//
// Remaining reports how much quiet time is left for a given key.
package cooldown
