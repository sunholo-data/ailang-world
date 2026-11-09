package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

// Row 220 M3: bounded, read-only keyset pages for the surface follower's
// episode overlay. Unlike ListSessions' single 500-row page, a caller can
// read every row and can tell a complete read from a bounded one.

// MaxOverlayPage caps one page of either read.
const MaxOverlayPage = 500

// SessionCursor is the keyset position after a page: (created_at,
// credential_id). The zero value starts at the beginning.
type SessionCursor struct {
	CreatedAt    int64
	CredentialID string
}

// ListSessionsAfter returns up to limit rows strictly after cursor, in
// ListSessions' order, and the cursor after the last row returned.
func (c *Credentials) ListSessionsAfter(ctx context.Context, cursor SessionCursor, limit int) ([]SessionRow, SessionCursor, error) {
	if err := requireDeadline(ctx); err != nil {
		return nil, cursor, err
	}
	if limit <= 0 || limit > MaxOverlayPage {
		limit = MaxOverlayPage
	}
	rows, err := c.db.QueryContext(ctx, sessionColumns(c.version)+` FROM session_credentials
WHERE created_at > ? OR (created_at = ? AND credential_id > ?)
ORDER BY created_at, credential_id LIMIT ?;`, cursor.CreatedAt, cursor.CreatedAt, cursor.CredentialID, limit)
	if err != nil {
		return nil, cursor, fmt.Errorf("store: list sessions after: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []SessionRow
	for rows.Next() {
		var r SessionRow
		var project, digest sql.NullString
		if err := rows.Scan(&r.CredentialID, &r.EpisodeID, &r.GrantsJSON, &r.ExpiresAt, &r.CreatedAt, &project, &digest); err != nil {
			return nil, cursor, fmt.Errorf("store: list sessions after: scan: %w", err)
		}
		r.Project, r.ProfileDigest = project.String, digest.String
		out = append(out, r)
		cursor = SessionCursor{CreatedAt: r.CreatedAt, CredentialID: r.CredentialID}
	}
	if err := rows.Err(); err != nil {
		return nil, cursor, fmt.Errorf("store: list sessions after: %w", err)
	}
	return out, cursor, nil
}

// EpisodeStartSummary is the narrow projection of one start record that the
// overlay needs. Preset bytes and paths, the remote, checks and launch
// history never leave SQLite: each field is selected with json_extract.
// Repo is the recorded source path; callers map it to a configured repo and
// must never serialize it.
type EpisodeStartSummary struct {
	Episode   string
	StartedAt string
	Repo      string
	BaseSHA   string
	Branch    string
	Agent     string
	Model     string
	TaskTitle string
	TTLS      int64
	Resumes   []EpisodeResume
}

// ListEpisodeStartSummaries returns up to limit starts with episode > after,
// ordered by episode. A store older than version 3 has no starts.
func (c *Credentials) ListEpisodeStartSummaries(ctx context.Context, after string, limit int) ([]EpisodeStartSummary, error) {
	if err := requireDeadline(ctx); err != nil {
		return nil, err
	}
	if c.version < 3 {
		return nil, nil
	}
	if limit <= 0 || limit > MaxOverlayPage {
		limit = MaxOverlayPage
	}
	rows, err := c.db.QueryContext(ctx, `SELECT episode,
  json_extract(start_json, '$.started_at'), json_extract(start_json, '$.repo'),
  json_extract(start_json, '$.base.sha'), json_extract(start_json, '$.branch'),
  json_extract(start_json, '$.agent'), json_extract(start_json, '$.model'),
  json_extract(start_json, '$.task.title'), json_extract(start_json, '$.ttl_s'),
  json_extract(start_json, '$.resumes')
FROM episode_starts WHERE episode > ? ORDER BY episode LIMIT ?;`, after, limit)
	if err != nil {
		return nil, fmt.Errorf("store: list episode starts: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []EpisodeStartSummary
	for rows.Next() {
		var s EpisodeStartSummary
		var startedAt, repo, base, branch, agent, model, title, resumes sql.NullString
		var ttl sql.NullInt64
		if err := rows.Scan(&s.Episode, &startedAt, &repo, &base, &branch, &agent, &model, &title, &ttl, &resumes); err != nil {
			return nil, fmt.Errorf("store: list episode starts: scan: %w", err)
		}
		s.StartedAt, s.Repo, s.BaseSHA, s.Branch = startedAt.String, repo.String, base.String, branch.String
		s.Agent, s.Model, s.TaskTitle, s.TTLS = agent.String, model.String, title.String, ttl.Int64
		if resumes.Valid && resumes.String != "" {
			if err := json.Unmarshal([]byte(resumes.String), &s.Resumes); err != nil {
				return nil, fmt.Errorf("store: list episode starts: %s resumes: %w", s.Episode, err)
			}
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: list episode starts: %w", err)
	}
	return out, nil
}
