package store

// Row 243 V3 sprint 3 M7 (design_docs/planned/w-world-address-bar-verbs-v3.md §2.4, §6 CS1-CS7):
// credentials schema 9's verb_requests. Every store is under t.TempDir() (D-WORLD-109); every
// request's bytes are approval.CanonicalVerb's.

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
)

var vrDecision = "sha256:" + strings.Repeat("d4", 32)

// vrRow is a pending verb request vr-<16 hex of n> for episode ep, created at created and expiring
// 600 s later: its bytes are approval.CanonicalVerb of a valid request, its digest their SHA-256.
func vrRow(t *testing.T, n int, ep string, created int64) VerbRequestRow {
	t.Helper()
	r := approval.VerbRequest{Schema: approval.VerbRequestSchema, Verb: approval.VerbFinish, Argv: []string{"ailang-worldd", "episode", "finish", ep, "--preset", "p"},
		Episode: ep, Prompt: fmt.Sprintf("committer: (repository git config)\nrun %d\n", n), Scope: "github-pr:o/r#scope", Title: "Fix the bar", Body: "body", Protected: []string{}}
	if err := r.Validate(); err != nil {
		t.Fatal(err)
	}
	b, err := approval.CanonicalVerb(r)
	if err != nil {
		t.Fatal(err)
	}
	return VerbRequestRow{RequestID: fmt.Sprintf("vr-%016x", n), Request: b, Digest: sha256Hex(b), Challenge: sha256Hex([]byte(fmt.Sprint(n))),
		EpisodeID: ep, CreatedAt: created, ExpiresAt: created + 600}
}

// vrStore is a fresh schema-9 store with ssh approver mark-ssh (key k1) and passkey x (counter 5)
// enrolled at 900.
func vrStore(t *testing.T) *Credentials {
	t.Helper()
	c := openTestCredentials(t)
	if err := c.EnrolApprover(testCtx(t), ApproverRow{Name: "mark-ssh", Format: "ssh", PublicKey: []byte("k1"), EnrolledAt: 900}); err != nil {
		t.Fatal(err)
	}
	a, p := testPasskey("x", 1)
	a.EnrolledAt, p.CredentialID, p.SignCount = 900, credX, 5
	if err := c.EnrolPasskey(testCtx(t), a, p); err != nil {
		t.Fatal(err)
	}
	return c
}

func vrPut(t *testing.T, c *Credentials, row VerbRequestRow, now int64) {
	t.Helper()
	if err := c.PutVerbRequest(testCtx(t), row, now); err != nil {
		t.Fatalf("put %s: %v", row.RequestID, err)
	}
}

func vrGet(t *testing.T, c *Credentials, id string) VerbRequestRow {
	t.Helper()
	r, ok, err := c.VerbRequest(testCtx(t), id)
	if err != nil || !ok {
		t.Fatalf("read %s: ok %v err %v", id, ok, err)
	}
	return r
}

// vrState is every row of verb_requests, the approvers' last use and revocation, and the passkeys' counters.
func vrState(t *testing.T, c *Credentials) string {
	t.Helper()
	var out string
	if err := c.db.QueryRow(`SELECT coalesce((SELECT group_concat(quote(request_id)||quote(request)||quote(consumed_at)||quote(outcome)||quote(approved_by)||quote(artifact)||quote(decision_ref), ';') FROM verb_requests), '')
		|| '|' || (SELECT group_concat(name||quote(last_used_at)||quote(revoked_at), ';') FROM approvers) || '|' || (SELECT group_concat(name||sign_count, ';') FROM webauthn_credentials)`).Scan(&out); err != nil {
		t.Fatal(err)
	}
	return out
}

// tableInfo is PRAGMA table_info of every table of db, in order.
func tableInfo(t *testing.T, db *sql.DB) (out []string) {
	t.Helper()
	rows, err := db.Query(`SELECT m.name, p.name, p.type, p."notnull", quote(p.dflt_value), p.pk FROM sqlite_master m, pragma_table_info(m.name) p WHERE m.type = 'table' ORDER BY m.name, p.cid`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var table, col, typ, def string
		var notNull, pk int
		if err := rows.Scan(&table, &col, &typ, &notNull, &def, &pk); err != nil {
			t.Fatal(err)
		}
		out = append(out, fmt.Sprintf("%s.%s %s notnull=%d default=%s pk=%d", table, col, typ, notNull, def, pk))
	}
	return out
}

