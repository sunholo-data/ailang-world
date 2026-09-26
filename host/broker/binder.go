package broker

import "github.com/sunholo-data/ailang-world/host/store"

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

// Bind validates and copies one descriptor authority envelope.
func (b *SessionBinder) Bind(m Manifest) (*BoundInvoker, error) { return b.s.Bind(m) }
