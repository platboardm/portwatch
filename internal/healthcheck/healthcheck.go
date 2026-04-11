// Package healthcheck exposes an HTTP endpoint that reports the current
// health and uptime of the portwatch daemon.
package healthcheck

import (
	"encoding/json"
	"net/http"
	"time"
)

// Status represents the JSON payload returned by the health endpoint.
type Status struct {
	OK      bool      `json:"ok"`
	Uptime  string    `json:"uptime"`
	Started time.Time `json:"started"`
}

// Handler returns an [http.Handler] that responds with a JSON health status.
// started is the time the daemon was launched; it is captured once and reused
// on every request.
func Handler(started time.Time) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status := Status{
			OK:      true,
			Uptime:  time.Since(started).Round(time.Second).String(),
			Started: started,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(status)
	})
}

// Server wraps an [http.Server] and exposes a single /healthz endpoint.
type Server struct {
	server *http.Server
}

// NewServer creates a Server that listens on addr (e.g. ":9090").
func NewServer(addr string, started time.Time) *Server {
	mux := http.NewServeMux()
	mux.Handle("/healthz", Handler(started))

	return &Server{
		server: &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
}

// ListenAndServe starts the HTTP server. It blocks until the server stops.
func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

// Close shuts the server down immediately.
func (s *Server) Close() error {
	return s.server.Close()
}
