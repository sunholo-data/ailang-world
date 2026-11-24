package approval

// Row 224 (design_docs/planned/w-standing-session-allowances.md §2.2-§2.3, D-WORLD-96):
// the pure domain, canonical bytes, domain-separated signed message, SSHSIG verification,
// and draw-admission law for standing per-preset session allowances.
// Pure: standard library only, no I/O, no clock, no randomness (coding standard S2).
// No production caller yet (Slice 2 stores it, Slice 3 wires CLI verbs).

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	// AllowanceSchema is the closed JSON schema of a stored standing allowance.
	AllowanceSchema = "world/standing-allowance/v1"

	// AllowancePrefix starts every signed standing-allowance message (design §2.2).
	AllowancePrefix = "world-allowance/v1"

	// AllowanceNamespace is the SSHSIG namespace every standing allowance is signed under,
	// domain-separated from mint ("world-approval@sunholo"), verb ("world-verb@sunholo"),
	// and operator action ("world-operator-action@sunholo").
	AllowanceNamespace = "world-allowance@sunholo"

	// AllowanceMessageLen is len(AllowancePrefix) + 1 ('\n') + 32 (digest) + 32 (challenge) + 8 (expiry).
	AllowanceMessageLen = len(AllowancePrefix) + 1 + 32 + 32 + 8

	// D-WORLD-96 hard ceilings.
	AllowanceMaxPerDay       = 10
	AllowanceMaxTTLSeconds   = 12 * 3600 // 43,200 s (12 h)
	AllowanceMaxValidDays    = 30
	AllowanceMaxValidSeconds = AllowanceMaxValidDays * 86400 // 2,592,000 s (30 d)

	// AllowanceMaxBytes bounds a standing allowance's canonical JSON bytes.
	AllowanceMaxBytes = 128 << 10

	// CallerMissionLoop identifies the excluded row 196 mission-loop caller class.
	CallerMissionLoop = "mission-loop"
)

const maxAllowancePresetNameLen = 51

// StandingAllowance is the canonical payload approved once by the operator for a launch preset.
// Every field that governs session authority is bound inside its canonical JSON digest.
type StandingAllowance struct {
	Schema             string     `json:"schema"`               // AllowanceSchema
	PresetName         string     `json:"preset_name"`          // 1..51 chars, episode grammar
	PresetSHA256       string     `json:"preset_sha256"`        // 64 lowercase hex of exact preset bytes
	Profile            ProfileRef `json:"profile"`              // both "" if unbound, both set if bound
	Repos              []string   `json:"repos"`                // non-empty, sorted, unique, absolute clean paths
	Grants             []Grant    `json:"grants"`               // non-empty preset grants in order
	TTLSeconds         int64      `json:"ttl_s"`                // 1..AllowanceMaxTTLSeconds (12 h)
	MaxPerDay          int        `json:"max_per_day"`          // 1..AllowanceMaxPerDay (10)
	ValidDays          int        `json:"valid_days"`           // 1..AllowanceMaxValidDays (30)
	WorldDB            string     `json:"world_db"`             // canonical absolute path
	CredentialsDB      string     `json:"credentials_db"`       // canonical absolute path
	WorkspaceRoot      string     `json:"workspace_root"`       // canonical absolute path
	ExcludeMissionLoop bool       `json:"exclude_mission_loop"` // must be true (D-WORLD-96)
}