// CS1, CS2, CS7(b): a real schema-8 store (makeV6Credentials, then the
// version 7 and 8 tables) opens at 8 and refuses every verb request method naming the serve
// restart; UpgradeCredentials takes it through 9 to the current version, and every table's columns and verb_requests'
// whole DDL, equal the fresh store's.
func TestCredentialsSchema9VerbRequests(t *testing.T) {
	fresh := openTestCredentials(t)
	if got, want := credentialsTables(9), []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "session_credentials", "surface_approvals", "verb_requests", "webauthn_credentials"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("credentialsTables(9) = %v", got)
	}
	if fresh.Version() != CredentialsSchemaVersion || MaxVerbRequestBytes != approval.MaxVerbRequestBytes {
		t.Fatalf("fresh version %d, constant %d", fresh.Version(), CredentialsSchemaVersion)
	}
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	makeV6Credentials(t, path, []SessionRow{oneBoundRow()})
	raw := rawDB(t, path)
	for _, stmt := range []string{surfaceApprovalsSQL, keyDecisionsSQL} {
		if _, err := raw.Exec(stmt); err != nil {
			t.Fatal(err)
		}
	}
	setPositiveFixtureVersion(t, raw, 8)
	v8, err := OpenCredentials(path)
	if err != nil || v8.Version() != 8 {
		t.Fatalf("open v8: %v", err)
	}
	ctx := testCtx(t)
	_, _, readErr := v8.VerbRequest(ctx, "vr-0000000000000001")
	_, listErr := v8.PendingVerbRequests(ctx, 1000, 0)
	for name, err := range map[string]error{"put": v8.PutVerbRequest(ctx, vrRow(t, 1, "w1", 1000), 1000), "read": readErr, "list": listErr,
		"consume": v8.ConsumeVerbRequest(ctx, "vr-0000000000000001", strings.Repeat("a", 64), kdUse, "ssh:mark-ssh", []byte("sig"), vrDecision, 1000)} {
		var oe *OlderCredentialsStoreError
		if !errors.As(err, &oe) || oe.Version != 8 || oe.Need != 9 || !strings.Contains(err.Error(), "verb requests need schema 9: restart `ailang-worldd serve`") {
			t.Errorf("v8 %s: %v, want the restart refusal naming schema 9", name, err)
		}
	}
	_ = v8.Close()
	if upgraded, err := UpgradeCredentials(path); err != nil || !upgraded || userVersion(t, path) != CredentialsSchemaVersion {
		t.Fatalf("UpgradeCredentials: %v %v, version %d", upgraded, err, userVersion(t, path))
	}
	if again, err := UpgradeCredentials(path); err != nil || again {
		t.Fatalf("second UpgradeCredentials: %v %v", again, err)
	}
	up := rawDB(t, path)
	requireExactTableNames(t, "upgraded v8 store", tableDDL(t, up), credentialsTables(CredentialsSchemaVersion))
	if got, want := tableInfo(t, up), tableInfo(t, fresh.db); !reflect.DeepEqual(got, want) || strings.Count(strings.Join(want, "\n"), "verb_requests.") != 12 {
		t.Fatalf("upgraded columns differ from a fresh store's:\n got: %v\nwant: %v", got, want)
	}
	if got, want := normalizeDDL(tableDDL(t, up)["verb_requests"]), normalizeDDL(tableDDL(t, fresh.db)["verb_requests"]); got != want {
		t.Fatalf("upgraded verb_requests DDL:\n got: %s\nwant: %s", got, want)
	}
}

// CS7(a): a store newer than this binary (user_version 11) is refused by both opens with the
// existing "newer than this binary's" text, now at 10, and the refusal changes nothing.
func TestCredentialsNewerThanTenIsRefused(t *testing.T) {
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	c, err := OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	_ = c.Close()
	db := rawDB(t, path)
	setPositiveFixtureVersion(t, db, 11)
	before, _ := schemaState(t, db)
	for name, open := range map[string]func(string) (*Credentials, error){"writable": OpenCredentials, "read-only": OpenCredentialsReadOnly, "upgrade": func(p string) (*Credentials, error) { _, err := UpgradeCredentials(p); return nil, err }} {
		c, err := open(path)
		if c != nil {
			_ = c.Close()
		}
		var se *CredentialsSchemaError
		if !errors.As(err, &se) || !strings.Contains(err.Error(), "user_version 11 is newer than this binary's 10; use a compatible binary") {
			t.Errorf("%s open of a v11 store: %v", name, err)
		}
	}
	if after, v := schemaState(t, db); !reflect.DeepEqual(after, before) || v != 11 {
		t.Fatalf("refused opens changed the store: %v -> %v (version %d)", before, after, v)
	}
}

