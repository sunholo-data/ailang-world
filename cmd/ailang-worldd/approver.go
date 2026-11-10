package main

// Row 201 A1b (design_docs/planned/w-world-approvals-by-key-a1b.md §2.10, §2.10a; sprint-2 plan
// §2.1): `ailang-worldd approver enroll|list|revoke`. Enrolment and revocation are typed-only
// (parent R3): each asks y at /dev/tty within 120 s, with no store handle open while it waits.

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
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
the credentials store (row 201). An enrolled, unrevoked key approves a mint
request with its signature: 'ailang-worldd approve <request-id> --signature
<file>' or '--sign-with <key>'. enroll and revoke each ask y at /dev/tty and wait at most 2m0s for it;
an agent cannot run them.

enroll   reads --pubkey once (one ssh-ed25519 or ecdsa-sha2-nistp256 key line,
         at most 8 KiB; a private key file is refused and never shown), shows
         the key's SHA256 fingerprint, and stores it after y. A name and a key
         are each enrolled at most once, ever, revoked ones included.
list     prints every approver, read-only: name, format, key type,
         fingerprint, enrolled, last used, and active or revoked.
revoke   marks the approver revoked after y. The row is kept: it never
         approves again, and its name and key can never be enrolled again.
         It then lists the sessions the approver approved (at most 64, newest
         first), each with its 'session revoke' line. It revokes none of them.

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
	case "list": // M3d
		return runApproverList(args[1:], stdout, stderr)
	case "enroll": // M3e
		return runApproverEnroll(args[1:], stdout, stderr, env)
	case "revoke": // M3g
		return runApproverRevoke(args[1:], stdout, stderr, env)
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

// ---- M3d ----

// runApproverList is `approver list` (design §2.10 L1-L3): read-only, no terminal, no lock.
func runApproverList(args []string, stdout, stderr io.Writer) int {
	const prefix = "ailang-worldd approver list: "
	db, creds, positional, code, done := approverFlags("list", args, stdout, stderr, nil)
	if done {
		return code
	}
	if len(positional) != 0 {
		fmt.Fprint(stderr, approverHelp)
		return exitUsage
	}
	credPath, refusal := approverStores(db, creds)
	if refusal != "" {
		return approverRefuse(stderr, prefix, "%s", refusal)
	}
	ctx, cancel := context.WithTimeout(context.Background(), approveReadBudget)
	defer cancel()
	ro, err := store.OpenCredentialsReadOnly(credPath)
	if err != nil {
		return approverRefuse(stderr, prefix, "cannot read the credentials store %s: %v", credPath, err)
	}
	rows, err := ro.ListApprovers(ctx)
	_ = ro.Close()
	if err != nil { // L1: below schema 5 names the serve restart
		return approverRefuse(stderr, prefix, "%v", err)
	}
	if len(rows) == 0 { // L2
		fmt.Fprintln(stdout, "no approvers")
		return exitOK
	}
	for _, r := range rows { // L3
		used := "never"
		if r.LastUsedAt != 0 {
			used = approverTime(r.LastUsedAt)
		}
		state := "active"
		if r.RevokedAt != 0 {
			state = "revoked " + approverTime(r.RevokedAt)
		}
		fmt.Fprintf(stdout, "%s  %s  %s  enrolled %s  last used %s  %s\n",
			r.Name, r.Format, approverKeyLabel(r.PublicKey), approverTime(r.EnrolledAt), used, state)
	}
	return exitOK
}

// approverKeyLabel is "<key type> <SHA256 fingerprint>" of a stored key blob; a blob the parser
// refuses (only a hand-edited row) shows as "invalid-key".
func approverKeyLabel(blob []byte) string {
	keyType := "invalid-key"
	if k, err := approval.KeyFromBlob(blob); err == nil {
		keyType = k.Type
	}
	return keyType + " " + approval.SSHPublicKey{Blob: blob}.Fingerprint()
}

// approverTime is a Unix time in RFC 3339, UTC.
func approverTime(v int64) string { return time.Unix(v, 0).UTC().Format(time.RFC3339) }

// ---- M3e ----

// approverEnrol is the enrol write, a test seam (plan A-D4): CA6 replaces it to inject a race or a
// failure; production never reassigns it.
var approverEnrol = (*store.Credentials).EnrolApprover

