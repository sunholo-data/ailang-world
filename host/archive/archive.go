// Package archive implements the Interpreter Artifact Archive from Decision 6 of
// the w-world-library-m1 design: the mechanism that pins the exact AILANG
// interpreter bytes used to write a world's log, so authoritative replay can
// resolve a log entry's interpreter HashRef back to a byte-identical executable.
//
// For a store database at "<store>.db", artifacts live in an adjacent tree:
//
//	<store>.db.artifacts/interpreters/<algo>/<digest>/ailang
//
// Alongside each executable a JSON sidecar manifest records the hash, the
// interpreter's `--version` STDOUT, the byte size, and the writing OS and
// architecture. The tree travels with database backups, giving the resolver a
// direct HashRef -> executable mapping.
//
// Archival is content-addressed and idempotent: re-archiving the same bytes is a
// no-op success, and any absent artifact, unsupported platform, hash mismatch,
// or exec failure produces a structured *ReplayError (mirroring
// store.ConflictError) so authoritative replay stops rather than proceeding on
// an unpinned or divergent binary.
package archive

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"

	"github.com/sunholo-data/ailang-world/host/childenv"
)

// executableName is the fixed leaf filename of every archived interpreter.
const executableName = "ailang"

// manifestName is the fixed leaf filename of the sidecar manifest written next
// to each archived executable.
const manifestName = "manifest.json"

// artifactsSuffix is appended to the store database path to locate the adjacent
// artifact tree: "<store>.db" -> "<store>.db.artifacts".
const artifactsSuffix = ".artifacts"

// archivedPerm is the permission set applied to an archived executable:
// read+execute for all, no write. Archived bytes are immutable by content
// address, so the file is made read-only after fsync and before rename.
const archivedPerm os.FileMode = 0o555

// probeTimeout bounds the ELAPSED TIME of one "<interpreter> --version"
// execution. The probe now also runs on the idempotent path (the conditional
// sidecar self-heal in Archive), which is reached on every daemon start, and
// nothing else bounds that sequence: daemon.New()'s startup block carries no
// deadline, and host/archive had zero bounded execution before this constant.
// 10 s is ~2 orders of magnitude above the measured probe (rc=0 in well under
// 1 s for 168 bytes of stdout) and far below the verify gate's enclosing test
// budgets (-timeout 8m race under a 600 s outer wait), so the labelled
// KindExecFailure always fires first.
//
// SCOPE OF THE BOUND (measured, not assumed): exec.CommandContext SIGKILLs the
// DIRECT child. For the archived ailang -- a single compiled binary with no
// forked children -- that closes the output pipes and Run() returns at the
// deadline. A hypothetical interpreter that forked a grandchild inheriting
// stdout could hold Wait() open for that grandchild's lifetime; closing that
// general case needs cmd.WaitDelay and is deliberately out of scope here.
const probeTimeout = 10 * time.Second

// versionPrefix is the well-formedness test for a recorded interpreter version
// string: the pinned interpreter's `--version` STDOUT begins "AILANG v". A
// stored value that fails this prefix is the signature of a sidecar written by
// a build that merged stderr into the version, and it is what the idempotent
// path's conditional self-heal re-probes. Deliberately a POSITIVE shape rather
// than a denylist of known pollutants: the Observatory warning is one instance
// of the stderr-chatter class, and matching it by name would leave the next
// instance undetected.
const versionPrefix = "AILANG v"

// ReplayError is the structured error that stops authoritative replay. It
// mirrors store.ConflictError's shape: a stable Kind plus context fields, an
// Error() method, and an errors.As target via IsReplayError. The archive returns
// it for every failure that must halt replay rather than silently proceed on an
// unpinned or divergent interpreter: an absent artifact, an unsupported
// platform, a hash mismatch against an existing sidecar/archive, or an exec
// failure while probing --version.
type ReplayError struct {
	// Kind is a stable machine-readable category (see the Kind* constants).
	Kind string
	// Ref is the interpreter HashRef in play, when known (zero otherwise).
	Ref hashref.HashRef
	// Path is the filesystem path in play, when relevant.
	Path string
	// Detail is a human-readable explanation.
	Detail string
	// Err is the underlying cause, when wrapping one (may be nil).
	Err error
}

// ReplayError kinds. These are stable identifiers callers may switch on.
const (
	// KindAbsentArtifact means the requested interpreter is not in the archive.
	KindAbsentArtifact = "absent_artifact"
	// KindUnsupportedPlatform means the source executable path is not a usable
	// regular file (absent, a directory, or otherwise not openable).
	KindUnsupportedPlatform = "unsupported_platform"
	// KindHashMismatch means freshly computed bytes disagree with an existing
	// archived artifact or its sidecar manifest.
	KindHashMismatch = "hash_mismatch"
	// KindExecFailure means invoking the interpreter (for --version) failed.
	KindExecFailure = "exec_failure"
)

