package broker

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// ---------------------------------------------------------------------------
// SM.D0 / AC25 / AC26 — the attended publish operation, driven end to end
//
// SAFETY, restated because this is the file that drives a real publish path:
// every handler here is built by newLoopbackRegistryPublishHandler, which is
// UNEXPORTED and refuses a non-loopback origin; every validator is an httptest
// server bound to 127.0.0.1; and the publisher subprocess is this test binary
// re-exec'd into TestRegistryPublishHelperProcess, never the real `ailang`.
// world/core@0.1.0 is not published by anything in this file.
//
// THIS IS THE ONLY PACKAGE THAT CAN DRIVE THE HAPPY PATH. cmd/world-publish's
// tests can reach the PRODUCTION constructor only, and only for its refusals —
// which is the point of D2 and the reason this file exists here rather than
// next to the command.
// ---------------------------------------------------------------------------

// attendedPlanFor builds the plan from a landed publish fixture. The times are
// the fixture's own defaults so the ordering refusals in validatePublishApproval
// are exercised in their satisfied direction rather than side-stepped.
func attendedPlanFor(f publishFixture, episode string) AttendedPublishPlan {
	times := defaultApprovalTimes()
	id := f.identity
	// The plan carries NO approval ref: it does not exist until the mint
	// returns it, and Payload is the only thing that binds one.
	id.ApprovalRef = hashref.HashRef{}
	return AttendedPublishPlan{
		Identity:    id,
		Hashes:      f.hashes,
		Requester:   "smd0-operator",
		DecidedBy:   "smd0-attended-operator",
		EpisodeID:   episode,
		RequestedAt: times.request,
		DecidedAt:   times.decide,
		PublishAt:   50,
		ExpiresAt:   times.expires,
	}
}

// attendedHandlerFor builds the loopback handler for a plan whose approval has
// already been minted. The approval stamp is defence in depth (the authority
// comes from the landed traversal), but the handler refuses to exist without
// one, so the minted ref has to reach it.
func attendedHandlerFor(
	t *testing.T, f publishFixture, validator *fakeValidator, mode string, approvalRef hashref.HashRef,
) *RegistryPublishHandler {
	t.Helper()
	approval := f.approval
	approval.ApprovalRef = approvalRef
	return loopbackHandler(t, RegistryPublishConfig{
		PublisherPath:   writePublisherScript(t, mode, nil),
		PackageDir:      f.dir,
		Manifest:        f.manifest,
		RegistryOrigin:  validator.origin(),
		ValidatorOrigin: validator.origin(),
		Credential:      writeCredentialFile(t, ac10Sentinel),
		Approval:        approval,
		ExecTimeout:     20 * time.Second,
	})
}

// decodeRecordAt reads an immutable EffectRecord back out of the store. The
// budget assertions below are stated against THIS object rather than against
// the Session's in-memory grant slice, because the record is what a reviewer,
// a replay and a recovery scan all read — and because a criterion satisfied by
// an in-memory counter is the exact failure MUT-SM-CLAIM-MEMORY was written for.
func decodeRecordAt(t *testing.T, s *store.Store, ref hashref.HashRef) EffectRecord {
	t.Helper()
	if ref.IsZero() {
		t.Fatal("instrument failure: asked to decode the ZERO record ref")
	}
	obj, ok, err := s.GetObject(context.Background(), ref)
	if err != nil || !ok {
		t.Fatalf("read effect record %s: ok=%v err=%v", ref, ok, err)
	}
	if obj.SemanticID != EffectRecordV1 {
		t.Fatalf("object %s is a %q, want %q", ref, obj.SemanticID, EffectRecordV1)
	}
	rec, err := DecodeRecord(obj.Payload)
	if err != nil {
		t.Fatalf("decode effect record %s: %v", ref, err)
	}
	return rec
}

// ---------------------------------------------------------------------------
// AC25 — mint through the LANDED traversal, spend EXACTLY ONCE, and prove the
// claim is durable across a store close/reopen with a FRESH budget
// ---------------------------------------------------------------------------

