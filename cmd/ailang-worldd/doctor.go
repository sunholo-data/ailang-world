package main

// Row 138 M4 (design_docs/planned/w-worldd-developer-cli.md §3.6): `doctor`
// names, in one read-only pass, every snag the attended ritual used to
// rediscover by hand (P6): the API-key guard, the two TTY fences, missing or
// drifted pins, a busy port, a held writer lock, a missing head or transition
// registry, the examples corpus and the workspace layout.
//
// READ-ONLY BY CONSTRUCTION. doctor never creates a store, a lock file or a
// directory: the writer lock is probed with store.WriterLockHeld (a shared,
// non-blocking flock on the lock file opened read-only; absent lock file =
// free) and the store is read through store.OpenReadOnly, which takes no lock
// and applies no schema. MUT-DOCTOR-OPENS-STORE (store.Open in place of
// either) is killed by TestDoctorLockProbeCreatesNothing. Pins are hashed,
// never executed.

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/pinfetch"
	"github.com/sunholo-data/ailang-world/host/store"
)

const doctorHelp = `usage: ailang-worldd [--addr <url>] doctor [--db <path>] [--workspace-root <dir>]
           [--interpreter-dir <dir>] [--tools-dir <dir>] [--examples-dir <dir>] [--online]

Checks this machine for everything World needs, read-only, and prints one
line per check: ✓ fine, ! worth knowing (with a fix), ✗ broken (with a fix).

  api key      AILANG_REGISTRY_API_KEY is unset (doctor could not run otherwise:
               every verb refuses while it is set)
  tty          /dev/tty opens, so the attended steps (world-publish, session
               new/mint) can read the confirmation you type
  pins         both pinned binaries present and hashing to the compiled-in
               digests (never executed); stale tarballs and old binaries noted
  daemon       what answers at --addr: a worldd (its store and interpreter),
               a foreign listener, or nothing
  store        exists (never created), writer lock held or free, world head
               and transition registry present
  examples     the examples corpus examples-search reads
  workspace    the workspace root; each child a git worktree with a valid
               episode name
  --online     the release .sha256 of both pins is reachable and matches

Defaults: --db ~/.ailang/world/world.db, --workspace-root ~/.ailang/world-ws,
the pin directories of 'setup', --examples-dir ~/.ailang/examples.
Exit: 0 no ✗; 1 at least one ✗ (or a usage error).
`

// doctorStoreBudget bounds doctor's read-only store reads.
const doctorStoreBudget = 3 * time.Second

// doctorProbeTimeout bounds the daemon health probe.
const doctorProbeTimeout = 2 * time.Second

// ttyObservation is what the attended fences look at: whether /dev/tty opens.
// Since 2026-10-06 that is the only terminal fact either fence decides on —
// world-publish, like session new/mint, reads the typed line from /dev/tty
// itself when stdin is something else — so stdin's identity is not observed.
type ttyObservation struct {
	cttyErr error
}

// observeTTY opens /dev/tty READ-ONLY and closes it at once. It is the only
// place doctor touches the terminal; tests replace it.
var observeTTY = func() ttyObservation {
	var o ttyObservation
	tty, err := os.OpenFile("/dev/tty", os.O_RDONLY, 0)
	if err != nil {
		o.cttyErr = err
		return o
	}
	_ = tty.Close()
	return o
}

type doctorMark string

const (
	markOK   doctorMark = "✓"
	markWarn doctorMark = "!"
	markFail doctorMark = "✗"
)

type doctorReport struct {
	w     io.Writer
	fails int
}

func (r *doctorReport) add(mark doctorMark, name, detail, fix string) {
	fmt.Fprintf(r.w, "%s %-10s %s\n", mark, name, detail)
	if fix != "" {
		fmt.Fprintf(r.w, "  %-10s fix: %s\n", "", fix)
	}
	if mark == markFail {
		r.fails++
	}
}

