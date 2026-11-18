package main

// Row 243 V3 sprint 5 M12 (design §2.7, §2.9 P1, §6 CF4–CF9, CF14): the confirm-on-surface parts
// one by one — the canonical argv, the request's refusals, the put's refusals, the wait's every end
// (a fake clock, a wall-clock test deadline), the re-check (a) and the read-back (b). Every store is
// under t.TempDir(); every key is approvaltest's (sigSigner enrolled as mark-mac, otherSigner never);
// an "approved by hand" row is raw SQL, as the design says; world objects are broker's own writes.

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

const sfT0 = int64(1_900_000_000)

// surfaceUnit is a credentials store with sigSigner enrolled as mark-mac, and a real finish input.
func surfaceUnit(t *testing.T) (string, surfaceFinish) {
	t.Helper()
	creds := filepath.Join(t.TempDir(), "credentials.db")
	seedApprover(t, creds, "mark-mac", sigSigner.PublicKey(), sfT0-100)
	scope := broker.FinishApprovalScope{Repository: "fixture/repo", Remote: "https://github.com/fixture/repo.git", Head: "custom-branch",
		SHA: strings.Repeat("a", 40), Base: "dev", TitleSHA256: finishSHA([]byte("Fixture title")), BodySHA256: finishSHA([]byte("body\n")),
		Draft: true, ExpiresAt: sfT0 + 900}
	return creds, surfaceFinish{credentials: creds, addr: "127.0.0.1:47330", ep: "ep1", sha: scope.SHA,
		argv: []string{"ailang-worldd", "episode", "finish", "ep1", "--confirm-on-surface"}, scope: scope,
		scopeText: broker.FormatFinishApprovalScope(scope), title: "Fixture title", body: "body\n"}
}

// storeOwn stores in's request as id with challenge bytes ch at now, through the production put.
func storeOwn(t *testing.T, in surfaceFinish, id string, ch byte, now int64) verbOwn {
	t.Helper()
	_, canon, err := finishVerbRequest(in)
	if err != nil {
		t.Fatal(err)
	}
	own := verbOwn{id: id, ep: in.ep, digest: sha256.Sum256(canon), challenge: [32]byte(bytes.Repeat([]byte{ch}, 32)), expires: now + 600}
	row := store.VerbRequestRow{RequestID: id, Request: canon, Digest: hex.EncodeToString(own.digest[:]), Challenge: hex.EncodeToString(own.challenge[:]),
		EpisodeID: in.ep, CreatedAt: now, ExpiresAt: own.expires}
	if refusal := putFinishVerbRequest(in.credentials, row, now); refusal != "" {
		t.Fatal(refusal)
	}
	return own
}

// verbSig is s's SSHSIG under namespace over the 86-byte message of digest, challenge and expiry.
func verbSig(t *testing.T, s *approvaltest.Signer, namespace string, digest, challenge [32]byte, expiry int64) []byte {
	t.Helper()
	m, err := approval.VerbMessage(digest, challenge, expiry)
	if err != nil {
		t.Fatal(err)
	}
	return s.Sign(namespace, m)
}

// credSQL runs raw statements on the credentials store (the design's "approved by hand").
func credSQL(t *testing.T, creds string, stmts ...string) {
	t.Helper()
	db, err := sql.Open("sqlite", creds)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			t.Fatalf("%s: %v", s, err)
		}
	}
}

// approveByHand marks id approved with artifact, approvedBy and a well-formed decision ref; set
// adds column overrides (raw SQL, as design §6 CF6).
func approveByHand(t *testing.T, creds, id string, artifact []byte, approvedBy string, at int64, set ...string) {
	t.Helper()
	cols := append([]string{fmt.Sprintf("consumed_at = %d", at), "outcome = 'approved'", "approved_by = '" + approvedBy + "'",
		"artifact = X'" + hex.EncodeToString(artifact) + "'", "decision_ref = 'sha256:" + strings.Repeat("d", 64) + "'"}, set...)
	credSQL(t, creds, "UPDATE verb_requests SET "+strings.Join(cols, ", ")+" WHERE request_id = '"+id+"'")
}

