package store

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestNextOperatorStart(t *testing.T) {
	c, ctx := openTestCredentials(t), testCtx(t)
	add := func(n int, kind string, created, expires int64, confirm bool) OperatorActionRow {
		t.Helper()
		req := operatorRequest(t, c, n, fmt.Sprintf("ep-%d", n), kind)
		req.CreatedAt, req.ExpiresAt = created, expires
		// Insert in reverse tie-break order to exercise the explicit sort.
		row := putOperator(t, c, req, created)
		if confirm {
			signOperatorPinActionAt(t, c, row, req.Action, created)
		}
		return getOperator(t, c, row.ID, created)
	}
	path := c.lock.dbPath
	check := func(now int64, want OperatorActionRow) {
		t.Helper()
		before, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		got, err := c.NextOperatorStart(ctx, now)
		if want.ID == "" {
			if !errors.Is(err, ErrOperatorMissing) || got.ID != "" {
				t.Fatalf("no eligible start at %d: %q %v", now, got.ID, err)
			}
		} else if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("next start at %d: %q %v, want %q", now, got.ID, err, want.ID)
		}
		assertGCLedgerBytes(t, path, before)
		if want.ID != "" && !bytes.Equal(getOperator(t, c, want.ID, now).Digest, want.Digest) {
			t.Fatal("read changed the action digest")
		}
	}
	add(1, "finish", 1000, 1600, true)
	add(2, "start", 1000, 1600, false)
	add(3, "start", 1000, 1100, true)
	check(1100, OperatorActionRow{})
	add(4, "start", 1200, 1600, true)
	check(1100, OperatorActionRow{})
	first := add(6, "start", 1100, 1600, true)
	winner := add(5, "start", 1100, 1600, true)
	check(1100, winner)
	check(1200, winner)
	operatorExec(t, c, "UPDATE operator_actions SET state='claimed', owner='runner' WHERE id=?", winner.ID)
	check(1100, first)
	check(1600, OperatorActionRow{})
}
