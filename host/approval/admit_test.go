package approval

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestAdmitLaw(t *testing.T) {
	base := AdmitInput{Now: 1000, CreatedAt: 900, ExpiresAt: 1500, StoredDigest: [32]byte{1}, BytesDigest: [32]byte{1}}
	arms := []struct {
		name   string
		mutate func(*AdmitInput)
		want   Verdict
	}{
		{"clean", func(in *AdmitInput) {}, Admitted},
		{"at-expiry", func(in *AdmitInput) { in.Now = 1500 }, RefuseExpired},
		{"expiry-minus-1", func(in *AdmitInput) { in.Now = 1499 }, Admitted},
		{"at-created", func(in *AdmitInput) { in.Now = 900 }, Admitted},
		{"clock-behind", func(in *AdmitInput) { in.Now = 899 }, RefuseClockBehind},
		{"consumed", func(in *AdmitInput) { in.Consumed = true }, RefuseConsumed},
		{"consumed+expired", func(in *AdmitInput) { in.Consumed = true; in.Now = 2000 }, RefuseConsumed},
		{"consumed+behind", func(in *AdmitInput) { in.Consumed = true; in.Now = 1 }, RefuseConsumed},
		{"digest", func(in *AdmitInput) { in.BytesDigest = [32]byte{2} }, RefuseDigestMismatch},
		{"digest+expired", func(in *AdmitInput) { in.BytesDigest = [32]byte{2}; in.Now = 1500 }, RefuseExpired},
		{"digest+behind", func(in *AdmitInput) { in.BytesDigest = [32]byte{2}; in.Now = 899 }, RefuseClockBehind},
	}
	for _, arm := range arms {
		in := base
		arm.mutate(&in)
		if got := Admit(in); got != arm.want {
			t.Errorf("%s: Admit = %s, want %s", arm.name, got, arm.want)
		}
	}

	names := map[Verdict]string{
		Admitted:             "admit",
		RefuseConsumed:       "consumed",
		RefuseClockBehind:    "clock-behind",
		RefuseExpired:        "expired",
		RefuseDigestMismatch: "digest-mismatch",
	}
	for v, want := range names {
		if v.String() != want {
			t.Errorf("Verdict(%d).String() = %q, want %q", int(v), v.String(), want)
		}
	}
	if Verdict(0).String() != "verdict(0)" {
		t.Errorf("Verdict(0).String() = %q, want \"verdict(0)\"", Verdict(0).String())
	}
	if Admitted == Verdict(0) {
		t.Fatal("the zero Verdict is an admission")
	}
}

// capJSON is the test-side grants_json element: field order matches FactGrant.
type capJSON struct {
	Effect    string
	Scope     string
	ExpiresAt int64
	Budget    int64
}

