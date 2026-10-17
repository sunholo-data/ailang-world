package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/canon"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/evidence"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/pins"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

type demoTestRig struct {
	db, session, pin, addr, manifest string
	mode                             string
	mu                               sync.Mutex
	commits                          []cliCommit
	headers                          []string
}

func newDemoTestRig(t *testing.T, mode string) *demoTestRig {
	t.Helper()
	pin := os.Getenv("AILANG_BIN")
	if pin == "" {
		t.Fatal("required AILANG_BIN absent")
	}
	dir := t.TempDir()
	r := &demoTestRig{mode: mode, db: filepath.Join(dir, "world.db"), session: filepath.Join(dir, "session"), pin: pin, manifest: filepath.Join(dir, "manifest.json")}
	st, err := store.Open(r.db)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().Unix()
	expiry := now + 3600
	grants := []broker.Capability{{Effect: "world.apply", Scope: "world", Budget: 10, ExpiresAt: expiry}, {Effect: "Workspace.Read", Scope: "worktree", Budget: 1, ExpiresAt: expiry}, {Effect: "Ailang.CLI", Scope: "worktree", Budget: 2, ExpiresAt: expiry}, {Effect: "Human.Approve", Scope: "worktree", Budget: 1, ExpiresAt: expiry}}
	if mode == "missing-grant" {
		grants = grants[1:]
	}
	if mode == "expired-grant" {
		grants[0].ExpiresAt = now - 10
	}
	if mode == "insufficient-budget" {
		grants[0].Budget = 9
	}
	mintTime := now
	if mode == "expired-session" {
		mintTime = now - 7200
	}
	token, _, _, err := authority.Mint(boundedTestContext(t), st, func() string {
		if mode == "wrong-episode" {
			return "other"
		}
		return "demo"
	}(), grants, 3600, mintTime, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if mode == "world-only" {
		h := func(s string) hashref.HashRef { return hashref.SumSHA256([]byte(s)) }
		if e := st.PutWorld(boundedTestContext(t), store.World{Ref: h("unselected-world"), Revision: 0, StateRoot: h("unselected-state"), LogHead: h("unselected-log")}); e != nil {
			t.Fatal(e)
		}
	}
	if mode == "spent-budget" {
		_, _, e := st.AppendNextEffectIntent(boundedTestContext(t), "demo", store.EffectIntent{EpisodeID: "demo", Effect: "Ailang.CLI", Scope: "worktree", Cost: 1, RequestRef: hashref.SumSHA256([]byte("prior-request")), LogicalTime: now})
		if e != nil {
			t.Fatal(e)
		}
	}
	if err = st.Close(); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(r.session, []byte(token), 0600); err != nil {
		t.Fatal(err)
	}
	d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: r.db, BindHost: daemon.DefaultBindHost, BindPort: 0, AilangBin: pin})
	if err != nil {
		t.Fatal(err)
	}
	if strings.HasPrefix(mode, "archive-") {
		pinBytes, e := os.ReadFile(pin)
		if e != nil {
			t.Fatal(e)
		}
		sidecar := filepath.Join(r.db+".artifacts", "interpreters", "sha256", hashref.SumSHA256(pinBytes).Digest(), "manifest.json")
		if mode == "archive-absent" {
			if e := os.Remove(sidecar); e != nil {
				t.Fatal(e)
			}
		} else {
			raw, e := os.ReadFile(sidecar)
			if e != nil {
				t.Fatal(e)
			}
			var manifest map[string]any
			if e = json.Unmarshal(raw, &manifest); e != nil {
				t.Fatal(e)
			}
			switch mode {
			case "archive-hash":
				manifest["hash"] = hashref.SumSHA256([]byte("other")).String()
			case "archive-size":
				manifest["size"] = 1
			case "archive-version", "archive-release":
				manifest["version"] = "AILANG v0.99.0\n"
			}
			raw, e = json.Marshal(manifest)
			if e != nil {
				t.Fatal(e)
			}
			if e = os.WriteFile(sidecar, raw, 0644); e != nil {
				t.Fatal(e)
			}
		}
	}
	h := d.Handler()
	headReads := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Synthetic HTTP arms exercise refusals only; acceptance always uses
		// the unmodified real daemon and its sole writer.
		if (mode == "epoch-decode" || mode == "epoch-candidate") && (req.URL.Path == "/v1/registry/"+store.EpochRegistryV1 || strings.HasPrefix(req.URL.Path, "/v1/objects/")) {
			payload := []byte("bad JSON")
			if mode == "epoch-candidate" {
				payload = []byte(`{"semantic_id":"world/epoch-registry/v1","epochs":[{"epoch":1,"candidates":["AILANG v0.99.0"]}]}`)
			}
			ref := hashref.SumSHA256(payload).String()
			if strings.HasPrefix(req.URL.Path, "/v1/registry/") {
				_ = json.NewEncoder(w).Encode(map[string]string{"head": ref})
			} else {
				_ = json.NewEncoder(w).Encode(wireObject{Hash: ref, SemanticID: store.EpochRegistryV1, Payload: payload})
			}
			return
		}
		if req.URL.Path == "/v1/health" && mode == "health-absent" {
			http.NotFound(w, req)
			return
		}
		if req.URL.Path == "/v1/health" && mode == "health-malformed" {
			_, _ = io.WriteString(w, "not json")
			return
		}
		if req.URL.Path == "/v1/health" && mode == "health-down" {
			http.Error(w, "unavailable", 503)
			return
		}
		if req.URL.Path == "/v1/registry/"+store.EpochRegistryV1 && mode == "missing-epoch" {
			http.NotFound(w, req)
			return
		}
		if strings.HasPrefix(req.URL.Path, "/v1/objects/") && mode == "epoch-hash" {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			var o wireObject
			_ = json.Unmarshal(rec.Body.Bytes(), &o)
			o.Payload = []byte("tampered epoch")
			_ = json.NewEncoder(w).Encode(o)
			return
		}
		if req.URL.Path == "/v1/health" && (mode == "readonly-db" || mode == "health-interpreter" || mode == "health-version" || mode == "archive-release") {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			var health daemon.HealthResponse
			if e := json.Unmarshal(rec.Body.Bytes(), &health); e != nil {
				t.Error(e)
			}
			if mode == "readonly-db" {
				health.DBPath = filepath.Join(filepath.Dir(r.db), "nonexistent.db")
			} else if mode == "health-interpreter" {
				health.InterpreterRef = hashref.SumSHA256([]byte("other")).String()
			} else {
				health.InterpreterVersion = "AILANG v0.99.0\n"
			}
			_ = json.NewEncoder(w).Encode(health)
			return
		}
		if req.URL.Path == "/v1/head" {
			headReads++
			if mode == "cas-change" && headReads == 2 {
				c := makeCLICommit(cliWorld{}, 0, "interloper")
				raw, _ := json.Marshal(c)
				inject := httptest.NewRequest("POST", "/v1/commit", bytes.NewReader(raw))
				inject.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()
				h.ServeHTTP(rec, inject)
				if rec.Code != 200 {
					t.Errorf("interloper commit failed: %d %s", rec.Code, rec.Body.String())
				}
			}
		}
		if req.Method == http.MethodPost && req.URL.Path == "/v1/commit" {
			raw, e := io.ReadAll(req.Body)
			if e != nil {
				t.Error(e)
			}
			req.Body = io.NopCloser(bytes.NewReader(raw))
			var c cliCommit
			if e = json.Unmarshal(raw, &c); e != nil {
				t.Error(e)
			}
			r.mu.Lock()
			r.commits = append(r.commits, c)
			r.headers = append(r.headers, req.Header.Get("Authorization"))
			r.mu.Unlock()
		}
		if req.Method == "POST" && req.URL.Path == "/v1/commit" && strings.HasPrefix(mode, "commit-") {
			if mode == "commit-denied" {
				w.WriteHeader(403)
				_ = json.NewEncoder(w).Encode(daemon.APIError{Error: daemon.APIErrorDetail{Class: "Denied", Message: "refusal control"}})
				return
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != 200 {
				t.Errorf("commit control real daemon: %d %s", rec.Code, rec.Body.String())
			}
			if mode == "commit-uncertain" {
				hijack, ok := w.(http.Hijacker)
				if !ok {
					t.Error("hijack unavailable")
					return
				}
				conn, _, e := hijack.Hijack()
				if e != nil {
					t.Error(e)
					return
				}
				_ = conn.Close()
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"selectedHead": "wrong"})
			return
		}
		h.ServeHTTP(w, req)
	}))
	r.addr = server.URL
	t.Cleanup(func() {
		server.Close()
		if e := d.Close(); e != nil {
			t.Error(e)
		}
	})
	return r
}
func (r *demoTestRig) args() []string {
	return []string{"demo", "seed", "--db", r.db, "--session", r.session, "--ailang-bin", r.pin, "--proof-key-file", filepath.Join(filepath.Dir(r.db), "private.key"), "--manifest", r.manifest}
}
func demoRun(t *testing.T, r *demoTestRig, args []string) (int, string, string) {
	t.Helper()
	cwd, e := os.Getwd()
	if e != nil {
		t.Fatal(e)
	}
	target := filepath.Join("..", "..")
	if r.mode == "source-absent" {
		target = t.TempDir()
	}
	if e = os.Chdir(target); e != nil {
		t.Fatal(e)
	}
	defer func() {
		if e := os.Chdir(cwd); e != nil {
			t.Fatal(e)
		}
	}()
	var out, err bytes.Buffer
	code := run(append([]string{"--addr", r.addr}, args...), &out, &err)
	r.mu.Lock()
	r.mu.Unlock() // synchronize completed HTTP recording before assertions
	return code, out.String(), err.String()
}
func TestDemoSeedAuthenticatedCommits(t *testing.T) {
	r := newDemoTestRig(t, "")
	code, out, errOut := demoRun(t, r, r.args())
	if code != exitUsage || !strings.Contains(errOut, "S2-S6 incomplete") {
		t.Fatalf("successful seed commit stage must reach explicit incomplete boundary: code=%d out=%s err=%s", code, out, errOut)
	}
	if out != "" {
		t.Fatalf("incomplete CLI printed success output: %q", out)
	}
	if len(r.commits) != 5 {
		t.Fatalf("domain commit count=%d want5", len(r.commits))
	}
	token, _ := os.ReadFile(r.session)
	reader := worldReader{ctx: boundedTestContext(t), c: newClient(r.addr)}
	source, err := os.ReadFile(filepath.Join("..", "..", "world", "logepoch.ail"))
	if err != nil {
		t.Fatal(err)
	}
	wantIDs := []string{"world/demo/goal/v1", "world/demo/proposal/v1", "world/demo/proposal/v1", "world/source/ailang/v1", "world/demo/authored-claim/v1"}
	pinBytes, e := os.ReadFile(r.pin)
	if e != nil {
		t.Fatal(e)
	}
	prev := ""
	prevEntry := ""
	for i, c := range r.commits {
		if r.headers[i] != "Bearer "+string(token) {
			t.Fatalf("commit%d missing session bearer", i)
		}
		if c.ObservedHead != prev || c.NextWorld.Revision != int64(i) {
			t.Fatalf("commit%d CAS/index mismatch", i)
		}
		if c.Entry.Header.Interpreter != hashref.SumSHA256(pinBytes).String() || c.Entry.Header.SemanticsEpoch != 1 || c.Entry.Header.WrittenBy != "demo" || c.Entry.Header.EntryIndex != int64(i) {
			t.Fatalf("commit%d frozen header identity mismatch", i)
		}
		if i > 0 && c.Entry.Header.PrevEntryHash != prevEntry {
			t.Fatalf("commit%d previous entry mismatch", i)
		}
		world, ok, e := reader.world(c.NextWorld.Ref)
		if e != nil || !ok || world.LogHead != c.Entry.EntryHash || world.Revision != int64(i) || world.StateRoot != c.NextWorld.StateRoot {
			t.Fatalf("commit%d real world mismatch: %+v %v", i, world, e)
		}
		if len(c.Objects) != 1 {
			t.Fatalf("commit%d objects=%d", i, len(c.Objects))
		}
		o := c.Objects[0]
		if c.Entry.Header.TransitionFn != o.Hash {
			t.Fatalf("commit%d transition fn not payload", i)
		}
		if o.InterfaceHash != hashref.SumSHA256([]byte(wantIDs[i])).String() || o.Provenance != "demo-authored" {
			t.Fatalf("commit%d interface/provenance mismatch", i)
		}
		if o.SemanticID != wantIDs[i] || hashref.SumSHA256(o.Payload).String() != o.Hash {
			t.Fatalf("commit%d kind/hash mismatch", i)
		}
		stored, ok, e := reader.object(o.Hash)
		if e != nil || !ok || !bytes.Equal(stored.Payload, o.Payload) {
			t.Fatalf("commit%d payload not actually stored: %v", i, e)
		}
		entry, ok, e := reader.entry(int64(i))
		if e != nil || !ok || entry.TransitionRef != o.Hash || entry.EntryHash != c.Entry.EntryHash {
			t.Fatalf("commit%d no real carrying entry: %+v %v", i, entry, e)
		}
		if i == 3 && !bytes.Equal(source, o.Payload) {
			t.Fatal("source differs from unchanged production logepoch")
		}
		if i == 4 && (!bytes.Contains(o.Payload, []byte("demo-authored")) || bytes.Contains(o.Payload, []byte("AiReview"))) {
			t.Fatal("commentary must be explicitly demo-authored")
		}
		prev = c.NextWorld.Ref
		prevEntry = c.Entry.EntryHash
	}
	if r.commits[1].Objects[0].Hash == r.commits[2].Objects[0].Hash {
		t.Fatal("two candidates must differ")
	}
	if _, e := os.Stat(r.manifest); !os.IsNotExist(e) {
		t.Fatal("incomplete stage wrote manifest")
	}
}
func TestDemoSeedPreflightRefusals(t *testing.T) {
	for _, mode := range []string{"manifest", "nonempty", "absent-session", "unknown-session", "expired-session", "missing-grant", "expired-grant", "wrong-db", "pin-mismatch", "health-interpreter", "health-version", "insufficient-budget", "spent-budget", "wrong-episode", "cas-change", "positional", "missing-db", "missing-pin", "missing-key", "missing-manifest", "health-absent", "health-malformed", "health-down", "missing-epoch", "epoch-hash", "readonly-db", "pin-absent", "world-only", "epoch-decode", "epoch-candidate", "source-absent", "unknown-verb", "archive-absent", "archive-hash", "archive-size", "archive-version", "archive-release"} {
		t.Run(mode, func(t *testing.T) {
			r := newDemoTestRig(t, mode)
			args := r.args()
			if mode == "unknown-verb" {
				args[1] = "unknown"
			}
			if mode == "manifest" {
				if e := os.WriteFile(r.manifest, []byte("existing"), 0600); e != nil {
					t.Fatal(e)
				}
			}
			if mode == "absent-session" {
				t.Setenv("WORLD_SESSION", "")
				args = append(args[:4], args[6:]...)
			}
			if mode == "unknown-session" {
				if e := os.WriteFile(r.session, []byte(strings.Repeat("a", 64)), 0600); e != nil {
					t.Fatal(e)
				}
			}
			if mode == "pin-absent" {
				args[7] = filepath.Join(t.TempDir(), "missing-pin")
			}
			if mode == "readonly-db" {
				args[3] = filepath.Join(filepath.Dir(r.db), "nonexistent.db")
			}
			if mode == "wrong-db" {
				args[3] = filepath.Join(t.TempDir(), "other.db")
			}
			if mode == "pin-mismatch" {
				// Hash identity must refuse before any version/manifest lookup.
				// A local regular byte file needs no separate workspace tool.
				args[7] = filepath.Join(t.TempDir(), "mismatched-interpreter")
				if e := os.WriteFile(args[7], []byte("mismatched pin fixture: "+r.db), 0600); e != nil {
					t.Fatal(e)
				}
			}
			if mode == "positional" {
				args = append(args, "unexpected")
			}
			for name, index := range map[string]int{"missing-db": 3, "missing-pin": 7, "missing-key": 9, "missing-manifest": 11} {
				if mode == name {
					args[index] = ""
				}
			}
			before := 0
			if mode == "nonempty" {
				c := makeCLICommit(cliWorld{}, 0, "preexisting")
				data, _ := json.Marshal(c)
				token, _ := os.ReadFile(r.session)
				status, body, e := newClient(r.addr).doAuth(boundedTestContext(t), "POST", "/v1/commit", bytes.NewReader(data), "Bearer "+string(token))
				if e != nil || status != 200 {
					t.Fatalf("setup commit %d %s %v", status, body, e)
				}
				before = len(r.commits)
			}
			reader := worldReader{ctx: boundedTestContext(t), c: newClient(r.addr)}
			beforeEntries, e := reader.page(0, 10)
			if e != nil {
				t.Fatal(e)
			}
			code, out, errOut := demoRun(t, r, args)
			afterEntries, e := reader.page(0, 10)
			if e != nil {
				t.Fatal(e)
			}
			wantEntries := len(beforeEntries)
			if mode == "cas-change" {
				wantEntries++
			}
			if len(afterEntries) != wantEntries {
				t.Fatalf("%s: refusal changed public inventory: actual log entries%d want%d", mode, len(afterEntries), wantEntries)
			}
			if len(r.commits) != before {
				t.Fatalf("%s: refusal changed public inventory: commits%d before%d", mode, len(r.commits), before)
			}
			if code != exitUsage || errOut == "" || strings.Contains(errOut, "S1 domain commits recorded; S2-S6 incomplete") {
				t.Fatalf("preflight must refuse before stage: %d %s %s", code, out, errOut)
			}
			expected := map[string]string{"manifest": "manifest must", "nonempty": "nonempty public", "absent-session": "no session credential", "unknown-session": "SessionUnknown", "expired-session": "SessionExpired", "missing-grant": "world.apply@world grant", "expired-grant": "world.apply@world grant", "wrong-db": "db/health identity", "pin-mismatch": "interpreter pin mismatch", "health-interpreter": "interpreter pin mismatch", "health-version": "interpreter version mismatch", "insufficient-budget": "world.apply@world grant", "spent-budget": "Ailang.CLI@worktree grant", "wrong-episode": "demo episode session required", "cas-change": "unexpected CAS head", "positional": "unexpected positional", "missing-db": "are required", "missing-pin": "are required", "missing-key": "are required", "missing-manifest": "are required", "health-absent": "required resource absent", "health-malformed": "malformed body", "health-down": "HTTP 503", "missing-epoch": "required resource absent", "epoch-hash": "epoch registry object mismatch", "readonly-db": "cannot open read-only authority", "pin-absent": "cannot read interpreter pin", "world-only": "nonempty public", "epoch-decode": "registry: decode", "epoch-candidate": "no nominated epoch", "source-absent": "prepare worktree containing", "unknown-verb": "usage:", "archive-absent": "interpreter archive manifest unavailable", "archive-hash": "interpreter archive identity mismatch", "archive-size": "interpreter archive identity mismatch", "archive-version": "interpreter version mismatch", "archive-release": "interpreter version mismatch"}[mode]
			if expected == "" || !strings.Contains(errOut, expected) {
				t.Fatalf("wrong refusal: want %q got %q", expected, errOut)
			}
			if len(r.commits) != before {
				t.Fatalf("%s: refusal changed public inventory: commits%d before%d", mode, len(r.commits), before)
			}
			if mode != "manifest" {
				if _, e := os.Stat(r.manifest); !os.IsNotExist(e) {
					t.Fatal("refusal wrote manifest")
				}
			}
			token, _ := os.ReadFile(r.session)
			if strings.Contains(out+errOut, string(token)) {
				t.Fatal("credential leaked")
			}
		})
	}
}
func TestDemoSeedHelpPreparation(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"demo", "seed", "--help"}, &out, &errOut); code != exitOK {
		t.Fatalf("demo help=%d %s", code, errOut.String())
	}
	// The proof pin and the tool release are the manifest's (pins/ailang-releases.json).
	for _, s := range []string{"--db", "--session", "--ailang-bin", "--proof-key-file", "--addr", "--manifest", pins.ReleaseFor(pins.RoleProof).Release, toolRelease(), "S2-S6"} {
		if !strings.Contains(out.String(), s) {
			t.Errorf("help missing %s", s)
		}
	}
	for _, args := range [][]string{{"demo", "--help"}, {"help", "demo"}, {"help", "demo", "seed"}} {
		var o, e bytes.Buffer
		if code := run(args, &o, &e); code != exitOK || !strings.Contains(o.String(), "S2-S6") {
			t.Fatalf("help form %v failed: %d %s", args, code, e.String())
		}
	}
	doc, e := os.ReadFile(filepath.Join("..", "..", "docs", "QUICKSTART.md"))
	if e != nil {
		t.Fatal(e)
	}
	for _, s := range []string{"demo seed", "world/logepoch.ail", "packages/world-demo/transitions.json", "typed", "git worktree add", "Human.Approve=worktree:1", "S2-S6"} {
		if !bytes.Contains(doc, []byte(s)) {
			t.Errorf("preparation docs missing %s", s)
		}
	}
	_, preparation, ok := strings.Cut(string(doc), "## 12. Prepare the release demo seed")
	if !ok {
		t.Fatal("demo preparation section absent")
	}
	for _, command := range []string{`export DEMO_BIN=$DEMO_DIR/ailang-worldd`, `go build -o "$DEMO_BIN" ./cmd/ailang-worldd`, `"$DEMO_BIN" session mint`, `"$DEMO_BIN" serve`, `"$DEMO_BIN" demo seed`} {
		if !strings.Contains(preparation, command) {
			t.Errorf("demo artifact command wiring missing %s", command)
		}
	}
}

