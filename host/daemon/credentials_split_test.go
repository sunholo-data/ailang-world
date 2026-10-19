package daemon

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/store"
)

// preSplitStore makes a schema-4 world store at path holding one credential
// for token: a fresh schema-5 store plus the v3/v4 session_credentials table,
// relabelled 4 — exactly the shape of every store an older build served.
func preSplitStore(t *testing.T, path, token, episode string) {
	t.Helper()
	s, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	_ = s.Close()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	grants, _ := json.Marshal(grantForAuth())
	sum := sha256.Sum256([]byte(token))
	for _, stmt := range []string{
		`CREATE TABLE session_credentials (
    credential_id TEXT PRIMARY KEY,
    episode_id   TEXT NOT NULL CHECK (episode_id <> ''),
    grants_json  TEXT NOT NULL,
    expires_at   INTEGER NOT NULL,
    created_at   INTEGER NOT NULL
)`,
		`PRAGMA user_version = 4`,
	} {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := db.Exec(`INSERT INTO session_credentials VALUES (?, ?, ?, ?, ?)`,
		hex.EncodeToString(sum[:]), episode, string(grants), time.Now().Unix()+3600, time.Now().Unix()); err != nil {
		t.Fatal(err)
	}
}

// TestServeMovesPreSplitCredentialsAndTheyStillAuthenticate: serve's startup
// moves a schema-4 store's credential out (logging it once), the world store
// keeps none, and the moved token authenticates a request.
func TestServeMovesPreSplitCredentialsAndTheyStillAuthenticate(t *testing.T) {
	dir := t.TempDir()
	world := filepath.Join(dir, "world.db")
	token := strings.Repeat("ab", 32)
	preSplitStore(t, world, token, "ep-moved")
	var log strings.Builder
	d, err := New(boundedTestContext(t), Config{DBPath: world, BindHost: DefaultBindHost, ErrorLog: &log})
	if err != nil {
		t.Fatalf("New on a schema-4 store: %v", err)
	}
	t.Cleanup(func() { _ = d.Close() })
	if !strings.Contains(log.String(), "moved 1 session credential(s)") {
		t.Fatalf("startup log = %q, want the one-time move reported", log.String())
	}
	if want := filepath.Join(dir, store.CredentialsFileName); !strings.HasSuffix(d.credentialsPath, store.CredentialsFileName) || filepath.Base(filepath.Dir(d.credentialsPath)) != filepath.Base(filepath.Dir(want)) {
		t.Fatalf("credentials path = %q, want it beside the world store (%s)", d.credentialsPath, want)
	}
	if v, err := store.WorldSchemaVersion(world); err != nil || v != 5 {
		t.Fatalf("world schema after serve = %d %v, want 5", v, err)
	}
	rec := postCommitAuth(t, d, "Bearer "+token)
	if rec.Code == http.StatusUnauthorized {
		t.Fatalf("the moved credential was refused: %s", rec.Body)
	}
}

// TestMintAndRevokeWhileServing is row 189's daemon-level proof: with the
// daemon holding the world writer lock, a credential minted into the
// credentials store authenticates the next request and a revoke refuses the
// one after — no restart. The world lock is still exclusive throughout.
func TestMintAndRevokeWhileServing(t *testing.T) {
	d := newHandlerDaemon(t)
	if s, err := store.Open(d.cfg.DBPath); err == nil {
		_ = s.Close()
		t.Fatal("a second world writer opened while the daemon serves")
	} else if !store.IsWriterAlreadyActive(err) {
		t.Fatalf("second world writer: %v", err)
	}
	creds := testCreds(t, d)
	tok, id, _, err := authority.Mint(boundedTestContext(t), creds, "ep-while-serving", grantForAuth(), 3600, time.Now().Unix(), nil)
	if err != nil {
		t.Fatalf("mint while serving: %v", err)
	}
	if rec := postCommitAuth(t, d, "Bearer "+tok); rec.Code != http.StatusBadRequest {
		t.Fatalf("fresh mint: status %d, want 400 from handleCommit (past the middleware); body=%s", rec.Code, rec.Body)
	}
	if err := authority.Revoke(boundedTestContext(t), creds, id); err != nil {
		t.Fatalf("revoke while serving: %v", err)
	}
	rec := postCommitAuth(t, d, "Bearer "+tok)
	if rec.Code != http.StatusUnauthorized || errorClass(t, rec.Body.Bytes()) != "SessionUnknown" {
		t.Fatalf("after revoke: status %d class %q, want 401 SessionUnknown", rec.Code, errorClass(t, rec.Body.Bytes()))
	}
}

func TestServeRefusesTheWorldStoreAsCredentialsStore(t *testing.T) {
	world := filepath.Join(t.TempDir(), "world.db")
	_, err := New(boundedTestContext(t), Config{DBPath: world, CredentialsPath: world, BindHost: DefaultBindHost})
	var se *StartupError
	if !errors.As(err, &se) || se.Stage != StageConfig {
		t.Fatalf("New with --credentials-db = --db: %v, want a config refusal", err)
	}
}
