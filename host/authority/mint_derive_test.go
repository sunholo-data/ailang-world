package authority

// Row 201 A1b sprint 3 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint3-plan.md §3
// M4c2, M4e): MintForRequest derives approved_by and binds the row's challenge; its caller names
// neither.

import (
	"encoding/hex"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
)

// au1Store is a fresh credentials store holding au1Request's row (created 1000, expires 1600).
func au1Store(t *testing.T) (*store.Credentials, approval.MintRequest, store.MintRequestRow) {
	t.Helper()
	st, err := store.OpenCredentials(filepath.Join(t.TempDir(), store.CredentialsFileName))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	mr, row := au1Request(t, st)
	return st, mr, row
}

// au1RequestMint is the RequestMint approve builds for au1Request's row: no approved_by.
func au1RequestMint(mr approval.MintRequest, row store.MintRequestRow) RequestMint {
	return RequestMint{Request: mr, RequestCreatedAt: 1000, Challenge: au1Challenge,
		Consume: store.RequestConsume{ID: row.RequestID, Digest: row.Digest}}
}

// au1Mint mints au1Request's session (episode ep, the profile binding, at 1010) through rm, the
// consume's clock reading 1020.
func au1Mint(t *testing.T, st *store.Credentials, rm RequestMint) error {
	t.Helper()
	good := ProfileBinding{Project: "proj", Digest: "sha256:" + strings.Repeat("c", 64)}
	rm.Consume.Clock = func() int64 { return 1020 }
	_, _, _, err := MintForRequest(boundedTestContext(t), st, "ep", good, au1Grants(5, 1610), 600, 1010, nil, rm)
	return err
}

// assertNotMinted fails unless row is still pending in st and st holds no session row.
func assertNotMinted(t *testing.T, st *store.Credentials, row store.MintRequestRow) {
	t.Helper()
	ctx := boundedTestContext(t)
	if got, ok, err := st.MintRequest(ctx, row.RequestID); err != nil || !ok || got.Consumed {
		t.Fatalf("request after a refused mint: %+v ok %v err %v, want it pending", got, ok, err)
	}
	if sessions, err := st.ListSessions(ctx, "", store.MaxListSessions); err != nil || len(sessions) != 0 {
		t.Fatalf("%d session rows after a refused mint (err %v), want 0", len(sessions), err)
	}
}

// AU1 (M4c2): a caller-set approved_by is refused (A1), a consume digest that is not 64 hex is
// refused (A2), the row's challenge is bound, and the typed mint records tty (A5).
func TestMintForRequestDerivesApprovedBy(t *testing.T) {
	for _, a := range []struct {
		name   string
		mutate func(*RequestMint)
		want   error
	}{
		{"caller-tty", func(rm *RequestMint) { rm.Consume.ApprovedBy = "tty" }, ErrApprovedBySetByCaller},
		{"caller-ssh", func(rm *RequestMint) { rm.Consume.ApprovedBy = "ssh:mallory" }, ErrApprovedBySetByCaller},
		{"digest-not-hex", func(rm *RequestMint) { rm.Consume.Digest = strings.Repeat("zz", 32) }, ErrBadConsumeDigest},
		{"digest-31-bytes", func(rm *RequestMint) { rm.Consume.Digest = strings.Repeat("ab", 31) }, ErrBadConsumeDigest},
		{"other-challenge", func(rm *RequestMint) { rm.Challenge[0] ^= 1 }, store.ErrRequestChallengeChanged},
	} {
		t.Run(a.name, func(t *testing.T) {
			st, mr, row := au1Store(t)
			rm := au1RequestMint(mr, row)
			a.mutate(&rm)
			if err := au1Mint(t, st, rm); !errors.Is(err, a.want) {
				t.Fatalf("err %v, want %v", err, a.want)
			}
			assertNotMinted(t, st, row)
		})
	}
	t.Run("typed-records-tty", func(t *testing.T) {
		st, mr, row := au1Store(t)
		if err := au1Mint(t, st, au1RequestMint(mr, row)); err != nil {
			t.Fatal(err)
		}
		got, ok, err := st.MintRequest(boundedTestContext(t), row.RequestID)
		if err != nil || !ok || !got.Consumed || got.ApprovedBy != "tty" {
			t.Fatalf("request after the typed mint: %+v ok %v err %v, want consumed by tty", got, ok, err)
		}
	})
}

// ---- M4e ----

// au1Signer is the approver key the AU tests sign with: approvaltest's fixed seed (D-WORLD-109).
var au1Signer = approvaltest.NewSigner(7)

