package coordinator

// Row 134 M5a: the packages/se-tools transitions, run through the REAL capsule
// runner under the pinned interpreter, must emit plans the Go plan law admits
// against the checked-in manifest's declaredEffects — and refuse a bad
// argument with a typed zero-effect L-ARGS refusal. The manifest is the one
// `world-publish transitions` publishes (cmd/world-publish's se-tools test
// publishes it; this test reads its effect declarations), so a module that
// drifts from its descriptor is red here before it can be published.

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/canon"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

const seToolsManifest = "packages/se-tools/transitions.json"

type seToolsRequirement struct {
	Effect string `json:"effect"`
	Scope  string `json:"scope"`
	Cost   int64  `json:"cost"`
}

// seToolsEntry is the world-publish manifest entry shape (cmd/world-publish
// manifestEntry), decoded STRICTLY: the verb's own decoder ignores unknown
// keys, so a misspelled field would otherwise publish silently without it.
type seToolsEntry struct {
	ID               string               `json:"id"`
	Title            string               `json:"title"`
	Description      string               `json:"description"`
	TransitionFnFile string               `json:"transitionFnFile"`
	InputSchema      json.RawMessage      `json:"inputSchema"`
	OutputSchema     json.RawMessage      `json:"outputSchema"`
	Access           seToolsRequirement   `json:"access"`
	DeclaredEffects  []seToolsRequirement `json:"declaredEffects"`
}

func seToolsRepoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func loadSeToolsManifest(t *testing.T) map[string]seToolsEntry {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(seToolsRepoRoot(t), seToolsManifest))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	var entries []seToolsEntry
	if err := dec.Decode(&entries); err != nil {
		t.Fatalf("strict manifest decode: %v", err)
	}
	byID := map[string]seToolsEntry{}
	for _, e := range entries {
		if _, dup := byID[e.ID]; dup {
			t.Fatalf("manifest repeats id %q", e.ID)
		}
		byID[e.ID] = e
	}
	return byID
}

func (e seToolsEntry) declared() []transitionreg.EffectRequirement {
	out := make([]transitionreg.EffectRequirement, len(e.DeclaredEffects))
	for i, d := range e.DeclaredEffects {
		out[i] = transitionreg.EffectRequirement{Effect: d.Effect, Scope: d.Scope, Cost: d.Cost}
	}
	return out
}

// seToolsRig runs a module the way Dispatch does: canonical source (the bytes
// publication stores), the archived pinned interpreter, input passed as ONE
// JSON string through --args-file.
type seToolsRig struct {
	t      *testing.T
	runner *capsule.Runner
	interp hashref.HashRef
	bin    string
}

func newSeToolsRig(t *testing.T) *seToolsRig {
	t.Helper()
	bin := os.Getenv("AILANG_BIN")
	if bin == "" {
		t.Fatal("AILANG_BIN unset: the se-tools plans must run under the pinned released interpreter; never skip")
	}
	out, err := exec.Command(bin, "--version").Output()
	if err != nil || !strings.HasPrefix(string(out), "AILANG v0.41.0") {
		t.Fatalf("AILANG_BIN %q is not the pinned v0.41.0 interpreter: err=%v output=%q", bin, err, out)
	}
	a := archive.New(filepath.Join(t.TempDir(), "world.db"))
	ref, err := a.Archive(bin)
	if err != nil {
		t.Fatalf("archive interpreter: %v", err)
	}
	return &seToolsRig{t: t, runner: capsule.New(a, capsule.Config{}), interp: ref, bin: bin}
}

func (r *seToolsRig) source(e seToolsEntry) []byte {
	r.t.Helper()
	raw, err := os.ReadFile(filepath.Join(seToolsRepoRoot(r.t), filepath.FromSlash(e.TransitionFnFile)))
	if err != nil {
		r.t.Fatalf("%s: read source: %v", e.ID, err)
	}
	src, err := canon.Source(raw)
	if err != nil {
		r.t.Fatalf("%s: canonicalise source: %v", e.ID, err)
	}
	return src
}

func (r *seToolsRig) run(e seToolsEntry, input string) []byte {
	r.t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := r.runner.RunContext(ctx, capsule.Entry{Interpreter: r.interp, Source: r.source(e), Args: mustJSON(input)})
	if err != nil {
		r.t.Fatalf("%s: capsule run of %s: %v", e.ID, input, err)
	}
	return bytes.TrimSuffix(res.Stdout, []byte("\n"))
}

