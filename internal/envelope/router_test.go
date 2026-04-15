package envelope_test

import (
	"testing"

	"github.com/user/portwatch/internal/envelope"
)

func TestNewRouter_PanicsOnEmptyFallback(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty fallback")
		}
	}()
	envelope.NewRouter(nil, "")
}

func TestRoute_FallbackWhenNoRules(t *testing.T) {
	r := envelope.NewRouter(nil, "default")
	e := envelope.New(makeAlert(), envelope.PriorityCritical, "")
	got := r.Route(e)
	if got.Channel != "default" {
		t.Errorf("expected channel default, got %q", got.Channel)
	}
}

func TestRoute_MatchesHighPriority(t *testing.T) {
	rules := []envelope.Rule{
		{MinPriority: envelope.PriorityCritical, Channel: "pagerduty"},
		{MinPriority: envelope.PriorityHigh, Channel: "slack"},
	}
	r := envelope.NewRouter(rules, "email")

	e := envelope.New(makeAlert(), envelope.PriorityCritical, "")
	got := r.Route(e)
	if got.Channel != "pagerduty" {
		t.Errorf("expected pagerduty, got %q", got.Channel)
	}
}

func TestRoute_MatchesLowerPriority(t *testing.T) {
	rules := []envelope.Rule{
		{MinPriority: envelope.PriorityCritical, Channel: "pagerduty"},
		{MinPriority: envelope.PriorityHigh, Channel: "slack"},
	}
	r := envelope.NewRouter(rules, "email")

	e := envelope.New(makeAlert(), envelope.PriorityHigh, "")
	got := r.Route(e)
	if got.Channel != "slack" {
		t.Errorf("expected slack, got %q", got.Channel)
	}
}

func TestRoute_FallbackForLowPriority(t *testing.T) {
	rules := []envelope.Rule{
		{MinPriority: envelope.PriorityCritical, Channel: "pagerduty"},
	}
	r := envelope.NewRouter(rules, "email")

	e := envelope.New(makeAlert(), envelope.PriorityLow, "")
	got := r.Route(e)
	if got.Channel != "email" {
		t.Errorf("expected email, got %q", got.Channel)
	}
}

func TestRoute_DoesNotMutateOriginal(t *testing.T) {
	r := envelope.NewRouter(nil, "default")
	e := envelope.New(makeAlert(), envelope.PriorityNormal, "original")
	_ = r.Route(e)
	if e.Channel != "original" {
		t.Error("Route must not mutate the original envelope")
	}
}

func TestRouter_String(t *testing.T) {
	r := envelope.NewRouter([]envelope.Rule{{MinPriority: envelope.PriorityHigh, Channel: "slack"}}, "email")
	s := r.String()
	if s == "" {
		t.Error("expected non-empty String()")
	}
}
