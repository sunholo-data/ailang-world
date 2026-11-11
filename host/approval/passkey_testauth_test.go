package approval_test

// Row 201 A2 sprint 2 M3b2 (design §6, TP1, XW1): the test authenticator approvaltest.Passkey (a
// fixed seed, never a real key: D-WORLD-109) makes files the production law accepts, in Chrome's
// layout. This is an external test package because approvaltest imports host/approval (sprint 1
// §1.4 P4).

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
)

const paOrigin, paRPID = "https://x.ts.net:8443", "x.ts.net"

var paDigest, paChallenge, paEnrol = [32]byte{1}, [32]byte{2}, [32]byte{3}

// paHash is the WebAuthn challenge of the PA request: its message hash (design §2.1).
func paHash(t *testing.T, digest, challenge [32]byte, expiry int64) [32]byte {
	t.Helper()
	h, err := approval.MessageHash(digest, challenge, expiry)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

// TP1: for both algorithms the test authenticator's registration enrols and its assertion approves,
// at a tailnet origin and at localhost; the approval names the credential and the canonical key.
func TestTestPasskeyVerifies(t *testing.T) {
	for _, alg := range []int64{approval.AlgES256, approval.AlgEdDSA} {
		for origin, rpID := range map[string]string{paOrigin: paRPID, "http://localhost:7644": "localhost"} {
			pk := approvaltest.NewPasskey(alg, 5)
			a, err := approval.NewPasskeyApprover(approval.RegistrationInput{Name: "mark-mac", Origin: origin, Challenge: paEnrol,
				Response: pk.Register(origin, rpID, paEnrol), Now: 900})
			if err != nil || !bytes.Equal(a.Passkey.CredentialID, pk.CredentialID()) || !bytes.Equal(a.Passkey.Key.Canonical(), pk.Key().Canonical()) {
				t.Fatalf("%d at %s: enrol: %v", alg, origin, err)
			}
			s, err := approval.ApproveByAssertion(approval.AssertionInput{
				Admit:     approval.AdmitInput{Now: 1010, CreatedAt: 1000, ExpiresAt: 1600, StoredDigest: paDigest, BytesDigest: paDigest},
				Challenge: paChallenge, Approver: a,
				Response: pk.Assert(origin, rpID, paHash(t, paDigest, paChallenge, 1600), 2),
			})
			if err != nil {
				t.Fatalf("%d at %s: assert: %v", alg, origin, err)
			}
			_, _, key, _ := s.Use()
			if id, n, ok := s.PasskeyUse(); !ok || n != 2 || !bytes.Equal(id, pk.CredentialID()) || !bytes.Equal(key, pk.Key().Canonical()) {
				t.Fatalf("%d at %s: PasskeyUse (%x, %d, %v), key %x", alg, origin, id, n, ok, key)
			}
		}
	}
}

// chromeParts is one recorded ceremony's envelope members and decoded response fields.
func chromeParts(t *testing.T, file []byte) (top, response []string, fields map[string][]byte) {
	t.Helper()
	var env, resp map[string]json.RawMessage
	if json.Unmarshal(file, &env) != nil || json.Unmarshal(env["response"], &resp) != nil {
		t.Fatalf("not an envelope: %s", file)
	}
	fields = map[string][]byte{}
	for k := range env {
		top = append(top, k)
	}
	for k, v := range resp {
		response = append(response, k)
		var s string
		if json.Unmarshal(v, &s) == nil {
			fields[k], _ = base64.RawURLEncoding.DecodeString(s)
		}
	}
	return top, response, fields
}

// memberOrder is the member names of a JSON object, in order.
func memberOrder(t *testing.T, b []byte) []string {
	t.Helper()
	d := json.NewDecoder(bytes.NewReader(b))
	var names []string
	if _, err := d.Token(); err != nil {
		t.Fatal(err)
	}
	for d.More() {
		k, err := d.Token()
		if err != nil {
			t.Fatal(err)
		}
		names = append(names, k.(string))
		var skip json.RawMessage
		if err := d.Decode(&skip); err != nil {
			t.Fatal(err)
		}
	}
	return names
}

// XW1: the test authenticator's files have Chrome's skeleton: the envelope members it writes are
// Chrome's, the clientDataJSON member order, the attestationObject's bytes outside authData, the
// flags, counters and AAGUID, the credential ID length, and the COSE key's bytes with x and y masked.
func TestTestPasskeyMatchesChromeLayout(t *testing.T) {
	for alg, name := range map[int64]string{approval.AlgES256: "es256", approval.AlgEdDSA: "eddsa"} {
		raw, err := os.ReadFile(filepath.Join("testdata", "webauthn", "chrome-"+name+".json"))
		if err != nil {
			t.Fatal(err)
		}
		var rec struct{ Registration, Assertion json.RawMessage }
		if err := json.Unmarshal(raw, &rec); err != nil {
			t.Fatal(err)
		}
		pk := approvaltest.NewPasskey(alg, 7)
		for _, c := range []struct {
			kind         string
			chrome, ours []byte
		}{
			{"registration", rec.Registration, pk.Register("http://localhost:56547", "localhost", [32]byte{})},
			{"assertion", rec.Assertion, pk.Assert("http://localhost:56547", "localhost", [32]byte{}, 2)},
		} {
			ctop, cresp, cf := chromeParts(t, c.chrome)
			otop, oresp, of := chromeParts(t, c.ours)
			for _, sub := range []struct{ ours, chrome []string }{{otop, ctop}, {oresp, cresp}} {
				for _, m := range sub.ours {
					if !strings.Contains(" "+strings.Join(sub.chrome, " ")+" ", " "+m+" ") {
						t.Errorf("%s %s: member %q is not one Chrome writes (%v)", name, c.kind, m, sub.chrome)
					}
				}
			}
			if got, want := strings.Join(memberOrder(t, of["clientDataJSON"]), ","), strings.Join(memberOrder(t, cf["clientDataJSON"]), ","); got != want {
				t.Errorf("%s %s: clientDataJSON members %s, Chrome's %s", name, c.kind, got, want)
			}
			if c.kind == "assertion" {
				if a, b := of["authenticatorData"], cf["authenticatorData"]; len(a) != 37 || len(b) != 37 || !bytes.Equal(a[32:], b[32:]) || !bytes.Equal(a[:32], b[:32]) {
					t.Errorf("%s assertion authData %x, Chrome's %x", name, a, b)
				}
				continue
			}
			oa, ca := of["attestationObject"], cf["attestationObject"]
			if len(oa) != len(ca) || !bytes.Equal(oa[:30], ca[:30]) { // fmt "none", attStmt {}, authData's head
				t.Fatalf("%s attestationObject head %x (%d bytes), Chrome's %x (%d bytes)", name, oa[:30], len(oa), ca[:30], len(ca))
			}
			oad, cad := oa[30:], ca[30:]
			mask := func(ad []byte) []byte {
				m := bytes.Clone(ad)
				copy(m[55:87], make([]byte, 32))  // the credential ID's bytes
				copy(m[97:129], make([]byte, 32)) // x
				if alg == approval.AlgES256 {
					copy(m[132:164], make([]byte, 32)) // y
				}
				return m
			}
			if !bytes.Equal(mask(oad), mask(cad)) {
				t.Errorf("%s registration authData skeleton:\n ours  %x\n chrome %x", name, mask(oad), mask(cad))
			}
		}
	}
}
