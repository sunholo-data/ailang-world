package store

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
)

func operatorRequest(t *testing.T, c *Credentials, n int, episode, kind string) OperatorRequest {
	t.Helper()
	a := approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: fmt.Sprintf("oa-%032x", n),
		Kind: kind, Episode: episode, Store: c.lock.dbPath, Review: "review",
		Options:    approval.OperatorOptions{Preset: "sol"},
		Resolution: []approval.OperatorResolved{{Name: "agent", Value: "codex"}}}
	if kind == "start" {
		a.Options.Headless, a.Task = true, "test task"
	}
	a.Argv = approval.OperatorArgv(kind, episode, a.Options)
	return OperatorRequest{Action: a, OperationID: fmt.Sprintf("%032x", n), CreatedAt: 1000, ExpiresAt: 1600}
}

func putOperator(t *testing.T, c *Credentials, r OperatorRequest, now int64) OperatorActionRow {
	t.Helper()
	row, err := c.RequestOperatorAction(testCtx(t), r, now)
	if err != nil {
		t.Fatal(err)
	}
	return row
}

func getOperator(t *testing.T, c *Credentials, id string, now int64) OperatorActionRow {
	t.Helper()
	row, err := c.OperatorAction(testCtx(t), id, now)
	if err != nil {
		t.Fatal(err)
	}
	return row
}

func operatorExec(t *testing.T, c *Credentials, query string, args ...any) {
	t.Helper()
	if _, err := c.db.Exec(query, args...); err != nil {
		t.Fatal(err)
	}
}

func operatorCounts(t *testing.T, c *Credentials, actions, steps int) {
	t.Helper()
	var a, s int
	err := c.db.QueryRow(`SELECT (SELECT count(*) FROM operator_actions), (SELECT count(*) FROM operator_action_steps)`).Scan(&a, &s)
	if err != nil || a != actions || s != steps {
		t.Fatalf("actions/steps: %d/%d, want %d/%d (err %v)", a, s, actions, steps, err)
	}
}

