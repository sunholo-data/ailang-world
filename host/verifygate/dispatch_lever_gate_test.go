package verifygate

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// onBlockTriggerKeys returns the trigger keys declared DIRECTLY under the on: block — the
// same depth as push:/pull_request: — so a workflow_dispatch folded into nested inputs: or
// buried in a comment does not count.
func onBlockTriggerKeys(t *testing.T, path, src string) []string {
	t.Helper()
	lines := strings.Split(src, "\n")
	start := -1
	for i, l := range lines {
		// EXACT byte equality at column 0 -- deliberately NOT strings.TrimSpace. A trimmed
		// comparison anchors on ANY nested `on:` (under a job step's `with:`/`inputs:`),
		// so the gate could scan a non-trigger block and pass vacuously. This matches
		// P14's verification method (see P14 -- its command was itself repaired to a
		// column-0 test, because `awk '{$1=$1}' ` strips leading whitespace and matched
		// an indented `  on:` too; measured, arms differ).
		if l == "on:" {
			start = i
			break
		}
	}
	if start == -1 {
		t.Fatalf("instrument failure: %s has no top-level `on:` trigger block", path)
	}
	// gpt5-6-sol R2: fail if the top-level `on:` is DUPLICATED, not merely absent.
	dupes := 0
	for _, l := range lines {
		if l == "on:" {
			dupes++
		}
	}
	if dupes != 1 {
		t.Fatalf("instrument failure: %s has %d top-level `on:` blocks, want exactly 1", path, dupes)
	}
	onLead := len(lines[start]) - len(strings.TrimLeft(lines[start], " "))
	triggerLead, keys := -1, []string{}
	for _, l := range lines[start+1:] {
		tok := strings.TrimSpace(l)
		if tok == "" || strings.HasPrefix(tok, "#") {
			continue
		}
		lead := len(l) - len(strings.TrimLeft(l, " "))
		if lead <= onLead {
			break // left the on: block
		}
		if triggerLead == -1 {
			triggerLead = lead
		}
		if lead == triggerLead {
			if kv := strings.SplitN(tok, ":", 2); len(kv) == 2 {
				// gpt5-6-sol R2: validate the VALUE. `workflow_dispatch` is legal as an
				// empty/null mapping key or as a mapping-valued trigger config -- never as a
				// bare scalar. `workflow_dispatch: garbage` must NOT count as declared.
				if kv[0] == "workflow_dispatch" && strings.TrimSpace(kv[1]) != "" {
					t.Errorf("%s: `workflow_dispatch:` has scalar value %q; want an empty key or a mapping",
						path, strings.TrimSpace(kv[1]))
					continue
				}
				keys = append(keys, kv[0])
			}
		}
	}
	return keys
}

// TestEveryWorkflowDeclaresDispatchLever pins queue-item 47's lever so it cannot be deleted
// silently. A push to dev whose webhook GitHub drops is PERMANENTLY unverifiable (never
// replayed, no run created), and the only API-driven lever is workflow_dispatch, which no
// workflow currently declares (P2). This gate asserts EVERY enumerated workflow file declares
// the lever as a trigger in its on: block.
//
// DECLARED RESIDUAL: this is a STATIC text scan over YAML. It proves the lever is DECLARED,
// never that a dispatch RUN is created or is green. And a workflow_dispatch run is NOT
// equivalent to the event it replaces: its checks do not satisfy PR branch protection
// (measured by this mission — required contexts can read success on the head SHA while
// gh pr checks --required still lists only the pull_request-event context with
// mergeStateStatus=BLOCKED). The lever buys A VERDICT ON A COMMIT, not A MERGEABLE PR. It
// also cannot see a workflow added OUTSIDE .github/workflows/ (a singular `workflow` dir, a
// root-level .yaml, a hidden file), a case-mismatched filename (the Glob is case-sensitive),
// or a nested subdirectory (which GitHub itself does not scan either).
func TestEveryWorkflowDeclaresDispatchLever(t *testing.T) {
	matches, err := filepath.Glob(filepath.Join(repoRoot, ".github", "workflows", "*"))
	if err != nil {
		t.Fatal(err)
	}
	// ANTI-VACUITY FLOOR: an empty enumeration FAILS LOUDLY rather than printing a checkmark.
	if len(matches) == 0 {
		t.Fatal("instrument failure: no workflow files enumerated under .github/workflows/ — an empty set proves nothing about the lever")
	}
	for _, m := range matches {
		raw, err := os.ReadFile(m)
		if err != nil {
			t.Fatal(err)
		}
		src := string(raw)
		// known-positive control: a scan that cannot see the on: block is reading the wrong file.
		if !strings.Contains(src, "on:") {
			t.Fatalf("instrument failure: %s does not contain known-positive control %q", m, "on:")
		}
		triggers := onBlockTriggerKeys(t, m, src)
		if !slices.Contains(triggers, "workflow_dispatch") {
			t.Errorf("%s: on-block triggers=%v lack workflow_dispatch — a dropped push to dev is permanently unverifiable; every workflow file must declare the lever", filepath.Base(m), triggers)
		}
	}
}
