// Package approvaltest is TEST SUPPORT for host/approval: an SSHSIG signer over a fixed Ed25519
// seed, so tests make keys and signatures with no randomness, no ssh-keygen and no ~/.ssh. Only
// _test.go files may import it (TestApprovalTestIsTestOnly); production code never does.
package approvaltest

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"strings"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// Signer is an Ed25519 key that signs as `ssh-keygen -Y sign` does (SSHSIG, sha512).
type Signer struct{ priv ed25519.PrivateKey }

// NewSigner returns the signer whose 32-byte Ed25519 seed is the byte seed, repeated.
func NewSigner(seed byte) *Signer {
	return &Signer{priv: ed25519.NewKeyFromSeed(bytes.Repeat([]byte{seed}, ed25519.SeedSize))}
}

func sshString(dst, s []byte) []byte {
	dst = binary.BigEndian.AppendUint32(dst, uint32(len(s)))
	return append(dst, s...)
}

func (s *Signer) blob() []byte {
	return sshString(sshString(nil, []byte("ssh-ed25519")), s.priv.Public().(ed25519.PublicKey))
}

// PublicKey is the signer's key as host/approval holds it.
func (s *Signer) PublicKey() approval.SSHPublicKey {
	return approval.SSHPublicKey{Type: "ssh-ed25519", Blob: s.blob()}
}

// AuthorizedKey is the signer's .pub file: "ssh-ed25519 <base64> approvaltest\n".
func (s *Signer) AuthorizedKey() []byte {
	return []byte("ssh-ed25519 " + base64.StdEncoding.EncodeToString(s.blob()) + " approvaltest\n")
}

// Sign returns the armored SSHSIG signature of msg under namespace, hashed with sha512, wrapped
// at 70 columns as ssh-keygen writes it.
func (s *Signer) Sign(namespace string, msg []byte) []byte {
	h := sha512.Sum512(msg)
	signed := []byte("SSHSIG")
	signed = sshString(signed, []byte(namespace))
	signed = sshString(signed, nil)
	signed = sshString(signed, []byte("sha512"))
	signed = sshString(signed, h[:])
	sig := sshString(sshString(nil, []byte("ssh-ed25519")), ed25519.Sign(s.priv, signed))
	blob := binary.BigEndian.AppendUint32([]byte("SSHSIG"), 1)
	blob = sshString(blob, s.blob())
	blob = sshString(blob, []byte(namespace))
	blob = sshString(blob, nil)
	blob = sshString(blob, []byte("sha512"))
	blob = sshString(blob, sig)
	b64 := base64.StdEncoding.EncodeToString(blob)
	var out strings.Builder
	out.WriteString("-----BEGIN SSH SIGNATURE-----\n")
	for len(b64) > 70 {
		out.WriteString(b64[:70] + "\n")
		b64 = b64[70:]
	}
	out.WriteString(b64 + "\n-----END SSH SIGNATURE-----\n")
	return []byte(out.String())
}
