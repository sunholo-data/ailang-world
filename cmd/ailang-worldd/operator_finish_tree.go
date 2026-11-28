package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
	"golang.org/x/sys/unix"
)

const operatorFinishDirtyLimit = 16 << 20

var (
	errOperatorFinishFence  = errors.New("worktree fence: operator input required")
	errOperatorFinishReview = errors.New("finish inputs changed; request fresh review")
	errOperatorCommitNoRef  = errors.New("local commit stopped before index or ref")
)

type operatorFinishGit struct {
	bin, dir, admin, common string
	env                     []string
	lock                    *os.File
}

func (g operatorFinishGit) git(ctx context.Context, input io.Reader, extra []string, args ...string) (string, error) {
	argv := []string{"--literal-pathspecs", "--git-dir=" + g.admin, "--work-tree=" + g.dir,
		"-c", "core.hooksPath=/dev/null", "-c", "core.fsmonitor=false", "-c", "core.attributesFile=/dev/null", "-c", "commit.gpgsign=false"}
	env := append(slices.Clone(g.env), extra...)
	var inherited []*os.File
	if g.lock != nil {
		inherited = append(inherited, g.lock)
	}
	return worktree.RunIsolated(ctx, g.bin, append(argv, args...), env, input, inherited...)
}

// Open beneath a held directory, and compare the file to its final pathname.
func operatorFinishRead(dir, name string, maxBytes int64) ([]byte, os.FileMode, error) {
	fd, err := unix.Open(dir, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, 0, err
	}
	parts := strings.Split(name, string(filepath.Separator))
	for i, part := range parts {
		flags := unix.O_RDONLY | unix.O_NOFOLLOW | unix.O_NONBLOCK | unix.O_CLOEXEC
		if i+1 < len(parts) {
			flags |= unix.O_DIRECTORY
		}
		next, err := unix.Openat(fd, part, flags, 0)
		unix.Close(fd)
		if err != nil {
			return nil, 0, err
		}
		fd = next
	}
	file := os.NewFile(uintptr(fd), name)
	defer file.Close()
	return readOperatorFinishFile(file, filepath.Join(dir, name), maxBytes)
}

func readOperatorFinishFile(file *os.File, path string, maxBytes int64) ([]byte, os.FileMode, error) {
	before, err := file.Stat()
	if err != nil || !before.Mode().IsRegular() || before.Size() > maxBytes {
		return nil, 0, errOperatorFinishFence
	}
	data, err := io.ReadAll(io.LimitReader(file, maxBytes+1))
	after, pathErr := os.Lstat(path)
	if err != nil || pathErr != nil || int64(len(data)) != before.Size() || int64(len(data)) > maxBytes ||
		!os.SameFile(before, after) || before.Mode() != after.Mode() ||
		before.Size() != after.Size() || before.ModTime() != after.ModTime() {
		return nil, 0, errOperatorFinishFence
	}
	return data, before.Mode(), nil
}

func operatorFinishPointer(dir, name, prefix string) (string, error) {
	b, _, err := operatorFinishRead(dir, name, 4096)
	path, ok := strings.CutPrefix(strings.TrimSuffix(string(b), "\n"), prefix)
	if err != nil || !ok || path == "" || strings.ContainsAny(path, "\x00\r\n") {
		return "", errOperatorFinishFence
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(dir, path)
	}
	return filepath.EvalSymlinks(path)
}

func openOperatorFinishGit(auth operatorFinishAuth, dir, repo string, lock *os.File) (operatorFinishGit, error) {
	var g operatorFinishGit
	dir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return g, errOperatorFinishFence
	}
	repo, err = filepath.EvalSymlinks(repo)
	if err != nil {
		return g, errOperatorFinishFence
	}
	admin, err := operatorFinishPointer(dir, ".git", "gitdir: ")
	if err != nil || filepath.Dir(admin) != filepath.Join(repo, ".git", "worktrees") {
		return g, errOperatorFinishFence
	}
	common, err := operatorFinishPointer(admin, "commondir", "")
	link, linkErr := operatorFinishPointer(admin, "gitdir", "")
	if err != nil || linkErr != nil || common != filepath.Join(repo, ".git") || link != filepath.Join(dir, ".git") {
		return g, errOperatorFinishFence
	}
	return operatorFinishGit{bin: auth.Git, dir: dir, admin: admin, common: common, lock: lock,
		env: []string{"HOME=" + auth.Home, "PATH=" + auth.Path, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null",
			"GIT_ATTR_NOSYSTEM=1", "GIT_OPTIONAL_LOCKS=0", "GIT_NO_REPLACE_OBJECTS=1", "GIT_NO_LAZY_FETCH=1", "GIT_TERMINAL_PROMPT=0"}}, nil
}