// CS3: a second request for an episode supersedes its pending one (consumed_at = now, outcome
// superseded); another episode's request stays pending.
func TestPutVerbRequestSupersedes(t *testing.T) {
	c := vrStore(t)
	a, b, n := vrRow(t, 1, "w1", 1000), vrRow(t, 2, "w2", 1000), vrRow(t, 3, "w1", 1010)
	vrPut(t, c, a, 1000)
	vrPut(t, c, b, 1000)
	vrPut(t, c, n, 1010)
	if got := vrGet(t, c, a.RequestID); got.Outcome != VerbOutcomeSuperseded || got.ConsumedAt != 1010 || got.ApprovedBy != "" || got.DecisionRef != "" || got.Artifact != nil {
		t.Fatalf("superseded row = %+v", got)
	}
	for _, id := range []string{b.RequestID, n.RequestID} {
		if got := vrGet(t, c, id); got.Outcome != "" || got.ConsumedAt != 0 {
			t.Fatalf("%s = %+v, want pending", id, got)
		}
	}
	got, err := c.PendingVerbRequests(testCtx(t), 1010, 0)
	if err != nil || len(got) != 2 || got[0].RequestID != n.RequestID || got[1].RequestID != b.RequestID || !bytes.Equal(got[0].Request, n.Request) {
		t.Fatalf("pending = %+v, %v", got, err)
	}
	vrPut(t, c, vrRow(t, 4, "w4", 1600+86400), 1600+86400) // a day after a's expiry: the prune takes pending rows only
	vrGet(t, c, a.RequestID)                               // the superseded row is kept as the record
}

