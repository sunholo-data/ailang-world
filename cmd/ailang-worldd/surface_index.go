package main

// Row 188 M0a (design_docs/planned/w-world-os-surface-sprint-plan.md, "The
// surface index"): `ailang-worldd surface-index` builds the index /surface
// reads from a repo's git history, and checks one.
//
// It is the thin git side of host/surface: it runs git log and git
// merge-base, and everything else (parsing, branes, the index, its checks)
// is host/surface's pure code. The file is written atomically (a temporary
// file in the same directory, then a rename) and never inside the repo.
//
// Row 188 M0c: with --rig and --addr it also reads the rig's episode records
// and prs.json, GETs each attached episode's report from the loopback
// daemon, and runs git diff, ls-tree and log per step, so surface.AttachRig
// can add steps[], files{}, reads{} and decisions[].

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

const surfaceIndexHelp = `usage: ailang-worldd surface-index --repo <dir> --charter <file> --out <file> [--rev dev] [--limit 200]
                                 [--rig <dir> --addr <url> [--project ailang-world]]
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

With --rig, the index also holds the episode facts (otherwise it is unchanged):
  steps[]      one per episode of --project whose PR (<rig>/trial/prs.json,
               "<project>#N") is a step in the index: task, agent, model, base
               (round 1), rounds, graders (the last round), agent minutes, and
               the report's calls and entries. Records are
               <rig>/evidence/<project>[-rN]/episodes/<ep>.json, one per
               round, ordered by started_at; the rest are listed by id in
               source.rig.unattached
  files{}      git diff --name-status --no-renames -z <base> <step>
  reads{}      the report's reads, each path classified at the step's base:
                 git ls-tree -r -z --name-only <base>      (once per step)
                 git log --first-parent -1 --format=%H %P <base> -- <path>
               new (not at base), step (a merge of a step in the index),
               commit (a non-merge commit) or outside (an older merge)
  decisions[]  each round's grade, the report's refusals, and the merge
The report is GET <addr>/v1/episodes/<ep>/report (404: no report, its
fields absent). No path field of a record is copied into the index.

--validate reads an index and checks it: the schema, unique ids, every parent
resolving (but the oldest shared state's first parent), every step's brane
listed, every step's merge pointing back at it, times in order along the shared
chain, the allowed kinds, and at most 5000 states; with a rig, every episode
fact names a listed step, holds allowed values and local paths. It exits 0, or
1 with the reason.

Refused: a missing flag, --validate with any build flag, a --limit outside
1..2500, a --rev starting with -, an --out inside --repo, --rig without --addr
(or --addr or --project without --rig), an --addr that is not loopback
(127.0.0.1, localhost or ::1), an episode record over 1 MiB or more than 2000
of them, and a report over 4 MiB or slower than 30s.

  --repo <dir>       the git repository to read (required to build)
  --charter <file>   the charter whose queue rows title the branes (required to build)
  --out <file>       the index to write (required to build)
  --rev <rev>        the branch whose first-parent merges are read (default dev)
  --limit N          read at most N merges (default 200)
  --rig <dir>        the trial rig whose episode records and prs.json are read
  --project <name>   the rig project whose episodes are read (default ailang-world)
  --addr <url>       the loopback daemon serving the episode reports (with --rig)
  --validate <file>  check an index instead of building one
`

// surfaceIndexMaxBytes bounds an index --validate reads (the M1 serving bound).
const surfaceIndexMaxBytes = 4 << 20

// surfaceIndexBudget bounds the git calls and report GETs of one build as a
// whole.
const surfaceIndexBudget = 5 * time.Minute

// The rig's bounds (M0c): each episode record and prs.json, the number of
// records, each report body, and each report GET.
const (
	surfaceRecordMaxBytes = 1 << 20
	surfaceRecordMaxFiles = 2000
	surfaceReportMaxBytes = 4 << 20
	surfaceReportTimeout  = 30 * time.Second
)

// surfaceReportTransport carries the report GETs, with no proxy; a test
// replaces it with a handler, so no socket is needed.
var surfaceReportTransport http.RoundTripper = &http.Transport{Proxy: nil}

