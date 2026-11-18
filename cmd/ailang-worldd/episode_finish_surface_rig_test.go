package main

// Row 243 V3 sprint 5 M12 (design §2.7, §2.9 P1, §6 CF1–CF3, CF5, CF6, CF10, CF12, CF13; S4 §1.4 F4,
// CX21's CLI half): episode finish --confirm-on-surface end to end on the finish rig, with a real
// daemon (daemon.New over the rig's world and credentials stores, served through its Handler on a
// loopback request) answering POST /surface/confirm/{id} and GET /v1/objects/{ref}. The card is
// played by POSTing sigSigner's SSHSIG over approval.VerbMessage of the row the CLI stored; the
// terminal answers y and counts its reads (the surface mode never reads it); the clock is the
// await rig's fake one, started at the wall clock the daemon reads.

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
)

// countingTerminal answers y and counts every read; it records what is written to it.
type countingTerminal struct {
	bytes.Buffer
	in    io.Reader
	reads int
}

func (t *countingTerminal) Read(p []byte) (int, error) { t.reads++; return t.in.Read(p) }
func (t *countingTerminal) Close() error               { return nil }

type finishSurfaceRig struct {
	t        *testing.T
	f        *finishFixture
	d        *daemon.Daemon
	creds    string
	tty      *countingTerminal
	now      int64
	start    int64
	waits    int64
	draws    byte
	onWait   func(n int64)
	objReads int
}

func newSurfaceRig(t *testing.T) *finishSurfaceRig {
	t.Helper()
	f := newFinishFixture(t)
	f.allowRevoked = true // a re-run records after the first run's revocation
	r := &finishSurfaceRig{t: t, f: f, creds: f.health.CredentialsDB, now: time.Now().Unix()}
	seedApprover(t, r.creds, "mark-mac", sigSigner.PublicKey(), r.now-100)
	d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: f.db, CredentialsPath: r.creds, BindHost: daemon.DefaultBindHost, ErrorLog: io.Discard})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = d.Close() })
	r.d = d
	finishClient = func(string) *client {
		return memoryDaemon(func(w http.ResponseWriter, req *http.Request) {
			if strings.HasPrefix(req.URL.Path, "/v1/objects/") { // the read-back reaches the real daemon
				r.objReads++
				rec := r.handle(req.Method, req.URL.RequestURI(), nil)
				w.WriteHeader(rec.Code)
				_, _ = w.Write(rec.Body.Bytes())
				return
			}
			f.serve(w, req)
		})
	}
	return r
}

// handle serves one loopback request on the real daemon, as the card's page sends it.
func (r *finishSurfaceRig) handle(method, target string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), http.LocalAddrContextKey, net.Addr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 7644})))
	req.Host = "localhost:7644"
	if method == http.MethodPost {
		for k, v := range map[string]string{"Origin": "http://localhost:7644", "Sec-Fetch-Site": "same-origin", "Content-Type": "application/json"} {
			req.Header.Set(k, v)
		}
	}
	rec := httptest.NewRecorder()
	r.d.Handler().ServeHTTP(rec, req)
	return rec
}

// env is the finish's env: the counting terminal, the fake clock, and a counter as randomness (each
// draw fills its bytes with the next byte value, so ids and challenges are known to the test).
func (r *finishSurfaceRig) env(answer string) sessionEnv {
	r.tty = &countingTerminal{in: strings.NewReader(answer)}
	return sessionEnv{
		openTerminal: func() (io.ReadWriteCloser, error) { return r.tty, nil },
		now:          func() int64 { return r.now },
		random: func(b []byte) error {
			r.draws++
			for i := range b {
				b[i] = r.draws
			}
			return nil
		},
		after: func(d time.Duration) <-chan time.Time {
			r.waits++
			r.now += int64(d / time.Second)
			if r.onWait != nil {
				r.onWait(r.waits)
			}
			ch := make(chan time.Time, 1)
			ch <- time.Unix(r.now, 0)
			return ch
		},
	}
}

