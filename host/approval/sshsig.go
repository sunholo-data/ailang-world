package approval

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// Namespace is the SSHSIG namespace every world approval is signed under.
const Namespace = "world-approval@sunholo"

const (
	armorBegin = "-----BEGIN SSH SIGNATURE-----"
	armorEnd   = "-----END SSH SIGNATURE-----"
	keyEd25519 = "ssh-ed25519"
	keyP256    = "ecdsa-sha2-nistp256"
)

var (
	ErrArmor         = errors.New("approval: sshsig: bad armor")
	ErrBase64        = errors.New("approval: sshsig: bad base64")
	ErrMagic         = errors.New("approval: sshsig: bad magic")
	ErrVersion       = errors.New("approval: sshsig: unsupported version")
	ErrTruncated     = errors.New("approval: sshsig: truncated")
	ErrTrailing      = errors.New("approval: sshsig: trailing bytes")
	ErrReserved      = errors.New("approval: sshsig: reserved field is not empty")
	ErrHashAlg       = errors.New("approval: sshsig: hash algorithm is not sha256 or sha512")
	ErrKeyType       = errors.New("approval: sshsig: key type refused")
	ErrNamespace     = errors.New("approval: sshsig: wrong namespace")
	ErrKeyMismatch   = errors.New("approval: sshsig: signature key is not the expected key")
	ErrSigType       = errors.New("approval: sshsig: signature type does not match the key")
	ErrBadSignature  = errors.New("approval: sshsig: signature does not verify")
	ErrAuthorizedKey = errors.New("approval: sshsig: malformed authorized key line")
)

// SSHPublicKey is an SSH public key: its key type and its SSH wire-encoded blob.
type SSHPublicKey struct {
	Type string
	Blob []byte
}

// SSHSig is a parsed SSHSIG signature.
type SSHSig struct {
	PublicKey          SSHPublicKey
	Namespace, HashAlg string
	SigType            string
	Sig                []byte
}

// sshReader reads SSH wire-format fields.
type sshReader struct{ b []byte }

func (r *sshReader) uint32() (uint32, error) {
	if len(r.b) < 4 {
		return 0, ErrTruncated
	}
	v := binary.BigEndian.Uint32(r.b[:4])
	r.b = r.b[4:]
	return v, nil
}

func (r *sshReader) readString() ([]byte, error) {
	n, err := r.uint32()
	if err != nil {
		return nil, err
	}
	if uint64(n) > uint64(len(r.b)) {
		return nil, ErrTruncated
	}
	s := r.b[:n:n]
	r.b = r.b[n:]
	return s, nil
}

func (r *sshReader) end() error {
	if len(r.b) != 0 {
		return ErrTrailing
	}
	return nil
}

// appendString appends an SSH string: uint32 BE length followed by s.
func appendString(dst, s []byte) []byte {
	dst = binary.BigEndian.AppendUint32(dst, uint32(len(s)))
	return append(dst, s...)
}

type keyMaterial struct {
	typ string
	ed  ed25519.PublicKey
	ec  *ecdsa.PublicKey
}

func parseKeyBlob(blob []byte) (keyMaterial, error) {
	r := sshReader{blob}
	tb, err := r.readString()
	if err != nil {
		return keyMaterial{}, err
	}
	t := string(tb)
	switch t {
	case keyEd25519:
		pk, err := r.readString()
		if err != nil {
			return keyMaterial{}, err
		}
		if err := r.end(); err != nil {
			return keyMaterial{}, err
		}
		if len(pk) != ed25519.PublicKeySize {
			return keyMaterial{}, fmt.Errorf("%w: ssh-ed25519 key is %d bytes, want 32", ErrKeyType, len(pk))
		}
		return keyMaterial{typ: t, ed: ed25519.PublicKey(bytes.Clone(pk))}, nil
	case keyP256:
		curve, err := r.readString()
		if err != nil {
			return keyMaterial{}, err
		}
		q, err := r.readString()
		if err != nil {
			return keyMaterial{}, err
		}
		if err := r.end(); err != nil {
			return keyMaterial{}, err
		}
		if string(curve) != "nistp256" {
			return keyMaterial{}, fmt.Errorf("%w: ecdsa-sha2-nistp256 curve is %q, want nistp256", ErrKeyType, curve)
		}
		pub, err := ecdsa.ParseUncompressedPublicKey(elliptic.P256(), q)
		if err != nil {
			return keyMaterial{}, fmt.Errorf("%w: ecdsa-sha2-nistp256 point: %v", ErrKeyType, err)
		}
		return keyMaterial{typ: t, ec: pub}, nil
	case "sk-ssh-ed25519@openssh.com", "sk-ecdsa-sha2-nistp256@openssh.com":
		return keyMaterial{}, fmt.Errorf("%w: %s security-key signatures are refused until an attended hardware fixture exists (row 201 Q1)", ErrKeyType, t)
	default:
		return keyMaterial{}, fmt.Errorf("%w: %q is not ssh-ed25519 or ecdsa-sha2-nistp256", ErrKeyType, t)
	}
}

