package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/agentrun"
	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

const browserStart = "episode start ep --preset p --launch --headless"
const browserTask = "# Task\n<script>alert(1)</script> ; $(literal) &"
const browserFakeChild = fakeLaunchMarker + "while [ ! -e release ]; do :; done\nexit 0\n"

func browserPreset(t *testing.T, r *worktreeStartRig, c runnerConfig) (store.OperatorStartReview, *resolvedPreset) {
	t.Helper()
	r.editPreset(t, func(p *launchPreset) {
		p.Model = "gpt-6.2-sol"
		p.Launch = &presetLaunch{Mode: "interactive", DeadlineS: 60}
		p.Checks = p.Checks[:1]
	})
	ctx, cancel := r.env.operatorContext(activationTestBudget)
	defer cancel()
	p, err := loadLaunchPreset(ctx, "p", r.src.presets, r.src.profiles, r.src.workspace)
	if err != nil {
		t.Fatal(err)
	}
	entry, ok := operatorCatalogueEntry(ctx, c, r.src.workspace, p, r.health)
	if !ok {
		t.Fatal("the fake installed preset did not enter the catalogue")
	}
	return entry, p
}

func publishBrowserCatalogue(t *testing.T, r *worktreeStartRig, c runnerConfig, h store.OperatorHeartbeat, entry store.OperatorStartReview) store.OperatorHeartbeat {
	t.Helper()
	catalog, err := json.Marshal(map[string]store.OperatorStartReview{"p": entry})
	if err != nil {
		t.Fatal(err)
	}
	h.Build = daemon.CurrentBuild().ExecutableSHA256
	h.Owner.Generation = 0 // Publish a new fake owner generation for the real HTTP daemon's build.
	r.env.operatorCatalogue = func(runnerConfig) []byte { return catalog }
	if err := r.env.runnerTick(c, &h); err != nil {
		t.Fatal(err)
	}
	fresh, err := runnerHeartbeat(c.Credentials, activationTestBudget)
	if err != nil || !fresh.Available(h.Build, "start", r.src.clock()) || fresh.CatalogueGeneration <= 0 {
		t.Fatalf("fake catalogue did not reach the read boundary: %v", err)
	}
	return fresh
}

