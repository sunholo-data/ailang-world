package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/verbs"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

type operatorWorktreeTarget struct {
	Action, Repo, Common, Path, Branch, SHA string
}

type operatorStartSource struct {
	presets, profiles, workspace, hosts string
	clock                               func() int64
}

// reserveOperatorStartWorktree borrows the caller-owned episode lock.
func reserveOperatorStartWorktree(ctx context.Context, credentials string, held *store.OperatorEpisodeLock, signed approval.OperatorAction,
	attempt store.OperatorMintAttempt, src operatorStartSource, checkpoint func(string) error, process ...*os.File) error {
	canonical, err := store.CanonicalPath(credentials)
	if err != nil {
		return err
	}
	if signed.Store != canonical || !verbs.EpisodeIDPattern.MatchString(signed.Episode) {
		return store.ErrOperatorConflict
	}
	if _, err := held.Verify(canonical, signed.Episode); err != nil {
		return err
	}
	in, err := reconstructOperatorStart(ctx, signed, attempt, src.presets, src.profiles, src.workspace, src.hosts, src.clock)
	if err != nil {
		return err
	}
	git, err := sessionGit()
	if err != nil {
		return err
	}
	target, err := newOperatorWorktreeTarget(ctx, git, signed, in.Start)
	if err != nil {
		return err
	}
	initial, err := target.inspect(ctx, git)
	if err != nil {
		return err
	}
	if initial.state != treeAbsent {
		return fmt.Errorf("%w: worktree path, branch or registration already exists", store.ErrOperatorConflict)
	}
	if err := checkpoint("before reservation"); err != nil {
		return err
	}
	identity, _ := json.Marshal(target) // string fields only
	var key store.OperatorStepKey
	err = runnerWrite(canonical, func(db *store.Credentials) error {
		row, err := db.OperatorAction(ctx, signed.ID, src.clock())
		if err != nil {
			return err
		}
		if len(row.Challenge) != 32 {
			return approval.ErrNotVerified
		}
		proof, err := authority.RecheckClaimedOperatorAction(ctx, db, signed, [32]byte(row.Challenge), row.ExpiresAt, src.clock(), row)
		if err != nil {
			return err
		}
		key, err = db.ReserveOperatorWorktree(ctx, in.Attempt, string(identity), proof.OperatorConfirmation().MintProof(), src.clock)
		return err
	})
	if err != nil {
		return err
	}
	if err := checkpoint("before add"); err != nil {
		return fmt.Errorf("%w; %s", err, reconcileWorktreeCommand(key.Action.ID))
	}
	fd, err := held.Verify(canonical, signed.Episode)
	if err != nil {
		return fmt.Errorf("%w; %s", err, reconcileWorktreeCommand(key.Action.ID))
	}
	addErr := worktree.Add(ctx, git, target.Repo, target.Path, target.Branch, target.SHA, append([]*os.File{fd}, process...)...)
	if _, err := held.Verify(canonical, signed.Episode); err != nil {
		return fmt.Errorf("%w; %s", err, reconcileWorktreeCommand(key.Action.ID))
	}
	if err := checkpoint("before readback"); err != nil {
		return fmt.Errorf("%w; %s", err, reconcileWorktreeCommand(key.Action.ID))
	}
	got, readErr := target.inspect(ctx, git)
	if err := checkpoint("before observe"); err != nil {
		return fmt.Errorf("%w; %s", err, reconcileWorktreeCommand(key.Action.ID))
	}

	created := addErr == nil && readErr == nil && got.state == treeExact
	disposition := store.OperatorWorktreeUncertain
	if addErr != nil && readErr == nil && got.state == treeAbsent {
		disposition = store.OperatorWorktreeAbsent
	}
	err = runnerWrite(canonical, func(db *store.Credentials) error {
		if created {
			return db.ObserveOperatorStep(ctx, key, []byte("created"))
		}
		return db.DisposeOperatorWorktree(ctx, key, disposition, src.clock())
	})
	if err != nil {
		return fmt.Errorf("%w; %s", err, reconcileWorktreeCommand(key.Action.ID))
	}
	if created {
		return nil
	}
	if disposition == store.OperatorWorktreeAbsent {
		return errors.New("worktree add failed; worktree and branch absent")
	}
	return fmt.Errorf("worktree outcome uncertain; %s", reconcileWorktreeCommand(key.Action.ID))
}

func operatorWorktreeFact(a approval.OperatorAction) string {
	for _, f := range a.Resolution {
		if f.Name == "worktree_path" {
			return f.Value
		}
	}
	return ""
}

func operatorPhysicalTarget(path string) (string, error) {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return "", store.ErrOperatorConflict
	}
	parent, err := canonicalLaunchPath(filepath.Dir(path))
	if err != nil {
		return "", err
	}
	return filepath.Join(parent, filepath.Base(path)), nil
}

