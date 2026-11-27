package store

import (
	"context"
	"database/sql"

	"github.com/sunholo-data/ailang-world/host/approval"
)

type OperatorStartWindow string

const (
	OperatorStartUntouched     OperatorStartWindow = "claimed without effects"
	OperatorStartGitPending    OperatorStartWindow = "worktree reserved"
	OperatorStartTreeOnly      OperatorStartWindow = "worktree observed"
	OperatorStartMintPending   OperatorStartWindow = "mint reserved without commit"
	OperatorStartMintUncertain OperatorStartWindow = "mint commit uncertain"
	OperatorStartMintOnly      OperatorStartWindow = "mint observed"
	OperatorStartFilesPending  OperatorStartWindow = "files reserved"
	OperatorStartPrepared      OperatorStartWindow = "prepared-awaiting-S6"
)

type OperatorStartRecovery struct {
	CAS            OperatorCAS
	Owner, Episode string
	Window         OperatorStartWindow
	Worktree       string
	steps          [3]startRecoveryStep
}

type startRecoveryStep struct {
	Identity, Phase, Receipt string
	Version                  int64
}

// ReadOperatorStartRecovery classifies a single transactional snapshot. An
// owned window never proves the runner, worktree child or terminal is retired.
func (c *Credentials) ReadOperatorStartRecovery(ctx context.Context, id string) (OperatorStartRecovery, error) {
	if err := c.operatorAccess(ctx); err != nil {
		return OperatorStartRecovery{}, err
	}
	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return OperatorStartRecovery{}, err
	}
	defer tx.Rollback()
	current, _, err := readOperatorStartRecovery(ctx, tx, id)
	return current, err
}

func readOperatorStartRecovery(ctx context.Context, tx *sql.Tx, id string) (OperatorStartRecovery, OperatorActionRow, error) {
	row, err := scanOperatorAction(tx.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions WHERE id=?`, id))
	if err != nil {
		return OperatorStartRecovery{}, row, err
	}
	current := OperatorStartRecovery{CAS: OperatorCAS{row.ID, [32]byte(row.Digest), row.Generation}, Owner: row.Owner, Episode: row.Episode}
	refuse := func() (OperatorStartRecovery, OperatorActionRow, error) {
		return OperatorStartRecovery{}, row, ErrOperatorConflict
	}
	if row.Kind != "start" || row.Owner == "" {
		return refuse()
	}
	steps, err := tx.QueryContext(ctx, `SELECT kind, identity, phase, owner, version, receipt
		FROM operator_action_steps WHERE action_id=? AND kind<>'operation' LIMIT 4`, id)
	if err != nil {
		return OperatorStartRecovery{}, row, err
	}
	defer steps.Close()
	for steps.Next() {
		var step startRecoveryStep
		var kind, owner string
		var receipt []byte
		if err := steps.Scan(&kind, &step.Identity, &step.Phase, &owner, &step.Version, &receipt); err != nil {
			return OperatorStartRecovery{}, row, err
		}
		var i int
		switch kind {
		case "worktree":
			i = 0
		case "mint":
			i = 1
		case "files":
			i = 2
		default:
			return refuse()
		}
		if owner != row.Owner || current.steps[i].Identity != "" || step.Identity == "" || step.Version < 1 {
			return refuse()
		}
		step.Receipt = string(receipt)
		current.steps[i] = step
	}
	if err := steps.Err(); err != nil {
		return OperatorStartRecovery{}, row, err
	}
	tree, mint, files := current.steps[0], current.steps[1], current.steps[2]
	var start, session bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM episode_starts WHERE episode=?),
		EXISTS(SELECT 1 FROM session_credentials WHERE episode_id=?)`,
		row.Episode, row.Episode).Scan(&start, &session)
	if err != nil {
		return OperatorStartRecovery{}, row, err
	}

	current.Worktree = tree.Identity
	switch {
	case row.State == "claimed" && tree.Phase == "" && mint.Phase == "" && files.Phase == "" && !start && !session:
		current.Window = OperatorStartUntouched
	case row.State != "running" && row.State != "uncertain":
		return refuse()
	case tree.Phase == "reserved" && mint.Phase == "" && files.Phase == "" && !start && !session:
		current.Window = OperatorStartGitPending
	case tree.Phase != "observed" || tree.Receipt != "created":
		return refuse()
	case mint.Phase == "" && files.Phase == "" && !start && !session:
		current.Window = OperatorStartTreeOnly
	case mint.Phase == "reserved" && mint.Identity == id && files.Phase == "":
		// CommitOperatorMint makes the start, session and observation one
		// transaction. No subset of those three is proof of rollback.
		if start || session {
			current.Window = OperatorStartMintUncertain
		} else {
			current.Window = OperatorStartMintPending
		}
	case mint.Phase == "observed" && mint.Identity == id:
		if _, _, err := readOperatorFiles(ctx, tx, id); err != nil {
			return refuse()
		}
		switch {
		case files.Phase == "":
			current.Window = OperatorStartMintOnly
		case files.Phase == "reserved":
			current.Window = OperatorStartFilesPending
		case files.Phase == "observed" && files.Receipt == OperatorFilesPublished:
			current.Window = OperatorStartPrepared
		default:
			return refuse()
		}
	default:
		return refuse()
	}
	return current, row, nil
}

