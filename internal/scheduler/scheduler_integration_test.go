package scheduler_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/portwatch/internal/scheduler"
)

// TestScheduler_TickCount is an integration-style test that verifies the
// scheduler fires approximately the expected number of times over a short
// window. We allow generous bounds to avoid flakiness on slow CI runners.
func TestScheduler_TickCount(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	const (
		interval  = 25 * time.Millisecond
		duration  = 200 * time.Millisecond
		minTicks  = 5  // conservative lower bound
		maxTicks  = 15 // generous upper bound
	)

	var count atomic.Int32

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	s := scheduler.New(interval, func(_ context.Context) {
		count.Add(1)
	})

	s.Run(ctx) // blocks until timeout

	n := int(count.Load())
	if n < minTicks || n > maxTicks {
		t.Fatalf("expected between %d and %d ticks, got %d", minTicks, maxTicks, n)
	}
}
