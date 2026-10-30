package main

// Row 138 M5 (design_docs/planned/w-worldd-developer-cli.md §3.7):
// `session new` provisions an episode — its git worktree under the workspace
// root and its session credential — in one ATTENDED act, and `session list`
// shows the credentials a store holds.
//
// session new keeps every fence `session mint` has, in this order:
//  1. validate: the episode grammar, an EXISTING --db (never created), a sane
//     workspace root that does not contain the store, the grants, and the
//     exec profile --profile names (prepareMint);
//  2. the unchanged mint fence, env.openTerminal(), BEFORE any side effect
//     (MUT-FENCE-BYPASS: no TTY -> no worktree and no row);
//  3. a y/N on that terminal naming the worktree, every grant and, with
//     --profile, the exec profile's project and digest;
//  4. the CREDENTIALS store's writer lock (row 189: its own file, so a
//     serving daemon does not block it), failing fast if another session
//     verb holds it;
//  5. host/worktree.Add (detached, or -b <branch>), unless the episode's
//     worktree already exists;
//  6. authority.MintBound through the existing core;
//  7. on a mint failure, remove the worktree THIS run created.
//
// Steps 2-7 are provisionEpisode, the ONE fence-and-mint path: `episode
// start` (episode.go, row 189 M3) runs it too, after its own validation, so
// the two verbs cannot drift apart and no other code creates a credential.
//
// Row 189 M2 (D-189-2 A): --profile <project> binds the credential to the
// exec profile <exec-profiles-dir>/<project>.json. The file is parsed (never
// probed) before the fence, its project and the sha256 of its exact bytes are
// named on the y/N line, and both are recorded on the credential; the daemon
// then runs the episode's workspace-exec under exactly those bytes, with no
// restart.
//
// Row 185 (w-episode-requires-git): an existing <workspace-root>/<episode>
// that is not a git worktree is refused, the refusal naming --allow-non-git.
// With that flag it is adopted as a DEGRADED episode: no own branch, no
// branch diff or provenance trailer, no reset to a base, and worktree deltas
// on the bounded-walk basis instead of git-status. Degraded is the fact
// worktree.IsLinked reads, so `session list` and /v1/health derive it from
// the worktree itself rather than from a column a schema bump would add.
// It binds only here, at mint; git episodes are unchanged.
//
// The se-tools preset is the union of packages/se-tools/transitions.json's
// declaredEffects (effect, scope); TestSessionNewPresetEqualsDeclaredEffects
// re-derives it from the manifest (MUT-PRESET-DRIFT).

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

