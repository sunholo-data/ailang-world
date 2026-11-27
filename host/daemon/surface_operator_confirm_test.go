package daemon

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"html"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sunholo-data/ailang-world/host/surfacepage"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/agentrun"
	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
)

func (f *operatorHTTP) row(id string) store.OperatorActionRow {
	f.t.Helper()
	ro, err := store.OpenCredentialsReadOnly(f.d.credentialsPath)
	if err != nil {
		f.t.Fatal(err)
	}
	defer ro.Close()
	row, err := ro.OperatorAction(boundedTestContext(f.t), id, 0)
	if err != nil {
		f.t.Fatal(err)
	}
	return row
}

func opMessage(t *testing.T, digest, nonce [32]byte, expires int64) []byte {
	t.Helper()
	m, err := approval.OperatorActionMessage(digest, nonce, expires)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestOperatorConfirmHTTP(t *testing.T) {
	f := newOperatorHTTP(t, time.Minute)
	testOperatorConfirm(t, f)
}

type confirmEvidence struct {
	row, counters, credentials, world  [32]byte
	actions, operations, countersCount int
	fixtureCount                       int64
}

// Observe raw durable state, so logical expiry and replay cannot mask a write.
func (f *operatorHTTP) confirmationSnapshot(db *sql.DB, id string) confirmEvidence {
	f.t.Helper()
	var out confirmEvidence
	var row, counters string
	var effects, starts, sessions int
	err := db.QueryRowContext(boundedTestContext(f.t), `SELECT
		(SELECT json_array(id, episode_id, kind, state, hex(review), hex(digest), hex(action), hex(challenge),
			created_at, expires_at, generation, owner, approved_by, hex(artifact), terminal_at, v3_id)
			FROM operator_actions WHERE id=?),
		(SELECT json_group_array(json_array(name, sign_count)) FROM
			(SELECT name, sign_count FROM webauthn_credentials ORDER BY name)),
		(SELECT count(*) FROM webauthn_credentials),
		coalesce((SELECT sign_count FROM webauthn_credentials WHERE name='fixture-passkey'), 0),
		(SELECT count(*) FROM operator_actions),
		(SELECT count(*) FROM operator_action_steps WHERE kind='operation'),
		(SELECT count(*) FROM operator_action_steps WHERE kind<>'operation'),
		(SELECT count(*) FROM episode_starts), (SELECT count(*) FROM session_credentials)`, id).
		Scan(&row, &counters, &out.countersCount, &out.fixtureCount, &out.actions, &out.operations, &effects, &starts, &sessions)
	if err != nil || effects != 0 || starts != 0 || sessions != 0 {
		f.t.Fatalf("confirmation crossed the effect boundary: effects=%d starts=%d sessions=%d (%v)", effects, starts, sessions, err)
	}
	if _, err := os.Lstat(f.fake); !os.IsNotExist(err) {
		f.t.Fatalf("unexpected fake child marker: %v", err)
	}
	hash := func(path string) [32]byte {
		f.t.Helper()
		b, err := os.ReadFile(path)
		if err != nil {
			f.t.Fatal(err)
		}
		return sha256.Sum256(b)
	}
	out.row, out.counters = sha256.Sum256([]byte(row)), sha256.Sum256([]byte(counters))
	out.credentials, out.world = hash(f.d.credentialsPath), hash(f.d.cfg.DBPath)
	return out
}

func (f *operatorHTTP) confirmSQL(query string, args ...any) {
	f.t.Helper()
	// These are fixture mutations. Confirmation still uses the normal writer DSN.
	dsn := url.URL{Scheme: "file", Path: f.d.credentialsPath, RawQuery: "_pragma=synchronous%28OFF%29"}
	db, err := sql.Open("sqlite", dsn.String())
	if err != nil {
		f.t.Fatal(err)
	}
	defer db.Close()
	if _, err = db.ExecContext(boundedTestContext(f.t), query, args...); err != nil {
		f.t.Fatal(err)
	}
}

func testOperatorConfirm(t *testing.T, f *operatorHTTP) {
	var lockCheck atomic.Bool
	var expireBetweenReads atomic.Int64
	f.d.approveNow = func() int64 {
		if lockCheck.Load() && f.d.humanMu.TryLock() {
			f.d.humanMu.Unlock()
			f.t.Error("confirm reached its store read outside the human writer lock")
		}
		if expireBetweenReads.Load() > 0 && expireBetweenReads.Add(-1) == 0 {
			return f.now.Load() + 1
		}
		return f.now.Load()
	}
	key := approvaltest.NewPasskey(approval.AlgEdDSA, 43)
	signer := approvaltest.NewSigner(45)
	f.enrolSSH("fixture-ssh", signer)
	f.enrolPasskey("fixture-passkey", key, opOrigin)
	s := f.request(opRequest(1))
	row := f.row(s.ID)
	digest, nonce := [32]byte(row.Digest), [32]byte(row.Challenge)
	d := hex.EncodeToString(row.Digest)
	m := opMessage(t, digest, nonce, row.ExpiresAt)
	challenge := sha256.Sum256(m)
	sig := func(key *approvaltest.Signer, namespace string, msg []byte) []byte {
		return sshBody(d, key.Sign(namespace, msg))
	}
	signed := sig(signer, approval.OperatorActionNamespace, m)
	passkey := approveBody(d, key.Assert(opOrigin, "localhost", challenge, 2))
	dbURL := url.URL{Scheme: "file", Path: f.d.credentialsPath, RawQuery: "mode=ro"}
	evidenceDB, err := sql.Open("sqlite", dbURL.String())
	if err != nil {
		t.Fatal(err)
	}
	defer evidenceDB.Close()
	before := f.confirmationSnapshot(evidenceDB, s.ID)
	post := func(body []byte) *httptest.ResponseRecorder {
		return f.do("POST", s.Href+"/confirm", body, func(r *http.Request) { r.Host = "localhost:7644" })
	}
	checkRefusal := func(r *httptest.ResponseRecorder, status int, class string) {
		f.t.Helper()
		requireRefusal(f.t, class, r, status, class)
		for k := range r.Header() {
			if k == "Location" || strings.HasPrefix(k, "Access-Control-") {
				f.t.Fatalf("confirmation refusal emitted %s", k)
			}
		}
		if f.confirmationSnapshot(evidenceDB, s.ID) != before {
			f.t.Fatalf("confirmation refusal %s changed the row, counter, or effect evidence", class)
		}
	}
	refuse := func(b []byte, status int, class string) { checkRefusal(post(b), status, class) }
	f.sub("Transport", func(t *testing.T) {
		for _, tc := range []struct {
			key, value, class string
			code              int
		}{
			{"Origin", "", "CrossOrigin", 403}, {"Origin", "null", "CrossOrigin", 403}, {"Origin", opOrigin + ".invalid", "CrossOrigin", 403},
			{"Sec-Fetch-Site", "same-site", "CrossOrigin", 403}, {"Sec-Fetch-Site", "Same-Origin", "CrossOrigin", 403},
			{"Sec-Fetch-Site", "same-origin, cross-site", "CrossOrigin", 403}, {"Content-Type", "text/plain", "UnsupportedMediaType", 415},
			{"Authorization", "Bearer y", "Forbidden", 403}, {"Cookie", "y=y", "Forbidden", 403},
		} {
			r := f.do("POST", s.Href+"/confirm", signed, func(r *http.Request) {
				if tc.value == "" {
					r.Header.Del(tc.key)
				} else {
					r.Header.Set(tc.key, tc.value)
				}
			})
			checkRefusal(r, tc.code, tc.class)
		}
		checkRefusal(f.do("POST", s.Href+"/confirm", signed, func(r *http.Request) { r.Host = "other.invalid:7644" }), 421, "MisdirectedRequest")
		for _, header := range []string{"Origin", "Sec-Fetch-Site", "Content-Type"} {
			code, class := 403, "CrossOrigin"
			if header == "Content-Type" {
				code, class = 415, "UnsupportedMediaType"
			}
			checkRefusal(f.do("POST", s.Href+"/confirm", signed, func(r *http.Request) {
				r.Header.Add(header, r.Header.Get(header))
			}), code, class)
		}
		req := fixtureConfirmRequest(s.Href+"/confirm", bytes.NewReader(signed))
		spy := &confirmReadSpy{ReadCloser: req.Body}
		req.Host, req.Body = "localhost:7644.evil", spy
		w := httptest.NewRecorder()
		f.handler.ServeHTTP(w, req)
		checkRefusal(w, 421, "MisdirectedRequest")
		if spy.reads != 0 {
			t.Fatalf("foreign Host read a valid proof body %d times", spy.reads)
		}
		for _, method := range []string{"GET", "OPTIONS", "PUT"} {
			checkRefusal(f.do(method, s.Href+"/confirm", signed, func(r *http.Request) {
				if method == "OPTIONS" {
					r.Header.Set("Origin", opOrigin)
					r.Header.Set("Access-Control-Request-Method", "POST")
					r.Header.Set("Access-Control-Request-Headers", "content-type")
				}
			}), 405, "MethodNotAllowed")
		}
		padded := append(bytes.Clone(signed), bytes.Repeat([]byte(" "), (20<<10)+1-len(signed))...)
		refuse(padded, 413, "TooLarge")
		for _, path := range []string{strings.TrimSuffix(s.Href, s.ID) + strings.ToUpper(s.ID),
			s.Href + "a", strings.TrimSuffix(s.Href, s.ID) + s.ID[:len(s.ID)-1], s.Href + "%0a",
			operatorPath + "//" + s.ID, operatorPath + "/x/../" + s.ID, "/surface/x/../operator-actions/" + s.ID,
			strings.Replace(s.Href, "/oa-", "/%6fa-", 1),
		} {
			checkRefusal(f.do("POST", path+"/confirm", signed, nil), 400, "BadRequest")
		}
		checkRefusal(f.do("POST", s.Href+"/confirm?next=elsewhere", signed, nil), 400, "BadRequest")
		checkRefusal(f.do("POST", s.Href+"/confirm/", signed, nil), 400, "BadRequest")
		checkRefusal(f.do("POST", s.Href+"/confirm/x", signed, nil), 404, "NotFound")
	})
	f.sub("Proof", func(t *testing.T) {
		for _, body := range [][]byte{[]byte("y"), opJSON(map[string]any{"v": 1, "digest": d}),
			opJSON(map[string]any{"v": 1, "digest": d, "y": true}), append([]byte(`{"v":1,`), passkey[1:]...), append([]byte(`{"sshsig":null,`), passkey[1:]...)} {
			refuse(body, 400, "BadRequest")
		}
		refuse(bytes.Replace(signed, []byte(`"v":1`), []byte(`"v":2`), 1), 400, "BadRequest")
		refuse(approveBody(d, []byte("null")), 400, "BadRequest")
		refuse(sshBody(d, key.Assert(opOrigin, "localhost", challenge, 2)), 400, "BadRequest")
		refuse(sshBody(d, []byte("malformed")), 422, "NotApproved")
		refuse(sshBody(strings.Repeat("0", 64), signer.Sign(approval.OperatorActionNamespace, m)), 409, "Conflict")
		other := nonce
		other[0] ^= 1
		refuse(sig(signer, approval.OperatorActionNamespace, opMessage(t, digest, other, row.ExpiresAt)), 422, "NotApproved")
	})
	f.sub("Authenticity", func(t *testing.T) {
		refuse(sig(approvaltest.NewSigner(47), approval.OperatorActionNamespace, m), 422, "NotApproved")
		forged := bytes.Clone(m)
		forged[0] ^= 1
		refuse(sig(signer, approval.OperatorActionNamespace, forged), 422, "NotApproved")
		refuse(sig(signer, approval.VerbNamespace, m), 422, "NotApproved")
		old, err := approval.VerbMessage(digest, nonce, row.ExpiresAt)
		if err != nil {
			t.Fatal(err)
		}
		refuse(sig(signer, approval.VerbNamespace, old), 422, "NotApproved")
		oldMint, err := approval.SignedMessage(digest, nonce, row.ExpiresAt)
		if err != nil {
			t.Fatal(err)
		}
		refuse(sig(signer, approval.Namespace, oldMint), 422, "NotApproved")
		for _, proof := range [][]byte{
			key.Assert(opOrigin, "localhost", challenge, 2, approvaltest.Edit{BadSignature: true}),
			key.Assert(opOrigin, "localhost", challenge, 2, approvaltest.Edit{Clear: 0x04}),
			key.Assert("http://localhost:7645", "localhost", challenge, 2),
			key.Assert(opOrigin, "elsewhere.invalid", challenge, 2),
			approvaltest.NewPasskey(approval.AlgEdDSA, 44).Assert(opOrigin, "localhost", challenge, 2),
		} {
			refuse(approveBody(d, proof), 422, "NotApproved")
		}
	})
	f.sub("Binding", func(t *testing.T) {
		var a map[string]json.RawMessage
		if err := json.Unmarshal(row.Action, &a); err != nil {
			t.Fatal(err)
		}
		for field, value := range map[string]any{"id": "oa-" + strings.Repeat("a", 32), "episode": "other", "kind": "finish",
			"store": "/other/credentials.db", "argv": []string{"different"}, "options": approval.OperatorOptions{Preset: "other", Headless: true}, "task": "changed task",
		} {
			match := string(opJSON(field)) + ":" + string(a[field])
			replace := string(opJSON(field)) + ":" + string(opJSON(value))
			alternate := bytes.Replace(row.Action, []byte(match), []byte(replace), 1)
			if bytes.Equal(row.Action, alternate) {
				t.Fatalf("binding control did not change %s", field)
			}
			refuse(sig(signer, approval.OperatorActionNamespace, opMessage(t, sha256.Sum256(alternate), nonce, row.ExpiresAt)), 422, "NotApproved")
		}
	})
	f.sub("Expiry", func(t *testing.T) {
		for _, now := range []int64{row.CreatedAt - 1, row.ExpiresAt, row.ExpiresAt + 1} {
			f.now.Store(now)
			refuse(signed, 409, "Expired")
		}
		f.now.Store(row.ExpiresAt - 1)
		expireBetweenReads.Store(2)
		refuse(signed, 409, "Expired")
		f.now.Store(apT0)
		db, err := sql.Open("sqlite", f.d.credentialsPath)
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		change := func(exp int64) {
			if _, err := db.ExecContext(boundedTestContext(t), "UPDATE operator_actions SET expires_at=? WHERE id=?", exp, s.ID); err != nil {
				t.Fatal(err)
			}
			before = f.confirmationSnapshot(evidenceDB, s.ID)
		}
		change(row.ExpiresAt + 1)
		refuse(sig(signer, approval.OperatorActionNamespace, opMessage(t, digest, nonce, row.ExpiresAt+1)), 409, "Expired")
		change(row.ExpiresAt)
	})
	f.sub("RevocationCounter", func(t *testing.T) {
		for _, count := range []uint32{0, 1} {
			refuse(approveBody(d, key.Assert(opOrigin, "localhost", challenge, count)), 422, "NotApproved")
		}
		revoked := approvaltest.NewSigner(46)
		f.enrolSSH("revoked-ssh", revoked)
		f.write(func(ctx context.Context, c *store.Credentials) error {
			return c.RevokeApprover(ctx, "revoked-ssh", apT0)
		})
		before = f.confirmationSnapshot(evidenceDB, s.ID)
		refuse(sig(revoked, approval.OperatorActionNamespace, m), 422, "NotApproved")
	})
	f.sub("UnsupportedKinds", func(t *testing.T) {
		for i, kind := range []string{"start", "finish", "foreign"} {
			a, err := approval.DecodeOperatorAction(row.Action)
			if err != nil {
				t.Fatal(err)
			}
			a.ID, a.Episode, a.Kind = "oa-"+opRequest(i+30).OperationID, "other-"+kind, kind
			a.Options.Headless = false
			if kind == "foreign" {
				a.Kind, a.Options.Headless = "start", true
			}
			if kind == "finish" {
				a.Task = ""
			}
			a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
			var stored store.OperatorActionRow
			f.write(func(ctx context.Context, c *store.Credentials) error {
				var err error
				stored, err = c.RequestOperatorAction(ctx, store.OperatorRequest{
					Action: a, OperationID: opRequest(i + 30).OperationID, Challenge: nonce, CreatedAt: apT0, ExpiresAt: apT0 + 600,
				}, apT0)
				return err
			})
			if kind == "foreign" {
				a.Store += ".foreign"
				wire, err := approval.CanonicalOperatorAction(a)
				if err != nil {
					t.Fatal(err)
				}
				digest := sha256.Sum256(wire)
				f.confirmSQL("UPDATE operator_actions SET action=?, digest=? WHERE id=?", wire, digest[:], a.ID)
				stored.Digest = digest[:]
			}
			before = f.confirmationSnapshot(evidenceDB, s.ID)
			body := sshBody(hex.EncodeToString(stored.Digest), signer.Sign(approval.OperatorActionNamespace,
				opMessage(t, [32]byte(stored.Digest), nonce, stored.ExpiresAt)))
			checkRefusal(f.do("POST", operatorPath+"/"+a.ID+"/confirm", body, nil), 409, "Conflict")
		}
	})
	f.sub("StoredIdentity", func(t *testing.T) {
		for _, edit := range []struct{ change, restore string }{
			{"episode_id='drift'", "episode_id=" + "'" + row.Episode + "'"},
			{"kind='finish'", "kind='start'"},
			{"digest=zeroblob(32)", "digest=X'" + d + "'"},
		} {
			f.confirmSQL("UPDATE operator_actions SET "+edit.change+" WHERE id=?", s.ID)
			before = f.confirmationSnapshot(evidenceDB, s.ID)
			proof := signed
			if edit.change == "digest=zeroblob(32)" {
				proof = sshBody(strings.Repeat("0", 64), signer.Sign(approval.OperatorActionNamespace, m))
			}
			refuse(proof, 409, "Conflict")
			f.confirmSQL("UPDATE operator_actions SET "+edit.restore+" WHERE id=?", s.ID)
		}
		a, err := approval.DecodeOperatorAction(row.Action)
		if err != nil {
			t.Fatal(err)
		}
		a.ID = "oa-" + strings.Repeat("b", 32)
		other, err := approval.CanonicalOperatorAction(a)
		if err != nil {
			t.Fatal(err)
		}
		different := sha256.Sum256(other)
		f.confirmSQL("UPDATE operator_actions SET action=?, digest=? WHERE id=?", other, different[:], s.ID)
		before = f.confirmationSnapshot(evidenceDB, s.ID)
		refuse(sshBody(hex.EncodeToString(different[:]), signer.Sign(approval.OperatorActionNamespace,
			opMessage(t, different, nonce, row.ExpiresAt))), 409, "Conflict")
		f.confirmSQL("UPDATE operator_actions SET action=?, digest=? WHERE id=?", row.Action, row.Digest, s.ID)
		f.confirmSQL("PRAGMA ignore_check_constraints=ON; UPDATE operator_actions SET challenge=zeroblob(31) WHERE id=?", s.ID)
		before = f.confirmationSnapshot(evidenceDB, s.ID)
		refuse(signed, 409, "Conflict")
		card := f.do("GET", s.Href, nil, nil)
		if card.Code != 200 || strings.Contains(card.Body.String(), "operator-ceremony") {
			t.Fatalf("card offered proof for an invalid challenge: %d %s", card.Code, card.Body)
		}
		f.confirmSQL("UPDATE operator_actions SET challenge=? WHERE id=?", row.Challenge, s.ID)
		before = f.confirmationSnapshot(evidenceDB, s.ID)
	})
	f.sub("Cancelled", func(t *testing.T) {
		other := f.request(opRequest(29))
		f.write(func(ctx context.Context, c *store.Credentials) error {
			return c.CancelUnstartedOperatorAction(ctx, other.ID, apT0+1)
		})
		before = f.confirmationSnapshot(evidenceDB, s.ID)
		cancelled := f.row(other.ID)
		checkRefusal(f.do("POST", other.Href+"/confirm", sshBody(hex.EncodeToString(cancelled.Digest), []byte("malformed")), nil), 409, "Conflict")
	})
	f.sub("MissingAndWriter", func(t *testing.T) {
		path := f.d.credentialsPath
		if err := os.Rename(path, path+".offline"); err != nil {
			t.Fatal(err)
		}
		r := post(signed)
		if err := os.Rename(path+".offline", path); err != nil {
			t.Fatal(err)
		}
		checkRefusal(r, 503, "Disconnected")
		checkRefusal(f.do("POST", operatorPath+"/oa-"+strings.Repeat("0", 32)+"/confirm", signed, nil), 404, "NotFound")
		held, err := store.OpenCredentials(f.d.credentialsPath)
		if err != nil {
			t.Fatal(err)
		}
		refuse(signed, 503, "Busy")
		if err := held.Close(); err != nil {
			t.Fatal(err)
		}
		f.confirmSQL("ALTER TABLE operator_actions RENAME COLUMN action TO unreadable_action")
		readFault := f.evidence()
		f.refuse(readFault, post(signed), 503, "Disconnected")
		f.confirmSQL("ALTER TABLE operator_actions RENAME COLUMN unreadable_action TO action")
		after := f.confirmationSnapshot(evidenceDB, s.ID)
		if before.row != after.row || before.counters != after.counters || before.countersCount != after.countersCount {
			t.Fatal("read failure changed the pending action or counters")
		}
		before = after
		f.confirmSQL("UPDATE approvers SET enrolled_at='unreadable' WHERE name='fixture-passkey'")
		before = f.confirmationSnapshot(evidenceDB, s.ID)
		refuse(passkey, 503, "Disconnected")
		f.confirmSQL("UPDATE approvers SET enrolled_at=? WHERE name='fixture-passkey'", apT0-10)
		before = f.confirmationSnapshot(evidenceDB, s.ID)
	})
	f.sub("Rollback", func(t *testing.T) {
		f.confirmSQL(`CREATE TRIGGER interrupt_confirm BEFORE UPDATE OF approved_by ON operator_actions
			BEGIN SELECT RAISE(ABORT, 'fixture: interrupted after passkey use'); END`)
		before = f.confirmationSnapshot(evidenceDB, s.ID)
		refuse(passkey, 503, "Disconnected")
		f.confirmSQL("DROP TRIGGER interrupt_confirm")
		before = f.confirmationSnapshot(evidenceDB, s.ID)
		f.confirmSQL(`CREATE TRIGGER concurrent_counter BEFORE UPDATE OF last_used_at ON approvers
			BEGIN UPDATE webauthn_credentials SET sign_count=2 WHERE name='fixture-passkey'; END`)
		before = f.confirmationSnapshot(evidenceDB, s.ID)
		refuse(passkey, 422, "NotApproved")
		f.confirmSQL("DROP TRIGGER concurrent_counter")
		f.confirmSQL(`CREATE TRIGGER concurrent_revocation BEFORE UPDATE OF last_used_at ON approvers
			WHEN NEW.name='fixture-passkey' BEGIN
				UPDATE approvers SET revoked_at=NEW.last_used_at WHERE name=NEW.name;
				SELECT RAISE(IGNORE);
			END`)
		before = f.confirmationSnapshot(evidenceDB, s.ID)
		refuse(passkey, 422, "NotApproved")
		f.confirmSQL("DROP TRIGGER concurrent_revocation")
		before = f.confirmationSnapshot(evidenceDB, s.ID)
	})
	f.sub("Ceremony", func(t *testing.T) {
		requireOperatorCeremony(t, f.do("GET", s.Href, nil, nil), s, d, m, key)
		f.now.Store(row.CreatedAt - 1)
		card := f.do("GET", s.Href, nil, nil)
		if card.Code != 200 || strings.Contains(card.Body.String(), "operator-ceremony") {
			t.Fatalf("card offered a ceremony before issue: %d %s", card.Code, card.Body)
		}
		f.now.Store(apT0)
		f.confirmSQL("UPDATE approvers SET enrolled_at='unreadable' WHERE name='fixture-passkey'")
		before = f.confirmationSnapshot(evidenceDB, s.ID)
		checkRefusal(f.do("GET", s.Href, nil, nil), 503, "Disconnected")
		f.confirmSQL("UPDATE approvers SET enrolled_at=? WHERE name='fixture-passkey'", apT0-10)
		f.confirmSQL("UPDATE operator_actions SET created_at=-2, expires_at=-1 WHERE id=?", s.ID)
		f.now.Store(-2)
		before = f.confirmationSnapshot(evidenceDB, s.ID)
		checkRefusal(f.do("GET", s.Href, nil, nil), 503, "Disconnected")
		f.confirmSQL("UPDATE operator_actions SET created_at=?, expires_at=? WHERE id=?", row.CreatedAt, row.ExpiresAt, s.ID)
		f.now.Store(apT0)
		before = f.confirmationSnapshot(evidenceDB, s.ID)
	})
	expect := func(r *httptest.ResponseRecorder, id, result, key string) {
		t.Helper()
		var out operatorConfirmationAnswer
		if err := json.Unmarshal(r.Body.Bytes(), &out); err != nil || r.Code != 200 || out.Result != result || out.ApprovedBy != key ||
			out.ID != id || out.State != "confirmed" {
			t.Fatalf("confirmation %s: %d %s (%v)", result, r.Code, r.Body, err)
		}
	}
	f.sub("ConcurrentReplay", func(t *testing.T) {
		f.confirmSQL("DELETE FROM operator_runner_state")
		before = f.confirmationSnapshot(evidenceDB, s.ID)
		f.now.Store(row.ExpiresAt - 1)
		lockCheck.Store(true)
		defer lockCheck.Store(false)
		// Two handlers start together. Both must return; only one may advance the counter.
		var wg sync.WaitGroup
		var replies [2]*httptest.ResponseRecorder
		ready := make(chan struct{})
		for i := range replies {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-ready
				replies[i] = post(passkey)
			}()
		}
		close(ready)
		wg.Wait()
		if !f.d.humanMu.TryLock() {
			t.Fatal("confirm retained the writer lock after responding")
		}
		f.d.humanMu.Unlock()
		if replies[0] == nil || replies[1] == nil {
			t.Fatal("confirmation request did not return")
		}
		if strings.Contains(replies[0].Body.String(), "already-confirmed") {
			replies[0], replies[1] = replies[1], replies[0]
		}
		expect(replies[0], s.ID, "confirmed", "webauthn:fixture-passkey")
		expect(replies[1], s.ID, "already-confirmed", "webauthn:fixture-passkey")
		before = f.confirmationSnapshot(evidenceDB, s.ID)
		expect(post(passkey), s.ID, "already-confirmed", "webauthn:fixture-passkey")
		if f.confirmationSnapshot(evidenceDB, s.ID) != before || before.fixtureCount != 2 {
			t.Fatal("replay changed durable evidence or advanced the counter twice")
		}
		refuse(signed, 409, "Conflict")
		after := f.row(s.ID)
		if after.State != "confirmed" || after.Generation != row.Generation || after.Owner != "" ||
			after.ApprovedBy != "webauthn:fixture-passkey" || f.evidence().operations != before.operations {
			t.Fatalf("handler took a runner step or failed to make the row claimable: %+v", after)
		}
		f.h.SeenAt = f.now.Load()
		f.h.Owner.Generation = 0 // Simulate a fresh runner registering after the POST.
		f.publish()
		f.write(func(ctx context.Context, c *store.Credentials) error {
			cas := store.OperatorCAS{ID: s.ID, Digest: digest, Generation: after.Generation}
			return c.ClaimOperatorAction(ctx, cas, f.h.Owner, f.h.SeenAt)
		})
		if got := f.row(s.ID); got.State != "claimed" || got.Owner != f.h.Owner.Nonce {
			t.Fatalf("a subsequent ready runner could not claim once: %+v", got)
		}
		f.evidence()
		f.now.Add(1)
	})
	f.sub("SSHSIG", func(t *testing.T) {
		second := f.request(opRequest(2))
		r2 := f.row(second.ID)
		body := sshBody(hex.EncodeToString(r2.Digest), signer.Sign(approval.OperatorActionNamespace,
			opMessage(t, [32]byte(r2.Digest), [32]byte(r2.Challenge), r2.ExpiresAt)))
		expect(f.do("POST", second.Href+"/confirm", body, nil), second.ID, "confirmed", "ssh:fixture-ssh")
		if !f.d.humanMu.TryLock() {
			t.Fatal("successful response retained the writer lock")
		}
		f.d.humanMu.Unlock()
		if card := f.do("GET", second.Href, nil, nil); card.Code != 200 || strings.Contains(card.Body.String(), "operator-ceremony") {
			t.Fatalf("confirmed card exposed an active ceremony: %d %s", card.Code, card.Body)
		}
		for _, text := range []string{f.do("GET", second.StatusHref, nil, nil).Body.String(), f.do("GET", "/surface/requests.json", nil, nil).Body.String(), f.log.String()} {
			if strings.Contains(text, "BEGIN SSH SIGNATURE") || strings.Contains(text, "alert(1)") || strings.Contains(text, hex.EncodeToString(m)) || strings.Contains(text, "authenticatorData") {
				t.Fatalf("private task or artifact reached a routine read or log: %s", text)
			}
		}
	})
}

