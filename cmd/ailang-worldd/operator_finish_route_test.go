package main

import (
	"bytes"
	"context"
	"crypto/sha256"
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
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/checkrun"
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
	effectTrace := filepath.Join(f.home, "git-effect-trace")
	tools := f.c.ToolDirs[0]
	f.c.FinishRemotes = []string{fixtureFinishRemote}
	writeFinishTestTool(t, &f.c, tools, "gh", "case \"$1 $2\" in 'auth status'|'label list') ;; 'repo view') printf 'WRITE\\n';; *) touch "+shellQuote(trace)+"; exit 96;; esac\n")
	writeFinishTestTool(t, &f.c, tools, "git", "for arg; do case \"$arg\" in push|fetch|clone|pull|send-pack|fetch-pack|receive-pack|upload-pack|remote|submodule) touch "+shellQuote(trace)+"; exit 96;; esac; done\n"+
		"case \" $* \" in *' ls-remote '*) exec "+shellQuote(repo.git)+" --git-dir=/dev/null ls-remote --symref -- "+shellQuote(bare)+" HEAD;; esac\n"+
		"for arg; do case \"$arg\" in commit-tree|update-ref|read-tree) printf '%s\\n' \"$arg\" >> "+shellQuote(effectTrace)+";; esac; done\n"+
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
	reportDown, reportMalformed := false, false
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
				if reportMalformed {
					raw, _ := json.Marshal(rep)
					fmt.Fprintf(w, `%s,"exec_cache":true}`, raw[:len(raw)-1])
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
		ID, Href, StatusHref, State, Kind, Reason, LocalCommit string
		Checks                                                 []struct{ Name, Result string }
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

	row, err = f.db.OperatorAction(ctx, status.ID, f.now)
	if err != nil || row.State != "confirmed" || row.Owner != "" {
		t.Fatalf("signed route did not leave an unclaimed row: %s %v", row.State, err)
	}

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

	probe := finishRechecks{f: f, start: start, report: &rep, row: row, sign: func(msg []byte) []byte { return key.Sign(approval.OperatorActionNamespace, msg) },
		run: run, bare: bare, trace: trace, change: change, effectTrace: effectTrace, startWire: startWire,
		checks: &checks, outcomes: &outcomes, malformed: &reportMalformed}
	probe.inputs(t)
	probe.catalogue(t)
	probe.claim(t, raw)
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
	probe.recovery(t, row, j, func(reason, sha string) {
		t.Helper()
		status.Reason, status.LocalCommit, status.Checks = "", "", nil
		request("GET", status.StatusHref, nil, &status)
		if status.Reason != reason || status.LocalCommit != sha || status.State != row.State ||
			status.ID != row.ID || status.Kind != "finish" || len(status.Checks) != 1 ||
			status.Checks[0].Name != start.Checks[0].Name || status.Checks[0].Result != "passed" {
			t.Fatalf("finish receipt status: %+v, want %q and SHA %q", status, reason, sha)
		}
		if card := request("GET", status.Href, nil, nil); !strings.Contains(card, reason) {
			t.Fatalf("finish card lost the projected reason %q", reason)
		}
	})
}

type finishRechecks struct {
	f                                *operatorLaunchRig
	start                            *store.EpisodeStart
	report                           *daemon.EpisodeReport
	row                              store.OperatorActionRow
	sign                             func([]byte) []byte
	run                              func() error
	bare, trace, change, effectTrace string
	startWire                        []byte
	checks, outcomes                 *int
	malformed                        *bool
}

// Observe the entire credentials store (including action generations, passkey
// counters, sessions and reservations), world store and both Git refs.
func (p finishRechecks) state(t *testing.T) [8]string {
	t.Helper()
	f := p.f
	if _, err := os.Lstat(p.trace); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("unexpected dispatch or remote write: %v", err)
	}
	bareHead := readStartFixtureFile(t, filepath.Join(p.bare, "HEAD"))
	bareRefs := serviceFileHash(filepath.Join(p.bare, "packed-refs"))
	if ref, symbolic := strings.CutPrefix(strings.TrimSpace(string(bareHead)), "ref: "); symbolic {
		bareRefs += serviceFileHash(filepath.Join(p.bare, ref))
	}
	return [8]string{serviceFileHash(f.creds), serviceFileHash(f.c.World),
		serviceFileHash(filepath.Join(f.f.r.f.repo, ".git", "refs", "heads", p.start.Branch)),
		serviceHash(bareHead), bareRefs,
		serviceFileHash(p.effectTrace), serviceFileHash(filepath.Join(f.bin, "ran")),
		fmt.Sprintf("%d/%d", *p.checks, *p.outcomes)}
}

func (p finishRechecks) unchanged(t *testing.T, want error, run func() error) {
	t.Helper()
	before := p.state(t)
	got := run()
	after := p.state(t)
	if !errors.Is(got, want) || before != after {
		t.Fatalf("want %v without changing row, journal, refs or counters, got %v", want, got)
	}
	switch want {
	case errOperatorFinishEligibility, errOperatorFinishReview, errOperatorFinishFence:
		if finishRefusal(got) != want.Error() {
			t.Fatalf("runner lost the refusal remedy for %v", want)
		}
	}
}

func (p finishRechecks) exec(t *testing.T, query string, args ...any) {
	t.Helper()
	ctx, stop := p.f.env.operatorContext(activationTestBudget)
	defer stop()
	if _, err := p.f.worktreeStartRig.db.ExecContext(ctx, query, args...); err != nil {
		t.Fatal(err)
	}
}

func (p finishRechecks) setStart(t *testing.T, start *store.EpisodeStart) {
	t.Helper()
	b, err := json.Marshal(start)
	if err != nil {
		t.Fatal(err)
	}
	p.exec(t, "UPDATE episode_starts SET start_json=? WHERE episode=?", string(b), start.Episode)
	p.report.Start = start
}

// Reuse the genuinely route-confirmed row and the verified episode flock.
// Each row can reach the exact claim preparation without repeating GitHub
// readiness, which the surrounding route turn exercises with fixture tools.
func (p finishRechecks) preclaim(t *testing.T, ctx context.Context) error {
	t.Helper()
	f := p.f
	held := f.borrow(t)
	defer held.Close()
	fd, err := held.Verify(f.creds, p.start.Episode)
	if err != nil {
		return err
	}
	row, err := f.db.OperatorAction(ctx, p.row.ID, f.now)
	if err != nil {
		return err
	}
	a, m, generation, err := finishManifestForRow(f.c, row)
	if err != nil {
		return err
	}
	_, err = f.env.finishWork(ctx, f.c, f.h, f.src, f.db, fd, row, operatorPreparationCAS(row), a, m, generation, false, f.now)
	return err
}

func (p finishRechecks) inputs(t *testing.T) {
	t.Helper()
	f := p.f
	originalReport, health, config := *p.report, f.health, f.c
	preset := f.f.file
	presetWire := readStartFixtureFile(t, preset)
	token := filepath.Join(f.f.r.sessions, p.start.Episode+".session")
	sessions, err := f.db.ListSessions(activationTestContext(t), p.start.Episode, 2)
	if err != nil || len(sessions) != 1 {
		t.Fatalf("original session: rows=%d err=%v", len(sessions), err)
	}
	launch := f.receipt(t).Receipt
	launchWire, _ := json.Marshal(launch)
	alterPreset := func(t *testing.T, edit func(*launchPreset, *presetDaemon)) {
		t.Helper()
		var l launchPreset
		var d presetDaemon
		if err := json.Unmarshal(presetWire, &l); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(l.Daemon, &d); err != nil {
			t.Fatal(err)
		}
		edit(&l, &d)
		l.Daemon, _ = json.Marshal(d)
		b, err := json.Marshal(l)
		if err != nil {
			t.Fatal(err)
		}
		seedStartFixtureFile(t, preset, b)
		changed := *p.start
		changed.Preset.Bytes, changed.Preset.SHA256 = string(b), serviceHash(b)
		p.setStart(t, &changed)
	}
	chmod := func(t *testing.T, path string, mode os.FileMode) {
		t.Helper()
		if err := os.Chmod(path, mode); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			want := os.FileMode(0600)
			if path == f.f.r.sessions {
				want = 0700
			}
			if err := os.Chmod(path, want); err != nil {
				t.Error(err)
			}
		})
	}
	input := func(t *testing.T) error {
		bound := activationTestContext(t)
		_, err := f.env.finishInputs(bound, f.c, f.src, f.db, p.start.Episode, nil)
		return err
	}
	p.unchanged(t, nil, func() error { return input(t) })
	for _, tc := range []struct {
		name string
		want error
		edit func(*testing.T)
	}{
		{"resumes", errOperatorFinishEligibility, func(t *testing.T) {
			changed := *p.start
			changed.Resumes = []store.EpisodeResume{{CredentialID: sessions[0].CredentialID}}
			p.setStart(t, &changed)
		}},
		{"preset-allowlist", errOperatorFinishEligibility, func(*testing.T) { f.c.PresetNames = nil }},
		{"remote-allowlist", errOperatorFinishEligibility, func(*testing.T) { f.c.FinishRemotes = nil }},
		{"launch-version", errOperatorFinishEligibility, func(t *testing.T) {
			old := launch
			old.Status = "launched"
			b, _ := json.Marshal(old)
			p.exec(t, "UPDATE operator_action_steps SET version=3, receipt=? WHERE kind='launch'", b)
		}},
		{"health-status", errOperatorFinishEligibility, func(*testing.T) { f.health.Status = "starting" }},
		{"build", errOperatorFinishEligibility, func(*testing.T) { f.health.Build.ExecutableSHA256 = "sha256:other-fixture" }},
		{"report-read", errOperatorFinishEligibility, func(*testing.T) { *p.malformed = true }},
		{"report-episode", errOperatorFinishEligibility, func(*testing.T) { p.report.Episode = "other-episode" }},
		{"report-intents", errOperatorFinishEligibility, func(*testing.T) { p.report.Uncommitted = 1 }},
		{"recorded-finish", errOperatorFinishEligibility, func(*testing.T) { p.report.Finish = &daemon.EpisodeFinish{} }},
		{"report-start", errOperatorFinishEligibility, func(*testing.T) {
			copy := *p.start
			copy.Model += "-other"
			p.report.Start = &copy
		}},
		{"preset-decode", errOperatorFinishReview, func(t *testing.T) {
			alterPreset(t, func(l *launchPreset, _ *presetDaemon) { l.Name = "other-preset" })
		}},
		{"workspace", errOperatorFinishReview, func(t *testing.T) {
			alterPreset(t, func(_ *launchPreset, d *presetDaemon) { d.WorkspaceRoot = t.TempDir() })
		}},
		{"world", errOperatorFinishReview, func(t *testing.T) {
			other := filepath.Join(t.TempDir(), "other-world.db")
			seedStartFixtureFile(t, other, readStartFixtureFile(t, f.c.World))
			alterPreset(t, func(_ *launchPreset, d *presetDaemon) { d.DB = other })
		}},
		{"store", errOperatorFinishReview, func(t *testing.T) {
			other := filepath.Join(t.TempDir(), "other-credentials.db")
			seedStartFixtureFile(t, other, readStartFixtureFile(t, f.creds))
			alterPreset(t, func(_ *launchPreset, d *presetDaemon) { d.CredentialsDB = other })
		}},
		{"profile", errOperatorFinishReview, func(t *testing.T) {
			p.exec(t, "UPDATE session_credentials SET profile_digest=? WHERE episode_id=?", "sha256:other", p.start.Episode)
		}},
		{"pinned-check", errOperatorFinishReview, func(t *testing.T) {
			changed := *p.start
			changed.Checks = append([]store.StartCheck(nil), p.start.Checks...)
			changed.Checks[0].Name = "different-pin"
			p.setStart(t, &changed)
		}},
		{"sessions", errOperatorFinishEligibility, func(t *testing.T) {
			p.exec(t, `INSERT INTO session_credentials SELECT 'zzz-fixture', episode_id, grants_json, expires_at,
				created_at+1, project, profile_digest FROM session_credentials WHERE episode_id=?`, p.start.Episode)
		}},
		{"project", errOperatorFinishEligibility, func(t *testing.T) {
			p.exec(t, "UPDATE session_credentials SET project=? WHERE episode_id=?", "other-project", p.start.Episode)
		}},
		{"session-file", errOperatorFinishEligibility, func(t *testing.T) { chmod(t, f.f.r.sessions, 0755) }},
		{"token-mode", errOperatorFinishEligibility, func(t *testing.T) { chmod(t, token, 0644) }},
		{"commit-identity", errOperatorFinishFence, func(t *testing.T) {
			alterPreset(t, func(l *launchPreset, _ *presetDaemon) { l.Committer = "" })
		}},
		{"reused-pid", nil, func(t *testing.T) {
			old := launch
			old.PID, old.Birth = os.Getpid(), "previous-process-birth"
			b, _ := json.Marshal(old)
			p.exec(t, "UPDATE operator_action_steps SET receipt=? WHERE kind='launch'", b)
		}},
	} {
		t.Run("input/"+tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				f.c, f.health, *p.report, *p.malformed = config, health, originalReport, false
				seedStartFixtureFile(t, preset, presetWire)
				p.exec(t, "UPDATE episode_starts SET start_json=? WHERE episode=?", string(p.startWire), p.start.Episode)
				p.exec(t, "UPDATE operator_action_steps SET version=4, receipt=? WHERE kind='launch'", launchWire)
				p.exec(t, "DELETE FROM session_credentials WHERE credential_id='zzz-fixture'")
				p.exec(t, "UPDATE session_credentials SET project=?,profile_digest=? WHERE credential_id=?",
					sessions[0].Project, sessions[0].ProfileDigest, sessions[0].CredentialID)
			})
			tc.edit(t)
			p.unchanged(t, tc.want, func() error { return input(t) })
			// Also observe the route-claimed worker's earlier admission
			// boundary. The isolated input assertion above distinguishes
			// the branch from a later frozen-manifest refusal.
			if tc.want != nil {
				want := tc.want
				if tc.name == "remote-allowlist" {
					want = approval.ErrNotVerified // frozen remote checked before live inputs
				}
				p.unchanged(t, want, func() error { return p.preclaim(t, activationTestContext(t)) })
			}
		})
	}
}

