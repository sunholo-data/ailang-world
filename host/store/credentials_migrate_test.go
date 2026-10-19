package store

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

// rigRows is eight credentials, the count on the operator's rig store.
func rigRows() []SessionRow {
	rows := make([]SessionRow, 0, 8)
	for i := 0; i < 8; i++ {
		rows = append(rows, SessionRow{
			CredentialID: fmt.Sprintf("%064x", i+1),
			EpisodeID:    fmt.Sprintf("w%d", 180+i),
			GrantsJSON:   fmt.Sprintf(`[{"Effect":"Workspace.Read","Scope":"worktree","ExpiresAt":%d,"Budget":%d}]`, 2000+i, 50+i),
			ExpiresAt:    int64(2000 + i),
			CreatedAt:    int64(1000 + i),
		})
	}
	return rows
}

// makeV4Store writes a schema-4 world store at path holding rows and one
// world object, the shape every pre-split store on disk has.
func makeV4Store(t *testing.T, path string, rows []SessionRow) Object {
	t.Helper()
	db := rawDB(t, path)
	if _, err := db.Exec(schemaV4SQL); err != nil {
		t.Fatal(err)
	}
	setPositiveFixtureVersion(t, db, 4)
	for _, r := range rows {
		if _, err := db.Exec(`INSERT INTO session_credentials (credential_id, episode_id, grants_json, expires_at, created_at) VALUES (?, ?, ?, ?, ?)`,
			r.CredentialID, r.EpisodeID, r.GrantsJSON, r.ExpiresAt, r.CreatedAt); err != nil {
			t.Fatal(err)
		}
	}
	o := obj("pre-split world object", "t/pre-split")
	if _, err := db.Exec(`INSERT INTO objects (hash_ref, interface_hash_ref, semantic_id, provenance, payload) VALUES (?, ?, ?, ?, ?)`,
		o.Hash.String(), o.InterfaceHash.String(), o.SemanticID, o.Provenance, o.Payload); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	return o
}

// worldCredentialRows counts session_credentials rows left in a world store,
// or -1 when the table is gone.
func worldCredentialRows(t *testing.T, path string) (version, rows int) {
	t.Helper()
	db := rawDB(t, path)
	defer func() { _ = db.Close() }()
	names, version := schemaState(t, db)
	for _, n := range names {
		if n == "session_credentials" {
			if err := db.QueryRow(`SELECT count(*) FROM session_credentials`).Scan(&rows); err != nil {
				t.Fatal(err)
			}
			return version, rows
		}
	}
	return version, -1
}