// run is one finish with --confirm-on-surface; its fake clock starts at the wall clock the daemon
// reads (a request created ahead of it would be "clock behind" at the POST).
func (r *finishSurfaceRig) run(extra ...string) (int, string, string) {
	r.waits, r.now = 0, time.Now().Unix()
	r.start = r.now
	return r.f.run(r.env("y\n"), append([]string{"--include-unattributed", "--confirm-on-surface"}, extra...)...)
}

// row reads one verb request from the credentials store.
func (r *finishSurfaceRig) row(id string) store.VerbRequestRow {
	r.t.Helper()
	ro, err := store.OpenCredentialsReadOnly(r.creds)
	if err != nil {
		r.t.Fatal(err)
	}
	defer func() { _ = ro.Close() }()
	row, ok, err := ro.VerbRequest(boundedTestContext(r.t), id)
	if err != nil || !ok {
		r.t.Fatalf("row %s: %v %v", id, ok, err)
	}
	return row
}

// card plays the card for id: s's SSHSIG over the stored row's approval.VerbMessage, POSTed to S4's
// POST /surface/confirm/{id}; it returns the status.
func (r *finishSurfaceRig) card(s *approvaltest.Signer, id string) int {
	r.t.Helper()
	row := r.row(id)
	dg, _ := hex.DecodeString(row.Digest)
	ch, _ := hex.DecodeString(row.Challenge)
	body, _ := json.Marshal(map[string]any{"v": 1, "digest": row.Digest, "sshsig": string(verbSig(r.t, s, approval.VerbNamespace, [32]byte(dg), [32]byte(ch), row.ExpiresAt))})
	return r.handle(http.MethodPost, "/surface/confirm/"+id, body).Code
}

// nextWallSecond waits (≤ 1 s, bounded at 3 s) for the wall clock the daemon reads to leave t0: a
// finish approval of the same scope in the same second is the same request object (design W59).
func nextWallSecond(t0 int64) {
	for limit := time.Now().Add(3 * time.Second); time.Now().Unix() == t0 && time.Now().Before(limit); {
		time.Sleep(20 * time.Millisecond)
	}
}

// object is the real daemon's GET /v1/objects/{ref} payload, decoded into out.
func (r *finishSurfaceRig) object(ref string, out any) {
	r.t.Helper()
	var o wireObject
	rec := r.handle(http.MethodGet, "/v1/objects/"+ref+"?payload=true", nil)
	if rec.Code != 200 || json.Unmarshal(rec.Body.Bytes(), &o) != nil || json.Unmarshal(o.Payload, out) != nil {
		r.t.Fatalf("object %s: %d %s", ref, rec.Code, rec.Body.String())
	}
}

// id is the request id the n-th run's first draw makes ("vr-" and 8 bytes of value n).
func drawnID(n byte) string { return "vr-" + strings.Repeat(fmt.Sprintf("%02x", n), 8) }

// kept fails unless a refusal kept the commit and dispatched nothing (design §2.7: every refusal).
func (r *finishSurfaceRig) kept(rc int, errs, want string) string {
	r.t.Helper()
	sha := strings.TrimSpace(r.f.gitOutput("-C", r.f.wt, "rev-parse", "HEAD"))
	if rc != exitUsage || !strings.Contains(errs, want) || !strings.Contains(errs, "commit kept at "+sha) || !strings.Contains(errs, "next: ailang-worldd episode finish ep1 --preset fixture") ||
		sha == r.f.before || len(r.f.seam.calls) != 0 || len(r.f.finishRecords) != 0 || len(r.f.scopes) != 0 {
		r.t.Fatalf("rc %d calls %v records %d\nstderr %s\nwant %q", rc, r.f.seam.calls, len(r.f.finishRecords), errs, want)
	}
	return sha
}