func (e *ReplayError) Error() string {
	msg := fmt.Sprintf("archive: replay error [%s]: %s", e.Kind, e.Detail)
	if !e.Ref.IsZero() {
		msg += fmt.Sprintf(" (ref %q)", e.Ref.String())
	}
	if e.Path != "" {
		msg += fmt.Sprintf(" (path %q)", e.Path)
	}
	if e.Err != nil {
		msg += ": " + e.Err.Error()
	}
	return msg
}

// Unwrap exposes the wrapped cause for errors.Is/As chains.
func (e *ReplayError) Unwrap() error { return e.Err }

// IsReplayError reports whether err is (or wraps) a *ReplayError, and returns it.
// It mirrors store.IsConflict's As-based helper style.
func IsReplayError(err error) (*ReplayError, bool) {
	var r *ReplayError
	if errors.As(err, &r) {
		return r, true
	}
	return nil, false
}

// Manifest is the JSON sidecar recorded next to each archived executable. It
// captures the pinning facts a later replay/verify needs for diagnostics: the
// content address, the interpreter's own version string, and the writing
// platform. Size lets a resolver cheaply sanity-check the executable on disk.
type Manifest struct {
	// Hash is the canonical "algo:digest" HashRef of the archived bytes.
	Hash string `json:"hash"`
	// Version is the verbatim `ailang --version` output captured at archival.
	Version string `json:"version"`
	// Size is the archived executable's byte size.
	Size int64 `json:"size"`
	// OS is the writing host operating system (runtime.GOOS at archival).
	OS string `json:"os"`
	// Arch is the writing host architecture (runtime.GOARCH at archival).
	Arch string `json:"arch"`
}

// Archive is a handle to the interpreter artifact tree adjacent to a store
// database. Construct it with New; it holds only the resolved root directory and
// is safe to reuse for many resolutions.
type Archive struct {
	// root is "<store>.db.artifacts": the parent of the "interpreters" tree.
	root string

	// probeTimeout is the per-probe wall-clock bound, seeded from the
	// probeTimeout constant by New. It is a field rather than a direct constant
	// reference for the same reason daemon.readDeadline is: the wiring stays
	// assertable and the value is shrinkable in tests, which is the only way the
	// SHIPPED timeout branch can be exercised without a ten-second test.
	probeTimeout time.Duration
}

// New returns an Archive for the store database at storeDBPath. The artifact
// tree lives at "<storeDBPath>.artifacts"; New does not create it (Archive
// creates directories lazily during Archive). storeDBPath may be relative or
// absolute and need not yet exist.
func New(storeDBPath string) *Archive {
	return &Archive{root: storeDBPath + artifactsSuffix, probeTimeout: probeTimeout}
}

// Root returns the artifact-tree root directory ("<store>.db.artifacts").
func (a *Archive) Root() string { return a.root }

// dirFor returns the per-artifact directory for a HashRef:
// "<root>/interpreters/<algo>/<digest>".
func (a *Archive) dirFor(ref hashref.HashRef) string {
	return filepath.Join(a.root, "interpreters", ref.Algo(), ref.Digest())
}

// pathFor returns the archived executable path for a HashRef.
func (a *Archive) pathFor(ref hashref.HashRef) string {
	return filepath.Join(a.dirFor(ref), executableName)
}

// manifestPathFor returns the sidecar manifest path for a HashRef.
func (a *Archive) manifestPathFor(ref hashref.HashRef) string {
	return filepath.Join(a.dirFor(ref), manifestName)
}

