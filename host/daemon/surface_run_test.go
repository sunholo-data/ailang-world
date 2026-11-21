package daemon

import (
	"path/filepath"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

func TestSurfaceRunUsesCompletedEvidence(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "world.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	d := &Daemon{store: s}
	ctx := boundedTestContext(t)
	put := func(kind, provenance, payload string) string {
		data := []byte(payload)
		object := store.Object{Hash: hashref.SumSHA256(data), InterfaceHash: hashref.SumSHA256([]byte(kind)), SemanticID: kind, Provenance: provenance, Payload: data}
		if err := s.PutObject(ctx, object); err != nil {
			t.Fatal(err)
		}
		return object.Hash.String()
	}
	const startTime = "2026-11-21T13:00:00Z"
	start := &store.EpisodeStart{Episode: "selected", Agent: "codex", StartedAt: startTime,
		Launches: []store.EpisodeLaunch{{StartedAt: startTime, EndedAt: "2026-11-21T13:01:00Z"}},
		Checks:   []store.StartCheck{{Name: "specified but never run"}}}
	actual := put("world/task-outcome/v1", "harness:floor", `{"episodeId":"selected","label":"episode-finish","task":"unit tests","outcome":{"class":"PASS"}}`)
	other := put("world/task-outcome/v1", "harness:floor", `{"episodeId":"other","label":"episode-finish","task":"other tests","outcome":{"class":"PASS"}}`)
	approval := put(broker.ApprovalDecisionV1, "operator", `{"decision":"approve","now":1795266060}`)
	finish := EpisodeFinish{Episode: "selected", FinishedAt: "2026-11-21T13:03:00Z", Checks: []string{actual, other}, Approval: approval}
	putFinish(t, d, finish, "untrusted")
	run, trusted := d.surfaceRun(ctx, start)
	if trusted != nil || len(run.Decisions) != 1 || len(run.Checks) != 0 || run.Elapsed != nil || run.CostUSD != nil || run.Turns != nil {
		t.Fatalf("untrusted finish %+v", run)
	}
	finish.Outcome = "done"
	putFinish(t, d, finish, finishProvenance)
	run, trusted = d.surfaceRun(ctx, start)
	if trusted == nil || run.Runtime == nil || *run.Runtime != 60 || run.Elapsed == nil || *run.Elapsed != 180 || len(run.Decisions) != 2 ||
		len(run.Checks) != 1 || run.Checks[0].Name != "unit tests" || run.Checks[0].Result != "PASS" {
		t.Fatalf("confirmed finish %+v", run)
	}
}