func (p finishRechecks) catalogue(t *testing.T) {
	t.Helper()
	f := p.f
	ctx := activationTestContext(t)
	t.Run("catalogue/approval-size", func(t *testing.T) {
		report := *p.report
		t.Cleanup(func() {
			p.setStart(t, p.start)
			*p.report = report
		})
		start := *p.start
		task := *start.Task
		task.Title = strings.Repeat("x", approval.MaxOperatorActionBytes)
		start.Task = &task
		p.setStart(t, &start)
		in, err := f.env.finishInputs(ctx, f.c, f.src, f.db, start.Episode, nil)
		if err != nil {
			t.Fatal(err)
		}
		m, err := in.frozen(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := m.Action(1, false); !errors.Is(err, approval.ErrOperatorTooLarge) {
			t.Fatalf("fixture did not reach approval byte bound: %v", err)
		}
		before := p.state(t)
		if _, ok := f.env.finishCatalogue(ctx, f.c, f.src)[start.Episode]; ok || before != p.state(t) {
			t.Fatal("unapprovable finish entered the runner catalogue")
		}
	})
	t.Run("catalogue/aggregate-size", func(t *testing.T) {
		c := f.c
		c.PresetNames = []string{"padding", p.start.Preset.Name}
		path := filepath.Join(f.src.presets, "padding.json")
		var preset launchPreset
		if err := json.Unmarshal([]byte(p.start.Preset.Bytes), &preset); err != nil {
			t.Fatal(err)
		}
		preset.Name = c.PresetNames[0]
		raw, _ := json.Marshal(preset)
		seedStartFixtureFile(t, path, raw)
		defer os.Remove(path)
		capBytes := 16 * approval.MaxOperatorActionBytes
		small := f.env.readOperatorCatalogue(c, f.src)
		var all map[string]json.RawMessage
		if err := json.Unmarshal(small, &all); err != nil || all["padding"] == nil || all[store.OperatorFinishCatalogueKey] == nil {
			t.Fatalf("fixture never published the start/finish catalogue: %d bytes, %v", len(small), err)
		}
		// Whitespace in the preset is part of its exact catalogued bytes.
		seedStartFixtureFile(t, path, append(raw, bytes.Repeat([]byte(" "), capBytes-len(small)+1)...))
		before := p.state(t)
		out := f.env.readOperatorCatalogue(c, f.src)
		all = nil
		if err := json.Unmarshal(out, &all); err != nil || len(out) > capBytes || all["padding"] == nil ||
			all[p.start.Preset.Name] == nil || all[store.OperatorFinishCatalogueKey] != nil || before != p.state(t) {
			t.Fatalf("over-limit finish catalogue did not preserve start entries: %d bytes; %v", len(out), err)
		}
	})
}

// Cancel at the lifecycle boundary using the same nonblocking lock as the
// caller. The signal is driven by lock ownership, never by elapsed time.
type finishLifecycleContext struct {
	context.Context
	cancel context.CancelFunc
	root   string
	mu     sync.Mutex
}

func (c *finishLifecycleContext) Err() error {
	// Child contexts may probe Err concurrently. Only production lock
	// ownership may trigger cancellation.
	c.mu.Lock()
	defer c.mu.Unlock()
	lock, err := runnerLifecycleLock(c.root)
	if lock != nil {
		_ = lock.Close()
	} else if errors.Is(err, syscall.EWOULDBLOCK) {
		c.cancel()
	}
	return c.Context.Err()
}

func (p finishRechecks) claim(t *testing.T, raw []byte) {
	t.Helper()
	f := p.f
	ctx := activationTestContext(t)
	h, err := f.db.OperatorHeartbeat(ctx)
	if err != nil {
		t.Fatal(err)
	}
	env := f.env
	launch := f.receipt(t).Receipt
	launchWire, _ := json.Marshal(launch)
	profile, profilePath, err := loadMintProfile(f.src.profiles, "proj")
	if err != nil {
		t.Fatal(err)
	}
	profileWire := readStartFixtureFile(t, profilePath)
	setLaunch := func(t *testing.T) {
		t.Helper()
		changed := launch
		changed.Birth += "-other"
		b, _ := json.Marshal(changed)
		p.exec(t, "UPDATE operator_action_steps SET receipt=? WHERE kind='launch'", b)
		t.Cleanup(func() { p.exec(t, "UPDATE operator_action_steps SET receipt=? WHERE kind='launch'", launchWire) })
	}
	table := []struct {
		name string
		want error
		edit func(*testing.T)
	}{
		{"catalogue-fresh", store.ErrOperatorRunner, func(t *testing.T) {
			p.exec(t, "UPDATE operator_runner_state SET catalogue_generation=catalogue_generation+1")
		}},
		{"heartbeat-owner", store.ErrOperatorRunner, func(*testing.T) { f.h.Owner.Generation++ }},
		{"heartbeat-ready", store.ErrOperatorRunner, func(*testing.T) { f.h.Build = "different-fixture-build" }},
		{"life-quiet", store.ErrOperatorRunner, func(*testing.T) {
			f.env.runnerProcesses = func() ([]worldProcess, error) { return nil, errors.New("inventory unavailable") }
		}},
		{"life-pair", store.ErrOperatorRunner, func(t *testing.T) {
			path := f.env.runnerPaths().Plist
			b := readStartFixtureFile(t, path)
			seedStartFixtureFile(t, path, append(bytes.Clone(b), '\n'))
			t.Cleanup(func() { seedStartFixtureFile(t, path, b) })
		}},
		{"manifest-protected", errOperatorFinishReview, func(t *testing.T) {
			var edited map[string]any
			if err := json.Unmarshal(profileWire, &edited); err != nil {
				t.Fatal(err)
			}
			edited["protected"] = []string{filepath.Base(p.change)}
			b, _ := json.Marshal(edited)
			seedStartFixtureFile(t, profilePath, b)
			next, _, err := loadMintProfile(f.src.profiles, "proj")
			if err != nil {
				t.Fatal(err)
			}
			p.exec(t, "UPDATE session_credentials SET profile_digest=? WHERE episode_id=?", next.Digest, p.start.Episode)
			t.Cleanup(func() {
				seedStartFixtureFile(t, profilePath, profileWire)
				p.exec(t, "UPDATE session_credentials SET profile_digest=? WHERE episode_id=?", profile.Digest, p.start.Episode)
			})
			p.signCurrentReview(t, true)
		}},
		{"manifest-unattributed", errOperatorFinishReview, func(t *testing.T) {
			old := p.report.Writes.Written
			p.report.Writes.Written = nil
			t.Cleanup(func() { p.report.Writes.Written = old })
			p.signCurrentReview(t, false)
		}},
	}
	for _, tc := range table {
		t.Run("claim/"+tc.name, func(t *testing.T) {
			original := f.h
			t.Cleanup(func() {
				f.h, f.env = original, env
				p.exec(t, "UPDATE operator_runner_state SET catalogue_generation=?", h.CatalogueGeneration)
			})
			tc.edit(t)
			p.unchanged(t, tc.want, func() error { return p.preclaim(t, activationTestContext(t)) })
		})
	}
	t.Run("claim/cancelled", func(t *testing.T) {
		bound, cancel := context.WithCancel(ctx)
		defer cancel()
		lifecycle := &finishLifecycleContext{Context: bound, cancel: cancel, root: filepath.Dir(env.paths.Record)}
		p.unchanged(t, store.ErrOperatorRunner, func() error {
			return p.preclaim(t, lifecycle)
		})
	})
	t.Run("claim/start-ready-loop", func(t *testing.T) {
		original := f.h
		f.h.Readiness = "start"
		f.env.operatorCatalogue = func(runnerConfig) []byte { return raw }
		if err := f.env.touchRunner(f.c, &f.h); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			f.h, f.env = original, env
			p.exec(t, "UPDATE operator_runner_state SET readiness=?, catalogue_generation=?", h.Readiness, h.CatalogueGeneration)
		})
		before := p.state(t)
		if log := f.tick(t); log != "" || before != p.state(t) {
			t.Fatalf("start-only runner attempted a finish: %q", log)
		}
	})

	t.Run("after-claim", func(t *testing.T) {
		held := f.borrow(t)
		defer held.Close()
		fd, err := held.Verify(f.creds, p.start.Episode)
		if err != nil {
			t.Fatal(err)
		}
		if err := runnerWrite(f.creds, func(db *store.Credentials) error {
			write, cancel := f.env.operatorContext(activationTestBudget)
			defer cancel()
			return db.ClaimOperatorAction(write, operatorPreparationCAS(p.row), h.Owner, f.now)
		}); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			p.exec(t, "DELETE FROM operator_action_steps WHERE action_id=? AND kind<>'operation'", p.row.ID)
			p.exec(t, "UPDATE operator_actions SET state=?,generation=?,owner=NULL WHERE id=?", p.row.State, p.row.Generation, p.row.ID)
		})
		row, err := f.db.OperatorAction(ctx, p.row.ID, f.now)
		if err != nil || row.Generation != p.row.Generation+1 || row.Owner != h.Owner.Nonce {
			t.Fatalf("expected exact CAS claim: %s/%d %v", row.Owner, row.Generation, err)
		}
		a, m, generation, err := finishManifestForRow(f.c, row)
		if err != nil {
			t.Fatal(err)
		}
		construct := func() (operatorFinishWork, error) {
			return f.env.finishWork(ctx, f.c, h, f.src, f.db, fd, row, operatorPreparationCAS(row), a, m, generation, false, f.now)
		}
		t.Run("review", func(t *testing.T) {
			setLaunch(t)
			p.unchanged(t, errOperatorFinishReview, func() error { _, err := construct(); return err })
		})
		work, err := construct()
		if err != nil {
			t.Fatal(err)
		}
		// The production check seam must be unreachable in every refusal case.
		work.run = func(context.Context, store.StartCheck) (checkrun.Result, error) {
			*p.checks++
			return checkrun.Result{}, errors.New("unexpected fixture check")
		}
		for _, tc := range []struct {
			name string
			want error
			edit func(*testing.T)
		}{
			{"clock", store.ErrOperatorRunner, func(t *testing.T) {
				p.exec(t, "UPDATE operator_runner_state SET seen_at=?", h.SeenAt-10)
				t.Cleanup(func() { p.exec(t, "UPDATE operator_runner_state SET seen_at=?", h.SeenAt) })
			}},
			{"inputs", errOperatorFinishReview, setLaunch},
			{"key", approval.ErrNotVerified, func(t *testing.T) {
				p.exec(t, "UPDATE approvers SET revoked_at=? WHERE name='finish-operator-fixture'", f.now)
				t.Cleanup(func() { p.exec(t, "UPDATE approvers SET revoked_at=NULL WHERE name='finish-operator-fixture'") })
			}},
		} {
			t.Run(tc.name, func(t *testing.T) {
				tc.edit(t)
				p.unchanged(t, tc.want, func() error { return runEpisodeFinishWithPrep(ctx, work) })
			})
		}
	})
}