func newOperatorWorktreeTarget(ctx context.Context, git string, a approval.OperatorAction, start store.EpisodeStart) (operatorWorktreeTarget, error) {
	var target operatorWorktreeTarget
	repo, err := filepath.EvalSymlinks(start.Repo)
	if err != nil {
		return target, err
	}
	path, err := operatorPhysicalTarget(operatorWorktreeFact(a))
	if err != nil {
		return target, err
	}
	common, err := operatorGitCommon(ctx, git, repo)
	if err != nil {
		return target, err
	}
	return operatorWorktreeTarget{a.ID, repo, common, path, start.Branch, start.Base.SHA}, nil
}

func operatorGitCommon(ctx context.Context, git, dir string) (string, error) {
	out, err := worktree.Output(ctx, git, []string{"-C", dir, "rev-parse", "--path-format=absolute", "--git-common-dir"})
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(strings.TrimSuffix(out, "\n"))
}

const (
	treeAbsent = iota
	treePartial
	treeExact
)

type operatorWorktreeEvidence struct {
	state                            int
	registration, refs, common, head string
	mode                             os.FileMode
	dev, inode                       uint64
}

// inspect checks the repository even when the target is absent. Path absence
// alone cannot distinguish a prunable registration or a branch-only result.
func (target operatorWorktreeTarget) inspect(ctx context.Context, git string) (operatorWorktreeEvidence, error) {
	var e operatorWorktreeEvidence
	common, err := operatorGitCommon(ctx, git, target.Repo)
	if err != nil || common != target.Common {
		return e, errors.New("worktree repository identity differs")
	}
	list, err := worktree.Output(ctx, git, []string{"-C", target.Repo, "worktree", "list", "--porcelain", "-z"})
	if err != nil {
		return e, err
	}
	ref := "refs/heads/" + target.Branch
	want := "worktree " + target.Path + "\x00HEAD " + target.SHA + "\x00branch " + ref
	exact := 0
	for _, entry := range strings.Split(strings.TrimSuffix(list, "\x00\x00"), "\x00\x00") {
		path, _, _ := strings.Cut(entry, "\x00")
		raw, ok := strings.CutPrefix(path, "worktree ")
		if !ok || !filepath.IsAbs(raw) {
			return e, errors.New("invalid Git worktree registration")
		}
		physical, err := canonicalLaunchPath(raw)
		if err != nil {
			return e, err
		}
		if physical == target.Path || raw == target.Path || strings.Contains("\x00"+entry+"\x00", "\x00branch "+ref+"\x00") {
			e.registration += entry + "\x00\x00"
			if "worktree "+physical+strings.TrimPrefix(entry, path) == want {
				exact++
			}
		}
	}
	refs, err := worktree.Output(ctx, git, []string{"-C", target.Repo, "for-each-ref", "--format=%(refname) %(objectname)", "refs/heads"})
	if err != nil {
		return e, err
	}
	for _, line := range strings.Split(refs, "\n") {
		name, _, _ := strings.Cut(line, " ")
		if name == ref || strings.HasPrefix(name, ref+"/") || strings.HasPrefix(ref, name+"/") {
			e.refs += line + "\n"
		}
	}
	info, err := os.Lstat(target.Path)
	if errors.Is(err, os.ErrNotExist) {
		if e.registration != "" || e.refs != "" {
			e.state = treePartial
		}
		return e, nil
	}
	if err != nil {
		return e, err
	}
	stat := info.Sys().(*syscall.Stat_t)
	e.mode, e.dev, e.inode = info.Mode(), uint64(stat.Dev), stat.Ino
	e.state = treePartial
	if !info.IsDir() || exact != 1 || strings.Count(e.registration, "\x00\x00") != 1 || e.refs != ref+" "+target.SHA+"\n" {
		return e, nil
	}
	physical, err := filepath.EvalSymlinks(target.Path)
	if err != nil || physical != target.Path {
		return e, errors.New("worktree physical path differs")
	}
	link, err := os.Lstat(filepath.Join(target.Path, ".git"))
	if errors.Is(err, os.ErrNotExist) {
		return e, nil
	}
	if err != nil {
		return e, err
	}
	if !link.Mode().IsRegular() || !worktree.IsLinked(target.Path) {
		return e, nil
	}
	e.common, err = operatorGitCommon(ctx, git, target.Path)
	if err != nil || e.common != target.Common {
		return e, nil
	}
	e.head, err = worktree.Output(ctx, git, []string{"-C", target.Path, "rev-parse", "HEAD", "--symbolic-full-name", "HEAD"})
	if err == nil && e.head == target.SHA+"\n"+ref+"\n" {
		e.state = treeExact
	}
	return e, nil
}
