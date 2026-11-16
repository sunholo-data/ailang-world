package why

// Row 243 V2a (design §6 VS2a, §2.1): the walk's pure parts and its imports.

import (
	"context"
	"errors"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// VS2a: host/why's production files import no effect (DIRECT imports only;
// host/broker and host/coordinator reach host/store transitively).
func TestWhyImportsNoEffect(t *testing.T) {
	files, _ := filepath.Glob("*.go")
	seen := map[string]bool{}
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		p, err := parser.ParseFile(token.NewFileSet(), f, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatal(err)
		}
		for _, i := range p.Imports {
			path, _ := strconv.Unquote(i.Path.Value)
			seen[path] = true
		}
	}
	if !seen["github.com/sunholo-data/ailang-world/host/broker"] { // the instrument's control
		t.Fatalf("imports not read: %v", seen)
	}
	for _, bad := range []string{"os", "os/exec", "net", "net/http", "io/fs", "database/sql", "github.com/sunholo-data/ailang-world/host/store"} {
		if seen[bad] {
			t.Errorf("host/why imports %s", bad)
		}
	}
}

// ParseTarget: each form, and the CLI's error texts (why.go at the base, :316-336).
func TestParseTargetForms(t *testing.T) {
	ref := "sha256:" + strings.Repeat("a", 64)
	for arg, want := range map[string]Target{"head": {Kind: "head"}, "0": {Kind: "index"}, "42": {Kind: "index", Index: 42}, "a2a:e:t": {Kind: "invocation", Text: "a2a:e:t"}, "mcp:x": {Kind: "invocation", Text: "mcp:x"}, "rest:x": {Kind: "invocation", Text: "rest:x"}, ref: {Kind: "ref", Text: ref}} {
		if got, err := ParseTarget(arg); err != nil || got.Kind != want.Kind || got.Index != want.Index || got.Text != want.Text {
			t.Errorf("ParseTarget(%q) = %+v, %v", arg, got, err)
		}
	}
	nt := " is not an entry index, head, sha256:<ref>, a2a:/mcp:/rest:<id>, or -"
	for arg, want := range map[string]string{"": "a target is required", "-1": `"-1"` + nt, "x": `"x"` + nt,
		"sha256:zz": `"sha256:zz" is not a content address: hashref: digest is not lowercase hexadecimal: "zz"`} {
		if _, err := ParseTarget(arg); err == nil || err.Error() != want {
			t.Errorf("ParseTarget(%q) error = %v, want %s", arg, err, want)
		}
	}
}

// Head and GetJSON: 200 is the answer, 404 is absent, any other status or error is an error.
func TestHeadAndGetJSONClasses(t *testing.T) {
	answer := func(status int, body string, err error) Getter {
		return func(context.Context, string) (int, []byte, error) { return status, []byte(body), err }
	}
	ctx := context.Background()
	if ref, ok, err := Head(ctx, answer(200, " sha256:w\n", nil)); ref != "sha256:w" || !ok || err != nil {
		t.Errorf("200: %q %v %v", ref, ok, err)
	}
	if ref, ok, err := Head(ctx, answer(404, "x", nil)); ref != "" || ok || err != nil {
		t.Errorf("404: %q %v %v", ref, ok, err)
	}
	if _, ok, err := Head(ctx, answer(500, "down\n", nil)); ok || err == nil || err.Error() != "GET /v1/head returned HTTP 500: down" {
		t.Errorf("500: %v %v", ok, err)
	}
	if _, ok, err := Head(ctx, answer(0, "", errors.New("boom"))); ok || err == nil || err.Error() != "boom" {
		t.Errorf("err: %v %v", ok, err)
	}
	var v struct{ A int }
	if ok, err := (Reader{Ctx: ctx, Get: answer(200, "{", nil)}).GetJSON("/v1/x", &v); ok || err == nil || !strings.HasPrefix(err.Error(), "GET /v1/x: malformed body:") {
		t.Errorf("malformed: %v %v", ok, err)
	}
	done, cancel := context.WithCancel(ctx)
	cancel()
	if ok, err := (Reader{Ctx: done, Get: answer(200, `{"A":1}`, nil)}).GetJSON("/v1/x", &v); ok || !errors.Is(err, context.Canceled) {
		t.Errorf("a done context still read: %v %v", ok, err)
	}
}

// Elide's bounds, as the base's elide (mcpclient.go:651) answered them (run there, plan §1.2 C21).
func TestElideBounds(t *testing.T) {
	a := strings.Repeat("a", 160)
	for in, want := range map[string]string{a: a, a + "a": a + "… (161 bytes)",
		strings.Repeat("a\n", 81): strings.Repeat(`a\n`, 53) + "a… (162 bytes)",
		strings.Repeat("é", 100):  strings.Repeat("é", 80) + "… (200 bytes)"} {
		if got := Elide(in); got != want {
			t.Errorf("Elide(%q) = %q, want %q", in, got, want)
		}
	}
}
