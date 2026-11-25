package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const allowanceSchemaVersion = 11

// Sentinel errors for standing allowance operations (row 224).
var (
	ErrAllowanceExists   = errors.New("store: an active standing allowance already exists for this preset")
	ErrAllowanceNotFound = errors.New("store: standing allowance not found")
	ErrAllowanceRevoked  = errors.New("store: standing allowance is revoked")
	ErrAllowanceDailyCap = errors.New("store: standing allowance daily cap reached")
)

// StandingAllowanceRow is one persisted row in standing_allowances (schema v11).
type StandingAllowanceRow struct {
	AllowanceID    string
	PresetName     string
	PresetSHA256   string
	ProfileProject string
	ProfileDigest  string
	Payload        []byte
	Digest         string
	Challenge      string
	MaxPerDay      int
	TTLS           int64
	ApprovedAt     int64
	ExpiresAt      int64
	ApprovedBy     string
	Artifact       []byte
	RevokedAt      *int64
	RevokedBy      string
	RevokedReason  string
}

// AllowanceDrawRow is one persisted row in allowance_draws (schema v11).
type AllowanceDrawRow struct {
	DrawID       string
	AllowanceID  string
	DayKey       string
	EpisodeID    string
	CredentialID string
	Repo         string
	DrawnAt      int64
	TTLS         int64
}

func (c *Credentials) requireAllowanceSchema() error {
	if c.version < allowanceSchemaVersion {
		return &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: allowanceSchemaVersion, Feature: "standing allowances"}
	}
	return nil
}

func (c *Credentials) requireAllowanceWrite(op string) error {
	if c.lock == nil {
		return fmt.Errorf("store: %s: credentials store %q is open read-only", op, c.display)
	}
	return c.requireAllowanceSchema()
}

