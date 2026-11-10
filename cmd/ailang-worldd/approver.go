package main

// Row 201 A1b (design_docs/planned/w-world-approvals-by-key-a1b.md §2.10, §2.10a; sprint-2 plan
// §2.1): `ailang-worldd approver enroll|list|revoke`. Enrolment and revocation are typed-only
// (parent R3): each asks y at /dev/tty within 120 s, with no store handle open while it waits.

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"time"
)

// ---- M3a ----

// approverConfirmBudget is the fixed y/N budget of approver enroll and revoke, counted from the
// moment the prompt is written (design §2.10a [r1], Q7: a constant, so a script cannot widen it).
const approverConfirmBudget = 120 * time.Second

// confirmApproverAct shows prompt at /dev/tty and returns exitOK only on a y typed within
// approverConfirmBudget, as a complete line with no read error (design §2.10a [r2]). Every other
// outcome returns its refusal code having written nothing. It is called with no store handle
// open, so no writer lock is held while it waits.
func confirmApproverAct(env sessionEnv, stderr io.Writer, prefix, prompt, nothing string) int {
	if env.after == nil { // N9a / K4a
		fmt.Fprintf(stderr, "%sinternal error: no prompt timer (sessionEnv.after)\n", prefix)
		return exitFatal
	}
	term, err := env.openTerminal()
	if err != nil { // N9 / K4
		fmt.Fprintf(stderr, "%srefusing: no controlling terminal (%v); %s. %s\n", prefix, err, nothing, sessionTTYFix)
		return exitUsage
	}
	defer func() { _ = term.Close() }() // closed on every return
	fmt.Fprint(term, prompt)            // the budget starts here
	line, rerr, ok := readTermLineErrWithin(term, approverConfirmBudget, env.after)
	if !ok { // N10a / K5a
		fmt.Fprintf(stderr, "%sno answer within %s; %s\n", prefix, approverConfirmBudget, nothing)
		return exitUsage
	}
	if rerr != nil || !strings.HasSuffix(line, "\n") || !isYes(line) { // N10 / K5: EOF, read error, partial line, not y
		fmt.Fprintf(stderr, "%saborted (not confirmed); %s. fix: run it again and type y at the prompt\n", prefix, nothing)
		return exitUsage
	}
	return exitOK
}

// readTermLineErrWithin is readTermLineWithin keeping the read's error and the RAW line
// (untrimmed, so a missing newline shows): ok is false when after(d) fires first. readTermLine
// drops the error and keeps a partial line (design R8), so the approver verbs use this instead.
// The one-shot rule applies: a timed-out read leaves its goroutine blocked until the process exits.
func readTermLineErrWithin(r io.Reader, d time.Duration, after func(time.Duration) <-chan time.Time) (line string, err error, ok bool) {
	type result struct {
		line string
		err  error
	}
	ch := make(chan result, 1)
	go func() {
		l, e := bufio.NewReader(r).ReadString('\n')
		ch <- result{l, e}
	}()
	select {
	case res := <-ch:
		return res.line, res.err, true
	case <-after(d):
		return "", nil, false
	}
}

// ---- M3b ----

// readRegularFileCapped reads path once: it opens non-blocking (so a FIFO cannot stall the
// open, plan P1), checks the open DESCRIPTOR (a regular file of at most max bytes), and reads
// through readCapped. Sprint 3's approve --signature reuses it.
func readRegularFileCapped(path string, max int64) ([]byte, error) {
	f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a regular file", path)
	}
	if info.Size() > max {
		return nil, fmt.Errorf("%s is %d bytes, over the %d-byte bound", path, info.Size(), max)
	}
	return readCapped(f, max)
}

// readCapped reads r to its end and refuses more than max bytes (a file that grew after its
// fstat).
func readCapped(r io.Reader, max int64) ([]byte, error) {
	b, err := io.ReadAll(io.LimitReader(r, max+1))
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > max {
		return nil, fmt.Errorf("grew past the %d-byte bound while being read", max)
	}
	return b, nil
}

