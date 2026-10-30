package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

type finishTransport func(*http.Request) (*http.Response, error)

func (f finishTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func memoryDaemon(handler http.HandlerFunc) *client {
	c := newClient("127.0.0.1:47330")
	c.http.Transport = finishTransport(func(r *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		handler(w, r)
		return w.Result(), nil
	})
	return c
}

type finishFixture struct {
	t                    *testing.T
	root, wt, git, dir   string
	sessions, db, before string
	rep                  daemon.EpisodeReport
	health               daemon.HealthResponse
	requests             []string
	checks, outcomes     int
	expired              bool
	failCheck            bool
	seam                 *fakeFinishDispatch
	finishRecords        []daemon.EpisodeFinish
	scopes               []string
	credentialID         string
	recordFailure        bool
	allowRevoked         bool
	recordHook           func()
}

func newFinishFixture(t *testing.T) *finishFixture {
	t.Helper()
	f := &finishFixture{t: t, root: t.TempDir(), rep: ac5Report()}
	f.git, _ = exec.LookPath("git")
	if f.git == "" {
		t.Skip("git unavailable")
	}
	f.dir = filepath.Join(f.root, "presets")
	f.sessions = filepath.Join(f.root, "sessions")
	f.db = filepath.Join(f.root, "world.db")
	f.wt = filepath.Join(f.root, "ws", "ep1")
	for _, dir := range []string{f.dir, f.sessions, filepath.Dir(f.wt)} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	repo := filepath.Join(f.root, "repo")
	f.gitOutput("init", "-q", repo)
	f.gitOutput("-C", repo, "config", "user.name", "Fixture Operator")
	f.gitOutput("-C", repo, "config", "user.email", "fixture@example.invalid")
	f.gitOutput("-C", repo, "config", "commit.gpgsign", "false")
	f.gitOutput("-C", repo, "commit", "--allow-empty", "-qm", "base")
	base := strings.TrimSpace(f.gitOutput("-C", repo, "rev-parse", "HEAD"))
	if err := worktree.Add(context.Background(), f.git, repo, f.wt, "custom-branch", base); err != nil {
		t.Fatal(err)
	}
	f.before = base
	f.rep.Start = &store.EpisodeStart{Episode: "ep1", Agent: "codex", Model: "fixture-model", Preset: store.StartPreset{Name: "fixture", Bytes: `{"profile":""}`},
		Branch: "custom-branch", Base: store.StartBase{SHA: base, Branch: "dev"}, Remote: store.StartRemote{Name: "origin", URL: "https://github.com/fixture/repo.git"}}
	f.gitOutput("-C", f.wt, "remote", "add", "origin", f.rep.Start.Remote.URL)
	credentials := filepath.Join(f.root, "credentials.db")
	f.health = daemon.HealthResponse{DBPath: f.db, CredentialsDB: credentials}
	st, err := store.OpenCredentials(credentials)
	if err != nil {
		t.Fatal(err)
	}
	_, f.credentialID, _, err = authority.Mint(boundedTestContext(t), st, "ep1", []broker.Capability{{Effect: "Workspace.Write", Scope: "world", Budget: 10, ExpiresAt: time.Now().Unix() + 3600}}, 3600, time.Now().Unix(), nil)
	st.Close()
	if err != nil {
		t.Fatal(err)
	}
	f.seam = &fakeFinishDispatch{label: true}
	oldFactory, oldHook := finishDispatchFactory, finishBeforeDispatch
	finishDispatchFactory = func() finishDispatch { return f.seam }
	t.Cleanup(func() { finishDispatchFactory, finishBeforeDispatch = oldFactory, oldHook })
	f.write("a.go", "agent\n")
	f.write("gen.txt", "produced\n")
	f.write("b.go", "hand\n")
	f.writePreset()
	os.WriteFile(filepath.Join(f.sessions, "ep1.session"), []byte(strings.Repeat("1", 64)), 0o600)
	oldClient, oldHealth := finishClient, fetchDaemonHealth
	t.Cleanup(func() { finishClient, fetchDaemonHealth = oldClient, oldHealth })
	finishClient = func(string) *client { return memoryDaemon(f.serve) }
	fetchDaemonHealth = func(string) (daemon.HealthResponse, error) {
		f.requests = append(f.requests, "health")
		return f.health, nil
	}
	return f
}

func (f *finishFixture) writePreset() {
	f.t.Helper()
	data, _ := json.Marshal(map[string]any{"schema": "world/preset/v1", "name": "fixture",
		"daemon": map[string]any{"addr": "127.0.0.1:47330", "db": f.db, "workspace_root": filepath.Dir(f.wt), "sessions_dir": f.sessions}})
	if err := os.WriteFile(filepath.Join(f.dir, "fixture.json"), data, 0o600); err != nil {
		f.t.Fatal(err)
	}
}

func (f *finishFixture) write(path, text string) {
	f.t.Helper()
	if err := os.WriteFile(filepath.Join(f.wt, path), []byte(text), 0o600); err != nil {
		f.t.Fatal(err)
	}
}

func (f *finishFixture) gitOutput(args ...string) string {
	f.t.Helper()
	out, err := worktree.Output(context.Background(), f.git, args)
	if err != nil {
		f.t.Fatal(err)
	}
	return out
}

func (f *finishFixture) serve(w http.ResponseWriter, r *http.Request) {
	f.requests = append(f.requests, r.URL.Path)
	switch r.URL.Path {
	case "/v1/episodes/ep1/report":
		if len(f.finishRecords) > 0 {
			latest := f.finishRecords[len(f.finishRecords)-1]
			f.rep.Finish = &latest
		}
		json.NewEncoder(w).Encode(f.rep)
	case "/v1/session/self":
		if f.expired {
			w.WriteHeader(401)
			return
		}
		if r.Header.Get("Authorization") != "Bearer "+strings.Repeat("1", 64) {
			f.t.Error("session credential missing")
		}
		json.NewEncoder(w).Encode(daemon.SessionSelf{Episode: "ep1", SecondsLeft: 100})
	case "/v1/episodes/ep1/finish":
		if r.Header.Get("Authorization") != "" || r.Header.Get("X-World-Operator-Proof") == "" {
			f.t.Error("operator proof required, no session credential")
		}
		var input daemon.FinishRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			f.t.Fatal(err)
		}
		path := filepath.Join(filepath.Dir(f.health.CredentialsDB), "operator-proofs", "ep1.json")
		data, err := os.ReadFile(path)
		if err != nil || !strings.Contains(string(data), finishSHA([]byte(r.Header.Get("X-World-Operator-Proof")))) {
			f.t.Error("proof file does not bind request")
		}
		os.Remove(path)
		if input.Phase == "approve" {
			f.scopes = append(f.scopes, input.Scope)
			json.NewEncoder(w).Encode(daemon.FinishApproval{Request: "sha256:" + strings.Repeat("2", 64), Approval: "sha256:" + strings.Repeat("3", 64)})
		} else {
			creds, err := store.OpenCredentialsReadOnly(f.health.CredentialsDB)
			if err != nil {
				f.t.Fatal(err)
			}
			_, exists, err := creds.ResolveSession(r.Context(), f.credentialID)
			creds.Close()
			if err != nil || (!exists && !f.allowRevoked) {
				f.t.Error("credential revoked before record")
			}
			if f.recordFailure {
				w.WriteHeader(500)
				return
			}
			f.finishRecords = append(f.finishRecords, *input.Record)
			if f.recordHook != nil {
				f.recordHook()
			}
			w.Write([]byte("{}"))
		}
	case "/v1/commit":
		if r.Header.Get("Authorization") != "Bearer "+strings.Repeat("1", 64) {
			f.t.Error("outcome credential missing")
		}
		var obj demoObject
		if err := json.NewDecoder(r.Body).Decode(&obj); err != nil || obj.SemanticID != taskOutcomeSemanticID {
			f.t.Errorf("outcome object: %+v %v", obj, err)
		}
		var payload taskOutcome
		json.Unmarshal(obj.Payload, &payload)
		if len(payload.Outcome.StderrHead) > 400 || payload.Arm != "world" {
			f.t.Errorf("bad outcome: %+v", payload)
		}
		f.outcomes++
		json.NewEncoder(w).Encode(map[string]int{"entry": 6000 + f.outcomes})
	case "/mcp/":
		f.checks++
		var rpc struct {
			Params struct {
				Name      string `json:"name"`
				Arguments struct {
					Command string   `json:"command"`
					Args    []string `json:"args"`
				} `json:"arguments"`
			} `json:"params"`
		}
		json.NewDecoder(r.Body).Decode(&rpc)
		if rpc.Params.Name != "workspace-exec" || rpc.Params.Arguments.Command != "pinned" || !reflect.DeepEqual(rpc.Params.Arguments.Args, []string{"test"}) {
			f.t.Errorf("unexpected pinned call: %+v", rpc)
		}
		rc := 0
		if f.failCheck && f.checks == 2 {
			rc = 7
		}
		payload, _ := json.Marshal(map[string]any{"exit_code": rc, "stderr": "fixture stderr"})
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprintf(w, "data: {\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"content\":[{\"type\":\"text\",\"text\":%q}]}}\n\n", string(payload))
	default:
		f.t.Errorf("unexpected request %s", r.URL.Path)
		w.WriteHeader(404)
	}
}