const sessionNewHelp = `usage: ailang-worldd session new <episode> [--db <path>] [--credentials-db <path>] [--workspace-root <dir>]
           [--repo <dir>] [--preset se-tools] [--grant EFFECT=SCOPE:BUDGET]...
           [--budget 50] [--ttl 3600] [--out <file>] [--branch <name>]
           [--allow-non-git] [--profile <project>] [--exec-profiles-dir <dir>]

ATTENDED. Provisions one episode: its git worktree at <workspace-root>/<episode>
and a session credential bound to it, after a y/N typed at /dev/tty (any
terminal window, an IDE terminal pane included; stdin is never read). With no
controlling terminal (an agent, cron, CI) it refuses before touching anything.

  <episode>            ^[a-z0-9][a-z0-9-]{0,63}$ (the daemon's grammar)
  --db <path>          an EXISTING world store (default ~/.ailang/world/world.db);
                       never created. A serving daemon may keep running.
  --credentials-db <path>
                       the session credentials store the session is minted
                       into (default credentials.db beside --db, so
                       ~/.ailang/world/credentials.db); it has its own writer
                       lock, which no daemon holds
  --workspace-root <dir>
                       default ~/.ailang/world-ws; must not contain the store
  --repo <dir>         the git repository the worktree is added from
                       (git worktree add --detach); not needed when
                       <workspace-root>/<episode> is already a worktree
  --branch <name>      add the worktree on a new branch instead of detached
  --allow-non-git      adopt an existing <workspace-root>/<episode> that is
                       not a git worktree (refused without it). The episode
                       is DEGRADED: no branch diff, provenance trailer or
                       reset to a base, and its worktree deltas use the
                       bounded walk (basis walk / walk-capped), not
                       git-status. session list and /v1/health show it.
  --preset se-tools    one grant per effect the se-tools transitions declare:
                       Workspace.Read, Workspace.Write, Ailang.Check, Ailang.Run,
                       Ailang.RunEnv, Ailang.RunNet, Ailang.Discover, Ailang.CLI,
                       Workspace.Exec, Workspace.Delete, each scope worktree
                       with --budget calls
  --grant EFFECT=SCOPE:BUDGET
                       an explicit grant (repeatable); it overrides a preset
                       grant with the same effect and scope
  --budget N           budget of each preset grant (default 50)
  --ttl N              session lifetime in seconds (default 3600); every grant
                       expires with the session
  --out <file>         write the token once to <file> (mode 0600); default:
                       print it to stdout once
  --profile <project>  bind the session to the exec profile
                       <exec-profiles-dir>/<project>.json: it is parsed (not
                       probed) before the terminal opens, the y/N names the
                       project and the sha256 of its exact bytes, and both are
                       recorded on the credential. The daemon runs the
                       episode's workspace-exec under that profile from the
                       next request, with no restart, and refuses every exec
                       if the profile it loaded has a different digest
                       (re-mint). The credentials store must be schema 3:
                       a serve of this build upgrades it once at startup.
  --exec-profiles-dir <dir>
                       where --profile looks (default
                       ~/.ailang-world-exec/profiles, serve's default)

The credential_id (the hash 'session revoke' takes) goes to stderr; the token
never does. If the mint fails, the worktree this run created is removed.
`

const sessionListHelp = `usage: ailang-worldd session list [--db <path>] [--credentials-db <path>] [--workspace-root <dir>] [--episode <ep>] [--json]

Lists the session credentials (at most 500): credential_id (the hash, never
the token), episode, live or expired, expiry, grants and, for a credential
minted with --profile, the exec profile binding (project and digest; "project"
and "profileDigest" with --json). It reads the credentials store
(--credentials-db, default credentials.db beside --db) through a read-only
handle, so it works while the daemon is running; for a world store from before
the split (schema 4) it reads the credentials still inside it, also read-only.
A --db that does not exist is refused, never created.

An episode whose <workspace-root>/<episode> exists but is not a git worktree
(adopted with session new --allow-non-git) is marked DEGRADED ("degraded":
true with --json). --workspace-root defaults to ~/.ailang/world-ws.
`

// sePresetGrants is the se-tools preset: the union of the (effect, scope)
// pairs packages/se-tools/transitions.json declares.
var sePresetGrants = []broker.Capability{
	{Effect: broker.EffectWorkspaceRead, Scope: "worktree"},
	{Effect: broker.EffectWorkspaceWrite, Scope: "worktree"},
	{Effect: broker.EffectAilangCheck, Scope: "worktree"},
	{Effect: broker.EffectAilangRun, Scope: "worktree"},
	{Effect: broker.EffectAilangRunEnv, Scope: "worktree"},
	{Effect: broker.EffectAilangRunNet, Scope: "worktree"},
	{Effect: broker.EffectAilangDiscover, Scope: "worktree"},
	{Effect: broker.EffectAilangCLI, Scope: "worktree"},
	{Effect: broker.EffectWorkspaceExec, Scope: "worktree"},
	{Effect: broker.EffectWorkspaceDelete, Scope: "worktree"},
}

// sessionNewBudget bounds the whole command after the confirmation: the
// worktree add (≤ worktree.Timeout), the 3 s mint, and a rollback remove.
const sessionNewBudget = 75 * time.Second

// sessionGit resolves the git binary; a test seam.
var sessionGit = func() (string, error) { return exec.LookPath("git") }

// sessionMint is authority.MintBound; a test seam for the rollback arm (AC5.3).
var sessionMint = authority.MintBound

// defaultExecProfilesDir is ~/.ailang-world-exec/profiles: where `serve
// --exec-profiles-dir` and `--profile` look by default (row 189 M2).
func defaultExecProfilesDir() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		home = "."
	}
	return filepath.Join(home, ".ailang-world-exec", "profiles")
}

