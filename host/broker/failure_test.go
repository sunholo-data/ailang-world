package broker

// Row 255 (design_docs/planned/w-exec-failure-carries-a-reason.md §5): the
// failure line's pure rule (TL), its record codec (TG), the classifier (TK),
// one Execute-driven case per reachable class through a live broker
// session (TX), and the store marker's order (TV).

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/procbound"
	"github.com/sunholo-data/ailang-world/host/store"
)

// TL: the bounding rule.
func TestFailureLineBoundsAndSanitises(t *testing.T) {
	if len(failureClasses) != 10 {
		t.Fatalf("failure classes = %d, want the closed set of 10", len(failureClasses))
	}
	long := strings.Repeat("é", 300)
	for _, tc := range []struct{ class, detail, want string }{
		{FailExecStackDrift, "/opt/node: digest changed", "exec-stack-drift: /opt/node: digest changed"},
		{FailTimeBudgetSpent, "", "time-budget-spent"},
		{FailSeedCache, "  GOCACHE:\n timed\tout\x00 after 9s \r", "seed-cache: GOCACHE: timed out after 9s"},
		{FailExecTools, "bad \xff byte", "exec-tools: bad ? byte"},
		{FailProjectRoot, "evil\u202egnp.txt\u2028next", "project-root: evil gnp.txt next"},
		{"not-a-class", "anything", ""},
		{"", "", ""},
	} {
		if got := FailureLine(tc.class, tc.detail); got != tc.want {
			t.Errorf("FailureLine(%q, %q) = %q, want %q", tc.class, tc.detail, got, tc.want)
		}
		if tc.want != "" && !ValidFailureLine(tc.want) {
			t.Errorf("ValidFailureLine(%q) = false for a line the rule produced", tc.want)
		}
	}
	cut := FailureLine(FailProjectRoot, long)
	if len(cut) > FailureLineMax || !strings.HasSuffix(cut, "…") || !utf8.ValidString(cut) || !ValidFailureLine(cut) {
		t.Fatalf("long line = %d bytes %q; want ≤ %d, ending …, valid UTF-8, a fixed point", len(cut), cut, FailureLineMax)
	}
	for _, bad := range []string{"", "nope: x", "exec-tools: a\nb", "exec-tools:  two spaces", "exec-tools: " + long} {
		if ValidFailureLine(bad) {
			t.Errorf("ValidFailureLine(%q) = true", bad)
		}
	}
}

// TG: the record codec. A failed record's line rides in an optional last
// field; every other record keeps its bytes (TestRecordGoldenBytes).
func TestRecordFailureGoldenBytes(t *testing.T) {
	rec := EffectRecord{
		Effect: "Workspace.Exec", Scope: "worktree", Cost: 1, BudgetBefore: 5, BudgetAfter: 4,
		Allowed: true, Failed: true, RequestRef: hashref.SumSHA256([]byte("request")),
		Failure: "seed-cache: GOCACHE: timed out after 9s",
	}
	const want = `{"effect":"Workspace.Exec","scope":"worktree","cost":1,"budgetBefore":5,"budgetAfter":4,"allowed":true,"failed":true,"denial":"","requestRef":"sha256:1f58b9145b24d108d7ac38887338b3ea3229833b9c1e418250343f907bfd1047","resultRef":"","failure":"seed-cache: GOCACHE: timed out after 9s"}`
	if got := string(EncodeRecord(rec)); got != want {
		t.Fatalf("golden bytes differ:\ngot  %s\nwant %s", got, want)
	}
	back, err := DecodeRecord([]byte(want))
	if err != nil || back != rec {
		t.Fatalf("round trip = %#v, %v", back, err)
	}
	for name, payload := range map[string]string{
		"on a success": strings.Replace(strings.Replace(want, `"failed":true`, `"failed":false`, 1), `"resultRef":""`,
			`"resultRef":"sha256:f6a214f7a5fcda0c2cee9660b7fc29f5649e3c68aad48e20e950137c98913a68"`, 1),
		"on a denial":       strings.Replace(strings.Replace(want, `"allowed":true,"failed":true,"denial":""`, `"allowed":false,"failed":false,"denial":"denied:budget"`, 1), `"budgetAfter":4`, `"budgetAfter":5`, 1),
		"not a fixed point": strings.Replace(want, `GOCACHE: timed`, `GOCACHE:  timed`, 1),
		"unknown class":     strings.Replace(want, `"failure":"seed-cache:`, `"failure":"seed-cash:`, 1),
	} {
		if _, err := DecodeRecord([]byte(payload)); err == nil || !strings.Contains(err.Error(), "failure") {
			t.Errorf("%s: DecodeRecord(%s) = %v, want a failure-line refusal", name, payload, err)
		}
	}
}

