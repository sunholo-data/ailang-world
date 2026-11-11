package main

// Row 201 A2 sprint 3 M6e (design_docs/planned/w-world-approvals-by-key-a2.md §2.12; CW4, CW6): the
// passkey enrolment's unanswered, declined and refused ceremonies, and its write's races. Every
// store is under t.TempDir(); every passkey is approvaltest.Passkey's fixed seed (D-WORLD-109).

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
)

// CW4: an unanswered code (WN6, after exactly 300 s), a code cut short (WN7), and a y/N declined,
// cut short or unanswered (N10, N10a, after 120 s) each refuse, close the terminal and enrol
// nothing; the writer lock is free at both prompts (enrol's helper checks it). Each arm runs under
// a 5 s guard, so an unbounded read fails the arm rather than the run.
func TestPasskeyEnrollDeclinedOrUnanswered(t *testing.T) {
	pk := approvaltest.NewPasskey(approval.AlgES256, 43)
	const aborted = "aborted (not confirmed); nothing was enrolled. fix: run it again and type y at the prompt"
	for _, arm := range []struct {
		name  string
		code  func(t *testing.T, file []byte) (string, error) // the code prompt's answer; nil: the file's code
		yes   func() (string, error)                          // the y/N's answer
		block int                                             // the prompt left unanswered (its timer fires): 0 none, 1 the code, 2 the y/N
		want  string
	}{
		{name: "WN6-code-timeout", block: 1, want: "ailang-worldd approver enroll: no code within 5m0s; nothing was enrolled\n"},
		{name: "WN7-code-eof", code: func(*testing.T, []byte) (string, error) { return "", io.EOF },
			want: "ailang-worldd approver enroll: aborted (no complete code line); nothing was enrolled. fix: run it again and type the code the page shows\n"},
		{name: "WN7-code-partial", code: func(t *testing.T, file []byte) (string, error) { return codeOf(t, file), io.EOF },
			want: "ailang-worldd approver enroll: aborted (no complete code line)"},
		{name: "WN7-code-read-error", code: func(t *testing.T, file []byte) (string, error) { return codeOf(t, file)[:7], errInjectedRead },
			want: "ailang-worldd approver enroll: aborted (no complete code line)"},
		{name: "N10-n", yes: func() (string, error) { return "n\n", nil }, want: aborted},
		{name: "N10-eof", yes: func() (string, error) { return "", io.EOF }, want: aborted},
		{name: "N10-y-then-eof", yes: func() (string, error) { return "y", io.EOF }, want: aborted},
		{name: "N10a-y-timeout", block: 2, want: "ailang-worldd approver enroll: no answer within 2m0s; nothing was enrolled\n"},
	} {
		t.Run(arm.name, func(t *testing.T) {
			r := newEnrolRig(t)
			var mu sync.Mutex
			var asked []time.Duration
			fire := make(chan time.Time) // closed by the unanswered prompt's read, once it waits
			env := enrolEnv(nil, nil)
			env.after = func(d time.Duration) <-chan time.Time {
				mu.Lock()
				defer mu.Unlock()
				asked = append(asked, d)
				if arm.block != 0 && d == []time.Duration{passkeyCeremonyBudget, approverConfirmBudget}[arm.block-1] {
					return fire
				}
				return make(chan time.Time)
			}
			var out, errw bytes.Buffer
			term := &ceremonyTerm{}
			term.answer = func(n int) (string, error) {
				if w, err := store.OpenCredentials(r.creds); err != nil { // the writer lock is free at each prompt
					t.Errorf("prompt %d: the credentials writer lock is held: %v", n, err)
				} else {
					_ = w.Close()
				}
				if n+1 == arm.block {
					close(fire)
					<-t.Context().Done()
					return "", io.EOF
				}
				if n == 0 {
					file := browserFor(pk)(t, linkFragment(t, out.String()))
					if err := os.WriteFile(r.reg, file, 0o600); err != nil {
						t.Error(err)
					}
					if arm.code != nil {
						return arm.code(t, file)
					}
					return codeOf(t, file) + "\n", nil
				}
				return arm.yes()
			}
			env.openTerminal = func() (io.ReadWriteCloser, error) { return term, nil }
			done := make(chan int, 1)
			go func() { done <- runApprover(r.args(), &out, &errw, env) }()
			var code int
			select {
			case code = <-done:
			case <-time.After(5 * time.Second):
				t.Fatal("enrol did not return within 5 s (an unbounded read)")
			}
			assertNothingEnrolled(t, r, enrolRun{code: code, stderr: errw.String(), term: term}, exitUsage, arm.want)
			mu.Lock()
			defer mu.Unlock()
			want := []time.Duration{passkeyCeremonyBudget}
			if arm.yes != nil || arm.block == 2 {
				want = append(want, approverConfirmBudget)
			}
			if !reflect.DeepEqual(asked, want) || passkeyCeremonyBudget != 300*time.Second {
				t.Fatalf("the prompt timers were %v, want %v (the code 5m0s, the y/N 2m0s)", asked, want)
			}
		})
	}
}

