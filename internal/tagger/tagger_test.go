package tagger_test

import (
	"testing"

	"github.com/example/portwatch/internal/tagger"
)

func TestNormalise_LowercasesAndTrims(t *testing.T) {
	input := []string{"  Web  ", "API", "DB"}
	got := tagger.Normalise(input)
	want := []string{"web", "api", "db"}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: got %q, want %q", i, got[i], w)
		}
	}
}

func TestValidate_AcceptsValidTags(t *testing.T) {
	tags := []string{"web", "api-gateway", "db1", "cache-01"}
	if err := tagger.Validate(tags); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_RejectsUppercase(t *testing.T) {
	if err := tagger.Validate([]string{"Web"}); err == nil {
		t.Fatal("expected error for uppercase tag, got nil")
	}
}

func TestValidate_RejectsSpaces(t *testing.T) {
	if err := tagger.Validate([]string{"my tag"}); err == nil {
		t.Fatal("expected error for tag with space, got nil")
	}
}

func TestValidate_RejectsLeadingHyphen(t *testing.T) {
	if err := tagger.Validate([]string{"-bad"}); err == nil {
		t.Fatal("expected error for leading hyphen, got nil")
	}
}

func TestDedupe_RemovesDuplicates(t *testing.T) {
	input := []string{"web", "api", "web", "db", "api"}
	got := tagger.Dedupe(input)
	if len(got) != 3 {
		t.Fatalf("got %d tags, want 3: %v", len(got), got)
	}
	if got[0] != "web" || got[1] != "api" || got[2] != "db" {
		t.Errorf("unexpected order or values: %v", got)
	}
}

func TestDedupe_EmptyInput(t *testing.T) {
	got := tagger.Dedupe(nil)
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestPrepare_NormalisesAndValidates(t *testing.T) {
	input := []string{" WEB ", "API", "web"}
	got, err := tagger.Prepare(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// After normalise + dedupe: ["web", "api"]
	if len(got) != 2 {
		t.Fatalf("got %d tags, want 2: %v", len(got), got)
	}
}

func TestPrepare_ReturnsErrorForInvalidTag(t *testing.T) {
	_, err := tagger.Prepare([]string{"valid", "INVALID"})
	if err == nil {
		t.Fatal("expected error for invalid tag, got nil")
	}
}

func TestErrInvalidTag_Message(t *testing.T) {
	err := tagger.ErrInvalidTag{Tag: "Bad Tag"}
	if err.Error() == "" {
		t.Fatal("expected non-empty error message")
	}
}
