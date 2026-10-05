package broker

// Row 140 M2 (design_docs/planned/w-workspace-exec-toolchain-effect.md §4.3):
// the operator's exec profile, read once into memory and digested (§4.6
// "Profile"). ParseExecProfile is the in-memory loader: the field refusals a
// handler needs to run a call honestly. The `serve` flags, the placement
// rules and the read_roots ancestor/realpath rules are M3's startup table
// (exec_startup.go); M3 adds the `probe` key (§4.6 arm 7).

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/sunholo-data/ailang-world/host/childenv"
)

// ExecProfileVersion is the only profile schema this host reads.
const ExecProfileVersion = "world/exec-profile/v1"

// ExecMaxTimeoutMS is HandlerCap − 1 000 ms (D-140-4 = A): a command's own
// deadline leaves the handler a second for verification and cleanup. Bound
// to coordinator.HandlerCap by a daemon test.
const ExecMaxTimeoutMS = 9000

// execMaxArgs is the plan's argument bound (§4.1), which max_args cannot exceed.
const execMaxArgs = 16

// Flag forms (row 140 M2, ruled in session 2026-10-05): how the emitted argv
// spells a valued flag. Matching is always on the normalized -f=v.
const (
	ExecFormEq  = "eq"  // one item: -f=v (the default)
	ExecFormSep = "sep" // two items: -f, v (argparse short options)
)

var (
	execProjectPattern   = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)
	execCommandIDPattern = regexp.MustCompile(`^[a-z][a-z0-9-]{0,31}$`)
	execEnvNamePattern   = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
)

// execWorldEnv is the child-environment set World always writes (§4.5), in
// rendered order; the profile may name none of them.
var execWorldEnv = []string{"HOME", "TMPDIR", "PATH", "LANG", "LC_ALL", "PYTHONPYCACHEPREFIX", "npm_config_cache", "CI", "NO_COLOR"}

// ExecProfile is one loaded, validated project profile (§4.3).
type ExecProfile struct {
	Project   string
	Root      string // the project root inside the worktree, "." or a relpath
	Path      []string
	Env       map[string]string
	ReadRoots []string
	Caches    map[string]ExecCache
	TimeoutMS int
	Commands  map[string]ExecProfileCommand
	// Probe is the startup probe's toolchain argv (§4.6 arm 7, row 140 M3):
	// the `probe` key, or by default the argv[0] of the first command (in
	// byte order of id) with --version. Its argv[0] is resolved like a
	// command's; it must exit 0 under the sandbox.
	Probe ExecProfileCommand
	// Digest is "sha256:<hex>" of the exact bytes read: the record's
	// profile.digest names what ran.
	Digest string
}

// ExecCache is one per-episode cache directory, seeded on first use from
// the operator's read-only Seed ("" = start empty).
type ExecCache struct {
	Seed string
}

// ExecProfileCommand is a command with its argv[0] resolved once at load.
type ExecProfileCommand struct {
	ExecCommand
	Argv0       string // absolute path of the executable argv[0] names
	Argv0SHA256 string // its sha256 at load; re-verified before every call
}

// ExecProfileError is a profile load refusal naming the field.
type ExecProfileError struct {
	Field string
	Why   string
}

func (e *ExecProfileError) Error() string { return fmt.Sprintf("exec profile: %s: %s", e.Field, e.Why) }

func profileErr(field, format string, args ...any) error {
	return &ExecProfileError{Field: field, Why: fmt.Sprintf(format, args...)}
}

// strictObject decodes one JSON object and refuses any key not in allowed.
// Exact names only: encoding/json would match "Project" to "project", so the
// struct decoder is never trusted with the key set.
func strictObject(field string, raw json.RawMessage, allowed ...string) (map[string]json.RawMessage, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil || m == nil {
		return nil, profileErr(field, "is not a JSON object")
	}
	for k := range m {
		if !slices.Contains(allowed, k) {
			return nil, profileErr(field, "unknown key %q", k)
		}
	}
	return m, nil
}

func decodeField(field string, raw json.RawMessage, v any) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return profileErr(field, "has the wrong type: %v", err)
	}
	return nil
}