func runDoctor(addr string, args []string, stdout, stderr io.Writer) int {
	def := defaultWorldPaths()
	fs := flag.NewFlagSet("ailang-worldd doctor", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	db := fs.String("db", def.db, "world store path")
	wsRoot := fs.String("workspace-root", def.workspaceRoot, "workspace root")
	interpDir := fs.String("interpreter-dir", def.interpDir, "interpreter pin directory")
	toolsDir := fs.String("tools-dir", def.toolsDir, "tool pin parent directory")
	examplesDir := fs.String("examples-dir", "", "examples corpus (default ~/.ailang/examples)")
	online := fs.Bool("online", false, "check the release .sha256 of both pins")
	if code, done := parseVerbFlags(fs, args, doctorHelp, stdout, stderr); done {
		return code
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(stderr, "ailang-worldd doctor: unexpected argument %q\n", fs.Arg(0))
		return exitUsage
	}
	r := &doctorReport{w: stdout}

	// API key: by construction. run() refused before reaching here if the
	// variable were set, so the guard stays exemption-free (§3.6, V8).
	r.add(markOK, "api key", "AILANG_REGISTRY_API_KEY is unset", "")

	doctorTTY(r, observeTTY())
	roles := doctorPins(r, *interpDir, *toolsDir)
	healthDB := doctorDaemon(r, addr, *db)
	doctorStore(r, *db, healthDB)
	doctorExamples(r, resolveExamplesDefault(*examplesDir, os.UserHomeDir))
	doctorWorkspace(r, *wsRoot, *db)
	if *online {
		doctorOnline(r, roles)
	}

	if r.fails > 0 {
		fmt.Fprintf(stdout, "%d check(s) failed\n", r.fails)
		return exitUsage
	}
	fmt.Fprintln(stdout, "no check failed")
	return exitOK
}

func doctorTTY(r *doctorReport, o ttyObservation) {
	if o.cttyErr != nil {
		r.add(markWarn, "tty", fmt.Sprintf("no controlling terminal (%v): world-publish and session new/mint will refuse here", o.cttyErr),
			"run the attended steps yourself from a terminal window (an IDE terminal pane works); an agent cannot run them")
		return
	}
	r.add(markOK, "tty", "/dev/tty opens: world-publish and session new/mint read your confirmation from it (no stdin redirect needed)", "")
}

// doctorPins hashes both pins against the table; it never executes them.
func doctorPins(r *doctorReport, interpDir, toolsDir string) []pinRole {
	roles := pinRoles(interpDir, toolsDir)
	for _, role := range roles {
		pin, err := lookupPin(role.release, setupPlatform)
		if err != nil {
			r.add(markFail, "pins", fmt.Sprintf("%s: %v", role.name, err), "")
			continue
		}
		label := fmt.Sprintf("%s AILANG %s", role.name, pin.Release)
		info, err := os.Lstat(role.dest)
		switch {
		case errors.Is(err, os.ErrNotExist):
			r.add(markFail, "pins", label+": missing at "+role.dest, "ailang-worldd setup")
			continue
		case err != nil:
			r.add(markFail, "pins", fmt.Sprintf("%s: %v", label, err), "")
			continue
		case !info.Mode().IsRegular():
			r.add(markFail, "pins", label+": "+role.dest+" is not a regular file", "ailang-worldd setup --replace")
			continue
		}
		sum, err := pinfetch.HashFile(role.dest)
		if err != nil {
			r.add(markFail, "pins", fmt.Sprintf("%s: hash %s: %v", label, role.dest, err), "")
			continue
		}
		if sum != pin.BinarySHA256 {
			r.add(markFail, "pins", fmt.Sprintf("%s: %s has sha256 %s, the pin is %s", label, role.dest, sum, pin.BinarySHA256),
				"ailang-worldd setup --replace (keeps the old file as ailang.prev-<sha8>)")
			continue
		}
		r.add(markOK, "pins", fmt.Sprintf("%s at %s (sha256 %s…)", label, role.dest, sum[:12]), "")
		if stale := staleFiles(filepath.Dir(role.dest)); len(stale) > 0 {
			r.add(markWarn, "pins", fmt.Sprintf("%s: stale files beside the pin: %s", role.name, strings.Join(stale, ", ")),
				"delete them; setup keeps no tarball")
		}
	}
	return roles
}

// staleFiles names what is neither the pin nor its pin.json.
func staleFiles(dir string) []string {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range ents {
		if e.Name() == "ailang" || e.Name() == pinfetch.RecordName || e.IsDir() {
			continue
		}
		out = append(out, e.Name())
	}
	sort.Strings(out)
	return out
}

// doctorDaemon reports what answers at addr; it returns the worldd's db_path
// ("" when no worldd answered).
func doctorDaemon(r *doctorReport, addr, db string) string {
	ctx, cancel := budgetContext(doctorProbeTimeout + time.Second)
	defer cancel()
	c := newClient(addr)
	c.timeout = doctorProbeTimeout
	status, body, err := c.do(ctx, http.MethodGet, "/v1/health", nil)
	if err != nil {
		r.add(markWarn, "daemon", "nothing answers at "+addr, "ailang-worldd serve --db … (when you need it running)")
		return ""
	}
	var h struct {
		Status             string `json:"status"`
		DaemonVersion      string `json:"daemon_version"`
		DBPath             string `json:"db_path"`
		InterpreterVersion string `json:"interpreter_version"`
	}
	if status != http.StatusOK || json.Unmarshal(body, &h) != nil || h.DaemonVersion == "" {
		r.add(markFail, "daemon", fmt.Sprintf("something that is not ailang-worldd answers at %s (HTTP %d)", addr, status),
			"stop it, or serve on another --bind and pass --addr")
		return ""
	}
	detail := fmt.Sprintf("ailang-worldd %s at %s, store %s", h.DaemonVersion, addr, h.DBPath)
	if h.InterpreterVersion != "" {
		detail += ", interpreter " + firstLine(h.InterpreterVersion)
	}
	if !sameStorePath(h.DBPath, db) {
		r.add(markWarn, "daemon", detail, "it serves a different store than --db "+db)
		return h.DBPath
	}
	r.add(markOK, "daemon", detail, "")
	return h.DBPath
}

func firstLine(s string) string {
	s, _, _ = strings.Cut(strings.TrimSpace(s), "\n")
	return s
}

func sameStorePath(a, b string) bool {
	ca, errA := filepath.Abs(a)
	cb, errB := filepath.Abs(b)
	if errA != nil || errB != nil {
		return a == b
	}
	if ra, err := filepath.EvalSymlinks(ca); err == nil {
		ca = ra
	}
	if rb, err := filepath.EvalSymlinks(cb); err == nil {
		cb = rb
	}
	return ca == cb
}

// doctorStore inspects the store without creating or locking anything.
func doctorStore(r *doctorReport, db, healthDB string) {
	// The lock probe runs first and on a missing store too: it creates
	// nothing (an absent lock file is "free"), which is exactly what
	// TestDoctorLockProbeCreatesNothing pins.
	held, lockErr := store.WriterLockHeld(db)
	info, err := os.Stat(db)
	if errors.Is(err, os.ErrNotExist) {
		r.add(markWarn, "store", db+" does not exist yet (doctor never creates it)",
			"world-publish transitions --store "+db+" … creates it (the attended publish)")
		return
	}
	if err != nil || !info.Mode().IsRegular() {
		r.add(markFail, "store", fmt.Sprintf("%s is not a readable regular file (%v)", db, err), "")
		return
	}
	switch {
	case lockErr != nil:
		r.add(markFail, "store", fmt.Sprintf("writer lock probe: %v", lockErr), "")
	case held && healthDB != "" && sameStorePath(healthDB, db):
		r.add(markWarn, "store", "writer lock held by the running daemon",
			"stop it before world-publish, session new/mint and session revoke (they need the writer lock)")
	case held:
		r.add(markWarn, "store", "writer lock held by another process",
			"stop that writer before world-publish, session new/mint and session revoke")
	default:
		r.add(markOK, "store", "writer lock free ("+db+")", "")
	}

	ro, err := store.OpenReadOnly(db)
	if err != nil {
		r.add(markFail, "store", fmt.Sprintf("cannot read %s: %v", db, err), "")
		return
	}
	defer func() { _ = ro.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), doctorStoreBudget)
	defer cancel()
	if head, ok, err := ro.SelectedHead(ctx); err != nil {
		r.add(markFail, "store", fmt.Sprintf("read the world head: %v", err), "")
	} else if !ok {
		r.add(markWarn, "store", "no world head yet: a tool call has nothing to commit on",
			"commit a genesis world (docs/QUICKSTART.md §9) after serve")
	} else {
		r.add(markOK, "store", "world head "+head.String(), "")
	}
	if ref, ok, err := ro.GetRegistryHead(ctx, store.TransitionRegistryV1); err != nil {
		r.add(markFail, "store", fmt.Sprintf("read the transition registry: %v", err), "")
	} else if !ok {
		r.add(markWarn, "store", "no transition registry: no tool is published in this store",
			"world-publish transitions --store "+db+" --manifest packages/se-tools/transitions.json …")
	} else {
		r.add(markOK, "store", "transition registry "+ref.String(), "")
	}
}