// surfaceIndexArgs is a parsed command line.
type surfaceIndexArgs struct {
	repo, charter, out, rev, validate string
	rig, project, addr                string
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
	fs.StringVar(&a.rig, "rig", "", "the trial rig whose episode records and prs.json are read")
	fs.StringVar(&a.project, "project", "ailang-world", "the rig project whose episodes are read")
	fs.StringVar(&a.addr, "addr", "", "the loopback daemon serving the episode reports")
	if err := fs.Parse(args); err != nil {
		return a, err
	}
	projectGiven := false
	fs.Visit(func(f *flag.Flag) { projectGiven = projectGiven || f.Name == "project" })
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
	case (a.rig == "") != (a.addr == ""):
		return a, errors.New("--rig and --addr go together: the rig's episodes need the daemon's reports")
	case projectGiven && a.rig == "":
		return a, errors.New("--project needs --rig")
	case a.rig != "" && (a.project == "" || strings.ContainsAny(a.project, `/\`) || strings.HasPrefix(a.project, ".")):
		return a, fmt.Errorf("--project %q is not a project name", a.project)
	}
	if a.addr != "" {
		base, err := surfaceLoopbackBase(a.addr)
		if err != nil {
			return a, err
		}
		a.addr = base
	}
	return a, nil
}

// surfaceLoopbackBase is --addr as a base URL, refused unless it is plain
// http(s) to 127.0.0.1, localhost or ::1 with no path, user or query.
func surfaceLoopbackBase(addr string) (string, error) {
	u, err := url.Parse(clientBaseURL(addr))
	if err != nil {
		return "", fmt.Errorf("--addr %q is not a URL: %v", addr, err)
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || u.RawQuery != "" || u.Fragment != "" ||
		(u.Path != "" && u.Path != "/") {
		return "", fmt.Errorf("--addr %q is not a daemon base URL (http://host:port)", addr)
	}
	switch u.Hostname() {
	case "127.0.0.1", "localhost", "::1":
	default:
		return "", fmt.Errorf("--addr %q is not loopback: the reports are read from 127.0.0.1, localhost or ::1 only", addr)
	}
	return u.Scheme + "://" + u.Host, nil
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
	var rig surface.RigInput
	if a.rig != "" {
		if rig, err = loadSurfaceRig(a.rig, a.project); err != nil {
			return fail("--rig: %v", err)
		}
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
	if err == nil && a.rig != "" {
		if rig.Reports, err = fetchSurfaceReports(ctx, a.addr, surfaceReportEpisodes(ix, rig)); err == nil {
			ix, err = surface.AttachRig(ix, rig, surfaceGitResolver{ctx: ctx, repo: repoAbs})
		}
	}
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
	counts := fmt.Sprintf("%d shared states, %d steps, %d branes", shared, len(ix.States)-shared, len(ix.Branes))
	if r := ix.Source.Rig; r != nil {
		counts += fmt.Sprintf("; rig: %d episodes, %d attached, %d decisions", r.Episodes, r.Attached, len(ix.Decisions))
	}
	return counts
}

// loadSurfaceRig reads project's episode records from every round directory
// of <rig>/evidence (surface.RoundDir) and its PRs from <rig>/trial/prs.json.
// Each file is read under surfaceRecordMaxBytes and at most
// surfaceRecordMaxFiles records are read; past either the build is refused.
func loadSurfaceRig(rig, project string) (surface.RigInput, error) {
	evidence := filepath.Join(rig, "evidence")
	dirs, err := os.ReadDir(evidence)
	if err != nil {
		return surface.RigInput{}, err
	}
	var records []surface.EpisodeRecord
	for _, d := range dirs {
		if !d.IsDir() || !surface.RoundDir(project, d.Name()) {
			continue
		}
		dir := filepath.Join(evidence, d.Name(), "episodes")
		files, err := os.ReadDir(dir)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return surface.RigInput{}, err
		}
		for _, f := range files {
			if !f.Type().IsRegular() || !strings.HasSuffix(f.Name(), ".json") {
				continue
			}
			if len(records) == surfaceRecordMaxFiles {
				return surface.RigInput{}, fmt.Errorf("more than %d episode records", surfaceRecordMaxFiles)
			}
			data, err := readFileCapped(filepath.Join(dir, f.Name()), surfaceRecordMaxBytes)
			if err != nil {
				return surface.RigInput{}, err
			}
			r, err := surface.ParseEpisodeRecord(data)
			if err != nil {
				return surface.RigInput{}, fmt.Errorf("%s/episodes/%s: %v", d.Name(), f.Name(), err)
			}
			records = append(records, r)
		}
	}
	data, err := readFileCapped(filepath.Join(rig, "trial", "prs.json"), surfaceRecordMaxBytes)
	if err != nil {
		return surface.RigInput{}, err
	}
	prs, err := surface.ParsePRs(data, project)
	if err != nil {
		return surface.RigInput{}, err
	}
	return surface.RigInput{Episodes: surface.GroupRounds(records), PRs: prs}, nil
}

// surfaceReportEpisodes lists the episodes whose PR is a step in the index:
// only their reports are fetched.
func surfaceReportEpisodes(ix surface.Index, rig surface.RigInput) []string {
	stepPRs := map[int]bool{}
	for _, s := range ix.States {
		if s.Kind == surface.KindStep && s.PR > 0 {
			stepPRs[s.PR] = true
		}
	}
	var out []string
	for _, ep := range rig.Episodes {
		if stepPRs[rig.PRs[ep.ID]] {
			out = append(out, ep.ID)
		}
	}
	sort.Strings(out)
	return out
}

// fetchSurfaceReports GETs each episode's report from the loopback daemon at
// base (already checked by surfaceLoopbackBase). A 404 is no report; any
// other non-200, a redirect, a body over surfaceReportMaxBytes or a request
// over surfaceReportTimeout refuses the build.
func fetchSurfaceReports(ctx context.Context, base string, episodes []string) (map[string]surface.Report, error) {
	hc := &http.Client{
		Timeout:       surfaceReportTimeout,
		Transport:     surfaceReportTransport,
		CheckRedirect: func(*http.Request, []*http.Request) error { return errors.New("redirect refused") },
	}
	out := map[string]surface.Report{}
	for _, ep := range episodes {
		path := "/v1/episodes/" + url.PathEscape(ep) + "/report"
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+path, nil)
		if err != nil {
			return nil, err
		}
		resp, err := hc.Do(req)
		if err != nil {
			return nil, fmt.Errorf("GET %s: %w", path, err)
		}
		body, err := io.ReadAll(io.LimitReader(resp.Body, surfaceReportMaxBytes+1))
		_ = resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("GET %s: %w", path, err)
		}
		if len(body) > surfaceReportMaxBytes {
			return nil, fmt.Errorf("GET %s: the report is over %d bytes", path, surfaceReportMaxBytes)
		}
		switch resp.StatusCode {
		case http.StatusOK:
			rep, err := surface.ParseReport(body)
			if err != nil {
				return nil, fmt.Errorf("GET %s: %v", path, err)
			}
			out[ep] = rep
		case http.StatusNotFound:
		default:
			return nil, fmt.Errorf("GET %s returned HTTP %d", path, resp.StatusCode)
		}
	}
	return out, nil
}

// surfaceGitResolver is surface.Git over the repo, through surfaceGit. Every
// revision it is given is a full commit hash from the index, never an option.
type surfaceGitResolver struct {
	ctx  context.Context
	repo string
}

func (g surfaceGitResolver) NameStatus(base, step string) (string, error) {
	return surfaceGit(g.ctx, g.repo, "diff", "--name-status", "--no-renames", "-z", base, step, "--")
}

func (g surfaceGitResolver) Tree(rev string) (string, error) {
	return surfaceGit(g.ctx, g.repo, "ls-tree", "-r", "-z", "--name-only", rev)
}

// LastChange reads the path literally (--literal-pathspecs), so a read of a
// file named with * or : matches only that file.
func (g surfaceGitResolver) LastChange(rev, path string) (string, error) {
	return surfaceGit(g.ctx, g.repo, "--literal-pathspecs", "log", "--first-parent", "-1", "--format=%H %P", rev, "--", path)
}

// readFileCapped reads path, refusing it over max bytes.
func readFileCapped(path string, bound int64) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, bound+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > bound {
		return nil, fmt.Errorf("%s is over the bound of %d bytes", filepath.Base(path), bound)
	}
	return data, nil
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
