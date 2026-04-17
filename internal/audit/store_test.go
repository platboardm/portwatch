package audit_test

import (
	"testing"
	"time"

	"portwatch/internal/audit"
)

func makeEvent(kind audit.Kind, target string) audit.Event {
	return audit.Event{At: time.Now().UTC(), Kind: kind, Target: target, Message: "test"}
}

func TestNewStore_PanicsOnZeroCap(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	audit.NewStore(0)
}

func TestStore_RecordAndLen(t *testing.T) {
	s := audit.NewStore(10)
	s.Record(makeEvent(audit.KindStartup, "portwatch"))
	s.Record(makeEvent(audit.KindShutdown, "portwatch"))
	if s.Len() != 2 {
		t.Fatalf("expected len 2, got %d", s.Len())
	}
}

func TestStore_EvictsOldest(t *testing.T) {
	s := audit.NewStore(3)
	for i := 0; i < 5; i++ {
		s.Record(makeEvent(audit.KindAlertSent, "svc"))
	}
	if s.Len() != 3 {
		t.Fatalf("expected len 3 after eviction, got %d", s.Len())
	}
}

func TestStore_EntriesIsImmutable(t *testing.T) {
	s := audit.NewStore(5)
	s.Record(makeEvent(audit.KindConfigReload, "portwatch"))
	e1 := s.Entries()
	e1[0].Message = "mutated"
	e2 := s.Entries()
	if e2[0].Message == "mutated" {
		t.Fatal("entries slice should be a copy")
	}
}

func TestStore_EntriesOrder(t *testing.T) {
	s := audit.NewStore(10)
	kinds := []audit.Kind{audit.KindStartup, audit.KindStateChange, audit.KindAlertSent}
	for _, k := range kinds {
		s.Record(makeEvent(k, "svc"))
	}
	entries := s.Entries()
	for i, k := range kinds {
		if entries[i].Kind != k {
			t.Errorf("pos %d: want %s, got %s", i, k, entries[i].Kind)
		}
	}
}
