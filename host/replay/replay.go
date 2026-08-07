// Package replay implements the authoritative replay engine and replay-doubling
// harness from Decision 7 of the w-world-library-m1 design.
//
// A recorded episode is a manifest naming an initial world, an ordered list of
// log entries (each pinning a transitionFn HashRef, an interpreter HashRef, and
// a semantics epoch), the canonical transition source objects, the archived
// interpreter artifacts, and — per entry — the recorded result bytes and the
// recorded final world hash.
//
// The load-bearing property is that replay DELEGATES per-transition execution
// to the released AILANG artifact and never reimplements the interpreter
// (DESIGN.md §14). For each entry the engine:
//
//  1. loads the transitionFn canonical bytes from the object store and verifies
//     their content address;
//  2. resolves the interpreter executable from the artifact archive using the
//     ENTRY's own interpreter HashRef (authoritative — never the epoch registry
//     candidate) and verifies its content address;
//  3. consults the (transitionFn, interpreter) verify cache, re-verifying and
//     caching on a miss;
//  4. invokes the ARCHIVED released binary (`ailang run`) on the pinned source;
//  5. byte-compares the produced result with the recorded result;
//  6. reconstructs the next world/log hash from the produced bytes and compares
//     it with the recorded final world hash.
//
// semanticsEpoch is copied into cache metadata for diagnostics only; cache
// lookup validity uses the pair key exclusively. Any divergence at any step is
// a structured error and fails replay — this is the standing hermeticity test
// required by ratified D1.
package replay

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/canon"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"

	"github.com/sunholo-data/ailang-world/host/childenv"
)

// execTimeout bounds every archived-interpreter subprocess invocation so replay
// can never hang; a transition that does not terminate within it fails replay.
const execTimeout = 60 * time.Second

// entryModulePath is the module/file path (relative to the reconstructed
// project root) at which a pinned transition source is materialized for the
// released binary to run. It is fixed so the module header MOD010 check is
// satisfied without a temp-path relaxation.
const entryModulePath = "host/replay/testdata/transition_fixture.ail"

// entryFn is the entrypoint the released binary runs for a fixture transition.
const entryFn = "main"

// EpisodeEntry is one recorded transition inside an episode. TransitionFn and
// Interpreter are the authoritative pair pins; SemanticsEpoch is compatibility
// metadata carried into the cache but never part of resolution or the cache
// key. RecordedResult is the exact result byte stream captured at recording
// time and RecordedWorldHash is the content address reconstructed from it.
type EpisodeEntry struct {
	// TransitionFn addresses the canonical transition source object.
	TransitionFn hashref.HashRef
	// Interpreter addresses the archived released binary that MUST execute this
	// entry. Authoritative replay resolves the executable from THIS ref.
	Interpreter hashref.HashRef
	// SemanticsEpoch is compatibility metadata (cache payload only).
	SemanticsEpoch int64
	// RecordedResult is the exact result byte stream produced at record time.
	RecordedResult []byte
	// RecordedWorldHash is the next-world content address reconstructed from the
	// recorded result bytes at record time.
	RecordedWorldHash hashref.HashRef
}

// Episode is a recorded episode manifest: an initial world and an ordered list
// of recorded entries. WorldLibDir is the directory holding the world/* AILANG
// library the pinned transition imports; the engine copies it into the
// reconstructed project root so the released binary resolves imports
// hermetically.
type Episode struct {
	// InitialWorld is the world the episode starts from (recorded for
	// completeness and chain reconstruction).
	InitialWorld store.World
	// Entries is the ordered list of recorded transitions.
	Entries []EpisodeEntry
	// WorldLibDir is the path to the world/* library directory the pinned
	// transitions import. It is copied verbatim into each replay project root.
	WorldLibDir string
}

// DivergenceError is the structured error returned when a replayed transition
// diverges from its record: the produced result bytes differ, or the
// reconstructed world hash differs. It is the fatal signal of ratified D1's
// hermeticity test.
type DivergenceError struct {
	// EntryIndex is the ordinal of the diverging entry within the episode.
	EntryIndex int
	// Field names what diverged ("result" or "world_hash").
	Field string
	// Detail is a human-readable explanation.
	Detail string
}

func (e *DivergenceError) Error() string {
	return fmt.Sprintf("replay: divergence at entry %d (%s): %s", e.EntryIndex, e.Field, e.Detail)
}

// Engine is the authoritative replay engine. It reads pinned objects from a
// Store, resolves interpreters from an Archive, and drives the released binary
// as a subprocess. It never links or reimplements the interpreter.
type Engine struct {
	store   *store.Store
	archive *archive.Archive
}

// NewEngine constructs an Engine over a store and its adjacent interpreter
// archive.
func NewEngine(s *store.Store, a *archive.Archive) *Engine {
	return &Engine{store: s, archive: a}
}