// CS4: the grammar (each refusal ErrBadVerbRequest, nothing written), the 16 cap (a supersede frees
// its episode's place first), and the prune of requests expired a day ago.
func TestPutVerbRequestGrammarCapPrune(t *testing.T) {
	c := vrStore(t)
	big := bytes.Repeat([]byte("x"), MaxVerbRequestBytes+1)
	for name, edit := range map[string]func(*VerbRequestRow){
		"id-upper": func(r *VerbRequestRow) { r.RequestID = "vr-000000000000000A" }, "id-mint": func(r *VerbRequestRow) { r.RequestID = "mr-0000000000000001" },
		"id-short":         func(r *VerbRequestRow) { r.RequestID = "vr-0001" },
		"request-empty":    func(r *VerbRequestRow) { r.Request, r.Digest = nil, sha256Hex(nil) },
		"request-128KiB+1": func(r *VerbRequestRow) { r.Request, r.Digest = big, sha256Hex(big) },
		"digest-upper":     func(r *VerbRequestRow) { r.Digest = strings.ToUpper(r.Digest) },
		"challenge-short":  func(r *VerbRequestRow) { r.Challenge = r.Challenge[:62] },
		"digest-not-sha":   func(r *VerbRequestRow) { r.Digest = r.Challenge },
		"episode-empty":    func(r *VerbRequestRow) { r.EpisodeID = "" },
		"expires-created":  func(r *VerbRequestRow) { r.ExpiresAt = r.CreatedAt },
		"outcome-set":      func(r *VerbRequestRow) { r.Outcome = VerbOutcomeApproved },
		"consumed-set":     func(r *VerbRequestRow) { r.ConsumedAt = 1000 },
		"approver-set":     func(r *VerbRequestRow) { r.ApprovedBy = "ssh:mark-ssh" },
		"artifact-set":     func(r *VerbRequestRow) { r.Artifact = []byte("sig") },
		"decision-set":     func(r *VerbRequestRow) { r.DecisionRef = vrDecision },
	} {
		r := vrRow(t, 1, "w1", 1000)
		edit(&r)
		if err := c.PutVerbRequest(testCtx(t), r, 1000); !errors.Is(err, ErrBadVerbRequest) {
			t.Errorf("%s: %v, want ErrBadVerbRequest", name, err)
		}
	}
	if s := vrState(t, c); strings.Contains(s, "'vr-") {
		t.Fatalf("a refused put wrote a row: %s", s)
	}
	max := vrRow(t, 99, "wmax", 1000)
	max.Request, max.Digest = bytes.Repeat([]byte("y"), approval.MaxVerbRequestBytes), sha256Hex(bytes.Repeat([]byte("y"), approval.MaxVerbRequestBytes))
	vrPut(t, c, max, 1000) // control: exactly 128 KiB, the approval law's bound
	old, kept := vrRow(t, 40, "wold", 1000-87000), vrRow(t, 41, "wkept", 1000-86999)
	vrPut(t, c, old, 1000-87000)  // expires 1000-86400: expired a day ago at 1000
	vrPut(t, c, kept, 1000-86999) // expires 1000-86399
	for i := 1; i <= 13; i++ {
		vrPut(t, c, vrRow(t, i, fmt.Sprintf("w%d", i), 1000), 1000)
	}
	if err := c.PutVerbRequest(testCtx(t), vrRow(t, 14, "w14", 1000), 1000); err != nil { // the 15th live (with wmax), old pruned
		t.Fatal(err)
	}
	if _, ok, _ := c.VerbRequest(testCtx(t), old.RequestID); ok {
		t.Fatal("a request expired a day ago was not pruned")
	}
	if _, ok, _ := c.VerbRequest(testCtx(t), kept.RequestID); !ok {
		t.Fatal("a request expired less than a day ago was pruned")
	}
	vrPut(t, c, vrRow(t, 15, "w15", 1000), 1000) // 16 live
	before := vrState(t, c)
	if err := c.PutVerbRequest(testCtx(t), vrRow(t, 16, "w16", 1000), 1000); !errors.Is(err, ErrTooManyPendingVerbRequests) {
		t.Fatalf("the 17th live request: %v", err)
	}
	if vrState(t, c) != before {
		t.Fatal("the refused 17th request wrote")
	}
	vrPut(t, c, vrRow(t, 17, "w1", 1000), 1000) // w1's request is superseded first: still 16 live
	if live, err := c.PendingVerbRequests(testCtx(t), 1000, 100); err != nil || len(live) != MaxPendingVerbRequests {
		t.Fatalf("live after a supersede = %d, %v", len(live), err)
	}
}