// au1Enrol enrols "mark-mac" with au1Signer's key in st, at 900.
func au1Enrol(t *testing.T, st *store.Credentials) {
	t.Helper()
	row := store.ApproverRow{Name: "mark-mac", Format: approval.FormatSSH, PublicKey: au1Signer.PublicKey().Blob, EnrolledAt: 900}
	if err := st.EnrolApprover(boundedTestContext(t), row); err != nil {
		t.Fatal(err)
	}
}

// au1Signed is au1Signer's verified approval of (digest, challenge, expiry), admitted at 1010 for a
// request created at 1000; it fails the test unless ApproveBySignature returns one.
func au1Signed(t *testing.T, digestHex string, challenge [32]byte, expiry int64) *approval.SignedApproval {
	t.Helper()
	d, err := hex.DecodeString(digestHex)
	if err != nil || len(d) != 32 {
		t.Fatalf("digest %q: %v", digestHex, err)
	}
	digest := [32]byte(d)
	msg, err := approval.SignedMessage(digest, challenge, expiry)
	if err != nil {
		t.Fatal(err)
	}
	s, err := approval.ApproveBySignature(approval.SignatureInput{
		Admit:     approval.AdmitInput{Now: 1010, CreatedAt: 1000, ExpiresAt: expiry, StoredDigest: digest, BytesDigest: digest},
		Challenge: challenge,
		Signature: au1Signer.Sign(approval.Namespace, msg),
		Approver:  approval.Approver{Name: "mark-mac", Format: approval.FormatSSH, Key: au1Signer.PublicKey(), EnrolledAt: 900},
	})
	if err != nil || s == nil {
		t.Fatalf("ApproveBySignature: (%v, %v), want a verified approval", s, err)
	}
	return s
}

// au1SignedMint is au1RequestMint carrying au1Signer's approval of the row as stored.
func au1SignedMint(t *testing.T, mr approval.MintRequest, row store.MintRequestRow) RequestMint {
	t.Helper()
	rm := au1RequestMint(mr, row)
	rm.ExpiresAt = row.ExpiresAt
	rm.Signed = au1Signed(t, row.Digest, au1Challenge, row.ExpiresAt)
	return rm
}

// approverLastUsed is mark-mac's last_used_at in st.
func approverLastUsed(t *testing.T, st *store.Credentials) int64 {
	t.Helper()
	a, ok, err := st.Approver(boundedTestContext(t), "mark-mac")
	if err != nil || !ok {
		t.Fatalf("approver mark-mac: ok %v err %v", ok, err)
	}
	return a.LastUsedAt
}

// AU1 (M4e): a verified signature approval records approved_by=ssh:<name> and marks the approver
// used at the consume's reading (A3).
func TestMintForRequestRecordsTheSigner(t *testing.T) {
	st, mr, row := au1Store(t)
	au1Enrol(t, st)
	if err := au1Mint(t, st, au1SignedMint(t, mr, row)); err != nil {
		t.Fatal(err)
	}
	got, ok, err := st.MintRequest(boundedTestContext(t), row.RequestID)
	if err != nil || !ok || !got.Consumed || got.ApprovedBy != "ssh:mark-mac" {
		t.Fatalf("request after the signed mint: %+v ok %v err %v, want consumed by ssh:mark-mac", got, ok, err)
	}
	if at := approverLastUsed(t, st); at != 1020 {
		t.Fatalf("mark-mac's last_used_at %d, want 1020 (the consume's reading)", at)
	}
}