func TestAttendedPublishMintsThroughTheLandedTraversalAndSpendsExactlyOnce(t *testing.T) {
	validator := newFakeValidator(t, "ok")
	base := openPublishStore(t)
	fixture := newPublishFixture(t, validator.origin(), "smd0-ac25")
	plan := attendedPlanFor(fixture, "smd0-ac25")

	// --- the mint -----------------------------------------------------------
	approvalRef, err := mintAttendedApproval(base, plan)
	if err != nil {
		t.Fatalf("MintAttendedApproval: %v", err)
	}
	if approvalRef.IsZero() {
		t.Fatal("MintAttendedApproval returned the zero ref")
	}
	// The minted ref must name a real ApprovalDecisionV1 that says "approve".
	// Without this the ref could be any digest and every assertion below would
	// be about a number rather than about an object.
	decisionObj, ok, err := base.GetObject(context.Background(), approvalRef)
	if err != nil || !ok {
		t.Fatalf("minted approval %s names no object: ok=%v err=%v", approvalRef, ok, err)
	}
	if decisionObj.SemanticID != ApprovalDecisionV1 {
		t.Fatalf("minted approval is a %q object, want %q", decisionObj.SemanticID, ApprovalDecisionV1)
	}
	var decision approvalDecisionWire
	if err := decodeApprovalJSON(decisionObj.Payload, &decision); err != nil {
		t.Fatal(err)
	}
	if decision.Decision != "approve" || decision.DecidedBy != plan.DecidedBy {
		t.Fatalf("minted decision = %q by %q, want \"approve\" by %q",
			decision.Decision, decision.DecidedBy, plan.DecidedBy)
	}
	// The decision's request must carry the plan's canonical publish-bound
	// scope. This is what makes the stamp bind BYTES rather than a name.
	requestRef, err := hashref.Parse(decision.RequestRef)
	if err != nil {
		t.Fatal(err)
	}
	requestObj, ok, err := base.GetObject(context.Background(), requestRef)
	if err != nil || !ok {
		t.Fatalf("minted decision references request %s which is absent", requestRef)
	}
	var request approvalRequestWire
	if err := decodeApprovalJSON(requestObj.Payload, &request); err != nil {
		t.Fatal(err)
	}
	if request.Scope != plan.ApprovalScope() {
		t.Fatalf("minted request scope =\n  %q\nwant\n  %q", request.Scope, plan.ApprovalScope())
	}
	// The mint is not a publish: nothing may have reached the validator yet.
	if got := validator.count(); got != 0 {
		t.Fatalf("the MINT produced %d validator request(s), want 0", got)
	}

	// --- the spend ----------------------------------------------------------
	handler := attendedHandlerFor(t, fixture, validator, "success", approvalRef)
	recording := &publishRecordingStore{base: base.Store}

	result, err := invokeAttendedPublish(context.Background(), recording, handler, plan, approvalRef)
	if err != nil {
		t.Fatalf("InvokeAttendedPublish: %v (counters %s)", err, readPublishCounters(validator, handler))
	}
	if result.Status != AttendedPublishSucceeded {
		t.Fatalf("result status = %q, want %q", result.Status, AttendedPublishSucceeded)
	}
	if result.Reconcilable() {
		t.Fatal("a succeeded publish reports itself reconcilable")
	}

	afterGrant := readPublishCounters(validator, handler)
	t.Logf("AC25 counters after the attended publish: %s", afterGrant)
	if afterGrant != (publishCounters{posts: 1, dispatches: 1, credentialLoads: 1}) {
		t.Fatalf("counters after the attended publish = %s, want POST=1 dispatches=1 credentialLoads=1",
			afterGrant)
	}

	// Budget 1 -> 0, read off the IMMUTABLE RECORD. MUT-D0-BUDGET-2 raises the
	// grant to 2 and reds precisely here: the record then reads 2 -> 1.
	rec := decodeRecordAt(t, base.Store, result.RecordRef)
	if rec.BudgetBefore != PublishCost || rec.BudgetAfter != 0 {
		t.Fatalf("record budget = %d -> %d, want %d -> 0: the attended grant is ONE irreversible attempt",
			rec.BudgetBefore, rec.BudgetAfter, PublishCost)
	}
	if !rec.Allowed || rec.Failed || rec.Denial != "" {
		t.Fatalf("record = allowed:%v failed:%v denial:%q, want an allowed success",
			rec.Allowed, rec.Failed, rec.Denial)
	}
	if rec.Scope != plan.PublishScope() {
		t.Fatalf("record scope = %q, want %q", rec.Scope, plan.PublishScope())
	}

	// The durable intent resolved, and a succeeded outcome was appended.
	if len(recording.effectIDs) != 1 {
		t.Fatalf("durable effect intents = %v, want exactly 1", recording.effectIDs)
	}
	receipt, ok, err := base.GetEffectReceipt(recording.effectIDs[0])
	if err != nil || !ok {
		t.Fatalf("effect receipt %s: ok=%v err=%v", recording.effectIDs[0], ok, err)
	}
	if receipt.State != store.ReceiptResolved {
		t.Fatalf("receipt state = %q, want %q", receipt.State, store.ReceiptResolved)
	}
	if receipt.EffectOutcome == nil || receipt.EffectOutcome.Status != "succeeded" {
		t.Fatalf("receipt outcome = %+v, want a succeeded outcome", receipt.EffectOutcome)
	}

	// --- the durable claim, across a CLOSE AND REOPEN -----------------------
	//
	// This half is what distinguishes a durable claim from an in-memory budget.
	// invokeAttendedPublish constructs a FRESH Capability with a FRESH budget on
	// every call, so nothing in the second attempt is refused because a counter
	// was already spent — the refusal has to come off the disk.
	reopened := base.reopen(t)
	if _, found, err := reopened.GetObject(context.Background(), approvalRef); err != nil || !found {
		t.Fatalf("the reopened store cannot read the minted approval: found=%v err=%v", found, err)
	}

	replan := plan
	replan.EpisodeID = "smd0-ac25-reopened"
	replan.PublishAt = 51
	_, reuseErr := invokeAttendedPublish(context.Background(), reopened.Store, handler, replan, approvalRef)
	if !errors.Is(reuseErr, store.ErrApprovalAlreadyConsumed) {
		t.Fatalf("reuse after reopen = %T %v, want store.ErrApprovalAlreadyConsumed", reuseErr, reuseErr)
	}
	var denied *DenialError
	if errors.As(reuseErr, &denied) {
		t.Fatalf("the reuse was refused by BUDGET (%s), not by the durable claim", denied.Decision.Label)
	}
	afterReuse := readPublishCounters(validator, handler)
	t.Logf("AC25 counters after the reopened reuse: %s", afterReuse)
	if afterReuse != afterGrant {
		t.Fatalf("counters moved across the refused reuse: %s -> %s; the refusal must land BEFORE "+
			"credential load and BEFORE any POST", afterGrant, afterReuse)
	}
	if got := validator.count(); got != 1 {
		t.Fatalf("total validator POST count = %d, want exactly 1 across BOTH attempts", got)
	}
}

