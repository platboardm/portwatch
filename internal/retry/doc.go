// Package retry implements a simple, configurable retry policy with
// exponential back-off for use throughout portwatch.
//
// # Overview
//
// Create a Policy with New, optionally tune its fields, then call Do with
// a context and the operation you want to retry:
//
//	p := retry.New(5)
//	p.BaseDelay = 100 * time.Millisecond
//	err := p.Do(ctx, func() error {
//		return sendWebhook(payload)
//	})
//
// Do returns nil on the first successful invocation. If every attempt
// fails it returns ErrExhausted. If the supplied context is cancelled
// between attempts the context error is returned immediately.
//
// # Back-off
//
// The sleep between attempt i and attempt i+1 is:
//
//	delay = min(BaseDelay * Multiplier^i, MaxDelay)
//
// Defaults: BaseDelay=200ms, MaxDelay=30s, Multiplier=2.0.
package retry
