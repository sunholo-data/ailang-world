package main

// Row 138 M5 (design_docs/planned/w-worldd-developer-cli.md §3.7):
// `session new` provisions an episode — its git worktree under the workspace
// root and its session credential — in one ATTENDED act, and `session list`
// shows the credentials a store holds.
//
// session new keeps every fence `session mint` has, in this order:
//  1. validate: the episode grammar, an EXISTING --db (never created), a sane
//     workspace root that does not contain the store, the grants;
//  2. the unchanged mint fence, env.openTerminal(), BEFORE any side effect
//     (MUT-FENCE-BYPASS: no TTY -> no worktree and no row);
//  3. a y/N on that terminal naming the worktree and every grant;
//  4. the CREDENTIALS store's writer lock (row 189: its own file, so a
//     serving daemon does not block it), failing fast if another session
//     verb holds it;
//  5. host/worktree.Add (detached, or -b <branch>), unless the episode's
//     worktree already exists;
//  6. authority.Mint through the existing core;
//  7. on a mint failure, remove the worktree THIS run created.
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
           [--allow-non-git]

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

The credential_id (the hash 'session revoke' takes) goes to stderr; the token
never does. If the mint fails, the worktree this run created is removed.
`

const sessionListHelp = `usage: ailang-worldd session list [--db <path>] [--credentials-db <path>] [--workspace-root <dir>] [--episode <ep>] [--json]

Lists the session credentials (at most 500): credential_id (the hash, never
the token), episode, live or expired, expiry and grants. It reads the
credentials store (--credentials-db, default credentials.db beside --db)
through a read-only handle, so it works while the daemon is running; for a
world store from before the split (schema 4) it reads the credentials still
inside it, also read-only. A --db that does not exist is refused, never created.

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

