package main

// Row 138 M3 (design_docs/planned/w-worldd-developer-cli.md §3.5): `setup`
// fetches and verifies both pins, then makes the store directory and the
// workspace root. Every byte it installs passed three tarball digests (the
// release's .sha256, the streamed hash, the compiled-in pin) and the
// compiled-in binary digest — host/pinfetch does the work. setup never
// executes a downloaded byte, and it never runs a human act (publish, mint).

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/pinfetch"
)

const setupHelp = `usage: ailang-worldd setup [--interpreter-dir <dir>] [--tools-dir <dir>]
           [--db <path>] [--workspace-root <dir>] [--from-dir <dir>] [--replace]

Installs the two pinned AILANG binaries World runs, verified byte for byte:
  interpreter  ` + interpreterRelease + `  -> <interpreter-dir>/ailang  (serve --ailang-bin)
  tool         <tool>   -> <tools-dir>/<tool>/ailang  (serve --tool-ailang-bin)
where <tool> is the tool-binary release this daemon is built for.

For each pin: a file already hashing to the pin is left alone (no network
request). Otherwise the release's .sha256 and tarball are fetched over https
from github.com/sunholo-data/ailang (the URL is compiled in), and the install
happens only when the release .sha256, the computed tarball sha256 and the
compiled-in digest all agree and the extracted ailang hashes to the
compiled-in binary digest. No tarball is kept; pin.json records what was
installed. A mismatching existing file is refused unless --replace, which
keeps it as ailang.prev-<sha8>. setup never runs a downloaded byte.

Then it creates the store directory and the workspace root (mode 0700),
refusing a workspace root that contains the store, and prints the attended
steps that follow.

  --interpreter-dir <dir>  default ~/.pinned-ailang
  --tools-dir <dir>        default ~/.pinned-ailang-tools
  --db <path>              world store (default ~/.ailang/world/world.db);
                           only its directory is created, never the store
  --workspace-root <dir>   default ~/.ailang/world-ws
  --from-dir <dir>         install offline from <dir>/<release>/<asset> and
                           its .sha256, verified the same way
  --replace                replace a mismatching existing binary

Platforms: darwin/arm64 and linux/amd64.
Exit: 0 ok; 1 usage, unsupported platform, network refusal or an existing
mismatching file; 3 a digest mismatch (nothing installed).
`

// setupBudget bounds the whole command: two tarballs at the 5-minute
// per-tarball budget plus the small assets.
const setupBudget = 12 * time.Minute

// Test seams: production never assigns them.
var (
	setupSource   = pinfetch.GitHub
	setupPlatform = runtime.GOOS + "/" + runtime.GOARCH
)

// worldPaths is the D-CLI-4 default layout under $HOME.
type worldPaths struct {
	interpDir, toolsDir, db, workspaceRoot string
}

func defaultWorldPaths() worldPaths {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		home = "."
	}
	return worldPaths{
		interpDir:     filepath.Join(home, ".pinned-ailang"),
		toolsDir:      filepath.Join(home, ".pinned-ailang-tools"),
		db:            filepath.Join(home, ".ailang", "world", "world.db"),
		workspaceRoot: filepath.Join(home, ".ailang", "world-ws"),
	}
}

