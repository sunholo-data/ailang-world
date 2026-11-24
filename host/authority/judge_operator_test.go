package authority

// Row 243 V4 sprint 2b (design §3.1 [r1], §5.2 ConfirmBothKeys and RunnerRechecks; judge-408 note 1): the
// operator judges on judge_test.go's rig (mark-mac's passkey, counter 1, and mark-ssh, both enrolled at
// 900) with approvaltest keys (D-WORLD-109). Every approver fact the law judges is the store's: each row
// edits the store, not an argument, and wants its refusal step, its error class and no proof.

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
)

var (
	ojStranger = approvaltest.NewSigner(91)
	ojOtherKey = approvaltest.NewPasskey(approval.AlgEdDSA, 92)
)

func ojAction() approval.OperatorAction {
	a := approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: "oa-0123456789abcdef0123456789abcdef", Kind: "start", Episode: "fixture-ep",
		Store: "/ws/state/credentials.db", Options: approval.OperatorOptions{Preset: "sol", Headless: true}, Task: "t", Review: "r",
		Resolution: []approval.OperatorResolved{{Name: "agent", Value: "codex"}}}
	a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
	return a
}

// ojIn is a requested action created at 1000, expiring at 1600, judged at 1010, with its SSHSIG by
// mark-ssh and its assertion (counter 2) by mark-mac.
func ojIn(t *testing.T) (OperatorJudgeInput, []byte, []byte) {
	t.Helper()
	a := ojAction()
	d, err := approval.OperatorSignatureDigest(a)
	if err != nil {
		t.Fatal(err)
	}
	in := OperatorJudgeInput{Action: a, Digest: d, BytesDigest: d, Challenge: [32]byte{0x6f}, CreatedAt: 1000, ExpiresAt: 1600, Now: 1010}
	return in, ojSSH(t, jgSigner, approval.OperatorActionNamespace, in), ojAssert(t, jgPasskey, in, 2)
}

func ojSSH(t *testing.T, s *approvaltest.Signer, ns string, in OperatorJudgeInput) []byte {
	t.Helper()
	m, err := approval.OperatorActionMessage(in.Digest, in.Challenge, in.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	return s.Sign(ns, m)
}

func ojAssert(t *testing.T, p *approvaltest.Passkey, in OperatorJudgeInput, count uint32) []byte {
	t.Helper()
	h, err := approval.OperatorActionMessageHash(in.Digest, in.Challenge, in.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	return p.Assert(jgOrigin, jgRPID, h, count)
}

// ojOpen is the rig after setup, read-only.
func ojOpen(t *testing.T, setup ...string) *store.Credentials {
	t.Helper()
	path, _ := jgRig(t)
	jgExec(t, path, setup...)
	ro, err := store.OpenCredentialsReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ro.Close() })
	return ro
}

// ojRefused wants no proof, a *JudgeError at step, and want: a binding field, an error class or an
// admission verdict.
func ojRefused(t *testing.T, name string, v *VerifiedOperatorAction, err error, step string, want any) {
	t.Helper()
	var je *JudgeError
	var na *approval.NotAdmittedError
	ok := v == nil && errors.As(err, &je) && je.Step == step
	var be *approval.BindingError
	switch w := want.(type) {
	case *approval.BindingError:
		ok = ok && errors.As(err, &be) && be.Field == w.Field
	case error:
		ok = ok && errors.Is(err, w)
	case approval.Verdict:
		ok = ok && errors.As(err, &na) && na.Verdict == w
	}
	if !ok {
		t.Errorf("%s: %v, %v; want step %s, %v and no proof", name, v, err, step, want)
	}
}

