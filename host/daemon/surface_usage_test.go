package daemon

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

const codexTurn = "{\"type\":\"turn.started\"}\n{\"type\":\"item.completed\",\"item\":{\"type\":\"agent_message\",\"text\":\"private fixture content\"}}\n{\"type\":\"turn.completed\",\"usage\":{\"input_tokens\":42}}\n"
const codexThread = "{\"type\":\"thread.started\",\"thread_id\":\"fixture\"}\n"
const claudeStart = "{\"type\":\"system\",\"subtype\":\"init\"}\n{\"type\":\"assistant\",\"message\":{\"private\":true}}\n"
const claudeResult = "{\"type\":\"result\",\"subtype\":\"success\",\"is_error\":false,\"num_turns\":3,\"total_cost_usd\":0.123}\n"
const launchSeparator = "{\"world_launch\": 2}\n"

func codexUsageTurn(usage string) string {
	return "{\"type\":\"turn.started\"}\n{\"type\":\"turn.completed\",\"usage\":" + usage + "}\n"
}

func claudeUsageResult(usage string) string {
	return strings.TrimSpace(claudeResult[:len(claudeResult)-2]) + ",\"usage\":" + usage + "}\n"
}

func TestSurfaceRunUsage(t *testing.T) {
	for _, tc := range []struct {
		name, agent string
		launches    []string
		wantTurns   int
		wantCost    float64
		wantTokens  string
		edit        func(*testing.T, *store.EpisodeStart, string)
	}{
		{name: "Codex audit counters preserve subsets and zero", agent: "codex", wantTurns: 1, wantCost: -1,
			launches:   []string{codexThread + codexUsageTurn(`{"input_tokens":258254,"cached_input_tokens":198784,"cache_write_input_tokens":0,"output_tokens":1883,"reasoning_output_tokens":24}`)},
			wantTokens: `{"input":258254,"output":1883,"cacheRead":198784,"cacheWrite":0,"reasoning":24}`},
		{name: "Codex uses final cumulative token evidence then sums fresh launches", agent: "codex", wantTurns: 3, wantCost: -1,
			launches: []string{
				codexThread + codexUsageTurn(`{"input_tokens":10,"output_tokens":2,"reasoning_output_tokens":1}`) +
					codexUsageTurn(`{"input_tokens":14,"output_tokens":3}`),
				codexThread + codexUsageTurn(`{"input_tokens":5,"output_tokens":0,"reasoning_output_tokens":0}`),
			}, wantTokens: `{"input":19,"output":3}`},
		{name: "Codex omits invalid subsets without losing input or output", agent: "codex", wantTurns: 1, wantCost: -1,
			launches:   []string{codexThread + codexUsageTurn(`{"input_tokens":14,"cached_input_tokens":15,"cache_write_input_tokens":-1,"output_tokens":3,"reasoning_output_tokens":4}`)},
			wantTokens: `{"input":14,"output":3}`},
		{name: "Claude input is disjoint ordinary tokens plus cache reads and writes", agent: "claude", wantTurns: 3, wantCost: 0.123,
			launches:   []string{claudeStart + claudeUsageResult(`{"input_tokens":10,"output_tokens":15,"cache_read_input_tokens":9,"cache_creation_input_tokens":2}`)},
			wantTokens: `{"input":21,"output":15,"cacheRead":9,"cacheWrite":2}`},
		{name: "Agy headless stream-json totals input, output and cacheRead", agent: "agy", wantTurns: 3, wantCost: 0.123,
			launches:   []string{claudeStart + claudeUsageResult(`{"input_tokens":10,"output_tokens":15,"cache_read_tokens":9}`)},
			wantTokens: `{"input":19,"output":15,"cacheRead":9}`},
		{name: "Claude invalid cache count omits input total and preserves independent metrics", agent: "claude", wantTurns: 3, wantCost: 0.123,
			launches:   []string{claudeStart + claudeUsageResult(`{"input_tokens":10,"output_tokens":0,"cache_read_input_tokens":"bogus","cache_creation_input_tokens":0}`)},
			wantTokens: `{"output":0,"cacheWrite":0}`},
		{name: "Claude malformed usage preserves turns and cost", agent: "claude", wantTurns: 3, wantCost: 0.123,
			launches: []string{claudeStart + claudeUsageResult(`[42]`)}, wantTokens: "null"},
		{name: "Claude aggregate overflow omits only affected token field", agent: "claude", wantTurns: 6, wantCost: 0.246,
			launches: []string{
				claudeStart + claudeUsageResult(`{"input_tokens":9007199254740991,"output_tokens":3,"cache_read_input_tokens":0,"cache_creation_input_tokens":0}`),
				claudeStart + claudeUsageResult(`{"input_tokens":10,"output_tokens":5,"cache_read_input_tokens":0,"cache_creation_input_tokens":0}`),
			}, wantTokens: `{"output":8,"cacheRead":0,"cacheWrite":0}`},
		{name: "Codex completed turns", agent: "codex", launches: []string{codexThread + codexTurn + codexTurn}, wantTurns: 2, wantCost: -1},
		{name: "Claude final result", agent: "claude", launches: []string{claudeStart + claudeResult}, wantTurns: 3, wantCost: 0.123},
		{name: "Claude totals across confirmed launches", agent: "claude", launches: []string{claudeStart + claudeResult, claudeStart + claudeResult}, wantTurns: 6, wantCost: 0.246},
		{name: "Codex totals across confirmed launches", agent: "codex", launches: []string{codexThread + codexTurn, codexThread + codexTurn + codexTurn}, wantTurns: 3, wantCost: -1},
		{name: "Claude omits missing cost", agent: "claude", launches: []string{claudeStart + strings.Replace(claudeResult, ",\"total_cost_usd\":0.123", "", 1)}, wantTurns: 3, wantCost: -1},
		{name: "Claude reported zero", agent: "claude", launches: []string{claudeStart + strings.NewReplacer("3,", "0,", "0.123", "0").Replace(claudeResult)}, wantTurns: 0, wantCost: 0},
		{name: "Claude partial result", agent: "claude", launches: []string{claudeStart}, wantTurns: -1, wantCost: -1},
		{name: "Claude duplicate result", agent: "claude", launches: []string{claudeStart + claudeResult + claudeResult}, wantTurns: -1, wantCost: -1},
		{name: "Codex partial final turn", agent: "codex", launches: []string{codexThread + codexTurn + "{\"type\":\"turn.started\"}\n"}, wantTurns: -1, wantCost: -1},
		{name: "Codex failed turn", agent: "codex", launches: []string{codexThread + strings.Replace(codexTurn, "turn.completed", "turn.failed", 1)}, wantTurns: -1, wantCost: -1},
		{name: "mixed formats", agent: "codex", launches: []string{codexThread + codexTurn, claudeStart + claudeResult}, wantTurns: -1, wantCost: -1},
		{name: "partial JSON line", agent: "codex", launches: []string{codexThread + codexTurn + "{\"type\":"}, wantTurns: -1, wantCost: -1},
		{name: "bytes differ from recorded digest", agent: "codex", launches: []string{codexThread + codexTurn}, wantTurns: -1, wantCost: -1,
			edit: func(t *testing.T, _ *store.EpisodeStart, path string) {
				if err := os.WriteFile(path, []byte(codexThread+codexTurn+codexTurn), 0o600); err != nil {
					t.Fatal(err)
				}
			}},
		{name: "incorrect earlier digest", agent: "codex", launches: []string{codexThread + codexTurn, codexThread + codexTurn}, wantTurns: -1, wantCost: -1,
			edit: func(_ *testing.T, start *store.EpisodeStart, _ string) {
				start.Launches[0].TranscriptSHA256 = start.Launches[1].TranscriptSHA256
			}},
		{name: "incorrect separator", agent: "codex", launches: []string{codexThread + codexTurn + strings.Replace(launchSeparator, "2", "3", 1) + codexThread + codexTurn}, wantTurns: -1, wantCost: -1},
		{name: "interactive launch", agent: "claude", launches: []string{claudeStart + claudeResult}, wantTurns: -1, wantCost: -1,
			edit: func(_ *testing.T, start *store.EpisodeStart, _ string) { start.Launches[0].Mode = "interactive" }},
		{name: "failed launch", agent: "claude", launches: []string{claudeStart + claudeResult}, wantTurns: -1, wantCost: -1,
			edit: func(_ *testing.T, start *store.EpisodeStart, _ string) {
				start.Launches[0].Exit, start.Launches[0].FinishReason = 1, "deadline"
			}},
		{name: "scrubbed launch", agent: "claude", launches: []string{claudeStart + claudeResult}, wantTurns: -1, wantCost: -1,
			edit: func(_ *testing.T, start *store.EpisodeStart, _ string) {
				start.Launches[0].TranscriptState = "scrubbed"
			}},
		{name: "record path is irrelevant", agent: "claude", launches: []string{claudeStart + claudeResult}, wantTurns: -1, wantCost: -1,
			edit: func(t *testing.T, start *store.EpisodeStart, path string) {
				start.Launches[0].Transcript = path + ".host"
				if err := os.Rename(path, start.Launches[0].Transcript); err != nil {
					t.Fatal(err)
				}
			}},
		{name: "symlink refused", agent: "claude", launches: []string{claudeStart + claudeResult}, wantTurns: -1, wantCost: -1,
			edit: func(t *testing.T, _ *store.EpisodeStart, path string) {
				if err := os.Rename(path, path+".real"); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(path+".real", path); err != nil {
					t.Fatal(err)
				}
			}},
		{name: "oversize refused", agent: "claude", launches: []string{claudeStart + strings.Repeat(" ", maxSurfaceUsageBytes) + claudeResult}, wantTurns: -1, wantCost: -1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			sessions := filepath.Join(dir, "sessions")
			if err := os.Mkdir(sessions, 0o700); err != nil {
				t.Fatal(err)
			}
			d := &Daemon{cfg: Config{DBPath: filepath.Join(dir, "world.db")}}
			start := &store.EpisodeStart{Episode: "wa1", Agent: tc.agent}
			var stream string
			for i, launch := range tc.launches {
				if i > 0 {
					stream += launchSeparator
				}
				stream += launch
				start.Launches = append(start.Launches, store.EpisodeLaunch{
					Mode: "headless", FinishReason: "exit", TranscriptState: "kept",
					TranscriptSHA256: hashref.SumSHA256([]byte(stream)).String(),
				})
			}
			path := filepath.Join(sessions, start.Episode+".transcript.jsonl")
			if err := os.WriteFile(path, []byte(stream), 0o600); err != nil {
				t.Fatal(err)
			}
			if tc.edit != nil {
				tc.edit(t, start, path)
			}
			turns, cost, tokens := d.surfaceRunUsage(start)
			gotTurns, gotCost := -1, -1.0
			if turns != nil {
				gotTurns = *turns
			}
			if cost != nil {
				gotCost = *cost
			}
			if gotTurns != tc.wantTurns || math.Abs(gotCost-tc.wantCost) > 1e-12 {
				t.Fatalf("turns and cost = %d, %v; want %d, %v (negative denotes unavailable)", gotTurns, gotCost, tc.wantTurns, tc.wantCost)
			}
			if tc.wantTokens != "" {
				got, err := json.Marshal(tokens)
				if err != nil || string(got) != tc.wantTokens {
					t.Fatalf("token evidence = %s (%v), want %s", got, err, tc.wantTokens)
				}
			} else if tc.wantTurns < 0 && tokens != nil {
				t.Fatalf("unverified transcript supplied token evidence: %+v", tokens)
			}
		})
	}
}

