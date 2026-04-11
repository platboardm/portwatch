// Package webhook provides an HTTP webhook notifier that posts alert
// payloads to a configured endpoint whenever a port status changes.
package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Payload is the JSON body sent to the webhook endpoint.
type Payload struct {
	Timestamp string `json:"timestamp"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Status    string `json:"status"`
	Severity  string `json:"severity"`
	Message   string `json:"message"`
}

// Notifier posts alert payloads to an HTTP endpoint.
type Notifier struct {
	endpoint string
	client   *http.Client
}

// New creates a Notifier that sends POST requests to endpoint.
// timeout controls the per-request deadline; zero uses 5 s.
func New(endpoint string, timeout time.Duration) *Notifier {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Notifier{
		endpoint: endpoint,
		client:   &http.Client{Timeout: timeout},
	}
}

// Notify serialises a into a JSON payload and POSTs it to the endpoint.
func (n *Notifier) Notify(ctx context.Context, a alert.Alert) error {
	p := Payload{
		Timestamp: a.OccurredAt.UTC().Format(time.RFC3339),
		Host:      a.Host,
		Port:      a.Port,
		Status:    a.Status.String(),
		Severity:  a.Severity.String(),
		Message:   a.String(),
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: POST %s: %w", n.endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: server returned %s", resp.Status)
	}
	return nil
}
