package main

// Row 201 A3 sprint 1 M3b (design_docs/planned/w-world-approvals-by-key-a3.md §2.6, §6 AW2-AW4,
// AW6, AW7): the minter re-judges every recorded approval, stops at the request's expiry, at
// --await-for and when the request is used elsewhere, and mints an SSH signature as ssh:<name>.

import (
	"bytes"
	"net"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// AW2: an assertion whose counter did not move is refused once (per received_at, however many
// polls see it), the wait goes on, and a replacement recorded later mints.
func TestApproveAwaitRejudgesAndWaits(t *testing.T) {
	a := passkeyRig(t)
	h := requestHashOf(t, a)
	env, c := awaitEnv(approveT0+1, nil, func(n, now int64) {
		switch n {
		case 1:
			recordSurface(t, a, "webauthn", "webauthn:mark-mac", passkeySigner.Assert(pkOrigin, pkRPID, h, 1), now)
		case 4:
			recordSurface(t, a, "webauthn", "webauthn:mark-mac", passkeySigner.Assert(pkOrigin, pkRPID, h, 2), now)
		}
	})
	var opened bool
	env.openTerminal = approveEnv(nil, &opened, nil).openTerminal
	var stdout, stderr bytes.Buffer
	if code := runApprove(a.args("--await"), &stdout, &stderr, env); code != exitOK || opened || c.waits != 4 {
		t.Fatalf("rc %d, terminal %v, %d waits; want exitOK, no terminal, 4 waits\nstderr %s", code, opened, c.waits, stderr.String())
	}
	want := "ailang-worldd approve: refused the surface's approval: the passkey assertion does not approve request " + a.id + ": " +
		approval.ErrSignCountRegression.Error() + ": enrolled 1, this assertion 1; if you did not cause this, run: ailang-worldd approver revoke mark-mac; waiting for another\n"
	if !strings.HasPrefix(stderr.String(), want) || strings.Count(stderr.String(), "refused the surface's approval") != 1 {
		t.Fatalf("stderr %q, want exactly one refusal, first: %q", stderr.String(), want)
	}
	requireMintedBy(t, a, "webauthn:mark-mac", stdout.String())
}

// AW3: with nothing recorded the wait ends exactly at the request's expiry (admission's text) or
// exactly at --await-for, exitUsage, nothing minted, the request still pending in the second case.
func TestApproveAwaitStopsAtExpiry(t *testing.T) {
	t.Run("expiry", func(t *testing.T) {
		a := newApproveRig(t, "ep1", approveT0)
		row := a.row(t)
		env, c := awaitEnv(approveT0+1, nil, nil)
		var stdout, stderr bytes.Buffer
		if code := runApprove(a.args("--await"), &stdout, &stderr, env); code != exitUsage {
			t.Fatalf("rc %d, want exitUsage\nstderr %s", code, stderr.String())
		}
		if c.now != row.ExpiresAt || c.waits != row.ExpiresAt-approveT0-1 {
			t.Fatalf("stopped at %d after %d waits, want exactly the expiry %d after %d", c.now, c.waits, row.ExpiresAt, row.ExpiresAt-approveT0-1)
		}
		if want := "ailang-worldd approve: " + expiredRefusal(row) + "\n"; stderr.String() != want {
			t.Fatalf("stderr %q, want %q", stderr.String(), want)
		}
		if a.row(t).Consumed || len(a.sessionRows(t)) != 0 {
			t.Fatal("an expired wait consumed the request or minted")
		}
	})
	t.Run("await-for", func(t *testing.T) {
		a := newApproveRig(t, "ep1", approveT0)
		env, c := awaitEnv(approveT0+1, nil, nil)
		var stdout, stderr bytes.Buffer
		if code := runApprove(a.args("--await", "--await-for", "30s"), &stdout, &stderr, env); code != exitUsage {
			t.Fatalf("rc %d, want exitUsage\nstderr %s", code, stderr.String())
		}
		if c.now != approveT0+31 || c.waits != 30 {
			t.Fatalf("stopped at %d after %d waits, want %d after 30", c.now, c.waits, approveT0+31)
		}
		want := "ailang-worldd approve: no approval for request " + a.id + " by " + approverTime(approveT0+31) +
			"; it stays pending until " + approverTime(a.row(t).ExpiresAt) + "\n"
		if stderr.String() != want {
			t.Fatalf("stderr %q, want %q", stderr.String(), want)
		}
		if a.row(t).Consumed || len(a.sessionRows(t)) != 0 {
			t.Fatal("a timed-out wait consumed the request or minted")
		}
	})
}

// AW4: a consume by another process mid-wait ends the wait with admission's "already used" text.
func TestApproveAwaitStopsWhenUsedElsewhere(t *testing.T) {
	a := newApproveRig(t, "ep1", approveT0)
	env, c := awaitEnv(approveT0+1, nil, func(n, now int64) {
		if n == 2 {
			raw := rawCredentials(t, a.credPath)
			defer func() { _ = raw.Close() }()
			if _, err := raw.Exec(`UPDATE mint_requests SET consumed_at = ?, approved_by = 'tty', credential_id = 'elsewhere' WHERE request_id = ?`, now, a.id); err != nil {
				t.Error(err)
			}
		}
	})
	var stdout, stderr bytes.Buffer
	if code := runApprove(a.args("--await"), &stdout, &stderr, env); code != exitUsage || c.waits != 2 {
		t.Fatalf("rc %d after %d waits, want exitUsage after 2\nstderr %s", code, c.waits, stderr.String())
	}
	if !strings.Contains(stderr.String(), "was already used at "+approverTime(approveT0+3)+" by tty (credential elsewhere)") {
		t.Fatalf("stderr %q lacks the already-used refusal", stderr.String())
	}
	if len(a.sessionRows(t)) != 0 {
		t.Fatal("the minter minted a request used elsewhere")
	}
}

// AW6: a recorded SSHSIG by an enrolled key mints as ssh:mark-mac, with no terminal; with no
// passkey enrolled the card's only address is the request's loopback daemon.
func TestApproveAwaitSSHSIG(t *testing.T) {
	a, _ := signedRig(t)
	sig := sigSigner.Sign(approval.Namespace, requestMessage(t, a))
	var opened bool
	env, _ := awaitEnv(approveT0+1, &opened, func(n, now int64) {
		if n == 1 {
			recordSurface(t, a, "ssh", "ssh:mark-mac", sig, now)
		}
	})
	var stdout, stderr bytes.Buffer
	if code := runApprove(a.args("--await"), &stdout, &stderr, env); code != exitOK || opened {
		t.Fatalf("rc %d, terminal opened %v; want exitOK with no terminal\nstderr %s", code, opened, stderr.String())
	}
	requireMintedBy(t, a, "ssh:mark-mac", stdout.String())
	mr, err := approval.Decode(a.row(t).Request)
	if err != nil {
		t.Fatal(err)
	}
	_, port, err := net.SplitHostPort(mr.Daemon.Addr)
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(stdout.String(), "/surface/approve/"); n != 1 || !strings.Contains(stdout.String(), "\n  open http://localhost:"+port+"/surface/approve/"+a.id+"\n") {
		t.Fatalf("want exactly the loopback card address:\n%s", stdout.String())
	}
}

// AW8 (planner): a store below schema 7 ends the wait at the first poll naming the serve restart
// (below schema 6, before it, naming the passkeys' restart), and a request deleted mid-wait ends
// it naming that; nothing is minted.
func TestApproveAwaitNeedsSchema7AndItsRequest(t *testing.T) {
	for _, arm := range []struct {
		name  string
		setup []string // before the run
		gone  bool     // delete the request after the first wait
		want  string
	}{
		{name: "v6-store", setup: []string{`DROP TABLE surface_approvals`, `PRAGMA user_version = 6`}, want: "surface approvals need schema 7: restart `ailang-worldd serve`"},
		{name: "v5-store", setup: []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`}, want: "passkeys need schema 6"},
		{name: "deleted", gone: true, want: "is no longer in the store"},
	} {
		t.Run(arm.name, func(t *testing.T) {
			a := newApproveRig(t, "ep1", approveT0)
			raw := rawCredentials(t, a.credPath)
			for _, s := range arm.setup {
				if _, err := raw.Exec(s); err != nil {
					t.Fatal(err)
				}
			}
			_ = raw.Close()
			var opened bool
			env, c := awaitEnv(approveT0+1, &opened, func(n, now int64) {
				if n == 1 && arm.gone {
					raw := rawCredentials(t, a.credPath)
					defer func() { _ = raw.Close() }()
					if _, err := raw.Exec(`DELETE FROM mint_requests WHERE request_id = ?`, a.id); err != nil {
						t.Error(err)
					}
				}
			})
			var stdout, stderr bytes.Buffer
			code := runApprove(a.args("--await"), &stdout, &stderr, env)
			if code != exitUsage || opened || !strings.Contains(stderr.String(), arm.want) {
				t.Fatalf("rc %d, terminal %v, stderr %q; want exitUsage naming %q", code, opened, stderr.String(), arm.want)
			}
			if want := map[bool]int64{false: 0, true: 1}[arm.gone]; c.waits != want {
				t.Fatalf("%d waits, want %d", c.waits, want)
			}
			if len(a.sessionRows(t)) != 0 {
				t.Fatal("minted")
			}
		})
	}
}

// AW7: a row the minter did not get from a verified POST — an unenrolled key's valid signature, a
// passkey assertion over another request, an approver column naming mark-mac over garbage — is
// judged, refused once, and never mints: no terminal, no session, the request pending.
func TestApproveAwaitNeverTrustsTheRow(t *testing.T) {
	for _, arm := range []struct {
		name, format, approver string
		artifact               func(t *testing.T, a approveRig) []byte
		want                   string
	}{
		{"unenrolled-key", "ssh", "ssh:mark-mac", func(t *testing.T, a approveRig) []byte {
			return otherSigner.Sign(approval.Namespace, requestMessage(t, a))
		}, "the signature's key " + otherSigner.PublicKey().Fingerprint() + " is not an enrolled approver"},
		{"other-request", "webauthn", "webauthn:mark-mac", func(t *testing.T, a approveRig) []byte {
			h := requestHashOf(t, a)
			h[31] ^= 1
			return passkeySigner.Assert(pkOrigin, pkRPID, h, 2)
		}, "the passkey assertion does not approve request "},
		{"garbage-named-mark", "webauthn", "webauthn:mark-mac", func(*testing.T, approveRig) []byte { return []byte(`{"id":"x"}`) },
			"the surface's approval is not a passkey assertion: " + approval.ErrPasskeyResponse.Error()},
	} {
		t.Run(arm.name, func(t *testing.T) {
			a := passkeyRig(t)
			seedApprover(t, a.credPath, "mark-ssh", sigSigner.PublicKey(), approveT0-100)
			art := arm.artifact(t, a)
			var opened bool
			env, _ := awaitEnv(approveT0+1, &opened, func(n, now int64) {
				if n == 1 {
					recordSurface(t, a, arm.format, arm.approver, art, now)
				}
			})
			var stdout, stderr bytes.Buffer
			code := runApprove(a.args("--await", "--await-for", "5s"), &stdout, &stderr, env)
			if code != exitUsage || opened {
				t.Fatalf("rc %d, terminal opened %v; want exitUsage, no terminal\nstderr %s", code, opened, stderr.String())
			}
			if n := strings.Count(stderr.String(), "refused the surface's approval: "); n != 1 || !strings.Contains(stderr.String(), arm.want) ||
				!strings.Contains(stderr.String(), "no approval for request "+a.id) {
				t.Fatalf("stderr %q: want one refusal naming %q, then the timeout", stderr.String(), arm.want)
			}
			if a.row(t).Consumed || len(a.sessionRows(t)) != 0 {
				t.Fatal("the minter trusted the row")
			}
		})
	}
}
