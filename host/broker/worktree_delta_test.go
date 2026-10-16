package broker

// Row 184 (w-toolchain-worktree-delta §3 AC1–AC4): the snapshot and the delta,
// on real git repositories and plain directories, then through the handlers.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

func shaHex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// gitRepo makes a committed repository: tracked a.txt, b.txt, a .gitignore
// ignoring venv/, and an untracked dirty u.txt present before any snapshot.
func gitRepo(t *testing.T) string {
	t.Helper()
	root := filepath.Join(canonicalTempDir(t), "ep1")
	mustMkdir(t, root)
	fixtureGit(t, root, "init", "-q", ".")
	mustWrite(t, filepath.Join(root, "a.txt"), "alpha\n")
	mustWrite(t, filepath.Join(root, "b.txt"), "beta\n")
	mustWrite(t, filepath.Join(root, ".gitignore"), "venv/\n")
	fixtureGit(t, root, "add", ".")
	fixtureGit(t, root, "commit", "-q", "-m", "init")
	mustWrite(t, filepath.Join(root, "u.txt"), "untracked before\n")
	return root
}

func deltaAround(t *testing.T, root string, change func()) *WorktreeDelta {
	t.Helper()
	p := beginWorktreeDelta(boundedTestContext(t), root)
	change()
	return p.finish(boundedTestContext(t))
}

func changeSet(d *WorktreeDelta) map[string]WorktreeChange {
	m := map[string]WorktreeChange{}
	for _, c := range d.Changes {
		m[c.Path] = c
	}
	return m
}

func assertChange(t *testing.T, d *WorktreeDelta, path, op, content string) {
	t.Helper()
	c, ok := changeSet(d)[path]
	if !ok {
		t.Fatalf("delta %+v lacks %s", d, path)
	}
	if c.Op != op {
		t.Fatalf("%s op = %q, want %q", path, c.Op, op)
	}
	if op == "deleted" {
		if c.Size != nil || c.SHA256 != "" {
			t.Fatalf("deleted %s carries size/sha: %+v", path, c)
		}
		return
	}
	if c.Size == nil || *c.Size != int64(len(content)) || c.SHA256 != shaHex(content) {
		t.Fatalf("%s = %+v, want size %d sha %s", path, c, len(content), shaHex(content))
	}
}

// AC1/AC2/AC4 (unit): created / modified / deleted on the git basis; an
// unchanged dirty path, a rewrite with identical bytes and a gitignored path
// are not changes.
func TestWorktreeDeltaGitBasis(t *testing.T) {
	root := gitRepo(t)
	d := deltaAround(t, root, func() {
		mustWrite(t, filepath.Join(root, "new.txt"), "fresh")
		mustWrite(t, filepath.Join(root, "a.txt"), "alpha changed\n")
		if err := os.Remove(filepath.Join(root, "b.txt")); err != nil {
			t.Fatal(err)
		}
		mustWrite(t, filepath.Join(root, "u.txt"), "untracked before\n") // same bytes
		mustMkdir(t, filepath.Join(root, "venv", "lib"))
		mustWrite(t, filepath.Join(root, "venv", "lib", "x.py"), "ignored")
		mustMkdir(t, filepath.Join(root, "sub dir"))
		mustWrite(t, filepath.Join(root, "sub dir", "q\"uote\ttab.txt"), "odd") // -z: no quoting
	})
	if d == nil || d.Basis != WorktreeBasisGit || d.Error != "" || d.Truncated || d.Total != 4 || len(d.Changes) != 4 {
		t.Fatalf("delta = %+v", d)
	}
	assertChange(t, d, "new.txt", "created", "fresh")
	assertChange(t, d, "a.txt", "modified", "alpha changed\n")
	assertChange(t, d, "b.txt", "deleted", "")
	assertChange(t, d, "sub dir/q\"uote\ttab.txt", "created", "odd")
	for _, c := range d.Changes {
		if strings.HasPrefix(c.Path, "venv") || c.Path == "u.txt" {
			t.Fatalf("listed %s: a gitignored or unchanged path is not a change", c.Path)
		}
	}
	// Sorted by path.
	for i := 1; i < len(d.Changes); i++ {
		if d.Changes[i-1].Path >= d.Changes[i].Path {
			t.Fatalf("changes not sorted: %v", d.Changes)
		}
	}
	// The wire shape (D1).
	b, _ := json.Marshal(d)
	if !bytes.Contains(b, []byte(`"basis":"git-status"`)) || !bytes.Contains(b, []byte(`"truncated":false`)) ||
		!bytes.Contains(b, []byte(`"total":4`)) {
		t.Fatalf("wire shape %s", b)
	}
}

