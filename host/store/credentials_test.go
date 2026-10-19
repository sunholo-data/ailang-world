package store

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

// canonicalCredentialsV1DDL is review-visible policy, hardcoded like
// canonicalTableDDL: the version-1 credentials table, byte-identical in
// columns to the world store's v3/v4 table so the move copies it unchanged.
var canonicalCredentialsV1DDL = `CREATE TABLE session_credentials (
    credential_id TEXT PRIMARY KEY,
    episode_id   TEXT NOT NULL CHECK (episode_id <> ''),
    grants_json  TEXT NOT NULL,
    expires_at   INTEGER NOT NULL,
    created_at   INTEGER NOT NULL
)`

// canonicalCredentialsDDL is the current (version 2, row 189 M2) table: the
// version-1 columns plus the nullable profile binding.
var canonicalCredentialsDDL = `CREATE TABLE session_credentials (
    credential_id TEXT PRIMARY KEY,
    episode_id   TEXT NOT NULL CHECK (episode_id <> ''),
    grants_json  TEXT NOT NULL,
    expires_at   INTEGER NOT NULL,
    created_at   INTEGER NOT NULL,
    project TEXT CHECK (project IS NULL OR project <> ''),
    profile_digest TEXT CHECK ((profile_digest IS NULL) = (project IS NULL) AND (profile_digest IS NULL OR profile_digest <> ''))
)`

func TestCredentialsSchemaMatchesCanonicalDDL(t *testing.T) {
	c := openTestCredentials(t)
	ddl := tableDDL(t, c.db)
	requireExactTableNames(t, "fresh credentials store", ddl, []string{"session_credentials"})
	if got, want := normalizeDDL(ddl["session_credentials"]), normalizeDDL(canonicalCredentialsDDL); got != want {
		t.Fatalf("credentials DDL:\n got: %s\nwant: %s", got, want)
	}
	// The pre-split world table is the version-1 table, column for column.
	v4 := rawDB(t, filepath.Join(t.TempDir(), "v4.db"))
	if _, err := v4.Exec(schemaV4SQL); err != nil {
		t.Fatal(err)
	}
	if got := normalizeDDL(tableDDL(t, v4)["session_credentials"]); got != normalizeDDL(canonicalCredentialsV1DDL) {
		t.Fatalf("v4 world table differs from the version-1 credentials table: %s", got)
	}
	var version int
	if err := c.db.QueryRow("PRAGMA user_version").Scan(&version); err != nil || version != credentialsSchemaVersion {
		t.Fatalf("credentials user_version = %d %v, want %d", version, err, credentialsSchemaVersion)
	}
}

// openTestCredentials opens a file-backed credentials store in a temp dir.
func openTestCredentials(t testing.TB) *Credentials {
	t.Helper()
	c, err := OpenCredentials(filepath.Join(t.TempDir(), CredentialsFileName))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Close() })
	return c
}

// TestCredentialsMintAndRevokeBesideALiveWorldWriter is row 189's store-level
// lock proof (D-WORLD-85 C): while a SECOND OS PROCESS holds the WORLD store's
// writer lock (what `serve` does), the credentials store's own writer lock is
// free, so a mint and a revoke both land — and the world lock still refuses a
// second world writer. Mutant: give the credentials store the world store's
// lock file and OpenCredentials below returns *WriterAlreadyActive.
func TestCredentialsMintAndRevokeBesideALiveWorldWriter(t *testing.T) {
	world := filepath.Join(t.TempDir(), "world.db")
	h := startWriterHelper(t, world)

	// The world writer lock still protects everything else.
	if s, err := Open(world); err == nil {
		_ = s.Close()
		t.Fatal("a second world writer opened beside the live one")
	} else if !IsWriterAlreadyActive(err) {
		t.Fatalf("second world writer: %v, want WriterAlreadyActive", err)
	}

	credPath, err := DefaultCredentialsPath(world)
	if err != nil {
		t.Fatal(err)
	}
	c, err := OpenCredentials(credPath)
	if err != nil {
		t.Fatalf("credentials writer beside a live world writer: %v", err)
	}
	row := SessionRow{CredentialID: "beside-live", EpisodeID: "ep", GrantsJSON: "[]", ExpiresAt: 100, CreatedAt: 1}
	if err := c.MintSession(boundedTestContext(t), row); err != nil {
		t.Fatalf("mint beside a live world writer: %v", err)
	}
	// The credentials lock is real: a second credentials writer is refused.
	if c2, err := OpenCredentials(credPath); err == nil {
		_ = c2.Close()
		t.Fatal("a second credentials writer opened beside the first")
	} else if !IsWriterAlreadyActive(err) {
		t.Fatalf("second credentials writer: %v, want WriterAlreadyActive", err)
	}
	ro, err := OpenCredentialsReadOnly(credPath)
	if err != nil {
		t.Fatal(err)
	}
	got, ok, err := ro.ResolveSession(boundedTestContext(t), row.CredentialID)
	if err != nil || !ok || got != row {
		t.Fatalf("read-only resolve = %+v %v %v, want the minted row", got, ok, err)
	}
	if err := c.RevokeSession(boundedTestContext(t), row.CredentialID); err != nil {
		t.Fatalf("revoke beside a live world writer: %v", err)
	}
	if _, ok, err := ro.ResolveSession(boundedTestContext(t), row.CredentialID); err != nil || ok {
		t.Fatalf("resolve after revoke = %v %v, want unknown", ok, err)
	}
	_ = ro.Close()
	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
	h.stop(t)
}

