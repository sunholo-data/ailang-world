package approval

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
	"testing"
)

func sampleAllowance() StandingAllowance {
	return StandingAllowance{
		Schema:       AllowanceSchema,
		PresetName:   "world-codex",
		PresetSHA256: strings.Repeat("a", 64),
		Profile:      ProfileRef{Project: "ailang-world", Digest: "sha256:" + strings.Repeat("b", 64)},
		Repos:        []string{"/Users/mark/dev/ailang-world", "/Users/mark/dev/sunholo-data"},
		Grants: []Grant{
			{Effect: "Workspace.Read", Scope: "workspace", Budget: 5000},
			{Effect: "Workspace.Write", Scope: "workspace", Budget: 5000},
		},
		TTLSeconds:         AllowanceMaxTTLSeconds,
		MaxPerDay:          AllowanceMaxPerDay,
		ValidDays:          AllowanceMaxValidDays,
		WorldDB:            "/Users/mark/.ailang/world/world.db",
		CredentialsDB:      "/Users/mark/.ailang/world/credentials.db",
		WorkspaceRoot:      "/Users/mark/.ailang/world-ws",
		ExcludeMissionLoop: true,
	}
}

// signTestSSHSIG constructs a valid armored Ed25519 SSHSIG over message under namespace.
func signTestSSHSIG(t *testing.T, priv ed25519.PrivateKey, namespace string, message []byte) ([]byte, SSHPublicKey) {
	t.Helper()
	pub := priv.Public().(ed25519.PublicKey)

	appendSSHString := func(dst, s []byte) []byte {
		var lenBuf [4]byte
		binary.BigEndian.PutUint32(lenBuf[:], uint32(len(s)))
		dst = append(dst, lenBuf[:]...)
		return append(dst, s...)
	}

	var pubBlob []byte
	pubBlob = appendSSHString(pubBlob, []byte("ssh-ed25519"))
	pubBlob = appendSSHString(pubBlob, pub)

	h := sha512.Sum512(message)
	var toSign []byte
	toSign = append(toSign, []byte("SSHSIG")...)
	toSign = appendSSHString(toSign, []byte(namespace))
	toSign = appendSSHString(toSign, nil) // reserved
	toSign = appendSSHString(toSign, []byte("sha512"))
	toSign = appendSSHString(toSign, h[:])

	rawSig := ed25519.Sign(priv, toSign)
	var sigBlob []byte
	sigBlob = appendSSHString(sigBlob, []byte("ssh-ed25519"))
	sigBlob = appendSSHString(sigBlob, rawSig)

	var blob []byte
	blob = append(blob, []byte("SSHSIG")...)
	var ver [4]byte
	binary.BigEndian.PutUint32(ver[:], 1)
	blob = append(blob, ver[:]...)
	blob = appendSSHString(blob, pubBlob)
	blob = appendSSHString(blob, []byte(namespace))
	blob = appendSSHString(blob, nil) // reserved
	blob = appendSSHString(blob, []byte("sha512"))
	blob = appendSSHString(blob, sigBlob)

	var armor bytes.Buffer
	armor.WriteString("-----BEGIN SSH SIGNATURE-----\n")
	b64 := base64.StdEncoding.EncodeToString(blob)
	for i := 0; i < len(b64); i += 70 {
		end := i + 70
		if end > len(b64) {
			end = len(b64)
		}
		armor.WriteString(b64[i:end])
		armor.WriteByte('\n')
	}
	armor.WriteString("-----END SSH SIGNATURE-----\n")
	return armor.Bytes(), SSHPublicKey{Type: "ssh-ed25519", Blob: pubBlob}
}

