package main

// Row 201 A1b sprint 2 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint2-plan.md):
// the approver verbs' bounded y/N, capped file read, dispatch and list.

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/broker"
)

// ---- M3a ----

// lockedTerm is a terminal safe to poll while the verb writes to it (plan A-D7): reads come from
// in; writes and Close are recorded under mu.
type lockedTerm struct {
	mu     sync.Mutex
	in     io.Reader
	out    bytes.Buffer
	closed bool
}

func (t *lockedTerm) Read(p []byte) (int, error) { return t.in.Read(p) }
func (t *lockedTerm) Write(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.out.Write(p)
}
func (t *lockedTerm) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.closed = true
	return nil
}

// Output is everything written to the terminal so far.
func (t *lockedTerm) Output() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.out.String()
}

// Closed reports whether Close was called.
func (t *lockedTerm) Closed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.closed
}

// partialReader returns data on its first Read and err on every later one: a terminal that yields
// a partial answer and then fails or ends.
type partialReader struct {
	data string
	err  error
	sent bool
}

func (r *partialReader) Read(p []byte) (int, error) {
	if r.sent {
		return 0, r.err
	}
	r.sent = true
	return copy(p, r.data), nil
}

var errInjectedRead = errors.New("injected read failure")

// neverFires is a timer that never fires.
func neverFires(time.Duration) <-chan time.Time { return make(chan time.Time) }

