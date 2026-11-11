package main

// Row 201 A2 sprint 3 M6g (design_docs/planned/w-world-approvals-by-key-a2.md §2.13, CW8): approve
// --webauthn-assertion refuses before any terminal and any network call, and F1 takes one of five
// modes. Every store is under t.TempDir(); every passkey is approvaltest.Passkey's fixed seed
// (D-WORLD-109).

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
)

// CW8: each refusal of --webauthn-assertion (WF2-WF6, and admission before the file is read) exits
// exitUsage with its text, opens no terminal, makes no daemon call, and leaves the request pending,
// no session and the counter at 1.
func TestApprovePasskeyRefusalsBeforeAnyEffect(t *testing.T) {
	good := func(t *testing.T, a approveRig) []byte {
		return passkeySigner.Assert(pkOrigin, pkRPID, requestHashOf(t, a), 2)
	}
	sql := func(stmts ...string) func(t *testing.T, a approveRig) {
		return func(t *testing.T, a approveRig) {
			raw := rawCredentials(t, a.credPath)
			defer func() { _ = raw.Close() }()
			for _, s := range stmts {
				if _, err := raw.Exec(s); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
	for _, arm := range []struct {
		name  string
		setup func(t *testing.T, a approveRig)
		file  func(t *testing.T, a approveRig) []byte // nil: path makes the path
		path  func(t *testing.T) string
		now   int64 // 0: approveT0+1
		want  string
	}{
		{name: "WF2-fifo", path: func(t *testing.T) string {
			p := filepath.Join(t.TempDir(), "assertion.fifo")
			if err := syscall.Mkfifo(p, 0o600); err != nil {
				t.Fatal(err)
			}
			return p
		}, want: "assertion.fifo is not a regular file"},
		{name: "WF2-over-16k", file: func(*testing.T, approveRig) []byte { return bytes.Repeat([]byte(" "), 16<<10+1) },
			want: "is 16385 bytes, over the 16384-byte bound"},
		{name: "WF2-missing", path: func(t *testing.T) string { return filepath.Join(t.TempDir(), "absent.json") }, want: "--webauthn-assertion: open "},
		{name: "WF3-not-an-assertion", file: func(*testing.T, approveRig) []byte { return []byte(`{"id":"x"}`) },
			want: "is not a passkey assertion: " + approval.ErrPasskeyResponse.Error()},
		{name: "WF4-v5-store", setup: sql(`DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`), file: good,
			want: "passkeys need schema 6"},
		{name: "WF5-unknown-credential", file: func(t *testing.T, a approveRig) []byte {
			return approvaltest.NewPasskey(approval.AlgEdDSA, 55).Assert(pkOrigin, pkRPID, requestHashOf(t, a), 2)
		}, want: "the assertion's credential cred " + passkeyCred(approvaltest.NewPasskey(approval.AlgEdDSA, 55).CredentialID()) + " is not an enrolled passkey"},
		{name: "WF5-stored-key-invalid", setup: sql(`UPDATE approvers SET public_key = x'a0' WHERE name = 'mark-mac'`), file: good,
			want: `approver "mark-mac"'s stored passkey key is not valid`},
		{name: "WF6-counter-regression", file: func(t *testing.T, a approveRig) []byte {
			return passkeySigner.Assert(pkOrigin, pkRPID, requestHashOf(t, a), 1)
		}, want: approval.ErrSignCountRegression.Error() + ": enrolled 1, this assertion 1"},
		{name: "WF6-revoked", setup: func(t *testing.T, a approveRig) { revokeSeeded(t, a.credPath, "mark-mac", approveT0-50) }, file: good,
			want: approval.ErrApproverRevoked.Error()},
		{name: "WF6-enrolled-after-now", setup: sql(`UPDATE approvers SET enrolled_at = 1800000050 WHERE name = 'mark-mac'`), file: good,
			want: approval.ErrApproverClockBehind.Error()},
		{name: "WF6-other-request", file: func(t *testing.T, a approveRig) []byte {
			h := requestHashOf(t, a)
			h[31] ^= 1
			return passkeySigner.Assert(pkOrigin, pkRPID, h, 2)
		}, want: approval.ErrChallenge.Error()},
		{name: "WF6-backup-eligible", setup: sql(`UPDATE webauthn_credentials SET backup_eligible = 1 WHERE name = 'mark-mac'`), file: good,
			want: approval.ErrBackupEligibility.Error()},
		{name: "WF6-another-origin", setup: sql(`UPDATE webauthn_credentials SET origin = 'http://localhost:7645' WHERE name = 'mark-mac'`), file: good,
			want: approval.ErrOriginMismatch.Error()},
		{name: "expired-first", path: func(t *testing.T) string { return filepath.Join(t.TempDir(), "absent.json") }, now: approveT0 + 600, want: "expired at"},
	} {
		t.Run(arm.name, func(t *testing.T) {
			a := passkeyRig(t)
			if arm.setup != nil {
				arm.setup(t, a)
			}
			var path string
			if arm.path != nil {
				path = arm.path(t)
			} else {
				path = writeAssertion(t, arm.file(t, a))
			}
			calls := countDaemonCalls(t)
			now := arm.now
			if now == 0 {
				now = approveT0 + 1
			}
			var opened bool
			var stdout, stderr bytes.Buffer
			code := runApprove(a.args("--webauthn-assertion", path), &stdout, &stderr, approveEnv(nil, &opened, atTime(now)))
			if code != exitUsage || opened || *calls != 0 {
				t.Fatalf("rc %d, terminal opened %v, daemon calls %d; want exitUsage, none, 0\nstderr %s", code, opened, *calls, stderr.String())
			}
			if !strings.Contains(stderr.String(), arm.want) {
				t.Fatalf("stderr %q lacks %q", stderr.String(), arm.want)
			}
			if row := a.row(t); row.Consumed || len(a.sessionRows(t)) != 0 {
				t.Fatal("a refused assertion consumed the request or minted")
			}
			if arm.name != "WF4-v5-store" {
				if p := passkeyOf(t, a.credPath, "mark-mac"); p.SignCount != 1 {
					t.Fatalf("sign count %d, want 1", p.SignCount)
				}
			}
		})
	}
}

// CW8 (F1, F2): any two of the five modes, and --list with a passkey mode, print the help and exit
// exitUsage before any store is read.
func TestApprovePasskeyFlagRules(t *testing.T) {
	db := dbPathOf(t)
	modes := [][]string{{"--signature", "s"}, {"--sign-with", "k"}, {"--write-message", "m"}, {"--webauthn-link"}, {"--webauthn-assertion", "f"}}
	var cases [][]string
	for i := range modes {
		for j := i + 1; j < len(modes); j++ {
			if i < 3 && j < 3 {
				continue // TestApproveFlagRules has the three SSH pairs
			}
			cases = append(cases, append(append(append([]string{firstRequestID}, modes[i]...), modes[j]...), "--db", db))
		}
	}
	cases = append(cases, []string{"--list", "--webauthn-link", "--db", db}, []string{"--list", "--webauthn-assertion", "f", "--db", db})
	if len(cases) != 9 {
		t.Fatalf("%d cases, want 7 pairs and 2 --list arms", len(cases))
	}
	for _, args := range cases {
		var stdout, stderr bytes.Buffer
		if code := runApprove(args, &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0))); code != exitUsage || stderr.String() != approveHelp {
			t.Fatalf("%q: rc %d, stderr %q; want exitUsage and exactly the help", args, code, stderr.String())
		}
		if _, err := os.Lstat(filepath.Join(filepath.Dir(db), "credentials.db")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("%q: a credentials store was touched (err %v)", args, err)
		}
	}
}
