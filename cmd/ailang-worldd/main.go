// Command ailang-worldd is the AILANG World local daemon and its CLI client
// (w-worldd-m2, Decision 5).
//
// One binary, subcommand style. `serve` takes SOLE writer authority over a world
// store and exposes the loopback-only worldd-native REST surface; every other
// verb is a thin REST client that dials that daemon. There is exactly one code
// path to the store, so CLI use continuously exercises the REST surface.
//
//	ailang-worldd serve --db <path> [--bind 127.0.0.1:7644] [--ailang-bin <path>]
//	                    [--workspace-root <dir> --tool-ailang-bin <path>]
//	                    [--examples-dir <dir>]
//	                    [--run-allow-caps Env,Net,Declassify]
//	                    [--run-net-allow 127.0.0.1:PORT ...] [--run-net-allow-http]
//	                    [--exec-profile <file> ... --exec-sandbox <dir>]
//	                    [--exec-episode-project EP=PROJECT ...] [--exec-node <path>]
//	                    [--exec-max-output-bytes N]
//	ailang-worldd [--addr http://127.0.0.1:7644] health
//	ailang-worldd [--addr http://127.0.0.1:7644] head
//	ailang-worldd [--addr http://127.0.0.1:7644] world get <ref>
//	ailang-worldd [--addr http://127.0.0.1:7644] object get <ref> [--payload]
//	ailang-worldd [--addr http://127.0.0.1:7644] object find <semanticId> [--after <ref>] [--limit N]
//	ailang-worldd [--addr http://127.0.0.1:7644] log get <index>
//	ailang-worldd [--addr http://127.0.0.1:7644] log range --from N [--limit M]
//	ailang-worldd [--addr http://127.0.0.1:7644] registry get <name>
//	ailang-worldd [--addr http://127.0.0.1:7644] commit --file <commit.json> [--session <file|token>]
//	ailang-worldd [--addr http://127.0.0.1:7644] log tail [--from N] [--follow] [--interval 1s] [--raw]
//	ailang-worldd [--addr http://127.0.0.1:7644] tools list [--session <file|token>] [--json]
//	ailang-worldd [--addr http://127.0.0.1:7644] call <tool> [--arg k=v]... [--arg-json k=<json>]... | --json <obj>|@file|- [--json-out] [--strict]
//	ailang-worldd [--addr http://127.0.0.1:7644] why <index|head|sha256:...|a2a:...|rest:...|-> | --result <file> [--scan N] [--json]
//	ailang-worldd [--addr http://127.0.0.1:7644] provenance [--since <entry>] [--episode <ep>]
//	ailang-worldd setup [--interpreter-dir <dir>] [--tools-dir <dir>] [--db <path>] [--workspace-root <dir>] [--from-dir <dir>] [--replace]
//	ailang-worldd [--addr http://127.0.0.1:7644] doctor [--db <path>] [--workspace-root <dir>] [--online] ...
//	ailang-worldd session new|list|revoke|mint ...
//	ailang-worldd help [<verb>]
//
// `--addr` is ONE GLOBAL CLIENT FLAG available to every client verb; it is not a
// `serve` flag, and passing it to `serve` is a usage error rather than a silently
// ignored argument.
//
// Exit codes: 0 success, 1 usage or client error, 2 fatal startup/runtime,
// 3 integrity refusal (a broken provenance link in why; call --strict on a
// committed ok:false; a setup digest mismatch).
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
)

// Exit codes (Decision 5).
const (
	exitOK    = 0
	exitUsage = 1
	exitFatal = 2
)

