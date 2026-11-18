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
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
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

// CF4 (+ CF5's and the exact-boundary rows, every way the wait ends): under the fake clock and a
// wall-clock test deadline, the wait ends at expires_at + verbLandGrace; a confirmation landing at
// exactly that instant is seen and one a second later is not; a supersession, a vanished row.
func TestAwaitVerbConfirmationEnds(t *testing.T) {
	type end struct {
		row     store.VerbRequestRow
		refusal string
		waits   int64
	}
	run := func(t *testing.T, creds string, own verbOwn, start int64, onWait func(n, now int64)) end {
		t.Helper()
		env, clock := awaitEnv(start, new(bool), onWait)
		done := make(chan end, 1)
		go func() {
			row, refusal := awaitVerbConfirmation(env, creds, own)
			done <- end{row, refusal, clock.waits}
		}()
		select {
		case e := <-done:
			return e
		case <-time.After(20 * time.Second):
			t.Fatal("the wait did not end: unbounded")
		}
		return end{}
	}
	t.Run("expiry+grace, bounded", func(t *testing.T) {
		creds, in := surfaceUnit(t)
		own := storeOwn(t, in, "vr-00000000000000a1", 0xa1, sfT0)
		deadline := own.expires + 10
		e := run(t, creds, own, own.expires-2, nil)
		if e.refusal != "no confirmation of request vr-00000000000000a1 by "+approverTime(deadline) || e.waits != 12 {
			t.Fatalf("%q after %d waits", e.refusal, e.waits)
		}
	})
	for _, c := range []struct {
		name string
		at   int64 // seconds after expires_at that the confirmation lands
		seen bool
	}{{"lands at expires_at+grace", 10, true}, {"lands a second later", 11, false}} {
		t.Run(c.name, func(t *testing.T) {
			creds, in := surfaceUnit(t)
			own := storeOwn(t, in, "vr-00000000000000a2", 0xa2, sfT0)
			art := verbSig(t, sigSigner, approval.VerbNamespace, own.digest, own.challenge, own.expires)
			e := run(t, creds, own, own.expires-2, func(_, now int64) {
				if now == own.expires+c.at {
					approveByHand(t, creds, own.id, art, "ssh:mark-mac", now)
				}
			})
			if seen := e.refusal == "" && e.row.Outcome == store.VerbOutcomeApproved && e.row.ApprovedBy == "ssh:mark-mac"; seen != c.seen {
				t.Fatalf("seen %v, want %v: %q", seen, c.seen, e.refusal)
			}
		})
	}
	t.Run("superseded", func(t *testing.T) {
		creds, in := surfaceUnit(t)
		own := storeOwn(t, in, "vr-00000000000000a3", 0xa3, sfT0)
		e := run(t, creds, own, sfT0, func(n, now int64) {
			if n == 1 {
				storeOwn(t, in, "vr-00000000000000a4", 0xa4, now)
			}
		})
		if e.refusal != "request vr-00000000000000a3 was superseded by a newer finish of ep1" || e.waits != 1 {
			t.Fatalf("%q after %d waits", e.refusal, e.waits)
		}
	})
	t.Run("a poll read fails", func(t *testing.T) {
		creds, in := surfaceUnit(t)
		own := storeOwn(t, in, "vr-00000000000000a5", 0xa5, sfT0)
		e := run(t, creds, own, sfT0, func(n, _ int64) {
			if n == 2 {
				credSQL(t, creds, `DELETE FROM verb_requests`)
			}
		})
		if e.refusal != "cannot read request vr-00000000000000a5: request vr-00000000000000a5 is no longer in the store" || e.waits != 2 {
			t.Fatalf("%q after %d waits", e.refusal, e.waits)
		}
	})
}

