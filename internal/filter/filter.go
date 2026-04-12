// Package filter provides target filtering based on tags and name patterns.
package filter

import (
	"strings"

	"github.com/user/portwatch/internal/config"
)

// Options holds the criteria used to select a subset of targets.
type Options struct {
	// Tags, when non-empty, keeps only targets that carry ALL listed tags.
	Tags []string
	// NamePrefix, when non-empty, keeps only targets whose name starts with
	// the given prefix (case-insensitive).
	NamePrefix string
}

// Apply returns the subset of targets that match every criterion in opts.
// If opts is the zero value, all targets are returned unchanged.
func Apply(targets []config.Target, opts Options) []config.Target {
	if len(opts.Tags) == 0 && opts.NamePrefix == "" {
		return targets
	}

	out := make([]config.Target, 0, len(targets))
	for _, t := range targets {
		if opts.NamePrefix != "" &&
			!strings.HasPrefix(strings.ToLower(t.Name), strings.ToLower(opts.NamePrefix)) {
			continue
		}
		if !hasAllTags(t.Tags, opts.Tags) {
			continue
		}
		out = append(out, t)
	}
	return out
}

// hasAllTags reports whether src contains every tag in required.
func hasAllTags(src, required []string) bool {
	if len(required) == 0 {
		return true
	}
	set := make(map[string]struct{}, len(src))
	for _, s := range src {
		set[strings.ToLower(s)] = struct{}{}
	}
	for _, r := range required {
		if _, ok := set[strings.ToLower(r)]; !ok {
			return false
		}
	}
	return true
}
