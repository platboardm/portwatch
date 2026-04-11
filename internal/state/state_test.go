package state_test

import (
	"sync"
	"testing"

	"portwatch/internal/state"
)

func TestStatusString(t *testing.T) {
	cases := []struct {
		s    state.Status
		want string
	}{
		{state.Up, "up"},
		{state.Down, "down"},
		{state.Unknown, "unknown"},
	}
	for _, tc := range cases {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("Status(%d).String() = %q, want %q", tc.s, got, tc.want)
		}
	}
}

func TestStore_GetUnknownKey(t *testing.T) {
	st := state.New()
	if got := st.Get("missing"); got != state.Unknown {
		t.Errorf("expected Unknown for missing key, got %v", got)
	}
}

func TestStore_SetAndGet(t *testing.T) {
	st := state.New()
	changed := st.Set("host:80", state.Up)
	if !changed {
		t.Error("first Set should report changed=true")
	}
	if got := st.Get("host:80"); got != state.Up {
		t.Errorf("Get after Set = %v, want Up", got)
	}
}

func TestStore_SetNoChange(t *testing.T) {
	st := state.New()
	st.Set("host:80", state.Down)
	changed := st.Set("host:80", state.Down)
	if changed {
		t.Error("Set with same value should report changed=false")
	}
}

func TestStore_SetTransition(t *testing.T) {
	st := state.New()
	st.Set("svc", state.Up)
	changed := st.Set("svc", state.Down)
	if !changed {
		t.Error("transition Up→Down should report changed=true")
	}
}

func TestStore_Snapshot(t *testing.T) {
	st := state.New()
	st.Set("a:1", state.Up)
	st.Set("b:2", state.Down)
	snap := st.Snapshot()
	if snap["a:1"] != state.Up || snap["b:2"] != state.Down {
		t.Errorf("unexpected snapshot: %v", snap)
	}
	// Mutating snapshot must not affect store.
	snap["a:1"] = state.Down
	if st.Get("a:1") != state.Up {
		t.Error("mutating snapshot should not affect store")
	}
}

func TestStore_ConcurrentAccess(t *testing.T) {
	st := state.New()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "host:80"
			st.Set(key, state.Up)
			_ = st.Get(key)
			_ = st.Snapshot()
		}(i)
	}
	wg.Wait()
}