func (f *finishFixture) run(env sessionEnv, extra ...string) (int, string, string) {
	var out, errw bytes.Buffer
	record := func(ctx context.Context, c *client, token string, obj demoObject, writtenBy, interpreter, op string) (int64, error) {
		if writtenBy != taskOutcomeWrittenBy {
			f.t.Error("wrong outcome writer")
		}
		data, _ := json.Marshal(obj)
		status, body, err := c.doAuth(ctx, http.MethodPost, "/v1/commit", bytes.NewReader(data), "Bearer "+token)
		if err != nil || status != 200 {
			return 0, fmt.Errorf("fixture outcome: %d %v", status, err)
		}
		var result struct {
			Entry int64
		}
		err = json.Unmarshal(body, &result)
		return result.Entry, err
	}
	args := append([]string{"ep1", "--preset", "fixture", "--presets-dir", f.dir}, extra...)
	rc := runEpisodeFinishWith(args, &out, &errw, env, record)
	return rc, out.String(), errw.String()
}

func TestEpisodeFinishFenceAC1(t *testing.T) {
	f := newFinishFixture(t)
	before := f.gitOutput("-C", f.wt, "status", "--porcelain")
	rc, _, errw := f.run(noTerminalEnv(1))
	if rc != 1 || !strings.Contains(errw, "refusing: no controlling terminal") || len(f.requests) != 0 {
		t.Fatalf("fence: %d %s requests %v", rc, errw, f.requests)
	}
	if got := strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")); got != f.before || f.gitOutput("-C", f.wt, "status", "--porcelain") != before {
		t.Fatal("fence changed git state")
	}
	var out, errs bytes.Buffer
	if rc := runEpisodeFinish([]string{"ep1", "--preset", "missing"}, &out, &errs, noTerminalEnv(1)); rc != 1 || !strings.Contains(errs.String(), "no controlling terminal") {
		t.Fatal("bad preset preceded fence")
	}
}