// TestObserveMintedDecisionDrivesLegThreeInBothDirections covers the refusal
// mintAttendedApproval adds that none of the eight inherited refusal families
// catches: the landed poll leg failing to observe back the decision the mint
// just made.
//
// The POSITIVE arm is produced by a REAL mint, so "the check passes" is a
// statement about the bytes Human.PollApproval actually returns rather than
// about a fixture shaped to satisfy it.
func TestObserveMintedDecisionDrivesLegThreeInBothDirections(t *testing.T) {
	base := openPublishStore(t)
	fixture := newPublishFixture(t, "http://127.0.0.1:1", "smd0-mint-legs")
	plan := attendedPlanFor(fixture, "smd0-mint-legs")

	decisionRef, err := mintAttendedApproval(base, plan)
	if err != nil {
		t.Fatalf("instrument failure: the control mint failed: %v", err)
	}
	decisionObj, ok, err := base.GetObject(context.Background(), decisionRef)
	if err != nil || !ok {
		t.Fatalf("minted decision %s is absent: ok=%v err=%v", decisionRef, ok, err)
	}
	// The positive control: the exact wire shape the landed poll returns.
	polled := mustApprovalJSON(observedDecisionWire{
		Status: "decided", Decision: decisionObj.Payload,
	})
	if err := observeMintedDecision(polled, decisionRef); err != nil {
		t.Fatalf("observeMintedDecision rejected a genuinely minted decision: %v", err)
	}

	for _, arm := range []struct {
		name string
		body []byte
		ref  hashref.HashRef
	}{
		{
			// The landed poll's own pending shape carries `requestRef`, which
			// observedDecisionWire does not declare, so it is refused by the
			// strict decoder BEFORE the status check. Reaching the status branch
			// therefore needs the decision-shaped wire with a pending status —
			// which is exactly what a poll would return if the head walk found
			// the request but no decision and the wire ever gained that field.
			"decision-shaped but still pending",
			mustApprovalJSON(observedDecisionWire{Status: "pending"}),
			decisionRef,
		},
		{
			"the landed pending wire (refused by the strict decoder)",
			pendingBytes(decisionRef),
			decisionRef,
		},
		{
			"a decision that hashes elsewhere",
			polled,
			hashref.SumSHA256([]byte("some other decision entirely")),
		},
		{
			"undecodable",
			[]byte("{not json"),
			decisionRef,
		},
	} {
		t.Run(arm.name, func(t *testing.T) {
			err := observeMintedDecision(arm.body, arm.ref)
			if err == nil {
				t.Fatal("observeMintedDecision accepted it")
			}
			t.Logf("refusal: %v", err)
		})
	}

	// The two content-bearing arms must be the TYPED refusal, so a caller can
	// tell "the traversal did not confirm the mint" from "the store errored".
	if !errors.Is(observeMintedDecision(mustApprovalJSON(observedDecisionWire{Status: "pending"}), decisionRef),
		ErrAttendedApprovalNotObserved) {
		t.Fatal("a pending poll is not reported as ErrAttendedApprovalNotObserved")
	}
	if !errors.Is(observeMintedDecision(polled, hashref.SumSHA256([]byte("other"))),
		ErrAttendedApprovalNotObserved) {
		t.Fatal("a mismatched decision hash is not reported as ErrAttendedApprovalNotObserved")
	}
}

