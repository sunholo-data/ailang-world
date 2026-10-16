package main

import (
	"bytes"
	"context"
	"encoding/json"
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
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/evidence"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
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
	for _, s := range []string{"--db", "--session", "--ailang-bin", "--proof-key-file", "--addr", "--manifest", "v0.41.0", "v0.52.1", "S2-S6"} {
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
}

func newDemoSeRig(t *testing.T, interp, tool string) demoSeRig {
	return newDemoSeRigForSeed(t, interp, tool, false)
}

func newDemoSeRigForSeed(t *testing.T, interp, tool string, empty bool) demoSeRig {
	t.Helper()
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
	cfg := daemon.Config{DBPath: db, BindHost: daemon.DefaultBindHost, AilangBin: interp, WorkspaceRoot: ws,
		ToolAilangBin: tool, ExamplesDir: corpus, ErrorLog: io.Discard}
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
	now := time.Now().Unix()
	grants := []broker.Capability{
		{Effect: "world.apply", Scope: "world", ExpiresAt: now + 7200, Budget: 100},
		{Effect: "Human.Approve", Scope: "worktree", ExpiresAt: now + 7200, Budget: 20},
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
	return demoSeRig{url: srv.URL, db: db, token: tok, ws: ws}
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
	// v0.52.1 at startup; the test records that archived release for this run.
	toolBytes, err := os.ReadFile(tool)
	if err != nil {
		t.Fatal(err)
	}
	toolRef := hashref.SumSHA256(toolBytes)
	if m, err := archive.New(r.db).ReadManifest(toolRef); err != nil || !strings.Contains(m.Version, "v0.52.1") {
		t.Fatalf("retained tool pin %s is not v0.52.1: %v %q", toolRef, err, m.Version)
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
	fakeA2A := func(t *testing.T, state string, meta, data map[string]any) (*client, *int) {
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
	goodData := map[string]any{"content": "x", "world": map[string]any{"effects": []any{map[string]any{"record": "sha256:" + strings.Repeat("a", 64)}}}}
	goodMeta := map[string]any{"invocation_id": "a2a:demo:t"}
	ctx := boundedTestContext(t)
	t.Run("a2a-positive", func(t *testing.T) {
		c, _ := fakeA2A(t, "completed", goodMeta, goodData)
		if _, _, err := demoA2ASend(ctx, c, "tok", "t", "ailang-read", map[string]any{}); err != nil {
			t.Fatalf("complete task refused: %v", err)
		}
	})
	t.Run("a2a-non-completed", func(t *testing.T) {
		c, _ := fakeA2A(t, "working", goodMeta, goodData)
		_, _, err := demoA2ASend(ctx, c, "tok", "t", "ailang-read", map[string]any{})
		if err == nil || !strings.Contains(err.Error(), "tasks/send refused") {
			t.Fatalf("non-completed state accepted: %v", err)
		}
	})
	t.Run("a2a-missing-invocation", func(t *testing.T) {
		c, _ := fakeA2A(t, "completed", map[string]any{}, goodData)
		_, _, err := demoA2ASend(ctx, c, "tok", "t", "ailang-read", map[string]any{})
		if err == nil || !strings.Contains(err.Error(), "invocation or effect record reference absent") {
			t.Fatalf("missing invocation id accepted: %v", err)
		}
	})
	t.Run("a2a-missing-effect-record", func(t *testing.T) {
		c, _ := fakeA2A(t, "completed", goodMeta, map[string]any{"content": "x"})
		_, _, err := demoA2ASend(ctx, c, "tok", "t", "ailang-read", map[string]any{})
		if err == nil || !strings.Contains(err.Error(), "invocation or effect record reference absent") {
			t.Fatalf("missing effect record accepted: %v", err)
		}
	})
	t.Run("op-allowlist-before-dispatch", func(t *testing.T) {
		c, calls := fakeA2A(t, "completed", goodMeta, goodData)
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
		// CI does not provision the v0.52.1 tool pin, so this skips there and
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
			var before wireWorld
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
				if writes != 1 || len(envelopes) != 1 || !strings.Contains(errOut, "S1-S4 recorded; S5-S6 incomplete; no success manifest written") {
					t.Fatalf("seed must record exactly one S4 envelope and refuse S5/S6 completion: writes=%d envelopes=%d err=%q", writes, len(envelopes), errOut)
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
				if writes != 1 || len(envelopes) != wantEnvelopes || !strings.Contains(errOut, "S4 refused:") || !strings.Contains(errOut, "inspect history") || !strings.Contains(errOut, "no success manifest written") {
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
