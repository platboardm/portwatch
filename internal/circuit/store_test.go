package circuit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/circuit"
)

func TestNewStore_PanicsOnBadThreshold(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for threshold=0")
		}
	}()
	circuit.NewStore(0, time.Second)
}

func TestStore_GetCreatesBreakerOnce(t *testing.T) {
	s := circuit.NewStore(3, time.Second)
	b1 := s.Get("svc-a")
	b2 := s.Get("svc-a")
	if b1 != b2 {
		t.Error("expected same Breaker instance for the same key")
	}
}

func TestStore_IndependentBreakers(t *testing.T) {
	s := circuit.NewStore(2, time.Second)
	a := s.Get("a")
	b := s.Get("b")
	a.RecordFailure()
	a.RecordFailure()
	if b.State() != circuit.StateClosed {
		t.Errorf("breaker 'b' should be unaffected; got %s", b.State())
	}
}

func TestStore_Keys(t *testing.T) {
	s := circuit.NewStore(1, time.Second)
	s.Get("x")
	s.Get("y")
	keys := s.Keys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}
