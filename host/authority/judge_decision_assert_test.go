package authority

// Row 201 A4 sprint 3 M6c (design §2.7 [r2], §6 JD1-JD3, the passkey half): JudgeDecisionAssertion's
// steps, causes and fields agree with the mint judge's (JudgeAssertion): parse, store, unknown key,
// stored key, law; it does not judge the counter (W15 is the writer's). Every passkey is
// approvaltest's (D-WORLD-109); the rig is judge_test.go's (mark-mac enrolled with counter 1).

import (
	"bytes"
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
)

// jdAssert is pk's assertion of in's decision message hash with counter count.
func jdAssert(t *testing.T, pk *approvaltest.Passkey, in DecisionJudgeInput, count uint32) []byte {
	t.Helper()
	h, err := approval.DecisionMessageHash(in.Verb, in.Packet, in.Request, in.Expiry)
	if err != nil {
		t.Fatal(err)
	}
	return pk.Assert(jgOrigin, jgRPID, h, count)
}

func TestJudgeDecisionAssertion(t *testing.T) {
	stranger := approvaltest.NewPasskey(approval.AlgEdDSA, 62)
	mint, _ := approval.MessageHash(jdIn.Request, jdIn.Packet, jdIn.Expiry)
	other := jdIn
	other.Verb = approval.VerbApprove
	for _, tc := range []struct {
		name     string
		setup    []string
		artifact []byte
		step     string
		cause    error
		rawID    []byte
		approver string
	}{
		{name: "JD-parse", artifact: []byte("{}"), step: JudgeStepParse, cause: approval.ErrPasskeyResponse},
		{name: "JD-store-v5", setup: []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`},
			artifact: jdAssert(t, jgPasskey, jdIn, 2), step: JudgeStepStore},
		{name: "JD-unknown-key", artifact: jdAssert(t, stranger, jdIn, 2), step: JudgeStepUnknownKey, cause: ErrNotEnrolled, rawID: stranger.CredentialID()},
		{name: "JD-stored-key", setup: []string{`UPDATE approvers SET public_key = X'00' WHERE name = 'mark-mac'`}, artifact: jdAssert(t, jgPasskey, jdIn, 2),
			step: JudgeStepStoredKey, approver: "mark-mac"},
		{name: "JD-law-mint-assertion", artifact: jgPasskey.Assert(jgOrigin, jgRPID, mint, 2), step: JudgeStepLaw, cause: approval.ErrChallenge},
		{name: "JD-law-other-verb", artifact: jdAssert(t, jgPasskey, other, 2), step: JudgeStepLaw, cause: approval.ErrChallenge},
		{name: "JD-law-revoked", setup: []string{`UPDATE approvers SET revoked_at = 950 WHERE name = 'mark-mac'`}, artifact: jdAssert(t, jgPasskey, jdIn, 2),
			step: JudgeStepLaw, cause: approval.ErrApproverRevoked},
		{name: "JD-law-backup-eligible", setup: []string{`UPDATE webauthn_credentials SET backup_eligible = 1 WHERE name = 'mark-mac'`}, artifact: jdAssert(t, jgPasskey, jdIn, 2),
			step: JudgeStepLaw, cause: approval.ErrBackupEligibility},
		{name: "JD-law-origin", setup: []string{`UPDATE webauthn_credentials SET origin = 'https://x.ts.net' WHERE name = 'mark-mac'`}, artifact: jdAssert(t, jgPasskey, jdIn, 2),
			step: JudgeStepLaw, cause: approval.ErrOriginMismatch},
	} {
		path, _ := jgRig(t)
		jgExec(t, path, tc.setup...)
		ro, err := store.OpenCredentialsReadOnly(path)
		if err != nil {
			t.Fatal(err)
		}
		d, err := JudgeDecisionAssertion(boundedTestContext(t), ro, jdIn, tc.artifact)
		_ = ro.Close()
		var je *JudgeError
		if d != nil || !errors.As(err, &je) || je.Step != tc.step || (tc.cause != nil && !errors.Is(err, tc.cause)) || !bytes.Equal(je.RawID, tc.rawID) || je.Approver != tc.approver {
			t.Errorf("%s: %v, %v; want step %s cause %v", tc.name, d, err, tc.step, tc.cause)
		}
	}
	path, _ := jgRig(t)
	ro, err := store.OpenCredentialsReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ro.Close() }()
	for _, count := range []uint32{2, 1} { // the stored counter is 1: the judge leaves 1 to the writer (W15 deferred)
		d, err := JudgeDecisionAssertion(boundedTestContext(t), ro, jdIn, jdAssert(t, jgPasskey, jdIn, count))
		id, n, ok := d.PasskeyUse()
		if err != nil || d.DecidedBy() != "webauthn:mark-mac" || d.Binds(jdIn.Verb, jdIn.Packet, jdIn.Request, jdIn.Expiry) != nil || !ok || n != count || !bytes.Equal(id, jgPasskey.CredentialID()) {
			t.Fatalf("counter %d: %v, %v", count, d, err)
		}
	}
}
