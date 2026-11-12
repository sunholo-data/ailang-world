package authority

// Row 201 A4 sprint 2 M4f (design §2.7, §6 JD1-JD3, SSH half): JudgeDecisionSSH's steps, causes and
// fields agree with the mint judge's (JudgeSSHSignature, JG2): parse, store, unknown key, law; a
// verified decision binds its subject. Every key is approvaltest's (D-WORLD-109).

import (
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
)

var jdIn = DecisionJudgeInput{Verb: approval.VerbDeny, Packet: [32]byte{0x50}, Request: [32]byte{0x51}, Expiry: 1760, Now: 1710}

// jdSign is s's SSHSIG over in's decision message under namespace.
func jdSign(t *testing.T, s *approvaltest.Signer, namespace string, in DecisionJudgeInput) []byte {
	t.Helper()
	msg, err := approval.DecisionMessage(in.Verb, in.Packet, in.Request, in.Expiry)
	if err != nil {
		t.Fatal(err)
	}
	return s.Sign(namespace, msg)
}

func TestJudgeDecisionSSH(t *testing.T) {
	good := jdSign(t, jgSigner, approval.DecisionNamespace, jdIn)
	stranger := approvaltest.NewSigner(65)
	mintMsg, _ := approval.SignedMessage(jdIn.Request, jdIn.Packet, jdIn.Expiry)
	for _, tc := range []struct {
		name     string
		setup    []string
		artifact []byte
		in       DecisionJudgeInput
		step     string
		cause    error
		field    string
	}{
		{name: "JD-parse", artifact: []byte("not a signature\n"), in: jdIn, step: JudgeStepParse, cause: approval.ErrArmor},
		{name: "JD-store-v4", setup: []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`},
			artifact: good, in: jdIn, step: JudgeStepStore},
		{name: "JD-unknown-key", artifact: jdSign(t, stranger, approval.DecisionNamespace, jdIn), in: jdIn, step: JudgeStepUnknownKey, cause: ErrNotEnrolled,
			field: stranger.PublicKey().Fingerprint()},
		{name: "JD-law-mint-signature", artifact: jgSigner.Sign(approval.Namespace, mintMsg), in: jdIn, step: JudgeStepLaw, cause: approval.ErrNamespace},
		{name: "JD-law-other-verb", artifact: good, in: DecisionJudgeInput{Verb: approval.VerbApprove, Packet: jdIn.Packet, Request: jdIn.Request, Expiry: jdIn.Expiry, Now: jdIn.Now},
			step: JudgeStepLaw, cause: approval.ErrBadSignature},
		{name: "JD-law-format", setup: []string{`UPDATE approvers SET format = 'webauthn' WHERE name = 'mark-ssh'`}, artifact: good, in: jdIn, step: JudgeStepLaw},
		{name: "JD-law-revoked", setup: []string{`UPDATE approvers SET revoked_at = 950 WHERE name = 'mark-ssh'`}, artifact: good, in: jdIn, step: JudgeStepLaw, cause: approval.ErrApproverRevoked},
		{name: "JD-law-clock-behind", artifact: good, in: DecisionJudgeInput{Verb: jdIn.Verb, Packet: jdIn.Packet, Request: jdIn.Request, Expiry: jdIn.Expiry, Now: 899},
			step: JudgeStepLaw, cause: approval.ErrApproverClockBehind},
	} {
		path, _ := jgRig(t)
		jgExec(t, path, tc.setup...)
		ro, err := store.OpenCredentialsReadOnly(path)
		if err != nil {
			t.Fatal(err)
		}
		d, err := JudgeDecisionSSH(boundedTestContext(t), ro, tc.in, tc.artifact)
		_ = ro.Close()
		var je *JudgeError
		if d != nil || !errors.As(err, &je) || je.Step != tc.step || (tc.cause != nil && !errors.Is(err, tc.cause)) || je.Fingerprint != tc.field {
			t.Errorf("%s: %v, %v; want step %s cause %v field %q", tc.name, d, err, tc.step, tc.cause, tc.field)
		}
	}
	path, _ := jgRig(t)
	ro, err := store.OpenCredentialsReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ro.Close() }()
	d, err := JudgeDecisionSSH(boundedTestContext(t), ro, jdIn, good)
	if err != nil || d.DecidedBy() != "ssh:mark-ssh" || d.Binds(jdIn.Verb, jdIn.Packet, jdIn.Request, jdIn.Expiry) != nil {
		t.Fatalf("control: %v, %v", d, err)
	}
}