func TestDemoSeedCommitFailures(t *testing.T) {
	for _, mode := range []string{"commit-denied", "commit-response", "commit-uncertain"} {
		t.Run(mode, func(t *testing.T) {
			r := newDemoTestRig(t, mode)
			code, out, errOut := demoRun(t, r, r.args())
			want := map[string]string{"commit-denied": "HTTP 403 Denied", "commit-response": "commit response mismatch", "commit-uncertain": "commit outcome uncertain"}[mode]
			if code != exitUsage || out != "" || !strings.Contains(errOut, want) {
				t.Fatalf("commit refusal boundary want %q: code%d out%q err%q", want, code, out, errOut)
			}
			if len(r.commits) != 1 {
				t.Fatalf("no retries: attempts%d", len(r.commits))
			}
			reader := worldReader{ctx: boundedTestContext(t), c: newClient(r.addr)}
			entries, e := reader.page(0, 10)
			wantEntries := 1
			if mode == "commit-denied" {
				wantEntries = 0
			}
			if e != nil || len(entries) != wantEntries {
				t.Fatalf("real commit control entries%d want%d: %v", len(entries), wantEntries, e)
			}
			if _, e := os.Stat(r.manifest); !os.IsNotExist(e) {
				t.Fatal("refused commit wrote manifest")
			}
		})
	}
}

// M2b (row 150): S2 read and S3 test/check through real /a2a/ on a daemon
// serving the published se-tools package, episode "demo".

// requireBinaryRelease fails (never skips) when a configured binary reports
// the wrong release.
func requireBinaryRelease(t *testing.T, bin, want string) {
	t.Helper()
	out, err := exec.CommandContext(boundedTestContext(t), bin, "--version").Output()
	if err != nil || !strings.Contains(string(out), want) {
		t.Fatalf("%s --version = %q (%v), want release %s", bin, out, err, want)
	}
}

type demoSeRig struct {
	url, db, token, ws string
	log                *bytes.Buffer
}

func newDemoSeRig(t *testing.T, interp, tool string) demoSeRig {
	return newDemoSeRigForSeed(t, interp, tool, false)
}

func newDemoSeRigForSeed(t *testing.T, interp, tool string, empty bool) demoSeRig {
	return newDemoApprovalRig(t, interp, tool, empty, false, 20)
}

func newDemoApprovalRig(t *testing.T, interp, tool string, empty, extra bool, budget int64) demoSeRig {
	t.Helper()
	// Production capsule/source-check scratch roots follow os.TempDir().
	// Keep those writes inside this test as well as the explicit fixture files.
	t.Setenv("TMPDIR", t.TempDir())
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	db, ws, corpus := filepath.Join(base, "state", "world.db"), filepath.Join(base, "ws"), filepath.Join(base, "outside", "corpus")
	for _, d := range []string{filepath.Dir(db), filepath.Join(ws, "demo", "world"), filepath.Join(corpus, "runnable")} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	source, err := os.ReadFile(filepath.Join("..", "..", demoLogepochPath))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ws, "demo", demoLogepochPath), source, 0o644); err != nil {
		t.Fatal(err)
	}
	writeFileT(t, filepath.Join(corpus, "manifest.json"), `{"examples":[]}`+"\n")
	log := new(bytes.Buffer)
	cfg := daemon.Config{DBPath: db, BindHost: daemon.DefaultBindHost, AilangBin: interp, WorkspaceRoot: ws,
		ToolAilangBin: tool, ExamplesDir: corpus, ErrorLog: log}
	d, err := daemon.New(boundedTestContext(t), cfg)
	if err != nil {
		t.Fatalf("daemon.New: %v", err)
	}
	srv := httptest.NewServer(d.Handler())
	_, body, err := newClient(srv.URL).get("/v1/health")
	srv.Close()
	if cerr := d.Close(); err != nil || cerr != nil {
		t.Fatalf("health %v close %v", err, cerr)
	}
	var h healthInfo
	if json.Unmarshal([]byte(body), &h) != nil || h.InterpreterRef == "" {
		t.Fatalf("health = %s", body)
	}
	st, err := store.Open(db)
	if err != nil {
		t.Fatal(err)
	}
	interpRef := hashref.MustParse(h.InterpreterRef)
	// A selected world is required before any tool call; a genesis commit
	// stands in for the demo's S1 domain commits (as in cliwalk_e2e_test).
	obj := store.Object{Hash: hashref.SumSHA256([]byte("genesis-state")), InterfaceHash: hashref.SumSHA256([]byte("test/genesis")),
		SemanticID: "test/genesis", Provenance: "demo_test", Payload: []byte("genesis-state")}
	entryHash := hashref.SumSHA256([]byte("genesis-entry"))
	w := store.World{Ref: hashref.SumSHA256([]byte("genesis-world")), Revision: 0, StateRoot: obj.Hash, LogHead: entryHash}
	if !empty {
		if err := st.Commit(boundedTestContext(t), store.Commit{Objects: []store.Object{obj}, NextWorld: w, Entry: store.LogEntry{Header: store.LogHeader{
			EntryIndex: 0, SemanticsEpoch: 1, TransitionFn: hashref.SumSHA256([]byte("genesis-fn")), Interpreter: interpRef,
			PrevEntryHash: hashref.SumSHA256([]byte("genesis-prev")), WrittenBy: "demo_test"}, EntryHash: entryHash, TransitionRef: obj.Hash}}); err != nil {
			t.Fatalf("genesis: %v", err)
		}
	}
	publishSeToolsCLI(t, st, db, interpRef)
	publishDemoApproval(t, st, db, interpRef, extra)
	now := time.Now().Unix()
	grants := []broker.Capability{
		{Effect: "world.apply", Scope: "world", ExpiresAt: now + 7200, Budget: 100},
		{Effect: "Human.Approve", Scope: "worktree", ExpiresAt: now + 7200, Budget: budget},
	}
	for _, e := range []string{broker.EffectWorkspaceRead, broker.EffectAilangCLI} {
		grants = append(grants, broker.Capability{Effect: e, Scope: broker.WorkspaceScope, ExpiresAt: now + 7200, Budget: 20})
	}
	tok, _, _, err := authority.Mint(boundedTestContext(t), st, "demo", grants, 3600, now, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	d, err = daemon.New(boundedTestContext(t), cfg)
	if err != nil {
		t.Fatalf("daemon.New (serve): %v", err)
	}
	srv = httptest.NewServer(d.Handler())
	t.Cleanup(func() {
		srv.Close()
		_ = d.Close()
	})
	return demoSeRig{url: srv.URL, db: db, token: tok, ws: ws, log: log}
}

