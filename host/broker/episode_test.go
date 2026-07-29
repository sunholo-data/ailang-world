package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

const episodeTransitionV1 = "world/test-episode-transition/v1"

type episodeEvidence struct {
	Kind      string `json:"kind"`
	RecordRef string `json:"recordRef"`
}

type episodeTransition struct {
	CapsuleOutput string            `json:"capsuleOutput"`
	Evidence      []episodeEvidence `json:"evidence"`
}

type episodeCall struct {
	request EffectRequest
	payload []byte
	result  []byte
	record  hashref.HashRef
	errType string
}

func episodeObject(semanticID string, payload []byte) store.Object {
	return store.Object{
		Hash:          hashref.SumSHA256(payload),
		InterfaceHash: hashref.SumSHA256([]byte(semanticID)),
		SemanticID:    semanticID,
		Provenance:    "host/broker/episode_test",
		Payload:       payload,
	}
}

func episodeGrants(path string) []Capability {
	return []Capability{
		{Effect: EffectFSRead, Scope: path, ExpiresAt: 100, Budget: 7},
		{Effect: EffectModelInfer, Scope: "episode-model", ExpiresAt: 100, Budget: 8},
		{Effect: EffectHumanApprove, Scope: "release", ExpiresAt: 100, Budget: 5},
		{Effect: EffectHumanPollApproval, Scope: "release", ExpiresAt: 100, Budget: 4},
		{Effect: "Episode.Fail", Scope: "failure", ExpiresAt: 100, Budget: 6},
	}
}

func runEpisodeCapsule(t *testing.T, binary string) (string, hashref.HashRef) {
	t.Helper()
	a := archive.New(filepath.Join(t.TempDir(), "archive.db"))
	interpreter, err := a.Archive(binary)
	if err != nil {
		t.Fatalf("archive pinned interpreter: %v", err)
	}
	source := []byte("module host/capsule/main\n\nexport func main() -> string { \"capsule-transition\" }\n")
	got, err := capsule.New(a, capsule.Config{}).Run(capsule.Entry{
		Interpreter: interpreter,
		Source:      source,
	})
	if err != nil {
		t.Fatalf("capsule transition: %v", err)
	}
	if string(got.Stdout) != "capsule-transition\n" || len(got.Stderr) != 0 {
		t.Fatalf("capsule output stdout=%q stderr=%q", got.Stdout, got.Stderr)
	}
	return string(got.Stdout), interpreter
}

func commitEpisode(
	t *testing.T,
	s *store.Store,
	capsuleOutput string,
	interpreter hashref.HashRef,
	records []hashref.HashRef,
) store.Commit {
	t.Helper()
	body := episodeTransition{CapsuleOutput: capsuleOutput}
	for _, ref := range records {
		body.Evidence = append(body.Evidence, episodeEvidence{
			Kind: "RecordedEffect", RecordRef: ref.String(),
		})
	}
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	transition := episodeObject(episodeTransitionV1, payload)
	entryHash := hashref.SumSHA256(append([]byte("episode-entry:"), payload...))
	commit := store.Commit{
		Objects: []store.Object{transition},
		NextWorld: store.World{
			Ref:       hashref.SumSHA256([]byte("episode-world")),
			Revision:  0,
			StateRoot: hashref.SumSHA256([]byte("episode-state")),
			LogHead:   entryHash,
		},
		Entry: store.LogEntry{
			Header: store.LogHeader{
				EntryIndex:     0,
				SemanticsEpoch: 1,
				TransitionFn:   hashref.SumSHA256([]byte("episode-transition-fn")),
				Interpreter:    interpreter,
				PrevEntryHash:  hashref.SumSHA256([]byte("episode-genesis-prev")),
				WrittenBy:      "broker-episode-test",
			},
			EntryHash:     entryHash,
			TransitionRef: transition.Hash,
		},
	}
	if err := s.Commit(commit); err != nil {
		t.Fatalf("commit episode: %v", err)
	}
	return commit
}

