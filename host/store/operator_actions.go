package store

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"regexp"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
)

var (
	ErrOperatorInput    = errors.New("store: invalid operator request")
	ErrOperatorConflict = errors.New("store: operator action conflicts with existing work")
	ErrOperatorCapacity = errors.New("store: sixteen operator actions already pending")
	ErrOperatorMissing  = errors.New("store: operator action missing")
	ErrOperatorExpired  = errors.New("store: operator action expired")
)

const operatorColumns = `id, episode_id, kind, review, digest, action, challenge, created_at, expires_at, state,
	generation, coalesce(owner, ''), coalesce(approved_by, ''), artifact, coalesce(terminal_at, 0), coalesce(v3_id, '')`

type OperatorActionRow struct {
	ID, Episode, Kind, State, Owner, ApprovedBy, V3ID string
	Review, Digest, Action, Challenge, Artifact       []byte
	CreatedAt, ExpiresAt, Generation, TerminalAt      int64
}

type OperatorRequest struct {
	Action               approval.OperatorAction
	OperationID          string
	Challenge            [32]byte
	CreatedAt, ExpiresAt int64
}

var operatorOperationID = regexp.MustCompile(`^[0-9a-f]{32}$`)

func (c *Credentials) operatorAccess(ctx context.Context) error {
	if err := requireDeadline(ctx); err != nil {
		return err
	}
	if c.version < 10 {
		return &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: 10, Feature: "operator actions"}
	}
	return nil
}

func scanOperatorAction(row *sql.Row) (a OperatorActionRow, err error) {
	err = row.Scan(&a.ID, &a.Episode, &a.Kind, &a.Review, &a.Digest, &a.Action, &a.Challenge,
		&a.CreatedAt, &a.ExpiresAt, &a.State, &a.Generation, &a.Owner, &a.ApprovedBy, &a.Artifact, &a.TerminalAt, &a.V3ID)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrOperatorMissing
	}
	return a, err
}

// OperatorAction reads the durable row and projects expiry for unclaimed work.
func (c *Credentials) OperatorAction(ctx context.Context, id string, now int64) (OperatorActionRow, error) {
	if err := c.operatorAccess(ctx); err != nil {
		return OperatorActionRow{}, err
	}
	row, err := scanOperatorAction(c.db.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions WHERE id = ?`, id))
	if (row.State == "requested" || row.State == "confirmed") && row.ExpiresAt <= now {
		row.State, row.TerminalAt = "expired", row.ExpiresAt
	}
	return row, err
}

// NextOperatorStart reads the oldest eligible start without assigning ownership.
// The caller must take the canonical episode lock and compare the row again
// before claiming it.
func (c *Credentials) NextOperatorStart(ctx context.Context, now int64) (OperatorActionRow, error) {
	if err := c.operatorAccess(ctx); err != nil {
		return OperatorActionRow{}, err
	}
	return scanOperatorAction(c.db.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions
        WHERE kind = 'start' AND state = 'confirmed' AND created_at <= ? AND expires_at > ?
        ORDER BY created_at, id LIMIT 1`, now, now))
}

