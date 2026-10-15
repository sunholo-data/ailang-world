package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
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