// CF14: interleaved spellings of one command give one argv; a flag left at its default (or a bool
// set false) never appears; --allow-failing once per value in the order given; bools bare.
func TestCanonicalFinishArgv(t *testing.T) {
	parse := func(args ...string) []string {
		fs := flag.NewFlagSet("finish", flag.ContinueOnError)
		fs.String("preset", "", "")
		fs.String("presets-dir", "/default", "")
		fs.String("title", "", "")
		fs.Bool("ready", false, "")
		fs.Bool("dry-run", false, "")
		fs.Bool("confirm-on-surface", false, "")
		var allowed multiFlag
		fs.Var(&allowed, "allow-failing", "")
		pos, err := verbs.ParseInterleaved(fs, args)
		if err != nil || len(pos) != 1 {
			t.Fatalf("%v: %v %v", args, pos, err)
		}
		return canonicalFinishArgv(fs, allowed, pos[0])
	}
	want := []string{"ailang-worldd", "episode", "finish", "ep1", "--allow-failing", "b", "--allow-failing", "a",
		"--confirm-on-surface", "--preset", "p", "--ready", "--title", "two words"}
	for _, args := range [][]string{
		{"ep1", "--preset", "p", "--allow-failing", "b", "--ready", "--title", "two words", "--allow-failing", "a", "--confirm-on-surface"},
		{"--confirm-on-surface", "--title=two words", "ep1", "--allow-failing=b", "--preset=p", "--allow-failing", "a", "--ready=true", "--dry-run=false"},
	} {
		if got := parse(args...); !reflect.DeepEqual(got, want) {
			t.Fatalf("%q\n got %q\nwant %q", args, got, want)
		}
	}
	if got := parse("ep1"); !reflect.DeepEqual(got, []string{"ailang-worldd", "episode", "finish", "ep1"}) {
		t.Fatalf("defaults appear: %q", got)
	}
}

// CF9: a request the law refuses (over 128 KiB), no randomness, failing randomness: each refuses
// before anything is stored or shown.
func TestFinishConfirmOnSurfaceRefusesBeforeTheStore(t *testing.T) {
	creds, in := surfaceUnit(t)
	big := in
	for i := 0; len(strings.Join(big.protected, "")) < approval.MaxVerbRequestBytes; i++ {
		big.protected = append(big.protected, fmt.Sprintf("p/%04d/%s", i, strings.Repeat("x", 200)))
	}
	missing := in
	missing.credentials = filepath.Join(t.TempDir(), "missing", "credentials.db")
	random := func(b []byte) error { return nil }
	for _, c := range []struct {
		name   string
		in     surfaceFinish
		random func([]byte) error
		want   string
	}{
		{"over 128 KiB", big, random, "the request is"},
		{"no randomness", in, nil, "internal error: no randomness source"},
		{"the id's draw fails", in, func(b []byte) error { return map[bool]error{true: errors.New("entropy gone")}[len(b) == 8] }, "randomness: entropy gone"},
		{"the challenge's draw fails", in, func(b []byte) error { return map[bool]error{true: errors.New("entropy gone")}[len(b) == 32] }, "randomness: entropy gone"},
		{"no credentials store", missing, random, "cannot read the credentials store " + missing.credentials},
		{"busy", in, random, "the credentials store is busy"},
	} {
		env, _ := awaitEnv(sfT0, new(bool), nil)
		env.random = c.random
		var tty, out bytes.Buffer
		if c.name == "busy" {
			st, err := store.OpenCredentials(creds)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = st.Close() }()
		}
		ref, refusal := finishConfirmOnSurface(boundedTestContext(t), nil, env, c.in, &tty, &out)
		if ref != "" || !strings.Contains(refusal, c.want) || tty.Len() != 0 || out.Len() != 0 {
			t.Fatalf("%s: %q %q tty %q out %q", c.name, ref, refusal, tty.String(), out.String())
		}
		if c.name == "over 128 KiB" && !strings.HasSuffix(refusal, "; confirm at the terminal") {
			t.Fatalf("%s: %q", c.name, refusal)
		}
		if n := pendingVerbRows(t, creds); n != 0 {
			t.Fatalf("%s: %d rows stored", c.name, n)
		}
	}
}