// mintOptions is what a provisioning verb asks for, before validation.
type mintOptions struct {
	verb        string // "session new" or "episode start": every message's prefix
	episode     string
	db          string
	credsDB     string
	wsRoot      string
	repo        string
	branch      string
	allowNonGit bool
	grants      []broker.Capability
	ttl         int64
	outPath     string
	// outDirCreate: outPath's directory may be absent and is then created
	// mode 0700 after the y/N (episode start's sessions dir); an existing one
	// must already be private.
	outDirCreate bool
	profile      string // --profile project, "" for an unbound session
	profilesDir  string
}

// mintRequest is a validated provisioning: everything provisionEpisode needs,
// computed with no side effect.
type mintRequest struct {
	mintOptions
	credPath string
	dir      string
	reuse    bool
	degraded bool
	git      string
	// prof is the parsed --profile (nil when unbound); profFile its path.
	prof     *broker.ExecProfile
	profFile string
	// afterToken runs once the token is written to outPath (episode start's
	// MCP config). A failure revokes the credential and rolls back like a
	// mint failure.
	afterToken     func(token string) error
	resume         *store.EpisodeResume
	beforeToken    func() error
	start          *store.EpisodeStart
	startPoint     string
	existingBranch bool
	presetPrompt   string
}

// prepareMint is step 1: it validates every input of a provisioning before
// the fence and touches nothing. A non-nil error is a usage refusal.
func prepareMint(o mintOptions) (*mintRequest, error) {
	if !workspaceEpisodePattern.MatchString(o.episode) {
		return nil, fmt.Errorf("episode %q does not match %s", o.episode, workspaceEpisodePattern)
	}
	if o.ttl < 0 {
		return nil, errors.New("--ttl and --budget must be non-negative")
	}
	if err := requireExistingStore(o.db); err != nil {
		return nil, err
	}
	credPath, err := credentialsFor(o.db, o.credsDB)
	if err != nil {
		return nil, err
	}
	if len(o.grants) == 0 {
		return nil, errors.New("no grants: pass --preset se-tools or at least one --grant EFFECT=SCOPE:BUDGET")
	}
	r := &mintRequest{mintOptions: o, credPath: credPath}
	root, err := filepath.Abs(o.wsRoot)
	if err != nil {
		return nil, fmt.Errorf("--workspace-root: %v", err)
	}
	if info, err := os.Lstat(root); err != nil || !info.IsDir() {
		return nil, fmt.Errorf("workspace root %s is not an existing directory (a symlink is refused too); ailang-worldd setup creates it", root)
	}
	storeDir, _ := filepath.Abs(filepath.Dir(o.db))
	if within(root, storeDir) {
		return nil, fmt.Errorf("the workspace root %s contains the store directory %s; serve refuses this layout", root, storeDir)
	}
	// An agent that can write the credentials store can mint without a human.
	// credPath is symlink-free, so compare it with the root's canonical form too.
	canonicalRoot, _ := filepath.EvalSymlinks(root)
	for _, rt := range []string{root, canonicalRoot} {
		if rt != "" && within(rt, filepath.Dir(credPath)) {
			return nil, fmt.Errorf("the workspace root %s contains the credentials store %s; serve refuses this layout", root, credPath)
		}
	}
	r.dir = filepath.Join(root, o.episode)
	if info, err := os.Lstat(r.dir); err == nil {
		switch {
		case !info.IsDir():
			return nil, fmt.Errorf("%s exists and is not a directory (a symlink is refused too); remove it or pick another episode", r.dir)
		case isGitWorktree(r.dir):
		case !o.allowNonGit:
			return nil, fmt.Errorf("%s exists and is not a git worktree: an episode without git has no own branch, "+
				"no branch diff or provenance trailer, no reset to a base and no git-status worktree deltas. "+
				"fix: remove it or pick another episode (session new --repo makes the worktree), "+
				"or pass --allow-non-git to adopt it as a DEGRADED episode", r.dir)
		default:
			r.degraded = true
		}
		r.reuse = true
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%s: %v", r.dir, err)
	}
	if o.outPath != "" {
		if _, err := os.Lstat(o.outPath); err == nil {
			return nil, fmt.Errorf("--out %s already exists; the token is written once, never over a file", o.outPath)
		}
		if err := checkOutDir(filepath.Dir(o.outPath), o.outDirCreate); err != nil {
			return nil, err
		}
	}
	if !r.reuse {
		if o.repo == "" {
			return nil, fmt.Errorf("--repo is required to create the worktree %s", r.dir)
		}
		if info, err := os.Stat(o.repo); err != nil || !info.IsDir() {
			return nil, fmt.Errorf("--repo %s is not a directory", o.repo)
		}
		if r.git, err = sessionGit(); err != nil {
			return nil, fmt.Errorf("git not found: %v", err)
		}
	}
	if o.profile != "" {
		if r.prof, r.profFile, err = loadMintProfile(o.profilesDir, o.profile); err != nil {
			return nil, err
		}
		if err := checkBindingFits(credPath, o.episode, r.prof); err != nil {
			return nil, err
		}
	}
	return r, nil
}

