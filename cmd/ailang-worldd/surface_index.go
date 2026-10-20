package main

// Row 188 M0a (design_docs/planned/w-world-os-surface-sprint-plan.md, "The
// surface index"): `ailang-worldd surface-index` builds the index /surface
// reads from a repo's git history, and checks one.
//
// It is the thin git side of host/surface: it runs git log and git
// merge-base, and everything else (parsing, branes, the index, its checks)
// is host/surface's pure code. The file is written atomically (a temporary
// file in the same directory, then a rename) and never inside the repo.

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

const surfaceIndexHelp = `usage: ailang-worldd surface-index --repo <dir> --charter <file> --out <file> [--rev dev] [--limit 200]
       ailang-worldd surface-index --validate <file>

Build the surface index (` + surface.Schema + `) that /surface reads, from a
repo's git history and the charter's queue rows, or check one.

Building runs, in --repo:
  git log --first-parent --merges --format=%H%x09%P%x09%cI%x09%s -n <limit> <rev>
  git merge-base <first parent> <second parent>    (once per merge)
Each merge is a shared state, and its second parent is a step (the PR's side).
A step's brane is the first row its PR branch names (trial/row189-…,
sprint/world-iter288-row156-…, design/row189-…, docs/row198-…,
trial/rows190-191-… → row 190); sprint/ branches are the loop's, the rest
attended, and a brane with both is mixed. A branch that names no row has no
brane: its step sits on the shared lane. Brane titles are the rows' slugs in
--charter. The index is validated, then written with a newline, atomically
(a temporary file beside --out, then a rename). It holds commit hashes, times,
PR numbers, branch names and row slugs, never a filesystem path.

--validate reads an index and checks it: the schema, unique ids, every parent
resolving (but the oldest shared state's first parent), every step's brane
listed, every step's merge pointing back at it, times in order along the shared
chain, the allowed kinds, and at most 5000 states. It exits 0, or 1 with the
reason.

Refused: a missing flag, --validate with any build flag, a --limit outside
1..2500, a --rev starting with -, and an --out inside --repo.

  --repo <dir>       the git repository to read (required to build)
  --charter <file>   the charter whose queue rows title the branes (required to build)
  --out <file>       the index to write (required to build)
  --rev <rev>        the branch whose first-parent merges are read (default dev)
  --limit N          read at most N merges (default 200)
  --validate <file>  check an index instead of building one
`

// surfaceIndexMaxBytes bounds an index --validate reads (the M1 serving bound).
const surfaceIndexMaxBytes = 4 << 20

// surfaceIndexBudget bounds the git calls of one build as a whole.
const surfaceIndexBudget = 2 * time.Minute

// surfaceIndexArgs is a parsed command line.
type surfaceIndexArgs struct {
	repo, charter, out, rev, validate string
	limit                             int
}