// GrantAllowance inserts a new standing allowance into standing_allowances.
func (c *Credentials) GrantAllowance(ctx context.Context, r StandingAllowanceRow) error {
	if err := c.requireAllowanceWrite("GrantAllowance"); err != nil {
		return err
	}
	var proj, pdig any
	if r.ProfileProject != "" {
		proj, pdig = r.ProfileProject, r.ProfileDigest
	}
	var artifact any
	if len(r.Artifact) > 0 {
		artifact = r.Artifact
	}
	_, err := c.db.ExecContext(ctx, `
		INSERT INTO standing_allowances (
			allowance_id, preset_name, preset_sha256, profile_project, profile_digest,
			payload, digest, challenge, max_per_day, ttl_s,
			approved_at, expires_at, approved_by, artifact
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.AllowanceID, r.PresetName, r.PresetSHA256, proj, pdig,
		r.Payload, r.Digest, r.Challenge, r.MaxPerDay, r.TTLS,
		r.ApprovedAt, r.ExpiresAt, r.ApprovedBy, artifact,
	)
	if err != nil {
		if strings.Contains(err.Error(), "standing_allowances_active_preset") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed: standing_allowances.preset_name") {
			return fmt.Errorf("%w: %s", ErrAllowanceExists, r.PresetName)
		}
		return fmt.Errorf("store: grant allowance %s: %w", r.AllowanceID, err)
	}
	return nil
}

func scanStandingAllowance(s interface{ Scan(dest ...any) error }) (StandingAllowanceRow, error) {
	var r StandingAllowanceRow
	var proj, pdig, revBy, revReason sql.NullString
	var revAt sql.NullInt64
	var artifact []byte
	if err := s.Scan(
		&r.AllowanceID, &r.PresetName, &r.PresetSHA256, &proj, &pdig,
		&r.Payload, &r.Digest, &r.Challenge, &r.MaxPerDay, &r.TTLS,
		&r.ApprovedAt, &r.ExpiresAt, &r.ApprovedBy, &artifact,
		&revAt, &revBy, &revReason,
	); err != nil {
		return StandingAllowanceRow{}, err
	}
	r.ProfileProject = proj.String
	r.ProfileDigest = pdig.String
	r.Artifact = artifact
	if revAt.Valid {
		v := revAt.Int64
		r.RevokedAt = &v
	}
	r.RevokedBy = revBy.String
	r.RevokedReason = revReason.String
	return r, nil
}

const standingAllowanceCols = `
	allowance_id, preset_name, preset_sha256, profile_project, profile_digest,
	payload, digest, challenge, max_per_day, ttl_s,
	approved_at, expires_at, approved_by, artifact,
	revoked_at, revoked_by, revoked_reason`

// GetActiveAllowance returns the unrevoked standing allowance for presetName, if any.
func (c *Credentials) GetActiveAllowance(ctx context.Context, presetName string) (StandingAllowanceRow, bool, error) {
	if err := c.requireAllowanceSchema(); err != nil {
		return StandingAllowanceRow{}, false, err
	}
	row := c.db.QueryRowContext(ctx,
		`SELECT `+standingAllowanceCols+` FROM standing_allowances WHERE preset_name = ? AND revoked_at IS NULL`,
		presetName,
	)
	got, err := scanStandingAllowance(row)
	if errors.Is(err, sql.ErrNoRows) {
		return StandingAllowanceRow{}, false, nil
	}
	if err != nil {
		return StandingAllowanceRow{}, false, fmt.Errorf("store: get active allowance %q: %w", presetName, err)
	}
	return got, true, nil
}

// GetAllowance returns a standing allowance by its allowance_id.
func (c *Credentials) GetAllowance(ctx context.Context, allowanceID string) (StandingAllowanceRow, bool, error) {
	if err := c.requireAllowanceSchema(); err != nil {
		return StandingAllowanceRow{}, false, err
	}
	row := c.db.QueryRowContext(ctx,
		`SELECT `+standingAllowanceCols+` FROM standing_allowances WHERE allowance_id = ?`,
		allowanceID,
	)
	got, err := scanStandingAllowance(row)
	if errors.Is(err, sql.ErrNoRows) {
		return StandingAllowanceRow{}, false, nil
	}
	if err != nil {
		return StandingAllowanceRow{}, false, fmt.Errorf("store: get allowance %q: %w", allowanceID, err)
	}
	return got, true, nil
}

// ListAllowances returns all standing allowances, newest first.
func (c *Credentials) ListAllowances(ctx context.Context, includeRevoked bool) ([]StandingAllowanceRow, error) {
	if err := c.requireAllowanceSchema(); err != nil {
		return nil, err
	}
	q := `SELECT ` + standingAllowanceCols + ` FROM standing_allowances`
	if !includeRevoked {
		q += ` WHERE revoked_at IS NULL`
	}
	q += ` ORDER BY approved_at DESC, allowance_id ASC`
	rows, err := c.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("store: list allowances: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []StandingAllowanceRow
	for rows.Next() {
		r, err := scanStandingAllowance(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// RevokeAllowance marks an active allowance revoked (matching allowance_id or preset_name).
func (c *Credentials) RevokeAllowance(ctx context.Context, idOrPreset string, revokedAt int64, revokedBy, reason string) (StandingAllowanceRow, error) {
	if err := c.requireAllowanceWrite("RevokeAllowance"); err != nil {
		return StandingAllowanceRow{}, err
	}
	var target StandingAllowanceRow
	var ok bool
	var err error
	if strings.HasPrefix(idOrPreset, "sa-") {
		target, ok, err = c.GetAllowance(ctx, idOrPreset)
	} else {
		target, ok, err = c.GetActiveAllowance(ctx, idOrPreset)
	}
	if err != nil {
		return StandingAllowanceRow{}, err
	}
	if !ok {
		return StandingAllowanceRow{}, fmt.Errorf("%w: %s", ErrAllowanceNotFound, idOrPreset)
	}
	if target.RevokedAt != nil {
		return target, fmt.Errorf("%w: %s", ErrAllowanceRevoked, target.AllowanceID)
	}
	var reasonVal any
	if reason != "" {
		reasonVal = reason
	}
	if _, err := c.db.ExecContext(ctx, `
		UPDATE standing_allowances
		SET revoked_at = ?, revoked_by = ?, revoked_reason = ?
		WHERE allowance_id = ? AND revoked_at IS NULL`,
		revokedAt, revokedBy, reasonVal, target.AllowanceID,
	); err != nil {
		return StandingAllowanceRow{}, fmt.Errorf("store: revoke allowance %s: %w", target.AllowanceID, err)
	}
	target.RevokedAt = &revokedAt
	target.RevokedBy = revokedBy
	target.RevokedReason = reason
	return target, nil
}

// CountAllowanceDraws returns how many sessions were drawn from allowanceID on dayKey (UTC YYYY-MM-DD).
func (c *Credentials) CountAllowanceDraws(ctx context.Context, allowanceID, dayKey string) (int, error) {
	if err := c.requireAllowanceSchema(); err != nil {
		return 0, err
	}
	var n int
	if err := c.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM allowance_draws WHERE allowance_id = ? AND day_key = ?`,
		allowanceID, dayKey,
	).Scan(&n); err != nil {
		return 0, fmt.Errorf("store: count allowance draws: %w", err)
	}
	return n, nil
}