func credentialRows(t *testing.T, credPath string) []SessionRow {
	t.Helper()
	db := rawDB(t, credPath)
	defer func() { _ = db.Close() }()
	got, err := listSessions(boundedTestContext(t), db, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func migratePaths(t *testing.T) (world, creds string) {
	dir := t.TempDir()
	return filepath.Join(dir, "world.db"), filepath.Join(dir, CredentialsFileName)
}

// TestMigrateCredentialsMovesEveryRowAndLeavesNone: all eight rows arrive
// field for field (ids, episodes, grants with budgets, expiry, creation), the
// world store is schema 5 with no session_credentials table at all, its world
// data is untouched, and it opens normally afterwards.
func TestMigrateCredentialsMovesEveryRowAndLeavesNone(t *testing.T) {
	world, creds := migratePaths(t)
	want := rigRows()
	o := makeV4Store(t, world, want)
	moved, err := MigrateCredentials(boundedTestContext(t), world, creds)
	if err != nil || moved != len(want) {
		t.Fatalf("MigrateCredentials = %d, %v; want %d, nil", moved, err, len(want))
	}
	if got := credentialRows(t, creds); !reflect.DeepEqual(got, want) {
		t.Fatalf("credentials store rows =\n%+v\nwant\n%+v", got, want)
	}
	if version, rows := worldCredentialRows(t, world); version != currentSchemaVersion || rows != -1 {
		t.Fatalf("world store after the move: version %d, session_credentials rows %d; want %d and no table", version, rows, currentSchemaVersion)
	}
	s, err := Open(world)
	if err != nil {
		t.Fatalf("Open after the move: %v", err)
	}
	defer func() { _ = s.Close() }()
	if got, ok, err := s.GetObject(boundedTestContext(t), o.Hash); err != nil || !ok || !reflect.DeepEqual(got, o) {
		t.Fatalf("world object after the move = %+v %v %v", got, ok, err)
	}
}

func TestMigrateCredentialsIsIdempotent(t *testing.T) {
	world, creds := migratePaths(t)
	want := rigRows()
	makeV4Store(t, world, want)
	if _, err := MigrateCredentials(boundedTestContext(t), world, creds); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		moved, err := MigrateCredentials(boundedTestContext(t), world, creds)
		if err != nil || moved != 0 {
			t.Fatalf("re-run %d = %d, %v; want 0, nil", i, moved, err)
		}
	}
	if got := credentialRows(t, creds); !reflect.DeepEqual(got, want) {
		t.Fatalf("rows after re-runs = %+v", got)
	}
}

// A fresh world store is born schema 5 without the table, and migrating it
// (or a missing one) is a no-op that creates no credentials file.
func TestMigrateCredentialsFreshStore(t *testing.T) {
	world, creds := migratePaths(t)
	if moved, err := MigrateCredentials(boundedTestContext(t), world, creds); err != nil || moved != 0 {
		t.Fatalf("missing world store: %d, %v", moved, err)
	}
	s, err := Open(world)
	if err != nil {
		t.Fatal(err)
	}
	_ = s.Close()
	if version, rows := worldCredentialRows(t, world); version != currentSchemaVersion || rows != -1 {
		t.Fatalf("fresh store: version %d rows %d; want %d and no table", version, rows, currentSchemaVersion)
	}
	if moved, err := MigrateCredentials(boundedTestContext(t), world, creds); err != nil || moved != 0 {
		t.Fatalf("fresh world store: %d, %v", moved, err)
	}
	if _, err := os.Stat(creds); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("a no-op move created %s: %v", creds, err)
	}
}

func TestMigrateCredentialsZeroCredentials(t *testing.T) {
	world, creds := migratePaths(t)
	makeV4Store(t, world, nil)
	moved, err := MigrateCredentials(boundedTestContext(t), world, creds)
	if err != nil || moved != 0 {
		t.Fatalf("empty move = %d, %v", moved, err)
	}
	if version, rows := worldCredentialRows(t, world); version != currentSchemaVersion || rows != -1 {
		t.Fatalf("world after empty move: version %d rows %d", version, rows)
	}
	if got := credentialRows(t, creds); len(got) != 0 {
		t.Fatalf("credentials store after empty move = %+v", got)
	}
}

// assertUntouchedV4 checks an interrupted move left the world store exactly
// pre-split (schema 4, every row) and the credentials store without copies.
func assertUntouchedV4(t *testing.T, world, creds string, want []SessionRow) {
	t.Helper()
	if version, rows := worldCredentialRows(t, world); version != preSplitSchemaVersion || rows != len(want) {
		t.Fatalf("world after an interrupted move: version %d rows %d; want 4 and %d", version, rows, len(want))
	}
	if _, err := os.Stat(creds); err == nil {
		if got := credentialRows(t, creds); len(got) != 0 {
			t.Fatalf("an interrupted move left %d copies in the credentials store", len(got))
		}
	}
}

