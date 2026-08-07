package broker

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// ---------------------------------------------------------------------------
// SM.B2b — the publish-bound approval scope
//
// The whole point of this milestone's approval work is that NOTHING in the
// landed approval surface changes. Human.Approve, Human.PollApproval,
// HumanHandler, DecideApproval, ApprovalRequestV1, ApprovalDecisionV1,
// approvalRequestWire and approvalDecisionWire are untouched; the publish
// binding is a canonical VALUE of the existing Scope string. These tests are
// what makes that claim checkable in both directions: publish-bound decisions
// must ride the landed wires, and pre-SM.B2b non-publish bytes must stay valid.
// ---------------------------------------------------------------------------

func canonicalPublishApprovalScope() PublishApprovalScope {
	return PublishApprovalScope{
		Publish:       PublishScope("https://registry.example", "world", "core", "0.1.0"),
		Effect:        EffectRegistryPublish,
		ManifestRef:   "sha256:" + strings.Repeat("11", 32),
		TarballSHA256: "sha256:" + strings.Repeat("22", 32),
		ContentHash:   "sha256:" + strings.Repeat("33", 32),
		InterfaceHash: "sha256:" + strings.Repeat("44", 32),
		ExpiresAt:     100,
	}
}

func TestPublishApprovalScopeRoundTripsThroughTheFrozenTermOrder(t *testing.T) {
	want := canonicalPublishApprovalScope()
	raw := FormatPublishApprovalScope(want)

	// The rendered suffix must present every frozen term, in the frozen order,
	// exactly once. Reading the order back out of the BYTES rather than out of
	// publishApprovalScopeTerms is what stops this being a tautology.
	mark := strings.Index(raw, publishApprovalScopeMark)
	if mark < 0 {
		t.Fatalf("rendered scope %q carries no %q mark", raw, publishApprovalScopeMark)
	}
	var keys []string
	for _, term := range strings.Split(raw[mark+len(publishApprovalScopeMark):], publishApprovalTermSep) {
		key, _, found := strings.Cut(term, publishApprovalKeySep)
		if !found {
			t.Fatalf("rendered term %q is not an assignment", term)
		}
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		t.Fatal("read zero terms out of the rendered scope; the instrument is broken")
	}
	if !sameStrings(keys, publishApprovalScopeTerms) {
		t.Fatalf("rendered term order = %v, want %v", keys, publishApprovalScopeTerms)
	}

	got, err := ParsePublishApprovalScope(raw)
	if err != nil {
		t.Fatalf("ParsePublishApprovalScope(%q): %v", raw, err)
	}
	if got != want {
		t.Fatalf("round trip = %#v, want %#v", got, want)
	}
	// And the head never carries the separator, which is what makes the mark
	// unambiguous no matter what a payload smuggles into its origin.
	if strings.Contains(got.Publish, publishApprovalScopeMark) {
		t.Fatalf("parsed publish grammar %q carries the scope mark", got.Publish)
	}
	// PublishApprovalScopeFor must agree with the hand-built canonical value.
	id := PublishIdentity{
		Vendor: "world", Name: "core", Version: "0.1.0",
		RegistryOrigin: "https://registry.example",
		ManifestRef:    mustHashRef(t, want.ManifestRef),
	}
	if built := PublishApprovalScopeFor(id, PublishHashes{
		TarballSHA256: want.TarballSHA256,
		ContentHash:   want.ContentHash,
		InterfaceHash: want.InterfaceHash,
	}, want.ExpiresAt); built != raw {
		t.Fatalf("PublishApprovalScopeFor = %q, want %q", built, raw)
	}
}

func mustHashRef(t *testing.T, raw string) hashref.HashRef {
	t.Helper()
	ref, err := hashref.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	return ref
}