func TestStandingAllowanceCanonicalAndValidate(t *testing.T) {
	a := sampleAllowance()
	raw, err := CanonicalAllowance(a)
	if err != nil {
		t.Fatalf("CanonicalAllowance valid: %v", err)
	}
	got, digest, err := DecodeAllowance(raw)
	if err != nil {
		t.Fatalf("DecodeAllowance valid: %v", err)
	}
	wantDigest, err := AllowanceDigest(a)
	if err != nil || digest != wantDigest {
		t.Fatalf("digest mismatch: got %x, want %x, err %v", digest, wantDigest, err)
	}
	if got.PresetName != a.PresetName || len(got.Repos) != 2 {
		t.Fatalf("decoded mismatch: %+v", got)
	}

	// Non-canonical whitespace must be refused.
	if _, _, err := DecodeAllowance(append(bytes.Clone(raw), ' ')); !errors.Is(err, ErrNotCanonical) {
		t.Fatalf("trailing space: got %v, want ErrNotCanonical", err)
	}

	cases := []struct {
		name   string
		mutate func(*StandingAllowance)
		want   error
	}{
		{"bad-schema", func(a *StandingAllowance) { a.Schema = "world/standing-allowance/v2" }, ErrSchema},
		{"bad-preset-name", func(a *StandingAllowance) { a.PresetName = "Bad_Name" }, ErrAllowancePresetName},
		{"bad-preset-digest", func(a *StandingAllowance) { a.PresetSHA256 = "not-hex" }, ErrAllowancePresetDigest},
		{"half-profile", func(a *StandingAllowance) { a.Profile.Digest = "" }, ErrHalfProfile},
		{"empty-repos", func(a *StandingAllowance) { a.Repos = nil }, ErrAllowanceRepos},
		{"unsorted-repos", func(a *StandingAllowance) { a.Repos = []string{"/b", "/a"} }, ErrAllowanceRepos},
		{"duplicate-repos", func(a *StandingAllowance) { a.Repos = []string{"/a", "/a"} }, ErrAllowanceRepos},
		{"relative-repo", func(a *StandingAllowance) { a.Repos = []string{"rel/repo"} }, ErrAllowanceRepos},
		{"no-grants", func(a *StandingAllowance) { a.Grants = nil }, ErrNoGrants},
		{"ttl-over-12h", func(a *StandingAllowance) { a.TTLSeconds = AllowanceMaxTTLSeconds + 1 }, ErrAllowanceBoundExceeded},
		{"max-per-day-over-10", func(a *StandingAllowance) { a.MaxPerDay = AllowanceMaxPerDay + 1 }, ErrAllowanceBoundExceeded},
		{"valid-days-over-30", func(a *StandingAllowance) { a.ValidDays = AllowanceMaxValidDays + 1 }, ErrAllowanceBoundExceeded},
		{"mission-loop-not-excluded", func(a *StandingAllowance) { a.ExcludeMissionLoop = false }, ErrAllowanceMissionLoopExcluded},
		{"relative-world-db", func(a *StandingAllowance) { a.WorldDB = "world.db" }, ErrPath},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bad := sampleAllowance()
			tc.mutate(&bad)
			if err := bad.Validate(); !errors.Is(err, tc.want) {
				t.Fatalf("Validate(%s): got %v, want %v", tc.name, err, tc.want)
			}
		})
	}
}

func TestApproveAllowanceBySignatureAndDomainSeparation(t *testing.T) {
	seed := bytes.Repeat([]byte{0x42}, ed25519.SeedSize)
	priv := ed25519.NewKeyFromSeed(seed)

	a := sampleAllowance()
	raw, err := CanonicalAllowance(a)
	if err != nil {
		t.Fatalf("CanonicalAllowance: %v", err)
	}
	digest, err := AllowanceDigest(a)
	if err != nil {
		t.Fatalf("AllowanceDigest: %v", err)
	}
	var challenge [32]byte
	copy(challenge[:], bytes.Repeat([]byte{0x11}, 32))
	now := int64(1760000000)
	expiresAt := now + AllowanceMaxValidSeconds

	msg, err := AllowanceMessage(digest, challenge, expiresAt)
	if err != nil {
		t.Fatalf("AllowanceMessage: %v", err)
	}
	if len(msg) != AllowanceMessageLen {
		t.Fatalf("AllowanceMessage len = %d, want %d", len(msg), AllowanceMessageLen)
	}

	sig, pub := signTestSSHSIG(t, priv, AllowanceNamespace, msg)
	approver := Approver{Name: "mark-mac", Format: FormatSSH, Key: pub, EnrolledAt: now - 100}

	signed, err := ApproveAllowanceBySignature(AllowanceSignatureInput{
		RawPayload: raw,
		Challenge:  challenge,
		Now:        now,
		ExpiresAt:  expiresAt,
		Signature:  sig,
		Approver:   approver,
	})
	if err != nil {
		t.Fatalf("ApproveAllowanceBySignature valid: %v", err)
	}
	if signed.ApprovedBy() != "ssh:mark-mac" {
		t.Fatalf("ApprovedBy = %q, want ssh:mark-mac", signed.ApprovedBy())
	}
	if err := signed.Binds(digest, challenge, expiresAt); err != nil {
		t.Fatalf("Binds: %v", err)
	}

	// Domain separation: an allowance signature MUST NOT verify as a mint approval or verb confirmation.
	if err := VerifyApproval(sig, digest, challenge, expiresAt, pub); err == nil {
		t.Fatalf("allowance signature unexpectedly verified as mint approval")
	}
	if err := verifyVerbSSHSIG(sig, digest, challenge, expiresAt, pub); err == nil {
		t.Fatalf("allowance signature unexpectedly verified as verb confirmation")
	}

	// A signature under the mint namespace MUST NOT verify as an allowance approval.
	wrongNS, _ := signTestSSHSIG(t, priv, Namespace, msg)
	if _, err := ApproveAllowanceBySignature(AllowanceSignatureInput{
		RawPayload: raw,
		Challenge:  challenge,
		Now:        now,
		ExpiresAt:  expiresAt,
		Signature:  wrongNS,
		Approver:   approver,
	}); !errors.Is(err, ErrNamespace) {
		t.Fatalf("wrong namespace: got %v, want ErrNamespace", err)
	}

	// Lifetime over 30 days refused.
	if _, err := ApproveAllowanceBySignature(AllowanceSignatureInput{
		RawPayload: raw,
		Challenge:  challenge,
		Now:        now,
		ExpiresAt:  now + AllowanceMaxValidSeconds + 1,
		Signature:  sig,
		Approver:   approver,
	}); !errors.Is(err, ErrAllowanceBoundExceeded) {
		t.Fatalf("over-30d expiry: got %v, want ErrAllowanceBoundExceeded", err)
	}
}

