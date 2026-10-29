package daemon

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
)

func TestEpisodeWritesSeeds(t *testing.T) {
	s := newReportSeed(t)
	call := func(ep, task, tool, path, effect string, result map[string]any) int64 {
		s.input = map[string]any{"path": path}
		defer func() { s.input = nil }()
		return s.call("mcp", ep, task, tool, []seededEffect{{effect: effect, scope: "worktree", cost: 1, before: 100, allowed: true, result: result}}, nil)
	}
	ok := map[string]any{"ok": true}
	first := call("writes", "w", "ailang-write", "a.go", broker.EffectWorkspaceWrite, ok)
	call("writes", "e", "ailang-edit", "a.go", broker.EffectWorkspaceWrite, ok)
	del := call("writes", "d", "ailang-delete", "gone.txt", broker.EffectWorkspaceDelete, ok)
	gen := call("writes", "g", "workspace-exec", "", broker.EffectWorkspaceExec, map[string]any{"worktree_delta": map[string]any{"changes": []map[string]any{{"path": "gen.txt", "op": "created"}}}})
	trunc := call("writes", "t", "ailang-check", "", broker.EffectAilangCheck, map[string]any{"worktree_delta": map[string]any{"truncated": true}})
	errEntry := call("writes", "err", "ailang-cli", "", broker.EffectAilangCLI, map[string]any{"worktree_delta": map[string]any{"error": "snapshot failed"}})
	call("writes", "no", "ailang-write", "refused.go", broker.EffectWorkspaceWrite, map[string]any{"ok": false})
	call("other", "w", "ailang-write", "other.go", broker.EffectWorkspaceWrite, ok)
	call("writes", "outside", "ailang-write", "/private/never-echo", broker.EffectWorkspaceWrite, ok)
	r, found, err := buildEpisodeReport(boundedTestContext(t), s.d.intents, s.d.reads, "writes")
	if err != nil || !found {
		t.Fatalf("report: %v %v", found, err)
	}
	w := r.Writes
	if len(w.Written) != 2 || w.Written[0].Path != "a.go" || w.Written[0].Op != "edit" || w.Written[0].Count != 2 || w.Written[0].FirstEntry != first || w.Written[1].Op != "delete" || w.Written[1].FirstEntry != del {
		t.Fatalf("written: %+v", w.Written)
	}
	if !reflect.DeepEqual(w.Produced, []ReportRead{{Path: "gen.txt", Count: 1, FirstEntry: gen}}) || !reflect.DeepEqual(w.DeltasTruncated, []int64{trunc}) || len(w.DeltaErrors) != 1 || w.DeltaErrors[0].Entry != errEntry || w.Unlisted != 1 {
		t.Fatalf("writes: %+v", w)
	}
	route := requestRecorder(t, s.d, http.MethodGet, "/v1/episodes/writes/report", nil)
	var served EpisodeReport
	if route.Code != http.StatusOK || json.Unmarshal(route.Body.Bytes(), &served) != nil || !reflect.DeepEqual(served.Writes.Produced, w.Produced) {
		t.Fatalf("real report route: %d %s", route.Code, route.Body)
	}
	if r.Notes["writes"] == "" {
		t.Fatal("missing evidence note")
	}
}

func TestEpisodeWritesCapsAndFailureEvidence(t *testing.T) {
	old := reportReadFiles
	reportReadFiles = 1
	defer func() { reportReadFiles = old }()
	s := newReportSeed(t)
	for i, path := range []string{"first.go", "second.go"} {
		s.input = map[string]any{"path": path}
		s.call("mcp", "caps", path, "ailang-write", []seededEffect{{effect: broker.EffectWorkspaceWrite, scope: "worktree", cost: 1, before: 5, allowed: true, result: map[string]any{"ok": true}}}, nil)
		s.input = nil
		result := map[string]any{"worktree_delta": map[string]any{"truncated": true, "error": strings.Repeat("x", 400), "changes": []map[string]any{{"path": path}, {"path": "cut", "path_truncated": true}}}}
		s.call("mcp", "caps", path+"-exec", "workspace-exec", []seededEffect{{effect: broker.EffectWorkspaceExec, scope: "worktree", cost: 1, before: int64(5 - i), allowed: true, result: result}}, nil)
	}
	s.input = map[string]any{"path": "failed.go"}
	s.call("mcp", "caps", "failed", "ailang-write", []seededEffect{{effect: broker.EffectWorkspaceWrite, scope: "worktree", cost: 1, before: 5, allowed: true, failed: true}}, nil)
	r := reportFor(t, s.d, "caps")
	w := r.Writes
	if len(w.Written) != 1 || w.Written[0].Path != "first.go" || w.WrittenTruncated != 1 || len(w.Produced) != 1 || w.ProducedTruncated != 1 || w.Unlisted != 2 {
		t.Fatalf("bounded paths: %+v", w)
	}
	if len(w.DeltasTruncated) != 1 || w.DeltaTruncatedOmitted != 1 || len(w.DeltaErrors) != 1 || w.DeltaErrorsOmitted != 1 || len(w.DeltaErrors[0].Error) > reportReasonMax+3 {
		t.Fatalf("bounded gaps: %+v", w)
	}
}
