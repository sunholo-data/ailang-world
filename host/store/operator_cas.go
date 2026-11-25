package store

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"

	"github.com/sunholo-data/ailang-world/host/approval"
)

var ErrOperatorRunner = errors.New("store: operator runner is absent, stale or replaced")

type OperatorCAS struct {
	ID         string
	Digest     [32]byte
	Generation int64
}

type OperatorOwner struct {
	Nonce      string
	Generation int64
}

// The interface transports a sealed law proof. A method implemented by a caller
// cannot supply the verified binding or approver-use fields as plain data.
type operatorConfirmation interface {
	OperatorConfirmation() *approval.SignedOperatorAction
}

func operatorCASRow(ctx context.Context, tx *sql.Tx, cas OperatorCAS) (OperatorActionRow, error) {
	row, err := scanOperatorAction(tx.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions WHERE id = ?`, cas.ID))
	if err == nil && (!bytes.Equal(row.Digest, cas.Digest[:]) || row.Generation != cas.Generation) {
		return OperatorActionRow{}, ErrOperatorConflict
	}
	return row, err
}

// ConfirmOperatorAction consumes an authority proof and persists its exact
// artifact and approver identity atomically with approver use and counter.
func (c *Credentials) ConfirmOperatorAction(ctx context.Context, cas OperatorCAS, artifact []byte, proof operatorConfirmation, now int64) (OperatorActionRow, error) {
	if err := c.operatorAccess(ctx); err != nil {
		return OperatorActionRow{}, err
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return OperatorActionRow{}, err
	}
	defer tx.Rollback()
	row, err := scanOperatorAction(tx.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions WHERE id = ?`, cas.ID))
	if err != nil {
		return OperatorActionRow{}, err
	}
	if !bytes.Equal(row.Digest, cas.Digest[:]) {
		return OperatorActionRow{}, ErrOperatorConflict
	}
	if row.ApprovedBy != "" {
		if !bytes.Equal(row.Artifact, artifact) {
			return OperatorActionRow{}, ErrOperatorConflict
		}
		return row, nil
	}
	if row.State != "requested" || row.Generation != cas.Generation {
		return OperatorActionRow{}, ErrOperatorConflict
	}
	if now >= row.ExpiresAt {
		return OperatorActionRow{}, ErrOperatorExpired
	}
	if now < row.CreatedAt || sha256.Sum256(row.Action) != cas.Digest || len(row.Challenge) != 32 {
		return OperatorActionRow{}, ErrOperatorConflict
	}
	a, err := approval.DecodeOperatorAction(row.Action)
	if err != nil {
		return OperatorActionRow{}, err
	}
	if proof == nil {
		return OperatorActionRow{}, approval.ErrNotVerified
	}
	signed := proof.OperatorConfirmation()
	if err := signed.Binds(a, [32]byte(row.Challenge), row.ExpiresAt); err != nil {
		return OperatorActionRow{}, err
	}
	binding, err := signed.ConfirmationFor(artifact)
	if err != nil || len(artifact) > MaxSurfaceApprovalArtifactBytes {
		return OperatorActionRow{}, approval.ErrNotVerified
	}
	if err := useApprover(ctx, tx, ApproverUse{Name: binding.Name, Format: binding.Format, PublicKey: []byte(binding.PublicKey)}, now); err != nil {
		return OperatorActionRow{}, err
	}
	if id, count, ok := signed.PasskeyUse(); ok {
		var rp, origin string
		var backup bool
		if err := tx.QueryRowContext(ctx, `SELECT rp_id, origin, backup_eligible FROM webauthn_credentials WHERE name = ?`, binding.Name).Scan(&rp, &origin, &backup); err != nil {
			return OperatorActionRow{}, err
		}
		if rp != binding.RPID || origin != binding.Origin || backup != binding.BackupEligible {
			return OperatorActionRow{}, ErrApproverKeyChanged
		}
		if err := usePasskey(ctx, tx, binding.Name, PasskeyUse{CredentialID: id, SignCount: int64(count)}); err != nil {
			return OperatorActionRow{}, err
		}
	}
	res, err := tx.ExecContext(ctx, `UPDATE operator_actions SET state = 'confirmed', approved_by = ?, artifact = ?
		WHERE id = ? AND digest = ? AND generation = ? AND state = 'requested' AND expires_at > ?`,
		signed.ApprovedBy(), artifact, cas.ID, cas.Digest[:], cas.Generation, now)
	if err := operatorChangedOne(res, err); err != nil {
		return OperatorActionRow{}, err
	}
	if err := tx.Commit(); err != nil {
		return OperatorActionRow{}, err
	}
	row.State, row.ApprovedBy, row.Artifact = "confirmed", signed.ApprovedBy(), bytes.Clone(artifact)
	return row, nil
}

// ClaimOperatorAction arbitrates one fresh claim. Runner effect readiness and
// the shared kernel episode lock are prerequisites of the future runner.
func (c *Credentials) ClaimOperatorAction(ctx context.Context, cas OperatorCAS, runner OperatorOwner, now int64) error {
	if err := c.operatorAccess(ctx); err != nil {
		return err
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	row, err := operatorCASRow(ctx, tx, cas)
	if err != nil {
		return err
	}
	if row.State != "confirmed" || row.Owner != "" || now < row.CreatedAt {
		return ErrOperatorConflict
	}
	if now >= row.ExpiresAt {
		return ErrOperatorExpired
	}
	var current bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM operator_runner_state WHERE
		generation = ? AND owner = ? AND seen_at > ? AND seen_at <= ? AND (readiness = 'both' OR readiness = ?))`,
		runner.Generation, runner.Nonce, now-10, now, row.Kind).Scan(&current); err != nil {
		return err
	}
	if !current {
		return ErrOperatorRunner
	}
	res, err := tx.ExecContext(ctx, `UPDATE operator_actions AS a SET state = 'claimed', owner = ?, generation = generation + 1
		WHERE id = ? AND digest = ? AND generation = ? AND state = 'confirmed' AND expires_at > ? AND owner IS NULL
		AND NOT EXISTS (SELECT 1 FROM operator_actions b WHERE b.episode_id = a.episode_id AND b.id <> a.id AND b.owner IS NOT NULL)`,
		runner.Nonce, cas.ID, cas.Digest[:], cas.Generation, now)
	if err := operatorChangedOne(res, err); err != nil {
		return err
	}
	return tx.Commit()
}

func operatorChangedOne(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return ErrOperatorConflict
	}
	return nil
}
