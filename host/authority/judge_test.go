package authority

// Row 201 A3 sprint 1 M2 (design_docs/planned/w-world-approvals-by-key-a3.md §2.5, §6 JG1, JG2): the
// shared judge, over the CLI's refusal matrix (WF3-WF6, F11-F14). The CLI's texts for each step are
// pinned, unedited, by cmd/ailang-worldd's TestApprovePasskeyRefusalsBeforeAnyEffect and
// TestApproveBySignatureRefusalsBeforeAnyEffect. Every key is approvaltest's (D-WORLD-109).

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
)

const jgOrigin, jgRPID = "http://localhost:7644", "localhost"

var (
	jgPasskey = approvaltest.NewPasskey(approval.AlgEdDSA, 61)
	jgSigner  = approvaltest.NewSigner(62)
)

// jgRig is a schema-7 store at path holding au1Request's row, jgPasskey enrolled as mark-mac
// (counter 1) and jgSigner as mark-ssh, both at 900.
func jgRig(t *testing.T) (string, store.MintRequestRow) {
	t.Helper()
	path := filepath.Join(t.TempDir(), store.CredentialsFileName)
	st, err := store.OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = st.Close() }()
	_, row := au1Request(t, st)
	k := jgPasskey.Key()
	ctx := boundedTestContext(t)
	if err := st.EnrolPasskey(ctx, store.ApproverRow{Name: "mark-mac", Format: approval.FormatWebAuthn, PublicKey: k.Canonical(), EnrolledAt: 900},
		store.PasskeyRow{Name: "mark-mac", RPID: jgRPID, Origin: jgOrigin, CredentialID: jgPasskey.CredentialID(), AAGUID: make([]byte, 16), Alg: k.Alg, SignCount: 1}); err != nil {
		t.Fatal(err)
	}
	if err := st.EnrolApprover(ctx, store.ApproverRow{Name: "mark-ssh", Format: "ssh", PublicKey: jgSigner.PublicKey().Blob, EnrolledAt: 900}); err != nil {
		t.Fatal(err)
	}
	return path, row
}

// jgExec runs stmts on the store at path.
func jgExec(t *testing.T, path string, stmts ...string) {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			t.Fatal(err)
		}
	}
}