// CF3 (+ CF13, CF14's rig half, judge-403 note 1's CLI half): the card confirms; the CLI re-checks,
// reads the decision back from the real daemon and opens the pull request with that decision. The
// terminal is never read; it shows exactly the stored prompt, which is FinishPromptText of the
// stored fields; the stored row is the CLI's own (id, challenge, times from the env).
func TestFinishConfirmOnSurfaceOpensThePR(t *testing.T) {
	r := newSurfaceRig(t)
	id := drawnID(1)
	r.onWait = func(n int64) {
		if n == 1 {
			if code := r.card(sigSigner, id); code != 200 {
				t.Errorf("card POST: %d", code)
			}
		}
	}
	r.f.rep.Start.Checks = []store.StartCheck{{Name: "pass", Argv: []string{"go", "test"}, Cwd: ".", TimeoutS: 1, Exec: "pinned", Args: []string{"test"}}}
	rc, out, errs := r.run("--allow-failing", "pass")
	if rc != exitOK {
		t.Fatalf("rc %d\nstdout %s\nstderr %s", rc, out, errs)
	}
	start, row := r.start, r.row(id)
	req, err := approval.DecodeVerb(row.Request)
	if err != nil {
		t.Fatal(err)
	}
	sha := strings.TrimSpace(r.f.gitOutput("-C", r.f.wt, "rev-parse", "HEAD"))
	if row.Challenge != strings.Repeat("02", 32) || row.EpisodeID != "ep1" || row.CreatedAt != start || row.ExpiresAt != start+600 || row.Digest != finishSHA(row.Request) ||
		row.Outcome != store.VerbOutcomeApproved || row.ApprovedBy != "ssh:mark-mac" {
		t.Fatalf("row %+v", row)
	}
	dg, _ := hex.DecodeString(row.Digest)
	h, _ := approval.VerbMessageHash([32]byte(dg), [32]byte(bytes.Repeat([]byte{2}, 32)), start+600)
	pr := r.f.finishRecords[0].PR.URL // the fake gh's own answer
	want := fmt.Sprintf("committed at %s\nbody: %s\nconfirm request %s on the surface (expires %s, message hash %s):\n  open http://localhost:47330/surface/confirm/%s\n"+
		"the typed y is not read in this mode; Ctrl-C stops waiting and keeps the commit at %s\nconfirmed by ssh:mark-mac (decision %s)\n", sha, filepath.Join(r.f.sessions, "ep1.pr-body.md"),
		id, approverTime(start+600), hex.EncodeToString(h[:8]), id, sha, row.DecisionRef)
	if !strings.HasPrefix(out, want) || !strings.HasSuffix(out, "\n"+pr+"\n") || pr == "" {
		t.Fatalf("stdout\n%s\nwant prefix\n%s\nand the PR %q last", out, want, pr)
	}
	if rec := r.f.finishRecords; len(rec) != 1 || rec[0].Approval != row.DecisionRef || rec[0].Outcome != "done" || len(r.f.scopes) != 0 || r.f.outcomes != 1 {
		t.Fatalf("records %+v, typed approve phases %d", rec, len(r.f.scopes))
	}
	if r.tty.reads != 0 || r.tty.String() != req.Prompt || r.objReads != 2 {
		t.Fatalf("terminal reads %d, object reads %d, shows\n%q\nwant the stored prompt\n%q", r.tty.reads, r.objReads, r.tty.String(), req.Prompt)
	}
	body, _ := os.ReadFile(filepath.Join(r.f.sessions, "ep1.pr-body.md"))
	scope, err := broker.ParseFinishApprovalScope(req.Scope)
	if err != nil || req.Prompt != broker.FinishPromptText(scope, req.Title, req.Body, req.Protected, req.Committer) || req.Body != string(body) ||
		req.Scope != r.f.finishRecords[0].Scope || scope.SHA != sha || !scope.Draft || req.Protected == nil {
		t.Fatalf("request %+v (%v)", req, err)
	}
	if wantArgv := []string{"ailang-worldd", "episode", "finish", "ep1", "--allow-failing", "pass", "--confirm-on-surface", "--include-unattributed", "--preset", "fixture", "--presets-dir", r.f.dir}; !reflect.DeepEqual(req.Argv, wantArgv) {
		t.Fatalf("argv %q, want %q", req.Argv, wantArgv)
	}
	var decision struct {
		Decision, DecidedBy string
		Now                 int64
	}
	r.object(row.DecisionRef, &decision)
	if decision.Decision != "approve" || decision.DecidedBy != "ssh:mark-mac" || decision.Now != row.ConsumedAt {
		t.Fatalf("decision %+v, row consumed at %d", decision, row.ConsumedAt)
	}
}

