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
//  4. store.Open, failing fast with "stop the daemon first";
//  5. host/worktree.Add (detached, or -b <branch>), unless the episode's
//     worktree already exists;
//  6. authority.Mint through the existing core;
//  7. on a mint failure, remove the worktree THIS run created.
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

const sessionNewHelp = `usage: ailang-worldd session new <episode> [--db <path>] [--workspace-root <dir>]
           [--repo <dir>] [--preset se-tools] [--grant EFFECT=SCOPE:BUDGET]...
           [--budget 50] [--ttl 3600] [--out <file>] [--branch <name>]

ATTENDED. Provisions one episode: its git worktree at <workspace-root>/<episode>
and a session credential bound to it, after a y/N typed at /dev/tty. With no
controlling terminal it refuses before touching anything (in an embedded
terminal, append </dev/tty).

  <episode>            ^[a-z0-9][a-z0-9-]{0,63}$ (the daemon's grammar)
  --db <path>          an EXISTING world store (default ~/.ailang/world/world.db);
                       never created. The daemon must be stopped.
  --workspace-root <dir>
                       default ~/.ailang/world-ws; must not contain the store
  --repo <dir>         the git repository the worktree is added from
                       (git worktree add --detach); not needed when
                       <workspace-root>/<episode> is already a worktree
  --branch <name>      add the worktree on a new branch instead of detached
  --preset se-tools    one grant per effect the se-tools transitions declare:
                       Workspace.Read, Workspace.Write, Ailang.Check, Ailang.Run,
                       Ailang.RunEnv, Ailang.RunNet, Ailang.Discover, Ailang.CLI,
                       Workspace.Exec, each scope worktree with --budget calls
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

const sessionListHelp = `usage: ailang-worldd session list [--db <path>] [--episode <ep>] [--json]

Lists the session credentials in a store (at most 500): credential_id (the
hash, never the token), episode, live or expired, expiry and grants. Reads
through a read-only handle, so it works while the daemon is running. A --db
that does not exist is refused, never created.
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
	wsRoot := fs.String("workspace-root", def.workspaceRoot, "workspace root")
	repo := fs.String("repo", "", "git repository the worktree is added from")
	branch := fs.String("branch", "", "add the worktree on this new branch")
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
	dir := filepath.Join(root, episode)
	reuse := false
	if info, err := os.Lstat(dir); err == nil {
		if !info.IsDir() || !isGitWorktree(dir) {
			return fail("%s exists and is not a git worktree; remove it or pick another episode", dir)
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
			"provisioning a session credential requires one human act at a terminal\n", err)
		return exitUsage
	}
	defer func() { _ = term.Close() }()

	// 3. y/N naming the worktree and every grant.
	now := env.now()
	for i := range grants {
		grants[i].ExpiresAt = now + *ttl
	}
	action := "Create worktree " + dir + " from " + *repo
	if reuse {
		action = "Use the existing worktree " + dir
	}
	fmt.Fprintf(term, "%s and mint a session for episode %s with %d grant(s) [%s], expiry +%ds? [y/N] ",
		action, episode, len(grants), grantList(grants), *ttl)
	if !isYes(readTermLine(term)) {
		fmt.Fprintln(stderr, "ailang-worldd session new: aborted (not confirmed)")
		return exitUsage
	}

	// 4. The writer lock, before the worktree exists.
	st, err := openStoreForSession(*db)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd session new: cannot open world store %q: %v\n", *db, err)
		return exitFatal
	}
	defer func() { _ = st.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), sessionNewBudget)
	defer cancel()

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
}

func runSessionList(args []string, stdout, stderr io.Writer, env sessionEnv) int {
	def := defaultWorldPaths()
	fs := flag.NewFlagSet("ailang-worldd session list", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	db := fs.String("db", def.db, "existing world store")
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
	ro, err := store.OpenReadOnly(*db)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd session list: cannot read %s: %v\n", *db, err)
		return exitFatal
	}
	defer func() { _ = ro.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := ro.ListSessions(ctx, *episode, store.MaxListSessions)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd session list: %v\n", err)
		return exitFatal
	}
	now := env.now()
	listing := make([]sessionListing, 0, len(rows))
	for _, r := range rows {
		var grants []broker.Capability
		_ = json.Unmarshal([]byte(r.GrantsJSON), &grants)
		sort.SliceStable(grants, func(i, j int) bool { return grants[i].Effect < grants[j].Effect })
		listing = append(listing, sessionListing{CredentialID: r.CredentialID, Episode: r.EpisodeID,
			Live: now < r.ExpiresAt, CreatedAt: r.CreatedAt, ExpiresAt: r.ExpiresAt, Grants: grants})
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
		fmt.Fprintf(stdout, "%s  %-12s %-7s expires %s  %s\n", l.CredentialID, l.Episode, state,
			time.Unix(l.ExpiresAt, 0).UTC().Format(time.RFC3339), grantList(l.Grants))
	}
	fmt.Fprintf(stdout, "%d session credential(s)\n", len(listing))
	return exitOK
}
