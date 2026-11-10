package main

// Row 201 A1b sprint 2 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint2-plan.md):
// the approver verbs' bounded y/N, capped file read, dispatch and list.

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
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
