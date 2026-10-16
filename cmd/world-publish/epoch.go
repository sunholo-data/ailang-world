package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/pins"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// Fences of `epoch nominate` (row 139 D2). The verb's refusals are STOP lines
// like every other fence here: nothing was written.
const (
	fenceRelease  = "release"
	fenceVerdict  = "verdict"
	fenceBinary   = "binary"
	fenceRegistry = "registry"
)

// compiledManifestName is how a refusal names the compiled-in manifest.
const compiledManifestName = "pins/ailang-releases.json (compiled in)"

// hostPlatform is the platform whose manifest digest a verdict and --binary
// must match: the nomination is about the binary THIS host would serve.
var hostPlatform = runtime.GOOS + "/" + runtime.GOARCH

// epochNominateBudget bounds the verb's store reads, its version probe and
// its one CAS write, each separately.
const epochNominateBudget = 10 * time.Second

// runEpochNominate is `world-publish epoch nominate` (row 139 D2,
// design_docs/planned/w-interpreter-epoch-upgrade-path.md): append ONE
// immutable epoch-registry revision naming --release, either as a candidate
// of the current epoch or as a new epoch.
//
// It refuses unless --gate-verdict is a world/ailang-release-verdict/v1
// verdict that names the same release, was produced on this platform,
// carries the manifest's binary and tarball digests for it, passed every
// check, and (for --candidate) reproduced every replay golden
// byte-identically; and unless --binary is a file whose sha256 is the
// manifest's (and the verdict's) binary digest and whose `--version` line is
// the release's.
//
// THE TRUST ANCHOR is the compiled-in manifest (pins.Current), which a human
// reviewed when its row merged. There is deliberately no flag naming another
// manifest (row 139 M2 judge r1): a headless caller able to choose the
// manifest could write a row and a matching verdict for any release. The
// verdict is self-attested; --binary is what ties the nomination to real
// bytes the reviewed manifest names.
//
// THE AUTHORITY SPLIT (D-WORLD-78 = a, ratified 2026-10-16): --candidate is
// headless-permitted, so the loop can nominate a same-semantics release whose
// gate passed; --new-epoch (semantics moved) is the same attended fence stack
// as `transitions` (automation tripwire, controlling terminal, a typed phrase
// naming this write) and so is a human's act.
func runEpochNominate(opts options, in io.Reader, out, errw io.Writer, env environment) int {
	return runEpochNominateWith(opts, pins.Current(), compiledManifestName, in, out, errw, env)
}