// ---------------------------------------------------------------------------
// AC26 — a typed-ambiguous dispatch appends NO outcome, is NOT retried, and is
// resolved read-only by reconciliation
// ---------------------------------------------------------------------------

func TestIndeterminatePublishAppendsNoOutcomeAndIsNeverRetried(t *testing.T) {
	// "reset" logs the request and THEN kills the connection: the request body
	// demonstrably left the publisher, which is exactly Decision 3's ambiguous
	// class rather than a definite failure.
	validator := newFakeValidator(t, "reset")
	base := openPublishStore(t)
	fixture := newPublishFixture(t, validator.origin(), "smd0-ac26")
	plan := attendedPlanFor(fixture, "smd0-ac26")

	approvalRef, err := mintAttendedApproval(base, plan)
	if err != nil {
		t.Fatalf("MintAttendedApproval: %v", err)
	}
	handler := attendedHandlerFor(t, fixture, validator, "reset", approvalRef)
	recording := &publishRecordingStore{base: base.Store}

	result, err := invokeAttendedPublish(context.Background(), recording, handler, plan, approvalRef)

	// The typed error, carrying the three fields a reconciliation pass needs.
	var indeterminate *IndeterminateEffectError
	if !errors.As(err, &indeterminate) {
		t.Fatalf("InvokeAttendedPublish error = %T %v, want *IndeterminateEffectError "+
			"(counters %s)", err, err, readPublishCounters(validator, handler))
	}
	if result.Status != AttendedPublishIndeterminate || !result.Reconcilable() {
		t.Fatalf("result status = %q (reconcilable %v), want %q",
			result.Status, result.Reconcilable(), AttendedPublishIndeterminate)
	}
	if result.InvocationID == "" || result.EpisodeID != plan.EpisodeID {
		t.Fatalf("result names invocation %q episode %q, want a non-empty invocation in episode %q",
			result.InvocationID, result.EpisodeID, plan.EpisodeID)
	}
	// The three fields must be MUTUALLY CONSISTENT, not merely non-empty: the
	// store derives the invocation ID from the episode and the ordinal, so
	// re-deriving it here catches a result that carried three fields copied out
	// of three different places. (Ordinals are 0-based; the first effect in an
	// episode is ordinal 0, MEASURED — store/journal.go:511 mints maxOrdinal+1
	// from a floor of -1 — so a `> 0` assertion here would be wrong, not strict.)
	if want := store.EffectInvocationID(result.EpisodeID, result.Ordinal); result.InvocationID != want {
		t.Fatalf("result names invocation %q but episode %q ordinal %d derives %q",
			result.InvocationID, result.EpisodeID, result.Ordinal, want)
	}
	if !result.RecordRef.IsZero() {
		t.Fatalf("an indeterminate attempt named an effect record %s; it must append none", result.RecordRef)
	}

	// RULE 3i. The observable is the FAKE VALIDATOR'S OWN request counter, which
	// is downstream of the dispatch and lives in another process's collaborator.
	// It is deliberately NOT handler.Dispatches(), which is incremented in the
	// same function as the dispatch it counts and therefore cannot fail for the
	// reason this criterion claims.
	afterAttempt := validator.count()
	t.Logf("AC26 validator POST count after the ambiguous attempt: %d", afterAttempt)
	if afterAttempt != 1 {
		t.Fatalf("validator saw %d request(s) after ONE ambiguous attempt, want exactly 1: "+
			"a retry here is the double-publish this design exists to prevent", afterAttempt)
	}

	// No EffectRecord and no outcome: the durable intent stays INDETERMINATE.
	if len(recording.records) != 0 {
		t.Fatalf("an indeterminate attempt appended %d effect record(s), want 0", len(recording.records))
	}
	if len(recording.effectIDs) != 1 {
		t.Fatalf("durable effect intents = %v, want exactly 1", recording.effectIDs)
	}
	receipt, ok, err := base.GetEffectReceipt(recording.effectIDs[0])
	if err != nil || !ok {
		t.Fatalf("effect receipt %s: ok=%v err=%v", recording.effectIDs[0], ok, err)
	}
	if receipt.State != store.ReceiptIndeterminate {
		t.Fatalf("receipt state = %q, want %q", receipt.State, store.ReceiptIndeterminate)
	}
	if receipt.EffectOutcome != nil {
		t.Fatalf("an indeterminate attempt appended an outcome: %+v", receipt.EffectOutcome)
	}

	// --- reconciliation resolves the receipt, READ-ONLY ---------------------
	//
	// Reconciliation reads the BUCKET, never the validator service. The two
	// origins are asserted to differ so "the POST count did not move" is a
	// statement about a server that reconciliation never addresses.
	bucket := newFakeBucket(t)
	if bucket.origin() == validator.origin() {
		t.Fatal("instrument failure: the bucket fake and the validator fake share an origin")
	}
	firingControl(bucket)
	bucket.put(targetKey(), bucketResponse{
		status: http.StatusOK,
		body:   metadataDocument(fixture.identity.Vendor+"/"+fixture.identity.Name, fixtureVersion, fixture.hashes),
	})
	cfg := reconcileCfg(bucket)
	cfg.Expected = fixture.hashes
	cfg.ObservedPublishStatus = 0
	reconciled := mustReconcile(t, cfg)
	t.Logf("AC26 reconciliation: %s", reconciled)
	if reconciled.State != ReconcileSucceededReconciled {
		t.Fatalf("reconciliation state = %q, want %q", reconciled.State, ReconcileSucceededReconciled)
	}

	// The whole point: reconciliation is read-only, so the irreversible POST
	// count is STILL exactly one across the ambiguous attempt AND its resolution.
	if got := validator.count(); got != 1 {
		t.Fatalf("validator POST count after reconciliation = %d, want exactly 1", got)
	}
	if got := handler.Dispatches(); got != 1 {
		t.Fatalf("publisher dispatches = %d, want exactly 1", got)
	}
	// Every reconciliation request must have been a GET against the bucket.
	requests := bucket.requests()
	if len(requests) == 0 {
		t.Fatal("instrument failure: reconciliation issued ZERO bucket requests")
	}
	for _, r := range requests {
		if r.Method != http.MethodGet {
			t.Fatalf("reconciliation issued a %s to %s; its only network verb is GET", r.Method, r.Path)
		}
	}
}

