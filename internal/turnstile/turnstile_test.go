package turnstile_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"portwatch/internal/turnstile"
)

func TestNew_PanicsOnZeroCapacity(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	turnstile.New(0)
}

func TestNew_PanicsOnNegativeCapacity(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	turnstile.New(-3)
}

func TestCap_ReturnsConfiguredCapacity(t *testing.T) {
	ts := turnstile.New(4)
	if ts.Cap() != 4 {
		t.Fatalf("expected cap 4, got %d", ts.Cap())
	}
}

func TestAvailable_StartsAtCap(t *testing.T) {
	ts := turnstile.New(3)
	if ts.Available() != 3 {
		t.Fatalf("expected 3 available, got %d", ts.Available())
	}
}

func TestAcquireRelease_UpdatesAvailable(t *testing.T) {
	ts := turnstile.New(2)
	ctx := context.Background()

	if err := ts.Acquire(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ts.Available() != 1 {
		t.Fatalf("expected 1 available, got %d", ts.Available())
	}
	ts.Release()
	if ts.Available() != 2 {
		t.Fatalf("expected 2 available after release, got %d", ts.Available())
	}
}

func TestAcquire_BlocksWhenFull(t *testing.T) {
	ts := turnstile.New(1)
	ctx := context.Background()

	if err := ts.Acquire(ctx); err != nil {
		t.Fatal(err)
	}

	ctx2, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := ts.Acquire(ctx2)
	if err == nil {
		t.Fatal("expected context deadline error")
	}
}

func TestAcquire_ConcurrentLimit(t *testing.T) {
	const cap = 3
	const workers = 9
	ts := turnstile.New(cap)

	var mu sync.Mutex
	peak := 0
	current := 0
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ts.Acquire(context.Background())
			mu.Lock()
			current++
			if current > peak {
				peak = current
			}
			mu.Unlock()
			time.Sleep(10 * time.Millisecond)
			mu.Lock()
			current--
			mu.Unlock()
			ts.Release()
		}()
	}
	wg.Wait()
	if peak > cap {
		t.Fatalf("peak concurrency %d exceeded cap %d", peak, cap)
	}
}