func TestParsePublishApprovalScopeIsStrictInEveryDirection(t *testing.T) {
	canonical := FormatPublishApprovalScope(canonicalPublishApprovalScope())
	// CONTROL FIRST: the unmodified value must parse, or every refusal below is
	// satisfied by an unrelated defect.
	if _, err := ParsePublishApprovalScope(canonical); err != nil {
		t.Fatalf("the control scope was refused: %v", err)
	}

	head, tail, _ := strings.Cut(canonical, publishApprovalScopeMark)
	terms := strings.Split(tail, publishApprovalTermSep)

	rebuild := func(head string, terms []string) string {
		return head + publishApprovalScopeMark + strings.Join(terms, publishApprovalTermSep)
	}
	swapped := append([]string(nil), terms...)
	swapped[0], swapped[1] = swapped[1], swapped[0]
	dropped := append([]string(nil), terms[:len(terms)-1]...)
	extra := append(append([]string(nil), terms...), "surplus=1")
	emptyValue := append([]string(nil), terms...)
	emptyValue[2] = "tarball="
	twoAssignments := append([]string(nil), terms...)
	twoAssignments[2] = "tarball=a=b"
	noAssignment := append([]string(nil), terms...)
	noAssignment[2] = "tarball"
	badExpiry := append([]string(nil), terms...)
	badExpiry[len(badExpiry)-1] = "expires=soon"

	for _, tc := range []struct {
		name     string
		raw      string
		contains string
	}{
		{"no mark at all", head, "carries no"},
		{"a second mark", canonical + publishApprovalScopeMark + "x", "more than one"},
		{"a dropped term", rebuild(head, dropped), "want exactly"},
		{"a surplus term", rebuild(head, extra), "want exactly"},
		{"a reordered term", rebuild(head, swapped), "the term order is frozen"},
		{"an empty value", rebuild(head, emptyValue), "is empty"},
		{"two assignments in one term", rebuild(head, twoAssignments), "exactly one"},
		{"a term that is not an assignment", rebuild(head, noAssignment), "exactly one"},
		{"a non-decimal expiry", rebuild(head, badExpiry), "not a decimal logical time"},
		{"an empty publish grammar", rebuild("", terms), "empty publish grammar"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParsePublishApprovalScope(tc.raw)
			if !errors.Is(err, ErrPublishApprovalMalformed) {
				t.Fatalf("%s error = %T %v, want ErrPublishApprovalMalformed", tc.name, err, err)
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("%s refusal %q does not mention %q", tc.name, err, tc.contains)
			}
		})
	}
}