// A path dirty before and reverted (or removed) during the call is a change
// even though git no longer lists it.
func TestWorktreeDeltaSeesPathsGitStopsListing(t *testing.T) {
	root := gitRepo(t)
	mustWrite(t, filepath.Join(root, "a.txt"), "dirty\n")
	d := deltaAround(t, root, func() {
		mustWrite(t, filepath.Join(root, "a.txt"), "alpha\n") // back to HEAD
		if err := os.Remove(filepath.Join(root, "u.txt")); err != nil {
			t.Fatal(err)
		}
	})
	if d == nil || d.Total != 2 {
		t.Fatalf("delta = %+v", d)
	}
	assertChange(t, d, "a.txt", "modified", "alpha\n")
	assertChange(t, d, "u.txt", "deleted", "")
}

func TestWorktreeDeltaNoChangeIsNil(t *testing.T) {
	root := gitRepo(t)
	if d := deltaAround(t, root, func() {}); d != nil {
		t.Fatalf("an unchanged git worktree gave %+v", d)
	}
	plain := canonicalTempDir(t)
	mustWrite(t, filepath.Join(plain, "x"), "x")
	if d := deltaAround(t, plain, func() {}); d != nil {
		t.Fatalf("an unchanged plain dir gave %+v", d)
	}
}

// AC3: over 256 changes ⇒ 256 entries, truncated, total exact.
func TestWorktreeDeltaTruncatesAt256(t *testing.T) {
	for _, git := range []bool{true, false} {
		root := canonicalTempDir(t)
		if git {
			root = gitRepo(t)
		}
		d := deltaAround(t, root, func() {
			for i := 0; i < 300; i++ {
				mustWrite(t, filepath.Join(root, fmt.Sprintf("f%03d.txt", i)), "x")
			}
		})
		if d == nil || len(d.Changes) != WorktreeDeltaMaxChanges || !d.Truncated || d.Total != 300 {
			t.Fatalf("git=%v: %d changes truncated=%v total=%d", git, len(d.Changes), d.Truncated, d.Total)
		}
		if d.Changes[0].Path != "f000.txt" || d.Changes[255].Path != "f255.txt" {
			t.Fatalf("git=%v: kept %s..%s, want the first 256 by path", git, d.Changes[0].Path, d.Changes[255].Path)
		}
	}
}

// AC4: outside git the walk is the basis; .git dirs are skipped; over the cap
// the delta says walk-capped instead of listing a partial answer.
func TestWorktreeDeltaWalkBasisAndCap(t *testing.T) {
	root := canonicalTempDir(t)
	mustWrite(t, filepath.Join(root, "keep.txt"), "k")
	mustWrite(t, filepath.Join(root, "mod.txt"), "m")
	mustWrite(t, filepath.Join(root, "gone.txt"), "g")
	d := deltaAround(t, root, func() {
		mustWrite(t, filepath.Join(root, "mod.txt"), "m2")
		if err := os.Remove(filepath.Join(root, "gone.txt")); err != nil {
			t.Fatal(err)
		}
		mustMkdir(t, filepath.Join(root, "d", "e"))
		mustWrite(t, filepath.Join(root, "d", "e", "new.txt"), "n")
	})
	if d == nil || d.Basis != WorktreeBasisWalk || d.Total != 3 {
		t.Fatalf("delta = %+v", d)
	}
	assertChange(t, d, "mod.txt", "modified", "m2")
	assertChange(t, d, "gone.txt", "deleted", "")
	assertChange(t, d, "d/e/new.txt", "created", "n")

	old := worktreeWalkCap
	worktreeWalkCap = 3
	t.Cleanup(func() { worktreeWalkCap = old })
	d = deltaAround(t, root, func() { mustWrite(t, filepath.Join(root, "more.txt"), "x") })
	if d == nil || d.Basis != WorktreeBasisWalkCapped || d.CappedAt != 3 || len(d.Changes) != 0 || d.Error != "" {
		t.Fatalf("over the cap: %+v, want walk-capped with nothing listed", d)
	}
	// Even an unchanged call says so (never silent).
	if d = deltaAround(t, root, func() {}); d == nil || d.Basis != WorktreeBasisWalkCapped {
		t.Fatalf("unchanged over the cap: %+v", d)
	}
}

