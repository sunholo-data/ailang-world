package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/sunholo-data/ailang-world/host/approval"
)

type OperatorLaunchReceipt struct {
	Run, Executable string
	PID             int
	Birth           string
	Status          string
	Exit            int
}

// An observed start with no terminal observation still excludes effects.
const operatorLaunchUnsettled = `a.state='succeeded' AND EXISTS (
	SELECT 1 FROM operator_action_steps ls WHERE ls.action_id=a.id
	AND ls.kind='launch' AND ls.phase='observed' AND ls.version=3)`

func operatorLaunchDigest(s string) bool {
	value, ok := strings.CutPrefix(s, "sha256:")
	return ok && operatorDigestString(value)
}

type OperatorLaunchSnapshot struct {
	Key     OperatorStepKey
	Phase   string
	Receipt OperatorLaunchReceipt
	NoSpawn bool
}

func operatorLaunchKey(s OperatorStartRecovery, handoff [32]byte) OperatorStepKey {
	return OperatorStepKey{Action: s.CAS, Owner: s.Owner, Kind: "launch", Identity: fmt.Sprintf("%s:%x", s.CAS.ID, sha256.Sum256(handoff[:])), Version: 1}
}

// ReserveOperatorLaunch is called under the episode and lifecycle locks. It
// consumes the exact prepared snapshot and rechecks the enrolled proof and
// current session inside the transaction that makes the reservation durable.
func (c *Credentials) ReserveOperatorLaunch(ctx context.Context, reviewed OperatorStartRecovery, handoff [32]byte, proof *approval.OperatorMintProof, clock func() int64) error {
	if err := c.operatorAccess(ctx); err != nil {
		return err
	}
	if clock == nil || c.lock == nil || handoff == ([32]byte{}) {
		return ErrOperatorInput
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	current, row, err := readOperatorStartRecovery(ctx, tx, reviewed.CAS.ID)
	now := clock()
	if err == nil && now >= row.ExpiresAt {
		return ErrOperatorExpired
	}
	if err != nil || current != reviewed || current.Window != OperatorStartPrepared {
		return ErrOperatorConflict
	}
	if _, _, err := operatorSpawnInput(ctx, tx, row, proof, now, false); err != nil {
		return err
	}
	if err := reserveOperatorStep(ctx, tx, operatorLaunchKey(current, handoff), now); err != nil {
		return err
	}
	return tx.Commit()
}

// ReadOperatorSpawnInput reads the frozen start and its exact credential for
// first dispatch. It leaves all persisted evidence as it found it.
func (c *Credentials) ReadOperatorSpawnInput(ctx context.Context, id string) (OperatorFileState, approval.OperatorAction, error) {
	if err := c.operatorAccess(ctx); err != nil {
		return OperatorFileState{}, approval.OperatorAction{}, err
	}
	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return OperatorFileState{}, approval.OperatorAction{}, err
	}
	defer tx.Rollback()
	s, _, err := readOperatorLaunch(ctx, tx, id)
	if err != nil || s.Phase != "reserved" || s.Key.Version != 1 {
		return OperatorFileState{}, approval.OperatorAction{}, ErrOperatorConflict
	}
	return readOperatorFilesAt(ctx, tx, id, true)
}

// AuthorizeOperatorSpawn commits one dispatch grant. The process that receives
// the private handoff must hold the credentials writer through cmd.Start.
func (c *Credentials) AuthorizeOperatorSpawn(ctx context.Context, id string, handoff [32]byte, proof *approval.OperatorMintProof, clock func() int64) error {
	if err := c.operatorAccess(ctx); err != nil {
		return err
	}
	if clock == nil || c.lock == nil || handoff == ([32]byte{}) {
		return ErrOperatorInput
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	s, row, err := readOperatorLaunch(ctx, tx, id)
	if err != nil || s.Phase != "reserved" || s.Key != operatorLaunchKey(OperatorStartRecovery{CAS: s.Key.Action, Owner: row.Owner}, handoff) {
		return ErrOperatorConflict
	}
	if _, _, err := operatorSpawnInput(ctx, tx, row, proof, clock(), true); err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx, `UPDATE operator_action_steps SET version=version+1
		WHERE action_id=? AND kind='launch' AND identity=? AND phase='reserved' AND version=1 AND owner=?`,
		id, s.Key.Identity, s.Key.Owner)
	if err := operatorChangedOne(res, err); err != nil {
		return err
	}
	return tx.Commit()
}

