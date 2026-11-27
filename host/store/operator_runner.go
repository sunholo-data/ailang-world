package store

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"os"
)

const OperatorRunnerProtocol = 1
const operatorSchemaVersion = 10

type OperatorHeartbeat struct {
	Owner                       OperatorOwner
	Protocol                    int
	Build, Readiness            string
	SeenAt, CatalogueGeneration int64
	Catalogue                   []byte
}

// Available applies the S4 read gate. An unavailable heartbeat leaves
// both the durable lease and all effect ownership exactly as recorded.
func (h OperatorHeartbeat) Available(build, kind string, now int64) bool {
	return h.Protocol == OperatorRunnerProtocol && h.Build == build &&
		h.Owner.Generation > 0 && h.Owner.Nonce != "" &&
		h.SeenAt <= now && h.SeenAt > now-10 &&
		(h.Readiness == "both" || h.Readiness == kind)
}

func (c *Credentials) OperatorHeartbeat(ctx context.Context) (h OperatorHeartbeat, err error) {
	if err := c.operatorAccess(ctx); err != nil {
		return h, err
	}
	err = c.db.QueryRowContext(ctx, `SELECT owner, generation, protocol, build, readiness, seen_at,
		catalogue_generation, catalogue FROM operator_runner_state WHERE singleton = 1`).
		Scan(&h.Owner.Nonce, &h.Owner.Generation, &h.Protocol, &h.Build, &h.Readiness, &h.SeenAt, &h.CatalogueGeneration, &h.Catalogue)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrOperatorRunner
	}
	return h, err
}

// PutOperatorHeartbeat publishes a complete generation in one statement.
// A new owner (generation zero) must hold the kernel runner lock for this store.
func (c *Credentials) PutOperatorHeartbeat(ctx context.Context, h OperatorHeartbeat) (int64, error) {
	if err := c.operatorAccess(ctx); err != nil {
		return 0, err
	}
	var generation int64
	if h.Owner.Generation == 0 {
		err := c.db.QueryRowContext(ctx, `INSERT INTO operator_runner_state
			(singleton,generation,owner,protocol,build,readiness,seen_at,catalogue_generation,catalogue)
			VALUES(1,1,?,?,?,?,?,1,?)
			ON CONFLICT(singleton) DO UPDATE SET generation = generation+1, owner = excluded.owner,
			protocol = excluded.protocol, build = excluded.build, readiness = excluded.readiness,
			seen_at = excluded.seen_at, catalogue_generation = catalogue_generation+1,
			catalogue = excluded.catalogue RETURNING generation`,
			h.Owner.Nonce, h.Protocol, h.Build, h.Readiness, h.SeenAt, h.Catalogue).Scan(&generation)
		return generation, err
	}
	// Refresh liveness without superseding a review of the same catalogue.
	err := c.db.QueryRowContext(ctx, `UPDATE operator_runner_state SET readiness=?, seen_at=?,
		catalogue_generation=catalogue_generation+CASE WHEN readiness IS ? AND catalogue IS ? THEN 0 ELSE 1 END,
        catalogue=? WHERE singleton=1 AND owner=? AND generation=? AND protocol=? AND build=? AND seen_at<=?
		RETURNING generation`, h.Readiness, h.SeenAt, h.Readiness, h.Catalogue, h.Catalogue, h.Owner.Nonce,
		h.Owner.Generation, h.Protocol, h.Build, h.SeenAt).Scan(&generation)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrOperatorRunner
	}
	return generation, err
}

// CancelUnstartedOperatorAction never removes or releases a reserved effect.
func (c *Credentials) CancelUnstartedOperatorAction(ctx context.Context, id string, now int64) error {
	if err := c.operatorAccess(ctx); err != nil {
		return err
	}
	result, err := c.db.ExecContext(ctx, `UPDATE operator_actions SET state='cancelled', terminal_at=?,
		generation=generation+1 WHERE id=? AND state IN ('requested','confirmed')
		AND owner IS NULL AND expires_at>? AND NOT EXISTS (SELECT 1 FROM operator_action_steps
			WHERE action_id=? AND kind<>'operation')`, now, id, now, id)
	return operatorChangedOne(result, err)
}

// ReadOperatorEpisodeOwner reads exclusion independently of lease liveness.
// The caller holds the episode lock or has just received a non-waiting
// contention error. A missing store or older schema has no operator rows.
func ReadOperatorEpisodeOwner(ctx context.Context, path, episode string) (string, error) {
	if err := requireDeadline(ctx); err != nil {
		return "", err
	}
	db, version, err := openOperatorReadOnly(ctx, path)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	defer db.Close()
	if version < operatorSchemaVersion {
		return "", nil
	}
	var id string
	err = db.QueryRowContext(ctx, `SELECT id FROM operator_actions a WHERE episode_id=? AND (owner IS NOT NULL OR (`+operatorLaunchUnsettled+`)) LIMIT 1`, episode).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return id, err
}

// ReadOperatorHeartbeat bounds opening and reading by the same context.
func ReadOperatorHeartbeat(ctx context.Context, path string) (OperatorHeartbeat, error) {
	if err := requireDeadline(ctx); err != nil {
		return OperatorHeartbeat{}, err
	}
	db, version, err := openOperatorReadOnly(ctx, path)
	if err != nil {
		return OperatorHeartbeat{}, err
	}
	defer db.Close()
	c := Credentials{db: db, version: version, display: path}
	return c.OperatorHeartbeat(ctx)
}

func openOperatorReadOnly(ctx context.Context, path string) (*sql.DB, int, error) {
	canonical, _, err := resolveDSN(path)
	if err != nil {
		return nil, 0, err
	}
	st, err := os.Stat(canonical)
	if err != nil {
		return nil, 0, err
	}
	if !st.Mode().IsRegular() {
		return nil, 0, errors.New("expected a regular operator credentials file")
	}
	db, err := sql.Open("sqlite", readOnlyDSN(canonical, url.Values{"_pragma": {"busy_timeout(0)"}}))
	if err != nil {
		return nil, 0, err
	}
	var version int
	if err := db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
		db.Close()
		return nil, 0, err
	}
	if version > credentialsSchemaVersion {
		db.Close()
		return nil, 0, &CredentialsSchemaError{Path: canonical, Detail: "unsupported operator credentials schema"}
	}
	return db, version, nil
}
