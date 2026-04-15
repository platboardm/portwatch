package dedup_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/dedup"
)

func makeAlert(target, host string, port int) alert.Alert {
	a, _ := alert.New(target, host, port, alert.SeverityWarning)
	return a
}

func TestNew_PanicsOnZeroWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero window")
		}
	}()
	dedup.New(0)
}

func TestNew_PanicsOnNegativeWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative window")
		}
	}()
	dedup.New(-time.Second)
}

func TestAllow_FirstCallAlwaysTrue(t *testing.T) {
	d := dedup.New(time.Minute)
	a := makeAlert("svc", "localhost", 8080)
	if !d.Allow(a) {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_DuplicateSuppressed(t *testing.T) {
	d := dedup.New(time.Minute)
	a := makeAlert("svc", "localhost", 8080)
	d.Allow(a)
	if d.Allow(a) {
		t.Fatal("expected duplicate to be suppressed")
	}
}

func TestAllow_DifferentTargetsIndependent(t *testing.T) {
	d := dedup.New(time.Minute)
	a1 := makeAlert("svc-a", "localhost", 8080)
	a2 := makeAlert("svc-b", "localhost", 9090)
	d.Allow(a1)
	if !d.Allow(a2) {
		t.Fatal("expected different target to be allowed")
	}
}

func TestAllow_PassesAfterWindowExpires(t *testing.T) {
	now := time.Now()
	d := dedup.New(50 * time.Millisecond)
	// Inject a controllable clock via the exported field (white-box via same package).
	// We test the observable behaviour by sleeping past the window instead.
	a := makeAlert("svc", "localhost", 8080)
	d.Allow(a)
	time.Sleep(60 * time.Millisecond)
	_ = now
	if !d.Allow(a) {
		t.Fatal("expected alert to be allowed after window expires")
	}
}

func TestLen_TracksEntries(t *testing.T) {
	d := dedup.New(time.Minute)
	if d.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", d.Len())
	}
	d.Allow(makeAlert("a", "localhost", 1))
	d.Allow(makeAlert("b", "localhost", 2))
	if d.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", d.Len())
	}
}
