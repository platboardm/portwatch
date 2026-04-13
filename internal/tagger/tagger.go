// Package tagger provides utilities for normalising and validating
// target tags used in filter expressions and configuration.
package tagger

import (
	"fmt"
	"regexp"
	"strings"
)

// validTag matches lowercase alphanumeric words separated by hyphens.
var validTag = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// ErrInvalidTag is returned when a tag does not conform to the allowed format.
type ErrInvalidTag struct {
	Tag string
}

func (e ErrInvalidTag) Error() string {
	return fmt.Sprintf("tagger: invalid tag %q: must be lowercase alphanumeric with optional hyphens", e.Tag)
}

// Normalise lower-cases each tag and trims surrounding whitespace.
// It does not validate; call Validate separately when strict checking is needed.
func Normalise(tags []string) []string {
	out := make([]string, len(tags))
	for i, t := range tags {
		out[i] = strings.ToLower(strings.TrimSpace(t))
	}
	return out
}

// Validate returns an error for the first tag that does not match the
// allowed pattern (lowercase alphanumeric, hyphen-separated).
func Validate(tags []string) error {
	for _, t := range tags {
		if !validTag.MatchString(t) {
			return ErrInvalidTag{Tag: t}
		}
	}
	return nil
}

// Dedupe returns a new slice with duplicate tags removed, preserving order.
func Dedupe(tags []string) []string {
	seen := make(map[string]struct{}, len(tags))
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		if _, ok := seen[t]; !ok {
			seen[t] = struct{}{}
			out = append(out, t)
		}
	}
	return out
}

// Prepare is a convenience function that normalises, deduplicates, and
// validates tags in one call. It returns the cleaned slice or an error.
func Prepare(tags []string) ([]string, error) {
	norm := Dedupe(Normalise(tags))
	if err := Validate(norm); err != nil {
		return nil, err
	}
	return norm, nil
}
