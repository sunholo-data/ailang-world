package replay

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/registry"
	"github.com/sunholo-data/ailang-world/host/store"
)

// ---------------------------------------------------------------------------
// Fixture wiring
// ---------------------------------------------------------------------------

// pinnedBinary resolves the SHIPPED released AILANG binary the replay engine
// drives as a subprocess. Per the sprint verify profile it MUST be the clean
// released v0.30.0 artifact (AILANG_BIN=/tmp/ailang-v0300/ailang), never the
// -dirty dev build on PATH. Tests that need it skip cleanly when it is absent
// so the package still builds/vets everywhere; the mission gate always sets
// AILANG_BIN, so a skip in CI is itself a red flag surfaced by the runner.
func pinnedBinary(t *testing.T) string {
	t.Helper()
	bin := os.Getenv("AILANG_BIN")
	if bin == "" {
		t.Skip("AILANG_BIN not set; replay requires the pinned released ailang binary")
	}
	info, err := os.Stat(bin)
	if err != nil || !info.Mode().IsRegular() {
		t.Skipf("AILANG_BIN %q is not a usable executable: %v", bin, err)
	}
	return bin
}

// repoDir returns the directory of this test file, used to locate the sibling
// testdata/ tree and the repo's world/ library directory without depending on
// the process working directory.
func repoDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine test file location")
	}
	return filepath.Dir(file) // .../host/replay
}

// worldLibDir returns the repo's world/* AILANG library directory.
func worldLibDir(t *testing.T) string {
	t.Helper()
	// host/replay -> host -> repo root -> world
	root := filepath.Dir(filepath.Dir(repoDir(t)))
	dir := filepath.Join(root, "world")
	if _, err := os.Stat(filepath.Join(dir, "transitions.ail")); err != nil {
		t.Fatalf("world library not found at %q: %v", dir, err)
	}
	return dir
}

// testdataPath joins the sibling testdata directory.
func testdataPath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(repoDir(t), "testdata", name)
}

// readTestdata reads a committed testdata file or fails the test.
func readTestdata(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(testdataPath(t, name))
	if err != nil {
		t.Fatalf("read testdata %q: %v", name, err)
	}
	return b
}

// fixtureEnv is one fully wired replay environment over a FRESH temp store: the
// store, the interpreter archive, the archived interpreter HashRef, the stored
// transition source object, and a one-entry episode referencing the committed
// golden recorded result and world hash. Every acceptance run builds a fresh
// one so replay-doubling replays from clean state.
type fixtureEnv struct {
	store    *store.Store
	archive  *archive.Archive
	engine   *Engine
	interp   hashref.HashRef
	srcObj   store.Object
	episode  Episode
	recorded []byte
	worldRef hashref.HashRef
	dbPath   string
}

// newFixtureEnv builds a fresh temp store, archives the pinned binary, stores
// the canonical fixture source object, and assembles the episode from the
// committed golden result bytes + world hash.
func newFixtureEnv(t *testing.T) *fixtureEnv {
	t.Helper()
	bin := pinnedBinary(t)

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "world.db")
	s, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })

	a := archive.New(dbPath)
	interp, err := a.Archive(bin)
	if err != nil {
		t.Fatalf("%s", archive.AttributeFailure("archive pinned binary", err))
	}

	rawSource := readTestdata(t, "transition_fixture.ail")
	srcObj, err := SourceObject(rawSource, "world/replay/fixture-transition", "m5-fixture-episode")
	if err != nil {
		t.Fatalf("build source object: %v", err)
	}
	if err := s.PutObject(srcObj); err != nil {
		t.Fatalf("put source object: %v", err)
	}

	recorded := readTestdata(t, "recorded_result.bytes")
	worldRef, err := hashref.Parse(strings.TrimSpace(string(readTestdata(t, "recorded_world_hash.txt"))))
	if err != nil {
		t.Fatalf("parse recorded world hash: %v", err)
	}

	ep := Episode{
		InitialWorld: store.World{
			Ref:       hashref.SumSHA256([]byte("genesis-world")),
			Revision:  0,
			StateRoot: hashref.MustParse("sha256:" + strings.Repeat("a", 64)),
			LogHead:   hashref.MustParse("sha256:" + strings.Repeat("0", 64)),
		},
		Entries: []EpisodeEntry{{
			TransitionFn:      srcObj.Hash,
			Interpreter:       interp,
			SemanticsEpoch:    1,
			RecordedResult:    recorded,
			RecordedWorldHash: worldRef,
		}},
		WorldLibDir: worldLibDir(t),
	}

	return &fixtureEnv{
		store:    s,
		archive:  a,
		engine:   NewEngine(s, a),
		interp:   interp,
		srcObj:   srcObj,
		episode:  ep,
		recorded: recorded,
		worldRef: worldRef,
		dbPath:   dbPath,
	}
}

