package broker

import (
	"context"
	"errors"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// ---------------------------------------------------------------------------
// SM.D0 — the attended publish OPERATION, as an artifact
//
// Everything this file calls was already landed by SM.A-SM.C. What did NOT
// exist was a caller: `grep -rhoE 'Publish|Approve' cmd/ --include='*.go'`
// measured ZERO at 6d1dce0, against a known-positive control of 27 `func `
// declarations, and the EXPORTED operator entry point DecideApproval had
// exactly one caller in the whole repository — a test. So the runbook's
// "mint the one-shot approval" and "invoke the publish effect exactly once"
// named Go functions with no command-line surface, and the attended publish
// could not be performed for want of an ARTIFACT, not for want of a decision.
//
// WHY THIS LIVES IN host/broker AND NOT IN cmd/world-publish
//
// newLoopbackRegistryPublishHandler is UNEXPORTED. No package outside
// host/broker can point a publisher at an httptest server. If this wiring lived
// in cmd/, its happy path could only ever be driven against a FAKE Handler, and
// this repository would ship a publish caller that had never once been driven
// through the real RegistryPublishHandler. The alternative — exporting a
// loopback door so cmd/ could reach it — would widen the production surface in
// order to test the fence guarding it. Neither is acceptable. Putting the
// wiring here is the third option: publish_op_test.go drives the ENTIRE path
// (mint, capability, Invoke, dispatch, outcome) through the real handler
// against a real loopback validator, and cmd/world-publish stays thin — flags,
// fences, and the PRODUCTION constructor.
//
// THERE IS NO AUTHORITY LOGIC HERE. Every refusal is inherited:
//
//   - origin validation                -> validatePublishOrigin
//   - ambient credential               -> AssertNoAmbientRegistryCredential
//   - the seven approval refusals      -> validatePublishApproval
//   - single use, durably              -> store.AppendClaimedEffectIntent
//   - effect/scope/expiry/budget       -> Decide
//   - typed-indeterminate suppression  -> Session.Invoke
//
// If a refusal appears below that cannot name something none of those catch,
// it is a duplicate and belongs deleted.
// ---------------------------------------------------------------------------

// AttendedPublishPlan is one attended publish stated once, so the mint and the
// spend cannot disagree about what is being published.
//
// The logical times are explicit because the broker never reads a wall clock
// and every ordering refusal in validatePublishApproval is stated in terms of
// them.
type AttendedPublishPlan struct {
	// Identity is the publication identity MINUS its approval reference: the
	// ref does not exist until MintAttendedApproval returns it, and binding it
	// here would invite a caller to supply one.
	Identity PublishIdentity
	// Hashes are the three digests recomputed from the projection directory.
	// They are what the approval STAMPS, so an approval minted against these
	// bytes authorizes no other bytes.
	Hashes PublishHashes

	// Requester and DecidedBy are recorded in the immutable approval objects.
	Requester string
	DecidedBy string

	// EpisodeID names the durable episode the publish intent is appended to.
	EpisodeID string

	// RequestedAt <= DecidedAt <= PublishAt < ExpiresAt. The inequalities are
	// enforced by validatePublishApproval, not restated here.
	RequestedAt int64
	DecidedAt   int64
	PublishAt   int64
	ExpiresAt   int64
}

// ApprovalScope renders the canonical publish-bound approval scope this plan
// mints and spends. It is exposed because an attended operator must be able to
// SEE the scope they are approving before it is minted, and because the mint
// and the spend must demonstrably use the same string.
func (p AttendedPublishPlan) ApprovalScope() string {
	return PublishApprovalScopeFor(p.Identity, p.Hashes, p.ExpiresAt)
}

// PublishScope renders the frozen capability scope for this plan.
func (p AttendedPublishPlan) PublishScope() string {
	return PublishScope(p.Identity.RegistryOrigin, p.Identity.Vendor, p.Identity.Name, p.Identity.Version)
}

// Payload renders the canonical publish payload for an approval that has
// already been minted.
func (p AttendedPublishPlan) Payload(approvalRef hashref.HashRef) []byte {
	id := p.Identity
	id.ApprovalRef = approvalRef
	return EncodePublishPayload(id, p.Hashes)
}

// ErrAttendedApprovalNotObserved reports that the approval was minted and
// decided, but the landed poll effect did not observe it back.
//
// It is not a paranoia check. The mint traverses three landed surfaces
// (Human.Approve -> DecideApproval -> Human.PollApproval); without reading the
// decision BACK OUT through the poll and hashing it, "the approval landed"
// would rest on the return value of the function that claims to have landed it.
var ErrAttendedApprovalNotObserved = errors.New(
	"broker: the minted approval was not observed back through Human.PollApproval")

// MintAttendedApproval performs the FIRST of the two attended acts: it mints
// the one-shot ApprovalRequestV1 / ApprovalDecisionV1 pair through the landed
// traversal and returns the decision's content hash, which is the approvalRef
// the publish payload will name and the key the durable single-use claim is
// taken on.
//
// It deliberately does NOT publish. Minting and spending are two separate
// invocations of world-publish (Decision D4) so that the operator can review
// the minted ref, and so that the durable claim is spent by a command that
// could not have created it.
func MintAttendedApproval(s *store.Store, plan AttendedPublishPlan) (hashref.HashRef, error) {
	return mintAttendedApproval(s, plan)
}

func mintAttendedApproval(s approvalStore, plan AttendedPublishPlan) (hashref.HashRef, error) {
	scope := plan.ApprovalScope()
	human := newHumanHandler(s)

	// Leg 1: the landed Human.Approve effect mints the immutable request.
	requestSession := newSession(s, plan.EpisodeID+"-approve", []Capability{
		{Effect: EffectHumanApprove, Scope: scope, ExpiresAt: plan.ExpiresAt, Budget: PublishCost},
	}, Registry{EffectHumanApprove: human}, Live, nil)
	pending, _, err := requestSession.Invoke(context.Background(), EffectRequest{
		Effect: EffectHumanApprove, Scope: scope, Cost: PublishCost, Now: plan.RequestedAt,
	}, mustApprovalJSON(approvalInputWire{Requester: plan.Requester}))
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf("broker: mint approval request: %w", err)
	}
	var observedPending pendingWire
	if err := decodeApprovalJSON(pending, &observedPending); err != nil {
		return hashref.HashRef{}, fmt.Errorf("broker: decode pending approval: %w", err)
	}
	requestRef, err := hashref.Parse(observedPending.RequestRef)
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf("broker: parse minted request ref: %w", err)
	}

	// Leg 2: the landed operator entry point mints the immutable decision.
	decisionRef, err := decideApproval(s, requestRef, "approve", plan.DecidedBy, plan.DecidedAt)
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf("broker: decide approval: %w", err)
	}

	// Leg 3: the landed poll effect OBSERVES the decision back, and the observed
	// bytes are hashed and compared with the ref the payload will carry. Without
	// this leg the returned ref would be evidence only of its own construction.
	pollSession := newSession(s, plan.EpisodeID+"-poll", []Capability{
		{Effect: EffectHumanPollApproval, Scope: scope, ExpiresAt: plan.ExpiresAt, Budget: PublishCost},
	}, Registry{EffectHumanPollApproval: human}, Live, nil)
	polled, _, err := pollSession.Invoke(context.Background(), EffectRequest{
		Effect: EffectHumanPollApproval, Scope: scope, Cost: PublishCost, Now: plan.DecidedAt,
	}, mustApprovalJSON(approvalInputWire{RequestRef: requestRef.String()}))
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf("broker: poll minted approval: %w", err)
	}
	if err := observeMintedDecision(polled, decisionRef); err != nil {
		return hashref.HashRef{}, err
	}
	return decisionRef, nil
}

