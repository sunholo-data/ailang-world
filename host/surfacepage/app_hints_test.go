package surfacepage

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestAddressHintsArePinned (row 243 V3 sprint 2 M0, design §9 (b), §6 VP12): the address bar's three
// hints. node evaluates app.js's own placeholder expression (each branch) and its KEYROWS literal, cut
// from the script's source; Go reads index.html's placeholder. It skips only when no node exists; with
// WORLD_EXEC_NODE set, a harness that cannot run fails.
func TestAddressHintsArePinned(t *testing.T) {
	const placeholder = `placeholder="row-189@6 · shared@0ffbb9c · back · episode report &lt;episode&gt; · why &lt;target&gt;">`
	if n := strings.Count(string(appIndex), placeholder); n != 1 || !strings.Contains(string(AppShell), placeholder) {
		t.Fatalf("index.html's address placeholder: %d of %q", n, placeholder)
	}
	node := os.Getenv("WORLD_EXEC_NODE")
	if node == "" {
		var err error
		node, err = exec.LookPath("node")
		if err != nil {
			t.Skip("node unavailable locally")
		}
	}
	dir := t.TempDir()
	for name, data := range map[string][]byte{"app.js": AppScript, "harness.js": []byte(hintsHarness)} {
		if err := os.WriteFile(filepath.Join(dir, name), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, node, filepath.Join(dir, "harness.js"))
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("node harness: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "ADDRESS HINTS PASS") {
		t.Fatalf("harness manifest missing: %s", out)
	}
}

const hintsHarness = `
const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const src = fs.readFileSync('app.js', 'utf8');
const one = (re) => { const m = src.match(re); assert(m && src.split(m[0]).length === 2, String(re)); return m[1]; };
// The placeholder: app.js's own expression, evaluated in each of its three branches.
const expr = one(/\n *\$\('addr'\)\.placeholder = (.*);\n/);
const at = (st, T) => vm.runInNewContext(expr, { st, T });
const tail = ' · episode report <episode> · why <target>';
assert.equal(at({ d: 1 }, null), '');
assert.equal(at({ d: 0 }, null), 'go to a place — row-189@6 · shared@0ffbb9c' + tail);
assert.equal(at({ d: 0 }, { GROUPS: [{ id: 'row' }, { id: 'shared' }] }), 'go to a place — row · shared · row:row-N@1' + tail);
// The keys popover: app.js's own KEYROWS literal; the bar's row names every form the bar takes.
const rows = JSON.parse(JSON.stringify(vm.runInNewContext(one(/\n *const KEYROWS = (\[.*\]);\n/)))); // out of the vm realm
assert.equal(rows.length, 13);
assert.deepEqual(rows[10], ['Enter', 'in the bar: a place, episode report E, why T, approve ID', 'the same']);
console.log('ADDRESS HINTS PASS');
`