// A symlink is stated, never followed: its sha is of the target string, and
// a path under a symlinked directory leading out is never read.
func TestWorktreeDeltaNeverFollowsSymlinks(t *testing.T) {
	base := canonicalTempDir(t)
	root, outside := filepath.Join(base, "root"), filepath.Join(base, "outside")
	mustMkdir(t, root, outside)
	mustWrite(t, filepath.Join(outside, "secret.txt"), "secret")
	d := deltaAround(t, root, func() {
		if err := os.Symlink(filepath.Join(outside, "secret.txt"), filepath.Join(root, "link")); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, filepath.Join(root, "dirlink")); err != nil {
			t.Fatal(err)
		}
	})
	if d == nil || d.Total != 2 {
		t.Fatalf("delta = %+v", d)
	}
	link := changeSet(d)["link"]
	target := filepath.Join(outside, "secret.txt")
	if link.Kind != "symlink" || link.SHA256 != shaHex(target) || *link.Size != int64(len(target)) {
		t.Fatalf("link = %+v", link)
	}
	if st := statePath(root, "dirlink/secret.txt"); st.kind != "outside" || st.sha != "" {
		t.Fatalf("a path through a symlinked dir was stated %+v, want outside and unread", st)
	}
}

// A failed snapshot is recorded, never swallowed and never a call failure.
func TestWorktreeDeltaSnapshotErrorIsRecorded(t *testing.T) {
	root := canonicalTempDir(t)
	// A .git pointer to nowhere: git refuses the repository.
	mustWrite(t, filepath.Join(root, ".git"), "gitdir: /nonexistent/x\n")
	d := deltaAround(t, root, func() {})
	if d == nil || d.Error == "" || len(d.Changes) != 0 || d.Basis != WorktreeBasisGit {
		t.Fatalf("delta = %+v, want a recorded error", d)
	}
	if len(d.Error) > worktreeDeltaMaxError {
		t.Fatalf("error is %d bytes, over the bound", len(d.Error))
	}
	// No git binary at all.
	old := worktreeGit
	worktreeGit = func() (string, error) { return "", errors.New("no git here") }
	t.Cleanup(func() { worktreeGit = old })
	repo := gitRepo(t)
	if d := deltaAround(t, repo, func() {}); d == nil || !strings.Contains(d.Error, "no git here") {
		t.Fatalf("delta = %+v, want the missing-git error", d)
	}
}

// ---------------------------------------------------------------------------
// through the handlers (binary-independent)
// ---------------------------------------------------------------------------

func decodeDelta(t *testing.T, out []byte) *WorktreeDelta {
	t.Helper()
	var m map[string]json.RawMessage
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatal(err)
	}
	raw, ok := m["worktree_delta"]
	if !ok {
		return nil
	}
	var d WorktreeDelta
	if err := json.Unmarshal(raw, &d); err != nil {
		t.Fatal(err)
	}
	return &d
}

// AC1 (fake tool): a run that writes a file records it; a run that writes
// nothing has no worktree_delta and is byte-identical to the pre-row-184
// result (pinned as a literal).
func TestAilangRunWorktreeDeltaOnTheFakeTool(t *testing.T) {
	h, tool, l := newFakeToolHandler(t, "admitted")
	run := func() []byte {
		out, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectAilangRun, Scope: WorkspaceScope, Cost: 1},
			[]byte(`{"path":"m.ail"}`))
		if err != nil {
			t.Fatal(err)
		}
		return out
	}
	want := `{"admitted":true,"exit_code":0,"decision":{"ok":true,"function":"main"},"limit":null,"stdout":"out\n",` +
		`"stderr":"policy: {\"ok\":true,\"decision\":{\"ok\":true,\"function\":\"main\"}}\n",` +
		`"policy":{"digest":"fakedigest","security_mode":"restricted","caps":["FS","IO"],"net_allow":[]}}`
	if got := run(); string(got) != want {
		t.Fatalf("unchanged run =\n%s\nwant (pre-row-184 bytes)\n%s", got, want)
	}
	tool.setMode(t, "writes")
	d := decodeDelta(t, run())
	if d == nil || d.Basis != WorktreeBasisWalk || d.Total != 1 {
		t.Fatalf("delta = %+v", d)
	}
	assertChange(t, d, "made.txt", "created", "made by run\n")
	_ = l
}

