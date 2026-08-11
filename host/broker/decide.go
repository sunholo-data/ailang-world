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

// Requirement is the authority requested by a registry descriptor. Expiry and
// budget belong to session grants, not registry metadata.
type Requirement struct {
	Effect string
	Scope  string
	Cost   int64
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

// decideOver applies the broker's one ranked selection mechanism to a ledger.
func decideOver(grants []Capability, r EffectRequest) (int, Decision) {
	if len(grants) == 0 {
		return -1, Decision{Label: LabelDeniedEffectName}
	}
	bestIndex := 0
	best := Decide(grants[0], r)
	for i := 1; i < len(grants); i++ {
		next := Decide(grants[i], r)
		if next.Allowed {
			return i, next
		}
		if denialRank(next.Label) > denialRank(best.Label) {
			bestIndex, best = i, next
		}
	}
	return bestIndex, best
}

// Allows evaluates a descriptor requirement over a captured ledger snapshot.
func Allows(snap CapabilitySnapshot, req Requirement) Decision {
	_, decision := decideOver(snap.grants, EffectRequest{
		Effect: req.Effect,
		Scope:  req.Scope,
		Cost:   req.Cost,
		Now:    snap.Now,
	})
	return decision
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
