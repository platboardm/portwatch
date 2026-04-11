package history_test

import (
	"testing"
	"time"

	"portwatch/internal/history"
)

func makeEntry(target, status string) history.Entry {
	return history.Entry{
		Target:    target,
		Host:      "localhost",
		Port:      8080,
		Status:    status,
		Timestamp: time.Now(),
	}
}

func TestNew_DefaultCapacity(t *testing.T) {
	r := history.New(0)
	if r.Len() != 0 {
		t.Fatalf("expected empty ring, got %d entries", r.Len())
	}
}

func TestAdd_AndLen(t *testing.T) {
	r := history.New(10)
	r.Add(makeEntry("svc-a", "down"))
	r.Add(makeEntry("svc-b", "up"))
	if r.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", r.Len())
	}
}

func TestEntries_Order(t *testing.T) {
	r := history.New(5)
	statuses := []string{"down", "up", "down"}
	for _, s := range statuses {
		r.Add(makeEntry("svc", s))
	}
	entries := r.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for i, e := range entries {
		if e.Status != statuses[i] {
			t.Errorf("entry[%d]: want status %q, got %q", i, statuses[i], e.Status)
		}
	}
}

func TestRing_Wraps(t *testing.T) {
	cap := 3
	r := history.New(cap)
	for i := 0; i < 5; i++ {
		status := "up"
		if i%2 == 0 {
			status = "down"
		}
		r.Add(makeEntry("svc", status))
	}
	if r.Len() != cap {
		t.Fatalf("expected ring capped at %d, got %d", cap, r.Len())
	}
	entries := r.Entries()
	// last 3 added: indices 2(down),3(up),4(down)
	expected := []string{"down", "up", "down"}
	for i, e := range entries {
		if e.Status != expected[i] {
			t.Errorf("entry[%d]: want %q, got %q", i, expected[i], e.Status)
		}
	}
}

func TestEntries_Empty(t *testing.T) {
	r := history.New(10)
	if got := r.Entries(); got != nil {
		t.Fatalf("expected nil from empty ring, got %v", got)
	}
}
