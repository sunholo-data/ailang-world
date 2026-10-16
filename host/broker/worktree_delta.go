package broker

// Row 184 (design_docs/planned/w-toolchain-worktree-delta.md §2 D1–D2): the
// files a toolchain call created, modified or deleted in the episode
// worktree, recorded in the handler result as `worktree_delta`.
//
// The broker snapshots the worktree host-side before the subprocess and again
// after it. In a git worktree the snapshot is what `git status --porcelain=v1
// -z --untracked-files=all --ignored=no` lists plus a sha256 of each listed
// path (D-184-1 = git-status basis): gitignored paths are out of scope, and
// the result says so (`basis`). Outside git it is a bounded lexical walk;
// over the cap the result says `walk-capped` and lists nothing — never a
// silent partial answer. A snapshot that fails never fails the call: the
// result carries `worktree_delta.error` instead.
//
// The git child is launched through runBounded (the one process lifecycle;
// the AC10 census drives it as driveWorktreeSnapshot) with a minimal
// environment: no registry credential, no GIT_* redirect from the daemon's
// environment, no system or global git config, and no fsmonitor hook.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"
)

const (
	// WorktreeDeltaMaxChanges bounds the listed changes (D1).
	WorktreeDeltaMaxChanges = 256
	// worktreeGitTimeout bounds each git invocation; worktreeGitMaxOutput
	// its stdout (20,000 paths of 4 KiB fit with room to spare).
	worktreeGitTimeout   = 30 * time.Second
	worktreeGitMaxOutput = 96 << 20
	// worktreeDeltaMaxError bounds a recorded snapshot error.
	worktreeDeltaMaxError = 512

	WorktreeBasisGit        = "git-status"
	WorktreeBasisWalk       = "walk"
	WorktreeBasisWalkCapped = "walk-capped"
)

// worktreeWalkCap bounds the non-git walk; worktreeGitCap bounds the paths
// git may list (each is hashed twice). Variables only so tests can lower them.
var (
	worktreeWalkCap = 20000
	worktreeGitCap  = 20000
)

// worktreeGit resolves the git binary the snapshot runs. Tests (and the AC10
// census) replace it.
var worktreeGit = func() (string, error) { return exec.LookPath("git") }

// WorktreeChange is one changed path. Size and SHA256 describe the after
// state and are absent for a deleted path; Kind is absent for a regular file
// and otherwise names what the path is (symlink: SHA256 is of the link
// target; dir: a nested repository git lists as one entry, not hashed;
// unreadable / outside: the path could not be hashed without leaving the
// worktree or failed to read — the change is still listed).
type WorktreeChange struct {
	Path   string `json:"path"`
	Op     string `json:"op"`
	Size   *int64 `json:"size,omitempty"`
	SHA256 string `json:"sha256,omitempty"`
	Kind   string `json:"kind,omitempty"`
}

// WorktreeDelta is the `worktree_delta` result field.
type WorktreeDelta struct {
	Changes   []WorktreeChange `json:"changes,omitempty"`
	Truncated bool             `json:"truncated"`
	Total     int              `json:"total"`
	Basis     string           `json:"basis"`
	// CappedAt is the walk's file cap when Basis is walk-capped.
	CappedAt int `json:"capped_at,omitempty"`
	// Error is a snapshot failure; nothing is listed with it.
	Error string `json:"error,omitempty"`
}

// pathState is one path's on-disk state at snapshot time.
type pathState struct {
	exists bool
	kind   string // "" regular file, "symlink", "dir", "unreadable", "outside"
	size   int64
	sha    string
}

type worktreeSnapshot struct {
	basis   string
	entries map[string]pathState
	codes   map[string]string // git XY per listed path (git basis only)
	err     error
}

// worktreeProbe is one call's before snapshot.
type worktreeProbe struct {
	root   string
	before worktreeSnapshot
}

// beginWorktreeDelta snapshots root before a toolchain subprocess.
func beginWorktreeDelta(ctx context.Context, root string) *worktreeProbe {
	return &worktreeProbe{root: root, before: snapshotWorktree(ctx, root, nil)}
}

// worktreeAfterBudget bounds the after snapshot's git invocations as a whole.
const worktreeAfterBudget = 2 * worktreeGitTimeout

// finish snapshots root again and returns the delta: nil when nothing
// changed. It runs outside the child's timeout and detached from the call's
// cancellation (the call already ran; its record should say what it did),
// bounded by worktreeAfterBudget. It makes no Store call.
func (p *worktreeProbe) finish(ctx context.Context) *WorktreeDelta {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), worktreeAfterBudget)
	defer cancel()
	if p.before.err != nil {
		return &WorktreeDelta{Basis: p.before.basis, Error: boundedDeltaError("before", p.before.err)}
	}
	also := make([]string, 0, len(p.before.entries))
	for path := range p.before.entries {
		also = append(also, path)
	}
	after := snapshotWorktree(ctx, p.root, also)
	return worktreeDeltaOf(p.before, after)
}

