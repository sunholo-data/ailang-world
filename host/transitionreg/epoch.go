package transitionreg

import (
	"context"
	"fmt"
	"strings"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/registry"
)

// EpochRegistryAbsentError reports a store with no epoch-registry head: the
// semantics epoch cannot be derived or validated, and this publisher never
// defaults it. The daemon owns epoch-1 bootstrapping (daemon.go:512).
type EpochRegistryAbsentError struct{}

func (e *EpochRegistryAbsentError) Error() string {
	return "publish transition registry: no epoch registry head exists (world/epoch-registry/v1); " +
		"run the daemon once so it bootstraps epoch 1 from the interpreter it serves, then retry"
}

// InterpreterEpochMismatchError reports a descriptor whose SemanticsEpoch is
// not one of the epochs the registry nominates for the pinned interpreter's
// release string: an unknown pair (Nominating empty) or a mismatched pair.
type InterpreterEpochMismatchError struct {
	ID         string
	Epoch      int64
	Release    string
	Nominating []int64
}

func (e *InterpreterEpochMismatchError) Error() string {
	if len(e.Nominating) == 0 {
		return fmt.Sprintf(
			"publish transition registry: entry %q pins semantics epoch %d, but no epoch in world/epoch-registry/v1 nominates interpreter release %q",
			e.ID, e.Epoch, e.Release)
	}
	return fmt.Sprintf(
		"publish transition registry: entry %q pins semantics epoch %d, but interpreter release %q is nominated by epochs %v only",
		e.ID, e.Epoch, e.Release, e.Nominating)
}

// releaseFromManifest reduces an archived interpreter manifest's verbatim
// `--version` output to the release string EXACTLY as the daemon's bootstrap
// does (host/daemon releaseFromVersion, daemon.go:613-624): the first
// non-blank line, trimmed; "unpinned" when there is none. Any other reduction
// false-mismatches the registry's Candidates.
func releaseFromManifest(version string) string {
	for _, line := range strings.Split(version, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			return trimmed
		}
	}
	return "unpinned"
}

// EpochsForInterpreter returns the semantics epochs whose Candidates name the
// pinned interpreter's release string (derived from the archived manifest),
// plus that release. Zero epochs means the registry does not nominate this
// interpreter release at all; several means only an explicit epoch statement
// can disambiguate.
func (r *StoreReader) EpochsForInterpreter(ctx context.Context, interpreter hashref.HashRef) ([]int64, string, error) {
	if r.arch == nil {
		return nil, "", &PublisherArchiveRequiredError{}
	}
	reg, err := r.epochRegistry(ctx)
	if err != nil {
		return nil, "", err
	}
	m, err := r.arch.ReadManifest(interpreter)
	if err != nil {
		return nil, "", fmt.Errorf("publish transition registry: read interpreter manifest %q: %w", interpreter.String(), err)
	}
	release := releaseFromManifest(m.Version)
	var epochs []int64
	for _, rec := range reg.Epochs {
		for _, candidate := range rec.Candidates {
			if candidate == release {
				epochs = append(epochs, rec.Epoch)
				break
			}
		}
	}
	return epochs, release, nil
}

// epochRegistry reads and decodes the epoch-registry head. It refuses
// (typed) when the head is absent — no default epoch is ever derived from
// absence.
func (r *StoreReader) epochRegistry(ctx context.Context) (registry.Registry, error) {
	head, ok, err := r.store.GetRegistryHead(ctx, registry.SemanticID)
	if err != nil {
		return registry.Registry{}, fmt.Errorf("publish transition registry: read epoch registry head: %w", err)
	}
	if !ok {
		return registry.Registry{}, &EpochRegistryAbsentError{}
	}
	object, found, err := r.store.GetObject(ctx, head)
	if err != nil {
		return registry.Registry{}, fmt.Errorf("publish transition registry: read epoch registry object: %w", err)
	}
	if !found {
		return registry.Registry{}, fmt.Errorf("publish transition registry: epoch registry head %q names an absent object", head.String())
	}
	reg, err := registry.Decode(object.Payload)
	if err != nil {
		return registry.Registry{}, fmt.Errorf("publish transition registry: decode epoch registry: %w", err)
	}
	return reg, nil
}

// verifyEpochs is the objection-A gate inside PublishSet: every descriptor's
// SemanticsEpoch must be one of the epochs the registry nominates for the
// descriptor's own pinned interpreter. A non-CLI caller cannot skip it — the
// same reader that publishes performs it, and a mismatched or defaulted
// epoch is a typed refusal, never a silent rewrite.
func (r *StoreReader) verifyEpochs(ctx context.Context, changes []Change) error {
	releases := make(map[hashref.HashRef]string)
	nominating := make(map[hashref.HashRef][]int64)
	for _, change := range changes {
		if change.Descriptor == nil {
			continue
		}
		d := change.Descriptor
		epochs, cached := nominating[d.Interpreter]
		if !cached {
			var err error
			epochs, releases[d.Interpreter], err = r.EpochsForInterpreter(ctx, d.Interpreter)
			if err != nil {
				return err
			}
			nominating[d.Interpreter] = epochs
		}
		if !epochIn(epochs, d.SemanticsEpoch) {
			return &InterpreterEpochMismatchError{
				ID: change.ID, Epoch: d.SemanticsEpoch,
				Release: releases[d.Interpreter], Nominating: epochs,
			}
		}
	}
	return nil
}

func epochIn(epochs []int64, want int64) bool {
	for _, e := range epochs {
		if e == want {
			return true
		}
	}
	return false
}
