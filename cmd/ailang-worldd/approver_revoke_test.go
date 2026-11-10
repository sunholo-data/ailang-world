package main

// Row 201 A1b sprint 2 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint2-plan.md §3
// M3g, CA8 and CA10): `approver revoke`. Every store is under t.TempDir(); every key is a
// checked-in fixture (D-WORLD-109). No test here runs t.Parallel: the seams are package vars.

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
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
)

// revokeRig is a fresh store with mark-mac enrolled (ed25519 fixture) at approveT0, and its rows.
func revokeRig(t *testing.T) (db, creds string, before []store.ApproverRow) {
	t.Helper()
	db, creds = approverStore(t)
	seedApprover(t, creds, "mark-mac", fixtureKeyOf(t, "ed25519.pub"), approveT0)
	return db, creds, approverRows(t, creds)
}

// CA8: revoke's arms. Each refusal before a write is exitUsage with its exact text and writes
// nothing; the write's own refusals and failure are exitFatal or exitUsage by constant.
func TestApproverRevoke(t *testing.T) {
	t.Run("yes", func(t *testing.T) {
		db, creds, before := revokeRig(t)
		var opened bool
		term := &lockedTerm{in: strings.NewReader("y\n")}
		var out, errw bytes.Buffer
		code := runApprover([]string{"revoke", "mark-mac", "--db", db}, &out, &errw, approveEnv(term, &opened, atTime(approveT0+100)))
		if code != exitOK {
			t.Fatalf("code %d, want exitOK; stderr %q", code, errw.String())
		}
		if want := "revoked approver mark-mac\nno session was approved by it\n"; out.String() != want { // K7, K9 (sprint 4 D3)
			t.Fatalf("stdout %q, want %q", out.String(), want)
		}
		wantPrompt := "Revoke approver mark-mac (ssh-ed25519 " + edFingerprint + ")? It will never approve again, and its name and key can never be enrolled again. [y/N] "
		if term.Output() != wantPrompt {
			t.Fatalf("prompt %q, want %q", term.Output(), wantPrompt)
		}
		if !term.Closed() {
			t.Fatal("the terminal was not closed")
		}
		rows := approverRows(t, creds)
		if len(rows) != 1 {
			t.Fatalf("%d rows, want 1 (revoke keeps the row)", len(rows))
		}
		if rows[0].RevokedAt != approveT0+100 {
			t.Fatalf("RevokedAt %d, want %d", rows[0].RevokedAt, approveT0+100)
		}
		got := rows[0]
		got.RevokedAt = before[0].RevokedAt
		if !reflect.DeepEqual(got, before[0]) {
			t.Fatalf("other fields changed: got %+v, before %+v", got, before[0])
		}
	})

	t.Run("arity-none", func(t *testing.T) {
		db, _, _ := revokeRig(t)
		assertRevokeRefusal(t, []string{"revoke", "--db", db}, approveEnv(nil, nil, atTime(approveT0)), exitUsage, "usage: ailang-worldd approver")
	})
	t.Run("arity-two", func(t *testing.T) {
		db, _, _ := revokeRig(t)
		assertRevokeRefusal(t, []string{"revoke", "a", "b", "--db", db}, approveEnv(nil, nil, atTime(approveT0)), exitUsage, "usage: ailang-worldd approver")
	})

	t.Run("unknown", func(t *testing.T) {
		db, _, _ := revokeRig(t)
		var opened bool
		assertRevokeRefusal(t, []string{"revoke", "nobody", "--db", db}, approveEnv(nil, &opened, atTime(approveT0+100)), exitUsage,
			`no approver named "nobody"`, "approver list shows them")
		if opened {
			t.Fatal("unknown: the terminal was opened")
		}
	})

	t.Run("already-revoked", func(t *testing.T) {
		db, creds, _ := revokeRig(t)
		revokeSeeded(t, creds, "mark-mac", approveT0+50)
		var opened bool
		assertRevokeRefusal(t, []string{"revoke", "mark-mac", "--db", db}, approveEnv(nil, &opened, atTime(approveT0+100)), exitUsage,
			`approver "mark-mac" cannot be revoked: approval: the approver is revoked`)
		if opened {
			t.Fatal("already-revoked: the terminal was opened")
		}
	})

	t.Run("clock-behind", func(t *testing.T) {
		db, _, _ := revokeRig(t)
		var opened bool
		assertRevokeRefusal(t, []string{"revoke", "mark-mac", "--db", db}, approveEnv(nil, &opened, atTime(approveT0-1)), exitUsage,
			"cannot be revoked: approval: the clock is behind")
		if opened {
			t.Fatal("clock-behind: the terminal was opened")
		}
	})

	t.Run("v4", func(t *testing.T) {
		db, _ := approverStoreV4(t)
		assertRevokeRefusal(t, []string{"revoke", "mark-mac", "--db", db}, approveEnv(nil, nil, atTime(approveT0+100)), exitUsage,
			"approvers need schema 5")
	})

	t.Run("no-credentials-store", func(t *testing.T) {
		db := dbPathOf(t)
		assertRevokeRefusal(t, []string{"revoke", "mark-mac", "--db", db}, approveEnv(nil, nil, atTime(approveT0+100)), exitUsage,
			"cannot read the credentials store")
	})

	t.Run("no-terminal", func(t *testing.T) {
		db, creds, before := revokeRig(t)
		env := noTerminalEnv(approveT0 + 100)
		env.after = time.After
		var out, errw bytes.Buffer
		code := runApprover([]string{"revoke", "mark-mac", "--db", db}, &out, &errw, env)
		if code != exitUsage {
			t.Fatalf("code %d, want exitUsage; stderr %q", code, errw.String())
		}
		for _, w := range []string{"nothing was revoked", sessionTTYFix} {
			if !strings.Contains(errw.String(), w) {
				t.Fatalf("stderr %q lacks %q", errw.String(), w)
			}
		}
		if rows := approverRows(t, creds); !reflect.DeepEqual(rows, before) {
			t.Fatalf("rows changed: %+v, want %+v", rows, before)
		}
	})

	t.Run("declined", func(t *testing.T) {
		const want = "ailang-worldd approver revoke: aborted (not confirmed); nothing was revoked. fix: run it again and type y at the prompt\n"
		arms := []struct {
			name string
			in   io.Reader
		}{
			{name: "n", in: strings.NewReader("n\n")},
			{name: "y-then-eof", in: &partialReader{data: "y", err: io.EOF}},
		}
		for _, arm := range arms {
			t.Run(arm.name, func(t *testing.T) {
				db, creds, before := revokeRig(t)
				var opened bool
				term := &lockedTerm{in: arm.in}
				var out, errw bytes.Buffer
				code := runApprover([]string{"revoke", "mark-mac", "--db", db}, &out, &errw, approveEnv(term, &opened, atTime(approveT0+100)))
				if code != exitUsage {
					t.Fatalf("code %d, want exitUsage; stderr %q", code, errw.String())
				}
				if errw.String() != want {
					t.Fatalf("stderr %q, want %q", errw.String(), want)
				}
				if !term.Closed() {
					t.Fatal("the terminal was not closed")
				}
				if rows := approverRows(t, creds); !reflect.DeepEqual(rows, before) {
					t.Fatalf("rows changed: %+v, want %+v", rows, before)
				}
			})
		}
	})

	t.Run("race-revoked", func(t *testing.T) {
		assertRevokeRace(t, fmt.Errorf("%w: %q at %d", store.ErrApproverRevoked, "mark-mac", 1))
	})
	t.Run("race-unknown", func(t *testing.T) {
		assertRevokeRace(t, fmt.Errorf("%w: %q", store.ErrApproverUnknown, "mark-mac"))
	})

	t.Run("failure", func(t *testing.T) {
		db, creds, before := revokeRig(t)
		orig := approverRevoke
		t.Cleanup(func() { approverRevoke = orig })
		approverRevoke = func(*store.Credentials, context.Context, string, int64) error { return errors.New("disk on fire") }
		term := &lockedTerm{in: strings.NewReader("y\n")}
		var out, errw bytes.Buffer
		code := runApprover([]string{"revoke", "mark-mac", "--db", db}, &out, &errw, approveEnv(term, nil, atTime(approveT0+100)))
		if code != exitFatal {
			t.Fatalf("code %d, want exitFatal; stderr %q", code, errw.String())
		}
		if !strings.Contains(errw.String(), `cannot revoke approver "mark-mac": disk on fire`) {
			t.Fatalf("stderr %q lacks the failure text", errw.String())
		}
		if rows := approverRows(t, creds); !reflect.DeepEqual(rows, before) {
			t.Fatalf("rows changed: %+v, want %+v", rows, before)
		}
	})

	t.Run("writer-busy", func(t *testing.T) {
		db, creds, before := revokeRig(t)
		hold, err := store.OpenCredentials(creds)
		if err != nil {
			t.Fatal(err)
		}
		term := &lockedTerm{in: strings.NewReader("y\n")}
		var out, errw bytes.Buffer
		code := runApprover([]string{"revoke", "mark-mac", "--db", db}, &out, &errw, approveEnv(term, nil, atTime(approveT0+100)))
		if err := hold.Close(); err != nil {
			t.Fatal(err)
		}
		if code != exitFatal {
			t.Fatalf("code %d, want exitFatal; stderr %q", code, errw.String())
		}
		if !strings.Contains(errw.String(), "another session new/mint/revoke is writing this credentials store") {
			t.Fatalf("stderr %q lacks the writer-busy text", errw.String())
		}
		if rows := approverRows(t, creds); !reflect.DeepEqual(rows, before) {
			t.Fatalf("rows changed: %+v, want %+v", rows, before)
		}
	})
}