func planInput(args string) string { return `{"phase":"plan","args":` + args + `}` }

func canonical(t *testing.T, raw string) string {
	t.Helper()
	b, err := transitionreg.CanonicalJSON([]byte(raw), 1<<20)
	if err != nil {
		t.Fatalf("canonical %s: %v", raw, err)
	}
	return string(b)
}

type seToolsCase struct {
	id      string
	args    string // an admitted call
	payload string // the effect payload it must plan
	finish  bool
	badArgs string // the same call plus one key outside the tool's schema
	badKey  string
}

var seToolsCases = []seToolsCase{
	{"ailang-read", `{"path":"src/main.ail"}`, `{"op":"read","path":"src/main.ail"}`, false,
		`{"path":"src/main.ail","limit":5}`, "limit"},
	{"ailang-write", `{"path":"src/new.ail","content":"module src/new\n"}`,
		`{"op":"write","path":"src/new.ail","content":"module src/new\n"}`, false,
		`{"path":"src/new.ail","content":"x","old_text":"y"}`, "old_text"},
	{"ailang-edit", `{"path":"data.txt","old_text":"hello","new_text":"howdy"}`,
		`{"op":"edit","path":"data.txt","old_text":"hello","new_text":"howdy"}`, false,
		`{"path":"data.txt","old":"hello","new":"howdy"}`, "old"},
	{"ailang-check", `{"path":"bad.ail"}`, `{"op":"ai_check","path":"bad.ail"}`, true,
		`{"path":"bad.ail","flags":{"json":""}}`, "flags"},
	{"ailang-run", `{"path":"args.ail","args_json":"\"data.txt\""}`, `{"path":"args.ail","args_json":"\"data.txt\""}`, false,
		`{"path":"args.ail","caps":"IO,Net"}`, "caps"},
	{"builtins-search", `{"query":"readFile","module":"std/fs"}`, `{"op":"builtins_list","flags":{"json":""}}`, true,
		`{"query":"readFile","limit":3}`, "limit"},
	{"examples-search", `{"query":"foldl"}`, `{"op":"examples_search","query":"foldl"}`, false,
		`{"query":"foldl","limit":2}`, "limit"},
	{"ailang-cli", `{"op":"iface","module":"std/json","flags":{"json":""}}`, `{"op":"iface","module":"std/json","flags":{"json":""}}`, false,
		`{"op":"check","path":"a.ail","content":"x"}`, "content"},
}

// TestSeToolsPlansAdmittedByPlanLaw is AC5 (M5a): per tool, the module run
// through the real capsule emits a plan parsePlan ADMITS against the
// manifest's declaredEffects — exactly one effect, the declared triple, the
// §4.2 payload, the finish flag — and a bad argument key yields a typed
// zero-effect L-ARGS refusal naming the key (which parsePlan also admits:
// a refusal is a lawful pure plan, L7).
func TestSeToolsPlansAdmittedByPlanLaw(t *testing.T) {
	manifest := loadSeToolsManifest(t)
	if len(manifest) != 8 || len(seToolsCases) != 8 {
		t.Fatalf("manifest has %d tools, cases %d; the ailang_only lane has 8", len(manifest), len(seToolsCases))
	}
	rig := newSeToolsRig(t)
	for _, c := range seToolsCases {
		t.Run(c.id, func(t *testing.T) {
			e, ok := manifest[c.id]
			if !ok {
				t.Fatalf("manifest has no %q", c.id)
			}
			declared := e.declared()
			if len(declared) != 1 || declared[0].Scope != "worktree" || declared[0].Cost != 1 ||
				e.Access.Effect != declared[0].Effect || e.Access.Scope != "worktree" || e.Access.Cost != 0 {
				t.Fatalf("descriptor triples: access %+v declared %+v; want the effect at cost 0 and cost 1 in scope worktree", e.Access, declared)
			}

			plan, err := parsePlan(rig.run(e, planInput(c.args)), declared)
			if err != nil {
				t.Fatalf("parsePlan refused the admitted call's plan: %v", err)
			}
			if len(plan.Effects) != 1 || plan.Result != nil {
				t.Fatalf("plan effects=%d result=%s; want one effect and a null result", len(plan.Effects), plan.Result)
			}
			got := plan.Effects[0]
			if got.Requirement != declared[0] || got.ID != "e1" {
				t.Fatalf("planned (%s, %+v), want (e1, %+v)", got.ID, got.Requirement, declared[0])
			}
			if string(got.Payload) != canonical(t, c.payload) {
				t.Fatalf("payload %s, want %s", got.Payload, canonical(t, c.payload))
			}
			if plan.Finish != c.finish {
				t.Fatalf("finish=%v, want %v", plan.Finish, c.finish)
			}

			refusal, err := parsePlan(rig.run(e, planInput(c.badArgs)), declared)
			if err != nil {
				t.Fatalf("parsePlan refused the L-ARGS refusal plan: %v", err)
			}
			if len(refusal.Effects) != 0 || refusal.Finish {
				t.Fatalf("bad argument %q planned %d effects (finish=%v); L-ARGS must refuse with zero effects",
					c.badKey, len(refusal.Effects), refusal.Finish)
			}
			var result struct {
				OK      *bool  `json:"ok"`
				Refused string `json:"refused"`
			}
			if err := json.Unmarshal(refusal.Result, &result); err != nil {
				t.Fatalf("refusal result %s: %v", refusal.Result, err)
			}
			want := `unknown argument "` + c.badKey + `"; ` + c.id + " admits only: "
			if result.OK == nil || *result.OK || !strings.HasPrefix(result.Refused, want) {
				t.Fatalf("refusal result %s, want ok:false and refused starting %q", refusal.Result, want)
			}
		})
	}
}

