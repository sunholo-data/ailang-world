// Package surfacefollow collects read-only pinned git projections. It publishes
// only to a caller-owned sink; it has no daemon, cache, or route dependency.
package surfacefollow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/surface"
)

const ConfigSchema = "world/surface-follow-config/v1"

// Config is loaded once at startup. Repo paths are canonical after LoadConfig.
type Config struct {
	Schema string `json:"schema"`
	Repos  []Repo `json:"repos"`
}

// Repo enrolls a local branch, never an arbitrary revision expression.
type Repo struct {
	ID      string `json:"id"`
	Repo    string `json:"repo"`
	Charter string `json:"charter"`
	Branch  string `json:"branch"`
	Limit   int    `json:"limit,omitempty"`
}

// LoadOptions supplies explicit path authority; no home lookup is implicit.
type LoadOptions struct {
	Home           string
	ForbiddenRoots []string
}

func configError(id, field, reason string) error {
	// Never repeat untrusted field contents or decoder errors (which can
	// contain tokens). Invalid IDs use the fixed label "config".
	if _, err := surface.QualifyID(id, "episode", "check"); err != nil {
		id = "config"
	}
	return fmt.Errorf("surface follow: repo %s field %s: %s; correct the named field; next: ailang-worldd serve --help", id, field, reason)
}

func expandHome(p, home string) string {
	if strings.HasPrefix(p, "~/") {
		return filepath.Join(home, p[2:])
	}
	return p
}

// owned uses the operating system's FileInfo owner metadata without introducing
// another process or platform-specific dependency. Unsupported metadata
// refuses closed.
func owned(info os.FileInfo) bool {
	v := reflect.ValueOf(info.Sys())
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return false
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return false
	}
	uid := v.FieldByName("Uid")
	return uid.IsValid() && uid.CanUint() && uid.Uint() == uint64(os.Getuid())
}

func validBranch(s string) bool {
	if s == "" || strings.HasPrefix(s, "-") || strings.HasPrefix(s, "/") ||
		strings.Contains(s, "..") || strings.Contains(s, "@{") || s == "@" ||
		strings.ContainsAny(s, "~^:?*[\\") || strings.HasSuffix(s, "/") {
		return false
	}
	for _, r := range s {
		if unicode.IsSpace(r) || unicode.IsControl(r) {
			return false
		}
	}
	for _, part := range strings.Split(s, "/") {
		if part == "" || strings.HasPrefix(part, ".") || strings.HasSuffix(part, ".lock") || strings.HasSuffix(part, ".") {
			return false
		}
	}
	return true
}

// LoadConfig validates permissions before reading bounded JSON. Refusals are
// deliberately independent of raw JSON, path values, and child stderr.
func LoadConfig(path string, opts LoadOptions) (Config, error) {
	var c Config
	fail := func(field, reason string) (Config, error) {
		return Config{}, configError("config", field, reason)
	}
	path = expandHome(path, opts.Home)
	if !filepath.IsAbs(path) {
		return fail("path", "must be absolute")
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return fail("path", "must be an existing regular non-symlink file")
	}
	if info.Mode().Perm() != 0o600 || !owned(info) {
		return fail("mode", "must be owned by the current user with mode 0600")
	}
	dir := filepath.Dir(path)
	di, err := os.Stat(dir)
	if err != nil || !di.IsDir() || di.Mode().Perm()&0o022 != 0 || !owned(di) {
		return fail("directory", "must be operator-owned and not group- or other-writable")
	}
	for _, root := range opts.ForbiddenRoots {
		for _, p := range []string{path, dir} {
			if broker.CheckPolicyOutsideRoot(p, root) != nil {
				return fail("path", "file and directory must be outside workspace roots")
			}
		}
	}
	f, err := os.Open(path)
	if err != nil {
		return fail("path", "cannot read config")
	}
	defer f.Close()
	actual, err := f.Stat()
	if err != nil || !os.SameFile(info, actual) {
		return fail("path", "config changed while opening")
	}
	data, err := io.ReadAll(io.LimitReader(f, surface.MaxBytes+1))
	if err != nil || len(data) > surface.MaxBytes {
		return fail("bytes", "config exceeds read bound")
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&c); err != nil {
		return fail("json", "invalid JSON or unknown field")
	}
	if err := d.Decode(new(any)); err != io.EOF {
		return fail("json", "trailing data")
	}
	if c.Schema != ConfigSchema {
		return fail("schema", "must be world/surface-follow-config/v1")
	}
	if len(c.Repos) < 1 || len(c.Repos) > 4 {
		return fail("repos", "must contain 1–4 repos")
	}
	var presence struct {
		Repos []map[string]json.RawMessage `json:"repos"`
	}
	_ = json.Unmarshal(data, &presence)
	ids, paths, total := map[string]bool{}, map[string]bool{}, 0
	for i := range c.Repos {
		r := &c.Repos[i]
		refuse := func(field, reason string) (Config, error) {
			return Config{}, configError(r.ID, field, reason)
		}
		if _, err := surface.QualifyID(r.ID, "episode", "check"); err != nil || ids[r.ID] {
			return refuse("id", "must be a unique slug")
		}
		ids[r.ID] = true
		if strings.HasPrefix(r.Repo, "~/") && !filepath.IsAbs(opts.Home) {
			return refuse("repo", "home must be absolute for tilde expansion")
		}
		r.Repo = expandHome(r.Repo, opts.Home)
		if !filepath.IsAbs(r.Repo) {
			return refuse("repo", "must be absolute")
		}
		canonical, err := filepath.EvalSymlinks(r.Repo)
		if err != nil {
			return refuse("repo", "cannot canonicalize repository")
		}
		ri, err := os.Stat(canonical)
		if err != nil || !ri.IsDir() || paths[canonical] {
			return refuse("repo", "must name a unique canonical directory")
		}
		r.Repo = canonical
		paths[canonical] = true
		if !filepath.IsLocal(r.Charter) || filepath.Clean(r.Charter) != r.Charter ||
			strings.Contains(r.Charter, "\\") {
			return refuse("charter", "must be a clean relative path without parent segments")
		}
		for _, part := range strings.Split(r.Charter, "/") {
			if part == ".." {
				return refuse("charter", "must not contain parent segments")
			}
		}
		if !validBranch(r.Branch) {
			return refuse("branch", "must name a local branch")
		}
		if r.Limit == 0 {
			if _, explicit := presence.Repos[i]["limit"]; explicit {
				return refuse("limit", "must be between 1 and 2500")
			}
			r.Limit = 200
		}
		if r.Limit < 1 || r.Limit > 2500 {
			return refuse("limit", "must be between 1 and 2500")
		}
		total += r.Limit
		if total > 2000 {
			return refuse("limit", "summed limits exceed 2000; lower this repo limit")
		}
	}
	return c, nil
}
