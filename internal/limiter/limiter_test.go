package limiter_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/limiter"
)

func TestNew_PanicsOnZeroMax(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	limiter.New(0)
}

func TestAcquire_FirstCallSucceeds(t *testing.T) {
	l := limiter.New(2)
	if err := l.Acquire("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAcquire_RespectsMax(t *testing.T) {
	l := limiter.New(2)
	_ = l.Acquire("svc")
	_ = l.Acquire("svc")
	if err := l.Acquire("svc"); err != limiter.ErrLimitReached {
		t.Fatalf("expected ErrLimitReached, got %v", err)
	}
}

func TestRelease_FreesSlot(t *testing.T) {
	l := limiter.New(1)
	_ = l.Acquire("svc")
	l.Release("svc")
	if err := l.Acquire("svc"); err != nil {
		t.Fatalf("unexpected error after release: %v", err)
	}
}

func TestRelease_NoopOnZero(t *testing.T) {
	l := limiter.New(2)
	l.Release("svc") // should not panic
}

func TestInflight_TracksCount(t *testing.T) {
	l := limiter.New(5)
	_ = l.Acquire("a")
	_ = l.Acquire("a")
	if got := l.Inflight("a"); got != 2 {
		t.Fatalf("want 2, got %d", got)
	}
}

func TestKeys_IndependentPerKey(t *testing.T) {
	l := limiter.New(3)
	_ = l.Acquire("x")
	_ = l.Acquire("y")
	keys := l.Keys()
	if len(keys) != 2 {
		t.Fatalf("want 2 keys, got %d", len(keys))
	}
}

func TestAcquire_Concurrent(t *testing.T) {
	l := limiter.New(10)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = l.Acquire("shared")
		}()
	}
	wg.Wait()
	if got := l.Inflight("shared"); got != 10 {
		t.Fatalf("want 10, got %d", got)
	}
}