const usage = `ailang-worldd — AILANG World local daemon (loopback only)

Usage:
  ailang-worldd serve --db <path> [--bind host:port] [--ailang-bin <path>]
                      [--workspace-root <dir> --tool-ailang-bin <path>]
                      [--examples-dir <dir>]
                      [--run-allow-caps Env,Net,Declassify]
                      [--run-net-allow 127.0.0.1:PORT ...] [--run-net-allow-http]
                      [--exec-profile <file> ... --exec-sandbox <dir>]
                      [--exec-episode-project EP=PROJECT ...] [--exec-node <path>]
                      [--exec-max-output-bytes N]
  ailang-worldd [--addr <url>] health
  ailang-worldd [--addr <url>] head
  ailang-worldd [--addr <url>] world get <ref>
  ailang-worldd [--addr <url>] object get <ref> [--payload]
  ailang-worldd [--addr <url>] object find <semanticId> [--after <ref>] [--limit N]
  ailang-worldd [--addr <url>] log get <index>
  ailang-worldd [--addr <url>] log range --from N [--limit M]
  ailang-worldd [--addr <url>] registry get <name>
  ailang-worldd [--addr <url>] log tail [--from N] [--follow] [--interval 1s] [--raw]
  ailang-worldd [--addr <url>] commit --file <commit.json> [--session <file|token>]
  ailang-worldd [--addr <url>] tools list [--session <file|token>] [--json]
  ailang-worldd [--addr <url>] call <tool> [--session <file|token>]
                    [--arg k=v]... [--arg-json k=<json>]... | --json <obj>|@file|-
                    [--json-out] [--strict]
  ailang-worldd [--addr <url>] why <index|head|sha256:<ref>|a2a:<id>|rest:<id>|->
                    [--result <file>] [--scan N] [--json]
  ailang-worldd [--addr <url>] provenance [--since <entry>] [--episode <ep>] [--scan N]
  ailang-worldd setup [--interpreter-dir <dir>] [--tools-dir <dir>] [--db <path>]
                    [--workspace-root <dir>] [--from-dir <dir>] [--replace]
  ailang-worldd [--addr <url>] doctor [--db <path>] [--workspace-root <dir>]
                    [--interpreter-dir <dir>] [--tools-dir <dir>]
                    [--examples-dir <dir>] [--online]
  ailang-worldd session new <episode> [--db <path>] [--workspace-root <dir>]
                    [--repo <dir>] [--preset se-tools] [--grant EFFECT=SCOPE:BUDGET]...
                    [--budget 50] [--ttl 3600] [--out <file>] [--branch <name>]
  ailang-worldd session list [--db <path>] [--episode <ep>] [--json]
  ailang-worldd session revoke --db <path> <credential_id-hash> | --episode <ep>
  ailang-worldd session mint --db <path> --episode <ep> --grant EFFECT=SCOPE:BUDGET...
                    [--ttl 3600] [--out <file>]
  ailang-worldd help [<verb>]

  '<verb> --help' and 'help <verb>' print a verb's help (exit 0).

Session credential (tools, call, commit): --session <file> (a file holding
the 64-hex token, mode 0600) or the token itself (warns: visible on argv),
else $WORLD_SESSION. The token is never printed.

Global client flag:
  --addr <url>   base URL of the daemon (default ` + daemon.DefaultAddr + `).
                 Applies to client verbs only; it is NOT a 'serve' flag.

serve flags:
  --db <path>          world store database (required)
  --bind host:port     loopback listen address (default ` + daemon.DefaultBind + `);
                       a non-loopback host is refused — there is no override
  --ailang-bin <path>  interpreter to archive and pin at startup (optional)
  --workspace-root <dir>
                       episode worktrees live at <dir>/<episode> (made by the
                       operator with git worktree add before session mint); it
                       must not contain the store, its archive, its rendered
                       policies or its tool cache — startup refuses otherwise
  --tool-ailang-bin <path>
                       AILANG binary the workspace tools run (must be
                       ` + daemon.ToolBinaryRelease + `); archived and hash-verified like
                       --ailang-bin. The Workspace.*/Ailang.* tools are served
                       only when both this and --workspace-root are set.
                       Workspace.Exec (workspace-exec) is bound with them; it
                       refuses every call "no exec profile configured" unless
                       --exec-profile is given
  --examples-dir <dir> AILANG examples corpus examples-search reads (passed
                       to the tool as AILANG_EXAMPLES; World never falls back
                       to a corpus in the binary). Default: ~/.ailang/examples
                       when it exists, else examples-search refuses "no
                       examples corpus configured". Must be outside
                       --workspace-root
  --run-allow-caps Env,Net,Declassify
                       extra capabilities an ailang-run may request beyond
                       IO and FS (default none). An Env run is the effect
                       Ailang.RunEnv, a Net run Ailang.RunNet; each needs
                       its own session grant
  --run-net-allow 127.0.0.1:PORT
                       a loopback IP:PORT a Net run may reach (repeatable;
                       required with Net). Every other host, port and
                       redirect hop is refused; a bare host, a name, or a
                       non-loopback address is refused at startup
  --run-net-allow-http allow plain http to the --run-net-allow pairs
  --exec-profile <file>
                       an operator exec profile (world/exec-profile/v1 JSON;
                       repeatable, one per project): the commands
                       workspace-exec may run, sandboxed by srt. Needs the
                       workspace tools and --exec-sandbox. Read once at
                       startup; it, srt, node and every cache seed must lie
                       outside --workspace-root and the state dir. Each
                       profile must pass the startup probe (seven arms,
                       arm1-write-inside … arm7-toolchain) or serve refuses
  --exec-episode-project EP=PROJECT
                       run episode EP's workspace-exec under the profile
                       whose project is PROJECT (repeatable; optional with
                       one profile, which every episode then uses)
  --exec-sandbox <dir> the node_modules holding @anthropic-ai/sandbox-runtime
                       ` + broker.ExecSandboxRelease + ` (its dist/cli.js must hash to the pin);
                       World archives the whole tree and re-verifies it
                       before every call. Required with --exec-profile
  --exec-node <path>   the node (>= 20.11) that runs srt and the probe
                       client (default: node on PATH, resolved once)
  --exec-max-output-bytes N
                       kill a workspace-exec command whose stdout or stderr
                       passes N bytes (default 67108864, 64 MiB)

Exit codes: 0 ok, 1 usage or client error, 2 fatal startup,
            3 integrity refusal (why: a broken link; call --strict: ok:false;
              setup: a digest mismatch, nothing installed).
`

