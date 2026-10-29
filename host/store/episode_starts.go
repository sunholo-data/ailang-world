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

// EpisodeStart is immutable attended mint history, independent of credentials.
// It contains no token and has no session-authorized writer.
type EpisodeStart struct {
	Schema        string       `json:"schema"`
	Episode       string       `json:"episode"`
	StartedAt     string       `json:"started_at"`
	Preset        StartPreset  `json:"preset"`
	Repo          string       `json:"repo"`
	Remote        StartRemote  `json:"remote"`
	Base          StartBase    `json:"base"`
	Branch        string       `json:"branch"`
	BranchExisted bool         `json:"branch_existed"`
	Agent         string       `json:"agent"`
	Model         string       `json:"model"`
	Overrides     []string     `json:"overrides"`
	Task          *StartTask   `json:"task,omitempty"`
	Checks        []StartCheck `json:"checks"`
	TTLS          int64        `json:"ttl_s"`
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

// EpisodeStart reads one immutable row, including from a read-only handle.
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