// AU2 (M4e): an approval made for another digest, challenge or expiry is refused by Binds (A4), a
// zero SignedApproval is ErrNotVerified, and a caller-set approver is refused (A1). Nothing is
// minted and the approver is not marked used.
func TestMintForRequestRefusesAnUnboundApproval(t *testing.T) {
	otherDigest := strings.Repeat("cd", 32)
	otherChallenge := au1Challenge
	otherChallenge[31] ^= 1
	bindingError := func(field string) func(error) bool {
		return func(err error) bool {
			var be *approval.BindingError
			return errors.As(err, &be) && be.Field == field
		}
	}
	is := func(target error) func(error) bool { return func(err error) bool { return errors.Is(err, target) } }
	t.Run("bound-control", func(t *testing.T) { // A4-EXP (sprint 4 T1): the unperturbed approval binds, so each arm below is refused for its own field
		st, mr, row := au1Store(t)
		au1Enrol(t, st)
		if err := au1Mint(t, st, au1SignedMint(t, mr, row)); err != nil {
			t.Fatalf("the bound approval: %v, want it minted", err)
		}
	})
	for _, a := range []struct {
		name   string
		mutate func(t *testing.T, rm *RequestMint, row store.MintRequestRow)
		want   func(error) bool
	}{
		{"other-digest", func(t *testing.T, rm *RequestMint, row store.MintRequestRow) {
			rm.Signed = au1Signed(t, otherDigest, au1Challenge, row.ExpiresAt)
		}, bindingError("digest")},
		{"other-challenge", func(t *testing.T, rm *RequestMint, row store.MintRequestRow) {
			rm.Signed = au1Signed(t, row.Digest, otherChallenge, row.ExpiresAt)
		}, bindingError("challenge")},
		{"other-expiry", func(t *testing.T, rm *RequestMint, row store.MintRequestRow) {
			rm.Signed = au1Signed(t, row.Digest, au1Challenge, row.ExpiresAt+1)
		}, bindingError("expiry")},
		{"expiry-not-passed", func(t *testing.T, rm *RequestMint, row store.MintRequestRow) { rm.ExpiresAt = 0 }, bindingError("expiry")},
		{"zero-approval", func(t *testing.T, rm *RequestMint, row store.MintRequestRow) { rm.Signed = &approval.SignedApproval{} }, is(approval.ErrNotVerified)},
		{"caller-approver", func(t *testing.T, rm *RequestMint, row store.MintRequestRow) {
			rm.Signed = nil
			rm.Consume.Approver = &store.ApproverUse{Name: "mark-mac", Format: approval.FormatSSH, PublicKey: au1Signer.PublicKey().Blob}
		}, is(ErrApprovedBySetByCaller)},
	} {
		t.Run(a.name, func(t *testing.T) {
			st, mr, row := au1Store(t)
			au1Enrol(t, st)
			rm := au1SignedMint(t, mr, row)
			a.mutate(t, &rm, row)
			if err := au1Mint(t, st, rm); !a.want(err) {
				t.Fatalf("err %v, want the %s refusal", err, a.name)
			}
			assertNotMinted(t, st, row)
			if at := approverLastUsed(t, st); at != 0 {
				t.Fatalf("mark-mac's last_used_at %d after a refused mint, want 0", at)
			}
		})
	}
}

// requestMintFields is where every RequestMint field comes from, and what binds it (RC5).
var requestMintFields = map[string]string{
	"Request":          "the decoded stored request; MintMatches checks the session row against it",
	"RequestCreatedAt": "the stored row's created_at; MintMatches bounds the session's created_at by it",
	"Challenge":        "the stored row's challenge (approve F4); bound in the consume's UPDATE (C6) and by Binds (A4)",
	"ExpiresAt":        "the stored row's expires_at; bound by Binds (A4)",
	"Signed":           "approval.ApproveBySignature's verified approval, or nil for the typed y; derives approved_by (A3, A5)",
	"Consume":          "the row's id and digest and env.now; approved_by and the approver are derived (A1)",
}

// RC5 (M4e): every RequestMint field has a source, and perturbing any one of them alone in a valid
// signed RequestMint refuses the mint.
func TestRequestMintFieldsAreSourced(t *testing.T) {
	typ := reflect.TypeOf(RequestMint{})
	for i := 0; i < typ.NumField(); i++ {
		if requestMintFields[typ.Field(i).Name] == "" {
			t.Errorf("RequestMint.%s has no source in requestMintFields", typ.Field(i).Name)
		}
	}
	if typ.NumField() != len(requestMintFields) {
		t.Fatalf("RequestMint has %d fields, requestMintFields %d", typ.NumField(), len(requestMintFields))
	}
	perturb := map[string]func(*RequestMint){
		"Request":          func(rm *RequestMint) { rm.Request.Episode = "other" },
		"RequestCreatedAt": func(rm *RequestMint) { rm.RequestCreatedAt = 1011 },
		"Challenge":        func(rm *RequestMint) { rm.Challenge[0] ^= 1 },
		"ExpiresAt":        func(rm *RequestMint) { rm.ExpiresAt++ },
		"Signed":           func(rm *RequestMint) { rm.Signed = &approval.SignedApproval{} },
		"Consume":          func(rm *RequestMint) { rm.Consume.ID = "mr-ffffffffffffffff" },
	}
	if len(perturb) != len(requestMintFields) {
		t.Fatalf("%d perturbations for %d fields", len(perturb), len(requestMintFields))
	}
	for field, p := range perturb {
		t.Run(field, func(t *testing.T) {
			st, mr, row := au1Store(t)
			au1Enrol(t, st)
			rm := au1SignedMint(t, mr, row)
			p(&rm)
			if err := au1Mint(t, st, rm); err == nil {
				t.Fatalf("a mint with %s perturbed succeeded", field)
			}
			assertNotMinted(t, st, row)
		})
	}
}