func boundedDeltaError(phase string, err error) string {
	msg := fmt.Sprintf("worktree snapshot %s the call failed: %v", phase, err)
	if len(msg) > worktreeDeltaMaxError {
		msg = msg[:worktreeDeltaMaxError-3] + "..."
	}
	return msg
}

// worktreeDeltaOf compares two snapshots of one root. A path listed both
// times in the same state is not a change.
func worktreeDeltaOf(before, after worktreeSnapshot) *WorktreeDelta {
	switch {
	case before.err != nil:
		return &WorktreeDelta{Basis: before.basis, Error: boundedDeltaError("before", before.err)}
	case after.err != nil:
		return &WorktreeDelta{Basis: after.basis, Error: boundedDeltaError("after", after.err)}
	case before.basis == WorktreeBasisWalkCapped || after.basis == WorktreeBasisWalkCapped:
		// Either side over the cap: nothing is listed, and the result says so.
		return &WorktreeDelta{Basis: WorktreeBasisWalkCapped, CappedAt: worktreeWalkCap}
	case before.basis != after.basis:
		return &WorktreeDelta{Basis: after.basis,
			Error: fmt.Sprintf("worktree snapshot basis changed during the call (%s before, %s after)", before.basis, after.basis)}
	}
	var changes []WorktreeChange
	add := func(path, op string, st pathState) {
		c := WorktreeChange{Path: path, Op: op, Kind: st.kind}
		if op != "deleted" {
			size := st.size
			c.Size, c.SHA256 = &size, st.sha
		}
		changes = append(changes, c)
	}
	for path, a := range after.entries {
		b, listed := before.entries[path]
		switch {
		case listed && b == a:
		case listed && !a.exists:
			if b.exists {
				add(path, "deleted", a)
			}
		case listed && !b.exists:
			add(path, "created", a)
		case listed:
			add(path, "modified", a)
		case !a.exists:
			// Listed only after and gone: a clean tracked file was deleted.
			add(path, "deleted", a)
		default:
			// Listed only after: untracked or newly staged means it did not
			// exist before; anything else was a clean tracked file.
			code := after.codes[path]
			if code == "??" || strings.HasPrefix(code, "A") || after.basis != WorktreeBasisGit {
				add(path, "created", a)
			} else {
				add(path, "modified", a)
			}
		}
	}
	if after.basis != WorktreeBasisGit {
		// The walk lists every path, so one listed only before is gone.
		for path, b := range before.entries {
			if _, still := after.entries[path]; !still && b.exists {
				add(path, "deleted", pathState{})
			}
		}
	}
	if len(changes) == 0 {
		return nil
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].Path < changes[j].Path })
	d := &WorktreeDelta{Total: len(changes), Basis: after.basis}
	if len(changes) > WorktreeDeltaMaxChanges {
		changes, d.Truncated = changes[:WorktreeDeltaMaxChanges], true
	}
	d.Changes = changes
	return d
}

// snapshotWorktree snapshots root: the git basis when root is the top of a
// git worktree (it has its own .git file or directory), else the walk. also
// names paths (root-relative, slash-separated) to state even when git no
// longer lists them — the before snapshot's paths, so a reverted or deleted
// path is seen.
func snapshotWorktree(ctx context.Context, root string, also []string) worktreeSnapshot {
	canon, err := filepath.EvalSymlinks(root)
	if err != nil {
		return worktreeSnapshot{basis: WorktreeBasisGit, err: fmt.Errorf("resolve root: %w", err)}
	}
	if _, err := os.Lstat(filepath.Join(canon, ".git")); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return worktreeSnapshot{basis: WorktreeBasisGit, err: err}
		}
		return walkSnapshot(canon)
	}
	return gitSnapshot(ctx, canon, also)
}

