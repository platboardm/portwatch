package healthcheck_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"portwatch/internal/healthcheck"
)

func TestHandler_ReturnsOK(t *testing.T) {
	started := time.Now().Add(-2 * time.Minute)
	h := healthcheck.Handler(started)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandler_ContentType(t *testing.T) {
	h := healthcheck.Handler(time.Now())
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	ct := rec.Header().Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Fatalf("unexpected Content-Type: %s", ct)
	}
}

func TestHandler_JSONFields(t *testing.T) {
	started := time.Now().Add(-90 * time.Second)
	h := healthcheck.Handler(started)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	var status healthcheck.Status
	if err := json.NewDecoder(rec.Body).Decode(&status); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !status.OK {
		t.Error("expected ok=true")
	}
	if status.Uptime == "" {
		t.Error("expected non-empty uptime")
	}
	if status.Started.IsZero() {
		t.Error("expected non-zero started time")
	}
}

func TestHandler_UptimeIncreases(t *testing.T) {
	started := time.Now().Add(-10 * time.Second)
	h := healthcheck.Handler(started)

	decode := func() healthcheck.Status {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
		var s healthcheck.Status
		_ = json.NewDecoder(rec.Body).Decode(&s)
		return s
	}

	s := decode()
	if s.Uptime == "0s" {
		t.Error("uptime should be > 0s for a daemon started 10s ago")
	}
}
