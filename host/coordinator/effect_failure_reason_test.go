package coordinator

// Row 255 (design_docs/planned/w-exec-failure-carries-a-reason.md §5, TC): a
// handler's failure line reaches the caller's world block and the effect
// record, never the cause's text, and the committed output replays byte for
// byte from the record alone. A failure with no line keeps today's bytes.

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

func TestEffectFailureReasonReachesTheWorldBlockAndReplays(t *testing.T) {
	const line = "seed-cache: GOCACHE: timed out after 9s"
	for _, tc := range []struct {
		name  string
		err   error
		want  string
		bytes string // a substring the output must (reason) or must not (no reason) carry
	}{
		{"classified", &broker.HandlerFailure{Class: broker.FailSeedCache, Detail: "GOCACHE: timed out after 9s",
			Err: errors.New("cp /Users/op/secret: boom")}, line, `"reason":"` + line + `"`},
		{"unclassified", errors.New("handler failed /Users/op/secret"), "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := fxRig(t, false)
			h := &probe{delay: func(context.Context) error { return tc.err }}
			c := r.fxCoordinator(r.st, &phaseRunner{plan: readPlan(false)}, broker.Registry{fxRead: h})
			res, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t1", fxGrants(5)))
			if err != nil {
				t.Fatalf("Dispatch: %v", err)
			}
			out := string(res.OutputBytes)
			if strings.Contains(out, "/Users/op/secret") {
				t.Fatalf("the cause's text reached the output: %s", out)
			}
			if tc.want != "" && !strings.Contains(out, tc.bytes) {
				t.Fatalf("output %s lacks %s", out, tc.bytes)
			}
			if tc.want == "" && strings.Contains(out, `"reason"`) {
				t.Fatalf("a failure with no line grew a reason key: %s", out)
			}
			wo := decodeWorld(t, res.OutputBytes)
			if len(wo.World.Effects) != 1 || wo.World.Effects[0].Status != StatusFailed {
				t.Fatalf("world block = %+v, want one failed effect", wo.World)
			}
			rec := r.recordV2(res.RecordRef)
			if er := r.effectRecord(hashref.MustParse(rec.Effects[0])); !er.Failed || er.Failure != tc.want {
				t.Fatalf("effect record failure = %q, want %q", er.Failure, tc.want)
			}
			replayed, err := c.Replay(boundedTestContext(t), res.InvocationID, r.replayBinder())
			if err != nil || string(replayed) != out {
				t.Fatalf("replay = %s, %v\nwant     %s", replayed, err, out)
			}
		})
	}
}
