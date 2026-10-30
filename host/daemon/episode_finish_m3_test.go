package daemon

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
)

func TestEpisodeFinishRecordReuseAndSingleUse(t *testing.T) {
	d := newHandlerDaemon(t)
	scopeText := testFinishScope()
	approve := func(text string) FinishApproval {
		t.Helper()
		proof := finishProofFile(t, d, "approve", time.Now().Unix()+100, 0o600)
		w := finishRoute(t, d, FinishRequest{Phase: "approve", Scope: text}, proof, "")
		var result FinishApproval
		if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &result) != nil {
			t.Fatalf("approve %d %s", w.Code, w.Body.String())
		}
		return result
	}
	approval := approve(scopeText)
	scope, _ := broker.ParseFinishApprovalScope(scopeText)
	record := EpisodeFinish{
		Episode: "ep1", Branch: scope.Head, Commit: scope.SHA, Remote: scope.Remote,
		Approval: approval.Approval, Outcome: "push-failed",
		PR: FinishPR{URL: "https://github.com/o/r/pull/1", Number: 1},
	}
	post := func(want int) {
		t.Helper()
		proof := finishProofFile(t, d, "record", time.Now().Unix()+100, 0o600)
		w := finishRoute(t, d, FinishRequest{Phase: "record", Record: &record}, proof, "")
		if w.Code != want {
			t.Fatalf("record %s: %d want %d: %s", record.Outcome, w.Code, want, w.Body.String())
		}
	}
	head := func() string {
		t.Helper()
		ref, exists, err := d.store.GetRegistryHead(boundedTestContext(t), finishHead("ep1"))
		if err != nil || !exists {
			t.Fatalf("head: %v %v", exists, err)
		}
		return ref.String()
	}
	post(200)
	first := head()
	record.Outcome = "done"
	post(200)
	done := head()
	if done == first {
		t.Fatal("done did not move head")
	}
	got, err := d.readEpisodeFinish(boundedTestContext(t), "ep1")
	if err != nil || got.Previous != first || got.Scope != scopeText {
		t.Fatalf("record %+v %v", got, err)
	}
	w := requestRecorder(t, d, "GET", "/v1/episodes/ep1/report", nil)
	var report EpisodeReport
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &report) != nil || report.Finish == nil || report.Finish.Scope != scopeText {
		t.Fatalf("report %d %s", w.Code, w.Body.String())
	}
	for _, outcome := range []string{"push-failed", "done"} {
		record.Outcome = outcome
		post(409)
		if head() != done {
			t.Fatal("used approval moved head")
		}
	}
	// A different approval is permitted for this episode after done.
	scope.TitleSHA256 = scope.BodySHA256
	newScope := broker.FormatFinishApprovalScope(scope)
	approval = approve(newScope)
	record.Approval, record.Scope = approval.Approval, scopeText
	post(403)
	if head() != done {
		t.Fatal("scope mismatch moved head")
	}
	record.Scope = newScope
	post(200)
	got, err = d.readEpisodeFinish(boundedTestContext(t), "ep1")
	if err != nil || got.Previous != done || got.Scope != newScope || got.Approval != approval.Approval {
		t.Fatalf("fix-up %+v %v", got, err)
	}
}
