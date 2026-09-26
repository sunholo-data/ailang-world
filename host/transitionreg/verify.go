package transitionreg

import (
	"context"
	"errors"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

// CanonicalSchema is the codec's single canonicalizer, exported for publisher
// callers (CLI verbs, coordinator paths, tests) so descriptor schemas are
// produced by exactly the code that will re-validate them at encode time.
func CanonicalSchema(raw []byte) ([]byte, error) { return canonicalSchema(raw) }

// NewPublisher returns a StoreReader that MAY PUBLISH: it carries the
// interpreter archive rooted next to the store, which PublishSet requires for
// its epoch and source-loadability verification. NewReader stays read-only —
// a bare reader's PublishSet refuses (PublisherArchiveRequiredError), so no
// caller can reach the write path without the verification inputs.
func NewPublisher(objects ObjectStore, arch *archive.Archive) *StoreReader {
	return &StoreReader{store: objects, cache: make(map[hashref.HashRef]Snapshot), arch: arch}
}

// PublisherArchiveRequiredError reports a publish attempted through a bare
// NewReader handle: the epoch and loadability checks need the interpreter
// archive, so the write path refuses rather than skipping them.
type PublisherArchiveRequiredError struct{}

func (e *PublisherArchiveRequiredError) Error() string {
	return "publish transition registry: the publisher needs the interpreter archive " +
		"(construct it with NewPublisher(db, archive.New(dbPath))); NewReader is read-only"
}

// TransitionSourceAbsentError reports a descriptor pinning a transition
// source object that the store does not hold. Publishing it would list a
// skill whose program cannot be resolved by execution or replay.
type TransitionSourceAbsentError struct{ Ref hashref.HashRef }

func (e *TransitionSourceAbsentError) Error() string {
	return fmt.Sprintf("publish transition registry: transition source %q is absent", e.Ref.String())
}

// TransitionSourceInvalidError reports canonical source bytes the pinned
// archived interpreter refuses to load: `check` exited non-zero (parse error,
// type error, or an unresolvable import). ID names the manifest entry when
// known. This is the objection-B gate: canon.Source (normalisation only)
// accepted these bytes; the interpreter did not.
type TransitionSourceInvalidError struct {
	ID     string
	Ref    hashref.HashRef
	Output string
}

func (e *TransitionSourceInvalidError) Error() string {
	who := "transition source"
	if e.ID != "" {
		who = "entry " + fmt.Sprintf("%q", e.ID)
	}
	return fmt.Sprintf("publish transition registry: %s's source %q is not loadable under its pinned interpreter: %s",
		who, e.Ref.String(), e.Output)
}

// EnsureSourceLoadable proves the canonical source bytes load under the
// archived interpreter the descriptor pins: host/archive.CheckSource stages
// the bytes in a scratch root and runs `<archived interpreter> check <file>`
// bounded (procbound reservation + reap bound + wall clock, the probeVersion
// pattern). `check` parses and type-checks; it never executes the module.
// Any refusal — parse error, type error, or an unresolvable import — is the
// typed TransitionSourceInvalidError carrying the interpreter's captured
// output.
func EnsureSourceLoadable(ctx context.Context, arch *archive.Archive, interpreter hashref.HashRef, source []byte) error {
	if arch == nil {
		return &PublisherArchiveRequiredError{}
	}
	result, err := arch.CheckSource(ctx, interpreter, source)
	if err != nil {
		return fmt.Errorf("publish transition registry: %w", err)
	}
	if !result.Passed {
		return &TransitionSourceInvalidError{Ref: hashref.SumSHA256(source), Output: result.Output}
	}
	return nil
}

// verifySources is the store-level honesty gate: a descriptor whose
// TransitionFn object is absent from the store cannot be resolved by
// execution or replay, so it must not be published.
func verifySources(ctx context.Context, objects ObjectStore, changes []Change) error {
	for _, change := range changes {
		if change.Descriptor == nil {
			continue
		}
		if _, found, err := objects.GetObject(ctx, change.Descriptor.TransitionFn); err != nil {
			return fmt.Errorf("publish transition registry: verify transition source %q: %w", change.Descriptor.TransitionFn.String(), err)
		} else if !found {
			return &TransitionSourceAbsentError{Ref: change.Descriptor.TransitionFn}
		}
	}
	return nil
}

// verifyLoadable is the objection-B gate inside PublishSet: it reads every
// descriptor's stored source payload and proves it loads under the
// descriptor's own pinned archived interpreter — BEFORE any head move. The
// CLI runs the same check earlier (before PutObject); this arm makes the
// guarantee unskippable for every caller.
func (r *StoreReader) verifyLoadable(ctx context.Context, changes []Change) error {
	for _, change := range changes {
		if change.Descriptor == nil {
			continue
		}
		d := change.Descriptor
		object, found, err := r.store.GetObject(ctx, d.TransitionFn)
		if err != nil {
			return fmt.Errorf("publish transition registry: verify transition source %q: %w", d.TransitionFn.String(), err)
		}
		if !found {
			return &TransitionSourceAbsentError{Ref: d.TransitionFn}
		}
		if err := EnsureSourceLoadable(ctx, r.arch, d.Interpreter, object.Payload); err != nil {
			var invalid *TransitionSourceInvalidError
			if errors.As(err, &invalid) && invalid.ID == "" {
				invalid.ID = change.ID
			}
			return err
		}
	}
	return nil
}