// demoSpend reads the episode's per-effect spend through the read-only view.
func demoSpend(t *testing.T, db string) map[store.EffectKey]int64 {
	t.Helper()
	ro, err := store.OpenReadOnly(db)
	if err != nil {
		t.Fatal(err)
	}
	defer ro.Close()
	spend, err := ro.EffectSpend(boundedTestContext(t), "demo")
	if err != nil {
		t.Fatal(err)
	}
	out := map[store.EffectKey]int64{}
	for k, v := range spend {
		out[k] = int64(v)
	}
	return out
}

func TestDemoSeedReadTestCheck(t *testing.T) {
	interp, tool := cliSeBins(t)
	requireBinaryRelease(t, interp, "v0.41.0")
	requireBinaryRelease(t, tool, daemon.ToolBinaryRelease)
	r := newDemoSeRig(t, interp, tool)
	source, err := os.ReadFile(filepath.Join("..", "..", demoLogepochPath))
	if err != nil {
		t.Fatal(err)
	}
	before := demoSpend(t, r.db)
	stages, err := demoRunToolStages(boundedTestContext(t), newClient(r.url), r.token, source)
	after := demoSpend(t, r.db)
	if stages.Read.Content != string(source) {
		t.Fatalf("S2 read bytes differ from the committed %s: got %d bytes %.160q (err %v)", demoLogepochPath, len(stages.Read.Content), stages.Read.Content, err)
	}
	if stages.Test.Op != "test" {
		t.Fatalf("retained TESTED invocation argv[op] is test: got %q", stages.Test.Op)
	}
	if err != nil {
		t.Fatalf("S3 stages: %v (test %+v)", err, stages.Test)
	}
	if !strings.Contains(strings.Join(stages.Test.Argv, " "), "test") || stages.Test.ExitCode != 0 {
		t.Fatalf("test argv/exit = %v %d", stages.Test.Argv, stages.Test.ExitCode)
	}
	rep, err := demoValidateTestReport(stages.Test.Stdout)
	if err != nil || rep.Passed == 0 {
		t.Fatalf("test report: %+v %v", rep, err)
	}
	if stages.Check.Op != "ai_check" || stages.Check.ExitCode != 0 {
		t.Fatalf("ai_check result = %+v", stages.Check)
	}
	if err := demoValidateAICheck(stages.Check.Stdout); err != nil {
		t.Fatalf("ai_check actual success: %v", err)
	}
	for _, run := range []demoToolRun{stages.Read, stages.Test, stages.Check} {
		if run.InvocationID == "" || run.EffectRecord == "" {
			t.Fatalf("%s lacks retained invocation/effect refs: %+v", run.Task, run)
		}
	}
	// Retained pin: the daemon archived the tool binary and refused anything but
	// the manifest's tool release at startup; the test records that archived
	// release for this run.
	toolBytes, err := os.ReadFile(tool)
	if err != nil {
		t.Fatal(err)
	}
	toolRef := hashref.SumSHA256(toolBytes)
	if m, err := archive.New(r.db).ReadManifest(toolRef); err != nil || !strings.Contains(m.Version, toolRelease()) {
		t.Fatalf("retained tool pin %s is not %s: %v %q", toolRef, toolRelease(), err, m.Version)
	}
	if d := after[store.EffectKey{Effect: broker.EffectWorkspaceRead, Scope: broker.WorkspaceScope}] - before[store.EffectKey{Effect: broker.EffectWorkspaceRead, Scope: broker.WorkspaceScope}]; d != 1 {
		t.Fatalf("Workspace.Read budget delta = %d, want 1", d)
	}
	if d := after[store.EffectKey{Effect: broker.EffectAilangCLI, Scope: broker.WorkspaceScope}] - before[store.EffectKey{Effect: broker.EffectAilangCLI, Scope: broker.WorkspaceScope}]; d != 2 {
		t.Fatalf("Ailang.CLI budget delta = %d, want 2 (test + ai_check)", d)
	}
}

// TestDemoSeedRejectsToolReportRefusals: the null-report control runs with no
// pin (unit control); the real arms need both pins and must refuse.
func TestDemoSeedRejectsToolReportRefusals(t *testing.T) {
	// Unit control: runs with no pin, so ordinary CI exercises it (PASS).
	t.Run("null-report-unit-control", func(t *testing.T) {
		_, err := demoValidateTestReport("null")
		if err == nil || !strings.Contains(err.Error(), "null report is refused") {
			t.Fatalf("null report is refused: got %v", err)
		}
	})
	// Real arms need both pins; ordinary CI skips this subtest only.
	t.Run("real-refusal-arms", testDemoRealRefusalArms)
}

func testDemoRealRefusalArms(t *testing.T) {
	interp, tool := cliSeBins(t)
	requireBinaryRelease(t, interp, "v0.41.0")
	requireBinaryRelease(t, tool, daemon.ToolBinaryRelease)
	r := newDemoSeRig(t, interp, tool)
	bad := filepath.Join(r.ws, "demo", "bad")
	if err := os.MkdirAll(bad, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFileT(t, filepath.Join(bad, "malformed.ail"), "module broken {\n")
	contracts, err := os.ReadFile(filepath.Join("..", "..", "world", "contracts.ail"))
	if err != nil {
		t.Fatal(err)
	}
	writeFileT(t, filepath.Join(bad, "nogen.ail"), string(contracts))
	ctx := boundedTestContext(t)
	c := newClient(r.url)
	if _, err := demoTestStage(ctx, c, r.token, demoTestTask, "ai_check", "bad/malformed.ail"); err == nil {
		t.Fatal("non-test op must be refused")
	}
	for _, p := range []string{"bad/malformed.ail", "bad/nogen.ail"} {
		// The validator decides on the real stdout even when the CLI exits
		// non-zero; a refusal must be a validator refusal, not an accepted run.
		run, _ := demoTestStage(ctx, c, r.token, demoTestTask, "test", p)
		_, verr := demoValidateTestReport(run.Stdout)
		t.Logf("%s: exit=%d validator=%v", p, run.ExitCode, verr)
		if verr == nil {
			t.Fatalf("%s: real refusal arm was accepted: %+v", p, run)
		}
	}
	// S2 mismatch arm: fresh task IDs; the error must name the S2 read-bytes
	// mismatch itself, not a later S3 refusal (the B1 control).
	stages, err := demoRunStages(ctx, c, r.token, []byte("not the committed source"), "-mismatch")
	if err == nil || !strings.Contains(err.Error(), "S2 read bytes differ") {
		t.Fatalf("S2 read-bytes mismatch must be named: err=%v", err)
	}
	if stages.Read.Content == "" || stages.Test.Task != "" {
		t.Fatalf("mismatch arm must stop after S2: read=%q test=%+v", stages.Read.Content, stages.Test)
	}
}

// TestDemoSeedToolGuardControls: unit controls for the report, ai_check,
// op-allowlist and A2A guards. They need no pins, so ordinary CI runs them.
func TestDemoSeedToolGuardControls(t *testing.T) {
	cases := []struct{ name, report, want string }{
		{"success-false", `{"success":false,"failed_tests":0,"passed_tests":3,"skipped_tests":0,"vacuous_skips":0}`, "success=false is refused"},
		{"failed-tests", `{"success":true,"failed_tests":1,"passed_tests":3,"skipped_tests":0,"vacuous_skips":0}`, "failed_tests=1 is refused"},
		{"skipped-tests", `{"success":true,"failed_tests":0,"passed_tests":3,"skipped_tests":2,"vacuous_skips":0}`, "skipped_tests=2 is refused"},
		{"vacuous-skips", `{"success":true,"failed_tests":0,"passed_tests":3,"skipped_tests":0,"vacuous_skips":4}`, "vacuous_skips=4 is refused"},
		{"zero-passes", `{"success":true,"failed_tests":0,"passed_tests":0,"skipped_tests":0,"vacuous_skips":0}`, "zero named passes"},
		{"null", `null`, "null report is refused"},
	}
	for _, c := range cases {
		t.Run("report-"+c.name, func(t *testing.T) {
			_, err := demoValidateTestReport(c.report)
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("report %s: err=%v, want %q", c.name, err, c.want)
			}
		})
	}
	t.Run("report-positive", func(t *testing.T) {
		if _, err := demoValidateTestReport(`{"success":true,"failed_tests":0,"passed_tests":10,"skipped_tests":0,"vacuous_skips":0}`); err != nil {
			t.Fatalf("clean report refused: %v", err)
		}
	})
	aiCases := []struct{ name, report, want string }{
		{"null", `null`, "null ai_check report is refused"},
		{"verified-zero", `{"check":{"passed":true,"error_count":0},"verify":{"available":true,"verified":0,"counterexample":0,"errors":0}}`, "verified=0 is refused"},
		{"check-failed", `{"check":{"passed":false,"error_count":0},"verify":{"available":true,"verified":2,"counterexample":0,"errors":0}}`, "check.passed=false is refused"},
		{"error-count", `{"check":{"passed":true,"error_count":3},"verify":{"available":true,"verified":2,"counterexample":0,"errors":0}}`, "error_count=3 is refused"},
	}
	for _, c := range aiCases {
		t.Run("ai-"+c.name, func(t *testing.T) {
			err := demoValidateAICheck(c.report)
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("ai_check %s: err=%v, want %q", c.name, err, c.want)
			}
		})
	}
	t.Run("ai-positive", func(t *testing.T) {
		if err := demoValidateAICheck(`{"check":{"passed":true,"error_count":0},"verify":{"available":true,"verified":2,"counterexample":0,"errors":0}}`); err != nil {
			t.Fatalf("clean ai_check refused: %v", err)
		}
	})

	// A fake /a2a/ that counts requests and answers with the given task.

	goodData := map[string]any{"content": "x", "world": map[string]any{"effects": []any{map[string]any{"record": "sha256:" + strings.Repeat("a", 64)}}}}
	goodMeta := map[string]any{"invocation_id": "a2a:demo:t", "entry_index": 0}
	ctx := boundedTestContext(t)
	t.Run("a2a-positive", func(t *testing.T) {
		c, _ := demoFakeA2A(t, "completed", goodMeta, goodData)
		if _, _, err := demoA2ASend(ctx, c, "tok", "t", "ailang-read", map[string]any{}); err != nil {
			t.Fatalf("complete task refused: %v", err)
		}
	})
	t.Run("a2a-non-completed", func(t *testing.T) {
		c, _ := demoFakeA2A(t, "working", goodMeta, goodData)
		_, _, err := demoA2ASend(ctx, c, "tok", "t", "ailang-read", map[string]any{})
		if err == nil || !strings.Contains(err.Error(), "tasks/send refused") {
			t.Fatalf("non-completed state accepted: %v", err)
		}
	})
	t.Run("a2a-missing-invocation", func(t *testing.T) {
		c, _ := demoFakeA2A(t, "completed", map[string]any{"entry_index": 0}, goodData)
		_, _, err := demoA2ASend(ctx, c, "tok", "t", "ailang-read", map[string]any{})
		if err == nil || !strings.Contains(err.Error(), "invocation or effect record reference absent") {
			t.Fatalf("missing invocation id accepted: %v", err)
		}
	})
	t.Run("a2a-missing-effect-record", func(t *testing.T) {
		c, _ := demoFakeA2A(t, "completed", goodMeta, map[string]any{"content": "x"})
		_, _, err := demoA2ASend(ctx, c, "tok", "t", "ailang-read", map[string]any{})
		if err == nil || !strings.Contains(err.Error(), "invocation or effect record reference absent") {
			t.Fatalf("missing effect record accepted: %v", err)
		}
	})
	t.Run("op-allowlist-before-dispatch", func(t *testing.T) {
		c, calls := demoFakeA2A(t, "completed", goodMeta, goodData)
		_, err := demoTestStage(ctx, c, "tok", demoTestTask, "shell", demoLogepochPath)
		if err == nil || !strings.Contains(err.Error(), "refuses unlisted op") || *calls != 0 {
			t.Fatalf("unlisted op: err=%v dispatches=%d", err, *calls)
		}
	})
}

