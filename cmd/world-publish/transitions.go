package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/canon"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// transitionSourceSemanticID names the object form a publisher stores for
// transition source bytes. The object convention follows store/journal's
// journalObject: InterfaceHash is the hash of the semantic ID itself.
const transitionSourceSemanticID = "world/transition-source/v1"

// manifestEntry is one descriptor in the --manifest file. The interpreter pin
// is deliberately NOT a manifest field: it is derived from the verified
// archive pin (--ailang-bin / --interpreter-ref), so a descriptor can never
// name an interpreter nobody verified. semanticsEpoch is OPTIONAL (revision
// 2): omitted, it is derived from the epoch registry when exactly one epoch
// nominates the interpreter's release, and refused as ambiguous otherwise —
// it is never defaulted.
type manifestEntry struct {
	ID               string            `json:"id"`
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	TransitionFn     string            `json:"transitionFn"`
	TransitionFnFile string            `json:"transitionFnFile"`
	SemanticsEpoch   int64             `json:"semanticsEpoch"`
	InputSchema      json.RawMessage   `json:"inputSchema"`
	OutputSchema     json.RawMessage   `json:"outputSchema"`
	Access           requirementWire   `json:"access"`
	DeclaredEffects  []requirementWire `json:"declaredEffects"`
}

type requirementWire struct {
	Effect string `json:"effect"`
	Scope  string `json:"scope"`
	Cost   int64  `json:"cost"`
}

// runTransitions is the local, attended, single-writer publish of the
// transition registry (queue row 107). It is NOT the network package publish:
// --live is refused, no credential, no registry origin. The write is the
// landed BuildNext/Publish CAS via transitionreg.PublishSet.
func runTransitions(opts options, in io.Reader, out, errw io.Writer, env environment) int {
	if serr := refuseLiveOnReadOnlyVerb("transitions", opts.live); serr != nil {
		return report(errw, serr)
	}
	if serr := requireStorePath(opts.store); serr != nil {
		return report(errw, serr)
	}
	if strings.TrimSpace(opts.manifest) == "" {
		fmt.Fprintln(errw, "world-publish transitions: --manifest names the descriptor manifest file")
		return exitUsage
	}

	// The store opens BEFORE the attended gate, like runPublish: an operator
	// who mistyped the path — or left the daemon holding writer authority —
	// learns it before being asked to type the confirmation phrase.
	db, err := store.Open(opts.store)
	if err != nil {
		if store.IsWriterAlreadyActive(err) {
			fmt.Fprintln(errw, "world-publish transitions: "+err.Error()+
				" — the transition registry is written under single-writer authority; stop the daemon first")
		}
		return report(errw, storeOpenFailed(opts.store, err))
	}
	defer func() { _ = db.Close() }()

	// THE GATE: the same human-in-the-loop fence as `approve` — refuse the
	// automation environment, require a controlling terminal, require the
	// typed confirmation. A headless loop must not be able to publish a
	// registry revision (clause 3: the write path is the operator's hands,
	// never an agent's).
	if serr := requireAttendedOperator(in, out, env.getenv, env.probe()); serr != nil {
		return report(errw, serr)
	}

	// One context roots every store call below (the Background root is the
	// pinned census row for runTransitions).
	ctx := context.Background()
	arch := archive.New(opts.store)
	interpreter, err := pinnedInterpreter(opts, arch, errw)
	if err != nil {
		fmt.Fprintln(errw, "world-publish transitions: "+err.Error())
		return exitError
	}
	entries, err := readManifest(opts.manifest, errw)
	if err != nil {
		fmt.Fprintln(errw, "world-publish transitions: "+err.Error())
		return exitError
	}
	changes, epochNote, err := buildChanges(ctx, db, arch, entries, interpreter)
	if err != nil {
		fmt.Fprintln(errw, "world-publish transitions: "+err.Error())
		return exitError
	}

	// The production publisher core: NewPublisher carries the interpreter
	// archive, so PublishSet re-derives and re-validates every descriptor's
	// epoch and loadability itself — the verb cannot skip those checks any
	// more than a non-CLI caller can.
	pub := transitionreg.NewPublisher(db, arch)
	res, err := pub.PublishSet(ctx, changes)
	if err != nil {
		fmt.Fprintln(errw, "world-publish transitions: "+publishSetErrorLine(err))
		return exitError
	}
	if epochNote != "" {
		fmt.Fprintln(out, epochNote)
	}
	if res.Unchanged {
		fmt.Fprintf(out, "transition registry UNCHANGED at revision %d (head %s); nothing written\n",
			res.Revision, res.Head)
		return exitOK
	}
	fmt.Fprintf(out, "published transition registry revision %d (head %s); the daemon's A2A card will list these skills on its next read\n",
		res.Revision, res.Head)
	return exitOK
}

