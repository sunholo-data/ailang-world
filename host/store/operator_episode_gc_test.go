package store

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestEpisodeHasUnresolvedOperatorSteps(t *testing.T) {
	for _, tc := range []struct {
		name, kind, phase, state string
		want                     bool
	}{
		{name: "none"},
		{name: "operation_only", state: "running"},
		{name: "reserved", phase: "reserved", state: "running", want: true},
		{name: "uncertain", phase: "reserved", state: "uncertain", want: true},
		{name: "observed_running", phase: "observed", state: "running", want: true},
		{name: "observed_uncertain", phase: "observed", state: "uncertain", want: true},
		{name: "observed_succeeded", phase: "observed", state: "succeeded"},
		{name: "observed_failed", phase: "observed", state: "failed"},
		{name: "observed_expired", phase: "observed", state: "expired"},
		{name: "observed_cancelled", phase: "observed", state: "cancelled"},
		{name: "terminal_with_reservation", phase: "reserved", state: "failed", want: true},
		{name: "files_reserved", kind: "files", phase: "reserved", state: "failed", want: true},
		{name: "files_observed_running", kind: "files", phase: "observed", state: "running", want: true},
		{name: "mint_reserved", kind: "mint", phase: "reserved", state: "failed", want: true},
		{name: "mint_observed_uncertain", kind: "mint", phase: "observed", state: "uncertain", want: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newOperatorPinFixture(t, false)
			if tc.phase != "" {
				var key OperatorStepKey
				if tc.kind == "" {
					key = f.reserveWorktree(t)
				} else {
					// Isolate the lookup by giving this episode exactly one effect kind.
					attempt := f.in.Attempt
					key = OperatorStepKey{attempt.CAS, attempt.Runner.Nonce, tc.kind, attempt.CAS.ID, 1}
					if err := f.c.ReserveOperatorStep(f.ctx, key, pinNow); err != nil {
						t.Fatal(err)
					}
				}
				if tc.phase == "observed" {
					if err := f.c.ObserveOperatorStep(f.ctx, key, []byte("created")); err != nil {
						t.Fatal(err)
					}
				}
			}
			// Materialize durable outcomes (including abnormal combinations)
			// subject to the same state and terminal_at CHECK as production.
			switch tc.state {
			case "running", "uncertain":
				operatorExec(t, f.c, "UPDATE operator_actions SET state=?", tc.state)
			case "succeeded", "failed", "expired", "cancelled":
				operatorExec(t, f.c, "UPDATE operator_actions SET state=?, owner=NULL, terminal_at=?", tc.state, pinNow+1)
			}
			episode := f.in.Start.Episode
			if tc.name == "none" {
				episode = "no-action"
			}
			// A second episode's reservation must never affect this result.
			other := putOperator(t, f.c, operatorRequest(t, f.c, 2, "other-episode", "finish"), 1000)
			operatorExec(t, f.c, `INSERT INTO operator_action_steps (kind,identity,action_id,phase,version,owner,created_at)
				VALUES ('cleanup','other-effect',?,'reserved',1,'other-runner',?)`, other.ID, pinNow)
			path := f.c.lock.dbPath
			before, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			got, err := EpisodeHasUnresolvedOperatorSteps(f.ctx, path, episode)
			if got != tc.want || err != nil {
				t.Fatalf("GC ledger lookup = %t, %v; want %t, nil", got, err, tc.want)
			}
			assertGCLedgerBytes(t, path, before)
		})
	}
}

func TestEpisodeOperatorStepLookupRefusals(t *testing.T) {
	t.Run("no_deadline", func(t *testing.T) {
		c := openTestCredentials(t)
		path := c.lock.dbPath
		before, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		got, err := EpisodeHasUnresolvedOperatorSteps(context.Background(), path, "ep")
		if got || !errors.Is(err, ErrNoDeadline) {
			t.Fatalf("unbounded lookup = %t %v", got, err)
		}
		assertGCLedgerBytes(t, path, before)
	})
	t.Run("missing_store", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), CredentialsFileName)
		got, err := EpisodeHasUnresolvedOperatorSteps(testCtx(t), path, "ep")
		if got || !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("missing store = %t %v", got, err)
		}
		if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("lookup created or replaced credentials: %v", err)
		}
	})
	t.Run("older_schema", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), CredentialsFileName)
		raw := rawDB(t, path)
		if _, err := raw.Exec(credentialsSchemaSQL); err != nil {
			t.Fatal(err)
		}
		setPositiveFixtureVersion(t, raw, 9)
		before, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		got, err := EpisodeHasUnresolvedOperatorSteps(testCtx(t), path, "ep")
		old, ok := errors.AsType[*OlderCredentialsStoreError](err)
		if got || !ok || old.Version != 9 || old.Need != operatorSchemaVersion {
			t.Fatalf("older schema = %t %v", got, err)
		}
		assertGCLedgerBytes(t, path, before)
	})
	t.Run("query_error", func(t *testing.T) {
		c := openTestCredentials(t)
		path := c.lock.dbPath
		operatorExec(t, c, "ALTER TABLE operator_action_steps RENAME TO unavailable_steps")
		before, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		got, err := EpisodeHasUnresolvedOperatorSteps(testCtx(t), path, "ep")
		if got || err == nil || !strings.Contains(err.Error(), "no such table") {
			t.Fatalf("missing ledger table = %t %v", got, err)
		}
		assertGCLedgerBytes(t, path, before)
	})
	t.Run("deadline", func(t *testing.T) {
		c := openTestCredentials(t)
		path := c.lock.dbPath
		before, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithDeadline(context.Background(), time.Unix(1, 0))
		defer cancel()
		got, err := EpisodeHasUnresolvedOperatorSteps(ctx, path, "ep")
		if got || !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expired lookup = %t %v", got, err)
		}
		assertGCLedgerBytes(t, path, before)
	})
}

func assertGCLedgerBytes(t *testing.T, path string, before []byte) {
	t.Helper()
	if after, err := os.ReadFile(path); err != nil || !bytes.Equal(before, after) {
		t.Fatalf("credentials changed during GC lookup: %v", err)
	}
}
