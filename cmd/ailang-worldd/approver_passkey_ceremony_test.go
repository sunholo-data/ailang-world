package main

// Row 201 A2 sprint 3 M6d (design_docs/planned/w-world-approvals-by-key-a2.md §2.12, §2.8 [r1];
// CW2, CW5): the passkey enrolment ceremony at the terminal, the code then y, and the substitute it
// refuses. Every store is under t.TempDir(); every passkey is approvaltest.Passkey's fixed seed
// (D-WORLD-109).

import (
	"bytes"
	"context"
	"encoding/base64"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
)

// ceremonyTerm is the operator's terminal in a passkey enrolment: its n-th Read returns answer(n),
// one line per read, as a terminal in canonical mode does; writes and Close are recorded.
type ceremonyTerm struct {
	lockedTerm
	reads  int
	answer func(n int) (string, error)
}

func (c *ceremonyTerm) Read(p []byte) (int, error) {
	s, err := c.answer(c.reads)
	c.reads++
	return copy(p, s), err
}

// browserFor is what the test's browser saves for a fragment: pk's registration for its challenge
// at its origin and rp.
func browserFor(pk *approvaltest.Passkey, e ...approvaltest.Edit) func(t *testing.T, f ceremonyLink) []byte {
	return func(t *testing.T, f ceremonyLink) []byte { return pk.Register(f.Origin, f.RP, challengeOf(t, f), e...) }
}

// enrolRun is one approver enroll --format webauthn run.
type enrolRun struct {
	code           int
	stdout, stderr string
	term           *ceremonyTerm
	file           []byte
}

// enrol runs approver enroll --format webauthn on r. At the code prompt the test's browser saves
// save(fragment) to the registration path (nil: it made the path itself) and the terminal answers
// typed(file); at the y/N it answers yes, after running atYes when not nil. The credentials writer
// lock must be free at both prompts.
func (r enrolRig) enrol(t *testing.T, atYes func(), save func(t *testing.T, f ceremonyLink) []byte, typed func(file []byte) string, yes string) enrolRun {
	t.Helper()
	var out, errw bytes.Buffer
	run := enrolRun{term: &ceremonyTerm{}}
	free := func(prompt string) {
		if w, err := store.OpenCredentials(r.creds); err != nil {
			t.Errorf("the credentials writer lock is held at the %s prompt: %v", prompt, err)
		} else {
			_ = w.Close()
		}
	}
	run.term.answer = func(n int) (string, error) {
		switch n {
		case 0:
			free("code")
			if run.file = save(t, linkFragment(t, out.String())); run.file != nil {
				if err := os.WriteFile(r.reg, run.file, 0o600); err != nil {
					t.Error(err)
				}
			}
			return typed(run.file), nil
		case 1:
			free("y/N")
			if atYes != nil {
				atYes()
			}
			return yes, nil
		}
		return "", os.ErrClosed
	}
	run.code = runApprover(r.args(), &out, &errw, enrolEnv(run.term, nil))
	run.stdout, run.stderr = out.String(), errw.String()
	return run
}

// typeCode answers the code prompt with the file's grouped code.
func typeCode(t *testing.T) func(file []byte) string {
	return func(file []byte) string { return codeOf(t, file) + "\n" }
}

// enrolled is the store's only approver and passkey, or a test failure.
func enrolled(t *testing.T, creds string) (store.ApproverRow, store.PasskeyRow) {
	t.Helper()
	rows := approverRows(t, creds)
	if len(rows) != 1 || passkeyCount(t, creds) != 1 {
		t.Fatalf("%d approvers and %d passkeys, want 1 and 1", len(rows), passkeyCount(t, creds))
	}
	ro, err := store.OpenCredentialsReadOnly(creds)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ro.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	p, ok, err := ro.Passkey(ctx, rows[0].Name)
	if err != nil || !ok {
		t.Fatalf("passkey of %q: %v %v", rows[0].Name, ok, err)
	}
	return rows[0], p
}