// Inspect Git configuration as data before running any worktree command.
func (g operatorFinishGit) config(ctx context.Context, start *store.EpisodeStart) (string, error) {
	for _, path := range []string{filepath.Join(g.admin, "config.worktree"), filepath.Join(g.common, "config.worktree"),
		filepath.Join(g.common, "objects", "info", "alternates"), filepath.Join(g.common, "info", "grafts")} {
		if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) {
			return "", errOperatorFinishFence
		}
	}
	data, _, err := operatorFinishRead(g.common, "config", 64<<10)
	if err != nil {
		return "", errOperatorFinishFence
	}
	raw, err := worktree.RunIsolated(ctx, g.bin, []string{"--git-dir=/dev/null", "config", "--file", "-",
		"--null", "--no-includes", "--list"}, g.env, strings.NewReader(string(data)))
	if err != nil {
		return "", errOperatorFinishFence
	}
	remote := 0
	for item := range strings.SplitSeq(strings.TrimSuffix(raw, "\x00"), "\x00") {
		name, value, _ := strings.Cut(item, "\n")
		switch {
		case name == "core.repositoryformatversion" && value == "0":
		case name == "core.bare" && value == "false":
		case slices.Contains([]string{"core.filemode", "core.logallrefupdates", "core.ignorecase", "core.precomposeunicode",
			"user.name", "user.email", "commit.gpgsign"}, name):
		case strings.HasPrefix(name, "branch.") && (strings.HasSuffix(name, ".remote") || strings.HasSuffix(name, ".merge")):
		case name == "remote."+start.Remote.Name+".fetch":
		case name == "remote."+start.Remote.Name+".url" && value == start.Remote.URL:
			remote++
		default:
			return "", errOperatorFinishFence
		}
	}
	if remote != 1 {
		return "", errOperatorFinishFence
	}
	current, _, err := operatorFinishRead(g.common, "config", 64<<10)
	if err != nil || !bytes.Equal(current, data) {
		return "", errOperatorFinishFence
	}
	return finishSHA(data), nil
}

