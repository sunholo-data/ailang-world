package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/pins"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// Fences of `epoch nominate` (row 139 D2). The verb's refusals are STOP lines
// like every other fence here: nothing was written.
const (
	fenceRelease  = "release"
	fenceVerdict  = "verdict"
	fenceRegistry = "registry"
)

// compiledManifestName is how a refusal names the compiled-in manifest.
const compiledManifestName = "pins/ailang-releases.json (compiled in)"

// hostPlatform is the platform whose manifest digest a verdict must match:
// the verdict is about the binary THIS host would serve.
var hostPlatform = runtime.GOOS + "/" + runtime.GOARCH

// epochNominateBudget bounds the verb's store reads and its one CAS write.
const epochNominateBudget = 10 * time.Second

// runEpochNominate is `world-publish epoch nominate` (row 139 D2,
// design_docs/planned/w-interpreter-epoch-upgrade-path.md): append ONE
// immutable epoch-registry revision naming --release, either as a candidate
// of the current epoch or as a new epoch.
//
// It refuses unless --gate-verdict is a world/ailang-release-verdict/v1
// verdict that names the same release, was produced on this platform,
// carries the manifest's binary and tarball digests for it, passed every
// check, and (for --candidate) reproduced every replay golden byte-identically.
//
// THE AUTHORITY SPLIT (D-WORLD-78 = a, ratified 2026-10-16): --candidate is
// headless-permitted, so the loop can nominate a same-semantics release whose
// gate passed; --new-epoch (semantics moved) is the same attended fence stack
// as `transitions` (automation tripwire, controlling terminal, a typed phrase
// naming this write) and so is a human's act.
func runEpochNominate(opts options, in io.Reader, out, errw io.Writer, env environment) int {
	if serr := refuseLiveOnReadOnlyVerb("epoch nominate", opts.live); serr != nil {
		return report(errw, serr)
	}
	if serr := requireStorePath(opts.store); serr != nil {
		return report(errw, serr)
	}
	mode, serr := nominationMode(opts.candidate, opts.newEpoch)
	if serr != nil {
		return report(errw, serr)
	}
	if strings.TrimSpace(opts.release) == "" {
		return report(errw, &stopError{Fence: fenceRelease, Reason: "absent",
			Detail: "--release names the AILANG release to nominate (vX.Y.Z, a row of the release manifest)"})
	}

	manifest, manifestName, err := loadReleaseManifest(opts.releaseManifest)
	if err != nil {
		fmt.Fprintln(errw, "world-publish epoch nominate: "+err.Error())
		return exitError
	}
	rel, ok := manifest.Lookup(opts.release)
	if !ok {
		return report(errw, &stopError{Fence: fenceRelease, Reason: "unknown",
			Detail: fmt.Sprintf("%s is not in the release manifest %s (releases: %v); a release is nominated only after its manifest row lands",
				opts.release, manifestName, manifest.Names())})
	}
	if serr := requireAdmittedVerdict(opts.gateVerdict, manifest, manifestName, opts.release, mode); serr != nil {
		return report(errw, serr)
	}

	// The store opens BEFORE the attended gate, like transitions: a mistyped
	// path, a running daemon or a refused nomination is learned before a
	// human is asked to type anything.
	db, err := store.Open(opts.store)
	if err != nil {
		if store.IsWriterAlreadyActive(err) {
			fmt.Fprintln(errw, "world-publish epoch nominate: "+err.Error()+
				" — the epoch registry is written under single-writer authority; stop the daemon first")
		}
		return report(errw, storeOpenFailed(opts.store, err))
	}
	defer func() { _ = db.Close() }()

	readCtx, cancelRead := epochNominateContext()
	defer cancelRead()
	cur, head, found, err := transitionreg.LoadEpochRegistry(readCtx, db)
	if err != nil {
		fmt.Fprintln(errw, "world-publish epoch nominate: read epoch registry: "+err.Error())
		return exitError
	}
	if !found {
		return report(errw, &stopError{Fence: fenceRegistry, Reason: "absent",
			Detail: "the store has no " + transitionreg.EpochRegistrySemanticID + " head; the daemon bootstraps epoch 1 from the interpreter it serves — start it once on this store, then nominate"})
	}
	// The preview is the same pure Next that Nominate applies, so an
	// already-nominated release is refused before the attended gate.
	preview, err := transitionreg.NextEpochRegistry(cur, rel.Version, mode)
	if err != nil {
		var already *transitionreg.EpochAlreadyNominatedError
		if errors.As(err, &already) {
			return report(errw, &stopError{Fence: fenceRegistry, Reason: "already-nominated", Detail: err.Error()})
		}
		fmt.Fprintln(errw, "world-publish epoch nominate: "+err.Error())
		return exitError
	}

	if mode == transitionreg.EpochNew {
		opened := preview.Epochs[len(preview.Epochs)-1].Epoch
		if serr := requireAttendedLocalWrite(in, out, env.getenv, env.probe(),
			fmt.Sprintf("Type exactly, to open a NEW semantics epoch %d (replay goldens differ) for %s in this store's epoch registry:", opened, rel.Version),
			newEpochPhrase(opened, opts.release, opts.store)); serr != nil {
			return report(errw, serr)
		}
	}

	writeCtx, cancelWrite := epochNominateContext()
	defer cancelWrite()
	next, nextHead, err := transitionreg.NominateEpoch(writeCtx, db, head, rel.Version, mode)
	if err != nil {
		if store.IsRegistryCASConflict(err) {
			fmt.Fprintln(errw, "world-publish epoch nominate: the epoch-registry head moved since it was read ("+err.Error()+"); nothing was written — re-run the command")
			return exitError
		}
		fmt.Fprintln(errw, "world-publish epoch nominate: "+err.Error())
		return exitError
	}
	fmt.Fprintf(out, "nominated %s (%s) in %s\n", rel.Version, mode, transitionreg.EpochRegistrySemanticID)
	fmt.Fprintf(out, "  registry revision head %s (previous %s)\n", nextHead, head)
	for _, rec := range next.Epochs {
		fmt.Fprintf(out, "  epoch %d candidates %q\n", rec.Epoch, rec.Candidates)
	}
	return exitOK
}