// RequestOperatorAction records inert work or permanently binds another tap's
// operation ID to the existing action. The active-episode index arbitrates starts
// and finishes together.
func (c *Credentials) RequestOperatorAction(ctx context.Context, request OperatorRequest, now int64) (OperatorActionRow, error) {
	if err := c.operatorAccess(ctx); err != nil {
		return OperatorActionRow{}, err
	}
	a := request.Action
	if !operatorOperationID.MatchString(request.OperationID) ||
		request.CreatedAt <= 0 || now < request.CreatedAt ||
		request.ExpiresAt-request.CreatedAt > int64(approval.OperatorActionLifetime/time.Second) ||
		c.lock == nil || a.Store != c.lock.dbPath {
		return OperatorActionRow{}, ErrOperatorInput
	}
	if now >= request.ExpiresAt {
		return OperatorActionRow{}, ErrOperatorExpired
	}
	action, err := approval.CanonicalOperatorAction(a)
	if err != nil {
		return OperatorActionRow{}, err
	}
	review, err := approval.OperatorReviewDigest(a)
	if err != nil {
		return OperatorActionRow{}, err
	}
	digest, err := approval.OperatorSignatureDigest(a)
	if err != nil {
		return OperatorActionRow{}, err
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return OperatorActionRow{}, err
	}
	defer tx.Rollback()
	if err := expireOperatorActions(ctx, tx, now); err != nil {
		return OperatorActionRow{}, err
	}
	row, err := scanOperatorAction(tx.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions WHERE id =
		(SELECT action_id FROM operator_action_steps WHERE kind = 'operation' AND identity = ?)`, request.OperationID))
	switch {
	case err == nil:
		if !bytes.Equal(row.Review, review[:]) {
			return OperatorActionRow{}, ErrOperatorConflict
		}
		return row, tx.Commit()
	case !errors.Is(err, ErrOperatorMissing):
		return OperatorActionRow{}, err
	}
	row, err = scanOperatorAction(tx.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions
		WHERE episode_id = ? AND kind = ? AND review = ? AND state NOT IN ('expired', 'cancelled')
		ORDER BY terminal_at IS NULL DESC, created_at DESC LIMIT 1`, a.Episode, a.Kind, review[:]))
	if errors.Is(err, ErrOperatorMissing) {
		row, err = insertOperatorAction(ctx, tx, request, action, digest, review)
	}
	if err != nil {
		return OperatorActionRow{}, err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO operator_action_steps
		(kind, identity, action_id, phase, version, created_at) VALUES ('operation', ?, ?, 'observed', 0, ?)`,
		request.OperationID, row.ID, now); err != nil {
		return OperatorActionRow{}, err
	}
	return row, tx.Commit()
}

func insertOperatorAction(ctx context.Context, tx *sql.Tx, r OperatorRequest, action []byte, digest, review [32]byte) (OperatorActionRow, error) {
	var active, pending int
	a := r.Action
	if err := tx.QueryRowContext(ctx, `SELECT count(CASE WHEN episode_id = ? AND terminal_at IS NULL OR id = ? THEN 1 END),
		count(CASE WHEN state IN ('requested', 'confirmed') THEN 1 END) FROM operator_actions`, a.Episode, a.ID).Scan(&active, &pending); err != nil {
		return OperatorActionRow{}, err
	}
	switch {
	case active != 0:
		return OperatorActionRow{}, ErrOperatorConflict
	case pending >= 16:
		return OperatorActionRow{}, ErrOperatorCapacity
	}
	return scanOperatorAction(tx.QueryRowContext(ctx, `INSERT INTO operator_actions
		(id, episode_id, kind, review, digest, action, challenge, created_at, expires_at, state, generation)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 'requested', 1) RETURNING `+operatorColumns,
		a.ID, a.Episode, a.Kind, review[:], digest[:], action, r.Challenge[:], r.CreatedAt, r.ExpiresAt))
}

func expireOperatorActions(ctx context.Context, tx *sql.Tx, now int64) error {
	_, err := tx.ExecContext(ctx, `UPDATE operator_actions SET state = 'expired', terminal_at = expires_at
		WHERE state IN ('requested', 'confirmed') AND expires_at <= ?`, now)
	return err
}

// PruneOperatorActions drops old presentation and proof bytes while retaining
// the action link, all operation/effect identities and any uncertain reservation.
func (c *Credentials) PruneOperatorActions(ctx context.Context, now int64) error {
	if err := c.operatorAccess(ctx); err != nil {
		return err
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := expireOperatorActions(ctx, tx, now); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE operator_actions SET action = NULL, challenge = NULL, artifact = NULL, approved_by = NULL
		WHERE terminal_at < ? AND NOT EXISTS
			(SELECT 1 FROM operator_action_steps WHERE action_id = operator_actions.id AND phase = 'reserved')`, now-30*86400); err != nil {
		return err
	}
	return tx.Commit()
}
