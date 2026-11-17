package authority

// Row 243 V3 sprint 1 M2 (design §2.3, §6 VJ1): JudgeVerbSSH, JudgeVerbAssertion and RecheckVerb
// name each step as the mint and decision judges do — parse, store, unknown key, stored key, law —
// in that order, on a read-only handle. Every key is approvaltest's (D-WORLD-109); the rig is
// judge_test.go's (mark-mac's passkey enrolled with counter 1, mark-ssh, both at 900).

import (
	"bytes"
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
)

var vjIn = VerbJudgeInput{Digest: [32]byte{0x76, 1}, Challenge: [32]byte{0x76, 2}, BytesDigest: [32]byte{0x76, 1}, CreatedAt: 1000, ExpiresAt: 1600, Now: 1010}

const vjDropV4 = `DROP TABLE key_decisions; DROP TABLE surface_approvals; DROP TABLE webauthn_credentials; DROP TABLE approvers; PRAGMA user_version = 4`
const vjDropV5 = `DROP TABLE key_decisions; DROP TABLE surface_approvals; DROP TABLE webauthn_credentials; PRAGMA user_version = 5`

// vjSign and vjAssert are a confirmation of in's request by s or pk.
func vjSign(t *testing.T, s *approvaltest.Signer, namespace string, in VerbJudgeInput) []byte {
	t.Helper()
	m, err := approval.VerbMessage(in.Digest, in.Challenge, in.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	return s.Sign(namespace, m)
}

func vjAssert(t *testing.T, pk *approvaltest.Passkey, in VerbJudgeInput, count uint32) []byte {
	t.Helper()
	h, err := approval.VerbMessageHash(in.Digest, in.Challenge, in.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	return pk.Assert(jgOrigin, jgRPID, h, count)
}

// vjCase is one refusal: setup runs on the rig, then judge(artifact) must fail at step with cause
// (when set) and the unknown-key / stored-key fields.
type vjCase struct {
	name     string
	setup    string
	in       VerbJudgeInput
	artifact []byte
	step     string
	cause    error
	field    string // the unknown-key SSHSIG fingerprint
	rawID    []byte // the unknown-key credential ID
	approver string // the stored-key approver
}

type vjJudge func(t *testing.T, ro *store.Credentials, in VerbJudgeInput, artifact []byte) (*approval.SignedVerb, error)

func vjRun(t *testing.T, judge vjJudge, cases []vjCase) {
	t.Helper()
	for _, tc := range cases {
		path, _ := jgRig(t)
		if tc.setup != "" {
			jgExec(t, path, tc.setup)
		}
		ro, err := store.OpenCredentialsReadOnly(path)
		if err != nil {
			t.Fatal(err)
		}
		v, err := judge(t, ro, tc.in, tc.artifact)
		_ = ro.Close()
		var je *JudgeError
		if v != nil || !errors.As(err, &je) || je.Step != tc.step || (tc.cause != nil && !errors.Is(err, tc.cause)) ||
			je.Fingerprint != tc.field || !bytes.Equal(je.RawID, tc.rawID) || je.Approver != tc.approver {
			t.Errorf("%s: %v, %v; want step %s cause %v", tc.name, v, err, tc.step, tc.cause)
		}
	}
}

// vjOK opens the rig after setup and returns judge's verified confirmation.
func vjOK(t *testing.T, judge vjJudge, setup string, in VerbJudgeInput, artifact []byte) *approval.SignedVerb {
	t.Helper()
	path, _ := jgRig(t)
	if setup != "" {
		jgExec(t, path, setup)
	}
	ro, err := store.OpenCredentialsReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ro.Close() }()
	v, err := judge(t, ro, in, artifact)
	if err != nil || v.Binds(in.Digest, in.Challenge, in.ExpiresAt) != nil {
		t.Fatalf("control: %v, %v", v, err)
	}
	return v
}

func TestJudgeVerbSSH(t *testing.T) {
	judge := func(t *testing.T, ro *store.Credentials, in VerbJudgeInput, a []byte) (*approval.SignedVerb, error) {
		return JudgeVerbSSH(boundedTestContext(t), ro, in, a)
	}
	good := vjSign(t, jgSigner, approval.VerbNamespace, vjIn)
	stranger := approvaltest.NewSigner(66)
	mintMsg, _ := approval.SignedMessage(vjIn.Digest, vjIn.Challenge, vjIn.ExpiresAt)
	consumed, otherChallenge, badBytes := vjIn, vjIn, vjIn
	consumed.Consumed = true
	otherChallenge.Challenge[0] ^= 1
	badBytes.BytesDigest[0] ^= 1 // the stored bytes do not hash to the digest column (P1)
	vjRun(t, judge, []vjCase{
		{name: "VJ-parse", in: vjIn, artifact: []byte("y\n"), step: JudgeStepParse, cause: approval.ErrArmor},
		{name: "VJ-parse before store", setup: vjDropV4, in: vjIn, artifact: []byte("y\n"), step: JudgeStepParse},
		{name: "VJ-store-v4", setup: vjDropV4, in: vjIn, artifact: good, step: JudgeStepStore},
		{name: "VJ-unknown-key", in: vjIn, artifact: vjSign(t, stranger, approval.VerbNamespace, vjIn), step: JudgeStepUnknownKey, cause: ErrNotEnrolled,
			field: stranger.PublicKey().Fingerprint()},
		{name: "VJ-unknown-key before law", in: consumed, artifact: vjSign(t, stranger, approval.VerbNamespace, vjIn), step: JudgeStepUnknownKey,
			field: stranger.PublicKey().Fingerprint()},
		{name: "VJ-law-consumed", in: consumed, artifact: good, step: JudgeStepLaw},
		{name: "VJ-law-bytes-digest", in: badBytes, artifact: good, step: JudgeStepLaw},
		{name: "VJ-law-mint-signature", in: vjIn, artifact: jgSigner.Sign(approval.Namespace, mintMsg), step: JudgeStepLaw, cause: approval.ErrNamespace},
		{name: "VJ-law-other-challenge", in: otherChallenge, artifact: good, step: JudgeStepLaw, cause: approval.ErrBadSignature},
		{name: "VJ-law-revoked", setup: `UPDATE approvers SET revoked_at = 950 WHERE name = 'mark-ssh'`, in: vjIn, artifact: good, step: JudgeStepLaw, cause: approval.ErrApproverRevoked},
	})
	var na *approval.NotAdmittedError
	path, _ := jgRig(t)
	ro, _ := store.OpenCredentialsReadOnly(path)
	if _, err := JudgeVerbSSH(boundedTestContext(t), ro, consumed, good); !errors.As(err, &na) || na.Verdict != approval.RefuseConsumed {
		t.Errorf("consumed: %v, want the P1 verdict through the law step", err)
	}
	_ = ro.Close()
	if v := vjOK(t, judge, "", vjIn, good); v.ApprovedBy() != "ssh:mark-ssh" {
		t.Fatalf("ApprovedBy = %q", v.ApprovedBy())
	}
}

func TestJudgeVerbAssertion(t *testing.T) {
	judge := func(t *testing.T, ro *store.Credentials, in VerbJudgeInput, a []byte) (*approval.SignedVerb, error) {
		return JudgeVerbAssertion(boundedTestContext(t), ro, in, a)
	}
	stranger := approvaltest.NewPasskey(approval.AlgEdDSA, 67)
	mint, _ := approval.MessageHash(vjIn.Digest, vjIn.Challenge, vjIn.ExpiresAt)
	expired := vjIn
	expired.Now = vjIn.ExpiresAt
	vjRun(t, judge, []vjCase{
		{name: "VJ-parse", in: vjIn, artifact: []byte("{}"), step: JudgeStepParse, cause: approval.ErrPasskeyResponse},
		{name: "VJ-store-v5", setup: vjDropV5, in: vjIn, artifact: vjAssert(t, jgPasskey, vjIn, 2), step: JudgeStepStore},
		{name: "VJ-unknown-key", in: vjIn, artifact: vjAssert(t, stranger, vjIn, 2), step: JudgeStepUnknownKey, cause: ErrNotEnrolled, rawID: stranger.CredentialID()},
		{name: "VJ-stored-key", setup: `UPDATE approvers SET public_key = X'00' WHERE name = 'mark-mac'`, in: vjIn, artifact: vjAssert(t, jgPasskey, vjIn, 2),
			step: JudgeStepStoredKey, approver: "mark-mac"},
		{name: "VJ-stored-key before law", setup: `UPDATE approvers SET public_key = X'00' WHERE name = 'mark-mac'`, in: expired, artifact: vjAssert(t, jgPasskey, vjIn, 2),
			step: JudgeStepStoredKey, approver: "mark-mac"},
		{name: "VJ-law-expired", in: expired, artifact: vjAssert(t, jgPasskey, vjIn, 2), step: JudgeStepLaw},
		{name: "VJ-law-mint-assertion", in: vjIn, artifact: jgPasskey.Assert(jgOrigin, jgRPID, mint, 2), step: JudgeStepLaw, cause: approval.ErrChallenge},
		{name: "VJ-law-W15", in: vjIn, artifact: vjAssert(t, jgPasskey, vjIn, 1), step: JudgeStepLaw, cause: approval.ErrSignCountRegression},
		{name: "VJ-law-origin", setup: `UPDATE webauthn_credentials SET origin = 'https://x.ts.net' WHERE name = 'mark-mac'`, in: vjIn, artifact: vjAssert(t, jgPasskey, vjIn, 2),
			step: JudgeStepLaw, cause: approval.ErrOriginMismatch},
	})
	v := vjOK(t, judge, "", vjIn, vjAssert(t, jgPasskey, vjIn, 2))
	if id, n, ok := v.PasskeyUse(); v.ApprovedBy() != "webauthn:mark-mac" || !ok || n != 2 || !bytes.Equal(id, jgPasskey.CredentialID()) {
		t.Fatalf("control: %q %x %d %v", v.ApprovedBy(), id, n, ok)
	}
}

// RecheckVerb runs after the consume: the row is consumed and may be past its expiry (no P1), and
// the stored counter is already the artifact's (no W15); a revocation since then is refused.
func TestRecheckVerb(t *testing.T) {
	judge := func(t *testing.T, ro *store.Credentials, in VerbJudgeInput, a []byte) (*approval.SignedVerb, error) {
		return RecheckVerb(boundedTestContext(t), ro, in, a)
	}
	after := vjIn
	after.Consumed, after.Now, after.BytesDigest = true, vjIn.ExpiresAt+5, [32]byte{}
	ssh := vjSign(t, jgSigner, approval.VerbNamespace, vjIn)
	asr := vjAssert(t, jgPasskey, vjIn, 2)
	moved := `UPDATE webauthn_credentials SET sign_count = 2 WHERE name = 'mark-mac'`
	otherDigest := after
	otherDigest.Digest[0] ^= 1
	vjRun(t, judge, []vjCase{
		{name: "VJ-parse neither", in: after, artifact: []byte("typed y\n"), step: JudgeStepParse, cause: ErrUnknownArtifact},
		{name: "VJ-parse an empty artifact", in: after, artifact: nil, step: JudgeStepParse, cause: ErrUnknownArtifact},
		{name: "VJ-store ssh", setup: vjDropV4, in: after, artifact: ssh, step: JudgeStepStore},
		{name: "VJ-store passkey", setup: vjDropV5, in: after, artifact: asr, step: JudgeStepStore},
		{name: "VJ-unknown-key ssh", in: after, artifact: vjSign(t, approvaltest.NewSigner(68), approval.VerbNamespace, vjIn), step: JudgeStepUnknownKey,
			cause: ErrNotEnrolled, field: approvaltest.NewSigner(68).PublicKey().Fingerprint()},
		{name: "VJ-unknown-key passkey", in: after, artifact: vjAssert(t, approvaltest.NewPasskey(approval.AlgES256, 69), vjIn, 2), step: JudgeStepUnknownKey,
			cause: ErrNotEnrolled, rawID: approvaltest.NewPasskey(approval.AlgES256, 69).CredentialID()},
		{name: "VJ-stored-key", setup: `UPDATE approvers SET public_key = X'00' WHERE name = 'mark-mac'`, in: after, artifact: asr, step: JudgeStepStoredKey, approver: "mark-mac"},
		{name: "VJ-law ssh revoked", setup: `UPDATE approvers SET revoked_at = 1601 WHERE name = 'mark-ssh'`, in: after, artifact: ssh, step: JudgeStepLaw, cause: approval.ErrApproverRevoked},
		{name: "VJ-law passkey revoked", setup: moved + `; UPDATE approvers SET revoked_at = 1601 WHERE name = 'mark-mac'`, in: after, artifact: asr, step: JudgeStepLaw,
			cause: approval.ErrApproverRevoked},
		{name: "VJ-law ssh another digest", in: otherDigest, artifact: ssh, step: JudgeStepLaw, cause: approval.ErrBadSignature},
		{name: "VJ-law passkey another digest", setup: moved, in: otherDigest, artifact: asr, step: JudgeStepLaw, cause: approval.ErrChallenge},
		{name: "VJ-law mint namespace", in: after, artifact: vjSign(t, jgSigner, approval.Namespace, vjIn), step: JudgeStepLaw, cause: approval.ErrNamespace},
	})
	if v := vjOK(t, judge, "", after, ssh); v.ApprovedBy() != "ssh:mark-ssh" {
		t.Fatalf("ssh: %q", v.ApprovedBy())
	}
	if v := vjOK(t, judge, moved, after, asr); v.ApprovedBy() != "webauthn:mark-mac" {
		t.Fatalf("passkey, counter moved: %q", v.ApprovedBy())
	}
	vjOK(t, judge, `UPDATE approvers SET enrolled_at = 1300 WHERE name = 'mark-ssh'`, after, ssh) // P5 at the caller's now (1605), not the row's creation (1000)
}
