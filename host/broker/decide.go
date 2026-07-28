// Package broker enforces the capability and budget law frozen in
// design_docs/sketches/effectbroker.ail.
package broker

import "strconv"

const (
	LabelDeniedEffectName = "denied:effect-name"
	LabelDeniedScope      = "denied:scope"
	LabelDeniedExpired    = "denied:expired"
	LabelDeniedBudget     = "denied:budget"
)

// Capability is one exact-effect, exact-scope grant.
type Capability struct {
	Effect    string
	Scope     string
	ExpiresAt int64
	Budget    int64
}

// EffectRequest is the authority-relevant portion of an invocation. Now is
// supplied by the caller; the broker never reads a wall clock.
type EffectRequest struct {
	Effect string
	Scope  string
	Cost   int64
	Now    int64
}

// Decision is the canonical projection of the frozen broker decision.
type Decision struct {
	Allowed   bool
	Remaining int64
	Label     string
}

// Decide mirrors sketches/effectbroker.ail. The order is authority-bearing:
// effect name, scope, expiry, then remaining budget.
func Decide(c Capability, r EffectRequest) Decision {
	if !effectNameMatches(c, r.Effect) {
		return Decision{Label: LabelDeniedEffectName}
	}
	if !scopeMatches(c, r.Scope) {
		return Decision{Label: LabelDeniedScope}
	}
	if !capabilityLive(c, r.Now) {
		return Decision{Label: LabelDeniedExpired}
	}
	if !withinEffectBudget(c.Budget, r.Cost) {
		return Decision{Label: LabelDeniedBudget}
	}
	remaining := c.Budget - r.Cost
	return Decision{Allowed: true, Remaining: remaining, Label: allowedLabel(remaining)}
}

func allowedLabel(remaining int64) string {
	return "allowed:" + strconv.FormatInt(remaining, 10)
}

func effectNameMatches(c Capability, effect string) bool { return c.Effect == effect }
func scopeMatches(c Capability, scope string) bool       { return c.Scope == scope }
func capabilityLive(c Capability, now int64) bool        { return now >= 0 && now < c.ExpiresAt }
func withinEffectBudget(budget, cost int64) bool {
	return cost >= 0 && budget >= 0 && cost <= budget
}
func debit(budget, cost int64) int64 { return budget - cost }
func effectAllowed(c Capability, r EffectRequest) bool {
	return effectNameMatches(c, r.Effect) &&
		scopeMatches(c, r.Scope) &&
		capabilityLive(c, r.Now) &&
		withinEffectBudget(c.Budget, r.Cost)
}