// runEpochNominateWith is runEpochNominate against a given manifest. In
// production only runEpochNominate calls it, with the compiled-in manifest;
// tests call it directly with a test-only manifest. No argv reaches the
// manifest argument.
func runEpochNominateWith(opts options, manifest pins.Manifest, manifestName string, in io.Reader, out, errw io.Writer, env environment) int {
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
	rel, ok := manifest.Lookup(opts.release)
	if !ok {
		return report(errw, &stopError{Fence: fenceRelease, Reason: "unknown",
			Detail: fmt.Sprintf("%s is not in the release manifest %s (releases: %v); a release is nominated only after its manifest row lands",
				opts.release, manifestName, manifest.Names())})
	}
	verdict, serr := requireAdmittedVerdict(opts.gateVerdict, manifest, manifestName, opts.release, mode)
	if serr != nil {
		return report(errw, serr)
	}
	if serr := requireReleaseBinary(opts.binary, rel, verdict, manifestName); serr != nil {
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
// the version probe, the pre-gate reads and the write, so a human slow to
// type the --new-epoch phrase never spends the write's budget.
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

// requireAdmittedVerdict reads --gate-verdict and runs the nomination gate.
func requireAdmittedVerdict(path string, m pins.Manifest, manifestName, release string, mode transitionreg.EpochMode) (pins.Verdict, *stopError) {
	if strings.TrimSpace(path) == "" {
		return pins.Verdict{}, &stopError{Fence: fenceVerdict, Reason: "absent",
			Detail: "--gate-verdict names the release's compatibility-gate verdict (" + pins.VerdictSchema + "); no release is nominated without one"}
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return pins.Verdict{}, &stopError{Fence: fenceVerdict, Reason: "unreadable", Detail: err.Error()}
	}
	v, err := pins.ParseVerdict(b)
	if err == nil {
		err = m.AdmitVerdict(v, release, hostPlatform, manifestName, mode == transitionreg.EpochCandidate)
	}
	if err != nil {
		var ve *pins.VerdictError
		if errors.As(err, &ve) {
			return pins.Verdict{}, &stopError{Fence: fenceVerdict, Reason: ve.Reason, Detail: path + ": " + ve.Detail}
		}
		return pins.Verdict{}, &stopError{Fence: fenceVerdict, Reason: pins.VerdictMalformed, Detail: path + ": " + err.Error()}
	}
	return v, nil
}

// requireReleaseBinary hashes --binary and probes its version. The sha256
// must equal the manifest's binary digest for this platform AND the
// verdict's binary_sha256, and the first non-blank `--version` line must be
// the release's version line. A forged verdict alone therefore cannot
// nominate: the release's real bytes must match the reviewed manifest.
func requireReleaseBinary(path string, rel pins.Release, v pins.Verdict, manifestName string) *stopError {
	if strings.TrimSpace(path) == "" {
		return &stopError{Fence: fenceBinary, Reason: "absent",
			Detail: "--binary names the release's ailang binary; it is hashed against the release manifest and probed for its version"}
	}
	// Resolve once, so the hash and the version probe name the same file: a bare name would
	// otherwise be opened relative to the cwd but executed via PATH (row 139 M2 judge r2).
	abs, aerr := filepath.Abs(path)
	if aerr != nil {
		return &stopError{Fence: fenceBinary, Reason: "unreadable", Detail: aerr.Error()}
	}
	path = abs
	f, err := os.Open(path)
	if err != nil {
		return &stopError{Fence: fenceBinary, Reason: "unreadable", Detail: err.Error()}
	}
	h := sha256.New()
	_, err = io.Copy(h, f)
	_ = f.Close()
	if err != nil {
		return &stopError{Fence: fenceBinary, Reason: "unreadable", Detail: path + ": " + err.Error()}
	}
	got := hex.EncodeToString(h.Sum(nil))
	asset, ok := rel.Assets[hostPlatform]
	if !ok || got != asset.BinarySHA256 || got != v.BinarySHA256 {
		return &stopError{Fence: fenceBinary, Reason: "digest-mismatch",
			Detail: fmt.Sprintf("%s has sha256 %s; the release manifest %s records %q for %s on %s, the verdict %s",
				path, got, manifestName, asset.BinarySHA256, rel.Release, hostPlatform, v.BinarySHA256)}
	}
	line, err := probeVersionLine(path)
	if err != nil {
		return &stopError{Fence: fenceBinary, Reason: "version-probe-failed", Detail: err.Error()}
	}
	if line != rel.Version {
		return &stopError{Fence: fenceBinary, Reason: "version-mismatch",
			Detail: fmt.Sprintf("%s --version prints %q, want %q", path, line, rel.Version)}
	}
	return nil
}

// probeVersionLine runs `<path> --version` through the archive's bounded,
// scrubbed version probe (World's one --version exec site; this command owns
// no subprocess site) and returns its first non-blank line, trimmed.
func probeVersionLine(path string) (string, error) {
	stdout, err := archive.ProbeVersion(path)
	if err != nil {
		return "", err
	}
	for _, l := range strings.Split(stdout, "\n") {
		if t := strings.TrimSpace(l); t != "" {
			return t, nil
		}
	}
	return "", fmt.Errorf("%s --version printed no version line", path)
}