// AC1 (Ailang.Check/CLI): policy-tool ops carry the delta; Workspace.Write
// (an explicit file effect) never does.
func TestPolicyToolWorktreeDeltaOnlyForToolchainOps(t *testing.T) {
	h, tool, l := newFakeToolHandler(t, "admitted")
	call := func(effect, body string) []byte {
		out, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: effect, Scope: WorkspaceScope, Cost: 1}, []byte(body))
		if err != nil {
			t.Fatalf("%s: %v", effect, err)
		}
		return out
	}
	if out := call(EffectAilangCheck, `{"op":"ai_check","path":"m.ail"}`); decodeDelta(t, out) != nil {
		t.Fatalf("an unchanged check carries a delta: %s", out)
	}
	tool.setMode(t, "ptwrites")
	d := decodeDelta(t, call(EffectAilangCheck, `{"op":"ai_check","path":"m.ail"}`))
	if d == nil {
		t.Fatal("a check that wrote the worktree has no delta")
	}
	assertChange(t, d, "pt-1.txt", "created", "pt\n")
	if d = decodeDelta(t, call(EffectAilangCLI, `{"op":"tree"}`)); d == nil {
		t.Fatal("an ailang-cli op that wrote the worktree has no delta")
	}
	assertChange(t, d, "pt-2.txt", "created", "pt\n")
	// The fake writes pt-3.txt on this call too, so the absence is not vacuous.
	if out := call(EffectWorkspaceWrite, `{"op":"write","path":"w.txt","content":"x"}`); decodeDelta(t, out) != nil {
		t.Fatalf("Workspace.Write carries a delta: %s", out)
	}
	if _, err := os.Stat(filepath.Join(l.root, "pt-3.txt")); err != nil {
		t.Fatalf("control: the Write call did not write the worktree: %v", err)
	}
	_ = l
}

// A snapshot failure is recorded on the result, and the call still succeeds.
func TestAilangRunWorktreeDeltaErrorIsRecorded(t *testing.T) {
	h, _, l := newFakeToolHandler(t, "admitted")
	mustWrite(t, filepath.Join(l.root, ".git"), "gitdir: /nonexistent/x\n")
	out, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectAilangRun, Scope: WorkspaceScope, Cost: 1},
		[]byte(`{"path":"m.ail"}`))
	if err != nil {
		t.Fatalf("a snapshot failure failed the call: %v", err)
	}
	if d := decodeDelta(t, out); d == nil || d.Error == "" {
		t.Fatalf("result %s does not record the snapshot error", out)
	}
}

// driveWorktreeSnapshot is the AC10 census driver for the snapshot's git
// child: the probe stands in for git; its rev-parse answer is not the root,
// so the snapshot records an error — the launch is the measurement.
func driveWorktreeSnapshot(t *testing.T, probe string) {
	t.Helper()
	old := worktreeGit
	worktreeGit = func() (string, error) { return probe, nil }
	t.Cleanup(func() { worktreeGit = old })
	t.Setenv("GIT_DIR", filepath.Join(t.TempDir(), "decoy.git"))
	root := canonicalTempDir(t)
	mustMkdir(t, filepath.Join(root, ".git"))
	if s := snapshotWorktree(boundedTestContext(t), root, nil); s.err == nil {
		t.Fatal("the probe's rev-parse answer was accepted as the worktree root")
	}
}

