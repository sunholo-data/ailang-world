package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// S7: the seed decisions are a pure table (row 257 §2.3, §2.4).
func TestSeedNextIsPure(t *testing.T) {
	next := []struct {
		attempts int
		running  bool
		verified bool
		want     seedAction
	}{
		{0, false, true, seedReady},
		{3, false, true, seedReady},
		{1, true, true, seedReady},
		{1, true, false, seedJoin},
		{3, true, false, seedJoin},
		{0, false, false, seedStart},
		{2, false, false, seedStart},
		{3, false, false, seedGiveUp},
		{4, false, false, seedGiveUp},
	}
	for _, c := range next {
		if got := execSeedNext(seedView{attempts: c.attempts, running: c.running}, c.verified); got != c.want {
			t.Errorf("execSeedNext(%d,%v,%v) = %v, want %v", c.attempts, c.running, c.verified, got, c.want)
		}
	}

	state := []struct {
		attempts            int
		running, ok, failed bool
		want                string
	}{
		{1, true, false, false, ExecSeedCopying},
		{0, false, true, false, ExecSeedReady},
		{1, false, true, false, ExecSeedReady},
		{1, false, false, true, ExecSeedFailed},
		{2, false, false, true, ExecSeedFailed},
		{3, false, false, true, ExecSeedGaveUp},
		{0, false, false, false, ""},
	}
	for _, c := range state {
		if got := seedStateName(c.attempts, c.running, c.ok, c.failed); got != c.want {
			t.Errorf("seedStateName(%d,%v,%v,%v) = %q, want %q", c.attempts, c.running, c.ok, c.failed, got, c.want)
		}
	}

	if got := seedAttemptDetail("GOCACHE", "the copy exited 1", 2); got != "GOCACHE: the copy exited 1 (attempt 2 of 3)" {
		t.Errorf("seedAttemptDetail = %q", got)
	}
	if got := seedStillCopyingDetail("GOCACHE", "bind", 4210*time.Millisecond); got != "GOCACHE: still copying since bind (4.2s); retry" {
		t.Errorf("seedStillCopyingDetail = %q", got)
	}
	if got := seedStagingName("proj", "GOCACHE", "1-2", 3); got != ".proj.GOCACHE.seeding-1-2-3" {
		t.Errorf("seedStagingName = %q", got)
	}
}

// S12: the marker is exact bytes, versioned, named, and size-capped (§2.1).
func TestSeedMarkerValidIsPure(t *testing.T) {
	if got := string(seedMarkerBytes("GOCACHE")); got != "world-exec-seed v1 GOCACHE\n" {
		t.Fatalf("seedMarkerBytes = %q", got)
	}
	if !seedMarkerValid("GOCACHE", []byte("world-exec-seed v1 GOCACHE\n")) {
		t.Fatal("the exact marker must be valid")
	}
	bad := []string{
		"world-exec-seed v1 GOCACHE\nx",
		"world-exec-seed v1 GOCACHE",
		"world-exec-seed v0 GOCACHE\n",
		"world-exec-seed v2 GOCACHE\n",
		"world-exec-seed v1 OTHER\n",
		"",
		strings.Repeat("x", 65),
	}
	for _, s := range bad {
		if seedMarkerValid("GOCACHE", []byte(s)) {
			t.Errorf("seedMarkerValid accepted %q", s)
		}
	}
	if seedMarkerValid("GOCACHE", nil) {
		t.Error("seedMarkerValid accepted nil")
	}
	n := strings.Repeat("n", 45)
	if got := len(seedMarkerBytes(n)); got != 65 {
		t.Fatalf("len(seedMarkerBytes(45-byte name)) = %d, want 65", got)
	}
	if seedMarkerValid(n, seedMarkerBytes(n)) {
		t.Error("the size cap must reject a 65-byte marker")
	}
}

// slowSeedCommand runs the default copy after 0.5 s, on every OS.
func slowSeedCommand(seed, dst string) (string, []string) {
	bin, args := defaultSeedCommand(seed, dst)
	return "/bin/sh", append([]string{"-c", `sleep 0.5; exec "$0" "$@"`, bin}, args...)
}