// CF6 (design §2.7 step 6(a)): an approved row is re-checked against this process's own digest,
// challenge and expiry and the enrolled key; each row marked approved by hand (raw SQL) that a key
// did not give over these bytes refuses "does not verify" — including an approver revoked after the
// confirmation (S1 Appendix C note 5: fail-safe, the commit is kept and a re-run asks afresh).
func TestPollVerbRequestRechecks(t *testing.T) {
	other := [32]byte(bytes.Repeat([]byte{0xee}, 32))
	for _, c := range []struct {
		name string
		art  func(own verbOwn) []byte
		by   string
		set  []string
		pre  func(t *testing.T, creds string)
		want string
	}{
		{"control: the key over these bytes", nil, "ssh:mark-mac", nil, nil, ""},
		{"the row's digest is not ours", nil, "ssh:mark-mac", []string{"digest = '" + hex.EncodeToString(other[:]) + "'"}, nil, "the stored request is not the one this finish made"},
		{"the row's challenge is not ours", nil, "ssh:mark-mac", []string{"challenge = '" + hex.EncodeToString(other[:]) + "'"}, nil, "the stored request is not the one this finish made"},
		{"the row's expiry is not ours", nil, "ssh:mark-mac", []string{"expires_at = expires_at + 1"}, nil, "the stored request is not the one this finish made"},
		{"signed over another digest", func(o verbOwn) []byte {
			return verbSig(t, sigSigner, approval.VerbNamespace, other, o.challenge, o.expires)
		}, "ssh:mark-mac", nil, nil, "does not verify: "},
		{"signed over another challenge", func(o verbOwn) []byte {
			return verbSig(t, sigSigner, approval.VerbNamespace, o.digest, other, o.expires)
		}, "ssh:mark-mac", nil, nil, "does not verify: "},
		{"signed over another expiry", func(o verbOwn) []byte {
			return verbSig(t, sigSigner, approval.VerbNamespace, o.digest, o.challenge, o.expires+1)
		}, "ssh:mark-mac", nil, nil, "does not verify: "},
		{"the mint namespace", func(o verbOwn) []byte {
			return verbSig(t, sigSigner, approval.Namespace, o.digest, o.challenge, o.expires)
		}, "ssh:mark-mac", nil, nil, "does not verify: "},
		{"a key never enrolled", func(o verbOwn) []byte {
			return verbSig(t, otherSigner, approval.VerbNamespace, o.digest, o.challenge, o.expires)
		}, "ssh:mark-mac", nil, nil, "does not verify: "},
		{"not an artifact", func(verbOwn) []byte { return []byte("y") }, "ssh:mark-mac", nil, nil, "does not verify: "},
		{"the approver revoked since", nil, "ssh:mark-mac", nil, func(t *testing.T, creds string) { revokeSeeded(t, creds, "mark-mac", sfT0+5) }, "revoked"},
		{"approved_by is not the key's", nil, "ssh:mark-other", nil, nil, "the key is ssh:mark-mac's, the row says ssh:mark-other"},
	} {
		t.Run(c.name, func(t *testing.T) {
			creds, in := surfaceUnit(t)
			own := storeOwn(t, in, "vr-00000000000000b1", 0xb1, sfT0)
			art := verbSig(t, sigSigner, approval.VerbNamespace, own.digest, own.challenge, own.expires)
			if c.art != nil {
				art = c.art(own)
			}
			approveByHand(t, creds, own.id, art, c.by, sfT0+1, c.set...)
			if c.pre != nil {
				c.pre(t, creds)
			}
			row, refusal, err := pollVerbRequest(creds, own, sfT0+10)
			if err != nil || (c.want == "") != (refusal == "") || !strings.Contains(refusal, c.want) {
				t.Fatalf("%q %v, want %q", refusal, err, c.want)
			}
			if c.want != "" && !strings.HasPrefix(refusal, "the confirmation of request vr-00000000000000b1 does not verify: ") {
				t.Fatalf("%q", refusal)
			}
			if c.want == "" && (row.Outcome != store.VerbOutcomeApproved || !bytes.Equal(row.Artifact, art)) {
				t.Fatalf("control row %+v", row)
			}
		})
	}
}