// ---------------------------------------------------------------------------
// Acceptance 1: recorded fixture episode replays bit-for-bit
// ---------------------------------------------------------------------------

func TestFixtureEpisodeReplaysBitForBit(t *testing.T) {
	env := newFixtureEnv(t)

	results, err := env.engine.ReplayEpisode(env.episode)
	if err != nil {
		t.Fatalf("replay episode: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 entry result, got %d", len(results))
	}
	r := results[0]
	if string(r.Produced) != string(env.recorded) {
		t.Fatalf("produced %q != recorded %q", string(r.Produced), string(env.recorded))
	}
	if r.WorldHash.String() != env.worldRef.String() {
		t.Fatalf("reconstructed world hash %q != recorded %q", r.WorldHash.String(), env.worldRef.String())
	}
}

// ---------------------------------------------------------------------------
// Acceptance 2: replay-doubling — A == B == recorded, divergence fails
// ---------------------------------------------------------------------------

func TestReplayDoubling(t *testing.T) {
	// Replay A and replay B each run from a FRESH temp store through the pinned
	// artifact. A must equal B must equal the recorded bytes; the same holds for
	// the reconstructed world hash. Any divergence fails the test.
	envA := newFixtureEnv(t)
	resA, err := envA.engine.ReplayEpisode(envA.episode)
	if err != nil {
		t.Fatalf("replay A: %v", err)
	}

	envB := newFixtureEnv(t)
	resB, err := envB.engine.ReplayEpisode(envB.episode)
	if err != nil {
		t.Fatalf("replay B: %v", err)
	}

	if len(resA) != 1 || len(resB) != 1 {
		t.Fatalf("expected 1 result each, got A=%d B=%d", len(resA), len(resB))
	}
	// A == B
	if string(resA[0].Produced) != string(resB[0].Produced) {
		t.Fatalf("replay A %q != replay B %q", string(resA[0].Produced), string(resB[0].Produced))
	}
	if resA[0].WorldHash.String() != resB[0].WorldHash.String() {
		t.Fatalf("replay A world hash %q != replay B %q", resA[0].WorldHash.String(), resB[0].WorldHash.String())
	}
	// A == recorded
	if string(resA[0].Produced) != string(envA.recorded) {
		t.Fatalf("replay A %q != recorded %q", string(resA[0].Produced), string(envA.recorded))
	}
	// B == recorded
	if string(resB[0].Produced) != string(envB.recorded) {
		t.Fatalf("replay B %q != recorded %q", string(resB[0].Produced), string(envB.recorded))
	}
}

// TestReplayDoublingDivergenceFails proves the harness actually FAILS on
// divergence: a recorded-result byte that does not match the produced bytes must
// surface as a *DivergenceError rather than passing silently.
func TestReplayDoublingDivergenceFails(t *testing.T) {
	env := newFixtureEnv(t)
	// Corrupt the recorded result so replay must diverge.
	corrupt := append([]byte(nil), env.recorded...)
	corrupt[0] ^= 0xFF
	env.episode.Entries[0].RecordedResult = corrupt

	_, err := env.engine.ReplayEpisode(env.episode)
	if err == nil {
		t.Fatal("expected divergence error, got nil")
	}
	var de *DivergenceError
	if !errors.As(err, &de) {
		t.Fatalf("expected *DivergenceError, got %T: %v", err, err)
	}
	if de.Field != "result" {
		t.Fatalf("expected result divergence, got field %q", de.Field)
	}
}