// The snapshot's git child sees exactly World's minimal environment: no
// GIT_* redirect from the daemon (a GIT_DIR decoy) and no global config.
func TestWorktreeSnapshotGitEnvIsMinimal(t *testing.T) {
	dump := filepath.Join(t.TempDir(), "env")
	driveWorktreeSnapshot(t, writeProbeScript(t, dump))
	var got []string
	for _, kv := range readEnvDump(t, dump) {
		switch name, _, _ := strings.Cut(kv, "="); name {
		case "PWD", "OLDPWD", "SHLVL", "_": // added by the probe's /bin/sh
		default:
			got = append(got, kv)
		}
	}
	sort.Strings(got)
	// The exact set, written out here (not gitEnv() compared with itself).
	want := []string{"GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1", "GIT_OPTIONAL_LOCKS=0",
		"GIT_TERMINAL_PROMPT=0", "LC_ALL=C", "PATH=/usr/bin:/bin"}
	if !slices.Equal(got, want) {
		t.Fatalf("git child env = %q, want exactly %q (no GIT_DIR decoy, no credential)", got, want)
	}
}

// The git argv carries the fsmonitor override, and a hostile core.fsmonitor
// hook in the repository's own config is not run by the snapshot (control:
// a plain `git status` in the same repository runs it).
func TestWorktreeSnapshotDoesNotRunAnFsmonitorHook(t *testing.T) {
	root := gitRepo(t)
	marker := filepath.Join(canonicalTempDir(t), "fsmonitor-ran")
	hook := filepath.Join(canonicalTempDir(t), "hook.sh")
	if err := os.WriteFile(hook, []byte("#!/bin/sh\n: > "+shellQuote(marker)+"\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	fixtureGit(t, root, "config", "core.fsmonitor", hook)
	fixtureGit(t, root, "status", "--porcelain")
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("control: plain git status did not run the fsmonitor hook (%v); the arm below would prove nothing", err)
	}
	if err := os.Remove(marker); err != nil {
		t.Fatal(err)
	}
	d := deltaAround(t, root, func() { mustWrite(t, filepath.Join(root, "new.txt"), "n") })
	if _, err := os.Stat(marker); err == nil {
		t.Fatal("the snapshot ran the repository's fsmonitor hook")
	}
	if d == nil || d.Error != "" || d.Total != 1 {
		t.Fatalf("delta = %+v", d)
	}
}

// r2 hardening: a configured filter.<x>.clean would be run by `git status`
// on a path .gitattributes selects (the judge's probe). The snapshot refuses
// instead: the command never runs, and the refusal is recorded.
func TestWorktreeSnapshotRefusesAConfiguredFilter(t *testing.T) {
	root := gitRepo(t)
	marker := filepath.Join(canonicalTempDir(t), "filter-ran")
	fixtureGit(t, root, "config", "filter.evil.clean", "sh -c ': > "+marker+"; cat'")
	mustWrite(t, filepath.Join(root, ".gitattributes"), "* filter=evil\n")
	mustWrite(t, filepath.Join(root, "a.txt"), "dirty for the control\n")
	fixtureGit(t, root, "status", "--porcelain")
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("control: plain git status did not run the clean filter (%v); the arm below would prove nothing", err)
	}
	if err := os.Remove(marker); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(root, "a.txt"), "dirty again, so status must clean it\n")
	d := deltaAround(t, root, func() { mustWrite(t, filepath.Join(root, "b.txt"), "changed") })
	if _, err := os.Stat(marker); err == nil {
		t.Fatal("the snapshot ran the configured clean filter")
	}
	if d == nil || len(d.Changes) != 0 || !strings.Contains(d.Error, WorktreeFilterRefused) ||
		!strings.Contains(d.Error, "filter.evil.clean") {
		t.Fatalf("delta = %+v, want the recorded filter refusal", d)
	}
	// Through a handler: recorded on the result, the call still succeeds.
	h, _, l := newFakeToolHandler(t, "admitted")
	fixtureGit(t, l.root, "init", "-q", ".")
	fixtureGit(t, l.root, "config", "filter.evil.process", "sh -c ': > "+marker+"'")
	out, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectAilangRun, Scope: WorkspaceScope, Cost: 1},
		[]byte(`{"path":"m.ail"}`))
	if err != nil {
		t.Fatal(err)
	}
	if d := decodeDelta(t, out); d == nil || !strings.Contains(d.Error, "filter.evil.process") {
		t.Fatalf("result %s does not record the filter refusal", out)
	}
}