// CF1, CF2: --dry-run refuses after the terminal and before any read; an unsandboxed pinned check
// refuses before the check loop — no check ran, no outcome recorded, no commit, the terminal unread.
func TestFinishConfirmOnSurfaceRefusesBeforeAnyCheck(t *testing.T) {
	r := newSurfaceRig(t)
	rc, _, errs := r.run("--dry-run")
	if rc != exitUsage || !strings.Contains(errs, "--confirm-on-surface cannot be combined with --dry-run\n  why: a dry run opens no pull request\n  next: drop one of them") || len(r.f.requests) != 0 || r.tty.reads != 0 {
		t.Fatalf("CF1: rc %d requests %v reads %d\n%s", rc, r.f.requests, r.tty.reads, errs)
	}
	r.f.rep.Start.Checks = []store.StartCheck{{Name: "host", Argv: []string{"/bin/sh", "-c", "echo ran > marker"}, Cwd: ".", TimeoutS: 1, Unsandboxed: true}}
	rc, _, errs = r.run()
	_, statErr := os.Stat(filepath.Join(r.f.wt, "marker"))
	if rc != exitUsage || !strings.Contains(errs, "check host is unsandboxed: it needs its own typed y\n  why: --confirm-on-surface replaces only the pull request's y\n  next: run finish without --confirm-on-surface") ||
		r.f.outcomes != 0 || !os.IsNotExist(statErr) || r.tty.reads != 0 || r.tty.Len() != 0 || strings.TrimSpace(r.f.gitOutput("-C", r.f.wt, "rev-parse", "HEAD")) != r.f.before || pendingVerbRows(t, r.creds) != 0 {
		t.Fatalf("CF2: rc %d outcomes %d marker %v reads %d\n%s", rc, r.f.outcomes, statErr, r.tty.reads, errs)
	}
}

// protectedProfile pins an exec profile protecting *_test.go and a committer for f's episode (as the
// typed path's TestEpisodeFinishProtected does) and changes a_test.go; it returns the profiles dir.
func protectedProfile(t *testing.T, f *finishFixture) string {
	t.Helper()
	profiles := filepath.Join(f.root, "profiles")
	raw, _ := json.Marshal(map[string]any{"profile": broker.ExecProfileVersion, "project": "fixture", "env": map[string]string{}, "path": []string{"/usr/bin", "/bin"},
		"timeout_ms": 9000, "protected": []string{"*_test.go"}, "commands": map[string]any{"check": map[string]any{"argv": []string{"true"}, "max_args": 1}}})
	profile, err := broker.ParseExecProfile(raw)
	if err != nil || os.MkdirAll(profiles, 0o700) != nil || os.WriteFile(filepath.Join(profiles, "fixture.json"), raw, 0o600) != nil {
		t.Fatal(err)
	}
	creds, err := store.OpenCredentials(f.health.CredentialsDB)
	if err != nil || authority.Revoke(boundedTestContext(t), creds, f.credentialID) != nil {
		t.Fatal(err)
	}
	_, f.credentialID, _, err = authority.MintBound(boundedTestContext(t), creds, "ep1", authority.ProfileBinding{Project: "fixture", Digest: profile.Digest},
		[]broker.Capability{{Effect: "Workspace.Write", Scope: "world", Budget: 10, ExpiresAt: 4000000000}}, 3600, time.Now().Unix(), nil)
	if creds.Close(); err != nil {
		t.Fatal(err)
	}
	f.rep.Start.Preset.Bytes = `{"profile":"fixture","committer":"Test Operator <operator@example.invalid>"}`
	f.write("a_test.go", "protected edit\n")
	return profiles
}

