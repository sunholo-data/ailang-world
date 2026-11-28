package appmanifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// exampleFile is the manifest docs/QUICKSTART.md copies; the test reads the same bytes.
func exampleFile(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "docs", "examples", "apps", "docparse-convert.json"))
	if err != nil {
		t.Fatalf("read the shipped example: %v", err)
	}
	return data
}

func exampleManifest(t *testing.T) Manifest {
	t.Helper()
	m, err := Parse(exampleFile(t))
	if err != nil {
		t.Fatalf("parse example: %v", err)
	}
	return m
}

func TestShippedExampleIsValid(t *testing.T) {
	m := exampleManifest(t)
	if ps := Validate(m); len(ps) != 0 {
		t.Fatalf("the shipped example must be valid, got %v", ps)
	}
	if m.Name != "docparse-convert" || len(m.Preset.Grants) == 0 || len(m.Outputs) == 0 {
		t.Fatalf("example lost its shape: %+v", m)
	}
}

func TestParseIsClosed(t *testing.T) {
	base := string(exampleFile(t))
	cases := map[string]string{
		"unknown top-level field": strings.Replace(base, `"schema"`, `"icon": "x", "schema"`, 1),
		"unsandboxed check":       strings.Replace(base, `"name": "ailang-test"`, `"name": "ailang-test", "unsandboxed": true`, 1),
		"repo path in template":   strings.Replace(base, `"agent"`, `"repo": "/Users/x/repo", "agent"`, 1),
		"trailing content":        base + "{}",
	}
	for name, src := range cases {
		if src == base {
			t.Fatalf("%s: the case did not change the example", name)
		}
		if _, err := Parse([]byte(src)); err == nil {
			t.Errorf("%s: Parse accepted it", name)
		}
	}
}

// TestValidateRefusals mutates one field of the valid example per case and requires the
// named field's refusal. Each case is a mutant of the example that the validator kills.
func TestValidateRefusals(t *testing.T) {
	cases := []struct {
		name  string
		field string
		edit  func(m *Manifest)
	}{
		{"schema", "schema", func(m *Manifest) { m.Schema = "world/preset/v1" }},
		{"uppercase name", "name", func(m *Manifest) { m.Name = "Docparse" }},
		{"leading hyphen", "name", func(m *Manifest) { m.Name = "-docparse" }},
		{"reserved name", "name", func(m *Manifest) { m.Name = "se-tools" }},
		{"long name", "name", func(m *Manifest) { m.Name = strings.Repeat("a", 52) }},
		{"empty title", "title", func(m *Manifest) { m.Title = "" }},
		{"long summary", "summary", func(m *Manifest) { m.Summary = strings.Repeat("s", 281) }},
		{"agent", "preset.agent", func(m *Manifest) { m.Preset.Agent = "pi" }},
		{"mode", "preset.launch.mode", func(m *Manifest) { m.Preset.Launch.Mode = "detached" }},
		{"deadline", "preset.launch.deadline_s", func(m *Manifest) { m.Preset.Launch.DeadlineS = 59 }},
		{"ttl", "preset.ttl_s", func(m *Manifest) { m.Preset.TTLS = 604801 }},
		{"no grants", "preset.grants", func(m *Manifest) { m.Preset.Grants = nil }},
		{"invented effect", "preset.grants[0].effect", func(m *Manifest) { m.Preset.Grants[0].Effect = "Host.Blender" }},
		{"model effect", "preset.grants[0].effect", func(m *Manifest) { m.Preset.Grants[0].Effect = "Model.Infer" }},
		{"scope", "preset.grants[0].scope", func(m *Manifest) { m.Preset.Grants[0].Scope = "/tmp" }},
		{"zero budget", "preset.grants[0].budget", func(m *Manifest) { m.Preset.Grants[0].Budget = 0 }},
		{"huge budget", "preset.grants[0].budget", func(m *Manifest) { m.Preset.Grants[0].Budget = 100001 }},
		{"repeated grant", "preset.grants[1]", func(m *Manifest) { m.Preset.Grants[1] = m.Preset.Grants[0] }},
		{"exec without profile", "preset.exec_profile", func(m *Manifest) { m.Preset.ExecProfile = "" }},
		{"check argv", "preset.checks[0].argv", func(m *Manifest) { m.Preset.Checks[0].Argv = nil }},
		{"check cwd escapes", "preset.checks[0].cwd", func(m *Manifest) { m.Preset.Checks[0].Cwd = "../x" }},
		{"check timeout", "preset.checks[0].timeout_s", func(m *Manifest) { m.Preset.Checks[0].TimeoutS = 3601 }},
		{"output kind", "outputs[0].kind", func(m *Manifest) { m.Outputs[0].Kind = "pr" }},
		{"absolute glob", "outputs[0].glob", func(m *Manifest) { m.Outputs[0].Glob = "/etc/*" }},
		{"escaping glob", "outputs[0].glob", func(m *Manifest) { m.Outputs[0].Glob = "out/../../x" }},
		{"bad pattern", "outputs[0].glob", func(m *Manifest) { m.Outputs[0].Glob = "out/[x" }},
		{"output names no check", "outputs[0].check", func(m *Manifest) { m.Outputs[0].Check = "missing" }},
		{"view outside catalog", "views[0]", func(m *Manifest) { m.Views[0] = "modal" }},
		{"repeated view", "views[1]", func(m *Manifest) { m.Views[1] = m.Views[0] }},
	}
	for _, c := range cases {
		m := exampleManifest(t)
		c.edit(&m)
		ps := Validate(m)
		found := false
		for _, p := range ps {
			if p.Field == c.field {
				found = true
			}
		}
		if !found {
			t.Errorf("%s: want a refusal on %s, got %v", c.name, c.field, ps)
		}
	}
}