func TestOperatorJudgeBothKeys(t *testing.T) {
	in, ssh, pk := ojIn(t)
	for by, art := range map[string][]byte{"ssh:mark-ssh": ssh, "webauthn:mark-mac": pk} {
		v, err := JudgeOperatorAction(boundedTestContext(t), ojOpen(t), in, art)
		if err != nil || v.ApprovedBy() != by || v.Binds(in.Action, in.Challenge, in.ExpiresAt) != nil || v.Rechecked() {
			t.Fatalf("%s: %v", by, err)
		}
		other := in.Action
		other.Task += "!"
		var be *approval.BindingError
		if err := v.Binds(other, in.Challenge, in.ExpiresAt); !errors.As(err, &be) || be.Field != "digest" {
			t.Errorf("%s: binds another action: %v", by, err)
		}
		cid, n, ok := v.PasskeyUse()
		if (by == "webauthn:mark-mac") != ok || (ok && (n != 2 || string(cid) != string(jgPasskey.CredentialID()))) {
			t.Errorf("%s: passkey use %x %d %v", by, cid, n, ok)
		}
	}
	// another enrolled key is judged under its own stored name
	ops := fmt.Sprintf(`INSERT INTO approvers (name, format, public_key, enrolled_at) VALUES ('ops-ssh', 'ssh', x'%x', 900)`, ojStranger.PublicKey().Blob)
	if v, err := JudgeOperatorAction(boundedTestContext(t), ojOpen(t, ops), in, ojSSH(t, ojStranger, approval.OperatorActionNamespace, in)); err != nil || v.ApprovedBy() != "ssh:ops-ssh" {
		t.Fatalf("ops-ssh: %v", err)
	}
}

// TestOperatorJudgeReadsTheStore: every approver fact is the store's (judge-408 note 1), and every
// admission input is the row's; each row is its own refusal.
func TestOperatorJudgeReadsTheStore(t *testing.T) {
	in, ssh, pk := ojIn(t)
	closed := ojOpen(t)
	_ = closed.Close()
	for _, c := range []struct {
		name  string
		setup []string
		ro    *store.Credentials
		edit  func(*OperatorJudgeInput)
		art   []byte
		step  string
		want  any
	}{
		{"typed y", nil, nil, nil, []byte("y\n"), JudgeStepParse, ErrUnknownArtifact},
		{"unknown ssh key", nil, nil, nil, ojSSH(t, ojStranger, approval.OperatorActionNamespace, in), JudgeStepUnknownKey, ErrNotEnrolled},
		{"unknown passkey credential", nil, nil, nil, ojAssert(t, ojOtherKey, in, 2), JudgeStepUnknownKey, ErrNotEnrolled},
		{"ssh store unreadable", nil, closed, nil, ssh, JudgeStepStore, nil},
		{"passkey store unreadable", nil, closed, nil, pk, JudgeStepStore, nil},
		{"ssh revoked", []string{`UPDATE approvers SET revoked_at = 1000 WHERE name = 'mark-ssh'`}, nil, nil, ssh, JudgeStepLaw, approval.ErrApproverRevoked},
		{"passkey revoked", []string{`UPDATE approvers SET revoked_at = 1000 WHERE name = 'mark-mac'`}, nil, nil, pk, JudgeStepLaw, approval.ErrApproverRevoked},
		{"ssh enrolled after now", []string{`UPDATE approvers SET enrolled_at = 1011 WHERE name = 'mark-ssh'`}, nil, nil, ssh, JudgeStepLaw, approval.ErrApproverClockBehind},
		{"passkey enrolled after now", []string{`UPDATE approvers SET enrolled_at = 1011 WHERE name = 'mark-mac'`}, nil, nil, pk, JudgeStepLaw, approval.ErrApproverClockBehind},
		{"ssh key stored as webauthn", []string{`UPDATE approvers SET format = 'webauthn' WHERE name = 'mark-ssh'`}, nil, nil, ssh, JudgeStepLaw, approval.ErrApproverFormat},
		{"passkey stored as ssh", []string{`UPDATE approvers SET format = 'ssh' WHERE name = 'mark-mac'`}, nil, nil, pk, JudgeStepLaw, approval.ErrApproverNotPasskey},
		{"stored counter at the assertion's", []string{`UPDATE webauthn_credentials SET sign_count = 2 WHERE name = 'mark-mac'`}, nil, nil, pk, JudgeStepLaw, approval.ErrSignCountRegression},
		{"stored key unparsable", []string{`UPDATE approvers SET public_key = x'00' WHERE name = 'mark-mac'`}, nil, nil, pk, JudgeStepStoredKey, nil},
		{"a verb namespace", nil, nil, nil, ojSSH(t, jgSigner, approval.VerbNamespace, in), JudgeStepLaw, approval.ErrNamespace},
		{"consumed", nil, nil, func(i *OperatorJudgeInput) { i.Consumed = true }, ssh, JudgeStepLaw, approval.RefuseConsumed},
		{"now at expiry", nil, nil, func(i *OperatorJudgeInput) { i.Now = 1600 }, ssh, JudgeStepLaw, approval.RefuseExpired},
		{"lifetime 601", nil, nil, func(i *OperatorJudgeInput) { i.CreatedAt = 999 }, ssh, JudgeStepLaw, approval.ErrOperatorLifetime},
		{"stored bytes differ", nil, nil, func(i *OperatorJudgeInput) { i.BytesDigest[0] ^= 1 }, ssh, JudgeStepLaw, approval.RefuseDigestMismatch},
		{"another action's digest column", nil, nil, func(i *OperatorJudgeInput) { i.Digest[0] ^= 1; i.BytesDigest = i.Digest }, ssh, JudgeStepLaw, &approval.BindingError{Field: "digest"}},
	} {
		ro, row := c.ro, in
		if ro == nil {
			ro = ojOpen(t, c.setup...)
		}
		if c.edit != nil {
			c.edit(&row)
		}
		v, err := JudgeOperatorAction(boundedTestContext(t), ro, row, c.art)
		ojRefused(t, c.name, v, err, c.step, c.want)
	}
}