func operatorSpawnInput(ctx context.Context, tx *sql.Tx, row OperatorActionRow, proof *approval.OperatorMintProof, now int64, reserved bool) (OperatorFileState, approval.OperatorAction, error) {
	files, a, err := readOperatorFilesAt(ctx, tx, row.ID, reserved)
	if err != nil {
		return files, a, ErrOperatorConflict
	}
	if err := admitOperatorFiles(ctx, tx, row, files, a, proof, now); err != nil {
		return files, a, err
	}
	var other bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(
		SELECT 1 FROM session_credentials WHERE episode_id=? AND credential_id<>?)`,
		row.Episode, files.Session.CredentialID).Scan(&other); err != nil {
		return files, a, err
	}
	if other || !operatorLaunchDigest(mintFact(a, "launch_binary_sha256")) {
		return files, a, ErrOperatorConflict
	}
	return files, a, nil
}

// ReadOperatorLaunch uses one read-only snapshot and never makes a replay
// decision from a PID, the heartbeat, or the current presence of a credential.
func (c *Credentials) ReadOperatorLaunch(ctx context.Context, id string) (OperatorLaunchSnapshot, error) {
	if err := c.operatorAccess(ctx); err != nil {
		return OperatorLaunchSnapshot{}, err
	}
	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return OperatorLaunchSnapshot{}, err
	}
	defer tx.Rollback()
	s, _, err := readOperatorLaunch(ctx, tx, id)
	return s, err
}

func readOperatorLaunch(ctx context.Context, tx *sql.Tx, id string) (OperatorLaunchSnapshot, OperatorActionRow, error) {
	var s OperatorLaunchSnapshot
	row, err := scanOperatorAction(tx.QueryRowContext(ctx, "SELECT "+operatorColumns+" FROM operator_actions WHERE id=?", id))
	if err != nil {
		return s, row, err
	}
	s.Key = OperatorStepKey{Action: OperatorCAS{row.ID, [32]byte(row.Digest), row.Generation}, Kind: "launch"}
	var n int
	var raw []byte
	err = tx.QueryRowContext(ctx, `SELECT identity,phase,version,owner,receipt,
		(SELECT count(*) FROM operator_action_steps WHERE action_id=? AND kind='launch')
		FROM operator_action_steps WHERE action_id=? AND kind='launch' LIMIT 1`, id, id).
		Scan(&s.Key.Identity, &s.Phase, &s.Key.Version, &s.Key.Owner, &raw, &n)
	if errors.Is(err, sql.ErrNoRows) && row.Kind == "start" {
		return s, row, ErrOperatorMissing
	}
	identity, prefixed := strings.CutPrefix(s.Key.Identity, id+":")
	if err != nil || n != 1 || row.Kind != "start" || s.Key.Owner == "" || !prefixed || !operatorDigestString(identity) {
		return s, row, ErrOperatorConflict
	}
	if s.Phase == "reserved" && (s.Key.Version == 1 || s.Key.Version == 2) && len(raw) == 0 && s.Key.Owner == row.Owner &&
		(row.State == "running" || row.State == "uncertain") {
		return s, row, nil
	}
	if s.Phase != "observed" || json.Unmarshal(raw, &s.Receipt) != nil || s.Receipt.Run != id ||
		!operatorLaunchDigest(s.Receipt.Executable) {
		return s, row, ErrOperatorConflict
	}
	if len(row.Action) != 0 {
		a, err := decodeOperatorMintAction(row)
		if err != nil || s.Receipt.Executable != mintFact(a, "launch_binary_sha256") {
			return s, row, ErrOperatorConflict
		}
	}
	switch s.Receipt.Status {
	case "launched":
		if s.Key.Version == 3 && row.State == "succeeded" && s.Receipt.PID > 0 && s.Receipt.Birth != "" {
			return s, row, nil
		}
	case "spawn-refused":
		if (s.Key.Version == 2 || s.Key.Version == 3) && row.State == "failed" && s.Receipt.PID == 0 && s.Receipt.Birth == "" {
			s.Receipt, s.NoSpawn = OperatorLaunchReceipt{}, true
			return s, row, nil
		}
	case "stopped", "failed", "capped":
		if s.Key.Version == 4 && row.State == "succeeded" && s.Receipt.PID > 0 && s.Receipt.Birth != "" {
			return s, row, nil
		}
	}
	return s, row, ErrOperatorConflict
}

// RecordOperatorLaunch compares the exact read-back and never replaces a
// process identity. Recording an admitted child's outcome ignores expiry.
func (c *Credentials) RecordOperatorLaunch(ctx context.Context, reviewed OperatorLaunchSnapshot, r OperatorLaunchReceipt, history *EpisodeLaunch, now int64) error {
	if err := c.operatorAccess(ctx); err != nil {
		return err
	}
	if c.lock == nil {
		return ErrOperatorInput
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	s, row, err := readOperatorLaunch(ctx, tx, reviewed.Key.Action.ID)
	a, decodeErr := decodeOperatorMintAction(row)
	if err != nil || decodeErr != nil || s != reviewed || r.Run != row.ID || r.Executable != mintFact(a, "launch_binary_sha256") {
		return ErrOperatorConflict
	}
	var terminal string
	switch s.Phase {
	case "reserved":
		switch {
		case s.Key.Version == 2 && r.Status == "launched" && r.PID > 0 && r.Birth != "" && history == nil:
			terminal = "succeeded"
		case r.Status == "spawn-refused" && r.PID == 0 && r.Birth == "" && history == nil:
			terminal = "failed"
		default:
			return ErrOperatorInput
		}
	case "observed":
		if s.Receipt.Status != "launched" || r.PID != s.Receipt.PID || r.Birth != s.Receipt.Birth ||
			history == nil || history.Mode != "headless" || history.Exit != r.Exit {
			return ErrOperatorConflict
		}
		switch r.Status {
		case "stopped", "failed", "capped":
		default:
			return ErrOperatorInput
		}
	}
	raw, err := json.Marshal(r)
	if err != nil || len(raw) > 2048 {
		return ErrOperatorInput
	}
	change, err := tx.ExecContext(ctx, `UPDATE operator_action_steps SET phase='observed',version=version+1,receipt=?
		WHERE action_id=? AND kind='launch' AND identity=? AND phase=? AND version=? AND owner=?`,
		raw, row.ID, s.Key.Identity, s.Phase, s.Key.Version, s.Key.Owner)
	if err := operatorChangedOne(change, err); err != nil {
		return err
	}
	if terminal != "" {
		change, err = tx.ExecContext(ctx, `UPDATE operator_actions SET state=?,owner=NULL,terminal_at=?,generation=generation+1
			WHERE id=? AND digest=? AND generation=? AND owner=? AND state='running'`,
			terminal, now, row.ID, row.Digest, row.Generation, s.Key.Owner)
		if err := operatorChangedOne(change, err); err != nil {
			return err
		}
	}
	if history != nil {
		if err := appendEpisodeLaunch(ctx, tx, row.Episode, *history); err != nil {
			return err
		}
	}
	return tx.Commit()
}