// TestPublishApprovalScopeRefusesASmuggledFragment is the test approve.go's
// publishApprovalScopeMark comment names. The '#' separator is safe for two
// independent reasons; this measures the second, which does not depend on
// validatePublishOrigin having run at all.
func TestPublishApprovalScopeRefusesASmuggledFragment(t *testing.T) {
	base := openPublishStore(t)
	fixture := newPublishFixture(t, "http://127.0.0.1:1", "fragment").
		landApproval(t, base, "approve", defaultApprovalTimes())
	req := EffectRequest{
		Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: 50,
	}

	// KNOWN-POSITIVE CONTROL: the unmodified payload is accepted, so the
	// refusal below is attributable to the fragment and not to the fixture.
	if _, err := validatePublishApproval(base, fixture.payload, req); err != nil {
		t.Fatalf("control payload was refused: %v", err)
	}

	// The mechanism is field-AGNOSTIC: it does not care which component of the
	// publish grammar the mark arrives in, because the parsed head cannot
	// contain one at all. Demonstrating it on every field the grammar
	// concatenates is what lets a future reader see that without redoing the
	// analysis.
	suffix := strings.SplitN(fixture.approvalScope, publishApprovalScopeMark, 2)[1]
	for _, tc := range []struct {
		field  string
		mutate func(*PublishIdentity)
	}{
		{"RegistryOrigin", func(id *PublishIdentity) { id.RegistryOrigin += publishApprovalScopeMark + suffix }},
		{"Vendor", func(id *PublishIdentity) { id.Vendor += publishApprovalScopeMark + suffix }},
		{"Name", func(id *PublishIdentity) { id.Name += publishApprovalScopeMark + suffix }},
		{"Version", func(id *PublishIdentity) { id.Version += publishApprovalScopeMark + suffix }},
	} {
		t.Run(tc.field, func(t *testing.T) {
			smuggled := fixture.identity
			tc.mutate(&smuggled)
			smuggledScope := PublishScope(
				smuggled.RegistryOrigin, smuggled.Vendor, smuggled.Name, smuggled.Version)
			if !strings.Contains(smuggledScope, publishApprovalScopeMark) {
				t.Fatal("the smuggled scope carries no mark; this arm would pass vacuously")
			}
			_, err := validatePublishApproval(base, EncodePublishPayload(smuggled, fixture.hashes),
				EffectRequest{
					Effect: EffectRegistryPublish, Scope: smuggledScope, Cost: PublishCost, Now: 50,
				})
			if !errors.Is(err, ErrPublishApprovalScope) {
				t.Fatalf("smuggled mark in %s: error = %T %v, want ErrPublishApprovalScope",
					tc.field, err, err)
			}
			// The refusal is the head comparison, not an incidental one: the
			// parsed head is mark-free by construction, so it can never equal a
			// wantPublish that carries one.
			if !strings.Contains(err.Error(), "approval stamps") {
				t.Errorf("smuggled mark in %s: refusal %q is not the publish-head comparison",
					tc.field, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// landed request/decision/poll compatibility
// ---------------------------------------------------------------------------

// TestLandedApprovalWiresCarryPublishBoundDecisionsUnchanged proves that a
// publish-bound approval is an ORDINARY landed approval: the same semantic IDs,
// the same two wires, decoded by the same DisallowUnknownFields codec, observed
// back through the same poll effect.
func TestLandedApprovalWiresCarryPublishBoundDecisionsUnchanged(t *testing.T) {
	base := openPublishStore(t)
	times := defaultApprovalTimes()
	fixture := newPublishFixture(t, "http://127.0.0.1:1", "wire").landApproval(t, base, "approve", times)

	requestObj, ok, err := base.GetObject(fixture.approvalRequestRef)
	if err != nil || !ok {
		t.Fatalf("approval request object = ok %v, err %v", ok, err)
	}
	if requestObj.SemanticID != ApprovalRequestV1 {
		t.Fatalf("publish-bound request semantic ID = %q, want %q",
			requestObj.SemanticID, ApprovalRequestV1)
	}
	var request approvalRequestWire
	if err := decodeApprovalJSON(requestObj.Payload, &request); err != nil {
		t.Fatalf("the landed request codec rejected a publish-bound request: %v", err)
	}
	wantRequest := approvalRequestWire{
		Effect: EffectHumanApprove, Scope: fixture.approvalScope, Cost: PublishCost,
		Requester: "sm-b2b-fixture", Now: times.request,
	}
	if request != wantRequest {
		t.Fatalf("publish-bound request = %+v, want %+v", request, wantRequest)
	}
	// Re-encoding must be byte-identical: a widened wire would show up here as
	// an extra key even though decoding still succeeded.
	if got := mustApprovalJSON(request); !bytes.Equal(got, requestObj.Payload) {
		t.Fatalf("re-encoded request = %s, want the landed bytes %s", got, requestObj.Payload)
	}

	decisionObj, ok, err := base.GetObject(fixture.identity.ApprovalRef)
	if err != nil || !ok {
		t.Fatalf("approval decision object = ok %v, err %v", ok, err)
	}
	if decisionObj.SemanticID != ApprovalDecisionV1 {
		t.Fatalf("publish-bound decision semantic ID = %q, want %q",
			decisionObj.SemanticID, ApprovalDecisionV1)
	}
	var decision approvalDecisionWire
	if err := decodeApprovalJSON(decisionObj.Payload, &decision); err != nil {
		t.Fatalf("the landed decision codec rejected a publish-bound decision: %v", err)
	}
	wantDecision := approvalDecisionWire{
		RequestRef: fixture.approvalRequestRef.String(), Decision: "approve",
		DecidedBy: "attended-operator", Now: times.decide,
	}
	if decision != wantDecision {
		t.Fatalf("publish-bound decision = %+v, want %+v", decision, wantDecision)
	}
	if got := mustApprovalJSON(decision); !bytes.Equal(got, decisionObj.Payload) {
		t.Fatalf("re-encoded decision = %s, want the landed bytes %s", got, decisionObj.Payload)
	}

	// The scope value the wire carried IS the publish binding, and it parses.
	scope, err := ParsePublishApprovalScope(request.Scope)
	if err != nil {
		t.Fatalf("the landed request's scope is not publish-grammatical: %v", err)
	}
	if scope.Effect != EffectRegistryPublish {
		t.Errorf("landed scope effect = %q, want %q", scope.Effect, EffectRegistryPublish)
	}
	if scope.TarballSHA256 != fixture.hashes.TarballSHA256 {
		t.Errorf("landed scope tarball = %q, want %q", scope.TarballSHA256, fixture.hashes.TarballSHA256)
	}
}

// TestOldNonPublishApprovalBytesRemainValid pins the other direction. The two
// literals below are the exact bytes the landed encoders produced BEFORE
// SM.B2b — an approval of "release", nothing to do with publishing.
func TestOldNonPublishApprovalBytesRemainValid(t *testing.T) {
	const landedRequest = `{"effect":"Human.Approve","scope":"release","cost":2,"requester":"agent-7","now":10}`
	const landedDecision = `{"requestRef":"sha256:aa11bb22cc33dd44ee55ff66007788990011223344556677889900aabbccddee",` +
		`"decision":"approve","decidedBy":"operator","now":11}`

	var request approvalRequestWire
	if err := decodeApprovalJSON([]byte(landedRequest), &request); err != nil {
		t.Fatalf("pre-SM.B2b request bytes were rejected: %v", err)
	}
	wantRequest := approvalRequestWire{
		Effect: EffectHumanApprove, Scope: "release", Cost: 2, Requester: "agent-7", Now: 10,
	}
	if request != wantRequest {
		t.Fatalf("decoded old request = %+v, want %+v", request, wantRequest)
	}
	if got := string(mustApprovalJSON(request)); got != landedRequest {
		t.Fatalf("re-encoded old request = %s, want the original bytes %s", got, landedRequest)
	}

	var decision approvalDecisionWire
	if err := decodeApprovalJSON([]byte(landedDecision), &decision); err != nil {
		t.Fatalf("pre-SM.B2b decision bytes were rejected: %v", err)
	}
	if got := string(mustApprovalJSON(decision)); got != landedDecision {
		t.Fatalf("re-encoded old decision = %s, want the original bytes %s", got, landedDecision)
	}

	// A non-publish scope is not accidentally publish-grammatical.
	if _, err := ParsePublishApprovalScope(request.Scope); !errors.Is(err, ErrPublishApprovalMalformed) {
		t.Fatalf("ParsePublishApprovalScope(%q) = %v, want ErrPublishApprovalMalformed", request.Scope, err)
	}

	// And the whole landed flow still works end to end for a non-publish scope,
	// on a live store: request, decide, poll.
	base := openPublishStore(t)
	requestRef, decisionRef := landNonPublishApproval(t, base, "release", 2,
		approvalTimes{request: 10, decide: 11, expires: 100})
	human := newHumanHandler(base)
	pollSession := newSession(base, "old-bytes-poll", []Capability{
		{Effect: EffectHumanPollApproval, Scope: "release", ExpiresAt: 100, Budget: 2},
	}, Registry{EffectHumanPollApproval: human}, Live, nil)
	polled, _, err := pollSession.Invoke(context.Background(), EffectRequest{
		Effect: EffectHumanPollApproval, Scope: "release", Cost: 1, Now: 12,
	}, mustApprovalJSON(approvalInputWire{RequestRef: requestRef.String()}))
	if err != nil {
		t.Fatalf("polling a non-publish approval: %v", err)
	}
	var observed observedDecisionWire
	if err := decodeApprovalJSON(polled, &observed); err != nil {
		t.Fatal(err)
	}
	if observed.Status != "decided" {
		t.Fatalf("non-publish poll status = %q, want \"decided\"", observed.Status)
	}
	if got := hashref.SumSHA256(observed.Decision); got != decisionRef {
		t.Fatalf("polled non-publish decision hashes to %s, want %s", got, decisionRef)
	}
}
