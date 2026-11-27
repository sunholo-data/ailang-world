package store

import (
	"context"
	"database/sql"
	"errors"
)

type OperatorStepKey struct {
	Action                OperatorCAS
	Owner, Kind, Identity string
	Version               int64
}

// ReserveOperatorStep records identity before a future caller crosses an
// effect boundary. The caller must first recheck its own frozen proof.
func (c *Credentials) ReserveOperatorStep(ctx context.Context, key OperatorStepKey, now int64) error {
	if err := c.operatorAccess(ctx); err != nil {
		return err
	}
	if (key.Kind == "operation" || key.Kind == "launch") || key.Identity == "" || key.Owner == "" || key.Version != 1 {
		return ErrOperatorInput
	}
	return reserveOperatorStep(ctx, c.db, key, now)
}

type operatorStepWriter interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

func reserveOperatorStep(ctx context.Context, db operatorStepWriter, key OperatorStepKey, now int64) error {
	result, err := db.ExecContext(ctx, `INSERT INTO operator_action_steps
		(kind, identity, action_id, phase, version, owner, created_at)
		SELECT ?, ?, id, 'reserved', ?, owner, ? FROM operator_actions
		WHERE id = ? AND digest = ? AND generation = ? AND owner = ?
		AND state IN ('claimed', 'running') AND expires_at > ? AND created_at <= ?
		AND (SELECT count(*) FROM operator_action_steps WHERE action_id = ? AND kind <> 'operation') < 32
		ON CONFLICT(kind, identity) DO NOTHING`, key.Kind, key.Identity, key.Version, now,
		key.Action.ID, key.Action.Digest[:], key.Action.Generation, key.Owner, now, now, key.Action.ID)
	return operatorChangedOne(result, err)
}

// ObserveOperatorStep records the outcome of a previously reserved effect.
// Expiry prevents new reservations but does not prevent exact observation.
func (c *Credentials) ObserveOperatorStep(ctx context.Context, key OperatorStepKey, receipt []byte) error {
	if err := c.operatorAccess(ctx); err != nil {
		return err
	}
	if key.Kind == "launch" || len(receipt) == 0 || len(receipt) > MaxSurfaceApprovalArtifactBytes {
		return ErrOperatorInput
	}
	return observeOperatorStep(ctx, c.db, key, receipt)
}

func observeOperatorStep(ctx context.Context, db operatorStepWriter, key OperatorStepKey, receipt []byte) error {
	result, err := db.ExecContext(ctx, `UPDATE operator_action_steps
		SET phase = 'observed', version = version + 1, receipt = ?
		WHERE action_id = ? AND kind = ? AND identity = ? AND version = ? AND owner = ? AND phase = 'reserved'
		AND EXISTS (SELECT 1 FROM operator_actions WHERE id = ? AND digest = ? AND generation = ? AND owner = ?)`,
		receipt, key.Action.ID, key.Kind, key.Identity, key.Version, key.Owner,
		key.Action.ID, key.Action.Digest[:], key.Action.Generation, key.Owner)
	return operatorChangedOne(result, err)
}

type OperatorStepReceipt struct {
	Phase, Owner string
	Version      int64
	Receipt      []byte
}

// OperatorStep reads a reservation or recorded observation across a restart.
func (c *Credentials) OperatorStep(ctx context.Context, key OperatorStepKey) (OperatorStepReceipt, error) {
	if err := c.operatorAccess(ctx); err != nil {
		return OperatorStepReceipt{}, err
	}
	var step OperatorStepReceipt
	err := c.db.QueryRowContext(ctx, `SELECT phase, coalesce(owner,''), version, receipt FROM operator_action_steps
		WHERE action_id = ? AND kind = ? AND identity = ?`, key.Action.ID, key.Kind, key.Identity).
		Scan(&step.Phase, &step.Owner, &step.Version, &step.Receipt)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrOperatorMissing
	}
	return step, err
}