// cliStdin is what `call --json -` and `why -` read.
var cliStdin io.Reader = os.Stdin

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

// run is main's testable body: it takes its arguments and streams explicitly so
// nothing about the CLI depends on process globals.
func run(args []string, stdout, stderr io.Writer) int {
	// STARTUP REFUSAL (w-self-mod-vertical Decision 4 / AC10). The public
	// AILANG registry is immutable, so an ambient AILANG_REGISTRY_API_KEY
	// hands unrecallable publish authority to every process this binary ever
	// forks — and to every agent and shell command that inherits from them.
	// The refusal is here, ahead of flag parsing, so it applies to `serve` and
	// to every client verb alike: there is no subcommand that is exempt.
	//
	// The message names the VARIABLE, never the value.
	//
	// The one exemption is a help request (keyguard.go): its text is resolved
	// from the argv alone and no verb runs, so there is no process to inherit
	// the key. Everything else refuses, and the refusal ends with the exact
	// `env -u` command to run instead.
	if err := broker.AssertNoAmbientRegistryCredential(os.Environ()); err != nil {
		if text, ok := helpOnly(args); ok {
			fmt.Fprint(stdout, text)
			return exitOK
		}
		fmt.Fprintf(stderr, "ailang-worldd: %v\n%s\n%s\n", err, registryKeyRationale, registryKeyFix(args))
		return exitFatal
	}

	globals := flag.NewFlagSet("ailang-worldd", flag.ContinueOnError)
	globals.SetOutput(stderr)
	globals.Usage = func() { fmt.Fprint(stderr, usage) }
	addr := globals.String("addr", daemon.DefaultAddr,
		"base URL of the daemon for client verbs (global client flag)")
	if err := globals.Parse(args); err != nil {
		return exitUsage
	}

	// Was --addr given explicitly? `serve` must reject it rather than accept a
	// flag that has no meaning for it — that is what makes "--addr is a client
	// flag, not a serve flag" structurally true instead of documentation.
	addrGiven := false
	globals.Visit(func(f *flag.Flag) {
		if f.Name == "addr" {
			addrGiven = true
		}
	})

	rest := globals.Args()
	if len(rest) == 0 {
		fmt.Fprint(stderr, usage)
		return exitUsage
	}
	// `help <verb...>` is `<verb...> --help` (row 138 AC5.7).
	if rest[0] == "help" && len(rest) > 1 {
		rest = append(append([]string(nil), rest[1:]...), "--help")
	}
	// Verbs with no help of their own print their usage lines (exit 0).
	if text, ok := genericVerbHelp(rest); ok {
		fmt.Fprint(stdout, text)
		return exitOK
	}

	switch verb := rest[0]; verb {
	case "serve":
		if addrGiven {
			fmt.Fprintln(stderr, "ailang-worldd: --addr is a client flag and is not valid for 'serve'; "+
				"use --bind to choose the listen address")
			return exitUsage
		}
		return runServe(rest[1:], stdout, stderr)

	case "health":
		return runClientGet(*addr, "/v1/health", rest[1:], stdout, stderr)

	case "head":
		return runClientGet(*addr, "/v1/head", rest[1:], stdout, stderr)

	case "world":
		return runWorld(*addr, rest[1:], stdout, stderr)

	case "object":
		return runObject(*addr, rest[1:], stdout, stderr)

	case "log":
		return runLog(*addr, rest[1:], stdout, stderr)

	case "registry":
		return runRegistry(*addr, rest[1:], stdout, stderr)

	case "commit":
		return runCommit(*addr, rest[1:], stdout, stderr)

	// Row 138 M1/M2 and D-WORLD-60: the MCP client, the provenance walk,
	// the log tail and the PR trailer.
	case "tools":
		return runTools(*addr, rest[1:], stdout, stderr)

	case "call":
		return runCall(*addr, rest[1:], cliStdin, stdout, stderr)

	case "why":
		return runWhy(*addr, rest[1:], cliStdin, stdout, stderr)

	case "provenance":
		return runProvenance(*addr, rest[1:], stdout, stderr)

	case "doctor":
		return runDoctor(*addr, rest[1:], stdout, stderr)

	case "setup":
		if addrGiven {
			fmt.Fprintln(stderr, "ailang-worldd: --addr is a client flag and is not valid for 'setup'")
			return exitUsage
		}
		return runSetup(rest[1:], stdout, stderr)

	case "session":
		// (w-session-authority D1/D4) session mint|revoke speak directly to the
		// store DB path via --db, never to a daemon HTTP endpoint, so --addr is
		// not a session flag — reject it exactly as serve rejects --addr.
		if addrGiven {
			fmt.Fprintln(stderr, "ailang-worldd: --addr is a client flag and is not valid for 'session'; "+
				"session mint/revoke act on --db directly, not on a running daemon")
			return exitUsage
		}
		return runSession(rest[1:], stdout, stderr)

	case "help":
		fmt.Fprint(stdout, usage)
		return exitOK

	default:
		fmt.Fprintf(stderr, "ailang-worldd: unknown command %q\n\n", verb)
		fmt.Fprint(stderr, usage)
		return exitUsage
	}
}

