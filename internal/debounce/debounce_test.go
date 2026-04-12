package debounce_test

import (
	"testing"

	"github.com/example/portwatch/internal/debounce"
)

func TestNew_PanicsOnZeroThreshold(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for threshold=0")
		}
	}()
	debounce.New(0)
}

func TestNew_PanicsOnNegativeThreshold(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for threshold=-1")
		}
	}()
	debounce.New(-1)
}

func TestConfirm_ThresholdOne_ImmediatelyTrue(t *testing.T) {
	d := debounce.New(1)
	if !d.Confirm("svc", true) {
		t.Fatal("expected true on first confirmation with threshold=1")
	}
}

func TestConfirm_NotConfirmedBeforeThreshold(t *testing.T) {
	d := debounce.New(3)
	for i := 0; i < 2; i++ {
		if d.Confirm("svc", true) {
			t.Fatalf("expected false on call %d (threshold not reached)", i+1)
		}
	}
}

func TestConfirm_ConfirmedAtThreshold(t *testing.T) {
	d := debounce.New(3)
	d.Confirm("svc", true)
	d.Confirm("svc", true)
	if !d.Confirm("svc", true) {
		t.Fatal("expected true on third confirmation with threshold=3")
	}
}

func TestConfirm_ResetsAfterConfirmation(t *testing.T) {
	d := debounce.New(2)
	d.Confirm("svc", true)
	d.Confirm("svc", true) // confirmed; counter resets
	// next single call should not confirm again
	if d.Confirm("svc", true) {
		t.Fatal("expected false after reset post-confirmation")
	}
	if d.Pending("svc") != 1 {
		t.Fatalf("expected pending=1 after reset, got %d", d.Pending("svc"))
	}
}

func TestConfirm_FalseResetsCounter(t *testing.T) {
	d := debounce.New(3)
	d.Confirm("svc", true)
	d.Confirm("svc", true)
	d.Confirm("svc", false) // interrupted
	if d.Pending("svc") != 0 {
		t.Fatalf("expected pending=0 after false, got %d", d.Pending("svc"))
	}
	// must accumulate again from scratch
	if d.Confirm("svc", true) {
		t.Fatal("expected false after counter was reset")
	}
}

func TestReset_ClearsState(t *testing.T) {
	d := debounce.New(3)
	d.Confirm("svc", true)
	d.Confirm("svc", true)
	d.Reset("svc")
	if d.Pending("svc") != 0 {
		t.Fatalf("expected pending=0 after Reset, got %d", d.Pending("svc"))
	}
}

func TestPending_ReturnsCount(t *testing.T) {
	d := debounce.New(5)
	for i := 1; i <= 3; i++ {
		d.Confirm("svc", true)
		if got := d.Pending("svc"); got != i {
			t.Fatalf("after %d calls: expected pending=%d, got %d", i, i, got)
		}
	}
}

func TestConfirm_IsolatesKeys(t *testing.T) {
	d := debounce.New(2)
	d.Confirm("a", true)
	d.Confirm("b", true)
	if d.Confirm("a", true) {
		t.Fatal("key 'a' should not be confirmed yet")
	}
	if !d.Confirm("b", true) {
		t.Fatal("key 'b' should be confirmed independently")
	}
}
