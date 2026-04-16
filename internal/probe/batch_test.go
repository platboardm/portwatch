package probe

import (
	"testing"
	"time"
)

func TestBatch_AllUp(t *testing.T) {
	addrs := make([]string, 3)
	for i := range addrs {
		addrs[i] = startTCP(t)
	}
	targets := []Target{
		{Name: "a", Addr: addrs[0]},
		{Name: "b", Addr: addrs[1]},
		{Name: "c", Addr: addrs[2]},
	}
	p := New(time.Second)
	results := p.Batch(targets)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Up {
			t.Errorf("target %q expected Up", r.Target)
		}
	}
}

func TestBatch_PreservesOrder(t *testing.T) {
	addr := startTCP(t)
	targets := []Target{
		{Name: "first", Addr: addr},
		{Name: "second", Addr: "127.0.0.1:1"},
	}
	p := New(200 * time.Millisecond)
	results := p.Batch(targets)
	if results[0].Target != "first" {
		t.Errorf("expected first, got %q", results[0].Target)
	}
	if results[1].Target != "second" {
		t.Errorf("expected second, got %q", results[1].Target)
	}
}

func TestBatch_Empty(t *testing.T) {
	p := New(time.Second)
	results := p.Batch(nil)
	if len(results) != 0 {
		t.Fatalf("expected empty results")
	}
}
