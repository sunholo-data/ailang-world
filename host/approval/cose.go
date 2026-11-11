package approval

// Row 201 A2 (design §2.4): COSE public keys, ES256 (P-256) and EdDSA (Ed25519) only, as an
// authenticator emits them inside authenticator data. What is stored is Canonical(), so one key
// has one encoding. No I/O.

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/sha256"
	"errors"
	"fmt"
)

const (
	AlgES256 int64 = -7
	AlgEdDSA int64 = -8
)

var (
	ErrCOSEKey = errors.New("approval: not an ES256 (P-256) or EdDSA (Ed25519) COSE key")
	ErrCOSEAlg = errors.New("approval: the passkey's algorithm is not ES256 (-7) or EdDSA (-8)")
)

// COSEKey is a parsed passkey public key.
type COSEKey struct {
	Alg int64  // AlgES256 or AlgEdDSA
	X   []byte // 32 bytes
	Y   []byte // 32 bytes for ES256; nil for EdDSA
}

func coseInt(n int64) cborValue { return cborValue{kind: cborInt, i: n} }

// ParseCOSEKey judges one decoded COSE key (Q1-Q8).
func ParseCOSEKey(v cborValue) (COSEKey, error) {
	if v.kind != cborMap { // Q1
		return COSEKey{}, fmt.Errorf("%w: not a map", ErrCOSEKey)
	}
	label := func(n int64) (cborValue, bool) { return v.get(coseInt(n)) }
	kty, ok1 := label(1)
	alg, ok2 := label(3)
	crv, ok3 := label(-1)
	if !ok1 || !ok2 || !ok3 || kty.kind != cborInt || alg.kind != cborInt || crv.kind != cborInt { // Q2
		return COSEKey{}, fmt.Errorf("%w: kty, alg and crv must be integers", ErrCOSEKey)
	}
	k := COSEKey{Alg: alg.i}
	allowed := map[int64]bool{1: true, 3: true, -1: true, -2: true}
	switch k.Alg {
	case AlgES256:
		if kty.i != 2 || crv.i != 1 { // Q4
			return COSEKey{}, fmt.Errorf("%w: ES256 needs kty 2 (EC2) and crv 1 (P-256)", ErrCOSEKey)
		}
		allowed[-3] = true
	case AlgEdDSA:
		if kty.i != 1 || crv.i != 6 { // Q5
			return COSEKey{}, fmt.Errorf("%w: EdDSA needs kty 1 (OKP) and crv 6 (Ed25519)", ErrCOSEKey)
		}
	default: // Q3
		return COSEKey{}, fmt.Errorf("%w (it is %d)", ErrCOSEAlg, k.Alg)
	}
	for j := 0; j < len(v.items); j += 2 {
		if key := v.items[j]; key.kind != cborInt || !allowed[key.i] { // Q8
			return COSEKey{}, fmt.Errorf("%w: an unexpected label", ErrCOSEKey)
		}
	}
	for _, c := range []struct {
		n   int64
		dst *[]byte
	}{{-2, &k.X}, {-3, &k.Y}} {
		if !allowed[c.n] {
			continue
		}
		b, ok := label(c.n)
		if !ok || b.kind != cborBytes || len(b.b) != 32 { // Q6
			return COSEKey{}, fmt.Errorf("%w: a coordinate is not 32 bytes", ErrCOSEKey)
		}
		*c.dst = append([]byte(nil), b.b...)
	}
	if k.Alg == AlgES256 {
		if _, err := ecdsa.ParseUncompressedPublicKey(elliptic.P256(), k.point()); err != nil { // Q7
			return COSEKey{}, fmt.Errorf("%w: %v", ErrCOSEKey, err)
		}
	}
	return k, nil
}

// ParseStoredCOSEKey parses a key as approvers.public_key holds it (row 201 A2 sprint 3: approve
// --webauthn-assertion rebuilds the enrolled passkey from its store row): exactly one COSE key,
// judged by Q1-Q8, with nothing after it.
func ParseStoredCOSEKey(b []byte) (COSEKey, error) {
	v, err := decodeCBORExact(b, attestationLimits)
	if err != nil {
		return COSEKey{}, fmt.Errorf("%w: %w", ErrCOSEKey, err)
	}
	return ParseCOSEKey(v)
}

func (k COSEKey) point() []byte { return append(append([]byte{4}, k.X...), k.Y...) }

// Canonical is the key's CTAP2 canonical encoding (labels 1, 3, -1, -2[, -3]): what is stored.
func (k COSEKey) Canonical() []byte {
	n, kty, crv := byte(0xa5), byte(2), byte(1)
	if k.Alg == AlgEdDSA {
		n, kty, crv = 0xa4, 1, 6
	}
	b := append([]byte{n, 0x01, kty, 0x03, byte(0x20 | (-1 - k.Alg)), 0x20, crv, 0x21, 0x58, byte(len(k.X))}, k.X...)
	if k.Alg == AlgES256 {
		b = append(append(b, 0x22, 0x58, byte(len(k.Y))), k.Y...)
	}
	return b
}

// AlgName is "ES256" or "EdDSA".
func (k COSEKey) AlgName() string {
	if k.Alg == AlgES256 {
		return "ES256"
	}
	return "EdDSA"
}

// verify checks sig over signed: ES256 is ASN.1 DER over SHA-256(signed), EdDSA is Ed25519.
func (k COSEKey) verify(signed, sig []byte) error {
	switch k.Alg {
	case AlgES256:
		pub, err := ecdsa.ParseUncompressedPublicKey(elliptic.P256(), k.point())
		h := sha256.Sum256(signed)
		if len(sig) > 72 || err != nil || !ecdsa.VerifyASN1(pub, h[:], sig) {
			return ErrBadSignature
		}
	case AlgEdDSA:
		if len(sig) != ed25519.SignatureSize || len(k.X) != ed25519.PublicKeySize || !ed25519.Verify(k.X, signed, sig) {
			return ErrBadSignature
		}
	default:
		return ErrCOSEAlg
	}
	return nil
}