func TestOperatorAdmissionBoundaries(t *testing.T) {
	c := openTestCredentials(t)
	ctx := testCtx(t)
	a := operatorRequest(t, c, 1, "ep-1", "start")
	accepted := putOperator(t, c, a, 1000)
	if accepted.State != "requested" || accepted.ApprovedBy != "" || len(accepted.Artifact) != 0 || accepted.Owner != "" {
		t.Fatalf("request is not inert: %+v", accepted)
	}

	cases := []struct {
		name string
		edit func(*OperatorRequest)
		now  int64
		want error
	}{
		{"operation ID", func(r *OperatorRequest) { r.OperationID = "not-128-bit-hex" }, 1000, ErrOperatorInput},
		{"other store", func(r *OperatorRequest) { r.Action.Store += ".other" }, 1000, ErrOperatorInput},
		{"creation zero", func(r *OperatorRequest) { r.CreatedAt, r.ExpiresAt = 0, 1 }, 0, ErrOperatorInput},
		{"future creation", func(r *OperatorRequest) { r.CreatedAt++ }, 1000, ErrOperatorInput},
		{"lifetime 601", func(r *OperatorRequest) { r.ExpiresAt++ }, 1000, ErrOperatorInput},
		{"at expiry", func(r *OperatorRequest) {}, 1600, ErrOperatorExpired},
		{"past expiry", func(r *OperatorRequest) {}, 1601, ErrOperatorExpired},
		{"action ID", func(r *OperatorRequest) { r.Action.ID = "invalid" }, 1000, approval.ErrOperatorAction},
		{"operation changed", func(r *OperatorRequest) { r.OperationID = a.OperationID }, 1000, ErrOperatorConflict},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := operatorRequest(t, c, 2, "ep-2", "start")
			tc.edit(&r)
			if _, err := c.RequestOperatorAction(ctx, r, tc.now); !errors.Is(err, tc.want) {
				t.Fatalf("refusal %v, want %v", err, tc.want)
			}
			operatorCounts(t, c, 1, 1)
			if got := getOperator(t, c, a.Action.ID, 1000); got.State != "requested" || !bytes.Equal(got.Action, accepted.Action) {
				t.Fatal("refusal changed the existing request")
			}
		})
	}
	finish := operatorRequest(t, c, 2, "ep-1", "finish")
	if _, err := c.RequestOperatorAction(ctx, finish, 1000); !errors.Is(err, ErrOperatorConflict) {
		t.Fatalf("cross-kind collision: %v", err)
	}
	// The DB must also exclude competing writes which reach it after preflight.
	if _, err := c.db.Exec(`INSERT INTO operator_actions
		(id, episode_id, kind, review, digest, action, challenge, created_at, expires_at, state, generation)
		SELECT ?, episode_id, 'finish', review, digest, action, challenge, created_at, expires_at, state, generation
		FROM operator_actions WHERE id = ?`, finish.Action.ID, a.Action.ID); err == nil {
		t.Fatal("database admitted a second active action")
	}
	operatorCounts(t, c, 1, 1)
	for i := 2; i <= 16; i++ {
		putOperator(t, c, operatorRequest(t, c, i, fmt.Sprintf("ep-%d", i), "start"), 1599)
	}
	if _, err := c.RequestOperatorAction(ctx, operatorRequest(t, c, 17, "ep-17", "start"), 1599); !errors.Is(err, ErrOperatorCapacity) {
		t.Fatalf("17th pending: %v", err)
	}
	operatorCounts(t, c, 16, 16)
	t.Run("PendingProjection", func(t *testing.T) {
		pending, err := c.PendingOperatorActions(testCtx(t), 1599)
		if err != nil || len(pending) != 16 {
			t.Fatalf("pending at capacity: %d, want 16 (err %v)", len(pending), err)
		}
		operatorExec(t, c, `UPDATE operator_actions SET state = 'confirmed', approved_by = 'ssh:fixture', artifact = action WHERE id = ?`, accepted.ID)
		pending, err = c.PendingOperatorActions(testCtx(t), 1599)
		if err != nil || len(pending) != 15 {
			t.Fatalf("pending after confirmation: %d, want 15 (err %v)", len(pending), err)
		}
		for _, row := range pending {
			if row.ID == accepted.ID {
				t.Fatalf("confirmed action still pending: %s", row.ID)
			}
		}
	})
	finish.CreatedAt, finish.ExpiresAt = 1600, 2200
	finish.Action.ID, finish.OperationID = "oa-00000000000000000000000000000011", "00000000000000000000000000000011"
	if got := putOperator(t, c, finish, 1600); got.Kind != "finish" || got.State != "requested" {
		t.Fatalf("finish shape or exact-expiry reclamation lost: %+v", got)
	}
	operatorCounts(t, c, 17, 17)
}

func TestOperatorDuplicateTap(t *testing.T) {
	c := openTestCredentials(t)
	first := operatorRequest(t, c, 1, "ep", "start")
	second := operatorRequest(t, c, 2, "ep", "start")
	second.Challenge[0] = 2
	row := putOperator(t, c, first, 1000)
	want, err := approval.OperatorReviewDigest(first.Action)
	if err != nil || !bytes.Equal(row.Review, want[:]) || bytes.Equal(row.Review, row.Digest) {
		t.Fatal("store must key on the review digest")
	}
	for _, r := range []OperatorRequest{first, second} {
		if got := putOperator(t, c, r, 1000); got.ID != row.ID || !bytes.Equal(got.Challenge, row.Challenge) {
			t.Fatal("duplicate did not return the original link and challenge")
		}
	}
	if _, err := c.db.Exec(`INSERT INTO operator_action_steps SELECT * FROM operator_action_steps WHERE identity = ?`, first.OperationID); err == nil {
		t.Fatal("database admitted a duplicate operation ID")
	}
	second.Action.Review += " changed"
	if _, err := c.RequestOperatorAction(testCtx(t), second, 1000); !errors.Is(err, ErrOperatorConflict) {
		t.Fatalf("second tap did not acquire durable id binding: %v", err)
	}
	operatorCounts(t, c, 1, 2)
}