// CF4's rig half, CF13, CF10: a request nobody confirms ends at expiry + grace with the commit kept;
// the typed path's prompt over the same commit is the stored prompt byte for byte (the typed run
// declines); a re-run stores a new request with the same prompt bytes, supersedes the old one at
// its own clock, and opens the pull request when that one is confirmed.
func TestFinishConfirmOnSurfaceRerunsAndTypedPrompt(t *testing.T) {
	r := newSurfaceRig(t)
	profiles := protectedProfile(t, r.f) // a protected path and a pinned committer in the prompt
	first := drawnID(1)
	r.onWait = func(n int64) {
		if n == 1 {
			r.now = r.row(first).ExpiresAt + 10 // the next poll is the last
		}
	}
	rc, _, errs := r.run("--exec-profiles-dir", profiles)
	sha := r.kept(rc, errs, "no confirmation of request "+first+" by ")
	old := r.row(first)
	stored, _ := approval.DecodeVerb(old.Request)
	if old.Outcome != "" || r.waits != 1 || !strings.HasPrefix(stored.Prompt, "committer: Test Operator <operator@example.invalid>\n") || !strings.Contains(stored.Prompt, "\nprotected paths changed: a_test.go") {
		t.Fatalf("first: %q %q after %d waits", old.Outcome, stored.Prompt, r.waits)
	}
	typed := &captureFinishTerminal{Reader: strings.NewReader("n\n")}
	env := testMintEnv("", 1)
	env.openTerminal = func() (io.ReadWriteCloser, error) { return typed, nil }
	if rc, _, errs := r.f.run(env, "--include-unattributed"); rc != exitUsage || !strings.Contains(errs, "pull request declined") {
		t.Fatalf("typed: %d %s", rc, errs)
	}
	if typed.String() != stored.Prompt+"Open this pull request? [y/N] " {
		t.Fatalf("CF13: the typed prompt\n%q\nis not the stored prompt\n%q", typed.String(), stored.Prompt)
	}
	second := drawnID(3)
	r.onWait = func(n int64) {
		if n == 1 && r.card(sigSigner, second) != 200 {
			t.Error("card POST")
		}
	}
	if rc, out, errs := r.run(); rc != exitOK || strings.TrimSpace(r.f.gitOutput("-C", r.f.wt, "rev-parse", "HEAD")) != sha {
		t.Fatalf("re-run: %d\n%s\n%s", rc, out, errs)
	}
	newer, old := r.row(second), r.row(first)
	again, _ := approval.DecodeVerb(newer.Request)
	if again.Prompt != stored.Prompt || newer.Outcome != store.VerbOutcomeApproved || old.Outcome != store.VerbOutcomeSuperseded || old.ConsumedAt != r.start {
		t.Fatalf("CF10: newer %+v old %+v (re-run at %d)", newer, old, r.start)
	}
}

// CF5 and CF6's rig rows (S1 Appendix C note 5): a newer finish supersedes the request; an approver
// revoked after the confirmation was recorded makes the CLI refuse it (fail-safe) — each keeps the
// commit and dispatches nothing; a re-run asks afresh (CF10).
func TestFinishConfirmOnSurfaceKeepsTheCommit(t *testing.T) {
	r := newSurfaceRig(t)
	in := surfaceFinish{credentials: r.creds, ep: "ep1", argv: []string{"ailang-worldd", "episode", "finish", "ep1"}, scopeText: "x", title: "t", body: "b"}
	r.onWait = func(n int64) {
		if n == 1 {
			storeOwn(t, in, "vr-00000000000000e1", 0xe1, r.now) // a newer finish of ep1
		}
	}
	rc, _, errs := r.run()
	r.kept(rc, errs, "request "+drawnID(1)+" was superseded by a newer finish of ep1")
	revoked := drawnID(3)
	r.onWait = func(n int64) {
		if n == 1 {
			if code := r.card(sigSigner, revoked); code != 200 {
				t.Errorf("card POST: %d", code)
			}
			revokeSeeded(t, r.creds, "mark-mac", r.now)
		}
	}
	rc, _, errs = r.run()
	r.kept(rc, errs, "the confirmation of request "+revoked+" does not verify: ")
	if !strings.Contains(errs, "revoked") || r.row(revoked).Outcome != store.VerbOutcomeApproved || r.objReads != 0 {
		t.Fatalf("revoked: %s (object reads %d)", errs, r.objReads)
	}
}