// publishSetErrorLine renders a PublishSet failure for the operator. The
// same-ID CAS-retry conflict is rendered with its ID named and an explicit
// "nothing was written", so a raced publish can never silently drop either
// publisher's bytes.
func publishSetErrorLine(err error) string {
	var absent *transitionreg.TransitionSourceAbsentError
	if errors.As(err, &absent) {
		return absent.Error() + " — publish the transition source object first, or name it via transitionFnFile"
	}
	var conflict *transitionreg.SameIDConflictError
	if errors.As(err, &conflict) {
		return fmt.Sprintf("both this publish and the registry's revision %d set %q with different bytes; nothing was written — republish the winner's descriptor, or remove the entry",
			conflict.Revision, conflict.ID)
	}
	if store.IsRegistryCASConflict(err) {
		return "the registry head kept moving during publish; re-run the command"
	}
	return err.Error()
}

// pinnedInterpreter resolves the descriptor interpreter pin from the archive
// rooted next to the store. --ailang-bin performs the SAME archival the
// daemon performs at startup (so the pin names the binary the daemon will
// run); --interpreter-ref must name an archived binary and a readable
// manifest. Exactly one of the two is required — the pin is verified or
// refused, never assumed.
func pinnedInterpreter(opts options, arch *archive.Archive, errw io.Writer) (hashref.HashRef, error) {
	switch {
	case opts.ailangBin != "" && opts.interpreterRef != "":
		return hashref.HashRef{}, errors.New("give exactly one of --ailang-bin or --interpreter-ref")
	case opts.ailangBin != "":
		ref, err := arch.Archive(opts.ailangBin)
		if err != nil {
			return hashref.HashRef{}, fmt.Errorf("archive interpreter: %w", err)
		}
		if _, err := arch.ReadManifest(ref); err != nil {
			return hashref.HashRef{}, fmt.Errorf("archived interpreter %q has no readable manifest: %w", ref, err)
		}
		return ref, nil
	case opts.interpreterRef != "":
		ref, err := hashref.Parse(strings.TrimSpace(opts.interpreterRef))
		if err != nil {
			return hashref.HashRef{}, fmt.Errorf("--interpreter-ref: %w", err)
		}
		if _, err := arch.Resolve(ref); err != nil {
			return hashref.HashRef{}, fmt.Errorf("interpreter %q is not archived next to the store: %w", ref, err)
		}
		if _, err := arch.ReadManifest(ref); err != nil {
			return hashref.HashRef{}, fmt.Errorf("archived interpreter %q has no readable manifest: %w", ref, err)
		}
		return ref, nil
	default:
		return hashref.HashRef{}, errors.New(
			"the interpreter pin is required: --ailang-bin B archives B now (as the daemon does at startup), " +
				"or --interpreter-ref R names an interpreter already archived next to the store")
	}
}

func readManifest(path string, errw io.Writer) ([]manifestEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read manifest: %w", err)
	}
	var entries []manifestEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parse manifest %s: %w", path, err)
	}
	if len(entries) == 0 {
		return nil, errors.New("manifest holds no descriptors")
	}
	return entries, nil
}

