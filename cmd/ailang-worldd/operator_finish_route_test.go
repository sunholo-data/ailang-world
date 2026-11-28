package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

func TestOperatorFinishRouteWorker(t *testing.T) {
	f := newOperatorLaunchRigWithRemote(t, fakeLaunchMarker+"exit 0\n", "private-task-fixture-marker", fixtureFinishRemote)
	if out := f.tick(t); !strings.Contains(out, ": stopped") {
		t.Fatalf("the fake child never completed: %s", out)
	}
	ctx := activationTestContext(t)
	start, err := f.db.EpisodeStart(ctx, f.action.Episode)
	if err != nil || start == nil || len(start.Checks) != 1 || start.Checks[0].Unsandboxed {
		t.Fatalf("completed start missing pinned checks: %v", err)
	}
	const filename = "private-path-fixture-marker"
	change := filepath.Join(f.path(), filename)
	seedStartFixtureFile(t, change, []byte("approved change\n"))
	repo := f.f.r.f
	bare := filepath.Join(t.TempDir(), "remote.git")
	presetGit(t, repo, "clone", "-q", "--local", "--bare", "--", repo.repo, bare)
	remoteHead := presetGit(t, repo, "--git-dir="+bare, "rev-parse", "HEAD")
	trace := filepath.Join(f.home, "unexpected-remote-write")
	tools := f.c.ToolDirs[0]
	f.c.FinishRemotes = []string{fixtureFinishRemote}
	writeFinishTestTool(t, &f.c, tools, "gh", "case \"$1 $2\" in 'auth status'|'label list') ;; 'repo view') printf 'WRITE\\n';; *) touch "+shellQuote(trace)+"; exit 96;; esac\n")
	writeFinishTestTool(t, &f.c, tools, "git", "for arg; do case \"$arg\" in push|fetch|clone|pull|send-pack|fetch-pack|receive-pack|upload-pack|remote|submodule) touch "+shellQuote(trace)+"; exit 96;; esac; done\n"+
		"case \" $* \" in *' ls-remote '*) exec "+shellQuote(repo.git)+" --git-dir=/dev/null ls-remote --symref -- "+shellQuote(bare)+" HEAD;; esac\n"+
		"exec "+shellQuote(repo.git)+" \"$@\"\n")
	paths := f.env.runnerPaths()
	config, err := json.Marshal(f.c)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.env.atomic(paths.Record, config, 0600); err != nil {
		t.Fatal(err)
	}
	if err := f.env.atomic(paths.Plist, []byte(runnerPlist(paths, f.c, config)), 0600); err != nil {
		t.Fatal(err)
	}

	const origin = "http://localhost:7644"
	log := &syncBuffer{}
	d, err := daemon.New(ctx, daemon.Config{DBPath: f.c.World, CredentialsPath: f.creds, BindHost: daemon.DefaultBindHost,
		SurfaceOrigins: []string{origin}, ErrorLog: log, ApprovalClock: f.src.clock})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = d.Close() })
	srv := httptest.NewServer(d.Handler())
	t.Cleanup(srv.Close)
	request := func(method, path string, input any, out any) string {
		t.Helper()
		var data []byte
		if input != nil {
			data, err = json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}
		}
		requestCtx, cancel := f.env.operatorContext(activationTestBudget)
		defer cancel()
		req, err := http.NewRequestWithContext(requestCtx, method, srv.URL+path, bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Origin", origin)
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("Content-Type", "application/json")
		res, err := srv.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		raw, err := io.ReadAll(res.Body)
		if err != nil || res.StatusCode != 200 || out != nil && json.Unmarshal(raw, out) != nil {
			t.Fatalf("%s %s: %d %s (%v)", method, path, res.StatusCode, raw, err)
		}
		return string(raw)
	}

	rep := daemon.EpisodeReport{Episode: start.Episode, Start: start, Entries: [2]int64{1, 1},
		Writes: daemon.ReportWrites{Written: []daemon.ReportWrite{{Path: filename, Op: "write", Count: 1}}}}
	h := hashref.SumSHA256([]byte("selected world")).String()
	checks, outcomes := 0, 0
	activeID, statusHref := "", ""
	reportDown := false
	original := finishClient
	t.Cleanup(func() { finishClient = original })
	finishClient = func(string) *client {
		return memoryDaemon(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/v1/episodes/" + start.Episode + "/report":
				if reportDown {
					w.WriteHeader(http.StatusServiceUnavailable)
					return
				}
				_ = json.NewEncoder(w).Encode(rep)
			case "/mcp/":
				if body := request("GET", statusHref, nil, nil); !strings.Contains(body, "Finish outcome uncertain") {
					t.Errorf("reserved check was missing from action status: %s", body)
				}
				terminal, lockErr := openTerminalEpisodeLock(f.creds, start.Episode, nil, os.O_CREATE, 5*time.Second)
				if terminal != nil {
					_ = terminal.Close()
				}
				if !errors.Is(lockErr, store.ErrOperatorConflict) || !strings.Contains(lockErr.Error(), "active action "+activeID) {
					t.Errorf("terminal did not see the active private finish: %v", lockErr)
				}
				var call struct {
					Params struct {
						Name      string
						Arguments struct {
							Command string
							Args    []string
						}
					}
				}
				if r.Header.Get("Authorization") != "Bearer "+f.token || json.NewDecoder(r.Body).Decode(&call) != nil ||
					call.Params.Name != "workspace-exec" || call.Params.Arguments.Command != start.Checks[0].Exec ||
					!reflect.DeepEqual(call.Params.Arguments.Args, start.Checks[0].Args) {
					t.Errorf("check used different argv or session")
				}
				checks++
				w.Header().Set("Content-Type", "text/event-stream")
				fmt.Fprintf(w, "data: {\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"content\":[{\"type\":\"text\",\"text\":%q}]}}\n\n", `{"exit_code":0}`)
			case "/v1/head":
				fmt.Fprint(w, h)
			case "/v1/worlds/" + h:
				_ = json.NewEncoder(w).Encode(wireWorld{Ref: h, Revision: 1, LogHead: h})
			case "/v1/log/1":
				_ = json.NewEncoder(w).Encode(wireEntry{Header: wireHeader{EntryIndex: 1, SemanticsEpoch: 1, Interpreter: h}})
			case "/v1/commit":
				var commit demoCommit
				if r.Header.Get("Authorization") != "Bearer "+f.token || json.NewDecoder(r.Body).Decode(&commit) != nil ||
					len(commit.Objects) != 1 || commit.Objects[0].SemanticID != taskOutcomeSemanticID {
					t.Errorf("check outcome writer crossed the pinned evidence seam")
				}
				outcomes++
				_ = json.NewEncoder(w).Encode(map[string]string{"selectedHead": commit.NextWorld.Ref})
			default:
				t.Errorf("unexpected remote, V3 or finish-record operation: %s %s", r.Method, r.URL.Path)
				w.WriteHeader(http.StatusBadRequest)
			}
		})
	}

	source := f.src
	source.workspace = "" // The real runner starts with the workspace supplied by its installed record.
	raw := f.env.readOperatorCatalogue(f.c, source)
	f.h.Build, f.h.Readiness, f.h.Catalogue, f.h.Owner.Generation = daemon.CurrentBuild().ExecutableSHA256, "both", raw, 0
	entry, err := f.h.FinishReview(start.Episode, start.Preset.Name)
	if err != nil || entry.Commit.Snapshot.Head != start.Base.SHA {
		t.Fatalf("completed launch was not captured from the installed workspace: %v", err)
	}
	if err := f.env.touchRunner(f.c, &f.h); err != nil {
		t.Fatal(err)
	}
	var review struct {
		Outcome             string
		Argv                []string
		CatalogueGeneration int64
	}
	command := "episode finish " + start.Episode + " --preset " + start.Preset.Name + " --confirm-on-surface"
	request("GET", "/surface/verb.json?q="+url.QueryEscape(command), nil, &review)
	if review.Outcome != "review" {
		t.Fatalf("missing review: %+v", review)
	}
	var status struct {
		ID, Href, StatusHref, State, LocalCommit string
		Checks                                   []struct{ Name, Result string }
	}
	request("POST", "/surface/operator-actions", map[string]any{
		"argv": review.Argv, "operationId": strings.Repeat("4", 32), "catalogueGeneration": review.CatalogueGeneration, "task": "",
	}, &status)
	row, err := f.db.OperatorAction(ctx, status.ID, f.now)
	if err != nil {
		t.Fatal(err)
	}
	card := request("GET", status.Href, nil, nil)
	if !strings.Contains(card, filename) || !strings.Contains(card, "Prepare local commit") || strings.Contains(card, `id="operator-task"`) {
		t.Fatal("signed review omitted its pinned manifest or rendered a task")
	}
	key := approvaltest.NewSigner(54)
	seedApprover(t, f.creds, "finish-operator-fixture", key.PublicKey(), f.now-1)
	message, err := approval.OperatorActionMessage([32]byte(row.Digest), [32]byte(row.Challenge), row.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	artifact := key.Sign(approval.OperatorActionNamespace, message)
	request("POST", status.Href+"/confirm", map[string]any{"v": 1, "digest": hex.EncodeToString(row.Digest), "sshsig": string(artifact)}, &status)

	activeID, statusHref = row.ID, status.StatusHref
	run := func() error {
		turn, stop := f.env.operatorContext(activationTestBudget)
		defer stop()
		return f.env.consumeFinish(turn, f.c, f.h, f.src)
	}
	before := serviceFileHash(f.creds)
	held, err := openTerminalEpisodeLock(f.creds, start.Episode, nil, os.O_CREATE, 5*time.Second)
	if err != nil || held == nil {
		t.Fatalf("terminal could not acquire the episode turn: %v", err)
	}
	blocked := run()
	if err := held.Close(); err != nil {
		t.Fatal(err)
	}
	if !errors.Is(blocked, store.ErrOperatorConflict) || !strings.Contains(blocked.Error(), row.ID) || before != serviceFileHash(f.creds) {
		t.Fatalf("lock contention crossed admission: %v", blocked)
	}
	var runnerLog syncBuffer
	refuse := func(name string, want error) {
		t.Helper()
		before := serviceFileHash(f.creds)
		err := run()
		_, remoteWrite := os.Lstat(trace)
		if !errors.Is(err, want) || before != serviceFileHash(f.creds) || checks+outcomes != 0 ||
			!errors.Is(remoteWrite, os.ErrNotExist) || finishRefusal(err) != want.Error() {
			t.Fatalf("%s changed evidence or lost the remedy before claim: %v", name, err)
		}
		fmt.Fprintln(&runnerLog, finishRefusal(err))
	}
	launch := f.receipt(t).Receipt
	savedLaunch, _ := json.Marshal(launch)
	launch.PID = os.Getpid()
	birth, live := operatorProcessIdentity(launch.PID)
	if !live {
		t.Fatal("live PID positive control unavailable")
	}
	launch.Birth = birth
	liveLaunch, _ := json.Marshal(launch)
	f.exec(t, "UPDATE operator_action_steps SET receipt=? WHERE kind='launch'", liveLaunch)
	refuse("active agent", errOperatorFinishEligibility)
	f.exec(t, "UPDATE operator_action_steps SET receipt=? WHERE kind='launch'", savedLaunch)
	startWire, _ := json.Marshal(start)
	f.exec(t, "DELETE FROM episode_starts WHERE episode=?", start.Episode)
	refuse("missing start", errOperatorFinishEligibility)
	f.exec(t, "INSERT INTO episode_starts(episode,start_json) VALUES (?,?)", start.Episode, string(startWire))
	originalHealth := f.health
	f.health.DegradedEpisodes = []string{start.Episode}
	refuse("degraded health", errOperatorFinishEligibility)
	f.health = originalHealth
	reportDown = true
	refuse("report unavailable", errOperatorFinishEligibility)
	reportDown = false
	if err := os.Remove(change); err != nil {
		t.Fatal(err)
	}
	refuse("empty worktree", errOperatorFinishReview)
	seedStartFixtureFile(t, change, []byte("different approved input\n"))
	refuse("untracked contents", errOperatorFinishReview)
	seedStartFixtureFile(t, change, []byte("approved change\n"))

	branch := "refs/heads/" + start.Branch
	baseBranch := "refs/heads/" + start.Base.Branch
	tree := presetGit(t, repo, "rev-parse", start.Base.SHA+"^{tree}")
	foreign := presetGit(t, repo, "-c", "user.name=Fixture", "-c", "user.email=fixture@invalid",
		"commit-tree", tree, "-p", start.Base.SHA, "-m", "new fixture HEAD")
	for _, ref := range []string{branch, baseBranch} {
		presetGit(t, repo, "update-ref", ref, foreign, start.Base.SHA)
		refuse(ref, errOperatorFinishReview)
		presetGit(t, repo, "update-ref", ref, start.Base.SHA, foreign)
	}
	presetGit(t, repo, "remote", "set-url", start.Remote.Name, "https://github.com/fixture/changed.git")
	refuse("remote", errOperatorFinishReview)
	presetGit(t, repo, "remote", "set-url", start.Remote.Name, start.Remote.URL)

	configPath := filepath.Join(repo.repo, ".git", "config")
	approvedConfig := readStartFixtureFile(t, configPath)
	seedStartFixtureFile(t, configPath, append(approvedConfig, []byte("\n[credential]\n\thelper = !touch "+shellQuote(trace)+"\n")...))
	refuse("repository helper", errOperatorFinishReview)
	seedStartFixtureFile(t, configPath, approvedConfig)

	_, profilePath, err := loadMintProfile(f.src.profiles, "proj")
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{f.f.file, profilePath} {
		pinned := readStartFixtureFile(t, path)
		seedStartFixtureFile(t, path, append(bytes.Clone(pinned), '\n'))
		refuse(filepath.Base(path), errOperatorFinishReview)
		seedStartFixtureFile(t, path, pinned)
	}
	token := filepath.Join(f.f.r.sessions, start.Episode+".session")
	saved := readStartFixtureFile(t, token)
	seedStartFixtureFile(t, token, []byte(strings.Repeat("0", len(f.token))))
	refuse("changed session", errOperatorFinishEligibility)
	seedStartFixtureFile(t, token, saved)
	f.exec(t, "UPDATE session_credentials SET expires_at=? WHERE episode_id=?", f.now, start.Episode)
	refuse("expired session", errOperatorFinishEligibility)
	f.exec(t, "UPDATE session_credentials SET expires_at=? WHERE episode_id=?", f.self.ExpiresAt, start.Episode)

	f.env.operatorCatalogue = func(runnerConfig) []byte { return raw }
	fmt.Fprint(&runnerLog, f.tick(t))
	row, err = f.db.OperatorAction(ctx, status.ID, f.now)
	if err != nil {
		t.Fatal(err)
	}
	j, err := f.db.ReadOperatorFinishJournal(ctx, operatorPreparationCAS(row))
	head, headErr := worktree.Output(ctx, repo.git, []string{"-C", f.path(), "rev-parse", "HEAD"})
	if err != nil || headErr != nil || j.Commit.SHA != strings.TrimSpace(head) || j.Commit.State != store.OperatorLocalCommitDone ||
		row.V3ID != "" || checks != 1 || outcomes != 1 || len(j.Checks) != 1 || j.Checks[0].Result != "passed" {
		t.Fatalf("local hand-off: %+v %v head %v, checks %d, outcomes %d", j, err, headErr, checks, outcomes)
	}
	if err := run(); !errors.Is(err, store.ErrOperatorMissing) || checks != 1 || outcomes != 1 {
		t.Fatalf("durable local hand-off replayed: %v", err)
	}
	request("GET", status.StatusHref, nil, &status)
	if row.Owner != f.h.Owner.Nonce || row.Generation == 0 || status.LocalCommit != j.Commit.SHA || len(status.Checks) != 1 ||
		status.Checks[0].Name != start.Checks[0].Name || status.Checks[0].Result != "passed" ||
		presetGit(t, repo, "-C", f.path(), "rev-list", "--count", start.Base.SHA+"..HEAD") != "1" {
		t.Fatalf("runner claim, exact commit count or projected result missing: %+v", status)
	}
	var v3 int
	if err := f.worktreeStartRig.db.QueryRowContext(ctx, "SELECT count(*) FROM verb_requests").Scan(&v3); err != nil || v3 != 0 {
		t.Fatalf("V3 put crossed local finish: %d %v", v3, err)
	}
	world, err := store.OpenReadOnly(f.c.World)
	if err != nil {
		t.Fatal(err)
	}
	defer world.Close()
	if _, found, err := world.GetRegistryHead(ctx, "world/episode-finish/"+start.Episode); err != nil || found {
		t.Fatalf("finish record crossed local finish: %v %v", found, err)
	}
	output := request("GET", status.StatusHref, nil, nil) + request("GET", "/surface/requests.json", nil, nil) +
		request("GET", "/surface/home-summary.json", nil, nil) + log.String() + runnerLog.String()
	for _, secret := range []string{filename, "private-task-fixture-marker", "BEGIN SSH SIGNATURE", f.token, "finish_manifest", j.Checks[0].Ref} {
		if strings.Contains(output, secret) {
			t.Fatalf("routine output exposed %q", secret)
		}
	}
	if got := presetGit(t, repo, "--git-dir="+bare, "rev-parse", "HEAD"); got != remoteHead {
		t.Fatalf("local runner changed the bare remote: %s", got)
	}
	if _, err := os.Lstat(trace); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("remote write call escaped the worker: %v", err)
	}
}
