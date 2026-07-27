// Command ailang-worldd is the AILANG World local daemon and its CLI client
// (w-worldd-m2, Decision 5).
//
// One binary, subcommand style. `serve` takes SOLE writer authority over a world
// store and exposes the loopback-only worldd-native REST surface; every other
// verb is a thin REST client that dials that daemon. There is exactly one code
// path to the store, so CLI use continuously exercises the REST surface.
//
//	ailang-worldd serve --db <path> [--bind 127.0.0.1:7644] [--ailang-bin <path>]
//	ailang-worldd [--addr http://127.0.0.1:7644] health
//	ailang-worldd [--addr http://127.0.0.1:7644] head
//
// `--addr` is ONE GLOBAL CLIENT FLAG available to every client verb; it is not a
// `serve` flag, and passing it to `serve` is a usage error rather than a silently
// ignored argument. The remaining client verbs of the frozen route table
// (world/object/log/registry/commit) land in M2.C alongside the real-subprocess
// end-to-end test; M2.A deliberately ships health and head only rather than
// half-implementing the rest.
//
// Exit codes: 0 success, 1 usage or client error, 2 fatal startup/runtime.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

// Exit codes (Decision 5).
const (
	exitOK    = 0
	exitUsage = 1
	exitFatal = 2
)

// maxClientResponseBytes bounds how much of a daemon response the client will
// buffer. Decision 7 bounds every request-driven allocation; a client reading an
// unbounded body from a wedged or hostile server would be the same class of
// defect on the other side of the socket. It matches the D7 commit-body cap.
const maxClientResponseBytes = 8388608

const usage = `ailang-worldd — AILANG World local daemon (loopback only)

Usage:
  ailang-worldd serve --db <path> [--bind host:port] [--ailang-bin <path>]
  ailang-worldd [--addr <url>] health
  ailang-worldd [--addr <url>] head

Global client flag:
  --addr <url>   base URL of the daemon (default ` + daemon.DefaultAddr + `).
                 Applies to client verbs only; it is NOT a 'serve' flag.

serve flags:
  --db <path>          world store database (required)
  --bind host:port     loopback listen address (default ` + daemon.DefaultBind + `);
                       a non-loopback host is refused — there is no override
  --ailang-bin <path>  interpreter to archive and pin at startup (optional)

Client verbs for the remaining routes (world/object/log/registry/commit) arrive
in milestone M2.C.

Exit codes: 0 ok, 1 usage or client error, 2 fatal startup.
`

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

// run is main's testable body: it takes its arguments and streams explicitly so
// nothing about the CLI depends on process globals.
func run(args []string, stdout, stderr io.Writer) int {
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := daemon.Config{DBPath: *dbPath, BindHost: host, BindPort: port, AilangBin: *ailangBin}
	if err := daemon.Run(ctx, cfg, stdout); err != nil {
		fmt.Fprintf(stderr, "ailang-worldd: %v\n", err)
		return exitFatal
	}
	return exitOK
}

// client is the CLI's REST caller. It owns one http.Client and derives a
// per-call deadline, so no client call can hang past the D7 budget. timeout is a
// struct field rather than a constant so M2.C's "a non-responding server yields
// a deadline error, not a hang" test does not have to wait 30 s.
type client struct {
	base    string
	timeout time.Duration
	http    *http.Client
}

func newClient(base string) *client {
	return &client{
		base:    strings.TrimSuffix(base, "/"),
		timeout: daemon.DefaultClientTimeout,
		http:    &http.Client{},
	}
}

// get performs one bounded GET and returns the status code and body text.
func (c *client) get(path string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+path, nil)
	if err != nil {
		return 0, "", fmt.Errorf("build request for %s%s: %w", c.base, path, err)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("GET %s%s: %w", c.base, path, err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxClientResponseBytes))
	if err != nil {
		return resp.StatusCode, "", fmt.Errorf("read %s%s response: %w", c.base, path, err)
	}
	return resp.StatusCode, string(body), nil
}

// runClientGet is the shared body of the read-only client verbs: dial, print the
// response, and map a non-2xx status onto the client-error exit code.
func runClientGet(addr, path string, args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		fmt.Fprintf(stderr, "ailang-worldd: unexpected argument %q\n", args[0])
		return exitUsage
	}
	code, body, err := newClient(addr).get(path)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd: %v\n", err)
		return exitUsage
	}
	if code < 200 || code >= 300 {
		fmt.Fprintf(stderr, "ailang-worldd: %s returned HTTP %d: %s\n", path, code, strings.TrimSpace(body))
		return exitUsage
	}
	fmt.Fprintln(stdout, strings.TrimSpace(body))
	return exitOK
}