func TestEpisodeFinishValidationAC2(t *testing.T) {
	for _, name := range []string{"degraded", "expired", "inflight", "unchanged", "no-start", "preset-mismatch", "argv-bound", "unknown-allow"} {
		t.Run(name, func(t *testing.T) {
			f := newFinishFixture(t)
			want := ""
			extra := []string{}
			switch name {
			case "degraded":
				f.health.DegradedEpisodes = []string{"ep1"}
				want = "is degraded"
			case "expired":
				f.expired = true
				want = "session has expired"
			case "inflight":
				f.rep.Uncommitted = 1
				want = "still in flight"
			case "unchanged":
				for _, p := range []string{"a.go", "b.go", "gen.txt"} {
					os.Remove(filepath.Join(f.wt, p))
				}
				want = "nothing to finish"
			case "no-start":
				f.rep.Start = nil
				want = "has no start record"
			case "preset-mismatch":
				f.rep.Start.Preset.Name = "other"
				want = "preset fixture does not match start preset other"
			case "argv-bound":
				argv := make([]string, 33)
				f.rep.Start.Checks = []store.StartCheck{{Name: "long", Argv: argv}}
				want = "argv needs 1..32"
			case "unknown-allow":
				extra = []string{"--allow-failing", "missing"}
				want = "missing is not a check"
			}
			rc, _, errw := f.run(testMintEnv("y\n", 1), extra...)
			if rc != 1 || !strings.Contains(errw, want) || f.checks != 0 || f.outcomes != 0 {
				t.Fatalf("%d %s checks %d outcomes %d", rc, errw, f.checks, f.outcomes)
			}
		})
	}
}