// TestDemoSeedProofStage executes the command from an empty public world, with
// real S1 commits, published se-tools S2/S3, and the production S4 adapters.
func TestDemoSeedProofStage(t *testing.T) {
	interp := requireProofPin(t)
	tool := os.Getenv(toolBinEnvVar)
	if tool == "" {
		// Same environment gate as every published-se-tools test (cliSeBins):
		// CI does not provision the tool pin, so this skips there and
		// the pin-free S4 controls (TestDemoProof*, TestDemoSeedRequiredIdentities)
		// carry CI coverage. A configured wrong version still fails below.
		t.Skipf("S4 seed integration composes published se-tools S2/S3 and needs %s (%s); unset, skipping", toolBinEnvVar, daemon.ToolBinaryRelease)
	}
	requireBinaryRelease(t, tool, daemon.ToolBinaryRelease)
	for _, arm := range []string{"success", "missing", "directory", "short", "long", "safe-refusal", "lost-response-after-commit"} {
		t.Run(arm, func(t *testing.T) {
			r := newDemoSeRigForSeed(t, interp, tool, true)
			ctx := boundedTestContext(t)
			c := demoProofClient(r.url)
			key := freshProofKey(t)
			private := t.TempDir()
			keyPath := filepath.Join(private, "proof.key")
			raw := append([]byte(nil), key[:]...)
			switch arm {
			case "directory":
				keyPath = private
			case "short":
				raw = raw[:31]
			case "long":
				raw = append(raw, raw[0])
			}
			if arm != "missing" && arm != "directory" {
				if err := os.WriteFile(keyPath, raw, 0o600); err != nil {
					t.Fatal(err)
				}
			}
			session, manifest := filepath.Join(private, "session"), filepath.Join(private, "manifest.json")
			if err := os.WriteFile(session, []byte(r.token), 0o600); err != nil {
				t.Fatal(err)
			}
			// Count only S4 envelope commits, observing the history immediately before
			// S4. The proxy forwards all real daemon responses byte-for-byte.
			writes := 0
			var proxyMu sync.Mutex
			var before, beforeS4 wireWorld
			proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				proxyMu.Lock()
				defer proxyMu.Unlock()
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Error(err)
					http.Error(w, "read", 500)
					return
				}
				var commit demoCommit
				if req.URL.Path == "/v1/commit" && json.Unmarshal(body, &commit) == nil && len(commit.Objects) == 1 && commit.Objects[0].SemanticID == evidence.ProofSemanticID {
					writes++
					beforeS4, _ = proofHeadEntry(t, ctx, c)
					if arm == "safe-refusal" {
						w.WriteHeader(http.StatusForbidden)
						_, _ = io.WriteString(w, `{"error":{"class":"Denied"}}`)
						return
					}
				}
				if req.URL.Path == "/a2a/" {
					before, _ = proofHeadEntry(t, ctx, c)
				}
				forward, err := http.NewRequestWithContext(req.Context(), req.Method, r.url+req.URL.RequestURI(), bytes.NewReader(body))
				if err != nil {
					t.Error(err)
					http.Error(w, "build", 500)
					return
				}
				forward.Header = req.Header.Clone()
				resp, err := http.DefaultClient.Do(forward)
				if err != nil {
					t.Error(err)
					http.Error(w, "forward", 500)
					return
				}
				defer resp.Body.Close()
				if writes > 0 && req.URL.Path == "/v1/commit" && arm == "lost-response-after-commit" {
					_, _ = io.Copy(io.Discard, resp.Body)
					if resp.StatusCode != http.StatusOK {
						t.Errorf("lost ack requires successful real commit: HTTP %d", resp.StatusCode)
					}
					conn, _, err := w.(http.Hijacker).Hijack()
					if err != nil {
						t.Error(err)
						return
					}
					_ = conn.Close()
					return
				}
				for name, vals := range resp.Header {
					w.Header()[name] = vals
				}
				w.WriteHeader(resp.StatusCode)
				_, _ = io.Copy(w, resp.Body)
			}))
			defer proxy.Close()
			options := &demoTestRig{addr: proxy.URL}
			args := []string{"demo", "seed", "--db", r.db, "--session", session, "--ailang-bin", interp, "--proof-key-file", keyPath, "--manifest", manifest}
			// Observe the actual seed call, then delegate unchanged to the real producer.
			produce := demoSeedProve
			defer func() { demoSeedProve = produce }()
			proofCalls, readerRequests, writerRequests := 0, 0, 0
			demoSeedProve = func(ctx context.Context, reader evidence.ObjectReader, writer evidence.ObjectWriter, h daemon.HealthResponse, key [32]byte, source hashref.HashRef, required []string) (hashref.HashRef, error) {
				proofCalls++
				t.Run("required-identities", func(t *testing.T) {
					want := []string{"sameRef", "servesEntry"}
					if !slices.Equal(demoSeedRequiredIdentities, want) || !slices.Equal(required, want) {
						t.Fatalf("seed-required-identities: S4 must require exactly [sameRef servesEntry]; configured=%v actual=%v", demoSeedRequiredIdentities, required)
					}
				})
				r, rok := reader.(demoProofReader)
				w, wok := writer.(demoProofWriter)
				t.Run("http-bounds", func(t *testing.T) {
					if !rok || !wok || r.c == nil || w.c == nil {
						t.Fatal("seed-s4-http-bounds: seed must construct real REST adapters")
					}
					for _, c := range []*client{r.c, w.c} {
						if c.timeout != demoProofHTTPTimeout || c.http.Timeout != demoProofHTTPTimeout {
							t.Fatalf("seed-s4-http-bounds: both adapters need 10s context/client bounds; got %s/%s", c.timeout, c.http.Timeout)
						}
					}
				})
				// Check the contexts at the transport too, for both source reads and
				// the envelope write. Preserve the seed's original forwarding transport.
				clients := []*client{r.c}
				if w.c != r.c {
					clients = append(clients, w.c)
				}
				for _, c := range clients {
					transport := c.http.Transport
					if transport == nil {
						transport = http.DefaultTransport
					}
					c.http.Transport = scriptRT(func(req *http.Request) (*http.Response, error) {
						deadline, ok := req.Context().Deadline()
						if !ok || time.Until(deadline) > demoProofHTTPTimeout || time.Until(deadline) < demoProofHTTPTimeout-time.Second {
							t.Error("seed-s4-request-context: actual REST request must carry the 10s deadline")
						}
						if req.Method == http.MethodGet && strings.HasPrefix(req.URL.Path, "/v1/objects/") {
							readerRequests++
						}
						if req.Method == http.MethodPost && req.URL.Path == "/v1/commit" {
							writerRequests++
						}
						return transport.RoundTrip(req)
					})
				}
				return produce(ctx, reader, writer, h, key, source, required)
			}
			code, out, errOut := demoRun(t, options, args)
			proxyMu.Lock()
			defer proxyMu.Unlock()
			assertDemoProofKeyPrivate(t, map[string]string{
				"stdout": out, "stderr": errOut, "manifest-path": manifest,
			}, key[:], raw)
			if contents, err := os.ReadFile(manifest); err == nil {
				assertDemoProofKeyPrivate(t, map[string]string{"manifest-contents": string(contents)}, key[:], raw)
			}
			if code != exitUsage || out != "" {
				t.Fatalf("seed must refuse completion without stdout: code=%d stdout=%q stderr=%q", code, out, errOut)
			}
			if _, err := os.Stat(manifest); !os.IsNotExist(err) {
				t.Fatal("seed through S4 must never create a success manifest")
			}
			entries, err := (worldReader{ctx: ctx, c: c}).page(0, 500)
			if err != nil {
				t.Fatal(err)
			}
			var envelopes []hashref.HashRef
			for _, entry := range entries {
				ref := hashref.MustParse(entry.TransitionRef)
				meta, _, err := (demoProofReader{c: c}).ReadObject(ctx, ref, 8<<20)
				if err != nil {
					t.Fatal(err)
				}
				if meta.SemanticID == evidence.ProofSemanticID {
					envelopes = append(envelopes, ref)
				}
			}
			if arm == "success" {
				if writes != 1 || len(envelopes) != 1 || !strings.Contains(errOut, "S1-S6 checked; roster checks M5c pending; no success manifest written") {
					t.Fatalf("seed must record exactly one S4 envelope and refuse S6 completion: writes=%d envelopes=%d err=%q", writes, len(envelopes), errOut)
				}
				if proofCalls != 1 || readerRequests == 0 || writerRequests != 1 {
					t.Fatalf("seed-s4-observed: need one real producer call with bounded reads and write; calls=%d reads=%d writes=%d", proofCalls, readerRequests, writerRequests)
				}
				source := demoProofSource(t)
				h := proofHealth(t, ctx, c)
				// This validator requires both identities, as --proof-required-identity does.
				requireDemoProven(t, ctx, proofValidator(t, c, h, key), envelopes[0], hashref.SumSHA256(source))
				_, payload, err := (demoProofReader{c: c}).ReadObject(ctx, envelopes[0], evidence.MaxBytes)
				if err != nil {
					t.Fatal(err)
				}
				envelope, err := evidence.DecodeAuthenticatedEnvelope(payload)
				if err != nil {
					t.Fatal(err)
				}
				report, err := evidence.DecodeProofReportV1(envelope.Report)
				if err != nil || !slices.Contains(report.Verified, "sameRef") || !slices.Contains(report.Verified, "servesEntry") {
					t.Fatalf("seed-verified-identities: envelope must verify both sameRef and servesEntry; got %v err=%v", report.Verified, err)
				}
			} else if arm == "safe-refusal" || arm == "lost-response-after-commit" {
				wantEnvelopes := 0
				if arm == "lost-response-after-commit" {
					wantEnvelopes = 1
				}
				after, _ := proofHeadEntry(t, ctx, c)
				unchanged := after == beforeS4
				if arm == "lost-response-after-commit" {
					unchanged = after.Revision == beforeS4.Revision+1
				}
				if !unchanged || writes != 1 || len(envelopes) != wantEnvelopes || !strings.Contains(errOut, "S4 refused:") || !strings.Contains(errOut, "inspect history") || !strings.Contains(errOut, "no success manifest written") {
					t.Fatalf("seed uncertain/refused S4 must stop with one write and no retry/reseed/manifest: writes=%d envelopes=%d err=%q", writes, len(envelopes), errOut)
				}
			} else {
				want := map[string]string{"missing": "proof key: open:", "directory": "not a regular file", "short": "holds 31 bytes", "long": "holds more than 32 bytes"}[arm]
				if !strings.Contains(errOut, want) || writes != 0 || len(envelopes) != 0 {
					t.Fatalf("key %s refusal must precede every S4 write: writes=%d envelopes=%d err=%q want=%q", arm, writes, len(envelopes), errOut, want)
				}
				after, _ := proofHeadEntry(t, ctx, c)
				// The last A2A call itself commits S3; it must be the final history entry.
				if after.Revision != before.Revision+1 {
					t.Fatalf("key refusal changed history after S3: before=%d after=%d", before.Revision, after.Revision)
				}
			}
		})
	}
}

func TestDemoSeedRequiredIdentities(t *testing.T) {
	if !slices.Equal(demoSeedRequiredIdentities, []string{"sameRef", "servesEntry"}) {
		t.Fatalf("seed-required-identities: seed configuration must equal [sameRef servesEntry]; got %v", demoSeedRequiredIdentities)
	}
}