// ParseExecProfile reads a profile's bytes once, validates every field this
// host relies on and resolves each command's argv[0] on the profile path.
func ParseExecProfile(data []byte) (*ExecProfile, error) {
	top, err := strictObject("profile", data, "profile", "project", "root", "path", "env", "read_roots", "caches",
		"timeout_ms", "commands", "probe")
	if err != nil {
		return nil, err
	}
	for _, k := range []string{"profile", "project", "path", "timeout_ms", "commands"} {
		if _, ok := top[k]; !ok {
			return nil, profileErr(k, "is required")
		}
	}
	var version string
	if err := decodeField("profile", top["profile"], &version); err != nil {
		return nil, err
	}
	if version != ExecProfileVersion {
		return nil, profileErr("profile", "is %q; this host reads %q", version, ExecProfileVersion)
	}
	sum := sha256.Sum256(data)
	p := &ExecProfile{Root: ".", Env: map[string]string{}, Caches: map[string]ExecCache{},
		Commands: map[string]ExecProfileCommand{}, Digest: "sha256:" + hex.EncodeToString(sum[:])}
	if err := decodeField("project", top["project"], &p.Project); err != nil {
		return nil, err
	}
	if !execProjectPattern.MatchString(p.Project) {
		return nil, profileErr("project", "%q does not match %s", p.Project, execProjectPattern)
	}
	if raw, ok := top["root"]; ok {
		if err := decodeField("root", raw, &p.Root); err != nil {
			return nil, err
		}
	}
	if !execPathOk(p.Root) || filepath.Clean(p.Root) != p.Root {
		return nil, profileErr("root", "%q is not a clean relative path inside the worktree", p.Root)
	}
	if err := decodeField("path", top["path"], &p.Path); err != nil {
		return nil, err
	}
	if len(p.Path) == 0 {
		return nil, profileErr("path", "is empty")
	}
	for _, d := range p.Path {
		if !cleanAbs(d) || strings.Contains(d, ":") {
			return nil, profileErr("path", "%q is not a clean absolute directory without ':'", d)
		}
	}
	if raw, ok := top["read_roots"]; ok {
		if err := decodeField("read_roots", raw, &p.ReadRoots); err != nil {
			return nil, err
		}
	}
	for _, r := range p.ReadRoots {
		if !cleanAbs(r) {
			return nil, profileErr("read_roots", "%q is not a clean absolute path", r)
		}
	}
	if raw, ok := top["caches"]; ok {
		if err := parseExecCaches(raw, p); err != nil {
			return nil, err
		}
	}
	if raw, ok := top["env"]; ok {
		if err := decodeField("env", raw, &p.Env); err != nil {
			return nil, err
		}
	}
	for name, value := range p.Env {
		if err := execEnvNameOk("env", name); err != nil {
			return nil, err
		}
		if _, isCache := p.Caches[name]; isCache {
			return nil, profileErr("env", "%q is a caches variable; World sets it", name)
		}
		if strings.ContainsRune(value, 0) {
			return nil, profileErr("env", "%q has a NUL byte", name)
		}
	}
	if mod, ok := p.Env["GOMODCACHE"]; ok && !underAny(mod, p.ReadRoots) {
		// §4.4: the module cache is read-only by construction only when it is
		// a read root and never in allowWrite.
		return nil, profileErr("env", "GOMODCACHE %q is not inside a read_roots entry", mod)
	}
	if err := decodeField("timeout_ms", top["timeout_ms"], &p.TimeoutMS); err != nil {
		return nil, err
	}
	if p.TimeoutMS < 1 || p.TimeoutMS > ExecMaxTimeoutMS {
		return nil, profileErr("timeout_ms", "%d is outside 1–%d (HandlerCap − 1 000 ms, D-140-4)", p.TimeoutMS, ExecMaxTimeoutMS)
	}
	var commands map[string]json.RawMessage
	if err := decodeField("commands", top["commands"], &commands); err != nil {
		return nil, err
	}
	if len(commands) == 0 {
		return nil, profileErr("commands", "is empty")
	}
	for id, raw := range commands {
		c, err := parseExecCommand(id, raw, p.Path)
		if err != nil {
			return nil, err
		}
		p.Commands[id] = c
	}
	if err := parseExecProbe(top["probe"], p); err != nil {
		return nil, err
	}
	return p, nil
}

// parseExecProbe reads the `probe` key (§4.6 arm 7) or derives its default:
// the first command's argv[0] with --version.
func parseExecProbe(raw json.RawMessage, p *ExecProfile) error {
	var argv []string
	if raw == nil {
		first := sortedKeys(p.Commands)[0]
		argv = []string{p.Commands[first].Argv[0], "--version"}
	} else {
		if err := decodeField("probe", raw, &argv); err != nil {
			return err
		}
		if len(argv) == 0 {
			return profileErr("probe", "is empty")
		}
		for _, a := range argv {
			if a == "" || strings.ContainsRune(a, 0) {
				return profileErr("probe", "has an empty or NUL-bearing item")
			}
		}
	}
	argv0, sum, err := resolveArgv0(argv[0], p.Path)
	if err != nil {
		return profileErr("probe", "%v", err)
	}
	p.Probe = ExecProfileCommand{ExecCommand: ExecCommand{Argv: argv}, Argv0: argv0, Argv0SHA256: sum}
	return nil
}