type confirmReadSpy struct {
	io.ReadCloser
	reads int
}

func (b *confirmReadSpy) Read(p []byte) (int, error) {
	b.reads++
	return b.ReadCloser.Read(p)
}

func fixtureConfirmRequest(target string, body io.Reader) *http.Request {
	req := httptest.NewRequest(http.MethodPost, opOrigin+target, body)
	req.Header.Set("Origin", opOrigin)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	return req
}

func requireOperatorCeremony(t *testing.T, card *httptest.ResponseRecorder, s operatorStatus, digest string, message []byte, key *approvaltest.Passkey) {
	t.Helper()
	body := card.Body.String()
	_, after, found := strings.Cut(body, `<pre id="operator-ceremony">`)
	wire, _, closed := strings.Cut(after, "</pre>")
	var got struct {
		Post, Digest, MessageHex, Namespace string
		Passkeys                            []surfacepage.CardPasskeys
	}
	if card.Code != 200 || card.Header().Get("Cache-Control") != "no-store" ||
		card.Header().Get("Content-Security-Policy") == "" || !found || !closed || strings.Contains(body, "<script>alert(1)</script>") ||
		!strings.Contains(body, "&lt;script&gt;alert(1)&lt;/script&gt;") ||
		json.Unmarshal([]byte(html.UnescapeString(wire)), &got) != nil ||
		got.Post != s.Href+"/confirm" || got.Digest != digest || got.MessageHex != hex.EncodeToString(message) ||
		got.Namespace != approval.OperatorActionNamespace || len(got.Passkeys) != 1 ||
		got.Passkeys[0].Origin != opOrigin || got.Passkeys[0].RP != "localhost" ||
		len(got.Passkeys[0].Allow) != 1 || got.Passkeys[0].Allow[0] != base64.RawURLEncoding.EncodeToString(key.CredentialID()) {
		t.Fatalf("card changed the task or key ceremony: %d %s", card.Code, body)
	}
}

