// Package approval is the pure core of row 201's approvals by key
// (design_docs/planned/w-world-approvals-by-key-a0a1.md): the mint request,
// its canonical bytes and digest, the signed message, the admission law, the
// mint-row check and the SSHSIG verifier. Standard library only; no I/O, no
// clock, no randomness (coding standard S2): every effect is the caller's.
package approval

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

const RequestSchema = "world/mint-request/v1"

// MessagePrefix starts every signed approval message (design §2.2).
const MessagePrefix = "world-approval/v1"

// RequestLifetime is a request's fixed lifetime: expires_at = created_at + 600.
const RequestLifetime = 10 * time.Minute

// MintRequest is the stored, single-use request an approval executes. Its
// canonical JSON (Canonical) is what is digested and what is stored.
type MintRequest struct {
	Schema        string     `json:"schema"`         // RequestSchema
	Episode       string     `json:"episode"`        // ^[a-z0-9][a-z0-9-]{0,63}$
	WorldDB       string     `json:"world_db"`       // [r2] canonical --db path
	WorkspaceRoot string     `json:"workspace_root"` // [r2] absolute --workspace-root
	Worktree      string     `json:"worktree"`       // absolute, clean: <workspace-root>/<episode>
	Reuse         bool       `json:"reuse"`          // the worktree exists already
	Degraded      bool       `json:"degraded"`       // non-git (always false on this path, §2.1a)
	Repo          string     `json:"repo"`           // [r2] absolute (filepath.Abs at create), or ""
	Git           string     `json:"git"`            // [r2] the resolved git binary; "" when Reuse
	Branch        string     `json:"branch"`
	Grants        []Grant    `json:"grants"`         // in order, at least one
	TTLSeconds    int64      `json:"ttl_s"`          // ≥ 0
	Profile       ProfileRef `json:"profile"`        // both "" when unbound
	ProfilesDir   string     `json:"profiles_dir"`   // [r2] absolute --exec-profiles-dir
	CredentialsDB string     `json:"credentials_db"` // canonical path the mint writes
	SessionFile   string     `json:"session_file"`   // absolute; the token's only file
	MCPConfig     string     `json:"mcp_config"`     // absolute for claude, "" for codex
	Agent         string     `json:"agent"`          // claude | codex
	Model         string     `json:"model"`
	Task          string     `json:"task"` // absolute or ""
	Daemon        DaemonRef  `json:"daemon"`
}

type Grant struct {
	Effect string `json:"effect"`
	Scope  string `json:"scope"`
	Budget int64  `json:"budget"`
}

type ProfileRef struct {
	Project string `json:"project"`
	Digest  string `json:"digest"`
}

type DaemonRef struct { // the requesting daemon's build identity (parent §2.2) and address
	Addr             string `json:"addr"`
	VCSRevision      string `json:"vcs_revision"`
	VCSModified      bool   `json:"vcs_modified"`
	ExecutableSHA256 string `json:"executable_sha256"`
	GoVersion        string `json:"go_version"`
}

var (
	ErrSchema         = errors.New("approval: bad request schema")
	ErrString         = errors.New("approval: string is not valid UTF-8 or holds NUL")
	ErrEpisode        = errors.New("approval: bad episode")
	ErrPath           = errors.New("approval: path is not absolute and clean")
	ErrGitReuse       = errors.New("approval: git must be set exactly when reuse is false")
	ErrNoGrants       = errors.New("approval: no grants")
	ErrGrantEmpty     = errors.New("approval: grant has an empty effect or scope")
	ErrNegative       = errors.New("approval: negative value")
	ErrHalfProfile    = errors.New("approval: profile is half set")
	ErrAgent          = errors.New("approval: agent is not claude or codex")
	ErrMCPConfig      = errors.New("approval: mcp_config does not match the agent")
	ErrDecode         = errors.New("approval: request does not decode")
	ErrNotCanonical   = errors.New("approval: request bytes are not canonical")
	ErrNegativeExpiry = errors.New("approval: negative expiry")
)

var episodePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

