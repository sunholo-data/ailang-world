package transitionreg

import (
	"context"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

// Binder is the descriptor-confinement seam implemented by a broker session.
type Binder interface {
	Bind(broker.Manifest) (*broker.BoundInvoker, error)
}

// CapabilitySource supplies one immutable capability reading.
type CapabilitySource interface {
	CapabilitySnapshot(now int64) broker.CapabilitySnapshot
}

// Proposal carries the authority-relevant descriptor pins into execution.
type Proposal struct {
	TransitionFn, Interpreter hashref.HashRef
	SemanticsEpoch            int64
	RequiredCaps              EffectRequirement
	ExpectedEffects           []EffectRequirement
}

// Bound couples a copied descriptor to its broker-confined dispatch surface.
type Bound struct {
	descriptor Descriptor
	invoker    *broker.BoundInvoker
}

// TransitionAbsentError reports an ID missing from the captured registry.
type TransitionAbsentError struct{ ID string }

func (e *TransitionAbsentError) Error() string {
	return fmt.Sprintf("transition registry: transition %q is absent", e.ID)
}

// AccessDeniedError reports the broker's exact denial label.
type AccessDeniedError struct {
	ID    string
	Label string
}

func (e *AccessDeniedError) Error() string {
	return fmt.Sprintf("transition registry: access to %q denied: %s", e.ID, e.Label)
}

// ProposalMismatchError reports one descriptor pin that disagrees.
type ProposalMismatchError struct{ Field string }

func (e *ProposalMismatchError) Error() string {
	return fmt.Sprintf("transition registry: proposal %s does not match bound descriptor", e.Field)
}

// Bind resolves one descriptor, admits its access requirement against the
// captured capabilities, and confines dispatch to its declared effects.
func Bind(snap Snapshot, id string, caps broker.CapabilitySnapshot, target Binder) (*Bound, error) {
	if snap.Head.IsZero() {
		return nil, fmt.Errorf("transition registry: cannot bind %q from a zero snapshot", id)
	}
	d, ok := snap.Lookup(id)
	if !ok {
		return nil, &TransitionAbsentError{ID: id}
	}
	decision := broker.Allows(caps, requirementOf(d.Access))
	if !decision.Allowed {
		return nil, &AccessDeniedError{ID: id, Label: decision.Label}
	}
	invoker, err := target.Bind(broker.Manifest{
		Access:   requirementOf(d.Access),
		Declared: requirementsOf(d.DeclaredEffects),
	})
	if err != nil {
		return nil, fmt.Errorf("transition registry: bind %q: %w", id, err)
	}
	return &Bound{descriptor: cloneDescriptor(d), invoker: invoker}, nil
}

// Descriptor returns a deep copy of the bound descriptor.
func (b *Bound) Descriptor() Descriptor { return cloneDescriptor(b.descriptor) }

// Check requires exact agreement on every authority-bearing proposal field.
func (b *Bound) Check(p Proposal) error {
	if p.TransitionFn != b.descriptor.TransitionFn {
		return &ProposalMismatchError{Field: "transition function"}
	}
	if p.Interpreter != b.descriptor.Interpreter {
		return &ProposalMismatchError{Field: "interpreter"}
	}
	if p.SemanticsEpoch != b.descriptor.SemanticsEpoch {
		return &ProposalMismatchError{Field: "semantics epoch"}
	}
	if p.RequiredCaps != b.descriptor.Access {
		return &ProposalMismatchError{Field: "required capabilities"}
	}
	if !equalRequirements(p.ExpectedEffects, b.descriptor.DeclaredEffects) {
		return &ProposalMismatchError{Field: "expected effects"}
	}
	return nil
}

// Request dispatches only through the broker's descriptor-bound invoker.
func (b *Bound) Request(ctx context.Context, req broker.EffectRequest, payload []byte) ([]byte, hashref.HashRef, error) {
	return b.invoker.Request(ctx, req, payload)
}

// Request is one consumer request with both mutable sources captured once.
type Request struct {
	Registry Snapshot
	Caps     broker.CapabilitySnapshot
}

// NewRequest captures exactly one registry and one capability reading.
func NewRequest(ctx context.Context, r Reader, caps CapabilitySource, now int64) (Request, error) {
	snap, err := r.ReadSnapshot(ctx)
	if err != nil {
		return Request{}, fmt.Errorf("transition registry: construct request: %w", err)
	}
	return Request{Registry: snap, Caps: caps.CapabilitySnapshot(now)}, nil
}

// Allowed returns detached descriptors admitted by the captured capabilities,
// preserving the registry snapshot's bytewise order.
func (q Request) Allowed() []Descriptor {
	entries := q.Registry.List()
	allowed := make([]Descriptor, 0, len(entries))
	for _, d := range entries {
		if broker.Allows(q.Caps, requirementOf(d.Access)).Allowed {
			allowed = append(allowed, cloneDescriptor(d))
		}
	}
	return allowed
}

func requirementOf(r EffectRequirement) broker.Requirement {
	return broker.Requirement{Effect: r.Effect, Scope: r.Scope, Cost: r.Cost}
}

func requirementsOf(rs []EffectRequirement) []broker.Requirement {
	out := make([]broker.Requirement, len(rs))
	for i, r := range rs {
		out[i] = requirementOf(r)
	}
	return out
}

func equalRequirements(a, b []EffectRequirement) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
