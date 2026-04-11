package webhook

import (
	"context"
	"fmt"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// RetryNotifier wraps a Notifier and retries transient failures with
// exponential back-off up to maxAttempts times.
type RetryNotifier struct {
	inner       *Notifier
	maxAttempts int
	baseDelay   time.Duration
}

// NewRetry creates a RetryNotifier. maxAttempts must be >= 1.
// baseDelay is the initial back-off; it doubles on each retry.
func NewRetry(n *Notifier, maxAttempts int, baseDelay time.Duration) *RetryNotifier {
	if maxAttempts < 1 {
		panic("webhook: maxAttempts must be >= 1")
	}
	if baseDelay <= 0 {
		baseDelay = 500 * time.Millisecond
	}
	return &RetryNotifier{inner: n, maxAttempts: maxAttempts, baseDelay: baseDelay}
}

// Notify attempts to deliver the alert, retrying on failure.
func (r *RetryNotifier) Notify(ctx context.Context, a alert.Alert) error {
	delay := r.baseDelay
	var lastErr error
	for attempt := 1; attempt <= r.maxAttempts; attempt++ {
		if err := r.inner.Notify(ctx, a); err == nil {
			return nil
		} else {
			lastErr = err
		}
		if attempt == r.maxAttempts {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
		delay *= 2
	}
	return fmt.Errorf("webhook: all %d attempts failed: %w", r.maxAttempts, lastErr)
}