// checkOutDir requires the token file's directory to exist — or, for
// episode start, to be creatable — and, when it is the private sessions dir,
// to be readable by its owner only.
func checkOutDir(dir string, create bool) error {
	info, err := os.Stat(dir)
	switch {
	case err == nil && info.IsDir():
		if create && info.Mode().Perm()&0o077 != 0 {
			return fmt.Errorf("the sessions directory %s is mode %v; it holds session tokens and must be private: chmod 700 %s", dir, info.Mode().Perm(), dir)
		}
		return nil
	case err == nil:
		return fmt.Errorf("--out %s: %s is not a directory", dir, dir)
	case create && errors.Is(err, os.ErrNotExist):
		return nil
	default:
		return fmt.Errorf("--out %s: its directory does not exist", dir)
	}
}

// loadMintProfile reads <dir>/<project>.json once and parses it with the
// daemon's own loader. It is never probed here: serve's startup probe stays
// the real check, and the digest names the exact bytes read.
func loadMintProfile(dir, project string) (*broker.ExecProfile, string, error) {
	if !workspaceEpisodePattern.MatchString(project) {
		return nil, "", fmt.Errorf("--profile %q is not a project name (%s)", project, workspaceEpisodePattern)
	}
	file := filepath.Join(dir, project+broker.ExecProfileSuffix)
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, "", fmt.Errorf("--profile %s: %v (profiles live in --exec-profiles-dir %s; a *.proposed.json is activated by renaming it)", project, err, dir)
	}
	p, err := broker.ParseExecProfile(data)
	if err != nil {
		return nil, "", fmt.Errorf("--profile %s: %s: %v", project, file, err)
	}
	if p.Project != project {
		return nil, "", fmt.Errorf("--profile %s: %s declares project %q; a profiles-dir file is <project>.json", project, file, p.Project)
	}
	return p, file, nil
}

