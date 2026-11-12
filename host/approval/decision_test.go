package approval

// Row 201 A4 sprint 1 (design §2.3, §6 DM1, DM2): the decision message's layout and refusals.

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"testing"
)

// DM1: 91 bytes; the prefix, 0x00, the verb at 18, P at 19, Q at 51, the big-endian expiry at 83;
// golden hex for two subjects (computed by an independent encoder, not by DecisionMessage).
func TestDecisionMessageLayout(t *testing.T) {
	p1, q1 := sha256.Sum256([]byte("packet-1")), sha256.Sum256([]byte("request-1"))
	var p2, q2 [32]byte
	for i := range p2 {
		p2[i], q2[i] = byte(i), byte(0xff-i)
	}
	for _, tc := range []struct {
		verb         byte
		p, q         [32]byte
		expiry       int64
		golden, hash string
	}{
		{VerbApprove, p1, q1, 1794414540,
			"776f726c642d6465636973696f6e2f76310001c48e8d2fd69d694f90962ad5c6da898c6e83f070c32804601d4c8b4b0285412719f1064b619d49d35392eac7261cd7266c720671fc594f4b226f32bf0bee74ba000000006af497cc",
			"efa6cc1b6c482505379b440b52639ae3db80351938d4e3addf1119553b28c4b7"},
		{VerbDeny, p2, q2, 0x0102030405060708,
			"776f726c642d6465636973696f6e2f76310002000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1ffffefdfcfbfaf9f8f7f6f5f4f3f2f1f0efeeedecebeae9e8e7e6e5e4e3e2e1e00102030405060708",
			"080bf79e5fa2f6a74c8726f47460ad398f2794235f5e92acbe1e51b35cc35ec2"},
	} {
		m, err := DecisionMessage(tc.verb, tc.p, tc.q, tc.expiry)
		if err != nil {
			t.Fatal(err)
		}
		if len(m) != 91 || DecisionMessageLen != 91 {
			t.Fatalf("len = %d (DecisionMessageLen %d), want 91", len(m), DecisionMessageLen)
		}
		if got := hex.EncodeToString(m); got != tc.golden {
			t.Errorf("message = %s, want %s", got, tc.golden)
		}
		if string(m[:17]) != "world-decision/v1" || m[17] != 0 || m[18] != tc.verb ||
			!bytes.Equal(m[19:51], tc.p[:]) || !bytes.Equal(m[51:83], tc.q[:]) || binary.BigEndian.Uint64(m[83:]) != uint64(tc.expiry) {
			t.Errorf("layout of %x", m)
		}
		h, err := DecisionMessageHash(tc.verb, tc.p, tc.q, tc.expiry)
		if err != nil || hex.EncodeToString(h[:]) != tc.hash {
			t.Errorf("hash = %x, %v; want %s", h, err, tc.hash)
		}
	}
	// Domain separation from the mint message (§2.3): a different prefix and length.
	mint, err := SignedMessage(p1, q1, 1794414540)
	if err != nil {
		t.Fatal(err)
	}
	m, _ := DecisionMessage(VerbApprove, p1, q1, 1794414540)
	if len(mint) == len(m) || bytes.HasPrefix(m, mint[:18]) || DecisionNamespace == Namespace || DecisionNamespace != "world-decision@sunholo" {
		t.Error("the decision message is not separated from the mint message")
	}
	if v, err := DecisionVerb("approve"); v != VerbApprove || err != nil {
		t.Errorf("approve = %d, %v", v, err)
	}
	if v, err := DecisionVerb("deny"); v != VerbDeny || err != nil {
		t.Errorf("deny = %d, %v", v, err)
	}
}

// DM2: verb 0 and 3 and a negative expiry are refused, by the message and by its hash.
func TestDecisionMessageRefusesBadInput(t *testing.T) {
	var p, q [32]byte
	for _, verb := range []byte{0, 3, 0xff} {
		if m, err := DecisionMessage(verb, p, q, 1); m != nil || !errors.Is(err, ErrDecisionVerb) {
			t.Errorf("verb %d: %x, %v", verb, m, err)
		}
		if _, err := DecisionMessageHash(verb, p, q, 1); !errors.Is(err, ErrDecisionVerb) {
			t.Errorf("hash verb %d: %v", verb, err)
		}
	}
	if m, err := DecisionMessage(VerbDeny, p, q, -1); m != nil || !errors.Is(err, ErrNegativeExpiry) {
		t.Errorf("negative expiry: %x, %v", m, err)
	}
	if m, err := DecisionMessage(VerbApprove, p, q, 0); err != nil || len(m) != 91 {
		t.Errorf("expiry 0 (control): %v", err)
	}
	for _, word := range []string{"", "Approve", "approve ", "yes"} {
		if _, err := DecisionVerb(word); !errors.Is(err, ErrDecisionVerb) {
			t.Errorf("word %q: %v", word, err)
		}
	}
}
