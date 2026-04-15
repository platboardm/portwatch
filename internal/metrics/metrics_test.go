package metrics_test

import (
	"sync"
	"testing"

	"portwatch/internal/metrics"
)

func TestCounter_StartsAtZero(t *testing.T) {
	r := metrics.New()
	if got := r.Counter("c").Value(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestCounter_Inc(t *testing.T) {
	r := metrics.New()
	c := r.Counter("checks")
	c.Inc()
	c.Inc()
	if got := c.Value(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestCounter_Add(t *testing.T) {
	r := metrics.New()
	c := r.Counter("alerts")
	c.Add(5)
	if got := c.Value(); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestCounter_SameNameReturnsSameInstance(t *testing.T) {
	r := metrics.New()
	a := r.Counter("x")
	b := r.Counter("x")
	a.Inc()
	if b.Value() != 1 {
		t.Fatal("expected same counter instance")
	}
}

func TestGauge_SetAndGet(t *testing.T) {
	r := metrics.New()
	g := r.Gauge("targets_up")
	g.Set(7)
	if got := g.Value(); got != 7 {
		t.Fatalf("expected 7, got %d", got)
	}
}

func TestGauge_IncDec(t *testing.T) {
	r := metrics.New()
	g := r.Gauge("online")
	g.Inc()
	g.Inc()
	g.Dec()
	if got := g.Value(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestRegistry_Snapshot(t *testing.T) {
	r := metrics.New()
	r.Counter("checks").Add(3)
	r.Gauge("up").Set(2)

	snap := r.Snapshot()
	if snap["checks"] != 3 {
		t.Errorf("checks: expected 3, got %d", snap["checks"])
	}
	if snap["up"] != 2 {
		t.Errorf("up: expected 2, got %d", snap["up"])
	}
}

func TestCounter_ConcurrentInc(t *testing.T) {
	r := metrics.New()
	c := r.Counter("concurrent")
	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()
	if got := c.Value(); got != goroutines {
		t.Fatalf("expected %d, got %d", goroutines, got)
	}
}
