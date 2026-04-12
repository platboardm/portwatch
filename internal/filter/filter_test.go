package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/filter"
)

func makeTargets() []config.Target {
	return []config.Target{
		{Name: "api-server", Host: "localhost", Port: 8080, Tags: []string{"http", "prod"}},
		{Name: "db-primary", Host: "localhost", Port: 5432, Tags: []string{"db", "prod"}},
		{Name: "cache", Host: "localhost", Port: 6379, Tags: []string{"cache", "staging"}},
		{Name: "api-internal", Host: "localhost", Port: 9090, Tags: []string{"http", "staging"}},
	}
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	targets := makeTargets()
	got := filter.Apply(targets, filter.Options{})
	if len(got) != len(targets) {
		t.Fatalf("expected %d targets, got %d", len(targets), len(got))
	}
}

func TestApply_FilterByTag(t *testing.T) {
	got := filter.Apply(makeTargets(), filter.Options{Tags: []string{"prod"}})
	if len(got) != 2 {
		t.Fatalf("expected 2 prod targets, got %d", len(got))
	}
	for _, tgt := range got {
		if tgt.Name != "api-server" && tgt.Name != "db-primary" {
			t.Errorf("unexpected target %q", tgt.Name)
		}
	}
}

func TestApply_FilterByMultipleTags(t *testing.T) {
	got := filter.Apply(makeTargets(), filter.Options{Tags: []string{"http", "prod"}})
	if len(got) != 1 {
		t.Fatalf("expected 1 target, got %d", len(got))
	}
	if got[0].Name != "api-server" {
		t.Errorf("expected api-server, got %q", got[0].Name)
	}
}

func TestApply_FilterByNamePrefix(t *testing.T) {
	got := filter.Apply(makeTargets(), filter.Options{NamePrefix: "api"})
	if len(got) != 2 {
		t.Fatalf("expected 2 api targets, got %d", len(got))
	}
}

func TestApply_NamePrefixCaseInsensitive(t *testing.T) {
	got := filter.Apply(makeTargets(), filter.Options{NamePrefix: "API"})
	if len(got) != 2 {
		t.Fatalf("expected 2 targets for prefix API, got %d", len(got))
	}
}

func TestApply_CombinedPrefixAndTag(t *testing.T) {
	got := filter.Apply(makeTargets(), filter.Options{
		NamePrefix: "api",
		Tags:       []string{"staging"},
	})
	if len(got) != 1 {
		t.Fatalf("expected 1 target, got %d", len(got))
	}
	if got[0].Name != "api-internal" {
		t.Errorf("expected api-internal, got %q", got[0].Name)
	}
}

func TestApply_NoMatch_ReturnsEmpty(t *testing.T) {
	got := filter.Apply(makeTargets(), filter.Options{Tags: []string{"nonexistent"}})
	if len(got) != 0 {
		t.Fatalf("expected 0 targets, got %d", len(got))
	}
}
