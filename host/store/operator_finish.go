package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"regexp"
)

type OperatorFinishCheck struct {
	Name   string `json:"name"`
	Result string `json:"result"`
	Ref    string `json:"ref"`
	Entry  int64  `json:"entry"`
}

type OperatorFinishLocalCommit struct {
	State string `json:"state"`
	SHA   string `json:"sha,omitempty"`
}

const OperatorLocalCommitDone = "local-commit-done"

var operatorLocalOID = regexp.MustCompile(`^[0-9a-f]{40}$`)

type OperatorFinishJournal struct {
	Checks      []OperatorFinishCheck
	Commit      OperatorFinishLocalCommit
	CommitKey   OperatorStepKey
	CommitPhase string
	OpenCheck   bool
}

// ReadOperatorFinishJournal projects only typed check and commit receipts.
// A check reservation without an observation is never permission to replay.
func (c *Credentials) ReadOperatorFinishJournal(ctx context.Context, action OperatorCAS) (OperatorFinishJournal, error) {
	var out OperatorFinishJournal
	if err := c.operatorAccess(ctx); err != nil {
		return out, err
	}
	rows, err := c.db.QueryContext(ctx, `SELECT s.kind,s.identity,s.phase,s.version,s.owner,s.receipt
		FROM operator_action_steps s JOIN operator_actions a ON a.id=s.action_id
		WHERE a.id=? AND a.generation=? AND a.digest=? AND a.kind='finish' AND s.kind IN ('checks','commit')
		ORDER BY s.rowid`, action.ID, action.Generation, action.Digest[:])
	if err != nil {
		return out, err
	}
	defer rows.Close()
	for rows.Next() {
		key := OperatorStepKey{Action: action}
		var phase string
		var raw []byte
		if err := rows.Scan(&key.Kind, &key.Identity, &phase, &key.Version, &key.Owner, &raw); err != nil {
			return out, err
		}
		switch {
		case key.Kind == "commit":
			if out.CommitPhase != "" {
				return OperatorFinishJournal{}, ErrOperatorConflict
			}
			if phase == "observed" {
				if json.Unmarshal(raw, &out.Commit) != nil ||
					!(out.Commit.State == OperatorLocalCommitDone && operatorLocalOID.MatchString(out.Commit.SHA) ||
						out.Commit.State == "commit-failed" && out.Commit.SHA == "") {
					return OperatorFinishJournal{}, ErrOperatorConflict
				}
			}
			out.CommitKey, out.CommitPhase = key, phase
		case phase == "reserved":
			out.OpenCheck = true
		default:
			var receipt OperatorFinishCheck
			if json.Unmarshal(raw, &receipt) != nil || receipt.Name == "" || receipt.Entry < 1 || receipt.Ref == "" ||
				(receipt.Result != "passed" && receipt.Result != "failed") {
				return OperatorFinishJournal{}, ErrOperatorConflict
			}
			out.Checks = append(out.Checks, receipt)
		}
	}
	return out, rows.Err()
}

func (c *Credentials) NextOperatorFinish(ctx context.Context, now int64) (OperatorActionRow, error) {
	if err := c.operatorAccess(ctx); err != nil {
		return OperatorActionRow{}, err
	}
	return scanOperatorAction(c.db.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions
		WHERE kind='finish' AND (state='confirmed' AND created_at<=? AND expires_at>? OR state IN ('claimed','running'))
		AND NOT EXISTS (SELECT 1 FROM operator_action_steps s WHERE s.action_id=operator_actions.id AND s.kind='commit' AND s.phase='observed')
		ORDER BY state='confirmed', created_at,id LIMIT 1`, now, now))
}

// OperatorCompletedLaunches returns bounded candidate IDs. The runner checks
// each exact typed receipt before publishing an eligible episode.
func (c *Credentials) OperatorCompletedLaunches(ctx context.Context) ([]string, error) {
	if err := c.operatorAccess(ctx); err != nil {
		return nil, err
	}
	rows, err := c.db.QueryContext(ctx, `SELECT a.id FROM operator_actions a JOIN operator_action_steps s ON s.action_id=a.id
		WHERE a.kind='start' AND a.state='succeeded' AND s.kind='launch' AND s.phase='observed' AND s.version=4
		ORDER BY a.created_at DESC,a.id LIMIT 16`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// OperatorFinishLaunch reads the sole completed start for an episode, with
// the same structural/identity checks as launch-status.
func (c *Credentials) OperatorFinishLaunch(ctx context.Context, episode string) (OperatorLaunchSnapshot, error) {
	if err := c.operatorAccess(ctx); err != nil {
		return OperatorLaunchSnapshot{}, err
	}
	id := ""
	err := c.db.QueryRowContext(ctx, `SELECT a.id FROM operator_actions a JOIN operator_action_steps s ON s.action_id=a.id
		WHERE a.episode_id=? AND a.kind='start' AND a.state='succeeded' AND s.kind='launch' AND s.phase='observed' AND s.version=4
		ORDER BY a.created_at DESC,a.id LIMIT 1`, episode).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrOperatorMissing
	}
	if err != nil {
		return OperatorLaunchSnapshot{}, err
	}
	return c.ReadOperatorLaunch(ctx, id)
}