func TestDefaultCredentialsPathIsBesideTheWorldStore(t *testing.T) {
	dir := t.TempDir()
	got, err := DefaultCredentialsPath(filepath.Join(dir, "world.db"))
	if err != nil {
		t.Fatal(err)
	}
	canonicalDir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(canonicalDir, "credentials.db"); got != want {
		t.Fatalf("DefaultCredentialsPath = %q, want %q", got, want)
	}
	dsn, err := DefaultCredentialsPath("file:" + filepath.Join(dir, "world.db") + "?_pragma=busy_timeout(100)")
	if err != nil || dsn != got {
		t.Fatalf("DSN form = %q %v, want %q", dsn, err, got)
	}
	if _, err := DefaultCredentialsPath(":memory:"); err == nil {
		t.Fatal("an in-memory world store got a default credentials path")
	}
}

// A world store passed as the credentials store is refused unmodified.
func TestOpenCredentialsRefusesAWorldStore(t *testing.T) {
	world := filepath.Join(t.TempDir(), "world.db")
	s, err := Open(world)
	if err != nil {
		t.Fatal(err)
	}
	_ = s.Close()
	db := rawDB(t, world)
	before, beforeVersion := schemaState(t, db)
	_ = db.Close()
	for _, open := range []func(string) (*Credentials, error){OpenCredentials, OpenCredentialsReadOnly} {
		c, err := open(world)
		if c != nil {
			_ = c.Close()
		}
		var schema *CredentialsSchemaError
		if !errors.As(err, &schema) {
			t.Fatalf("open world store as credentials: %v, want CredentialsSchemaError", err)
		}
	}
	after, afterVersion := schemaState(t, rawDB(t, world))
	if !reflect.DeepEqual(after, before) || afterVersion != beforeVersion {
		t.Fatalf("refused open changed the world store: %v/%d -> %v/%d", before, beforeVersion, after, afterVersion)
	}
}

func TestOpenCredentialsReadOnlyNeverCreates(t *testing.T) {
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	if c, err := OpenCredentialsReadOnly(path); err == nil {
		_ = c.Close()
		t.Fatal("read-only open of a missing credentials store succeeded")
	}
	for _, p := range []string{path, path + writerLockSuffix} {
		if _, err := os.Stat(p); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("read-only open created %s: %v", p, err)
		}
	}
}

// A credentials store whose table was dropped is refused, never recreated
// empty (that would silently revoke everything).
func TestDroppedCredentialsTableIsRefused(t *testing.T) {
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	c, err := OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	_ = c.Close()
	db := rawDB(t, path)
	if _, err := db.Exec(`DROP TABLE session_credentials`); err != nil {
		t.Fatal(err)
	}
	_ = db.Close()
	if c, err := OpenCredentials(path); err == nil {
		_ = c.Close()
		t.Fatal("a credentials store missing its table was reopened")
	}
	if names, _ := schemaState(t, rawDB(t, path)); len(names) != 0 {
		t.Fatalf("refused open recreated %v", names)
	}
}

func TestCredentialsMethodsRequireADeadline(t *testing.T) {
	c := openTestCredentials(t)
	ctx := context.Background()
	row := SessionRow{CredentialID: "x", EpisodeID: "ep", GrantsJSON: "[]"}
	errs := map[string]error{
		"MintSession":   c.MintSession(ctx, row),
		"RevokeSession": c.RevokeSession(ctx, "x"),
	}
	_, _, errs["ResolveSession"] = c.ResolveSession(ctx, "x")
	_, errs["ListSessions"] = c.ListSessions(ctx, "", 1)
	for name, err := range errs {
		if !errors.Is(err, ErrNoDeadline) {
			t.Errorf("%s(Background) = %v, want ErrNoDeadline", name, err)
		}
	}
	bounded, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := c.MintSession(bounded, row); err != nil {
		t.Fatalf("control: bounded mint: %v", err)
	}
}