func gitEnv() []string {
	return []string{"PATH=/usr/bin:/bin", "LC_ALL=C", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_OPTIONAL_LOCKS=0", "GIT_TERMINAL_PROMPT=0"}
}

func runWorktreeGit(ctx context.Context, git, root string, args ...string) ([]byte, error) {
	stderr := newBoundedSink(4096)
	argv := append([]string{"-C", root, "-c", "core.fsmonitor=false", "-c", "core.quotePath=false"}, args...)
	out, err := runBounded(ctx, handlerBounds{execTimeout: worktreeGitTimeout, maxOutputBytes: worktreeGitMaxOutput},
		handlerCommand{path: git, args: argv, dir: root, env: gitEnv(), stderr: stderr})
	if err != nil {
		return nil, fmt.Errorf("git %s: %w (stderr %q)", args[0], err, stderr.Bytes())
	}
	return out, nil
}

func gitSnapshot(ctx context.Context, root string, also []string) worktreeSnapshot {
	s := worktreeSnapshot{basis: WorktreeBasisGit, entries: map[string]pathState{}, codes: map[string]string{}}
	git, err := worktreeGit()
	if err != nil {
		s.err = fmt.Errorf("root has a .git entry but no git binary: %w", err)
		return s
	}
	top, err := runWorktreeGit(ctx, git, root, "rev-parse", "--show-toplevel")
	if err != nil {
		s.err = err
		return s
	}
	if t, err := filepath.EvalSymlinks(strings.TrimSpace(string(top))); err != nil || t != root {
		s.err = fmt.Errorf("git toplevel %q is not the worktree root %q", strings.TrimSpace(string(top)), root)
		return s
	}
	out, err := runWorktreeGit(ctx, git, root, "status", "--porcelain=v1", "-z", "--untracked-files=all",
		"--ignored=no", "--no-renames")
	if err != nil {
		s.err = err
		return s
	}
	records := strings.Split(strings.TrimSuffix(string(out), "\x00"), "\x00")
	if len(out) == 0 {
		records = nil
	}
	for i := 0; i < len(records); i++ {
		rec := records[i]
		if len(rec) < 4 || rec[2] != ' ' {
			s.err = fmt.Errorf("git status record %q is not porcelain v1", rec)
			return s
		}
		code, path := rec[:2], strings.TrimSuffix(rec[3:], "/")
		s.codes[path] = code
		if code[0] == 'R' || code[0] == 'C' {
			// The source path follows as its own NUL-terminated field.
			if i+1 >= len(records) {
				s.err = fmt.Errorf("git status rename record %q has no source path", rec)
				return s
			}
			i++
			if _, seen := s.codes[records[i]]; !seen {
				s.codes[records[i]] = "  "
			}
		}
	}
	if len(s.codes) > worktreeGitCap {
		s.err = fmt.Errorf("git status lists %d paths, more than the %d the snapshot hashes", len(s.codes), worktreeGitCap)
		return s
	}
	for path := range s.codes {
		s.entries[path] = statePath(root, path)
	}
	for _, path := range also {
		if _, done := s.entries[path]; !done {
			s.entries[path] = statePath(root, path)
		}
	}
	return s
}

// walkSnapshot states every file and symlink under root in lexical order,
// never descending into a .git directory or through a symlink. Over the cap
// it lists nothing and says walk-capped.
func walkSnapshot(root string) worktreeSnapshot {
	s := worktreeSnapshot{basis: WorktreeBasisWalk, entries: map[string]pathState{}}
	var paths []string
	errCapped := errors.New("capped")
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if p == root {
			return nil
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return fs.SkipDir
			}
			return nil
		}
		if len(paths) == worktreeWalkCap {
			return errCapped
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		paths = append(paths, filepath.ToSlash(rel))
		return nil
	})
	switch {
	case errors.Is(err, errCapped):
		return worktreeSnapshot{basis: WorktreeBasisWalkCapped}
	case err != nil:
		s.err = err
		return s
	}
	for _, path := range paths {
		s.entries[path] = statePath(root, path)
	}
	return s
}

// statePath states one root-relative path without following a symlink: the
// final component is Lstat'ed and opened O_NOFOLLOW, and a path whose parent
// resolves outside root is reported as "outside" and never read.
func statePath(root, rel string) pathState {
	full := filepath.Join(root, filepath.FromSlash(rel))
	if !within(full, root) || full == root {
		return pathState{exists: true, kind: "outside"}
	}
	parent, err := filepath.EvalSymlinks(filepath.Dir(full))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return pathState{}
		}
		return pathState{exists: true, kind: "unreadable"}
	}
	if !within(parent, root) {
		return pathState{exists: true, kind: "outside"}
	}
	full = filepath.Join(parent, filepath.Base(full))
	info, err := os.Lstat(full)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return pathState{}
	case err != nil:
		return pathState{exists: true, kind: "unreadable"}
	case info.Mode()&fs.ModeSymlink != 0:
		target, err := os.Readlink(full)
		if err != nil {
			return pathState{exists: true, kind: "unreadable"}
		}
		sum := sha256.Sum256([]byte(target))
		return pathState{exists: true, kind: "symlink", size: int64(len(target)), sha: hex.EncodeToString(sum[:])}
	case info.IsDir():
		return pathState{exists: true, kind: "dir"}
	case !info.Mode().IsRegular():
		return pathState{exists: true, kind: "unreadable"}
	}
	f, err := os.OpenFile(full, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		return pathState{exists: true, kind: "unreadable"}
	}
	defer func() { _ = f.Close() }()
	h := sha256.New()
	n, err := io.Copy(h, f)
	if err != nil {
		return pathState{exists: true, kind: "unreadable"}
	}
	return pathState{exists: true, size: n, sha: hex.EncodeToString(h.Sum(nil))}
}