// CS5: the table's CHECKs, by raw INSERTs (backstops behind the Go grammar).
func TestVerbRequestsChecks(t *testing.T) {
	c := openTestCredentials(t)
	id := 0
	ins := func(reqID string, request []byte, consumed, outcome, approvedBy, artifact, decision any) error {
		id++
		if reqID == "" {
			reqID = fmt.Sprintf("vr-%016x", id)
		}
		_, err := c.db.Exec(`INSERT INTO verb_requests (request_id, request, digest, challenge, episode_id, created_at, expires_at, consumed_at, outcome, approved_by, artifact, decision_ref) VALUES (?, ?, ?, ?, 'w1', 1000, 1600, ?, ?, ?, ?, ?)`,
			reqID, request, strings.Repeat("a", 64), strings.Repeat("b", 64), consumed, outcome, approvedBy, artifact, decision)
		return err
	}
	r := []byte("r")
	for _, tc := range []struct {
		name                                     string
		reqID                                    string
		request                                  []byte
		consumed, outcome, by, artifact, decided any
		ok                                       bool
	}{
		{"pending", "", r, nil, nil, nil, nil, nil, true},
		{"superseded", "", r, 1100, "superseded", nil, nil, nil, true},
		{"approved", "", r, 1100, "approved", "ssh:m", []byte("sig"), vrDecision, true},
		{"id-mint", "mr-0000000000000001", r, nil, nil, nil, nil, nil, false},
		{"id-length", "vr-01", r, nil, nil, nil, nil, nil, false},
		{"request-empty", "", []byte{}, nil, nil, nil, nil, nil, false},
		{"request-128KiB+1", "", make([]byte, 131073), nil, nil, nil, nil, nil, false},
		{"outcome-unknown", "", r, 1100, "denied", nil, nil, nil, false},
		{"outcome-unconsumed", "", r, nil, "superseded", nil, nil, nil, false},
		{"consumed-no-outcome", "", r, 1100, nil, nil, nil, nil, false},
		{"approved-no-by", "", r, 1100, "approved", nil, []byte("sig"), vrDecision, false},
		{"approved-no-artifact", "", r, 1100, "approved", "ssh:m", nil, vrDecision, false},
		{"approved-no-decision", "", r, 1100, "approved", "ssh:m", []byte("sig"), nil, false},
		// Design §6 CS5 says a decision_ref on a superseded row is refused; the design's DDL (kept
		// verbatim) does not refuse it: the approved-iff CHECK holds when all three are not set. Only
		// ConsumeVerbRequest writes decision_ref, with outcome approved (plan §1.4 F3).
		{"superseded-with-decision-accepted", "", r, 1100, "superseded", nil, nil, vrDecision, true},
		{"superseded-with-approver", "", r, 1100, "superseded", "ssh:m", []byte("sig"), vrDecision, false},
		{"by-empty", "", r, 1100, "approved", "", []byte("sig"), vrDecision, false},
		{"artifact-16385", "", r, 1100, "approved", "ssh:m", make([]byte, 16385), vrDecision, false},
		{"decision-sha512", "", r, 1100, "approved", "ssh:m", []byte("sig"), "sha512:" + strings.Repeat("d4", 32), false},
		{"decision-70", "", r, 1100, "approved", "ssh:m", []byte("sig"), vrDecision[:70], false},
	} {
		if err := ins(tc.reqID, tc.request, tc.consumed, tc.outcome, tc.by, tc.artifact, tc.decided); (err == nil) != tc.ok {
			t.Errorf("%s: err %v, want ok %v", tc.name, err, tc.ok)
		}
	}
}

