// Package turnstile implements a counting semaphore (concurrency gate) for
// portwatch's internal components.
//
// A Turnstile limits the number of goroutines that may execute a section of
// code concurrently. It is useful when portwatch fans out probe work across
// many targets and needs to cap the number of simultaneous outbound TCP dials
// to avoid exhausting file descriptors or overwhelming the local network stack.
//
// Basic usage:
//
//	ts := turnstile.New(10) // at most 10 concurrent passes
//
//	if err := ts.Acquire(ctx); err != nil {
//	    return err // context cancelled or timed out
//	}
//	defer ts.Release()
//	// … do limited-concurrency work …
package turnstile