func TestNarrowGrants(t *testing.T) {
	tmpl := []Grant{{"Workspace.Read", "worktree", 2000}, {"Workspace.Write", "worktree", 500}}
	if err := NarrowGrants(tmpl, []Grant{{"Workspace.Read", "worktree", 2000}, {"Workspace.Write", "worktree", 1}}); err != nil {
		t.Fatalf("narrowing within the template refused: %v", err)
	}
	for name, req := range map[string][]Grant{
		"wider budget":  {{"Workspace.Write", "worktree", 501}},
		"zero budget":   {{"Workspace.Read", "worktree", 0}},
		"new effect":    {{"Workspace.Exec", "worktree", 1}},
		"another scope": {{"Workspace.Read", "/tmp", 1}},
	} {
		if err := NarrowGrants(tmpl, req); err == nil {
			t.Errorf("%s: NarrowGrants accepted a widening", name)
		}
	}
}

func TestReadCatalog(t *testing.T) {
	if got, err := ReadCatalog(filepath.Join(t.TempDir(), "absent")); err != nil || len(got) != 0 {
		t.Fatalf("a missing directory is an empty catalog, got %v %v", got, err)
	}
	dir := t.TempDir()
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	write("docparse-convert.json", string(exampleFile(t)))
	write("broken.json", "{")
	write("renamed.json", string(exampleFile(t)))
	write("notes.txt", "ignored")
	if err := os.Mkdir(filepath.Join(dir, "sub.json"), 0o700); err != nil {
		t.Fatal(err)
	}
	got, err := ReadCatalog(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got[0].Name != "broken" || got[1].Name != "docparse-convert" || got[2].Name != "renamed" {
		t.Fatalf("catalog rows %+v", got)
	}
	if got[0].Error == "" || got[0].Valid() {
		t.Errorf("broken.json must carry a parse error: %+v", got[0])
	}
	if !got[1].Valid() || !strings.HasPrefix(got[1].Digest, "sha256:") || len(got[1].Digest) != 71 {
		t.Errorf("the example must be valid with a sha256 digest: %+v", got[1])
	}
	if got[2].Valid() || len(got[2].Problems) != 1 || got[2].Problems[0].Field != "name" {
		t.Errorf("a file named unlike its manifest must be refused on name: %+v", got[2])
	}
}
