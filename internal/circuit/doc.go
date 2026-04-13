// Package circuit provides a lightweight per-target circuit breaker for
// portwatch.
//
// A Breaker tracks consecutive check failures for a single target. Once the
// failure count reaches the configured threshold the breaker opens and
// subsequent Allow calls return false, preventing redundant alert noise while
// a target is known to be down.
//
// After a configurable reset window the breaker transitions to half-open,
// allowing a single probe attempt. A successful probe closes the breaker and
// resets the failure counter; another failure re-opens it.
//
// Use Store to manage a collection of Breakers keyed by target name without
// requiring callers to handle synchronisation themselves.
//
// Example:
//
//	store := circuit.NewStore(3, 30*time.Second)
//	b := store.Get("api:8080")
//	if b.Allow() {
//	    if err := check(); err != nil {
//	        b.RecordFailure()
//	    } else {
//	        b.RecordSuccess()
//	    }
//	}
package circuit