// CX21's CLI half (S4 §1.4 F4): C15 fails after C14 (a verb_requests UPDATE trigger) → 500, one
// decision no row names; the identical body a second later → 200, a second decision the row names;
// the waiting finish dispatches with the row's decision, never the orphan.
func TestFinishConfirmOnSurfaceDispatchesTheNamedDecision(t *testing.T) {
	r := newSurfaceRig(t)
	id := drawnID(1)
	r.onWait = func(n int64) {
		if n != 1 {
			return
		}
		credSQL(t, r.creds, `CREATE TRIGGER fail_consume BEFORE UPDATE ON verb_requests WHEN NEW.outcome = 'approved' BEGIN SELECT RAISE(ABORT, 'injected consume failure'); END`)
		t0 := time.Now().Unix()
		if code := r.card(sigSigner, id); code != 500 || r.row(id).Outcome != "" {
			t.Errorf("first POST: %d", code)
		}
		credSQL(t, r.creds, `DROP TRIGGER fail_consume`)
		nextWallSecond(t0) // a new decision, not the same one again (W59)
		if code := r.card(sigSigner, id); code != 200 {
			t.Errorf("retry POST: %d", code)
		}
	}
	if rc, out, errs := r.run(); rc != exitOK {
		t.Fatalf("rc %d\n%s\n%s", rc, out, errs)
	}
	named := r.row(id).DecisionRef
	db, err := sql.Open("sqlite", r.f.db)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	var all, isNamed int
	if err := db.QueryRow(`SELECT count(*), count(*) FILTER (WHERE hash_ref = ?) FROM objects WHERE semantic_id = ?`, named, broker.ApprovalDecisionV1).Scan(&all, &isNamed); err != nil {
		t.Fatal(err)
	}
	if all != 2 || isNamed != 1 || r.f.finishRecords[0].Approval != named {
		t.Fatalf("%d decisions (%d named %s); the record names %s", all, isNamed, named, r.f.finishRecords[0].Approval)
	}
}

// CF12 (S1 Appendix C note 4): a surface-confirmed finish whose push failed is resumed by a re-run
// under its recorded approval — mode resume-dispatch reuses it (the unchanged branch, design §2.7
// step 3): no new request, no wait, no re-check or read-back (CF6/CF7 are skipped), F6's re-check of
// HEAD, remote, title and body still runs; the pull request opens under the first decision.
func TestFinishConfirmOnSurfaceResumeReusesTheApproval(t *testing.T) {
	r := newSurfaceRig(t)
	id := drawnID(1)
	r.onWait = func(n int64) {
		if n == 1 && r.card(sigSigner, id) != 200 {
			t.Error("card POST")
		}
	}
	r.f.seam.fail = "push"
	if rc, _, errs := r.run(); rc != exitFatal || len(r.f.finishRecords) != 1 || r.f.finishRecords[0].Outcome != "push-failed" {
		t.Fatalf("first: %d %s", rc, errs)
	}
	decision := r.row(id).DecisionRef
	r.f.seam.fail, r.objReads = "", 0
	r.onWait = func(int64) { t.Error("the resume waited") }
	rc, out, errs := r.run()
	if rc != exitOK || pendingVerbRows(t, r.creds) != 1 || r.objReads != 0 || r.tty.reads != 0 || !strings.HasPrefix(r.tty.String(), "resuming under approval "+decision+", expires in ") {
		t.Fatalf("resume: %d rows %d object reads %d tty %q\n%s\n%s", rc, pendingVerbRows(t, r.creds), r.objReads, r.tty.String(), out, errs)
	}
	if rec := r.f.finishRecords; len(rec) != 2 || rec[1].Approval != decision || rec[1].Outcome != "done" {
		t.Fatalf("records %+v", rec)
	}
}
