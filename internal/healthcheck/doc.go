// Package healthcheck provides a lightweight HTTP health endpoint for the
// portwatch daemon.
//
// # Overview
//
// The package exposes a single HTTP handler that responds to GET /healthz with
// a JSON document containing:
//
//   - ok      – always true while the process is running
//   - uptime  – human-readable duration since the daemon started
//   - started – RFC 3339 timestamp of the daemon start time
//
// # Usage
//
//	started := time.Now()
//	srv := healthcheck.NewServer(":9090", started)
//	go srv.ListenAndServe()
//	// later …
//	srv.Close()
package healthcheck
