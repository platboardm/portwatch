package retry_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/retry"
)

// TestRetry_EventualSuccess verifies that a transiently failing operation
// ultimately succeeds within a realistic time budget.
func TestRetry_EventualSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	p := retry.New(5)
	p.BaseDelay = 20 * time.Millisecond
	p.MaxDelay = 200 * time.Millisecond
	p.Multiplier = 2.0

	var calls int32
	start := time.Now()
	err := p.Do(context.Background(), func() error {
		if atomic.AddInt32(&calls, 1) < 4 {
			return errors.New("not ready")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
	elapsed := time.Since(start)
	// 20 + 40 + 80 = 140 ms of sleep across 3 back-offs
	if elapsed < 100*time.Millisecond {
		t.Fatalf("completed suspiciously fast: %v", elapsed)
	}
	if calls != 4 {
		t.Fatalf("expected 4 calls, got %d", calls)
	}
}
