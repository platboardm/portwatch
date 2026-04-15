package envelope

import "fmt"

// Rule maps a priority range to a named channel.
type Rule struct {
	MinPriority Priority
	Channel     string
}

// Router assigns a channel to an envelope based on a priority-ordered list of
// rules. The first rule whose MinPriority is less than or equal to the
// envelope priority is selected. If no rule matches, the fallback channel is
// used.
type Router struct {
	rules    []Rule
	fallback string
}

// NewRouter returns a Router with the given rules and fallback channel.
// rules should be ordered from highest to lowest MinPriority so that the most
// specific rule is matched first. Panics if fallback is empty.
func NewRouter(rules []Rule, fallback string) *Router {
	if fallback == "" {
		panic("envelope: NewRouter fallback channel must not be empty")
	}
	return &Router{rules: rules, fallback: fallback}
}

// Route returns a copy of e with Channel set according to the configured
// rules. The original envelope is not modified.
func (r *Router) Route(e Envelope) Envelope {
	for _, rule := range r.rules {
		if e.Priority >= rule.MinPriority {
			e.Channel = rule.Channel
			return e
		}
	}
	e.Channel = r.fallback
	return e
}

// String returns a human-readable description of the router configuration.
func (r *Router) String() string {
	return fmt.Sprintf("Router{rules:%d fallback:%q}", len(r.rules), r.fallback)
}
