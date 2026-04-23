package fanout_test

import (
	"errors"
	"sync/atomic"
	"testing"

	"portwatch/internal/alert"
	"portwatch/internal/fanout"
)

// stubSender records calls and optionally returns an error.
type stubSender struct {
	calls atomic.Int32
	err   error
}

func (s *stubSender) Send(_ alert.Alert) error {
	s.calls.Add(1)
	return s.err
}

func makeAlert() alert.Alert {
	return alert.New("web", "localhost", 8080, alert.SeverityWarning)
}

func TestNew_PanicsOnEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for empty senders")
		}
	}()
	fanout.New()
}

func TestNew_PanicsOnNilSender(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil sender")
		}
	}()
	fanout.New((*stubSender)(nil))
}

func TestSend_ReachesAllSenders(t *testing.T) {
	a, b := &stubSender{}, &stubSender{}
	f := fanout.New(a, b)

	if err := f.Send(makeAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.calls.Load() != 1 {
		t.Errorf("sender a: want 1 call, got %d", a.calls.Load())
	}
	if b.calls.Load() != 1 {
		t.Errorf("sender b: want 1 call, got %d", b.calls.Load())
	}
}

func TestSend_CollectsErrors(t *testing.T) {
	errA := errors.New("sink A failed")
	errB := errors.New("sink B failed")
	a := &stubSender{err: errA}
	b := &stubSender{err: errB}
	f := fanout.New(a, b)

	err := f.Send(makeAlert())
	if err == nil {
		t.Fatal("expected combined error, got nil")
	}
	if !errors.Is(err, errA) {
		t.Errorf("expected errA in combined error")
	}
	if !errors.Is(err, errB) {
		t.Errorf("expected errB in combined error")
	}
}

func TestSend_PartialFailureDeliversBoth(t *testing.T) {
	good := &stubSender{}
	bad := &stubSender{err: errors.New("oops")}
	f := fanout.New(good, bad)

	_ = f.Send(makeAlert())

	if good.calls.Load() != 1 {
		t.Errorf("good sender should still be called, got %d calls", good.calls.Load())
	}
}

func TestLen_ReturnsCount(t *testing.T) {
	f := fanout.New(&stubSender{}, &stubSender{}, &stubSender{})
	if f.Len() != 3 {
		t.Errorf("want Len 3, got %d", f.Len())
	}
}
