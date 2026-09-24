package daemon

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// value_demonstration_test.go — w-prove-1-0-phase-a (row 92, Phase A) durability test.
//
// The COMMITTED fixtures in design_docs/verification/w-1-0-value-demonstration/commits/
// carry the plan's harvest (one chained commit per taken question, Q1 genesis -> ...).
// This test replays them through the REAL daemon stack (httptest + session-gated POST)
// and then performs the provenance walk IN-PROCESS (scan log range, follow transitionRefs,
// read incident + evidence objects as a walker would), asserting that the walk, and only
// the walk, yields the recorded answers and evidence — i.e. that the checked-in record is
// load-bearing, not decorative. This is the M3 durability guard that survives while M4 runs
// the live timed walk against a scratch DB.
//
// It is also the mutation-battery's sole-killer target:
//   - MU-1 / MU-3 -> TestValueDemo_AnswersMatchRecord
//   - MU-2       -> TestValueDemo_EvidenceChainComplete

// ExpectedAnswers mirrors the generator's FACTS: each taken question's recorded answer.
// Copied verbatim from gen_fixtures.py / the plan's harvest table (never paraphrased).
var expectedAnswers = map[string]struct {
	answer   string
	excerpts []string
}{
	"iter171-index-row": {
		answer: "The pinned v0.30.0 binary lacks the mission command group: `ailang.v0.30.0 mission rotate-log world --keep 31` prints `Error: unknown command 'mission'`; the index hole is visible in git history of world-mission-index.md (the row exists today because it was added manually after the iter-172 rule).",
		excerpts: []string{
			"The pinned v0.30.0 binary has no `ailang mission` command at all: `~/.pinned-ailang/ailang mission rotate-log world --keep 31` prints **`Error: unknown command 'mission'`**.",
			"iteration **171** had a full log entry and **no index row** (`grep -c '^| 171 |'` -> **0**, control `170` -> **1**).",
			"the failure is SILENT in the direction that matters: Gate 2's first instruction is to grep the index before picking, so a missing row does not produce an error.",
		},
	},
	"iter154-unfenced-pi": {
		answer: "The fleet's scripts/mission_pi_run.sh invokes `pi --mode json --no-session --model \"$MODEL\" < \"$DIRECTIVE\"` with no -e flag and no PI_FENCE_ROOT anywhere in the file, so the shared skill's two sandbox extensions are never wired; issue sunholo-data/ailang#1043 is OPEN.",
		excerpts: []string{
			"measured first-party at `scripts/mission_pi_run.sh:155`: the invocation is `pi --mode json --no-session --model \"$MODEL\" < \"$DIRECTIVE\"`, with **no `-e` flag at all** and no `PI_FENCE_ROOT`.",
			"The recipe's invocation block is explicit — `-e \"$REPO/tools/pi-extensions/sandbox/index.ts\" -e \"$REPO/tools/pi-extensions/worktree-fence.ts\"`.",
			"the code half is filed upstream as [`ailang#1043`](https://github.com/sunholo-data/ailang/issues/1043).",
		},
	},
	"iter155-faithfulness-proof": {
		answer: "The prescribed manifest covers the **final tree**, identical whether the split is correct or collapsed; `git add design_docs host` staged disk state, MS1 swallowed MS3's files, MS3 committed nothing, and `shasum -c` returned OK on every file.",
		excerpts: []string{
			"**That manifest is over the FINAL TREE.** Bisectability is a property of the SPLIT, and the final tree is identical whether the split is correct or whether commit 1 swallowed everything and commits 2-3 are empty.",
			"**My first commit reconstruction was wrong and I rebuilt it.** `git add design_docs host` staged whatever was on disk rather than the snapshot's named files, so MS1's commit swallowed MS3's two doc files.",
			"zsh does not word-split unquoted expansions, so all three paths arrived as ONE argument, `cp` failed, and the boundary gates still printed `vet=0 verifygate=0`.",
		},
	},
	"iter181-ci-bench-401": {
		answer: "a036062 is an ancestor of base and its diff adds `req.Header.Set(\"Authorization\", auth)` to bench_test.go — the bench predated the session gate and POSTed to the now-protected /v1/commit unauthenticated, 401 (the middleware working as designed).",
		excerpts: []string{
			"`a036062` is an ancestor of base and its diff adds `req.Header.Set(\"Authorization\", auth)` to `bench_test.go` — the bench predated the session gate and 401'd.",
			"remote CI on PR #141 then caught what every local gate missed — `BenchmarkRESTCommit` POSTing to the now-protected `/v1/commit` unauthenticated (401: the middleware working as designed; the bench predates the sprint and was not in the plan's local gate list).",
		},
	},
}