// ---- M3c ----

const approverHelp = `usage: ailang-worldd approver enroll --name <name> --format ssh --pubkey <file> [--db <path>] [--credentials-db <path>]
       ailang-worldd approver list [--db <path>] [--credentials-db <path>]
       ailang-worldd approver revoke <name> [--db <path>] [--credentials-db <path>]

ATTENDED (enroll, revoke). An approver is a name and one SSH public key, kept in
the credentials store (row 201). Approving a mint request with its signature
arrives with 'approve --signature' (row 201 A1b, a later sprint); until then an
enrolled key approves nothing. enroll and revoke each ask y at /dev/tty and
wait at most 2m0s for it; an agent cannot run them.

enroll   reads --pubkey once (one ssh-ed25519 or ecdsa-sha2-nistp256 key line,
         at most 8 KiB; a private key file is refused and never shown), shows
         the key's SHA256 fingerprint, and stores it after y. A name and a key
         are each enrolled at most once, ever, revoked ones included.
list     prints every approver, read-only: name, format, key type,
         fingerprint, enrolled, last used, and active or revoked.
revoke   marks the approver revoked after y. The row is kept: it never
         approves again, and its name and key can never be enrolled again.

  --name <name>             1-48 of a-z 0-9 . _ - (the first a letter or digit);
                            its approvals record approved_by=ssh:<name>
  --format ssh              the only format until row 201 A2 (webauthn)
  --pubkey <file>           the .pub file (compare: ssh-keygen -l -f <file>)
  --db <path>               the world store (default ~/.ailang/world/world.db)
  --credentials-db <path>   default credentials.db beside --db
The global --addr is refused: approvers live in the credentials store.
`

// runApprover dispatches `approver enroll|list|revoke` (design §2.10: D1-D3).
func runApprover(args []string, stdout, stderr io.Writer, env sessionEnv) int {
	if len(args) == 0 { // D1
		fmt.Fprint(stderr, approverHelp)
		return exitUsage
	}
	switch args[0] {
	case "--help", "-help", "-h", "help": // D3
		fmt.Fprint(stdout, approverHelp)
		return exitOK
	default: // D2
		fmt.Fprintf(stderr, "ailang-worldd approver: unknown subcommand %q\n\n%s", args[0], approverHelp)
		return exitUsage
	}
}

// approverFlags parses one subcommand's flags (--db, --credentials-db, and the extra ones extra
// declares). done is true when it already answered: help (exitOK) or misuse (exitUsage).
func approverFlags(name string, args []string, stdout, stderr io.Writer, extra func(*flag.FlagSet)) (db, creds string, positional []string, code int, done bool) {
	fs := flag.NewFlagSet("ailang-worldd approver "+name, flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	dbFlag := fs.String("db", defaultWorldPaths().db, "the world store (default ~/.ailang/world/world.db)")
	credsFlag := fs.String("credentials-db", "", credentialsDBUsage)
	if extra != nil {
		extra(fs)
	}
	positional, err := parseInterleaved(fs, args)
	if errors.Is(err, flag.ErrHelp) {
		fmt.Fprint(stdout, approverHelp)
		return "", "", nil, exitOK, true
	}
	if err != nil {
		fmt.Fprint(stderr, approverHelp)
		return "", "", nil, exitUsage, true
	}
	return *dbFlag, *credsFlag, positional, 0, false
}

// approverRefuse prints prefix and the message to stderr and returns exitUsage.
func approverRefuse(stderr io.Writer, prefix, format string, a ...any) int {
	fmt.Fprintf(stderr, "%s"+format+"\n", append([]any{prefix}, a...)...)
	return exitUsage
}

// approverStores resolves the credentials store a subcommand acts on, refusing a missing world
// store (N3): "" and the refusal text, or the path and "".
func approverStores(db, creds string) (string, string) {
	if err := requireExistingStore(db); err != nil {
		return "", err.Error()
	}
	path, err := credentialsFor(db, creds)
	if err != nil {
		return "", err.Error()
	}
	return path, ""
}