// Archive performs the Decision 6 startup archival for the interpreter at
// execPath and returns its content-addressed HashRef.
//
// Mechanics (exactly per Decision 6):
//
//  1. Open the configured executable once.
//  2. Stream the opened bytes through SHA-256 WHILE copying to a temp file in
//     the destination directory (one pass over one opened byte stream).
//  3. Compare the computed digest with any pre-existing archive: identical bytes
//     are an idempotent no-op success; a digest that disagrees with an existing
//     sidecar is a KindHashMismatch ReplayError.
//  4. fsync the temp file, set it read-only executable (0o555), and atomically
//     os.Rename it into place.
//  5. Write the sidecar manifest (hash, --version STDOUT, size, OS, arch).
//
// Errors: an absent/irregular source is KindUnsupportedPlatform; a --version
// probe failure is KindExecFailure; a mismatch against an existing artifact is
// KindHashMismatch. All are *ReplayError.
func (a *Archive) Archive(execPath string) (hashref.HashRef, error) {
	// Step 1: open the configured executable exactly once. Reject anything that
	// is not a usable regular file as an unsupported platform/artifact.
	info, err := os.Stat(execPath)
	if err != nil {
		return hashref.HashRef{}, &ReplayError{
			Kind:   KindUnsupportedPlatform,
			Path:   execPath,
			Detail: "configured interpreter is not accessible",
			Err:    err,
		}
	}
	if !info.Mode().IsRegular() {
		return hashref.HashRef{}, &ReplayError{
			Kind:   KindUnsupportedPlatform,
			Path:   execPath,
			Detail: "configured interpreter is not a regular file",
		}
	}
	src, err := os.Open(execPath)
	if err != nil {
		return hashref.HashRef{}, &ReplayError{
			Kind:   KindUnsupportedPlatform,
			Path:   execPath,
			Detail: "cannot open configured interpreter",
			Err:    err,
		}
	}
	defer func() { _ = src.Close() }()

	// The destination directory must exist before we create a temp file in it,
	// so the atomic rename stays on the same filesystem.
	destDir := filepath.Join(a.root, "interpreters")
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return hashref.HashRef{}, fmt.Errorf("archive: create artifact dir %q: %w", destDir, err)
	}

	tmp, err := os.CreateTemp(destDir, ".ailang-*.tmp")
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf("archive: create temp file in %q: %w", destDir, err)
	}
	tmpPath := tmp.Name()
	// Best-effort cleanup: if we do not rename the temp away, remove it. After a
	// successful rename tmpPath no longer exists, so Remove is a harmless no-op.
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
	}()

	// Step 2: hash and copy consume the SAME opened byte stream in one pass. The
	// TeeReader feeds every byte read from src into the hasher as it is copied to
	// the temp file; src is never opened or read a second time.
	hasher := sha256.New()
	size, err := io.Copy(tmp, io.TeeReader(src, hasher))
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf("archive: stream %q: %w", execPath, err)
	}
	digest := hex.EncodeToString(hasher.Sum(nil))
	ref, err := hashref.New(hashref.AlgoSHA256, digest)
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf("archive: form hashref: %w", err)
	}

	// Step 3: idempotence and mismatch detection against any pre-existing
	// artifact for this exact content address.
	finalPath := a.pathFor(ref)
	if existing, err := a.readManifest(ref); err == nil {
		// A sidecar exists for this digest. Because the archive path is keyed by
		// the digest itself, a recorded hash that differs from the just-computed
		// one is corruption we must refuse rather than overwrite.
		if existing.Hash != ref.String() {
			return hashref.HashRef{}, &ReplayError{
				Kind:   KindHashMismatch,
				Ref:    ref,
				Path:   finalPath,
				Detail: fmt.Sprintf("existing sidecar records hash %q for this artifact slot", existing.Hash),
			}
		}
		// Identical bytes already archived: idempotent no-op success. The temp
		// file is discarded by the deferred cleanup.
		//
		// One exception, and it is CONDITIONAL on purpose: a sidecar written by
		// a build that merged the interpreter's stderr into its version string
		// is polluted forever, because this path never re-probes. Heal it -- but
		// ONLY when the stored string fails the well-formedness test, so a
		// healthy or already-healed artifact performs ZERO process executions
		// and this really is the no-op the comment above claims.
		//
		// The predicate is POSITIVE-SHAPE, not a denylist for "Observatory":
		// the log line is one instance of a class (any stderr chatter), and
		// matching the known pollutant would leave the next one undetected.
		if !strings.HasPrefix(existing.Version, versionPrefix) {
			fresh, err := a.probeVersion(finalPath)
			if err != nil {
				return hashref.HashRef{}, &ReplayError{
					Kind:   KindExecFailure,
					Ref:    ref,
					Path:   finalPath,
					Detail: "cannot obtain --version from archived interpreter while healing its sidecar",
					Err:    err,
				}
			}
			// The inner compare is NOT redundant with the outer guard: the outer
			// one buys the zero-execution guarantee, this one buys convergence.
			// --version stdout is deterministic for fixed binary bytes, so at
			// most one rewrite ever happens per artifact and there is no rewrite
			// churn once the probe agrees.
			if fresh != existing.Version {
				existing.Version = fresh
				if err := a.writeManifest(ref, existing); err != nil {
					return hashref.HashRef{}, err
				}
			}
		}
		return ref, nil
	}

	// Step 4: fsync, make read-only executable, and atomically rename into place.
	if err := tmp.Sync(); err != nil {
		return hashref.HashRef{}, fmt.Errorf("archive: fsync temp file: %w", err)
	}
	if err := tmp.Chmod(archivedPerm); err != nil {
		return hashref.HashRef{}, fmt.Errorf("archive: set executable perms: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return hashref.HashRef{}, fmt.Errorf("archive: close temp file: %w", err)
	}
	if err := os.MkdirAll(a.dirFor(ref), 0o755); err != nil {
		return hashref.HashRef{}, fmt.Errorf("archive: create artifact slot: %w", err)
	}
	if err := os.Rename(tmpPath, finalPath); err != nil {
		return hashref.HashRef{}, fmt.Errorf("archive: atomic rename into %q: %w", finalPath, err)
	}

	// Step 5: probe --version and write the sidecar manifest. The version probe
	// runs against the just-archived read-only executable so the recorded string
	// reflects the pinned bytes.
	version, err := a.probeVersion(finalPath)
	if err != nil {
		return hashref.HashRef{}, &ReplayError{
			Kind:   KindExecFailure,
			Ref:    ref,
			Path:   finalPath,
			Detail: "cannot obtain --version from archived interpreter",
			Err:    err,
		}
	}
	m := Manifest{
		Hash:    ref.String(),
		Version: version,
		Size:    size,
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
	}
	if err := a.writeManifest(ref, m); err != nil {
		return hashref.HashRef{}, err
	}
	return ref, nil
}

