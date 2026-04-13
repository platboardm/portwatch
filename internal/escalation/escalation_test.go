package escalation_test

import (
	"testing"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/escalation"
)

func makeAlert(target string, sev alert.Severity) alert.Alert {
	return alert.New(target, "127.0.0.1", 8080, sev)
}

func TestNew_PanicsOnZeroThreshold(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero threshold")
		}
	}()
	escalation.New(0)
}

func TestNew_PanicsOnNegativeThreshold(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative threshold")
		}
	}()
	escalation.New(-time.Second)
}

func TestEvaluate_NotEscalatedOnFirstDown(t *testing.T) {
	e := escalation.New(5 * time.Minute)
	a := makeAlert("svc-a", alert.SeverityCritical)
	if e.Evaluate(a) {
		t.Fatal("should not escalate immediately on first down event")
	}
}

func TestEvaluate_EscalatesAfterThreshold(t *testing.T) {
	e := escalation.New(10 * time.Millisecond)
	a := makeAlert("svc-b", alert.SeverityCritical)

	// Seed the tracker.
	e.Evaluate(a)

	time.Sleep(20 * time.Millisecond)

	if !e.Evaluate(a) {
		t.Fatal("expected escalation after threshold elapsed")
	}
}

func TestEvaluate_RecoveryClears(t *testing.T) {
	e := escalation.New(10 * time.Millisecond)
	down := makeAlert("svc-c", alert.SeverityCritical)
	up := makeAlert("svc-c", alert.SeverityInfo)

	e.Evaluate(down)
	time.Sleep(20 * time.Millisecond)

	// Recovery should clear the tracker.
	e.Evaluate(up)

	_, tracked := e.DownSince("svc-c")
	if tracked {
		t.Fatal("target should no longer be tracked after recovery")
	}
}

func TestReset_ClearsTarget(t *testing.T) {
	e := escalation.New(time.Minute)
	a := makeAlert("svc-d", alert.SeverityCritical)
	e.Evaluate(a)

	_, tracked := e.DownSince("svc-d")
	if !tracked {
		t.Fatal("expected target to be tracked after down event")
	}

	e.Reset("svc-d")
	_, tracked = e.DownSince("svc-d")
	if tracked {
		t.Fatal("expected target to be cleared after Reset")
	}
}

func TestEvaluate_IndependentTargets(t *testing.T) {
	e := escalation.New(10 * time.Millisecond)
	a1 := makeAlert("svc-x", alert.SeverityCritical)
	a2 := makeAlert("svc-y", alert.SeverityCritical)

	e.Evaluate(a1)
	time.Sleep(20 * time.Millisecond)
	e.Evaluate(a2) // svc-y just started going down

	if !e.Evaluate(a1) {
		t.Fatal("svc-x should be escalated")
	}
	if e.Evaluate(a2) {
		t.Fatal("svc-y should not yet be escalated")
	}
}
