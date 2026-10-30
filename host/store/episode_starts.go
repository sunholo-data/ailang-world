package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
)

const episodeStartsSQL = `CREATE TABLE episode_starts (
    episode TEXT PRIMARY KEY CHECK (episode <> ''),
    start_json TEXT NOT NULL
)`

// EpisodeStart holds immutable attended mint fields, independent of credentials.
// Mint fields never change. Launches and Resumes are appended only by
// AppendEpisodeLaunch and MintSessionWithResume respectively. FailEpisodeResume
// may mark a resume failed without removing its history.
// It contains no token and has no session-authorized writer.
type EpisodeStart struct {
	Schema        string          `json:"schema"`
	Episode       string          `json:"episode"`
	StartedAt     string          `json:"started_at"`
	Preset        StartPreset     `json:"preset"`
	Repo          string          `json:"repo"`
	Remote        StartRemote     `json:"remote"`
	Base          StartBase       `json:"base"`
	Branch        string          `json:"branch"`
	BranchExisted bool            `json:"branch_existed"`
	Agent         string          `json:"agent"`
	Model         string          `json:"model"`
	Overrides     []string        `json:"overrides"`
	Task          *StartTask      `json:"task,omitempty"`
	Checks        []StartCheck    `json:"checks"`
	TTLS          int64           `json:"ttl_s"`
	Resumes       []EpisodeResume `json:"resumes,omitempty"`
	Launches      []EpisodeLaunch `json:"launches,omitempty"`
}
type StartPreset struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
	Bytes  string `json:"bytes"`
}
type StartRemote struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type StartBase struct {
	Branch string `json:"branch"`
	SHA    string `json:"sha"`
}
type StartTask struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
	Title  string `json:"title"`
}
type StartCheck struct {
	Name        string   `json:"name"`
	Argv        []string `json:"argv"`
	Cwd         string   `json:"cwd"`
	TimeoutS    float64  `json:"timeout_s"`
	Unsandboxed bool     `json:"unsandboxed,omitempty"`
	Sandbox     string   `json:"sandbox"`
	Exec        string   `json:"exec,omitempty"`
	Args        []string `json:"args"`
}

// EpisodeStart reads mint fields and launch history, including from a read-only handle.
func (c *Credentials) EpisodeStart(ctx context.Context, episode string) (*EpisodeStart, error) {
	if err := requireDeadline(ctx); err != nil {
		return nil, err
	}
	if c.version < 3 {
		return nil, nil
	}
	var data string
	err := c.db.QueryRowContext(ctx, "SELECT start_json FROM episode_starts WHERE episode = ?", episode).Scan(&data)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var start EpisodeStart
	if err := json.Unmarshal([]byte(data), &start); err != nil {
		return nil, err
	}
	return &start, nil
}

type EpisodeLaunch struct {
	Mode             string   `json:"mode"`
	StartedAt        string   `json:"started_at"`
	EndedAt          string   `json:"ended_at"`
	Exit             int      `json:"exit"`
	FinishReason     string   `json:"finish_reason"`
	Transcript       string   `json:"transcript"`
	TranscriptSHA256 string   `json:"transcript_sha256"`
	TranscriptState  string   `json:"transcript_state"`
	Host             string   `json:"host"`
	ArgvRedacted     []string `json:"argv_redacted"`
	HostSessionID    string   `json:"host_session_id,omitempty"`
	HostTranscript   string   `json:"host_transcript,omitempty"`
	TranscriptReason string   `json:"transcript_reason,omitempty"`
	MCPConfigDeleted *bool    `json:"mcp_config_deleted,omitempty"`
}

// AppendEpisodeLaunch is the sole launch writer. OpenCredentials owns the
// credentials writer lock; no session route exposes this operation. Raw JSON
// preserves every mint field's value, including fields unknown to this build.
func (c *Credentials) AppendEpisodeLaunch(ctx context.Context, episode string, launch EpisodeLaunch) error {
	if err := requireDeadline(ctx); err != nil {
		return err
	}
	if c.lock == nil {
		return errors.New("store: append launch requires the credentials writer lock")
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	var data string
	if err := tx.QueryRowContext(ctx, "SELECT start_json FROM episode_starts WHERE episode = ?", episode).Scan(&data); err != nil {
		return err
	}
	var row map[string]json.RawMessage
	if err := json.Unmarshal([]byte(data), &row); err != nil {
		return err
	}
	var launches []json.RawMessage
	if raw := row["launches"]; raw != nil {
		if err := json.Unmarshal(raw, &launches); err != nil {
			return err
		}
	}
	entry, err := json.Marshal(launch)
	if err != nil {
		return err
	}
	launches = append(launches, entry)
	row["launches"], err = json.Marshal(launches)
	if err != nil {
		return err
	}
	encoded, err := json.Marshal(row)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "UPDATE episode_starts SET start_json = ? WHERE episode = ?", string(encoded), episode); err != nil {
		return err
	}
	return tx.Commit()
}

