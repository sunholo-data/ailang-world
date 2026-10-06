package broker

// Row 140 M2 (w-workspace-exec-toolchain-effect §4.3): the exec profile type,
// its in-memory loader and the emitted argv's per-flag form. The `serve`
// flags, the placement rules and the read_roots ancestor rules are M3.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// fakeToolDir holds an executable per name, for argv0 resolution.
func fakeToolDir(t *testing.T, names ...string) string {
	t.Helper()
	dir := t.TempDir()
	for _, n := range names {
		if err := os.WriteFile(filepath.Join(dir, n), []byte("#!/bin/sh\necho "+n+"\n"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	real, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}
	return real
}

// profileJSON renders a profile object with one substituted path list.
func profileJSON(t *testing.T, bin string, mutate func(m map[string]any)) []byte {
	t.Helper()
	m := map[string]any{
		"profile":    ExecProfileVersion,
		"project":    "ailang-compiler",
		"root":       ".",
		"path":       []any{bin, "/usr/bin", "/bin"},
		"env":        map[string]any{"GOFLAGS": "-mod=readonly", "GOPROXY": "off", "GOTOOLCHAIN": "go1.26.6", "GOTELEMETRY": "off"},
		"read_roots": []any{"/opt/modcache"},
		"caches":     map[string]any{"GOCACHE": map[string]any{"seed": "/opt/seeds/go-build"}},
		"timeout_ms": 9000,
		"commands": map[string]any{
			"test": map[string]any{"argv": []any{"go", "test", "-count=1"},
				"flags":      map[string]any{"-run": "regex", "-v": "bool", "-short": "bool"},
				"positional": "pkgpattern", "max_args": 8},
			"vet": map[string]any{"argv": []any{"go", "vet"}, "positional": "pkgpattern", "max_args": 4},
		},
	}
	if mutate != nil {
		mutate(m)
	}
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestParseExecProfileDesignExample(t *testing.T) {
	bin := fakeToolDir(t, "go")
	data := profileJSON(t, bin, nil)
	p, err := ParseExecProfile(data)
	if err != nil {
		t.Fatalf("ParseExecProfile: %v", err)
	}
	sum := sha256.Sum256(data)
	if p.Digest != "sha256:"+hex.EncodeToString(sum[:]) {
		t.Fatalf("Digest = %q, want the sha256 of the bytes read", p.Digest)
	}
	goBytes, _ := os.ReadFile(filepath.Join(bin, "go"))
	goSum := sha256.Sum256(goBytes)
	test := p.Commands["test"]
	if test.Argv0 != filepath.Join(bin, "go") || test.Argv0SHA256 != hex.EncodeToString(goSum[:]) {
		t.Fatalf("test argv0 = %q %q, want %q and its sha256", test.Argv0, test.Argv0SHA256, filepath.Join(bin, "go"))
	}
	if !reflect.DeepEqual(test.Argv, []string{"go", "test", "-count=1"}) || test.MaxArgs != 8 || test.Flags["-run"] != ExecClassRegex {
		t.Fatalf("test command = %+v", test)
	}
	if p.Project != "ailang-compiler" || p.Root != "." || p.TimeoutMS != 9000 || p.Caches["GOCACHE"].Seed != "/opt/seeds/go-build" ||
		p.Env["GOPROXY"] != "off" || !reflect.DeepEqual(p.ReadRoots, []string{"/opt/modcache"}) {
		t.Fatalf("profile = %+v", p)
	}
}

func TestParseExecProfileRefusals(t *testing.T) {
	bin := fakeToolDir(t, "go")
	cmd := func(m map[string]any) map[string]any { return m["commands"].(map[string]any)["test"].(map[string]any) }
	for _, tc := range []struct {
		name, want string
		mutate     func(m map[string]any)
	}{
		{"unknown top-level key", `unknown key "network"`, func(m map[string]any) { m["network"] = []any{} }},
		{"case-aliased key", `unknown key "Project"`, func(m map[string]any) { m["Project"] = "x" }},
		{"unknown command key", `unknown key "cwd"`, func(m map[string]any) { cmd(m)["cwd"] = "/" }},
		{"unknown flag-object key", `unknown key "default"`, func(m map[string]any) {
			cmd(m)["flags"].(map[string]any)["-k"] = map[string]any{"class": "regex", "default": "x"}
		}},
		{"wrong profile version", "profile", func(m map[string]any) { m["profile"] = "world/exec-profile/v2" }},
		{"bad project", "project", func(m map[string]any) { m["project"] = "Ailang" }},
		{"root escapes", "root", func(m map[string]any) { m["root"] = "../x" }},
		{"root absolute", "root", func(m map[string]any) { m["root"] = "/x" }},
		{"path relative", "path", func(m map[string]any) { m["path"] = []any{"bin"} }},
		{"env bad name", "env", func(m map[string]any) { m["env"] = map[string]any{"A-B": "1"} }},
		{"env HOME", `"HOME"`, func(m map[string]any) { m["env"] = map[string]any{"HOME": "/x"} }},
		{"env TMPDIR", `"TMPDIR"`, func(m map[string]any) { m["env"] = map[string]any{"TMPDIR": "/x"} }},
		{"env PATH", `"PATH"`, func(m map[string]any) { m["env"] = map[string]any{"PATH": "/x"} }},
		{"env proxy", `"http_proxy"`, func(m map[string]any) { m["env"] = map[string]any{"http_proxy": "x"} }},
		{"env NO_PROXY", `"NO_PROXY"`, func(m map[string]any) { m["env"] = map[string]any{"NO_PROXY": "x"} }},
		{"env GIT_", `"GIT_DIR"`, func(m map[string]any) { m["env"] = map[string]any{"GIT_DIR": "x"} }},
		{"env SANDBOX_RUNTIME", `"SANDBOX_RUNTIME"`, func(m map[string]any) { m["env"] = map[string]any{"SANDBOX_RUNTIME": "0"} }},
		{"env pycache", `"PYTHONPYCACHEPREFIX"`, func(m map[string]any) { m["env"] = map[string]any{"PYTHONPYCACHEPREFIX": "x"} }},
		{"env registry credential", `"AILANG_REGISTRY_API_KEY"`, func(m map[string]any) {
			m["env"] = map[string]any{"AILANG_REGISTRY_API_KEY": "x"}
		}},
		{"env names a cache variable", `"GOCACHE"`, func(m map[string]any) { m["env"] = map[string]any{"GOCACHE": "/x"} }},
		{"GOMODCACHE outside read_roots", "GOMODCACHE", func(m map[string]any) { m["env"] = map[string]any{"GOMODCACHE": "/elsewhere"} }},
		{"cache name World-owned", `"HOME"`, func(m map[string]any) { m["caches"] = map[string]any{"HOME": map[string]any{}} }},
		{"cache seed relative", "seed", func(m map[string]any) { m["caches"] = map[string]any{"GOCACHE": map[string]any{"seed": "seeds"}} }},
		{"read_root relative", "read_roots", func(m map[string]any) { m["read_roots"] = []any{"mod"} }},
		{"timeout zero", "timeout_ms", func(m map[string]any) { m["timeout_ms"] = 0 }},
		{"timeout over the handler cap", "timeout_ms", func(m map[string]any) { m["timeout_ms"] = 9001 }},
		{"command id", `"Test"`, func(m map[string]any) { m["commands"] = map[string]any{"Test": cmd(m)} }},
		{"no commands", "commands", func(m map[string]any) { m["commands"] = map[string]any{} }},
		{"empty argv", "argv", func(m map[string]any) { cmd(m)["argv"] = []any{} }},
		{"argv0 not on path", "does not resolve", func(m map[string]any) { cmd(m)["argv"] = []any{"cargo", "test"} }},
		{"argv0 relative with a slash", "argv", func(m map[string]any) { cmd(m)["argv"] = []any{"./go"} }},
		{"flag class", "class", func(m map[string]any) { cmd(m)["flags"].(map[string]any)["-x"] = "float" }},
		{"flag name not a flag", "flag", func(m map[string]any) { cmd(m)["flags"].(map[string]any)["run"] = "regex" }},
		{"flag name is --", "flag", func(m map[string]any) { cmd(m)["flags"].(map[string]any)["--"] = "bool" }},
		{"unknown form", `"space"`, func(m map[string]any) {
			cmd(m)["flags"].(map[string]any)["-k"] = map[string]any{"class": "regex", "form": "space"}
		}},
		{"form on a bool flag", "form", func(m map[string]any) {
			cmd(m)["flags"].(map[string]any)["-v"] = map[string]any{"class": "bool", "form": "sep"}
		}},
		{"positional class", "positional", func(m map[string]any) { cmd(m)["positional"] = "glob" }},
		{"testfile without suffixes", "suffixes", func(m map[string]any) { cmd(m)["positional"] = "testfile" }},
		{"max_args over 16", "max_args", func(m map[string]any) { cmd(m)["max_args"] = 17 }},
		{"passthrough other than --", "passthrough", func(m map[string]any) { cmd(m)["passthrough"] = "++" }},
		{"not an object", "", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := profileJSON(t, bin, tc.mutate)
			if tc.mutate == nil {
				data = []byte(`["not","an","object"]`)
			}
			p, err := ParseExecProfile(data)
			if err == nil {
				t.Fatalf("ParseExecProfile accepted it: %+v", p)
			}
			var perr *ExecProfileError
			if !errors.As(err, &perr) || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("err = %v, want an *ExecProfileError naming %q", err, tc.want)
			}
		})
	}
}

// [M2, ruled in session 2026-10-05] Matching stays on the normalized form;
// the EMITTED argv renders each valued flag per its form: "eq" (the default)
// as one item -f=v, "sep" as two items -f, v. argparse reads `-k=expr` as the
// value "=expr", so a pytest -k must be profiled "sep".
func TestEmitExecArgsForms(t *testing.T) {
	bin := fakeToolDir(t, "pytest", "go")
	data := profileJSON(t, bin, func(m map[string]any) {
		m["commands"] = map[string]any{
			"test-file": map[string]any{"argv": []any{"pytest", "-q"},
				"flags": map[string]any{"-k": map[string]any{"class": "regex", "form": "sep"}, "-x": "bool",
					"--maxfail": map[string]any{"class": "int", "form": "eq"}},
				"positional": "testfile", "suffixes": []any{".py"}, "max_args": 6},
			"test": map[string]any{"argv": []any{"go", "test"}, "flags": map[string]any{"-run": "regex"},
				"positional": "pkgpattern", "max_args": 4, "passthrough": "--"},
		}
	})
	p, err := ParseExecProfile(data)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		command               string
		args, normalized, out []string
	}{
		{"test-file", []string{"-k", "foo and not bar", "-x", "tests/test_a.py"},
			[]string{"-k=foo and not bar", "-x", "tests/test_a.py"}, []string{"-k", "foo and not bar", "-x", "tests/test_a.py"}},
		{"test-file", []string{"-k=a=b", "--maxfail", "2"},
			[]string{"-k=a=b", "--maxfail=2"}, []string{"-k", "a=b", "--maxfail=2"}},
		{"test", []string{"-run", "TestX", "./a/..."},
			[]string{"-run=TestX", "./a/..."}, []string{"-run=TestX", "./a/..."}},
		{"test", []string{"--", "./b"}, []string{"--", "./b"}, []string{"--", "./b"}},
	} {
		c := p.Commands[tc.command]
		norm, err := MatchExecArgs(c.ExecCommand, tc.args)
		if err != nil {
			t.Fatalf("%s %q: %v", tc.command, tc.args, err)
		}
		if !reflect.DeepEqual(norm, tc.normalized) {
			t.Fatalf("%s %q normalized = %q, want %q", tc.command, tc.args, norm, tc.normalized)
		}
		if got := EmitExecArgs(c.ExecCommand, norm); !reflect.DeepEqual(got, tc.out) {
			t.Fatalf("%s %q emitted = %q, want %q", tc.command, tc.args, got, tc.out)
		}
	}
}