// runServe parses the serve flags and drives the daemon lifecycle until
// SIGINT/SIGTERM. The resolved listen address is announced on stdout by
// daemon.Run once the socket is bound.
func runServe(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("ailang-worldd serve", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { fmt.Fprint(stderr, usage) }
	dbPath := fs.String("db", "", "world store database (required)")
	bind := fs.String("bind", daemon.DefaultBind, "loopback listen address host:port")
	ailangBin := fs.String("ailang-bin", "", "interpreter to archive and pin at startup")
	workspaceRoot := fs.String("workspace-root", "", "directory holding one worktree per episode")
	toolAilangBin := fs.String("tool-ailang-bin", "", "AILANG binary the workspace tools run")
	examplesDir := fs.String("examples-dir", "", "AILANG examples corpus for examples-search (default ~/.ailang/examples when it exists)")
	var runAllowCaps, runNetAllow []string
	fs.Func("run-allow-caps", "extra capabilities an ailang-run may request (Env,Net,Declassify)", func(v string) error {
		for _, name := range strings.Split(v, ",") {
			if name = strings.TrimSpace(name); name == "" {
				return fmt.Errorf("empty capability name in %q", v)
			}
			runAllowCaps = append(runAllowCaps, strings.TrimSpace(name))
		}
		return nil
	})
	fs.Func("run-net-allow", "a loopback IP:PORT a Net run may reach (repeatable)", func(v string) error {
		runNetAllow = append(runNetAllow, v)
		return nil
	})
	runNetAllowHTTP := fs.Bool("run-net-allow-http", false, "allow plain http to the --run-net-allow pairs")
	var execProfiles, execEpisodeProjects []string
	fs.Func("exec-profile", "an operator exec profile (repeatable, one per project)", func(v string) error {
		execProfiles = append(execProfiles, v)
		return nil
	})
	fs.Func("exec-episode-project", "EP=PROJECT: the profile an episode's workspace-exec runs under (repeatable)", func(v string) error {
		execEpisodeProjects = append(execEpisodeProjects, v)
		return nil
	})
	execSandbox := fs.String("exec-sandbox", "", "node_modules holding the pinned @anthropic-ai/sandbox-runtime")
	execNode := fs.String("exec-node", "", "the node that runs srt (default: node on PATH)")
	var execMaxOutput int64
	fs.Func("exec-max-output-bytes", "per-stream output kill for workspace-exec (default 67108864)", func(v string) error {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n < 1 {
			return fmt.Errorf("%q is not a positive byte count", v)
		}
		execMaxOutput = n
		return nil
	})
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if extra := fs.Args(); len(extra) > 0 {
		fmt.Fprintf(stderr, "ailang-worldd serve: unexpected argument %q\n", extra[0])
		return exitUsage
	}
	if *dbPath == "" {
		fmt.Fprintln(stderr, "ailang-worldd serve: --db is required")
		return exitUsage
	}

	host, portText, err := net.SplitHostPort(*bind)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd serve: --bind %q is not host:port: %v\n", *bind, err)
		return exitUsage
	}
	port, err := strconv.Atoi(portText)
	if err != nil || port < 0 || port > 65535 {
		fmt.Fprintf(stderr, "ailang-worldd serve: --bind %q has an invalid port %q\n", *bind, portText)
		return exitUsage
	}

	// SIGINT/SIGTERM cancel the context, which starts the D7 bounded drain.
	ctx, stop := serveSignalContext()
	defer stop()

	cfg := daemon.Config{DBPath: *dbPath, BindHost: host, BindPort: port, AilangBin: *ailangBin,
		WorkspaceRoot: *workspaceRoot, ToolAilangBin: *toolAilangBin,
		ExamplesDir:  resolveExamplesDefault(*examplesDir, os.UserHomeDir),
		RunAllowCaps: runAllowCaps, RunNetAllow: runNetAllow, RunNetAllowHTTP: *runNetAllowHTTP,
		ExecProfiles: execProfiles, ExecEpisodeProjects: execEpisodeProjects, ExecSandbox: *execSandbox,
		ExecNode: *execNode, ExecMaxOutputBytes: execMaxOutput}
	return serveResult(runDaemon(ctx, cfg, stdout), stderr)
}