// epochNominateContext is the verb's one context root: a fresh budget for
// the pre-gate reads and another for the write, so a human slow to type the
// --new-epoch phrase never spends the write's budget.
func epochNominateContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), epochNominateBudget)
}

// nominationMode refuses an ambiguous intent: exactly one of --candidate and
// --new-epoch, never a default.
func nominationMode(candidate, newEpoch bool) (transitionreg.EpochMode, *stopError) {
	switch {
	case candidate && newEpoch:
		return "", &stopError{Fence: fenceMode, Reason: "both",
			Detail: "--candidate and --new-epoch are mutually exclusive"}
	case candidate:
		return transitionreg.EpochCandidate, nil
	case newEpoch:
		return transitionreg.EpochNew, nil
	default:
		return "", &stopError{Fence: fenceMode, Reason: "none",
			Detail: "epoch nominate requires exactly one of --candidate (goldens byte-identical) or --new-epoch (semantics moved; attended)"}
	}
}

// newEpochPhrase names the epoch, the release and the store, so it cannot be
// reused for another nomination.
func newEpochPhrase(epoch int64, release, storePath string) string {
	return fmt.Sprintf("open epoch %d for %s in %s", epoch, release, storePath)
}

// loadReleaseManifest returns the compiled-in manifest, or the one at path
// validated by the same pins.Parse.
func loadReleaseManifest(path string) (pins.Manifest, string, error) {
	if strings.TrimSpace(path) == "" {
		return pins.Current(), compiledManifestName, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return pins.Manifest{}, "", fmt.Errorf("read release manifest: %w", err)
	}
	m, err := pins.Parse(b)
	if err != nil {
		return pins.Manifest{}, "", fmt.Errorf("release manifest %s: %w", path, err)
	}
	return m, path, nil
}

// requireAdmittedVerdict reads --gate-verdict and runs the nomination gate.
func requireAdmittedVerdict(path string, m pins.Manifest, manifestName, release string, mode transitionreg.EpochMode) *stopError {
	if strings.TrimSpace(path) == "" {
		return &stopError{Fence: fenceVerdict, Reason: "absent",
			Detail: "--gate-verdict names the release's compatibility-gate verdict (" + pins.VerdictSchema + "); no release is nominated without one"}
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return &stopError{Fence: fenceVerdict, Reason: "unreadable", Detail: err.Error()}
	}
	v, err := pins.ParseVerdict(b)
	if err == nil {
		err = m.AdmitVerdict(v, release, hostPlatform, manifestName, mode == transitionreg.EpochCandidate)
	}
	if err != nil {
		var ve *pins.VerdictError
		if errors.As(err, &ve) {
			return &stopError{Fence: fenceVerdict, Reason: ve.Reason, Detail: path + ": " + ve.Detail}
		}
		return &stopError{Fence: fenceVerdict, Reason: pins.VerdictMalformed, Detail: path + ": " + err.Error()}
	}
	return nil
}