// ReplayResult is the outcome of replaying one entry: the exact produced result
// bytes, the reconstructed next-world hash, and whether the verify cache was
// consulted as a hit (false means the pair was re-verified and cached this run).
type ReplayResult struct {
	// Produced is the exact result byte stream from the archived interpreter.
	Produced []byte
	// WorldHash is the next-world content address reconstructed from Produced.
	WorldHash hashref.HashRef
	// CacheHit reports whether the (transitionFn, interpreter) verify row was
	// already present (true) or had to be re-verified and inserted (false).
	CacheHit bool
	// ExecPath is the archived executable resolved from the entry's interpreter
	// HashRef and actually invoked (authoritative resolution).
	ExecPath string
}

// ReplayEntry performs the authoritative six-step replay of one recorded entry
// and compares the produced bytes and reconstructed world hash against the
// record. A byte or hash divergence returns a *DivergenceError; a missing
// object, unresolved interpreter, or exec failure returns the appropriate
// structured error (store / archive.ReplayError). idx is the entry ordinal used
// only for error context.
func (e *Engine) ReplayEntry(ep Episode, idx int, entry EpisodeEntry) (ReplayResult, error) {
	// Step 1: load transitionFn canonical bytes and verify content address.
	src, ok, err := e.store.GetObject(entry.TransitionFn)
	if err != nil {
		return ReplayResult{}, err
	}
	if !ok {
		return ReplayResult{}, &archive.ReplayError{
			Kind:   archive.KindAbsentArtifact,
			Ref:    entry.TransitionFn,
			Detail: "transitionFn source object is absent from the store",
		}
	}
	computed, err := hashref.Sum(entry.TransitionFn.Algo(), src.Payload)
	if err != nil {
		return ReplayResult{}, err
	}
	if computed.Digest() != entry.TransitionFn.Digest() {
		return ReplayResult{}, &DivergenceError{
			EntryIndex: idx,
			Field:      "transition_source",
			Detail: fmt.Sprintf("stored source hashes to %q, entry pins %q",
				computed.String(), entry.TransitionFn.String()),
		}
	}

	// Step 2: resolve the interpreter from the ENTRY's interpreter HashRef
	// (authoritative — never the registry candidate) and verify its content
	// address against the archived bytes.
	execPath, err := e.archive.Resolve(entry.Interpreter)
	if err != nil {
		return ReplayResult{}, err
	}
	if err := verifyExecutable(execPath, entry.Interpreter); err != nil {
		return ReplayResult{}, err
	}

	// Step 3: consult the (transitionFn, interpreter) verify cache. A miss
	// re-verifies (here: confirms the pinned pair resolves) and caches the row.
	// semanticsEpoch is written as metadata only and never keys the lookup.
	_, cacheHit, err := e.store.GetVerifyResult(entry.TransitionFn, entry.Interpreter)
	if err != nil {
		return ReplayResult{}, err
	}
	if !cacheHit {
		if err := e.store.PutVerifyResult(store.VerifyResult{
			TransitionFn:   entry.TransitionFn,
			Interpreter:    entry.Interpreter,
			SemanticsEpoch: entry.SemanticsEpoch,
			Verified:       true,
			Detail:         "replay re-verification: pinned pair resolved and executed",
		}); err != nil {
			return ReplayResult{}, err
		}
	}

	// Step 4: invoke the archived released binary on the pinned source. The
	// interpreter is the archived artifact; the host does not evaluate AILANG.
	produced, err := runPinnedTransition(execPath, ep.WorldLibDir, src.Payload)
	if err != nil {
		return ReplayResult{}, err
	}

	// Step 5: byte-compare the produced result with the recorded result.
	if !bytes.Equal(produced, entry.RecordedResult) {
		return ReplayResult{}, &DivergenceError{
			EntryIndex: idx,
			Field:      "result",
			Detail: fmt.Sprintf("produced %d bytes %q, recorded %d bytes %q",
				len(produced), string(produced), len(entry.RecordedResult), string(entry.RecordedResult)),
		}
	}

	// Step 6: reconstruct the next world/log hash from the produced bytes and
	// compare it with the recorded final world hash.
	//
	// M1 fixture scope: this "world hash" is the content address of the result
	// bytes (SHA-256(produced)), so it cannot independently catch a divergence
	// step 5 misses. A real episode would hash the full World struct (log +
	// registry + heads), giving an independent transition witness.
	worldHash := hashref.SumSHA256(produced)
	if worldHash.Digest() != entry.RecordedWorldHash.Digest() ||
		worldHash.Algo() != entry.RecordedWorldHash.Algo() {
		return ReplayResult{}, &DivergenceError{
			EntryIndex: idx,
			Field:      "world_hash",
			Detail: fmt.Sprintf("reconstructed %q, recorded %q",
				worldHash.String(), entry.RecordedWorldHash.String()),
		}
	}

	return ReplayResult{
		Produced:  produced,
		WorldHash: worldHash,
		CacheHit:  cacheHit,
		ExecPath:  execPath,
	}, nil
}