func TestEpisodeFinishChecksAC3(t *testing.T) {
	for _, allow := range []bool{false, true} {
		t.Run(fmt.Sprint(allow), func(t *testing.T) {
			f := newFinishFixture(t)
			f.failCheck = true
			for _, name := range []string{"pass", "fail"} {
				f.rep.Start.Checks = append(f.rep.Start.Checks, store.StartCheck{Name: name, Argv: []string{"go", "test"}, Cwd: ".", TimeoutS: 1, Exec: "pinned", Args: []string{"test"}})
			}
			extra := []string{"--include-unattributed"}
			if allow {
				extra = append(extra, "--allow-failing", "fail")
			}
			rc, _, errw := f.run(testMintEnv("y\n", 1), extra...)
			if f.checks != 2 || f.outcomes != 2 {
				t.Fatalf("checks %d outcomes %d: %s", f.checks, f.outcomes, errw)
			}
			if !allow {
				if rc != 1 || !strings.Contains(errw, "check fail failed (rc 7), outcome 6002") || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != f.before {
					t.Fatalf("%d %s", rc, errw)
				}
			} else {
				body, _ := os.ReadFile(filepath.Join(f.sessions, "ep1.pr-body.md"))
				if rc != 0 || !strings.Contains(string(body), "fail: FAILED (sandboxed), outcome entry 6002") {
					t.Fatalf("%d %s\n%s", rc, errw, body)
				}
			}
		})
	}
}

func TestEpisodeFinishUnsandboxedAC4(t *testing.T) {
	for _, answer := range []string{"n\n", "y\n"} {
		t.Run(strings.TrimSpace(answer), func(t *testing.T) {
			f := newFinishFixture(t)
			// Names alone are safe to set; the child must never inherit the token.
			t.Setenv("WORLD_SESSION", "fixture-sentinel")
			f.rep.Start.Checks = []store.StartCheck{{Name: "host", Argv: []string{"/bin/sh", "-c", "test -z \"$WORLD_SESSION\" && test -z \"$AILANG_REGISTRY_API_KEY\" && echo ran > marker"}, Cwd: ".", TimeoutS: 1, Unsandboxed: true}}
			rc, _, errw := f.run(testMintEnv(answer, 1), "--include-unattributed", "--dry-run")
			_, err := os.Stat(filepath.Join(f.wt, "marker"))
			if answer == "n\n" {
				if rc != 1 || !os.IsNotExist(err) || f.outcomes != 0 {
					t.Fatalf("declined check ran: %d %v %s", rc, err, errw)
				}
			} else if rc != 0 || err != nil || f.outcomes != 1 {
				t.Fatalf("host check: %d %v %s", rc, err, errw)
			}
		})
	}
}