func parseExecCaches(raw json.RawMessage, p *ExecProfile) error {
	var caches map[string]json.RawMessage
	if err := decodeField("caches", raw, &caches); err != nil {
		return err
	}
	for name, obj := range caches {
		if err := execEnvNameOk("caches", name); err != nil {
			return err
		}
		fields, err := strictObject("caches."+name, obj, "seed")
		if err != nil {
			return err
		}
		var c ExecCache
		if s, ok := fields["seed"]; ok {
			if err := decodeField("caches."+name+".seed", s, &c.Seed); err != nil {
				return err
			}
			if !cleanAbs(c.Seed) {
				return profileErr("caches."+name+".seed", "%q is not a clean absolute path", c.Seed)
			}
		}
		p.Caches[name] = c
	}
	return nil
}

// execEnvNameOk refuses a variable name outside the grammar or one World
// owns (§4.3): its own set, srt's, the proxies, git's redirects and every
// registry variable childenv strips.
func execEnvNameOk(field, name string) error {
	upper := strings.ToUpper(name)
	switch {
	case !execEnvNamePattern.MatchString(name):
		return profileErr(field, "%q is not a variable name (%s)", name, execEnvNamePattern)
	case slices.Contains(execWorldEnv, name), name == "SANDBOX_RUNTIME", upper == "NO_PROXY",
		strings.HasSuffix(upper, "_PROXY"), strings.HasPrefix(upper, "GIT_"), slices.Contains(childenv.RegistryVariables, name):
		return profileErr(field, "%q is World-owned; a profile cannot set it", name)
	}
	return nil
}

func parseExecCommand(id string, raw json.RawMessage, path []string) (ExecProfileCommand, error) {
	field := "commands." + id
	if !execCommandIDPattern.MatchString(id) {
		return ExecProfileCommand{}, profileErr("commands", "command id %q does not match %s", id, execCommandIDPattern)
	}
	m, err := strictObject(field, raw, "argv", "flags", "positional", "suffixes", "max_args", "passthrough")
	if err != nil {
		return ExecProfileCommand{}, err
	}
	var c ExecProfileCommand
	if err := decodeField(field+".argv", m["argv"], &c.Argv); err != nil {
		return c, err
	}
	if len(c.Argv) == 0 {
		return c, profileErr(field+".argv", "is empty")
	}
	for _, a := range c.Argv {
		if a == "" || strings.ContainsRune(a, 0) {
			return c, profileErr(field+".argv", "has an empty or NUL-bearing item")
		}
	}
	if raw, ok := m["flags"]; ok {
		if err := parseExecFlags(field+".flags", raw, &c.ExecCommand); err != nil {
			return c, err
		}
	}
	for _, k := range []struct {
		key string
		dst any
	}{{"positional", &c.Positional}, {"suffixes", &c.Suffixes}, {"max_args", &c.MaxArgs}, {"passthrough", &c.Passthrough}} {
		if raw, ok := m[k.key]; ok {
			if err := decodeField(field+"."+k.key, raw, k.dst); err != nil {
				return c, err
			}
		}
	}
	switch c.Positional {
	case "", ExecPositionalRelpath, ExecPositionalPkgPattern, ExecPositionalTestfile:
	default:
		return c, profileErr(field+".positional", "%q is not relpath, pkgpattern or testfile", c.Positional)
	}
	if (c.Positional == ExecPositionalTestfile) != (len(c.Suffixes) > 0) {
		return c, profileErr(field+".suffixes", "are required with, and only with, positional testfile")
	}
	for _, s := range c.Suffixes {
		if s == "" {
			return c, profileErr(field+".suffixes", "has an empty suffix")
		}
	}
	if c.MaxArgs < 0 || c.MaxArgs > execMaxArgs {
		return c, profileErr(field+".max_args", "%d is outside 0–%d", c.MaxArgs, execMaxArgs)
	}
	if c.Passthrough != "" && c.Passthrough != ExecPassthrough {
		return c, profileErr(field+".passthrough", "%q is not %q", c.Passthrough, ExecPassthrough)
	}
	c.Argv0, c.Argv0SHA256, err = resolveArgv0(c.Argv[0], path)
	if err != nil {
		return c, profileErr(field+".argv", "%v", err)
	}
	return c, nil
}