func (g operatorFinishGit) snapshot(ctx context.Context, start *store.EpisodeStart) (store.OperatorFinishTree, [][]byte, error) {
	var tree store.OperatorFinishTree
	config, err := g.config(ctx, start)
	if err != nil {
		return tree, nil, err
	}
	branch, err := g.git(ctx, nil, nil, "symbolic-ref", "HEAD")
	if err != nil || strings.TrimSpace(branch) != "refs/heads/"+start.Branch {
		return tree, nil, errOperatorFinishFence
	}
	head, err := g.git(ctx, nil, nil, "rev-parse", "--verify", "HEAD")
	if err != nil || !operatorFinishSHA.MatchString(strings.TrimSpace(head)) {
		return tree, nil, errOperatorFinishFence
	}
	index, _, err := operatorFinishRead(g.admin, "index", 64<<20)
	if err != nil {
		return tree, nil, errOperatorFinishFence
	}
	staged, err := g.git(ctx, nil, nil, "diff", "--cached", "--name-only", "--no-ext-diff", "--no-textconv", "HEAD", "--")
	if err != nil || staged != "" {
		return tree, nil, errOperatorFinishFence
	}
	flags, err := g.git(ctx, nil, nil, "ls-files", "-v", "-z")
	if err != nil {
		return tree, nil, errOperatorFinishFence
	}
	trackedPaths := make(map[string]bool)
	for entry := range strings.SplitSeq(strings.TrimSuffix(flags, "\x00"), "\x00") {
		if entry == "" {
			continue
		}
		if !strings.HasPrefix(entry, "H ") {
			return tree, nil, errOperatorFinishFence
		}
		trackedPaths[entry[2:]] = true
	}
	tracked, err := g.git(ctx, nil, nil, "diff", "--name-only", "--no-renames", "--no-ext-diff", "--no-textconv", "-z", "HEAD", "--")
	if err != nil {
		return tree, nil, errOperatorFinishFence
	}
	var names []string
	if tracked != "" {
		names = strings.Split(strings.TrimSuffix(tracked, "\x00"), "\x00")
	}
	root, err := os.OpenRoot(g.dir)
	if err != nil {
		return tree, nil, errOperatorFinishFence
	}
	defer root.Close()
	names, err = operatorFinishUntracked(ctx, root.FS(), trackedPaths, names)
	if err != nil {
		return tree, nil, err
	}
	slices.Sort(names)
	names = slices.Compact(names)
	if len(names) > 257 {
		return tree, nil, errOperatorFinishFence
	}
	tree = store.OperatorFinishTree{Head: strings.TrimSpace(head), Branch: strings.TrimSpace(branch), Index: finishSHA(index), Config: config}
	var contents [][]byte
	budget := int64(operatorFinishDirtyLimit)
	for _, path := range names {
		if path == "" {
			continue
		}
		if !filepath.IsLocal(path) || filepath.Clean(path) != path || len(path) > 1024 || !utf8.ValidString(path) {
			return store.OperatorFinishTree{}, nil, errOperatorFinishFence
		}
		for part := range strings.SplitSeq(path, "/") {
			if strings.EqualFold(part, ".git") {
				return store.OperatorFinishTree{}, nil, errOperatorFinishFence
			}
		}
		b, mode, err := operatorFinishRead(g.dir, path, budget)
		p := store.OperatorFinishPath{Path: path, Mode: "100644", Digest: finishSHA(b)}
		switch {
		case errors.Is(err, os.ErrNotExist):
			p.Mode, p.Digest = "0", ""
		case err != nil:
			return store.OperatorFinishTree{}, nil, errOperatorFinishFence
		case mode&0100 != 0:
			p.Mode = "100755"
		}
		tree.Paths = append(tree.Paths, p)
		contents = append(contents, b)
		budget -= int64(len(b))
	}
	return tree, contents, nil
}

func operatorFinishUntracked(ctx context.Context, root fs.FS, trackedPaths map[string]bool, names []string) ([]string, error) {
	visited := 0
	err := fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) error {
		visited++
		if err != nil || ctx.Err() != nil || visited > 16384 {
			return errOperatorFinishFence
		}
		if path != ".git" && !d.IsDir() && !trackedPaths[path] {
			names = append(names, path)
		}
		return nil
	})
	if err != nil {
		return nil, errOperatorFinishFence
	}
	return names, nil
}

func (g operatorFinishGit) stage(ctx context.Context, spec store.OperatorLocalCommit, contents [][]byte, env []string) (string, error) {
	if _, err := g.git(ctx, nil, env, "read-tree", spec.Snapshot.Head); err != nil {
		return "", err
	}
	var entries strings.Builder
	for i, p := range spec.Snapshot.Paths {
		oid := strings.Repeat("0", 40)
		if p.Mode != "0" {
			out, err := g.git(ctx, strings.NewReader(string(contents[i])), env, "hash-object", "-w", "--no-filters", "--stdin")
			oid = strings.TrimSpace(out)
			if err != nil || !operatorFinishSHA.MatchString(oid) {
				return "", errOperatorFinishFence
			}
		}
		fmt.Fprintf(&entries, "%s %s\t%s\x00", p.Mode, oid, p.Path)
	}
	if _, err := g.git(ctx, strings.NewReader(entries.String()), env, "update-index", "-z", "--index-info"); err != nil {
		return "", err
	}
	out, err := g.git(ctx, nil, env, "write-tree")
	if !operatorFinishSHA.MatchString(strings.TrimSpace(out)) {
		return "", errOperatorFinishFence
	}
	return strings.TrimSpace(out), err
}