func grantsJSON(t *testing.T, exp int64, gs ...Grant) string {
	t.Helper()
	caps := make([]capJSON, 0, len(gs))
	for _, g := range gs {
		caps = append(caps, capJSON{Effect: g.Effect, Scope: g.Scope, ExpiresAt: exp, Budget: g.Budget})
	}
	b, err := json.Marshal(caps)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestMintMatchesEveryColumn(t *testing.T) {
	r := fixtureRequest()
	sum := sha256.Sum256([]byte(strings.Repeat("aa", 32)))
	th := hex.EncodeToString(sum[:])
	clean := func() MintFacts {
		return MintFacts{
			CredentialID: th, EpisodeID: "fixture-ep", GrantsJSON: grantsJSON(t, 4600, r.Grants...),
			ExpiresAt: 4600, CreatedAt: 1000, Project: "ailang-world", ProfileDigest: strings.Repeat("ab", 32),
			TokenHash: th, RequestCreatedAt: 900, ConsumedAt: 1010,
		}
	}
	if err := MintMatches(r, clean()); err != nil {
		t.Fatalf("clean facts refused: %v", err)
	}
	boundaries := []struct {
		name string
		f    MintFacts
	}{
		{"created-at-lower-bound", func() MintFacts {
			f := clean()
			f.CreatedAt, f.ExpiresAt = 900, 4500
			f.GrantsJSON = grantsJSON(t, 4500, r.Grants...)
			return f
		}()},
		{"created-at-upper-bound", func() MintFacts {
			f := clean()
			f.CreatedAt, f.ExpiresAt = 1010, 4610
			f.GrantsJSON = grantsJSON(t, 4610, r.Grants...)
			return f
		}()},
	}
	for _, b := range boundaries {
		if err := MintMatches(r, b.f); err != nil {
			t.Errorf("%s refused: %v", b.name, err)
		}
	}

	extraGrants := append(append([]Grant{}, r.Grants...), r.Grants[1])
	arms := []struct {
		name   string
		mutate func(*MintFacts, *MintRequest)
		column string
	}{
		{"cred-id", func(f *MintFacts, _ *MintRequest) {
			sum := sha256.Sum256([]byte("other"))
			f.CredentialID = hex.EncodeToString(sum[:])
		}, "credential_id"},
		{"episode", func(f *MintFacts, _ *MintRequest) { f.EpisodeID = "other-ep" }, "episode_id"},
		{"grants-undecodable", func(f *MintFacts, _ *MintRequest) { f.GrantsJSON = "not json" }, "grants_json"},
		{"grants-unknown-key", func(f *MintFacts, _ *MintRequest) {
			f.GrantsJSON = `[{"Effect":"Workspace.Read","Scope":"fixture-ep","ExpiresAt":4600,"Budget":0,"Extra":1},{"Effect":"Workspace.Write","Scope":"fixture-ep","ExpiresAt":4600,"Budget":25}]`
		}, "grants_json"},
		{"grants-extra", func(f *MintFacts, _ *MintRequest) { f.GrantsJSON = grantsJSON(t, 4600, extraGrants...) }, "grants_json"},
		{"grants-reordered", func(f *MintFacts, _ *MintRequest) {
			f.GrantsJSON = grantsJSON(t, 4600, r.Grants[1], r.Grants[0])
		}, "grants_json"},
		{"grants-budget", func(f *MintFacts, _ *MintRequest) {
			g := append([]Grant{}, r.Grants...)
			g[1].Budget = 26
			f.GrantsJSON = grantsJSON(t, 4600, g...)
		}, "grants_json"},
		{"grants-scope", func(f *MintFacts, _ *MintRequest) {
			g := append([]Grant{}, r.Grants...)
			g[0].Scope = "other"
			f.GrantsJSON = grantsJSON(t, 4600, g...)
		}, "grants_json"},
		{"grant-exp-zero", func(f *MintFacts, _ *MintRequest) { f.GrantsJSON = grantsJSON(t, 0, r.Grants...) }, "grants_json.ExpiresAt"},
		{"grant-exp-plus1", func(f *MintFacts, _ *MintRequest) { f.GrantsJSON = grantsJSON(t, 4601, r.Grants...) }, "grants_json.ExpiresAt"},
		{"expires-off-by-one", func(f *MintFacts, _ *MintRequest) {
			f.ExpiresAt = 4601
			f.GrantsJSON = grantsJSON(t, 4601, r.Grants...)
		}, "expires_at"},
		{"created-before", func(f *MintFacts, _ *MintRequest) {
			f.CreatedAt, f.ExpiresAt = 899, 4499
			f.GrantsJSON = grantsJSON(t, 4499, r.Grants...)
		}, "created_at"},
		{"created-after", func(f *MintFacts, _ *MintRequest) {
			f.CreatedAt, f.ExpiresAt = 1011, 4611
			f.GrantsJSON = grantsJSON(t, 4611, r.Grants...)
		}, "created_at"},
		{"other-digest", func(f *MintFacts, _ *MintRequest) { f.ProfileDigest = strings.Repeat("ef", 32) }, "profile_digest"},
		{"bound-req-unbound-row", func(f *MintFacts, _ *MintRequest) { f.Project, f.ProfileDigest = "", "" }, "project"},
		{"unbound-req-bound-row", func(_ *MintFacts, req *MintRequest) { req.Profile = ProfileRef{} }, "project"},
	}
	for _, arm := range arms {
		f := clean()
		req := fixtureRequest()
		arm.mutate(&f, &req)
		err := MintMatches(req, f)
		var mm *MintMismatchError
		if !errors.As(err, &mm) {
			t.Errorf("%s: err = %v, want *MintMismatchError", arm.name, err)
			continue
		}
		if mm.Column != arm.column {
			t.Errorf("%s: Column = %q, want %q", arm.name, mm.Column, arm.column)
		}
	}
}
