package broker

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

type handlerRecordingStore struct {
	base    *store.Store
	records []store.Object
	objects []store.Object
	deleted map[hashref.HashRef]bool
}

func (s *handlerRecordingStore) PutObject(obj store.Object) error {
	if err := s.base.PutObject(obj); err != nil {
		return err
	}
	s.objects = append(s.objects, obj)
	if obj.SemanticID == EffectRecordV1 {
		s.records = append(s.records, obj)
	}
	return nil
}

func (s *handlerRecordingStore) GetObject(ref hashref.HashRef) (store.Object, bool, error) {
	if s.deleted[ref] {
		return store.Object{}, false, nil
	}
	return s.base.GetObject(ref)
}

func (s *handlerRecordingStore) SetRegistryHead(name string, ref hashref.HashRef) error {
	return s.base.SetRegistryHead(name, ref)
}

func (s *handlerRecordingStore) GetRegistryHead(name string) (hashref.HashRef, bool, error) {
	return s.base.GetRegistryHead(name)
}

func handlerSession(t *testing.T, effect, scope string, handler Handler) (*Session, *handlerRecordingStore) {
	t.Helper()
	recording := &handlerRecordingStore{base: openTestStore(t)}
	session := newSession(recording, []Capability{{
		Effect: effect, Scope: scope, ExpiresAt: 100, Budget: 7,
	}}, Registry{effect: handler}, Live, nil)
	return session, recording
}

func assertHandlerFailureRecord(
	t *testing.T,
	session *Session,
	recording *handlerRecordingStore,
	ref hashref.HashRef,
	err error,
	cost int64,
	underlying error,
) {
	t.Helper()
	var failed *EffectFailedError
	if !errors.As(err, &failed) {
		t.Errorf("Invoke error = %v, want *EffectFailedError", err)
	}
	if underlying != nil && !errors.Is(err, underlying) {
		t.Errorf("Invoke error = %v, want errors.Is(..., %v)", err, underlying)
	}
	if failed != nil && (ref.IsZero() || failed.RecordRef != ref) {
		t.Errorf("failure refs = %s and %s, want same non-zero ref", ref, failed.RecordRef)
	}
	if got := len(recording.records); got != 1 {
		t.Fatalf("failure record count = %d, want exactly 1", got)
	}
	rec, decodeErr := DecodeRecord(recording.records[0].Payload)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}
	if !rec.Allowed || !rec.Failed || rec.BudgetAfter != rec.BudgetBefore-cost ||
		!rec.ResultRef.IsZero() || !RecordConsistent(rec) {
		t.Errorf("failure record arm = %#v", rec)
	}
	if got, want := session.grants[0].Budget, rec.BudgetAfter; got != want {
		t.Fatalf("standing debit budget = %d, want %d", got, want)
	}
}