func readEpisodeEvidence(t *testing.T, s *store.Store, ref hashref.HashRef) []hashref.HashRef {
	t.Helper()
	obj, ok, err := s.GetObject(ref)
	if err != nil || !ok {
		t.Fatalf("transition object %s: ok=%v err=%v", ref, ok, err)
	}
	var body episodeTransition
	if err := json.Unmarshal(obj.Payload, &body); err != nil {
		t.Fatalf("decode transition: %v", err)
	}
	refs := make([]hashref.HashRef, len(body.Evidence))
	for i, item := range body.Evidence {
		if item.Kind != "RecordedEffect" {
			t.Fatalf("evidence[%d].kind = %q", i, item.Kind)
		}
		refs[i], err = hashref.Parse(item.RecordRef)
		if err != nil {
			t.Fatalf("evidence[%d]: %v", i, err)
		}
		if _, found, getErr := s.GetObject(refs[i]); getErr != nil || !found {
			t.Fatalf("evidence record %s: found=%v err=%v", refs[i], found, getErr)
		}
	}
	return refs
}

func TestEpisodeLiveReplayThreeArmsAndEvidence(t *testing.T) {
	binary := os.Getenv("AILANG_BIN")
	if binary == "" {
		t.Fatal("AILANG_BIN must name the pinned released interpreter")
	}
	inputPath := filepath.Join(t.TempDir(), "episode.txt")
	if err := os.WriteFile(inputPath, []byte("episode-fs-bytes"), 0o600); err != nil {
		t.Fatal(err)
	}
	s, err := store.Open(filepath.Join(t.TempDir(), "world.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	capsuleOutput, interpreter := runEpisodeCapsule(t, binary)
	model, err := NewModelHandler(ModelHandlerConfig{AILANGPath: binary, Stub: true})
	if err != nil {
		t.Fatal(err)
	}
	human := NewHumanHandler(s)
	liveDispatches := 0
	failHandler := HandlerFunc(func(context.Context, EffectRequest, []byte) ([]byte, error) {
		liveDispatches++
		return nil, errors.New("deliberate episode failure")
	})
	live := NewSession(s, episodeGrants(inputPath), Registry{
		EffectFSRead:            FSHandler{},
		EffectModelInfer:        model,
		EffectHumanApprove:      human,
		EffectHumanPollApproval: human,
		"Episode.Fail":          failHandler,
	})
	ctx := context.Background()
	calls := []episodeCall{
		{request: EffectRequest{Effect: EffectFSRead, Scope: inputPath, Cost: 2, Now: 10}},
		{request: EffectRequest{Effect: EffectModelInfer, Scope: "episode-model", Cost: 3, Now: 11},
			payload: []byte("episode prompt")},
		{request: EffectRequest{Effect: EffectHumanApprove, Scope: "release", Cost: 2, Now: 12},
			payload: mustApprovalJSON(approvalInputWire{Requester: "episode-agent"})},
	}
	for i := range calls {
		calls[i].result, calls[i].record, err = live.Invoke(ctx, calls[i].request, calls[i].payload)
		if err != nil {
			t.Fatalf("live call %d: %v", i, err)
		}
	}
	requestRef := decodePendingRef(t, calls[2].result)
	if _, err := DecideApproval(s, requestRef, "approve", "episode-operator", 13); err != nil {
		t.Fatalf("decide approval: %v", err)
	}
	calls = append(calls, episodeCall{
		request: EffectRequest{Effect: EffectHumanPollApproval, Scope: "release", Cost: 1, Now: 14},
		payload: mustApprovalJSON(approvalInputWire{RequestRef: requestRef.String()}),
	})
	calls[3].result, calls[3].record, err = live.Invoke(ctx, calls[3].request, calls[3].payload)
	if err != nil {
		t.Fatalf("live poll: %v", err)
	}
	calls = append(calls,
		episodeCall{request: EffectRequest{Effect: "Episode.Fail", Scope: "failure", Cost: 2, Now: 15},
			errType: "failed"},
		episodeCall{request: EffectRequest{Effect: "Episode.Denied", Scope: "denied", Cost: 1, Now: 16},
			errType: "denied"},
	)
	_, calls[4].record, err = live.Invoke(ctx, calls[4].request, nil)
	var failed *EffectFailedError
	if !errors.As(err, &failed) {
		t.Fatalf("live failed arm error = %T %v", err, err)
	}
	_, calls[5].record, err = live.Invoke(ctx, calls[5].request, nil)
	var denied *DenialError
	if !errors.As(err, &denied) {
		t.Fatalf("live denied arm error = %T %v", err, err)
	}
	if liveDispatches != 1 {
		t.Fatalf("live failure dispatches = %d, want 1", liveDispatches)
	}

	records := make([]hashref.HashRef, len(calls))
	for i := range calls {
		records[i] = calls[i].record
		obj, ok, getErr := s.GetObject(calls[i].record)
		if getErr != nil || !ok {
			t.Fatalf("live record %d: ok=%v err=%v", i, ok, getErr)
		}
		rec, decodeErr := DecodeRecord(obj.Payload)
		if decodeErr != nil || !RecordConsistent(rec) {
			t.Fatalf("live record %d inconsistent: %#v err=%v", i, rec, decodeErr)
		}
		if i == 4 && (!rec.Allowed || !rec.Failed || !rec.ResultRef.IsZero() ||
			rec.BudgetBefore != 6 || rec.BudgetAfter != 4) {
			t.Fatalf("failed record did not retain debit and zero result: %#v", rec)
		}
		if i == 5 && (rec.Allowed || rec.Failed || !rec.ResultRef.IsZero()) {
			t.Fatalf("denied record has wrong arm: %#v", rec)
		}
	}
	commit := commitEpisode(t, s, capsuleOutput, interpreter, records)
	liveEvidence := readEpisodeEvidence(t, s, commit.Entry.TransitionRef)
	if len(liveEvidence) != len(records) {
		t.Fatalf("live evidence count = %d, want %d", len(liveEvidence), len(records))
	}

	replayDispatches := 0
	countingStub := HandlerFunc(func(context.Context, EffectRequest, []byte) ([]byte, error) {
		replayDispatches++
		return []byte("stub-result-must-be-ignored"), nil
	})
	stubs := Registry{}
	for _, call := range calls {
		stubs[call.request.Effect] = countingStub
	}
	replay := NewReplaySession(s, episodeGrants(inputPath), stubs, records)
	for i, call := range calls {
		got, gotRef, replayErr := replay.Invoke(ctx, call.request, call.payload)
		switch call.errType {
		case "failed":
			var gotFailed *EffectFailedError
			if !errors.As(replayErr, &gotFailed) {
				t.Fatalf("replay failed arm = %T %v", replayErr, replayErr)
			}
		case "denied":
			var gotDenied *DenialError
			if !errors.As(replayErr, &gotDenied) {
				t.Fatalf("replay denied arm = %T %v", replayErr, replayErr)
			}
		default:
			if replayErr != nil {
				t.Fatalf("replay call %d: %v", i, replayErr)
			}
			if !bytes.Equal(got, call.result) {
				t.Errorf("replay byte identity call %d: got %q want %q", i, got, call.result)
			}
		}
		if gotRef != call.record {
			t.Errorf("replay record ref call %d = %s, want %s", i, gotRef, call.record)
		}
	}
	if replayDispatches != 0 {
		t.Fatalf("replay handler dispatches = %d, want 0", replayDispatches)
	}

	mismatch := NewReplaySession(s, episodeGrants(inputPath), stubs, records[:1])
	wrong := calls[0].request
	wrong.Cost++
	_, _, err = mismatch.Invoke(ctx, wrong, calls[0].payload)
	var mismatchGap *ReplayGapError
	if !errors.As(err, &mismatchGap) {
		t.Errorf("replay mismatch error = %T %v, want *ReplayGapError", err, err)
	}

	firstRecord, _, _ := s.GetObject(records[0])
	firstDecoded, err := DecodeRecord(firstRecord.Payload)
	if err != nil {
		t.Fatal(err)
	}
	gapStore := &handlerRecordingStore{
		base: s, deleted: map[hashref.HashRef]bool{firstDecoded.ResultRef: true},
	}
	gapDispatches := 0
	gapStub := HandlerFunc(func(context.Context, EffectRequest, []byte) ([]byte, error) {
		gapDispatches++
		return nil, fmt.Errorf("live fallback")
	})
	gap := newSession(gapStore, episodeGrants(inputPath),
		Registry{EffectFSRead: gapStub}, Replay, records[:1])
	_, _, err = gap.Invoke(ctx, calls[0].request, calls[0].payload)
	var replayGap *ReplayGapError
	if !errors.As(err, &replayGap) {
		t.Errorf("missing result error = %T %v, want *ReplayGapError", err, err)
	}
	if gapDispatches != 0 {
		t.Fatalf("gap live fallback dispatches = %d, want 0", gapDispatches)
	}

	replayedEvidence := readEpisodeEvidence(t, s, commit.Entry.TransitionRef)
	for i := range records {
		if replayedEvidence[i] != liveEvidence[i] || replayedEvidence[i] != records[i] {
			t.Errorf("evidence[%d] = %s, live %s, record %s",
				i, replayedEvidence[i], liveEvidence[i], records[i])
		}
	}
}