func doctorExamples(r *doctorReport, dir string) {
	if dir == "" {
		r.add(markWarn, "examples", "no examples corpus (~/.ailang/examples is absent): examples-search will refuse",
			"pass serve --examples-dir <dir>")
		return
	}
	n := 0
	_ = filepath.WalkDir(dir, func(_ string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(d.Name(), ".ail") {
			n++
		}
		return nil
	})
	if n == 0 {
		r.add(markWarn, "examples", dir+" holds no .ail file: examples-search will find nothing", "")
		return
	}
	r.add(markOK, "examples", fmt.Sprintf("%d .ail files in %s", n, dir), "")
}

// workspaceEpisodePattern is the daemon's episode grammar (host/daemon
// workspace.go episodeIDPattern; TestDoctorEpisodeGrammarMatchesTheDaemon).
var workspaceEpisodePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

func doctorWorkspace(r *doctorReport, root, db string) {
	info, err := os.Lstat(root)
	if errors.Is(err, os.ErrNotExist) {
		r.add(markWarn, "workspace", root+" does not exist", "ailang-worldd setup (creates it, mode 0700)")
		return
	}
	if err != nil || !info.IsDir() {
		r.add(markFail, "workspace", root+" is not a directory (a symlink is refused by serve too)", "")
		return
	}
	absRoot, _ := filepath.Abs(root)
	absStore, _ := filepath.Abs(filepath.Dir(db))
	if within(absRoot, absStore) {
		r.add(markFail, "workspace", fmt.Sprintf("%s contains the store directory %s: serve refuses this", root, absStore),
			"move the store or the workspace root apart")
		return
	}
	ents, err := os.ReadDir(root)
	if err != nil {
		r.add(markFail, "workspace", fmt.Sprintf("read %s: %v", root, err), "")
		return
	}
	var good []string
	var bad []string
	for _, e := range ents {
		name := e.Name()
		switch {
		case !workspaceEpisodePattern.MatchString(name):
			bad = append(bad, name+" (not a valid episode name)")
		case e.Type()&fs.ModeSymlink != 0 || !e.IsDir():
			bad = append(bad, name+" (not a real directory)")
		case !isGitWorktree(filepath.Join(root, name)):
			bad = append(bad, name+" (not a git worktree)")
		default:
			good = append(good, name)
		}
	}
	r.add(markOK, "workspace", fmt.Sprintf("%s: %d episode worktree(s)%s", root, len(good), listSuffix(good)), "")
	if len(bad) > 0 {
		r.add(markWarn, "workspace", "not usable as episodes: "+strings.Join(bad, ", "),
			"an episode is <root>/<ep>, a git worktree (session new makes one)")
	}
}

