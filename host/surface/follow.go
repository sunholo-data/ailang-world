package surface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

// FollowSchema wraps unchanged v1 partitions. The episode overlay is absent
// until M3 defines its validation contract.
const FollowSchema = "world/surface-follow-index/v1"

var repoSlug = regexp.MustCompile(`^[a-z][a-z0-9-]{0,31}$`)

// FollowIndex preserves configuration order; generation identifies the whole set.
type FollowIndex struct {
	Schema     string            `json:"schema"`
	Generation string            `json:"generation"`
	Repos      []FollowPartition `json:"repos"`
}

// FollowFailure contains a short, path-free cause, never raw child stderr.
type FollowFailure struct {
	Class string `json:"class"`
	Text  string `json:"text"`
}

// FollowPartition is one accepted, pinned v1 projection and its freshness.
type FollowPartition struct {
	ID          string         `json:"id"`
	Head        string         `json:"head"`
	BuiltAt     string         `json:"built_at"`
	CheckedAt   string         `json:"checked_at"`
	LastFailure *FollowFailure `json:"last_failure,omitempty"`
	Index       Index          `json:"index"`
}

// QualifyID uses the same three-part form for all kinds, including branes.
func QualifyID(repo, kind, local string) (string, error) {
	if !repoSlug.MatchString(repo) {
		return "", fmt.Errorf("qualified id: repo must be a slug")
	}
	if kind != "commit" && kind != "episode" && kind != "brane" {
		return "", fmt.Errorf("qualified id: unknown kind")
	}
	if local == "" || (kind != "brane" && strings.Contains(local, ":")) {
		return "", fmt.Errorf("qualified id: invalid local part")
	}
	if kind == "commit" && !commitHash.MatchString(local) {
		return "", fmt.Errorf("qualified id: commit must be a full commit hash")
	}
	return repo + ":" + kind + ":" + local, nil
}

// ParseQualifiedID reverses QualifyID, validating every component.
func ParseQualifiedID(id string) (repo, kind, local string, err error) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("qualified id: want repo:kind:local")
	}
	if _, err = QualifyID(parts[0], parts[1], parts[2]); err != nil {
		return "", "", "", err
	}
	return parts[0], parts[1], parts[2], nil
}

// ValidateFollow does not weaken any v1 invariant.
func ValidateFollow(ix FollowIndex) error {
	if ix.Schema != FollowSchema {
		return fmt.Errorf("surface follow: schema must be %s (legacy schema is %s)", FollowSchema, Schema)
	}
	if ix.Generation == "" {
		return fmt.Errorf("surface follow: generation is empty")
	}
	if len(ix.Repos) < 1 || len(ix.Repos) > 4 {
		return fmt.Errorf("surface follow: repos must contain 1–4 partitions")
	}
	seen, total := map[string]bool{}, 0
	for i, p := range ix.Repos {
		if !repoSlug.MatchString(p.ID) || seen[p.ID] {
			return fmt.Errorf("surface follow: repos[%d].id must be a unique slug", i)
		}
		seen[p.ID] = true
		if !commitHash.MatchString(p.Head) || p.Head != p.Index.Source.Head {
			return fmt.Errorf("surface follow: repo %s head disagrees with index source.head or is invalid", p.ID)
		}
		for _, field := range []struct{ name, value string }{{"built_at", p.BuiltAt}, {"checked_at", p.CheckedAt}} {
			if _, err := time.Parse(time.RFC3339, field.value); err != nil {
				return fmt.Errorf("surface follow: repo %s %s must be RFC 3339", p.ID, field.name)
			}
		}
		if f := p.LastFailure; f != nil {
			if !repoSlug.MatchString(f.Class) || f.Text == "" || len(f.Text) > 240 ||
				strings.ContainsAny(f.Text, "/\\\r\n\x00") {
				return fmt.Errorf("surface follow: repo %s last_failure must be short and path-free", p.ID)
			}
		}
		if err := Validate(p.Index); err != nil {
			return fmt.Errorf("surface follow: repo %s index: %w", p.ID, err)
		}
		total += len(p.Index.States)
	}
	if total > MaxStates {
		return fmt.Errorf("surface follow: summed states exceed %d; lower a repo limit", MaxStates)
	}
	return nil
}

// DecodeFollow bounds and strictly decodes an envelope, including trailing data.
func DecodeFollow(data []byte) (FollowIndex, error) {
	var ix FollowIndex
	if len(data) > MaxBytes {
		return ix, fmt.Errorf("surface follow: bytes exceed %d", MaxBytes)
	}
	var header struct {
		Schema string `json:"schema"`
	}
	if err := json.Unmarshal(data, &header); err != nil {
		return ix, fmt.Errorf("surface follow: invalid JSON")
	}
	if header.Schema != FollowSchema {
		return ix, fmt.Errorf("surface follow: schema must be %s (legacy schema is %s)", FollowSchema, Schema)
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&ix); err != nil {
		return FollowIndex{}, fmt.Errorf("surface follow: decode: %w", err)
	}
	if err := d.Decode(new(any)); err != io.EOF {
		return FollowIndex{}, fmt.Errorf("surface follow: trailing data")
	}
	if err := ValidateFollow(ix); err != nil {
		return FollowIndex{}, err
	}
	return ix, nil
}