// RL1: readTermLineErrWithin returns the raw line with its read error, or ok=false on the timer.
func TestReadTermLineErrWithin(t *testing.T) {
	for _, c := range []struct {
		name string
		r    io.Reader
		line string
		err  error
	}{
		{"y-newline", strings.NewReader("y\n"), "y\n", nil},
		{"raw-untrimmed", strings.NewReader(" y \n"), " y \n", nil},
		{"y-eof", strings.NewReader("y"), "y", io.EOF},
		{"y-error", &partialReader{data: "y", err: errInjectedRead}, "y", errInjectedRead},
		{"empty-eof", strings.NewReader(""), "", io.EOF},
	} {
		line, err, ok := readTermLineErrWithin(c.r, time.Second, neverFires)
		if !ok || line != c.line || !errors.Is(err, c.err) || (c.err == nil) != (err == nil) {
			t.Errorf("%s: (%q, %v, %v), want (%q, %v, true)", c.name, line, err, ok, c.line, c.err)
		}
	}
	pr, pw := io.Pipe()
	t.Cleanup(func() { _ = pw.Close() })
	fired := func(time.Duration) <-chan time.Time { c := make(chan time.Time); close(c); return c }
	type result struct {
		line string
		err  error
		ok   bool
	}
	done := make(chan result, 1)
	go func() {
		line, err, ok := readTermLineErrWithin(pr, time.Second, fired)
		done <- result{line, err, ok}
	}()
	select {
	case got := <-done:
		if got.ok || got.line != "" || got.err != nil {
			t.Fatalf("timer-first: (%q, %v, %v), want (\"\", nil, false)", got.line, got.err, got.ok)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timer-first: readTermLineErrWithin did not return on the timer")
	}
}

// CF1: confirmApproverAct confirms only a complete "y" line within the budget; every other
// outcome refuses with its own text, and the terminal is closed whenever it was opened.
func TestConfirmApproverAct(t *testing.T) {
	const aborted = "pfx: aborted (not confirmed); nothing was done. fix: run it again and type y at the prompt\n"
	type arm struct {
		name       string
		in         func(t *testing.T) io.Reader // nil: no terminal
		fire       bool                         // the timer fires at once
		nilAfter   bool
		wantCode   int
		wantStderr string // exact (when contains is nil)
		contains   []string
	}
	arms := []arm{
		{name: "y", in: func(*testing.T) io.Reader { return strings.NewReader("y\n") }, wantCode: exitOK},
		{name: "yes", in: func(*testing.T) io.Reader { return strings.NewReader("yes\n") }, wantCode: exitOK},
		{name: "n", in: func(*testing.T) io.Reader { return strings.NewReader("n\n") }, wantCode: exitUsage, wantStderr: aborted},
		{name: "empty-line", in: func(*testing.T) io.Reader { return strings.NewReader("\n") }, wantCode: exitUsage, wantStderr: aborted},
		{name: "eof", in: func(*testing.T) io.Reader { return strings.NewReader("") }, wantCode: exitUsage, wantStderr: aborted},
		{name: "read-error", in: func(*testing.T) io.Reader { return &partialReader{err: errInjectedRead} }, wantCode: exitUsage, wantStderr: aborted},
		{name: "y-then-eof", in: func(*testing.T) io.Reader { return &partialReader{data: "y", err: io.EOF} }, wantCode: exitUsage, wantStderr: aborted},
		{name: "y-then-error", in: func(*testing.T) io.Reader { return &partialReader{data: "y", err: errInjectedRead} }, wantCode: exitUsage, wantStderr: aborted},
		{name: "timeout", in: func(t *testing.T) io.Reader {
			pr, pw := io.Pipe()
			t.Cleanup(func() { _ = pw.Close() })
			return pr
		}, fire: true, wantCode: exitUsage, wantStderr: "pfx: no answer within 2m0s; nothing was done\n"},
		{name: "no-terminal", wantCode: exitUsage, contains: []string{"pfx: refusing: no controlling terminal", "nothing was done", sessionTTYFix}},
		{name: "nil-after", in: func(*testing.T) io.Reader { return strings.NewReader("y\n") }, nilAfter: true, wantCode: exitFatal,
			wantStderr: "pfx: internal error: no prompt timer (sessionEnv.after)\n"},
	}
	for _, a := range arms {
		t.Run(a.name, func(t *testing.T) {
			var term *lockedTerm
			opened := false
			var mu sync.Mutex
			var waited time.Duration
			env := sessionEnv{now: atTime(approveT0), after: func(d time.Duration) <-chan time.Time {
				mu.Lock()
				waited = d
				mu.Unlock()
				if a.fire {
					c := make(chan time.Time)
					close(c)
					return c
				}
				return neverFires(d)
			}}
			env.openTerminal = func() (io.ReadWriteCloser, error) {
				opened = true
				if a.in == nil {
					return nil, os.ErrNotExist
				}
				term = &lockedTerm{in: a.in(t)}
				return term, nil
			}
			if a.nilAfter {
				env.after = nil
			}
			var stderr bytes.Buffer
			done := make(chan int, 1)
			go func() { done <- confirmApproverAct(env, &stderr, "pfx: ", "PROMPT? ", "nothing was done") }()
			var code int
			select {
			case code = <-done:
			case <-time.After(5 * time.Second):
				t.Fatal("confirmApproverAct did not return within 5 s")
			}
			if code != a.wantCode {
				t.Fatalf("code %d, want %d (stderr %q)", code, a.wantCode, stderr.String())
			}
			if a.wantStderr != "" || a.contains == nil {
				if stderr.String() != a.wantStderr {
					t.Errorf("stderr %q, want %q", stderr.String(), a.wantStderr)
				}
			}
			for _, w := range a.contains {
				if !strings.Contains(stderr.String(), w) {
					t.Errorf("stderr %q lacks %q", stderr.String(), w)
				}
			}
			switch {
			case a.nilAfter:
				if opened {
					t.Error("the terminal was opened with no prompt timer")
				}
			case term != nil:
				if term.Output() != "PROMPT? " {
					t.Errorf("terminal got %q, want the prompt exactly", term.Output())
				}
				if !term.Closed() {
					t.Error("the terminal was not closed")
				}
				mu.Lock()
				d := waited
				mu.Unlock()
				if d != 120*time.Second {
					t.Errorf("the prompt was bounded by %v, want 2m0s", d)
				}
			}
		})
	}
}

// ---- M3b ----

// FC1: readRegularFileCapped reads a regular file of at most max bytes once, and refuses
// anything else by the open descriptor, without blocking on a FIFO.
func TestReadRegularFileCapped(t *testing.T) {
	dir := t.TempDir()
	write := func(name string, n int) string {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, bytes.Repeat([]byte("k"), n), 0o600); err != nil {
			t.Fatal(err)
		}
		return p
	}
	if b, err := readRegularFileCapped(write("at-max", 8192), 8192); err != nil || len(b) != 8192 {
		t.Errorf("an 8192-byte file at max 8192: %d bytes, %v; want all of it", len(b), err)
	}
	if b, err := readRegularFileCapped(write("small", 3), 8192); err != nil || string(b) != "kkk" {
		t.Errorf("a 3-byte file: %q, %v", b, err)
	}
	over := write("over", 8193)
	if _, err := readRegularFileCapped(over, 8192); err == nil || err.Error() != over+" is 8193 bytes, over the 8192-byte bound" {
		t.Errorf("an 8193-byte file: %v, want the fstat size refusal", err)
	}
	if _, err := readRegularFileCapped(filepath.Join(dir, "missing"), 8192); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("a missing file: %v, want os.ErrNotExist", err)
	}
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := readRegularFileCapped(sub, 8192); err == nil || err.Error() != sub+" is not a regular file" {
		t.Errorf("a directory: %v, want the not-regular refusal", err)
	}
	fifo := filepath.Join(dir, "fifo")
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Fatal(err)
	}
	// A mutant that blocks in open(2) is released when the test ends: a writer that opens the
	// FIFO non-blocking succeeds only while a reader is waiting.
	t.Cleanup(func() {
		if w, err := os.OpenFile(fifo, os.O_WRONLY|syscall.O_NONBLOCK, 0); err == nil {
			_ = w.Close()
		}
	})
	done := make(chan error, 1)
	go func() { _, err := readRegularFileCapped(fifo, 8192); done <- err }()
	select {
	case err := <-done:
		if err == nil || err.Error() != fifo+" is not a regular file" {
			t.Errorf("a FIFO: %v, want the not-regular refusal", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("a FIFO: readRegularFileCapped blocked for 5 s")
	}
}

// FC2: readCapped refuses a stream that is longer than max (the file grew after its fstat).
func TestReadCappedRefusesGrowth(t *testing.T) {
	if b, err := readCapped(strings.NewReader(strings.Repeat("x", 10)), 10); err != nil || len(b) != 10 {
		t.Errorf("10 bytes at max 10: %d bytes, %v", len(b), err)
	}
	if _, err := readCapped(strings.NewReader(strings.Repeat("x", 11)), 10); err == nil || err.Error() != "grew past the 10-byte bound while being read" {
		t.Errorf("11 bytes at max 10: %v, want the growth refusal", err)
	}
}

// CA1 (sprint-2 plan §3 M3c): the verb's help, dispatch and the --addr refusal.
func TestApproverDispatchAndHelp(t *testing.T) {
	t.Setenv(broker.RegistryCredentialVariable, "")
	var out, errw bytes.Buffer
	if code := runApprover(nil, &out, &errw, sessionEnv{}); code != exitUsage {
		t.Fatalf("no args: code %d, want %d", code, exitUsage)
	}
	if errw.String() != approverHelp {
		t.Errorf("no args: stderr is not approverHelp:\n%s", errw.String())
	}
	if out.Len() != 0 {
		t.Errorf("no args: stdout %q, want empty", out.String())
	}

	out.Reset()
	errw.Reset()
	if code := runApprover([]string{"frob"}, &out, &errw, sessionEnv{}); code != exitUsage {
		t.Fatalf("frob: code %d, want %d", code, exitUsage)
	}
	if want := "ailang-worldd approver: unknown subcommand \"frob\"\n\n" + approverHelp; errw.String() != want {
		t.Errorf("frob: stderr %q, want %q", errw.String(), want)
	}

	for _, arg := range []string{"--help", "-h", "-help", "help"} {
		out.Reset()
		errw.Reset()
		if code := runApprover([]string{arg}, &out, &errw, sessionEnv{}); code != exitOK {
			t.Fatalf("%s: code %d, want %d", arg, code, exitOK)
		}
		if out.String() != approverHelp {
			t.Errorf("%s: stdout is not approverHelp", arg)
		}
		if errw.Len() != 0 {
			t.Errorf("%s: stderr %q, want empty", arg, errw.String())
		}
	}

	out.Reset()
	errw.Reset()
	if code := run([]string{"--addr", "http://127.0.0.1:1", "approver", "list"}, &out, &errw); code != exitUsage {
		t.Fatalf("--addr: code %d, want %d", code, exitUsage)
	}
	const wantAddr = "ailang-worldd: --addr is a client flag and is not valid for 'approver'; approvers live in the credentials store beside --db, not on a running daemon\n"
	if errw.String() != wantAddr {
		t.Errorf("--addr: stderr %q, want %q", errw.String(), wantAddr)
	}

	out.Reset()
	errw.Reset()
	if code := run([]string{"approver", "--help"}, &out, &errw); code != exitOK {
		t.Fatalf("approver --help: code %d, want %d", code, exitOK)
	}
	if out.String() != approverHelp {
		t.Errorf("approver --help: stdout is not approverHelp")
	}
}

// CA7 (sprint-2 plan §3 M3d): `approver list` is read-only and prints every row, empty, v4 and misuse
// arms included.
func TestApproverList(t *testing.T) {
	t.Setenv(broker.RegistryCredentialVariable, "")

	// populated: an active row, a revoked row, and a row whose stored blob the parser refuses.
	db, creds := approverStore(t)
	seedApprover(t, creds, "a", fixtureKeyOf(t, "ed25519.pub"), approveT0)
	seedApprover(t, creds, "b", fixtureKeyOf(t, "p256.pub"), approveT0)
	revokeSeeded(t, creds, "b", approveT0+200)
	raw := rawCredentials(t, creds)
	if _, err := raw.Exec(`UPDATE approvers SET last_used_at = 1800000300 WHERE name = 'a'`); err != nil {
		t.Fatal(err)
	}
	if _, err := raw.Exec(`INSERT INTO approvers (name, format, public_key, enrolled_at) VALUES ('c', 'ssh', x'00000001', 1800000000)`); err != nil {
		t.Fatal(err)
	}
	_ = raw.Close()

	var out, errw bytes.Buffer
	if code := runApprover([]string{"list", "--db", db}, &out, &errw, sessionEnv{}); code != exitOK {
		t.Fatalf("populated: code %d, want %d; stderr %q", code, exitOK, errw.String())
	}
	want := "a  ssh  ssh-ed25519 " + edFingerprint + "  enrolled 2027-01-15T08:00:00Z  last used 2027-01-15T08:05:00Z  active\n" +
		"b  ssh  ecdsa-sha2-nistp256 " + p256Fingerprint + "  enrolled 2027-01-15T08:00:00Z  last used never  revoked 2027-01-15T08:03:20Z\n" +
		"c  ssh  invalid-key " + approval.SSHPublicKey{Blob: []byte{0, 0, 0, 1}}.Fingerprint() + "  enrolled 2027-01-15T08:00:00Z  last used never  active\n"
	if out.String() != want {
		t.Errorf("populated stdout:\n%q\nwant:\n%q", out.String(), want)
	}

	// empty (L2)
	dbEmpty, _ := approverStore(t)
	out.Reset()
	errw.Reset()
	if code := runApprover([]string{"list", "--db", dbEmpty}, &out, &errw, sessionEnv{}); code != exitOK {
		t.Fatalf("empty: code %d, want %d", code, exitOK)
	}
	if out.String() != "no approvers\n" {
		t.Errorf("empty stdout %q, want %q", out.String(), "no approvers\n")
	}

	// v4 store (L1)
	dbV4, _ := approverStoreV4(t)
	out.Reset()
	errw.Reset()
	if code := runApprover([]string{"list", "--db", dbV4}, &out, &errw, sessionEnv{}); code != exitUsage {
		t.Fatalf("v4: code %d, want %d", code, exitUsage)
	}
	if !strings.HasPrefix(errw.String(), "ailang-worldd approver list: ") ||
		!strings.Contains(errw.String(), "approvers need schema 5") ||
		!strings.Contains(errw.String(), "restart `ailang-worldd serve` ") {
		t.Errorf("v4 stderr %q lacks the list prefix, the schema 5 need or the serve restart", errw.String())
	}

	// help
	out.Reset()
	errw.Reset()
	if code := runApprover([]string{"list", "--help"}, &out, &errw, sessionEnv{}); code != exitOK {
		t.Fatalf("list --help: code %d, want %d", code, exitOK)
	}
	if out.String() != approverHelp {
		t.Errorf("list --help: stdout is not approverHelp")
	}

	// misuse: an unknown flag and a stray positional
	out.Reset()
	errw.Reset()
	if code := runApprover([]string{"list", "--bogus", "--db", db}, &out, &errw, sessionEnv{}); code != exitUsage {
		t.Fatalf("--bogus: code %d, want %d", code, exitUsage)
	}
	if !strings.Contains(errw.String(), "usage: ailang-worldd approver") {
		t.Errorf("--bogus stderr %q lacks the usage line", errw.String())
	}
	out.Reset()
	errw.Reset()
	if code := runApprover([]string{"list", "extra", "--db", db}, &out, &errw, sessionEnv{}); code != exitUsage {
		t.Fatalf("positional: code %d, want %d", code, exitUsage)
	}
	if !strings.Contains(errw.String(), "usage: ailang-worldd approver") {
		t.Errorf("positional stderr %q lacks the usage line", errw.String())
	}

	// missing world store
	out.Reset()
	errw.Reset()
	absent := filepath.Join(t.TempDir(), "absent.db")
	if code := runApprover([]string{"list", "--db", absent}, &out, &errw, sessionEnv{}); code != exitUsage {
		t.Fatalf("absent world store: code %d, want %d", code, exitUsage)
	}
	if !strings.Contains(errw.String(), "does not exist; refusing to create one") {
		t.Errorf("absent world store stderr %q", errw.String())
	}

	// missing credentials store
	out.Reset()
	errw.Reset()
	dbNoCreds := dbPathOf(t)
	if code := runApprover([]string{"list", "--db", dbNoCreds}, &out, &errw, sessionEnv{}); code != exitUsage {
		t.Fatalf("no credentials store: code %d, want %d", code, exitUsage)
	}
	if !strings.Contains(errw.String(), "cannot read the credentials store") {
		t.Errorf("no credentials store stderr %q", errw.String())
	}
}