type EpisodeResume struct {
	CredentialID     string `json:"credential_id"`
	Failed           bool   `json:"failed,omitempty"`
	ResumedFrom      string `json:"resumed_from"`
	ResumedAt        string `json:"resumed_at"`
	TTLS             int64  `json:"ttl_s"`
	PresetFileSHA256 string `json:"preset_file_sha256"`
}

// CurrentEpisodeCredential returns the last successful resume's credential,
// or the original start credential. Failed attempts never advance the chain.
func CurrentEpisodeCredential(start *EpisodeStart, rows []SessionRow) string {
	for i := len(start.Resumes) - 1; i >= 0; i-- {
		if !start.Resumes[i].Failed {
			return start.Resumes[i].CredentialID
		}
	}
	if len(start.Resumes) > 0 {
		return start.Resumes[0].ResumedFrom
	}
	var oldest *SessionRow
	for i := range rows {
		if oldest == nil || rows[i].CreatedAt < oldest.CreatedAt {
			oldest = &rows[i]
		}
	}
	if oldest != nil {
		return oldest.CredentialID
	}
	return ""
}

// checkResume runs before insertion in the mint transaction. CreatedAt is the
// caller's mint time, never a store clock read.
func checkResume(ctx context.Context, tx *sql.Tx, row SessionRow, resume EpisodeResume) error {
	var data string
	if err := tx.QueryRowContext(ctx, "SELECT start_json FROM episode_starts WHERE episode = ?", row.EpisodeID).Scan(&data); err != nil {
		return err
	}
	var start EpisodeStart
	if err := json.Unmarshal([]byte(data), &start); err != nil {
		return err
	}
	rs, err := tx.QueryContext(ctx, "SELECT credential_id, created_at, expires_at FROM session_credentials WHERE episode_id = ?", row.EpisodeID)
	if err != nil {
		return err
	}
	var rows []SessionRow
	for rs.Next() {
		var r SessionRow
		if err := rs.Scan(&r.CredentialID, &r.CreatedAt, &r.ExpiresAt); err != nil {
			_ = rs.Close()
			return err
		}
		rows = append(rows, r)
	}
	err = rs.Err()
	_ = rs.Close()
	if err != nil {
		return err
	}
	current := CurrentEpisodeCredential(&start, rows)
	found := false
	for _, r := range rows {
		if r.CredentialID == current {
			found = true
		}
		if r.ExpiresAt >= row.CreatedAt {
			return errors.New("episode has a live session (F12)")
		}
	}
	if !found {
		return errors.New("no credential to resume (revoked or finished): revoke ends an episode")
	}
	if current != resume.ResumedFrom {
		return errors.New("episode current credential changed (F12); retry resume")
	}
	return nil
}

// MintSessionWithResume atomically mints a credential and appends resume history.
func (c *Credentials) MintSessionWithResume(ctx context.Context, row SessionRow, resume EpisodeResume) error {
	if c.lock == nil {
		return errors.New("store: resume requires the credentials writer lock")
	}
	return c.mintSession(ctx, row, nil, &resume)
}

// FailEpisodeResume keeps history while returning current authority to the
// previous credential. It changes only the failed entry's status.
func (c *Credentials) FailEpisodeResume(ctx context.Context, episode, credential string) error {
	if err := requireDeadline(ctx); err != nil {
		return err
	}
	if c.lock == nil {
		return errors.New("store: fail resume requires the credentials writer lock")
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	var data string
	if err := tx.QueryRowContext(ctx, "SELECT start_json FROM episode_starts WHERE episode = ?", episode).Scan(&data); err != nil {
		return err
	}
	var start map[string]json.RawMessage
	if err := json.Unmarshal([]byte(data), &start); err != nil {
		return err
	}
	var entries []map[string]json.RawMessage
	if err := json.Unmarshal(start["resumes"], &entries); err != nil {
		return err
	}
	for _, entry := range entries {
		var id string
		if err := json.Unmarshal(entry["credential_id"], &id); err != nil {
			return err
		}
		if id == credential {
			entry["failed"] = json.RawMessage("true")
		}
	}
	start["resumes"], err = json.Marshal(entries)
	if err != nil {
		return err
	}
	encoded, err := json.Marshal(start)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "UPDATE episode_starts SET start_json = ? WHERE episode = ?", string(encoded), episode); err != nil {
		return err
	}
	return tx.Commit()
}

func appendResume(ctx context.Context, tx *sql.Tx, episode string, resume EpisodeResume) error {
	var data string
	if err := tx.QueryRowContext(ctx, "SELECT start_json FROM episode_starts WHERE episode = ?", episode).Scan(&data); err != nil {
		return err
	}
	var row map[string]json.RawMessage
	if err := json.Unmarshal([]byte(data), &row); err != nil {
		return err
	}
	var entries []json.RawMessage
	if raw := row["resumes"]; raw != nil {
		if err := json.Unmarshal(raw, &entries); err != nil {
			return err
		}
	}
	entry, err := json.Marshal(resume)
	if err != nil {
		return err
	}
	row["resumes"], err = json.Marshal(append(entries, entry))
	if err != nil {
		return err
	}
	encoded, err := json.Marshal(row)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "UPDATE episode_starts SET start_json = ? WHERE episode = ?", string(encoded), episode)
	return err
}