// TestSessionLevelRetryOfAnIndeterminatePublishIsRefusedByTheDurableClaim
// records a MEASURED second layer, so a reviewer does not read AC26's "no
// retry" as resting on one guard.
//
// If InvokeAttendedPublish were mutated to retry by calling Session.Invoke
// again, the retry would NOT reach the validator: validatePublishApproval and
// AppendClaimedEffectIntent refuse the already-consumed stamp before the
// handler runs. The dangerous retry — and therefore the one
// MUT-D0-INDETERMINATE-RETRY performs — is a retry at the HANDLER layer, which
// bypasses the claim entirely. This test measures the session-layer half so the
// distinction is on the record rather than in a comment.
func TestSessionLevelRetryOfAnIndeterminatePublishIsRefusedByTheDurableClaim(t *testing.T) {
	validator := newFakeValidator(t, "reset")
	base := openPublishStore(t)
	fixture := newPublishFixture(t, validator.origin(), "smd0-retry-layers")
	plan := attendedPlanFor(fixture, "smd0-retry-layers")

	approvalRef, err := mintAttendedApproval(base, plan)
	if err != nil {
		t.Fatal(err)
	}
	handler := attendedHandlerFor(t, fixture, validator, "reset", approvalRef)

	if _, err := invokeAttendedPublish(context.Background(), base.Store, handler, plan, approvalRef); err == nil {
		t.Fatal("the ambiguous arm returned no error")
	}
	if got := validator.count(); got != 1 {
		t.Fatalf("validator POST count after the first attempt = %d, want 1", got)
	}

	// The session-layer retry a careless fix would write.
	retryPlan := plan
	retryPlan.EpisodeID = "smd0-retry-layers-again"
	retryPlan.PublishAt = 51
	_, retryErr := invokeAttendedPublish(context.Background(), base.Store, handler, retryPlan, approvalRef)
	if !errors.Is(retryErr, store.ErrApprovalAlreadyConsumed) {
		t.Fatalf("session-layer retry error = %T %v, want store.ErrApprovalAlreadyConsumed", retryErr, retryErr)
	}
	if got := validator.count(); got != 1 {
		t.Fatalf("a session-layer retry reached the validator: POST count = %d, want 1", got)
	}

	// And the handler-layer retry DOES reach it — which is why
	// MUT-D0-INDETERMINATE-RETRY is written at that layer. Driving it here once,
	// deliberately, is the known-positive control that proves the counter above
	// can still move; without it "the count stayed at 1" is also what a dead
	// validator would report.
	if _, err := handler.Execute(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: plan.PublishScope(), Cost: PublishCost, Now: 52,
	}, plan.Payload(approvalRef)); err == nil {
		t.Fatal("the handler-layer control returned no error from the reset validator")
	}
	if got := validator.count(); got != 2 {
		t.Fatalf("known-positive control: validator POST count = %d, want 2 — the counter is stuck "+
			"and every zero-delta assertion above is void", got)
	}
	t.Log("MEASURED: a session-layer retry is refused by the durable claim (POST stays 1); " +
		"a handler-layer retry reaches the validator (POST 1 -> 2). MUT-D0-INDETERMINATE-RETRY " +
		"is therefore written at the handler layer.")
}