// Boundaries: exactly 256 changes is not truncated; exactly the walk cap is
// still a walk.
func TestWorktreeDeltaBoundaries(t *testing.T) {
	root := canonicalTempDir(t)
	d := deltaAround(t, root, func() {
		for i := 0; i < WorktreeDeltaMaxChanges; i++ {
			mustWrite(t, filepath.Join(root, fmt.Sprintf("f%03d", i)), "x")
		}
	})
	if d == nil || len(d.Changes) != 256 || d.Truncated || d.Total != 256 {
		t.Fatalf("256 changes: %d listed truncated=%v total=%d", len(d.Changes), d.Truncated, d.Total)
	}
	old := worktreeWalkCap
	worktreeWalkCap = 257
	t.Cleanup(func() { worktreeWalkCap = old })
	d = deltaAround(t, root, func() { mustWrite(t, filepath.Join(root, "f256"), "x") }) // 256 -> 257 = the cap
	if d == nil || d.Basis != WorktreeBasisWalk || d.Total != 1 {
		t.Fatalf("exactly the cap: %+v, want a walk", d)
	}
	d = deltaAround(t, root, func() { mustWrite(t, filepath.Join(root, "f257"), "x") }) // 258 > the cap
	if d == nil || d.Basis != WorktreeBasisWalkCapped {
		t.Fatalf("one over the cap: %+v, want walk-capped", d)
	}
}

// The after snapshot's own failure is recorded as such.
func TestWorktreeDeltaAfterSnapshotErrorIsRecorded(t *testing.T) {
	root := gitRepo(t)
	gitDir := filepath.Join(canonicalTempDir(t), "moved.git")
	d := deltaAround(t, root, func() {
		if err := os.Rename(filepath.Join(root, ".git"), gitDir); err != nil {
			t.Fatal(err)
		}
		mustWrite(t, filepath.Join(root, ".git"), "gitdir: /nonexistent/x\n")
	})
	if d == nil || !strings.HasPrefix(d.Error, "worktree snapshot after the call failed: ") || len(d.Changes) != 0 {
		t.Fatalf("delta = %+v, want the after-snapshot error", d)
	}
}

func TestWorktreeDeltaErrorIsBoundedToExactly512(t *testing.T) {
	msg := boundedDeltaError("after", errors.New(strings.Repeat("x", 2000)))
	if len(msg) != worktreeDeltaMaxError || worktreeDeltaMaxError != 512 || !strings.HasSuffix(msg, "...") {
		t.Fatalf("bounded error is %d bytes (suffix %q), want exactly 512 ending ...", len(msg), msg[len(msg)-3:])
	}
	if short := boundedDeltaError("before", errors.New("boom")); short != "worktree snapshot before the call failed: boom" {
		t.Fatalf("short error = %q", short)
	}
}

// The after snapshot is detached from the call's cancellation: a call whose
// context is cancelled by the time the child exits still records its delta.
func TestWorktreeDeltaAfterSnapshotIgnoresCancellation(t *testing.T) {
	root := gitRepo(t)
	ctx, cancel := context.WithCancel(boundedTestContext(t))
	p := beginWorktreeDelta(ctx, root)
	mustWrite(t, filepath.Join(root, "new.txt"), "n")
	cancel()
	d := p.finish(ctx)
	if d == nil || d.Error != "" || d.Total != 1 {
		t.Fatalf("delta after cancellation = %+v, want the change recorded", d)
	}
}

// The walk never descends a nested .git directory.
func TestWorktreeWalkSkipsNestedGitDirs(t *testing.T) {
	root := canonicalTempDir(t)
	mustMkdir(t, filepath.Join(root, "vendor", "dep", ".git"))
	mustWrite(t, filepath.Join(root, "vendor", "dep", ".git", "HEAD"), "ref: a\n")
	d := deltaAround(t, root, func() {
		mustWrite(t, filepath.Join(root, "vendor", "dep", ".git", "HEAD"), "ref: b\n")
		mustWrite(t, filepath.Join(root, "vendor", "dep", ".git", "index"), "i")
		mustWrite(t, filepath.Join(root, "vendor", "dep", "x.go"), "x")
	})
	if d == nil || d.Basis != WorktreeBasisWalk || d.Total != 1 || d.Changes[0].Path != "vendor/dep/x.go" {
		t.Fatalf("delta = %+v, want only vendor/dep/x.go", d)
	}
}

