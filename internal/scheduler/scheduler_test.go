package scheduler_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/portwatch/internal/scheduler"
)

func TestNew_PanicsOnZeroInterval(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero interval")
		}
	}()
	scheduler.New(0, func(_ context.Context) {})
}

func TestNew_PanicsOnNilFunc(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil fn")
		}
	}()
	scheduler.New(time.Second, nil)
}

func TestScheduler_Interval(t *testing.T) {
	s := scheduler.New(42*time.Millisecond, func(_ context.Context) {})
	if s.Interval() != 42*time.Millisecond {
		t.Fatalf("expected 42ms, got %v", s.Interval())
	}
}

// TestScheduler_CallsImmediately verifies that fn is called before the first
// tick so the daemon reacts without waiting a full interval.
func TestScheduler_CallsImmediately(t *testing.T) {
	var count atomic.Int32

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	called := make(chan struct{}, 1)
	s := scheduler.New(10*time.Second, func(_ context.Context) {
		count.Add(1)
		select {
		case called <- struct{}{}:
		default:
		}
	})

	go s.Run(ctx)

	select {
	case <-called:
	// good — fn was called before the 10 s tick
	case <-time.After(500 * time.Millisecond):
		t.Fatal("fn was not called immediately")
	}
}

// TestScheduler_StopsOnCancel verifies that Run returns when the context is
// cancelled and does not call fn after cancellation.
func TestScheduler_StopsOnCancel(t *testing.T) {
	var count atomic.Int32

	ctx, cancel := context.WithCancel(context.Background())

	s := scheduler.New(20*time.Millisecond, func(_ context.Context) {
		count.Add(1)
	})

	done := make(chan struct{})
	go func() {
		s.Run(ctx)
		close(done)
	}()

	// Let it tick a couple of times.
	time.Sleep(60 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Run did not return after context cancellation")
	}

	snap := count.Load()
	time.Sleep(40 * time.Millisecond)
	if after := count.Load(); after != snap {
		t.Fatalf("fn called %d times after cancel (was %d before)", after-snap, snap)
	}
}