// TestReplayWorldHashDivergenceFails proves a mismatched recorded world hash
// (with matching bytes) still fails at step 6.
func TestReplayWorldHashDivergenceFails(t *testing.T) {
	env := newFixtureEnv(t)
	env.episode.Entries[0].RecordedWorldHash = hashref.SumSHA256([]byte("not the real world"))

	_, err := env.engine.ReplayEpisode(env.episode)
	var de *DivergenceError
	if !errors.As(err, &de) {
		t.Fatalf("expected *DivergenceError, got %T: %v", err, err)
	}
	if de.Field != "world_hash" {
		t.Fatalf("expected world_hash divergence, got field %q", de.Field)
	}
}

// ---------------------------------------------------------------------------
// Acceptance 3: authoritative replay resolves from the ENTRY interpreter hash
// ---------------------------------------------------------------------------

func TestAuthoritativeResolutionUsesEntryInterpreter(t *testing.T) {
	env := newFixtureEnv(t)
	results, err := env.engine.ReplayEpisode(env.episode)
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	// The executable actually invoked must be the archive path for the ENTRY's
	// interpreter HashRef, not any other artifact.
	wantPath, err := env.archive.Resolve(env.episode.Entries[0].Interpreter)
	if err != nil {
		t.Fatalf("resolve entry interpreter: %v", err)
	}
	if results[0].ExecPath != wantPath {
		t.Fatalf("invoked %q, expected the entry-pinned interpreter %q", results[0].ExecPath, wantPath)
	}
}

// TestUnknownEntryInterpreterFailsAbsent proves that when an entry pins an
// interpreter HashRef the archive does not hold, authoritative replay fails with
// a KindAbsentArtifact ReplayError — it never falls back to any other binary.
func TestUnknownEntryInterpreterFailsAbsent(t *testing.T) {
	env := newFixtureEnv(t)
	// A well-formed but unarchived interpreter ref.
	env.episode.Entries[0].Interpreter = hashref.SumSHA256([]byte("some other binary"))

	_, err := env.engine.ReplayEpisode(env.episode)
	re, ok := archive.IsReplayError(err)
	if !ok {
		t.Fatalf("expected *archive.ReplayError, got %T: %v", err, err)
	}
	if re.Kind != archive.KindAbsentArtifact {
		t.Fatalf("expected KindAbsentArtifact, got %q", re.Kind)
	}
}

// ---------------------------------------------------------------------------
// Acceptance 4: epoch-registry candidate changes cannot redirect replay
// ---------------------------------------------------------------------------