// buildChanges turns manifest entries into Change values, filling every
// descriptor field HONESTLY: the semantics epoch is DERIVED from the epoch
// registry when the entry omits it (and the derivation is unique) and
// validated against the registry when the entry states it — never defaulted;
// transitionFn is either an existing store object (existence re-verified by
// PublishSet) or a source file canonicalised by host/canon, CHECKED for
// loadability under the pinned interpreter BEFORE it is stored, and stored as
// an object; the interpreter is the verified archive pin; schemas go through
// the codec's own canonicalizer.
func buildChanges(ctx context.Context, db *store.Store, arch *archive.Archive,
	entries []manifestEntry, interpreter hashref.HashRef) ([]transitionreg.Change, string, error) {
	epochs, release, err := transitionreg.NewPublisher(db, arch).EpochsForInterpreter(ctx, interpreter)
	if err != nil {
		return nil, "", err
	}
	if len(epochs) == 0 {
		return nil, "", fmt.Errorf(
			"interpreter release %q is nominated by NO epoch in world/epoch-registry/v1; the daemon bootstraps epoch 1 from the interpreter it serves — run the daemon with this binary once, or publish against a served store",
			release)
	}
	var epochNote string
	changes := make([]transitionreg.Change, 0, len(entries))
	for _, e := range entries {
		epoch := e.SemanticsEpoch
		if epoch == 0 {
			if len(epochs) > 1 {
				return nil, "", fmt.Errorf(
					"entry %q: semanticsEpoch is required — interpreter release %q is nominated by more than one epoch (%v); state one of them",
					e.ID, release, epochs)
			}
			epoch = epochs[0]
			epochNote = fmt.Sprintf("semantics epoch %d derived from world/epoch-registry/v1 for interpreter release %q", epoch, release)
		}
		if !epochInList(epochs, epoch) {
			return nil, "", fmt.Errorf(
				"entry %q: semanticsEpoch %d is not one of the epochs nominating interpreter release %q (%v)",
				e.ID, epoch, release, epochs)
		}
		fn, err := pinnedTransitionFn(ctx, db, arch, interpreter, e)
		if err != nil {
			return nil, "", fmt.Errorf("entry %q: %w", e.ID, err)
		}
		input, err := transitionreg.CanonicalSchema(e.InputSchema)
		if err != nil {
			return nil, "", fmt.Errorf("entry %q: input schema: %w", e.ID, err)
		}
		output, err := transitionreg.CanonicalSchema(e.OutputSchema)
		if err != nil {
			return nil, "", fmt.Errorf("entry %q: output schema: %w", e.ID, err)
		}
		d := transitionreg.Descriptor{
			ID: e.ID, TransitionFn: fn, Interpreter: interpreter,
			SemanticsEpoch: epoch, InputSchema: input, OutputSchema: output,
			Access:          requirementOf(e.Access),
			DeclaredEffects: requirementsOf(e.DeclaredEffects),
			Title:           e.Title, Description: e.Description,
		}
		changes = append(changes, transitionreg.Change{ID: e.ID, Descriptor: &d})
	}
	return changes, epochNote, nil
}

func epochInList(epochs []int64, want int64) bool {
	for _, e := range epochs {
		if e == want {
			return true
		}
	}
	return false
}

// pinnedTransitionFn resolves the transition source pin: exactly one of
// transitionFn (an existing object ref; PublishSet re-verifies existence) or
// transitionFnFile (canonical AILANG source bytes stored as an object). The
// file path is CHECKED FIRST: the canonical bytes must load under the pinned
// archived interpreter (a bounded `check` subprocess) BEFORE PutObject, so a
// broken source never even becomes a stored object.
func pinnedTransitionFn(ctx context.Context, db *store.Store, arch *archive.Archive,
	interpreter hashref.HashRef, e manifestEntry) (hashref.HashRef, error) {
	switch {
	case e.TransitionFn != "" && e.TransitionFnFile != "":
		return hashref.HashRef{}, errors.New("give exactly one of transitionFn or transitionFnFile")
	case e.TransitionFn != "":
		return hashref.Parse(strings.TrimSpace(e.TransitionFn))
	case e.TransitionFnFile != "":
		raw, err := os.ReadFile(e.TransitionFnFile)
		if err != nil {
			return hashref.HashRef{}, fmt.Errorf("read transition source: %w", err)
		}
		source, err := canon.Source(raw)
		if err != nil {
			return hashref.HashRef{}, fmt.Errorf("canonicalise transition source: %w", err)
		}
		if err := transitionreg.EnsureSourceLoadable(ctx, arch, interpreter, source); err != nil {
			var invalid *transitionreg.TransitionSourceInvalidError
			if errors.As(err, &invalid) {
				invalid.ID = e.ID
			}
			return hashref.HashRef{}, err
		}
		obj := store.Object{
			Hash: hashref.SumSHA256(source), InterfaceHash: hashref.SumSHA256([]byte(transitionSourceSemanticID)),
			SemanticID: transitionSourceSemanticID, Provenance: "cmd/world-publish", Payload: source,
		}
		if err := db.PutObject(obj); err != nil {
			return hashref.HashRef{}, fmt.Errorf("store transition source: %w", err)
		}
		return obj.Hash, nil
	default:
		return hashref.HashRef{}, errors.New("a transition source pin is required: transitionFn (existing object ref) or transitionFnFile (AILANG source file)")
	}
}

func requirementOf(r requirementWire) transitionreg.EffectRequirement {
	return transitionreg.EffectRequirement{Effect: r.Effect, Scope: r.Scope, Cost: r.Cost}
}

func requirementsOf(rs []requirementWire) []transitionreg.EffectRequirement {
	out := make([]transitionreg.EffectRequirement, len(rs))
	for i, r := range rs {
		out[i] = requirementOf(r)
	}
	return out
}
