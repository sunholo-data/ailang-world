package approval

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Verdict is the outcome of the admission law. The zero Verdict is not an
// admission, so an unset verdict fails closed.
type Verdict int

const (
	Admitted Verdict = iota + 1 // the zero Verdict is not an admission
	RefuseConsumed
	RefuseClockBehind
	RefuseExpired
	RefuseDigestMismatch
)

func (v Verdict) String() string {
	switch v {
	case Admitted:
		return "admit"
	case RefuseConsumed:
		return "consumed"
	case RefuseClockBehind:
		return "clock-behind"
	case RefuseExpired:
		return "expired"
	case RefuseDigestMismatch:
		return "digest-mismatch"
	}
	return fmt.Sprintf("verdict(%d)", int(v))
}

type AdmitInput struct {
	Now, CreatedAt, ExpiresAt int64
	Consumed                  bool
	StoredDigest, BytesDigest [32]byte // the digest column; SHA-256 of the stored bytes
}

// Admit is the admission law of design §2.3, in its fixed order.
func Admit(in AdmitInput) Verdict {
	if in.Consumed {
		return RefuseConsumed
	}
	if in.Now < in.CreatedAt {
		return RefuseClockBehind
	}
	if in.Now >= in.ExpiresAt {
		return RefuseExpired
	}
	if in.StoredDigest != in.BytesDigest {
		return RefuseDigestMismatch
	}
	return Admitted
}

// FactGrant is one element of session_credentials.grants_json. Its field names
// and order match broker.Capability, which has no json tags (design V34).
type FactGrant struct {
	Effect    string
	Scope     string
	ExpiresAt int64
	Budget    int64
}

// MintFacts is the session_credentials row the mint built (its seven columns),
// plus the token's hash and the two times that bound the mint (design §2.3).
type MintFacts struct {
	CredentialID     string // credential_id
	EpisodeID        string // episode_id
	GrantsJSON       string // grants_json, raw; MintMatches decodes it
	ExpiresAt        int64  // expires_at
	CreatedAt        int64  // created_at
	Project          string // project
	ProfileDigest    string // profile_digest
	TokenHash        string // hashHexToken(the token handed to the writer)
	RequestCreatedAt int64  // the request row's created_at: lower bound for CreatedAt
	ConsumedAt       int64  // the in-transaction consume reading: upper bound for CreatedAt
}

type MintMismatchError struct{ Column, Want, Got string }

func (e *MintMismatchError) Error() string {
	return fmt.Sprintf("approval: the minted row's %s is %s, the request needs %s", e.Column, e.Got, e.Want)
}

// MintMatches checks that the row the mint built is the row the request
// asked for. It returns *MintMismatchError naming the first column that differs.
func MintMatches(r MintRequest, f MintFacts) error {
	if f.CredentialID != f.TokenHash {
		return &MintMismatchError{Column: "credential_id", Want: f.TokenHash, Got: f.CredentialID}
	}
	if f.EpisodeID != r.Episode {
		return &MintMismatchError{Column: "episode_id", Want: r.Episode, Got: f.EpisodeID}
	}
	var got []FactGrant
	dec := json.NewDecoder(strings.NewReader(f.GrantsJSON))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&got); err != nil {
		return &MintMismatchError{Column: "grants_json", Want: "the request's grants", Got: "undecodable"}
	}
	if len(got) != len(r.Grants) {
		return &MintMismatchError{Column: "grants_json", Want: fmt.Sprint(len(r.Grants)) + " grants", Got: fmt.Sprint(len(got)) + " grants"}
	}
	for i, g := range got {
		want := r.Grants[i]
		if g.Effect != want.Effect || g.Scope != want.Scope || g.Budget != want.Budget {
			return &MintMismatchError{Column: "grants_json", Want: fmt.Sprintf("grants[%d]", i), Got: fmt.Sprintf("%+v", g)}
		}
	}
	for _, g := range got {
		if g.ExpiresAt != f.ExpiresAt {
			return &MintMismatchError{Column: "grants_json.ExpiresAt", Want: fmt.Sprint(f.ExpiresAt), Got: fmt.Sprint(g.ExpiresAt)}
		}
	}
	if f.ExpiresAt != f.CreatedAt+r.TTLSeconds {
		return &MintMismatchError{Column: "expires_at", Want: fmt.Sprint(f.CreatedAt + r.TTLSeconds), Got: fmt.Sprint(f.ExpiresAt)}
	}
	if f.CreatedAt < f.RequestCreatedAt {
		return &MintMismatchError{Column: "created_at", Want: "at or after " + fmt.Sprint(f.RequestCreatedAt), Got: fmt.Sprint(f.CreatedAt)}
	}
	if f.CreatedAt > f.ConsumedAt {
		return &MintMismatchError{Column: "created_at", Want: "at or before " + fmt.Sprint(f.ConsumedAt), Got: fmt.Sprint(f.CreatedAt)}
	}
	if f.Project != r.Profile.Project {
		return &MintMismatchError{Column: "project", Want: r.Profile.Project, Got: f.Project}
	}
	if f.ProfileDigest != r.Profile.Digest {
		return &MintMismatchError{Column: "profile_digest", Want: r.Profile.Digest, Got: f.ProfileDigest}
	}
	return nil
}