// makeSeed writes <dir>/ab/entry = "seeded\n" under a fresh canonicalTempDir and returns the seed dir.
func makeSeed(t *testing.T) string {
	t.Helper()
	dir := filepath.Join(canonicalTempDir(t), "seed")
	if err := os.MkdirAll(filepath.Join(dir, "ab"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ab", "entry"), []byte("seeded\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

// seedFixture is newExecFixture with caches.GOCACHE.seed = seed.
func seedFixture(t *testing.T, seed string) *execFixture {
	t.Helper()
	return newExecFixture(t, execFixtureOpts{profile: func(m map[string]any) {
		m["caches"] = map[string]any{"GOCACHE": map[string]any{"seed": seed}}
	}})
}

// waitSeedIdle waits (≤ 20 s, else t.Fatal) until h.seed.running is nil, re-checking after each close.
func waitSeedIdle(t *testing.T, h *ExecHandler) {
	t.Helper()
	deadline := time.Now().Add(20 * time.Second)
	for {
		h.seed.mu.Lock()
		done := h.seed.running
		h.seed.mu.Unlock()
		if done == nil {
			return
		}
		select {
		case <-done:
		case <-time.After(time.Until(deadline)):
			t.Fatal("the seed copy did not end within 20s")
		}
	}
}

// seedState returns the single SeedStatus entry (t.Fatal if len != 1).
func seedState(t *testing.T, h *ExecHandler) ExecSeedStatus {
	t.Helper()
	st := h.SeedStatus()
	if len(st) != 1 {
		t.Fatalf("SeedStatus = %+v, want one entry", st)
	}
	return st[0]
}

// stagingLeft lists the entries of filepath.Dir(f.cache) whose name contains ".seeding-".
func stagingLeft(t *testing.T, f *execFixture) []string {
	t.Helper()
	ents, err := os.ReadDir(filepath.Dir(f.cache))
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, e := range ents {
		if strings.Contains(e.Name(), ".seeding-") {
			out = append(out, e.Name())
		}
	}
	return out
}

// invokeLine runs one Workspace.Exec "test a" through a live broker session (as TestExecFailureClassesReachTheRecord does) and returns EffectFailedError.Reason ("" on success).
func invokeLine(t *testing.T, h Handler) string {
	t.Helper()
	scope := WorkspaceScope
	st := openTestStore(t)
	s := newSession(st, "ep1", []Capability{{Effect: EffectWorkspaceExec, Scope: scope, ExpiresAt: 100, Budget: 5}},
		Registry{EffectWorkspaceExec: h}, Live, nil)
	payload, _ := json.Marshal(map[string]any{"command": "test", "args": []string{"a"}})
	_, _, err := s.Invoke(boundedTestContext(t), EffectRequest{Effect: EffectWorkspaceExec, Scope: scope, Cost: 1, Now: 1}, payload)
	var failed *EffectFailedError
	if !errors.As(err, &failed) {
		t.Fatalf("Invoke = %v, want *EffectFailedError", err)
	}
	return failed.Reason
}

// S1: a copy started at bind serves the first call without a copy in the call.
func TestSeedAtBindFirstCallSpawnsNoCopy(t *testing.T) {
	f := seedFixture(t, makeSeed(t))
	var log bytes.Buffer
	f.h.log = &log
	f.h.StartSeed()
	waitSeedIdle(t, f.h)
	want := ExecSeedStatus{Project: "proj", Cache: "GOCACHE", State: ExecSeedReady, Attempts: 1, By: "bind"}
	if got := seedState(t, f.h); got.State != want.State || got.Attempts != want.Attempts || got.By != want.By || got.Failure != "" {
		t.Fatalf("seed state = %+v, want %+v", got, want)
	}
	n := f.h.Spawns()
	f.mustRun(t, "test", "a")
	if f.h.Spawns() != n+1 {
		t.Fatalf("the call spawned %d processes, want 1 (srt only)", f.h.Spawns()-n)
	}
	if got := f.stubCalls(t); got != 1 {
		t.Fatalf("the stub ran %d times, want 1", got)
	}
	data, err := os.ReadFile(filepath.Join(f.cache, "GOCACHE", "ab", "entry"))
	if err != nil || string(data) != "seeded\n" {
		t.Fatalf("seeded entry = %q, %v", data, err)
	}
	marker, err := os.ReadFile(filepath.Join(f.cache, "GOCACHE", seedMarkerName))
	if err != nil || string(marker) != "world-exec-seed v1 GOCACHE\n" {
		t.Fatalf("marker = %q, %v", marker, err)
	}
	if left := stagingLeft(t, f); len(left) != 0 {
		t.Fatalf("staging left behind: %v", left)
	}
	if !strings.Contains(log.String(), `exec cache for episode "ep1" (proj): GOCACHE seeded at bind `) {
		t.Fatalf("log = %q", log.String())
	}
}

// S2: a failed bind copy is what the first call reports.
func TestSeedAtBindFailureReachesTheCall(t *testing.T) {
	f := seedFixture(t, filepath.Join(canonicalTempDir(t), "no-such-seed"))
	var log bytes.Buffer
	f.h.log = &log
	f.h.StartSeed()
	waitSeedIdle(t, f.h)
	first := "seed-cache: GOCACHE: the copy exited 1 (attempt 1 of 3)"
	if got := seedState(t, f.h); got.State != ExecSeedFailed || got.Attempts != 1 || got.By != "bind" || got.Failure != first {
		t.Fatalf("seed state = %+v", got)
	}
	if !strings.Contains(log.String(), "seeding at bind failed: "+first) {
		t.Fatalf("log = %q", log.String())
	}
	second := "seed-cache: GOCACHE: the copy exited 1 (attempt 2 of 3)"
	if got := invokeLine(t, f.h); got != second {
		t.Fatalf("call failure = %q, want %q", got, second)
	}
	if f.stubCalls(t) != 0 {
		t.Fatal("a failed seed must not reach srt")
	}
}

// S3: the attempt cap holds; the third refusal spawns nothing.
func TestSeedAttemptsAreCapped(t *testing.T) {
	f := seedFixture(t, filepath.Join(canonicalTempDir(t), "no-such-seed"))
	f.h.StartSeed()
	waitSeedIdle(t, f.h)
	for k := 2; k <= 3; k++ {
		_, err := f.call(t, "test", "a")
		var hf *HandlerFailure
		want := "seed-cache: GOCACHE: the copy exited 1 (attempt " + strconv.Itoa(k) + " of 3)"
		if !errors.As(err, &hf) || hf.Line() != want {
			t.Fatalf("call %d: err = %v, want line %q", k-1, err, want)
		}
	}
	n := f.h.Spawns()
	_, err := f.call(t, "test", "a")
	var hf *HandlerFailure
	want := "seed-cache: GOCACHE: the copy exited 1 (attempt 3 of 3); 3 of 3 attempts failed: fix the seed, then restart serve"
	if !errors.As(err, &hf) || hf.Line() != want {
		t.Fatalf("call 3: err = %v, want line %q", err, want)
	}
	if f.h.Spawns() != n {
		t.Fatalf("the refused call spawned %d processes, want 0", f.h.Spawns()-n)
	}
	if got := seedState(t, f.h); got.State != ExecSeedGaveUp || got.Attempts != 3 {
		t.Fatalf("seed state = %+v, want gave-up with 3 attempts", got)
	}
}

// S4: a partial copy is never accepted, and a crash window's staging is swept.
func TestSeedPartialCopyIsNeverAccepted(t *testing.T) {
	seed := makeSeed(t)
	f2 := filepath.Join(seed, "ab", "locked")
	if err := os.WriteFile(f2, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(f2, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(f2, 0o644) })
	f := seedFixture(t, seed)
	f.h.StartSeed()
	waitSeedIdle(t, f.h)
	if got := seedState(t, f.h); got.State != ExecSeedFailed || !strings.HasPrefix(got.Failure, "seed-cache: GOCACHE: the copy exited ") {
		t.Fatalf("seed state = %+v, want failed with a copy exit", got)
	}
	if _, err := os.Lstat(filepath.Join(f.cache, "GOCACHE")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Lstat(published GOCACHE) = %v, want ENOENT", err)
	}
	if left := stagingLeft(t, f); len(left) != 0 {
		t.Fatalf("staging left after the failed copy: %v", left)
	}
	// A crash window from a dead incarnation: a staging dir with a partial copy.
	old := filepath.Join(filepath.Dir(f.cache), seedStagingName("proj", "GOCACHE", "dead-1", 1))
	if err := os.MkdirAll(filepath.Join(old, "copy"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(old, "copy", "partial"), []byte("half"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(f2, 0o644); err != nil {
		t.Fatal(err)
	}
	f.mustRun(t, "test", "a")
	if _, err := os.Lstat(old); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("the dead staging dir survived: %v", err)
	}
	if !seedCacheVerified(filepath.Join(f.cache, "GOCACHE"), "GOCACHE") {
		t.Fatal("GOCACHE is not verified after the second attempt")
	}
	_ = filepath.Walk(filepath.Join(f.cache, "GOCACHE"), func(p string, _ os.FileInfo, err error) error {
		if err == nil && filepath.Base(p) == "partial" {
			t.Errorf("partial copy reached the cache: %s", p)
		}
		return nil
	})
	if left := stagingLeft(t, f); len(left) != 0 {
		t.Fatalf("staging left after the second attempt: %v", left)
	}
}

// slowCall runs Execute "test a" with a 100 ms context, the caller waiting on a slow copy.
func slowCall(t *testing.T, f *execFixture) error {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	data, _ := json.Marshal(map[string]any{"command": "test", "args": []string{"a"}})
	_, err := f.h.Execute(ctx, EffectRequest{Effect: EffectWorkspaceExec, Scope: WorkspaceScope, Cost: 1}, data)
	return err
}

// S5: a copy outlives the call that waits for it; its context is not the call's.
func TestSeedOutlivesTheCallThatWaits(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		f := seedFixture(t, makeSeed(t))
		f.h.seedCommand = slowSeedCommand
		var hf *HandlerFailure
		err := slowCall(t, f)
		if !errors.As(err, &hf) {
			t.Fatalf("err = %v, want *HandlerFailure", err)
		}
		if l := hf.Line(); !strings.HasPrefix(l, "seed-cache: GOCACHE: still copying since call (") || !strings.HasSuffix(l, "); retry") {
			t.Fatalf("line = %q", l)
		}
		waitSeedIdle(t, f.h)
		if got := seedState(t, f.h); got.State != ExecSeedReady || got.Attempts != 1 {
			t.Fatalf("seed state = %+v, want ready with 1 attempt", got)
		}
		n := f.h.Spawns()
		f.mustRun(t, "test", "a")
		if f.h.Spawns() != n+1 {
			t.Fatalf("the call spawned %d processes, want 1 (srt only)", f.h.Spawns()-n)
		}
	})
	t.Run("bind", func(t *testing.T) {
		f := seedFixture(t, makeSeed(t))
		f.h.seedCommand = slowSeedCommand
		f.h.StartSeed()
		var hf *HandlerFailure
		err := slowCall(t, f)
		if !errors.As(err, &hf) {
			t.Fatalf("err = %v, want *HandlerFailure", err)
		}
		if l := hf.Line(); !strings.HasPrefix(l, "seed-cache: GOCACHE: still copying since bind (") {
			t.Fatalf("line = %q", l)
		}
		waitSeedIdle(t, f.h)
		if got := seedState(t, f.h); got.State != ExecSeedReady {
			t.Fatalf("seed state = %+v, want ready", got)
		}
	})
}

// S6: the copy has its own bound, not the call's.
func TestSeedHasItsOwnBound(t *testing.T) {
	f := seedFixture(t, makeSeed(t))
	f.h.seedBudget = time.Millisecond
	f.h.seedCommand = slowSeedCommand
	f.h.StartSeed()
	waitSeedIdle(t, f.h)
	want := "seed-cache: GOCACHE: timed out after 1ms (attempt 1 of 3)"
	if got := seedState(t, f.h); got.State != ExecSeedFailed || got.Failure != want {
		t.Fatalf("seed state = %+v, want failed with %q", got, want)
	}
}

// S11: a seed that carries World's marker is refused before any copy runs.
func TestSeedHoldingTheMarkerIsRefused(t *testing.T) {
	seed := makeSeed(t)
	if err := os.WriteFile(filepath.Join(seed, seedMarkerName), seedMarkerBytes("GOCACHE"), 0o644); err != nil {
		t.Fatal(err)
	}
	f := seedFixture(t, seed)
	n := f.h.Spawns()
	f.h.StartSeed()
	waitSeedIdle(t, f.h)
	want := "seed-cache: GOCACHE: the seed holds World's marker; remove it from the seed (attempt 1 of 3)"
	if got := seedState(t, f.h); got.Failure != want {
		t.Fatalf("Failure = %q, want %q", got.Failure, want)
	}
	if _, err := os.Lstat(filepath.Join(f.cache, "GOCACHE")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Lstat(GOCACHE) = %v, want ENOENT", err)
	}
	if f.h.Spawns() != n {
		t.Fatalf("a marked seed spawned %d processes, want 0", f.h.Spawns()-n)
	}
}

// plantCache writes dir/ab/entry = "planted\n", plus dir/<extra> = "x" when extra != "", plus
// dir/.world-seed-complete = marker when marker != nil.
func plantCache(t *testing.T, dir string, marker []byte, extra string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, "ab"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ab", "entry"), []byte("planted\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if extra != "" {
		if err := os.WriteFile(filepath.Join(dir, extra), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if marker != nil {
		if err := os.WriteFile(filepath.Join(dir, seedMarkerName), marker, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// requireSeededEntry asserts GOCACHE is verified and its entry is the copied seed.
func requireSeededEntry(t *testing.T, f *execFixture) {
	t.Helper()
	gc := filepath.Join(f.cache, "GOCACHE")
	if !seedCacheVerified(gc, "GOCACHE") {
		t.Fatal("GOCACHE is not verified")
	}
	data, err := os.ReadFile(filepath.Join(gc, "ab", "entry"))
	if err != nil || string(data) != "seeded\n" {
		t.Fatalf("ab/entry = %q, %v; want the copied seed", data, err)
	}
}

// S8: a restart reuses a cache only when its marker is valid.
func TestStartSeedReusesOnlyAVerifiedCache(t *testing.T) {
	t.Run("verified", func(t *testing.T) {
		f := seedFixture(t, makeSeed(t))
		plantCache(t, filepath.Join(f.cache, "GOCACHE"), seedMarkerBytes("GOCACHE"), "")
		n := f.h.Spawns()
		f.h.StartSeed()
		if f.h.Spawns() != n {
			t.Fatalf("StartSeed spawned %d processes over a verified cache, want 0", f.h.Spawns()-n)
		}
		if got := seedState(t, f.h); got.State != ExecSeedReady || got.Attempts != 0 || got.By != "bind" {
			t.Fatalf("seed state = %+v, want ready with 0 attempts by bind", got)
		}
		f.mustRun(t, "test", "a")
		if f.h.Spawns() != n+1 {
			t.Fatalf("the call spawned %d processes, want 1 (srt only)", f.h.Spawns()-n)
		}
		data, err := os.ReadFile(filepath.Join(f.cache, "GOCACHE", "ab", "entry"))
		if err != nil || string(data) != "planted\n" {
			t.Fatalf("the verified cache was replaced: %q, %v", data, err)
		}
	})
	t.Run("stale marker", func(t *testing.T) {
		f := seedFixture(t, makeSeed(t))
		plantCache(t, filepath.Join(f.cache, "GOCACHE"), []byte("world-exec-seed v0 GOCACHE\n"), "")
		n := f.h.Spawns()
		f.h.StartSeed()
		waitSeedIdle(t, f.h)
		if f.h.Spawns() != n+1 {
			t.Fatalf("the stale cache spawned %d processes, want 1 (the copy)", f.h.Spawns()-n)
		}
		if got := seedState(t, f.h); got.State != ExecSeedReady || got.Attempts != 1 {
			t.Fatalf("seed state = %+v, want ready with 1 attempt", got)
		}
		requireSeededEntry(t, f)
	})
}

// S9: a legacy cache with no marker is removed and reseeded, at bind or by a call.
func TestLegacyPartialCacheIsRemovedAndReseeded(t *testing.T) {
	t.Run("bind", func(t *testing.T) {
		f := seedFixture(t, makeSeed(t))
		plantCache(t, filepath.Join(f.cache, "GOCACHE"), nil, "stray")
		f.h.StartSeed()
		waitSeedIdle(t, f.h)
		if got := seedState(t, f.h); got.State != ExecSeedReady {
			t.Fatalf("seed state = %+v, want ready", got)
		}
		if _, err := os.Lstat(filepath.Join(f.cache, "GOCACHE", "stray")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("the legacy stray file survived: %v", err)
		}
		requireSeededEntry(t, f)
	})
	t.Run("call", func(t *testing.T) {
		f := seedFixture(t, makeSeed(t))
		plantCache(t, filepath.Join(f.cache, "GOCACHE"), nil, "stray")
		n := f.h.Spawns()
		f.mustRun(t, "test", "a")
		if f.h.Spawns() != n+2 {
			t.Fatalf("the call spawned %d processes, want 2 (the copy, then srt)", f.h.Spawns()-n)
		}
		if _, err := os.Lstat(filepath.Join(f.cache, "GOCACHE", "stray")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("the legacy stray file survived: %v", err)
		}
		requireSeededEntry(t, f)
	})
}

// S10: no crash window is ever trusted; each one is swept or refused.
func TestSeedCrashWindowsAreNeverTrusted(t *testing.T) {
	setup := func(t *testing.T) (*execFixture, string) {
		t.Helper()
		f := seedFixture(t, makeSeed(t))
		if err := os.MkdirAll(filepath.Dir(f.cache), 0o700); err != nil {
			t.Fatal(err)
		}
		return f, filepath.Join(filepath.Dir(f.cache), seedStagingName("proj", "GOCACHE", "dead-1", 1))
	}
	sweptAndCopied := func(t *testing.T, f *execFixture, old string, n int64) {
		t.Helper()
		f.h.StartSeed()
		waitSeedIdle(t, f.h)
		if f.h.Spawns() != n+1 {
			t.Fatalf("spawned %d processes, want 1 (the copy)", f.h.Spawns()-n)
		}
		if _, err := os.Lstat(old); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("the dead staging dir survived: %v", err)
		}
		requireSeededEntry(t, f)
		if _, err := os.Lstat(filepath.Join(f.cache, "GOCACHE", "planted-only")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("a dead staging entry reached GOCACHE: %v", err)
		}
	}
	t.Run("partial staging", func(t *testing.T) {
		f, old := setup(t)
		if err := os.MkdirAll(filepath.Join(old, "copy"), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(old, "copy", "planted-only"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		sweptAndCopied(t, f, old, f.h.Spawns())
	})
	t.Run("full staging, no marker", func(t *testing.T) {
		f, old := setup(t)
		plantCache(t, filepath.Join(old, "copy"), nil, "planted-only")
		sweptAndCopied(t, f, old, f.h.Spawns())
	})
	t.Run("full staging, marker", func(t *testing.T) {
		f, old := setup(t)
		plantCache(t, filepath.Join(old, "copy"), seedMarkerBytes("GOCACHE"), "planted-only")
		sweptAndCopied(t, f, old, f.h.Spawns())
	})
	t.Run("power loss", func(t *testing.T) {
		f, _ := setup(t)
		plantCache(t, filepath.Join(f.cache, "GOCACHE"), []byte{}, "")
		n := f.h.Spawns()
		f.h.StartSeed()
		waitSeedIdle(t, f.h)
		if f.h.Spawns() != n+1 {
			t.Fatalf("spawned %d processes, want 1 (the copy)", f.h.Spawns()-n)
		}
		if got := seedState(t, f.h); got.State != ExecSeedReady {
			t.Fatalf("seed state = %+v, want ready", got)
		}
		requireSeededEntry(t, f)
	})
	t.Run("unremovable", func(t *testing.T) {
		f, _ := setup(t)
		plantCache(t, filepath.Join(f.cache, "GOCACHE"), nil, "")
		if err := os.Chmod(f.cache, 0o500); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chmod(f.cache, 0o700) })
		f.h.StartSeed()
		waitSeedIdle(t, f.h)
		want := "seed-cache: GOCACHE: an unverified copy cannot be removed (attempt 1 of 3)"
		if got := seedState(t, f.h); got.State != ExecSeedFailed || got.Failure != want {
			t.Fatalf("seed state = %+v, want failed with %q", got, want)
		}
	})
}

// S13: a surviving old copy never reaches the new attempt; cleanup is bounded.
func TestSeedRestartIgnoresASurvivingOldCopy(t *testing.T) {
	t.Run("survivor", func(t *testing.T) {
		f := seedFixture(t, makeSeed(t))
		epDir := filepath.Dir(f.cache)
		if err := os.MkdirAll(epDir, 0o700); err != nil {
			t.Fatal(err)
		}
		old := filepath.Join(epDir, seedStagingName("proj", "GOCACHE", "old-1", 1))
		cmd := exec.Command("/bin/sh", "-c", `i=0; while :; do mkdir -p "$0/copy"; : > "$0/copy/old-$i"; i=$((i+1)); sleep 0.01; done`, old)
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		stopped := false
		stop := func() {
			if !stopped {
				stopped = true
				_ = cmd.Process.Kill()
				_ = cmd.Wait()
			}
		}
		t.Cleanup(stop)
		deadline := time.Now().Add(5 * time.Second)
		for {
			if _, err := os.Lstat(filepath.Join(old, "copy", "old-0")); err == nil {
				break
			}
			if time.Now().After(deadline) {
				t.Fatal("the old process never wrote its first file")
			}
			time.Sleep(10 * time.Millisecond)
		}
		f.h.StartSeed()
		waitSeedIdle(t, f.h)
		if got := seedState(t, f.h); got.State != ExecSeedReady {
			t.Fatalf("seed state = %+v, want ready", got)
		}
		requireSeededEntry(t, f)
		walkErr := filepath.Walk(filepath.Join(f.cache, "GOCACHE"), func(p string, _ os.FileInfo, err error) error {
			if err == nil && strings.HasPrefix(filepath.Base(p), "old-") {
				t.Errorf("an old-process file reached GOCACHE: %s", p)
			}
			return nil
		})
		if walkErr != nil {
			t.Fatal(walkErr)
		}
		stop()
		if err := os.MkdirAll(old+"/copy", 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(old, "copy", "old-final"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(filepath.Join(f.cache, "GOCACHE", seedMarkerName)); err != nil {
			t.Fatal(err)
		}
		f.mustRun(t, "test", "a")
		if _, err := os.Lstat(old); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("the old staging dir survived the restart: %v", err)
		}
		if !seedCacheVerified(filepath.Join(f.cache, "GOCACHE"), "GOCACHE") {
			t.Fatal("the marker is not valid again after attempt 2")
		}
	})
	t.Run("bounded cleanup", func(t *testing.T) {
		f := seedFixture(t, makeSeed(t))
		epDir := filepath.Dir(f.cache)
		if err := os.MkdirAll(epDir, 0o700); err != nil {
			t.Fatal(err)
		}
		for i := 0; i < 10; i++ {
			d := filepath.Join(epDir, seedStagingName("proj", "GOCACHE", fmt.Sprintf("old%02d", i), 1))
			if err := os.Mkdir(d, 0o700); err != nil {
				t.Fatal(err)
			}
		}
		other := filepath.Join(epDir, ".proj.OTHER.seeding-x-1")
		if err := os.Mkdir(other, 0o700); err != nil {
			t.Fatal(err)
		}
		f.h.StartSeed()
		waitSeedIdle(t, f.h)
		if got := seedState(t, f.h); got.State != ExecSeedReady {
			t.Fatalf("seed state = %+v, want ready", got)
		}
		ents, err := os.ReadDir(epDir)
		if err != nil {
			t.Fatal(err)
		}
		var left []string
		for _, e := range ents {
			if strings.HasPrefix(e.Name(), seedStagingPrefix("proj", "GOCACHE")) {
				left = append(left, e.Name())
			}
		}
		want := []string{seedStagingName("proj", "GOCACHE", "old08", 1), seedStagingName("proj", "GOCACHE", "old09", 1)}
		if strings.Join(left, ",") != strings.Join(want, ",") {
			t.Fatalf("remaining staging = %v, want %v", left, want)
		}
		if _, err := os.Lstat(other); err != nil {
			t.Fatalf("the other project's staging was touched: %v", err)
		}
	})
}