// Validate refuses a request that the mint must never act on. Each refusal
// names its JSON path in the error text.
func (r MintRequest) Validate() error {
	if r.Schema != RequestSchema {
		return fmt.Errorf("%w: schema is %q, want %q", ErrSchema, r.Schema, RequestSchema)
	}
	if err := stringLeaves(reflect.ValueOf(r), "", func(path, s string) error {
		if !utf8.ValidString(s) || strings.ContainsRune(s, 0) {
			return fmt.Errorf("%w: %s", ErrString, path)
		}
		return nil
	}); err != nil {
		return err
	}
	if !episodePattern.MatchString(r.Episode) {
		return fmt.Errorf("%w: episode %q", ErrEpisode, r.Episode)
	}
	for _, p := range []struct{ name, val string }{
		{"world_db", r.WorldDB},
		{"workspace_root", r.WorkspaceRoot},
		{"worktree", r.Worktree},
		{"profiles_dir", r.ProfilesDir},
		{"credentials_db", r.CredentialsDB},
		{"session_file", r.SessionFile},
	} {
		if !absClean(p.val) {
			return fmt.Errorf("%w: %s %q", ErrPath, p.name, p.val)
		}
	}
	for _, p := range []struct{ name, val string }{
		{"repo", r.Repo},
		{"git", r.Git},
		{"mcp_config", r.MCPConfig},
		{"task", r.Task},
	} {
		if p.val != "" && !absClean(p.val) {
			return fmt.Errorf("%w: %s %q", ErrPath, p.name, p.val)
		}
	}
	if (r.Git != "") == r.Reuse {
		return fmt.Errorf("%w: git is %q with reuse %t", ErrGitReuse, r.Git, r.Reuse)
	}
	if len(r.Grants) == 0 {
		return fmt.Errorf("%w: grants", ErrNoGrants)
	}
	for i, g := range r.Grants {
		if g.Effect == "" || g.Scope == "" {
			return fmt.Errorf("%w: grants[%d]", ErrGrantEmpty, i)
		}
		if g.Budget < 0 {
			return fmt.Errorf("%w: grants[%d].budget is %d", ErrNegative, i, g.Budget)
		}
	}
	if r.TTLSeconds < 0 {
		return fmt.Errorf("%w: ttl_s is %d", ErrNegative, r.TTLSeconds)
	}
	if (r.Profile.Project == "") != (r.Profile.Digest == "") {
		return fmt.Errorf("%w: profile needs both project and digest", ErrHalfProfile)
	}
	if r.Agent != "claude" && r.Agent != "codex" {
		return fmt.Errorf("%w: agent %q", ErrAgent, r.Agent)
	}
	if (r.Agent == "claude" && r.MCPConfig == "") || (r.Agent == "codex" && r.MCPConfig != "") {
		return fmt.Errorf("%w: mcp_config is %q for agent %q", ErrMCPConfig, r.MCPConfig, r.Agent)
	}
	return nil
}

// Canonical is the compact JSON of r in struct field order, with HTML
// escaping off and no trailing newline. It does not call Validate.
func Canonical(r MintRequest) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(r); err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil
}

// Decode is the strict inverse of Canonical. It refuses unknown fields,
// trailing data, and any input that does not re-encode to the same bytes. It
// does not call Validate; callers do.
func Decode(b []byte) (MintRequest, error) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	var r MintRequest
	if err := dec.Decode(&r); err != nil {
		return MintRequest{}, fmt.Errorf("%w: %v", ErrDecode, err)
	}
	c, err := Canonical(r)
	if err != nil {
		return MintRequest{}, err
	}
	if !bytes.Equal(c, b) {
		return MintRequest{}, ErrNotCanonical
	}
	return r, nil
}

// Digest is the SHA-256 of Canonical(r).
func Digest(r MintRequest) ([32]byte, error) {
	c, err := Canonical(r)
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(c), nil
}

// Diff names the first JSON path at which a and b differ, or "" when equal.
func Diff(a, b MintRequest) string {
	return diffValue(reflect.ValueOf(a), reflect.ValueOf(b), "")
}

// SignedMessage is the 90-byte message every approval signs (design §2.2):
// MessagePrefix ‖ 0x00 ‖ digest ‖ challenge ‖ expiry (uint64, big-endian).
func SignedMessage(digest, challenge [32]byte, expiry int64) ([]byte, error) {
	if expiry < 0 {
		return nil, ErrNegativeExpiry
	}
	m := make([]byte, 0, 90)
	m = append(m, MessagePrefix...)
	m = append(m, 0)
	m = append(m, digest[:]...)
	m = append(m, challenge[:]...)
	m = binary.BigEndian.AppendUint64(m, uint64(expiry))
	return m, nil
}

// MessageHash is the SHA-256 of SignedMessage(digest, challenge, expiry).
func MessageHash(digest, challenge [32]byte, expiry int64) ([32]byte, error) {
	m, err := SignedMessage(digest, challenge, expiry)
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(m), nil
}

func jsonName(f reflect.StructField) string {
	return strings.Split(f.Tag.Get("json"), ",")[0]
}

// stringLeaves calls fn for every string leaf, with its JSON path.
func stringLeaves(v reflect.Value, path string, fn func(path, s string) error) error {
	switch v.Kind() {
	case reflect.String:
		return fn(path, v.String())
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			p := jsonName(t.Field(i))
			if path != "" {
				p = path + "." + p
			}
			if err := stringLeaves(v.Field(i), p, fn); err != nil {
				return err
			}
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if err := stringLeaves(v.Index(i), fmt.Sprintf("%s[%d]", path, i), fn); err != nil {
				return err
			}
		}
	}
	return nil
}

// diffValue returns the JSON path of the first differing leaf, or "".
func diffValue(a, b reflect.Value, path string) string {
	switch a.Kind() {
	case reflect.Struct:
		t := a.Type()
		for i := 0; i < a.NumField(); i++ {
			p := jsonName(t.Field(i))
			if path != "" {
				p = path + "." + p
			}
			if d := diffValue(a.Field(i), b.Field(i), p); d != "" {
				return d
			}
		}
		return ""
	case reflect.Slice:
		if a.Len() != b.Len() {
			return path
		}
		for i := 0; i < a.Len(); i++ {
			if d := diffValue(a.Index(i), b.Index(i), fmt.Sprintf("%s[%d]", path, i)); d != "" {
				return d
			}
		}
		return ""
	default:
		if a.Interface() != b.Interface() {
			return path
		}
		return ""
	}
}

// absClean reports whether p is absolute and already in clean form.
func absClean(p string) bool {
	return filepath.IsAbs(p) && filepath.Clean(p) == p
}