func TestEpochRegistryCandidateCannotRedirect(t *testing.T) {
	env := newFixtureEnv(t)

	// Bootstrap the epoch registry naming the pinned interpreter version, then
	// mutate the registry head to a DIFFERENT candidate revision. Authoritative
	// replay must ignore the registry entirely and still resolve + execute the
	// entry-pinned interpreter, producing the identical recorded bytes.
	if _, _, err := registry.Bootstrap(env.store, "AILANG v0.30.0 (pinned)"); err != nil {
		t.Fatalf("bootstrap registry: %v", err)
	}
	// Record what the interpreter-pinned replay produced BEFORE mutating.
	before, err := env.engine.ReplayEpisode(env.episode)
	if err != nil {
		t.Fatalf("replay before registry mutation: %v", err)
	}

	// Store a divergent registry revision and repoint the head at it.
	rogue := registry.Registry{
		SemanticID: registry.SemanticID,
		Epochs: []registry.EpochRecord{
			{Epoch: 1, Candidates: []string{"SOME OTHER BINARY v9.9.9"}},
		},
	}
	rogueBytes, err := rogue.Encode()
	if err != nil {
		t.Fatalf("encode rogue registry: %v", err)
	}
	rogueObj := store.Object{
		Hash:          hashref.SumSHA256(rogueBytes),
		InterfaceHash: hashref.SumSHA256([]byte(registry.SemanticID)),
		SemanticID:    registry.SemanticID,
		Provenance:    "rogue-candidate",
		Payload:       rogueBytes,
	}
	if err := env.store.PutObject(rogueObj); err != nil {
		t.Fatalf("put rogue registry object: %v", err)
	}
	if err := env.store.SetRegistryHead(registry.SemanticID, rogueObj.Hash); err != nil {
		t.Fatalf("repoint registry head: %v", err)
	}

	// Replay again: the registry candidate changed, but authoritative resolution
	// is by entry interpreter hash, so the result is byte-identical.
	after, err := env.engine.ReplayEpisode(env.episode)
	if err != nil {
		t.Fatalf("replay after registry mutation: %v", err)
	}
	if string(after[0].Produced) != string(before[0].Produced) {
		t.Fatalf("registry mutation redirected replay: before %q, after %q",
			string(before[0].Produced), string(after[0].Produced))
	}
	if after[0].ExecPath != before[0].ExecPath {
		t.Fatalf("registry mutation changed resolved interpreter: before %q, after %q",
			before[0].ExecPath, after[0].ExecPath)
	}
	if string(after[0].Produced) != string(env.recorded) {
		t.Fatalf("post-mutation replay %q != recorded %q", string(after[0].Produced), string(env.recorded))
	}
}

// ---------------------------------------------------------------------------
// Acceptance 5: replacing either pair member causes a cache miss + re-verify
// ---------------------------------------------------------------------------

func TestPairMemberChangeCausesCacheMiss(t *testing.T) {
	env := newFixtureEnv(t)

	// First replay: cache miss (fresh store), re-verifies and caches the pair.
	first, err := env.engine.ReplayEpisode(env.episode)
	if err != nil {
		t.Fatalf("first replay: %v", err)
	}
	if first[0].CacheHit {
		t.Fatal("first replay on a fresh store should be a cache MISS")
	}

	// Second replay of the SAME pair: cache HIT (row is present).
	second, err := env.engine.ReplayEpisode(env.episode)
	if err != nil {
		t.Fatalf("second replay: %v", err)
	}
	if !second[0].CacheHit {
		t.Fatal("second replay of the same pair should be a cache HIT")
	}

	// Now change the transitionFn member of the pair (store a second, distinct
	// source object). The cache key is exactly (transitionFn, interpreter), so a
	// different transitionFn must MISS even though the interpreter is unchanged.
	altRaw := append([]byte("-- alternate fixture variant\n"), readTestdata(t, "transition_fixture.ail")...)
	altObj, err := SourceObject(altRaw, "world/replay/fixture-transition", "m5-fixture-episode-alt")
	if err != nil {
		t.Fatalf("build alt source object: %v", err)
	}
	if altObj.Hash.String() == env.srcObj.Hash.String() {
		t.Fatal("alternate source must hash differently for the pair-change test")
	}
	if err := env.store.PutObject(altObj); err != nil {
		t.Fatalf("put alt source object: %v", err)
	}

	// Re-record the alternate episode's golden bytes by running the pinned
	// interpreter once on the alternate source (recording step), then replay.
	altEpisode := env.episode
	altEpisode.Entries = []EpisodeEntry{{
		TransitionFn:   altObj.Hash,
		Interpreter:    env.interp,
		SemanticsEpoch: 1,
	}}
	// Verify the transitionFn member is a cache MISS before replay.
	if _, hit, err := env.store.GetVerifyResult(altObj.Hash, env.interp); err != nil {
		t.Fatalf("cache lookup for changed transitionFn: %v", err)
	} else if hit {
		t.Fatal("changed transitionFn must be a cache MISS")
	}

	// Run the pinned interpreter to obtain the recorded bytes for the alt source,
	// then replay it and assert it re-verifies (cache miss -> hit populated).
	altExec, err := env.archive.Resolve(env.interp)
	if err != nil {
		t.Fatalf("resolve interpreter: %v", err)
	}
	altBytes, err := runPinnedTransition(altExec, env.episode.WorldLibDir, altObj.Payload)
	if err != nil {
		t.Fatalf("record alt transition: %v", err)
	}
	altEpisode.Entries[0].RecordedResult = altBytes
	altEpisode.Entries[0].RecordedWorldHash = hashref.SumSHA256(altBytes)

	altRes, err := env.engine.ReplayEpisode(altEpisode)
	if err != nil {
		t.Fatalf("replay alt episode: %v", err)
	}
	if altRes[0].CacheHit {
		t.Fatal("changed transitionFn replay should be a cache MISS (re-verification)")
	}
	// After the alt replay the changed pair is now cached.
	if _, hit, err := env.store.GetVerifyResult(altObj.Hash, env.interp); err != nil {
		t.Fatalf("post-replay cache lookup: %v", err)
	} else if !hit {
		t.Fatal("changed pair must be cached after re-verification")
	}

	// The ORIGINAL pair row must be untouched: changing one member does not
	// evict or alter the other pair's cached row (independent keys).
	if _, hit, err := env.store.GetVerifyResult(env.srcObj.Hash, env.interp); err != nil {
		t.Fatalf("original pair cache lookup: %v", err)
	} else if !hit {
		t.Fatal("original pair row must survive an unrelated pair change")
	}
}