// TestSeToolsPathPredicateRefusesBeforeAnyEffect: every path-taking tool
// refuses a hyphen-leading path (V58's flag-injection vector) and a `..`
// escape in its pure plan — zero effects, so nothing reaches a handler.
func TestSeToolsPathPredicateRefusesBeforeAnyEffect(t *testing.T) {
	manifest := loadSeToolsManifest(t)
	rig := newSeToolsRig(t)
	calls := map[string]string{
		"ailang-read":  `{"path":%q}`,
		"ailang-write": `{"path":%q,"content":"x"}`,
		"ailang-edit":  `{"path":%q,"old_text":"a","new_text":"b"}`,
		"ailang-check": `{"path":%q}`,
		"ailang-run":   `{"path":%q}`,
		"ailang-cli":   `{"op":"check","path":%q}`,
	}
	for id, shape := range calls {
		for _, p := range []string{"-policy.ail", "../ailang/go.mod", "/etc/hosts"} {
			e := manifest[id]
			args := strings.Replace(shape, "%q", `"`+p+`"`, 1)
			plan, err := parsePlan(rig.run(e, planInput(args)), e.declared())
			if err != nil {
				t.Fatalf("%s %s: parsePlan: %v", id, p, err)
			}
			if len(plan.Effects) != 0 || !strings.Contains(string(plan.Result), `refused: it must be relative`) {
				t.Fatalf("%s path %q: effects=%d result=%s; want the path-predicate refusal", id, p, len(plan.Effects), plan.Result)
			}
		}
	}
}

func effectResultsInput(args, effect, output string) string {
	return `{"phase":"finish","args":` + args + `,"results":[{"id":"e1","effect":"` + effect +
		`","status":"ok","record":"sha256:` + strings.Repeat("0", 64) + `","output":` + output + `}]}`
}