// TK: the classifier, with no process run. A typed cause overrides the step;
// the line names only disclosed paths; err's text never reaches it.
func TestExecFailureClassifies(t *testing.T) {
	d := execDisclosed{archive: "/state/archive/x/node_modules", paths: []string{"/opt/node", "", "/usr/bin/go", ""}}
	secret := errors.New("open /Users/op/.secret/token: permission denied")
	for _, tc := range []struct {
		name          string
		class, detail string
		err           error
		want          string
	}{
		{"archive file", FailExecStackDrift, "unreadable",
			&ExecStackDriftError{Path: "/state/archive/x/node_modules/@anthropic-ai/sandbox-runtime/dist/cli.js", Why: "changed since it was archived"},
			"exec-stack-drift: srt archive @anthropic-ai/sandbox-runtime/dist/cli.js: changed since it was archived"},
		{"node digest", FailExecStackDrift, "unreadable", &ExecStackDriftError{Path: "/opt/node", Why: "sha256 sha256:aa, verified as sha256:bb"},
			"exec-stack-drift: /opt/node: digest changed"},
		{"undisclosed path, os error", FailExecStackDrift, "unreadable", &ExecStackDriftError{Path: "/Users/op/lib/libnode.dylib", Why: secret.Error()},
			"exec-stack-drift: libnode.dylib: unreadable"},
		{"drift through resolve", FailExecTools, "the per-call re-check of the exec tools failed",
			fmt.Errorf("wrapped: %w", &ExecStackDriftError{Path: "/usr/bin/go", Why: "missing from the archive"}),
			"exec-stack-drift: /usr/bin/go: missing from the archive"},
		{"settings mode", FailExecSettingsDrift, "unreadable", &ExecSettingsError{Path: "/state/exec/ep1.proj.srt.json", Why: "mode -rw-r--r--, want a regular 0600 file"},
			"exec-settings-drift: ep1.proj.srt.json: not a regular 0600 file"},
		{"already classified", FailExecCache, "the episode cache cannot be prepared",
			&HandlerFailure{Class: FailSeedCache, Detail: "GOCACHE: " + seedWhat(&HandlerTimeoutError{Timeout: 9200 * time.Millisecond}), Err: secret},
			"seed-cache: GOCACHE: timed out after 9.2s"},
		{"cleanup incomplete", FailCleanupIncomplete, "the command's processes were not all reaped",
			fmt.Errorf("run: %w", procbound.ErrCleanupIncomplete), "cleanup-incomplete: the command's processes were not all reaped"},
		{"untyped cause", FailSandboxStart, "the sandbox could not be started", secret, "sandbox-start: the sandbox could not be started"},
	} {
		got := execFailure(tc.class, tc.detail, tc.err, d)
		if got.Line() != tc.want || strings.Contains(got.Line(), "/Users/op") {
			t.Errorf("%s: line %q, want %q", tc.name, got.Line(), tc.want)
		}
		if !errors.Is(got, tc.err) && got != tc.err {
			t.Errorf("%s: the failure does not unwrap to its cause", tc.name)
		}
	}
	for _, tc := range []struct {
		err  error
		want string
	}{
		{&HandlerTimeoutError{Timeout: time.Millisecond}, "timed out after 1ms"},
		{errors.Join(&HandlerTimeoutError{Timeout: 9 * time.Second}, procbound.ErrCleanupIncomplete), "timed out after 9s"},
		{procbound.ErrCleanupIncomplete, "the copy was not reaped"},
		{errors.New("start: no such file"), "the copy could not run"},
	} {
		if got := seedWhat(tc.err); got != tc.want {
			t.Errorf("seedWhat(%v) = %q, want %q", tc.err, got, tc.want)
		}
	}
}

