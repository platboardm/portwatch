package escalation_test

import (
	"testing"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/escalation"
)

// TestEscalator_FullLifecycle exercises a realistic down → escalate → recover
// sequence using a very short threshold so the test runs quickly.
func TestEscalator_FullLifecycle(t *testing.T) {
	const threshold = 30 * time.Millisecond
	e := escalation.New(threshold)
	target := "db-primary"

	down := alert.New(target, "10.0.0.1", 5432, alert.SeverityCritical)
	up := alert.New(target, "10.0.0.1", 5432, alert.SeverityInfo)

	// First down event — not yet escalated.
	if e.Evaluate(down) {
		t.Fatal("should not escalate on first down event")
	}

	_, tracked := e.DownSince(target)
	if !tracked {
		t.Fatal("target should be tracked after first down event")
	}

	// Wait for threshold to elapse.
	time.Sleep(threshold + 10*time.Millisecond)

	if !e.Evaluate(down) {
		t.Fatal("expected escalation after threshold")
	}

	// Recovery clears the tracker.
	e.Evaluate(up)

	_, tracked = e.DownSince(target)
	if tracked {
		t.Fatal("target should not be tracked after recovery")
	}

	// A subsequent down event restarts the clock — should not escalate yet.
	if e.Evaluate(down) {
		t.Fatal("should not escalate immediately after recovery + re-down")
	}
}