// TestSeToolsFinishPhases runs the two finish phases on REAL tool output from
// the pinned binary (ai-check of a broken module; the builtin inventory) and
// requires parseFinish to admit what they return (L6).
func TestSeToolsFinishPhases(t *testing.T) {
	manifest := loadSeToolsManifest(t)
	rig := newSeToolsRig(t)

	t.Run("ailang-check", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "bad.ail"), []byte("module bad\n\nexport func main() -> int = \"oops\"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command(rig.bin, "ai-check", "bad.ail")
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "AILANG_CACHE_DIR="+filepath.Join(dir, "cache"))
		stdout, _ := cmd.Output() // rc 1 on a type error; the JSON report is on stdout
		handlerOut, _ := json.Marshal(map[string]any{"ok": false, "argv": []string{"ai-check", "bad.ail"},
			"exit_code": 1, "stdout": string(stdout), "tool": "sha256:x", "policy_digest": "d"})
		out, err := parseFinish(rig.run(manifest["ailang-check"], effectResultsInput(`{"path":"bad.ail"}`, "Ailang.Check", string(handlerOut))))
		if err != nil {
			t.Fatalf("parseFinish: %v", err)
		}
		var rep struct {
			OK         bool `json:"ok"`
			Passed     *bool
			ErrorCount int `json:"error_count"`
			Errors     []struct{ Code, Message, File string }
		}
		if err := json.Unmarshal(out, &rep); err != nil {
			t.Fatal(err)
		}
		if !rep.OK || rep.Passed == nil || *rep.Passed || rep.ErrorCount != 1 || len(rep.Errors) != 1 ||
			rep.Errors[0].File != "bad.ail" || !strings.Contains(rep.Errors[0].Message, "bad.ail:3") {
			t.Fatalf("check report %s; want ok, passed:false, one error in bad.ail with its line in message", out)
		}
	})

	t.Run("builtins-search", func(t *testing.T) {
		inventory, err := exec.Command(rig.bin, "builtins", "list", "--json").Output()
		if err != nil {
			t.Fatalf("builtins list: %v", err)
		}
		handlerOut, _ := json.Marshal(map[string]any{"ok": true, "argv": []string{"builtins", "list", "--json"},
			"stdout": string(inventory), "tool": "sha256:x", "policy_digest": "d"})
		start := time.Now()
		out, err := parseFinish(rig.run(manifest["builtins-search"],
			effectResultsInput(`{"query":"READFILE","module":"std/fs"}`, "Ailang.Discover", string(handlerOut))))
		if err != nil {
			t.Fatalf("parseFinish: %v", err)
		}
		t.Logf("builtins finish over a %d-byte inventory: %s", len(inventory), time.Since(start))
		var rep struct {
			Count   int
			Matches []struct{ Name, Module string }
		}
		if err := json.Unmarshal(out, &rep); err != nil {
			t.Fatal(err)
		}
		found := false
		for _, m := range rep.Matches {
			if m.Module != "std/fs" && !strings.Contains(m.Module, "std/fs") {
				t.Fatalf("match %+v escapes the module filter", m)
			}
			found = found || m.Name == "readFile" || strings.Contains(strings.ToLower(m.Name), "readfile")
		}
		if rep.Count != len(rep.Matches) || rep.Count == 0 || rep.Count > 10 || !found {
			t.Fatalf("builtins report %s; want 1..10 std/fs matches including readFile", out)
		}
	})
}

// TestSeToolsSharedCoreIsByteIdentical: publication forbids a shared local
// module (self-contained sources, row 107 Residual 7), so the contracted path
// predicate, the convention-v2 key law and the plan builders are COPIED into
// every module. This holds the copies identical, so a fix to one cannot leave
// seven stale.
func TestSeToolsSharedCoreIsByteIdentical(t *testing.T) {
	manifest := loadSeToolsManifest(t)
	blocks := []struct{ name, start, end string }{
		{"path core", "-- ── Contracted pure core", "\n-- L-ARGS (§4.1, §4.2)"},
		{"inputKeyAllowed", "export func inputKeyAllowed", "\n}\n"},
		{"plan construction", "-- ── Plan construction", "\nfunc firstUnknownArg"},
	}
	var first map[string]string
	firstID := ""
	for _, c := range seToolsCases {
		raw, err := os.ReadFile(filepath.Join(seToolsRepoRoot(t), filepath.FromSlash(manifest[c.id].TransitionFnFile)))
		if err != nil {
			t.Fatal(err)
		}
		src := string(raw)
		got := map[string]string{}
		for _, b := range blocks {
			i := strings.Index(src, b.start)
			if i < 0 || strings.Count(src, b.start) != 1 {
				t.Fatalf("%s: block %q start marker found %d times", c.id, b.name, strings.Count(src, b.start))
			}
			j := strings.Index(src[i:], b.end)
			if j < 0 {
				t.Fatalf("%s: block %q has no end marker", c.id, b.name)
			}
			got[b.name] = src[i : i+j]
		}
		if first == nil {
			first, firstID = got, c.id
			continue
		}
		for _, b := range blocks {
			if got[b.name] != first[b.name] {
				t.Fatalf("%s's %s differs from %s's", c.id, b.name, firstID)
			}
		}
	}
}