// runApproverEnroll is `approver enroll` (design §2.10 N1-N14, plan A-D6 order).
func runApproverEnroll(args []string, stdout, stderr io.Writer, env sessionEnv) int {
	const prefix = "ailang-worldd approver enroll: "
	var name, format, pubkey *string
	db, creds, positional, code, done := approverFlags("enroll", args, stdout, stderr, func(fs *flag.FlagSet) {
		name = fs.String("name", "", "the approver's name")
		format = fs.String("format", "", "ssh")
		pubkey = fs.String("pubkey", "", "the public key file (.pub)")
	})
	if done {
		return code
	}
	if len(positional) != 0 { // N1
		fmt.Fprint(stderr, approverHelp)
		return exitUsage
	}
	if *name == "" || *pubkey == "" { // N2
		return approverRefuse(stderr, prefix, "--name and --pubkey are required")
	}
	credPath, refusal := approverStores(db, creds) // N3
	if refusal != "" {
		return approverRefuse(stderr, prefix, "%s", refusal)
	}
	keyFile, err := readRegularFileCapped(*pubkey, approval.MaxKeyFileBytes) // N4: read once
	if err != nil {
		return approverRefuse(stderr, prefix, "--pubkey: %v", err)
	}
	a, err := approval.NewApprover(*name, *format, keyFile, env.now()) // N5
	if err != nil {
		return approverRefuse(stderr, prefix, "%v", err)
	}

	// Pre-checks on a read-only handle, closed before the prompt (N6-N8).
	pre, cancelPre := context.WithTimeout(context.Background(), approveReadBudget)
	ro, err := store.OpenCredentialsReadOnly(credPath)
	if err != nil {
		cancelPre()
		return approverRefuse(stderr, prefix, "cannot read the credentials store %s: %v", credPath, err)
	}
	refusal = approverEnrolPrecheck(pre, ro, a)
	_ = ro.Close()
	cancelPre()
	if refusal != "" {
		return approverRefuse(stderr, prefix, "%s", refusal)
	}

	// The fence: y within approverConfirmBudget, no store handle open (N9a, N9, N10a, N10).
	prompt := approval.RenderEnrolment(a) + fmt.Sprintf("Enrol approver %s? [y/N] ", a.Name)
	if code := confirmApproverAct(env, stderr, prefix, prompt, "nothing was enrolled"); code != exitOK {
		return code
	}

	// The write, under the credentials writer lock taken only now (N11-N13).
	ctx, cancel := context.WithTimeout(context.Background(), approveReadBudget)
	defer cancel()
	st, err := openCredentialsForSession(ctx, db, credPath)
	if err != nil { // N11
		fmt.Fprintf(stderr, "%scannot open the credentials store for writing: %v\n", prefix, err)
		return exitFatal
	}
	defer func() { _ = st.Close() }()
	err = approverEnrol(st, ctx, store.ApproverRow{Name: a.Name, Format: a.Format, PublicKey: a.Key.Blob, EnrolledAt: a.EnrolledAt})
	switch {
	case errors.Is(err, store.ErrApproverNameTaken), errors.Is(err, store.ErrApproverKeyEnrolled): // N12
		return approverRefuse(stderr, prefix, "%v (another enrol ran first); nothing was enrolled", err)
	case err != nil: // N13
		fmt.Fprintf(stderr, "%scannot enrol approver %q: %v\n", prefix, a.Name, err)
		return exitFatal
	}
	fmt.Fprintf(stdout, "enrolled approver %s: %s %s; approvals by it record approved_by=%s:%s\n", // N14
		a.Name, a.Key.Type, a.Key.Fingerprint(), a.Format, a.Name)
	return exitOK
}

// approverEnrolPrecheck is enroll's read-only refusal before any terminal: "" when the store is
// at schema 5 (N6) and neither the name (N7) nor the key (N8) is enrolled, revoked rows included.
func approverEnrolPrecheck(ctx context.Context, ro *store.Credentials, a approval.Approver) string {
	_, taken, err := ro.Approver(ctx, a.Name)
	if err != nil { // N6: below schema 5 names the serve restart
		return err.Error()
	}
	if taken { // N7
		return fmt.Sprintf("an approver named %q already exists (a name is never reused, even after revoke); fix: choose another --name", a.Name)
	}
	holder, enrolled, err := ro.ApproverByKey(ctx, a.Key.Blob)
	if err != nil {
		return err.Error()
	}
	if enrolled { // N8
		return fmt.Sprintf("key %s is already enrolled as approver %q (a key is never enrolled twice, even after revoke)", a.Key.Fingerprint(), holder.Name)
	}
	return ""
}

// ---- M3g ----

// approverRevoke is the revoke write, a test seam (plan A-D4): CA8 replaces it to inject a race or
// a failure; production never reassigns it.
var approverRevoke = (*store.Credentials).RevokeApprover

