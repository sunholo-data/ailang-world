package authority

// Row 201 A1b sprint 3 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint3-plan.md §3
// M4c2, M4e): MintForRequest derives approved_by and binds the row's challenge; its caller names
// neither.

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
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
