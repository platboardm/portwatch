package watchdog_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"portwatch/internal/watchdog"
)

func TestNew_PanicsOnZeroTimeout(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero timeout")
		}
	}()
	watchdog.New(0, func(time.Duration) {})
}

func TestNew_PanicsOnNilHandler(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil handler")
		}
	}()
	watchdog.New(time.Second, nil)
}

func TestWatchdog_FiresWhenStalled(t *testing.T) {
	var fired atomic.Int32
	timeout := 50 * time.Millisecond
	w := watchdog.New(timeout, func(d time.Duration) {
		fired.Add(1)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	go w.Run(ctx)

	// Do NOT ping — watchdog should fire.
	<-ctx.Done()

	if fired.Load() == 0 {
		t.Fatal("expected handler to be called when no ping received")
	}
}

func TestWatchdog_SilentWhenPinged(t *testing.T) {
	var fired atomic.Int32
	timeout := 80 * time.Millisecond
	w := watchdog.New(timeout, func(d time.Duration) {
		fired.Add(1)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go w.Run(ctx)

	// Keep pinging faster than the timeout.
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	deadline := time.After(200 * time.Millisecond)
loop:
	for {
		select {
		case <-ticker.C:
			w.Ping()
		case <-deadline:
			break loop
		}
	}

	if fired.Load() != 0 {
		t.Fatalf("expected no handler calls while pinging, got %d", fired.Load())
	}
}

func TestWatchdog_StopsOnContextCancel(t *testing.T) {
	w := watchdog.New(50*time.Millisecond, func(time.Duration) {})
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		w.Run(ctx)
		close(done)
	}()

	cancel()
	select {
	case <-done:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Run did not return after context cancel")
	}
}
