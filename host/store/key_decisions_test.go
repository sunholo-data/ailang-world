package store

// Row 201 A4 sprint 2 M4e (design §2.7, §6 KD1, KD5, KD6; KD3, KD4, KD7 and KD9 are M4f's, in
// key_decisions_record_test.go; KD2's counter is the passkey arm's, sprint 3): credentials schema
// 8's key_decisions. Every store is under t.TempDir() (D-WORLD-109).

import (
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// canonicalKeyDecisionsDDL is the current (version 8, row 201 A4) key_decisions table: the
// design's DDL, both sha256 CHECKs verbatim (KD5, KD6).
var canonicalKeyDecisionsDDL = `CREATE TABLE key_decisions (
    request_ref   TEXT PRIMARY KEY CHECK (length(request_ref) = 71 AND request_ref GLOB 'sha256:*'),
    packet_ref    TEXT NOT NULL CHECK (length(packet_ref) = 71 AND packet_ref GLOB 'sha256:*'),
    verb          TEXT NOT NULL CHECK (verb IN ('approve', 'deny')),
    format        TEXT NOT NULL CHECK (format IN ('ssh', 'webauthn')),
    approver      TEXT NOT NULL CHECK (approver <> ''),
    artifact      BLOB NOT NULL CHECK (length(artifact) BETWEEN 1 AND 16384),
    expiry        INTEGER NOT NULL,
    decided_at    INTEGER NOT NULL,
    decision_ref  TEXT CHECK (decision_ref IS NULL OR (length(decision_ref) = 71 AND decision_ref GLOB 'sha256:*'))
)`

var kdQ, kdP, kdD = "sha256:" + strings.Repeat("a1", 32), "sha256:" + strings.Repeat("b2", 32), "sha256:" + strings.Repeat("c3", 32)

// kdStore is a schema-8 store with mark-ssh (an ssh approver, key "k1") enrolled at 900.
func kdStore(t *testing.T) *Credentials {
	t.Helper()
	c := openTestCredentials(t)
	if err := c.EnrolApprover(testCtx(t), ApproverRow{Name: "mark-ssh", Format: "ssh", PublicKey: []byte("k1"), EnrolledAt: 900}); err != nil {
		t.Fatal(err)
	}
	return c
}

var kdUse = ApproverUse{Name: "mark-ssh", Format: "ssh", PublicKey: []byte("k1")}

func kdRow() KeyDecisionRow {
	return KeyDecisionRow{RequestRef: kdQ, PacketRef: kdP, Verb: "approve", Format: "ssh", Approver: "ssh:mark-ssh", Artifact: []byte("sig-1"), Expiry: 1760}
}

func kdRead(t *testing.T, c *Credentials) (KeyDecisionRow, bool) {
	t.Helper()
	r, ok, err := c.KeyDecision(testCtx(t), kdQ)
	if err != nil {
		t.Fatal(err)
	}
	return r, ok
}

// KD1: a fresh store is schema 8 with key_decisions as credentials_schema.sql and the upgrade
// write it; a schema-7 store is upgraded by UpgradeCredentials only, and before that every key
// decision method refuses naming the serve restart.
func TestCredentialsSchema8KeyDecisions(t *testing.T) {
	c := openTestCredentials(t)
	for name, got := range map[string]string{"schema file": tableDDL(t, c.db)["key_decisions"], "keyDecisionsSQL": keyDecisionsSQL} {
		if normalizeDDL(got) != normalizeDDL(canonicalKeyDecisionsDDL) {
			t.Fatalf("key_decisions DDL (%s):\n got: %s\nwant: %s", name, got, canonicalKeyDecisionsDDL)
		}
	}
	if got, want := credentialsTables(8), []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "session_credentials", "surface_approvals", "webauthn_credentials"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("credentialsTables(8) = %v", got)
	}
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	makeV6Credentials(t, path, []SessionRow{oneBoundRow()})
	raw := rawDB(t, path)
	if _, err := raw.Exec(surfaceApprovalsSQL); err != nil {
		t.Fatal(err)
	}
	setPositiveFixtureVersion(t, raw, 7)
	v7, err := OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	var oe *OlderCredentialsStoreError
	_, _, readErr := v7.KeyDecision(testCtx(t), kdQ)
	for name, err := range map[string]error{"record": v7.RecordKeyDecision(testCtx(t), kdRow(), kdUse, 1000), "land": v7.LandKeyDecision(testCtx(t), kdQ, kdD), "read": readErr} {
		if !errors.As(err, &oe) || oe.Version != 7 || oe.Need != 8 || !strings.Contains(err.Error(), "key decisions need schema 8: restart `ailang-worldd serve`") {
			t.Errorf("v7 %s: %v, want the restart refusal naming schema 8", name, err)
		}
	}
	_ = v7.Close()
	if upgraded, err := UpgradeCredentials(path); err != nil || !upgraded || userVersion(t, path) != 8 {
		t.Fatalf("UpgradeCredentials: %v %v, version %d", upgraded, err, userVersion(t, path))
	}
	up := tableDDL(t, rawDB(t, path))
	requireExactTableNames(t, "upgraded v7 store", up, credentialsTables(8))
	if normalizeDDL(up["key_decisions"]) != normalizeDDL(canonicalKeyDecisionsDDL) {
		t.Fatalf("upgraded key_decisions DDL: %s", up["key_decisions"])
	}
}

// KD5, KD6: the DDL's ref CHECKs, by raw INSERTs: packet_ref and decision_ref of length 71 not
// starting sha256:, of length 70, and decision_ref "" are constraint errors; sha256: + 64 hex and
// a NULL decision_ref are accepted (the controls).
func TestKeyDecisionsRefChecks(t *testing.T) {
	c := openTestCredentials(t)
	ins := func(request, packet string, decision any) error {
		_, err := c.db.Exec(`INSERT INTO key_decisions (request_ref, packet_ref, verb, format, approver, artifact, expiry, decided_at, decision_ref) VALUES (?, ?, 'approve', 'ssh', 'ssh:m', x'01', 1, 1, ?)`, request, packet, decision)
		return err
	}
	sha512 := "sha512:" + strings.Repeat("ab", 32)
	for i, tc := range []struct {
		packet   string
		decision any
		ok       bool
	}{
		{kdP, nil, true}, {kdP, kdD, true}, // controls
		{sha512, nil, false}, {kdP[:70], nil, false}, // KD5
		{kdP, sha512, false}, {kdP, kdD[:70], false}, {kdP, "", false}, // KD6
	} {
		request := "sha256:" + strings.Repeat(string(rune('0'+i)), 64)
		if err := ins(request, tc.packet, tc.decision); (err == nil) != tc.ok {
			t.Errorf("packet %q decision %v: err %v, want ok %v", tc.packet, tc.decision, err, tc.ok)
		}
	}
}
