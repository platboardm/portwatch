package aggregator_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/aggregator"
	"github.com/user/portwatch/internal/alert"
)

func makeAlert(target string, sev alert.Severity) alert.Alert {
	al, _ := alert.New(target, "localhost", 8080, sev)
	return al
}

func TestNew_PanicsOnZeroCap(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for recentCap=0")
		}
	}()
	aggregator.New(0)
}

func TestSummarise_Empty(t *testing.T) {
	a := aggregator.New(10)
	s := a.Summarise()
	if s.Total != 0 || s.Up != 0 || s.Down != 0 {
		t.Fatalf("expected zero summary, got %+v", s)
	}
}

func TestRecord_UpdatesStatus(t *testing.T) {
	a := aggregator.New(10)
	a.Record(makeAlert("svc-a", alert.SeverityCritical))
	a.Record(makeAlert("svc-b", alert.SeverityInfo))

	s := a.Summarise()
	if s.Total != 2 {
		t.Fatalf("want Total=2, got %d", s.Total)
	}
	if s.Down != 1 {
		t.Fatalf("want Down=1, got %d", s.Down)
	}
	if s.Up != 1 {
		t.Fatalf("want Up=1, got %d", s.Up)
	}
}

func TestRecord_OverwritesSameTarget(t *testing.T) {
	a := aggregator.New(10)
	a.Record(makeAlert("svc-a", alert.SeverityCritical))
	a.Record(makeAlert("svc-a", alert.SeverityInfo))

	s := a.Summarise()
	if s.Total != 1 {
		t.Fatalf("want Total=1, got %d", s.Total)
	}
	if s.Up != 1 || s.Down != 0 {
		t.Fatalf("want Up=1 Down=0, got %+v", s)
	}
}

func TestRecord_RecentCap(t *testing.T) {
	a := aggregator.New(3)
	for i := 0; i < 5; i++ {
		a.Record(makeAlert("svc", alert.SeverityCritical))
	}
	s := a.Summarise()
	if len(s.Alerts) != 3 {
		t.Fatalf("want 3 recent alerts, got %d", len(s.Alerts))
	}
}

func TestSummarise_Timestamp(t *testing.T) {
	a := aggregator.New(5)
	before := time.Now().UTC()
	s := a.Summarise()
	after := time.Now().UTC()
	if s.At.Before(before) || s.At.After(after) {
		t.Fatalf("timestamp %v out of range [%v, %v]", s.At, before, after)
	}
}