// TestOperatorRunnerRechecks: the runner trusts no row digest, stored approver or daemon status.
func TestOperatorRunnerRechecks(t *testing.T) {
	t.Run("claimed store row", testOperatorClaimedRecheck)
	in, ssh, pk := ojIn(t)
	moved := `UPDATE webauthn_credentials SET sign_count = 2 WHERE name = 'mark-mac'`
	for by, c := range map[string]struct {
		setup []string
		art   []byte
	}{"ssh:mark-ssh": {nil, ssh}, "webauthn:mark-mac": {[]string{moved}, pk}} {
		v, err := RecheckOperatorAction(boundedTestContext(t), ojOpen(t, c.setup...), in.Action, in.Challenge, in.ExpiresAt, 1599, c.art, by)
		if err != nil || !v.Rechecked() || v.ApprovedBy() != by || v.Binds(in.Action, in.Challenge, in.ExpiresAt) != nil {
			t.Fatalf("%s control: %v", by, err)
		}
	}
	forged := in.Action
	forged.Task = "another task"
	for name, c := range map[string]struct {
		setup     []string
		action    approval.OperatorAction
		challenge [32]byte
		art       []byte
		by        string
		now       int64
		step      string
		want      any
	}{
		"a forged row's action":     {nil, forged, in.Challenge, ssh, "ssh:mark-ssh", 1500, JudgeStepLaw, approval.ErrBadSignature},
		"another challenge":         {nil, in.Action, [32]byte{0x70}, ssh, "ssh:mark-ssh", 1500, JudgeStepLaw, approval.ErrBadSignature},
		"another recorded approver": {nil, in.Action, in.Challenge, ssh, "webauthn:mark-mac", 1500, JudgeStepLaw, ErrOperatorApprover},
		"no recorded approver":      {nil, in.Action, in.Challenge, ssh, "", 1500, JudgeStepLaw, ErrOperatorApprover},
		"revoked since":             {[]string{`UPDATE approvers SET revoked_at = 1400 WHERE name = 'mark-ssh'`}, in.Action, in.Challenge, ssh, "ssh:mark-ssh", 1500, JudgeStepLaw, approval.ErrApproverRevoked},
		"expired at now":            {nil, in.Action, in.Challenge, ssh, "ssh:mark-ssh", 1600, JudgeStepLaw, approval.RefuseExpired},
		"counter never recorded":    {nil, in.Action, in.Challenge, pk, "webauthn:mark-mac", 1500, JudgeStepLaw, approval.ErrOperatorCounter},
		"unknown key":               {nil, in.Action, in.Challenge, ojSSH(t, ojStranger, approval.OperatorActionNamespace, in), "ssh:mark-ssh", 1500, JudgeStepUnknownKey, ErrNotEnrolled},
		"typed y":                   {nil, in.Action, in.Challenge, []byte("y\n"), "ssh:mark-ssh", 1500, JudgeStepParse, ErrUnknownArtifact},
	} {
		v, err := RecheckOperatorAction(boundedTestContext(t), ojOpen(t, c.setup...), c.action, c.challenge, in.ExpiresAt, c.now, c.art, c.by)
		ojRefused(t, name, v, err, c.step, c.want)
	}
}