// fixtureDir resolves the checked-in commit fixtures from the package's test working
// directory (host/daemon), up two levels to the repo root, into the verification tree.
func fixtureDir() string {
	return filepath.Join("..", "..", "design_docs",
		"verification", "w-1-0-value-demonstration", "commits")
}

type incidentWalk struct {
	Question    string   `json:"question"`
	Answer      string   `json:"answer"`
	DiagnosedAt string   `json:"diagnosed_at"`
	Sources     []string `json:"sources"`
}

type evidenceWalk struct {
	Kind    string `json:"kind"`
	Ref     string `json:"ref"`
	Check   string `json:"check"`
	Excerpt string `json:"excerpt"`
}

// valueDemoSetup replays every committed fixture through a fresh daemon stack and returns
// the running test server (closed via srv we defer), its daemon, the auth header, and the
// ordered fixture paths. Called by both named tests so the replay is not duplicated.
func valueDemoSetup(t *testing.T) (*httptest.Server, *Daemon, string, []string) {
	t.Helper()
	d := newHandlerDaemon(t)
	auth := authHeader(t, d)
	srv := httptest.NewServer(d.Handler())

	fixtures, err := filepath.Glob(filepath.Join(fixtureDir(), "q*.json"))
	if err != nil {
		srv.Close()
		t.Fatalf("glob fixtures: %v", err)
	}
	if len(fixtures) < 3 {
		srv.Close()
		t.Fatalf("expected >=3 commit fixtures, found %d", len(fixtures))
	}
	sort.Strings(fixtures)

	for _, f := range fixtures {
		raw, err := os.ReadFile(f)
		if err != nil {
			srv.Close()
			t.Fatalf("read fixture %s: %v", f, err)
		}
		req, err := http.NewRequest(http.MethodPost, srv.URL+"/v1/commit", bytes.NewReader(raw))
		if err != nil {
			srv.Close()
			t.Fatalf("build commit request: %v", err)
		}
		req.Header.Set("Authorization", auth) // POST /v1/commit is session-gated
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			srv.Close()
			t.Fatalf("POST fixture %s: %v", f, err)
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			srv.Close()
			t.Fatalf("POST fixture %s: status=%d body=%s", f, resp.StatusCode, body)
		}
	}
	return srv, d, auth, fixtures
}

// walkIncidents performs the provenance walk a walker would: scan log range from 0, follow
// each transitionRef's object metadata, and collect incidents whose semanticId is
// world/mission/incident/<slug>. Returns slug -> incident payload + slug -> []evidence.
func walkIncidents(t *testing.T, srv *httptest.Server) (map[string]incidentWalk, map[string][]evidenceWalk) {
	t.Helper()
	res, err := http.Get(srv.URL + "/v1/log?from=0&limit=100") // GET routes are open
	if err != nil {
		t.Fatalf("walk log range: %v", err)
	}
	var page logRangeResponse
	if err := json.NewDecoder(res.Body).Decode(&page); err != nil {
		_ = res.Body.Close()
		t.Fatalf("decode log range: %v", err)
	}
	_ = res.Body.Close()
	if len(page.Items) == 0 {
		t.Fatalf("log range returned 0 entries after committing fixtures")
	}

	incidents := map[string]incidentWalk{}
	evidence := map[string][]evidenceWalk{}
	for _, item := range page.Items {
		incObject := getObject(t, srv, item.TransitionRef, true)
		var meta objectResponse
		_ = json.Unmarshal([]byte(incObject), &meta)
		if !strings.HasPrefix(meta.SemanticID, "world/mission/incident/") {
			continue // an evidence or non-mission object; skip
		}
		slug := strings.TrimPrefix(meta.SemanticID, "world/mission/incident/")
		inc := decodeIncident(t, meta, incObject)
		incidents[slug] = inc
		var evs []evidenceWalk
		for _, src := range inc.Sources {
			evObject := getObject(t, srv, src, true)
			evs = append(evs, decodeEvidence(t, evObject))
		}
		evidence[slug] = evs
	}
	return incidents, evidence
}