func TestAdmitAllowanceDraw(t *testing.T) {
	a := sampleAllowance()
	raw, err := CanonicalAllowance(a)
	if err != nil {
		t.Fatalf("CanonicalAllowance: %v", err)
	}
	digest, err := AllowanceDigest(a)
	if err != nil {
		t.Fatalf("AllowanceDigest: %v", err)
	}
	approvedAt := int64(1760000000)
	expiresAt := approvedAt + AllowanceMaxValidSeconds

	baseInput := func() AllowanceDrawInput {
		return AllowanceDrawInput{
			RawPayload:      raw,
			StoredDigestHex: hex.EncodeToString(digest[:]),
			ApprovedAt:      approvedAt,
			ExpiresAt:       expiresAt,
			Now:             approvedAt + 3600,
			Caller:          "herdr",
			PresetName:      a.PresetName,
			PresetSHA256:    a.PresetSHA256,
			Profile:         a.Profile,
			Repo:            a.Repos[0],
			Grants:          append([]Grant(nil), a.Grants...),
			TTLSeconds:      a.TTLSeconds,
			DrawsToday:      9, // 10th draw of the day (0-indexed count 9 < 10)
		}
	}

	if err := AdmitAllowanceDraw(baseInput()); err != nil {
		t.Fatalf("10th draw should succeed: %v", err)
	}

	revoked := approvedAt + 10
	cases := []struct {
		name   string
		mutate func(*AllowanceDrawInput)
		want   error
	}{
		{"tampered-digest", func(in *AllowanceDrawInput) { in.StoredDigestHex = strings.Repeat("0", 64) }, ErrAllowanceDigestMismatch},
		{"mission-loop-excluded", func(in *AllowanceDrawInput) { in.Caller = CallerMissionLoop }, ErrAllowanceMissionLoopExcluded},
		{"revoked", func(in *AllowanceDrawInput) { in.RevokedAt = &revoked }, ErrAllowanceRevoked},
		{"approver-revoked", func(in *AllowanceDrawInput) { in.ApproverRevoked = true }, ErrAllowanceApproverRevoked},
		{"clock-behind", func(in *AllowanceDrawInput) { in.Now = approvedAt - 1 }, ErrAllowanceClockBehind},
		{"expired", func(in *AllowanceDrawInput) { in.Now = expiresAt }, ErrAllowanceExpired},
		{"preset-modified", func(in *AllowanceDrawInput) { in.PresetSHA256 = strings.Repeat("f", 64) }, ErrAllowancePresetMismatch},
		{"preset-name-mismatch", func(in *AllowanceDrawInput) { in.PresetName = "other" }, ErrAllowancePresetMismatch},
		{"profile-modified", func(in *AllowanceDrawInput) { in.Profile.Digest = "sha256:" + strings.Repeat("c", 64) }, ErrAllowanceProfileMismatch},
		{"repo-not-listed", func(in *AllowanceDrawInput) { in.Repo = "/Users/mark/dev/unlisted" }, ErrAllowanceRepoNotAllowed},
		{"grants-widened", func(in *AllowanceDrawInput) { in.Grants[0].Budget = 9999 }, ErrAllowanceGrantsMismatch},
		{"ttl-exceeded", func(in *AllowanceDrawInput) { in.TTLSeconds = AllowanceMaxTTLSeconds + 1 }, ErrAllowanceTTLExceeded},
		{"daily-cap-reached", func(in *AllowanceDrawInput) { in.DrawsToday = 10 }, ErrAllowanceDailyCapReached},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := baseInput()
			tc.mutate(&in)
			if err := AdmitAllowanceDraw(in); !errors.Is(err, tc.want) {
				t.Fatalf("AdmitAllowanceDraw(%s): got %v, want %v", tc.name, err, tc.want)
			}
		})
	}
}