// jgHash and jgMessage are what an approval of row signs.
func jgHash(t *testing.T, row store.MintRequestRow) [32]byte {
	t.Helper()
	d, _ := hex.DecodeString(row.Digest)
	c, _ := hex.DecodeString(row.Challenge)
	h, err := approval.MessageHash([32]byte(d), [32]byte(c), row.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func jgMessage(t *testing.T, row store.MintRequestRow) []byte {
	t.Helper()
	d, _ := hex.DecodeString(row.Digest)
	c, _ := hex.DecodeString(row.Challenge)
	m, err := approval.SignedMessage([32]byte(d), [32]byte(c), row.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

// judgeArm is one refusal arm: what to change, the artifact, and the step and cause it must name.
type judgeArm struct {
	name     string
	setup    []string                    // SQL on the store before the judgement
	edit     func(*store.MintRequestRow) // the row the judge is handed
	artifact func(t *testing.T, row store.MintRequestRow) []byte
	now      int64 // 0: 1010
	step     string
	cause    error  // errors.Is target; nil: only the step
	field    string // the JudgeError field the CLI prints, when the step has one
}

// runJudgeArms judges each arm on a fresh rig through judge, and requires the step, the cause and
// the CLI's field; a verified control on the same rig binds the row.
func runJudgeArms(t *testing.T, judge func(ro *store.Credentials, row store.MintRequestRow, artifact []byte, now int64) (*approval.SignedApproval, error),
	arms []judgeArm, field func(*JudgeError) string, good func(t *testing.T, row store.MintRequestRow) []byte, approvedBy string) {
	t.Helper()
	for _, arm := range arms {
		t.Run(arm.name, func(t *testing.T) {
			path, row := jgRig(t)
			jgExec(t, path, arm.setup...)
			art := arm.artifact(t, row)
			if arm.edit != nil {
				arm.edit(&row)
			}
			now := arm.now
			if now == 0 {
				now = 1010
			}
			ro, err := store.OpenCredentialsReadOnly(path)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = ro.Close() }()
			signed, err := judge(ro, row, art, now)
			var je *JudgeError
			if signed != nil || !errors.As(err, &je) || je.Step != arm.step {
				t.Fatalf("judge = %v, %v; want a %q JudgeError", signed, err, arm.step)
			}
			if arm.cause != nil && !errors.Is(err, arm.cause) {
				t.Fatalf("judge = %v, want it to wrap %v", err, arm.cause)
			}
			if got := field(je); got != arm.field {
				t.Fatalf("the CLI's field = %q, want %q", got, arm.field)
			}
		})
	}
	t.Run("verified", func(t *testing.T) {
		path, row := jgRig(t)
		ro, err := store.OpenCredentialsReadOnly(path)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = ro.Close() }()
		signed, err := judge(ro, row, good(t, row), 1010)
		if err != nil || signed == nil || signed.ApprovedBy() != approvedBy {
			t.Fatalf("judge = %v, %v; want approved by %s", signed, err, approvedBy)
		}
		d, _ := hex.DecodeString(row.Digest)
		c, _ := hex.DecodeString(row.Challenge)
		if err := signed.Binds([32]byte(d), [32]byte(c), row.ExpiresAt); err != nil {
			t.Fatalf("the verified approval does not bind the row: %v", err)
		}
	})
}

// requestArms are the row arms both judges share: a row whose digest or challenge is not 64 hex is
// refused before any parse or lookup; a consumed row, or one whose bytes are not its digest, is
// refused by the law's admission.
func requestArms(good func(t *testing.T, row store.MintRequestRow) []byte) []judgeArm {
	return []judgeArm{
		{name: "request-digest", edit: func(r *store.MintRequestRow) { r.Digest = strings.Repeat("zz", 32) }, artifact: good, step: JudgeStepRequest},
		{name: "request-challenge", edit: func(r *store.MintRequestRow) { r.Challenge = "abcd" }, artifact: good, step: JudgeStepRequest},
		{name: "law-consumed", edit: func(r *store.MintRequestRow) { r.Consumed = true }, artifact: good, step: JudgeStepLaw},
		{name: "law-bytes-changed", edit: func(r *store.MintRequestRow) { r.Request = append(append([]byte(nil), r.Request...), ' ') }, artifact: good, step: JudgeStepLaw},
	}
}

// JG1: JudgeAssertion's step, cause and field for each of the CLI's WF3-WF6 refusals, and a
// verified assertion that binds the row.
func TestJudgeAssertionIsTheCLIJudge(t *testing.T) {
	good := func(t *testing.T, row store.MintRequestRow) []byte {
		return jgPasskey.Assert(jgOrigin, jgRPID, jgHash(t, row), 2)
	}
	stranger := approvaltest.NewPasskey(approval.AlgEdDSA, 63)
	arms := append(requestArms(good),
		judgeArm{name: "WF3-not-an-assertion", artifact: func(*testing.T, store.MintRequestRow) []byte { return []byte(`{"id":"x"}`) },
			step: JudgeStepParse, cause: approval.ErrPasskeyResponse},
		judgeArm{name: "WF4-v5-store", setup: []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`}, artifact: good, step: JudgeStepStore},
		judgeArm{name: "WF5-unknown-credential", artifact: func(t *testing.T, row store.MintRequestRow) []byte {
			return stranger.Assert(jgOrigin, jgRPID, jgHash(t, row), 2)
		}, step: JudgeStepUnknownKey, cause: ErrNotEnrolled, field: string(stranger.CredentialID())},
		judgeArm{name: "WF5-stored-key-invalid", setup: []string{`UPDATE approvers SET public_key = x'a0' WHERE name = 'mark-mac'`}, artifact: good,
			step: JudgeStepStoredKey, field: "mark-mac"},
		judgeArm{name: "WF6-counter-regression", artifact: func(t *testing.T, row store.MintRequestRow) []byte {
			return jgPasskey.Assert(jgOrigin, jgRPID, jgHash(t, row), 1)
		}, step: JudgeStepLaw, cause: approval.ErrSignCountRegression},
		judgeArm{name: "WF6-other-request", artifact: func(t *testing.T, row store.MintRequestRow) []byte {
			h := jgHash(t, row)
			h[31] ^= 1
			return jgPasskey.Assert(jgOrigin, jgRPID, h, 2)
		}, step: JudgeStepLaw, cause: approval.ErrChallenge},
		judgeArm{name: "WF6-revoked", setup: []string{`UPDATE approvers SET revoked_at = 950 WHERE name = 'mark-mac'`}, artifact: good,
			step: JudgeStepLaw, cause: approval.ErrApproverRevoked},
		judgeArm{name: "WF6-another-origin", setup: []string{`UPDATE webauthn_credentials SET origin = 'http://localhost:7645' WHERE name = 'mark-mac'`}, artifact: good,
			step: JudgeStepLaw, cause: approval.ErrOriginMismatch},
		judgeArm{name: "law-at-now-expired", artifact: good, now: 1600, step: JudgeStepLaw},
	)
	field := func(je *JudgeError) string {
		switch je.Step {
		case JudgeStepUnknownKey:
			return string(je.RawID)
		case JudgeStepStoredKey:
			return je.Approver
		}
		return ""
	}
	runJudgeArms(t, func(ro *store.Credentials, row store.MintRequestRow, artifact []byte, now int64) (*approval.SignedApproval, error) {
		return JudgeAssertion(boundedTestContext(t), ro, row, artifact, now)
	}, arms, field, good, "webauthn:mark-mac")
}

// JG2: JudgeSSHSignature's step, cause and field for each of the CLI's F11-F14 refusals, and a
// verified signature that binds the row.
func TestJudgeSSHSignatureIsTheCLIJudge(t *testing.T) {
	good := func(t *testing.T, row store.MintRequestRow) []byte {
		return jgSigner.Sign(approval.Namespace, jgMessage(t, row))
	}
	stranger := approvaltest.NewSigner(64)
	arms := append(requestArms(good),
		judgeArm{name: "F11-garbage", artifact: func(*testing.T, store.MintRequestRow) []byte { return []byte("not a signature\n") },
			step: JudgeStepParse, cause: approval.ErrArmor},
		judgeArm{name: "F12-v4-store", setup: []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`},
			artifact: good, step: JudgeStepStore},
		judgeArm{name: "F13-unenrolled-key", artifact: func(t *testing.T, row store.MintRequestRow) []byte {
			return stranger.Sign(approval.Namespace, jgMessage(t, row))
		}, step: JudgeStepUnknownKey, cause: ErrNotEnrolled, field: stranger.PublicKey().Fingerprint()},
		judgeArm{name: "F14-other-request", artifact: func(t *testing.T, row store.MintRequestRow) []byte {
			m := jgMessage(t, row)
			m[len(m)-9] ^= 1
			return jgSigner.Sign(approval.Namespace, m)
		}, step: JudgeStepLaw, cause: approval.ErrBadSignature},
		judgeArm{name: "F14-git-namespace", artifact: func(t *testing.T, row store.MintRequestRow) []byte { return jgSigner.Sign("git", jgMessage(t, row)) },
			step: JudgeStepLaw, cause: approval.ErrNamespace},
		judgeArm{name: "F14-revoked", setup: []string{`UPDATE approvers SET revoked_at = 950 WHERE name = 'mark-ssh'`}, artifact: good,
			step: JudgeStepLaw, cause: approval.ErrApproverRevoked},
		judgeArm{name: "F14-format-not-ssh", setup: []string{`UPDATE approvers SET format = 'webauthn' WHERE name = 'mark-ssh'`}, artifact: good, step: JudgeStepLaw},
		judgeArm{name: "law-at-now-expired", artifact: good, now: 1600, step: JudgeStepLaw},
	)
	field := func(je *JudgeError) string {
		if je.Step == JudgeStepUnknownKey {
			return je.Fingerprint
		}
		return ""
	}
	runJudgeArms(t, func(ro *store.Credentials, row store.MintRequestRow, artifact []byte, now int64) (*approval.SignedApproval, error) {
		return JudgeSSHSignature(boundedTestContext(t), ro, row, artifact, now)
	}, arms, field, good, "ssh:mark-ssh")
}
