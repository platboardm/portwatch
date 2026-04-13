package sink_test

import (
	"context"
	"errors"
	"sort"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/sink"
)

// stubSender is a Sender that returns a preset error.
type stubSender struct{ err error }

func (s *stubSender) Notify(_ context.Context, _ alert.Alert) error { return s.err }

func makeAlert() alert.Alert {
	a, _ := alert.New("web", alert.SeverityWarning, alert.StatusDown)
	return a
}

func TestNew_PanicsOnEmptySenders(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	sink.New(map[string]sink.Sender{})
}

func TestDispatch_AllSucceed(t *testing.T) {
	s := sink.New(map[string]sink.Sender{
		"a": &stubSender{},
		"b": &stubSender{},
	})
	errs := s.Dispatch(context.Background(), makeAlert())
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestDispatch_PartialFailure(t *testing.T) {
	s := sink.New(map[string]sink.Sender{
		"ok":  &stubSender{},
		"bad": &stubSender{err: errors.New("boom")},
	})
	errs := s.Dispatch(context.Background(), makeAlert())
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Name != "bad" {
		t.Errorf("expected error from 'bad', got %q", errs[0].Name)
	}
}

func TestDispatch_AllFail(t *testing.T) {
	s := sink.New(map[string]sink.Sender{
		"x": &stubSender{err: errors.New("x-err")},
		"y": &stubSender{err: errors.New("y-err")},
	})
	errs := s.Dispatch(context.Background(), makeAlert())
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
	names := []string{errs[0].Name, errs[1].Name}
	sort.Strings(names)
	if names[0] != "x" || names[1] != "y" {
		t.Errorf("unexpected error names: %v", names)
	}
}

func TestSendError_ErrorString(t *testing.T) {
	e := sink.SendError{Name: "webhook", Err: errors.New("timeout")}
	got := e.Error()
	if got != `sink "webhook": timeout` {
		t.Errorf("unexpected error string: %q", got)
	}
}