func listSuffix(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return ": " + strings.Join(names, ", ")
}

// isGitWorktree reports whether dir has the `.git` FILE a linked worktree
// carries ("gitdir: …").
func isGitWorktree(dir string) bool {
	f, err := os.Open(filepath.Join(dir, ".git"))
	if err != nil {
		return false
	}
	defer func() { _ = f.Close() }()
	head := make([]byte, 8)
	n, _ := io.ReadFull(f, head)
	return string(head[:n]) == "gitdir: "
}

func doctorOnline(r *doctorReport, roles []pinRole) {
	ctx, cancel := budgetContext(time.Minute)
	defer cancel()
	src := setupSource()
	for _, role := range roles {
		pin, err := lookupPin(role.release, setupPlatform)
		if err != nil {
			continue
		}
		sum, err := pinfetch.ReleaseChecksum(ctx, src, pin)
		switch {
		case err != nil:
			r.add(markFail, "online", fmt.Sprintf("%s AILANG %s: %v", role.name, pin.Release, err), "")
		case sum != pin.TarballSHA256:
			r.add(markFail, "online", fmt.Sprintf("%s AILANG %s: the release .sha256 says %s, the compiled-in pin is %s",
				role.name, pin.Release, sum, pin.TarballSHA256), "do not install; report it")
		default:
			r.add(markOK, "online", fmt.Sprintf("%s AILANG %s: release .sha256 reachable and matches the pin", role.name, pin.Release), "")
		}
	}
}