// TestInterpreterMemberChangeCausesCacheMiss changes the OTHER pair member (the
// interpreter) and asserts a cache miss for the new pair, keeping the original
// pair's row intact. It exercises the cache key directly via a
// synthetic-but-well-formed interpreter ref inserted alongside the real one.
//
// NOTE (NB2 boundary): this is a cache-key-only assertion — it never drives a
// full replay through the second interpreter, because `otherInterp` is a bare
// hash with no archived executable. The genuine end-to-end interpreter-member
// replay (a SECOND, hash-distinct, WORKING archived interpreter driven through
// all six replay steps) is
// TestInterpreterMemberChangeDrivesRealReplayEndToEnd below.
func TestInterpreterMemberChangeCausesCacheMiss(t *testing.T) {
	env := newFixtureEnv(t)

	// Populate the original pair via a real replay.
	if _, err := env.engine.ReplayEpisode(env.episode); err != nil {
		t.Fatalf("seed replay: %v", err)
	}
	if _, hit, err := env.store.GetVerifyResult(env.srcObj.Hash, env.interp); err != nil {
		t.Fatalf("original pair lookup: %v", err)
	} else if !hit {
		t.Fatal("original pair must be cached after replay")
	}

	// A different interpreter member: same transitionFn, different interpreter
	// HashRef. The cache key is the pair, so this must MISS.
	otherInterp := hashref.SumSHA256([]byte("a second interpreter artifact"))
	if otherInterp.String() == env.interp.String() {
		t.Fatal("synthetic interpreter must differ")
	}
	if _, hit, err := env.store.GetVerifyResult(env.srcObj.Hash, otherInterp); err != nil {
		t.Fatalf("changed-interpreter cache lookup: %v", err)
	} else if hit {
		t.Fatal("changed interpreter member must be a cache MISS")
	}
}