// parseExecFlags reads "flags": each value is a class string, or an object
// {"class": …, "form": "eq"|"sep"}.
func parseExecFlags(field string, raw json.RawMessage, c *ExecCommand) error {
	var flags map[string]json.RawMessage
	if err := decodeField(field, raw, &flags); err != nil {
		return err
	}
	c.Flags = map[string]string{}
	for name, v := range flags {
		if !strings.HasPrefix(name, "-") || name == "--" || strings.Contains(name, "=") || len(name) < 2 {
			return profileErr(field, "%q is not a flag name", name)
		}
		class, form := "", ""
		if json.Unmarshal(v, &class) != nil {
			obj, err := strictObject(field+"."+name, v, "class", "form")
			if err != nil {
				return err
			}
			if err := decodeField(field+"."+name+".class", obj["class"], &class); err != nil {
				return err
			}
			if f, ok := obj["form"]; ok {
				if err := decodeField(field+"."+name+".form", f, &form); err != nil {
					return err
				}
				if form != ExecFormEq && form != ExecFormSep {
					return profileErr(field+"."+name+".form", "%q is not %q or %q", form, ExecFormEq, ExecFormSep)
				}
				if class == ExecClassBool {
					return profileErr(field+"."+name+".form", "a bool flag takes no value, so it has no form")
				}
			}
		}
		if !execClassOk(class) {
			return profileErr(field+"."+name+".class", "%q is not bool, regex, int, relpath or enum:[…]", class)
		}
		c.Flags[name] = class
		if form == ExecFormSep {
			if c.Forms == nil {
				c.Forms = map[string]string{}
			}
			c.Forms[name] = form
		}
	}
	return nil
}

func execClassOk(class string) bool {
	switch class {
	case ExecClassBool, ExecClassRegex, ExecClassInt, ExecClassRelpath:
		return true
	}
	if !strings.HasPrefix(class, "enum:[") || !strings.HasSuffix(class, "]") {
		return false
	}
	words := strings.Split(class[len("enum:["):len(class)-1], ",")
	for _, w := range words {
		if w == "" {
			return false
		}
	}
	return true
}

// EmitExecArgs renders normalized agent args (MatchExecArgs's output) as the
// argv the command receives: a valued flag whose form is "sep" becomes two
// items, everything else is unchanged. After a passthrough "--" nothing is a
// flag.
func EmitExecArgs(c ExecCommand, normalized []string) []string {
	out := make([]string, 0, len(normalized))
	afterDashDash := false
	for _, a := range normalized {
		if !afterDashDash && a == "--" {
			afterDashDash = true
		} else if !afterDashDash && strings.HasPrefix(a, "-") && strings.Contains(a, "=") && c.Forms[flagNameOf(a)] == ExecFormSep {
			out = append(out, flagNameOf(a), flagValueOf(a))
			continue
		}
		out = append(out, a)
	}
	return out
}

// resolveArgv0 resolves argv[0] once: an absolute path as given, a bare name
// on the profile path, to a regular executable file, plus its sha256. The
// trampoline's `env -i K=V…` would read a path holding '=' as an assignment.
func resolveArgv0(name string, path []string) (string, string, error) {
	var candidates []string
	switch {
	case filepath.IsAbs(name):
		candidates = []string{filepath.Clean(name)}
	case strings.Contains(name, "/"):
		return "", "", fmt.Errorf("argv[0] %q is neither a bare name nor an absolute path", name)
	default:
		for _, d := range path {
			candidates = append(candidates, filepath.Join(d, name))
		}
	}
	for _, c := range candidates {
		info, err := os.Stat(c)
		if err != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0o111 == 0 {
			continue
		}
		if strings.Contains(c, "=") {
			return "", "", fmt.Errorf("argv[0] resolves to %q, which holds '='", c)
		}
		sum, err := sha256File(c)
		if err != nil {
			return "", "", err
		}
		return c, sum, nil
	}
	return "", "", fmt.Errorf("argv[0] %q does not resolve on path %v to an executable", name, path)
}

func sha256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func cleanAbs(p string) bool { return filepath.IsAbs(p) && filepath.Clean(p) == p }

// underAny reports p equal to or inside one of roots (lexically).
func underAny(p string, roots []string) bool {
	for _, r := range roots {
		if p == r || strings.HasPrefix(p, strings.TrimSuffix(r, "/")+"/") {
			return true
		}
	}
	return false
}

// sortedKeys returns m's keys in byte order.
func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