// runDaemon is daemon.Run; a test replaces it to see the Config serve built.
var runDaemon = daemon.Run

// resolveExamplesDefault is --examples-dir's default, resolved once at
// startup: an explicit flag wins; otherwise the operator's own
// ~/.ailang/examples (what `ailang examples download` populates) when it is a
// directory; otherwise "" (no corpus, examples-search refuses). The tool's
// HOME is the per-episode cache, so the binary's own ~/.ailang fallback can
// never see the operator's corpus (V65).
func resolveExamplesDefault(flagValue string, home func() (string, error)) string {
	if flagValue != "" {
		return flagValue
	}
	h, err := home()
	if err != nil || h == "" {
		return ""
	}
	dir := filepath.Join(h, ".ailang", "examples")
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		return dir
	}
	return ""
}

func serveResult(err error, stderr io.Writer) int {
	if err == nil {
		return exitOK
	}
	fmt.Fprintf(stderr, "ailang-worldd: %v\n", err)
	var startup *daemon.StartupError
	if errors.As(err, &startup) && startup.Settled != nil {
		<-startup.Settled
	}
	return exitFatal
}

// ownHelpVerbs handle --help themselves (each prints its own help, exit 0).
var ownHelpVerbs = map[string]bool{
	"tools": true, "call": true, "why": true, "provenance": true,
	"setup": true, "doctor": true, "session": true,
}

// genericVerbHelp answers `<verb> ... --help` for the verbs that have no help
// text of their own, from the usage block itself, so the two cannot drift.
func genericVerbHelp(rest []string) (string, bool) {
	verb := rest[0]
	if ownHelpVerbs[verb] || (verb == "log" && len(rest) > 1 && rest[1] == "tail") {
		return "", false
	}
	wants := false
	for _, a := range rest[1:] {
		if a == "--help" || a == "-help" || a == "-h" {
			wants = true
		}
	}
	if !wants {
		return "", false
	}
	lines := usageLinesFor(verb)
	if len(lines) == 0 {
		return "", false
	}
	text := "usage:\n" + strings.Join(lines, "\n") + "\n"
	if verb == "serve" {
		if i := strings.Index(usage, "serve flags:"); i >= 0 {
			flags := usage[i:]
			if j := strings.Index(flags, "\nExit codes:"); j >= 0 {
				flags = flags[:j+1]
			}
			text += "\n" + flags
		}
	} else {
		text += "\n--addr <url> is the daemon's base URL (default " + daemon.DefaultAddr + ").\n"
	}
	return text, true
}

// usageLinesFor returns the usage lines (with continuation lines) of verb.
func usageLinesFor(verb string) []string {
	var out []string
	in := false
	for _, l := range strings.Split(usage, "\n") {
		if strings.HasPrefix(l, "  ailang-worldd ") {
			cmd := strings.TrimPrefix(strings.TrimPrefix(l, "  ailang-worldd "), "[--addr <url>] ")
			in = cmd == verb || strings.HasPrefix(cmd, verb+" ")
			if in {
				out = append(out, l)
			}
			continue
		}
		if in && strings.HasPrefix(l, "                    ") {
			out = append(out, l)
			continue
		}
		in = false
	}
	return out
}
