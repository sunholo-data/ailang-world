package main

// Row 243 V2a (design §6 VW2–VW5): host/why's Run is the CLI's why over any Getter; bodies are recorded from the rig.

import (
	"context"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/verbs"
	"github.com/sunholo-data/ailang-world/host/why"
)

type recorded struct {
	status int
	body   []byte
}

// recordingGetter reads addr through the CLI's client and keeps every answer.
func recordingGetter(addr string, into map[string]recorded) why.Getter {
	get := clientGetter(newClient(addr))
	return func(ctx context.Context, path string) (int, []byte, error) {
		status, body, err := get(ctx, path)
		if err == nil {
			into[path] = recorded{status, body}
		}
		return status, body, err
	}
}

// replayGetter answers from a recording; edit may replace one answer.
func replayGetter(rec map[string]recorded, edit func(path string, r recorded) (recorded, error)) why.Getter {
	return func(_ context.Context, path string) (int, []byte, error) {
		r, ok := rec[path]
		if !ok {
			return 0, nil, errors.New("unrecorded " + path)
		}
		if edit != nil {
			var err error
			if r, err = edit(path, r); err != nil {
				return 0, nil, err
			}
		}
		return r.status, r.body, nil
	}
}

func whyAnswer(get why.Getter, arg string, o why.Options) why.Answer {
	tg, _ := why.ParseTarget(arg)
	return why.Run(why.Reader{Ctx: context.Background(), Get: get}, tg, o)
}

// VW2: Run's exit and streams are the CLI's, live and replayed.
func TestWhyRunReplaysTheCLI(t *testing.T) {
	r := newWalkRig(t)
	inv := r.invoke("ep1", "ailang-read", "task-1", "hello\n")
	for _, c := range []struct {
		addr, arg string
		o         why.Options
		exit      int
		want      string
	}{
		{r.url, "1", why.Options{Scan: 500}, exitOK, "all 7 link(s) verified\n"},
		{r.url, "head", why.Options{Scan: 500, JSON: true}, exitOK, `"matchedBy": "head"`},
		{tamperProxy(t, r.url, inv.input, []byte("forged")), "1", why.Options{Scan: 500}, exitIntegrity, "BROKEN — integrity refusal (exit 3)\n"},
		{r.url, "9", why.Options{Scan: 500}, exitUsage, "ailang-worldd why: log entry 9 not found\n"},
		{r.url, "a2a:nobody:x", why.Options{Scan: 1}, exitUsage, "ailang-worldd why: not found in the last 1 entries (invocationId)\n"},
	} {
		args := []string{c.arg, "--scan", strconv.Itoa(c.o.Scan)}
		if c.o.JSON {
			args = append(args, "--json")
		}
		code, out, errOut := whyRun(t, c.addr, "", args...)
		rec := map[string]recorded{}
		live := whyAnswer(recordingGetter(c.addr, rec), c.arg, c.o)
		replay := whyAnswer(replayGetter(rec, nil), c.arg, c.o)
		cli := why.Answer{Stdout: out, Stderr: errOut, Exit: code}
		if live != cli || replay != cli || code != c.exit || !strings.Contains(out+errOut, c.want) ||
			(out == "") != (code == exitUsage) || (errOut == "") != (code != exitUsage) {
			t.Fatalf("why %v: cli %+v\nlive %+v\nreplay %+v", args, cli, live, replay)
		}
	}
}

// VW3: a 404 is an answer (✗ "does not resolve"); an error or 5xx is the CLI's error line; a deadline is "not found".
func TestWhyRunGetterClasses(t *testing.T) {
	r := newWalkRig(t)
	inv := r.invoke("ep1", "ailang-read", "task-1", "hello\n")
	rec := map[string]recorded{}
	if a := whyAnswer(recordingGetter(r.url, rec), "1", why.Options{Scan: 500}); a.Exit != exitOK {
		t.Fatalf("recording walk: %+v", a)
	}
	input := "/v1/objects/" + inv.input + "?payload=true"
	for _, c := range []struct {
		path   string
		answer recorded
		err    error
		exit   int
		want   string
	}{
		{input, recorded{404, []byte(`{"error":"NotFound"}`)}, nil, exitIntegrity, "✗ input   " + inv.input + "   — ✗ does not resolve\n"},
		{"/v1/log/1", recorded{}, errors.New("boom"), exitUsage, "ailang-worldd why: boom\n"},
		{"/v1/log/1", recorded{500, []byte("down\n")}, nil, exitUsage, "ailang-worldd why: GET /v1/log/1 returned HTTP 500: down\n"},
		{"/v1/log/1", recorded{404, nil}, nil, exitUsage, "ailang-worldd why: log entry 1 not found\n"},
		{"/v1/health", recorded{}, errors.New("health down"), exitUsage, "ailang-worldd why: health down\n"},
		{"/v1/log/1", recorded{}, context.DeadlineExceeded, exitUsage, "ailang-worldd why: not found in the last 500 entries (entry index)\n"},
	} {
		a := whyAnswer(replayGetter(rec, func(path string, r recorded) (recorded, error) {
			if path != c.path {
				return r, nil
			}
			return c.answer, c.err
		}), "1", why.Options{Scan: 500})
		if a.Exit != c.exit || !strings.Contains(a.Stdout+a.Stderr, c.want) {
			t.Fatalf("%s %+v %v: %+v", c.path, c.answer, c.err, a)
		}
	}
}

// VW4: the moved constants, as literals, so a change on both sides still fails.
func TestWhyConstantsMatchTheCLI(t *testing.T) {
	if why.MaxResponseBytes != maxClientResponseBytes || maxClientResponseBytes != 8388608 ||
		whyScanDefault != 500 || whyScanMax != 5000 || logPageMax != 500 || why.ElideAt != 160 ||
		verbs.WhyScanDefault != why.ScanDefault || verbs.WhyScanMax != why.ScanMax {
		t.Fatalf("constants drifted: %d %d %d %d %d", why.MaxResponseBytes, whyScanDefault, whyScanMax, logPageMax, why.ElideAt)
	}
}

// VW5: one /v1/head reader, and one elide: cmd's are single delegations to host/why.
func TestWhyHeadIsOneReader(t *testing.T) {
	f, err := parser.ParseFile(token.NewFileSet(), "mcpclient.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"headState": "Head", "elide": "Elide"}
	for _, d := range f.Decls {
		fn, ok := d.(*ast.FuncDecl)
		if !ok || want[fn.Name.Name] == "" {
			continue
		}
		var sel *ast.SelectorExpr
		if len(fn.Body.List) == 1 {
			if ret, ok := fn.Body.List[0].(*ast.ReturnStmt); ok && len(ret.Results) == 1 {
				if call, ok := ret.Results[0].(*ast.CallExpr); ok {
					sel, _ = call.Fun.(*ast.SelectorExpr)
				}
			}
		}
		if sel == nil || sel.Sel.Name != want[fn.Name.Name] || ast.Unparen(sel.X).(*ast.Ident).Name != "why" {
			t.Fatalf("%s is not one return of why.%s", fn.Name.Name, want[fn.Name.Name])
		}
		delete(want, fn.Name.Name)
	}
	if len(want) != 0 {
		t.Fatalf("not found in mcpclient.go: %v", want)
	}
}
