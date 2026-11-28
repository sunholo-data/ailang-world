// Package appmanifest reads and validates World app manifests (queue row 274 M1,
// design_docs/planned/w-world-applications.md §2).
//
// An app is an extension package with an app.json manifest (schema world/app/v1): a
// portable preset template of existing effects only, declared outputs, and views drawn
// from row 151's closed catalog. The semantic source of truth is
// design_docs/sketches/appmanifest.ail, whose laws carry Z3-proven contracts. This
// file mirrors each law under the same name, and drift_test.go replays every inline
// test row of the sketch through these mirrors.
//
// The package is pure apart from ReadCatalog, which reads one directory. It holds no
// authority: a manifest only describes a launch, and M2 resolves it into a preset that
// the existing start path reviews and confirms.
package appmanifest

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// Schema is the manifest's schema string (schemaOk).
const Schema = "world/app/v1"

// DefaultCheckTimeoutS is the preset grader default a check without timeout_s gets.
const DefaultCheckTimeoutS = 900

// Grant is one template grant: an existing effect, the worktree scope and a budget.
type Grant struct {
	Effect string `json:"effect"`
	Scope  string `json:"scope"`
	Budget int64  `json:"budget"`
}

// Launch is the template's launch mode and deadline.
type Launch struct {
	Mode      string `json:"mode"`
	DeadlineS int64  `json:"deadline_s"`
}

// Check is a grader in the preset shape, without `unsandboxed`.
type Check struct {
	Name     string   `json:"name"`
	Argv     []string `json:"argv"`
	Cwd      string   `json:"cwd,omitempty"`
	TimeoutS int64    `json:"timeout_s,omitempty"`
}

// Template is the portable part of a preset. repo, remote, base, profile and daemon
// come from the operator's repo profile at launch, never from the app.
type Template struct {
	Agent       string  `json:"agent"`
	Launch      Launch  `json:"launch"`
	TTLS        int64   `json:"ttl_s"`
	Grants      []Grant `json:"grants"`
	ExecProfile string  `json:"exec_profile,omitempty"`
	Checks      []Check `json:"checks,omitempty"`
}

// Output names a file the app is for, and optionally the check that judges it.
type Output struct {
	Kind  string `json:"kind"`
	Glob  string `json:"glob"`
	Check string `json:"check,omitempty"`
}

// Manifest is one app.json.
type Manifest struct {
	Schema   string   `json:"schema"`
	Name     string   `json:"name"`
	Title    string   `json:"title"`
	Summary  string   `json:"summary"`
	RepoKind string   `json:"repo_kind,omitempty"`
	Preset   Template `json:"preset"`
	Outputs  []Output `json:"outputs,omitempty"`
	Views    []string `json:"views,omitempty"`
}

// Problem is one refusal: the field it is about and why.
type Problem struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func (p Problem) String() string { return p.Field + ": " + p.Reason }

// ── Mirrors of design_docs/sketches/appmanifest.ail (same names, same bounds) ──

func schemaOk(s string) bool      { return s == Schema }
func nameLenOk(n int) bool        { return n >= 1 && n <= 51 }
func nameReserved(s string) bool  { return s == "se-tools" }
func titleLenOk(n int) bool       { return n >= 1 && n <= 80 }
func summaryLenOk(n int) bool     { return n >= 1 && n <= 280 }
func agentOk(a string) bool       { return a == "claude" || a == "codex" || a == "agy" }
func modeOk(m string) bool        { return m == "headless" || m == "interactive" }
func deadlineOk(s int64) bool     { return s >= 60 && s <= 86400 }
func ttlOk(s int64) bool          { return s >= 60 && s <= 604800 }
func scopeOk(s string) bool       { return s == "worktree" }
func budgetOk(n int64) bool       { return n >= 1 && n <= 100000 }
func grantCountOk(n int) bool     { return n >= 1 && n <= 16 }
func listCountOk(n int) bool      { return n >= 0 && n <= 16 }
func checkNameLenOk(n int) bool   { return n >= 1 && n <= 40 }
func checkTimeoutOk(s int64) bool { return s >= 1 && s <= 3600 }
func outputKindOk(k string) bool  { return k == "file" }
func viewCountOk(n int) bool      { return n >= 0 && n <= 8 }
func grantNarrows(templateBudget, requested int64) bool {
	return requested >= 1 && requested <= templateBudget
}

