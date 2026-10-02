package broker

import (
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// SessionBinder is the only production handle on a session outside
// host/broker (w-transition-invocation-coordinator M2). It can Bind a
// descriptor manifest and nothing else: the raw Session, and therefore raw
// dispatch, stays unreachable to every other package (TR.C). A consumer
// dispatches only through the returned BoundInvoker, which refuses every
// undeclared effect/scope/cost triple.
type SessionBinder struct{ s *Session }

// OpenBinder opens a live, session-scoped binder for one resolved session.
// The handler registry is empty: slice-1 invocations declare no effects.
func OpenBinder(s *store.Store, episodeID string, grants []Capability) *SessionBinder {
	return &SessionBinder{s: newSession(s, episodeID, grants, nil, Live, nil)}
}

// OpenReplayBinder opens a replay binder over the ordered effect records of a
// committed invocation (row 134 §4.1 Replay). Its handler registry is nil and
// its mode is Replay: every request is answered from recordRefs and no
// handler is ever dispatched (V10).
func OpenReplayBinder(s *store.Store, grants []Capability, recordRefs []hashref.HashRef) *SessionBinder {
	return &SessionBinder{s: newSession(s, "", grants, nil, Replay, recordRefs)}
}

// ReplayGrant returns a capability under which the frozen decision law
// reproduces rec's recorded decision for a request of rec's effect, scope and
// cost at logical time now: the grant's budget is the record's BudgetBefore
// (replay seeds from each record's BudgetBefore, row 134 §4.4), and a denial
// label is reproduced by the one grant field that law checks for it.
func ReplayGrant(rec EffectRecord, now int64) Capability {
	c := Capability{Effect: rec.Effect, Scope: rec.Scope, ExpiresAt: now + 1, Budget: rec.BudgetBefore}
	switch rec.Denial {
	case LabelDeniedEffectName:
		c.Effect += "\x00replay"
	case LabelDeniedScope:
		c.Scope += "\x00replay"
	case LabelDeniedExpired:
		c.ExpiresAt = now
	}
	return c
}

// Bind validates and copies one descriptor authority envelope.
func (b *SessionBinder) Bind(m Manifest) (*BoundInvoker, error) { return b.s.Bind(m) }
