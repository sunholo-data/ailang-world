package daemon

// Row 201 A3 sprint 2: the approval routes' test rig. Every store is under t.TempDir(); every key
// is approvaltest's (D-WORLD-109); no listener: each request carries the local address
// 127.0.0.1:7644 in its context, as net/http sets it for a connection.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// apT0 is the approval routes' pinned clock; apLoopback the Host of a request on the loopback port.
const apT0, apLoopback = int64(1800000000), "127.0.0.1:7644"

type approvalRig struct {
	t *testing.T
	d *Daemon
}

func newApprovalRig(t *testing.T) approvalRig {
	t.Helper()
	d, err := New(boundedTestContext(t), Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost, ErrorLog: io.Discard})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = d.Close() })
	d.approveNow = func() int64 { return apT0 }
	return approvalRig{t: t, d: d}
}

// write runs f on the credentials writer, opened and closed around it.
func (a approvalRig) write(f func(ctx context.Context, c *store.Credentials) error) {
	a.t.Helper()
	c, err := store.OpenCredentials(a.d.credentialsPath)
	if err != nil {
		a.t.Fatal(err)
	}
	defer func() { _ = c.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := f(ctx, c); err != nil {
		a.t.Fatal(err)
	}
}

// putRequest stores a pending, canonical request id for ep, created at created, expiring 600 s
// later; its challenge is the SHA-256 of its id.
func (a approvalRig) putRequest(id, ep string, created int64) store.MintRequestRow {
	a.t.Helper()
	b, err := approval.Canonical(approval.MintRequest{Schema: approval.RequestSchema, Episode: ep, WorldDB: "/ws/world.db", WorkspaceRoot: "/ws/root",
		Worktree: "/ws/root/" + ep, Git: "/usr/bin/git", Branch: "world/" + ep, Grants: []approval.Grant{{Effect: "Workspace.Read", Scope: ep}, {Effect: "Workspace.Write", Scope: ep, Budget: 5}},
		TTLSeconds: 3600, ProfilesDir: "/ws/profiles", CredentialsDB: "/ws/credentials.db", SessionFile: "/ws/sessions/" + ep + "/session.token",
		MCPConfig: "/ws/sessions/" + ep + "/mcp.json", Agent: "claude", Model: "claude-opus-5-5", Task: "/ws/task.md", Daemon: approval.DaemonRef{Addr: apLoopback}})
	if err != nil {
		a.t.Fatal(err)
	}
	d, c := sha256.Sum256(b), sha256.Sum256([]byte(id))
	row := store.MintRequestRow{RequestID: id, Request: b, Digest: hex.EncodeToString(d[:]), Challenge: hex.EncodeToString(c[:]), EpisodeID: ep, CreatedAt: created, ExpiresAt: created + 600}
	a.write(func(ctx context.Context, c *store.Credentials) error { return c.PutMintRequest(ctx, row, created) })
	return row
}

// consume consumes row by the typed y at `at`, minting session credential cred.
func (a approvalRig) consume(row store.MintRequestRow, at int64, cred string) {
	a.t.Helper()
	a.write(func(ctx context.Context, c *store.Credentials) error {
		return c.MintSessionForRequest(ctx, store.SessionRow{CredentialID: cred, EpisodeID: row.EpisodeID, GrantsJSON: "[]", ExpiresAt: at + 3600, CreatedAt: at},
			store.RequestConsume{ID: row.RequestID, Digest: row.Digest, Challenge: row.Challenge, ApprovedBy: "tty",
				Clock: func() int64 { return at }, Check: func(int64) error { return nil }})
	})
}

// enrolPasskey enrols pk as name at origin (rpID: the origin's host) at apT0-10: after the tests'
// requests (apT0-30), so a judge at the request's creation instead of now refuses.
func (a approvalRig) enrolPasskey(name string, pk *approvaltest.Passkey, origin string) {
	a.t.Helper()
	rp, err := approval.ParseCeremonyOrigin(origin)
	if err != nil {
		a.t.Fatal(err)
	}
	k := pk.Key()
	a.write(func(ctx context.Context, c *store.Credentials) error {
		return c.EnrolPasskey(ctx, store.ApproverRow{Name: name, Format: approval.FormatWebAuthn, PublicKey: k.Canonical(), EnrolledAt: apT0 - 10},
			store.PasskeyRow{Name: name, RPID: rp, Origin: origin, CredentialID: pk.CredentialID(), AAGUID: make([]byte, 16), Alg: k.Alg, SignCount: 1})
	})
}

// serve sends one request through the daemon's handler as if it arrived on 127.0.0.1:7644, with
// Host host and the header pairs hdr.
func (a approvalRig) serve(method, target, host string, body []byte, hdr ...string) *httptest.ResponseRecorder {
	a.t.Helper()
	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), http.LocalAddrContextKey, net.Addr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 7644})))
	req.Host = host
	for i := 0; i+1 < len(hdr); i += 2 {
		req.Header.Set(hdr[i], hdr[i+1])
	}
	rec := httptest.NewRecorder()
	a.d.Handler().ServeHTTP(rec, req)
	return rec
}

// approvalHeaders are an approval route's headers: JSON, no-store, nosniff, and no CORS header.
func approvalHeaders(rec *httptest.ResponseRecorder) error {
	for k, v := range map[string]string{"Content-Type": "application/json", "Cache-Control": "no-store", "X-Content-Type-Options": "nosniff"} {
		if got := rec.Header().Get(k); got != v {
			return fmt.Errorf("%s = %q, want %q", k, got, v)
		}
	}
	for k := range rec.Header() {
		if strings.HasPrefix(strings.ToLower(k), "access-control-") {
			return fmt.Errorf("a CORS header %s", k)
		}
	}
	return nil
}

// requireRefusal fails unless rec is an inert refusal with status and class.
func requireRefusal(t *testing.T, what string, rec *httptest.ResponseRecorder, status int, class string) surfacepage.ApproveRefusal {
	t.Helper()
	var r surfacepage.ApproveRefusal
	if err := json.Unmarshal(rec.Body.Bytes(), &r); err != nil || rec.Code != status || r.V != 1 || r.Error != class || r.Reason == "" {
		t.Fatalf("%s: %d %s, want %d %s", what, rec.Code, rec.Body.String(), status, class)
	}
	if err := approvalHeaders(rec); err != nil {
		t.Fatalf("%s: %v", what, err)
	}
	return r
}

// setRequestColumn sets one column of request id straight in the credentials file (a hand edit,
// or a row no store method writes).
func setRequestColumn(t *testing.T, path, id, column string, v any) {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	if res, err := db.Exec(`UPDATE mint_requests SET `+column+` = ? WHERE request_id = ?`, v, id); err != nil {
		t.Fatal(err)
	} else if n, _ := res.RowsAffected(); n != 1 {
		t.Fatalf("%d rows updated", n)
	}
}

// setVersion relabels the credentials store at path to an older schema by hand (DROP TABLE, PRAGMA).
func setVersion(t *testing.T, path string, statements ...string) {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	for _, s := range statements {
		if _, err := db.Exec(s); err != nil {
			t.Fatal(err)
		}
	}
}