// MintSessionFromAllowance atomically verifies the allowance is still active and
// below maxPerDay for draw.DayKey, then inserts session_credentials,
// episode_starts (with Allowance provenance), and allowance_draws in one
// BEGIN IMMEDIATE transaction. Returns the 1-based draw number for the day.
func (c *Credentials) MintSessionFromAllowance(
	ctx context.Context,
	session SessionRow,
	start EpisodeStart,
	draw AllowanceDrawRow,
	maxPerDay int,
) (int, error) {
	if err := c.requireAllowanceWrite("MintSessionFromAllowance"); err != nil {
		return 0, err
	}
	conn, err := c.db.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("store: mint from allowance: conn: %w", err)
	}
	defer func() { _ = conn.Close() }()
	if _, err := conn.ExecContext(ctx, "BEGIN IMMEDIATE"); err != nil {
		return 0, fmt.Errorf("store: mint from allowance: begin: %w", err)
	}
	rollback := true
	defer func() {
		if rollback {
			_, _ = conn.ExecContext(ctx, "ROLLBACK")
		}
	}()

	var revAt sql.NullInt64
	var approvedBy string
	var storedMax int
	if err := conn.QueryRowContext(ctx,
		`SELECT revoked_at, approved_by, max_per_day FROM standing_allowances WHERE allowance_id = ?`,
		draw.AllowanceID,
	).Scan(&revAt, &approvedBy, &storedMax); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%w: %s", ErrAllowanceNotFound, draw.AllowanceID)
		}
		return 0, fmt.Errorf("store: mint from allowance: read allowance: %w", err)
	}
	if revAt.Valid {
		return 0, fmt.Errorf("%w: %s", ErrAllowanceRevoked, draw.AllowanceID)
	}
	capLimit := storedMax
	if maxPerDay > 0 && maxPerDay < capLimit {
		capLimit = maxPerDay
	}
	var drawnToday int
	if err := conn.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM allowance_draws WHERE allowance_id = ? AND day_key = ?`,
		draw.AllowanceID, draw.DayKey,
	).Scan(&drawnToday); err != nil {
		return 0, fmt.Errorf("store: mint from allowance: count draws: %w", err)
	}
	if drawnToday >= capLimit {
		return 0, fmt.Errorf("%w: %d/%d on %s", ErrAllowanceDailyCap, drawnToday, capLimit, draw.DayKey)
	}

	drawNum := drawnToday + 1
	start.Allowance = &StartAllowance{
		AllowanceID: draw.AllowanceID,
		DrawID:      draw.DrawID,
		ApprovedBy:  approvedBy,
		DayKey:      draw.DayKey,
		DrawNumber:  drawNum,
	}
	startJSON, err := json.Marshal(start)
	if err != nil {
		return 0, fmt.Errorf("store: mint from allowance: encode start: %w", err)
	}

	var proj, pdig any
	if session.Project != "" {
		proj, pdig = session.Project, session.ProfileDigest
	}
	if _, err := conn.ExecContext(ctx, `
		INSERT INTO session_credentials (
			credential_id, episode_id, grants_json, expires_at, created_at, project, profile_digest
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		session.CredentialID, session.EpisodeID, session.GrantsJSON,
		session.ExpiresAt, session.CreatedAt, proj, pdig,
	); err != nil {
		return 0, fmt.Errorf("store: mint from allowance: insert session: %w", err)
	}
	if _, err := conn.ExecContext(ctx,
		`INSERT INTO episode_starts (episode, start_json) VALUES (?, ?)`,
		start.Episode, string(startJSON),
	); err != nil {
		return 0, fmt.Errorf("store: mint from allowance: insert start: %w", err)
	}
	if _, err := conn.ExecContext(ctx, `
		INSERT INTO allowance_draws (
			draw_id, allowance_id, day_key, episode_id, credential_id, repo, drawn_at, ttl_s
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		draw.DrawID, draw.AllowanceID, draw.DayKey, draw.EpisodeID,
		draw.CredentialID, draw.Repo, draw.DrawnAt, draw.TTLS,
	); err != nil {
		return 0, fmt.Errorf("store: mint from allowance: insert draw: %w", err)
	}
	if _, err := conn.ExecContext(ctx, "COMMIT"); err != nil {
		return 0, fmt.Errorf("store: mint from allowance: commit: %w", err)
	}
	rollback = false
	return drawNum, nil
}