// Publish the checked-in descriptor through the same production publisher as
// se-tools. Only the negative control changes source, in a private temp module.
func publishDemoApproval(t *testing.T, st *store.Store, db string, interp hashref.HashRef, extra bool) {
	t.Helper()
	raw, err := os.ReadFile("../../packages/world-demo/transitions.json")
	if err != nil {
		t.Fatal(err)
	}
	var entries []struct {
		ID, Title, Description, TransitionFnFile string
		InputSchema, OutputSchema                json.RawMessage
		Access                                   transitionreg.EffectRequirement
		DeclaredEffects                          []transitionreg.EffectRequirement
	}
	if err := json.Unmarshal(raw, &entries); err != nil || len(entries) != 1 {
		t.Fatalf("demo descriptor: %v (%d entries)", err, len(entries))
	}
	e := entries[0]
	if e.Access.Cost != 0 || len(e.DeclaredEffects) != 1 || e.DeclaredEffects[0].Cost != 1 {
		t.Fatalf("descriptor access cost must be 0, declared approval cost 1: %+v", e)
	}
	src, err := os.ReadFile(filepath.Join("../..", e.TransitionFnFile))
	if err != nil {
		t.Fatal(err)
	}
	if extra {
		one := `jo([kv("id", js("e1")), kv("effect", js("Human.Approve")),
      kv("scope", js("worktree")), kv("cost", jint(1)),
      kv("payload", jo([kv("requester", js("demo"))]))])`
		two := one + ", " + strings.Replace(one, `js("e1")`, `js("e2")`, 1)
		if strings.Count(string(src), one) != 1 {
			t.Fatal("extra-effect control must replace exactly one plan effect")
		}
		src = []byte(strings.Replace(string(src), one, two, 1))
		path := filepath.Join(t.TempDir(), "approve.ail")
		if err := os.WriteFile(path, src, 0600); err != nil {
			t.Fatal(err)
		}
		cmd := exec.CommandContext(boundedTestContext(t), os.Getenv("AILANG_BIN"), "check", "approve.ail")
		cmd.Dir = filepath.Dir(path)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("valid two-effect test module must typecheck: %v %s", err, out)
		}
	}
	src, err = canon.Source(src)
	if err != nil {
		t.Fatal(err)
	}
	obj := store.Object{Hash: hashref.SumSHA256(src), InterfaceHash: hashref.SumSHA256([]byte("world/transition-source/v1")),
		SemanticID: "world/transition-source/v1", Provenance: "demo_test", Payload: src}
	ctx := boundedTestContext(t)
	if err := st.PutObject(ctx, obj); err != nil {
		t.Fatal(err)
	}
	in, err := transitionreg.CanonicalSchema(e.InputSchema)
	if err != nil {
		t.Fatal(err)
	}
	out, err := transitionreg.CanonicalSchema(e.OutputSchema)
	if err != nil {
		t.Fatal(err)
	}
	pub := transitionreg.NewPublisher(st, archive.New(db))
	epochs, _, err := pub.EpochsForInterpreter(ctx, interp)
	if err != nil || len(epochs) != 1 {
		t.Fatalf("demo epochs = %v %v", epochs, err)
	}
	desc := transitionreg.Descriptor{ID: e.ID, TransitionFn: obj.Hash, Interpreter: interp, SemanticsEpoch: epochs[0],
		InputSchema: in, OutputSchema: out, Access: e.Access, DeclaredEffects: e.DeclaredEffects, Title: e.Title, Description: e.Description}
	if _, err := pub.PublishSet(ctx, []transitionreg.Change{{ID: e.ID, Descriptor: &desc}}); err != nil {
		t.Fatal(err)
	}
}

func TestDemoSeedApprovalControls(t *testing.T) {
	proof := hashref.SumSHA256([]byte("preceding S4 envelope"))
	roster := demoApprovalRoster("test-record", "read-effect", "authored-object", proof)
	t.Run("roster-from-stage-refs", func(t *testing.T) {
		want := []demoEvidenceRef{{"test-invocation", "test-record"}, {"recorded-effect", "read-effect"},
			{"authored-claim", "authored-object"}, {"proof-envelope", proof.String()}}
		if !slices.Equal(roster, want) {
			t.Fatalf("roster must preserve preceding stage refs, including S4 envelope: got %+v", roster)
		}
	})
	t.Run("wrong-roster-length", func(t *testing.T) {
		bad := append(append([]demoEvidenceRef(nil), roster...), roster[0])
		_, err := demoApprovalInput("goal", []string{"a", "b"}, bad)
		if err == nil || !strings.Contains(err.Error(), "exactly four") {
			t.Fatalf("wrong roster length must refuse: %v", err)
		}
	})
	t.Run("wrong-roster-kind", func(t *testing.T) {
		bad := append([]demoEvidenceRef(nil), roster...)
		bad[0].Kind = "PROVEN"
		_, err := demoApprovalInput("goal", []string{"a", "b"}, bad)
		if err == nil || !strings.Contains(err.Error(), "roster kind") {
			t.Fatalf("wrong roster kind must refuse: %v", err)
		}
	})
	t.Run("invalid-roster-before-dispatch", func(t *testing.T) {
		c, calls := demoFakeA2A(t, "working", nil, nil)
		err := demoRequestApproval(boundedTestContext(t), c, "token", "goal", []string{"a", "b"}, nil)
		if *calls != 0 || err == nil || !strings.Contains(err.Error(), "exactly four") {
			t.Fatalf("invalid roster must refuse before dispatch: calls=%d err=%v", *calls, err)
		}
	})
	t.Run("missing-test-entry-before-dispatch", func(t *testing.T) {
		calls := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				calls++
			}
			http.NotFound(w, r)
		}))
		defer srv.Close()
		err := demoSeedApproval(boundedTestContext(t), newClient(srv.URL), "token", demoDomainObjects(nil), demoToolStages{}, proof)
		if calls != 0 || err == nil || !strings.Contains(err.Error(), "test invocation unavailable") {
			t.Fatalf("missing S3 entry must refuse before dispatch: calls=%d err=%v", calls, err)
		}
	})
	t.Run("approval-refusal-stops", func(t *testing.T) {
		c, calls := demoFakeA2A(t, "working", nil, nil)
		err := demoRequestApproval(boundedTestContext(t), c, "token", "goal", []string{"a", "b"}, roster)
		if *calls != 1 || err == nil || !strings.Contains(err.Error(), "S5 refused:") {
			t.Fatalf("S5 refusal must stop without retry: calls=%d err=%v", *calls, err)
		}
	})
	t.Run("s6-incomplete-refusal", func(t *testing.T) {
		err := demoApprovalIncomplete()
		if err == nil || !strings.Contains(err.Error(), "S1-S5 recorded; S6 incomplete; no success manifest written") {
			t.Fatalf("S6 incomplete must refuse success: %v", err)
		}
	})
}

func demoRESTObjects(t *testing.T, r worldReader, semantic string) []wireObject {
	t.Helper()
	var page struct{ Items []wireObject }
	if err := demoGet(r, "/v1/objects/by-semantic-id/"+semantic+"?limit=500", &page); err != nil {
		t.Fatal(err)
	}
	for i, obj := range page.Items {
		page.Items[i] = demoRESTObject(t, r, obj.Hash, semantic)
	}
	return page.Items
}

func demoRESTObject(t *testing.T, r worldReader, ref, semantic string) wireObject {
	t.Helper()
	o, ok, reason := r.checkedObject(ref, semantic)
	if !ok || o.Hash != ref {
		t.Fatalf("REST content hash/semantic mismatch %s: %s", ref, reason)
	}
	return o
}

// The fixture mints one grant with this initial budget. Read all immutable
// accounting records over the daemon REST read path on BOTH sides of S5;
// subtract spend independently, then cross-check the broker's before/after.
// No live database handle or A2A return value supplies these observations.
func demoRESTApprovalBudget(t *testing.T, r worldReader, initial int64) int64 {
	t.Helper()
	remaining := initial
	for _, obj := range demoRESTObjects(t, r, broker.EffectRecordV1) {
		rec, err := broker.DecodeRecord(obj.Payload)
		if err != nil {
			t.Fatal(err)
		}
		if rec.Effect == broker.EffectHumanApprove && rec.Allowed {
			remaining -= rec.Cost
		}
	}
	return remaining
}

// Pure assertion predicate has independent negative controls: widening the
// exact delta to >= must fail even if the real daemon still spends exactly 1.
func demoApprovalBudgetDelta(before, after int64) bool { return before-after == 1 }

func TestDemoSeedApprovalBudgetControl(t *testing.T) {
	for _, d := range []int64{0, 2, -1} {
		if demoApprovalBudgetDelta(20, 20-d) {
			t.Fatalf("budget delta must equal exactly 1; accepted delta %d", d)
		}
	}
	if !demoApprovalBudgetDelta(20, 19) {
		t.Fatal("budget delta 1 must pass")
	}
}