func TestOperatorReconnect(t *testing.T) {
	first := openTestCredentials(t)
	path := first.display
	req := operatorRequest(t, first, 1, "restart", "start")
	want := putOperator(t, first, req, 1000)
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	if got := putOperator(t, second, req, 1010); got.ID != want.ID || !bytes.Equal(got.Action, want.Action) || !bytes.Equal(got.Challenge, want.Challenge) {
		t.Fatal("lost request on reopen")
	}
	operatorCounts(t, second, 1, 1)
}

func TestOperatorPrune(t *testing.T) {
	c := openTestCredentials(t)
	const day30 = 30 * 86400
	r := operatorRequest(t, c, 1, "expired", "start")
	putOperator(t, c, r, 1000)
	if getOperator(t, c, r.Action.ID, 1599).State != "requested" ||
		getOperator(t, c, r.Action.ID, 1600).State != "expired" {
		t.Fatal("logical expiry boundary")
	}
	var stored string
	if err := c.db.QueryRow(`SELECT state FROM operator_actions WHERE id = ?`, r.Action.ID).Scan(&stored); err != nil || stored != "requested" {
		t.Fatalf("logical expiry read changed durable state: %q %v", stored, err)
	}
	for i, state := range []string{"uncertain", "failed", "succeeded"} {
		q := operatorRequest(t, c, i+2, state, "start")
		putOperator(t, c, q, 1000)
		if state == "uncertain" {
			operatorExec(t, c, `UPDATE operator_actions SET state = ?, owner = 'runner', approved_by = 'ssh:mark', artifact = x'01' WHERE id = ?`, state, q.Action.ID)
		} else {
			operatorExec(t, c, `UPDATE operator_actions SET state = ?, terminal_at = 1600 WHERE id = ?`, state, q.Action.ID)
		}
		phase := "reserved"
		if state == "succeeded" {
			phase = "observed"
		}
		operatorExec(t, c, `INSERT INTO operator_action_steps VALUES ('launch', ?, ?, ?, 1, 'runner', 1000, NULL)`, state, q.Action.ID, phase)
	}
	for _, now := range []int64{1600 + day30 - 1, 1600 + day30, 1600 + day30 + 1} {
		if err := c.PruneOperatorActions(testCtx(t), now); err != nil {
			t.Fatal(err)
		}
		for i, state := range []string{"expired", "uncertain", "failed", "succeeded"} {
			got := getOperator(t, c, operatorRequest(t, c, i+1, state, "start").Action.ID, now)
			pruned := now > 1600+day30 && (state == "expired" || state == "succeeded")
			if got.State != state || (len(got.Action) == 0) != pruned {
				t.Fatalf("prune %s at %d: %+v", state, now, got)
			}
		}
		operatorCounts(t, c, 4, 7)
	}
	retry := operatorRequest(t, c, 10, "expired", "start")
	retry.OperationID, retry.CreatedAt, retry.ExpiresAt = r.OperationID, 1600+day30, 2200+day30
	if got := putOperator(t, c, retry, retry.CreatedAt); got.ID != r.Action.ID || got.State != "expired" {
		t.Fatal("prune lost the retired operation identity")
	}
	retry.OperationID = "00000000000000000000000000000010"
	if got := putOperator(t, c, retry, retry.CreatedAt); got.ID != retry.Action.ID || got.State != "requested" {
		t.Fatal("fresh review could not supersede expired intent")
	}
}

func TestOperatorTerminalExclusion(t *testing.T) {
	c := openTestCredentials(t)
	path := c.display
	r := operatorRequest(t, c, 1, "lock", "start")
	putOperator(t, c, r, 1000)
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	contender := make(chan error, 1)
	go func() {
		other, err := OpenCredentials(path)
		if other != nil {
			other.Close()
		}
		contender <- err
	}()
	select {
	case err := <-contender:
		if !IsWriterAlreadyActive(err) {
			t.Fatalf("second writer: %v", err)
		}
	case <-ctx.Done():
		_ = c.Close()
		t.Fatal("credentials writer open waited on the incumbent")
	}
	operatorCounts(t, c, 1, 1)
	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	if _, err := second.RequestOperatorAction(testCtx(t), operatorRequest(t, second, 2, "lock", "finish"), 1001); !errors.Is(err, ErrOperatorConflict) {
		t.Fatalf("exclusive episode after the writer reopened: %v", err)
	}
	operatorCounts(t, second, 1, 1)
}

