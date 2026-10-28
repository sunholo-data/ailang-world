package main

// Row 150 M8a (w-marketing-capture-sprint-plan §M8a, orchestration §5): the
// reproduction docs, captions, website page and capture-manifest schema are
// wired to the six-shot inventory in png.go. The docs gate proves wiring only;
// real asset resolution runs against a promoted directory (--assets) or, in
// release mode, against the attended capture-manifest.json and its bytes.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

var docsFlagToken = regexp.MustCompile(`^--[a-z][a-z0-9-]*$`)

// docsRepoFile reads one repository file relative to the module root.
func docsRepoFile(t *testing.T, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", filepath.FromSlash(rel)))
	pngOK(t, err)
	return string(b)
}

// docsOverlay reads the real repository, except for the in-memory edits.
func docsOverlay(edits map[string]string) func(string) ([]byte, error) {
	return func(rel string) ([]byte, error) {
		if text, ok := edits[rel]; ok {
			return []byte(text), nil
		}
		return os.ReadFile(filepath.Join("..", "..", filepath.FromSlash(rel)))
	}
}

// docsWantRefusal requires a typed refusal with the code and, when given, a detail substring.
func docsWantRefusal(t *testing.T, err error, code, detail string) {
	t.Helper()
	var r *Refusal
	if !errors.As(err, &r) || r.Code != code || !strings.Contains(r.Detail, detail) {
		t.Fatalf("want refusal %s with detail %q, got %v", code, detail, err)
	}
}

// docsHashes returns the raw SHA-256 of every PNG in dir, keyed by name.
func docsHashes(t *testing.T, dir string) map[string]string {
	t.Helper()
	hashes := map[string]string{}
	for _, s := range Shots {
		b, err := os.ReadFile(filepath.Join(dir, s.Name))
		pngOK(t, err)
		sum := sha256.Sum256(b)
		hashes[s.Name] = hex.EncodeToString(sum[:])
	}
	return hashes
}

// docsSection returns the body of the QUICKSTART section whose heading line is
// heading: the lines after it, up to the next "## " heading or EOF.
func docsSection(t *testing.T, raw, heading string) string {
	t.Helper()
	start := strings.Index(raw, "\n"+heading+"\n")
	if start < 0 {
		t.Fatalf("QUICKSTART has no section %q", heading)
	}
	body := raw[start+len(heading)+2:]
	if end := strings.Index(body, "\n## "); end >= 0 {
		body = body[:end]
	}
	return body
}