func runSetup(args []string, stdout, stderr io.Writer) int {
	def := defaultWorldPaths()
	fs := flag.NewFlagSet("ailang-worldd setup", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	interpDir := fs.String("interpreter-dir", def.interpDir, "interpreter pin directory")
	toolsDir := fs.String("tools-dir", def.toolsDir, "tool pin parent directory")
	db := fs.String("db", def.db, "world store path (only its directory is created)")
	wsRoot := fs.String("workspace-root", def.workspaceRoot, "workspace root")
	fromDir := fs.String("from-dir", "", "install offline from this directory")
	replace := fs.Bool("replace", false, "replace a mismatching existing binary")
	if code, done := parseVerbFlags(fs, args, setupHelp, stdout, stderr); done {
		return code
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(stderr, "ailang-worldd setup: unexpected argument %q\n", fs.Arg(0))
		return exitUsage
	}

	// Validate everything before the first side effect.
	pins := make([]pinfetch.Pin, 0, 2)
	roles := pinRoles(*interpDir, *toolsDir)
	for _, role := range roles {
		pin, err := lookupPin(role.release, setupPlatform)
		if err != nil {
			fmt.Fprintf(stderr, "ailang-worldd setup: %v\n", err)
			return exitUsage
		}
		pins = append(pins, pin)
	}
	storeDir, err := filepath.Abs(filepath.Dir(*db))
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd setup: --db: %v\n", err)
		return exitUsage
	}
	root, err := filepath.Abs(*wsRoot)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd setup: --workspace-root: %v\n", err)
		return exitUsage
	}
	if within(root, storeDir) {
		fmt.Fprintf(stderr, "ailang-worldd setup: refusing: the workspace root %s contains the store directory %s; "+
			"an agent's worktree must never reach the store (serve refuses this layout too)\n", root, storeDir)
		return exitUsage
	}

	src := setupSource()
	if *fromDir != "" {
		src = pinfetch.FromDir(*fromDir)
	}
	ctx, cancel := context.WithTimeout(context.Background(), setupBudget)
	defer cancel()
	for i, role := range roles {
		pin := pins[i]
		res, err := pinfetch.Install(ctx, src, pin, role.dest, pinfetch.Options{Replace: *replace})
		if err != nil {
			fmt.Fprintf(stderr, "✗ %s AILANG %s (%s): %v\n", role.name, pin.Release, pin.Platform, err)
			if pinfetch.IsIntegrity(err) {
				fmt.Fprintln(stderr, "  nothing was installed; the release does not match the compiled-in pin")
				return exitIntegrity
			}
			return exitUsage
		}
		switch res.Status {
		case pinfetch.StatusPresent:
			fmt.Fprintf(stdout, "✓ %-11s AILANG %s present   %s (sha256 %s…)\n", role.name, pin.Release, role.dest, pin.BinarySHA256[:12])
		case pinfetch.StatusReplaced:
			fmt.Fprintf(stdout, "✓ %-11s AILANG %s installed %s (sha256 %s…); kept the old file as %s\n",
				role.name, pin.Release, role.dest, pin.BinarySHA256[:12], res.Kept)
		default:
			fmt.Fprintf(stdout, "✓ %-11s AILANG %s installed %s (sha256 %s…)\n", role.name, pin.Release, role.dest, pin.BinarySHA256[:12])
		}
	}
	for _, d := range []struct{ name, path string }{{"store directory", storeDir}, {"workspace root", root}} {
		if err := os.MkdirAll(d.path, 0o700); err != nil {
			fmt.Fprintf(stderr, "✗ %s %s: %v\n", d.name, d.path, err)
			return exitFatal
		}
		fmt.Fprintf(stdout, "✓ %-15s %s\n", d.name, d.path)
	}
	printNextSteps(stdout, roles, *db, root)
	return exitOK
}

// within reports whether path is dir or lies beneath it.
func within(dir, path string) bool {
	rel, err := filepath.Rel(dir, path)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}

func printNextSteps(w io.Writer, roles []pinRole, db, root string) {
	interp, tool := roles[0].dest, roles[1].dest
	fmt.Fprintf(w, `
next (the attended steps; docs/QUICKSTART.md §9 has the full walk):
  1. publish the se-tools transitions (a human act at a terminal), from the repo root:
       world-publish transitions --store %[1]s \
         --manifest packages/se-tools/transitions.json --ailang-bin %[3]s
  2. mint a session for an episode worktree under %[2]s
       (ailang-worldd session mint --db %[1]s --episode <ep> --grant ... --out <file>)
  3. ailang-worldd serve --db %[1]s --ailang-bin %[3]s \
       --workspace-root %[2]s --tool-ailang-bin %[4]s
`, db, root, interp, tool)
}