// CS6: ConsumeVerbRequest's refusals each write nothing (rows, approvers' last use, counters); a
// confirmation sets the five columns, decision_ref included, marks the approver used and moves a
// passkey's counter forward, in one transaction.
func TestConsumeVerbRequest(t *testing.T) {
	sshUse := kdUse
	keyUse := ApproverUse{Name: "x", Format: "webauthn", PublicKey: []byte("cose-x"), Passkey: &PasskeyUse{CredentialID: credX, SignCount: 6}}
	type call struct {
		id, digest, by string
		use            ApproverUse
		artifact       []byte
		decision       string
		now            int64
	}
	setup := func(t *testing.T) (*Credentials, VerbRequestRow, call) {
		c := vrStore(t)
		r := vrRow(t, 1, "w1", 1000)
		vrPut(t, c, r, 1000)
		return c, r, call{r.RequestID, r.Digest, "ssh:mark-ssh", sshUse, []byte("sig"), vrDecision, 1100}
	}
	consume := func(c *Credentials, k call) error {
		return c.ConsumeVerbRequest(testCtx(t), k.id, k.digest, k.use, k.by, k.artifact, k.decision, k.now)
	}
	for _, tc := range []struct {
		name string
		prep func(*testing.T, *Credentials, VerbRequestRow)
		edit func(*call)
		want error
	}{
		{"unknown", nil, func(k *call) { k.id = "vr-00000000000000ff" }, ErrVerbRequestUnknown},
		{"consumed", func(t *testing.T, c *Credentials, r VerbRequestRow) {
			if err := c.ConsumeVerbRequest(testCtx(t), r.RequestID, r.Digest, sshUse, "ssh:mark-ssh", []byte("sig"), vrDecision, 1050); err != nil {
				t.Fatal(err)
			}
		}, nil, ErrVerbRequestConsumed},
		{"superseded", func(t *testing.T, c *Credentials, _ VerbRequestRow) { vrPut(t, c, vrRow(t, 2, "w1", 1050), 1050) }, nil, ErrVerbRequestSuperseded},
		{"clock-behind", nil, func(k *call) { k.now = 999 }, ErrVerbRequestClockBehind},
		{"expired", nil, func(k *call) { k.now = 1600 }, ErrVerbRequestExpired},
		{"digest-other", nil, func(k *call) { k.digest = strings.Repeat("e", 64) }, ErrVerbRequestDigestChanged},
		{"bytes-edited", func(t *testing.T, c *Credentials, r VerbRequestRow) {
			if _, err := c.db.Exec(`UPDATE verb_requests SET request = ? WHERE request_id = ?`, append(r.Request, ' '), r.RequestID); err != nil {
				t.Fatal(err)
			}
		}, nil, ErrVerbRequestDigestChanged},
		{"approver-revoked", func(t *testing.T, c *Credentials, _ VerbRequestRow) {
			if err := c.RevokeApprover(testCtx(t), "mark-ssh", 1050); err != nil {
				t.Fatal(err)
			}
		}, nil, ErrApproverRevoked},
		{"approver-key-changed", nil, func(k *call) { k.use.PublicKey = []byte("k2") }, ErrApproverKeyChanged},
		{"counter-not-advancing", nil, func(k *call) {
			k.use, k.by = keyUse, "webauthn:x"
			k.use.Passkey = &PasskeyUse{CredentialID: credX, SignCount: 5}
		}, ErrSignCountRegression},
		{"by-not-the-approver", nil, func(k *call) { k.by = "ssh:other" }, ErrBadVerbRequest},
		{"webauthn-without-passkey", nil, func(k *call) { k.use, k.by = keyUse, "webauthn:x"; k.use.Passkey = nil }, ErrBadVerbRequest},
		{"artifact-empty", nil, func(k *call) { k.artifact = nil }, ErrBadVerbRequest},
		{"artifact-16385", nil, func(k *call) { k.artifact = make([]byte, 16385) }, ErrBadVerbRequest},
		{"decision-upper", nil, func(k *call) { k.decision = strings.ToUpper(vrDecision) }, ErrBadVerbRequest},
		{"id-bad", nil, func(k *call) { k.id = "mr-0000000000000001" }, ErrBadVerbRequest},
		{"digest-bad", nil, func(k *call) { k.digest = "x" }, ErrBadVerbRequest},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c, r, k := setup(t)
			if tc.prep != nil {
				tc.prep(t, c, r)
			}
			if tc.edit != nil {
				tc.edit(&k)
			}
			before := vrState(t, c)
			if err := consume(c, k); !errors.Is(err, tc.want) || tc.name == "consumed" && !strings.Contains(err.Error(), "by ssh:mark-ssh") {
				t.Fatalf("err %v, want %v", err, tc.want)
			}
			if after := vrState(t, c); after != before {
				t.Fatalf("a refusal wrote:\n%s->\n%s", before, after)
			}
		})
	}
	for _, tc := range []struct {
		name, by   string
		use        ApproverUse
		count, now int64
	}{{"ssh", "ssh:mark-ssh", sshUse, 5, 1000}, {"webauthn", "webauthn:x", keyUse, 6, 1599}} { // at creation; a second before expiry
		t.Run("confirm-"+tc.name, func(t *testing.T) {
			c, r, k := setup(t)
			k.use, k.by, k.now = tc.use, tc.by, tc.now
			if err := consume(c, k); err != nil {
				t.Fatal(err)
			}
			got := vrGet(t, c, r.RequestID)
			if got.ConsumedAt != tc.now || got.Outcome != VerbOutcomeApproved || got.ApprovedBy != tc.by || !bytes.Equal(got.Artifact, []byte("sig")) || got.DecisionRef != vrDecision {
				t.Fatalf("confirmed row = %+v", got)
			}
			if p, _, _ := c.Passkey(testCtx(t), "x"); p.SignCount != tc.count {
				t.Fatalf("counter %d, want %d", p.SignCount, tc.count)
			}
			if a, _, _ := c.Approver(testCtx(t), tc.use.Name); a.LastUsedAt != tc.now {
				t.Fatalf("%s last used %d", tc.use.Name, a.LastUsedAt)
			}
			if live, _ := c.PendingVerbRequests(testCtx(t), tc.now, 0); len(live) != 0 {
				t.Fatalf("still pending: %+v", live)
			}
			vrPut(t, c, vrRow(t, 2, "w1", tc.now), tc.now) // a later finish of w1 supersedes only a pending request
			if again := vrGet(t, c, r.RequestID); !reflect.DeepEqual(again, got) {
				t.Fatalf("the confirmed row changed: %+v", again)
			}
		})
	}
}
