package approval

// Row 201 A2 (design §2.5): authenticator data, laid out as WebAuthn L3 §6.1 says: rpIdHash (32),
// flags (1), signCount (4, big-endian), then for a registration the attested credential data,
// then any extensions map. No I/O.

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// MaxCredentialIDBytes is L3 §7.1 step 25's bound; there is no lower bound but 1 (design §2.5 [r1]).
const MaxCredentialIDBytes = 1023

var (
	ErrAuthDataShort      = errors.New("approval: authenticator data is truncated")
	ErrAuthDataFlags      = errors.New("approval: authenticator data has inconsistent flags")
	ErrCredentialID       = errors.New("approval: a credential ID must be 1 to 1023 bytes")
	ErrAuthDataExtensions = errors.New("approval: authenticator data extensions are not one CBOR map")
	ErrAuthDataTrailing   = errors.New("approval: authenticator data has trailing bytes")
)

// AuthData is parsed authenticator data. The last three fields are a registration's only.
type AuthData struct {
	RPIDHash       [32]byte
	UP, UV, BE, BS bool // flags bits 0, 2, 3, 4
	SignCount      uint32
	CredentialID   []byte
	AAGUID         [16]byte
	Key            COSEKey
}

// ParseAuthData judges authenticator data (D1-D8); registration says whether attested
// credential data must be present (a registration) or absent (an assertion).
func ParseAuthData(b []byte, registration bool) (AuthData, error) {
	if len(b) < 37 { // D1
		return AuthData{}, ErrAuthDataShort
	}
	var a AuthData
	copy(a.RPIDHash[:], b[:32])
	f := b[32]
	a.UP, a.UV, a.BE, a.BS = f&0x01 != 0, f&0x04 != 0, f&0x08 != 0, f&0x10 != 0
	if a.BS && !a.BE { // D2
		return AuthData{}, fmt.Errorf("%w: backup state without backup eligibility", ErrAuthDataFlags)
	}
	if (f&0x40 != 0) != registration { // D3
		return AuthData{}, fmt.Errorf("%w: attested credential data is for a registration only", ErrAuthDataFlags)
	}
	a.SignCount = binary.BigEndian.Uint32(b[33:37])
	rest := b[37:]
	if registration {
		if len(rest) < 18 { // D4
			return AuthData{}, ErrAuthDataShort
		}
		copy(a.AAGUID[:], rest[:16])
		n := int(binary.BigEndian.Uint16(rest[16:18]))
		if n < 1 || n > MaxCredentialIDBytes { // D5
			return AuthData{}, ErrCredentialID
		}
		if len(rest) < 18+n { // D4
			return AuthData{}, ErrAuthDataShort
		}
		a.CredentialID = append([]byte(nil), rest[18:18+n]...)
		v, r, err := decodeCBOR(rest[18+n:], attestationLimits)
		if err == nil {
			a.Key, err = ParseCOSEKey(v)
		}
		if err != nil { // D6
			return AuthData{}, fmt.Errorf("approval: the passkey's public key: %w", err)
		}
		rest = r
	}
	if f&0x80 != 0 {
		v, r, err := decodeCBOR(rest, attestationLimits)
		if err != nil || v.kind != cborMap { // D7
			return AuthData{}, ErrAuthDataExtensions
		}
		rest = r
	}
	if len(rest) != 0 { // D8
		return AuthData{}, ErrAuthDataTrailing
	}
	return a, nil
}