func TestOperatorConfirmProjection(t *testing.T) {
	f := newOperatorSurface(t, nil, time.Minute)
	raw := f.postOperator(opRequest(1))
	var s operatorStatus
	if err := json.Unmarshal(raw.Body.Bytes(), &s); err != nil || raw.Code != 200 || s.ID == "" {
		t.Fatalf("projection request fixture: %d %s", raw.Code, raw.Body)
	}
	read := func(path string) *httptest.ResponseRecorder { return f.serve("GET", path, apLoopback, nil) }
	original := f.row(s.ID).Digest
	f.confirmSQL("UPDATE operator_actions SET digest=zeroblob(32) WHERE id=?", s.ID)
	f.refuse(f.evidence(), read(s.Href), 409, "Conflict")
	f.confirmSQL("UPDATE operator_actions SET digest=? WHERE id=?", original, s.ID)
	if r := read(s.Href); r.Code != 200 || !strings.Contains(r.Body.String(), "operator-ceremony") {
		t.Fatalf("unmodified review did not render: %d %s", r.Code, r.Body)
	}

	// Publish the shape of an S6 launch receipt in the migrated fixture. Only
	// the identity read is live. Polling must have no write or child effect.
	executable := "sha256:" + strings.Repeat("a", 64)
	generation := f.row(s.ID).Generation
	f.confirmSQL("UPDATE operator_actions SET action=NULL, state='failed', terminal_at=? WHERE id=?", apT0, s.ID)
	receipt := store.OperatorLaunchReceipt{Run: s.ID, Executable: executable, Status: "spawn-refused"}
	f.confirmSQL(`INSERT INTO operator_action_steps(action_id,kind,identity,phase,version,owner,receipt,created_at)
		VALUES(?, 'launch', ?, 'observed', 2, 'fixture', ?, ?)`, s.ID, s.ID+":"+strings.Repeat("a", 64), opJSON(receipt), apT0)
	status := func(want string) {
		t.Helper()
		before := f.row(s.ID)
		r := read(s.StatusHref)
		var got operatorStatus
		if err := json.Unmarshal(r.Body.Bytes(), &got); err != nil || r.Code != 200 || got.State != want ||
			got.ReportHref != "/v1/episodes/"+s.Episode+"/report" || got.ID != s.ID ||
			!bytes.Equal(opJSON(f.row(s.ID)), opJSON(before)) || before.Generation != generation {
			t.Fatalf("status readback want %s: %d %s", want, r.Code, r.Body)
		}
	}
	status("spawn-refused")
	receipt.PID = os.Getpid()
	birth, live := agentrun.ProcessIdentity(receipt.PID)
	if !live || birth == "" {
		t.Skip("live birth unavailable on this platform")
	}
	f.confirmSQL("UPDATE operator_actions SET state='succeeded' WHERE id=?", s.ID)
	receipt.Birth, receipt.Status = birth, "launched"
	f.confirmSQL("UPDATE operator_action_steps SET receipt=?, version=3 WHERE kind='launch'", opJSON(receipt))
	status("launched")
	receipt.Birth += "-foreign"
	f.confirmSQL("UPDATE operator_action_steps SET receipt=? WHERE kind='launch'", opJSON(receipt))
	status("uncertain")
}

