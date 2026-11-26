package store

import "context"

// EpisodeHasUnresolvedOperatorSteps reads the durable effect ledger without a
// writer lock. An observation alone does not dispose of an effect. A reserved
// row still protects the worktree even if its action was marked terminal.
func EpisodeHasUnresolvedOperatorSteps(ctx context.Context, path, episode string) (bool, error) {
	if err := requireDeadline(ctx); err != nil {
		return false, err
	}
	db, version, err := openOperatorReadOnly(ctx, path)
	if err != nil {
		return false, err
	}
	defer db.Close()
	if version < operatorSchemaVersion {
		return false, &OlderCredentialsStoreError{Path: path, Version: version, Need: operatorSchemaVersion, Feature: "operator actions"}
	}
	var unresolved bool
	err = db.QueryRowContext(ctx, `SELECT EXISTS (
		SELECT 1 FROM operator_actions a JOIN operator_action_steps s ON s.action_id = a.id
		WHERE a.episode_id = ? AND s.kind <> 'operation'
		AND (s.phase = 'reserved' OR a.terminal_at IS NULL))`, episode).Scan(&unresolved)
	return unresolved, err
}
