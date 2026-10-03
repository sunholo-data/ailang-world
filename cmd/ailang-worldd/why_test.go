package main

// Row 138 M2 (w-worldd-developer-cli §5 AC2.1–AC2.4) and the D-WORLD-60
// provenance verb, on a real in-process daemon whose log holds a REST
// genesis followed by coordinator-shaped invocation entries committed over
// POST /v1/commit. The shapes mirror host/coordinator/plan.go; the arms in
// cliwalk_e2e_test.go repeat the walk on entries the real coordinator wrote,
// so a drift between this mirror and the coordinator reds there.

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

type walkRig struct {
	t     *testing.T
	url   string
	token string
	tip   cliWorld
	next  int64
	invs  map[int64]synthInv
}

// synthInv is what one synthetic invocation committed.
type synthInv struct {
	id, episode                                 string
	record, input, output, plan, effect, result string
	outputText                                  string
}

// newWalkRig serves a fresh store from daemon.New (no interpreter: the
// coordinator is absent, the read surface and POST /v1/commit are whole),
// commits a REST genesis, and returns the rig.
func newWalkRig(t *testing.T) *walkRig {
	t.Helper()
	db := filepath.Join(t.TempDir(), "world.db")
	token := testCLISession(t, db)
	d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: db, BindHost: daemon.DefaultBindHost, ErrorLog: io.Discard})
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(d.Handler())
	t.Cleanup(func() {
		srv.Close()
		_ = d.Close()
	})
	r := &walkRig{t: t, url: srv.URL, token: token, invs: map[int64]synthInv{}}
	genesis := makeCLICommit(cliWorld{}, 0, "genesis")
	r.post(genesis)
	r.tip, r.next = genesis.NextWorld, 1
	return r
}

func (r *walkRig) post(c cliCommit) {
	r.t.Helper()
	body, err := json.Marshal(c)
	if err != nil {
		r.t.Fatal(err)
	}
	req, _ := http.NewRequest(http.MethodPost, r.url+"/v1/commit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		r.t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		r.t.Fatalf("commit %d: %d %s", c.Entry.Header.EntryIndex, resp.StatusCode, b)
	}
}

func synthObject(semantic, provenance string, payload []byte) cliObject {
	return cliObject{Hash: hashref.SumSHA256(payload).String(), InterfaceHash: hashref.SumSHA256([]byte(semantic)).String(),
		SemanticID: semantic, Provenance: provenance, Payload: payload}
}

