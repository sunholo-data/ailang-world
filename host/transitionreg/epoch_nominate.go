package transitionreg

import (
	"context"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/registry"
	"github.com/sunholo-data/ailang-world/host/store"
)

// The epoch-registry surface `world-publish epoch nominate` (row 139 D2)
// uses. This package is already the publisher's only view of
// world/epoch-registry/v1 (EpochsForInterpreter, verifyEpochs), and
// cmd/world-publish may not import host/registry directly (host/boundary
// allowlist), so the verb reaches the registry through these re-exports.
// They are aliases and one-line delegations: no behaviour lives here.

// EpochMode is registry.Mode.
type EpochMode = registry.Mode

// The two nomination shapes.
const (
	EpochCandidate = registry.ModeCandidate
	EpochNew       = registry.ModeNewEpoch
)

// EpochRegistrySemanticID is registry.SemanticID.
const EpochRegistrySemanticID = registry.SemanticID

// EpochRegistry is registry.Registry.
type EpochRegistry = registry.Registry

// EpochAlreadyNominatedError is registry.AlreadyNominatedError.
type EpochAlreadyNominatedError = registry.AlreadyNominatedError

// LoadEpochRegistry is registry.Load.
func LoadEpochRegistry(ctx context.Context, s *store.Store) (EpochRegistry, hashref.HashRef, bool, error) {
	return registry.Load(ctx, s)
}

// NextEpochRegistry is registry.Next (pure).
func NextEpochRegistry(cur EpochRegistry, release string, mode EpochMode) (EpochRegistry, error) {
	return registry.Next(cur, release, mode)
}

// NominateEpoch is registry.Nominate.
func NominateEpoch(ctx context.Context, s *store.Store, head hashref.HashRef, release string, mode EpochMode) (EpochRegistry, hashref.HashRef, error) {
	return registry.Nominate(ctx, s, head, release, mode)
}