func TestDemoSeedApprovalPacket(t *testing.T) {
	interp := requireProofPin(t)
	_, tool := cliSeBins(t) // unset tool pin skips; a configured wrong release fails.
	requireBinaryRelease(t, tool, daemon.ToolBinaryRelease)
	for _, arm := range []string{"single-approval", "no-human-budget", "extra-effect-refused"} {
		t.Run(arm, func(t *testing.T) {
			budget := int64(20)
			if arm == "no-human-budget" {
				budget = 0
			}
			rig := newDemoApprovalRig(t, interp, tool, true, arm == "extra-effect-refused", budget)
			ctx := boundedTestContext(t)
			r := worldReader{ctx: ctx, c: newClient(rig.url)}
			private := t.TempDir()
			key := freshProofKey(t)
			keyPath, session, manifest := filepath.Join(private, "proof.key"), filepath.Join(private, "session"), filepath.Join(private, "manifest.json")
			if err := os.WriteFile(keyPath, key[:], 0600); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(session, []byte(rig.token), 0600); err != nil {
				t.Fatal(err)
			}
			before := demoRESTApprovalBudget(t, r, budget)
			beforeHead, beforeHasHead, err := headState(ctx, r.c)
			if err != nil {
				t.Fatal(err)
			}
			calls, decisions := 0, 0
			var preS5 wireWorld
			proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Error(err)
					return
				}
				if req.Method == http.MethodPost && strings.Contains(string(body), `"demo-approval"`) {
					calls++
					if calls == 1 {
						preS5, _ = proofHeadEntry(t, ctx, r.c)
						before = demoRESTApprovalBudget(t, r, budget)
					}
				}
				if req.Method != http.MethodGet && req.URL.Path != "/a2a/" && req.URL.Path != "/v1/commit" {
					decisions++
				}
				forward, err := http.NewRequestWithContext(req.Context(), req.Method, rig.url+req.URL.RequestURI(), bytes.NewReader(body))
				if err != nil {
					t.Error(err)
					return
				}
				forward.Header = req.Header.Clone()
				resp, err := http.DefaultClient.Do(forward)
				if err != nil {
					t.Error(err)
					return
				}
				defer resp.Body.Close()
				for k, v := range resp.Header {
					w.Header()[k] = v
				}
				w.WriteHeader(resp.StatusCode)
				_, _ = io.Copy(w, resp.Body)
			}))
			defer proxy.Close()
			code, out, stderr := demoRun(t, &demoTestRig{addr: proxy.URL}, []string{"demo", "seed", "--db", rig.db,
				"--session", session, "--ailang-bin", interp, "--proof-key-file", keyPath, "--manifest", manifest})
			if code != exitUsage || out != "" {
				t.Fatalf("seed must refuse completion: code=%d out=%q err=%q", code, out, stderr)
			}
			if _, err := os.Stat(manifest); !os.IsNotExist(err) {
				t.Fatal("S6 incomplete: no success manifest may be written")
			}
			packets := demoRESTObjects(t, r, coordinator.DecisionPacketV1)
			entries, err := r.page(0, 500)
			if err != nil {
				t.Fatal(err)
			}
			var approvalEntries []wireEntry
			var approvalRecords []*invRecord
			stageRecords := map[string]*invRecord{}
			stageRefs := map[string]string{}
			for _, entry := range entries {
				rec, obj, err := r.record(entry)
				if err != nil {
					t.Fatal(err)
				}
				if rec != nil {
					demoRESTObject(t, r, entry.TransitionRef, obj.SemanticID)
					stageRecords[rec.InvocationID] = rec
					stageRefs[rec.InvocationID] = entry.TransitionRef
					if rec.SkillID == "demo-approve" {
						approvalEntries = append(approvalEntries, entry)
						approvalRecords = append(approvalRecords, rec)
					}
				}
			}
			if arm == "no-human-budget" {
				afterHead, afterHasHead, err := headState(ctx, r.c)
				if len(packets) != 0 || calls != 0 || !strings.Contains(stderr, "Human.Approve@worktree grant") ||
					err != nil || afterHead != beforeHead || afterHasHead != beforeHasHead || len(entries) != 0 {
					t.Fatalf("missing Human.Approve budget must refuse before any packet or dispatch and preserve REST history: packets=%d calls=%d entries=%d head=%s before=%s read=%v err=%s", len(packets), calls, len(entries), afterHead, beforeHead, err, stderr)
				}
				return
			}
			if arm == "extra-effect-refused" {
				after, _ := proofHeadEntry(t, ctx, r.c)
				if calls != 1 || len(packets) != 0 || len(approvalEntries) != 0 || after != preS5 ||
					demoRESTApprovalBudget(t, r, budget) != before || !strings.Contains(rig.log.String(), "effect plan violates L1") {
					t.Fatalf("valid two-effect module must be refused without packet/spend/commit: packets=%d entries=%d calls=%d err=%s coordinator=%s", len(packets), len(approvalEntries), calls, stderr, rig.log.String())
				}
				t.Logf("coordinator refusal: %s", strings.TrimSpace(rig.log.String()))
				return
			}
			// This is deliberately first: the extra-effect mutant must red here,
			// with the real coordinator refusal included in the diagnostic.
			if len(approvalEntries) != 1 || len(packets) != 1 || len(approvalRecords[0].Effects) != 1 {
				t.Fatalf("exactly one S5 carrying entry/effect/packet required: entries=%d packets=%d; seed=%s coordinator=%s", len(approvalEntries), len(packets), stderr, rig.log.String())
			}
			rec := approvalRecords[0]
			var registryHead struct{ Head string }
			if err := demoGet(r, "/v1/registry/"+transitionreg.SemanticIDV1, &registryHead); err != nil {
				t.Fatal(err)
			}
			registryObj := demoRESTObject(t, r, registryHead.Head, transitionreg.SemanticIDV1)
			revision, err := transitionreg.DecodeRevision(registryObj.Payload)
			if err != nil {
				t.Fatal(err)
			}
			foundDescriptor := false
			for _, descriptor := range revision.Entries {
				if descriptor.ID == "demo-approve" {
					foundDescriptor = true
					if descriptor.Access.Cost != 0 || descriptor.TransitionFn.String() != rec.TransitionFn || descriptor.Interpreter.String() != rec.Interpreter {
						t.Fatalf("published descriptor must have cost 0 and bind S5 source/interpreter: %+v", descriptor)
					}
				}
			}
			if !foundDescriptor {
				t.Fatal("published registry lacks demo-approve descriptor")
			}
			demoRESTObject(t, r, rec.TransitionFn, "world/transition-source/v1")
			packet, err := coordinator.DecodeDecisionPacket(packets[0].Payload)
			if err != nil || packet.ProposalHash.String() != rec.Plan || approvalEntries[0].Header.EntryIndex != preS5.Revision+1 {
				t.Fatalf("packet proposal must be S5 carrying entry's plan: %v %+v", err, packet)
			}
			plan := demoRESTObject(t, r, rec.Plan, coordinator.EffectPlanV1)
			var planDoc struct {
				Effects []struct {
					ID, Effect, Scope string
					Cost              int64
					Payload           json.RawMessage
				}
				Finish bool
				Result any
			}
			if err := json.Unmarshal(plan.Payload, &planDoc); err != nil || len(planDoc.Effects) != 1 || planDoc.Finish || planDoc.Result != nil {
				t.Fatalf("real single-effect plan required: %s %v", plan.Payload, err)
			}
			e := planDoc.Effects[0]
			if e.ID != "e1" || e.Effect != broker.EffectHumanApprove || e.Scope != "worktree" || e.Cost != 1 || string(e.Payload) != `{"requester":"demo"}` {
				t.Fatalf("e1 Human.Approve plan mismatch: %s", plan.Payload)
			}
			effect := demoRESTObject(t, r, rec.Effects[0], broker.EffectRecordV1)
			accounting, err := broker.DecodeRecord(effect.Payload)
			if err != nil || !broker.RecordConsistent(accounting) || !accounting.Allowed || accounting.Failed {
				t.Fatalf("successful real broker record required: %+v %v", accounting, err)
			}
			effectRequest := demoRESTObject(t, r, accounting.RequestRef.String(), broker.EffectRequestV1)
			request := demoRESTObject(t, r, packet.RequestRef.String(), broker.ApprovalRequestV1)
			approval, err := broker.DecodeApprovalRequest(request.Payload)
			if err != nil || approval.Requester != "demo" || approval.Effect != e.Effect || approval.Scope != e.Scope || approval.Cost != e.Cost {
				t.Fatalf("packet approval request mismatch: %+v %v", approval, err)
			}
			wantRequest := append([]byte(fmt.Sprintf("%d:%s%d:%s%d:%d:", len(e.Effect), e.Effect, len(e.Scope), e.Scope, e.Cost, approval.Now)), e.Payload...)
			if !bytes.Equal(effectRequest.Payload, wantRequest) || approval.Now != packet.CreatedAt || accounting.Effect != e.Effect || accounting.Scope != e.Scope || accounting.Cost != e.Cost {
				t.Fatalf("effect request bytes must bind the plan payload, requirements and packet time: got %q want %q", effectRequest.Payload, wantRequest)
			}
			result := demoRESTObject(t, r, accounting.ResultRef.String(), broker.EffectResultV1)
			pending, err := broker.ParsePendingApproval(result.Payload)
			if err != nil || pending != packet.RequestRef {
				t.Fatalf("effect result must link pending request: %s %v", result.Payload, err)
			}
			output := demoRESTObject(t, r, rec.Output, coordinator.OutputV1)
			var outputDoc struct {
				Status, RequestRef string
				World              struct {
					Plan    string
					Effects []struct{ ID, Status, Record string }
				}
			}
			if err := json.Unmarshal(output.Payload, &outputDoc); err != nil || outputDoc.Status != "pending" || outputDoc.RequestRef != pending.String() || outputDoc.World.Plan != rec.Plan ||
				len(outputDoc.World.Effects) != 1 || outputDoc.World.Effects[0].ID != "e1" || outputDoc.World.Effects[0].Status != "ok" || outputDoc.World.Effects[0].Record != effect.Hash {
				t.Fatalf("committed output/e1/request/plan hash links mismatch: %s %v", output.Payload, err)
			}
			input := demoRESTObject(t, r, rec.Input, coordinator.InputV1)
			var inputDoc struct {
				Goal       string
				Candidates []string
				Roster     []demoEvidenceRef
			}
			if err := json.Unmarshal(input.Payload, &inputDoc); err != nil {
				t.Fatal(err)
			}
			testID := coordinator.InvocationID(coordinator.SurfaceA2A, "demo", demoTestTask)
			readID := coordinator.InvocationID(coordinator.SurfaceA2A, "demo", demoReadTask)
			envelopes := demoRESTObjects(t, r, evidence.ProofSemanticID)
			claims := demoRESTObjects(t, r, "world/demo/authored-claim/v1")
			if len(envelopes) != 1 || len(claims) != 1 || stageRecords[readID] == nil || stageRecords[testID] == nil {
				t.Fatal("preceding real S1-S4 refs must exist")
			}
			wantRoster := []demoEvidenceRef{{"test-invocation", stageRefs[testID]}, {"recorded-effect", stageRecords[readID].Effects[0]},
				{"authored-claim", claims[0].Hash}, {"proof-envelope", envelopes[0].Hash}}
			if !slices.Equal(inputDoc.Roster, wantRoster) {
				t.Fatalf("four roster refs must equal preceding real stages: got %+v want %+v", inputDoc.Roster, wantRoster)
			}
			goals := demoRESTObjects(t, r, "world/demo/goal/v1")
			proposals := demoRESTObjects(t, r, "world/demo/proposal/v1")
			var committedGoal struct{ Goal string }
			_ = json.Unmarshal(goals[0].Payload, &committedGoal)
			var descriptions []string
			for _, obj := range proposals {
				var p struct{ Description string }
				_ = json.Unmarshal(obj.Payload, &p)
				descriptions = append(descriptions, p.Description)
			}
			slices.Sort(descriptions)
			gotDescriptions := append([]string(nil), inputDoc.Candidates...)
			slices.Sort(gotDescriptions)
			if inputDoc.Goal != committedGoal.Goal || !slices.Equal(gotDescriptions, descriptions) {
				t.Fatalf("S5 must use S1 committed descriptions: %+v", inputDoc)
			}
			after := demoRESTApprovalBudget(t, r, budget)
			if !demoApprovalBudgetDelta(before, after) || accounting.BudgetBefore != before || accounting.BudgetAfter != after {
				t.Fatalf("Human.Approve remaining budget must decrease EXACTLY 1: before=%d after=%d broker=%+v", before, after, accounting)
			}
			t.Run("no-decision-verb-packet-stays-pending", func(t *testing.T) {
				status, page, err := r.c.do(ctx, http.MethodGet, "/workbench?object="+packets[0].Hash, nil)
				wantCarrier := fmt.Sprintf("decision packet validated: committed with its invocation record at entry #%d;", approvalEntries[0].Header.EntryIndex)
				if err != nil || status != http.StatusOK || !strings.Contains(string(page), wantCarrier) || !strings.Contains(string(page), "approvals chain: pending") {
					t.Fatalf("daemon packet validator must confirm S5 carrying entry and live pending approvals-chain anchor: status=%d err=%v page=%s", status, err, page)
				}
				if decisions != 0 || len(demoRESTObjects(t, r, broker.ApprovalDecisionV1)) != 0 || outputDoc.Status != "pending" {
					t.Fatalf("no decision verb may run; packet must stay pending: calls=%d", decisions)
				}
			})
			t.Run("seed-invokes-s5-once", func(t *testing.T) {
				if calls != 1 {
					t.Fatalf("seed must invoke demo-approval exactly once: calls=%d", calls)
				}
			})
			if !strings.Contains(stderr, "S1-S6 checked; roster checks M5c pending; no success manifest written") {
				t.Fatalf("successful S5 must still refuse S6 completion: %s", stderr)
			}
		})
	}
}

func demoFakeA2A(t *testing.T, state string, meta, data map[string]any) (*client, *int) {
	t.Helper()
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		calls++
		body := map[string]any{"jsonrpc": "2.0", "id": 9, "result": map[string]any{
			"status": map[string]any{"state": state}, "metadata": meta,
			"artifacts": []any{map[string]any{"parts": []any{map[string]any{"type": "data", "data": data}}}}}}
		_ = json.NewEncoder(w).Encode(body)
	}))
	t.Cleanup(srv.Close)
	return newClient(srv.URL), &calls
}

// All omission controls start with JSON emitted by the real CLI against a
// temporary daemon, so the validator cannot grow its own provenance dialect.
func demoM4bWhyFixture(t *testing.T) (*walkRig, whyChain, []byte) {
	t.Helper()
	t.Setenv("TMPDIR", t.TempDir())
	rig := newWalkRig(t)
	rig.invoke("demo", "ailang-read", "m4b", "actual stored result")
	code, out, stderr := whyRun(t, rig.url, "", "1", "--json")
	if code != exitOK {
		t.Fatalf("real why fixture: %d %s", code, stderr)
	}
	var ch whyChain
	if err := json.Unmarshal([]byte(out), &ch); err != nil {
		t.Fatal(err)
	}
	if err := demoValidateWhy(ch); err != nil {
		t.Fatalf("real why chain accepted: %v", err)
	}
	return rig, ch, []byte(out)
}

func demoM4bRefusal(t *testing.T, err error, assertion string) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), assertion) {
		t.Fatalf("%s: %v", assertion, err)
	}
}

func TestDemoSeedWhyNamedLinks(t *testing.T) {
	t.Run("controls", func(t *testing.T) {
		rig, ch, _ := demoM4bWhyFixture(t)

		t.Run("accepted-report", func(t *testing.T) {
			accepted := ch // M4b-link stimulus: a decoded real report
			if err := demoValidateWhy(accepted); err != nil {
				t.Fatal(err)
			}
		})
		t.Run("approval-stimulus", func(t *testing.T) {
			rig.invoke("demo", "demo-approve", "approval-stimulus", "stored approval stimulus") // M4b-plain-rest stimulus
			r := worldReader{ctx: boundedTestContext(t), c: newClient(rig.url)}
			chain, _, err := demoWhy(r, rig.next-1)
			if err != nil {
				t.Fatal(err)
			}
			if err := demoValidateWhy(chain); err != nil {
				t.Fatal(err)
			}
		})
		for _, name := range []string{"coordinator", "ok", "zero-links", "effect-absent", "world-absent", "entry-absent", "record-absent", "input-absent", "plan-absent", "output-absent", "effect-false"} {
			t.Run(name, func(t *testing.T) {
				bad := ch
				bad.Links = append([]whyLink(nil), ch.Links...)
				want := ""
				switch name {
				case "coordinator":
					bad.Coordinator = false
					bad.OK = false
					bad.Links = nil
					want = "coordinator must be true"
				case "ok":
					bad.OK = false
					want = "ok must be true"
				case "zero-links":
					bad.Links = nil
					want = "zero links refused"
				default:
					link := strings.Split(name, "-")[0]
					if strings.HasSuffix(name, "-absent") {
						bad.Links = slices.DeleteFunc(bad.Links, func(l whyLink) bool { return l.Name == link })
					} else {
						for i := range bad.Links {
							if bad.Links[i].Name == link {
								bad.Links[i].OK = false
							}
						}
					}
					want = "expected " + link + " link must exist and be true"
				}
				demoM4bRefusal(t, demoValidateWhy(bad), want)
			})
		}
		t.Run("plain-rest", func(t *testing.T) {
			// A valid ordinary object replaces the approval stimulus. Read it back,
			// and confirm the preceding coordinator entry survived unchanged.
			r := worldReader{ctx: boundedTestContext(t), c: newClient(rig.url)}
			before, _, err := r.entry(1)
			if err != nil {
				t.Fatal(err)
			}
			commit := makeCLICommit(rig.tip, rig.next, "ordinary demo object")
			rig.post(commit)
			actual, _, err := r.entry(rig.next)
			if err != nil || actual.EntryHash != commit.Entry.EntryHash {
				t.Fatalf("ordinary REST commit readback mismatch: %v", err)
			}
			after, _, err := r.entry(1)
			if err != nil || after != before {
				t.Fatal("ordinary REST commit changed preceding entry")
			}
			ordinary, _, err := demoWhy(r, rig.next)
			if err != nil {
				t.Fatal(err)
			}
			// This calls the actual validator before any packet count.
			demoM4bRefusal(t, demoValidateWhy(ordinary), "coordinator must be true")
			_, err = demoSeedS6(r, "", nil, demoToolStages{Read: demoToolRun{EntryIndex: 1}, Test: demoToolRun{EntryIndex: 1}, Check: demoToolRun{EntryIndex: 1}}, demoToolRun{EntryIndex: rig.next}, nil)
			demoM4bRefusal(t, err, "coordinator must be true")
		})
	})
	t.Run("real-stack", func(t *testing.T) {
		r, m := demoM4bRealSeed(t)
		if len(m.Why) != 4 {
			t.Fatal("read/test/check/approval must each have a why chain")
		}
		for _, w := range m.Why {
			ch, raw, err := demoWhy(r, w.Entry)
			if err != nil {
				t.Fatal(err)
			}
			if err := demoValidateWhy(ch); err != nil {
				t.Fatalf("real named links accepted: %v", err)
			}
			if sum(raw) != w.SHA256 {
				t.Fatal("raw why JSON hash retained")
			}
			for _, name := range linkNames {
				found := false
				for _, l := range ch.Links {
					if l.Name == name && l.OK {
						found = true
					}
				}
				if !found {
					t.Fatalf("expected %s link must exist and be true", name)
				}
			}
		}
	})
}