// CW6: a missing, failed or all-zero challenge source is an internal error before any terminal
// (WN5); a FIFO at the registration path (WN8), a registration the law refuses (WN9), and a key or
// credential another enrol stored while the code prompt waited (WN11) refuse before the y/N; and
// the write's races (N12) refuse, its failures (N11, N13) are fatal, each enrolling nothing.
func TestPasskeyEnrollRaceAndRefusals(t *testing.T) {
	pk := approvaltest.NewPasskey(approval.AlgEdDSA, 46)
	t.Run("WN5", func(t *testing.T) {
		for name, random := range map[string]func([]byte) error{
			"nil": nil,
			"failed": func(b []byte) error {
				copy(b, bytes.Repeat([]byte{0xff}, 32))
				return errors.New("entropy unavailable")
			},
			"zero": func(b []byte) error { clear(b); return nil },
		} {
			r := newEnrolRig(t)
			var opened bool
			env := enrolEnv(nil, &opened)
			env.random = random
			var out, errw bytes.Buffer
			code := runApprover(r.args(), &out, &errw, env)
			if code != exitFatal || opened || out.Len() != 0 || !strings.Contains(errw.String(), "internal error: no fresh challenge") {
				t.Fatalf("%s: code %d, terminal opened %v, stdout %q, stderr %q; want exitFatal before any terminal", name, code, opened, out.String(), errw.String())
			}
		}
	})
	for _, arm := range []struct {
		name string
		save func(t *testing.T, r enrolRig, f ceremonyLink) []byte // the browser's step; nil: it made the path itself
		want string
		pre  int // approvers another enrol stored
	}{
		{name: "WN8-fifo", want: "is not a regular file; nothing was enrolled", save: func(t *testing.T, r enrolRig, f ceremonyLink) []byte {
			if err := syscall.Mkfifo(r.reg, 0o600); err != nil {
				t.Error(err)
			}
			return nil
		}},
		{name: "WN8-over-16k", want: "is 16385 bytes, over the 16384-byte bound; nothing was enrolled", save: func(t *testing.T, r enrolRig, f ceremonyLink) []byte {
			if err := os.WriteFile(r.reg, bytes.Repeat([]byte(" "), 16<<10+1), 0o600); err != nil {
				t.Error(err)
			}
			return nil
		}},
		{name: "WN9-uv-clear", want: "does not enrol: " + approval.ErrUserVerification.Error(), save: func(t *testing.T, r enrolRig, f ceremonyLink) []byte {
			return browserFor(pk, approvaltest.Edit{Clear: 0x04})(t, f)
		}},
		{name: "WN9-another-challenge", want: "does not enrol: " + approval.ErrChallenge.Error(), save: func(t *testing.T, r enrolRig, f ceremonyLink) []byte {
			return pk.Register(f.Origin, f.RP, [32]byte{9})
		}},
		{name: "WN11-key", pre: 1, want: `this passkey's key is already enrolled as approver "other"`, save: func(t *testing.T, r enrolRig, f ceremonyLink) []byte {
			seedPasskeyRows(t, r.creds, "other", pk.Key().Canonical(), approvaltest.NewPasskey(approval.AlgEdDSA, 47).CredentialID())
			return browserFor(pk)(t, f)
		}},
		{name: "WN11-credential", pre: 1, want: "credential cred " + passkeyCred(pk.CredentialID()) + ` is already enrolled as approver "other"`,
			save: func(t *testing.T, r enrolRig, f ceremonyLink) []byte {
				seedPasskeyRows(t, r.creds, "other", approvaltest.NewPasskey(approval.AlgEdDSA, 47).Key().Canonical(), pk.CredentialID())
				return browserFor(pk)(t, f)
			}},
	} {
		t.Run(arm.name, func(t *testing.T) {
			r := newEnrolRig(t)
			run := r.enrol(t, nil, func(t *testing.T, f ceremonyLink) []byte { return arm.save(t, r, f) }, func(file []byte) string {
				if file == nil { // the browser saved no registration to show a code for
					return "AAAAA-AAAAA-AAAAA-AAAAA-AAAAAA\n"
				}
				return codeOf(t, file) + "\n"
			}, "y\n")
			if run.code != exitUsage || !strings.Contains(run.stderr, arm.want) || strings.Contains(run.term.Output(), "[y/N]") {
				t.Fatalf("code %d, stderr %q, terminal %q; want exitUsage naming %q before the y/N", run.code, run.stderr, run.term.Output(), arm.want)
			}
			if n, m := len(approverRows(t, r.creds)), passkeyCount(t, r.creds); n != arm.pre || m != arm.pre {
				t.Fatalf("%d approvers and %d passkeys, want %d (another enrol's)", n, m, arm.pre)
			}
		})
	}
	for _, arm := range []struct {
		name string
		err  error // the write's error (the seam)
		hold bool  // another writer holds the lock at the write
		code int
		want string
	}{
		{name: "N12-name", err: store.ErrApproverNameTaken, code: exitUsage, want: "(another enrol ran first); nothing was enrolled"},
		{name: "N12-key", err: store.ErrApproverKeyEnrolled, code: exitUsage, want: "(another enrol ran first); nothing was enrolled"},
		{name: "N12-credential", err: store.ErrPasskeyCredentialEnrolled, code: exitUsage, want: "(another enrol ran first); nothing was enrolled"},
		{name: "N13-failure", err: errors.New("disk on fire"), code: exitFatal, want: `cannot enrol approver "mark-mac": disk on fire`},
		{name: "N11-writer-busy", hold: true, code: exitFatal, want: "cannot open the credentials store for writing"},
	} {
		t.Run(arm.name, func(t *testing.T) {
			if arm.err != nil {
				orig := approverEnrolPasskey
				t.Cleanup(func() { approverEnrolPasskey = orig })
				err := fmt.Errorf("%w: approver %q holds it", arm.err, "x")
				approverEnrolPasskey = func(*store.Credentials, context.Context, store.ApproverRow, store.PasskeyRow) error { return err }
			}
			r := newEnrolRig(t)
			var hold *store.Credentials
			atYes := func() {
				if arm.hold {
					var err error
					if hold, err = store.OpenCredentials(r.creds); err != nil {
						t.Error(err)
					}
				}
			}
			run := r.enrol(t, atYes, browserFor(pk), typeCode(t), "y\n")
			if hold != nil {
				_ = hold.Close()
			}
			assertNothingEnrolled(t, r, run, arm.code, arm.want)
		})
	}
}

// seedPasskeyRows stores approver name with key and credential straight through the store.
func seedPasskeyRows(t *testing.T, creds, name string, key, credentialID []byte) {
	t.Helper()
	st, err := store.OpenCredentials(creds)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() { _ = st.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := st.EnrolPasskey(ctx, store.ApproverRow{Name: name, Format: approval.FormatWebAuthn, PublicKey: key, EnrolledAt: approveT0 - 10},
		store.PasskeyRow{Name: name, RPID: pkRPID, Origin: pkOrigin, CredentialID: credentialID, AAGUID: make([]byte, 16), Alg: approval.AlgEdDSA}); err != nil {
		t.Error(err)
	}
}