// Fingerprint returns the OpenSSH-style SHA256 fingerprint of the key blob.
func (k SSHPublicKey) Fingerprint() string {
	sum := sha256.Sum256(k.Blob)
	return "SHA256:" + base64.RawStdEncoding.EncodeToString(sum[:])
}

// ParseAuthorizedKey parses one authorized_keys-style line: "<type> <base64> [comment]".
func ParseAuthorizedKey(line []byte) (SSHPublicKey, error) {
	s := strings.TrimSuffix(string(line), "\n")
	if s == "" || strings.ContainsAny(s, "\r\n") {
		return SSHPublicKey{}, ErrAuthorizedKey
	}
	f := strings.Fields(s)
	if len(f) < 2 {
		return SSHPublicKey{}, ErrAuthorizedKey
	}
	blob, err := base64.StdEncoding.Strict().DecodeString(f[1])
	if err != nil {
		return SSHPublicKey{}, fmt.Errorf("%w: %v", ErrBase64, err)
	}
	km, err := parseKeyBlob(blob)
	if err != nil {
		return SSHPublicKey{}, err
	}
	if f[0] != km.typ {
		return SSHPublicKey{}, fmt.Errorf("%w: the line says %q, the key is %q", ErrKeyType, f[0], km.typ)
	}
	return SSHPublicKey{Type: km.typ, Blob: blob}, nil
}

// ParseSSHSIG parses an armored SSHSIG signature. It never returns a non-nil *SSHSig together
// with an error.
func ParseSSHSIG(armored []byte) (*SSHSig, error) {
	if len(armored) == 0 {
		return nil, ErrArmor
	}
	s := strings.TrimSuffix(string(armored), "\n")
	lines := strings.Split(s, "\n")
	if len(lines) < 3 || lines[0] != armorBegin || lines[len(lines)-1] != armorEnd {
		return nil, ErrArmor
	}
	body := lines[1 : len(lines)-1]
	for _, l := range body {
		if l == "" {
			return nil, ErrArmor
		}
	}
	blob, err := base64.StdEncoding.Strict().DecodeString(strings.Join(body, ""))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBase64, err)
	}
	return parseSSHSIGBlob(blob)
}

func parseSSHSIGBlob(b []byte) (*SSHSig, error) {
	if len(b) < 6 {
		return nil, ErrTruncated
	}
	if !bytes.Equal(b[:6], []byte("SSHSIG")) {
		return nil, ErrMagic
	}
	r := sshReader{b[6:]}
	v, err := r.uint32()
	if err != nil {
		return nil, err
	}
	if v != 1 {
		return nil, fmt.Errorf("%w: %d", ErrVersion, v)
	}
	var fields [5][]byte
	for i := range fields {
		if fields[i], err = r.readString(); err != nil {
			return nil, err
		}
	}
	if err := r.end(); err != nil {
		return nil, err
	}
	pk, ns, reserved, hashAlg, sig := fields[0], fields[1], fields[2], fields[3], fields[4]
	if len(reserved) != 0 {
		return nil, ErrReserved
	}
	if h := string(hashAlg); h != "sha256" && h != "sha512" {
		return nil, ErrHashAlg
	}
	km, err := parseKeyBlob(pk)
	if err != nil {
		return nil, err
	}
	sr := sshReader{sig}
	sigType, err := sr.readString()
	if err != nil {
		return nil, err
	}
	sigBlob, err := sr.readString()
	if err != nil {
		return nil, err
	}
	if err := sr.end(); err != nil {
		return nil, err
	}
	return &SSHSig{
		PublicKey: SSHPublicKey{Type: km.typ, Blob: bytes.Clone(pk)},
		Namespace: string(ns),
		HashAlg:   string(hashAlg),
		SigType:   string(sigType),
		Sig:       bytes.Clone(sigBlob),
	}, nil
}