// A policy-tool response that itself carries worktree_delta is refused: the
// key is the handler's, never the tool's.
func TestPolicyToolResponseCannotForgeWorktreeDelta(t *testing.T) {
	h, _, _ := newFakeToolHandler(t, "ptreserved")
	_, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectAilangCheck, Scope: WorkspaceScope, Cost: 1},
		[]byte(`{"op":"ai_check","path":"m.ail"}`))
	if err == nil || !strings.Contains(err.Error(), `reserved key "worktree_delta"`) {
		t.Fatalf("a forged worktree_delta was accepted: %v", err)
	}
}

// r2 option 4: a file over the hash cap is listed with its size and
// hashed:false (no sha256), and a same-size rewrite is still seen (mtime).
func TestWorktreeDeltaHashCap(t *testing.T) {
	old := worktreeHashMaxBytes
	worktreeHashMaxBytes = 8
	t.Cleanup(func() { worktreeHashMaxBytes = old })
	root := canonicalTempDir(t)
	big := filepath.Join(root, "big.bin")
	d := deltaAround(t, root, func() {
		mustWrite(t, big, "0123456789")
		mustWrite(t, filepath.Join(root, "small"), "01234567") // exactly the cap: hashed
	})
	c, s := changeSet(d)["big.bin"], changeSet(d)["small"]
	if c.Hashed == nil || *c.Hashed || c.SHA256 != "" || c.Size == nil || *c.Size != 10 {
		t.Fatalf("big = %+v, want size 10, hashed:false, no sha256", c)
	}
	if s.Hashed != nil || s.SHA256 != shaHex("01234567") {
		t.Fatalf("small = %+v, want hashed", s)
	}
	b, _ := json.Marshal(c)
	if !strings.Contains(string(b), `"hashed":false`) {
		t.Fatalf("wire %s lacks hashed:false", b)
	}
	d = deltaAround(t, root, func() {
		mustWrite(t, big, "9876543210")
		future := time.Now().Add(time.Hour)
		if err := os.Chtimes(big, future, future); err != nil {
			t.Fatal(err)
		}
	})
	if d == nil || changeSet(d)["big.bin"].Op != "modified" {
		t.Fatalf("same-size rewrite of an unhashed file: %+v", d)
	}
}

// r2 option 4: each path is bounded once JSON-encoded, so even 256 hostile
// paths (long, all escape-expanding control bytes) stay far under the 1 MiB
// coordinator output cap.
func TestWorktreeDeltaWorstCaseFitsTheOutputCap(t *testing.T) {
	before := worktreeSnapshot{basis: WorktreeBasisWalk, entries: map[string]pathState{}}
	after := worktreeSnapshot{basis: WorktreeBasisWalk, entries: map[string]pathState{}}
	for i := 0; i < 300; i++ {
		p := fmt.Sprintf("%03d/", i) + strings.Repeat("\x01é", 4000)
		after.entries[p] = pathState{exists: true, size: 1 << 40, sha: strings.Repeat("f", 64), kind: "unreadable"}
	}
	d := worktreeDeltaOf(before, after)
	b, _ := json.Marshal(d)
	if len(d.Changes) != 256 || len(b) > 512<<10 {
		t.Fatalf("worst case: %d changes, %d bytes, want 256 within 512 KiB", len(d.Changes), len(b))
	}
	for _, c := range d.Changes {
		pb, _ := json.Marshal(c.Path)
		if !c.PathTruncated || len(pb) > worktreeDeltaMaxPathJSON || !utf8.ValidString(c.Path) {
			t.Fatalf("path %d JSON bytes, truncated=%v, valid=%v", len(pb), c.PathTruncated, utf8.ValidString(c.Path))
		}
	}
	if p, cut := boundPathJSON("ok/path.txt"); cut || p != "ok/path.txt" {
		t.Fatalf("a short path was cut: %q", p)
	}
	t.Logf("worst-case delta: %d bytes", len(b))
}