// checkBindingFits refuses, before the fence, a bound mint the credentials
// store cannot take or the daemon would refuse: a store still at schema 1
// (only serve upgrades it), and an episode that already holds a live
// credential bound to another profile or other bytes. It only reads.
func checkBindingFits(credPath, episode string, p *broker.ExecProfile) error {
	if _, err := os.Stat(credPath); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	ro, err := store.OpenCredentialsReadOnly(credPath)
	if err != nil {
		var schema *store.CredentialsSchemaError
		if errors.As(err, &schema) {
			return err
		}
		return nil // a store this can't read yet (pre-split): the mint moves it first
	}
	defer func() { _ = ro.Close() }()
	if ro.Version() < store.CredentialsSchemaVersion {
		return &store.UnboundCredentialsStoreError{Path: credPath}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := ro.ListSessions(ctx, episode, store.MaxListSessions)
	if err != nil {
		return err
	}
	now := time.Now().Unix()
	for _, r := range rows {
		if now > r.ExpiresAt || !r.Bound() || (r.Project == p.Project && r.ProfileDigest == p.Digest) {
			continue
		}
		return fmt.Errorf("episode %s already holds the live credential %s bound to exec profile %s %s; "+
			"one episode runs under one profile: revoke it first (ailang-worldd session revoke %s), then mint again",
			episode, r.CredentialID, r.Project, r.ProfileDigest, r.CredentialID)
	}
	return nil
}

// provisionEpisode is steps 2-7: the fence, the y/N, the credentials lock,
// the worktree, the mint and the rollback, in that order, for both session
// new and episode start. It returns the exit code.
func provisionEpisode(r *mintRequest, stdout, stderr io.Writer, env sessionEnv) int {
	prefix := "ailang-worldd " + r.verb + ": "
	// 2. The mint fence, unchanged, before any side effect.
	term, err := env.openTerminal()
	if err != nil {
		fmt.Fprintf(stderr, "%srefusing: no controlling terminal (%v); "+
			"provisioning a session credential requires one human act at a terminal. %s\n", prefix, err, sessionTTYFix)
		return exitUsage
	}
	defer func() { _ = term.Close() }()

	// 3. y/N naming the worktree, every grant and the exec profile binding.
	now := env.now()
	grants := append([]broker.Capability(nil), r.grants...)
	for i := range grants {
		grants[i].ExpiresAt = now + r.ttl
	}
	action := "Create worktree " + r.dir + " from " + r.repo
	switch {
	case r.degraded:
		action = "Use the existing NON-GIT directory " + r.dir + " as a DEGRADED episode (no branch diff, " +
			"provenance trailer or reset to a base; walk-basis worktree deltas)"
	case r.reuse:
		action = "Use the existing worktree " + r.dir
	}
	if !r.reuse && r.branch != "" && !r.existingBranch {
		action += " on new branch " + r.branch
	}
	binding := authority.ProfileBinding{}
	bound := ""
	if r.prof != nil {
		binding = authority.ProfileBinding{Project: r.prof.Project, Digest: r.prof.Digest}
		bound = fmt.Sprintf(", bound to exec profile %s %s (%s)", r.prof.Project, r.prof.Digest, r.profFile)
	}
	if r.presetPrompt != "" {
		fmt.Fprint(term, r.presetPrompt)
	}
	if r.existingBranch {
		action += " on existing branch " + r.branch
	}
	fmt.Fprintf(term, "%s and mint a session for episode %s with %d grant(s) [%s]%s, expiry +%ds? [y/N] ",
		action, r.episode, len(grants), grantList(grants), bound, r.ttl)
	if !isYes(readTermLine(term)) {
		fmt.Fprintf(stderr, "%saborted (not confirmed). fix: run it again and type y at the prompt\n", prefix)
		return exitUsage
	}

	// 4. The credentials writer lock, before the worktree exists.
	ctx, cancel := context.WithTimeout(context.Background(), sessionNewBudget)
	defer cancel()
	st, err := openCredentialsForSession(ctx, r.db, r.credPath)
	if err != nil {
		fmt.Fprintf(stderr, "%scannot open the credentials store %q: %v\n", prefix, r.credPath, err)
		return exitFatal
	}
	defer func() { _ = st.Close() }()

	// 5. The worktree.
	created := false
	if !r.reuse {
		var addErr error
		if r.existingBranch {
			addErr = worktree.AddExisting(ctx, r.git, r.repo, r.dir, r.branch)
		} else {
			addErr = worktree.Add(ctx, r.git, r.repo, r.dir, r.branch, r.startPoint)
		}
		if err := addErr; err != nil {
			fmt.Fprintf(stderr, "%s%v\n", prefix, err)
			return exitFatal
		}
		created = true
		fmt.Fprintf(stderr, "%sworktree %s created\n", prefix, r.dir)
	}

	// 6. The mint, through the existing core.
	var out io.Writer
	if r.outPath == "" {
		out = stdout
	}
	mintCtx, cancelMint := context.WithTimeout(ctx, 3*time.Second)
	var rawToken, credentialID string
	var row store.SessionRow
	if r.resume != nil {
		rawToken, credentialID, row, err = authority.MintResumed(mintCtx, st, r.episode, binding, grants, r.ttl, now, out, *r.resume)
	} else if r.start != nil {
		rawToken, credentialID, row, err = authority.MintStarted(mintCtx, st, r.episode, binding, grants, r.ttl, now, out, r.start)
	} else {
		rawToken, credentialID, row, err = sessionMint(mintCtx, st, r.episode, binding, grants, r.ttl, now, out)
	}
	cancelMint()
	minted := err == nil
	if err == nil && r.beforeToken != nil {
		err = r.beforeToken()
	}
	if err == nil && r.outPath != "" {
		if r.outDirCreate {
			err = os.MkdirAll(filepath.Dir(r.outPath), 0o700)
		}
		if err == nil {
			err = writeTokenFile(r.outPath, rawToken)
		}
		if err == nil && r.afterToken != nil {
			if err = r.afterToken(rawToken); err != nil {
				_ = os.Remove(r.outPath)
			}
		}
	}
	if err != nil && minted {
		// The row exists but nobody holds its token: revoke it too.
		if rerr := authority.Revoke(ctx, st, credentialID); rerr != nil {
			fmt.Fprintf(stderr, "%scould not revoke the unwritten credential %s: %v\n", prefix, credentialID, rerr)
		}
	}
	if err != nil {
		fmt.Fprintf(stderr, "%s%v\n", prefix, err)
		// 7. Roll back the worktree this run created.
		if created {
			if rerr := worktree.Remove(ctx, r.git, r.repo, r.dir); rerr != nil {
				fmt.Fprintf(stderr, "%scould not remove the worktree it created: %v\n", prefix, rerr)
			} else {
				fmt.Fprintf(stderr, "%sremoved the worktree %s it created\n", prefix, r.dir)
			}
		}
		return exitFatal
	}
	if r.outPath != "" {
		fmt.Fprintf(stdout, "session for episode %s: %d grant(s), expires epoch %d, token written to %s\n",
			r.episode, len(grants), row.ExpiresAt, r.outPath)
	} else {
		io.WriteString(stdout, "\n")
	}
	fmt.Fprintf(stderr, "%scredential_id=%s (episode %s, worktree %s, %d grant(s), expires %d)\n",
		prefix, credentialID, r.episode, r.dir, len(grants), row.ExpiresAt)
	if r.prof != nil {
		fmt.Fprintf(stderr, "%sbound to exec profile %s %s\n", prefix, r.prof.Project, r.prof.Digest)
	}
	if r.degraded {
		fmt.Fprintf(stderr, "%sepisode %s is DEGRADED: %s is not a git worktree (--allow-non-git); "+
			"session list and /v1/health show it\n", prefix, r.episode, r.dir)
	}
	return exitOK
}

func runSessionNew(args []string, stdout, stderr io.Writer, env sessionEnv) int {
	def := defaultWorldPaths()
	fs := flag.NewFlagSet("ailang-worldd session new", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	db := fs.String("db", def.db, "existing world store")
	credsDB := fs.String("credentials-db", "", credentialsDBUsage)
	wsRoot := fs.String("workspace-root", def.workspaceRoot, "workspace root")
	repo := fs.String("repo", "", "git repository the worktree is added from")
	branch := fs.String("branch", "", "add the worktree on this new branch")
	allowNonGit := fs.Bool("allow-non-git", false, "adopt an existing non-git episode directory, marked degraded")
	preset := fs.String("preset", "", "grant preset (se-tools)")
	var grantSpecs multiFlag
	fs.Var(&grantSpecs, "grant", "EFFECT=SCOPE:BUDGET, repeatable")
	budget := fs.Int64("budget", 50, "budget of each preset grant")
	ttl := fs.Int64("ttl", 3600, "session lifetime in seconds")
	outPath := fs.String("out", "", "write the token once to this file (0600)")
	profile := fs.String("profile", "", "bind the session to this exec profile's project")
	profilesDir := fs.String("exec-profiles-dir", defaultExecProfilesDir(), "where --profile looks")
	positional, err := parseInterleaved(fs, args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(stdout, sessionNewHelp)
			return exitOK
		}
		fmt.Fprint(stderr, sessionNewHelp)
		return exitUsage
	}
	fail := func(format string, a ...any) int {
		fmt.Fprintf(stderr, "ailang-worldd session new: "+format+"\n", a...)
		return exitUsage
	}

	// 1. Validate everything before the fence.
	if len(positional) != 1 {
		fmt.Fprint(stderr, sessionNewHelp)
		return exitUsage
	}
	if !workspaceEpisodePattern.MatchString(positional[0]) {
		return fail("episode %q does not match %s", positional[0], workspaceEpisodePattern)
	}
	if *ttl < 0 || *budget < 0 {
		return fail("--ttl and --budget must be non-negative")
	}
	grants, err := sessionNewGrants(*preset, grantSpecs, *budget)
	if err != nil {
		return fail("%v", err)
	}
	req, err := prepareMint(mintOptions{verb: "session new", episode: positional[0], db: *db, credsDB: *credsDB,
		wsRoot: *wsRoot, repo: *repo, branch: *branch, allowNonGit: *allowNonGit, grants: grants, ttl: *ttl,
		outPath: *outPath, profile: *profile, profilesDir: *profilesDir})
	if err != nil {
		return fail("%v", err)
	}
	if req.prof != nil && !hasGrant(grants, broker.EffectWorkspaceExec) {
		fmt.Fprintf(stderr, "ailang-worldd session new: note: --profile %s binds workspace-exec, but no Workspace.Exec grant is minted\n", req.prof.Project)
	}
	return provisionEpisode(req, stdout, stderr, env)
}

// hasGrant reports whether grants hold effect at any scope.
func hasGrant(grants []broker.Capability, effect string) bool {
	for _, g := range grants {
		if g.Effect == effect {
			return true
		}
	}
	return false
}

// writeTokenFile writes the token once at mode 0600, refusing to overwrite.
func writeTokenFile(path, token string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return fmt.Errorf("write the token to %s: %w", path, errors.Unwrap(err))
	}
	if _, err := io.WriteString(f, token); err != nil {
		_ = f.Close()
		return fmt.Errorf("write the token to %s: %w", path, err)
	}
	return f.Close()
}