// sessionMint is authority.Mint; a test seam for the rollback arm (AC5.3).
var sessionMint = authority.Mint

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
	episode := positional[0]
	if !workspaceEpisodePattern.MatchString(episode) {
		return fail("episode %q does not match %s", episode, workspaceEpisodePattern)
	}
	if *ttl < 0 || *budget < 0 {
		return fail("--ttl and --budget must be non-negative")
	}
	if err := requireExistingStore(*db); err != nil {
		return fail("%v", err)
	}
	credPath, err := credentialsFor(*db, *credsDB)
	if err != nil {
		return fail("%v", err)
	}
	grants, err := sessionNewGrants(*preset, grantSpecs, *budget)
	if err != nil {
		return fail("%v", err)
	}
	root, err := filepath.Abs(*wsRoot)
	if err != nil {
		return fail("--workspace-root: %v", err)
	}
	if info, err := os.Lstat(root); err != nil || !info.IsDir() {
		return fail("workspace root %s is not an existing directory (a symlink is refused too); ailang-worldd setup creates it", root)
	}
	storeDir, _ := filepath.Abs(filepath.Dir(*db))
	if within(root, storeDir) {
		return fail("the workspace root %s contains the store directory %s; serve refuses this layout", root, storeDir)
	}
	// An agent that can write the credentials store can mint without a human.
	// credPath is symlink-free, so compare it with the root's canonical form too.
	canonicalRoot, _ := filepath.EvalSymlinks(root)
	for _, r := range []string{root, canonicalRoot} {
		if r != "" && within(r, filepath.Dir(credPath)) {
			return fail("the workspace root %s contains the credentials store %s; serve refuses this layout", root, credPath)
		}
	}
	dir := filepath.Join(root, episode)
	reuse, degraded := false, false
	if info, err := os.Lstat(dir); err == nil {
		switch {
		case !info.IsDir():
			return fail("%s exists and is not a directory (a symlink is refused too); remove it or pick another episode", dir)
		case isGitWorktree(dir):
		case !*allowNonGit:
			return fail("%s exists and is not a git worktree: an episode without git has no own branch, "+
				"no branch diff or provenance trailer, no reset to a base and no git-status worktree deltas. "+
				"fix: remove it or pick another episode (session new --repo makes the worktree), "+
				"or pass --allow-non-git to adopt it as a DEGRADED episode", dir)
		default:
			degraded = true
		}
		reuse = true
	} else if !errors.Is(err, os.ErrNotExist) {
		return fail("%s: %v", dir, err)
	}
	if *outPath != "" {
		if _, err := os.Lstat(*outPath); err == nil {
			return fail("--out %s already exists; the token is written once, never over a file", *outPath)
		}
		if info, err := os.Stat(filepath.Dir(*outPath)); err != nil || !info.IsDir() {
			return fail("--out %s: its directory does not exist", *outPath)
		}
	}
	git := ""
	if !reuse {
		if *repo == "" {
			return fail("--repo is required to create the worktree %s", dir)
		}
		if info, err := os.Stat(*repo); err != nil || !info.IsDir() {
			return fail("--repo %s is not a directory", *repo)
		}
		if git, err = sessionGit(); err != nil {
			return fail("git not found: %v", err)
		}
	}

	// 2. The mint fence, unchanged, before any side effect.
	term, err := env.openTerminal()
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd session new: refusing: no controlling terminal (%v); "+
			"provisioning a session credential requires one human act at a terminal. %s\n", err, sessionTTYFix)
		return exitUsage
	}
	defer func() { _ = term.Close() }()

	// 3. y/N naming the worktree and every grant.
	now := env.now()
	for i := range grants {
		grants[i].ExpiresAt = now + *ttl
	}
	action := "Create worktree " + dir + " from " + *repo
	switch {
	case degraded:
		action = "Use the existing NON-GIT directory " + dir + " as a DEGRADED episode (no branch diff, " +
			"provenance trailer or reset to a base; walk-basis worktree deltas)"
	case reuse:
		action = "Use the existing worktree " + dir
	}
	fmt.Fprintf(term, "%s and mint a session for episode %s with %d grant(s) [%s], expiry +%ds? [y/N] ",
		action, episode, len(grants), grantList(grants), *ttl)
	if !isYes(readTermLine(term)) {
		fmt.Fprintln(stderr, "ailang-worldd session new: aborted (not confirmed). fix: run it again and type y at the prompt")
		return exitUsage
	}

	// 4. The credentials writer lock, before the worktree exists.
	ctx, cancel := context.WithTimeout(context.Background(), sessionNewBudget)
	defer cancel()
	st, err := openCredentialsForSession(ctx, *db, credPath)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd session new: cannot open the credentials store %q: %v\n", credPath, err)
		return exitFatal
	}
	defer func() { _ = st.Close() }()

	// 5. The worktree.
	created := false
	if !reuse {
		if err := worktree.Add(ctx, git, *repo, dir, *branch); err != nil {
			fmt.Fprintf(stderr, "ailang-worldd session new: %v\n", err)
			return exitFatal
		}
		created = true
		fmt.Fprintf(stderr, "ailang-worldd session new: worktree %s created\n", dir)
	}

	// 6. The mint, through the existing core.
	var out io.Writer
	if *outPath == "" {
		out = stdout
	}
	mintCtx, cancelMint := context.WithTimeout(ctx, 3*time.Second)
	rawToken, credentialID, row, err := sessionMint(mintCtx, st, episode, grants, *ttl, now, out)
	cancelMint()
	if err == nil && *outPath != "" {
		if err = writeTokenFile(*outPath, rawToken); err != nil {
			// The row exists but nobody holds its token: revoke it too.
			if rerr := authority.Revoke(ctx, st, credentialID); rerr != nil {
				fmt.Fprintf(stderr, "ailang-worldd session new: could not revoke the unwritten credential %s: %v\n", credentialID, rerr)
			}
		}
	}
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd session new: %v\n", err)
		// 7. Roll back the worktree this run created.
		if created {
			if rerr := worktree.Remove(ctx, git, *repo, dir); rerr != nil {
				fmt.Fprintf(stderr, "ailang-worldd session new: could not remove the worktree it created: %v\n", rerr)
			} else {
				fmt.Fprintf(stderr, "ailang-worldd session new: removed the worktree %s it created\n", dir)
			}
		}
		return exitFatal
	}
	if *outPath != "" {
		fmt.Fprintf(stdout, "session for episode %s: %d grant(s), expires epoch %d, token written to %s\n",
			episode, len(grants), row.ExpiresAt, *outPath)
	} else {
		io.WriteString(stdout, "\n")
	}
	fmt.Fprintf(stderr, "ailang-worldd session new: credential_id=%s (episode %s, worktree %s, %d grant(s), expires %d)\n",
		credentialID, episode, dir, len(grants), row.ExpiresAt)
	if degraded {
		fmt.Fprintf(stderr, "ailang-worldd session new: episode %s is DEGRADED: %s is not a git worktree (--allow-non-git); "+
			"session list and /v1/health show it\n", episode, dir)
	}
	return exitOK
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
			Degraded: episodeDegraded(root, r.EpisodeID)})
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
		if l.Degraded {
			mark = "  DEGRADED (not a git worktree)"
		}
		fmt.Fprintf(stdout, "%s  %-12s %-7s expires %s  %s%s\n", l.CredentialID, l.Episode, state,
			time.Unix(l.ExpiresAt, 0).UTC().Format(time.RFC3339), grantList(l.Grants), mark)
	}
	fmt.Fprintf(stdout, "%d session credential(s)\n", len(listing))
	return exitOK
}
