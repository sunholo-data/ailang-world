package surfacepage

// Row 201 A3 sprint 2 M4a (design §2.3, §2.4, §6 AM1-AM3): the approval routes' JSON, pure.

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
)

const amCreated, amExpires = int64(1800000000), int64(1800000600)

var amChallenge = strings.Repeat("5a", 32)

// amRequest is a stored request whose bytes are canonical and hash to its digest.
func amRequest(t *testing.T) StoredRequest {
	t.Helper()
	b, err := approval.Canonical(approval.MintRequest{Schema: approval.RequestSchema, Episode: "w201x",
		WorldDB: "/ws/world.db", WorkspaceRoot: "/ws/root", Worktree: "/ws/root/w201x", Git: "/usr/bin/git", Branch: "world/w201x",
		Grants:     []approval.Grant{{Effect: "Workspace.Read", Scope: "w201x"}, {Effect: "Workspace.Write", Scope: "w201x", Budget: 5}, {Effect: "Exec", Scope: "w201x"}},
		TTLSeconds: 43200, ProfilesDir: "/ws/profiles", CredentialsDB: "/ws/credentials.db", SessionFile: "/ws/s/session.token",
		MCPConfig: "/ws/s/mcp.json", Agent: "claude", Model: "claude-opus<5>&x", Task: "/ws/task.md", Daemon: approval.DaemonRef{Addr: "127.0.0.1:7644"}})
	if err != nil {
		t.Fatal(err)
	}
	d := sha256.Sum256(b)
	return StoredRequest{ID: "mr-0123456789abcdef", Episode: "w201x", Request: b, Digest: hex.EncodeToString(d[:]), Challenge: amChallenge,
		Created: amCreated, Expires: amExpires}
}

// jsonKeys is the sorted top-level keys of v's JSON.
func jsonKeys(t *testing.T, v any) []string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// AM1: a pending request's line is decoded from its bytes; hash16 is the CLI's message hash
// prefix; recorded is null until an approval is recorded; a decided line truncates the credential.
func TestApproveModelSummaries(t *testing.T) {
	r := amRequest(t)
	p := Pending(r, "")
	d, _ := hex32(r.Digest)
	h, err := approval.MessageHash(d, [32]byte(bytes.Repeat([]byte{0x5a}, 32)), amExpires)
	if err != nil {
		t.Fatal(err)
	}
	if p.ID != r.ID || p.Episode != "w201x" || p.Agent != "claude" || p.Model != "claude-opus<5>&x" || p.Grants != 3 ||
		p.TTL == nil || *p.TTL != 43200 || p.Expires != amExpires || p.Hash16 != hex.EncodeToString(h[:8]) || p.Recorded != nil || p.Refused != "" {
		t.Fatalf("pending = %+v", p)
	}
	b, _ := json.Marshal(p)
	if !strings.Contains(string(b), `"recorded":null`) {
		t.Fatalf("an unrecorded line is %s, want recorded null", b)
	}
	if got := Pending(r, "webauthn:mark"); got.Recorded == nil || *got.Recorded != "webauthn:mark" {
		t.Fatalf("recorded = %v", got.Recorded)
	}
	if want := []string{"agent", "episode", "expires", "grants", "hash16", "id", "model", "recorded", "ttl_s"}; !reflect.DeepEqual(jsonKeys(t, p), want) {
		t.Fatalf("pending keys %v, want %v (no bytes, challenge, artifact or path)", jsonKeys(t, p), want)
	}
	r.Consumed, r.ConsumedAt, r.ApprovedBy, r.CredentialID = true, amCreated+9, "webauthn:mark", "cred-0123456789abcdefXYZ"
	if got := Decided(r); got != (DecidedRequest{ID: r.ID, Episode: "w201x", ApprovedBy: "webauthn:mark", At: amCreated + 9, Credential16: "cred-0123456789a"}) {
		t.Fatalf("decided = %+v", got)
	}
	r.CredentialID = "short"
	if got := Decided(r); got.Credential16 != "short" {
		t.Fatalf("a short credential id became %q", got.Credential16)
	}
}

// AM2: a row whose bytes do not hash to its digest, do not decode canonically, or whose challenge
// is not hex is listed with refused and no summary, never hidden.
func TestApproveModelRefusesWhatDoesNotDecode(t *testing.T) {
	for name, edit := range map[string]func(*StoredRequest){
		"digest": func(r *StoredRequest) { r.Digest = strings.Repeat("0", 64) },
		"canonical": func(r *StoredRequest) {
			r.Request = append(r.Request, ' ')
			d := sha256.Sum256(r.Request)
			r.Digest = hex.EncodeToString(d[:])
		},
		"challenge":       func(r *StoredRequest) { r.Challenge = "zz" },
		"challenge-short": func(r *StoredRequest) { r.Challenge = amChallenge[2:] }, // 31 bytes
		"expiry":          func(r *StoredRequest) { r.Expires = -1 },
	} {
		r := amRequest(t)
		edit(&r)
		p := Pending(r, "ssh:mark")
		if p.Refused == "" || p.ID != r.ID || p.Episode != "w201x" || p.Agent != "" || p.Grants != 0 || p.TTL != nil || p.Hash16 != "" || p.Recorded == nil {
			t.Errorf("%s: %+v, want refused with no summary", name, p)
		}
		if want := []string{"episode", "expires", "id", "recorded", "refused"}; !reflect.DeepEqual(jsonKeys(t, p), want) {
			t.Errorf("%s: keys %v, want %v", name, jsonKeys(t, p), want)
		}
	}
}

// AM2 (the card): the stored bytes as base64url, the row's digest, nonce and times, the active
// passkeys grouped by origin, the SSH namespace, and no message-hash field.
func TestApproveModelCard(t *testing.T) {
	r := amRequest(t)
	c := RequestCard(r, []CardPasskey{{Origin: "https://s.tail0.ts.net", RPID: "s.tail0.ts.net", CredentialID: []byte{1}},
		{Origin: "http://localhost:7644", RPID: "localhost", CredentialID: []byte{2}}, {Origin: "https://s.tail0.ts.net", RPID: "s.tail0.ts.net", CredentialID: []byte{3}}})
	if b, err := base64.RawURLEncoding.DecodeString(c.Request); err != nil || !bytes.Equal(b, r.Request) {
		t.Fatalf("request does not decode to the stored bytes: %v", err)
	}
	want := []CardPasskeys{{Origin: "https://s.tail0.ts.net", RP: "s.tail0.ts.net", Allow: []string{"AQ", "Aw"}}, {Origin: "http://localhost:7644", RP: "localhost", Allow: []string{"Ag"}}}
	if c.V != 1 || c.ID != r.ID || c.Digest != r.Digest || c.Challenge != amChallenge || c.Created != amCreated || c.Expires != amExpires ||
		c.Consumed || c.ApprovedBy != "" || !reflect.DeepEqual(c.Passkeys, want) || c.SSH.Namespace != "world-approval@sunholo" {
		t.Fatalf("card = %+v", c)
	}
	if want := []string{"approvedBy", "challenge", "consumed", "created", "digest", "expires", "id", "passkeys", "request", "ssh", "v"}; !reflect.DeepEqual(jsonKeys(t, c), want) {
		t.Fatalf("card keys %v, want %v (no message hash)", jsonKeys(t, c), want)
	}
	if c := RequestCard(r, nil); c.Passkeys == nil || len(c.Passkeys) != 0 {
		t.Fatalf("no passkeys = %#v, want []", c.Passkeys)
	}
}