// Executes the same S1-S5 composition as TestDemoSeedApprovalPacket, including
// the production proof producer and HumanHandler; then independently reassembles.
func demoM4bRealSeed(t *testing.T) (worldReader, demoPublicManifest) {
	t.Helper()
	t.Setenv("TMPDIR", t.TempDir())
	interp := requireProofPin(t)
	_, tool := cliSeBins(t)
	requireBinaryRelease(t, tool, daemon.ToolBinaryRelease)
	rig := newDemoApprovalRig(t, interp, tool, true, false, 20)
	private := t.TempDir()
	key := freshProofKey(t)
	keyPath, session, manifest := filepath.Join(private, "proof.key"), filepath.Join(private, "session"), filepath.Join(private, "manifest.json")
	if err := os.WriteFile(keyPath, key[:], 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(session, []byte(rig.token), 0600); err != nil {
		t.Fatal(err)
	}
	code, out, stderr := demoRun(t, &demoTestRig{addr: rig.url}, []string{"demo", "seed", "--db", rig.db, "--session", session, "--ailang-bin", interp, "--proof-key-file", keyPath, "--manifest", manifest})
	// Keep the plain-REST mutant at the coordinator assertion before any
	// generic completion or packet-count assertion can mask it.
	if strings.Contains(stderr, "coordinator must be true") {
		t.Fatalf("coordinator must be true for every real seed entry: %s", stderr)
	}
	if strings.Contains(stderr, "manifest completeness:") {
		t.Fatalf("manifest completeness: real seed assembly must contain every required role: %s", stderr)
	}
	if code != exitUsage || out != "" || !strings.Contains(stderr, "S1-S6 checked; roster checks M5c pending; no success manifest written") {
		t.Fatalf("S1-S6 checked fail-closed boundary: code=%d out=%q err=%s", code, out, stderr)
	}
	if _, err := os.Stat(manifest); !os.IsNotExist(err) {
		t.Fatal("S6 never writes success manifest")
	}
	r := worldReader{ctx: boundedTestContext(t), c: newClient(rig.url)}
	entries, err := r.page(0, 500)
	if err != nil {
		t.Fatal(err)
	}
	var stages demoToolStages
	var approval demoToolRun
	var roster []demoEvidenceRef
	for _, e := range entries {
		rec, _, err := r.record(e)
		if err != nil {
			t.Fatal(err)
		}
		if rec == nil {
			continue
		}
		run := demoToolRun{EntryIndex: e.Header.EntryIndex, InvocationID: rec.InvocationID}
		switch rec.SkillID {
		case "ailang-read":
			stages.Read = run
		case "ailang-cli":
			if strings.HasSuffix(rec.InvocationID, ":"+demoTestTask) {
				stages.Test = run
			} else {
				stages.Check = run
			}
		case "demo-approve":
			approval = run
			obj := demoRESTObject(t, r, rec.Input, coordinator.InputV1)
			var in struct{ Roster []demoEvidenceRef }
			if err := json.Unmarshal(obj.Payload, &in); err != nil {
				t.Fatal(err)
			}
			roster = in.Roster
		}
	}
	m, err := demoSeedS6(r, rig.token, demoDomainObjects(demoProofSource(t)), stages, approval, roster)
	if err != nil {
		t.Fatalf("real S6 assembly: %v", err)
	}
	return r, m
}

// The fixture commits real REST content to a unit daemon. Only health (the
// synthetic interpreter), A2A completion and workbench HTML are proxy responses;
// all manifest object/world/entry bytes come from the daemon's store.
type demoM4bManifestRig struct {
	r        worldReader
	objects  []demoObject
	stages   demoToolStages
	approval demoToolRun
	roster   []demoEvidenceRef
	token    string
}

func demoM4bCanonicalCommit(c cliCommit) cliCommit {
	h := c.Entry.Header
	c.Entry.EntryHash = sumJSON(entryWire{h.EntryIndex, h.SemanticsEpoch, h.TransitionFn, h.Interpreter, h.PrevEntryHash, h.WrittenBy, c.Entry.TransitionRef})
	c.NextWorld.LogHead = c.Entry.EntryHash
	c.NextWorld.Ref = sumJSON(worldWire{c.NextWorld.Revision, c.NextWorld.StateRoot, c.NextWorld.LogHead})
	return c
}

func demoM4bManifestFixture(t *testing.T, packets int, anchor string) demoM4bManifestRig {
	t.Helper()
	db := filepath.Join(t.TempDir(), "world.db")
	token := testCLISession(t, db)
	d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: db, BindHost: daemon.DefaultBindHost, ErrorLog: io.Discard})
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(d.Handler())
	t.Cleanup(func() { srv.Close(); _ = d.Close() })
	rig := &walkRig{t: t, url: srv.URL, token: token, invs: map[int64]synthInv{}}
	genesis := demoM4bCanonicalCommit(makeCLICommit(cliWorld{}, 0, "genesis"))
	rig.post(genesis)
	rig.tip, rig.next = genesis.NextWorld, 1

	skills := []string{"ailang-read", "ailang-cli", "ailang-cli", "demo-approve"}
	for i, skill := range skills {
		rig.invoke("demo", skill, fmt.Sprintf("manifest-%d", i), fmt.Sprintf("stored result %d", i))
	}
	objects := demoDomainObjects([]byte("unit daemon source"))
	packet := synthObject(coordinator.DecisionPacketV1, "host/coordinator", coordinator.EncodeDecisionPacket(coordinator.DecisionPacket{
		ProposalHash: hashref.MustParse(rig.invs[4].plan), RequestRef: hashref.MustParse(rig.invs[4].input), CreatedAt: 1700, DeadlineAt: 1760, Policy: "Cancel",
	}))
	items := []cliObject{}
	for i := 0; i < packets; i++ {
		item := packet
		if i > 0 {
			item = synthObject(coordinator.DecisionPacketV1, "host/coordinator", coordinator.EncodeDecisionPacket(coordinator.DecisionPacket{ProposalHash: hashref.MustParse(rig.invs[4].plan), RequestRef: hashref.MustParse(rig.invs[4].input), CreatedAt: 1700 + int64(i), DeadlineAt: 1760, Policy: "Cancel"}))
		}
		items = append(items, item)
	}
	commit := makeCLICommit(rig.tip, rig.next, "manifest fixture content")
	for _, o := range objects {
		commit.Objects = append(commit.Objects, cliObject{Hash: o.Hash, InterfaceHash: o.InterfaceHash, SemanticID: o.SemanticID, Provenance: o.Provenance, Payload: o.Payload})
	}
	for _, skill := range []string{"ailang-read", "ailang-cli", "demo-approve"} {
		commit.Objects = append(commit.Objects, synthObject("fixture/transition", "fixture", []byte("fn-"+skill)))
	}
	rig.post(demoM4bCanonicalCommit(commit))
	// The public writer correctly forbids minting coordinator-only packets.
	// Seed this immutable codec fixture through the test store, then resolve it
	// exclusively over REST in assembly and validation.
	srv.Close()
	if err := d.Close(); err != nil {
		t.Fatal(err)
	}
	st, err := store.Open(db)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range items {
		err = st.PutObject(boundedTestContext(t), store.Object{Hash: hashref.MustParse(item.Hash), InterfaceHash: hashref.MustParse(item.InterfaceHash), SemanticID: item.SemanticID, Provenance: item.Provenance, Payload: item.Payload})
		if err != nil {
			_ = st.Close()
			t.Fatal(err)
		}
	}
	closeErr := st.Close()
	if err != nil {
		t.Fatal(err)
	}
	if closeErr != nil {
		t.Fatal(closeErr)
	}
	d, err = daemon.New(boundedTestContext(t), daemon.Config{DBPath: db, BindHost: daemon.DefaultBindHost, ErrorLog: io.Discard})
	if err != nil {
		t.Fatal(err)
	}
	srv = httptest.NewServer(d.Handler())
	rig.url = srv.URL
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/v1/health":
			_ = json.NewEncoder(w).Encode(healthInfo{InterpreterRef: sum([]byte("interpreter"))})
			return
		case "/v1/objects/by-semantic-id/" + coordinator.DecisionPacketV1:

			_ = json.NewEncoder(w).Encode(map[string]any{"items": items})
			return
		case "/workbench":
			_, _ = io.WriteString(w, anchor)
			return
		case "/a2a/":
			_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": 9, "result": map[string]any{
				"status": map[string]any{"state": "completed"}, "metadata": map[string]any{"invocation_id": rig.invs[4].id, "entry_index": 4},
				"artifacts": []any{map[string]any{"parts": []any{map[string]any{"type": "data", "data": map[string]any{"world": map[string]any{"effects": []any{map[string]any{"record": rig.invs[4].effect}}}}}}}},
			}})
			return
		}
		resp, err := http.Get(rig.url + req.URL.RequestURI())
		if err != nil {
			t.Error(err)
			http.Error(w, "upstream", 500)
			return
		}
		defer resp.Body.Close()
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	}))
	t.Cleanup(proxy.Close)
	stages := demoToolStages{Read: demoToolRun{EntryIndex: 1, EffectRecord: rig.invs[1].effect}, Test: demoToolRun{EntryIndex: 2}, Check: demoToolRun{EntryIndex: 3}}
	return demoM4bManifestRig{worldReader{ctx: boundedTestContext(t), c: newClient(proxy.URL)}, objects, stages, demoToolRun{EntryIndex: 4}, demoApprovalRoster(rig.invs[2].record, rig.invs[1].effect, objects[4].Hash, hashref.MustParse(rig.invs[3].output)), rig.token}
}

// This expectation deliberately does not call production completeness or derive
// required roles from the assembled manifest. Each fixture/real run has one effect.
func demoM4bAssertManifestRoles(t *testing.T, r worldReader, m demoPublicManifest) {
	t.Helper()
	expected := map[string]int{}
	for _, role := range strings.Fields("frozen-head goal candidate-a candidate-b source claim test-invocation recorded-effect authored-claim proof-envelope packet request proposal") {
		expected[role] = 1
	}
	for _, role := range strings.Fields("invocation transition interpreter world record input plan effect output effect-request effect-result") {
		expected[role] = 4
	}
	entries, err := r.page(0, 500)
	if err != nil {
		t.Fatal(err)
	}
	expected["entry"], expected["object"] = len(entries), len(entries)
	actual := map[string]int{}
	for _, ref := range m.References {
		actual[ref.Role]++
	}
	for role, n := range expected {
		if actual[role] != n {
			t.Fatalf("manifest completeness: role %s requires %d references, got %d", role, n, actual[role])
		}
	}
	for role, n := range actual {
		if expected[role] != n {
			t.Fatalf("manifest completeness: extra role %s: %d", role, n)
		}
	}
	why := map[string]int{}
	for _, w := range m.Why {
		why[w.Role]++
	}
	for _, role := range []string{"read", "test", "check", "approval"} {
		if why[role] != 1 {
			t.Fatalf("manifest completeness: role why-%s requires exactly one hash", role)
		}
	}
	if len(m.Why) != 4 {
		t.Fatal("manifest completeness: exactly four why hashes required")
	}
}