// TestMigrateCredentialsInterruptedRollsBack fails the move inside its one
// transaction, after the copy and after the drop: neither leaves a lost or a
// duplicated credential, and the re-run moves each row exactly once. Mutant:
// commit the copy in its own transaction and the "copied" arm finds copies.
func TestMigrateCredentialsInterruptedRollsBack(t *testing.T) {
	for _, stage := range []string{"copied", "dropped"} {
		t.Run(stage, func(t *testing.T) {
			world, creds := migratePaths(t)
			want := rigRows()
			makeV4Store(t, world, want)
			induced := errors.New("induced interruption at " + stage)
			prev := migrateCredentialsHook
			migrateCredentialsHook = func(s string) error {
				if s == stage {
					return induced
				}
				return nil
			}
			_, err := MigrateCredentials(boundedTestContext(t), world, creds)
			migrateCredentialsHook = prev
			if !errors.Is(err, induced) {
				t.Fatalf("interrupted move = %v, want the induced error", err)
			}
			assertUntouchedV4(t, world, creds, want)
			moved, err := MigrateCredentials(boundedTestContext(t), world, creds)
			if err != nil || moved != len(want) {
				t.Fatalf("re-run = %d, %v", moved, err)
			}
			if got := credentialRows(t, creds); !reflect.DeepEqual(got, want) {
				t.Fatalf("rows after re-run = %+v", got)
			}
		})
	}
}

const (
	migrateCrashWorldEnv = "WORLDD_MIGRATE_CRASH_WORLD"
	migrateCrashCredsEnv = "WORLDD_MIGRATE_CRASH_CREDS"
	migrateCrashStageEnv = "WORLDD_MIGRATE_CRASH_STAGE"
	migrateCrashExit     = 9
)

// TestMigrateCredentialsCrashHelper is the body of the subprocess the crash
// test kills mid-move: it exits the process inside the move's transaction.
func TestMigrateCredentialsCrashHelper(t *testing.T) {
	world := os.Getenv(migrateCrashWorldEnv)
	if world == "" {
		t.Skip("subprocess helper body for TestMigrateCredentialsSurvivesProcessDeath")
	}
	stage := os.Getenv(migrateCrashStageEnv)
	migrateCredentialsHook = func(s string) error {
		if s == stage {
			os.Exit(migrateCrashExit)
		}
		return nil
	}
	_, err := MigrateCredentials(boundedTestContext(t), world, os.Getenv(migrateCrashCredsEnv))
	fmt.Fprintf(os.Stderr, "helper: the move finished without reaching %s: %v\n", stage, err)
	os.Exit(1)
}

// TestMigrateCredentialsSurvivesProcessDeath kills a real process inside the
// move (no deferred cleanup, hot journals left on disk): the next open rolls
// both files back, nothing is lost or duplicated, and the re-run completes.
func TestMigrateCredentialsSurvivesProcessDeath(t *testing.T) {
	for _, stage := range []string{"copied", "dropped"} {
		t.Run(stage, func(t *testing.T) {
			world, creds := migratePaths(t)
			want := rigRows()
			makeV4Store(t, world, want)
			cmd := exec.Command(os.Args[0], "-test.run=^TestMigrateCredentialsCrashHelper$")
			cmd.Env = append(os.Environ(), migrateCrashWorldEnv+"="+world, migrateCrashCredsEnv+"="+creds, migrateCrashStageEnv+"="+stage)
			out, err := cmd.CombinedOutput()
			var exit *exec.ExitError
			if !errors.As(err, &exit) || exit.ExitCode() != migrateCrashExit {
				t.Fatalf("helper did not die inside the move (err %v); output:\n%s", err, out)
			}
			assertUntouchedV4(t, world, creds, want)
			moved, err := MigrateCredentials(boundedTestContext(t), world, creds)
			if err != nil || moved != len(want) {
				t.Fatalf("re-run after the crash = %d, %v", moved, err)
			}
			if got := credentialRows(t, creds); !reflect.DeepEqual(got, want) {
				t.Fatalf("rows after re-run = %+v", got)
			}
			if version, rows := worldCredentialRows(t, world); version != currentSchemaVersion || rows != -1 {
				t.Fatalf("world after re-run: version %d rows %d", version, rows)
			}
		})
	}
}