// TestInterpreterMemberChangeDrivesRealReplayEndToEnd resolves the NB2
// carry-forward (M4/M5): re-verify the replay path END-TO-END when the
// INTERPRETER itself is the changed member of the (transitionFn, interpreter)
// verify-cache key — not just the cache lookup in isolation.
//
// The env is "constrained by one archived binary" only if a second interpreter
// must be a DISTINCT UPSTREAM RELEASE. It need not be: a second WORKING
// interpreter that is byte-distinct (hence a distinct HashRef and a distinct
// cache-key member) is constructible from the single pinned binary as a thin
// exec wrapper. Running it produces byte-identical AILANG execution output, so
// the same recorded goldens replay bit-for-bit through BOTH interpreters. This
// drives all six replay steps (content check, entry-interpreter resolution,
// cache consult, ARCHIVED-binary invocation, byte compare, world-hash
// reconstruction) through the second interpreter and asserts:
//
//	(1) the (transitionFn, interpreter2) pair is a genuine end-to-end cache MISS
//	    that re-verifies and populates ITS OWN row;
//	(2) the executable actually invoked is the archive slot for interpreter2's
//	    hash (authoritative resolution follows the changed member);
//	(3) the (transitionFn, interpreter1) row is untouched by the change; and
//	(4) both interpreters produce byte-identical recorded results and world
//	    hashes (the wrapper is a faithful second interpreter, not a stub).
//
// STILL ENV-CONSTRAINED (documented, not faked): asserting that two
// SEMANTICALLY-DIVERGENT interpreter releases produce DIFFERENT replay bytes for
// the same source needs >=2 distinct upstream AILANG releases in the archive,
// which this single-binary CI env cannot supply. That negative case is the
// upstream multi-release integration test, out of M1 scope.
func TestInterpreterMemberChangeDrivesRealReplayEndToEnd(t *testing.T) {
	env := newFixtureEnv(t)
	bin := pinnedBinary(t) // absolute path to the real pinned interpreter

	// Seed and cache the ORIGINAL (transitionFn, interpreter1) pair via a real
	// replay so we can later prove the change leaves this row intact.
	if _, err := env.engine.ReplayEpisode(env.episode); err != nil {
		t.Fatalf("seed replay of original pair: %v", err)
	}
	if _, hit, err := env.store.GetVerifyResult(env.srcObj.Hash, env.interp); err != nil {
		t.Fatalf("original pair lookup: %v", err)
	} else if !hit {
		t.Fatal("original pair must be cached after the seed replay")
	}

	// Build a SECOND, byte-distinct, WORKING interpreter: a thin exec wrapper
	// around the real pinned binary. It is a different byte stream (=> different
	// content hash => different cache-key member) but a faithful interpreter,
	// producing identical AILANG execution output.
	wrapperDir := t.TempDir()
	wrapperPath := filepath.Join(wrapperDir, "ailang")
	wrapper := "#!/usr/bin/env bash\nexec " + strconv.Quote(bin) + " \"$@\"\n"
	if err := os.WriteFile(wrapperPath, []byte(wrapper), 0o755); err != nil {
		t.Fatalf("write wrapper interpreter: %v", err)
	}

	interp2, err := env.archive.Archive(wrapperPath)
	if err != nil {
		t.Fatalf("%s", archive.AttributeFailure("archive second (wrapper) interpreter", err))
	}
	if interp2.String() == env.interp.String() {
		t.Fatal("wrapper interpreter must be a DISTINCT HashRef from the real binary")
	}

	// Pre-condition (1): the (transitionFn, interpreter2) pair is a cache MISS
	// before any replay drives it.
	if _, hit, err := env.store.GetVerifyResult(env.srcObj.Hash, interp2); err != nil {
		t.Fatalf("interpreter2 pair pre-lookup: %v", err)
	} else if hit {
		t.Fatal("changed-interpreter pair must be a cache MISS before its first replay")
	}

	// Drive the FULL replay through interpreter2 using the SAME committed goldens
	// (the wrapper produces byte-identical output to interpreter1).
	ep2 := env.episode
	ep2.Entries = []EpisodeEntry{{
		TransitionFn:      env.srcObj.Hash,
		Interpreter:       interp2,
		SemanticsEpoch:    1,
		RecordedResult:    env.recorded,
		RecordedWorldHash: env.worldRef,
	}}
	res2, err := env.engine.ReplayEpisode(ep2)
	if err != nil {
		t.Fatalf("end-to-end replay through interpreter2: %v", err)
	}

	// (1) genuine end-to-end cache MISS on this run (re-verified, not a hit).
	if res2[0].CacheHit {
		t.Fatal("first end-to-end replay of the changed interpreter pair must be a cache MISS")
	}
	// ... and the row is now populated for the NEW pair.
	if _, hit, err := env.store.GetVerifyResult(env.srcObj.Hash, interp2); err != nil {
		t.Fatalf("interpreter2 pair post-lookup: %v", err)
	} else if !hit {
		t.Fatal("changed-interpreter pair must be cached after its end-to-end re-verification")
	}

	// (2) authoritative resolution followed the CHANGED member: the invoked
	// executable is interpreter2's archive slot, not interpreter1's.
	want2, err := env.archive.Resolve(interp2)
	if err != nil {
		t.Fatalf("resolve interpreter2: %v", err)
	}
	if res2[0].ExecPath != want2 {
		t.Fatalf("end-to-end replay invoked %q, expected the changed interpreter %q", res2[0].ExecPath, want2)
	}
	if want1, err := env.archive.Resolve(env.interp); err != nil {
		t.Fatalf("resolve interpreter1: %v", err)
	} else if res2[0].ExecPath == want1 {
		t.Fatal("changed-member replay must NOT resolve back to interpreter1")
	}

	// (3) the ORIGINAL (transitionFn, interpreter1) row is untouched.
	if _, hit, err := env.store.GetVerifyResult(env.srcObj.Hash, env.interp); err != nil {
		t.Fatalf("original pair post-change lookup: %v", err)
	} else if !hit {
		t.Fatal("original pair row must survive the interpreter-member change")
	}

	// (4) the second interpreter is faithful: byte-identical result + world hash.
	if string(res2[0].Produced) != string(env.recorded) {
		t.Fatalf("interpreter2 produced %q != recorded %q", string(res2[0].Produced), string(env.recorded))
	}
	if res2[0].WorldHash.String() != env.worldRef.String() {
		t.Fatalf("interpreter2 world hash %q != recorded %q", res2[0].WorldHash.String(), env.worldRef.String())
	}
}

