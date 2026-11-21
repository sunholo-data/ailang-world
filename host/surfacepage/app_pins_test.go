package surfacepage

// Row 243 V4 S1 (design §5.6): the three judge-405 pins on app.js's Needs you and the bar's navigate
// allowlist, each a node harness over the daemon's own answers (requests-confirmations.golden.json,
// written by host/daemon TestRequestsJSONConfirmationsGolden; verb-route.golden.json, checked by
// TestSurfaceVerbOutcomes). Each harness pairs its hostile rows with the real row they are cut from.

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// runPinHarness runs one harness under node beside app.js and the two goldens, and requires its
// manifest line. It skips only when no node exists; with WORLD_EXEC_NODE set, a harness that cannot
// run fails.
func runPinHarness(t *testing.T, harness, manifest string) {
	t.Helper()
	node := os.Getenv("WORLD_EXEC_NODE")
	if node == "" {
		var err error
		node, err = exec.LookPath("node")
		if err != nil {
			t.Skip("node unavailable locally")
		}
	}
	dir := t.TempDir()
	files := map[string][]byte{"app.js": AppScript, "harness.js": []byte(pinPrelude + harness)}
	for name, src := range map[string]string{"confirm.json": "requests-confirmations.golden.json", "verb.json": "verb-route.golden.json"} {
		b, err := os.ReadFile(filepath.Join("testdata", src))
		if err != nil {
			t.Fatal(err)
		}
		files[name] = b
	}
	for name, data := range files {
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
	if !strings.Contains(string(out), manifest) {
		t.Fatalf("harness manifest missing: %s", out)
	}
}

const pinPrelude = `
const assert = require('node:assert/strict');
const fs = require('node:fs');
const app = require('./app.js');
const G = JSON.parse(fs.readFileSync('confirm.json', 'utf8')), V = JSON.parse(fs.readFileSync('verb.json', 'utf8'));
const HEAD = (n) => '<div class="nsec" data-needs><div class="nh" style="color:var(--' + (n ? 'youShell' : 'sub') + ')">Needs you · ' + n + '</div>';
const real = G.v9.confirmations[1];   // the handler's summarised draft row, vr-00000000000000c4
const Q = 'ab'.repeat(32), dec = {request: 'sha256:' + Q, scope: 'worktree:1', requester: 'agent-1', deadline: 50520, reserved: '', refusal: ''};
const decLink = '<div class="nl"><a id="decide-' + Q + '" href="/surface/decide/sha256:' + Q + '">Decide · worktree:1 · requested by agent-1 · deadline 14:02</a></div>';
`

// TestNeedsYouConfirmationEmbeddedNewline (judge-405 J-JS-RE-NL): a confirmation id that is a valid
// id plus a line break and more, either side, is not a confirmation id — no link; nor is the id in
// another case (each of the regex's flags and anchors on its own); the valid row and a valid decision
// beside it still link.
func TestNeedsYouConfirmationEmbeddedNewline(t *testing.T) {
	runPinHarness(t, `
assert.equal(real.id, 'vr-00000000000000c4');
const ok = app.needsYouHTML({confirmations: [real], decisions: [dec]}, 0);
assert(ok.includes(decLink) && ok.includes('<a id="confirm-' + real.id + '" href="/surface/confirm/' + real.id + '#activity">'), ok);
for (const id of [real.id + '\n', real.id + '\njunk', 'junk\n' + real.id, '\n' + real.id, real.id + '\r', real.id + '\r\nx',
  real.id + '\u2028x', 'x\u2029' + real.id, real.id.toUpperCase(), 'VR' + real.id.slice(2), real.id.slice(0, 17) + 'C4']) {
  const h = app.needsYouHTML({confirmations: [{...real, id}], decisions: [dec]}, 0);
  assert.equal(h, HEAD(2) + decLink + '<div class="nl sub">Confirm · ' + app.esc(id) + ' · refused: not a confirmation id</div></div>', JSON.stringify(id));
  assert(!h.includes('/surface/confirm/') && !h.includes('confirm-vr'), JSON.stringify(id));
}
console.log('NEEDS YOU NEWLINE PASS');
`, "NEEDS YOU NEWLINE PASS")
}

// TestNeedsYouConfirmationsNonArray (judge-405 J-JS-CS-NOGUARD): a confirmations field that is not an
// array — an object, an array-like, a string, null, a number, a boolean, a thing with a map — lists
// nothing and throws nothing; the mint request and the decision of the same answer still render, exactly
// as with no confirmations; the real array still lists.
func TestNeedsYouConfirmationsNonArray(t *testing.T) {
	runPinHarness(t, `
const pending = [{id: 'mr-1', episode: 'e', agent: 'a', grants: 1, expires: 60}];
const base = {...G.v9, pending, decisions: [dec]};
const none = app.needsYouHTML({...base, confirmations: []}, 0);
assert(none.startsWith(HEAD(2) + '<div class="nl"><a id="needs-mr-1"') && none.includes(decLink), none);
assert(app.needsYouHTML(base, 0).startsWith(HEAD(5)), 'the real array lists');
for (const bad of [{}, {0: real, length: 1}, real.id, '', null, 0, 1, true, {map: () => ['<a href="/surface/confirm/x">']}]) {
  assert.equal(app.needsYouHTML({...base, confirmations: bad}, 0), none, JSON.stringify(bad));
}
const { confirmations: _, ...absent } = base;
assert.equal(app.needsYouHTML(absent, 0), none);
console.log('NEEDS YOU NON-ARRAY PASS');
`, "NEEDS YOU NON-ARRAY PASS")
}

// TestVerbDecidePrefixBoundary (judge-405 J-JS-ALLOW-DECIDE): a navigate answer opens a decision card
// only when its href STARTS with /surface/decide/sha256: — a foreign origin, a path that contains the
// prefix, the prefix without its colon, or another case is refused; the same boundaries hold for the
// approve and confirm prefixes. The daemon's own decide answer navigates.
func TestVerbDecidePrefixBoundary(t *testing.T) {
	runPinHarness(t, `
const REF = 'approve sha256:' + Q, MR = 'approve mr-0123456789abcdef', VR = 'approve vr-0123456789abcdef';
for (const q of [REF, MR, VR]) assert.deepEqual(app.verbOutcome(200, V[q]), {assign: JSON.parse(V[q]).href}, q);
const b = JSON.parse(V[REF]), H = 'sha256:' + Q;
assert.equal(b.href, '/surface/decide/' + H);
for (const href of ['//evil/surface/decide/' + H, 'https://evil/surface/decide/' + H, '/x/surface/decide/' + H,
  '/surface/approve/x/surface/decide/' + H, ' /surface/decide/' + H, '/surface/decide/sha256' + Q, '/surface/decide/sha256x' + Q,
  '/surface/decide/sha256/' + Q, '/surface/decide/SHA256:' + Q, '/Surface/decide/' + H,
  '/surface/approve/mr0123456789abcdef', '/surface/approve/mrx0123456789abcdef', '/surface/approve/MR-0123456789abcdef',
  '/x/surface/approve/mr-0123456789abcdef', '/SURFACE/confirm/vr-0123456789abcdef']) {
  assert.deepEqual(app.verbOutcome(200, JSON.stringify({...b, href})), {say: 'refused: unexpected target'}, href);
}
console.log('DECIDE PREFIX PASS');
`, "DECIDE PREFIX PASS")
}
