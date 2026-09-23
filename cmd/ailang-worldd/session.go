package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

// sessionEnv is the injectable human-act seam a session verb performs. It lets
// the mint flow, which is otherwise gated on a real controlling terminal, be
// driven from a test without a pty (AC-M1-1). The production path builds it from
// /dev/tty (D1, mirroring cmd/world-publish/tty.go).
type sessionEnv struct {
	// openTerminal opens the controlling-terminal device. A non-nil error means
	// there is no controlling terminal, and the mint verb refuses (the tty
	// fence). The returned object's reader feeds the y/N confirmation; its
	// writer carries the prompt.
	openTerminal func() (io.ReadWriteCloser, error)
	// now returns the current Unix time (seconds). Injectable for determinism.
	now func() int64
}

// realSessionEnv is the production seam: /dev/tty and the wall clock. It is
// the ONLY place the mint path touches the terminal.
func realSessionEnv() sessionEnv {
	return sessionEnv{
		openTerminal: func() (io.ReadWriteCloser, error) {
			return os.OpenFile("/dev/tty", os.O_RDWR, 0)
		},
		now: func() int64 { return time.Now().Unix() },
	}
}

// runSession dispatches the `session` verb family (D1 mint / D4 revoke). It is
// invoked from main.go's `case "session"` with the remaining args. Mint speaks
// directly to the store DB path (`--db`), never to a daemon HTTP endpoint
// (residual R13); revoke deletes the mapping row by credential_id hash.
func runSession(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "ailang-worldd session: usage: session mint | session revoke <credential_id-hash>")
		return exitUsage
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case "mint":
		return runSessionMint(rest, stdout, stderr, realSessionEnv())
	case "revoke":
		return runSessionRevoke(rest, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "ailang-worldd session: unknown subcommand %q\n", sub)
		return exitUsage
	}
}

// multiFlag is a repeatable string flag value (used by the repeatable --grant).
type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, ",") }
func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

// runSessionMint implements `session mint --episode <id> [--grant
// EFFECT=SCOPE:BUDGET]... [--ttl 3600] [--out <path>]`. It is tty-fenced (D1):
// it refuses when there is no controlling terminal, and asks a one-line y/N
// confirmation read from that terminal before minting. The raw 64-hex token
// prints/writes EXACTLY ONCE (stdout or `--out` at mode 0600) and is never
// logged; the credential_id hash is emitted separately as a diagnostic so the
// operator can revoke later.
//
// env is injected so the fence and confirmation can be driven in a test; run()
// always passes realSessionEnv().
func runSessionMint(args []string, stdout, stderr io.Writer, env sessionEnv) int {
	fs := flag.NewFlagSet("ailang-worldd session mint", flag.ContinueOnError)
	fs.SetOutput(stderr)
	episode := fs.String("episode", "", "episode the session binds to (required)")
	var grantSpecs multiFlag
	fs.Var(&grantSpecs, "grant", "EFFECT=SCOPE:BUDGET, repeatable, at least one")
	ttl := fs.Int64("ttl", 3600, "lifetime in whole seconds (default 3600)")
	outPath := fs.String("out", "", "write the credential ONCE to <path> at mode 0600 (default: print to stdout exactly once)")
	// `--db` default is "" — matching serve --db exactly (serve's --db is
	// required, so an empty value is refused below). A literal path default
	// would invent a store serve does not define; a credential mint must never
	// silently target a surprise store.
	dbPath := fs.String("db", "", "world store database (required)")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if extra := fs.Args(); len(extra) > 0 {
		fmt.Fprintf(stderr, "ailang-worldd session mint: unexpected argument %q\n", extra[0])
		return exitUsage
	}
	if *episode == "" {
		fmt.Fprintln(stderr, "ailang-worldd session mint: --episode is required")
		return exitUsage
	}
	if *dbPath == "" {
		fmt.Fprintln(stderr, "ailang-worldd session mint: --db is required (matching serve --db)")
		return exitUsage
	}
	if len(grantSpecs) == 0 {
		fmt.Fprintln(stderr, "ailang-worldd session mint: at least one --grant EFFECT=SCOPE:BUDGET is required")
		return exitUsage
	}
	grants, err := parseGrants(grantSpecs)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd session mint: %v\n", err)
		return exitUsage
	}
	if *ttl < 0 {
		fmt.Fprintf(stderr, "ailang-worldd session mint: --ttl must be non-negative, got %d\n", *ttl)
		return exitUsage
	}

	// TTY FENCE (D1): mint is one human act. A process with no controlling
	// terminal (a pasted/redirected drive) collapses "one human act" into
	// "whoever runs the script gets a minted session credential", so it refuses.
	term, err := env.openTerminal()
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd session mint: refusing: no controlling terminal (%v); "+
			"minting a session credential requires one human act at a terminal\n", err)
		return exitUsage
	}
	defer func() { _ = term.Close() }()

	now := env.now()
	fmt.Fprintf(term, "Confirm mint for episode %s (%d grant(s), expiry +%ds) with a session credential? [y/N] ",
		*episode, len(grants), *ttl)
	if !isYes(readTermLine(term)) {
		fmt.Fprintln(stderr, "ailang-worldd session mint: aborted (not confirmed)")
		return exitUsage
	}

	st, err := store.Open(*dbPath)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd session mint: cannot open world store %q: %v\n", *dbPath, err)
		return exitFatal
	}
	defer func() { _ = st.Close() }()

	// stdout mode: Mint writes the raw token to stdout exactly once.
	// --out mode: out=nil, we write the token to the file ourselves at 0600.
	var out io.Writer
	if *outPath == "" {
		out = stdout
	}
	rawToken, credentialID, row, err := authority.Mint(context.Background(), st, *episode, grants, *ttl, now, out)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd session mint: %v\n", err)
		return exitFatal
	}

	if *outPath != "" {
		if err := os.WriteFile(*outPath, []byte(rawToken), 0o600); err != nil {
			fmt.Fprintf(stderr, "ailang-worldd session mint: write credential to %s: %v\n", *outPath, err)
			return exitFatal
		}
		// One-line confirmation, never the credential (D1/F8).
		fmt.Fprintf(stdout, "minted session credential for episode %s: %d grant(s), expires epoch %d, written to %s\n",
			*episode, len(grants), row.ExpiresAt, *outPath)
	} else {
		// The raw token was printed to stdout exactly once by Mint; end the line.
		io.WriteString(stdout, "\n")
	}

	// The stored credential_id (the hash), never the raw token, is the
	// operator-facing diagnostic the revoke verb consumes. D3/F9: only the hash
	// leaves the mint path.
	fmt.Fprintf(stderr, "ailang-worldd session mint: credential_id=%s (episode %s, %d grant(s), expires %d)\n",
		credentialID, *episode, len(grants), row.ExpiresAt)
	return exitOK
}