// expiredCtx answers Execute with an already-expired context, so the budget
// and seed classes are reached through a live broker session (whose own
// store writes need a live context).
type expiredCtx struct{ h *ExecHandler }

func (e expiredCtx) Execute(_ context.Context, req EffectRequest, payload []byte) ([]byte, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	return e.h.Execute(ctx, req, payload)
}

// TX: one Execute-driven case per reachable class. Each failure is recorded
// with its line, returned as EffectFailedError.Reason, and keeps the cause's
// own error text; no line names the state dir or the operator HOME.
func TestExecFailureClassesReachTheRecord(t *testing.T) {
	seedProfile := func(seed string) func(m map[string]any) {
		return func(m map[string]any) { m["caches"] = map[string]any{"GOCACHE": map[string]any{"seed": seed}} }
	}
	for _, tc := range []struct {
		name    string
		opts    func(t *testing.T) execFixtureOpts
		build   func(t *testing.T) *execFixture // instead of newExecFixture(opts)
		mutate  func(t *testing.T, f *execFixture) Handler
		scope   string
		want    func(f *execFixture) string
		errText func(f *execFixture) string // a substring the cause's own text keeps
	}{
		{name: "exec-misrouted", scope: "/",
			want: func(*execFixture) string { return "exec-misrouted: Workspace.Exec@/" }},
		{name: "exec-tools", build: func(t *testing.T) *execFixture {
			m := newMovingExec(t, execFixtureOpts{})
			m.state.Close()
			return m.f
		}, want: func(*execFixture) string { return "exec-tools: the per-call re-check of the exec tools failed" },
			errText: func(*execFixture) string { return "closed" }},
		{name: "exec-stack-drift archive", mutate: func(t *testing.T, f *execFixture) Handler {
			flipByte(t, f.sandbox.CLI())
			return f.h
		}, want: func(*execFixture) string {
			return "exec-stack-drift: srt archive @anthropic-ai/sandbox-runtime/dist/cli.js: changed since it was archived"
		}, errText: func(f *execFixture) string { return f.sandbox.CLI() }},
		{name: "exec-stack-drift node", mutate: func(t *testing.T, f *execFixture) Handler {
			flipByte(t, f.sandbox.Node)
			return f.h
		}, want: func(f *execFixture) string { return "exec-stack-drift: " + f.sandbox.Node + ": digest changed" },
			errText: func(f *execFixture) string { return f.sandbox.Node }},
		{name: "exec-settings-drift", mutate: func(t *testing.T, f *execFixture) Handler {
			if err := os.WriteFile(filepath.Join(f.state, "exec", "ep1.proj.srt.json"), []byte(`{"filesystem":{}}`), 0o600); err != nil {
				t.Fatal(err)
			}
			return f.h
		}, want: func(*execFixture) string {
			return "exec-settings-drift: ep1.proj.srt.json: bytes differ from World's rendering"
		}, errText: func(f *execFixture) string { return filepath.Join(f.state, "exec", "ep1.proj.srt.json") }},
		{name: "project-root", opts: func(*testing.T) execFixtureOpts {
			return execFixtureOpts{profile: func(m map[string]any) { m["root"] = "missing" }}
		}, want: func(*execFixture) string {
			return `project-root: root "missing" is not a directory inside the worktree`
		}},
		{name: "exec-cache", mutate: func(t *testing.T, f *execFixture) Handler {
			if err := os.MkdirAll(filepath.Dir(f.cache), 0o700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(f.cache, []byte("a file where the cache dir goes"), 0o600); err != nil {
				t.Fatal(err)
			}
			return f.h
		}, want: func(*execFixture) string { return "exec-cache: home: cannot be created" }},
		{name: "seed-cache exited", opts: func(t *testing.T) execFixtureOpts {
			return execFixtureOpts{profile: seedProfile(filepath.Join(canonicalTempDir(t), "no-such-seed"))}
		}, want: func(*execFixture) string { return "seed-cache: GOCACHE: the copy exited 1 (attempt 1 of 3)" }},
		{name: "seed-cache timed out", opts: func(t *testing.T) execFixtureOpts {
			seed := filepath.Join(canonicalTempDir(t), "seed")
			if err := os.MkdirAll(seed, 0o755); err != nil {
				t.Fatal(err)
			}
			return execFixtureOpts{profile: seedProfile(seed)}
		}, mutate: func(_ *testing.T, f *execFixture) Handler {
			f.h.seedBudget = time.Millisecond
			f.h.seedCommand = slowSeedCommand
			return f.h
		},
			want: func(*execFixture) string { return "seed-cache: GOCACHE: timed out after 1ms (attempt 1 of 3)" }},
		{name: "time-budget-spent", mutate: func(_ *testing.T, f *execFixture) Handler { return expiredCtx{f.h} },
			want: func(*execFixture) string {
				return `time-budget-spent: the call's time was spent before "test" could start`
			}, errText: func(*execFixture) string { return "budget is spent" }},
		{name: "sandbox-start", mutate: func(t *testing.T, f *execFixture) Handler {
			if err := os.Chmod(f.sandbox.Node, 0o644); err != nil {
				t.Fatal(err)
			}
			return f.h
		}, want: func(*execFixture) string { return "sandbox-start: the sandbox could not be started" },
			errText: func(*execFixture) string { return "start handler subprocess" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			opts := execFixtureOpts{}
			if tc.opts != nil {
				opts = tc.opts(t)
			}
			var f *execFixture
			if tc.build != nil {
				f = tc.build(t)
			} else {
				f = newExecFixture(t, opts)
			}
			var h Handler = f.h
			if tc.mutate != nil {
				h = tc.mutate(t, f)
			}
			scope := WorkspaceScope
			if tc.scope != "" {
				scope = tc.scope
			}
			st := openTestStore(t)
			recording := &cfj2RecordingStore{base: st}
			s := newSession(recording, "ep1", []Capability{{Effect: EffectWorkspaceExec, Scope: scope, ExpiresAt: 100, Budget: 5}},
				Registry{EffectWorkspaceExec: h}, Live, nil)
			payload, _ := json.Marshal(map[string]any{"command": "test", "args": []string{"a"}})
			_, ref, err := s.Invoke(boundedTestContext(t), EffectRequest{Effect: EffectWorkspaceExec, Scope: scope, Cost: 1, Now: 1}, payload)
			var failed *EffectFailedError
			if !errors.As(err, &failed) {
				t.Fatalf("Invoke = %v, want *EffectFailedError", err)
			}
			want := tc.want(f)
			if failed.Reason != want {
				t.Fatalf("EffectFailedError.Reason = %q, want %q (cause: %v)", failed.Reason, want, errors.Unwrap(failed))
			}
			if len(recording.records) != 1 || recording.records[0].Hash != ref {
				t.Fatalf("records %d, want the one failure record %s", len(recording.records), ref)
			}
			rec, err := DecodeRecord(recording.records[0].Payload)
			if err != nil || !rec.Allowed || !rec.Failed || rec.Failure != want || !RecordConsistent(rec) {
				t.Fatalf("record = %#v, %v; want failed with failure %q", rec, err, want)
			}
			if strings.Contains(want, f.state) || strings.Contains(want, f.home) {
				t.Fatalf("line %q names the state dir or the operator HOME", want)
			}
			if tc.errText != nil && !strings.Contains(errors.Unwrap(failed).Error(), tc.errText(f)) {
				t.Fatalf("cause text %q lost %q", errors.Unwrap(failed).Error(), tc.errText(f))
			}
		})
	}
}

type lineFailure struct {
	err   error
	calls int
}

func (h *lineFailure) Execute(context.Context, EffectRequest, []byte) ([]byte, error) {
	h.calls++
	return nil, h.err
}

// markRefusingStore is a store whose schema version cannot be raised.
type markRefusingStore struct {
	*cfj2RecordingStore
	markErr  error
	outcomes int
}

func (s *markRefusingStore) MarkFailureLines(context.Context) error { return s.markErr }

func (s *markRefusingStore) AppendEffectOutcome(ctx context.Context, id string, outcome store.EffectOutcome) (int64, hashref.HashRef, error) {
	s.outcomes++
	return s.cfj2RecordingStore.AppendEffectOutcome(ctx, id, outcome)
}

// orderStore logs the order of the mark and the effect-record put.
type orderStore struct {
	*cfj2RecordingStore
	order []string
}

func (o *orderStore) MarkFailureLines(ctx context.Context) error {
	o.order = append(o.order, "mark")
	return o.cfj2RecordingStore.MarkFailureLines(ctx)
}

func (o *orderStore) PutObject(ctx context.Context, obj store.Object) error {
	if obj.SemanticID == EffectRecordV1 {
		o.order = append(o.order, "record")
	}
	return o.cfj2RecordingStore.PutObject(ctx, obj)
}

// TV (§3.7): the store is raised to schema 6 before the first lined record;
// a failure with no line leaves it at 5; a failed mark returns a write error
// without a record or outcome, rerun or refund.
func TestFailureLineMarksTheStoreFirst(t *testing.T) {
	classified := &HandlerFailure{Class: FailSeedCache, Detail: "GOCACHE: timed out after 9s", Err: errors.New("boom")}
	for _, tc := range []struct {
		name        string
		err         error
		refuse      bool
		wantVersion int
		wantLine    string
	}{
		{"line marks", classified, false, 6, "seed-cache: GOCACHE: timed out after 9s"},
		{"no line, no mark", errors.New("plain"), false, 5, ""},
		{"mark refused, write error", classified, true, 5, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "world.db")
			st, err := store.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			defer st.Close()
			recording := &cfj2RecordingStore{base: st}
			ordered := &orderStore{cfj2RecordingStore: recording}
			markErr := errors.New("injected mark error")
			refusing := &markRefusingStore{cfj2RecordingStore: recording, markErr: markErr}
			var backing objectStore = ordered
			if tc.refuse {
				backing = refusing
			}
			handler := &lineFailure{err: tc.err}
			s := newSession(backing, "ep1", []Capability{{Effect: EffectWorkspaceExec, Scope: WorkspaceScope, ExpiresAt: 100, Budget: 5}},
				Registry{EffectWorkspaceExec: handler}, Live, nil)
			result, ref, err := s.Invoke(boundedTestContext(t), EffectRequest{Effect: EffectWorkspaceExec, Scope: WorkspaceScope, Cost: 1, Now: 1}, []byte(`{}`))
			if handler.calls != 1 {
				t.Fatalf("handler calls = %d, want exactly 1", handler.calls)
			}
			if got := s.grants[0].Budget; got != 4 {
				t.Fatalf("grant budget = %d, want 4 (debit stands)", got)
			}
			if v, err := store.WorldSchemaVersion(path); err != nil || v != tc.wantVersion {
				t.Fatalf("store user_version = %d, %v; want %d", v, err, tc.wantVersion)
			}
			if tc.refuse {
				var writeErr *OutcomeWriteError
				if !errors.As(err, &writeErr) || !errors.Is(err, markErr) {
					t.Fatalf("Invoke = %v, want *OutcomeWriteError wrapping mark error", err)
				}
				if writeErr.Effect != EffectWorkspaceExec || writeErr.Scope != WorkspaceScope || !writeErr.RecordRef.IsZero() ||
					!strings.Contains(writeErr.Error(), "broker: mark store for failure lines: ") {
					t.Fatalf("write error = %#v", writeErr)
				}
				if result != nil || !ref.IsZero() || len(recording.records) != 0 || refusing.outcomes != 0 {
					t.Fatalf("failed mark: result=%q ref=%s records=%d outcomes=%d; want no record or outcome", result, ref, len(recording.records), refusing.outcomes)
				}
				return
			}
			if len(recording.records) != 1 {
				t.Fatalf("record count = %d, want 1", len(recording.records))
			}
			var failed *EffectFailedError
			if !errors.As(err, &failed) {
				t.Fatalf("Invoke = %v, want *EffectFailedError", err)
			}
			if failed.Reason != tc.wantLine {
				t.Fatalf("EffectFailedError.Reason = %q, want %q", failed.Reason, tc.wantLine)
			}
			rec, err := DecodeRecord(recording.records[0].Payload)
			if err != nil || rec.Failure != tc.wantLine {
				t.Fatalf("record failure = %q, %v; want %q", rec.Failure, err, tc.wantLine)
			}
			if tc.wantLine != "" && strings.Join(ordered.order, ",") != "mark,record" {
				t.Fatalf("call order = %v, want the mark before the lined record", ordered.order)
			}
		})
	}
}