// TestOperatorUnjudgedProofBindsNothing: what a refusal returns (nil) and a proof no judge built bind nothing.
func TestOperatorUnjudgedProofBindsNothing(t *testing.T) {
	for name, v := range map[string]*VerifiedOperatorAction{"nil": nil, "zero": {}, "unsigned, flagged": {rechecked: true}} {
		_, _, ok := v.PasskeyUse()
		if ok || !errors.Is(v.Binds(ojAction(), [32]byte{}, 0), approval.ErrNotVerified) || v.ApprovedBy() != "" || v.Rechecked() {
			t.Errorf("%s: an unjudged proof binds something", name)
		}
	}
}

// The proof can be constructed only by verifiedOperatorAction, called only by the two judges, and
// neither judge accepts an approver from its caller (the law's own input does: the positive control).
func TestVerifiedOperatorActionOnlyFromTheJudges(t *testing.T) {
	for _, typ := range []reflect.Type{reflect.TypeOf(VerifiedOperatorAction{}), reflect.TypeOf(approval.SignedOperatorAction{})} {
		for i := 0; i < typ.NumField(); i++ {
			if typ.Field(i).IsExported() {
				t.Errorf("%s.%s is exported", typ.Name(), typ.Field(i).Name)
			}
		}
	}
	approverT := reflect.TypeOf(approval.Approver{})
	holds := func(typ reflect.Type) bool {
		for i := 0; typ.Kind() == reflect.Struct && i < typ.NumField(); i++ {
			if typ.Field(i).Type == approverT {
				return true
			}
		}
		return typ == approverT
	}
	if !holds(reflect.TypeOf(approval.OperatorConfirmInput{})) {
		t.Fatal("control: the law's input holds no approver")
	}
	for _, f := range []any{JudgeOperatorAction, RecheckOperatorAction} {
		for ft, i := reflect.TypeOf(f), 0; i < ft.NumIn(); i++ {
			if holds(ft.In(i)) {
				t.Errorf("a judge takes an approver from its caller: %v", ft)
			}
		}
	}
	pkgs, err := parser.ParseDir(token.NewFileSet(), ".", func(fi fs.FileInfo) bool { return !strings.HasSuffix(fi.Name(), "_test.go") }, 0)
	if err != nil {
		t.Fatal(err)
	}
	lits, calls := 0, map[string]int{}
	for _, f := range pkgs["authority"].Files {
		for _, d := range f.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok {
				continue
			}
			ast.Inspect(fd, func(n ast.Node) bool {
				if cl, ok := n.(*ast.CompositeLit); ok {
					if id, ok := cl.Type.(*ast.Ident); ok && id.Name == "VerifiedOperatorAction" {
						lits++
						if fd.Name.Name != "verifiedOperatorAction" {
							t.Errorf("%s builds a VerifiedOperatorAction", fd.Name.Name)
						}
					}
				}
				if ce, ok := n.(*ast.CallExpr); ok {
					if id, ok := ce.Fun.(*ast.Ident); ok && id.Name == "verifiedOperatorAction" {
						calls[fd.Name.Name]++
					}
				}
				return true
			})
		}
	}
	if lits != 1 || len(calls) != 2 || calls["JudgeOperatorAction"] != 1 || calls["RecheckOperatorAction"] != 1 {
		t.Fatalf("constructor literals %d, callers %v; want 1 and the two judges (the positive control)", lits, calls)
	}
}