var (
	ErrAllowancePresetName          = errors.New("approval: allowance: bad preset_name")
	ErrAllowancePresetDigest        = errors.New("approval: allowance: preset_sha256 must be 64 lowercase hex")
	ErrAllowanceRepos               = errors.New("approval: allowance: repos must be non-empty, sorted, unique, absolute clean paths")
	ErrAllowanceBoundExceeded       = errors.New("approval: allowance: bound exceeds D-WORLD-96 ceiling")
	ErrAllowanceMissionLoopExcluded = errors.New("approval: allowance: mission loop (row 196) is excluded by D-WORLD-96")
	ErrAllowanceTooLarge            = errors.New("approval: allowance: payload exceeds maximum bytes")

	// Draw-admission refusals (D1-D8).
	ErrAllowanceDigestMismatch  = errors.New("approval: allowance: stored bytes do not match stored digest")
	ErrAllowanceRevoked         = errors.New("approval: allowance: standing allowance has been revoked")
	ErrAllowanceApproverRevoked = errors.New("approval: allowance: signing approver key has been revoked")
	ErrAllowanceClockBehind     = errors.New("approval: allowance: clock is behind approved_at")
	ErrAllowanceExpired         = errors.New("approval: allowance: standing allowance has expired")
	ErrAllowancePresetMismatch  = errors.New("approval: allowance: preset name or file digest does not match allowance")
	ErrAllowanceProfileMismatch = errors.New("approval: allowance: exec profile binding does not match allowance")
	ErrAllowanceRepoNotAllowed  = errors.New("approval: allowance: repo is not listed in allowance")
	ErrAllowanceGrantsMismatch  = errors.New("approval: allowance: requested grants do not match allowance")
	ErrAllowanceTTLExceeded     = errors.New("approval: allowance: requested ttl exceeds allowance bound")
	ErrAllowanceDailyCapReached = errors.New("approval: allowance: daily session cap reached")
)