// observeMintedDecision is leg 3's check, extracted so it can be driven
// directly. Reaching its two refusals through mintAttendedApproval would need
// the approvals chain severed BETWEEN the decide and the poll, which no caller
// can do — and a branch that can only be reached by a caller that does not
// exist is a branch no test can red.
//
// It hashes the bytes the poll RETURNED and compares them with the ref
// decideApproval produced. The two are computed by different code from
// different sources; without the comparison, the returned ref would be evidence
// only of its own construction.
func observeMintedDecision(polled []byte, decisionRef hashref.HashRef) error {
	var observed observedDecisionWire
	if err := decodeApprovalJSON(polled, &observed); err != nil {
		return fmt.Errorf("broker: decode polled approval: %w", err)
	}
	if observed.Status != "decided" {
		return fmt.Errorf(
			"%w: poll reports status %q", ErrAttendedApprovalNotObserved, observed.Status)
	}
	if got := hashref.SumSHA256(observed.Decision); got != decisionRef {
		return fmt.Errorf(
			"%w: polled decision hashes to %s, want %s",
			ErrAttendedApprovalNotObserved, got.String(), decisionRef.String())
	}
	return nil
}

// Publish outcome classes as the OPERATOR sees them. They are the three
// dispositions of Decision 3's table that a human has to act on differently,
// and nothing more: done, stop, or reconcile.
const (
	AttendedPublishSucceeded     = "succeeded"
	AttendedPublishFailed        = "failed"
	AttendedPublishIndeterminate = "indeterminate"
)