func nameByteOk(b int, first bool) bool {
	return (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') || (!first && b == '-')
}

// allowedEffects is effectAllowed's set. allowed_effects_test.go binds every entry to a
// broker effect constant, so the list cannot name an effect the broker lacks.
var allowedEffects = map[string]bool{
	"Workspace.Read": true, "Workspace.Write": true, "Workspace.Delete": true,
	"Workspace.Exec": true, "Ailang.Check": true, "Ailang.Run": true, "Ailang.RunEnv": true,
	"Ailang.RunNet": true, "Ailang.Discover": true, "Ailang.CLI": true,
}

func effectAllowed(e string) bool { return allowedEffects[e] }

// viewCatalog is viewOk's set: row 151's closed catalog.
var viewCatalog = map[string]bool{
	"text": true, "row": true, "column": true, "list": true, "table": true, "card": true,
	"checklist": true, "timeline": true, "chart": true, "media": true, "output-link": true,
	"object-link": true, "grade-badge": true, "choice": true, "action-button": true,
}

func viewOk(v string) bool { return viewCatalog[v] }

func noDotDotSegment(p string) bool {
	return p != ".." && !strings.HasPrefix(p, "../") && !strings.HasSuffix(p, "/..") &&
		!strings.Contains(p, "/../")
}

func outputGlobOk(g string) bool {
	n := len(g)
	return n >= 1 && n <= 256 && !strings.HasPrefix(g, "/") && !strings.HasPrefix(g, "-") &&
		!strings.Contains(g, `\`) && noDotDotSegment(g)
}

// NameOk applies nameLenOk, nameByteOk to every byte, and nameReserved.
func NameOk(name string) bool {
	if !nameLenOk(len(name)) || nameReserved(name) {
		return false
	}
	for i := 0; i < len(name); i++ {
		if !nameByteOk(int(name[i]), i == 0) {
			return false
		}
	}
	return true
}

// ── Decoding and validation ──

// Parse decodes one manifest strictly: unknown fields anywhere are refused, and so is
// trailing content after the object.
func Parse(data []byte) (Manifest, error) {
	var m Manifest
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&m); err != nil {
		return Manifest{}, fmt.Errorf("app.json: %w", err)
	}
	if _, err := d.Token(); !errors.Is(err, io.EOF) {
		return Manifest{}, errors.New("app.json: content after the manifest object")
	}
	return m, nil
}

// Validate applies every law to m and returns every problem found, in field order.
// An empty result means the manifest is valid.
func Validate(m Manifest) []Problem {
	var ps []Problem
	add := func(field, reason string) { ps = append(ps, Problem{field, reason}) }

	if !schemaOk(m.Schema) {
		add("schema", fmt.Sprintf("is %q, not %s", m.Schema, Schema))
	}
	if !NameOk(m.Name) {
		add("name", fmt.Sprintf("%q must match ^[a-z0-9][a-z0-9-]*$, be 1 to 51 characters, and not be se-tools", m.Name))
	}
	if !titleLenOk(len(m.Title)) {
		add("title", "must be 1 to 80 bytes")
	}
	if !summaryLenOk(len(m.Summary)) {
		add("summary", "must be 1 to 280 bytes")
	}
	if m.RepoKind != "" && !NameOk(m.RepoKind) {
		add("repo_kind", fmt.Sprintf("%q must use the name grammar", m.RepoKind))
	}

	t := m.Preset
	if !agentOk(t.Agent) {
		add("preset.agent", fmt.Sprintf("%q is not claude, codex or agy", t.Agent))
	}
	if !modeOk(t.Launch.Mode) {
		add("preset.launch.mode", fmt.Sprintf("%q is not headless or interactive", t.Launch.Mode))
	}
	if !deadlineOk(t.Launch.DeadlineS) {
		add("preset.launch.deadline_s", "must be 60 to 86400 seconds")
	}
	if !ttlOk(t.TTLS) {
		add("preset.ttl_s", "must be 60 to 604800 seconds")
	}
	if !grantCountOk(len(t.Grants)) {
		add("preset.grants", "must hold 1 to 16 grants")
	}
	seenGrant := map[string]bool{}
	hasExec := false
	for i, g := range t.Grants {
		f := fmt.Sprintf("preset.grants[%d]", i)
		if !effectAllowed(g.Effect) {
			add(f+".effect", fmt.Sprintf("%q is not an existing episode effect", g.Effect))
		}
		if !scopeOk(g.Scope) {
			add(f+".scope", fmt.Sprintf("%q is not worktree", g.Scope))
		}
		if !budgetOk(g.Budget) {
			add(f+".budget", "must be 1 to 100000")
		}
		key := g.Effect + "=" + g.Scope
		if seenGrant[key] {
			add(f, "repeats "+key)
		}
		seenGrant[key] = true
		if g.Effect == "Workspace.Exec" {
			hasExec = true
		}
	}
	if hasExec && t.ExecProfile == "" {
		add("preset.exec_profile", "Workspace.Exec needs an exec profile name")
	}
	if t.ExecProfile != "" && !NameOk(t.ExecProfile) {
		add("preset.exec_profile", fmt.Sprintf("%q must use the name grammar", t.ExecProfile))
	}

	if !listCountOk(len(t.Checks)) {
		add("preset.checks", "must hold at most 16 checks")
	}
	checks := map[string]bool{}
	for i, c := range t.Checks {
		f := fmt.Sprintf("preset.checks[%d]", i)
		if !checkNameLenOk(len(c.Name)) || !NameOk(c.Name) {
			add(f+".name", fmt.Sprintf("%q must use the name grammar, 1 to 40 characters", c.Name))
		}
		if checks[c.Name] {
			add(f+".name", "repeats "+c.Name)
		}
		checks[c.Name] = true
		if len(c.Argv) == 0 {
			add(f+".argv", "must be a non-empty list")
		}
		for j, a := range c.Argv {
			if a == "" {
				add(fmt.Sprintf("%s.argv[%d]", f, j), "must not be empty")
			}
		}
		if c.Cwd != "" && (!outputGlobOk(c.Cwd) || strings.ContainsAny(c.Cwd, "*?[")) {
			add(f+".cwd", fmt.Sprintf("%q must be a relative path inside the worktree", c.Cwd))
		}
		timeout := c.TimeoutS
		if timeout == 0 {
			timeout = DefaultCheckTimeoutS
		}
		if !checkTimeoutOk(timeout) {
			add(f+".timeout_s", "must be 1 to 3600 seconds")
		}
	}

	if !listCountOk(len(m.Outputs)) {
		add("outputs", "must hold at most 16 outputs")
	}
	for i, o := range m.Outputs {
		f := fmt.Sprintf("outputs[%d]", i)
		if !outputKindOk(o.Kind) {
			add(f+".kind", fmt.Sprintf("%q is not file", o.Kind))
		}
		if !outputGlobOk(o.Glob) {
			add(f+".glob", fmt.Sprintf("%q must be a relative pattern inside the worktree", o.Glob))
		} else if _, err := path.Match(strings.ReplaceAll(o.Glob, "**", "*"), ""); err != nil {
			add(f+".glob", fmt.Sprintf("%q is not a valid pattern: %v", o.Glob, err))
		}
		if o.Check != "" && !checks[o.Check] {
			add(f+".check", fmt.Sprintf("%q names no check in preset.checks", o.Check))
		}
	}

	if !viewCountOk(len(m.Views)) {
		add("views", "must name at most 8 views")
	}
	seenView := map[string]bool{}
	for i, v := range m.Views {
		f := fmt.Sprintf("views[%d]", i)
		if !viewOk(v) {
			add(f, fmt.Sprintf("%q is not in row 151's catalog", v))
		}
		if seenView[v] {
			add(f, "repeats "+v)
		}
		seenView[v] = true
	}
	return ps
}

// NarrowGrants checks that every requested grant narrows a template grant: the same
// effect and scope, with a budget from 1 to the template's. It is M2's launch law.
func NarrowGrants(template, requested []Grant) error {
	byKey := map[string]int64{}
	for _, g := range template {
		byKey[g.Effect+"="+g.Scope] = g.Budget
	}
	for _, r := range requested {
		key := r.Effect + "=" + r.Scope
		tb, ok := byKey[key]
		if !ok {
			return fmt.Errorf("grant %s is not in the app's template; an app cannot widen its grants", key)
		}
		if !grantNarrows(tb, r.Budget) {
			return fmt.Errorf("grant %s budget %d must be 1 to the template's %d", key, r.Budget, tb)
		}
	}
	return nil
}

// ── The local catalog ──

// Entry is one row of the local app catalog.
type Entry struct {
	Name     string    `json:"name"`
	File     string    `json:"file"`
	Digest   string    `json:"digest,omitempty"`
	Manifest *Manifest `json:"manifest,omitempty"`
	Problems []Problem `json:"problems,omitempty"`
	Error    string    `json:"error,omitempty"`
}

// Valid reports whether the entry parsed and passed every law.
func (e Entry) Valid() bool { return e.Error == "" && len(e.Problems) == 0 && e.Manifest != nil }

// DefaultDir is the local app catalog: ~/.ailang-world/apps.
func DefaultDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ailang-world", "apps")
}

// ReadCatalog reads every <name>.json in dir, sorted by name. A missing directory is an
// empty catalog. A file that fails to read or parse becomes an entry with Error; a file
// whose name differs from its manifest's name, or that breaks a law, carries Problems.
func ReadCatalog(dir string) ([]Entry, error) {
	files, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return []Entry{}, nil
	}
	if err != nil {
		return nil, err
	}
	out := []Entry{}
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		out = append(out, ReadEntry(filepath.Join(dir, f.Name())))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// ReadEntry reads and validates one catalog file.
func ReadEntry(file string) Entry {
	e := Entry{Name: strings.TrimSuffix(filepath.Base(file), ".json"), File: file}
	data, err := os.ReadFile(file)
	if err != nil {
		e.Error = err.Error()
		return e
	}
	sum := sha256.Sum256(data)
	e.Digest = "sha256:" + hex.EncodeToString(sum[:])
	m, err := Parse(data)
	if err != nil {
		e.Error = err.Error()
		return e
	}
	e.Manifest = &m
	e.Problems = Validate(m)
	if m.Name != e.Name {
		e.Problems = append(e.Problems, Problem{"name", fmt.Sprintf("the file is %s.json but the manifest says %q", e.Name, m.Name)})
	}
	return e
}