// Resolve maps an interpreter HashRef to the archived executable path, used by
// later replay/verify. It returns a KindAbsentArtifact *ReplayError when the
// archive holds no executable for ref.
func (a *Archive) Resolve(ref hashref.HashRef) (string, error) {
	if ref.IsZero() {
		return "", &ReplayError{
			Kind:   KindAbsentArtifact,
			Detail: "cannot resolve the zero interpreter HashRef",
		}
	}
	p := a.pathFor(ref)
	info, err := os.Stat(p)
	if err != nil || !info.Mode().IsRegular() {
		return "", &ReplayError{
			Kind:   KindAbsentArtifact,
			Ref:    ref,
			Path:   p,
			Detail: "no archived interpreter for this HashRef",
			Err:    err,
		}
	}
	return p, nil
}

// ReadManifest returns the sidecar manifest for ref. It returns a
// KindAbsentArtifact *ReplayError when no manifest is present.
func (a *Archive) ReadManifest(ref hashref.HashRef) (Manifest, error) {
	m, err := a.readManifest(ref)
	if err != nil {
		return Manifest{}, &ReplayError{
			Kind:   KindAbsentArtifact,
			Ref:    ref,
			Path:   a.manifestPathFor(ref),
			Detail: "no sidecar manifest for this HashRef",
			Err:    err,
		}
	}
	return m, nil
}

// readManifest reads and decodes the sidecar for ref, returning the raw error
// (used internally where a missing sidecar is an expected control-flow signal).
func (a *Archive) readManifest(ref hashref.HashRef) (Manifest, error) {
	data, err := os.ReadFile(a.manifestPathFor(ref))
	if err != nil {
		return Manifest{}, err
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return Manifest{}, err
	}
	return m, nil
}

// writeManifest writes the sidecar manifest for ref as indented JSON.
func (a *Archive) writeManifest(ref hashref.HashRef, m Manifest) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("archive: encode manifest: %w", err)
	}
	data = append(data, '\n')
	p := a.manifestPathFor(ref)
	if err := os.WriteFile(p, data, 0o644); err != nil {
		return fmt.Errorf("archive: write manifest %q: %w", p, err)
	}
	return nil
}

// probeVersion runs "<execPath> --version" under a.probeTimeout and returns its
// STDOUT verbatim (trailing whitespace preserved as emitted by the
// interpreter). stderr is captured separately and never joins the returned
// string: the recorded version must be a function of the pinned bytes alone,
// and an interpreter that writes unrelated chatter to stderr (this rig's
// released ailang logs an Observatory size warning there on every invocation)
// would otherwise stamp wall-clock-varying text into a content-addressed
// artifact's sidecar. Both streams are labelled in the error, so nothing is
// lost on the failure path.
//
// The caller wraps failures as a KindExecFailure ReplayError. On deadline
// expiry the error wraps context.DeadlineExceeded explicitly, and because
// ReplayError unwraps, errors.Is(err, context.DeadlineExceeded) stays true
// through the public error chain.
func (a *Archive) probeVersion(execPath string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.probeTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, execPath, "--version")
	// Decision 4 of design_docs/planned/w-self-mod-vertical.md: every
	// World-launched process other than the one live publish dispatch runs
	// under the stripped environment. Without this the version probe inherits
	// the process environment wholesale, so an ambient registry credential
	// would reach a child that has no use for it.
	cmd.Env = childenv.Scrubbed(os.Environ())
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("%s --version: timed out after %v: %w (stdout: %q, stderr: %q)",
				execPath, a.probeTimeout, context.DeadlineExceeded, stdout.String(), stderr.String())
		}
		return "", fmt.Errorf("%s --version: %w (stdout: %q, stderr: %q)",
			execPath, err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}