func isLowerHex64(s string) bool {
	if len(s) != 64 {
		return false
	}
	for i := 0; i < 64; i++ {
		c := s[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}

// Validate checks every structural and D-WORLD-96 bound on a StandingAllowance.
func (a StandingAllowance) Validate() error {
	if a.Schema != AllowanceSchema {
		return fmt.Errorf("%w: schema is %q, want %q", ErrSchema, a.Schema, AllowanceSchema)
	}
	if err := stringLeaves(reflect.ValueOf(a), "", func(path, s string) error {
		if !utf8.ValidString(s) || strings.ContainsRune(s, 0) {
			return fmt.Errorf("%w: %s", ErrString, path)
		}
		return nil
	}); err != nil {
		return err
	}
	if len(a.PresetName) == 0 || len(a.PresetName) > maxAllowancePresetNameLen || !episodePattern.MatchString(a.PresetName) {
		return fmt.Errorf("%w: %q", ErrAllowancePresetName, a.PresetName)
	}
	if !isLowerHex64(a.PresetSHA256) {
		return fmt.Errorf("%w: %q", ErrAllowancePresetDigest, a.PresetSHA256)
	}
	if (a.Profile.Project == "") != (a.Profile.Digest == "") {
		return fmt.Errorf("%w: profile needs both project and digest", ErrHalfProfile)
	}
	if len(a.Repos) == 0 {
		return fmt.Errorf("%w: empty", ErrAllowanceRepos)
	}
	for i, r := range a.Repos {
		if !absClean(r) {
			return fmt.Errorf("%w: repos[%d] %q", ErrAllowanceRepos, i, r)
		}
		if i > 0 && a.Repos[i-1] >= r {
			return fmt.Errorf("%w: repos must be strictly sorted and unique at index %d", ErrAllowanceRepos, i)
		}
	}
	if len(a.Grants) == 0 {
		return fmt.Errorf("%w: grants", ErrNoGrants)
	}
	for i, g := range a.Grants {
		if g.Effect == "" || g.Scope == "" {
			return fmt.Errorf("%w: grants[%d]", ErrGrantEmpty, i)
		}
		if g.Budget < 0 {
			return fmt.Errorf("%w: grants[%d].budget is %d", ErrNegative, i, g.Budget)
		}
	}
	if a.TTLSeconds <= 0 || a.TTLSeconds > AllowanceMaxTTLSeconds {
		return fmt.Errorf("%w: ttl_s %d not in 1..%d", ErrAllowanceBoundExceeded, a.TTLSeconds, AllowanceMaxTTLSeconds)
	}
	if a.MaxPerDay <= 0 || a.MaxPerDay > AllowanceMaxPerDay {
		return fmt.Errorf("%w: max_per_day %d not in 1..%d", ErrAllowanceBoundExceeded, a.MaxPerDay, AllowanceMaxPerDay)
	}
	if a.ValidDays <= 0 || a.ValidDays > AllowanceMaxValidDays {
		return fmt.Errorf("%w: valid_days %d not in 1..%d", ErrAllowanceBoundExceeded, a.ValidDays, AllowanceMaxValidDays)
	}
	for _, p := range []struct{ name, val string }{
		{"world_db", a.WorldDB},
		{"credentials_db", a.CredentialsDB},
		{"workspace_root", a.WorkspaceRoot},
	} {
		if !absClean(p.val) {
			return fmt.Errorf("%w: %s %q", ErrPath, p.name, p.val)
		}
	}
	if !a.ExcludeMissionLoop {
		return ErrAllowanceMissionLoopExcluded
	}
	return nil
}

// CanonicalAllowance encodes a as compact JSON in struct field order with HTML escaping disabled
// and no trailing newline. It validates a first.
func CanonicalAllowance(a StandingAllowance) ([]byte, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(a); err != nil {
		return nil, err
	}
	b := bytes.TrimSuffix(buf.Bytes(), []byte("\n"))
	if len(b) > AllowanceMaxBytes {
		return nil, fmt.Errorf("%w: %d > %d", ErrAllowanceTooLarge, len(b), AllowanceMaxBytes)
	}
	return b, nil
}

// AllowanceDigest returns the SHA-256 of CanonicalAllowance(a).
func AllowanceDigest(a StandingAllowance) ([32]byte, error) {
	b, err := CanonicalAllowance(a)
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(b), nil
}

// DecodeAllowance strictly decodes canonical standing-allowance bytes, verifies canonical form
// and D-WORLD-96 bounds, and returns the decoded struct and its SHA-256 digest.
func DecodeAllowance(raw []byte) (StandingAllowance, [32]byte, error) {
	if len(raw) == 0 || len(raw) > AllowanceMaxBytes {
		return StandingAllowance{}, [32]byte{}, fmt.Errorf("%w: length %d", ErrDecode, len(raw))
	}
	var a StandingAllowance
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&a); err != nil {
		return StandingAllowance{}, [32]byte{}, fmt.Errorf("%w: %v", ErrDecode, err)
	}
	if dec.More() {
		return StandingAllowance{}, [32]byte{}, fmt.Errorf("%w: trailing tokens", ErrDecode)
	}
	canon, err := CanonicalAllowance(a)
	if err != nil {
		return StandingAllowance{}, [32]byte{}, err
	}
	if !bytes.Equal(raw, canon) {
		return StandingAllowance{}, [32]byte{}, ErrNotCanonical
	}
	return a, sha256.Sum256(raw), nil
}

// AllowanceMessage builds the 91-byte domain-separated message signed for a standing allowance:
// "world-allowance/v1\n" || digest (32) || challenge (32) || big-endian uint64(expiresAt) (8).
func AllowanceMessage(digest, challenge [32]byte, expiresAt int64) ([]byte, error) {
	if expiresAt < 0 {
		return nil, fmt.Errorf("%w: %d", ErrNegativeExpiry, expiresAt)
	}
	msg := make([]byte, AllowanceMessageLen)
	n := copy(msg, AllowancePrefix)
	msg[n] = '\n'
	n++
	n += copy(msg[n:], digest[:])
	n += copy(msg[n:], challenge[:])
	binary.BigEndian.PutUint64(msg[n:], uint64(expiresAt))
	return msg, nil
}

// VerifyAllowanceSSHSIG parses an armored SSHSIG and verifies it over AllowanceMessage(digest,
// challenge, expiresAt) under AllowanceNamespace with key.
func VerifyAllowanceSSHSIG(armored []byte, digest, challenge [32]byte, expiresAt int64, key SSHPublicKey) error {
	s, err := ParseSSHSIG(armored)
	if err != nil {
		return err
	}
	msg, err := AllowanceMessage(digest, challenge, expiresAt)
	if err != nil {
		return err
	}
	return VerifySSHSIG(s, AllowanceNamespace, msg, key)
}

// SignedAllowance is a verified cryptographic approval of a standing allowance. Its fields are
// unexported so only ApproveAllowanceBySignature can construct a verified instance.
type SignedAllowance struct {
	ok        bool
	digest    [32]byte
	challenge [32]byte
	expiresAt int64
	name      string
	format    string
	key       []byte
}

// ApprovedBy returns "<format>:<name>" for a verified SignedAllowance, or "" otherwise.
func (s *SignedAllowance) ApprovedBy() string {
	if s == nil || !s.ok {
		return ""
	}
	return s.format + ":" + s.name
}

// Binds reports whether s is verified and matches (digest, challenge, expiresAt).
func (s *SignedAllowance) Binds(digest, challenge [32]byte, expiresAt int64) error {
	if s == nil || !s.ok {
		return ErrNotVerified
	}
	if s.digest != digest {
		return &BindingError{Field: "digest"}
	}
	if s.challenge != challenge {
		return &BindingError{Field: "challenge"}
	}
	if s.expiresAt != expiresAt {
		return &BindingError{Field: "expiry"}
	}
	return nil
}

// AllowanceSignatureInput is the input to ApproveAllowanceBySignature.
type AllowanceSignatureInput struct {
	RawPayload []byte
	Challenge  [32]byte
	Now        int64
	ExpiresAt  int64
	Signature  []byte
	Approver   Approver
}

// ApproveAllowanceBySignature verifies an SSHSIG approval for a standing allowance:
//  1. Decode & validate canonical payload and D-WORLD-96 bounds (including expiry <= valid_days <= 30d).
//  2. Verify approver state (name, format, not revoked, enrolled <= Now).
//  3. Verify SSHSIG under AllowanceNamespace over AllowanceMessage(digest, challenge, expiresAt).
func ApproveAllowanceBySignature(in AllowanceSignatureInput) (*SignedAllowance, error) {
	a, digest, err := DecodeAllowance(in.RawPayload)
	if err != nil {
		return nil, err
	}
	if in.ExpiresAt <= in.Now {
		return nil, ErrAllowanceExpired
	}
	maxLifetime := int64(a.ValidDays) * 86400
	if in.ExpiresAt-in.Now > maxLifetime || in.ExpiresAt-in.Now > AllowanceMaxValidSeconds {
		return nil, fmt.Errorf("%w: lifetime %ds exceeds %ds", ErrAllowanceBoundExceeded, in.ExpiresAt-in.Now, maxLifetime)
	}
	admit := AdmitInput{Now: in.Now, CreatedAt: in.Now, ExpiresAt: in.ExpiresAt, StoredDigest: digest, BytesDigest: digest}
	if err := checkApprover(admit, in.Approver, FormatSSH); err != nil {
		return nil, err
	}
	if err := VerifyAllowanceSSHSIG(in.Signature, digest, in.Challenge, in.ExpiresAt, in.Approver.Key); err != nil {
		return nil, err
	}
	return &SignedAllowance{
		ok:        true,
		digest:    digest,
		challenge: in.Challenge,
		expiresAt: in.ExpiresAt,
		name:      in.Approver.Name,
		format:    in.Approver.Format,
		key:       bytes.Clone(in.Approver.Key.Blob),
	}, nil
}

// AllowanceDayKey returns the UTC calendar day key ("YYYY-MM-DD") for unixSec (design R3).
func AllowanceDayKey(unixSec int64) string {
	return time.Unix(unixSec, 0).UTC().Format("2006-01-02")
}

// AllowanceDrawInput carries every fact needed to judge whether a session may be drawn
// from a standing allowance at Now without prompting the operator.
type AllowanceDrawInput struct {
	// Stored allowance record facts.
	RawPayload      []byte
	StoredDigestHex string
	ApprovedAt      int64
	ExpiresAt       int64
	RevokedAt       *int64
	ApproverRevoked bool // true if signing approver key has since been revoked (design R6)

	// Draw request facts.
	Now          int64
	Caller       string // "" or host adapter name; "mission-loop" is always refused
	PresetName   string
	PresetSHA256 string
	Profile      ProfileRef
	Repo         string
	Grants       []Grant
	TTLSeconds   int64
	DrawsToday   int // count of existing draws for this allowance on AllowanceDayKey(Now)
}

// AdmitAllowanceDraw is the pure draw-admission law (D1-D8, in order). It returns nil iff
// the draw is strictly within the standing allowance and every D-WORLD-96 bound.
func AdmitAllowanceDraw(in AllowanceDrawInput) error {
	// D1: Integrity & canonical decode.
	a, computedDigest, err := DecodeAllowance(in.RawPayload)
	if err != nil {
		return err
	}
	if hex.EncodeToString(computedDigest[:]) != in.StoredDigestHex {
		return ErrAllowanceDigestMismatch
	}

	// D2: Mission loop exclusion (D-WORLD-96).
	if in.Caller == CallerMissionLoop || !a.ExcludeMissionLoop {
		return ErrAllowanceMissionLoopExcluded
	}

	// D3: Revocation (allowance or its signing approver key).
	if in.RevokedAt != nil {
		return ErrAllowanceRevoked
	}
	if in.ApproverRevoked {
		return ErrAllowanceApproverRevoked
	}

	// D4: Clock & 30-day expiry ceiling.
	if in.Now < in.ApprovedAt {
		return ErrAllowanceClockBehind
	}
	if in.Now >= in.ExpiresAt {
		return ErrAllowanceExpired
	}
	if in.ExpiresAt-in.ApprovedAt > AllowanceMaxValidSeconds || in.ExpiresAt-in.ApprovedAt > int64(a.ValidDays)*86400 {
		return fmt.Errorf("%w: stored validity window %ds exceeds ceiling", ErrAllowanceBoundExceeded, in.ExpiresAt-in.ApprovedAt)
	}

	// D5: Preset name, exact preset bytes digest, and exec profile digest binding.
	if in.PresetName != a.PresetName || in.PresetSHA256 != a.PresetSHA256 {
		return ErrAllowancePresetMismatch
	}
	if in.Profile != a.Profile {
		return ErrAllowanceProfileMismatch
	}

	// D6: Repo listed in preset / allowance.
	if !slices.Contains(a.Repos, in.Repo) {
		return fmt.Errorf("%w: %q", ErrAllowanceRepoNotAllowed, in.Repo)
	}

	// D7: Exact preset grants and TTL within 12 h / allowance bound.
	if !slices.Equal(in.Grants, a.Grants) {
		return ErrAllowanceGrantsMismatch
	}
	if in.TTLSeconds <= 0 || in.TTLSeconds > a.TTLSeconds || in.TTLSeconds > AllowanceMaxTTLSeconds {
		return fmt.Errorf("%w: requested %ds, max %ds", ErrAllowanceTTLExceeded, in.TTLSeconds, a.TTLSeconds)
	}

	// D8: Daily cap (at most 10 sessions/day).
	if in.DrawsToday < 0 {
		return fmt.Errorf("%w: draws_today %d", ErrNegative, in.DrawsToday)
	}
	if in.DrawsToday >= a.MaxPerDay || in.DrawsToday >= AllowanceMaxPerDay {
		return fmt.Errorf("%w: %d/%d used on %s", ErrAllowanceDailyCapReached, in.DrawsToday, a.MaxPerDay, AllowanceDayKey(in.Now))
	}
	return nil
}