// assertRevokeRefusal runs argv under env, wants code and every substring in stderr, and wants
// nothing on stdout.
func assertRevokeRefusal(t *testing.T, argv []string, env sessionEnv, code int, want ...string) {
	t.Helper()
	var out, errw bytes.Buffer
	if got := runApprover(argv, &out, &errw, env); got != code {
		t.Fatalf("code %d, want %d; stderr %q", got, code, errw.String())
	}
	for _, w := range want {
		if !strings.Contains(errw.String(), w) {
			t.Fatalf("stderr %q lacks %q", errw.String(), w)
		}
	}
	if out.Len() != 0 {
		t.Fatalf("stdout %q, want empty", out.String())
	}
}

// assertRevokeRace seeds the write seam with err, answers y, and wants exitUsage with the
// race text and no row change.
func assertRevokeRace(t *testing.T, err error) {
	t.Helper()
	db, creds, before := revokeRig(t)
	orig := approverRevoke
	t.Cleanup(func() { approverRevoke = orig })
	approverRevoke = func(*store.Credentials, context.Context, string, int64) error { return err }
	term := &lockedTerm{in: strings.NewReader("y\n")}
	var out, errw bytes.Buffer
	code := runApprover([]string{"revoke", "mark-mac", "--db", db}, &out, &errw, approveEnv(term, nil, atTime(approveT0+100)))
	if code != exitUsage {
		t.Fatalf("code %d, want exitUsage; stderr %q", code, errw.String())
	}
	if !strings.Contains(errw.String(), "(another revoke ran first); nothing was revoked") {
		t.Fatalf("stderr %q lacks the race text", errw.String())
	}
	if rows := approverRows(t, creds); !reflect.DeepEqual(rows, before) {
		t.Fatalf("rows changed: %+v, want %+v", rows, before)
	}
}

