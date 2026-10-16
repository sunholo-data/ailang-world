package registry

import (
	"context"
	"fmt"
	"strings"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/pins"
	"github.com/sunholo-data/ailang-world/host/store"
)

// Mode is the shape of a nomination (row 139 D2,
// design_docs/planned/w-interpreter-epoch-upgrade-path.md).
type Mode string

const (
	// ModeCandidate adds a release to the CURRENT (highest) epoch's
	// Candidates: same semantics, every replay golden byte-identical.
	ModeCandidate Mode = "candidate"
	// ModeNewEpoch appends epoch N+1 with Candidates [release]: the goldens
	// differ, semantics moved. Old entries keep their epoch (V5).
	ModeNewEpoch Mode = "new-epoch"
)

// EpochsNominating returns, in stored (ascending) order, the epochs whose
// Candidates contain release.
func (r Registry) EpochsNominating(release string) []int64 {
	var out []int64
	for _, rec := range r.Epochs {
		for _, c := range rec.Candidates {
			if c == release {
				out = append(out, rec.Epoch)
				break
			}
		}
	}
	return out
}

// Summary renders the registry's epochs and candidates on one line, e.g.
// `epoch 1: ["AILANG a" "AILANG b"]; epoch 2: ["AILANG c"]`.
func (r Registry) Summary() string {
	parts := make([]string, 0, len(r.Epochs))
	for _, rec := range r.Epochs {
		parts = append(parts, fmt.Sprintf("epoch %d: %q", rec.Epoch, rec.Candidates))
	}
	if len(parts) == 0 {
		return "no epochs"
	}
	return strings.Join(parts, "; ")
}

// NotNominatedError is the startup refusal (D3): the store has a registry
// head and no epoch in it nominates the served release.
type NotNominatedError struct {
	Release  string
	Head     hashref.HashRef
	Registry Registry
}

func (e *NotNominatedError) Error() string {
	return fmt.Sprintf("registry: existing head %q nominates release %q in no epoch (%s)",
		e.Head.String(), e.Release, e.Registry.Summary())
}

// NominateCommand is the exact operator command that would let a store at
// storePath start on the interpreter whose first `--version` line is
// servedRelease ("AILANG vX.Y.Z") and whose binary is binPath: a CANDIDATE
// nomination (D-WORLD-78 = a, the loop may run it) backed by the release's
// compatibility-gate verdict and the binary itself (hashed against the
// compiled-in release manifest). A served release that is not a
// manifest-shaped release ("unpinned", a dev build) has no nomination: ok is
// false.
func NominateCommand(storePath, servedRelease, binPath string) (string, bool) {
	rel, ok := pins.ReleaseOfVersionLine(servedRelease)
	if !ok {
		return "", false
	}
	return fmt.Sprintf("world-publish epoch nominate --store %s --release %s --candidate --gate-verdict <verdict.json> --binary %s",
		storePath, rel, binPath), true
}

// AlreadyNominatedError refuses a nomination of a release some epoch already
// names: through Nominate a release belongs to exactly one epoch.
type AlreadyNominatedError struct {
	Release string
	Epochs  []int64
}

func (e *AlreadyNominatedError) Error() string {
	return fmt.Sprintf("registry: release %q is already nominated by epoch(s) %v; nothing was written", e.Release, e.Epochs)
}

// validate checks the structural invariants of a revision Nominate extends:
// the v1 semantic ID, at least one epoch, epochs numbered from 1 contiguously
// in ascending order, every epoch with at least one non-blank trimmed
// candidate, and no release named by two epochs.
func (r Registry) validate() error {
	if r.SemanticID != SemanticID {
		return fmt.Errorf("registry: semantic ID %q, want %q", r.SemanticID, SemanticID)
	}
	if len(r.Epochs) == 0 {
		return fmt.Errorf("registry: revision holds no epochs")
	}
	seen := map[string]int64{}
	for i, rec := range r.Epochs {
		if want := firstEpoch + int64(i); rec.Epoch != want {
			return fmt.Errorf("registry: epoch at position %d is %d, want %d (epochs are contiguous from %d)", i, rec.Epoch, want, firstEpoch)
		}
		if len(rec.Candidates) == 0 {
			return fmt.Errorf("registry: epoch %d has no candidates", rec.Epoch)
		}
		for _, c := range rec.Candidates {
			if strings.TrimSpace(c) == "" || c != strings.TrimSpace(c) {
				return fmt.Errorf("registry: epoch %d has a blank or untrimmed candidate %q", rec.Epoch, c)
			}
			if prev, dup := seen[c]; dup && prev != rec.Epoch {
				return fmt.Errorf("registry: release %q is nominated by epochs %d and %d", c, prev, rec.Epoch)
			}
			seen[c] = rec.Epoch
		}
	}
	return nil
}