func TestOperatorFakePresetEligibility(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	r, c, h, _ := operatorActivationFixtureWithTask(t, browserFakeChild, browserTask)
	entry, p := browserPreset(t, r, c)
	h = publishBrowserCatalogue(t, r, c, h, entry)
	var out, stderr syncBuffer
	if rc := runPreset([]string{"check", "--presets-dir", r.src.presets, "--exec-profiles-dir", r.src.profiles, "p"}, &out, &stderr); rc != exitOK {
		t.Fatalf("fake preset check: %d %s", rc, &stderr)
	}
	t.Logf("preset check: %s", out.String())
	if !strings.Contains(out.String(), "check short: sandboxed: ") || strings.Contains(out.String(), "UNSANDBOXED") ||
		!strings.Contains(out.String(), "base dev @ "+p.BaseSHA[:7]) {
		t.Fatal("resolved preset parity failed")
	}
	mode, _, err := launchMode(p.Preset, false, false)
	if err != nil || mode != "interactive" {
		t.Fatalf("preset default: %s %v", mode, err)
	}
	words, err := verbs.Split(browserStart)
	if err != nil {
		t.Fatal(err)
	}
	policy := verbs.OperatorPolicy{Start: "A", Finish: "A"}
	explicit, implicit := verbs.AdmitOperator(words, policy), verbs.AdmitOperator(words[:len(words)-1], policy)
	if explicit.Intent == nil || !explicit.Intent.Headless || implicit.Class != verbs.OperatorTerminal || implicit.Intent != nil {
		t.Fatal("only explicit headless may override the interactive fixture")
	}
	mode, wall, err := launchMode(p.Preset, explicit.Intent.Headless, false)
	if err != nil || mode != "headless" || wall != time.Minute {
		t.Fatalf("explicit override: %s %s %v", mode, wall, err)
	}
	a, err := entry.Action("ep", browserTask, h.CatalogueGeneration)
	if err != nil {
		t.Fatal(err)
	}
	a.ID = fixtureActionID
	wire, err := approval.CanonicalOperatorAction(a)
	var decoded approval.OperatorAction
	if err == nil {
		decoded, err = approval.DecodeOperatorAction(wire)
	}
	if err != nil || !decoded.Options.Headless || !slices.Equal(decoded.Argv, explicit.Intent.Argv) ||
		!strings.Contains(decoded.Review, "Headless codex gpt-6.2-sol, cap 1m0s") || !strings.Contains(string(wire), `"headless":true`) {
		t.Fatalf("signed review mode: %v %v", err, decoded.Options)
	}
	digest, err := approval.OperatorSignatureDigest(decoded)
	if err != nil {
		t.Fatal(err)
	}
	challenge := [32]byte{1}
	msg, err := approval.OperatorActionMessage(digest, challenge, 1600)
	if err != nil {
		t.Fatal(err)
	}
	key := approvaltest.NewSigner(89)
	proof, err := approval.CheckOperatorAction(approval.OperatorCheckInput{Action: decoded, Challenge: challenge, Expiry: 1600, Now: 1010,
		Artifact: key.Sign(approval.OperatorActionNamespace, msg), Format: approval.FormatSSH,
		Approver: approval.Approver{Name: "fixture", Format: approval.FormatSSH, Key: key.PublicKey(), EnrolledAt: 900}})
	if err != nil || proof.Binds(decoded, challenge, 1600) != nil {
		t.Fatalf("the mode was not bound by fixture proof: %v", err)
	}

	ctx, cancel := r.env.operatorContext(activationTestBudget)
	defer cancel()
	if _, err := reconstructOperatorStart(ctx, decoded, store.OperatorMintAttempt{Catalogue: h.CatalogueGeneration},
		r.src.presets, r.src.profiles, r.src.workspace, r.src.hosts, r.src.clock); err != nil {
		t.Fatalf("same signed headless start did not re-resolve: %v", err)
	}
	t.Log("mode=interactive; explicit=headless; cap=1m0s; headless is in canonical signed bytes")
}

func TestOperatorStartWalkthroughInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	t.Setenv("HOME", t.TempDir())
	r, c, h, runner := operatorActivationFixtureWithTask(t, browserFakeChild, browserTask)
	entry, _ := browserPreset(t, r, c)
	h = publishBrowserCatalogue(t, r, c, h, entry)
	writeCtx, cancel := r.env.operatorContext(activationTestBudget)
	_, err := r.db.ExecContext(writeCtx, "DELETE FROM operator_action_steps; DELETE FROM operator_actions")
	cancel()
	if err != nil {
		t.Fatal(err)
	}
	log := &syncBuffer{}
	d, err := daemon.New(activationTestContext(t), daemon.Config{DBPath: c.World, CredentialsPath: c.Credentials,
		BindHost: daemon.DefaultBindHost, ErrorLog: log, ApprovalClock: r.src.clock})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Close() })
	srv := httptest.NewServer(d.Handler())
	t.Cleanup(func() { srv.CloseClientConnections(); srv.Close() })
	origin := strings.Replace(srv.URL, "127.0.0.1", "localhost", 1)
	tab := map[string]any{"authenticator": "uv", "inject": "getRecorder", "viewport": map[string]any{"width": 390, "height": 844}}
	b.do(t, nil, "tab", tab)
	virtualPasskey(t, b, c.Credentials, origin, "fixture-sol", 93, r.src.clock()-10)
	wait := func(js string) { t.Helper(); b.do(t, nil, "waitFor", map[string]any{"expression": js, "ms": 10000}) }
	open := func(path string) {
		t.Helper()
		b.do(t, nil, "openCard", map[string]any{"url": origin + path, "raw": true})
	}
	tap := func(id string) { t.Helper(); b.do(t, nil, "tap", map[string]any{"id": id}) }
	fill := func(id, text string) {
		t.Helper()
		wait(fmt.Sprintf(`(() => { const e=document.getElementById(%q); e.value=%q; e.dispatchEvent(new Event('input')); return true; })()`, id, text))
	}
	key := func(k string) { t.Helper(); b.do(t, nil, "key", map[string]any{"key": k}) }
	enter := func(line string) {
		t.Helper()
		wait(fmt.Sprintf(`(() => { const e=document.getElementById('addr'); e.focus(); e.value=%q; return document.activeElement===e; })()`, line))
		key("Enter")
	}
	counts := func(want [4]int) {
		t.Helper()
		var got [4]int
		err := r.db.QueryRowContext(activationTestContext(t), `SELECT (SELECT count(*) FROM operator_actions),
  (SELECT count(*) FROM episode_starts),(SELECT count(*) FROM session_credentials),
  (SELECT sign_count FROM webauthn_credentials WHERE name='fixture-sol')`).Scan(&got[0], &got[1], &got[2], &got[3])
		if err != nil || got != want {
			t.Fatalf("actions/starts/sessions/counter: %v want %v (%v)", got, want, err)
		}
	}
	open("/surface")
	wait(`Boolean(window.World)`)
	noEffects := activationEvidence(t, r)
	enter(strings.TrimSuffix(browserStart, " --headless"))
	wait(`Boolean(document.querySelector('#verb .verb-body')?.textContent.includes('terminal') || document.getElementById('hmsg').textContent)`)
	if !reflect.DeepEqual(activationEvidence(t, r), noEffects) {
		t.Fatal("the preset's implicit interactive mode inserted an action")
	}
	doc, err := os.ReadFile(filepath.Join("..", "..", "docs", "QUICKSTART.md"))
	if err != nil {
		t.Fatal(err)
	}
	_, example, present := strings.Cut(string(doc), "<!-- operator-browser-fixture -->\n")
	line, rest, lineOK := strings.Cut(strings.TrimPrefix(example, "```text\n"), "\n```")
	_, rest, taskStart := strings.Cut(rest, "```text\n")
	task, _, taskEnd := strings.Cut(rest, "\n```")
	if !present || !lineOK || !taskStart || !taskEnd || line != browserStart || task != browserTask {
		t.Fatal("QUICKSTART fixture did not match the actual preset and literal task")
	}
	enter(line)
	wait(`document.activeElement?.id==='operator-task' && document.querySelector('#verb .verb-body')?.textContent.includes('Headless codex gpt-6.2-sol')`)
	counts([4]int{0, 0, 0, 1})
	key("n")
	wait(`document.activeElement?.id==='operator-task' && document.getElementById('operator-task').value==='n'`)
	fill("operator-task", "")
	if !reflect.DeepEqual(activationEvidence(t, r), noEffects) {
		t.Fatal("Enter changed the start evidence")
	}
	wait(`innerWidth===390 && innerHeight===844 && document.documentElement.scrollWidth===390 && document.getElementById('operator-task').getBoundingClientRect().right<=390`)
	key("Tab")
	wait(`document.activeElement?.id==='operator-request'`)
	key("Enter")
	wait(`document.getElementById('operator-request-info')?.textContent.includes('1–8192 UTF-8 bytes')`)
	counts([4]int{0, 0, 0, 1})
	fill("operator-task", task)
	wait(`(() => {
  const original=window.fetch.bind(window);
  window.fetch=async (path,options) => {
    const response=await original(path,options);
    if(path==='/surface/operator-actions' && options?.method==='POST') {
      window.__fixtureTaskFrozen = document.getElementById('operator-task').readOnly;
      throw Error('fixture dropped POST response');
    }
    return response;
  };
  return true;
 })()`)
	tap("operator-request")
	wait(`document.getElementById('operator-request-info')?.textContent.includes('Request outcome unavailable') && window.__fixtureTaskFrozen===true`)
	counts([4]int{1, 0, 0, 1})
	open("/surface")
	wait(`document.activeElement?.id==='operator-request' && document.getElementById('operator-request').textContent==='Reconnect' && document.getElementById('operator-task').readOnly`)
	tap("operator-request")
	wait(`document.activeElement?.id==='operator-title' && document.getElementById('operator-key')?.disabled===false && !document.getElementById('operator-controls').hidden`)
	counts([4]int{1, 0, 0, 1})
	ctx := activationTestContext(t)
	ro, err := store.OpenCredentialsReadOnly(c.Credentials)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ro.Close() })
	var id string
	var operations int
	if err := r.db.QueryRowContext(ctx, "SELECT id,(SELECT count(*) FROM operator_action_steps WHERE kind='operation') FROM operator_actions").Scan(&id, &operations); err != nil || operations != 1 {
		t.Fatalf("reconnect operation count=%d: %v", operations, err)
	}
	row, err := ro.OperatorAction(ctx, id, r.src.clock())
	if err != nil {
		t.Fatal(err)
	}
	r.action, err = approval.DecodeOperatorAction(row.Action)
	if err != nil || r.action.Task != browserTask || !r.action.Options.Headless {
		t.Fatalf("literal signed task or mode: %v", err)
	}
	wait(fmt.Sprintf(`document.getElementById('operator-task').textContent===%q && !document.getElementById('operator-task').querySelector('script') && document.documentElement.scrollWidth<=innerWidth`, browserTask))
	key("Tab")
	wait(`document.activeElement?.id==='operator-key'`)
	wait(`document.querySelector('#operator-controls details').open = true`)
	fill("operator-signature", "y")
	before, err := os.ReadFile(c.Credentials)
	if err != nil {
		t.Fatal(err)
	}
	tap("operator-sign")
	wait(`document.getElementById('operator-answer')?.textContent.includes('NotApproved') && !document.getElementById('operator-key').disabled`)
	counts([4]int{1, 0, 0, 1})
	after, err := os.ReadFile(c.Credentials)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("plain y changed the pending row or key counter")
	}
	tap("operator-key")
	wait(`document.getElementById('operator-answer')?.textContent==='Confirmed by webauthn:fixture-sol'`)
	counts([4]int{1, 0, 0, 2})
	var web struct {
		Gets       []cardGet
		Violations []string
	}
	b.do(t, &web, "waitFor", map[string]any{"expression": `({gets:window.__gets||[],violations:window.__cspv||[]})`})
	challenge, err := approval.OperatorActionMessageHash([32]byte(row.Digest), [32]byte(row.Challenge), row.ExpiresAt)
	if err != nil || len(web.Gets) != 1 || web.Gets[0].UserVerification != "required" ||
		web.Gets[0].Challenge != base64.RawURLEncoding.EncodeToString(challenge[:]) || len(web.Violations) != 0 {
		t.Fatalf("browser ceremony: %+v", web)
	}
	b.do(t, nil, "tab", map[string]any{"viewport": tab["viewport"]}) // close the approved page before running a single fake launch

	f := &operatorLaunchRig{worktreeStartRig: r, c: c, h: h, runner: runner, ctx: ctx, db: ro, now: r.src.clock()}
	if out := f.tick(t); out != "" {
		t.Fatalf("prepare: %s", out)
	}
	f.readPrepared(t)
	health := r.env.client.Transport
	r.env.client.Transport = serviceTransport(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/v1/session/self" {
			return health.RoundTrip(req)
		}
		if req.URL.Host != c.Address || req.Header.Get("Authorization") != "Bearer "+f.token {
			t.Error("session read escaped its fixture binding")
		}
		data, _ := json.Marshal(f.self)
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(data))}, nil
	})
	// Only the fixture launcher is entered. Keep the fake child alive until the stored birth is read.
	r.env.operatorSupervisor = func(job operatorLaunchJob, out io.Writer) error { return r.env.executeOperatorLaunch(ctx, job, out) }
	r.env.operatorAgentRun = func(ctx context.Context, argv []string, dir string, env []string, out, errs *os.File, wall time.Duration, hooks agentrun.LaunchHooks) (agentrun.Result, error) {
		prompts := 0
		for _, arg := range argv {
			if strings.Contains(arg, browserTask) {
				prompts++
			}
		}
		if !slices.Contains(argv, "read-only") || !strings.HasPrefix(argv[len(argv)-1], browserTask) || prompts != 1 ||
			!slices.Contains(env, "WORLD_SESSION="+f.token) || wall != time.Minute {
			t.Error("fake launch did not preserve confinement or the literal prompt")
		}
		start := hooks.Started
		hooks.Started = func(pid int) error {
			defer func() {
				if err := os.WriteFile(filepath.Join(dir, "release"), nil, 0600); err != nil {
					t.Error(err)
				}
			}()
			if err := start(pid); err != nil {
				return err
			}
			if got := f.receipt(t).ProcessStatus(agentrun.ProcessIdentity); got != "launched" {
				t.Errorf("live birth readback: %s", got)
			}
			return nil
		}
		return agentrun.RunObserved(ctx, argv, dir, env, out, errs, wall, hooks)
	}
	if out := f.tick(t); !strings.Contains(out, ": launched") || !strings.Contains(out, ": stopped") {
		t.Fatalf("fake launch: %s", out)
	}
	p := "/surface/operator-actions/" + id
	open(p)
	wait(`document.getElementById('operator-state')?.textContent==='stopped' && document.getElementById('operator-report')?.getAttribute('href')==='/v1/episodes/ep/report' &&
   document.activeElement?.id==='operator-title' && !document.getElementById('operator-refresh').disabled && document.documentElement.scrollWidth<=innerWidth`)
	tap("operator-refresh")
	wait(`document.getElementById('operator-state')?.textContent==='stopped' && !document.getElementById('operator-refresh').disabled`)
	counts([4]int{1, 1, 1, 2})
	worktrees := presetGit(t, r.f.r.f, "worktree", "list", "--porcelain")
	physical, err := filepath.EvalSymlinks(r.path())
	ran, readErr := os.ReadFile(filepath.Join(r.bin, "ran"))
	matches := 0
	for _, line := range strings.Split(worktrees, "\n") {
		path, ok := strings.CutPrefix(line, "worktree ")
		if !ok {
			continue
		}
		linked, e := filepath.EvalSymlinks(path)
		if e != nil {
			t.Fatal(e)
		}
		if linked == physical {
			matches++
		}
	}
	if err != nil || readErr != nil || string(ran) != "x" || strings.Count(worktrees, "worktree ") != 2 ||
		matches != 1 || f.receipt(t).Receipt.Run != id {
		t.Fatalf("durable worktree or launch: %v %v %q (%d matching registrations)", err, readErr, ran, matches)
	}
	for _, endpoint := range []string{"/surface/requests.json", p + "/status.json"} {
		var text string
		b.do(t, nil, "waitFor", map[string]any{"expression": fmt.Sprintf("(window.__read=undefined, fetch(%q).then(r=>r.text()).then(x=>{window.__read=x}), true)", endpoint)})
		b.do(t, &text, "waitFor", map[string]any{"expression": "window.__read ?? null"})
		if strings.Contains(text, "alert(1)") || strings.Contains(text, "authenticatorData") || strings.Contains(text, "messageHex") {
			t.Fatal("routine read revealed task or proof")
		}
	}
	if strings.Contains(log.String(), "alert(1)") || strings.Contains(log.String(), "authenticatorData") {
		t.Fatal("daemon log revealed task or proof")
	}
	tap("operator-back")
	wait(`location.pathname==='/surface' && !location.search && !location.hash && Boolean(document.getElementById('addr'))`)
	b.do(t, nil, "tab", map[string]any{})
}