// Seed the same canonical approval and genuine fixture signature used by the
// confirm route, bypassing request admission to exercise the later re-check.
func (p finishRechecks) signCurrentReview(t *testing.T, protected bool) {
	t.Helper()
	f := p.f
	ctx := activationTestContext(t)
	i, err := f.env.finishInputs(ctx, f.c, f.src, f.db, p.start.Episode, nil)
	if err != nil {
		t.Fatal(err)
	}
	m, err := i.frozen(ctx)
	if err != nil || protected && (len(m.Protected) != 1 || len(m.Unattributed) != 0) ||
		!protected && (len(m.Unattributed) != 1 || len(m.Protected) != 0) {
		t.Fatalf("fixture did not freeze exact change classification: %v %v %v", m.Protected, m.Unattributed, err)
	}
	_, _, generation, err := finishManifestForRow(f.c, p.row)
	if err != nil {
		t.Fatal(err)
	}
	a, err := m.Action(generation, false)
	if err != nil {
		t.Fatal(err)
	}
	a.ID = p.row.ID
	raw, err := approval.CanonicalOperatorAction(a)
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(raw)
	message, err := approval.OperatorActionMessage(digest, [32]byte(p.row.Challenge), p.row.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	const update = "UPDATE operator_actions SET action=?,digest=?,artifact=? WHERE id=?"
	p.exec(t, update, raw, digest[:], p.sign(message), p.row.ID)
	t.Cleanup(func() { p.exec(t, update, p.row.Action, p.row.Digest, p.row.Artifact, p.row.ID) })
}

func (p finishRechecks) recovery(t *testing.T, row store.OperatorActionRow, complete store.OperatorFinishJournal, status func(string, string)) {
	t.Helper()
	f := p.f
	ctx := activationTestContext(t)
	original, _ := json.Marshal(complete.Commit)
	const reset = "UPDATE operator_action_steps SET phase=?,version=?,receipt=? WHERE action_id=? AND kind='commit'"
	before := p.state(t)
	t.Cleanup(func() {
		p.exec(t, reset, "observed", complete.CommitKey.Version, original, row.ID)
	})
	p.exec(t, reset, "reserved", complete.CommitKey.Version-1, nil, row.ID)
	p.unchanged(t, nil, func() error {
		status("Finish outcome uncertain; inspect durable receipts", "")
		return nil
	})
	now := f.now
	f.now = row.ExpiresAt
	err := p.run()
	f.now = now
	j, readErr := f.db.ReadOperatorFinishJournal(ctx, operatorPreparationCAS(row))
	after := p.state(t)
	if err != nil || readErr != nil || !reflect.DeepEqual(j, complete) || !reflect.DeepEqual(before[1:], after[1:]) {
		t.Fatalf("expired reserved commit did not observe the exact existing object: %v %v %+v", err, readErr, j)
	}

	failed, _ := json.Marshal(store.OperatorFinishLocalCommit{State: "commit-failed"})
	p.exec(t, reset, "observed", complete.CommitKey.Version, failed, row.ID)
	p.unchanged(t, nil, func() error {
		status("Local commit failed; check receipts retained", "")
		return nil
	})
}