// runApproverRevoke is `approver revoke <name>` (design §2.10 K1-K7).
func runApproverRevoke(args []string, stdout, stderr io.Writer, env sessionEnv) int {
	const prefix = "ailang-worldd approver revoke: "
	db, creds, positional, code, done := approverFlags("revoke", args, stdout, stderr, nil)
	if done {
		return code
	}
	if len(positional) != 1 { // K1
		fmt.Fprint(stderr, approverHelp)
		return exitUsage
	}
	name := positional[0]
	credPath, refusal := approverStores(db, creds)
	if refusal != "" {
		return approverRefuse(stderr, prefix, "%s", refusal)
	}

	// The read-only lookup, closed before the prompt (K2, K3).
	pre, cancelPre := context.WithTimeout(context.Background(), approveReadBudget)
	ro, err := store.OpenCredentialsReadOnly(credPath)
	if err != nil {
		cancelPre()
		return approverRefuse(stderr, prefix, "cannot read the credentials store %s: %v", credPath, err)
	}
	row, ok, err := ro.Approver(pre, name)
	_ = ro.Close()
	cancelPre()
	if err != nil { // below schema 5 names the serve restart
		return approverRefuse(stderr, prefix, "%v", err)
	}
	if !ok { // K2
		return approverRefuse(stderr, prefix, "no approver named %q in %s; fix: ailang-worldd approver list shows them", name, credPath)
	}
	if err := approval.CheckRevoke(approval.Approver{EnrolledAt: row.EnrolledAt, RevokedAt: row.RevokedAt}, env.now()); err != nil { // K3
		return approverRefuse(stderr, prefix, "approver %q cannot be revoked: %v", name, err)
	}

	// The fence (K4a, K4, K5a, K5).
	prompt := fmt.Sprintf("Revoke approver %s (%s)? It will never approve again, and its name and key can never be enrolled again. [y/N] ",
		name, approverKeyLabel(row.PublicKey))
	if code := confirmApproverAct(env, stderr, prefix, prompt, "nothing was revoked"); code != exitOK {
		return code
	}

	// The write (K5b, K6, K7). Its root and its writer end before K8 reads.
	ctx, cancel := context.WithTimeout(context.Background(), approveReadBudget)
	st, err := openCredentialsForSession(ctx, db, credPath)
	if err != nil { // K5b
		cancel()
		fmt.Fprintf(stderr, "%scannot open the credentials store for writing: %v\n", prefix, err)
		return exitFatal
	}
	err = approverRevoke(st, ctx, name, env.now())
	_ = st.Close()
	cancel()
	switch {
	case errors.Is(err, store.ErrApproverRevoked), errors.Is(err, store.ErrApproverUnknown): // K6
		return approverRefuse(stderr, prefix, "%v (another revoke ran first); nothing was revoked", err)
	case err != nil:
		fmt.Fprintf(stderr, "%scannot revoke approver %q: %v\n", prefix, name, err)
		return exitFatal
	}
	fmt.Fprintf(stdout, "revoked approver %s\n", name) // K7

	// K8, K9 (design §2.10 [r2]): what it approved, under its own read-only root. Q4: it revokes
	// no session; the operator decides.
	list, cancelList := context.WithTimeout(context.Background(), approveReadBudget)
	defer cancelList()
	ro, err = store.OpenCredentialsReadOnly(credPath)
	if err != nil {
		fmt.Fprintf(stderr, "%srevoked, but cannot read the credentials store %s to list what it approved: %v; fix: ailang-worldd approve --list --decided\n", prefix, credPath, err)
		return exitFatal
	}
	rows, err := ro.DecidedMintRequests(list, row.Format+":"+name, 0)
	_ = ro.Close()
	if err != nil {
		fmt.Fprintf(stderr, "%srevoked, but cannot list what it approved: %v; fix: ailang-worldd approve --list --decided\n", prefix, err)
		return exitFatal
	}
	if len(rows) == 0 { // K9
		fmt.Fprintln(stdout, "no session was approved by it")
		return exitOK
	}
	fmt.Fprintf(stdout, "it approved %d session(s), newest first; each stays valid until it expires or is revoked:\n", len(rows)) // K8
	for _, r := range rows {
		fmt.Fprintf(stdout, "  %s  episode %s  at %s  credential %s\n    ailang-worldd session revoke --db %s --credentials-db %s %s\n",
			r.RequestID, r.EpisodeID, approverTime(r.ConsumedAt), credentialPrefix(r.CredentialID), shellQuote(db), shellQuote(credPath), r.CredentialID)
	}
	return exitOK
}