// TestAttendedPublishPlanRendersTheFrozenScopes pins the two strings the mint
// and the spend must agree on. They are rendered by production formatters, so
// this asserts the PLAN wires them together rather than re-deriving either.
func TestAttendedPublishPlanRendersTheFrozenScopes(t *testing.T) {
	base := openPublishStore(t)
	_ = base
	fixture := newPublishFixture(t, "http://127.0.0.1:1", "smd0-scopes")
	plan := attendedPlanFor(fixture, "smd0-scopes")

	if got, want := plan.PublishScope(), fixture.scope; got != want {
		t.Fatalf("plan publish scope = %q, want %q", got, want)
	}
	want := PublishApprovalScopeFor(fixture.identity, fixture.hashes, plan.ExpiresAt)
	// The fixture's identity carries a synthetic approval ref; the plan's does
	// not. Neither is an input to the approval scope, so the two must agree.
	if got := plan.ApprovalScope(); got != want {
		t.Fatalf("plan approval scope =\n  %q\nwant\n  %q", got, want)
	}
	parsed, err := ParsePublishApprovalScope(plan.ApprovalScope())
	if err != nil {
		t.Fatalf("the plan's approval scope does not parse: %v", err)
	}
	if parsed.Effect != EffectRegistryPublish {
		t.Fatalf("plan approval scope stamps effect %q, want %q", parsed.Effect, EffectRegistryPublish)
	}
	for _, arm := range []struct{ name, got, want string }{
		{"tarball", parsed.TarballSHA256, fixture.hashes.TarballSHA256},
		{"content", parsed.ContentHash, fixture.hashes.ContentHash},
		{"interface", parsed.InterfaceHash, fixture.hashes.InterfaceHash},
	} {
		if arm.got != arm.want {
			t.Errorf("plan approval scope %s hash = %s, want %s", arm.name, arm.got, arm.want)
		}
	}

	// The payload the plan renders must carry the approval ref it is handed and
	// nothing else from the fixture's synthetic one.
	ref := hashref.SumSHA256([]byte("smd0-scope-probe"))
	id, hashes, err := DecodePublishPayload(plan.Payload(ref))
	if err != nil {
		t.Fatal(err)
	}
	if id.ApprovalRef != ref {
		t.Fatalf("payload approval ref = %s, want %s", id.ApprovalRef, ref)
	}
	if hashes != fixture.hashes {
		t.Fatalf("payload hashes = %+v, want %+v", hashes, fixture.hashes)
	}
	if !strings.Contains(plan.ApprovalScope(), plan.PublishScope()) {
		t.Fatal("the approval scope does not carry the publish scope as its head")
	}
}
