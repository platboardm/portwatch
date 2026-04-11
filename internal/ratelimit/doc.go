// Package ratelimit implements a per-target cooldown mechanism for alert
// suppression in portwatch.
//
// When a monitored port transitions to a down or up state, the monitor may
// fire multiple alerts in quick succession if the service is flapping.
// Wrapping the notifier with a Limiter ensures that repeated alerts for the
// same target are suppressed until the configured cooldown window expires.
//
// Basic usage:
//
//	limiter := ratelimit.New(5 * time.Minute)
//
//	if limiter.Allow(target.Key()) {
//		notifier.Notify(event)
//	}
//
// Reset can be called to clear the cooldown for a key, for example when
// the target transitions between distinct states (down → up).
package ratelimit