// AttendedPublishResult is the verbatim result of ONE attempt.
type AttendedPublishResult struct {
	Status string
	// RecordRef is the immutable EffectRecord, present for a resolved outcome
	// and ZERO for an indeterminate one — because an indeterminate attempt
	// deliberately appends no record.
	RecordRef hashref.HashRef
	// InvocationID, EpisodeID and Ordinal name the durable intent a
	// reconciliation pass must resolve. They are populated only for the
	// indeterminate class, which is the only class that owes a reconciliation.
	InvocationID string
	EpisodeID    string
	Ordinal      int64
	// Detail is the redacted handler detail, or the failure reason.
	Detail string
}

// Reconcilable reports whether this result leaves a receipt a human must
// resolve read-only before anything else may happen.
func (r AttendedPublishResult) Reconcilable() bool {
	return r.Status == AttendedPublishIndeterminate
}

// InvokeAttendedPublish performs the SECOND attended act: it constructs the
// one-shot capability and calls Session.Invoke EXACTLY ONCE.
//
// "Exactly once" is structural, not a promise: there is no loop and no retry in
// this function, and MUT-D0-INDETERMINATE-RETRY — which adds one — is the
// mutation that reds AC26 by driving the fake validator's own request counter
// to 2. A retry on an ambiguous dispatch is the double-publish this entire
// design exists to prevent: the request body may already have reached a
// registry that cannot be asked to take it back.
//
// The handler is INJECTED rather than constructed here, and that is the whole
// content of D2: cmd/world-publish hands in the PRODUCTION constructor's
// handler, publish_op_test.go hands in the loopback one, and neither can reach
// the other's door.
func InvokeAttendedPublish(
	ctx context.Context,
	s *store.Store,
	handler Handler,
	plan AttendedPublishPlan,
	approvalRef hashref.HashRef,
) (AttendedPublishResult, error) {
	return invokeAttendedPublish(ctx, s, handler, plan, approvalRef)
}

func invokeAttendedPublish(
	ctx context.Context,
	s objectStore,
	handler Handler,
	plan AttendedPublishPlan,
	approvalRef hashref.HashRef,
) (AttendedPublishResult, error) {
	scope := plan.PublishScope()
	// The grant is exact in all four dimensions and its budget is exactly one
	// irreversible attempt. MUT-D0-BUDGET-2 raises it and reds AC25: with a
	// budget of two the in-memory debit no longer refuses the second attempt,
	// so the refusal has to come from the DURABLE claim — which is precisely
	// the property AC25 is stated to distinguish.
	grant := Capability{
		Effect:    EffectRegistryPublish,
		Scope:     scope,
		ExpiresAt: plan.ExpiresAt,
		Budget:    PublishCost,
	}
	session := newSession(s, plan.EpisodeID, []Capability{grant},
		Registry{EffectRegistryPublish: handler}, Live, nil)

	_, recordRef, err := session.Invoke(ctx, EffectRequest{
		Effect: EffectRegistryPublish, Scope: scope, Cost: PublishCost, Now: plan.PublishAt,
	}, plan.Payload(approvalRef))

	// THERE IS NO SECOND Invoke BELOW THIS LINE. Everything that follows is
	// classification of the one attempt that was made.
	var indeterminate *IndeterminateEffectError
	switch {
	case err == nil:
		return AttendedPublishResult{
			Status: AttendedPublishSucceeded, RecordRef: recordRef, EpisodeID: plan.EpisodeID,
		}, nil
	case errors.As(err, &indeterminate):
		return AttendedPublishResult{
			Status:       AttendedPublishIndeterminate,
			InvocationID: indeterminate.InvocationID,
			EpisodeID:    indeterminate.EpisodeID,
			Ordinal:      indeterminate.Ordinal,
			Detail:       indeterminate.Detail,
		}, err
	default:
		return AttendedPublishResult{
			Status: AttendedPublishFailed, RecordRef: recordRef,
			EpisodeID: plan.EpisodeID, Detail: err.Error(),
		}, err
	}
}