// ---------------------------------------------------------------------------
// M4 carry-forward 1: KindExecFailure replay branch
// ---------------------------------------------------------------------------

// TestExecFailureReplayBranch exercises the archive.KindExecFailure replay path:
// when the entry-pinned interpreter is archived but is NOT a working AILANG
// binary (a non-executable stub), invoking it fails and replay surfaces a
// KindExecFailure ReplayError rather than proceeding on a broken interpreter.
func TestExecFailureReplayBranch(t *testing.T) {
	_ = pinnedBinary(t) // ensure the env is a replay-capable one (skip otherwise)

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "world.db")
	s, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	a := archive.New(dbPath)
	engine := NewEngine(s, a)

	// Hand-place a non-working "interpreter" into the archive tree at the slot
	// its own content hash addresses, so Resolve succeeds and verifyExecutable
	// passes, but running it fails (it is not a real binary).
	stub := []byte("#!/nonexistent/interpreter\nnot a real binary\n")
	stubRef := hashref.SumSHA256(stub)
	slot := filepath.Join(a.Root(), "interpreters", stubRef.Algo(), stubRef.Digest())
	if err := os.MkdirAll(slot, 0o755); err != nil {
		t.Fatalf("mkdir stub slot: %v", err)
	}
	if err := os.WriteFile(filepath.Join(slot, "ailang"), stub, 0o555); err != nil {
		t.Fatalf("write stub interpreter: %v", err)
	}

	rawSource := readTestdata(t, "transition_fixture.ail")
	srcObj, err := SourceObject(rawSource, "world/replay/fixture-transition", "m5-fixture-episode")
	if err != nil {
		t.Fatalf("source object: %v", err)
	}
	if err := s.PutObject(srcObj); err != nil {
		t.Fatalf("put source: %v", err)
	}

	ep := Episode{
		Entries: []EpisodeEntry{{
			TransitionFn:      srcObj.Hash,
			Interpreter:       stubRef,
			SemanticsEpoch:    1,
			RecordedResult:    []byte("unused"),
			RecordedWorldHash: hashref.SumSHA256([]byte("unused")),
		}},
		WorldLibDir: worldLibDir(t),
	}

	_, err = engine.ReplayEpisode(ep)
	re, ok := archive.IsReplayError(err)
	if !ok {
		t.Fatalf("expected *archive.ReplayError, got %T: %v", err, err)
	}
	if re.Kind != archive.KindExecFailure {
		t.Fatalf("expected KindExecFailure, got %q", re.Kind)
	}
}