func TestDemoSeedManifestReferences(t *testing.T) {
	t.Run("controls", func(t *testing.T) {
		fixture := demoM4bManifestFixture(t, 1, "decision packet validated: committed with its invocation record at entry #4;")
		r := fixture.r
		// Use an input ref for the existing corruption controls.
		e, _, err := r.entry(1)
		if err != nil {
			t.Fatal(err)
		}
		rec, _, err := r.record(e)
		if err != nil {
			t.Fatal(err)
		}
		inv := struct{ input string }{rec.Input}

		// Obtain packet refs from the real REST payload, independently of assembly.
		packets := demoRESTObjects(t, r, coordinator.DecisionPacketV1)
		packet, err := coordinator.DecodeDecisionPacket(packets[0].Payload)
		if err != nil {
			t.Fatal(err)
		}
		m, err := demoAssembleManifest(r, fixture.objects, []demoToolRun{fixture.stages.Read, fixture.stages.Test, fixture.stages.Check, fixture.approval}, fixture.roster, packets[0].Hash, packet.RequestRef.String(), packet.ProposalHash.String())
		if err != nil {
			t.Fatal(err)
		}
		demoM4bAssertManifestRoles(t, r, m)
		if err := demoValidateManifest(r, m, fixture.token); err != nil {
			t.Fatalf("complete manifest accepted: %v", err)
		}
		t.Run("completeness", func(t *testing.T) {
			for _, ref := range m.References {
				if ref.Role == "entry" || ref.Role == "object" {
					continue
				}
				ref := ref
				t.Run(fmt.Sprintf("missing-%s-%d", ref.Role, ref.Index), func(t *testing.T) {
					bad := m
					bad.References = slices.DeleteFunc(slices.Clone(m.References), func(r demoManifestRef) bool { return r.Role == ref.Role && r.Index == ref.Index })
					demoM4bRefusal(t, demoValidateManifest(demoM4bFreshReader(t, r), bad, fixture.token), "manifest completeness: role "+ref.Role)
				})
				t.Run(fmt.Sprintf("extra-%s-%d", ref.Role, ref.Index), func(t *testing.T) {
					bad := m
					bad.References = append(slices.Clone(m.References), ref)
					demoM4bRefusal(t, demoValidateManifest(demoM4bFreshReader(t, r), bad, fixture.token), "manifest completeness: role "+ref.Role)
				})
			}
			for _, w := range m.Why {
				for _, extra := range []bool{false, true} {
					t.Run(fmt.Sprintf("why-%s-extra-%t", w.Role, extra), func(t *testing.T) {
						bad := m
						bad.Why = slices.Clone(m.Why)
						if extra {
							bad.Why = append(bad.Why, w)
						} else {
							bad.Why = slices.DeleteFunc(bad.Why, func(v demoManifestWhy) bool { return v.Role == w.Role })
						}
						demoM4bRefusal(t, demoValidateManifest(demoM4bFreshReader(t, r), bad, fixture.token), "manifest completeness: role why-"+w.Role)
					})
				}
			}
		})
		t.Run("absent-ref", func(t *testing.T) {
			bad := m
			bad.References = append([]demoManifestRef(nil), m.References...)
			bad.References[0].Ref = sum([]byte("valid absent manifest reference"))
			demoM4bRefusal(t, demoValidateManifest(demoM4bFreshReader(t, r), bad, fixture.token), "every declared manifest reference resolves to matching content")
		})
		t.Run("mismatched-content", func(t *testing.T) {
			// Corrupt the read response only; compare immutable daemon bytes again after.
			before := demoRESTObject(t, demoM4bFreshReader(t, r), inv.input, coordinator.InputV1)
			proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if strings.HasPrefix(req.URL.Path, "/v1/objects/"+inv.input) {
					bad := before
					bad.Payload = []byte("mismatched bytes")
					_ = json.NewEncoder(w).Encode(bad)
					return
				}
				resp, err := http.Get(r.c.base + req.URL.RequestURI())
				if err != nil {
					t.Error(err)
					return
				}
				defer resp.Body.Close()
				w.WriteHeader(resp.StatusCode)
				_, _ = io.Copy(w, resp.Body)
			}))
			defer proxy.Close()
			badReader := worldReader{ctx: boundedTestContext(t), c: newClient(proxy.URL)}
			demoM4bRefusal(t, demoValidateManifest(badReader, m, fixture.token), "every declared manifest reference resolves to matching content")
			after := demoRESTObject(t, demoM4bFreshReader(t, r), inv.input, coordinator.InputV1)
			if !bytes.Equal(before.Payload, after.Payload) {
				t.Fatal("mismatch control changed committed bytes")
			}
		})
		for _, kind := range []string{"object", "world", "entry", "pin", "unknown"} {
			t.Run("absent-"+kind, func(t *testing.T) {
				bad := m
				bad.References = slices.Clone(m.References)
				bad.References[0] = demoManifestRef{Role: "input", Kind: kind, Ref: sum([]byte("absent " + kind)), Index: 999}
				demoM4bRefusal(t, demoValidateManifest(demoM4bFreshReader(t, r), bad, fixture.token), "every declared manifest reference resolves to matching content")
			})
		}
		for _, field := range []string{"schema", "version", "proof", "tool"} {
			t.Run("schema-"+field, func(t *testing.T) {
				bad := m
				switch field {
				case "schema":
					bad.Schema = "other"
				case "version":
					bad.Version = 2
				case "proof":
					bad.ProofVersion = "other"
				case "tool":
					bad.ToolVersion = "other"
				}
				demoM4bRefusal(t, demoValidateManifest(demoM4bFreshReader(t, r), bad, fixture.token), "public manifest schema/version/pins must match")
			})
		}
		t.Run("role", func(t *testing.T) {
			bad := m
			bad.References = append([]demoManifestRef(nil), m.References...)
			bad.References[0].Role = "arbitrary"
			demoM4bRefusal(t, demoValidateManifest(demoM4bFreshReader(t, r), bad, fixture.token), "public manifest role must be a declared reference name")
		})
		t.Run("why-hash", func(t *testing.T) {
			bad := m
			bad.Why = slices.Clone(m.Why)
			bad.Why[0].SHA256 = sum([]byte("wrong why"))
			demoM4bRefusal(t, demoValidateManifest(demoM4bFreshReader(t, r), bad, fixture.token), "why raw JSON hash mismatch")
		})
		for name, value := range map[string]string{"session-token": fixture.token, "absolute-path": "/opt/private/capture", "timestamp": "2026-10-16T12:34:56Z", "key-material": "private-key: hidden"} {
			t.Run(name, func(t *testing.T) {
				bad := m
				bad.Schema = value
				demoM4bRefusal(t, demoValidateManifest(demoM4bFreshReader(t, r), bad, fixture.token), "public bytes contain")
			})
		}
	})
	t.Run("real-stack", func(t *testing.T) {
		r, m := demoM4bRealSeed(t)
		demoM4bAssertManifestRoles(t, r, m)
		if err := demoValidateManifest(r, m, ""); err != nil {
			t.Fatalf("every declared manifest reference resolves to matching content: %v", err)
		}
		for _, badValue := range []string{sum([]byte("absent real manifest ref"))} {
			bad := m
			bad.References = append([]demoManifestRef(nil), m.References...)
			bad.References[0].Ref = badValue
			demoM4bRefusal(t, demoValidateManifest(r, bad, ""), "every declared manifest reference resolves to matching content")
		}
	})
}

func TestDemoSeedM4bBoundaryControls(t *testing.T) {
	t.Run("s1-decode", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			t.Run(fmt.Sprint(i), func(t *testing.T) {
				objects := demoDomainObjects(nil)
				objects[i].Payload = []byte("{invalid")
				srv := httptest.NewServer(http.NotFoundHandler())
				defer srv.Close()
				err := demoSeedApproval(boundedTestContext(t), newClient(srv.URL), "", objects, demoToolStages{}, hashref.HashRef{})
				demoM4bRefusal(t, err, fmt.Sprintf("S1 payload %d decode refused", i))
			})
		}
	})
	t.Run("test-index", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = io.WriteString(w, `{"header":{},"transitionRef":"present"}`)
		}))
		defer srv.Close()
		demoM4bRefusal(t, demoSeedApproval(boundedTestContext(t), newClient(srv.URL), "", demoDomainObjects(nil), demoToolStages{}, hashref.HashRef{}), "test entry entryIndex must be present and match")
	})
	t.Run("metadata-index", func(t *testing.T) {
		c, _ := demoFakeA2A(t, "completed", map[string]any{"invocation_id": "actual-id"}, map[string]any{"world": map[string]any{"effects": []any{map[string]any{"record": "present"}}}})
		_, _, err := demoA2ASend(boundedTestContext(t), c, "token", "task", "skill", nil)
		demoM4bRefusal(t, err, "entry_index must be present and integral")
	})
	t.Run("s6-fail-closed", func(t *testing.T) {
		fixture := demoM4bManifestFixture(t, 1, "decision packet validated: committed with its invocation record at entry #4;")
		// First establish that S6 reaches its successful end, not an earlier guard.
		if _, err := demoSeedS6(fixture.r, fixture.token, fixture.objects, fixture.stages, fixture.approval, fixture.roster); err != nil {
			t.Fatalf("unit S6 prerequisite: %v", err)
		}
		err := demoSeedApproval(fixture.r.ctx, fixture.r.c, fixture.token, fixture.objects, fixture.stages, hashref.MustParse(fixture.roster[3].Ref))
		demoM4bRefusal(t, err, "S1-S6 checked; roster checks M5c pending; no success manifest written")
	})
}

func TestDemoSeedPacketControls(t *testing.T) {
	t.Run("two-packets", func(t *testing.T) {
		f := demoM4bManifestFixture(t, 2, "decision packet validated: committed with its invocation record at entry #4;")
		_, err := demoSeedS6(f.r, f.token, f.objects, f.stages, f.approval, f.roster)
		demoM4bRefusal(t, err, "entry 4: exactly one decision packet required")
	})
	for _, anchor := range []string{"", "decision packet validated: committed with its invocation record at entry #99;"} {
		name := "missing-anchor"
		if anchor != "" {
			name = "wrong-anchor"
		}
		t.Run(name, func(t *testing.T) {
			f := demoM4bManifestFixture(t, 1, anchor)
			_, err := demoSeedS6(f.r, f.token, f.objects, f.stages, f.approval, f.roster)
			demoM4bRefusal(t, err, "entry 4: packet carrying entry/proposal mismatch")
		})
	}

	t.Run("packet-count", func(t *testing.T) {
		rig, _, _ := demoM4bWhyFixture(t)
		r := worldReader{ctx: boundedTestContext(t), c: newClient(rig.url)}
		run := demoToolRun{EntryIndex: 1}
		_, err := demoSeedS6(r, "", nil, demoToolStages{Read: run, Test: run, Check: run}, run, nil)
		demoM4bRefusal(t, err, "entry 1: exactly one decision packet required")
	})
	for _, name := range []string{"carrying-entry", "proposal"} {
		t.Run(name, func(t *testing.T) {
			carried, proposal := true, "plan"
			if name == "carrying-entry" {
				carried = false
			} else {
				proposal = "other"
			}
			demoM4bRefusal(t, demoValidatePacketAnchor(7, carried, proposal, "plan"), "entry 7: packet carrying entry/proposal mismatch")
		})
	}
	if err := demoValidatePacketAnchor(7, true, "plan", "plan"); err != nil {
		t.Fatal(err)
	}
}

func TestDemoSeedPacketObjectControls(t *testing.T) {
	for _, name := range []string{"packet-object", "packet-decode"} {
		t.Run(name, func(t *testing.T) {
			rig, _, _ := demoM4bWhyFixture(t)
			packet := wireObject{SemanticID: coordinator.DecisionPacketV1, Payload: []byte("invalid canonical packet")}
			packet.Hash = sum(packet.Payload)
			proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				switch req.URL.Path {
				case "/v1/objects/by-semantic-id/" + coordinator.DecisionPacketV1:
					_ = json.NewEncoder(w).Encode(map[string]any{"items": []wireObject{packet}})
					return
				case "/v1/objects/" + packet.Hash:
					if name == "packet-object" {
						http.NotFound(w, req)
					} else {
						_ = json.NewEncoder(w).Encode(packet)
					}
					return
				}
				resp, err := http.Get(rig.url + req.URL.RequestURI())
				if err != nil {
					t.Error(err)
					return
				}
				defer resp.Body.Close()
				w.WriteHeader(resp.StatusCode)
				_, _ = io.Copy(w, resp.Body)
			}))
			defer proxy.Close()
			r := worldReader{ctx: boundedTestContext(t), c: newClient(proxy.URL)}
			run := demoToolRun{EntryIndex: 1}
			_, err := demoSeedS6(r, "", nil, demoToolStages{Read: run, Test: run, Check: run}, run, nil)
			want := "packet: does not resolve"
			if name == "packet-decode" {
				want = "not a canonical world/decision-packet/v1 encoding"
			}
			demoM4bRefusal(t, err, want)
		})
	}
}

// demoM4bFreshReader gives each table-driven control its own bounded
// deadline. Every validation re-reads the whole manifest over REST, so the
// parent's single 30 s budget, shared across ~120 subtests, expires under
// -race on a 2-core runner (PR #289 run, iteration 275).
func demoM4bFreshReader(t *testing.T, r worldReader) worldReader {
	t.Helper()
	r.ctx = boundedTestContext(t)
	return r
}
