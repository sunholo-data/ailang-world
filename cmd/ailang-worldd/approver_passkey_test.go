package main

// Row 201 A2 sprint 3 M6c (design_docs/planned/w-world-approvals-by-key-a2.md §2.12, CW1, CW3):
// approver enroll --format webauthn refuses before any terminal. Every store is under t.TempDir();
// every passkey is approvaltest.Passkey's fixed seed (D-WORLD-109).

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
)

// enrolRig is a fresh store and the (absent) path the browser saves the registration to.
type enrolRig struct{ db, creds, reg string }

func newEnrolRig(t *testing.T) enrolRig {
	t.Helper()
	db, creds := approverStore(t)
	return enrolRig{db: db, creds: creds, reg: filepath.Join(t.TempDir(), "world-registration-mark-mac.json")}
}

// args is approver enroll --name mark-mac --format webauthn --origin pkOrigin --registration reg --db db.
func (r enrolRig) args() []string {
	return []string{"enroll", "--name", "mark-mac", "--format", "webauthn", "--origin", pkOrigin, "--registration", r.reg, "--db", r.db}
}

// passkeyCount is the number of webauthn_credentials rows in creds.
func passkeyCount(t *testing.T, creds string) int {
	t.Helper()
	raw := rawCredentials(t, creds)
	defer func() { _ = raw.Close() }()
	var n int
	if err := raw.QueryRow(`SELECT count(*) FROM sqlite_master WHERE name = 'webauthn_credentials'`).Scan(&n); err != nil || n == 0 {
		return 0 // a schema-5 store has none
	}
	if err := raw.QueryRow(`SELECT count(*) FROM webauthn_credentials`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// linkFragment decodes the fragment of the first ceremony link in out.
func linkFragment(t *testing.T, out string) ceremonyLink {
	t.Helper()
	_, rest, ok := strings.Cut(out, "/surface/passkey#")
	if !ok {
		t.Fatalf("no ceremony link in %q", out)
	}
	frag, _, _ := strings.Cut(rest, "\n")
	b, err := base64.RawURLEncoding.Strict().DecodeString(frag)
	if err != nil {
		t.Fatal(err)
	}
	var c ceremonyLink
	if err := json.Unmarshal(b, &c); err != nil {
		t.Fatal(err)
	}
	return c
}

// challengeOf is a fragment's challenge.
func challengeOf(t *testing.T, c ceremonyLink) [32]byte {
	t.Helper()
	b, err := base64.RawURLEncoding.Strict().DecodeString(c.Challenge)
	if err != nil || len(b) != 32 {
		t.Fatalf("challenge %q: %v", c.Challenge, err)
	}
	return [32]byte(b)
}

// codeOf is the grouped enrolment code of a registration file, as the page shows it.
func codeOf(t *testing.T, file []byte) string {
	t.Helper()
	r, err := approval.ParseRegistrationResponse(file)
	if err != nil {
		t.Fatal(err)
	}
	return approval.GroupEnrolmentCode(approval.EnrolmentCode(r.ClientDataJSON, r.AttestationObject))
}

// enrolEnv is approveEnv with seqRandom as the challenge source.
func enrolEnv(term io.ReadWriteCloser, opened *bool) sessionEnv {
	env := approveEnv(term, opened, atTime(approveT0))
	env.random = seqRandom()
	return env
}

// CW1, CW3: every refusal before the terminal is exitUsage, names its reason, opens no terminal,
// prints no link and enrols nothing; each arm runs under a 5 s guard. Without a terminal, enrol
// refuses with the tty fix and prints no link.
func TestPasskeyEnrollFlagRules(t *testing.T) {
	ed := fixturePub("ed25519.pub")
	for _, arm := range []struct {
		name  string
		args  func(t *testing.T, r enrolRig) []string
		setup func(t *testing.T, r enrolRig)
		want  []string
		rows  int // approvers seeded
		pks   int // passkeys seeded
	}{
		{name: "WN1-webauthn-pubkey", args: func(t *testing.T, r enrolRig) []string { return append(r.args(), "--pubkey", ed) },
			want: []string{"--pubkey is for --format ssh; --origin and --registration are for --format webauthn", "usage: ailang-worldd approver"}},
		{name: "WN1-ssh-origin", args: func(t *testing.T, r enrolRig) []string {
			return append(enrollArgs(r.db, "mark-mac", ed), "--origin", pkOrigin)
		}, want: []string{"--origin and --registration are for --format webauthn"}},
		{name: "WN1-ssh-registration", args: func(t *testing.T, r enrolRig) []string {
			return append(enrollArgs(r.db, "mark-mac", ed), "--registration", r.reg)
		}, want: []string{"--origin and --registration are for --format webauthn"}},
		{name: "WN2-name", args: func(t *testing.T, r enrolRig) []string { return dropFlag(r.args(), "--name") },
			want: []string{"--name, --origin and --registration are required with --format webauthn"}},
		{name: "WN2-origin", args: func(t *testing.T, r enrolRig) []string { return dropFlag(r.args(), "--origin") },
			want: []string{"--name, --origin and --registration are required with --format webauthn"}},
		{name: "WN2-registration", args: func(t *testing.T, r enrolRig) []string { return dropFlag(r.args(), "--registration") },
			want: []string{"--name, --origin and --registration are required with --format webauthn"}},
		{name: "WN2a-name-grammar", args: func(t *testing.T, r enrolRig) []string { return replaceFlag(r.args(), "--name", "Mark") },
			want: []string{`--name "Mark": approval: approver name must match`}},
		{name: "N3-no-world-store", args: func(t *testing.T, r enrolRig) []string {
			return replaceFlag(r.args(), "--db", filepath.Join(t.TempDir(), "absent.db"))
		}, want: []string{"does not exist; refusing to create one"}},
		{name: "WN3-ip", args: func(t *testing.T, r enrolRig) []string {
			return replaceFlag(r.args(), "--origin", "http://127.0.0.1:7644")
		},
			want: []string{`--origin "http://127.0.0.1:7644": `, "use http://localhost:<port>"}},
		{name: "WN3-other-host", args: func(t *testing.T, r enrolRig) []string {
			return replaceFlag(r.args(), "--origin", "https://evil.example")
		}, want: []string{`--origin "https://evil.example": `, "localhost or a tailnet *.ts.net name"}},
		{name: "WN4-exists", setup: func(t *testing.T, r enrolRig) {
			if err := os.WriteFile(r.reg, []byte("{}"), 0o600); err != nil {
				t.Fatal(err)
			}
		}, args: func(t *testing.T, r enrolRig) []string { return r.args() }, want: []string{"world-registration-mark-mac.json already exists or cannot be checked; pass a new path"}},
		{name: "WN4-dangling-symlink", setup: func(t *testing.T, r enrolRig) {
			if err := os.Symlink(filepath.Join(t.TempDir(), "absent"), r.reg); err != nil {
				t.Fatal(err)
			}
		}, args: func(t *testing.T, r enrolRig) []string { return r.args() }, want: []string{"already exists or cannot be checked"}},
		{name: "WN4-unsearchable-dir", args: func(t *testing.T, r enrolRig) []string {
			dir := filepath.Join(t.TempDir(), "locked")
			if err := os.Mkdir(dir, 0o000); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })
			return replaceFlag(r.args(), "--registration", filepath.Join(dir, "r.json"))
		}, want: []string{"r.json already exists or cannot be checked"}},
		{name: "N6-v5", setup: func(t *testing.T, r enrolRig) {
			raw := rawCredentials(t, r.creds)
			defer func() { _ = raw.Close() }()
			for _, stmt := range []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
				if _, err := raw.Exec(stmt); err != nil {
					t.Fatal(err)
				}
			}
		}, args: func(t *testing.T, r enrolRig) []string { return r.args() }, want: []string{"passkeys need schema 6", "restart `ailang-worldd serve`"}},
		{name: "N7-name-taken", rows: 1, setup: func(t *testing.T, r enrolRig) {
			seedApprover(t, r.creds, "mark-mac", fixtureKeyOf(t, "ed25519.pub"), approveT0-10)
		}, args: func(t *testing.T, r enrolRig) []string { return r.args() }, want: []string{`an approver named "mark-mac" already exists`}},
		{name: "N7-name-taken-revoked-passkey", rows: 1, pks: 1, setup: func(t *testing.T, r enrolRig) {
			seedPasskey(t, r.creds, "mark-mac", approvaltest.NewPasskey(approval.AlgEdDSA, 31), approveT0-10)
			revokeSeeded(t, r.creds, "mark-mac", approveT0-5)
		}, args: func(t *testing.T, r enrolRig) []string { return r.args() }, want: []string{`an approver named "mark-mac" already exists`}},
	} {
		t.Run(arm.name, func(t *testing.T) {
			r := newEnrolRig(t)
			if arm.setup != nil {
				arm.setup(t, r)
			}
			args := arm.args(t, r)
			var out, errw bytes.Buffer
			var opened bool
			done := make(chan int, 1)
			go func() { done <- runApprover(args, &out, &errw, enrolEnv(nil, &opened)) }()
			var code int
			select {
			case code = <-done:
			case <-time.After(5 * time.Second):
				t.Fatal("enrol did not refuse within 5 s")
			}
			if code != exitUsage || opened || strings.Contains(out.String(), "/surface/passkey") {
				t.Fatalf("code %d, terminal opened %v, stdout %q; want exitUsage, none, no link\nstderr %q", code, opened, out.String(), errw.String())
			}
			for _, w := range arm.want {
				if !strings.Contains(errw.String(), w) {
					t.Fatalf("stderr %q lacks %q", errw.String(), w)
				}
			}
			if n, m := len(approverRows(t, r.creds)), passkeyCount(t, r.creds); n != arm.rows || m != arm.pks {
				t.Fatalf("%d approvers and %d passkeys, want %d and %d (the seeds)", n, m, arm.rows, arm.pks)
			}
		})
	}

	t.Run("CW3-no-terminal", func(t *testing.T) {
		r := newEnrolRig(t)
		env := noTerminalEnv(approveT0)
		env.after, env.random = time.After, seqRandom()
		var out, errw bytes.Buffer
		if code := runApprover(r.args(), &out, &errw, env); code != exitUsage || out.Len() != 0 {
			t.Fatalf("code %d, stdout %q; want exitUsage and no link", code, out.String())
		}
		for _, w := range []string{"refusing: no controlling terminal (", "nothing was enrolled", sessionTTYFix} {
			if !strings.Contains(errw.String(), w) {
				t.Fatalf("stderr %q lacks %q", errw.String(), w)
			}
		}
		if n, m := len(approverRows(t, r.creds)), passkeyCount(t, r.creds); n != 0 || m != 0 {
			t.Fatalf("%d approvers, %d passkeys; want 0", n, m)
		}
	})
}

// dropFlag removes flag and its value from args.
func dropFlag(args []string, flag string) []string {
	var out []string
	for i := 0; i < len(args); i++ {
		if args[i] == flag {
			i++
			continue
		}
		out = append(out, args[i])
	}
	return out
}

// replaceFlag sets flag's value in args.
func replaceFlag(args []string, flag, value string) []string {
	out := append([]string(nil), args...)
	for i := range out {
		if out[i] == flag {
			out[i+1] = value
		}
	}
	return out
}