// CF7 (design §2.7 step 6(b)): the decision the row names is read back through worldReader.object
// over GET /v1/objects; the decision and its request are broker's own writes (RecordOperatorApproval
// on a test world store), each refused row the same object with one field changed.
func TestReadBackFinishDecision(t *testing.T) {
	_, in := surfaceUnit(t)
	st, err := store.Open(filepath.Join(t.TempDir(), "world.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = st.Close() }()
	ctx := boundedTestContext(t)
	qref, dref, err := broker.RecordOperatorApproval(ctx, st, in.scopeText, "operator:episode-finish", "ssh:mark-mac", sfT0)
	if err != nil {
		t.Fatal(err)
	}
	real := map[string]wireObject{}
	for _, ref := range []string{qref.String(), dref.String()} {
		obj, ok, err := st.GetObject(ctx, hashref.MustParse(ref))
		if err != nil || !ok {
			t.Fatal(ref, ok, err)
		}
		real[ref] = wireObject{Hash: ref, SemanticID: obj.SemanticID, Payload: obj.Payload}
	}
	// edit is the real object at ref with one payload field (or the semantic id) changed.
	edit := func(ref, field string, value any) map[string]wireObject {
		objs := map[string]wireObject{}
		for k, v := range real {
			objs[k] = v
		}
		o := objs[ref]
		if field == "semanticId" {
			o.SemanticID = value.(string)
		} else {
			var m map[string]any
			if err := json.Unmarshal(o.Payload, &m); err != nil {
				t.Fatal(err)
			}
			m[field] = value
			o.Payload, _ = json.Marshal(m)
		}
		objs[ref] = o
		return objs
	}
	q, d := qref.String(), dref.String()
	for _, c := range []struct {
		name string
		objs map[string]wireObject
		want string
	}{
		{"control", real, ""},
		{"a deny", edit(d, "decision", "deny"), `decision "deny"`},
		{"decided by another", edit(d, "decidedBy", "ssh:mark-other"), `by "ssh:mark-other"`},
		{"not a decision", edit(d, "semanticId", broker.ApprovalRequestV1), "not a finish approval decision"},
		{"its request is missing", edit(d, "requestRef", "sha256:"+strings.Repeat("0", 64)), "is not readable"},
		{"its request is not a request", edit(q, "semanticId", broker.ApprovalDecisionV1), "is not readable"},
		{"another effect", edit(q, "effect", "Workspace.Write"), "over another scope"},
		{"another requester", edit(q, "requester", "operator:terminal"), "over another scope"},
		{"another scope", edit(q, "scope", strings.Replace(in.scopeText, "draft=true", "draft=false", 1)), "over another scope"},
	} {
		c := c
		r := worldReader{ctx: ctx, c: memoryDaemon(func(w http.ResponseWriter, req *http.Request) {
			o, ok := c.objs[strings.TrimPrefix(req.URL.Path, "/v1/objects/")]
			if !ok || req.URL.Query().Get("payload") != "true" {
				http.NotFound(w, req)
				return
			}
			_ = json.NewEncoder(w).Encode(o)
		})}
		got := readBackFinishDecision(r, d, "ssh:mark-mac", in.scopeText)
		if (c.want == "") != (got == "") || !strings.Contains(got, c.want) || (got != "" && !strings.HasPrefix(got, "the world decision "+d+" does not approve this finish (F6): ")) {
			t.Fatalf("%s: %q, want %q", c.name, got, c.want)
		}
	}
}

// §2.9 P1 "decision does not approve this scope", at the wiring: a confirmation that re-checks but
// whose decision the daemon does not hold is refused (F6), and no decision is returned.
func TestFinishConfirmOnSurfaceRefusesAnUnrecordedDecision(t *testing.T) {
	creds, in := surfaceUnit(t)
	env, _ := awaitEnv(sfT0, new(bool), func(n, now int64) {
		if n == 1 {
			ro, _ := store.OpenCredentialsReadOnly(creds)
			rows, err := ro.PendingVerbRequests(boundedTestContext(t), now, 16)
			_ = ro.Close()
			if err != nil || len(rows) != 1 {
				t.Fatalf("pending %v %v", rows, err)
			}
			d, _ := hex.DecodeString(rows[0].Digest)
			c, _ := hex.DecodeString(rows[0].Challenge)
			approveByHand(t, creds, rows[0].RequestID, verbSig(t, sigSigner, approval.VerbNamespace, [32]byte(d), [32]byte(c), rows[0].ExpiresAt), "ssh:mark-mac", now)
		}
	})
	env.random = func(b []byte) error { return nil }
	var tty, out bytes.Buffer
	ref, refusal := finishConfirmOnSurface(boundedTestContext(t), memoryDaemon(http.NotFound), env, in, &tty, &out)
	if ref != "" || !strings.HasPrefix(refusal, "the world decision sha256:"+strings.Repeat("d", 64)+" does not approve this finish (F6): ") || strings.Contains(out.String(), "confirmed by") {
		t.Fatalf("%q %q\n%s", ref, refusal, out.String())
	}
}