// CA10: an unanswered revoke prompt times out after approverConfirmBudget, the writer lock is
// free while it waits, and a missing prompt timer is an internal error that opens nothing.
func TestApproverRevokeUnansweredPromptTimesOut(t *testing.T) {
	db, creds, before := revokeRig(t)
	pr, pw := io.Pipe()
	t.Cleanup(func() { _ = pw.Close() })
	term := &lockedTerm{in: pr}
	fire := make(chan time.Time)
	var mu sync.Mutex
	var waited time.Duration
	env := sessionEnv{openTerminal: func() (io.ReadWriteCloser, error) { return term, nil }, now: atTime(approveT0 + 100),
		after: func(d time.Duration) <-chan time.Time { mu.Lock(); waited = d; mu.Unlock(); return fire }}
	type res struct {
		code   int
		stderr string
	}
	done := make(chan res, 1)
	go func() {
		var out, errw bytes.Buffer
		code := runApprover([]string{"revoke", "mark-mac", "--db", db}, &out, &errw, env)
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
		t.Fatal("revoke did not return within 5 s of the timer")
	}
	if r.code != exitUsage {
		t.Fatalf("code %d, want exitUsage; stderr %q", r.code, r.stderr)
	}
	if want := "ailang-worldd approver revoke: no answer within 2m0s; nothing was revoked\n"; r.stderr != want {
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
	if rows := approverRows(t, creds); !reflect.DeepEqual(rows, before) {
		t.Fatalf("rows changed: %+v, want %+v (RevokedAt must stay 0)", rows, before)
	}

	// Arm 2: no prompt timer at all is an internal error, and no terminal is opened.
	var opened bool
	var out, errw bytes.Buffer
	noTimer := approveEnv(nil, &opened, atTime(approveT0+100))
	noTimer.after = nil
	if code := runApprover([]string{"revoke", "mark-mac", "--db", db}, &out, &errw, noTimer); code != exitFatal {
		t.Fatalf("no timer: code %d, want exitFatal; stderr %q", code, errw.String())
	}
	if !strings.Contains(errw.String(), "internal error: no prompt timer (sessionEnv.after)") {
		t.Fatalf("no timer: stderr %q lacks the internal error", errw.String())
	}
	if opened {
		t.Fatal("no timer: the terminal was opened")
	}
	if rows := approverRows(t, creds); !reflect.DeepEqual(rows, before) {
		t.Fatalf("no timer: rows changed: %+v", rows)
	}
}

// ---- row 201 A2 sprint 1 M0 (K8-ERRCODE): K8's read failing after a committed revoke is
// exitFatal (a failure after the write), with the revoke already printed.

// assertRevokedThenFatal runs a confirmed revoke and requires K7's line, exitFatal and want.
func assertRevokedThenFatal(t *testing.T, db, want string) {
	t.Helper()
	term := &lockedTerm{in: strings.NewReader("y\n")}
	var out, errw bytes.Buffer
	code := runApprover([]string{"revoke", "mark-mac", "--db", db}, &out, &errw, approveEnv(term, nil, atTime(approveT0+100)))
	if code != exitFatal {
		t.Fatalf("code %d, want exitFatal; stderr %q", code, errw.String())
	}
	if out.String() != "revoked approver mark-mac\n" {
		t.Fatalf("stdout %q, want only K7's line", out.String())
	}
	if !strings.Contains(errw.String(), want) || !strings.Contains(errw.String(), "fix: ailang-worldd approve --list --decided") {
		t.Fatalf("stderr %q lacks %q and the fix", errw.String(), want)
	}
}

// K8a: the real revoke commits, then the credentials file is moved away, so K8's read-only
// reopen fails.
func TestApproverRevokeListingReopenFailureIsFatal(t *testing.T) {
	db, creds, _ := revokeRig(t)
	orig := approverRevoke
	t.Cleanup(func() { approverRevoke = orig })
	approverRevoke = func(c *store.Credentials, ctx context.Context, name string, at int64) error {
		if err := orig(c, ctx, name, at); err != nil {
			return err
		}
		return os.Rename(creds, creds+".away")
	}
	assertRevokedThenFatal(t, db, "revoked, but cannot read the credentials store")
	if _, err := os.Stat(creds + ".away"); err != nil {
		t.Fatalf("the seam did not run: %v", err)
	}
}

// K8b: K8's listing query fails after the revoke committed.
func TestApproverRevokeListingQueryFailureIsFatal(t *testing.T) {
	db, creds, _ := revokeRig(t)
	orig := approverDecided
	t.Cleanup(func() { approverDecided = orig })
	approverDecided = func(*store.Credentials, context.Context, string, int) ([]store.MintRequestRow, error) {
		return nil, errors.New("disk on fire")
	}
	assertRevokedThenFatal(t, db, "revoked, but cannot list what it approved: disk on fire")
	if rows := approverRows(t, creds); len(rows) != 1 || rows[0].RevokedAt != approveT0+100 {
		t.Fatalf("rows %+v, want mark-mac revoked at %d", rows, approveT0+100)
	}
}