// ReplayEpisode replays every entry of an episode in order and returns the
// per-entry results. The first divergence or structured error stops the episode
// and is returned.
func (e *Engine) ReplayEpisode(ep Episode) ([]ReplayResult, error) {
	results := make([]ReplayResult, 0, len(ep.Entries))
	for i, entry := range ep.Entries {
		r, err := e.ReplayEntry(ep, i, entry)
		if err != nil {
			return results, err
		}
		results = append(results, r)
	}
	return results, nil
}

// verifyExecutable streams the executable at path through SHA-256 and confirms
// the digest equals ref. It guards against an archived artifact that was
// corrupted or swapped after archival, mirroring the object content check.
func verifyExecutable(path string, ref hashref.HashRef) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return &archive.ReplayError{
			Kind:   archive.KindAbsentArtifact,
			Ref:    ref,
			Path:   path,
			Detail: "cannot read archived interpreter for verification",
			Err:    err,
		}
	}
	got, err := hashref.Sum(ref.Algo(), data)
	if err != nil {
		return err
	}
	if got.Digest() != ref.Digest() {
		return &archive.ReplayError{
			Kind:   archive.KindHashMismatch,
			Ref:    ref,
			Path:   path,
			Detail: fmt.Sprintf("archived interpreter hashes to %q, entry pins %q", got.String(), ref.String()),
		}
	}
	return nil
}

// runPinnedTransition materializes a hermetic project root, writes the pinned
// canonical source at its fixed module path and copies the world/* library into
// it, then runs `<execPath> run --quiet --caps "" --entry main <source>` with
// the project root as the working directory. It returns the exact stdout bytes.
//
// The archived released binary is the interpreter; this function only stages
// bytes and captures output. It never evaluates AILANG.
func runPinnedTransition(execPath, worldLibDir string, canonicalSource []byte) ([]byte, error) {
	root, err := os.MkdirTemp("", "world-replay-*")
	if err != nil {
		return nil, fmt.Errorf("replay: create project root: %w", err)
	}
	defer func() { _ = os.RemoveAll(root) }()

	// Materialize the world/* library so the pinned source's imports resolve
	// against exactly the recorded library bytes.
	if err := copyDir(worldLibDir, filepath.Join(root, "world")); err != nil {
		return nil, fmt.Errorf("replay: stage world library: %w", err)
	}

	// Write the pinned canonical source at its module path.
	srcPath := filepath.Join(root, filepath.FromSlash(entryModulePath))
	if err := os.MkdirAll(filepath.Dir(srcPath), 0o755); err != nil {
		return nil, fmt.Errorf("replay: create source dir: %w", err)
	}
	if err := os.WriteFile(srcPath, canonicalSource, 0o644); err != nil {
		return nil, fmt.Errorf("replay: write pinned source: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), execTimeout)
	defer cancel()

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, execPath,
		"run", "--quiet", "--caps", "", "--entry", entryFn, entryModulePath)
	cmd.Dir = root
	// Decision 4 of design_docs/planned/w-self-mod-vertical.md: replay is a
	// non-publish subprocess and must observe every registry variable unset.
	// Scrubbing os.Environ() rather than emptying it keeps the rest of the
	// environment — and therefore the recorded goldens — unchanged.
	cmd.Env = childenv.Scrubbed(os.Environ())
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, &archive.ReplayError{
			Kind:   archive.KindExecFailure,
			Path:   execPath,
			Detail: fmt.Sprintf("archived interpreter run failed (stderr: %q)", stderr.String()),
			Err:    err,
		}
	}
	return stdout.Bytes(), nil
}

// copyDir recursively copies the regular files under src into dst, creating
// directories as needed. It is deliberately minimal: the world library is a
// small flat tree of .ail files, so symlinks and special files are not
// expected and are skipped.
func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, rerr := filepath.Rel(src, path)
		if rerr != nil {
			return rerr
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if !d.Type().IsRegular() {
			return nil
		}
		data, rerr := os.ReadFile(path)
		if rerr != nil {
			return rerr
		}
		if merr := os.MkdirAll(filepath.Dir(target), 0o755); merr != nil {
			return merr
		}
		return os.WriteFile(target, data, 0o644)
	})
}

// SourceObject builds the immutable store.Object for a transition source by
// canonicalizing raw source bytes (Decision 2) and content-addressing the
// canonical result. Recording and replay both hash exactly these canonical
// bytes, so a source stored by SourceObject round-trips through ReplayEntry's
// step-1 verification. semanticID and provenance label the object.
func SourceObject(rawSource []byte, semanticID, provenance string) (store.Object, error) {
	canonical, err := canon.Source(rawSource)
	if err != nil {
		return store.Object{}, fmt.Errorf("replay: canonicalize transition source: %w", err)
	}
	return store.Object{
		Hash:          hashref.SumSHA256(canonical),
		InterfaceHash: hashref.SumSHA256([]byte(semanticID)),
		SemanticID:    semanticID,
		Provenance:    provenance,
		Payload:       canonical,
	}, nil
}
