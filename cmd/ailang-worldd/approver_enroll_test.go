package main

// Row 201 A1b sprint 2 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint2-plan.md §3
// M3e, M3f): `approver enroll` refusals, the stored key, the declined arms and the races. Every
// store is under t.TempDir(); every key is a checked-in fixture (D-WORLD-109).

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
)

// hookTerm is a terminal that runs hook once, on its first Read, then reads from in: CA2 rewrites
// the key file there, after the prompt was shown and before the answer is read.
type hookTerm struct {
	lockedTerm
	hook   func()
	hooked bool
}

func (t *hookTerm) Read(p []byte) (int, error) {
	if !t.hooked {
		t.hooked = true
		t.hook()
	}
	return t.lockedTerm.Read(p)
}

// CA2: the key enrolled is the key shown, even when the file changes under the prompt.
func TestApproverEnrollStoresWhatItShowed(t *testing.T) {
	db, creds := approverStore(t)
	pub := filepath.Join(t.TempDir(), "k.pub")
	edBytes, err := os.ReadFile(fixturePub("ed25519.pub"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(pub, edBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	p256Bytes, err := os.ReadFile(fixturePub("p256.pub"))
	if err != nil {
		t.Fatal(err)
	}
	var opened bool
	term := &hookTerm{
		lockedTerm: lockedTerm{in: strings.NewReader("y\n")},
		hook: func() {
			if err := os.WriteFile(pub, p256Bytes, 0o600); err != nil {
				t.Error(err)
			}
		},
	}
	var out, errw bytes.Buffer
	code := runApprover(enrollArgs(db, "mark-mac", pub), &out, &errw, approveEnv(term, &opened, atTime(approveT0)))
	if code != exitOK {
		t.Fatalf("enrol: code %d, want exitOK; stderr %q", code, errw.String())
	}
	if want := "enrolled approver mark-mac: ssh-ed25519 " + edFingerprint + "; approvals by it record approved_by=ssh:mark-mac\n"; out.String() != want {
		t.Fatalf("stdout %q, want %q", out.String(), want)
	}
	shown := term.Output()
	if !strings.HasSuffix(shown, "Enrol approver mark-mac? [y/N] ") {
		t.Fatalf("the terminal output does not end with the prompt: %q", shown)
	}
	for _, want := range []string{edFingerprint, "records   approved_by=ssh:mark-mac", approval.EnrolmentWarning} {
		if !strings.Contains(shown, want) {
			t.Fatalf("the terminal output lacks %q:\n%s", want, shown)
		}
	}
	if !term.Closed() {
		t.Fatal("the terminal was not closed")
	}
	rows := approverRows(t, creds)
	want := []store.ApproverRow{{Name: "mark-mac", Format: "ssh", PublicKey: fixtureKeyOf(t, "ed25519.pub").Blob, EnrolledAt: approveT0}}
	if !reflect.DeepEqual(rows, want) {
		t.Fatalf("approver rows %+v, want %+v (the first key, not the p256 now on disk)", rows, want)
	}
	now, err := os.ReadFile(pub)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(now, p256Bytes) {
		t.Fatal("control: the hook did not rewrite the key file")
	}
}

// CA3: with no controlling terminal, enrol refuses before anything is written.
func TestApproverEnrollWithoutATerminalEnrolsNothing(t *testing.T) {
	db, creds := approverStore(t)
	env := noTerminalEnv(approveT0)
	env.after = time.After // P9: the no-terminal arm needs the timer seam too
	var out, errw bytes.Buffer
	code := runApprover(enrollArgs(db, "mark-mac", fixturePub("ed25519.pub")), &out, &errw, env)
	if code != exitUsage {
		t.Fatalf("code %d, want exitUsage; stderr %q", code, errw.String())
	}
	const prefix = "ailang-worldd approver enroll: refusing: no controlling terminal ("
	if !strings.HasPrefix(errw.String(), prefix) {
		t.Fatalf("stderr %q does not start with %q", errw.String(), prefix)
	}
	for _, want := range []string{"nothing was enrolled", sessionTTYFix} {
		if !strings.Contains(errw.String(), want) {
			t.Fatalf("stderr %q lacks %q", errw.String(), want)
		}
	}
	if out.Len() != 0 {
		t.Fatalf("stdout %q, want nothing", out.String())
	}
	if rows := approverRows(t, creds); len(rows) != 0 {
		t.Fatalf("%d approver rows, want 0", len(rows))
	}
}

// CA4: a declined prompt (n, or y cut short by end of input) enrols nothing and closes the terminal.
func TestApproverEnrollDeclinedEnrolsNothing(t *testing.T) {
	const want = "ailang-worldd approver enroll: aborted (not confirmed); nothing was enrolled. fix: run it again and type y at the prompt\n"
	arms := []struct {
		name string
		in   io.Reader
	}{
		{name: "n", in: strings.NewReader("n\n")},
		{name: "y-then-eof", in: &partialReader{data: "y", err: io.EOF}},
	}
	for _, arm := range arms {
		t.Run(arm.name, func(t *testing.T) {
			db, creds := approverStore(t)
			var opened bool
			term := &lockedTerm{in: arm.in}
			var out, errw bytes.Buffer
			code := runApprover(enrollArgs(db, "mark-mac", fixturePub("ed25519.pub")), &out, &errw, approveEnv(term, &opened, atTime(approveT0)))
			if code != exitUsage {
				t.Fatalf("code %d, want exitUsage; stderr %q", code, errw.String())
			}
			if errw.String() != want {
				t.Fatalf("stderr %q, want %q", errw.String(), want)
			}
			if !term.Closed() {
				t.Fatal("the terminal was not closed")
			}
			if rows := approverRows(t, creds); len(rows) != 0 {
				t.Fatalf("%d approver rows, want 0", len(rows))
			}
		})
	}
}

// CA5: every refusal before the prompt is exitUsage, names its reason, opens no terminal, and
// writes nothing beyond the arm's own seeds. Each arm runs under a 5 s guard, so a blocking open
// or a nil terminal fails the arm rather than the run.
func TestApproverEnrollRefusals(t *testing.T) {
	std := func(t *testing.T) (string, string) { return approverStore(t) }
	ed := fixturePub("ed25519.pub")
	noise := func() string {
		s := sha512.Sum512([]byte("cli-noise"))
		return base64.StdEncoding.EncodeToString(bytes.Repeat(s[:], 4))
	}()
	priv := "-----BEGIN OPENSSH PRIVATE KEY-----\n" + noise + "\n-----END OPENSSH PRIVATE KEY-----\n"
	arms := []struct {
		name      string
		setup     func(t *testing.T) (db, creds string)
		args      func(t *testing.T, db string) []string
		seed      func(t *testing.T, creds string)
		want      []string
		forbid    string
		checkRows bool
		wantRows  int
	}{
		{name: "flag", setup: std, want: []string{"usage: ailang-worldd approver"}, checkRows: true,
			args: func(t *testing.T, db string) []string {
				return append([]string{"enroll", "--bogus"}, enrollArgs(db, "mark-mac", ed)[1:]...)
			}},
		{name: "positional", setup: std, want: []string{"usage: ailang-worldd approver"}, checkRows: true,
			args: func(t *testing.T, db string) []string { return append(enrollArgs(db, "mark-mac", ed), "extra") }},
		{name: "no-name", setup: std, want: []string{"--name and --pubkey are required"}, checkRows: true,
			args: func(t *testing.T, db string) []string {
				return []string{"enroll", "--format", "ssh", "--pubkey", ed, "--db", db}
			}},
		{name: "no-pubkey", setup: std, want: []string{"--name and --pubkey are required"}, checkRows: true,
			args: func(t *testing.T, db string) []string {
				return []string{"enroll", "--name", "mark-mac", "--format", "ssh", "--db", db}
			}},
		{name: "no-world-store", want: []string{"does not exist; refusing to create one"},
			setup: func(t *testing.T) (string, string) { return filepath.Join(t.TempDir(), "absent.db"), "" },
			args:  func(t *testing.T, db string) []string { return enrollArgs(db, "mark-mac", ed) }},
		{name: "fifo", setup: std, want: []string{"is not a regular file"}, checkRows: true,
			args: func(t *testing.T, db string) []string {
				fifo := filepath.Join(t.TempDir(), "fifo")
				if err := syscall.Mkfifo(fifo, 0o600); err != nil {
					t.Fatal(err)
				}
				// The FC1 writer: a mutant that blocks in open(2) is released when the test ends.
				t.Cleanup(func() {
					if w, err := os.OpenFile(fifo, os.O_WRONLY|syscall.O_NONBLOCK, 0); err == nil {
						_ = w.Close()
					}
				})
				return enrollArgs(db, "mark-mac", fifo)
			}},
		{name: "oversize", setup: std, want: []string{"is 8193 bytes, over the 8192-byte bound"}, checkRows: true,
			args: func(t *testing.T, db string) []string {
				p := filepath.Join(t.TempDir(), "big.pub")
				if err := os.WriteFile(p, bytes.Repeat([]byte("k"), 8193), 0o600); err != nil {
					t.Fatal(err)
				}
				return enrollArgs(db, "mark-mac", p)
			}},
		{name: "private-key", setup: std, want: []string{"that is a PRIVATE key file"}, forbid: noise[:24], checkRows: true,
			args: func(t *testing.T, db string) []string {
				p := filepath.Join(t.TempDir(), "id.pub")
				if err := os.WriteFile(p, []byte(priv), 0o600); err != nil {
					t.Fatal(err)
				}
				return enrollArgs(db, "mark-mac", p)
			}},
		{name: "webauthn", setup: std, want: []string{"row 201 A2"}, checkRows: true,
			args: func(t *testing.T, db string) []string {
				return []string{"enroll", "--name", "mark-mac", "--format", "webauthn", "--pubkey", ed, "--db", db}
			}},
		{name: "no-format", setup: std, want: []string{"approver format must be ssh"}, checkRows: true,
			args: func(t *testing.T, db string) []string {
				return []string{"enroll", "--name", "mark-mac", "--pubkey", ed, "--db", db}
			}},
		{name: "bad-name", setup: std, want: []string{"approver name must match"}, checkRows: true,
			args: func(t *testing.T, db string) []string { return enrollArgs(db, "Mark", ed) }},
		{name: "no-credentials-store", want: []string{"cannot read the credentials store"},
			setup: func(t *testing.T) (string, string) { return dbPathOf(t), "" },
			args:  func(t *testing.T, db string) []string { return enrollArgs(db, "mark-mac", ed) }},
		{name: "v4", want: []string{"approvers need schema 5", "restart `ailang-worldd serve`"},
			setup: func(t *testing.T) (string, string) { return approverStoreV4(t) },
			args:  func(t *testing.T, db string) []string { return enrollArgs(db, "mark-mac", ed) }},
		{name: "name-taken", setup: std, checkRows: true, wantRows: 1,
			seed: func(t *testing.T, creds string) {
				seedApprover(t, creds, "mark-mac", fixtureKeyOf(t, "p256.pub"), approveT0-10)
			},
			args: func(t *testing.T, db string) []string { return enrollArgs(db, "mark-mac", ed) },
			want: []string{`an approver named "mark-mac" already exists (a name is never reused, even after revoke); fix: choose another --name`}},
		{name: "name-taken-revoked", setup: std, checkRows: true, wantRows: 1,
			seed: func(t *testing.T, creds string) {
				seedApprover(t, creds, "mark-mac", fixtureKeyOf(t, "p256.pub"), approveT0-10)
				revokeSeeded(t, creds, "mark-mac", approveT0-5)
			},
			args: func(t *testing.T, db string) []string { return enrollArgs(db, "mark-mac", ed) },
			want: []string{`an approver named "mark-mac" already exists (a name is never reused, even after revoke); fix: choose another --name`}},
		{name: "key-enrolled", setup: std, checkRows: true, wantRows: 1,
			seed: func(t *testing.T, creds string) {
				seedApprover(t, creds, "other", fixtureKeyOf(t, "ed25519.pub"), approveT0-10)
			},
			args: func(t *testing.T, db string) []string { return enrollArgs(db, "mark-mac", ed) },
			want: []string{edFingerprint + ` is already enrolled as approver "other" (a key is never enrolled twice, even after revoke)`}},
		{name: "key-enrolled-revoked", setup: std, checkRows: true, wantRows: 1,
			seed: func(t *testing.T, creds string) {
				seedApprover(t, creds, "other", fixtureKeyOf(t, "ed25519.pub"), approveT0-10)
				revokeSeeded(t, creds, "other", approveT0-5)
			},
			args: func(t *testing.T, db string) []string { return enrollArgs(db, "mark-mac", ed) },
			want: []string{edFingerprint + ` is already enrolled as approver "other" (a key is never enrolled twice, even after revoke)`}},
	}
	for _, arm := range arms {
		t.Run(arm.name, func(t *testing.T) {
			db, creds := arm.setup(t)
			if arm.seed != nil {
				arm.seed(t, creds)
			}
			args := arm.args(t, db)
			var out, errw bytes.Buffer
			var opened bool
			done := make(chan int, 1)
			go func() {
				done <- runApprover(args, &out, &errw, approveEnv(nil, &opened, atTime(approveT0)))
			}()
			var code int
			select {
			case code = <-done:
			case <-time.After(5 * time.Second):
				t.Fatal("enrol did not refuse within 5 s (a blocking open, or a terminal reached)")
			}
			if code != exitUsage {
				t.Fatalf("code %d, want exitUsage; stderr %q", code, errw.String())
			}
			if opened {
				t.Fatal("the terminal was opened for a refused enrol")
			}
			for _, w := range arm.want {
				if !strings.Contains(errw.String(), w) {
					t.Fatalf("stderr %q lacks %q", errw.String(), w)
				}
			}
			if arm.forbid != "" && strings.Contains(errw.String(), arm.forbid) {
				t.Fatalf("stderr echoes the private key: %q", errw.String())
			}
			if arm.checkRows {
				if n := len(approverRows(t, creds)); n != arm.wantRows {
					t.Fatalf("%d approver rows, want %d (the seeds only)", n, arm.wantRows)
				}
			}
		})
	}
}

// CA6: a race the store reports, a failure after the write opened, and a busy writer lock each
// refuse or fail, say so, and enrol nothing. The seam arms replace approverEnrol, restored on
// cleanup.
func TestApproverEnrollRaceAndFailure(t *testing.T) {
	arms := []struct {
		name    string
		seamErr func() error
		hold    bool
		code    int
		want    []string
	}{
		{name: "name-race", code: exitUsage, want: []string{"(another enrol ran first); nothing was enrolled"},
			seamErr: func() error { return fmt.Errorf("%w: %q", store.ErrApproverNameTaken, "mark-mac") }},
		{name: "key-race", code: exitUsage, want: []string{"(another enrol ran first); nothing was enrolled"},
			seamErr: func() error { return fmt.Errorf("%w: it is approver %q's key", store.ErrApproverKeyEnrolled, "x") }},
		{name: "failure", code: exitFatal, want: []string{`cannot enrol approver "mark-mac": disk on fire`},
			seamErr: func() error { return errors.New("disk on fire") }},
		{name: "writer-busy", code: exitFatal, hold: true,
			want: []string{"cannot open the credentials store for writing", "another session new/mint/revoke is writing this credentials store"}},
	}
	for _, arm := range arms {
		t.Run(arm.name, func(t *testing.T) {
			db, creds := approverStore(t)
			if arm.seamErr != nil {
				orig := approverEnrol
				t.Cleanup(func() { approverEnrol = orig })
				seamErr := arm.seamErr()
				approverEnrol = func(*store.Credentials, context.Context, store.ApproverRow) error { return seamErr }
			}
			var hold *store.Credentials
			if arm.hold {
				h, err := store.OpenCredentials(creds)
				if err != nil {
					t.Fatal(err)
				}
				hold = h
			}
			term := &lockedTerm{in: strings.NewReader("y\n")}
			var out, errw bytes.Buffer
			code := runApprover(enrollArgs(db, "mark-mac", fixturePub("ed25519.pub")), &out, &errw, approveEnv(term, nil, atTime(approveT0)))
			if hold != nil {
				if err := hold.Close(); err != nil {
					t.Fatal(err)
				}
			}
			if code != arm.code {
				t.Fatalf("code %d, want %d; stderr %q", code, arm.code, errw.String())
			}
			for _, w := range arm.want {
				if !strings.Contains(errw.String(), w) {
					t.Fatalf("stderr %q lacks %q", errw.String(), w)
				}
			}
			if n := len(approverRows(t, creds)); n != 0 {
				t.Fatalf("%d approver rows, want 0", n)
			}
		})
	}
}

// CA9: an unanswered prompt times out after approverConfirmBudget, the writer lock is free while
// it waits, and a missing prompt timer is an internal error that opens nothing.
func TestApproverEnrollUnansweredPromptTimesOut(t *testing.T) {
	db, creds := approverStore(t)
	pr, pw := io.Pipe()
	t.Cleanup(func() { _ = pw.Close() })
	term := &lockedTerm{in: pr}
	fire := make(chan time.Time)
	var mu sync.Mutex
	var waited time.Duration
	env := sessionEnv{openTerminal: func() (io.ReadWriteCloser, error) { return term, nil }, now: atTime(approveT0),
		after: func(d time.Duration) <-chan time.Time { mu.Lock(); waited = d; mu.Unlock(); return fire }}
	type res struct {
		code   int
		stderr string
	}
	done := make(chan res, 1)
	go func() {
		var out, errw bytes.Buffer
		code := runApprover(enrollArgs(db, "mark-mac", fixturePub("ed25519.pub")), &out, &errw, env)
		done <- res{code, errw.String()}
	}()
	deadline := time.Now().Add(5 * time.Second)
	for !strings.Contains(term.Output(), "[y/N]") {
		if time.Now().After(deadline) {
			t.Fatal("no prompt within 5 s")
		}
		time.Sleep(10 * time.Millisecond)
	}
	w, err := store.OpenCredentials(creds) // the writer lock is free while the prompt waits
	if err != nil {
		t.Fatalf("the credentials writer lock is held during the prompt: %v", err)
	}
	_ = w.Close()
	close(fire)
	var r res
	select {
	case r = <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("enroll did not return within 5 s of the timer")
	}
	if r.code != exitUsage {
		t.Fatalf("code %d, want exitUsage; stderr %q", r.code, r.stderr)
	}
	if want := "ailang-worldd approver enroll: no answer within 2m0s; nothing was enrolled\n"; r.stderr != want {
		t.Fatalf("stderr %q, want %q", r.stderr, want)
	}
	mu.Lock()
	got := waited
	mu.Unlock()
	if got != 120*time.Second {
		t.Fatalf("the prompt timer was set to %s, want 2m0s", got)
	}
	if !term.Closed() {
		t.Fatal("the terminal was not closed")
	}
	if n := len(approverRows(t, creds)); n != 0 {
		t.Fatalf("%d approver rows, want 0", n)
	}

	// Arm 2: no prompt timer at all is an internal error, and no terminal is opened.
	var opened bool
	var out, errw bytes.Buffer
	noTimer := approveEnv(nil, &opened, atTime(approveT0))
	noTimer.after = nil
	if code := runApprover(enrollArgs(db, "mark-mac", fixturePub("ed25519.pub")), &out, &errw, noTimer); code != exitFatal {
		t.Fatalf("no timer: code %d, want exitFatal; stderr %q", code, errw.String())
	}
	if !strings.Contains(errw.String(), "internal error: no prompt timer (sessionEnv.after)") {
		t.Fatalf("no timer: stderr %q lacks the internal error", errw.String())
	}
	if opened {
		t.Fatal("no timer: the terminal was opened")
	}
	if n := len(approverRows(t, creds)); n != 0 {
		t.Fatalf("no timer: %d approver rows, want 0", n)
	}
}