// ---------------------------------------------------------------------------
// M4 carry-forward 2: sidecar-present / executable-absent absent-artifact edge
// ---------------------------------------------------------------------------

// TestSidecarPresentExecutableAbsentResolvesAbsent covers the archive.Resolve()
// idempotence-recovery edge flagged by the M4 evaluator: a slot that has a
// sidecar manifest but whose executable file is missing must Resolve() to a
// KindAbsentArtifact ReplayError (Resolve keys on the executable file, not the
// sidecar), and authoritative replay must therefore fail absent rather than
// treat the slot as usable.
func TestSidecarPresentExecutableAbsentResolvesAbsent(t *testing.T) {
	bin := pinnedBinary(t)

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "world.db")
	s, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	a := archive.New(dbPath)

	// Archive the real binary, then delete ONLY its executable file, leaving the
	// sidecar manifest in place — the exact "sidecar present, executable absent"
	// state the M4 evaluator called out.
	interp, err := a.Archive(bin)
	if err != nil {
		t.Fatalf("%s", archive.AttributeFailure("archive", err))
	}
	if _, err := a.ReadManifest(interp); err != nil {
		t.Fatalf("manifest should be present: %v", err)
	}
	slot := filepath.Join(a.Root(), "interpreters", interp.Algo(), interp.Digest())
	if err := os.Remove(filepath.Join(slot, "ailang")); err != nil {
		t.Fatalf("remove executable: %v", err)
	}

	// Resolve now fails absent even though the sidecar remains.
	if _, err := a.Resolve(interp); err != nil {
		re, ok := archive.IsReplayError(err)
		if !ok || re.Kind != archive.KindAbsentArtifact {
			t.Fatalf("expected KindAbsentArtifact from Resolve, got %v", err)
		}
	} else {
		t.Fatal("Resolve must fail when the executable file is absent")
	}

	// And replay fails absent through the same path.
	rawSource := readTestdata(t, "transition_fixture.ail")
	srcObj, err := SourceObject(rawSource, "world/replay/fixture-transition", "m5-fixture-episode")
	if err != nil {
		t.Fatalf("source object: %v", err)
	}
	if err := s.PutObject(srcObj); err != nil {
		t.Fatalf("put source: %v", err)
	}
	ep := Episode{
		Entries: []EpisodeEntry{{
			TransitionFn:      srcObj.Hash,
			Interpreter:       interp,
			SemanticsEpoch:    1,
			RecordedResult:    []byte("unused"),
			RecordedWorldHash: hashref.SumSHA256([]byte("unused")),
		}},
		WorldLibDir: worldLibDir(t),
	}
	_, err = NewEngine(s, a).ReplayEpisode(ep)
	re, ok := archive.IsReplayError(err)
	if !ok || re.Kind != archive.KindAbsentArtifact {
		t.Fatalf("expected KindAbsentArtifact from replay, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// Guard: absent transition source object fails absent (not a silent pass)
// ---------------------------------------------------------------------------

func TestAbsentTransitionSourceFails(t *testing.T) {
	env := newFixtureEnv(t)
	// Point the entry at a transitionFn the store does not hold.
	env.episode.Entries[0].TransitionFn = hashref.SumSHA256([]byte("no such source"))
	_, err := env.engine.ReplayEpisode(env.episode)
	re, ok := archive.IsReplayError(err)
	if !ok || re.Kind != archive.KindAbsentArtifact {
		t.Fatalf("expected KindAbsentArtifact for absent source, got %v", err)
	}
}