func TestOperatorSchema10Constraints(t *testing.T) {
	c := openTestCredentials(t)
	r := operatorRequest(t, c, 1, "schema", "start")
	putOperator(t, c, r, 1000)
	operatorExec(t, c, `INSERT INTO operator_runner_state VALUES (1, 7, 'runner', 1, 'test', 1000, 'start', 1, NULL)`)
	cases := []string{
		`INSERT INTO operator_runner_state SELECT * FROM operator_runner_state`,
		`INSERT INTO operator_actions SELECT NULL, 'other', kind, review, digest, action, challenge,
			created_at, expires_at, state, generation, owner, approved_by, artifact, terminal_at, v3_id FROM operator_actions`,
		`UPDATE operator_actions SET kind = 'other'`,
		`UPDATE operator_actions SET review = zeroblob(31)`,
		`UPDATE operator_actions SET digest = zeroblob(31)`,
		`UPDATE operator_actions SET action = NULL`,
		`UPDATE operator_actions SET challenge = zeroblob(31)`,
		`UPDATE operator_actions SET expires_at = created_at`,
		`UPDATE operator_actions SET state = 'other'`,
		`UPDATE operator_actions SET artifact = action`,
		`UPDATE operator_actions SET approved_by = 'ssh:fixture', artifact = action`,
		`UPDATE operator_actions SET state = 'confirmed'`,
		`UPDATE operator_actions SET state = 'claimed', approved_by = 'ssh:fixture', artifact = action`,
		`UPDATE operator_actions SET terminal_at = 1600`,
		`UPDATE operator_actions SET v3_id = 'vr-0000000000000001'`,
		`UPDATE operator_actions SET kind = 'finish', state = 'awaiting-pr-confirmation',
			approved_by = 'ssh:fixture', artifact = action, owner = 'runner'`,
		`UPDATE operator_action_steps SET action_id = 'missing'`,
		`INSERT INTO operator_actions SELECT id, 'other', kind, review, digest, action, challenge,
			created_at, expires_at, state, generation, owner, approved_by, artifact, terminal_at, v3_id FROM operator_actions`,
		`UPDATE operator_action_steps SET kind = 'unknown', owner = 'runner'`,
		`UPDATE operator_action_steps SET kind = 'launch', owner = 'runner', phase = 'unknown'`,
		`UPDATE operator_action_steps SET owner = 'x'`,
		`UPDATE operator_action_steps SET kind = 'launch'`,
		`UPDATE operator_action_steps SET kind = 'launch', phase = 'reserved', owner = 'x', receipt = zeroblob(1)`,
		`INSERT INTO operator_runner_state VALUES (2, 7, 'runner', 1, 'test', 1000, 'ready', 1, NULL)`,
	}
	for _, column := range []string{"episode_id", "kind", "review", "digest", "created_at", "expires_at", "state", "generation"} {
		cases = append(cases, "UPDATE operator_actions SET "+column+" = NULL")
	}
	for _, column := range []string{"action_id", "identity", "kind", "phase", "version", "created_at"} {
		cases = append(cases, "UPDATE operator_action_steps SET "+column+" = NULL")
	}
	for _, column := range []string{"generation", "owner", "protocol", "build", "seen_at", "readiness", "catalogue_generation"} {
		cases = append(cases, "UPDATE operator_runner_state SET "+column+" = NULL")
	}
	for _, query := range cases {
		if _, err := c.db.Exec(query); err == nil || !strings.Contains(err.Error(), "constraint failed") {
			t.Errorf("%s: %v, want SQLite constraint failure", query, err)
		}
	}
	operatorCounts(t, c, 1, 1)
}