// sessionNewGrants builds the preset grants (each with budget), then applies
// explicit --grant specs: one with the same effect and scope replaces the
// preset's, others are appended.
func sessionNewGrants(preset string, specs []string, budget int64) ([]broker.Capability, error) {
	var grants []broker.Capability
	switch preset {
	case "":
	case "se-tools":
		for _, g := range sePresetGrants {
			g.Budget = budget
			grants = append(grants, g)
		}
	default:
		return nil, fmt.Errorf("unknown --preset %q (known: se-tools)", preset)
	}
	extra, err := parseGrants(specs)
	if err != nil {
		return nil, err
	}
	for _, e := range extra {
		replaced := false
		for i := range grants {
			if grants[i].Effect == e.Effect && grants[i].Scope == e.Scope {
				grants[i] = e
				replaced = true
			}
		}
		if !replaced {
			grants = append(grants, e)
		}
	}
	if len(grants) == 0 {
		return nil, errors.New("no grants: pass --preset se-tools or at least one --grant EFFECT=SCOPE:BUDGET")
	}
	return grants, nil
}

func grantList(grants []broker.Capability) string {
	parts := make([]string, 0, len(grants))
	for _, g := range grants {
		parts = append(parts, fmt.Sprintf("%s=%s:%d", g.Effect, g.Scope, g.Budget))
	}
	return strings.Join(parts, " ")
}