// docsCommands returns the logical lines of the section's fenced bash blocks,
// with backslash continuations joined and comment lines dropped.
func docsCommands(section string) []string {
	var cmds []string
	var cur strings.Builder
	inBlock := false
	for _, line := range strings.Split(section, "\n") {
		t := strings.TrimRight(line, " \t\r")
		if !inBlock {
			if strings.TrimSpace(t) == "```bash" {
				inBlock = true
			}
			continue
		}
		if strings.TrimSpace(t) == "```" {
			inBlock = false
			continue
		}
		if strings.HasSuffix(t, `\`) {
			cur.WriteString(strings.TrimSuffix(t, `\`))
			cur.WriteByte(' ')
			continue
		}
		cur.WriteString(t)
		if s := strings.TrimSpace(cur.String()); s != "" && !strings.HasPrefix(s, "#") {
			cmds = append(cmds, s)
		}
		cur.Reset()
	}
	return cmds
}

// docsWithPrefix returns the commands that start with prefix.
func docsWithPrefix(cmds []string, prefix string) []string {
	var out []string
	for _, c := range cmds {
		if strings.HasPrefix(c, prefix) {
			out = append(out, c)
		}
	}
	return out
}

// docsFlags is the set of --flag tokens on one command line.
func docsFlags(cmd string) map[string]bool {
	set := map[string]bool{}
	for _, tok := range strings.Fields(cmd) {
		if docsFlagToken.MatchString(tok) {
			set[tok] = true
		}
	}
	return set
}

// docsExpand replaces the documented shell variables with the given paths and
// returns the argv. Any variable left over is a test failure.
func docsExpand(t *testing.T, cmd string, vars map[string]string) []string {
	t.Helper()
	cmd = strings.ReplaceAll(cmd, `"`, "")
	for name, value := range vars {
		cmd = strings.ReplaceAll(cmd, "$"+name, value)
	}
	if strings.Contains(cmd, "$") {
		t.Fatalf("unexpanded shell variable in documented line %q", cmd)
	}
	return strings.Fields(cmd)
}

// docsServePreview returns the one documented `"$DEMO_BIN" serve` line of a section.
func docsServePreview(t *testing.T, section string) string {
	t.Helper()
	lines := docsWithPrefix(docsCommands(section), `"$DEMO_BIN" serve `)
	if len(lines) != 1 {
		t.Fatalf("§13 documents %d `\"$DEMO_BIN\" serve` lines, want 1", len(lines))
	}
	return lines[0]
}

// docsShellServeFlags returns the flags of the serve argv that the capture shell starts.
func docsShellServeFlags(t *testing.T) map[string]bool {
	t.Helper()
	raw := docsRepoFile(t, "scripts/capture_marketing.sh")
	lines := strings.Split(raw, "\n")
	for i, line := range lines {
		if !strings.Contains(line, `"$worldd_bin" serve`) {
			continue
		}
		var argv strings.Builder
		for j := i; j < len(lines); j++ {
			argv.WriteString(lines[j])
			argv.WriteByte(' ')
			if !strings.HasSuffix(strings.TrimSpace(lines[j]), `\`) {
				break
			}
		}
		return docsFlags(argv.String())
	}
	t.Fatal("capture_marketing.sh has no daemon serve argv")
	return nil
}

// docsHelpFlags returns the --flag names in the configuration block of capture_marketing.sh --help.
func docsHelpFlags(t *testing.T) map[string]bool {
	t.Helper()
	out, err := exec.Command("bash", filepath.Join("..", "..", "scripts", "capture_marketing.sh"), "--help").CombinedOutput()
	pngOK(t, err)
	set := map[string]bool{}
	inBlock := false
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "Usage:") {
			inBlock = true
			continue
		}
		if inBlock && strings.TrimSpace(line) == "" {
			break
		}
		if inBlock && strings.HasPrefix(line, "  --") {
			set[strings.Fields(line)[0]] = true
		}
	}
	if len(set) == 0 {
		t.Fatalf("capture_marketing.sh --help has no configuration block:\n%s", out)
	}
	return set
}

// docsBuildDaemon builds ailang-worldd into dir; a failed build is never run.
func docsBuildDaemon(t *testing.T, dir string) string {
	t.Helper()
	bin := filepath.Join(dir, "ailang-worldd")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "go", "build", "-o", bin, "../../cmd/ailang-worldd").CombinedOutput()
	if err != nil {
		t.Fatalf("go build ailang-worldd failed (%v): %s", err, out)
	}
	return bin
}

// docsServeErr returns nil only when the documented serve argv selects still mode and
// parses: it must contain --workbench-still, carry the shell's serve flags, document
// each flag in serve --help, and reach the bind check after flag parsing.
func docsServeErr(t *testing.T, bin string, vars map[string]string, line string) error {
	t.Helper()
	if !docsFlags(line)["--workbench-still"] {
		return fmt.Errorf("serve preview lacks --workbench-still: %q", line)
	}
	preview := docsFlags(line)
	shell := docsShellServeFlags(t)
	for flag := range shell {
		if !preview[flag] {
			return fmt.Errorf("serve preview lacks shell flag %s", flag)
		}
	}
	for flag := range preview {
		if !shell[flag] {
			return fmt.Errorf("serve preview has flag %s the shell does not pass", flag)
		}
	}
	helpCtx, cancelHelp := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelHelp()
	help, err := exec.CommandContext(helpCtx, bin, "serve", "--help").CombinedOutput()
	if err != nil && helpCtx.Err() != nil {
		return fmt.Errorf("serve --help timed out")
	}
	for flag := range preview {
		if !strings.Contains(string(help), flag) {
			return fmt.Errorf("serve --help does not document %s", flag)
		}
	}
	argv := docsExpand(t, line, vars)
	if len(argv) < 2 || argv[1] != "serve" {
		return fmt.Errorf("serve preview is not a serve argv: %q", line)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, append(append([]string(nil), argv[1:]...), "--bind", "not-a-hostport")...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	runErr := cmd.Run()
	if ctx.Err() != nil {
		return fmt.Errorf("serve parse check timed out")
	}
	code := 0
	if runErr != nil {
		var ee *exec.ExitError
		if !errors.As(runErr, &ee) {
			return runErr
		}
		code = ee.ExitCode()
	}
	// ailang-worldd's usage exit is 1 (exitUsage); 2 is exitFatal. The bind check is a usage error.
	if code != 1 || !strings.Contains(stderr.String(), "is not host:port") {
		return fmt.Errorf("serve argv %q: exit %d, stderr %q", strings.Join(argv[1:], " "), code, stderr.String())
	}
	return nil
}

func TestMarketingDocsLinks(t *testing.T) {
	t.Run("real-repo", func(t *testing.T) {
		pngOK(t, CheckDocs("../..", nil))
		v1 := filepath.Join("..", "..", "website", "static", "img", "v1")
		manifest := filepath.Join(v1, "capture-manifest.json")
		if _, err := os.Stat(manifest); err == nil {
			pngOK(t, CheckPromotedAssets(v1, manifest, true))
			return
		}
		entries, err := os.ReadDir(v1)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		}
		for _, e := range entries {
			if filepath.Ext(e.Name()) == ".png" || e.Name() == "capture-manifest.json" {
				t.Fatalf("unattended tree holds %s before the attended capture", e.Name())
			}
		}
	})
	t.Run("promoted-dir", func(t *testing.T) {
		bs, _ := fixtureSet(t)
		dir := writeSet(t, bs)
		pngOK(t, CheckPromotedAssets(dir, "", false))
		m, err := BuildCaptureManifest(demoFixture(), docsHashes(t, dir))
		pngOK(t, err)
		pngOK(t, CheckPromotedAssets(dir, putJSON(t, t.TempDir(), "capture-manifest.json", m), true))
		placeholder := m
		placeholder.Shots = append([]CaptureRow(nil), m.Shots...)
		placeholder.Shots[0].SHA256 = strings.Repeat("0", 64)
		docsWantRefusal(t, CheckPromotedAssets(dir, putJSON(t, t.TempDir(), "capture-manifest.json", placeholder), true), "manifest-mismatch", "")
		dims := m
		dims.Shots = append([]CaptureRow(nil), m.Shots...)
		dims.Shots[2].Width++
		docsWantRefusal(t, CheckPromotedAssets(dir, putJSON(t, t.TempDir(), "capture-manifest.json", dims), true), "manifest-mismatch", "")
		missing := writeSet(t, bs)
		pngOK(t, os.Remove(filepath.Join(missing, "phone.png")))
		docsWantRefusal(t, CheckPromotedAssets(missing, "", false), "missing-shot", "")
	})
	t.Run("M8a-filename", func(t *testing.T) {
		page := docsRepoFile(t, docsPagePath)
		mutant := strings.Replace(page, "/img/v1/packet.png", "/img/v1/packet-v2.png", 1)
		if mutant == page {
			t.Fatal("the packet.png mutation did not apply")
		}
		docsWantRefusal(t, CheckDocs("../..", docsOverlay(map[string]string{docsPagePath: mutant})), "docs-image", "documented image target must match promoted inventory")
	})
	t.Run("refuse-readme-embed", func(t *testing.T) {
		readme := docsRepoFile(t, docsReadmePath) + "\n![decision packet](website/static/img/v1/packet.png)\n"
		docsWantRefusal(t, CheckDocs("../..", docsOverlay(map[string]string{docsReadmePath: readme})), "docs-readme", "")
	})
	t.Run("refuse-alt-route", func(t *testing.T) {
		captions := strings.Replace(docsRepoFile(t, docsCaptionsPath), "Alt: ", "Alt: /workbench?object=", 1)
		docsWantRefusal(t, CheckDocs("../..", docsOverlay(map[string]string{docsCaptionsPath: captions})), "docs-alt", "")
	})
	t.Run("refuse-captions-disclosure", func(t *testing.T) {
		captions := strings.ReplaceAll(docsRepoFile(t, docsCaptionsPath), "static capture of recorded history", "")
		docsWantRefusal(t, CheckDocs("../..", docsOverlay(map[string]string{docsCaptionsPath: captions})), "docs-captions", "")
	})
	t.Run("refuse-page-disclosure", func(t *testing.T) {
		for _, d := range []string{"static capture of recorded history", "demo-authored", "v0.54.15", "1440x3700"} {
			page := strings.ReplaceAll(docsRepoFile(t, docsPagePath), d, "")
			docsWantRefusal(t, CheckDocs("../..", docsOverlay(map[string]string{docsPagePath: page})), "docs-disclosure", "required disclosure is absent from the website page")
		}
	})
	t.Run("refuse-readme-live-claim", func(t *testing.T) {
		readme := strings.Replace(docsRepoFile(t, docsReadmePath), "are static captures of recorded history", "are live screenshots", 1)
		docsWantRefusal(t, CheckDocs("../..", docsOverlay(map[string]string{docsReadmePath: readme})), "docs-readme", "the Screenshots section must say the images are static captures of recorded history")
	})
	t.Run("refuse-schema-hash", func(t *testing.T) {
		schema := strings.Replace(docsRepoFile(t, docsSchemaPath), "{", "{\n  \"x-note\": \""+strings.Repeat("ab", 32)+"\",", 1)
		docsWantRefusal(t, CheckDocs("../..", docsOverlay(map[string]string{docsSchemaPath: schema})), "docs-schema", "")
	})
	t.Run("refuse-readme-broken-link", func(t *testing.T) {
		readme := strings.Replace(docsRepoFile(t, docsReadmePath), "### Screenshots\n", "### Screenshots\n\n[broken](docs/NOPE.md)\n", 1)
		docsWantRefusal(t, CheckDocs("../..", docsOverlay(map[string]string{docsReadmePath: readme})), "docs-link", "")
	})
	t.Run("refuse-page-markdown-image", func(t *testing.T) {
		page := docsRepoFile(t, docsPagePath) + "\n![x](/img/v1/graph.png)\n"
		docsWantRefusal(t, CheckDocs("../..", docsOverlay(map[string]string{docsPagePath: page})), "docs-image", "")
	})
	t.Run("cli", func(t *testing.T) {
		if rc := run([]string{"docs"}); rc != 2 {
			t.Fatalf("docs without --root: rc %d, want 2", rc)
		}
		if rc := run([]string{"docs", "--root", "../.."}); rc != 0 {
			t.Fatalf("docs --root ../..: rc %d, want 0", rc)
		}
		for _, args := range [][]string{
			{"docs", "--root", "../..", "--require-assets"},
			{"docs", "--root", "../..", "--capture-manifest", "x.json"},
			{"docs", "--root", "../..", "--bogus"},
			{"docs", "--root", "../..", "extra"},
		} {
			if rc := run(args); rc != 2 {
				t.Fatalf("docs %v: rc %d, want 2", args, rc)
			}
		}
		if rc := run([]string{"docs", "--help"}); rc != 0 {
			t.Fatalf("docs --help: rc %d, want 0", rc)
		}
	})
}

func TestMarketingRunbookMatchesHelp(t *testing.T) {
	raw := docsRepoFile(t, docsQuickstartPath)
	s12 := docsSection(t, raw, "## 12. Prepare the release demo seed (incremental tooling)")
	s13 := docsSection(t, raw, docsRunbookHeading)
	cmds13 := docsCommands(s13)
	if strings.Contains(s13, "/tmp/ailang-worldd") {
		t.Fatal("§13 must use $DEMO_BIN, not /tmp/ailang-worldd")
	}
	bin := docsBuildDaemon(t, t.TempDir())
	dir := t.TempDir()
	vars := map[string]string{"DEMO_DIR": dir, "DEMO_WS": filepath.Join(dir, "workspaces"), "DEMO_BIN": bin, "PIN": filepath.Join(dir, "pin"), "TOOLPIN": filepath.Join(dir, "toolpin")}

	t.Run("capture", func(t *testing.T) {
		lines := docsWithPrefix(cmds13, "scripts/capture_marketing.sh ")
		if len(lines) != 1 {
			t.Fatalf("§13 documents %d capture lines, want 1", len(lines))
		}
		want := docsHelpFlags(t)
		got := docsFlags(lines[0])
		for flag := range want {
			if !got[flag] {
				t.Errorf("documented capture line lacks %s (help configuration flag)", flag)
			}
		}
		for flag := range got {
			if !want[flag] {
				t.Errorf("documented capture line has %s, which capture_marketing.sh --help does not list", flag)
			}
		}
		if len(want) != 13 {
			t.Errorf("help lists %d configuration flags, want 13 (twelve value flags and --verify-twice)", len(want))
		}
	})
	t.Run("M8a-flag", func(t *testing.T) {
		line := docsServePreview(t, s13)
		if err := docsServeErr(t, bin, vars, line); err != nil {
			t.Fatalf("the unmutated preview must pass first: %v", err)
		}
		mutant := strings.Replace(line, "--workbench-still", "--workbench-static", 1)
		if mutant == line {
			t.Fatal("the serve preview lacks --workbench-still")
		}
		if err := docsServeErr(t, bin, vars, mutant); err == nil {
			t.Fatal("the --workbench-static mutant parsed as still mode")
		}
	})
	t.Run("serve-preview", func(t *testing.T) {
		if err := docsServeErr(t, bin, vars, docsServePreview(t, s13)); err != nil {
			t.Fatalf("documented serve argv must parse and select still mode: %v", err)
		}
	})
	t.Run("marketingcheck-lines", func(t *testing.T) {
		lines := docsWithPrefix(cmds13, "go run ./scripts/marketingcheck ")
		if len(lines) == 0 {
			t.Fatal("§13 documents no `go run ./scripts/marketingcheck` line")
		}
		for _, line := range lines {
			fields := strings.Fields(line)
			sub := fields[3]
			usageLine := ""
			for _, u := range strings.Split(usageText, "\n") {
				if strings.Contains(u, "marketingcheck "+sub+" ") {
					usageLine = u
				}
			}
			if usageLine == "" {
				t.Fatalf("usageText has no line for `marketingcheck %s`", sub)
			}
			for flag := range docsFlags(line) {
				if !strings.Contains(usageLine, flag) {
					t.Errorf("`marketingcheck %s` documents %s, absent from its usage line", sub, flag)
				}
			}
		}
	})
	t.Run("seed", func(t *testing.T) {
		lines := docsWithPrefix(docsCommands(s12), `"$DEMO_BIN" demo seed `)
		if len(lines) != 1 {
			t.Fatalf("§12 documents %d `demo seed` lines, want 1", len(lines))
		}
		helpCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		help, err := exec.CommandContext(helpCtx, bin, "demo", "seed", "--help").CombinedOutput()
		if err != nil && helpCtx.Err() != nil {
			t.Fatal("demo seed --help timed out")
		}
		for flag := range docsFlags(lines[0]) {
			if !strings.Contains(string(help), flag) {
				t.Errorf("documented demo seed flag %s is absent from demo seed --help", flag)
			}
		}
	})
	t.Run("attended-only", func(t *testing.T) {
		for _, cmd := range append(docsCommands(s12), cmds13...) {
			if strings.Contains(cmd, "yes |") || strings.Contains(cmd, "echo y |") || strings.Contains(cmd, "<<<") {
				t.Errorf("documented line pipes a confirmation: %q", cmd)
			}
			if (strings.Contains(cmd, "session mint") || strings.Contains(cmd, "world-publish")) && (strings.Contains(cmd, "|") || strings.Contains(cmd, "<<<")) {
				t.Errorf("documented line pipes or redirects session mint or world-publish: %q", cmd)
			}
		}
		if !strings.Contains(s13, "typed by a person") || !strings.Contains(s13, "§12") {
			t.Error("§13 must say the session mint and descriptor publish are typed by a person, and refer to §12 for them")
		}
	})
	t.Run("unminted-seed-refuses", func(t *testing.T) {
		seed := docsWithPrefix(docsCommands(s12), `"$DEMO_BIN" demo seed `)[0]
		argv := docsExpand(t, seed, vars)
		manifest := filepath.Join(dir, "manifest.json")
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
		devNull, err := os.Open(os.DevNull)
		pngOK(t, err)
		defer devNull.Close()
		cmd.Stdin = devNull
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		runErr := cmd.Run()
		if ctx.Err() != nil {
			t.Fatalf("unminted seed did not return within 60 s: %s", stderr.String())
		}
		if runErr == nil {
			t.Fatal("unminted seed exited 0")
		}
		if _, statErr := os.Stat(manifest); !errors.Is(statErr, os.ErrNotExist) {
			t.Fatalf("unminted seed created %s", manifest)
		}
		if docsHex64.MatchString(stderr.String()) {
			t.Fatalf("unminted seed printed a 64-hex string: %s", stderr.String())
		}
		t.Logf("unminted seed refusal (%v): %s", runErr, strings.TrimSpace(stderr.String()))
	})
}
