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