// sessionListing is one `session list --json` row.
type sessionListing struct {
	CredentialID string              `json:"credentialId"`
	Episode      string              `json:"episode"`
	Live         bool                `json:"live"`
	CreatedAt    int64               `json:"createdAt"`
	ExpiresAt    int64               `json:"expiresAt"`
	Grants       []broker.Capability `json:"grants"`
	// Degraded: the episode's directory exists and is not a git worktree
	// (row 185, session new --allow-non-git).
	Degraded bool `json:"degraded"`
	// Project and ProfileDigest are the exec profile binding (row 189 M2);
	// absent for an unbound credential.
	Project       string `json:"project,omitempty"`
	ProfileDigest string `json:"profileDigest,omitempty"`
}

// episodeDegraded reports whether <root>/<episode> exists as a real
// directory that is not a git worktree. An absent directory is not degraded:
// there is nothing to judge.
func episodeDegraded(root, episode string) bool {
	if root == "" || !workspaceEpisodePattern.MatchString(episode) {
		return false
	}
	dir := filepath.Join(root, episode)
	info, err := os.Lstat(dir)
	return err == nil && info.IsDir() && !isGitWorktree(dir)
}

// listSessionRows reads the credentials, read-only and lock-free on every
// path: from the credentials store, or — for a world store from before the
// split (schema 4, e.g. beside a daemon of an older build) — from the table
// still inside it. A current world store whose credentials store does not
// exist yet simply has none.
func listSessionRows(ctx context.Context, db, credPath, episode string, stderr io.Writer) ([]store.SessionRow, error) {
	version, err := store.WorldSchemaVersion(db)
	if err != nil {
		return nil, err
	}
	if version == store.PreSplitSchemaVersion {
		fmt.Fprintf(stderr, "ailang-worldd session list: %s is a pre-split store (schema 4): listing the credentials still inside it; "+
			"the first serve of this build moves them to %s\n", db, credPath)
		return store.LegacySessions(ctx, db, episode, store.MaxListSessions)
	}
	if _, err := os.Stat(credPath); errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	ro, err := store.OpenCredentialsReadOnly(credPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", credPath, err)
	}
	defer func() { _ = ro.Close() }()
	return ro.ListSessions(ctx, episode, store.MaxListSessions)
}