// parseSurfaceIndexArgs parses and checks the flags; help is flag.ErrHelp.
func parseSurfaceIndexArgs(args []string, stderr io.Writer) (surfaceIndexArgs, error) {
	var a surfaceIndexArgs
	fs := flag.NewFlagSet("ailang-worldd surface-index", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	fs.StringVar(&a.repo, "repo", "", "the git repository to read")
	fs.StringVar(&a.charter, "charter", "", "the charter whose queue rows title the branes")
	fs.StringVar(&a.out, "out", "", "the index to write")
	fs.StringVar(&a.rev, "rev", "dev", "the branch whose first-parent merges are read")
	fs.IntVar(&a.limit, "limit", 200, "read at most N merges")
	fs.StringVar(&a.validate, "validate", "", "check an index instead of building one")
	if err := fs.Parse(args); err != nil {
		return a, err
	}
	switch {
	case fs.NArg() > 0:
		return a, fmt.Errorf("unexpected argument %q", fs.Arg(0))
	case a.validate != "":
		build := false
		fs.Visit(func(f *flag.Flag) { build = build || f.Name != "validate" })
		if build {
			return a, errors.New("--validate takes no other flag")
		}
	case a.repo == "" || a.charter == "" || a.out == "":
		return a, errors.New("--repo, --charter and --out are required (or --validate <file>)")
	case a.limit < 1 || 2*a.limit > surface.MaxStates:
		return a, fmt.Errorf("--limit %d is outside 1..%d", a.limit, surface.MaxStates/2)
	case a.rev == "" || strings.HasPrefix(a.rev, "-"):
		return a, fmt.Errorf("--rev %q is not a revision", a.rev)
	}
	return a, nil
}

func runSurfaceIndex(args []string, stdout, stderr io.Writer) int {
	a, err := parseSurfaceIndexArgs(args, stderr)
	if errors.Is(err, flag.ErrHelp) {
		_, _ = io.WriteString(stdout, surfaceIndexHelp) // it holds git's %x09, which vet reads as a verb
		return exitOK
	}
	fail := func(format string, v ...any) int {
		fmt.Fprintf(stderr, "ailang-worldd surface-index: refused: "+format+"\n", v...)
		return exitUsage
	}
	if err != nil {
		return fail("%v", err)
	}
	if a.validate != "" {
		ix, err := readSurfaceIndex(a.validate)
		if err == nil {
			err = surface.Validate(ix)
		}
		if err != nil {
			return fail("%s: %v", a.validate, err)
		}
		fmt.Fprintf(stdout, "%s: valid %s (%s)\n", a.validate, surface.Schema, surfaceIndexCounts(ix))
		return exitOK
	}

	repoAbs, err := filepath.Abs(a.repo)
	if err == nil {
		var info os.FileInfo
		if info, err = os.Stat(repoAbs); err == nil && !info.IsDir() {
			err = errors.New("not a directory")
		}
	}
	if err != nil {
		return fail("--repo %s: %v", a.repo, err)
	}
	outAbs, err := filepath.Abs(a.out)
	if err != nil {
		return fail("--out %s: %v", a.out, err)
	}
	if outInsideRepo(outAbs, repoAbs) {
		return fail("--out %s is inside --repo %s: the index is a projection, never a file in the repo it reads", outAbs, repoAbs)
	}
	charter, err := os.ReadFile(a.charter)
	if err != nil {
		return fail("--charter: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), surfaceIndexBudget)
	defer cancel()
	in, err := gitSurfaceInput(ctx, repoAbs, a.rev, a.limit)
	if err != nil {
		return fail("%v", err)
	}
	in.Titles = surface.ChartRowTitles(string(charter))
	in.BuiltAt = time.Now()
	ix, err := surface.Build(in)
	if err == nil {
		err = surface.Validate(ix)
	}
	if err != nil {
		return fail("%v", err)
	}
	data, err := surface.Marshal(ix)
	if err != nil {
		return fail("%v", err)
	}
	if err := writeFileAtomic(outAbs, data); err != nil {
		return fail("%v", err)
	}
	fmt.Fprintf(stdout, "wrote %s: %s at %s (%s)\n", outAbs, a.rev, ix.Source.Head, surfaceIndexCounts(ix))
	return exitOK
}

func surfaceIndexCounts(ix surface.Index) string {
	shared := 0
	for _, s := range ix.States {
		if s.Kind == surface.KindShared {
			shared++
		}
	}
	return fmt.Sprintf("%d shared states, %d steps, %d branes", shared, len(ix.States)-shared, len(ix.Branes))
}

// outInsideRepo reports out lying in repo, as given or once both are resolved
// (out's file need not exist yet, so its directory is resolved).
func outInsideRepo(outAbs, repoAbs string) bool {
	if pathInsideOrEqual(outAbs, repoAbs) {
		return true
	}
	outReal := filepath.Join(realpathOrSelf(filepath.Dir(outAbs)), filepath.Base(outAbs))
	return pathInsideOrEqual(outReal, realpathOrSelf(repoAbs))
}

func readSurfaceIndex(path string) (surface.Index, error) {
	f, err := os.Open(path)
	if err != nil {
		return surface.Index{}, err
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, surfaceIndexMaxBytes+1))
	if err != nil {
		return surface.Index{}, err
	}
	if len(data) > surfaceIndexMaxBytes {
		return surface.Index{}, fmt.Errorf("over the bound of %d bytes", surfaceIndexMaxBytes)
	}
	return surface.Decode(data)
}

// gitSurfaceInput reads the merges, the head and each merge's merge-base.
func gitSurfaceInput(ctx context.Context, repo, rev string, limit int) (surface.BuildInput, error) {
	// rev never starts with "-" (parseSurfaceIndexArgs), so git cannot read
	// it as an option.
	head, err := surfaceGit(ctx, repo, "rev-parse", "--verify", rev+"^{commit}")
	if err != nil {
		return surface.BuildInput{}, err
	}
	log, err := surfaceGit(ctx, repo, "log", "--first-parent", "--merges",
		"--format=%H%x09%P%x09%cI%x09%s", "-n", strconv.Itoa(limit), rev)
	if err != nil {
		return surface.BuildInput{}, err
	}
	merges, err := surface.ParseMerges(log)
	if err != nil {
		return surface.BuildInput{}, err
	}
	bases := map[string]string{}
	for _, m := range merges {
		base, err := surfaceGit(ctx, repo, "merge-base", m.Parents[0], m.Parents[1])
		if err != nil {
			return surface.BuildInput{}, err
		}
		bases[m.Hash] = strings.TrimSpace(base)
	}
	return surface.BuildInput{Merges: merges, Bases: bases, Rev: rev, Head: strings.TrimSpace(head)}, nil
}

// surfaceGit runs `git -C repo args...` through host/worktree's one bounded,
// scrubbed git site, so no registry credential or GIT_* redirect reaches git.
func surfaceGit(ctx context.Context, repo string, args ...string) (string, error) {
	git, err := exec.LookPath("git")
	if err != nil {
		return "", err
	}
	return worktree.Output(ctx, git, append([]string{"-C", repo}, args...))
}

// writeFileAtomic writes data to a temporary file beside path and renames it
// over path, so a reader sees the old index or the new one, never part of one.
func writeFileAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	_, err = tmp.Write(data)
	if err == nil {
		err = tmp.Sync()
	}
	if cerr := tmp.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Chmod(name, 0o644)
	}
	if err == nil {
		err = os.Rename(name, path)
	}
	if err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