func TestOperatorConfirmationAdmissionNeighbours(t *testing.T) {
	f := newOperatorSurface(t, nil, time.Minute)
	checked := 0
	f.d.approveNow = func() int64 {
		checked++
		if f.d.humanMu.TryLock() {
			f.d.humanMu.Unlock()
			t.Error("request reached its store read outside the human writer lock")
		}
		return apT0
	}
	bad := opRequest(1)
	bad.CatalogueGeneration = 0
	before := f.evidence()
	f.refuse(before, f.postOperator(bad), 400, "BadRequest")
	if r := f.postOperator(opRequest(1)); r.Code != 200 || checked < 2 {
		t.Fatalf("request positive control did not exercise lock: %d (%d checks)", r.Code, checked)
	}
	f.d.approveNow = f.now.Load
	ctx := boundedTestContext(t)
	ro, err := store.OpenCredentialsReadOnly(f.d.credentialsPath)
	if err != nil {
		t.Fatal(err)
	}
	defer ro.Close()
	if _, _, err = f.d.operatorReview(ctx, ro, "sol", apT0); err != nil {
		t.Fatalf("catalogue control: %v", err)
	}
	f.confirmSQL("UPDATE operator_runner_state SET catalogue_generation=0")
	if h, _, err := f.d.operatorReview(ctx, ro, "sol", apT0); !errors.Is(err, store.ErrOperatorRunner) || h.CatalogueGeneration != 0 {
		t.Fatalf("invalid generation crossed review: %v, %v", h, err)
	}
	bar := "/surface/verb.json?q=" + url.QueryEscape("episode start ep-2 --preset sol --launch --headless")
	if r := f.serve("GET", bar, apLoopback, nil); r.Code != 200 || !strings.Contains(r.Body.String(), `"outcome":"terminal"`) {
		t.Fatalf("invalid generation reached the bar review: %d %s", r.Code, r.Body)
	}
	summary := func() surfacepage.HomeSummary {
		t.Helper()
		r := f.serve("GET", "/surface/home-summary.json", apLoopback, nil)
		var out surfacepage.HomeSummary
		if r.Code != 200 || json.Unmarshal(r.Body.Bytes(), &out) != nil || len(out.Items) != 1 || out.Items[0].Label != "Needs you" {
			t.Fatalf("home summary: %d %s", r.Code, r.Body)
		}
		return out
	}
	if got := summary().Items[0].Value; got != "1" {
		t.Fatalf("healthy pending read did not show one row: %s", got)
	}
	f.confirmSQL("UPDATE operator_actions SET expires_at='unreadable' WHERE episode_id='ep-1'")
	before = f.evidence()
	r := f.serve("GET", "/surface/requests.json", apLoopback, nil)
	var requests surfacepage.RequestsJSON
	if r.Code != 200 || json.Unmarshal(r.Body.Bytes(), &requests) != nil || requests.Unavailable != "" || requests.OperatorActionsUnavailable == "" {
		t.Fatalf("fixture must fail only the operator pending read: %d %s", r.Code, r.Body)
	}
	if got := summary().Items[0].Value; got != "unknown" || before != f.evidence() {
		t.Fatalf("unavailable operator read looked like an empty queue or changed evidence: %s", got)
	}
}