// runSessionRevoke implements `session revoke <credential_id-hash>` (D4): it
// DELETEs the mapping row so the next resolve of that credential is an UNKNOWN
// credential. No terminal fence is required — revocation is an ordinary
// operator command, not a one-human-act mint.
func runSessionRevoke(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("ailang-worldd session revoke", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dbPath := fs.String("db", "", "world store database (required)")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	rest := fs.Args()
	if len(rest) != 1 {
		fmt.Fprintln(stderr, "ailang-worldd session revoke: usage: session revoke <credential_id-hash> [--db <path>]")
		return exitUsage
	}
	if *dbPath == "" {
		fmt.Fprintln(stderr, "ailang-worldd session revoke: --db is required (matching serve --db)")
		return exitUsage
	}
	credID := rest[0]
	if !isHex64(credID) {
		fmt.Fprintf(stderr, "ailang-worldd session revoke: %q is not a 64-hex credential_id hash\n", credID)
		return exitUsage
	}

	st, err := store.Open(*dbPath)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd session revoke: cannot open world store %q: %v\n", *dbPath, err)
		return exitFatal
	}
	defer func() { _ = st.Close() }()

	if err := authority.Revoke(context.Background(), st, credID); err != nil {
		fmt.Fprintf(stderr, "ailang-worldd session revoke: %v\n", err)
		return exitFatal
	}
	fmt.Fprintf(stdout, "revoked session credential %s\n", credID)
	return exitOK
}

// parseGrants converts the repeatable --grant EFFECT=SCOPE:BUDGET specs into
// []broker.Capability, preserving order.
func parseGrants(specs []string) ([]broker.Capability, error) {
	grants := make([]broker.Capability, 0, len(specs))
	for _, spec := range specs {
		c, err := parseGrant(spec)
		if err != nil {
			return nil, err
		}
		grants = append(grants, c)
	}
	return grants, nil
}

// parseGrant parses one "EFFECT=SCOPE:BUDGET" grant (D1 unified form), e.g.
// "fs.read=/tmp/b:1". The budget is the last ':'-separated segment so a scope
// containing ':' is tolerated; the first '=' separates EFFECT.
func parseGrant(spec string) (broker.Capability, error) {
	eq := strings.Index(spec, "=")
	if eq < 0 || eq == len(spec)-1 {
		return broker.Capability{}, fmt.Errorf("invalid --grant %q: want EFFECT=SCOPE:BUDGET", spec)
	}
	effect, rest := spec[:eq], spec[eq+1:]
	colon := strings.LastIndex(rest, ":")
	if colon < 0 || colon == len(rest)-1 {
		return broker.Capability{}, fmt.Errorf("invalid --grant %q: want EFFECT=SCOPE:BUDGET", spec)
	}
	scope, budgetText := rest[:colon], rest[colon+1:]
	if effect == "" || scope == "" {
		return broker.Capability{}, fmt.Errorf("invalid --grant %q: EFFECT and SCOPE must be non-empty", spec)
	}
	budget, err := strconv.ParseInt(budgetText, 10, 64)
	if err != nil {
		return broker.Capability{}, fmt.Errorf("invalid --grant %q: budget %q is not an integer", spec, budgetText)
	}
	if budget < 0 {
		return broker.Capability{}, fmt.Errorf("invalid --grant %q: budget must be non-negative", spec)
	}
	return broker.Capability{Effect: effect, Scope: scope, Budget: budget}, nil
}

// readTermLine reads one line from the controlling terminal (the y/N answer).
func readTermLine(r io.Reader) string {
	line, _ := bufio.NewReader(r).ReadString('\n')
	return strings.TrimSpace(line)
}

// isYes reports whether a confirmation answer is affirmative (y/yes).
func isYes(answer string) bool {
	switch strings.ToLower(strings.TrimSpace(answer)) {
	case "y", "yes":
		return true
	default:
		return false
	}
}

// isHex64 reports whether s is exactly 64 hex characters.
func isHex64(s string) bool {
	if len(s) != 64 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}