// ---------------------------------------------------------------------------
// Workspace.Exec (AC2, AC3, AC4)
// ---------------------------------------------------------------------------

func deltaCommands(m map[string]any) {
	addCommands(m, map[string]any{
		"mod":  map[string]any{"argv": []any{"/bin/sh", "-c", "printf changed > tracked.txt"}},
		"del":  map[string]any{"argv": []any{"/bin/sh", "-c", "rm tracked.txt"}},
		"many": map[string]any{"argv": []any{"/bin/sh", "-c", "i=0; while [ $i -lt 300 ]; do : > m$i.txt; i=$((i+1)); done"}},
		"ign":  map[string]any{"argv": []any{"/bin/sh", "-c", "mkdir -p venv && printf x > venv/x"}},
	})
}

func gitWorktreeSetup(t *testing.T) func(f *execFixture) {
	return func(f *execFixture) {
		fixtureGit(t, f.worktree, "init", "-q", ".")
		mustWrite(t, filepath.Join(f.worktree, "tracked.txt"), "original\n")
		mustWrite(t, filepath.Join(f.worktree, ".gitignore"), "venv/\n")
		fixtureGit(t, f.worktree, "add", ".")
		fixtureGit(t, f.worktree, "commit", "-q", "-m", "init")
	}
}

func execDeltaArms(t *testing.T, f *execFixture) {
	t.Helper()
	// An unchanged call: the §4.5 key set, no worktree_delta (byte-identical).
	out, err := f.call(t, "test", "a")
	if err != nil {
		t.Fatal(err)
	}
	if d := decodeDelta(t, out); d != nil {
		t.Fatalf("an unchanged exec carries %+v", d)
	}
	r := f.mustRun(t, "mod")
	if r.WorktreeDelta == nil || r.WorktreeDelta.Basis != WorktreeBasisGit || r.WorktreeDelta.Total != 1 {
		t.Fatalf("mod: delta %+v (exit %s stderr %q)", r.WorktreeDelta, r.code(), r.Stderr)
	}
	assertChange(t, r.WorktreeDelta, "tracked.txt", "modified", "changed")
	r = f.mustRun(t, "del")
	if r.WorktreeDelta == nil || r.WorktreeDelta.Total != 1 {
		t.Fatalf("del: delta %+v", r.WorktreeDelta)
	}
	assertChange(t, r.WorktreeDelta, "tracked.txt", "deleted", "")
	r = f.mustRun(t, "ign")
	if r.WorktreeDelta != nil {
		t.Fatalf("a gitignored write is listed: %+v", r.WorktreeDelta)
	}
	if fileState(filepath.Join(f.worktree, "venv", "x")) != "x" {
		t.Fatal("control: the ignored write did not happen")
	}
	r = f.mustRun(t, "many")
	if d := r.WorktreeDelta; d == nil || len(d.Changes) != WorktreeDeltaMaxChanges || !d.Truncated || d.Total != 300 {
		t.Fatalf("many: delta %+v", d)
	}
}

// AC2/AC3/AC4 on the pass-through stub (the command runs unsandboxed in the
// worktree; the snapshot is the handler's own, host-side).
func TestExecHandlerWorktreeDelta(t *testing.T) {
	f := newExecFixture(t, execFixtureOpts{profile: deltaCommands, setup: gitWorktreeSetup(t)})
	execDeltaArms(t, f)
	// A non-git worktree uses the walk and says so.
	g := newExecFixture(t, execFixtureOpts{profile: deltaCommands, setup: func(f *execFixture) {
		mustWrite(t, filepath.Join(f.worktree, "tracked.txt"), "original\n")
	}})
	r := g.mustRun(t, "mod")
	if r.WorktreeDelta == nil || r.WorktreeDelta.Basis != WorktreeBasisWalk {
		t.Fatalf("non-git: %+v", r.WorktreeDelta)
	}
	assertChange(t, r.WorktreeDelta, "tracked.txt", "modified", "changed")
}

// The same arms through the real, pinned srt (runs in CI's go job).
func TestExecSrtWorktreeDelta(t *testing.T) {
	f := newSrtFixture(t, func(m map[string]any, _ *execFixture) { deltaCommands(m) }, gitWorktreeSetup(t))
	execDeltaArms(t, f)
}