// getObject fetches /v1/objects/{ref}[?payload=true] and returns the raw JSON envelope.
func getObject(t *testing.T, srv *httptest.Server, ref string, payload bool) string {
	t.Helper()
	path := srv.URL + "/v1/objects/" + ref
	if payload {
		path += "?payload=true"
	}
	res, err := http.Get(path)
	if err != nil {
		t.Fatalf("walk object %s: %v", ref, err)
	}
	body, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("walk object %s: status=%d body=%s", ref, res.StatusCode, body)
	}
	return string(body)
}

func decodeIncident(t *testing.T, meta objectResponse, raw string) incidentWalk {
	t.Helper()
	if meta.Payload == nil {
		t.Fatalf("incident object %s returned no payload (payload=true expected)", meta.SemanticID)
	}
	var inc incidentWalk
	if err := json.Unmarshal(*meta.Payload, &inc); err != nil {
		t.Fatalf("decode incident %s payload: %v (raw=%s)", meta.SemanticID, err, raw)
	}
	if len(inc.Sources) == 0 {
		t.Fatalf("incident %s has empty sources[]", meta.SemanticID)
	}
	return inc
}

func decodeEvidence(t *testing.T, raw string) evidenceWalk {
	t.Helper()
	var obj objectResponse
	if err := json.Unmarshal([]byte(raw), &obj); err != nil {
		t.Fatalf("decode evidence envelope: %v", err)
	}
	if obj.Payload == nil {
		t.Fatalf("evidence object %s returned no payload", obj.SemanticID)
	}
	var ev evidenceWalk
	if err := json.Unmarshal(*obj.Payload, &ev); err != nil {
		t.Fatalf("decode evidence payload: %v", err)
	}
	return ev
}

// TestValueDemo_AnswersMatchRecord replays the fixtures and asserts the walk's incident
// answer byte-equals the recorded answer for every taken question. Sole killer of MU-1
// (wrong expected answer) and MU-3 (corrupt fixture payload).
func TestValueDemo_AnswersMatchRecord(t *testing.T) {
	srv, _, _, _ := valueDemoSetup(t)
	defer srv.Close()
	incidents, _ := walkIncidents(t, srv)

	for slug, want := range expectedAnswers {
		inc, ok := incidents[slug]
		if !ok {
			t.Errorf("walk located no incident object for slug %q", slug)
			continue
		}
		if inc.Answer != want.answer {
			t.Errorf("answer mismatch for %q:\n walk: %s\n want: %s", slug, inc.Answer, want.answer)
		}
	}
}

// TestValueDemo_EvidenceChainComplete asserts every sources[] hash in every incident
// resolves to an evidence object whose excerpt byte-equals the recorded excerpt. This is
// the end-to-end chain: a walk that follows incident.answer but refuses to resolve
// evidence (MU-2 collapsing the fetch loop) must fail here.
func TestValueDemo_EvidenceChainComplete(t *testing.T) {
	srv, _, _, _ := valueDemoSetup(t)
	defer srv.Close()
	_, evidence := walkIncidents(t, srv)

	for slug, want := range expectedAnswers {
		evs, ok := evidence[slug]
		if !ok {
			t.Errorf("walk resolved no evidence objects for slug %q", slug)
			continue
		}
		if len(evs) != len(want.excerpts) {
			t.Errorf("evidence count for %q = %d, want %d", slug, len(evs), len(want.excerpts))
		}
		for i, ev := range evs {
			if i >= len(want.excerpts) {
				break
			}
			if ev.Excerpt != want.excerpts[i] {
				t.Errorf("evidence[%d].excerpt mismatch for %q:\n  walk: %s\n  want: %s",
					i, slug, ev.Excerpt, want.excerpts[i])
			}
			if ev.Kind == "" || ev.Ref == "" || ev.Check == "" {
				t.Errorf("evidence[%d] for %q missing kind/ref/check: %+v", i, slug, ev)
			}
		}
	}
}
