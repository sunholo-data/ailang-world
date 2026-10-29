package surfacefollow

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigRefusals(t *testing.T) {
	// LoadConfig canonicalises repo paths, and on macOS t.TempDir() is under
	// /var, a symlink to /private/var: resolve it so the expected path matches.
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	configDir := filepath.Join(root, "operator")
	home := filepath.Join(root, "home")
	workspace := filepath.Join(root, "workspace")
	for _, p := range []string{configDir, filepath.Join(home, "one"), filepath.Join(home, "two"), workspace} {
		if err := os.MkdirAll(p, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	path := filepath.Join(configDir, "follow.json")
	base := func() map[string]any {
		return map[string]any{"schema": ConfigSchema, "repos": []any{
			map[string]any{"id": "world", "repo": "~/one", "charter": "docs/charter.md", "branch": "dev"},
			map[string]any{"id": "tools", "repo": filepath.Join(home, "two"), "charter": "charter.md", "branch": "main"},
		}}
	}
	write := func(c map[string]any) {
		t.Helper()
		b, _ := json.Marshal(c)
		if err := os.WriteFile(path, b, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	opts := LoadOptions{Home: home, ForbiddenRoots: []string{workspace}}
	write(base())
	c, err := LoadConfig(path, opts)
	if err != nil || c.Repos[0].Limit != 200 || c.Repos[0].Repo != filepath.Join(home, "one") {
		t.Fatalf("valid two-repo tilde config: %v", err)
	}
	type refusal struct {
		name, field string
		mutate      func(map[string]any)
	}
	cases := []refusal{
		{"schema", "schema", func(c map[string]any) { c["schema"] = "secret-token" }},
		{"unknown", "json", func(c map[string]any) { c["secret-token"] = "secret-token" }},
		{"zero", "repos", func(c map[string]any) { c["repos"] = []any{} }},
		{"five", "repos", func(c map[string]any) {
			r := c["repos"].([]any)
			c["repos"] = append(r, r[0], r[0], r[0])
		}},
	}
	repoCase := func(name, field string, value any) {
		cases = append(cases, refusal{name, field, func(c map[string]any) {
			c["repos"].([]any)[0].(map[string]any)[field] = value
		}})
	}
	repoCase("bad-id", "id", "secret-token!")
	cases = append(cases, refusal{"duplicate-id", "id", func(c map[string]any) {
		c["repos"].([]any)[1].(map[string]any)["id"] = "world"
	}})
	repoCase("relative-repo", "repo", "secret-token")
	repoCase("missing-repo", "repo", filepath.Join(home, "missing"))
	cases = append(cases, refusal{"duplicate-repo", "repo", func(c map[string]any) {
		c["repos"].([]any)[1].(map[string]any)["repo"] = filepath.Join(home, "one")
	}})
	link := filepath.Join(home, "alias")
	if err := os.Symlink(filepath.Join(home, "two"), link); err != nil {
		t.Fatal(err)
	}
	repoCase("canonical-duplicate", "repo", link)
	for _, s := range []string{"", "../charter", "a/../charter", "/charter", "a//charter", "./charter", "a\\b"} {
		repoCase("charter-"+s, "charter", s)
	}
	for _, s := range []string{"", "-dev", "a..b", "a@{b", "a b", "a\tb", "a\x00b", "a~b", "a^b", "a:b", "a?b", "a*b", "a[b", "a\\b", "dev.lock", "dev/", "a//b", ".a", "a."} {
		repoCase("branch-"+s, "branch", s)
	}
	for _, n := range []int{-1, 0, 2501, 1900} {
		repoCase("limit", "limit", n)
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			x := base()
			tc.mutate(x)
			write(x)
			_, err := LoadConfig(path, opts)
			if err == nil || !strings.Contains(err.Error(), "field "+tc.field) ||
				!strings.Contains(err.Error(), "ailang-worldd serve --help") || strings.Contains(err.Error(), "secret-token") {
				t.Fatalf("unsafe or incomplete refusal: %v", err)
			}
			t.Log(err)
		})
	}
	write(base())
	for _, mode := range []os.FileMode{0o644, 0o660, 0o400} {
		_ = os.Chmod(path, mode)
		if _, err := LoadConfig(path, opts); err == nil {
			t.Fatal("unsafe mode accepted")
		}
	}
	_ = os.Chmod(path, 0o600)
	_ = os.Chmod(configDir, 0o770)
	if _, err := LoadConfig(path, opts); err == nil {
		t.Fatal("writable directory accepted")
	}
	_ = os.Chmod(configDir, 0o700)
	sym := filepath.Join(configDir, "link")
	_ = os.Symlink(path, sym)
	if _, err := LoadConfig(sym, opts); err == nil {
		t.Fatal("symlink config accepted")
	}
	if _, err := LoadConfig(configDir, opts); err == nil {
		t.Fatal("directory config accepted")
	}
	if _, err := LoadConfig(filepath.Join(configDir, "missing"), opts); err == nil {
		t.Fatal("missing file accepted")
	}
	if _, err := LoadConfig("relative", opts); err == nil {
		t.Fatal("relative config accepted")
	}
	for _, forbidden := range []string{configDir, root} {
		if _, err := LoadConfig(path, LoadOptions{Home: home, ForbiddenRoots: []string{forbidden}}); err == nil {
			t.Fatal("forbidden file/directory accepted")
		}
	}
	parentAlias := filepath.Join(root, "operator-alias")
	_ = os.Symlink(configDir, parentAlias)
	if _, err := LoadConfig(filepath.Join(parentAlias, "follow.json"),
		LoadOptions{Home: home, ForbiddenRoots: []string{configDir}}); err == nil {
		t.Fatal("symlink parent into workspace accepted")
	}
	// Ownership is checked against real OS metadata. A foreign uid requires
	// chown authority; exercise it where available and explicitly report otherwise.
	if err := os.Chown(path, os.Getuid()+1, -1); err == nil {
		if _, err := LoadConfig(path, opts); err == nil {
			t.Fatal("foreign owner accepted")
		}
		_ = os.Chown(path, os.Getuid(), -1)
	} else {
		t.Log("foreign-owner fixture unavailable without chown authority; owned() is still enforced")
	}
}