// assertNothingEnrolled fails unless the run refused with code, wrote want to stderr, closed the
// terminal and left both tables empty.
func assertNothingEnrolled(t *testing.T, r enrolRig, run enrolRun, code int, want ...string) {
	t.Helper()
	if run.code != code {
		t.Fatalf("code %d, want %d; stderr %q", run.code, code, run.stderr)
	}
	for _, w := range want {
		if !strings.Contains(run.stderr, w) {
			t.Fatalf("stderr %q lacks %q", run.stderr, w)
		}
	}
	if run.term.reads > 0 && !run.term.Closed() {
		t.Fatal("the terminal was not closed")
	}
	if n, m := len(approverRows(t, r.creds)), passkeyCount(t, r.creds); n != 0 || m != 0 {
		t.Fatalf("%d approvers and %d passkeys, want none", n, m)
	}
}

// CW2: the right code and y enrol exactly the registration the browser made for the printed
// challenge: one approvers row (webauthn, the canonical COSE key) and one passkey row; the y/N
// shows RenderPasskeyEnrolment; N14 names the passkey.
func TestPasskeyEnrollEnrolsWhatTheBrowserMade(t *testing.T) {
	for _, alg := range []int64{approval.AlgES256, approval.AlgEdDSA} {
		pk := approvaltest.NewPasskey(alg, 41)
		r := newEnrolRig(t)
		run := r.enrol(t, nil, browserFor(pk), typeCode(t), "y\n")
		if run.code != exitOK {
			t.Fatalf("alg %d: code %d; stderr %q", alg, run.code, run.stderr)
		}
		f := linkFragment(t, run.stdout)
		seq := make([]byte, 32)
		_ = seqRandom()(seq)
		if want := (ceremonyLink{V: 1, Op: "create", RP: pkRPID, Origin: pkOrigin, Name: "mark-mac", Challenge: base64.RawURLEncoding.EncodeToString(seq)}); !reflect.DeepEqual(f, want) {
			t.Fatalf("fragment %+v, want %+v (the challenge from env.random)", f, want)
		}
		if !strings.Contains(run.stdout, "\n  "+pkOrigin+"/surface/passkey#") {
			t.Fatalf("stdout %q lacks the link at the origin", run.stdout)
		}
		a, err := approval.NewPasskeyApprover(approval.RegistrationInput{Name: "mark-mac", Origin: pkOrigin, Challenge: [32]byte(seq), Response: run.file, Now: approveT0})
		if err != nil {
			t.Fatal(err)
		}
		shown := run.term.Output()
		for _, w := range []string{"Enrolment code shown by the browser (XXXXX-XXXXX-XXXXX-XXXXX-XXXXXX): ", approval.RenderPasskeyEnrolment(a)} {
			if !strings.Contains(shown, w) {
				t.Fatalf("the terminal lacks %q:\n%s", w, shown)
			}
		}
		if !strings.HasSuffix(shown, "Enrol approver mark-mac? [y/N] ") || !run.term.Closed() {
			t.Fatalf("the terminal does not end with the y/N, or was not closed: %q", shown)
		}
		label := passkeyLabel(alg, pkRPID, pk.CredentialID())
		if want := "enrolled approver mark-mac: " + label + "; approvals by it record approved_by=webauthn:mark-mac\n"; !strings.HasSuffix(run.stdout, want) {
			t.Fatalf("stdout %q lacks %q", run.stdout, want)
		}
		row, p := enrolled(t, r.creds)
		wantRow := store.ApproverRow{Name: "mark-mac", Format: "webauthn", PublicKey: pk.Key().Canonical(), EnrolledAt: approveT0}
		wantP := store.PasskeyRow{Name: "mark-mac", RPID: pkRPID, Origin: pkOrigin, CredentialID: pk.CredentialID(),
			AAGUID: bytes.Repeat([]byte{1, 2, 3, 4, 5, 6, 7, 8}, 2), Alg: alg, SignCount: 1}
		if !reflect.DeepEqual(row, wantRow) || !reflect.DeepEqual(p, wantP) {
			t.Fatalf("stored %+v %+v, want %+v %+v", row, p, wantRow, wantP)
		}
	}
	// A backup-eligible passkey is stored as one (W13 compares against it).
	r := newEnrolRig(t)
	if run := r.enrol(t, nil, browserFor(approvaltest.NewPasskey(approval.AlgEdDSA, 42), approvaltest.Edit{Set: 0x08}), typeCode(t), "y\n"); run.code != exitOK {
		t.Fatalf("BE: code %d; stderr %q", run.code, run.stderr)
	}
	if _, p := enrolled(t, r.creds); !p.BackupEligible {
		t.Fatal("a backup-eligible passkey was stored as not")
	}
}