func mustJSONT(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// invoke commits one coordinator-shaped effectful invocation (record v2, one
// allowed Workspace.Read effect) on the tip, exactly as planEffectInvocation
// lays it out, and returns what it committed.
func (r *walkRig) invoke(episode, skill, task, content string) synthInv {
	t := r.t
	t.Helper()
	in := synthObject(coordinator.InputV1, "coordinator:a2a", mustJSONT(t, map[string]string{"path": "data.txt"}))
	plan := synthObject(coordinator.EffectPlanV1, "coordinator:a2a",
		[]byte(`{"effects":[{"cost":1,"effect":"Workspace.Read","id":"e1","payload":{"op":"read","path":"data.txt"},"scope":"worktree"}],"finish":true,"plan":"world/effect-plan/v1","result":null}`))
	reqObj := synthObject(broker.EffectRequestV1, "host/broker", []byte(`{"op":"read","path":"data.txt","task":"`+task+`"}`))
	resObj := synthObject(broker.EffectResultV1, "host/broker", mustJSONT(t, map[string]string{"content": content}))
	rec := broker.EffectRecord{Effect: "Workspace.Read", Scope: "worktree", Cost: 1, BudgetBefore: 20, BudgetAfter: 19, Allowed: true,
		RequestRef: hashref.MustParse(reqObj.Hash), ResultRef: hashref.MustParse(resObj.Hash)}
	recObj := synthObject(broker.EffectRecordV1, "host/broker", broker.EncodeRecord(rec))
	outputText := string(mustJSONT(t, map[string]any{"ok": true, "content": content,
		"world": map[string]any{"plan": plan.Hash, "effects": []map[string]any{{"id": "e1", "status": "ok", "record": recObj.Hash}}}}))
	out := synthObject(coordinator.OutputV1, "coordinator:a2a", []byte(outputText))
	id := coordinator.InvocationID(episode, task)
	fn := hashref.SumSHA256([]byte("fn-" + skill)).String()
	interp := hashref.SumSHA256([]byte("interpreter")).String()
	record := synthObject(coordinator.RecordV2, "coordinator:a2a", mustJSONT(t, struct {
		InvocationID   string   `json:"invocationId"`
		EpisodeID      string   `json:"episodeId"`
		SkillID        string   `json:"skillId"`
		TransitionFn   string   `json:"transitionFn"`
		Interpreter    string   `json:"interpreter"`
		SemanticsEpoch int64    `json:"semanticsEpoch"`
		Input          string   `json:"input"`
		Output         string   `json:"output"`
		Plan           string   `json:"plan"`
		Effects        []string `json:"effects"`
	}{id, episode, skill, fn, interp, 1, in.Hash, out.Hash, plan.Hash, []string{recObj.Hash}}))
	hdr := cliHeader{EntryIndex: r.next, SemanticsEpoch: 1, TransitionFn: fn, Interpreter: interp, PrevEntryHash: r.tip.LogHead, WrittenBy: "coordinator:a2a"}
	entryHash := hashref.SumSHA256(mustJSONT(t, struct {
		EntryIndex     int64  `json:"entryIndex"`
		SemanticsEpoch int64  `json:"semanticsEpoch"`
		TransitionFn   string `json:"transitionFn"`
		Interpreter    string `json:"interpreter"`
		PrevEntryHash  string `json:"prevEntryHash"`
		WrittenBy      string `json:"writtenBy"`
		TransitionRef  string `json:"transitionRef"`
	}{hdr.EntryIndex, 1, fn, interp, hdr.PrevEntryHash, hdr.WrittenBy, record.Hash})).String()
	next := cliWorld{Revision: r.next, StateRoot: out.Hash, LogHead: entryHash}
	next.Ref = hashref.SumSHA256(mustJSONT(t, struct {
		Revision  int64  `json:"revision"`
		StateRoot string `json:"stateRoot"`
		LogHead   string `json:"logHead"`
	}{next.Revision, next.StateRoot, next.LogHead})).String()
	r.post(cliCommit{ObservedHead: r.tip.Ref, Objects: []cliObject{in, out, plan, reqObj, resObj, recObj, record}, NextWorld: next,
		Entry: cliEntry{Header: hdr, EntryHash: entryHash, TransitionRef: record.Hash}})
	inv := synthInv{id: id, episode: episode, record: record.Hash, input: in.Hash, output: out.Hash, plan: plan.Hash,
		effect: recObj.Hash, result: resObj.Hash, outputText: outputText}
	r.invs[r.next] = inv
	r.tip = next
	r.next++
	return inv
}

// whyRun runs why with stdin.
func whyRun(t *testing.T, addr string, stdin string, args ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runWhy(addr, args, strings.NewReader(stdin), &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

var linkNames = []string{"world", "entry", "record", "input", "plan", "effect", "output"}

func assertChainOK(t *testing.T, label string, code int, stdout, stderr string, entry int64) {
	t.Helper()
	if code != exitOK || !strings.HasPrefix(stdout, fmt.Sprintf("why: entry %d ", entry)) || strings.Contains(stdout, "✗") ||
		!strings.Contains(stdout, "all 7 link(s) verified") {
		t.Fatalf("%s: exit %d, stdout:\n%s\nstderr: %s", label, code, stdout, stderr)
	}
	for _, name := range linkNames {
		if !strings.Contains(stdout, "✓ "+name) {
			t.Fatalf("%s: link %s is not ✓:\n%s", label, name, stdout)
		}
	}
}

// TestWhyTargetFormsResolveToOneEntry is AC2.1: the five target forms — an
// entry index, head, a sha256 object (several kinds), an a2a: invocation id
// and a --json-out result on stdin (and via --result) — resolve to the same
// entry, with every link ✓.
func TestWhyTargetFormsResolveToOneEntry(t *testing.T) {
	r := newWalkRig(t)
	inv := r.invoke("ep1", "ailang-read", "task-1", "hello from ep1\n")
	r.invoke("ep1", "ailang-read", "task-2", "second\n")
	target := r.invoke("ep2", "ailang-read", "task-3", "third\n") // head
	_ = inv
	resultFile := filepath.Join(t.TempDir(), "result.json")
	writeFileT(t, resultFile, target.outputText+"\n")
	for _, c := range []struct {
		label, stdin string
		args         []string
	}{
		{"index", "", []string{"3"}},
		{"head", "", []string{"head"}},
		{"sha256 output", "", []string{target.output}},
		{"sha256 record", "", []string{target.record}},
		{"sha256 input", "", []string{target.input}},
		{"sha256 plan", "", []string{target.plan}},
		{"sha256 effect record", "", []string{target.effect}},
		{"sha256 effect result", "", []string{target.result}},
		{"sha256 world ref", "", []string{r.tip.Ref}},
		{"a2a invocation id", "", []string{target.id}},
		{"- (call --json-out on stdin)", target.outputText + "\n", []string{"-"}},
		{"--result file", "", []string{"--result", resultFile}},
	} {
		code, stdout, stderr := whyRun(t, r.url, c.stdin, c.args...)
		assertChainOK(t, c.label, code, stdout, stderr, 3)
		if c.label == "- (call --json-out on stdin)" && !strings.Contains(stdout, "your result is these exact bytes") {
			t.Fatalf("why - did not confirm the result bytes:\n%s", stdout)
		}
	}
	// An older entry by index and by its a2a id: the same entry 1.
	for _, arg := range []string{"1", inv.id, inv.output} {
		code, stdout, stderr := whyRun(t, r.url, "", arg)
		assertChainOK(t, arg, code, stdout, stderr, 1)
	}
	// --json prints the chain.
	code, stdout, _ := whyRun(t, r.url, "", "head", "--json")
	var ch whyChain
	if code != exitOK || json.Unmarshal([]byte(stdout), &ch) != nil || !ch.OK || ch.Entry != 3 || len(ch.Links) != 7 {
		t.Fatalf("--json exit %d: %s", code, stdout)
	}
	// The REST genesis is not a coordinator entry: rendered, and said why.
	code, stdout, _ = whyRun(t, r.url, "", "0")
	if code != exitOK || !strings.Contains(stdout, "is not a coordinator invocation") {
		t.Fatalf("why 0: exit %d\n%s", code, stdout)
	}
}

func writeFileT(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

// tamperProxy serves the daemon, but rewrites the payload of one object.
func tamperProxy(t *testing.T, daemonURL, ref string, payload []byte) string {
	t.Helper()
	target, _ := url.Parse(daemonURL)
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ModifyResponse = func(resp *http.Response) error {
		if resp.Request.URL.Path != "/v1/objects/"+ref || resp.Request.URL.Query().Get("payload") != "true" {
			return nil
		}
		var obj map[string]any
		b, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		_ = json.Unmarshal(b, &obj)
		obj["payload"] = base64.StdEncoding.EncodeToString(payload)
		nb, _ := json.Marshal(obj)
		resp.Body = io.NopCloser(bytes.NewReader(nb))
		resp.ContentLength = int64(len(nb))
		resp.Header.Del("Content-Length")
		return nil
	}
	srv := httptest.NewServer(proxy)
	t.Cleanup(srv.Close)
	return srv.URL
}

// TestWhyBrokenLinkIsIntegrityRefusal is AC2.2 (kills MUT-WHY-NOVERIFY): a
// read surface that serves a crafted output whose bytes do not hash to the
// record's output ref is a ✗ on the output link and exit 3 — for every
// target form that reaches the entry.
func TestWhyBrokenLinkIsIntegrityRefusal(t *testing.T) {
	r := newWalkRig(t)
	inv := r.invoke("ep1", "ailang-read", "task-1", "hello from ep1\n")
	crafted := strings.Replace(inv.outputText, "hello from ep1", "hello from EVE", 1)
	addr := tamperProxy(t, r.url, inv.output, []byte(crafted))
	for _, args := range [][]string{{"1"}, {"head"}, {inv.id}} {
		code, stdout, stderr := whyRun(t, addr, "", args...)
		if code != exitIntegrity || !strings.Contains(stdout, "✗ output") || !strings.Contains(stdout, "BROKEN") {
			t.Fatalf("why %v over a tampered output: exit %d\n%s%s", args, code, stdout, stderr)
		}
		if strings.Contains(stdout, "✗ world") || strings.Contains(stdout, "✗ record") {
			t.Fatalf("only the output link is broken:\n%s", stdout)
		}
	}
	// A tampered effect result: ✗ on the effect link.
	addr = tamperProxy(t, r.url, inv.result, []byte(`{"content":"forged"}`))
	if code, stdout, _ := whyRun(t, addr, "", "1"); code != exitIntegrity || !strings.Contains(stdout, "✗ effect") {
		t.Fatalf("tampered effect result: exit %d\n%s", code, stdout)
	}
	// Control: the untampered surface verifies.
	code, stdout, stderr := whyRun(t, r.url, "", "1")
	assertChainOK(t, "control", code, stdout, stderr, 1)
	// A result that is not the committed output bytes (one byte changed) is
	// found by the world.plan fallback and refused.
	if code, stdout, _ := whyRun(t, r.url, crafted, "-"); code != exitIntegrity || !strings.Contains(stdout, "world.plan (fallback") ||
		!strings.Contains(stdout, "is not the committed output bytes") {
		t.Fatalf("why - of altered bytes: exit %d\n%s", code, stdout)
	}
}

// TestWhyScanBound is AC2.3: a target older than the scan window is not
// found, exit 1; the same target with the default window is found.
func TestWhyScanBound(t *testing.T) {
	r := newWalkRig(t)
	old := r.invoke("ep1", "ailang-read", "task-1", "one\n")
	r.invoke("ep1", "ailang-read", "task-2", "two\n")
	code, stdout, stderr := whyRun(t, r.url, "", old.id, "--scan", "1")
	if code != exitUsage || stdout != "" || !strings.Contains(stderr, "not found in the last 1 entries") {
		t.Fatalf("--scan 1: exit %d stdout %q stderr %q", code, stdout, stderr)
	}
	code, stdout, stderr = whyRun(t, r.url, "", old.id)
	assertChainOK(t, "default scan", code, stdout, stderr, 1)
	if code, _, stderr := whyRun(t, r.url, "", "a2a:ep9:nope"); code != exitUsage || !strings.Contains(stderr, "not found in the last 500 entries") {
		t.Fatalf("absent id: exit %d %q", code, stderr)
	}
	if code, _, _ := whyRun(t, r.url, "", "x", "--scan", "5001"); code != exitUsage {
		t.Fatalf("--scan 5001 accepted")
	}
}

// TestLogTailFollowPrintsEachCommitOnce is AC2.4 (kills MUT-FOLLOW-SKIP):
// with --follow running from the head, three commits made one at a time are
// each printed exactly once, in order.
func TestLogTailFollowPrintsEachCommitOnce(t *testing.T) {
	r := newWalkRig(t)
	r.invoke("ep1", "ailang-read", "task-0", "zero\n")
	ctx, cancel := context.WithCancel(boundedTestContext(t))
	defer cancel()
	var stdout, stderr syncBuffer
	done := make(chan int, 1)
	go func() {
		done <- logTail(ctx, newClient(r.url), tailOpts{from: -1, follow: true, interval: 20 * time.Millisecond}, &stdout, &stderr)
	}()
	waitFor := func(want string) {
		t.Helper()
		deadline := time.Now().Add(10 * time.Second)
		for !strings.Contains(stdout.String(), want) {
			if time.Now().After(deadline) {
				t.Fatalf("log tail never printed %q; stdout:\n%s\nstderr: %s", want, stdout.String(), stderr.String())
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	waitFor("#1 ")
	for i := 2; i <= 4; i++ {
		time.Sleep(60 * time.Millisecond) // let at least one empty poll pass
		r.invoke("ep1", "ailang-read", fmt.Sprintf("task-%d", i), "more\n")
		waitFor(fmt.Sprintf("#%d ", i))
	}
	time.Sleep(100 * time.Millisecond)
	cancel()
	if code := <-done; code != exitOK {
		t.Fatalf("log tail exit %d", code)
	}
	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	want := []string{"#0 ", "#1 ", "#2 ", "#3 ", "#4 "}
	if len(lines) != len(want) {
		t.Fatalf("log tail printed %d lines, want %d:\n%s", len(lines), len(want), stdout.String())
	}
	for i, l := range lines {
		if !strings.HasPrefix(l, want[i]) {
			t.Fatalf("line %d = %q, want prefix %q", i, l, want[i])
		}
	}
	if !strings.Contains(lines[2], "coordinator:a2a ep1 ailang-read [Workspace.Read ok]") {
		t.Fatalf("coordinator line = %q", lines[2])
	}
	// Without --follow the verb prints and exits; --raw prints entry JSON.
	code, out, _ := runCLI(t, r.url, "log", "tail", "--from", "3", "--raw")
	var e wireEntry
	if code != exitOK || len(strings.Split(strings.TrimSpace(out), "\n")) != 2 || json.Unmarshal([]byte(strings.Split(out, "\n")[0]), &e) != nil || e.Header.EntryIndex != 3 {
		t.Fatalf("log tail --from 3 --raw: exit %d %q", code, out)
	}
}

// TestProvenanceTrailer is D-WORLD-60's verb: the trailer names the store,
// the episode and its first..last coordinator entries since --since.
func TestProvenanceTrailer(t *testing.T) {
	r := newWalkRig(t)
	r.invoke("ep1", "ailang-read", "t1", "a\n") // 1
	r.invoke("ep2", "ailang-read", "t2", "b\n") // 2
	r.invoke("ep1", "ailang-read", "t3", "c\n") // 3
	r.invoke("ep1", "ailang-read", "t4", "d\n") // 4
	var h healthInfo
	_, body, _ := newClient(r.url).get("/v1/health")
	if json.Unmarshal([]byte(body), &h) != nil || h.DBPath == "" {
		t.Fatalf("health %s", body)
	}
	code, stdout, stderr := runCLI(t, r.url, "provenance", "--episode", "ep1")
	if want := "World-Provenance: store=" + h.DBPath + " episode=ep1 entries=1-4\n"; code != exitOK || stdout != want {
		t.Fatalf("provenance --episode ep1: exit %d %q %q, want %q", code, stdout, stderr, want)
	}
	code, stdout, _ = runCLI(t, r.url, "provenance", "--since", "2", "--episode", "ep1")
	if code != exitOK || !strings.HasSuffix(stdout, " episode=ep1 entries=3-4\n") {
		t.Fatalf("--since 2: exit %d %q", code, stdout)
	}
	code, stdout, _ = runCLI(t, r.url, "provenance")
	if code != exitOK || stdout != "World-Provenance: store="+h.DBPath+" episode=ep1 entries=1-4\nWorld-Provenance: store="+h.DBPath+" episode=ep2 entries=2-2\n" {
		t.Fatalf("all episodes: exit %d %q", code, stdout)
	}
	if code, _, stderr := runCLI(t, r.url, "provenance", "--episode", "ep9"); code != exitUsage || !strings.Contains(stderr, "no coordinator entries for episode ep9") {
		t.Fatalf("absent episode: exit %d %q", code, stderr)
	}
	if code, _, _ := runCLI(t, r.url, "provenance", "--since", "9"); code != exitUsage {
		t.Fatal("--since past head accepted")
	}
	// Each trailer entry walks back with why.
	code, stdout, stderr = whyRun(t, r.url, "", "4")
	assertChainOK(t, "trailer entry", code, stdout, stderr, 4)
}
