package approval_test

// Row 243 V4 sprint 2 M3 (design §3.1, §5.2 ConfirmBothKeys, ConfirmProof, ConfirmAuthenticity,
// ConfirmRevocationCounter, ConfirmBinding, ConfirmExpiry, and the pure half of RunnerRechecks): the
// operator confirmation law with approvaltest keys and passkeys (D-WORLD-109), every check refused on
// its own with its own error and no proof, against a valid control in the same test.

import (
	"bytes"
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
)

var (
	opSigner  = approvaltest.NewSigner(81)
	opPasskey = approvaltest.NewPasskey(approval.AlgEdDSA, 82)
	opC       = [32]byte{0x6f, 0x63}
)

const opCreated, opExp, opNow = int64(1000), int64(1600), int64(1010)

func opDigest(t *testing.T, a approval.OperatorAction) [32]byte {
	t.Helper()
	d, err := approval.OperatorSignatureDigest(a)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

// opSSH and opAssert sign action a's message (challenge opC, expiry exp).
func opSSH(t *testing.T, s *approvaltest.Signer, ns string, a approval.OperatorAction, exp int64) []byte {
	t.Helper()
	m, err := approval.OperatorActionMessage(opDigest(t, a), opC, exp)
	if err != nil {
		t.Fatal(err)
	}
	return s.Sign(ns, m)
}

func opAssert(t *testing.T, a approval.OperatorAction, exp int64, count uint32, e ...approvaltest.Edit) []byte {
	t.Helper()
	h, err := approval.OperatorActionMessageHash(opDigest(t, a), opC, exp)
	if err != nil {
		t.Fatal(err)
	}
	return opPasskey.Assert(paOrigin, paRPID, h, count, e...)
}

// opIn is a valid confirmation of the start action: mark-ssh (opSigner) or mark-mac (opPasskey, stored
// counter 1, assertion 2), created 1000, expires 1600, now 1010.
func opIn(t *testing.T, format string) approval.OperatorConfirmInput {
	t.Helper()
	a := opAction("start")
	d := opDigest(t, a)
	in := approval.OperatorConfirmInput{Action: a, Challenge: opC, Format: format,
		Admit: approval.AdmitInput{Now: opNow, CreatedAt: opCreated, ExpiresAt: opExp, StoredDigest: d, BytesDigest: d}}
	if format == approval.FormatSSH {
		in.Artifact, in.Approver = opSSH(t, opSigner, approval.OperatorActionNamespace, a, opExp), approval.Approver{Name: "mark-ssh",
			Format: approval.FormatSSH, Key: opSigner.PublicKey(), EnrolledAt: 900}
	} else {
		in.Artifact, in.Approver = opAssert(t, a, opExp, 2), paApprover(t, opPasskey)
	}
	return in
}

// opRefused runs the law on in and wants err class want and no proof.
func opRefused(t *testing.T, name string, in approval.OperatorConfirmInput, want error) {
	t.Helper()
	s, err := approval.ConfirmOperatorAction(in)
	if s != nil || !errors.Is(err, want) {
		t.Errorf("%s: %v, %v; want %v and no proof", name, s, err, want)
	}
}

func TestOperatorConfirmBothKeys(t *testing.T) {
	a := opAction("start")
	for format, by := range map[string]string{approval.FormatSSH: "ssh:mark-ssh", approval.FormatWebAuthn: "webauthn:mark-mac"} {
		s, err := approval.ConfirmOperatorAction(opIn(t, format))
		if err != nil || s.ApprovedBy() != by || s.Binds(a, opC, opExp) != nil {
			t.Fatalf("%s: %v %q", format, err, s.ApprovedBy())
		}
		cid, n, ok := s.PasskeyUse()
		if (format == approval.FormatWebAuthn) != ok || (ok && (n != 2 || !bytes.Equal(cid, opPasskey.CredentialID()))) {
			t.Fatalf("%s: passkey use %x %d %v", format, cid, n, ok)
		}
		var be *approval.BindingError
		other := a
		other.Task += "!"
		for field, err := range map[string]error{"digest": s.Binds(other, opC, opExp), "challenge": s.Binds(a, [32]byte{1}, opExp),
			"expiry": s.Binds(a, opC, opExp+1)} {
			if !errors.As(err, &be) || be.Field != field {
				t.Errorf("%s: Binds another %s: %v", format, field, err)
			}
		}
		if err := s.Binds(approval.OperatorAction{}, opC, opExp); !errors.Is(err, approval.ErrSchema) {
			t.Errorf("%s: Binds an invalid action: %v", format, err)
		}
	}
	var zero *approval.SignedOperatorAction
	for _, s := range []*approval.SignedOperatorAction{zero, {}} {
		_, _, ok := s.PasskeyUse()
		if !errors.Is(s.Binds(a, opC, opExp), approval.ErrNotVerified) || s.ApprovedBy() != "" || ok {
			t.Fatal("an unverified proof binds something")
		}
	}
}

func TestOperatorConfirmProof(t *testing.T) {
	ssh, pk := opIn(t, approval.FormatSSH), opIn(t, approval.FormatWebAuthn)
	edit := func(in approval.OperatorConfirmInput, f func(*approval.OperatorConfirmInput)) approval.OperatorConfirmInput {
		f(&in)
		return in
	}
	opRefused(t, "typed y", edit(ssh, func(in *approval.OperatorConfirmInput) { in.Artifact = []byte("y\n") }), approval.ErrArmor)
	opRefused(t, "no proof (ssh)", edit(ssh, func(in *approval.OperatorConfirmInput) { in.Artifact = nil }), approval.ErrArmor)
	opRefused(t, "no proof (passkey)", edit(pk, func(in *approval.OperatorConfirmInput) { in.Artifact = nil }), approval.ErrPasskeyResponse)
	opRefused(t, "malformed armor", edit(ssh, func(in *approval.OperatorConfirmInput) { in.Artifact = in.Artifact[:len(in.Artifact)/2] }), approval.ErrArmor)
	opRefused(t, "malformed assertion", edit(pk, func(in *approval.OperatorConfirmInput) { in.Artifact = []byte("{}") }), approval.ErrPasskeyResponse)
	opRefused(t, "an assertion as ssh", edit(ssh, func(in *approval.OperatorConfirmInput) { in.Artifact = pk.Artifact }), approval.ErrArmor)
	opRefused(t, "an SSHSIG as passkey", edit(pk, func(in *approval.OperatorConfirmInput) { in.Artifact = ssh.Artifact }), approval.ErrPasskeyResponse)
	opRefused(t, "no format", edit(ssh, func(in *approval.OperatorConfirmInput) { in.Format = "" }), approval.ErrVerbFormat)
	opRefused(t, "typed format", edit(ssh, func(in *approval.OperatorConfirmInput) { in.Format = "typed" }), approval.ErrVerbFormat)
	opRefused(t, "bad challenge (ssh)", edit(ssh, func(in *approval.OperatorConfirmInput) { in.Challenge[0] ^= 1 }), approval.ErrBadSignature)
	opRefused(t, "bad challenge (passkey)", edit(pk, func(in *approval.OperatorConfirmInput) { in.Challenge[0] ^= 1 }), approval.ErrChallenge)
	opRefused(t, "an invalid action", edit(ssh, func(in *approval.OperatorConfirmInput) { in.Action.ID = "" }), approval.ErrOperatorAction)
}

func TestOperatorConfirmAuthenticity(t *testing.T) {
	a := opAction("start")
	ssh, pk := opIn(t, approval.FormatSSH), opIn(t, approval.FormatWebAuthn)
	for _, tc := range []struct {
		name string
		in   approval.OperatorConfirmInput
		want error
	}{
		{"another key's SSHSIG", func() approval.OperatorConfirmInput {
			in := ssh
			in.Artifact = opSSH(t, approvaltest.NewSigner(83), approval.OperatorActionNamespace, a, opExp)
			return in
		}(), approval.ErrKeyMismatch},
		{"forged assertion", func() approval.OperatorConfirmInput {
			in := pk
			in.Artifact = opAssert(t, a, opExp, 2, approvaltest.Edit{BadSignature: true})
			return in
		}(), approval.ErrBadSignature},
		{"unknown credential", func() approval.OperatorConfirmInput {
			in := pk
			in.Artifact = opAssert(t, a, opExp, 2, approvaltest.Edit{RawID: []byte("another-credential")})
			return in
		}(), approval.ErrCredentialIDMismatch},
	} {
		opRefused(t, tc.name, tc.in, tc.want)
	}
	for _, ns := range []string{approval.VerbNamespace, approval.Namespace, approval.DecisionNamespace, "git", "world-operator-action@sunholo.x"} {
		in := ssh
		in.Artifact = opSSH(t, opSigner, ns, a, opExp)
		opRefused(t, "namespace "+ns, in, approval.ErrNamespace)
	}
	h, _ := approval.OperatorActionMessageHash(opDigest(t, a), opC, opExp)
	for _, tc := range []struct {
		name, origin, rpID string
		e                  approvaltest.Edit
		want               error
	}{
		{"no UV", paOrigin, paRPID, approvaltest.Edit{Clear: 0x04}, approval.ErrUserVerification},
		{"no UP", paOrigin, paRPID, approvaltest.Edit{Clear: 0x01}, approval.ErrUserPresence},
		{"another RP id", paOrigin, "evil.ts.net", approvaltest.Edit{}, approval.ErrRPIDHash},
		{"another origin", "https://x.ts.net:8444", paRPID, approvaltest.Edit{}, approval.ErrOriginMismatch},
		{"cross-origin", paOrigin, paRPID, approvaltest.Edit{CrossOrigin: true}, approval.ErrCrossOrigin},
		{"a registration", paOrigin, paRPID, approvaltest.Edit{Type: "webauthn.create"}, approval.ErrClientDataType},
	} {
		in := pk
		in.Artifact = opPasskey.Assert(tc.origin, tc.rpID, h, 2, tc.e)
		opRefused(t, tc.name, in, tc.want)
	}
}

func TestOperatorConfirmRevocationCounter(t *testing.T) {
	ssh, pk := opIn(t, approval.FormatSSH), opIn(t, approval.FormatWebAuthn)
	a := opAction("start")
	revoked := ssh
	revoked.Approver.RevokedAt = 950
	opRefused(t, "revoked ssh", revoked, approval.ErrApproverRevoked)
	pr := opIn(t, approval.FormatWebAuthn)
	pr.Approver.RevokedAt = 950
	opRefused(t, "revoked passkey", pr, approval.ErrApproverRevoked)
	late := ssh
	late.Approver.EnrolledAt = opNow + 1
	opRefused(t, "enrolled after now", late, approval.ErrApproverClockBehind)
	if late.Approver.EnrolledAt = opCreated + 5; func() error { _, err := approval.ConfirmOperatorAction(late); return err }() != nil {
		t.Error("enrolled after the request, before now: refused") // P5 at now, not at the row's creation
	}
	asSSH := pk
	asSSH.Approver.Format = approval.FormatSSH
	opRefused(t, "a passkey approver that is not a passkey", asSSH, approval.ErrApproverNotPasskey)
	for _, c := range []struct {
		stored, asserted uint32
		ok               bool
	}{{1, 1, false}, {5, 2, false}, {5, 0, false}, {1, 2, true}, {0, 0, true}, {0, 1, true}} {
		in := opIn(t, approval.FormatWebAuthn)
		in.Approver.Passkey.SignCount = c.stored
		in.Artifact = opAssert(t, a, opExp, c.asserted)
		s, err := approval.ConfirmOperatorAction(in)
		if c.ok != (err == nil) || (!c.ok && (s != nil || !errors.Is(err, approval.ErrSignCountRegression))) {
			t.Errorf("counter stored %d asserted %d: %v", c.stored, c.asserted, err)
		}
	}
}

func TestOperatorConfirmBinding(t *testing.T) {
	edits := opFieldEdits()
	edits["id"] = func(a *approval.OperatorAction) { a.ID = "oa-ffffffffffffffffffffffffffffffff" }
	edits["kind"] = func(a *approval.OperatorAction) { *a = opAction("finish") }
	for _, format := range []string{approval.FormatSSH, approval.FormatWebAuthn} {
		want := map[string]error{approval.FormatSSH: approval.ErrBadSignature, approval.FormatWebAuthn: approval.ErrChallenge}[format]
		for name, edit := range edits {
			in := opIn(t, format) // signs the start action
			edit(&in.Action)
			d := opDigest(t, in.Action)
			in.Admit.StoredDigest, in.Admit.BytesDigest = d, d // a consistent row for another action
			opRefused(t, format+" another "+name, in, want)
		}
		in := opIn(t, format)
		in.Admit.StoredDigest = opDigest(t, opAction("finish"))
		in.Admit.BytesDigest = in.Admit.StoredDigest
		var be *approval.BindingError
		if s, err := approval.ConfirmOperatorAction(in); s != nil || !errors.As(err, &be) || be.Field != "digest" {
			t.Errorf("%s: the row's digest is another action's: %v", format, err)
		}
		in = opIn(t, format)
		if format == approval.FormatSSH {
			in.Artifact = opSSH(t, opSigner, approval.OperatorActionNamespace, opAction("start"), opExp+1)
		} else {
			in.Artifact = opAssert(t, opAction("start"), opExp+1, 2)
		}
		opRefused(t, format+" another expiry", in, want)
	}
	a := opAction("start")
	d := opDigest(t, a)
	vm, _ := approval.VerbMessage(d, opC, opExp)
	mm, _ := approval.SignedMessage(d, opC, opExp)
	for name, sig := range map[string][]byte{"a V3 verb SSHSIG": opSigner.Sign(approval.VerbNamespace, vm), "an A3 mint SSHSIG": opSigner.Sign(approval.Namespace, mm)} {
		in := opIn(t, approval.FormatSSH)
		in.Artifact = sig
		opRefused(t, name, in, approval.ErrNamespace)
	}
	in := opIn(t, approval.FormatSSH)
	in.Artifact = opSigner.Sign(approval.OperatorActionNamespace, vm)
	opRefused(t, "the verb message under the operator namespace", in, approval.ErrBadSignature)
	vh, _ := approval.VerbMessageHash(d, opC, opExp)
	in = opIn(t, approval.FormatWebAuthn)
	in.Artifact = opPasskey.Assert(paOrigin, paRPID, vh, 2)
	opRefused(t, "a V3 verb assertion", in, approval.ErrChallenge)
}

func TestOperatorConfirmExpiry(t *testing.T) {
	a := opAction("start")
	for _, format := range []string{approval.FormatSSH, approval.FormatWebAuthn} {
		for _, tc := range []struct {
			name                string
			now, created, exp   int64
			consumed, otherByte bool
			want                error
		}{
			{"now at expiry", opExp, opCreated, opExp, false, false, &approval.NotAdmittedError{Verdict: approval.RefuseExpired}},
			{"now past expiry", opExp + 1, opCreated, opExp, false, false, &approval.NotAdmittedError{Verdict: approval.RefuseExpired}},
			{"now one before expiry", opExp - 1, opCreated, opExp, false, false, nil},
			{"now before creation", opCreated - 1, opCreated, opExp, false, false, &approval.NotAdmittedError{Verdict: approval.RefuseClockBehind}},
			{"now at creation", opCreated, opCreated, opExp, false, false, nil},
			{"lifetime 600", opNow, opCreated, opCreated + 600, false, false, nil},
			{"lifetime 300", opNow, opCreated, opCreated + 300, false, false, nil},
			{"lifetime 601", opNow, opCreated, opCreated + 601, false, false, approval.ErrOperatorLifetime},
			{"consumed", opNow, opCreated, opExp, true, false, &approval.NotAdmittedError{Verdict: approval.RefuseConsumed}},
			{"stored bytes differ", opNow, opCreated, opExp, false, true, &approval.NotAdmittedError{Verdict: approval.RefuseDigestMismatch}},
		} {
			in := opIn(t, format)
			in.Admit.Now, in.Admit.CreatedAt, in.Admit.ExpiresAt, in.Admit.Consumed = tc.now, tc.created, tc.exp, tc.consumed
			if tc.otherByte {
				in.Admit.BytesDigest[0] ^= 1
			}
			if format == approval.FormatSSH {
				in.Artifact = opSSH(t, opSigner, approval.OperatorActionNamespace, a, tc.exp)
			} else {
				in.Artifact = opAssert(t, a, tc.exp, 2)
			}
			s, err := approval.ConfirmOperatorAction(in)
			var na, wna *approval.NotAdmittedError
			switch {
			case tc.want == nil:
				if err != nil || s.Binds(a, opC, tc.exp) != nil {
					t.Errorf("%s %s: %v", format, tc.name, err)
				}
			case errors.As(tc.want, &wna):
				if s != nil || !errors.As(err, &na) || na.Verdict != wna.Verdict {
					t.Errorf("%s %s: %v, want %v", format, tc.name, err, wna.Verdict)
				}
			default:
				if s != nil || !errors.Is(err, tc.want) {
					t.Errorf("%s %s: %v, want %v", format, tc.name, err, tc.want)
				}
			}
		}
	}
}

// A negative expiry has no message: refused as such, never verified over an empty message.
func TestOperatorConfirmNegativeExpiry(t *testing.T) {
	for format, want := range map[string]error{approval.FormatSSH: approval.ErrNegativeExpiry, approval.FormatWebAuthn: approval.ErrChallenge} {
		in := opIn(t, format)
		in.Admit.Now, in.Admit.CreatedAt, in.Admit.ExpiresAt, in.Approver.EnrolledAt = -50, -100, -5, -100
		opRefused(t, format+" negative expiry", in, want)
	}
}

// The runner's recheck from its own frozen action: no admission row, the counter already recorded,
// the key and the expiry judged at now.
func TestOperatorRunnerRecheckLaw(t *testing.T) {
	a := opAction("start")
	check := func(format string) approval.OperatorCheckInput {
		in := opIn(t, format)
		c := approval.OperatorCheckInput{Action: a, Challenge: opC, Expiry: opExp, Now: opExp - 1, Artifact: in.Artifact, Format: format, Approver: in.Approver}
		if format == approval.FormatWebAuthn {
			c.Approver.Passkey.SignCount = 2 // the confirmation's counter, recorded
		}
		return c
	}
	for _, format := range []string{approval.FormatSSH, approval.FormatWebAuthn} {
		if s, err := approval.CheckOperatorAction(check(format)); err != nil || s.Binds(a, opC, opExp) != nil {
			t.Fatalf("%s control: %v", format, err)
		}
		for _, tc := range []struct {
			name string
			edit func(*approval.OperatorCheckInput)
			want error
		}{
			{"now at expiry", func(c *approval.OperatorCheckInput) { c.Now = opExp }, &approval.NotAdmittedError{}},
			{"revoked since", func(c *approval.OperatorCheckInput) { c.Approver.RevokedAt = 1500 }, approval.ErrApproverRevoked},
			{"another frozen action", func(c *approval.OperatorCheckInput) { c.Action.Task += "!" }, map[string]error{approval.FormatSSH: approval.ErrBadSignature, approval.FormatWebAuthn: approval.ErrChallenge}[format]},
			{"another challenge", func(c *approval.OperatorCheckInput) { c.Challenge[1] ^= 1 }, map[string]error{approval.FormatSSH: approval.ErrBadSignature, approval.FormatWebAuthn: approval.ErrChallenge}[format]},
			{"an invalid action", func(c *approval.OperatorCheckInput) { c.Action.Store = "x" }, approval.ErrPath},
			{"no format", func(c *approval.OperatorCheckInput) { c.Format = "" }, approval.ErrVerbFormat},
		} {
			c := check(format)
			tc.edit(&c)
			var na *approval.NotAdmittedError
			s, err := approval.CheckOperatorAction(c)
			if _, isNA := tc.want.(*approval.NotAdmittedError); s != nil || (isNA && (!errors.As(err, &na) || na.Verdict != approval.RefuseExpired)) || (!isNA && !errors.Is(err, tc.want)) {
				t.Errorf("%s %s: %v, %v", format, tc.name, s, err)
			}
		}
	}
	for _, c := range []struct {
		stored uint32
		ok     bool
	}{{1, false}, {2, true}, {3, true}} {
		in := check(approval.FormatWebAuthn)
		in.Approver.Passkey.SignCount = c.stored
		s, err := approval.CheckOperatorAction(in)
		if c.ok != (err == nil) || (!c.ok && (s != nil || !errors.Is(err, approval.ErrOperatorCounter))) {
			t.Errorf("recheck, stored counter %d against the assertion's 2: %v", c.stored, err)
		}
	}
}
