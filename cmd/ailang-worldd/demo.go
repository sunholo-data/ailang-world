package main

import (
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/sunholo-data/ailang-world/host/pins"
)

var demoHelp = `usage: ailang-worldd [--addr URL] demo seed --db PATH --session FILE
    --ailang-bin PIN --proof-key-file KEY --manifest FILE [--addr URL]

Prepare a disposable demo worktree and publish descriptors, then mint the demo
session at a controlling terminal (QUICKSTART section 12). Seed never mints,
publishes, creates a worktree or opens another writable database.
--db identifies the daemon's store and permits read-only grant inspection.
--session is a private session file (default WORLD_SESSION).
--ailang-bin is the capsule/proof pin (` + pins.ReleaseFor(pins.RoleProof).Release + `); serve's separate
--tool-ailang-bin must be ` + toolRelease() + ` for the workspace tools.
--proof-key-file is the private operator key required by later proof stages.
--manifest is a new output path; existing files are refused.
--addr is the daemon URL, also accepted as the global client flag.

This incremental build implements preflight and S1 domain commits only.
S2-S6 incomplete: the command refuses completion and writes no success manifest.
Do not retry a partial seed; inspect its recorded history first.
`

type demoOptions struct{ db, session, pin, key, manifest, addr string }

func runDemo(addr string, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] != "seed" {
		fmt.Fprint(stderr, demoHelp)
		return exitUsage
	}
	fs := flag.NewFlagSet("ailang-worldd demo seed", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	var o demoOptions
	fs.StringVar(&o.db, "db", "", "daemon database identity")
	fs.StringVar(&o.session, "session", "", "private session file")
	fs.StringVar(&o.pin, "ailang-bin", "", "capsule/proof pin "+pins.ReleaseFor(pins.RoleProof).Release)
	fs.StringVar(&o.key, "proof-key-file", "", "private proof key for later stages")
	fs.StringVar(&o.manifest, "manifest", "", "new success manifest path")
	fs.StringVar(&o.addr, "addr", addr, "daemon URL")
	if code, done := parseVerbFlags(fs, args[1:], demoHelp, stdout, stderr); done {
		return code
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(stderr, "ailang-worldd demo seed: unexpected positional argument")
		return exitUsage
	}
	if o.db == "" || o.pin == "" || o.key == "" || o.manifest == "" {
		fmt.Fprintln(stderr, "ailang-worldd demo seed: --db, --ailang-bin, --proof-key-file and --manifest are required")
		return exitUsage
	}
	token, err := resolveSession(o.session, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd demo seed: %v\n", err)
		return exitUsage
	}
	ctx, cancel := budgetContext(90 * time.Second)
	defer cancel()
	if err = demoSeedDomain(ctx, o, token); err != nil {
		fmt.Fprintf(stderr, "ailang-worldd demo seed: %v\n", err)
		return exitUsage
	}
	fmt.Fprintln(stderr, "ailang-worldd demo seed: S1 domain commits recorded; S2-S6 incomplete; no success manifest written. Inspect history; do not blindly retry.")
	return exitUsage
}
