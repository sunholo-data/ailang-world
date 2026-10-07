package coordinator

// Row 153 M3a: the se-tools plans and finishes run warm from a publication-built
// compile-cache template and produce the bytes a cold run does.

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/capsule"
)

// runFull is run with stderr: the interpreter warns MOD010 on stderr only when
// it COMPILES a module whose header differs from the staging path, which every
// se-tools source does (module se_tools/*). That warning is the deterministic
// warmth witness (P-M4); output equality alone cannot see a stale template
// (P-M3: the interpreter silently recompiles).
func (r *seToolsRig) runFull(e seToolsEntry, input string) capsule.Result {
	r.t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := r.runner.RunContext(ctx, capsule.Entry{Interpreter: r.interp, Source: r.source(e), Args: mustJSON(input)})
	if err != nil {
		r.t.Fatalf("%s: capsule run of %s: %v (stderr %q)", e.ID, input, err, res.Stderr)
	}
	return res
}

func (r *seToolsRig) dropTemplate(e seToolsEntry) {
	r.t.Helper()
	if err := os.RemoveAll(r.arch.CapsuleTemplateDir(r.interp, r.source(e))); err != nil {
		r.t.Fatal(err)
	}
}

func (r *seToolsRig) buildTemplate(e seToolsEntry) {
	r.t.Helper()
	res, err := r.arch.CheckSource(context.Background(), r.interp, r.source(e))
	if err != nil || !res.Passed {
		r.t.Fatalf("%s: CheckSource = (%+v, %v), want a pass that built the template", e.ID, res, err)
	}
}

func compiledStderr(res capsule.Result) bool { return bytes.Contains(res.Stderr, []byte("MOD010")) }

// AC3.3 (MUT-STALE-KEY): per se-tools phase, a run from the publication-built
// template emits the same bytes as a cold run AND does not recompile; a cold
// run does. The cross-source arm is the tooth for a template key that ignores
// the source: a template built for another tool must not warm this one.
func TestCapsuleWarmEqualsColdForSeTools(t *testing.T) {
	manifest := loadSeToolsManifest(t)
	rig := newSeToolsRig(t)

	type phase struct{ name, id, input string }
	var phases []phase
	for _, c := range seToolsCases {
		phases = append(phases, phase{c.id + "/plan", c.id, planInput(c.args)})
	}
	// A finish phase: the canned handler output is the check tool's report shape.
	phases = append(phases, phase{"ailang-check/finish", "ailang-check",
		effectResultsInput(`{"path":"bad.ail"}`, "Ailang.Check", `{"ok":true,"argv":["ai-check","bad.ail"],"exit_code":0,"stdout":"{}","tool":"sha256:x","policy_digest":"d"}`)})

	for _, p := range phases {
		t.Run(p.name, func(t *testing.T) {
			e := manifest[p.id]
			rig.dropTemplate(e)
			coldBefore := rig.runner.ColdRuns()
			cold := rig.runFull(e, p.input)
			if !compiledStderr(cold) || rig.runner.ColdRuns() != coldBefore+1 {
				t.Fatalf("cold arm: MOD010=%v ColdRuns %d -> %d; want a compiling run counted cold", compiledStderr(cold), coldBefore, rig.runner.ColdRuns())
			}
			// The cold run lazily promoted its own cache. Drop it so the warm
			// arm runs from the PUBLICATION-built template, not the run-promoted
			// one (judge N2: otherwise a hollow check-built template is invisible).
			rig.dropTemplate(e)
			rig.buildTemplate(e)
			coldBefore = rig.runner.ColdRuns()
			warm := rig.runFull(e, p.input)
			if compiledStderr(warm) || rig.runner.ColdRuns() != coldBefore {
				t.Fatalf("warm arm: MOD010=%v ColdRuns %d -> %d (stderr %q); want a template-served run", compiledStderr(warm), coldBefore, rig.runner.ColdRuns(), warm.Stderr)
			}
			if !bytes.Equal(warm.Stdout, cold.Stdout) {
				t.Fatalf("warm stdout differs from cold:\nwarm %s\ncold %s", warm.Stdout, cold.Stdout)
			}
		})
	}

	t.Run("cross-source", func(t *testing.T) {
		read, exec := manifest["ailang-read"], manifest["workspace-exec"]
		rig.dropTemplate(read)
		rig.dropTemplate(exec)
		rig.buildTemplate(read)
		rig.buildTemplate(exec)
		res := rig.runFull(exec, planInput(`{"command":"test","args":["-run=X"]}`))
		if compiledStderr(res) {
			t.Fatalf("workspace-exec compiled although its own template was built (stderr %q): the template key ignores the source", res.Stderr)
		}
	})
}
