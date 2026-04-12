// Package aggregator collects check results across multiple targets and
// produces a rolled-up summary suitable for dashboards or status endpoints.
package aggregator

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Summary holds a point-in-time snapshot of overall service health.
type Summary struct {
	Total   int       `json:"total"`
	Up      int       `json:"up"`
	Down    int       `json:"down"`
	Unknown int       `json:"unknown"`
	Alerts  []alert.Alert `json:"alerts,omitempty"`
	At      time.Time `json:"at"`
}

// Aggregator accumulates alerts and tracks per-target status.
type Aggregator struct {
	mu     sync.RWMutex
	status map[string]alert.Severity
	recent []alert.Alert
	cap    int
}

// New returns an Aggregator that retains at most recentCap recent alerts.
// It panics if recentCap is less than 1.
func New(recentCap int) *Aggregator {
	if recentCap < 1 {
		panic("aggregator: recentCap must be >= 1")
	}
	return &Aggregator{
		status: make(map[string]alert.Severity),
		cap:    recentCap,
	}
}

// Record ingests an alert, updating the target's current severity and the
// recent-alerts ring.
func (a *Aggregator) Record(al alert.Alert) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.status[al.Target] = al.Severity

	a.recent = append(a.recent, al)
	if len(a.recent) > a.cap {
		a.recent = a.recent[len(a.recent)-a.cap:]
	}
}

// Summarise returns a snapshot of current health across all known targets.
func (a *Aggregator) Summarise() Summary {
	a.mu.RLock()
	defer a.mu.RUnlock()

	s := Summary{
		At:     time.Now().UTC(),
		Total:  len(a.status),
		Alerts: make([]alert.Alert, len(a.recent)),
	}
	copy(s.Alerts, a.recent)

	for _, sev := range a.status {
		switch sev {
		case alert.SeverityCritical:
			s.Down++
		case alert.SeverityInfo:
			s.Up++
		default:
			s.Unknown++
		}
	}
	return s
}