// A different row already under one of the ids refuses the whole move.
func TestMigrateCredentialsConflictRefusesEverything(t *testing.T) {
	world, creds := migratePaths(t)
	want := rigRows()
	makeV4Store(t, world, want)
	c, err := OpenCredentials(creds)
	if err != nil {
		t.Fatal(err)
	}
	clash := want[3]
	clash.GrantsJSON = `[{"Effect":"Workspace.Write","Scope":"worktree","Budget":999}]`
	if err := c.MintSession(boundedTestContext(t), clash); err != nil {
		t.Fatal(err)
	}
	_ = c.Close()
	if _, err := MigrateCredentials(boundedTestContext(t), world, creds); err == nil {
		t.Fatal("a conflicting move succeeded")
	}
	if version, rows := worldCredentialRows(t, world); version != preSplitSchemaVersion || rows != len(want) {
		t.Fatalf("world after a refused move: version %d rows %d", version, rows)
	}
	if got := credentialRows(t, creds); len(got) != 1 || got[0] != clash {
		t.Fatalf("credentials after a refused move = %+v", got)
	}
}

// A held world writer lock (a daemon of an older build) is a named refusal,
// and the store is untouched until the lock is free.
func TestMigrateCredentialsRefusesAHeldWorldLock(t *testing.T) {
	world, creds := migratePaths(t)
	want := rigRows()
	makeV4Store(t, world, want)
	canonical, _, err := resolveDSN(world)
	if err != nil {
		t.Fatal(err)
	}
	lock, err := acquireWriterLock(canonical)
	if err != nil {
		t.Fatal(err)
	}
	_, err = MigrateCredentials(boundedTestContext(t), world, creds)
	var notMoved *CredentialsNotMovedError
	if !errors.As(err, &notMoved) || !IsWriterAlreadyActive(err) {
		t.Fatalf("move under a held lock = %v, want CredentialsNotMovedError wrapping WriterAlreadyActive", err)
	}
	assertUntouchedV4(t, world, creds, want)
	_ = lock.release()
	if moved, err := MigrateCredentials(boundedTestContext(t), world, creds); err != nil || moved != len(want) {
		t.Fatalf("move after release = %d, %v", moved, err)
	}
}

func TestMigrateCredentialsRefusesTheWorldStoreAsTarget(t *testing.T) {
	world, _ := migratePaths(t)
	makeV4Store(t, world, rigRows())
	if _, err := MigrateCredentials(boundedTestContext(t), world, world); err == nil {
		t.Fatal("moving credentials into the world store itself succeeded")
	}
	if version, rows := worldCredentialRows(t, world); version != preSplitSchemaVersion || rows != 8 {
		t.Fatalf("world after refusal: version %d rows %d", version, rows)
	}
}

// A schema-4 store opens read-only (world tables unchanged) and its rows list
// through LegacySessions; a writable Open is refused with the named fix and
// changes nothing.
func TestPreSplitStoreOpenModes(t *testing.T) {
	world, _ := migratePaths(t)
	want := rigRows()
	makeV4Store(t, world, want)
	if s, err := Open(world); err == nil {
		_ = s.Close()
		t.Fatal("writable Open of a schema-4 store succeeded")
	} else {
		var pre *PreSplitStoreError
		if !errors.As(err, &pre) {
			t.Fatalf("writable Open = %v, want PreSplitStoreError", err)
		}
	}
	ro, err := OpenReadOnly(world)
	if err != nil {
		t.Fatalf("read-only Open of a schema-4 store: %v", err)
	}
	_ = ro.Close()
	got, err := LegacySessions(boundedTestContext(t), world, "", 0)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("LegacySessions = %+v, %v", got, err)
	}
	if v, err := WorldSchemaVersion(world); err != nil || v != preSplitSchemaVersion {
		t.Fatalf("WorldSchemaVersion = %d, %v", v, err)
	}
	if version, rows := worldCredentialRows(t, world); version != preSplitSchemaVersion || rows != len(want) {
		t.Fatalf("pre-split store changed: version %d rows %d", version, rows)
	}
}
