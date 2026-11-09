package broker

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
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