// VerifySSHSIG checks s over message under namespace with the expected key. The steps run in a
// fixed order (design §2.7): namespace, then key identity, then the key type, then the signature.
func VerifySSHSIG(s *SSHSig, namespace string, message []byte, key SSHPublicKey) error {
	if s == nil {
		return ErrBadSignature
	}
	if s.Namespace != namespace {
		return ErrNamespace
	}
	if !bytes.Equal(s.PublicKey.Blob, key.Blob) {
		return ErrKeyMismatch
	}
	km, err := parseKeyBlob(key.Blob)
	if err != nil {
		return err
	}
	if s.SigType != km.typ {
		return ErrSigType
	}
	sd, err := signedData(namespace, s.HashAlg, message)
	if err != nil {
		return err
	}
	switch km.typ {
	case keyEd25519:
		if len(s.Sig) != ed25519.SignatureSize || !ed25519.Verify(km.ed, sd, s.Sig) {
			return ErrBadSignature
		}
	case keyP256:
		sr := sshReader{s.Sig}
		rb, err := sr.readString()
		if err != nil {
			return fmt.Errorf("%w: %v", ErrBadSignature, err)
		}
		sb, err := sr.readString()
		if err != nil {
			return fmt.Errorf("%w: %v", ErrBadSignature, err)
		}
		if err := sr.end(); err != nil {
			return fmt.Errorf("%w: %v", ErrBadSignature, err)
		}
		rv, err := parseMpint(rb)
		if err != nil {
			return err
		}
		sv, err := parseMpint(sb)
		if err != nil {
			return err
		}
		h := sha256.Sum256(sd)
		if !ecdsa.Verify(km.ec, h[:], rv, sv) {
			return ErrBadSignature
		}
	default:
		return ErrKeyType
	}
	return nil
}

// signedData is the byte string an SSHSIG signature covers.
func signedData(namespace, hashAlg string, message []byte) ([]byte, error) {
	var h []byte
	switch hashAlg {
	case "sha512":
		d := sha512.Sum512(message)
		h = d[:]
	case "sha256":
		d := sha256.Sum256(message)
		h = d[:]
	default:
		return nil, ErrHashAlg
	}
	out := []byte("SSHSIG")
	out = appendString(out, []byte(namespace))
	out = appendString(out, nil)
	out = appendString(out, []byte(hashAlg))
	out = appendString(out, h)
	return out, nil
}

// parseMpint parses an SSH mpint, refusing negative and non-minimal encodings.
func parseMpint(b []byte) (*big.Int, error) {
	if len(b) > 0 && b[0]&0x80 != 0 {
		return nil, fmt.Errorf("%w: negative mpint", ErrBadSignature)
	}
	if len(b) > 0 && b[0] == 0 && (len(b) == 1 || b[1]&0x80 == 0) {
		return nil, fmt.Errorf("%w: non-minimal mpint", ErrBadSignature)
	}
	return new(big.Int).SetBytes(b), nil
}

// VerifyApproval parses an armored SSHSIG over the approval message for (digest, challenge,
// expiry) and verifies it under Namespace with key. It returns the first error.
func VerifyApproval(armored []byte, digest, challenge [32]byte, expiry int64, key SSHPublicKey) error {
	s, err := ParseSSHSIG(armored)
	if err != nil {
		return err
	}
	msg, err := SignedMessage(digest, challenge, expiry)
	if err != nil {
		return err
	}
	return VerifySSHSIG(s, Namespace, msg, key)
}
