package webhook_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/webhook"
)

func makeAlert(t *testing.T) alert.Alert {
	t.Helper()
	a, err := alert.New("localhost", 8080, alert.SeverityCritical)
	if err != nil {
		t.Fatalf("alert.New: %v", err)
	}
	return a
}

func TestNotify_PostsJSON(t *testing.T) {
	var got webhook.Payload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected Content-Type: %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := webhook.New(srv.URL, 0)
	a := makeAlert(t)
	if err := n.Notify(context.Background(), a); err != nil {
		t.Fatalf("Notify: %v", err)
	}

	if got.Host != "localhost" {
		t.Errorf("host = %q, want %q", got.Host, "localhost")
	}
	if got.Port != 8080 {
		t.Errorf("port = %d, want 8080", got.Port)
	}
	if got.Timestamp == "" {
		t.Error("timestamp is empty")
	}
}

func TestNotify_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	n := webhook.New(srv.URL, 0)
	if err := n.Notify(context.Background(), makeAlert(t)); err == nil {
		t.Error("expected error for 500 response, got nil")
	}
}

func TestNotify_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := webhook.New(srv.URL, 50*time.Millisecond)
	if err := n.Notify(context.Background(), makeAlert(t)); err == nil {
		t.Error("expected timeout error, got nil")
	}
}

func TestNew_DefaultTimeout(t *testing.T) {
	// Zero timeout should not panic and should produce a usable notifier.
	n := webhook.New("http://example.com", 0)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}
