package store

// Row 201 A1b sprint 4 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint4-plan.md §3
// D1): DecidedMintRequests, the consumed rows approve --list --decided and approver revoke's K8
// read. Every store is under t.TempDir().

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// putDecided inserts request n as consumed at `at` by approvedBy, with credential id "%064x" of n.
func putDecided(t *testing.T, db *sql.DB, n int, at int64, approvedBy string) MintRequestRow {
	t.Helper()
	r := testMintRequestRow(n, at-10)
	r.Consumed, r.ConsumedAt, r.ApprovedBy, r.CredentialID = true, at, approvedBy, fmt.Sprintf("%064x", n)
	if _, err := db.Exec(`INSERT INTO mint_requests (`+mintRequestColumns+`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.RequestID, r.Request, r.Digest, r.Challenge, r.EpisodeID, r.CreatedAt, r.ExpiresAt, at, approvedBy, r.CredentialID); err != nil {
		t.Fatal(err)
	}
	return r
}

// CS6: consumed rows only, newest consumed_at first and request_id descending among equal times,
// every column returned; approvedBy keeps only its own rows; the cap holds; schema 4 is read, and
// a store below it is refused.
func TestDecidedMintRequests(t *testing.T) {
	if MaxDecidedMintRequests != 64 {
		t.Fatalf("MaxDecidedMintRequests = %d, want 64 (design §2.4 [r2])", MaxDecidedMintRequests)
	}
	ctx := testCtx(t)
	c := openTestCredentials(t)
	if err := c.PutMintRequest(ctx, testMintRequestRow(100, 1000), 1000); err != nil { // pending: never listed
		t.Fatal(err)
	}
	if got, err := c.DecidedMintRequests(ctx, "", 0); err != nil || len(got) != 0 {
		t.Fatalf("no decided rows: (%+v, %v), want none", got, err)
	}
	r1 := putDecided(t, c.db, 1, 2000, "tty")
	r2 := putDecided(t, c.db, 2, 3000, "ssh:x") // r2 and r3 share consumed_at: request_id DESC puts r3 first
	r3 := putDecided(t, c.db, 3, 3000, "tty")
	r4 := putDecided(t, c.db, 4, 1500, "ssh:x")
	got, err := c.DecidedMintRequests(ctx, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	if want := []MintRequestRow{r3, r2, r1, r4}; !reflect.DeepEqual(got, want) {
		t.Fatalf("decided rows\n got %+v\nwant %+v", got, want)
	}
	if got, err := c.DecidedMintRequests(ctx, "ssh:x", 0); err != nil || !reflect.DeepEqual(got, []MintRequestRow{r2, r4}) {
		t.Fatalf("approved by ssh:x: (%+v, %v), want r2 and r4", got, err)
	}
	if got, err := c.DecidedMintRequests(ctx, "ssh:nobody", 0); err != nil || len(got) != 0 {
		t.Fatalf("approved by ssh:nobody: (%+v, %v), want none", got, err)
	}
	if got, err := c.DecidedMintRequests(ctx, "", 2); err != nil || !reflect.DeepEqual(got, []MintRequestRow{r3, r2}) {
		t.Fatalf("limit 2: (%+v, %v), want r3 and r2", got, err)
	}

	t.Run("cap", func(t *testing.T) {
		c := openTestCredentials(t)
		for n := 1; n <= MaxDecidedMintRequests+6; n++ {
			putDecided(t, c.db, n, int64(2000+n), "tty")
		}
		for _, limit := range []int{0, -1, MaxDecidedMintRequests + 1} {
			if got, err := c.DecidedMintRequests(ctx, "", limit); err != nil || len(got) != MaxDecidedMintRequests {
				t.Fatalf("limit %d: %d rows (%v), want %d", limit, len(got), err, MaxDecidedMintRequests)
			}
		}
	})
	t.Run("v4", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), CredentialsFileName)
		makeV4Credentials(t, path, nil)
		c4, err := OpenCredentials(path)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = c4.Close() }()
		r := putDecided(t, c4.db, 7, 2000, "tty")
		if got, err := c4.DecidedMintRequests(ctx, "", 0); err != nil || !reflect.DeepEqual(got, []MintRequestRow{r}) {
			t.Fatalf("v4: (%+v, %v), want the one decided row", got, err)
		}
	})
	t.Run("v3", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), CredentialsFileName)
		makeV3Credentials(t, path, nil)
		c3, err := OpenCredentials(path)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = c3.Close() }()
		var older *OlderCredentialsStoreError
		if _, err := c3.DecidedMintRequests(ctx, "", 0); !errors.As(err, &older) || older.Need != 4 || !strings.Contains(err.Error(), "restart `ailang-worldd serve`") {
			t.Fatalf("v3: err %v, want the mint-requests refusal (need 4)", err)
		}
	})
	t.Run("no-deadline", func(t *testing.T) {
		if _, err := c.DecidedMintRequests(context.Background(), "", 0); !errors.Is(err, ErrNoDeadline) {
			t.Fatalf("err %v, want ErrNoDeadline", err)
		}
	})
}