// The request's every field is finish's: the prompt is FinishPromptText of the same fields
// (committer and protected paths included), protected an array even when empty.
func TestFinishVerbRequestFields(t *testing.T) {
	_, in := surfaceUnit(t)
	for _, c := range []struct {
		committer string
		protected []string
	}{{"Mark <m@example.invalid>", []string{"go.mod"}}, {"", nil}} {
		in.committer, in.protected = c.committer, c.protected
		r, canon, err := finishVerbRequest(in)
		want := approval.VerbRequest{Schema: approval.VerbRequestSchema, Verb: approval.VerbFinish, Argv: in.argv, Episode: "ep1",
			Prompt: broker.FinishPromptText(in.scope, in.title, in.body, c.protected, c.committer), Scope: in.scopeText, Title: in.title, Body: in.body,
			Committer: c.committer, Protected: append([]string{}, c.protected...)}
		if back, derr := approval.DecodeVerb(canon); err != nil || derr != nil || !reflect.DeepEqual(r, want) || !reflect.DeepEqual(back, want) {
			t.Fatalf("%+v (%v %v)\nwant %+v", r, err, derr, want)
		}
	}
}

// writeAwaitLinks' address rule: one link per distinct passkey origin, then the daemon's loopback
// unless an origin already is it.
func TestSurfaceCardLinks(t *testing.T) {
	creds, _ := surfaceUnit(t)
	if got, err := surfaceCardLinks(boundedTestContext(t), creds, "/p", "127.0.0.1:7644"); err != nil || !reflect.DeepEqual(got, []string{"http://localhost:7644/p"}) {
		t.Fatalf("no passkey: %q %v", got, err)
	}
	seedPasskey(t, creds, "pk-a", approvaltest.NewPasskey(approval.AlgEdDSA, 31), sfT0)
	seedPasskey(t, creds, "pk-b", approvaltest.NewPasskey(approval.AlgEdDSA, 32), sfT0)
	for addr, want := range map[string][]string{"127.0.0.1:7644": {pkOrigin + "/p"}, "127.0.0.1:47330": {pkOrigin + "/p", "http://localhost:47330/p"}, "no-port": {pkOrigin + "/p"}} {
		if got, err := surfaceCardLinks(boundedTestContext(t), creds, "/p", addr); err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("%s: %q %v", addr, got, err)
		}
	}
}

// pendingVerbRows counts creds' verb_requests rows.
func pendingVerbRows(t *testing.T, creds string) int {
	t.Helper()
	db, err := sql.Open("sqlite", creds)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	var n int
	if err := db.QueryRow(`SELECT count(*) FROM verb_requests`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// CF8: the put's refusals — a busy writer, a store below schema 9 (naming the restart), 16 pending.
func TestPutFinishVerbRequestRefusals(t *testing.T) {
	creds, in := surfaceUnit(t)
	st, err := store.OpenCredentials(creds)
	if err != nil {
		t.Fatal(err)
	}
	_, canon, _ := finishVerbRequest(in)
	row := store.VerbRequestRow{RequestID: "vr-00000000000000f1", Request: canon, Digest: finishSHA(canon), Challenge: strings.Repeat("c", 64),
		EpisodeID: "ep1", CreatedAt: sfT0, ExpiresAt: sfT0 + 600}
	if got := putFinishVerbRequest(filepath.Join(t.TempDir(), "missing", "c.db"), row, sfT0); !strings.HasPrefix(got, "cannot open the credentials store ") {
		t.Fatalf("no store: %q", got)
	}
	if got := putFinishVerbRequest(creds, row, sfT0); got != "the credentials store is busy (another writer holds it); run finish again" {
		t.Fatalf("busy: %q", got)
	}
	_ = st.Close()
	for i := 0; i < store.MaxPendingVerbRequests; i++ {
		storeOwn(t, surfaceFinish{credentials: creds, addr: in.addr, ep: fmt.Sprintf("ep-%02d", i), argv: []string{"ailang-worldd", "episode", "finish", fmt.Sprintf("ep-%02d", i)}, scope: in.scope,
			scopeText: in.scopeText, title: in.title, body: in.body}, fmt.Sprintf("vr-00000000000001%02x", i), 1, sfT0)
	}
	if got := putFinishVerbRequest(creds, row, sfT0); got != "cannot store request vr-00000000000000f1: "+store.ErrTooManyPendingVerbRequests.Error() {
		t.Fatalf("cap: %q", got)
	}
	credSQL(t, creds, `DROP TABLE verb_requests`, `PRAGMA user_version = 8`)
	if got := putFinishVerbRequest(creds, row, sfT0); !strings.Contains(got, "verb requests need schema 9") || !strings.Contains(got, "restart `ailang-worldd serve`") {
		t.Fatalf("v8: %q", got)
	}
}