func TestSurfaceUsageMemoUsesOnlyUnchangedVerifiedFile(t *testing.T) {
	dir := t.TempDir()
	sessions := filepath.Join(dir, "sessions")
	if err := os.Mkdir(sessions, 0o700); err != nil {
		t.Fatal(err)
	}
	d := &Daemon{cfg: Config{DBPath: filepath.Join(dir, "world.db")}}
	stream := codexThread + codexTurn
	start := store.EpisodeStart{Episode: "wa1", Agent: "codex", Repo: "/repo/one", Launches: []store.EpisodeLaunch{{
		Mode: "headless", FinishReason: "exit", TranscriptState: "kept",
		TranscriptSHA256: hashref.SumSHA256([]byte(stream)).String(),
	}}}
	path := filepath.Join(sessions, "wa1.transcript.jsonl")
	if got, _, _ := d.surfaceRunUsage(&start); got != nil {
		t.Fatal("missing transcript supplied usage")
	}
	if err := os.WriteFile(path, []byte(stream), 0o600); err != nil {
		t.Fatal(err)
	}
	original, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	turns, _, tokens := d.surfaceRunUsage(&start)
	if turns == nil || *turns != 1 || tokens == nil || tokens.Input == nil || *tokens.Input != 42 {
		t.Fatal("transcript did not supply usage after it appeared")
	}
	if _, cacheable := surfaceUsageChangeTime(original); !cacheable {
		t.Skip("this platform does not cache usage")
	}
	*turns = 100
	*tokens.Input = 100
	d.surfaceUsage.parsing = true
	if got, _, cached := d.surfaceRunUsage(&start); got == nil || *got != 1 || cached == nil || cached.Input == nil || *cached.Input != 42 {
		t.Fatal("cache hit was mutable or skipped while another file parsed")
	}
	other := start
	other.Repo = "/repo/two"
	if got, _, _ := d.surfaceRunUsage(&other); got != nil {
		t.Fatal("another repository reused the cached record during a busy miss")
	}
	ineligible := start
	ineligible.Launches = append([]store.EpisodeLaunch(nil), start.Launches...)
	ineligible.Launches[0].Mode = "interactive"
	if got, _, _ := d.surfaceRunUsage(&ineligible); got != nil {
		t.Fatal("ineligible metadata reused cached usage")
	}
	d.surfaceUsage.parsing = false

	// A same-length rewrite can reset mtime but cannot reset the inode's ctime.
	tampered := strings.Replace(stream, "42", "43", 1)
	if err := os.WriteFile(path, []byte(tampered), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, original.ModTime(), original.ModTime()); err != nil {
		t.Fatal(err)
	}
	if got, _, _ := d.surfaceRunUsage(&start); got != nil {
		t.Fatal("a same-size rewrite with restored mtime returned stale usage")
	}
	if err := os.WriteFile(path, []byte(stream), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, _, _ := d.surfaceRunUsage(&start); got == nil || *got != 1 {
		t.Fatal("a rejected mismatch prevented a later verified read")
	}
}

func TestSurfaceUsageMemoBoundsEntries(t *testing.T) {
	var memo surfaceUsageMemo
	count := 1
	entry := surfaceUsageEntry{file: surfaceUsageFileVersion{cacheable: true}, turns: &count}
	for i := range maxSurfaceUsageEntries + 1 {
		memo.finish(surfaceUsageEpisode{episode: strconv.Itoa(i)}, entry)
	}
	if len(memo.entries) != maxSurfaceUsageEntries {
		t.Fatalf("retained %d entries, limit %d", len(memo.entries), maxSurfaceUsageEntries)
	}
}

func TestSurfacePermissionsSeparatesExecutedAndReservedDebits(t *testing.T) {
	cfg, _ := followFixture(t, 1)
	d := newFollowDaemon(t, cfg)
	ctx := boundedTestContext(t)
	credentials, err := store.OpenCredentials(d.credentialsPath)
	if err != nil {
		t.Fatal(err)
	}
	defer credentials.Close()
	start := &store.EpisodeStart{Episode: "selected", Schema: "world/episode-start/v1", StartedAt: "1970-01-01T00:00:10Z"}
	caps := func(budget, expiry int64) string {
		return string(mustMarshal(t, []broker.Capability{
			{Effect: "read", Scope: "worktree", Budget: budget, ExpiresAt: expiry},
			{Effect: "exec", Scope: "/private/project-one", Budget: 11, ExpiresAt: expiry},
			{Effect: "exec", Scope: "/private/project-two", Budget: 3, ExpiresAt: expiry},
		}))
	}
	first := store.SessionRow{CredentialID: "secret-first", EpisodeID: start.Episode, GrantsJSON: caps(80, 99), CreatedAt: 10, ExpiresAt: 99}
	second := store.SessionRow{CredentialID: "secret-second", EpisodeID: start.Episode, GrantsJSON: caps(10, 200), CreatedAt: 20, ExpiresAt: 200}
	if err := credentials.MintSessionWithStart(ctx, first, start); err != nil {
		t.Fatal(err)
	}
	if err := credentials.MintSession(ctx, second); err != nil {
		t.Fatal(err)
	}
	start.Resumes = []store.EpisodeResume{{ResumedFrom: first.CredentialID, CredentialID: second.CredentialID, ResumedAt: "1970-01-01T00:00:20Z"}}

	req := reportObject(broker.EffectRequestV1, []byte("private request body"))
	if err := d.store.PutObject(ctx, req); err != nil {
		t.Fatal(err)
	}
	appendEffect := func(episode string, cost int64) string {
		t.Helper()
		id, _, err := d.store.AppendNextEffectIntent(ctx, episode, store.EffectIntent{
			EpisodeID: episode, Effect: "read", Scope: "worktree", Cost: cost, RequestRef: req.Hash, LogicalTime: 50,
		})
		if err != nil {
			t.Fatal(err)
		}
		return id
	}
	executed := appendEffect(start.Episode, 17)
	rec := reportObject(broker.EffectRecordV1, broker.EncodeRecord(broker.EffectRecord{
		Effect: "read", Scope: "worktree", Cost: 17, BudgetBefore: 80, BudgetAfter: 63,
		Allowed: true, RequestRef: req.Hash, ResultRef: hashref.SumSHA256([]byte("private result")),
	}))
	if err := d.store.PutObject(ctx, rec); err != nil {
		t.Fatal(err)
	}
	if _, _, err := d.store.AppendEffectOutcome(ctx, executed, store.EffectOutcome{
		InvocationID: executed, Status: "succeeded", RecordRef: rec.Hash, LogicalTime: 50,
	}); err != nil {
		t.Fatal(err)
	}
	appendEffect(start.Episode, 2)      // A durable reservation still debits spendable balance.
	appendEffect("selected:other", 123) // Even nested key-range episodes never contribute.

	out := d.surfacePermissions(ctx, start, 100)
	if out == nil || len(out.Sessions) != 2 || out.Truncated {
		t.Fatalf("permissions: %+v", out)
	}
	if out.Sessions[0].State != "expired" || out.Sessions[0].Current ||
		out.Sessions[1].State != "active" || !out.Sessions[1].Current {
		t.Fatalf("session states: %+v", out.Sessions)
	}
	old, active := out.Sessions[0].Grants[0], out.Sessions[1].Grants[0]
	if old.ScopeKey != active.ScopeKey || *old.Budget != 80 || *active.Budget != 10 || *active.Used != 17 ||
		*active.Reserved != 2 || *active.Remaining != 0 || *old.Remaining != 0 || *old.LedgerRemaining != 61 {
		t.Fatalf("incorrect projected budget: old=%+v active=%+v", old, active)
	}
	unused, otherUnused := out.Sessions[1].Grants[1], out.Sessions[1].Grants[2]
	if *unused.Used != 0 || *unused.Remaining != 11 || unused.Scope != otherUnused.Scope || unused.ScopeKey == otherUnused.ScopeKey {
		t.Fatalf("unused/colliding private permission: %+v %+v", unused, otherUnused)
	}
	wire := string(mustMarshal(t, out))
	for _, private := range []string{"secret-first", "secret-second", "/private/", "private request body", "private result"} {
		if strings.Contains(wire, private) {
			t.Fatalf("private permission data %q escaped", private)
		}
	}
	for _, row := range []store.SessionRow{first, second} {
		if err := credentials.RevokeSession(ctx, row.CredentialID); err != nil {
			t.Fatal(err)
		}
	}
	out = d.surfacePermissions(ctx, start, 100)
	if len(out.Observed) != 1 || out.Observed[0].State != "observed" || *out.Observed[0].Budget != 80 ||
		*out.Observed[0].Used != 17 || *out.Observed[0].Reserved != 2 || *out.Observed[0].LedgerRemaining != 61 ||
		out.Observed[0].Remaining != nil || out.Sessions[1].State != "revoked-or-unavailable" {
		t.Fatalf("missing credentials erased or strengthened observed evidence: %+v", out)
	}
}