func TestEpisodeFinishAttributionCommitDryRunAC5AC6(t *testing.T) {
	f := newFinishFixture(t)
	paths := []string{"a.go", "b.go", "gen.txt"}
	if got := finishUnattributed(paths, f.rep); !reflect.DeepEqual(got, []string{"b.go"}) {
		t.Fatal(got)
	}
	rc, _, errw := f.run(testMintEnv("", 1))
	if rc != 1 || !strings.Contains(errw, "unattributed paths (F5): b.go") {
		t.Fatalf("%d %s", rc, errw)
	}
	before := map[string]string{}
	for _, p := range paths {
		b, _ := os.ReadFile(filepath.Join(f.wt, p))
		before[p] = string(b)
	}
	rc, out, errw := f.run(testMintEnv("", 1), "--dry-run", "--include-unattributed")
	if rc != 0 || !strings.Contains(out, "### Not made through World\n\n- b.go") || !strings.Contains(out, "3 files changed") {
		t.Fatalf("dry-run %d %s\n%s", rc, errw, out)
	}
	t.Log("AC5 --dry-run --include-unattributed:\n" + out)
	if strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != f.before {
		t.Fatal("dry-run changed HEAD")
	}
	for p, want := range before {
		b, _ := os.ReadFile(filepath.Join(f.wt, p))
		if string(b) != want {
			t.Fatal("dry-run changed bytes of " + p)
		}
	}
	rc, out, errw = f.run(testMintEnv("y\n", 1), "--include-unattributed", "--title", "Fixture finish")
	if rc != 0 {
		t.Fatalf("%d %s", rc, errw)
	}
	message := f.gitOutput("-C", f.wt, "log", "-1", "--format=%B")
	t.Log("AC6 commit message:\n" + message)
	if !strings.HasSuffix(strings.TrimSpace(message), reportTrailer(f.rep, f.db)) || !strings.Contains(message, "entries=1-6000") {
		t.Fatal(message)
	}
	identity := f.gitOutput("-C", f.wt, "log", "-1", "--format=%an <%ae>|%cn <%ce>")
	if strings.TrimSpace(identity) != "codex via AILANG World <ailang-world@invalid>|Fixture Operator <fixture@example.invalid>" {
		t.Fatal(identity)
	}
	got := f.gitOutput("-C", f.wt, "show", "--format=", "--name-only", "HEAD")
	if strings.TrimSpace(got) != "a.go\nb.go\ngen.txt" {
		t.Fatal(got)
	}
	bodyPath := filepath.Join(f.sessions, "ep1.pr-body.md")
	info, _ := os.Stat(bodyPath)
	if info == nil || info.Mode().Perm() != 0o600 || !strings.Contains(out, "https://github.com/fixture/repo/pull/1") {
		t.Fatal(out)
	}
	head := strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD"))
	// M3: a re-run after `done` is the finished state, not M1's "already committed" refusal.
	records, finishPosts := len(f.finishRecords), countRequests(f.requests, "/v1/episodes/ep1/finish")
	rc, out, errw = f.run(testMintEnv("", 1), "--include-unattributed")
	if rc != 0 || !strings.Contains(out, "episode ep1 is already finished: https://github.com/fixture/repo/pull/1") || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != head {
		t.Fatalf("rerun %d %s %s", rc, out, errw)
	}
	if len(f.finishRecords) != records || countRequests(f.requests, "/v1/episodes/ep1/finish") != finishPosts {
		t.Fatalf("rerun after done posted again: records %d -> %d", records, len(f.finishRecords))
	}
}

func TestEpisodeFinishHelp(t *testing.T) {
	var out, errw bytes.Buffer
	if rc := runEpisodeFinish([]string{"--help"}, &out, &errw, noTerminalEnv(1)); rc != 0 {
		t.Fatalf("help rc %d: %s", rc, errw.String())
	}
	for _, flag := range []string{"--preset", "--presets-dir", "--title", "--allow-failing", "--include-unattributed", "--dry-run", "--ready", "--exec-profiles-dir"} {
		if !strings.Contains(out.String(), flag) {
			t.Error("missing " + flag)
		}
	}
	for _, flag := range []string{"--yes", "--renew", "--check"} {
		if strings.Contains(out.String(), flag) {
			t.Error("out-of-scope flag " + flag)
		}
	}
}

func countRequests(paths []string, want string) int {
	n := 0
	for _, p := range paths {
		if p == want {
			n++
		}
	}
	return n
}
