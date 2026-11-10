package main

// Row 201 A1b sprint 2 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint2-plan.md §3
// M3e, M3f): `approver enroll` refusals, the stored key, the declined arms and the races. Every
// store is under t.TempDir(); every key is a checked-in fixture (D-WORLD-109).

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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