// Next is the PURE half of Nominate: the revision that nominating release in
// mode produces from cur. cur is never mutated. It refuses an invalid cur, a
// blank or untrimmed release, a release already nominated by any epoch, and
// an unknown mode.
func Next(cur Registry, release string, mode Mode) (Registry, error) {
	if err := cur.validate(); err != nil {
		return Registry{}, err
	}
	if strings.TrimSpace(release) == "" || release != strings.TrimSpace(release) {
		return Registry{}, fmt.Errorf("registry: release %q is blank or untrimmed", release)
	}
	if epochs := cur.EpochsNominating(release); len(epochs) > 0 {
		return Registry{}, &AlreadyNominatedError{Release: release, Epochs: epochs}
	}
	next := Registry{SemanticID: cur.SemanticID, Epochs: make([]EpochRecord, len(cur.Epochs))}
	for i, rec := range cur.Epochs {
		next.Epochs[i] = EpochRecord{Epoch: rec.Epoch, Candidates: append([]string(nil), rec.Candidates...)}
	}
	last := len(next.Epochs) - 1
	switch mode {
	case ModeCandidate:
		next.Epochs[last].Candidates = append(next.Epochs[last].Candidates, release)
	case ModeNewEpoch:
		next.Epochs = append(next.Epochs, EpochRecord{Epoch: next.Epochs[last].Epoch + 1, Candidates: []string{release}})
	default:
		return Registry{}, fmt.Errorf("registry: unknown nomination mode %q (want %q or %q)", mode, ModeCandidate, ModeNewEpoch)
	}
	return next, nil
}

type headReader interface {
	GetRegistryHead(context.Context, string) (hashref.HashRef, bool, error)
	GetObject(context.Context, hashref.HashRef) (store.Object, bool, error)
}

type nominateStore interface {
	headReader
	PutObject(context.Context, store.Object) error
	CompareAndSetRegistryHead(context.Context, string, hashref.HashRef, hashref.HashRef) error
}

// Load reads and decodes the revision the registry head names. ok is false
// (and the error nil) when the store has no head.
// Pass the caller's intended read ctx; this API supplies no default timeout.
func Load(ctx context.Context, s *store.Store) (Registry, hashref.HashRef, bool, error) {
	return loadHead(ctx, s)
}

func loadHead(ctx context.Context, s headReader) (Registry, hashref.HashRef, bool, error) {
	head, ok, err := s.GetRegistryHead(ctx, SemanticID)
	if err != nil || !ok {
		return Registry{}, hashref.HashRef{}, false, err
	}
	reg, err := load(ctx, s, head)
	if err != nil {
		return Registry{}, hashref.HashRef{}, false, err
	}
	return reg, head, true, nil
}

func load(ctx context.Context, s interface {
	GetObject(context.Context, hashref.HashRef) (store.Object, bool, error)
}, head hashref.HashRef) (Registry, error) {
	obj, found, err := s.GetObject(ctx, head)
	if err != nil {
		return Registry{}, err
	}
	if !found {
		return Registry{}, fmt.Errorf("registry: head %q names an absent object", head.String())
	}
	if obj.SemanticID != SemanticID {
		return Registry{}, fmt.Errorf("registry: head %q names a %q object, not %q", head.String(), obj.SemanticID, SemanticID)
	}
	return Decode(obj.Payload)
}

// Nominate appends one immutable registry revision (row 139 D2): the
// revision Next computes from the one head names, stored through the same
// object mechanism as Bootstrap (PutObject) and named by a COMPARE-AND-SET of
// the registry head from head to it. A stale head (the registry moved since
// the caller read it) is refused before anything is written, and a race lost
// between that read and the swap is refused by the store's CAS; both are a
// *store.RegistryCASConflict and leave the head unchanged. The previous
// revision object is never modified.
//
// Nominate does not decide WHO may nominate or whether a release qualified:
// its caller (world-publish epoch nominate) checks the gate verdict and the
// D-WORLD-78 authority split first.
// Pass the caller's intended ctx (it must carry a deadline: the CAS is a
// durable write); this API supplies no default timeout.
func Nominate(ctx context.Context, s *store.Store, head hashref.HashRef, release string, mode Mode) (Registry, hashref.HashRef, error) {
	return nominate(ctx, s, head, release, mode)
}

func nominate(ctx context.Context, s nominateStore, head hashref.HashRef, release string, mode Mode) (Registry, hashref.HashRef, error) {
	if head.IsZero() {
		return Registry{}, hashref.HashRef{}, fmt.Errorf("registry: nominate needs the current registry head; a fresh store is bootstrapped by the daemon, never nominated")
	}
	actual, ok, err := s.GetRegistryHead(ctx, SemanticID)
	if err != nil {
		return Registry{}, hashref.HashRef{}, err
	}
	if !ok || actual != head {
		return Registry{}, hashref.HashRef{}, &store.RegistryCASConflict{Name: SemanticID, Expected: head, Actual: actual, HadHead: ok}
	}
	cur, err := load(ctx, s, head)
	if err != nil {
		return Registry{}, hashref.HashRef{}, err
	}
	next, err := Next(cur, release, mode)
	if err != nil {
		return Registry{}, hashref.HashRef{}, err
	}
	obj, err := object(next, provenanceNominate)
	if err != nil {
		return Registry{}, hashref.HashRef{}, err
	}
	if err := s.PutObject(ctx, obj); err != nil {
		return Registry{}, hashref.HashRef{}, err
	}
	if err := s.CompareAndSetRegistryHead(ctx, SemanticID, head, obj.Hash); err != nil {
		return Registry{}, hashref.HashRef{}, err
	}
	return next, obj.Hash, nil
}