// DisposeOperatorStartGap applies either a fresh release of an untouched
// claim or an attended stop. The caller proves process quiescence under the
// episode lock and inspects any retained worktree before calling.
func (c *Credentials) DisposeOperatorStartGap(ctx context.Context, reviewed OperatorStartRecovery, release *approval.OperatorMintProof, now int64) error {
	if err := c.operatorAccess(ctx); err != nil {
		return err
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	current, row, err := readOperatorStartRecovery(ctx, tx, reviewed.CAS.ID)
	if err != nil {
		return err
	}
	if current != reviewed {
		return ErrOperatorConflict
	}
	a, err := decodeOperatorMintAction(row)
	if err != nil || c.lock == nil || a.Store != c.lock.dbPath {
		return ErrOperatorConflict
	}
	switch current.Window {
	case OperatorStartUntouched, OperatorStartTreeOnly, OperatorStartMintPending:
	default:
		return ErrOperatorConflict
	}
	state, terminalAt := "failed", any(now)
	if release != nil {
		if current.Window != OperatorStartUntouched {
			return ErrOperatorConflict
		}
		if now >= row.ExpiresAt {
			return ErrOperatorExpired
		}
		binding, err := release.Bind(a, [32]byte(row.Challenge), row.ExpiresAt, row.Artifact)
		if err != nil || row.ApprovedBy != binding.Format+":"+binding.Name {
			return ErrOperatorConflict
		}
		if err := operatorMintApprover(ctx, tx, release, binding, now); err != nil {
			return err
		}
		state, terminalAt = "confirmed", nil
	}
	result, err := tx.ExecContext(ctx, `UPDATE operator_actions SET state=?, terminal_at=?, owner=NULL, generation=generation+1
		WHERE id=? AND digest=? AND generation=? AND owner=?`, state, terminalAt,
		current.CAS.ID, current.CAS.Digest[:], current.CAS.Generation, current.Owner)
	if err := operatorChangedOne(result, err); err != nil {
		return err
	}
	return tx.Commit()
}

// NextOwnedOperatorStart scans only owned starts. A cursor lets the runner
// inspect one action per tick without starving behind a prepared checkpoint.
func (c *Credentials) NextOwnedOperatorStart(ctx context.Context, after string) (OperatorActionRow, error) {
	if err := c.operatorAccess(ctx); err != nil {
		return OperatorActionRow{}, err
	}
	return scanOperatorAction(c.db.QueryRowContext(ctx, `SELECT `+operatorColumns+
		` FROM operator_actions a WHERE kind='start' AND (owner IS NOT NULL OR (`+operatorLaunchUnsettled+`)) AND id>? ORDER BY id LIMIT 1`, after))
}

// MarkOperatorMintUncertain preserves every receipt and the old owner while
// recording the ambiguity found after an interrupted mint.
func (c *Credentials) MarkOperatorMintUncertain(ctx context.Context, reviewed OperatorStartRecovery) error {
	if err := c.operatorAccess(ctx); err != nil {
		return err
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	current, row, err := readOperatorStartRecovery(ctx, tx, reviewed.CAS.ID)
	if err != nil || current != reviewed || current.Window != OperatorStartMintUncertain {
		return ErrOperatorConflict
	}
	if row.State == "uncertain" {
		return nil
	}
	res, err := tx.ExecContext(ctx, `UPDATE operator_actions SET state='uncertain', generation=generation+1
		WHERE id=? AND digest=? AND generation=? AND owner=? AND state='running'`,
		current.CAS.ID, current.CAS.Digest[:], current.CAS.Generation, current.Owner)
	if err := operatorChangedOne(res, err); err != nil {
		return err
	}
	return tx.Commit()
}