// CW5 [r1]: read-then-substitute. The attacker reads the genuine registration G, computes its code,
// and saves its own registration X for the same printed challenge; the operator types G's code.
// WN10 refuses, the y/N is never shown, nothing is enrolled. A code differing from G's only in its
// last character, its first group or its length, or holding a 0, is refused on G's own file; G's
// code in lowercase with dashes, spaces and a tab, or as one string, enrols G.
func TestPasskeyEnrollRefusesASubstituteEvenWithTheGenuineCode(t *testing.T) {
	genuine, attacker := approvaltest.NewPasskey(approval.AlgES256, 44), approvaltest.NewPasskey(approval.AlgES256, 45)
	var genuineCode string
	substitute := func(t *testing.T, f ceremonyLink) []byte {
		genuineCode = codeOf(t, genuine.Register(f.Origin, f.RP, challengeOf(t, f)))
		return attacker.Register(f.Origin, f.RP, challengeOf(t, f))
	}
	r := newEnrolRig(t)
	run := r.enrol(t, nil, substitute, func([]byte) string { return genuineCode + "\n" }, "y\n")
	assertNothingEnrolled(t, r, run, exitUsage, "the code you typed is not the code of "+r.reg+"; nothing was enrolled")
	if genuineCode == codeOf(t, run.file) || strings.Contains(run.term.Output(), "[y/N]") {
		t.Fatalf("the substitute's code equals the genuine one, or the y/N was shown: %q", run.term.Output())
	}

	other := func(c byte) byte {
		if c == 'A' {
			return 'B'
		}
		return 'A'
	}
	for _, arm := range []struct {
		name  string
		typed func(code string) string // code is the grouped genuine code
		ok    bool
	}{
		{name: "last-character", typed: func(c string) string { return c[:len(c)-1] + string(other(c[len(c)-1])) }},
		{name: "first-group", typed: func(c string) string { return string(other(c[0])) + c[1:] }},
		{name: "25-characters", typed: func(c string) string { return c[:len(c)-1] }},
		{name: "27-characters", typed: func(c string) string { return c + "A" }},
		{name: "a-zero", typed: func(c string) string { return "0" + c[1:] }},
		{name: "lowercase-dashes-spaces-tab", ok: true, typed: func(c string) string {
			return strings.ToLower(strings.Replace(strings.Replace(c, "-", " - ", 1), "-", "\t", 2))
		}},
		{name: "one-string", ok: true, typed: func(c string) string { return strings.ReplaceAll(c, "-", "") }},
	} {
		t.Run(arm.name, func(t *testing.T) {
			r := newEnrolRig(t)
			run := r.enrol(t, nil, browserFor(genuine), func(file []byte) string { return arm.typed(codeOf(t, file)) + "\n" }, "y\n")
			if arm.ok {
				if run.code != exitOK {
					t.Fatalf("code %d; stderr %q", run.code, run.stderr)
				}
				enrolled(t, r.creds)
				return
			}
			assertNothingEnrolled(t, r, run, exitUsage, "the code you typed is not the code of ")
		})
	}
}
