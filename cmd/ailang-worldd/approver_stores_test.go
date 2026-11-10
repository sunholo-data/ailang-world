package main

// Row 201 A1b sprint 3 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint3-plan.md
// G1): approver enroll|list|revoke and approve --list act on the stores their flags name, and on
// the default --db when none is given. Sprint 2's judge found that no test passed
// --credentials-db or relied on the default --db (mutants F-CREDS, AS-CREDFOR and F-DB survived).
// Every store is under t.TempDir(), and HOME is a temp dir (D-WORLD-109).

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
)

// elsewhereStores is an existing world store and a fresh schema-5 credentials store in ANOTHER
// directory: a verb that ignored --credentials-db would look beside db and find nothing there.
func elsewhereStores(t *testing.T) (db, creds string) {
	t.Helper()
	db = dbPathOf(t)
	creds = filepath.Join(t.TempDir(), "elsewhere-credentials.db")
	st, err := store.OpenCredentials(creds)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	return db, creds
}

// homeStores points HOME at a temp dir and makes the default world store
// ($HOME/.ailang/world/world.db) and the credentials store beside it.
func homeStores(t *testing.T) (db, creds string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	db = filepath.Join(home, ".ailang", "world", "world.db")
	if got := defaultWorldPaths().db; got != db {
		t.Fatalf("control: the default --db is %q, want %q", got, db)
	}
	if err := os.MkdirAll(filepath.Dir(db), 0o700); err != nil {
		t.Fatal(err)
	}
	s, err := store.Open(db)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	creds, err = credentialsFor(db, "")
	if err != nil {
		t.Fatal(err)
	}
	c, err := store.OpenCredentials(creds)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
	return db, creds
}

// assertNoStoreBeside fails if a credentials store exists beside db (the verb used the default).
func assertNoStoreBeside(t *testing.T, db string) {
	t.Helper()
	beside, err := credentialsFor(db, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(beside); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("a credentials store exists beside the world store (%s, err %v): --credentials-db was ignored", beside, err)
	}
}

// G1: each approver subcommand reads and writes the --credentials-db it is given, and the default
// --db's store when no --db is given.
func TestApproverVerbsUseTheNamedStores(t *testing.T) {
	ed := fixtureKeyOf(t, "ed25519.pub")
	wantLine := "mark-mac  ssh  ssh-ed25519 " + edFingerprint + "  enrolled 2027-01-15T08:00:00Z  last used never  active\n"
	yes := func() sessionEnv {
		return approveEnv(&lockedTerm{in: strings.NewReader("y\n")}, nil, atTime(approveT0+100))
	}
	type arm struct {
		name   string
		stores func(t *testing.T) (db, creds string)
		flags  func(db, creds string) []string // the store flags this arm passes
		beside bool                            // assert no store was created beside db
	}
	arms := []arm{
		{"credentials-db", elsewhereStores, func(db, creds string) []string { return []string{"--db", db, "--credentials-db", creds} }, true},
		{"default-db", homeStores, func(string, string) []string { return nil }, false},
	}
	for _, a := range arms {
		t.Run("list/"+a.name, func(t *testing.T) {
			db, creds := a.stores(t)
			seedApprover(t, creds, "mark-mac", ed, approveT0)
			var out, errw bytes.Buffer
			code := runApprover(append([]string{"list"}, a.flags(db, creds)...), &out, &errw, approveEnv(nil, nil, atTime(approveT0)))
			if code != exitOK || out.String() != wantLine {
				t.Fatalf("code %d, stdout %q, want exitOK and %q; stderr %q", code, out.String(), wantLine, errw.String())
			}
			if a.beside {
				assertNoStoreBeside(t, db)
			}
		})
		t.Run("enroll/"+a.name, func(t *testing.T) {
			db, creds := a.stores(t)
			args := append([]string{"enroll", "--name", "mark-mac", "--format", "ssh", "--pubkey", fixturePub("ed25519.pub")}, a.flags(db, creds)...)
			var out, errw bytes.Buffer
			if code := runApprover(args, &out, &errw, yes()); code != exitOK {
				t.Fatalf("code %d, want exitOK; stderr %q", code, errw.String())
			}
			rows := approverRows(t, creds)
			if len(rows) != 1 || rows[0].Name != "mark-mac" || !bytes.Equal(rows[0].PublicKey, ed.Blob) {
				t.Fatalf("approver rows in %s: %+v, want mark-mac's ed25519 key", creds, rows)
			}
			if a.beside {
				assertNoStoreBeside(t, db)
			}
		})
		t.Run("revoke/"+a.name, func(t *testing.T) {
			db, creds := a.stores(t)
			seedApprover(t, creds, "mark-mac", ed, approveT0)
			args := append([]string{"revoke", "mark-mac"}, a.flags(db, creds)...)
			var out, errw bytes.Buffer
			if code := runApprover(args, &out, &errw, yes()); code != exitOK {
				t.Fatalf("code %d, want exitOK; stderr %q", code, errw.String())
			}
			if rows := approverRows(t, creds); len(rows) != 1 || rows[0].RevokedAt != approveT0+100 {
				t.Fatalf("approver rows in %s: %+v, want mark-mac revoked at %d", creds, rows, approveT0+100)
			}
			if a.beside {
				assertNoStoreBeside(t, db)
			}
		})
	}
}

// putPendingRequest stores one pending request for episode ep1 (id mr-0102030405060708) in creds.
func putPendingRequest(t *testing.T, creds string) {
	t.Helper()
	st, err := store.OpenCredentials(creds)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = st.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := store.MintRequestRow{RequestID: firstRequestID, Request: []byte("{}"), Digest: strings.Repeat("ab", 32),
		Challenge: seqChallenge(), EpisodeID: "ep1", CreatedAt: approveT0, ExpiresAt: approveT0 + 600}
	if err := st.PutMintRequest(ctx, row, approveT0); err != nil {
		t.Fatal(err)
	}
}

// G2: approve --list reads the --credentials-db it is given, and the default --db's store when no
// --db is given.
func TestApproveListUsesTheNamedStores(t *testing.T) {
	want := firstRequestID + "  episode ep1  digest abababababababab  expires 2027-01-15T08:10:00Z\n"
	t.Run("credentials-db", func(t *testing.T) {
		db, creds := elsewhereStores(t)
		putPendingRequest(t, creds)
		var out, errw bytes.Buffer
		code := runApprove([]string{"--list", "--db", db, "--credentials-db", creds}, &out, &errw, approveEnv(nil, nil, atTime(approveT0+1)))
		if code != exitOK || out.String() != want {
			t.Fatalf("code %d, stdout %q, want exitOK and %q; stderr %q", code, out.String(), want, errw.String())
		}
		assertNoStoreBeside(t, db)
	})
	t.Run("default-db", func(t *testing.T) {
		_, creds := homeStores(t)
		putPendingRequest(t, creds)
		var out, errw bytes.Buffer
		code := runApprove([]string{"--list"}, &out, &errw, approveEnv(nil, nil, atTime(approveT0+1)))
		if code != exitOK || out.String() != want {
			t.Fatalf("code %d, stdout %q, want exitOK and %q; stderr %q", code, out.String(), want, errw.String())
		}
	})
}