func (g operatorFinishGit) planTree(ctx context.Context, spool string, spec store.OperatorLocalCommit, contents [][]byte) (string, error) {
	dir, err := os.MkdirTemp(spool, "finish-tree.")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	env := []string{"GIT_INDEX_FILE=" + filepath.Join(dir, "index"), "GIT_OBJECT_DIRECTORY=" + dir,
		"GIT_ALTERNATE_OBJECT_DIRECTORIES=" + strconv.Quote(filepath.Join(g.common, "objects"))}
	return g.stage(ctx, spec, contents, env)
}

func (g operatorFinishGit) inspectCommit(ctx context.Context, spec store.OperatorLocalCommit, tree, agent string) (string, error) {
	branch, err := g.git(ctx, nil, nil, "symbolic-ref", "HEAD")
	if err != nil || strings.TrimSpace(branch) != spec.Snapshot.Branch {
		return "", errOperatorFinishFence
	}
	head, err := g.git(ctx, nil, nil, "rev-parse", "--verify", "HEAD")
	head = strings.TrimSpace(head)
	if err != nil || !operatorFinishSHA.MatchString(head) || head == spec.Snapshot.Head {
		return "", errOperatorFinishFence
	}
	raw, err := g.git(ctx, nil, nil, "cat-file", "commit", head)
	if err != nil || !operatorFinishCommitMatches(raw, spec, tree, agent) {
		return "", errOperatorFinishFence
	}
	index, err := g.git(ctx, nil, nil, "diff", "--cached", "--name-only", "--no-ext-diff", "--no-textconv", head, "--")
	if err != nil || index != "" {
		return "", errOperatorFinishFence
	}
	return head, nil
}

func operatorFinishCommitMatches(raw string, spec store.OperatorLocalCommit, tree, agent string) bool {
	header, message, found := strings.Cut(raw, "\n\n")
	lines := strings.Split(header, "\n")
	if !found || len(lines) != 4 || lines[0] != "tree "+tree ||
		lines[1] != "parent "+spec.Snapshot.Head || message != spec.Message {
		return false
	}
	author, authorTime, ok := strings.Cut(lines[2], "> ")
	if !ok || author+">" != "author "+agent+" via AILANG World <ailang-world@invalid>" || len(strings.Fields(authorTime)) != 2 {
		return false
	}
	committer, timeFields, ok := strings.Cut(lines[3], "> ")
	return ok && committer+">" == "committer "+spec.Committer && len(strings.Fields(timeFields)) == 2
}

func (g operatorFinishGit) commit(ctx context.Context, spool string, spec store.OperatorLocalCommit, contents [][]byte, tree string, start *store.EpisodeStart) error {
	if config, err := g.config(ctx, start); err != nil || config != spec.Snapshot.Config {
		return errOperatorFinishReview
	}
	name, email, err := worktree.ParseCommitter(spec.Committer)
	if err != nil || name == "" || strings.ContainsAny(start.Agent, "\r\n<>") {
		return errOperatorCommitNoRef
	}
	dir, err := os.MkdirTemp(spool, "finish-commit.")
	if err != nil {
		return errOperatorCommitNoRef
	}
	defer os.RemoveAll(dir)
	got, err := g.stage(ctx, spec, contents, []string{"GIT_INDEX_FILE=" + filepath.Join(dir, "index")})
	if err != nil || got != tree {
		return errOperatorCommitNoRef
	}
	env := []string{"GIT_COMMITTER_NAME=" + name, "GIT_COMMITTER_EMAIL=" + email,
		"GIT_AUTHOR_NAME=" + start.Agent + " via AILANG World", "GIT_AUTHOR_EMAIL=ailang-world@invalid"}
	out, err := g.git(ctx, strings.NewReader(spec.Message), env, "commit-tree", tree, "-p", spec.Snapshot.Head, "-F", "-")
	head := strings.TrimSpace(out)
	if err != nil || !operatorFinishSHA.MatchString(head) {
		return errOperatorCommitNoRef
	}
	current, _, err := g.snapshot(ctx, start)
	if err != nil || !reflect.DeepEqual(current, spec.Snapshot) {
		return errOperatorFinishReview
	}
	if _, err := g.git(ctx, nil, nil, "read-tree", tree); err != nil {
		return err
	}
	_, err = g.git(ctx, nil, nil, "update-ref", spec.Snapshot.Branch, head, spec.Snapshot.Head)
	return err
}