func writeExecutable(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "handler")
	if err := os.WriteFile(path, []byte("#!/bin/sh\nset -eu\n"+body+"\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	return path
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(git, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v: %s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

func TestGitCommitRoundTripScrubsHostileHome(t *testing.T) {
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	repo := t.TempDir()
	runGit(t, repo, "init", "-q")
	if err := os.WriteFile(filepath.Join(repo, "tracked.txt"), []byte("content"), 0o600); err != nil {
		t.Fatal(err)
	}
	runGit(t, repo, "add", "tracked.txt")

	hostileHome := t.TempDir()
	hookDir := filepath.Join(hostileHome, "hooks")
	if err := os.MkdirAll(hookDir, 0o700); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(repo, "hostile-hook-ran")
	hook := filepath.Join(hookDir, "pre-commit")
	if err := os.WriteFile(hook, []byte("#!/bin/sh\n: > \""+marker+"\"\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	config := "[user]\n\tname = HOSTILE NAME\n\temail = hostile@example.invalid\n" +
		"[commit]\n\tgpgsign = false\n[core]\n\thooksPath = " + hookDir + "\n"
	if err := os.WriteFile(filepath.Join(hostileHome, ".gitconfig"), []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", hostileHome)

	handler, err := NewGitHandler(GitHandlerConfig{GitPath: git})
	if err != nil {
		t.Fatal(err)
	}
	session, recording := handlerSession(t, EffectGitCommit, repo, handler)
	result, ref, err := session.Invoke(context.Background(),
		EffectRequest{Effect: EffectGitCommit, Scope: repo, Cost: 2, Now: 1}, []byte("broker commit"))
	if err != nil {
		t.Fatal(err)
	}
	if ref.IsZero() || len(result) == 0 || len(recording.records) != 1 {
		t.Fatalf("Invoke = result %q, ref %s, records %d", result, ref, len(recording.records))
	}
	if got := runGit(t, repo, "rev-parse", "HEAD"); got == "" {
		t.Fatal("commit does not exist")
	}
	if got := runGit(t, repo, "show", "-s", "--format=%an <%ae>|%cn <%ce>", "HEAD"); got !=
		"AILANG World <ailang-world@invalid>|AILANG World <ailang-world@invalid>" {
		t.Fatalf("commit identity = %q, hostile HOME leaked", got)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("hostile hooksPath became observable: marker stat error = %v", err)
	}
}

// pinnedAILANG resolves the pinned interpreter the way the landed replay and
// capsule tests do. Hardcoding /tmp/ailang-v0300/ailang passes on the rig and
// fails in CI, where the pinned v0.30.0 binary is installed under $HOME — a
// green local gate that cannot see the failure. verify_go.sh mandates
// AILANG_BIN, and CI exports it, so a skip here is itself the alarm.
func pinnedAILANG(t *testing.T) string {
	t.Helper()
	bin := os.Getenv("AILANG_BIN")
	if bin == "" {
		t.Skip("AILANG_BIN not set; Model.Infer requires the pinned released ailang binary")
	}
	info, err := os.Stat(bin)
	if err != nil || !info.Mode().IsRegular() {
		t.Skipf("AILANG_BIN %q is not a usable executable: %v", bin, err)
	}
	return bin
}

// A Model.Infer payload is arbitrary caller bytes. The prompt is handed to the
// interpreter as a JSON argument, so every byte a caller can send must survive
// that encoding. Control bytes (0x00-0x1f) are the case a hand-rolled escaper
// silently gets wrong: the subprocess then fails to parse its own arguments and
// — post-M3.B0 — the caller is charged a standing debit for what is really a
// host encoding bug, not a failed effect.
func TestModelPromptEncodesControlBytes(t *testing.T) {
	ailang := pinnedAILANG(t)
	handler, err := NewModelHandler(ModelHandlerConfig{AILANGPath: ailang, Stub: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, payload := range [][]byte{{0x01}, {0x08}, {0x1f}, []byte("a\x00b")} {
		session, recording := handlerSession(t, EffectModelInfer, "model-scope", handler)
		_, _, invokeErr := session.Invoke(context.Background(),
			EffectRequest{Effect: EffectModelInfer, Scope: "model-scope", Cost: 2, Now: 1},
			payload)
		if invokeErr != nil {
			t.Errorf("payload %q: invoke error = %v, want success", payload, invokeErr)
			continue
		}
		if len(recording.records) != 1 {
			t.Fatalf("payload %q: record count = %d, want 1", payload, len(recording.records))
		}
		rec, decodeErr := DecodeRecord(recording.records[0].Payload)
		if decodeErr != nil {
			t.Fatal(decodeErr)
		}
		if !rec.Allowed || rec.Failed || !RecordConsistent(rec) {
			t.Errorf("payload %q: want the success arm, got %#v", payload, rec)
		}
	}
}

func TestModelStubRoundTripDeterministicRecordedBytes(t *testing.T) {
	ailang := pinnedAILANG(t)
	handler, err := NewModelHandler(ModelHandlerConfig{AILANGPath: ailang, Stub: true})
	if err != nil {
		t.Fatal(err)
	}
	var results [][]byte
	var records [][]byte
	for i := 0; i < 2; i++ {
		session, recording := handlerSession(t, EffectModelInfer, "model-scope", handler)
		result, _, invokeErr := session.Invoke(context.Background(),
			EffectRequest{Effect: EffectModelInfer, Scope: "model-scope", Cost: 2, Now: 1},
			[]byte("choose the next action"))
		if invokeErr != nil {
			t.Errorf("invoke %d error = %v, want success", i, invokeErr)
		}
		if len(recording.records) != 1 {
			t.Fatalf("invoke %d record count = %d, want 1", i, len(recording.records))
		}
		rec, decodeErr := DecodeRecord(recording.records[0].Payload)
		if decodeErr != nil {
			t.Fatal(decodeErr)
		}
		if !rec.Allowed || rec.Failed || !RecordConsistent(rec) {
			t.Errorf("invoke %d success record arm = %#v", i, rec)
		}
		results = append(results, result)
		records = append(records, recording.records[0].Payload)
	}
	if !bytes.Equal(results[0], results[1]) || !bytes.Equal(records[0], records[1]) {
		t.Fatalf("model stub differs: results %q / %q, records %q / %q",
			results[0], results[1], records[0], records[1])
	}
	if string(results[0]) != "{\"kind\":\"Wait\"}\n" {
		t.Fatalf("model stub result = %q, want deterministic Wait bytes", results[0])
	}
}

func TestGitHandlerTimeoutWritesFailureRecord(t *testing.T) {
	fake := writeExecutable(t, "sleep 5\nprintf 'finished\\n'")
	handler, err := NewGitHandler(GitHandlerConfig{
		GitPath: fake, ExecTimeout: 40 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	scope := t.TempDir()
	session, recording := handlerSession(t, EffectGitCommit, scope, handler)
	start := time.Now()
	_, ref, invokeErr := session.Invoke(context.Background(),
		EffectRequest{Effect: EffectGitCommit, Scope: scope, Cost: 2, Now: 1}, []byte("timeout"))
	elapsed := time.Since(start)
	// 2s tolerates process-spawn latency on a shared CI runner while still
	// reding hard if the deadline is ignored (the subprocess sleeps 5s).
	if elapsed > 2*time.Second {
		t.Errorf("Git timeout elapsed = %s, want <= 2s for a 40ms bound", elapsed)
	}
	assertHandlerFailureRecord(t, session, recording, ref, invokeErr, 2, ErrHandlerTimeout)
}

func TestModelHandlerTimeoutWritesFailureRecord(t *testing.T) {
	fake := writeExecutable(t, "sleep 5\nprintf 'finished\\n'")
	handler, err := NewModelHandler(ModelHandlerConfig{
		AILANGPath: fake, Stub: true, ExecTimeout: 40 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	session, recording := handlerSession(t, EffectModelInfer, "model-scope", handler)
	start := time.Now()
	_, ref, invokeErr := session.Invoke(context.Background(),
		EffectRequest{Effect: EffectModelInfer, Scope: "model-scope", Cost: 2, Now: 1}, []byte("timeout"))
	elapsed := time.Since(start)
	if elapsed > 2*time.Second {
		t.Errorf("Model timeout elapsed = %s, want <= 2s for a 40ms bound", elapsed)
	}
	assertHandlerFailureRecord(t, session, recording, ref, invokeErr, 2, ErrHandlerTimeout)
}

func TestGitHandlerOutputCapWritesFailureRecord(t *testing.T) {
	fake := writeExecutable(t, "i=0\nwhile [ \"$i\" -lt 256 ]; do printf 0123456789abcdef; i=$((i+1)); done")
	handler, err := NewGitHandler(GitHandlerConfig{
		GitPath: fake, MaxOutputBytes: 64,
	})
	if err != nil {
		t.Fatal(err)
	}
	scope := t.TempDir()
	session, recording := handlerSession(t, EffectGitCommit, scope, handler)
	_, ref, invokeErr := session.Invoke(context.Background(),
		EffectRequest{Effect: EffectGitCommit, Scope: scope, Cost: 2, Now: 1}, []byte("overflow"))
	assertHandlerFailureRecord(t, session, recording, ref, invokeErr, 2, ErrHandlerOverflow)
}

func TestModelHandlerOutputCapWritesFailureRecord(t *testing.T) {
	fake := writeExecutable(t, "i=0\nwhile [ \"$i\" -lt 256 ]; do printf 0123456789abcdef; i=$((i+1)); done")
	handler, err := NewModelHandler(ModelHandlerConfig{
		AILANGPath: fake, Stub: true, MaxOutputBytes: 64,
	})
	if err != nil {
		t.Fatal(err)
	}
	session, recording := handlerSession(t, EffectModelInfer, "model-scope", handler)
	_, ref, invokeErr := session.Invoke(context.Background(),
		EffectRequest{Effect: EffectModelInfer, Scope: "model-scope", Cost: 2, Now: 1}, []byte("overflow"))
	assertHandlerFailureRecord(t, session, recording, ref, invokeErr, 2, ErrHandlerOverflow)
}

func TestGitHandlerNonZeroExitWritesFailureRecord(t *testing.T) {
	fake := writeExecutable(t, "printf 'diagnostic\\n'\nexit 17")
	handler, err := NewGitHandler(GitHandlerConfig{GitPath: fake})
	if err != nil {
		t.Fatal(err)
	}
	scope := t.TempDir()
	session, recording := handlerSession(t, EffectGitCommit, scope, handler)
	_, ref, invokeErr := session.Invoke(context.Background(),
		EffectRequest{Effect: EffectGitCommit, Scope: scope, Cost: 2, Now: 1}, nil)
	var exitErr *HandlerExitError
	if !errors.As(invokeErr, &exitErr) {
		t.Fatalf("Invoke error = %v, want *HandlerExitError through EffectFailedError", invokeErr)
	}
	assertHandlerFailureRecord(t, session, recording, ref, invokeErr, 2, nil)
}

func TestGitHandlerMissingRepoWritesFailureRecord(t *testing.T) {
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	handler, err := NewGitHandler(GitHandlerConfig{GitPath: git})
	if err != nil {
		t.Fatal(err)
	}
	scope := filepath.Join(t.TempDir(), "missing")
	session, recording := handlerSession(t, EffectGitCommit, scope, handler)
	_, ref, invokeErr := session.Invoke(context.Background(),
		EffectRequest{Effect: EffectGitCommit, Scope: scope, Cost: 2, Now: 1}, nil)
	assertHandlerFailureRecord(t, session, recording, ref, invokeErr, 2, nil)
}

func TestRefusedFSPathWritesFailureRecord(t *testing.T) {
	scope := "relative.txt"
	session, recording := handlerSession(t, EffectFSRead, scope, FSHandler{})
	_, ref, invokeErr := session.Invoke(context.Background(),
		EffectRequest{Effect: EffectFSRead, Scope: scope, Cost: 2, Now: 1}, nil)
	var pathErr *FSPathError
	if !errors.As(invokeErr, &pathErr) {
		t.Fatalf("Invoke error = %v, want *FSPathError through EffectFailedError", invokeErr)
	}
	assertHandlerFailureRecord(t, session, recording, ref, invokeErr, 2, nil)
}

func approvalSession(t *testing.T) (*Session, *handlerRecordingStore, *HumanHandler) {
	t.Helper()
	recording := &handlerRecordingStore{base: openTestStore(t), deleted: make(map[hashref.HashRef]bool)}
	human := newHumanHandler(recording)
	session := newSession(recording, []Capability{
		{Effect: EffectHumanApprove, Scope: "release", ExpiresAt: 100, Budget: 5},
		{Effect: EffectHumanPollApproval, Scope: "release", ExpiresAt: 100, Budget: 4},
	}, Registry{
		EffectHumanApprove: human, EffectHumanPollApproval: human,
	}, Live, nil)
	return session, recording, human
}

func decodePendingRef(t *testing.T, payload []byte) hashref.HashRef {
	t.Helper()
	var pending pendingWire
	if err := decodeApprovalJSON(payload, &pending); err != nil {
		t.Fatal(err)
	}
	if pending.Status != "pending" {
		t.Fatalf("pending status = %q, want pending", pending.Status)
	}
	ref, err := hashref.Parse(pending.RequestRef)
	if err != nil {
		t.Fatal(err)
	}
	return ref
}

func TestApprovalFlowImmutableRecordAndSeparateDecision(t *testing.T) {
	session, recording, _ := approvalSession(t)
	approveReq := EffectRequest{
		Effect: EffectHumanApprove, Scope: "release", Cost: 2, Now: 10,
	}
	pending, approveRecordRef, err := session.Invoke(
		context.Background(), approveReq, mustApprovalJSON(approvalInputWire{Requester: "agent-7"}))
	if err != nil {
		t.Fatal(err)
	}
	requestRef := decodePendingRef(t, pending)
	if len(recording.records) != 1 {
		t.Fatalf("approve record count = %d, want exactly 1", len(recording.records))
	}
	approveRecordBefore := append([]byte(nil), recording.records[0].Payload...)
	if recording.records[0].Hash != approveRecordRef {
		t.Fatalf("approve record hash = %s, want %s", recording.records[0].Hash, approveRecordRef)
	}
	approveRecord, err := DecodeRecord(approveRecordBefore)
	if err != nil {
		t.Fatal(err)
	}
	resultObj, ok, err := recording.GetObject(approveRecord.ResultRef)
	if err != nil || !ok || !bytes.Equal(resultObj.Payload, pending) {
		t.Fatalf("Pending result object = ok %v payload %q err %v", ok, resultObj.Payload, err)
	}
	requestBefore, ok, err := recording.GetObject(requestRef)
	if err != nil || !ok || requestBefore.SemanticID != ApprovalRequestV1 {
		t.Fatalf("request object = %#v, ok %v, err %v", requestBefore, ok, err)
	}

	decisionRef, err := decideApproval(recording, requestRef, "approve", "operator", 11)
	if err != nil {
		t.Fatal(err)
	}
	decisionObj, ok, err := recording.GetObject(decisionRef)
	if err != nil || !ok || decisionObj.SemanticID != ApprovalDecisionV1 {
		t.Fatalf("decision object = %#v, ok %v, err %v", decisionObj, ok, err)
	}
	if len(recording.records) != 1 {
		t.Fatalf("record count after decision = %d, want unchanged at 1", len(recording.records))
	}
	headRef, ok, err := recording.GetRegistryHead(ApprovalsV1)
	if err != nil || !ok {
		t.Fatalf("approval head = %s, ok %v, err %v", headRef, ok, err)
	}
	headObj, ok, err := recording.GetObject(headRef)
	if err != nil || !ok {
		t.Fatalf("approval head object = ok %v, err %v", ok, err)
	}
	var head approvalHeadWire
	if err := decodeApprovalJSON(headObj.Payload, &head); err != nil {
		t.Fatal(err)
	}
	if head.DecisionRef != decisionRef.String() || head.RequestRef != requestRef.String() {
		t.Fatalf("moved head = %#v, want request and separate decision refs", head)
	}
	requestAfter, ok, err := recording.GetObject(requestRef)
	if err != nil || !ok || requestAfter.Hash != requestBefore.Hash ||
		!bytes.Equal(requestAfter.Payload, requestBefore.Payload) {
		t.Fatalf("approval request changed after decision")
	}
	approveRecordAfter, ok, err := recording.GetObject(approveRecordRef)
	if err != nil || !ok || approveRecordAfter.Hash != approveRecordRef ||
		!bytes.Equal(approveRecordAfter.Payload, approveRecordBefore) {
		t.Fatalf("approve effect record changed after decision")
	}

	pollReq := EffectRequest{
		Effect: EffectHumanPollApproval, Scope: "release", Cost: 1, Now: 12,
	}
	pollResult, _, err := session.Invoke(context.Background(), pollReq,
		mustApprovalJSON(approvalInputWire{RequestRef: requestRef.String()}))
	if err != nil {
		t.Fatal(err)
	}
	var observed observedDecisionWire
	if err := decodeApprovalJSON(pollResult, &observed); err != nil {
		t.Fatal(err)
	}
	if observed.Status != "decided" || !bytes.Equal(observed.Decision, decisionObj.Payload) {
		t.Fatalf("poll result = %q, want recorded decision observation", pollResult)
	}
	if len(recording.records) != 2 {
		t.Fatalf("record count after poll = %d, want 2", len(recording.records))
	}
	if got := session.grants[0].Budget; got != 3 {
		t.Fatalf("Human.Approve budget = %d, want 3 (debited at request time)", got)
	}
	if got := session.grants[1].Budget; got != 3 {
		t.Fatalf("Human.PollApproval budget = %d, want 3 (own debit)", got)
	}
	for _, obj := range recording.records {
		if got := hashref.SumSHA256(obj.Payload); got != obj.Hash {
			t.Errorf("effect record stored hash %s != hash(bytes) %s", obj.Hash, got)
		}
	}
}

func TestHumanApproveSynchronousPending(t *testing.T) {
	session, recording, _ := approvalSession(t)
	result, _, err := session.Invoke(context.Background(),
		EffectRequest{Effect: EffectHumanApprove, Scope: "release", Cost: 2, Now: 1},
		mustApprovalJSON(approvalInputWire{Requester: "agent"}))
	if err != nil {
		t.Fatal(err)
	}
	requestRef := decodePendingRef(t, result)
	if requestRef.IsZero() || len(recording.records) != 1 {
		t.Fatalf("Human.Approve request = %s, records = %d", requestRef, len(recording.records))
	}
	rec, err := DecodeRecord(recording.records[0].Payload)
	if err != nil {
		t.Fatal(err)
	}
	if !rec.Allowed || rec.Failed || rec.ResultRef.IsZero() || session.grants[0].Budget != 3 {
		t.Fatalf("synchronous Pending record = %#v, budget = %d", rec, session.grants[0].Budget)
	}
}

func TestApprovalRecordIntegritySweep(t *testing.T) {
	session, recording, _ := approvalSession(t)
	pending, _, err := session.Invoke(context.Background(),
		EffectRequest{Effect: EffectHumanApprove, Scope: "release", Cost: 2, Now: 1},
		mustApprovalJSON(approvalInputWire{Requester: "agent"}))
	if err != nil {
		t.Fatal(err)
	}
	requestRef := decodePendingRef(t, pending)
	if _, err := decideApproval(recording, requestRef, "approve", "operator", 2); err != nil {
		t.Fatal(err)
	}
	if _, _, err := session.Invoke(context.Background(),
		EffectRequest{Effect: EffectHumanPollApproval, Scope: "release", Cost: 1, Now: 3},
		mustApprovalJSON(approvalInputWire{RequestRef: requestRef.String()})); err != nil {
		t.Fatal(err)
	}
	if len(recording.records) != 2 {
		t.Fatalf("effect record count = %d, want approve and poll only", len(recording.records))
	}
	for _, obj := range recording.records {
		if got := hashref.SumSHA256(obj.Payload); got != obj.Hash {
			t.Errorf("stored effect-record hash %s != hash(bytes) %s", obj.Hash, got)
		}
	}
}

func TestDecideApprovalBeforeRequestRejected(t *testing.T) {
	_, recording, _ := approvalSession(t)
	ref := hashref.SumSHA256([]byte("not-a-request"))
	_, err := decideApproval(recording, ref, "approve", "operator", 1)
	if !errors.Is(err, ErrApprovalRequestNotFound) {
		t.Fatalf("DecideApproval error = %v, want ErrApprovalRequestNotFound", err)
	}
}

func TestPollApprovalDeniedWithoutOwnGrant(t *testing.T) {
	recording := &handlerRecordingStore{base: openTestStore(t), deleted: make(map[hashref.HashRef]bool)}
	human := newHumanHandler(recording)
	session := newSession(recording, []Capability{
		{Effect: EffectHumanApprove, Scope: "release", ExpiresAt: 100, Budget: 5},
	}, Registry{EffectHumanApprove: human, EffectHumanPollApproval: human}, Live, nil)
	pending, _, err := session.Invoke(context.Background(),
		EffectRequest{Effect: EffectHumanApprove, Scope: "release", Cost: 2, Now: 1},
		mustApprovalJSON(approvalInputWire{Requester: "agent"}))
	if err != nil {
		t.Fatal(err)
	}
	requestRef := decodePendingRef(t, pending)
	_, _, err = session.Invoke(context.Background(),
		EffectRequest{Effect: EffectHumanPollApproval, Scope: "release", Cost: 1, Now: 2},
		mustApprovalJSON(approvalInputWire{RequestRef: requestRef.String()}))
	var denial *DenialError
	if !errors.As(err, &denial) || denial.Decision.Label != LabelDeniedEffectName {
		t.Fatalf("poll error = %v, want denied effect-name without own grant", err)
	}
	if got := session.grants[0].Budget; got != 3 {
		t.Fatalf("approve budget after denied poll = %d, want 3", got)
	}
	if len(recording.records) != 2 {
		t.Fatalf("record count = %d, want approve plus denied poll", len(recording.records))
	}
}

type failingApprovalStore struct {
	*handlerRecordingStore
	failSemantic string
	failHead     bool
	failed       bool
}

var errApprovalStoreInjected = errors.New("injected approval store failure")

func (s *failingApprovalStore) PutObject(obj store.Object) error {
	if !s.failed && obj.SemanticID == s.failSemantic {
		s.failed = true
		return errApprovalStoreInjected
	}
	return s.handlerRecordingStore.PutObject(obj)
}

func (s *failingApprovalStore) SetRegistryHead(name string, ref hashref.HashRef) error {
	if !s.failed && s.failHead && name == ApprovalsV1 {
		s.failed = true
		return errApprovalStoreInjected
	}
	return s.handlerRecordingStore.SetRegistryHead(name, ref)
}

func TestApprovalFailuresKeepStandingAttentionDebit(t *testing.T) {
	for _, tc := range []struct {
		name         string
		failSemantic string
		failHead     bool
	}{
		{name: "request PutObject", failSemantic: ApprovalRequestV1},
		{name: "head move", failHead: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			base := &handlerRecordingStore{base: openTestStore(t), deleted: make(map[hashref.HashRef]bool)}
			failing := &failingApprovalStore{
				handlerRecordingStore: base, failSemantic: tc.failSemantic, failHead: tc.failHead,
			}
			human := newHumanHandler(failing)
			session := newSession(failing, []Capability{{
				Effect: EffectHumanApprove, Scope: "release", ExpiresAt: 100, Budget: 5,
			}}, Registry{EffectHumanApprove: human}, Live, nil)
			_, ref, invokeErr := session.Invoke(context.Background(),
				EffectRequest{Effect: EffectHumanApprove, Scope: "release", Cost: 2, Now: 1},
				mustApprovalJSON(approvalInputWire{Requester: "agent"}))
			var failed *EffectFailedError
			if !errors.As(invokeErr, &failed) || !errors.Is(invokeErr, errApprovalStoreInjected) {
				t.Fatalf("Invoke error = %v, want *EffectFailedError wrapping injected error", invokeErr)
			}
			assertHandlerFailureRecord(t, session, base, ref, invokeErr, 2, errApprovalStoreInjected)
		})
	}
}

func TestApprovalReplayContract(t *testing.T) {
	session, recording, _ := approvalSession(t)
	approveReq := EffectRequest{Effect: EffectHumanApprove, Scope: "release", Cost: 2, Now: 10}
	approvePayload := mustApprovalJSON(approvalInputWire{Requester: "agent"})
	pending, approveRecordRef, err := session.Invoke(context.Background(), approveReq, approvePayload)
	if err != nil {
		t.Fatal(err)
	}
	requestRef := decodePendingRef(t, pending)
	decisionRef, err := decideApproval(recording, requestRef, "deny", "operator", 11)
	if err != nil {
		t.Fatal(err)
	}
	pollReq := EffectRequest{Effect: EffectHumanPollApproval, Scope: "release", Cost: 1, Now: 12}
	pollPayload := mustApprovalJSON(approvalInputWire{RequestRef: requestRef.String()})
	observation, pollRecordRef, err := session.Invoke(context.Background(), pollReq, pollPayload)
	if err != nil {
		t.Fatal(err)
	}

	// Delete the decision object from this store view. Replay must use only the
	// recorded poll result object and therefore remain byte-identical.
	recording.deleted[decisionRef] = true
	dispatches := 0
	stub := HandlerFunc(func(context.Context, EffectRequest, []byte) ([]byte, error) {
		dispatches++
		return nil, errors.New("replay dispatched a handler")
	})
	replay := newSession(recording, []Capability{
		{Effect: EffectHumanApprove, Scope: "release", ExpiresAt: 100, Budget: 5},
		{Effect: EffectHumanPollApproval, Scope: "release", ExpiresAt: 100, Budget: 4},
	}, Registry{EffectHumanApprove: stub, EffectHumanPollApproval: stub}, Replay,
		[]hashref.HashRef{approveRecordRef, pollRecordRef})
	replayedPending, gotApproveRef, err := replay.Invoke(context.Background(), approveReq, approvePayload)
	if err != nil {
		t.Fatal(err)
	}
	replayedObservation, gotPollRef, err := replay.Invoke(context.Background(), pollReq, pollPayload)
	if err != nil {
		t.Fatal(err)
	}
	if dispatches != 0 {
		t.Fatalf("replay handler dispatches = %d, want 0", dispatches)
	}
	if gotApproveRef != approveRecordRef || gotPollRef != pollRecordRef ||
		!bytes.Equal(replayedPending, pending) || !bytes.Equal(replayedObservation, observation) {
		t.Fatalf("approval replay was not byte-identical")
	}
}

// The wall-clock bound must kill the whole process TREE, not just the handler's
// direct child. A forked grandchild inherits the stdout pipe, so if it survives
// the kill the capped read blocks until the GRANDCHILD exits and the bound is not
// enforced. Linux CI measured 5.002s against a 40ms bound — exactly the runtime of
// the grandchild — while darwin reported 42ms for the same code, so this is a real
// guarantee that one platform hides.
//
// Only elapsed is asserted. Checking that the grandchild is DEAD afterwards is
// vacuous: when the kill misses it, Invoke blocks on the inherited pipe until the
// grandchild exits on its own, so by the time a test could look it has always died.
func TestHandlerTimeoutKillsTheWholeProcessGroup(t *testing.T) {
	fake := writeExecutable(t, "sleep 5 &\nwait")
	handler, err := NewGitHandler(GitHandlerConfig{
		GitPath: fake, ExecTimeout: 100 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	scope := t.TempDir()
	session, _ := handlerSession(t, EffectGitCommit, scope, handler)
	start := time.Now()
	_, _, _ = session.Invoke(context.Background(),
		EffectRequest{Effect: EffectGitCommit, Scope: scope, Cost: 2, Now: 1}, []byte("group"))
	elapsed := time.Since(start)
	if elapsed > 2*time.Second {
		t.Errorf("Invoke took %s for a 100ms bound: the kill missed the forked "+
			"grandchild, which held the stdout pipe until it exited on its own", elapsed)
	}
}
