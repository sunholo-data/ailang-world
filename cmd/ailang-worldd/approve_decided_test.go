package main

// Row 201 A1b sprint 4 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint4-plan.md §3
// D2, D3): approve --list --decided, and what approver revoke lists. Every store is under
// t.TempDir(); every key is approvaltest's fixed seed (D-WORLD-109).

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// decidedRig is one store with three requests: a approved by the typed y at +1, b approved by
// mark-mac's signature at +2, and c still pending.
func decidedRig(t *testing.T) (a, b, c approveRig) {
	t.Helper()
	r := newEpisodeRun(t)
	a = newApproveRigIn(t, r, "ep1", approveT0)
	b = newApproveRigIn(t, r, "ep2", approveT0)
	c = newApproveRigIn(t, r, "ep3", approveT0)
	seedApprover(t, b.credPath, "mark-mac", sigSigner.PublicKey(), approveT0-100)
	var stdout, stderr bytes.Buffer
	if code := runApprove(a.args(), &stdout, &stderr, approveEnv(&promptTerm{in: strings.NewReader("y\n")}, nil, atTime(approveT0+1))); code != exitOK {
		t.Fatalf("typed approve of a: rc %d\nstderr %s", code, stderr.String())
	}
	sig := writeSignature(t, sigSigner.Sign(approval.Namespace, requestMessage(t, b)))
	if code := runApprove(b.args("--signature", sig), &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0+2))); code != exitOK {
		t.Fatalf("signed approve of b: rc %d\nstderr %s", code, stderr.String())
	}
	return a, b, c
}

// CR7: --list --decided prints the used requests only, newest first, each with who approved it,
// when, and its credential's first 16 hex (F16); with none it says so (F15).
func TestApproveDecidedList(t *testing.T) {
	a, b, _ := decidedRig(t)
	ar, br := a.row(t), b.row(t)
	want := b.id + "  episode ep2  approved by ssh:mark-mac  at 2027-01-15T08:00:02Z  credential " + br.CredentialID[:16] + "\n" +
		a.id + "  episode ep1  approved by tty  at 2027-01-15T08:00:01Z  credential " + ar.CredentialID[:16] + "\n"
	var stdout, stderr bytes.Buffer
	args := []string{"--list", "--decided", "--db", a.r.f.db, "--credentials-db", a.credPath}
	if code := runApprove(args, &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0+3))); code != exitOK || stdout.String() != want {
		t.Fatalf("rc %d, stdout\n%s\nwant\n%s\nstderr %s", code, stdout.String(), want, stderr.String())
	}

	t.Run("none", func(t *testing.T) {
		db, creds := approverStore(t)
		var stdout, stderr bytes.Buffer
		code := runApprove([]string{"--list", "--decided", "--db", db, "--credentials-db", creds}, &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0)))
		if code != exitOK || stdout.String() != "no decided mint requests\n" {
			t.Fatalf("rc %d, stdout %q; want exitOK and %q\nstderr %s", code, stdout.String(), "no decided mint requests\n", stderr.String())
		}
	})
}

// ---- D3 ----

// CR12: approver revoke lists the sessions the approver approved, each with a session revoke line
// that works, and revokes none of them (K8, Q4); an approver that approved nothing says so (K9).
func TestApproverRevokeListsWhatItApproved(t *testing.T) {
	a, b, _ := decidedRig(t)
	seedApprover(t, a.credPath, "idle", otherSigner.PublicKey(), approveT0-100)
	br := b.row(t)
	revoke := func(name string) string {
		t.Helper()
		var stdout, stderr bytes.Buffer
		args := []string{"revoke", name, "--db", a.r.f.db, "--credentials-db", a.credPath}
		if code := runApprover(args, &stdout, &stderr, approveEnv(&lockedTerm{in: strings.NewReader("y\n")}, nil, atTime(approveT0+10))); code != exitOK {
			t.Fatalf("revoke %s: rc %d\nstderr %s", name, code, stderr.String())
		}
		return stdout.String()
	}
	line := "ailang-worldd session revoke --db " + shellQuote(a.r.f.db) + " --credentials-db " + shellQuote(a.credPath) + " " + br.CredentialID
	want := "revoked approver mark-mac\n" +
		"it approved 1 session(s), newest first; each stays valid until it expires or is revoked:\n" +
		"  " + b.id + "  episode ep2  at 2027-01-15T08:00:02Z  credential " + br.CredentialID[:16] + "\n" +
		"    " + line + "\n"
	if got := revoke("mark-mac"); got != want {
		t.Fatalf("stdout\n%s\nwant\n%s", got, want)
	}
	if got, want := revoke("idle"), "revoked approver idle\nno session was approved by it\n"; got != want {
		t.Fatalf("stdout %q, want %q", got, want)
	}
	if n := len(a.sessionRows(t)); n != 2 {
		t.Fatalf("%d session rows after the revokes, want 2: approver revoke revokes no session", n)
	}
	var stdout, stderr bytes.Buffer
	if code := runSessionRevoke(strings.Fields(line)[3:], &stdout, &stderr); code != exitOK {
		t.Fatalf("the printed line %q: rc %d\nstderr %s", line, code, stderr.String())
	}
	rows := a.sessionRows(t)
	if len(rows) != 1 || rows[0].CredentialID == br.CredentialID {
		t.Fatalf("session rows after the printed line: %+v, want only a's", rows)
	}
}
