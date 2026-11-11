package approvaltest

// Row 201 A2 (design §2.15): a test authenticator. It makes registration and assertion files in
// Chrome 157's measured layout (TestTestPasskeyMatchesChromeLayout) from a fixed seed: no browser,
// no real authenticator, no randomness in the key (D-WORLD-109). Edit makes each refusal arm.

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"strconv"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// Passkey is an ES256 or EdDSA test authenticator with one 32-byte credential.
type Passkey struct {
	ed ed25519.PrivateKey
	ec *ecdsa.PrivateKey
	id []byte
}

// NewPasskey returns the authenticator whose 32-byte private key is the byte seed (1-254), repeated.
func NewPasskey(alg int64, seed byte) *Passkey {
	raw := bytes.Repeat([]byte{seed}, 32)
	id := sha256.Sum256(append([]byte("approvaltest-credential"), byte(alg), seed))
	p := &Passkey{id: id[:]}
	switch alg {
	case approval.AlgEdDSA:
		p.ed = ed25519.NewKeyFromSeed(raw)
	case approval.AlgES256:
		k, err := ecdsa.ParseRawPrivateKey(elliptic.P256(), raw)
		if err != nil {
			panic(err)
		}
		p.ec = k
	default:
		panic("approvaltest: NewPasskey needs approval.AlgES256 or approval.AlgEdDSA")
	}
	return p
}

// CredentialID is a copy of the credential ID.
func (p *Passkey) CredentialID() []byte { return bytes.Clone(p.id) }

// Key is the public key as host/approval parses it.
func (p *Passkey) Key() approval.COSEKey {
	if p.ed != nil {
		return approval.COSEKey{Alg: approval.AlgEdDSA, X: bytes.Clone(p.ed.Public().(ed25519.PublicKey))}
	}
	b, err := p.ec.PublicKey.Bytes()
	if err != nil {
		panic(err)
	}
	return approval.COSEKey{Alg: approval.AlgES256, X: b[1:33], Y: b[33:65]}
}

// Edit changes one ceremony from Chrome's layout; the zero Edit is Chrome's.
type Edit struct {
	Type         string // clientDataJSON type ("" is the ceremony's own)
	CrossOrigin  bool
	TopOrigin    string // a topOrigin member when not ""
	ClientData   []byte // replaces the whole clientDataJSON when not nil
	Set, Clear   byte   // authData flag bits set, then cleared (UP 0x01, UV 0x04, BE 0x08, BS 0x10, AT 0x40)
	Extensions   bool   // ED and the extensions map {"credProtect": 2}
	UserHandle   []byte // an assertion's userHandle member when not nil (Chrome sends none)
	RawID        []byte // replaces id and rawId when not nil
	BadSignature bool   // flips the assertion signature's last byte
}

func one(e []Edit) Edit {
	if len(e) == 0 {
		return Edit{}
	}
	return e[0]
}

// Register is a registration file for challenge at origin (fmt "none", flags UP|UV|AT, counter 1).
func (p *Passkey) Register(origin, rpID string, challenge [32]byte, e ...Edit) []byte {
	ed := one(e)
	att := append(bytes.Repeat([]byte{1, 2, 3, 4, 5, 6, 7, 8}, 2), byte(len(p.id)>>8), byte(len(p.id)))
	att = append(append(att, p.id...), p.Key().Canonical()...)
	ad := authData(rpID, 0x45, 1, ed, att)
	obj := append([]byte("\xa3\x63fmt\x64none\x67attStmt\xa0\x68authData\x58"), byte(len(ad)))
	return envelope(p, ed, map[string]string{"clientDataJSON": b64(clientData("webauthn.create", challenge, origin, ed)),
		"attestationObject": b64(append(obj, ad...))})
}

// Assert is an assertion file signing challenge at origin with counter (flags UP|UV).
func (p *Passkey) Assert(origin, rpID string, challenge [32]byte, counter uint32, e ...Edit) []byte {
	ed := one(e)
	ad := authData(rpID, 0x05, counter, ed, nil)
	cdj := clientData("webauthn.get", challenge, origin, ed)
	h := sha256.Sum256(cdj)
	sig := p.sign(append(bytes.Clone(ad), h[:]...))
	if ed.BadSignature {
		sig[len(sig)-1] ^= 1
	}
	r := map[string]string{"clientDataJSON": b64(cdj), "authenticatorData": b64(ad), "signature": b64(sig)}
	if ed.UserHandle != nil {
		r["userHandle"] = b64(ed.UserHandle)
	}
	return envelope(p, ed, r)
}

func (p *Passkey) sign(msg []byte) []byte {
	if p.ed != nil {
		return ed25519.Sign(p.ed, msg)
	}
	h := sha256.Sum256(msg)
	sig, err := ecdsa.SignASN1(rand.Reader, p.ec, h[:])
	if err != nil {
		panic(err)
	}
	return sig
}

func authData(rpID string, flags byte, counter uint32, ed Edit, attested []byte) []byte {
	f := (flags | ed.Set) &^ ed.Clear
	if ed.Extensions {
		f |= 0x80
	}
	h := sha256.Sum256([]byte(rpID))
	b := append(binary.BigEndian.AppendUint32(append(h[:], f), counter), attested...)
	if ed.Extensions {
		b = append(b, "\xa1\x6bcredProtect\x02"...)
	}
	return b
}

func clientData(typ string, challenge [32]byte, origin string, ed Edit) []byte {
	if ed.ClientData != nil {
		return ed.ClientData
	}
	if ed.Type != "" {
		typ = ed.Type
	}
	q := func(s string) string { b, _ := json.Marshal(s); return string(b) }
	s := `{"type":` + q(typ) + `,"challenge":` + q(b64(challenge[:])) + `,"origin":` + q(origin) + `,"crossOrigin":` + strconv.FormatBool(ed.CrossOrigin)
	if ed.TopOrigin != "" {
		s += `,"topOrigin":` + q(ed.TopOrigin)
	}
	return []byte(s + "}")
}

func b64(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }

func envelope(p *Passkey, ed Edit, response map[string]string) []byte {
	id := p.id
	if ed.RawID != nil {
		id = ed.RawID
	}
	b, err := json.Marshal(map[string]any{"id": b64(id), "rawId": b64(id), "type": "public-key", "response": response})
	if err != nil {
		panic(err)
	}
	return b
}