func runSessionList(args []string, stdout, stderr io.Writer, env sessionEnv) int {
	def := defaultWorldPaths()
	fs := flag.NewFlagSet("ailang-worldd session list", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	db := fs.String("db", def.db, "existing world store")
	credsDB := fs.String("credentials-db", "", credentialsDBUsage)
	wsRoot := fs.String("workspace-root", def.workspaceRoot, "workspace root (marks degraded episodes)")
	episode := fs.String("episode", "", "only this episode")
	asJSON := fs.Bool("json", false, "print JSON")
	positional, err := parseInterleaved(fs, args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(stdout, sessionListHelp)
			return exitOK
		}
		fmt.Fprint(stderr, sessionListHelp)
		return exitUsage
	}
	if len(positional) > 0 {
		fmt.Fprintf(stderr, "ailang-worldd session list: unexpected argument %q\n", positional[0])
		return exitUsage
	}
	if err := requireExistingStore(*db); err != nil {
		fmt.Fprintf(stderr, "ailang-worldd session list: %v\n", err)
		return exitUsage
	}
	credPath, err := credentialsFor(*db, *credsDB)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd session list: %v\n", err)
		return exitUsage
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := listSessionRows(ctx, *db, credPath, *episode, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd session list: %v\n", err)
		return exitFatal
	}
	root, err := filepath.Abs(*wsRoot)
	if err != nil {
		root = ""
	}
	now := env.now()
	listing := make([]sessionListing, 0, len(rows))
	for _, r := range rows {
		var grants []broker.Capability
		_ = json.Unmarshal([]byte(r.GrantsJSON), &grants)
		sort.SliceStable(grants, func(i, j int) bool { return grants[i].Effect < grants[j].Effect })
		listing = append(listing, sessionListing{CredentialID: r.CredentialID, Episode: r.EpisodeID,
			Live: now < r.ExpiresAt, CreatedAt: r.CreatedAt, ExpiresAt: r.ExpiresAt, Grants: grants,
			Degraded: episodeDegraded(root, r.EpisodeID), Project: r.Project, ProfileDigest: r.ProfileDigest})
	}
	if *asJSON {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(listing)
		return exitOK
	}
	for _, l := range listing {
		state := "live"
		if !l.Live {
			state = "expired"
		}
		mark := ""
		if l.Project != "" {
			mark += "  profile " + l.Project + " " + l.ProfileDigest
		}
		if l.Degraded {
			mark += "  DEGRADED (not a git worktree)"
		}
		fmt.Fprintf(stdout, "%s  %-12s %-7s expires %s  %s%s\n", l.CredentialID, l.Episode, state,
			time.Unix(l.ExpiresAt, 0).UTC().Format(time.RFC3339), grantList(l.Grants), mark)
	}
	fmt.Fprintf(stdout, "%d session credential(s)\n", len(listing))
	return exitOK
}
