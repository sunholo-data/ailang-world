package broker

import (
	"context"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// Manifest is one descriptor's authority envelope, copied out of the registry.
type Manifest struct {
	Access   Requirement
	Declared []Requirement
}

// BoundInvoker is the only dispatch surface intended for descriptor-bound
// consumers. It retains a private copy of the declaration set.
type BoundInvoker struct {
	s        *Session
	declared []Requirement
}

// UndeclaredEffectError reports a request outside the complete declared
// effect/scope/cost triples in the manifest.
type UndeclaredEffectError struct {
	Effect string
	Scope  string
	Cost   int64
}

func (e *UndeclaredEffectError) Error() string {
	return fmt.Sprintf("broker: undeclared effect request: effect %q scope %q cost %d", e.Effect, e.Scope, e.Cost)
}

// Bind validates and copies a descriptor authority envelope.
func (s *Session) Bind(m Manifest) (*BoundInvoker, error) {
	if m.Access.Cost < 0 {
		return nil, fmt.Errorf("broker: manifest access cost must be non-negative: %d", m.Access.Cost)
	}
	declared := append([]Requirement(nil), m.Declared...)
	for i, req := range declared {
		if req.Cost < 0 {
			return nil, fmt.Errorf("broker: manifest declared cost at index %d must be non-negative: %d", i, req.Cost)
		}
		for j := 0; j < i; j++ {
			if declared[j] == req {
				return nil, fmt.Errorf("broker: duplicate declared requirement at indexes %d and %d: effect %q scope %q cost %d", j, i, req.Effect, req.Scope, req.Cost)
			}
		}
	}
	return &BoundInvoker{s: s, declared: declared}, nil
}

// Declared returns a fresh copy of the bound declaration set.
func (b *BoundInvoker) Declared() []Requirement {
	return append([]Requirement(nil), b.declared...)
}

// Request refuses undeclared triples before entering the broker pipeline.
func (b *BoundInvoker) Request(
	ctx context.Context,
	req EffectRequest,
	payload []byte,
) ([]byte, hashref.HashRef, error) {
	for _, declared := range b.declared {
		if declared.Effect == req.Effect && declared.Scope == req.Scope && declared.Cost == req.Cost {
			return b.s.invoke(ctx, req, payload)
		}
	}
	return nil, hashref.HashRef{}, &UndeclaredEffectError{Effect: req.Effect, Scope: req.Scope, Cost: req.Cost}
}
