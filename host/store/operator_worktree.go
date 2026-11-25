package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// ReserveOperatorWorktree admits the first effect of a claimed start. The
// caller holds episode exclusion and has reconstructed the reviewed inputs.
func (c *Credentials) ReserveOperatorWorktree(ctx context.Context, at OperatorMintAttempt, identity string, proof *approval.OperatorMintProof, clock func() int64) (OperatorStepKey, error) {
	var key OperatorStepKey
	if err := c.operatorAccess(ctx); err != nil {
		return key, err
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return key, err
	}
	defer tx.Rollback()
	now := clock()
	row, err := operatorCASRow(ctx, tx, at.CAS)
	if err != nil {
		return key, err
	}
	if now >= row.ExpiresAt {
		return key, ErrOperatorExpired
	}
	if row.State != "claimed" || row.Owner != at.Runner.Nonce || len(row.Challenge) != 32 || identity == "" {
		return key, ErrOperatorConflict
	}
	action, err := decodeOperatorMintAction(row)
	if err != nil {
		return key, err
	}
	if c.lock == nil || action.Store != c.lock.dbPath || action.Kind != "start" || !action.Options.Headless {
		return key, ErrOperatorConflict
	}
	b, err := proof.Bind(action, [32]byte(row.Challenge), row.ExpiresAt, row.Artifact)
	if err != nil {
		return key, err
	}
	if row.ApprovedBy != b.Format+":"+b.Name {
		return key, ErrOperatorConflict
	}
	if err := operatorMintApprover(ctx, tx, proof, b, now); err != nil {
		return key, err
	}
	res, err := tx.ExecContext(ctx, `UPDATE operator_actions SET state = 'running', generation = generation + 1
		WHERE id = ? AND digest = ? AND generation = ? AND owner = ? AND state = 'claimed'
		AND EXISTS (SELECT 1 FROM operator_runner_state WHERE generation = ? AND owner = ? AND catalogue_generation = ?
			AND protocol = ? AND seen_at > ? AND seen_at <= ? AND readiness IN ('start','both'))
		AND NOT EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id = ? AND kind <> 'operation')
		AND NOT EXISTS (SELECT 1 FROM episode_starts WHERE episode = ?)
		AND NOT EXISTS (SELECT 1 FROM session_credentials WHERE episode_id = ?)`,
		row.ID, at.CAS.Digest[:], row.Generation, row.Owner, at.Runner.Generation, at.Runner.Nonce,
		at.Catalogue, OperatorRunnerProtocol, now-10, now, row.ID, row.Episode, row.Episode)
	if err := operatorChangedOne(res, err); err != nil {
		return key, err
	}
	key = OperatorStepKey{Action: at.CAS, Owner: row.Owner, Kind: "worktree", Identity: identity, Version: 1}
	key.Action.Generation++
	if err := reserveOperatorStep(ctx, tx, key, now); err != nil {
		return OperatorStepKey{}, err
	}
	return key, tx.Commit()
}

// ReadOperatorWorktree returns only an unresolved worktree reservation. The
// CLI verifies its effect identity and re-reads it under the episode lock.
func (c *Credentials) ReadOperatorWorktree(ctx context.Context, id string) (OperatorStepKey, approval.OperatorAction, error) {
	var key OperatorStepKey
	var action approval.OperatorAction
	row, err := c.OperatorAction(ctx, id, 0)
	if err != nil {
		return key, action, err
	}
	if row.State != "running" && row.State != "uncertain" {
		return key, action, ErrOperatorConflict
	}
	action, err = decodeOperatorMintAction(row)
	if err != nil {
		return key, action, err
	}
	key = OperatorStepKey{Action: OperatorCAS{id, [32]byte(row.Digest), row.Generation}, Owner: row.Owner, Kind: "worktree"}
	err = c.db.QueryRowContext(ctx, `SELECT identity, version FROM operator_action_steps
		WHERE action_id = ? AND kind = 'worktree' AND owner = ? AND phase = 'reserved'
		AND NOT EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id = ? AND kind NOT IN ('operation','worktree'))
		AND (SELECT count(*) FROM operator_action_steps WHERE action_id = ? AND kind = 'worktree') = 1`,
		id, row.Owner, id, id).Scan(&key.Identity, &key.Version)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrOperatorMissing
	}
	return key, action, err
}

const (
	OperatorWorktreeUncertain = ""
	OperatorWorktreeAbsent    = "stopped: absent"
	OperatorWorktreeRetained  = "stopped: retained"
)

// DisposeOperatorWorktree checks the current action and the exact reserved
// step independently. Uncertainty retains ownership and leaves the step
// reserved at its new version. A stop leaves permanent evidence.
func (c *Credentials) DisposeOperatorWorktree(ctx context.Context, key OperatorStepKey, disposition string, now int64) error {
	if err := c.operatorAccess(ctx); err != nil {
		return err
	}
	if key.Kind != "worktree" {
		return ErrOperatorInput
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	switch disposition {
	case OperatorWorktreeUncertain:
		res, err := tx.ExecContext(ctx, `UPDATE operator_action_steps SET version = version + 1
			WHERE action_id = ? AND kind = 'worktree' AND identity = ? AND owner = ? AND version = ? AND phase = 'reserved'`,
			key.Action.ID, key.Identity, key.Owner, key.Version)
		if err := operatorChangedOne(res, err); err != nil {
			return err
		}
	case OperatorWorktreeAbsent, OperatorWorktreeRetained:
		if err := observeOperatorStep(ctx, tx, key, []byte(disposition)); err != nil {
			return err
		}
	default:
		return ErrOperatorInput
	}
	res, err := tx.ExecContext(ctx, `UPDATE operator_actions SET state = CASE WHEN ? THEN 'uncertain' ELSE 'failed' END,
		owner = CASE WHEN ? THEN owner END, terminal_at = CASE WHEN NOT ? THEN ? END, generation = generation + 1
		WHERE id = ? AND digest = ? AND generation = ? AND owner = ? AND kind = 'start' AND state IN ('running','uncertain')
		AND NOT EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id = ? AND kind NOT IN ('operation','worktree'))
		AND (SELECT count(*) FROM operator_action_steps WHERE action_id = ? AND kind = 'worktree') = 1`,
		disposition == OperatorWorktreeUncertain, disposition == OperatorWorktreeUncertain, disposition == OperatorWorktreeUncertain, now,
		key.Action.ID, key.Action.Digest[:], key.Action.Generation, key.Owner, key.Action.ID, key.Action.ID)
	if err := operatorChangedOne(res, err); err != nil {
		return err
	}
	return tx.Commit()
}
