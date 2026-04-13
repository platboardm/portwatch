// Package backoff implements an exponential back-off strategy for use when
// retrying transient failures such as unreachable ports or failed webhook
// deliveries.
//
// # Basic usage
//
//	s := backoff.New(500*time.Millisecond, 30*time.Second)
//	for attempt := 0; attempt < maxAttempts; attempt++ {
//	    if err := tryOperation(); err == nil {
//	        break
//	    }
//	    time.Sleep(s.Delay(attempt))
//	}
//
// The Delay method returns the wait duration for a given attempt index
// (0-indexed). Attempt 0 always returns the base duration; each subsequent
// attempt multiplies the previous delay by Multiplier (default 2.0) until
// the configured maximum is reached.
//
// # Jitter
//
// To avoid thundering-herd problems when many goroutines retry simultaneously,
// use DelayWithJitter instead of Delay. It adds a random fraction of the
// computed delay so that retries are spread out over time:
//
//	time.Sleep(s.DelayWithJitter(attempt))
package backoff
